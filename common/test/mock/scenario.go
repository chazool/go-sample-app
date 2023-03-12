package mock

import (
	"encoding/json"
	"os"

	"github.com/chazool/go-sample-app/common/pkg/utils"
	"github.com/chazool/go-sample-app/common/pkg/utils/constant"
	"go.uber.org/zap"
)

type Service struct {
	ServiceName      string              `json:"serviceName"`
	FilePath         string              `json:"filePath"`
	ScenarioFileName string              `json:"scenarioFileName"`
	IsGetService     bool                `json:"isGetService"`
	ServiceUrl       string              `json:"serviceUrl"`
	HttpMethodType   constant.HttpMethod `json:"httpMethodType"`
	WorkFlowService  string              `json:"workFlowService"`
	IsDeleteService  bool                `json:"isDeleteService"`
	ISNetUtilsUsed   bool                `json:"iSNetUtilsUsed"`
}

type Case struct {
	Name                          string     `json:"name"`
	Description                   string     `json:"description"`
	IsMockRequired                bool       `json:"isMockRequired"`
	RequestURL                    string     `json:"requestURL"`
	RequestFileName               string     `json:"requestFileName"`
	ResponseFileName              string     `json:"responseFileName"`
	ExpectedStatusCode            string     `json:"expectedStatusCode"`
	IsIgnoreCompare               string     `json:"isIgnoreCompare"`
	FuncMockSteps                 []FuncStep `json:"funcMockSteps"`
	MarkDowExampleName            string     `json:"markDowExampleName"`
	MarkDownExampleDescription    string     `json:"markDownExampleDescription"`
	IsOmitFromServiceMarkdownDocs bool       `json:"isOmitFromServiceMarkdownDocs"`
}

type FuncStep struct {
	MockFunction string     `json:"method"`
	ReturnData   ReturnData `json:"returnData"`
	Params       Params     `json:"params"`
	IsMock       bool       `json:"isMock"`
}

type Params struct {
	CommonLogFields []zap.Field            `json:"commonLogFields"`
	Args            map[string]interface{} `json:"args"`
}

type ReturnData struct {
	Error   string                 `json:"error"`
	Outputs map[string]interface{} `json:"outputs"`
}

func GetTestCase(scenarioPath string) []Case {
	b, err := os.ReadFile(scenarioPath)
	if err != nil {
		utils.Logger.Error("Invalid Data or File", zap.Error(err))
	}

	var testCase []Case
	err = json.Unmarshal(b, &testCase)
	if err != nil {
		utils.Logger.Error("Invalid Data or File", zap.Error(err))
	}
	return testCase
}
