package response

import (
	"github.com/gin-gonic/gin"
	"github.com/zhitoo/golang-web-api/app/validations"
)

type BaseHttpResponse struct {
	Result           any                            `json:"result"`
	Success          bool                           `json:"success"`
	HttpSatatusCode  int                            `json:"http_status_code"`
	ValidationErrors *[]validations.ValidationError `json:"validation_errors,omitempty"`
	Error            any                            `json:"error,omitempty"`
}

func NewReponse() *BaseHttpResponse {
	res := new(BaseHttpResponse)
	res.Success = true
	res.HttpSatatusCode = 200
	res.Error = nil
	res.ValidationErrors = nil
	return res
}

func (bh *BaseHttpResponse) SetResult(result any) *BaseHttpResponse {
	bh.Result = result

	return bh
}

func (bh *BaseHttpResponse) SetStatus(success bool) *BaseHttpResponse {
	bh.Success = success

	return bh
}

func (bh *BaseHttpResponse) SetHttpStatusCode(httpStatusCode int) *BaseHttpResponse {
	bh.HttpSatatusCode = httpStatusCode

	return bh
}

func (bh *BaseHttpResponse) SetError(err error) *BaseHttpResponse {
	bh.Error = err

	bh.ValidationErrors = validations.GetValidationErrors(err)

	return bh
}

func (bh *BaseHttpResponse) Json(c *gin.Context) {
	c.JSON(bh.HttpSatatusCode, bh)
}
