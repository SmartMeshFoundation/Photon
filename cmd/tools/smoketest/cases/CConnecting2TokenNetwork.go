package cases

import (
	"net/http"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/smoketest/models"
)

// Connecting2TokenNetworkTest : test case for connect token network
func Connecting2TokenNetworkTest(env *models.RaidenEnvReader, allowFail bool) {
	case1 := &APITestCase{
		CaseName:  "Connecting2TokenNetwork",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "Connecting2TokenNetwork",
			FullURL: env.RandomNode().Host + "/api/1/connections/" + env.RandomToken().Address,
			Method:  http.MethodPut,
			Payload: "{\"funds\": 10}",
			Timeout: time.Second * 180,
		},
		TargetStatusCode: 201,
	}
	case1.Run()
}
