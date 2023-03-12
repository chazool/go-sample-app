package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/chazool/go-sample-app/common/pkg/utils/constant"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Request struct {
	URL     string              `json:"URL"`
	Method  constant.HttpMethod `json:"Method"`
	Headers map[string]string   `json:"Headers"`
	Body    interface{}         `json:"Body"`
	TimeOut time.Duration       `json:"TimeOut"`
}

type Response struct {
	Code    int
	Body    []byte
	Headers map[string]string
}

type NetUtils interface {
	LookupIP(endpoint string, commonLogFields []zapcore.Field) ([]net.IP, error)
	Dial(url string, commonLogFields []zapcore.Field) (net.Conn, error)
	HttpRequest(request Request, commonLogFields []zapcore.Field) (response Response, err error)
}

type NetUtilsService struct{}

type NetUtilsImpl struct {
	NetUtils NetUtils
}

func (utils *NetUtilsService) LookupIP(endpoint string, commonLogFields []zapcore.Field) ([]net.IP, error) {
	Logger.Debug(fmt.Sprintf(constant.TraceMsgFuncStart, "LookupIp"), commonLogFields...)
	defer Logger.Debug(fmt.Sprintf(constant.TraceMsgFuncEnd, "LookupIp"), commonLogFields...)

	ips, err := net.LookupIP(endpoint)

	for _, ip := range ips {
		Logger.Debug(fmt.Sprintf("endpointIpResolution", "LookupIp"), append(commonLogFields, []zap.Field{zap.String("ip", ip.String())}...)...)
	}
	return ips, err
}

func (utils *NetUtilsService) Dial(url string, commonLogFields []zapcore.Field) (net.Conn, error) {
	Logger.Debug(fmt.Sprintf(constant.TraceMsgFuncStart, "Dial"), commonLogFields...)
	defer Logger.Debug(fmt.Sprintf(constant.TraceMsgFuncEnd, "Dial"), commonLogFields...)
	return net.Dial("tcp", url)
}

func (utils *NetUtilsService) HttpRequest(request Request, commonLogFields []zapcore.Field) (response Response, err error) {
	Logger.Debug(fmt.Sprintf(constant.TraceMsgFuncStart, "Dial"), commonLogFields...)

	var agent = fiber.AcquireAgent()

	defer func() {
		agent.ConnectionClose()
		Logger.Debug(fmt.Sprintf(constant.TraceMsgFuncEnd, "Dial"), commonLogFields...)
	}()

	// set request body
	if request.Body != nil {
		agent.JSON(request.Body)
	}

	// request timeout
	if request.TimeOut > 0 {
		agent = agent.Timeout(request.TimeOut)
	}

	req := agent.Request()
	req.SetRequestURI(request.URL)
	req.Header.SetMethod(string(request.Method))

	// request headers
	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	Logger.Debug("Parse initializes URI and HostClient HttpRequest", append(commonLogFields, []zap.Field{zap.Any("request", request)}...)...)
	if err := agent.Parse(); err != nil {
		Logger.Error("Error Input provided is invalid & unable to parse", append(commonLogFields, []zap.Field{zap.Any("request", request), zap.Error(err)}...)...)
		return response, err
	}

	// response
	resp := fiber.AcquireResponse()
	agent.SetResponse(resp)
	code, body, errs := agent.SetResponse(resp).Bytes()
	if errs != nil {
		Logger.Error("Error Input provided is invalid & unable to parse", append(commonLogFields, []zap.Field{zap.Any("request", request), zap.Errors("responseErrors", errs)}...)...)
		err := errors.New("Error")
		for _, e := range errs {
			err = errors.Wrap(err, e.Error())
		}
		return response, err
	}

	//response headers
	var (
		reader          = bytes.NewReader(resp.Header.Header())
		bReader         = bufio.NewReader(reader)
		responseHeaders = make(map[string]string)
	)

	for {
		line, err := bReader.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err != io.EOF {
				Logger.Error(fmt.Sprintf(constant.ErrorOccurredFromService, "reading response headers"), append(commonLogFields, []zap.Field{zap.Any("responseHeaders ", responseHeaders), zap.Error(err)}...)...)
				return response, err
			} else {
				Logger.Debug("response headers reading complete", append(commonLogFields, []zap.Field{zap.Any("responseHeaders ", responseHeaders)}...)...)
				break
			}
		}

		if strings.Contains(line, ":") {
			data := strings.Split(line, ":")
			responseHeaders[strings.TrimSpace(data[0])] = strings.TrimSpace(data[1])
		}
	}

	response = Response{Code: code, Body: body, Headers: responseHeaders}
	Logger.Debug(constant.TranceMsgResponse, append(commonLogFields, []zapcore.Field{zap.Any("Response", response)}...)...)
	return response, nil
}
