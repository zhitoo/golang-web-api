package user

type LoginRequest struct {
	Mobile string `json:"mobile" binding:"required,irmobile"`
	OTP    string `json:"otp" binding:"required,numeric"`
}
