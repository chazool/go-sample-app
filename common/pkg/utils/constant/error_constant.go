package constant

//error message for log
const (
	PanicOccurred            = "Error or panic occurred with following stacktrace"
	StackTrace               = "Error stacktrace"
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
	TestConnectionFilCode          = "0009"
)

const (
	InvalidRequestErrorMessage   = "Invalid Request validation Error Occurred"
	UnexpectedErrorMessage       = "Unexpected Error occurred at %s"
	UnexpectedWhenUnmarshalError = "Unexpected Error occurred when Unmarshal the data &s with identifier %s"
	UnexpectedFileCreateError    = "Unexpected Error occurred when Create the % file"
	TestConnectionFilMessage     = "Error Occurred when TestConnection"
	InvalidDataOrFile            = "Invalid Data or File "
	InvalidJsonTestData          = "Invalid Json test data"
)
