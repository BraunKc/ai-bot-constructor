package httpserver

import (
	httphandlers "github.com/braunkc/ai-bot-constructor/auth-gateway/internal/interfaces/http/handlers"
	httproutes "github.com/braunkc/ai-bot-constructor/auth-gateway/internal/interfaces/http/routes"
	"github.com/gin-gonic/gin"
)

func New(httpHandlers httphandlers.HTTPHandlers) *gin.Engine {
	r := gin.Default()
	httproutes.Init(r, httpHandlers)
	return r
}
