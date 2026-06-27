package user

type LoginRequest struct {
	Mobile string `json:"mobile" binding:"required_without=Email,omitempty,irmobile"`
	Email  string `json:"email" binding:"required_without=Mobile,omitempty,email"`
	OTP    string `json:"otp" binding:"required,numeric"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
