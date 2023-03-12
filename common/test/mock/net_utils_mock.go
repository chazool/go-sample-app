package mock

import (
	"encoding/json"
	"net"

	util "github.com/chazool/go-sample-app/common/pkg/utils"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type NetUtilsMock struct {
	mock.Mock
}

var netutils = util.NetUtilsImpl{NetUtils: &util.NetUtilsService{}}.NetUtils

func (utils *NetUtilsMock) LookupIP(endpoint string, commonLogFields []zapcore.Field) ([]net.IP, error) {
	args := utils.Called(endpoint, commonLogFields)
	var ips []net.IP
	return ips, args.Error(1)
}

func (utils *NetUtilsMock) Dial(url string, commonLogFields []zapcore.Field) (net.Conn, error) {
	args := utils.Called(url, commonLogFields)
	_, client := net.Pipe()
	return client, args.Error(1)
}

func (utils *NetUtilsMock) HttpRequest(request util.Request, commonLogFields []zapcore.Field) (response util.Response, err error) {

	args := utils.Called(request, commonLogFields)

	if args.Bool(2) {
		var (
			resMap = args.Get(0).(map[string]interface{})
			//headers    = resMap["headers"].(map[string]string)
			code    = resMap["code"].(int)
			boBytes = resMap["Body"].([]byte)
			body    = &[]byte{}
		)
		json.Unmarshal(boBytes, body)
		response = util.Response{Code: code, Body: *body}
	} else {
		response, err = netutils.HttpRequest(request, commonLogFields)
	}
	util.Logger.Debug("HttpRequest netutilsMock details", zap.Any("response", response), zap.Error(err))
	return response, err
}
