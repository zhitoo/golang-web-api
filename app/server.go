package app

import (
	"fmt"

	"github.com/zhitoo/golang-web-api/app/middlewares"
	"github.com/zhitoo/golang-web-api/app/modules/health"
	"github.com/zhitoo/golang-web-api/app/modules/otp"
	"github.com/zhitoo/golang-web-api/app/modules/test"
	"github.com/zhitoo/golang-web-api/app/modules/user"
	"github.com/zhitoo/golang-web-api/app/validations"
	"github.com/zhitoo/golang-web-api/config"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/zhitoo/golang-web-api/docs"
)

func InitServer(cfg *config.Config) {
	r := gin.New()

	RegisterValidatores()

	r.Use(middlewares.DefaultStructuredLogger(cfg), gin.Logger(), gin.Recovery())

	RegisterRoutes(r)
	RegisterSwagger(r, cfg)

	r.Run(fmt.Sprintf(":%s", cfg.Server.Port))
}

func RegisterValidatores() {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		val.RegisterValidation("irmobile", validations.IranianMobileNumberValidator, true)
	}
}
func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			health.Routes(v1.Group("/health"))
			test.Routes(v1.Group("/test"))
			otp.Routes(v1.Group("/otp"))
			user.Routes(v1.Group("/users"))
		}
	}
}

func RegisterSwagger(r *gin.Engine, cfg *config.Config) {
	docs.SwaggerInfo.Title = "Web API"
	docs.SwaggerInfo.Description = "Web API Docs"
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.Schemes = []string{"http"}
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = "localhost:" + cfg.Server.Port
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
