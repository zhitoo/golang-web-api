package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhitoo/golang-web-api/app/modules/otp"
	"github.com/zhitoo/golang-web-api/app/modules/user/models"
	"github.com/zhitoo/golang-web-api/app/response"
	"github.com/zhitoo/golang-web-api/database/db"
)

func Routes(r *gin.RouterGroup) {

	r.GET("", func(c *gin.Context) {
		response.NewReponse().SetResult("users list").Json(c)
	})

	r.POST("/login", func(c *gin.Context) {

		request := new(LoginRequest)

		err := c.ShouldBindJSON(request)
		if err != nil {
			response.NewReponse().SetError(err).SetHttpStatusCode(http.StatusUnprocessableEntity).Json(c)
			return
		}

		err = otp.VerifyOTP(request.Mobile, request.OTP)
		if err != nil {
			response.NewReponse().SetError(err).SetHttpStatusCode(http.StatusUnauthorized).Json(c)
			return
		}

		gorm := db.GetDb()

		var user models.User
		result := gorm.Where("mobile = ?", request.Mobile).First(&user)

		if result.Error != nil {
			// user not found, create new one
			user = models.User{MobileNumber: request.Mobile, Username: request.Mobile, Password: "createAndHashARandomPassword"}
			if err := gorm.Create(&user).Error; err != nil {
				response.NewReponse().SetError(err).SetHttpStatusCode(http.StatusInternalServerError).Json(c)
				return
			}
		}

		//send auth token to user

		response.NewReponse().SetResult("welcom").Json(c)
	})
}
