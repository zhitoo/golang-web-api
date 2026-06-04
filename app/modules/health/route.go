package health

import (
	"net/http"

	"github.com/zhitoo/golang-web-api/app/helper"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	r.Use(TestMiddleware())
	r.GET("", func(c *gin.Context) {
		HealthCheck()
		c.JSON(http.StatusOK, helper.GenerateBaseResponse(":)", true, http.StatusOK))
	})
}
