package validator

import (
	"fmt"

	"github.com/chazool/go-sample-app/common/app/routes/dto"
	"github.com/chazool/go-sample-app/common/pkg/common"
	"github.com/chazool/go-sample-app/common/pkg/utils"
	"github.com/chazool/go-sample-app/common/pkg/utils/constant"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ValidateTestConnection(requestId string, ctx *fiber.Ctx) (dto.HostPost, *common.ErrorResult) {

	commonLogFields := []zap.Field{zap.String(constant.TraceMsgReqId, requestId)}
	utils.Logger.Debug(fmt.Sprintf(constant.TraceMsgFuncStart, ValidateTestConnectionMethod), commonLogFields...)
	defer utils.Logger.Debug(fmt.Sprintf(constant.TraceMsgFuncEnd, ValidateTestConnectionMethod), commonLogFields...)

	var (
		request dto.HostPost
		body    = string(ctx.Body())
		err     = ctx.QueryParser(&request)
	)

	if err != nil {
		utils.Logger.Error(constant.InvalidInputAndPassErr, append(commonLogFields, []zap.Field{zap.String(constant.ErrorRequestBody, body), zap.Error(err)}...)...)
		errorResult := common.BuildBadReqErrResult(constant.BindingErrorMessage, constant.BindingErrorMessage, err.Error())
		return request, &errorResult
	}

	err = validate.Struct(request)
	if err != nil {
		utils.Logger.Error(constant.InvalidInputAndPassErr, append(commonLogFields, []zap.Field{zap.String(constant.ErrorRequestBody, body), zap.Error(err)}...)...)
		return request, BuildValidationErrorResponse(requestId, err)
	}

	return request, nil
}
