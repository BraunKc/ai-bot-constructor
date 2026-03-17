package httpserver

import (
	botusecase "github.com/braunkc/ai-bot-constructor/orchestrator-gateway/interntal/application/usecase/bot"
	httphandlers "github.com/braunkc/ai-bot-constructor/orchestrator-gateway/interntal/interfaces/http/handlers"
	httpmiddlewares "github.com/braunkc/ai-bot-constructor/orchestrator-gateway/interntal/interfaces/http/middlewares"
	httproutes "github.com/braunkc/ai-bot-constructor/orchestrator-gateway/interntal/interfaces/http/routes"
	"github.com/gin-gonic/gin"
)

func New(botUsecase botusecase.BotUsecase) *gin.Engine {
	r := gin.Default()
	r.Use(httpmiddlewares.AuthMiddleware())

	h := httphandlers.NewBotHandlers(botUsecase)
	httproutes.Init(r, h)

	return r
}
