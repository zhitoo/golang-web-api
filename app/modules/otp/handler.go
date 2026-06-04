package otp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhitoo/golang-web-api/app/response"
	"github.com/zhitoo/golang-web-api/config"
)

var cfg *config.Config

func init() {
	if cfg == nil {
		cfg = config.GetConfig()
	}
}

func sendOTPHandler(c *gin.Context) {
	request := new(SendOTPRequest)

	if err := c.ShouldBindJSON(request); err != nil {
		response.NewReponse().SetStatus(false).SetError(err).SetHttpStatusCode(http.StatusUnprocessableEntity).Json(c)
		return
	}

	otp, err := SendOTP(request.Mobile)
	if err != nil {
		response.NewReponse().SetStatus(false).SetError(err).SetHttpStatusCode(http.StatusBadRequest).Json(c)
		return
	}

	result := gin.H{
		"message": "OTP sent successfully",
	}
	if cfg.App.Env == "local" {
		result = gin.H{
			"message": "OTP sent successfully",
			"otp":     otp,
		}
	}

	response.NewReponse().SetResult(result).Json(c)
}

func verifyOTPHandler(c *gin.Context) {
	request := new(VerifyOTPRequest)
	if err := c.ShouldBindJSON(request); err != nil {
		response.NewReponse().SetStatus(false).SetError(err).SetHttpStatusCode(http.StatusUnprocessableEntity).Json(c)
		return
	}
	err := VerifyOTP(request.Mobile, request.OTP)
	if err != nil {
		response.NewReponse().SetStatus(false).SetError(err).SetHttpStatusCode(http.StatusUnauthorized).Json(c)
		return
	}
	response.NewReponse().SetResult(":)").Json(c)
}
