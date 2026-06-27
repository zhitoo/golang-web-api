package user

type LoginRequest struct {
	Mobile   string `json:"mobile" binding:"required_without=Email,omitempty,irmobile"`
	Email    string `json:"email" binding:"required_without=Mobile,omitempty,email"`
	OTP      string `json:"otp" binding:"omitempty,numeric"`
	Password string `json:"password" binding:"omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type ResetPasswordRequest struct {
	Mobile   string `json:"mobile" binding:"required_without=Email,omitempty,irmobile"`
	Email    string `json:"email" binding:"required_without=Mobile,omitempty,email"`
	OTP      string `json:"otp" binding:"required,numeric"`
	Password string `json:"password" binding:"required,min=8"`
}
