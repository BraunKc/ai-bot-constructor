package httproutes

import (
	httphandlers "github.com/braunkc/ai-bot-constructor/auth-gateway/internal/interfaces/http/handlers"
	httpmiddlewares "github.com/braunkc/ai-bot-constructor/auth-gateway/internal/interfaces/http/middlewares"
	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine, httpHandlers httphandlers.HTTPHandlers) {
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.POST("/register", httpHandlers.Register())
			v1.POST("/login", httpHandlers.Login())
			user := v1.Group("/user")
			user.Use(httpmiddlewares.AuthMiddleware())
			{
				user.GET("", httpHandlers.GetUser())
				user.PATCH("", httpHandlers.UpdateUser())
				user.DELETE("", httpHandlers.DeleteUser())
			}
		}
	}
}
