package constant

//error message for log
const (
	PanicOccurred            = "Error or panic occurred with following stacktrace"
	stackTrace               = "Error stacktrace"
	ErrorRequestBody         = "Error Request"
	InvalidInputAndPassErr   = "Error input provided is invalid & unable to parse"
	MissingRequiredField     = "Missing Required Field"
	ErrorOccurredFromService = "Missing Occurred from  %s"
	ErrorOccurredFromMethod  = "Error Occurred from  Method %s"
	BindingErrorMessage      = "Error Occurred when bind the request context"
)

//error codes
const (
	UnexpectedErrorCode            = "0000"
	BindingErrorCode               = "0001"
	MissingRequiredFieldErrorCode  = "0002"
	MissingRequireWithoutFieldCode = "0003"
	MissingRequireWithFieldCode    = "0004"
	MinLengthFieldCode             = "0005"
	MaxLengthFieldCode             = "0006"
	GreaterValueFieldCode          = "0007"
	PatternErrorCode               = "0008"
)
