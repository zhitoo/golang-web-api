package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhitoo/golang-web-api/app/response"
	"github.com/zhitoo/golang-web-api/pkg/logging"
)

type userResponse struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Mobile    string `json:"mobile"`
	Email     string `json:"email"`
}

func loginHandler(svc *UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(LoginRequest)
		if err := c.ShouldBindJSON(req); err != nil {
			log.Warn("invalid login request", map[logging.ExtraKey]any{logging.ErrorMessage: err.Error()})
			response.NewResponse().SetError(err).SetHttpStatusCode(http.StatusUnprocessableEntity).Json(c)
			return
		}

		token, err := svc.Login(req.Mobile, req.Email, req.OTP)
		if err != nil {
			log.Warn("login failed", map[logging.ExtraKey]any{logging.ErrorMessage: err.Error()})
			response.NewResponse().SetError(err).SetHttpStatusCode(http.StatusUnauthorized).Json(c)
			return
		}

		response.NewResponse().SetResult(token).Json(c)
	}
}

func refreshTokenHandler(svc *UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(RefreshTokenRequest)
		if err := c.ShouldBindJSON(req); err != nil {
			response.NewResponse().SetError(err).SetHttpStatusCode(http.StatusUnprocessableEntity).Json(c)
			return
		}

		token, err := svc.RefreshToken(req.RefreshToken)
		if err != nil {
			log.Warn("token refresh failed", map[logging.ExtraKey]any{logging.ErrorMessage: err.Error()})
			response.NewResponse().SetError(err).SetHttpStatusCode(http.StatusUnauthorized).Json(c)
			return
		}

		response.NewResponse().SetResult(token).Json(c)
	}
}

func meHandler(svc *UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdVal, exists := c.Get("UserId")
		if !exists {
			response.NewResponse().SetError(fmt.Errorf("unauthorized")).SetHttpStatusCode(http.StatusUnauthorized).Json(c)
			return
		}

		userId, err := strconv.Atoi(userIdVal.(string))
		if err != nil {
			response.NewResponse().SetError(fmt.Errorf("invalid user id")).SetHttpStatusCode(http.StatusBadRequest).Json(c)
			return
		}

		u, err := svc.GetById(userId)
		if err != nil {
			log.Error("failed to get user", map[logging.ExtraKey]any{logging.ErrorMessage: err.Error()})
			response.NewResponse().SetError(fmt.Errorf("user not found")).SetHttpStatusCode(http.StatusNotFound).Json(c)
			return
		}

		response.NewResponse().SetResult(userResponse{
			Id:        u.Id,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Mobile:    u.Mobile,
			Email:     u.Email,
		}).Json(c)
	}
}
