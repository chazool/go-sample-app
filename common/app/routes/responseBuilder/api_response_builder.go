package responsebuilder

import (
	"net/http"

	"github.com/chazool/go-sample-app/common/pkg/common"
	"github.com/chazool/go-sample-app/common/pkg/utils"
	"github.com/chazool/go-sample-app/common/pkg/utils/constant"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// APIResponse use to define api response
type APIResponse struct {
	_             struct{}
	Ctx           *fiber.Ctx
	Response      interface{}
	RequestID     string
	HttpStatus    int
	ErrorResponse common.ErrorResult
}

func (response *APIResponse) BuildAPIResponse() {
	commonLogFields := []zap.Field{zap.String(constant.TraceMsgReqId, response.RequestID)}
	utils.Logger.Debug(constant.TraceMsgAPIResponse, commonLogFields...)

	if response.ErrorResponse.IsError {
		utils.Logger.Debug(constant.TraceMsgAPISuccess, append(commonLogFields, zap.Any(constant.TraceMsgReqBody, response.ErrorResponse))...)
		if response.HttpStatus == 0 {
			response.HttpStatus = http.StatusInternalServerError
		}
	} else {
		utils.Logger.Debug(constant.TraceMsgAPISuccess, commonLogFields...)
		response.Ctx.Status(response.HttpStatus).JSON(response.Response)
	}
}
