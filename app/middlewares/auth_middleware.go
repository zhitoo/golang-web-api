package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	authSvc "github.com/zhitoo/golang-web-api/app/modules/auth"
	"github.com/zhitoo/golang-web-api/app/response"
	"github.com/zhitoo/golang-web-api/config"
)

func Authentication(cfg *config.Config) gin.HandlerFunc {
	svc := authSvc.NewAuthService(cfg)
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.NewResponse().SetError(fmt.Errorf("missing authorization header")).SetHttpStatusCode(http.StatusUnauthorized).Json(c)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := svc.GetClaims(tokenStr)
		if err != nil {
			response.NewResponse().SetError(fmt.Errorf("invalid or expired token")).SetHttpStatusCode(http.StatusUnauthorized).Json(c)
			c.Abort()
			return
		}

		if userId, ok := claims["user_id"].(string); ok {
			c.Set("UserId", userId)
		}
		if username, ok := claims["username"].(string); ok {
			c.Set("Username", username)
		}
		if roles, ok := claims["roles"].([]any); ok {
			roleStrs := make([]string, 0, len(roles))
			for _, r := range roles {
				if s, ok := r.(string); ok {
					roleStrs = append(roleStrs, s)
				}
			}
			c.Set("Roles", roleStrs)
		}

		c.Next()
	}
}
