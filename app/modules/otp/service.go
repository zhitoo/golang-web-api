package otp

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/zhitoo/golang-web-api/database/cache"
)

func SendOTP(mobile string) (string, error) {
	// Generate a random 6-digit OTP
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	//send otp via sms

	// Store the OTP in Redis with an expiration time (e.g., 5 minutes)
	return otp, cache.SetValue("otp:"+mobile, otp, 2*time.Minute)

}

func VerifyOTP(mobile string, otp string) error {
	cachedOTP, err := cache.GetValue[string]("otp:" + mobile)
	if err != nil {
		return err
	}
	if cachedOTP != otp {
		return fmt.Errorf("invalid OTP")
	}
	cache.DeleteValue("otp:" + mobile)
	return nil

}
