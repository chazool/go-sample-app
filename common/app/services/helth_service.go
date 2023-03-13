package services

import (
	"context"
	"fmt"
	"runtime/debug"
	"strconv"

	"github.com/chazool/go-sample-app/common/app/routes/dto"
	"github.com/chazool/go-sample-app/common/pkg/common"
	"github.com/chazool/go-sample-app/common/pkg/utils"
	"github.com/chazool/go-sample-app/common/pkg/utils/constant"
	"go.uber.org/zap"
)

type HealthService struct {
	_              struct{}
	serviceContext ServiceContext
	netutils       *utils.NetUtilsImpl
}

var (
	netutils = utils.NetUtilsImpl{NetUtils: &utils.NetUtilsService{}}
)

func CreateHealthService(parentCtx context.Context, requestID string, request string) HealthService {
	serviceContext := CreateServiceContext(parentCtx, requestID, request)
	return HealthService{serviceContext: serviceContext, netutils: &netutils}
}

// getNetUtils used to return the Netutils from implemented service
func (service *HealthService) getNetUtils() utils.NetUtils {
	return service.netutils.NetUtils
}

func (service *HealthService) ReadyzService() (isReady bool, errResult *common.ErrorResult) {
	commonLogFields := []zap.Field{zap.String(constant.TraceMsgReqId, service.serviceContext.RequestID)}
	childSpan, _ := utils.CreateChildSpan(service.serviceContext.ParentCtx, ReadyzServiceMethod)
	utils.Logger.Debug(fmt.Sprintf(constant.TraceMsgFuncStart, ReadyzServiceMethod), commonLogFields...)

	defer func() {
		utils.CloseSpan(childSpan)
		utils.Logger.Debug(fmt.Sprintf(constant.TraceMsgFuncEnd, ReadyzServiceMethod), commonLogFields...)
	}()

	isReady = true

	return isReady, nil
}

func (service *HealthService) TestConnection(request dto.HostPost) (response *string, errResult *common.ErrorResult) {
	commonLogFields := []zap.Field{zap.String(constant.TraceMsgReqId, service.serviceContext.RequestID)}
	utils.Logger.Debug(fmt.Sprintf(constant.TraceMsgFuncStart, TestConnectionMethod), commonLogFields...)

	var (
		childSpan, _ = utils.CreateChildSpan(service.serviceContext.ParentCtx, TestConnectionMethod)
		port         = strconv.Itoa(request.Port)
		url          = request.Host + ":" + port
	)

	defer func() {
		utils.CloseSpan(childSpan)
		if r := recover(); r != nil {
			utils.Logger.Error(constant.PanicOccurred, append(commonLogFields, []zap.Field{zap.String(constant.ErrorRequestBody, service.serviceContext.RequestBody), zap.String(constant.StackTrace, string(debug.Stack()))}...)...)
			errRes := common.BuildInternalServerErrResult(constant.UnexpectedErrorCode, fmt.Sprintf(constant.UnexpectedErrorMessage, TestConnectionMethod), "")
			errResult = &errRes
		}
		utils.Logger.Debug(fmt.Sprintf(constant.TraceMsgFuncEnd, TestConnectionMethod), commonLogFields...)
	}()

	utils.Logger.Debug(fmt.Sprintf(constant.TraceMsgBeforeInvoke, constant.Dial), append(commonLogFields, zap.String(serverAndPort, url))...)
	connection, err := service.getNetUtils().Dial(url, commonLogFields)
	utils.Logger.Debug(fmt.Sprintf(constant.TraceMsgBeforeInvoke, constant.Dial), commonLogFields...)

	if err != nil {
		utils.Logger.Error(CouldNotGetConnectionToSocket, append(commonLogFields, []zap.Field{zap.String(serverAndPort, url), zap.Error(err)}...)...)
		errRes := common.BuildBadReqErrResult(constant.TestConnectionFilCode, constant.TestConnectionFilMessage, err.Error())
		return nil, &errRes
	}
	defer connection.Close()

	res := fmt.Sprintf(TestConnectionSuccess, port, request.Host)
	return &res, nil
}
