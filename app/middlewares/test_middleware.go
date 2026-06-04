package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TestMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("x-api-key")

		if apiKey == "hehe" {
			ctx.Next()
		} else {

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "request not authenticated",
			})
		}

	}
}
