package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zhitoo/golang-web-api/config"
)

type AuthService struct {
	cfg *config.Config
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{cfg: cfg}
}

type tokenDto struct {
	UserId   string   `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

type AuthToken struct {
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	AccessExpireAt  int64  `json:"access_expire_at"`
	RefreshExpireAt int64  `json:"refresh_expire_at"`
}

func (s *AuthService) GenerateToken(token *tokenDto) (*AuthToken, error) {
	authToken := &AuthToken{}
	authToken.AccessExpireAt = time.Now().Add(s.cfg.JWT.AccessTokenExpireDuration * time.Second).Unix()
	authToken.RefreshExpireAt = time.Now().Add(s.cfg.JWT.RefreshTokenExpireDuration * time.Second).Unix()

	accessTokenClaims := jwt.MapClaims{}

	accessTokenClaims["user_id"] = token.UserId
	accessTokenClaims["username"] = token.Username
	accessTokenClaims["roles"] = token.Roles
	accessTokenClaims["exp"] = authToken.AccessExpireAt

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	var err error
	authToken.AccessToken, err = at.SignedString([]byte(s.cfg.JWT.Secret))

	if err != nil {
		return nil, err
	}

	refreshTokenClaims := jwt.MapClaims{}

	refreshTokenClaims["user_id"] = token.UserId
	refreshTokenClaims["username"] = token.Username
	refreshTokenClaims["roles"] = token.Roles
	refreshTokenClaims["exp"] = authToken.RefreshExpireAt

	rft := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	authToken.RefreshToken, err = rft.SignedString([]byte(s.cfg.JWT.RefreshSecret))

	if err != nil {
		return nil, err
	}

	return authToken, nil
}

func (s *AuthService) VerifyToken(token string) (*jwt.Token, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, fmt.Errorf("unexpected error duraing verify auth token")
		}

		return []byte(s.cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *AuthService) GetClaims(token string) (claimMap map[string]any, err error) {
	claimMap = map[string]any{}
	verifiedToken, err := s.VerifyToken(token)

	if err != nil {
		return
	}

	claims, ok := verifiedToken.Claims.(jwt.MapClaims)

	if ok && verifiedToken.Valid {
		for k, v := range claims {
			claimMap[k] = v
		}
		return
	}

	return nil, fmt.Errorf("claim not found")

}
