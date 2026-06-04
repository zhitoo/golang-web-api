package otp

type SendOTPRequest struct {
	Mobile string `json:"mobile" binding:"required,irmobile"`
}

type VerifyOTPRequest struct {
	Mobile string `json:"mobile" binding:"required,irmobile"`
	OTP    string `json:"otp" binding:"required"`
}
