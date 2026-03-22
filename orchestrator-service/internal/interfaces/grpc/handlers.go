package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	orchestratorpb "github.com/braunkc/ai-bot-constructor/orchestrator-service/api/orchestrator-service/v1"
	botdto "github.com/braunkc/ai-bot-constructor/orchestrator-service/internal/application/dto/bot"
	botdomain "github.com/braunkc/ai-bot-constructor/orchestrator-service/internal/domain/bot"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *OrchestratorServiceServer) CreateBot(ctx context.Context, req *orchestratorpb.CreateBotReq) (*orchestratorpb.Bot, error) {
	s.log.Debug("received request to create bot", slog.String("bot_name", req.Name))

	userID, err := s.userIDFromContext(ctx)
	if err != nil {
		return nil, s.grpcError(err)
	}

	bot, err := s.botUsecase.CreateBot(ctx, userID, &botdto.CreateBotReq{
		Name:         req.Name,
		ApiKey:       req.ApiKey,
		SystemPrompt: req.SystemPrompt,
	})
	if err != nil {
		return nil, s.grpcError(err)
	}

	return s.botDTOToPBModel(bot), nil
}

func (s *OrchestratorServiceServer) GetBot(ctx context.Context, req *orchestratorpb.GetBotReq) (*orchestratorpb.Bot, error) {
	s.log.Debug("received request to get bot", slog.String("bot_id", req.Id))

	userID, err := s.userIDFromContext(ctx)
	if err != nil {
		return nil, s.grpcError(err)
	}

	botID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid bot_id")
	}

	bot, err := s.botUsecase.GetBot(ctx, userID, &botdto.GetBotReq{
		ID: botID,
	})
	if err != nil {
		return nil, s.grpcError(err)
	}

	return s.botDTOToPBModel(bot), nil
}

func (s *OrchestratorServiceServer) GetAllBots(ctx context.Context, _ *emptypb.Empty) (*orchestratorpb.GetAllBotsResp, error) {
	s.log.Debug("received request to get all bots")

	userID, err := s.userIDFromContext(ctx)
	if err != nil {
		return nil, s.grpcError(err)
	}

	botsDTO, err := s.botUsecase.GetAllBots(ctx, userID)
	if err != nil {
		return nil, s.grpcError(err)
	}

	bots := make([]*orchestratorpb.Bot, 0, len(botsDTO))
	for _, bot := range botsDTO {
		bots = append(bots, s.botDTOToPBModel(bot))
	}

	return &orchestratorpb.GetAllBotsResp{
		Bots: bots,
	}, nil
}

func (s *OrchestratorServiceServer) StopBot(ctx context.Context, req *orchestratorpb.StopBotReq) (*orchestratorpb.Bot, error) {
	s.log.Debug("received request to stop bot", slog.String("bot_id", req.Id))

	userID, err := s.userIDFromContext(ctx)
	if err != nil {
		return nil, s.grpcError(err)
	}

	botID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid bot_id")
	}

	bot, err := s.botUsecase.StopBot(ctx, userID, &botdto.StopBotReq{
		ID: botID,
	})
	if err != nil {
		return nil, s.grpcError(err)
	}

	return s.botDTOToPBModel(bot), nil
}

func (s *OrchestratorServiceServer) StopBots(ctx context.Context, req *orchestratorpb.StopBotsReq) (*orchestratorpb.StopBotsResp, error) {
	s.log.Debug("received request to stop bots", slog.Any("bots_ids", req.Ids))

	userID, err := s.userIDFromContext(ctx)
	if err != nil {
		return nil, s.grpcError(err)
	}

	ids, err := s.parseUUIDs(req.Ids)
	if err != nil {
		return nil, err
	}

	resp, err := s.botUsecase.StopBots(ctx, userID, &botdto.StopBotsReq{
		IDs: ids,
	})
	if err != nil {
		return nil, s.grpcError(err)
	}

	bots := make([]*orchestratorpb.Bot, 0, len(resp.Bots))
	for _, bot := range resp.Bots {
		bots = append(bots, s.botDTOToPBModel(&bot))
	}

	return &orchestratorpb.StopBotsResp{
		Bots:         bots,
		AllSucceeded: resp.AllSucceeded,
	}, nil
}

func (s *OrchestratorServiceServer) StartBot(ctx context.Context, req *orchestratorpb.StartBotReq) (*orchestratorpb.Bot, error) {
	s.log.Debug("received request to start bot", slog.String("bot_id", req.Id))

	userID, err := s.userIDFromContext(ctx)
	if err != nil {
		return nil, s.grpcError(err)
	}

	botID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid bot_id")
	}

	bot, err := s.botUsecase.StartBot(ctx, userID, &botdto.StartBotReq{
		ID: botID,
	})
	if err != nil {
		return nil, s.grpcError(err)
	}

	return s.botDTOToPBModel(bot), nil
}

func (s *OrchestratorServiceServer) StartBots(ctx context.Context, req *orchestratorpb.StartBotsReq) (*orchestratorpb.StartBotsResp, error) {
	s.log.Debug("received request to start bots", slog.Any("bots_ids", req.Ids))

	userID, err := s.userIDFromContext(ctx)
	if err != nil {
		return nil, s.grpcError(err)
	}

	ids, err := s.parseUUIDs(req.Ids)
	if err != nil {
		return nil, err
	}

	resp, err := s.botUsecase.StartBots(ctx, userID, &botdto.StartBotsReq{
		IDs: ids,
	})
	if err != nil {
		return nil, s.grpcError(err)
	}

	bots := make([]*orchestratorpb.Bot, 0, len(resp.Bots))
	for _, bot := range resp.Bots {
		bots = append(bots, s.botDTOToPBModel(&bot))
	}

	return &orchestratorpb.StartBotsResp{
		Bots:         bots,
		AllSucceeded: resp.AllSucceeded,
	}, nil
}

func (s *OrchestratorServiceServer) RestartBot(ctx context.Context, req *orchestratorpb.RestartBotReq) (*orchestratorpb.Bot, error) {
	s.log.Debug("received request to restart bot", slog.String("bot_id", req.Id))

	userID, err := s.userIDFromContext(ctx)
	if err != nil {
		return nil, s.grpcError(err)
	}

	botID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid bot_id")
	}

	bot, err := s.botUsecase.RestartBot(ctx, userID, &botdto.RestartBotReq{
		ID: botID,
	})
	if err != nil {
		return nil, s.grpcError(err)
	}

	return s.botDTOToPBModel(bot), nil
}

func (s *OrchestratorServiceServer) DeleteBot(ctx context.Context, req *orchestratorpb.DeleteBotReq) (*emptypb.Empty, error) {
	s.log.Debug("received request to delete bot", slog.String("bot_id", req.Id))

	userID, err := s.userIDFromContext(ctx)
	if err != nil {
		return nil, s.grpcError(err)
	}

	botID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid bot_id")
	}

	if err := s.botUsecase.DeleteBot(ctx, userID, &botdto.DeleteBotReq{
		ID: botID,
	}); err != nil {
		return nil, s.grpcError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *OrchestratorServiceServer) DeleteBots(ctx context.Context, req *orchestratorpb.DeleteBotsReq) (*orchestratorpb.DeleteBotsResp, error) {
	s.log.Debug("received request to delete bots", slog.Any("bots_ids", req.Ids))

	userID, err := s.userIDFromContext(ctx)
	if err != nil {
		return nil, s.grpcError(err)
	}

	ids, err := s.parseUUIDs(req.Ids)
	if err != nil {
		return nil, err
	}

	resp, err := s.botUsecase.DeleteBots(ctx, userID, &botdto.DeleteBotsReq{
		IDs: ids,
	})
	if err != nil {
		return nil, s.grpcError(err)
	}

	return &orchestratorpb.DeleteBotsResp{
		AllSucceeded: resp.AllSucceeded,
	}, nil
}

func (s *OrchestratorServiceServer) DeleteAllBots(ctx context.Context, _ *emptypb.Empty) (*orchestratorpb.DeleteAllBotsResp, error) {
	s.log.Debug("received request to delete all user bots")

	userID, err := s.userIDFromContext(ctx)
	if err != nil {
		return nil, s.grpcError(err)
	}

	resp, err := s.botUsecase.DeleteAllBots(ctx, userID)
	if err != nil {
		return nil, s.grpcError(err)
	}

	return &orchestratorpb.DeleteAllBotsResp{
		AllSucceeded: resp.AllSucceeded,
	}, nil
}

func (s *OrchestratorServiceServer) userIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return uuid.Nil, errors.New("no user_id in context")
	}

	return uuid.Parse(userID)
}

func (s *OrchestratorServiceServer) parseUUIDs(ids []string) ([]uuid.UUID, error) {
	result := make([]uuid.UUID, 0, len(ids))
	for _, idStr := range ids {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid bot id: %s", idStr))
		}
		result = append(result, id)
	}
	return result, nil
}

func (s *OrchestratorServiceServer) botDTOToPBModel(bot *botdto.Bot) *orchestratorpb.Bot {
	return &orchestratorpb.Bot{
		Id:        bot.ID.String(),
		UserId:    bot.UserID.String(),
		Status:    orchestratorpb.BotStatus(bot.BotStatus),
		Name:      bot.Name,
		LastError: bot.LastError,
		CreatedAt: timestamppb.New(bot.CreatedAt),
		UpdatedAt: timestamppb.New(bot.UpdatedAt),
	}
}

func (s *OrchestratorServiceServer) grpcError(err error) error {
	switch {
	case errors.Is(err, botdomain.ErrBotNameMustBeLonger),
		errors.Is(err, botdomain.ErrBotNameMustBeShorter),
		errors.Is(err, botdomain.ErrInvalidApiKey),
		errors.Is(err, botdomain.ErrInvalidBotStatus):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, botdomain.ErrDuplicatedKey):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, botdomain.ErrNotEnoughRights):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, botdomain.ErrRecordNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, botdomain.ErrInvalidStorageData):
		s.log.Error("internal error", slog.Any("err", err))

		return status.Error(codes.Internal, "internal error")
	default:
		s.log.Error("internal error", slog.Any("err", err))

		return status.Error(codes.Internal, "internal error")
	}
}
