package validator

import (
	"github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate         *validator.Validate
	trans            ut.Translator
	GeneralErrorCode map[string]string
)

func init() {
	validate = validator.New()
	SetTransLatorForStructError(validate)

}

// SetTransLatorForStructError set the translator for the struct error
func SetTransLatorForStructError(validate *validator.Validate) (ut.Translator, error) {
	uni := ut.New(en_US.New())
	translator, _ := uni.GetTranslator("en_US")
	validationErr := translations.RegisterDefaultTranslations(validate, translator)
	return translator, validationErr
}
