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
	NewChannel     bool   `json:"new_channel"`
}

func newOpenChannelPayload(partnerAddress, tokenAddress string, balance, settleTimeout int32, newChannel bool) []byte {
	var n OpenChannelPayload
	n.PartnerAddress = partnerAddress
	n.TokenAddress = tokenAddress
	n.Balance = balance
	n.SettleTimeout = settleTimeout
	n.NewChannel = newChannel
	payload, err := json.Marshal(n)
	if err != nil {
		panic(err)
	}
	return payload
}

// OpenChannelTest : test case for open channel
func OpenChannelTest(env *models.PhotonEnvReader, allowFail bool) {
	// prepare data
	payload := newOpenChannelPayload("0x000000000000000000000000000000000FfffFfF",
		env.RandomToken().Address,
		50,
		350, true)
	// run case
	case1 := &APITestCase{
		CaseName:  "OpenChannel",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "OpenChannel",
			FullURL: env.RandomNode().Host + "/api/1/deposit",
			Method:  http.MethodPut,
			Payload: string(payload),
			Timeout: time.Second * 180,
		},
		TargetStatusCode: 200,
	}
	case1.Run()
}
