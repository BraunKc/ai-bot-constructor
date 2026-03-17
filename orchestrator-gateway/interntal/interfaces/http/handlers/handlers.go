package httphandlers

import (
	"errors"
	"net/http"

	botdto "github.com/braunkc/ai-bot-constructor/orchestrator-gateway/interntal/application/dto/bot"
	botusecase "github.com/braunkc/ai-bot-constructor/orchestrator-gateway/interntal/application/usecase/bot"
	orchestratorerrors "github.com/braunkc/ai-bot-constructor/orchestrator-gateway/interntal/infra/grpc/orchestrator/errors"
	"github.com/gin-gonic/gin"
)

type HTTPHandlers interface {
	CreateBot(ctx *gin.Context)
	GetBot(ctx *gin.Context)
	GetAllBots(ctx *gin.Context)
	StartBot(ctx *gin.Context)
	StartBots(ctx *gin.Context)
	StopBot(ctx *gin.Context)
	StopBots(ctx *gin.Context)
	RestartBot(ctx *gin.Context)
	DeleteBot(ctx *gin.Context)
	DeleteBots(ctx *gin.Context)
	DeleteAllBots(ctx *gin.Context)
}

type httpHandlers struct {
	botUsecase botusecase.BotUsecase
}

func NewBotHandlers(botUsecase botusecase.BotUsecase) HTTPHandlers {
	return &httpHandlers{
		botUsecase: botUsecase,
	}
}

func (hh *httpHandlers) CreateBot(ctx *gin.Context) {
	var req botdto.CreateBotRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	createdBot, err := hh.botUsecase.CreateBot(ctx.Request.Context(), &req)
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

	ctx.JSON(http.StatusCreated, createdBot)
}

func (hh *httpHandlers) GetBot(ctx *gin.Context) {
	id := ctx.Param("id")
	req := botdto.GetBotRequest{ID: id}

	b, err := hh.botUsecase.GetBot(ctx.Request.Context(), &req)
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

	ctx.JSON(http.StatusOK, b)
}

func (hh *httpHandlers) GetAllBots(ctx *gin.Context) {
	bots, err := hh.botUsecase.GetAllBots(ctx.Request.Context())
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

	ctx.JSON(http.StatusOK, bots)
}

func (hh *httpHandlers) StartBot(ctx *gin.Context) {
	id := ctx.Param("id")
	req := botdto.GetBotRequest{ID: id}

	b, err := hh.botUsecase.StartBot(ctx.Request.Context(), &req)
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

	ctx.JSON(http.StatusOK, b)
}

func (hh *httpHandlers) StartBots(ctx *gin.Context) {
	var req botdto.IDsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	resp, err := hh.botUsecase.StartBots(ctx.Request.Context(), &req)
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

	ctx.JSON(http.StatusOK, resp)
}

func (hh *httpHandlers) StopBot(ctx *gin.Context) {
	id := ctx.Param("id")
	req := botdto.GetBotRequest{ID: id}

	b, err := hh.botUsecase.StopBot(ctx.Request.Context(), &req)
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

	ctx.JSON(http.StatusOK, b)
}

func (hh *httpHandlers) StopBots(ctx *gin.Context) {
	var req botdto.IDsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	resp, err := hh.botUsecase.StopBots(ctx.Request.Context(), &req)
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

	ctx.JSON(http.StatusOK, resp)
}

func (hh *httpHandlers) RestartBot(ctx *gin.Context) {
	id := ctx.Param("id")
	req := botdto.GetBotRequest{ID: id}

	b, err := hh.botUsecase.RestartBot(ctx.Request.Context(), &req)
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

	ctx.JSON(http.StatusOK, b)
}

func (hh *httpHandlers) DeleteBot(ctx *gin.Context) {
	id := ctx.Param("id")
	req := botdto.GetBotRequest{ID: id}

	if err := hh.botUsecase.DeleteBot(ctx.Request.Context(), &req); err != nil {
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

func (hh *httpHandlers) DeleteBots(ctx *gin.Context) {
	var req botdto.IDsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	succeeded, err := hh.botUsecase.DeleteBots(ctx.Request.Context(), &req)
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

	ctx.JSON(http.StatusOK, gin.H{"all_succeeded": succeeded})
}

func (hh *httpHandlers) DeleteAllBots(ctx *gin.Context) {
	succeeded, err := hh.botUsecase.DeleteAllBots(ctx.Request.Context())
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

	ctx.JSON(http.StatusOK, gin.H{"all_succeeded": succeeded})
}

func (hh *httpHandlers) appError(err error) (*orchestratorerrors.AppError, bool) {
	if err == nil {
		return nil, false
	}

	var appError *orchestratorerrors.AppError
	if ok := errors.As(err, &appError); ok {
		return appError, true
	}

	return nil, false
}
