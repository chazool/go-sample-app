package validator

import (
	"github.com/chazool/go-sample-app/common/pkg/common"
	"github.com/chazool/go-sample-app/common/pkg/utils/constant"
	"github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate         *validator.Validate
	trans            ut.Translator
	generalErrorCode map[string]string
)

func init() {
	validate = validator.New()
	trans, _ = SetTransLatorForStructError(validate)
	RegisterCustomValidation(validate)
	RegisterCustomTranslation(validate, trans)
	generalErrorCode = BuildGeneralErrorCode()
}

func BuildValidationErrorResponse(requestId string, validationError error) *common.ErrorResult {
	if validationError != nil {
		errorList := []common.ErrorInfo{}
		for _, validationErrorsTranslation := range validationError.(validator.ValidationErrors) {
			errorList = append(errorList, common.BuildErrorInfo(generalErrorCode[validationErrorsTranslation.Tag()], validationErrorsTranslation.Translate(trans), ""))
		}
		err := common.BuildBadReqErrResultWithList(errorList...)
		return &err
	}
	return nil
}

// SetTransLatorForStructError set the translator for the struct error
func SetTransLatorForStructError(validate *validator.Validate) (ut.Translator, error) {
	uni := ut.New(en_US.New())
	translator, _ := uni.GetTranslator("en_US")
	validationErr := translations.RegisterDefaultTranslations(validate, translator)
	return translator, validationErr
}

// BuildGeneralErrorCode used to validate general error code
func BuildGeneralErrorCode() map[string]string {

	commonErrorMap := make(map[string]string)
	commonErrorMap["required"] = constant.MissingRequiredFieldErrorCode
	commonErrorMap["required_without"] = constant.MissingRequireWithoutFieldCode
	commonErrorMap["required_with"] = constant.MissingRequireWithFieldCode
	commonErrorMap["min"] = constant.MinLengthFieldCode
	commonErrorMap["max"] = constant.MaxLengthFieldCode
	commonErrorMap[alpha] = constant.PatternErrorCode
	commonErrorMap[alphaNumeric] = constant.PatternErrorCode

	return commonErrorMap
}
