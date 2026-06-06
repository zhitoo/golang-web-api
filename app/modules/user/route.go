package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhitoo/golang-web-api/app/modules/otp"
	"github.com/zhitoo/golang-web-api/app/modules/user/models"
	"github.com/zhitoo/golang-web-api/app/response"
	"github.com/zhitoo/golang-web-api/config"
	"github.com/zhitoo/golang-web-api/database/db"
	"github.com/zhitoo/golang-web-api/pkg/logging"
)

var log logging.ScopedLogger

func init() {
	cfg := config.GetConfig()
	log = logging.NewLogger(cfg).With(logging.Internal, logging.Api)
}

func Routes(r *gin.RouterGroup) {

	r.GET("", func(c *gin.Context) {
		response.NewReponse().SetResult("users list").Json(c)
	})

	r.POST("/login", func(c *gin.Context) {

		request := new(LoginRequest)
		if err := c.ShouldBindJSON(request); err != nil {
			log.Warn("invalid login request", map[logging.ExtraKey]any{logging.ErrorMessage: err.Error()})
			response.NewReponse().SetError(err).SetHttpStatusCode(http.StatusUnprocessableEntity).Json(c)
			return
		}

		if err := otp.VerifyOTP(request.Mobile, request.OTP); err != nil {
			log.Warn("OTP verification failed on login", map[logging.ExtraKey]any{logging.ErrorMessage: err.Error()})
			response.NewReponse().SetError(err).SetHttpStatusCode(http.StatusUnauthorized).Json(c)
			return
		}

		gorm := db.GetDb()

		var user models.User
		result := gorm.Where("mobile = ?", request.Mobile).First(&user)

		if result.Error != nil {
			log.Info("user not found, creating new user", nil)
			user = models.User{MobileNumber: request.Mobile, Username: request.Mobile, Password: "createAndHashARandomPassword"}
			if err := gorm.Create(&user).Error; err != nil {
				log.Error("failed to create user", map[logging.ExtraKey]any{logging.ErrorMessage: err.Error()})
				response.NewReponse().SetError(err).SetHttpStatusCode(http.StatusInternalServerError).Json(c)
				return
			}
		}

		// send auth token to user

		response.NewReponse().SetResult("welcom").Json(c)
	})
}
