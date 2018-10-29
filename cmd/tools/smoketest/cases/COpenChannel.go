package cases

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/smoketest/models"
)

// OpenChannelPayload :
type OpenChannelPayload struct {
	PartnerAddress string `json:"partner_address"`
	TokenAddress   string `json:"token_address"`
	Balance        int32  `json:"balance"`
	SettleTimeout  int32  `json:"settle_timeout"`
}

// OpenChannelTest : test case for open channel
func OpenChannelTest(env *models.PhotonEnvReader, allowFail bool) {
	// prepare data
	var newchannel OpenChannelPayload
	newchannel.PartnerAddress = "0x000000000000000000000000000000000FfffFfF"
	newchannel.TokenAddress = env.RandomToken().Address
	//newchannel.Balance = 50
	newchannel.SettleTimeout = 35
	payload, err := json.Marshal(newchannel)
	if err != nil {
		panic(err)
	}
	// run case
	case1 := &APITestCase{
		CaseName:  "OpenChannel",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "OpenChannel",
			FullURL: env.RandomNode().Host + "/api/1/channels",
			Method:  http.MethodPut,
			Payload: string(payload),
			Timeout: time.Second * 180,
		},
		TargetStatusCode: 200,
	}
	case1.Run()
}
