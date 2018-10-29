package cases

import (
	"errors"
	"log"
	"net/http"
	"time"

	"encoding/json"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/smoketest/models"
)

// TransferPayload API  http body
type TransferPayload struct {
	Amount   int32  `json:"amount"`
	Fee      int64  `json:"fee"`
	IsDirect bool   `json:"is_direct"`
	Secret   string `json:"secret"`
}

type testTransferParams struct {
	Env          *models.PhotonEnvReader
	AllowFail    bool
	CaseName     string
	PrepareData  func(env *models.PhotonEnvReader) (node1 *models.PhotonNode, node2 *models.PhotonNode, token *models.Token, err error)
	IsDirect     bool
	TargetStatus int
}

// InitiatingTransferTest : test case for InitiatingTransfer
func InitiatingTransferTest(env *models.PhotonEnvReader, allowFail bool) {

	// test transfer between two nodes who have direct opened channel
	testTransfer(&testTransferParams{
		Env:          env,
		AllowFail:    allowFail,
		CaseName:     "DirectTransfer A-B isDirect=true",
		PrepareData:  prepareDataForDirectTransfer,
		IsDirect:     true,
		TargetStatus: 200,
	})
	testTransfer(&testTransferParams{
		Env:          env,
		AllowFail:    allowFail,
		CaseName:     "DirectTransfer A-B isDirect=false",
		PrepareData:  prepareDataForDirectTransfer,
		IsDirect:     false,
		TargetStatus: 200,
	})
	// test transfer between two nodes who doesn't have direct opened channel
	testTransfer(&testTransferParams{
		Env:          env,
		AllowFail:    allowFail,
		CaseName:     "IndirectTransfer A-B-C isDirect=true",
		PrepareData:  prepareDataForIndirectTransfer,
		IsDirect:     true,
		TargetStatus: 409,
	})
	// test transfer between two nodes who doesn't have direct opened channel
	testTransfer(&testTransferParams{
		Env:          env,
		AllowFail:    allowFail,
		CaseName:     "IndirectTransfer A-B-C isDirect=false",
		PrepareData:  prepareDataForIndirectTransfer,
		IsDirect:     false,
		TargetStatus: 200,
	})
}

func testTransfer(param *testTransferParams) {
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
	p, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	// run case
	case1 := &APITestCase{
		CaseName:  param.CaseName,
		AllowFail: param.AllowFail,
		Req: &models.Req{
			APIName: "InitiatingTransfer",
			FullURL: sender.Host + "/api/1/transfers/" + token.Address + "/" + receiver.AccountAddress,
			Method:  http.MethodPost,
			Payload: string(p),
			Timeout: time.Second * 60,
		},
		TargetStatusCode: param.TargetStatus,
	}
	case1.Run()
}

// find a opened channel from env, if there is none, create one
func prepareDataForDirectTransfer(env *models.PhotonEnvReader) (sender *models.PhotonNode, receiver *models.PhotonNode, token *models.Token, err error) {
	if len(env.PhotonNodes) < 2 {
		err = errors.New("no enough photon node")
		return
	}
	sender, receiver = env.PhotonNodes[0], env.PhotonNodes[1]
	for _, t := range env.Tokens {
		if env.HasOpenedChannelBetween(sender, receiver, t) {
			token = t
			break
		}
	}
	if token == nil {
		err = fmt.Errorf("no opened channel between %s and %s", sender.AccountAddress, receiver.AccountAddress)
		return
	}
	return
}

// find a enable route from env, if there is none, create one
func prepareDataForIndirectTransfer(env *models.PhotonEnvReader) (sender *models.PhotonNode, receiver *models.PhotonNode, token *models.Token, err error) {
	if len(env.PhotonNodes) < 3 {
		err = errors.New("no enough photon node")
		return
	}
	sender, mid, receiver := env.PhotonNodes[0], env.PhotonNodes[1], env.PhotonNodes[2]
	for _, t := range env.Tokens {
		if env.HasOpenedChannelBetween(sender, mid, t) && env.HasOpenedChannelBetween(mid, receiver, t) {
			token = t
			break
		}
	}
	if token == nil {
		err = fmt.Errorf("no enable route between %s and %s", sender.AccountAddress, receiver.AccountAddress)
		return
	}
	return
}
