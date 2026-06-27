package user

import (
	"fmt"
	"strconv"

	authSvc "github.com/zhitoo/golang-web-api/app/modules/auth"
	otpSvc "github.com/zhitoo/golang-web-api/app/modules/otp"
	"github.com/zhitoo/golang-web-api/app/modules/user/models"
	"github.com/zhitoo/golang-web-api/config"
	"github.com/zhitoo/golang-web-api/database/db"
	"github.com/zhitoo/golang-web-api/pkg/utils"
)

type UserService struct {
	cfg         *config.Config
	authService *authSvc.AuthService
}

func NewUserService(cfg *config.Config) *UserService {
	return &UserService{
		cfg:         cfg,
		authService: authSvc.NewAuthService(cfg),
	}
}

func (s *UserService) Login(mobile, email, otpCode, password string) (*authSvc.AuthToken, error) {
	if otpCode == "" && password == "" {
		return nil, fmt.Errorf("otp or password is required")
	}

	gorm := db.GetDb()
	var u models.User

	if mobile != "" {
		if password != "" {
			if err := gorm.Where("mobile = ?", mobile).First(&u).Error; err != nil {
				return nil, fmt.Errorf("invalid credentials")
			}
			if !utils.VerifyPassword(password, u.Password) {
				return nil, fmt.Errorf("invalid credentials")
			}
		} else {
			if err := otpSvc.VerifyOTP(mobile, otpCode); err != nil {
				return nil, err
			}
			if err := gorm.Where("mobile = ?", mobile).First(&u).Error; err != nil {
				pwd, err := utils.HashPassword(utils.RandomString(16))
				if err != nil {
					return nil, err
				}
				u = models.User{Mobile: mobile, Password: pwd}
				if err := gorm.Create(&u).Error; err != nil {
					return nil, err
				}
			}
		}
	} else {
		if err := gorm.Where("email = ?", email).First(&u).Error; err != nil {
			return nil, fmt.Errorf("invalid credentials")
		}
		if password != "" {
			if !utils.VerifyPassword(password, u.Password) {
				return nil, fmt.Errorf("invalid credentials")
			}
		} else {
			if u.Mobile == "" {
				return nil, fmt.Errorf("no mobile number linked to this account")
			}
			if err := otpSvc.VerifyOTP(u.Mobile, otpCode); err != nil {
				return nil, err
			}
		}
	}

	roles := []string{}
	if u.UserRoles != nil {
		for _, r := range *u.UserRoles {
			roles = append(roles, r.Role.Name)
		}
	}

	return s.authService.GenerateToken(&authSvc.TokenDto{
		UserId:   strconv.Itoa(u.Id),
		Username: u.Mobile,
		Roles:    roles,
	})
}

func (s *UserService) RefreshToken(refreshToken string) (*authSvc.AuthToken, error) {
	claims, err := s.authService.GetRefreshClaims(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	userId, ok := claims["user_id"].(string)
	if !ok || userId == "" {
		return nil, fmt.Errorf("invalid token claims")
	}

	roles := []string{}
	if r, ok := claims["roles"].([]any); ok {
		for _, v := range r {
			if s, ok := v.(string); ok {
				roles = append(roles, s)
			}
		}
	}

	username, _ := claims["username"].(string)

	return s.authService.GenerateToken(&authSvc.TokenDto{
		UserId:   userId,
		Username: username,
		Roles:    roles,
	})
}

func (s *UserService) ChangePassword(userId int, oldPassword, newPassword string) error {
	gorm := db.GetDb()
	var u models.User
	if err := gorm.First(&u, userId).Error; err != nil {
		return fmt.Errorf("user not found")
	}
	if !utils.VerifyPassword(oldPassword, u.Password) {
		return fmt.Errorf("incorrect current password")
	}
	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	return gorm.Model(&u).Update("password", hash).Error
}

func (s *UserService) ResetPassword(mobile, email, otpCode, newPassword string) error {
	gorm := db.GetDb()
	var u models.User

	var mobileToVerify string
	if mobile != "" {
		if err := gorm.Where("mobile = ?", mobile).First(&u).Error; err != nil {
			return fmt.Errorf("user not found")
		}
		mobileToVerify = mobile
	} else {
		if err := gorm.Where("email = ?", email).First(&u).Error; err != nil {
			return fmt.Errorf("user not found")
		}
		if u.Mobile == "" {
			return fmt.Errorf("no mobile number linked to this account")
		}
		mobileToVerify = u.Mobile
	}

	if err := otpSvc.VerifyOTP(mobileToVerify, otpCode); err != nil {
		return err
	}

	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	return gorm.Model(&u).Update("password", hash).Error
}

func (s *UserService) GetById(id int) (*models.User, error) {
	gorm := db.GetDb()
	var u models.User
	if err := gorm.Preload("UserRoles.Role").First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
