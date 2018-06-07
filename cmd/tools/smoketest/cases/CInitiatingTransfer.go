package cases

import (
	"log"
	"net/http"
	"time"

	"encoding/json"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/smoketest/models"
	"github.com/go-errors/errors"
)

// TransferPayload API  http body
type TransferPayload struct {
	Amount   int32 `json:"amount"`
	Fee      int64 `json:"fee"`
	IsDirect bool  `json:"is_direct"`
}

type testTransferParams struct {
	Env         *models.RaidenEnvReader
	AllowFail   bool
	CaseName    string
	PrepareData func(env *models.RaidenEnvReader) (node1 *models.RaidenNode, node2 *models.RaidenNode, token *models.Token, err error)
	IsDirect    bool
}

// InitiatingTransferTest : test case for InitiatingTransfer
func InitiatingTransferTest(env *models.RaidenEnvReader, allowFail bool) {

	// test transfer between two nodes who have direct opened channel
	testTransfer(&testTransferParams{
		Env:         env,
		AllowFail:   allowFail,
		CaseName:    "DirectTransfer A-B isDirect=true",
		PrepareData: prepareDataForDirectTransfer,
		IsDirect:    true,
	}, 200)
	testTransfer(&testTransferParams{
		Env:         env,
		AllowFail:   allowFail,
		CaseName:    "DirectTransfer A-B isDirect=false",
		PrepareData: prepareDataForDirectTransfer,
		IsDirect:    false,
	}, 200)
	// test transfer between two nodes who doesn't have direct opened channel
	testTransfer(&testTransferParams{
		Env:         env,
		AllowFail:   allowFail,
		CaseName:    "IndirectTransfer A-B-C isDirect=true",
		PrepareData: prepareDataForIndirectTransfer,
		IsDirect:    true,
	}, 500)
	// test transfer between two nodes who doesn't have direct opened channel
	testTransfer(&testTransferParams{
		Env:         env,
		AllowFail:   allowFail,
		CaseName:    "IndirectTransfer A-B-C isDirect=false",
		PrepareData: prepareDataForIndirectTransfer,
		IsDirect:    false,
	}, 200)
}

func testTransfer(param *testTransferParams, targetStatus int) {
	// prepare data
	sender, receiver, token, err := param.PrepareData(param.Env)
	if err != nil {
		log.Printf("Case [%-40s] FAILED because no suitable env : %s", param.CaseName, err.Error())
		Logger.Printf("Case [%-40s] FAILED because no suitable env : %s", param.CaseName, err.Error())
		if !param.AllowFail {
			Logger.Println("allowFail = false,exit")
			panic("allowFail = false,exit")
		}
		return
	}
	var payload TransferPayload
	payload.Amount = 5
	payload.Fee = 0
	payload.IsDirect = param.IsDirect
	p, _ := json.Marshal(payload)
	// run case
	case1 := &APITestCase{
		CaseName:  param.CaseName,
		AllowFail: param.AllowFail,
		Req: &models.Req{
			APIName: "InitiatingTransfer",
			FullURL: sender.Host + "/api/1/transfers/" + token.Address + "/" + receiver.AccountAddress,
			Method:  http.MethodPost,
			Payload: string(p),
			Timeout: time.Second * 180,
		},
		TargetStatusCode: targetStatus,
	}
	case1.Run()
}

// find a opened channel from env, if there is none, create one
func prepareDataForDirectTransfer(env *models.RaidenEnvReader) (sender *models.RaidenNode, receiver *models.RaidenNode, token *models.Token, err error) {
	if len(env.RaidenNodes) < 2 {
		err = errors.New("no enough raiden node")
		return
	}
	sender, receiver = env.RaidenNodes[0], env.RaidenNodes[1]
	for _, t := range env.Tokens {
		if env.HasOpenedChannelBetween(sender, receiver, t) {
			token = t
			break
		}
	}
	if token == nil {
		err = errors.New(fmt.Errorf("no opened channel between %s and %s", sender.AccountAddress, receiver.AccountAddress))
		return
	}
	return
}

// find a enable route from env, if there is none, create one
func prepareDataForIndirectTransfer(env *models.RaidenEnvReader) (sender *models.RaidenNode, receiver *models.RaidenNode, token *models.Token, err error) {
	if len(env.RaidenNodes) < 3 {
		err = errors.New("no enough raiden node")
		return
	}
	sender, mid, receiver := env.RaidenNodes[0], env.RaidenNodes[1], env.RaidenNodes[2]
	for _, t := range env.Tokens {
		if env.HasOpenedChannelBetween(sender, mid, t) && env.HasOpenedChannelBetween(mid, receiver, t) {
			token = t
			break
		}
	}
	if token == nil {
		err = errors.New(fmt.Errorf("no enable route between %s and %s", sender.AccountAddress, receiver.AccountAddress))
		return
	}
	return
}
