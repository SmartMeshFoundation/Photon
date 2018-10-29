package cases

import (
	"net/http"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/smoketest/models"
)

// PrepareUpdateTest :
func PrepareUpdateTest(env *models.PhotonEnvReader, allowFail bool) {
	// run case
	case1 := &APITestCase{
		CaseName:  "PrepareUpdate",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "PrepareUpdate",
			FullURL: env.PhotonNodes[len(env.PhotonNodes)-1].Host + "/api/1/prepare-update",
			Method:  http.MethodPost,
			Payload: "",
			Timeout: queryTimeOut,
		},
		TargetStatusCode: 0,
	}
	case1.Run()
}
