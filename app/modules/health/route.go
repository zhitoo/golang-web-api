package health

import (
	"carsale/app/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	r.Use(TestMiddleware())
	r.GET("", func(c *gin.Context) {
		HealthCheck()
		c.JSON(http.StatusOK, helper.GenerateBaseResponse(":)", true, http.StatusOK))
	})
}
