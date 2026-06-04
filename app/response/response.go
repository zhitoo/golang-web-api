package response

import (
	"github.com/gin-gonic/gin"
	"github.com/zhitoo/golang-web-api/app/validations"
)

type baseHttpResponse struct {
	Result           any                            `json:"result"`
	Success          bool                           `json:"success"`
	HttpSatatusCode  int                            `json:"http_status_code"`
	ValidationErrors *[]validations.ValidationError `json:"validation_errors,omitempty"`
	Error            any                            `json:"error,omitempty"`
}

func NewReponse() *baseHttpResponse {
	res := new(baseHttpResponse)
	res.Success = true
	res.HttpSatatusCode = 200
	res.Error = nil
	res.ValidationErrors = nil
	return res
}

func (bh *baseHttpResponse) SetResult(result any) *baseHttpResponse {
	bh.Result = result

	return bh
}

func (bh *baseHttpResponse) SetStatus(success bool) *baseHttpResponse {
	bh.Success = success

	return bh
}

func (bh *baseHttpResponse) SetHttpStatusCode(httpStatusCode int) *baseHttpResponse {
	bh.HttpSatatusCode = httpStatusCode

	return bh
}

func (bh *baseHttpResponse) SetError(err error) *baseHttpResponse {
	bh.Error = err.Error()

	bh.ValidationErrors = validations.GetValidationErrors(err)

	return bh
}

func (bh *baseHttpResponse) Json(c *gin.Context) {
	c.JSON(bh.HttpSatatusCode, bh)
}
