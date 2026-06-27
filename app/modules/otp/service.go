package otp

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/zhitoo/golang-web-api/database/cache"
	"github.com/zhitoo/golang-web-api/pkg/job"
	"github.com/zhitoo/golang-web-api/pkg/logging"
)

func SendOTP(mobile string) (string, error) {
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	if err := job.Dispatch(&SendSmsJob{
		Mobile: mobile,
		Body:   fmt.Sprintf("Your OTP code is: %s", otp),
	}); err != nil {
		log.Error("failed to dispatch SMS job", map[logging.ExtraKey]any{logging.ErrorMessage: err.Error()})
	}

	if err := cache.SetValue("otp:"+mobile, otp, 2*time.Minute); err != nil {
		log.Error("failed to cache OTP", map[logging.ExtraKey]any{logging.ErrorMessage: err.Error()})
		return "", err
	}

	log.Info("OTP sent", nil)
	return otp, nil
}

func VerifyOTP(mobile string, otp string) error {
	cachedOTP, err := cache.GetValue[string]("otp:" + mobile)
	if err != nil {
		log.Error("failed to get OTP from cache", map[logging.ExtraKey]any{logging.ErrorMessage: err.Error()})
		return err
	}
	if cachedOTP != otp {
		return fmt.Errorf("invalid OTP")
	}
	err = cache.DeleteValue("otp:" + mobile)
	if err != nil {
		return err
	}
	log.Info("OTP verified", nil)
	return nil
}
