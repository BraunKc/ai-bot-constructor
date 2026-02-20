package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	authpb "github.com/braunkc/ai-bot-constructor/auth-service/api/auth-service/v1"
	"github.com/braunkc/ai-bot-constructor/auth-service/config"
	userusecase "github.com/braunkc/ai-bot-constructor/auth-service/internal/application/usecase/user"
	userdomain "github.com/braunkc/ai-bot-constructor/auth-service/internal/domain/user"
	"github.com/braunkc/ai-bot-constructor/auth-service/internal/infra/database"
	orchestratorgrpc "github.com/braunkc/ai-bot-constructor/auth-service/internal/infra/grpc/orchestrator"
	hasherinfra "github.com/braunkc/ai-bot-constructor/auth-service/internal/infra/hasher"
	"github.com/braunkc/ai-bot-constructor/auth-service/internal/infra/jwt"
	userpersistence "github.com/braunkc/ai-bot-constructor/auth-service/internal/infra/persistence/user"
	grpcserver "github.com/braunkc/ai-bot-constructor/auth-service/internal/interfaces/grpc"
	grpcinterceptors "github.com/braunkc/ai-bot-constructor/auth-service/internal/interfaces/grpc/interceptors"
	"github.com/braunkc/ai-bot-constructor/auth-service/pkg/log"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type App struct {
	grpcServer         *grpc.Server
	db                 *gorm.DB
	userRepo           userdomain.UserRepo
	orchestratorClient userusecase.OrchestratorClient
	tokenManager       *jwt.TokenManager
	userUsecase        userusecase.UserUsecase
	cfg                *config.Config
	log                *slog.Logger
}

func (app *App) initConfig(path string) error {
	if err := godotenv.Load(path); err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	cfg, err := config.New(os.Getenv("CONFIG_PATH"))
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}
	app.cfg = cfg
	slog.Debug("cfg created")

	return nil
}

func (app *App) initLogger() error {
	var loggerOutputType log.OutputType
	switch app.cfg.Logger.OutputType {
	case config.Console:
		loggerOutputType = log.Console
	case config.File:
		loggerOutputType = log.File
	case config.Both:
		loggerOutputType = log.Both
	}

	var loggerLevel slog.Level
	switch app.cfg.Logger.Level {
	case config.Debug:
		loggerLevel = slog.LevelDebug
	case config.Info:
		loggerLevel = slog.LevelInfo
	case config.Warn:
		loggerLevel = slog.LevelWarn
	case config.Error:
		loggerLevel = slog.LevelError
	}

	logCfg := log.Config{
		Service:    app.cfg.Logger.Service,
		OutputType: loggerOutputType,
		Level:      loggerLevel,
	}
	h, err := log.NewHandler(&logCfg)
	if err != nil {
		return fmt.Errorf("failed to create log handler: %w", err)
	}
	app.log = slog.New(h)

	return nil
}

func (app *App) initDB() error {
	db, err := database.New(&app.cfg.DB)
	if err != nil {
		return fmt.Errorf("failed to create database repository: %w", err)
	}

	if app.cfg.Env == config.Develop {
		if err := db.AutoMigrate(&userpersistence.User{}); err != nil {
			return fmt.Errorf("failed to migrate tables: %w", err)
		}
	}
	app.db = db
	app.log.Info("database inited")

	return nil
}

func (app *App) initUserRepo() {
	app.userRepo = userpersistence.NewRepo(app.db, app.log)
	app.log.Info("user repository created")
}

func (app *App) initOrchestratorClient() error {
	orchestratorClient, err := orchestratorgrpc.NewClient(&app.cfg.GRPC.OrchestratorService, app.log)
	if err != nil {
		return fmt.Errorf("failed to create orchestrator grpc client: %w", err)
	}
	app.orchestratorClient = orchestratorClient
	app.log.Info("orchestrator client created")

	return nil
}

func (app *App) initTokenManager() error {
	tokenManager, err := jwt.NewTokenManager(os.Getenv("SECRET_KEY"))
	if err != nil {
		return fmt.Errorf("failed to create token manager: %w", err)
	}
	app.tokenManager = tokenManager
	app.log.Info("token manager created")

	return nil
}

func (app *App) initUserUsecase() {
	app.userUsecase = userusecase.New(app.userRepo, app.orchestratorClient, app.tokenManager, &hasherinfra.Hasher{Cost: bcrypt.DefaultCost}, app.log)
	app.log.Info("user usecase created")
}

func (app *App) initGRPCServer() {
	var protectedMethods = map[string]struct{}{
		authpb.Auth_GetUser_FullMethodName:    {},
		authpb.Auth_UpdateUser_FullMethodName: {},
		authpb.Auth_DeleteUser_FullMethodName: {},
	}

	authInterceptor := grpcinterceptors.NewAuthInterceptor(app.tokenManager, protectedMethods)

	app.grpcServer = grpcserver.New(authInterceptor, app.userUsecase, app.log)
	app.log.Info("grpc server created")
}

func New(envPath string) (*App, error) {
	var app App
	if err := app.initConfig(envPath); err != nil {
		return nil, err
	}
	if err := app.initLogger(); err != nil {
		return nil, err
	}
	if err := app.initDB(); err != nil {
		return nil, err
	}
	app.initUserRepo()
	if err := app.initOrchestratorClient(); err != nil {
		return nil, err
	}
	if err := app.initTokenManager(); err != nil {
		return nil, err
	}
	app.initUserUsecase()
	app.initGRPCServer()

	return &app, nil
}

func (app *App) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", app.cfg.GRPC.Host, app.cfg.GRPC.Port))
	if err != nil {
		return fmt.Errorf("failed to create grpc listener: %w", err)
	}

	app.log.Info("grpc server serve",
		slog.String("host", app.cfg.GRPC.Host),
		slog.String("port", app.cfg.GRPC.Port),
	)
	if err := app.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve grpc server: %w", err)
	}

	return nil
}

func (app *App) Stop() {
	app.log.Info("graceful shutdown")

	doneCh := make(chan any, 1)

	go func() {
		defer close(doneCh)

		app.grpcServer.GracefulStop()
		if err := app.userRepo.Close(); err != nil {
			app.log.Error("failed to close db gracefully")
		}
		if err := app.orchestratorClient.Close(); err != nil {
			app.log.Error("failed to close orchestrator conn gracefully")
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case <-doneCh:
		app.log.Info("successfully graceful shutdown")
	case <-ctx.Done():
		app.log.Error("failed graceful shutdown, terminate")
	}
}
