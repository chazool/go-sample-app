package validator

import (
	"net"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// RegistorCustomValidation use add custom validator
func RegistorCustomValidation(validate *validator.Validate) {

	validate.RegisterValidation(alpha, func(fl validator.FieldLevel) bool {
		return containsOnly(fl.Field().String(), alfaRegex)
	})

	validate.RegisterValidation(alphaNumaric, func(fl validator.FieldLevel) bool {
		return containsOnly(fl.Field().String(), alphaNumaricRegex)
	})

	validate.RegisterValidation(host, func(fl validator.FieldLevel) bool {
		if ip := net.ParseIP(fl.Field().String()); ip == nil {
			return false
		}
		return true
	})

	validate.RegisterValidation(port, func(fl validator.FieldLevel) bool {
		port := fl.Field().Int()
		if port < 0 || port > 65535 {
			return false
		}
		return true
	})

}

// RegisterCustomTranlation use add custom validater translation
func RegisterCustomTranlation(validate *validator.Validate, trans ut.Translator) {

	validate.RegisterTranslation(alpha, trans, func(ut ut.Translator) error {
		return ut.Add(alpha, "{0} can only contain Alpha characters", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(alpha, fe.Field())
		return t
	})

	validate.RegisterTranslation(alphaNumaric, trans, func(ut ut.Translator) error {
		return ut.Add(alphaNumaric, "{0} can only contain Alpha-Numaric  characters", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(alphaNumaric, fe.Field())
		return t
	})

	validate.RegisterTranslation(host, trans, func(ut ut.Translator) error {
		return ut.Add(host, "{0} is invalid Host", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(host, fe.Field())
		return t
	})

	validate.RegisterTranslation(port, trans, func(ut ut.Translator) error {
		return ut.Add(port, "{0} is invalid Port, can only only contain 1 to 65535 values only", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(port, fe.Field())
		return t
	})
 
}
