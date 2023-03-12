package mock

import (
	"encoding/json"
	"fmt"

	"github.com/chazool/go-sample-app/common/pkg/utils"
	"github.com/chazool/go-sample-app/common/pkg/utils/constant"
	"github.com/fatih/structs"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func LookupIP(utils *NetUtilsMock, funcStep FuncStep) {
	var (
		params, returnData = funcStep.Params, funcStep.ReturnData
		err                = getError(returnData.Error)
		url                = params.Args["EndPoint"].(string)
	)
	utils.On(lookupIP, url, mock.Anything).Return(mock.Anything, err).Once()
}

func Dial(utils *NetUtilsMock, funcStep FuncStep) {
	var (
		params, returnData = funcStep.Params, funcStep.ReturnData
		err                = getError(returnData.Error)
		url                = params.Args["ServerAndPort"].(string)
	)

	utils.On(dial, url, mock.Anything).Return(mock.Anything, err).Once()
}

func HttpRequest(mockutils *NetUtilsMock, funcStep FuncStep) {
	var (
		err                error
		code               int
		body               []byte
		returnData, isMock = funcStep.ReturnData, funcStep.IsMock
	)

	if isMock {
		err = getError(returnData.Error)
		code = int(returnData.Outputs["code"].(float64))
		b := returnData.Outputs["body"]
		utils.Logger.Debug("Body Details", zap.Any("Body", b))
		body, _ = json.Marshal(b)
	}

	response := utils.Response{Code: code, Body: body}
	utils.Logger.Debug(fmt.Sprintf(mockDetails, httpRequest), zap.Any("Response", response), zap.Any(constant.ErrorNote, err))
	mockutils.On(httpRequest, mock.Anything, mock.Anything).Return(structs.Map(response), err, isMock).Once()
}
