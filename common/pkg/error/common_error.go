package error

import (
	"net/http"
	"strings"
)

// ErrorResult used to define error result of the response
type ErrorResult struct {
	ErrorList  []ErrorInfo
	StatusCode int  `json:"_"`
	IsError    bool `json:"_"`
}

// ErrorInfo use to define error information of ther ErrorResult
type ErrorInfo struct {
	ErrorCode    string `json:"ErrorCode" example:"ER0001`
	ErrorMessage string `json:"ErrorMessage" example:"Records not found`
	ErrorDetail  string `json:"ErrorDetail" example:"XYZ dataa not availble in db`
}

// BilddErrorInfo used to build error information
func BuildErrorInfo(errCode, errMessage, errDetail string) ErrorInfo {
	return ErrorInfo{
		ErrorCode:    errCode,
		ErrorMessage: errMessage,
		ErrorDetail:  errDetail,
	}
}

// BuildErrResultWithSuccessStatus used to build ErrorResult with success code
func BuildErrResultWithSuccessStatus(errCode, errMessage, errDetail string) ErrorResult {
	errList := []ErrorInfo{BuildErrorInfo(errCode, errMessage, errDetail)}
	return ErrorResult{
		ErrorList:  errList,
		IsError:    false,
		StatusCode: http.StatusOK,
	}
}

// BuilBbadReqErrResultWirhList used to build ErrorResult with bErrorinfo list and bad request code
func BuilBbadReqErrResultWirhList(errInfo ErrorInfo) ErrorResult {
	return ErrorResult{
		ErrorList:  []ErrorInfo{},
		IsError:    false,
		StatusCode: http.StatusBadRequest,
	}
}

// BuildErrResultWithSuccessStatus used to build ErrorResult with bad request code
func BuilBbadReqErrResult(errCode, errMessage, errDetail string) ErrorResult {
	errList := []ErrorInfo{BuildErrorInfo(errCode, errMessage, errDetail)}
	return ErrorResult{
		ErrorList:  errList,
		IsError:    false,
		StatusCode: http.StatusBadRequest,
	}
}

// BuildNotFoundErrResult used to build ErrorResult with bad request code
func BuildNotFoundErrResult(errCode, errMessage, errDetail string) ErrorResult {
	errList := []ErrorInfo{BuildErrorInfo(errCode, errMessage, errDetail)}
	return ErrorResult{
		ErrorList:  errList,
		IsError:    false,
		StatusCode: http.StatusNotFound,
	}
}

// BuildErrResultWithSuccessStatus used to build ErrorResult with bad request code
func BuildInternalServerErrResult(errCode, errMessage, errDetail string) ErrorResult {
	errList := []ErrorInfo{BuildErrorInfo(errCode, errMessage, errDetail)}
	return ErrorResult{
		ErrorList:  errList,
		IsError:    false,
		StatusCode: http.StatusInternalServerError,
	}
}

// GetErrorMessage use to retun error message from ErrorList
func GetErrorMessage(errorResult *ErrorResult) string {
	var errMessages []string
	for _, err := range errorResult.ErrorList {
		errMessages = append(errMessages, err.ErrorMessage)
	}
	return strings.Join(errMessages, ",")
}
