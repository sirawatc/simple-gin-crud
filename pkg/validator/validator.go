package validator

import (
	english "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) Validate(i interface{}) []string {
	err := v.validate.Struct(i)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return v.TranslateErrors(validationErrors)
		}
		return []string{err.Error()}
	}
	return nil
}

func (v *Validator) TranslateErrors(validationErrors validator.ValidationErrors) []string {
	eng := english.New()
	uni := ut.New(eng, eng)
	trans, _ := uni.GetTranslator("en")
	err := en.RegisterDefaultTranslations(v.validate, trans)
	if err != nil {
		return []string{err.Error()}
	}
	errors := []string{}

	for _, validationError := range validationErrors {
		errors = append(errors, validationError.Translate(trans))
	}
	return errors
}
