package cases

import (
	"net/http"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/smoketest/models"
)

// LeavingTokenNetworkTest : test case for leave token network
func LeavingTokenNetworkTest(env *models.RaidenEnvReader, allowFail bool) {
	case1 := &APITestCase{
		CaseName:  "LeavingTokenNetwork",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "LeavingTokenNetwork",
			FullURL: env.RandomNode().Host + "/api/1/connections/" + env.RandomToken().Address,
			Method:  http.MethodDelete,
			Payload: "{\"only_receiving_channels\": false}",
			Timeout: time.Second * 180,
		},
		TargetStatusCode: 201,
	}
	case1.Run()
}
