package handler

import (
	"fmt"
	"net/http"

	"github.com/chazool/go-sample-app/common/app/routes/handler/validator"
	responsebuilder "github.com/chazool/go-sample-app/common/app/routes/responseBuilder"
	"github.com/chazool/go-sample-app/common/app/services"
	"github.com/chazool/go-sample-app/common/pkg/common"
	"github.com/chazool/go-sample-app/common/pkg/utils"
	"github.com/chazool/go-sample-app/common/pkg/utils/constant"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Lives(ctx *fiber.Ctx) (err error) {
	var (
		requestID     = GetRequestID(ctx)
		parentSpan, _ = SetTracing(ctx, LivezMethod, requestID)
	)

	defer utils.CloseSpan(parentSpan)

	responseBuilder := responsebuilder.APIResponse{
		Ctx:        ctx,
		HttpStatus: fiber.StatusOK,
		Response:   nil,
		RequestID:  requestID,
	}
	responseBuilder.BuildAPIResponse()
	return nil
}

func Readyz(ctx *fiber.Ctx) (err error) {
	var (
		statusCode            int
		requestID             = GetRequestID(ctx)
		parentSpan, parentCtx = SetTracing(ctx, LivezMethod, requestID)
	)
	defer utils.CloseSpan(parentSpan)

	service := services.CreateHealthService(parentCtx, requestID, "")
	isReady, errRes := service.ReadyzService()

	if errRes != nil {
		statusCode = errRes.StatusCode
		errRes.IsError = true
	} else {
		statusCode = http.StatusOK
	}

	responseBuilder := responsebuilder.APIResponse{
		Ctx:           ctx,
		HttpStatus:    statusCode,
		ErrorResponse: *errRes,
		Response:      isReady,
		RequestID:     requestID,
	}
	responseBuilder.BuildAPIResponse()
	return nil
}

func TestConnection(ctx *fiber.Ctx) (err error) {
	var (
		statusCode      int
		response        *string
		errRes          *common.ErrorResult
		requestID       = GetRequestID(ctx)
		commonLogFields = []zap.Field{zap.String(constant.TraceMsgReqId, requestID)}
	)
	utils.Logger.Debug(fmt.Sprintf(constant.TraceMsgFuncStart, TestConnectionMethod), commonLogFields...)

	parentSpan, parentCtx := SetTracing(ctx, LivezMethod, requestID)

	defer func() {
		utils.CloseSpan(parentSpan)
		utils.Logger.Debug(fmt.Sprintf(constant.TraceMsgFuncEnd, TestConnectionMethod), commonLogFields...)
	}()

	//validate
	utils.Logger.Debug(fmt.Sprintf(constant.TraceMsgBeforeInvoke, validator.ValidateTestConnectionMethod), commonLogFields...)
	request, errRes := validator.ValidateTestConnection(requestID, ctx)
	utils.Logger.Debug(fmt.Sprintf(constant.TraceMsgAfterInvoke, validator.ValidateTestConnectionMethod), commonLogFields...)
	if errRes == nil {
		service := services.CreateHealthService(parentCtx, requestID, "")
		utils.Logger.Debug(fmt.Sprintf(constant.TraceMsgBeforeInvoke, services.TestConnectionMethod), commonLogFields...)
		response, errRes = service.TestConnection(request)
		utils.Logger.Debug(fmt.Sprintf(constant.TraceMsgAfterInvoke, services.TestConnectionMethod), append(commonLogFields, zap.Any(constant.Response, response), zap.Any(constant.ErrorNote, errRes))...)
	}

	if errRes != nil {
		utils.Logger.Error(fmt.Sprintf(constant.TraceMsgAfterInvoke, services.TestConnectionMethod), append(commonLogFields, zap.Any(constant.Response, response), zap.Any(constant.ErrorNote, errRes))...)
		statusCode = errRes.StatusCode
		errRes.IsError = true
	} else {
		statusCode = http.StatusOK
	}

	responseBuilder := responsebuilder.APIResponse{
		Ctx:           ctx,
		HttpStatus:    statusCode,
		ErrorResponse: *errRes,
		Response:      *response,
		RequestID:     requestID,
	}
	responseBuilder.BuildAPIResponse()
	return nil
}
