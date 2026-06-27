package health

import (
	"github.com/zhitoo/golang-web-api/app/response"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	//r.Use(TestMiddleware())
	r.GET("", healthHandler)
}

// @Summary      Health Check
// @Description  returns ok if the service is healthy
// @Tags         health
// @Produce      json
// @Success      200  {object}  response.BaseHttpResponse
// @Router       /v1/health [get]
func healthHandler(c *gin.Context) {
	HealthCheck()
	response.NewResponse().SetResult(":)").Json(c)
}
