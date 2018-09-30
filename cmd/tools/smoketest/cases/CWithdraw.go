package cases

import (
	"log"
	"net/http"
	"time"

	"math/big"

	"encoding/json"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/smoketest/models"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
)

// WithdrawTest :
func WithdrawTest(env *models.RaidenEnvReader, allowFail bool) {
	// prepare data for this case
	node := env.RaidenNodes[len(env.RaidenNodes)-2]
	channels := env.GetChannelsOfNodeByState(node.AccountAddress, contracts.ChannelStateOpened)
	if channels == nil || len(channels) == 0 {
		Logger.Println("Current env can not afford this case !!!")
		if !allowFail {
			Logger.Println("allowFail = false,exit")
			panic("allowFail = false,exit")
		}
		log.Println("Case [CancelTransferTest] FAILED because no suitable env !!!")
		Logger.Println("Case [CancelTransferTest] FAILED because no suitable env !!!")
	}
	// 1. n1 withdraw
	type Req struct {
		Amount *big.Int
		Op     string
	}
	const OpPrepareWithdraw = "preparewithdraw"
	const OpCancelPrepare = "cancelprepare"

	var payload Req
	payload.Amount = big.NewInt(1)
	p, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	case1 := &APITestCase{
		CaseName:  "Withdraw",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "Withdraw",
			FullURL: node.Host + "/api/1/withdraw/" + channels[0].ChannelIdentifier,
			Method:  http.MethodPut,
			Payload: string(p),
			Timeout: time.Second * 180,
		},
		TargetStatusCode: 200,
	}
	case1.Run()
}
