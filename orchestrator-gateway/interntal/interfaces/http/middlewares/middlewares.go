package httpmiddlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const tokenCtxKey = "authorization"

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("access_token")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})

			return
		}

		if strings.TrimSpace(token) == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})

			return
		}

		c := context.WithValue(ctx.Request.Context(), tokenCtxKey, token)
		ctx.Request = ctx.Request.WithContext(c)

		ctx.Next()
	}
}
