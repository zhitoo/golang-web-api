package health

import (
	"github.com/zhitoo/golang-web-api/app/response"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	//r.Use(TestMiddleware())
	r.GET("", func(c *gin.Context) {
		HealthCheck()
		response.NewResponse().SetResult(":)").Json(c)
	})
}
