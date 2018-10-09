package cases

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/smoketest/models"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

// RegisterSecretPayload :
type RegisterSecretPayload struct {
	Secret       string `json:"secret"`
	TokenAddress string `json:"token_address"`
}

// RegisterSecretTest : test case for open channel
func RegisterSecretTest(env *models.RaidenEnvReader, allowFail bool) {
	// prepare data
	var payload RegisterSecretPayload
	payload.Secret = utils.NewRandomHash().String()
	payload.TokenAddress = env.RandomToken().Address
	payloadStr, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	// run case
	case1 := &APITestCase{
		CaseName:  "RegisterSecret",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "RegisterSecret",
			FullURL: env.RandomNode().Host + "/api/1/registersecret",
			Method:  http.MethodPost,
			Payload: string(payloadStr),
			Timeout: time.Second * 180,
		},
		TargetStatusCode: 500,
	}
	case1.Run()
}
