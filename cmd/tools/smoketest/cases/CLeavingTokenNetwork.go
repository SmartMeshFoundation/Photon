package cases

import (
	"net/http"
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/smoketest/models"
)

// LeavingTokenNetworkTest : test case for leave token network
func LeavingTokenNetworkTest(env *models.PhotonEnvReader, allowFail bool) {
	case1 := &APITestCase{
		CaseName:  "LeavingTokenNetwork",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "LeavingTokenNetwork",
			FullURL: env.RandomNode().Host + "/api/1/connections/" + env.RandomToken().Address,
			Method:  http.MethodDelete,
			Payload: "{\"only_receiving_channels\": false}",
			Timeout: time.Second * 360,
		},
		TargetStatusCode: 200,
	}
	case1.Run()
}
