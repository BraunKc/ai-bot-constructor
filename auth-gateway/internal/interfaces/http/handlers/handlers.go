package httphandlers

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	userdto "github.com/braunkc/ai-bot-constructor/auth-gateway/internal/application/dto/user"
	userusecase "github.com/braunkc/ai-bot-constructor/auth-gateway/internal/application/usecase/user"
	autherrors "github.com/braunkc/ai-bot-constructor/auth-gateway/internal/infra/grpc/auth/errors"
	"github.com/gin-gonic/gin"
)

type HTTPHandlers interface {
	Register() gin.HandlerFunc
	Login() gin.HandlerFunc
	GetUser() gin.HandlerFunc
	UpdateUser() gin.HandlerFunc
	DeleteUser() gin.HandlerFunc
}

type httpHandlers struct {
	userUsecase userusecase.UserUsecase
	log         *slog.Logger
}

func New(userUsecase userusecase.UserUsecase, log *slog.Logger) HTTPHandlers {
	return &httpHandlers{
		userUsecase: userUsecase,
		log:         log,
	}
}

func (hh *httpHandlers) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hh.log.Debug("received register request")

		var req userdto.AuthReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})

			return
		}

		resp, err := hh.userUsecase.Register(ctx.Request.Context(), &req)
		if err != nil {
			if appError, ok := hh.appError(err); ok {
				ctx.AbortWithStatusJSON(appError.HTTPStatus, gin.H{
					"error": appError.Message,
				})

				return
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "internal error",
			})

			return
		}

		ctx.SetCookie("access_token", resp.Token, int(24*time.Hour.Seconds()), "/", "", false, true)
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "registered successfully",
		})
	}
}

func (hh *httpHandlers) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hh.log.Debug("received login request")

		var req userdto.AuthReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})

			return
		}

		resp, err := hh.userUsecase.Login(ctx.Request.Context(), &req)
		if err != nil {
			if appError, ok := hh.appError(err); ok {
				ctx.AbortWithStatusJSON(appError.HTTPStatus, gin.H{
					"error": appError.Message,
				})

				return
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "internal error",
			})

			return
		}

		ctx.SetCookie("access_token", resp.Token, int(24*time.Hour.Seconds()), "/", "", false, true)
		ctx.JSON(http.StatusOK, gin.H{
			"message": "login successfully",
		})
	}
}

func (hh *httpHandlers) GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hh.log.Debug("received get user request")

		resp, err := hh.userUsecase.GetUser(ctx.Request.Context())
		if err != nil {
			if appError, ok := hh.appError(err); ok {
				ctx.AbortWithStatusJSON(appError.HTTPStatus, gin.H{
					"error": appError.Message,
				})

				return
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "internal error",
			})

			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"user": resp,
		})
	}
}

func (hh *httpHandlers) UpdateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hh.log.Debug("received update user request")

		var req userdto.UpdateUserReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})

			return
		}

		resp, err := hh.userUsecase.UpdateUser(ctx.Request.Context(), &req)
		if err != nil {
			if appError, ok := hh.appError(err); ok {
				ctx.AbortWithStatusJSON(appError.HTTPStatus, gin.H{
					"error": appError.Message,
				})

				return
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "internal error",
			})

			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"user": resp,
		})
	}
}

func (hh *httpHandlers) DeleteUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hh.log.Debug("received delete user request")

		if err := hh.userUsecase.DeleteUser(ctx.Request.Context()); err != nil {
			if appError, ok := hh.appError(err); ok {
				ctx.AbortWithStatusJSON(appError.HTTPStatus, gin.H{
					"error": appError.Message,
				})

				return
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "internal error",
			})

			return
		}

		ctx.Status(http.StatusNoContent)
	}
}

func (hh *httpHandlers) appError(err error) (*autherrors.AppError, bool) {
	if err == nil {
		return nil, false
	}

	var appError *autherrors.AppError
	if ok := errors.As(err, &appError); ok {
		return appError, true
	}

	return nil, false
}
