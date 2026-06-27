package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhitoo/golang-web-api/app/validations"
)

type BaseHttpResponse struct {
	Result           any                            `json:"result"`
	Success          bool                           `json:"success"`
	HttpStatusCode   int                            `json:"-"`
	ValidationErrors *[]validations.ValidationError `json:"validation_errors,omitempty"`
	Error            any                            `json:"error,omitempty"`
}

func NewResponse() *BaseHttpResponse {
	res := new(BaseHttpResponse)
	res.Success = true
	res.HttpStatusCode = 200
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
	bh.HttpStatusCode = httpStatusCode

	return bh
}

func (bh *BaseHttpResponse) SetError(err error) *BaseHttpResponse {
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

func (bh *BaseHttpResponse) Json(c *gin.Context) {
	if bh.Error != nil {
		c.AbortWithStatusJSON(bh.HttpStatusCode, bh)
		return
	}
	c.JSON(bh.HttpStatusCode, bh)
}
