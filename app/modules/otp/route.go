package otp

import (
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {

	r.POST("/verify", verifyOTPHandler)
	r.Use(LimitOTP())
	r.POST("/send", sendOTPHandler)
}
