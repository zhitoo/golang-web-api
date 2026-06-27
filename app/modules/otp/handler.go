package otp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhitoo/golang-web-api/app/response"
	"github.com/zhitoo/golang-web-api/config"
	"github.com/zhitoo/golang-web-api/pkg/logging"
)

var cfg *config.Config
var log logging.ScopedLogger

func init() {
	if cfg == nil {
		cfg = config.GetConfig()
		log = logging.NewLogger(cfg).With(logging.Internal, logging.Api)
	}
}

func sendOTPHandler(c *gin.Context) {
	request := new(SendOTPRequest)

	if err := c.ShouldBindJSON(request); err != nil {
		log.Warn("invalid send OTP request", map[logging.ExtraKey]any{logging.ErrorMessage: err.Error()})
		response.NewResponse().SetStatus(false).SetError(err).SetHttpStatusCode(http.StatusUnprocessableEntity).Json(c)
		return
	}

	otp, err := SendOTP(request.Mobile)
	if err != nil {
		log.Error("failed to send OTP", map[logging.ExtraKey]any{logging.ErrorMessage: err.Error()})
		response.NewResponse().SetStatus(false).SetError(err).SetHttpStatusCode(http.StatusBadRequest).Json(c)
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

	response.NewResponse().SetResult(result).Json(c)
}

func verifyOTPHandler(c *gin.Context) {
	request := new(VerifyOTPRequest)
	if err := c.ShouldBindJSON(request); err != nil {
		log.Warn("invalid verify OTP request", map[logging.ExtraKey]any{logging.ErrorMessage: err.Error()})
		response.NewResponse().SetStatus(false).SetError(err).SetHttpStatusCode(http.StatusUnprocessableEntity).Json(c)
		return
	}
	err := VerifyOTP(request.Mobile, request.OTP)
	if err != nil {
		log.Warn("OTP verification failed", map[logging.ExtraKey]any{logging.ErrorMessage: err.Error()})
		response.NewResponse().SetStatus(false).SetError(err).SetHttpStatusCode(http.StatusUnauthorized).Json(c)
		return
	}
	response.NewResponse().SetResult(":)").Json(c)
}
