package user

import (
	"github.com/gin-gonic/gin"
	"github.com/zhitoo/golang-web-api/app/middlewares"
	"github.com/zhitoo/golang-web-api/config"
	"github.com/zhitoo/golang-web-api/pkg/logging"
)

var log logging.ScopedLogger

func init() {
	cfg := config.GetConfig()
	log = logging.NewLogger(cfg).With(logging.Internal, logging.Api)
}

func Routes(r *gin.RouterGroup) {
	cfg := config.GetConfig()
	svc := NewUserService(cfg)

	r.POST("/login", loginHandler(svc))
	r.POST("/refresh-token", refreshTokenHandler(svc))
	r.POST("/reset-password", resetPasswordHandler(svc))

	authenticated := r.Group("", middlewares.Authentication(cfg))
	{
		authenticated.GET("/me", meHandler(svc))
		authenticated.POST("/change-password", changePasswordHandler(svc))
	}
}
