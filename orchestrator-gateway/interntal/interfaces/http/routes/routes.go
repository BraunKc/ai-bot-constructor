package httproutes

import (
	httphandlers "github.com/braunkc/ai-bot-constructor/orchestrator-gateway/interntal/interfaces/http/handlers"
	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine, h httphandlers.HTTPHandlers) {
	api := r.Group("/api/v1")
	{
		bots := api.Group("/bots")
		{
			bots.POST("", h.CreateBot)
			bots.GET("", h.GetAllBots)
			bots.DELETE("", h.DeleteAllBots)

			bots.POST("/start-batch", h.StartBots)
			bots.POST("/stop-batch", h.StopBots)
			bots.DELETE("/batch", h.DeleteBots)

			bots.GET("/:id", h.GetBot)
			bots.POST("/:id/start", h.StartBot)
			bots.POST("/:id/stop", h.StopBot)
			bots.POST("/:id/restart", h.RestartBot)
			bots.DELETE("/:id", h.DeleteBot)
		}
	}
}
