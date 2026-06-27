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

// @Summary      Send OTP
// @Description  sends a 6-digit OTP to the given mobile number (stored in Redis for 2 minutes); in local env the OTP is returned in the response
// @Tags         otp
// @Accept       json
// @Produce      json
// @Param        body  body      SendOTPRequest  true  "Mobile number"
// @Success      200   {object}  response.BaseHttpResponse
// @Failure      422   {object}  response.BaseHttpResponse
// @Failure      400   {object}  response.BaseHttpResponse
// @Router       /v1/otp/send [post]
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

// @Summary      Verify OTP
// @Description  verifies the OTP sent to the given mobile number and deletes it from Redis on success
// @Tags         otp
// @Accept       json
// @Produce      json
// @Param        body  body      VerifyOTPRequest  true  "Mobile number and OTP"
// @Success      200   {object}  response.BaseHttpResponse
// @Failure      422   {object}  response.BaseHttpResponse
// @Failure      401   {object}  response.BaseHttpResponse
// @Router       /v1/otp/verify [post]
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
