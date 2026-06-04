package helper

import "carsale/app/validations"

type BaseHttpResponse struct {
	Result           any                            `json:"result"`
	Success          bool                           `json:"success"`
	ResultCode       int                            `json:"result_code"`
	ValidationErrors *[]validations.ValidationError `json:"validation_errors,omitempty"`
	Error            any                            `json:"error,omitempty"`
}

func GenerateBaseResponse(result any, success bool, resultCode int) *BaseHttpResponse {
	return &BaseHttpResponse{
		Result:     result,
		Success:    success,
		ResultCode: resultCode,
	}
}

func GenerateBaseResponseWithError(result any, success bool, resultCode int, err error) *BaseHttpResponse {
	output := GenerateBaseResponse(result, success, resultCode)
	output.Error = err
	return output
}
func GenerateBaseResponseWithValidationError(result any, success bool, resultCode int, err error) *BaseHttpResponse {
	output := GenerateBaseResponseWithError(result, success, resultCode, err)
	output.Error = err
	output.ValidationErrors = validations.GetValidationErrors(err)
	return output
}
