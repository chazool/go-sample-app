package constant

const (
	DatadogTracingSink = "datadog"
)

// special charactes
const (
	colon    = ":"
	Basepath = "./"
	Empty    = ""
	Hyphen   = "-"
)

// file constant
const (
	DocumentHtml  = "document.html"
	StaticHtml    = "static.html"
	Doc           = "doc"
	Static        = "static"
	HTML          = "html"
	DotJson       = ".json"
	DotGoldenJson = ".golden.json"
)

// Logger Message
const (
	TraceMsgBeforeFetching   = "Before Fetching %v"
	TraceMsgAfterFetching    = "After Fetching %v"
	TraceMsgBeforeInsertion  = "Before Creating %v"
	TraceMsgAfterInsertion   = "After Creating %v"
	TraceMsgBeforeUpdate     = "Before Update %v"
	TraceMsgAfterUpdate      = "After Update %v"
	TraceMsgFuncEnd          = "%v End here"
	TraceMsgFuncStart        = "%v Start here"
	TraceMsgRequestInitiated = "%v request initiated"
	TraceMsgReqId            = "Request Id"
	TraceMsgReqBody          = "Request Body"
	TraceMsgRequestHeader    = "Request Header"
	TranceMsgResponse        = "Response"
	TraceMsgBeforeInvoke     = "Before Call %v"
	TraceMsgAfterInvoke      = "After Call %v"
	TraceMsgAPIResponse      = "Build API Response"
	TraceMsgResponseDetails  = "Response Details"
	TraceMsgAPISuccess       = "Success Response"
	TraceMsgAPIErrorResponse = "Error Response"
	TraceMsgErrorDetails     = "Error Details"
	MethodInput              = "Method Input %v"
	Result                   = "Result"
	DebugNote                = "Debug workflow"
	ErrorNote                = "Error Note"
	HTMLPassErr              = "HTML Template pass Error"
	Response                 = "Response"
)

// utils func
const (
	Dial = "Dial"
)

type HttpMethod string

const (
	Get, Post, Patch, Delete HttpMethod = "GET", "POST", "PATCH", "DELETE"
)
