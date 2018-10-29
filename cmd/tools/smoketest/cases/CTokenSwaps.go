package cases

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/smoketest/models"
	"github.com/SmartMeshFoundation/Photon/utils"
)

// TokenSwapsPayload :
type TokenSwapsPayload struct {
	Role            string `json:"role"`
	SendingAmount   int32  `json:"sending_amount"`
	SendingToken    string `json:"sending_token"`
	ReceivingAmount int32  `json:"receiving_amount"`
	ReceivingToken  string `json:"receiving_token"`
	Secret          string `json:"secret"`
}

type testTokenSwapParams struct {
	Env         *models.PhotonEnvReader
	AllowFail   bool
	CaseName    string
	PrepareData func(env *models.PhotonEnvReader) (node1 *models.PhotonNode, node2 *models.PhotonNode, token1 *models.Token, token2 *models.Token, err error)
}

// TokenSwapsTest : test case for TokenSwap
func TokenSwapsTest(env *models.PhotonEnvReader, allowFail bool) {
	// test TokenSwap between two nodes who have direct opened channel
	testTokenSwap(&testTokenSwapParams{
		Env:         env,
		AllowFail:   allowFail,
		CaseName:    "DirectTokenSwap A-B",
		PrepareData: prepareDataForDirectTokenSwap,
	})
	// test TokenSwap between two nodes who doesn't have direct opened channel
	testTokenSwap(&testTokenSwapParams{
		Env:         env,
		AllowFail:   allowFail,
		CaseName:    "IndirectTokenSwap A-B-C",
		PrepareData: prepareDataForIndirectTokenSwap,
	})
}

func testTokenSwap(param *testTokenSwapParams) {

	// prepare data
	node1, node2, token1, token2, err := param.PrepareData(param.Env)
	if err != nil {
		log.Printf("Case [%-40s] FAILED because no suitable env : %s", param.CaseName, err.Error())
		Logger.Printf("Case [%-40s] FAILED because no suitable env : %s", param.CaseName, err.Error())
		if !param.AllowFail {
			Logger.Println("allowFail = false,exit")
			panic("allowFail = false,exit")
		}
		return
	}

	// run case
	secret, lockSecretHash := getRandomSecret()
	invokeTokenSwap(node1, node2, token1, token2, 1, 2, "taker", param.CaseName, param.AllowFail, lockSecretHash, "")
	invokeTokenSwap(node2, node1, token2, token1, 2, 1, "maker", param.CaseName, param.AllowFail, lockSecretHash, secret)
}

func prepareDataForDirectTokenSwap(env *models.PhotonEnvReader) (sender *models.PhotonNode, receiver *models.PhotonNode, token1 *models.Token, token2 *models.Token, err error) {
	if len(env.PhotonNodes) < 2 {
		err = errors.New("no enough photon node")
		return
	}
	if len(env.Tokens) < 2 {
		err = errors.New("no enough registered token ")
		return
	}
	sender, receiver = env.PhotonNodes[0], env.PhotonNodes[1]
	token1, token2 = env.Tokens[0], env.Tokens[1]
	if !env.HasOpenedChannelBetween(sender, receiver, token1) {
		err = fmt.Errorf("no opened channel on token [%s] between %s and %s", token1.Address, sender.AccountAddress, receiver.AccountAddress)
		return
	}
	if !env.HasOpenedChannelBetween(sender, receiver, token2) {
		err = fmt.Errorf("no opened channel on token [%s] between %s and %s", token2.Address, sender.AccountAddress, receiver.AccountAddress)
		return
	}
	return
}

func prepareDataForIndirectTokenSwap(env *models.PhotonEnvReader) (sender *models.PhotonNode, receiver *models.PhotonNode, token1 *models.Token, token2 *models.Token, err error) {
	if len(env.PhotonNodes) < 3 {
		err = errors.New("no enough photon node")
		return
	}
	if len(env.Tokens) < 2 {
		err = errors.New("no enough registered token ")
		return
	}
	sender, mid, receiver := env.PhotonNodes[0], env.PhotonNodes[1], env.PhotonNodes[2]
	token1, token2 = env.Tokens[0], env.Tokens[1]
	if !env.HasOpenedChannelBetween(sender, mid, token1) {
		err = fmt.Errorf("no opened channel on token [%s] between %s and %s", token1.Address, sender.AccountAddress, mid.AccountAddress)
		return
	}
	if !env.HasOpenedChannelBetween(mid, receiver, token1) {
		err = fmt.Errorf("no opened channel on token [%s] between %s and %s", token1.Address, mid.AccountAddress, receiver.AccountAddress)
		return
	}
	if !env.HasOpenedChannelBetween(sender, mid, token2) {
		err = fmt.Errorf("no opened channel on token [%s] between %s and %s", token2.Address, sender.AccountAddress, mid.AccountAddress)
		return
	}
	if !env.HasOpenedChannelBetween(mid, receiver, token2) {
		err = fmt.Errorf("no opened channel on token [%s] between %s and %s", token2.Address, mid.AccountAddress, receiver.AccountAddress)
		return
	}
	return
}

func invokeTokenSwap(node1 *models.PhotonNode, node2 *models.PhotonNode, token1 *models.Token, token2 *models.Token, amount1 int32, amount2 int32, role string, caseName string, allowFail bool, lockSecretHash string, secret string) {
	payload := TokenSwapsPayload{
		Role:            role,
		SendingToken:    token1.Address,
		SendingAmount:   amount1,
		ReceivingToken:  token2.Address,
		ReceivingAmount: amount2,
		Secret:          secret,
	}
	p, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	// run case
	case1 := &APITestCase{
		CaseName:  caseName + " " + role,
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "TokenSwap",
			FullURL: node1.Host + "/api/1/token_swaps/" + node2.AccountAddress + "/" + lockSecretHash,
			Method:  http.MethodPut,
			Payload: string(p),
			Timeout: time.Second * 60,
		},
		TargetStatusCode: 201,
	}
	case1.Run()
}

func getRandomSecret() (string, string) {
	t := utils.RandomString(5)
	secret := utils.ShaSecret([]byte(t))
	lockSecretHash := utils.ShaSecret(secret.Bytes())
	return secret.String(), lockSecretHash.String()
}
