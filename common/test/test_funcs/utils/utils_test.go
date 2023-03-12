package utils

import (
	"testing"

	"github.com/chazool/go-sample-app/common/pkg/utils"
	"github.com/chazool/go-sample-app/common/pkg/utils/constant"
	test "github.com/chazool/go-sample-app/common/test/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type testHealth struct {
	Name     string `json:"name"`
	EndPoint string `json:"endPoint"`
	Expected string `json:"expected"`
}

var (
	commonLogFields = []zap.Field{zap.String(constant.TraceMsgReqId, "00000000-0000-0000-0000-000-00000000000")}
	netutils        = utils.NetUtilsImpl{NetUtils: &utils.NetUtilsService{}}.NetUtils
)

func initutilsAndCommonFieldData() {
	logConfig := zap.NewDevelopmentEncoderConfig()
	logConfig.FunctionKey = "F"
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   ".../.../data/utils/genarated/utils_test.log",
		MaxSize:    100,
		MaxBackups: 1,
		MaxAge:     1,
		Compress:   true,
	})
	utils.Logger = zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(logConfig), w, zapcore.DebugLevel))
}

func init() {
	initutilsAndCommonFieldData()
}

func TestLookIp(t *testing.T) {
	tetCases := test.GetTestCase("")

	for _, testCase := range tetCases {
		t.Run(testCase.Name, func(t *testing.T) {
			assert.NotNil(t, "look")
		})
	}

}
