package dto

import (
	english "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type BaseResponse struct {
	Code    Code        `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func BuildBaseResponse(code Code, data interface{}) *BaseResponse {
	return &BaseResponse{
		Code:    code,
		Message: CodeMessage[code],
		Data:    data,
	}
}

func BuildValidationErrorResponse(err error) *BaseResponse {
	eng := english.New()
	uni := ut.New(eng, eng)
	trans, _ := uni.GetTranslator("en")

	errors := []string{}
	for _, validationError := range err.(validator.ValidationErrors) {
		errors = append(errors, validationError.Translate(trans))
	}

	return &BaseResponse{
		Code:    BadRequest,
		Message: CodeMessage[BadRequest],
		Data:    errors,
	}
}
