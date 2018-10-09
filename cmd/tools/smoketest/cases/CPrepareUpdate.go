package cases

import (
	"net/http"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/smoketest/models"
)

// PrepareUpdateTest :
func PrepareUpdateTest(env *models.RaidenEnvReader, allowFail bool) {
	// run case
	case1 := &APITestCase{
		CaseName:  "PrepareUpdate",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "PrepareUpdate",
			FullURL: env.RaidenNodes[len(env.RaidenNodes)-1].Host + "/api/1/prepare-update",
			Method:  http.MethodPost,
			Payload: "",
			Timeout: queryTimeOut,
		},
		TargetStatusCode: 0,
	}
	case1.Run()
}
