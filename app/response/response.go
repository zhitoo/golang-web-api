package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhitoo/golang-web-api/app/validations"
)

type baseHttpResponse struct {
	Result           any                            `json:"result"`
	Success          bool                           `json:"success"`
	HttpStatusCode   int                            `json:"-"`
	ValidationErrors *[]validations.ValidationError `json:"validation_errors,omitempty"`
	Error            any                            `json:"error,omitempty"`
}

func NewResponse() *baseHttpResponse {
	res := new(baseHttpResponse)
	res.Success = true
	res.HttpStatusCode = 200
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
	bh.HttpStatusCode = httpStatusCode

	return bh
}

func (bh *baseHttpResponse) SetError(err error) *baseHttpResponse {
	bh.Error = err.Error()

	bh.ValidationErrors = validations.GetValidationErrors(err)

	if bh.Success {
		bh.Success = false
	}

	if bh.HttpStatusCode == http.StatusOK {
		bh.HttpStatusCode = http.StatusBadRequest
	}

	return bh
}

func (bh *baseHttpResponse) Json(c *gin.Context) {
	if bh.Error != nil {
		c.AbortWithStatusJSON(bh.HttpStatusCode, bh)
		return
	}
	c.JSON(bh.HttpStatusCode, bh)
}
