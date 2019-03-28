package cases

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/smoketest/models"
	"github.com/SmartMeshFoundation/Photon/utils"
)

// CancelTransferTest :
func CancelTransferTest(env *models.PhotonEnvReader, allowFail bool) {
	// prepare data for this case
	sender, receiver, token, err := prepareDataForDirectTransfer(env)
	if err != nil {
		Logger.Println("Current env can not afford this case !!!")
		if !allowFail {
			Logger.Println("allowFail = false,exit")
			panic("allowFail = false,exit")
		}
		log.Println("Case [CancelTransferTest] FAILED because no suitable env !!!")
		Logger.Println("Case [CancelTransferTest] FAILED because no suitable env !!!")
		return
	}

	// 1. n1 transfer to n2 with secret
	secret := utils.NewRandomHash()
	lockSecretHash := utils.ShaSecret(secret[:])
	var payload TransferPayload
	payload.Amount = 1
	payload.Fee = 0
	payload.IsDirect = false
	payload.Secret = secret.String()
	p, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	case1 := &APITestCase{
		CaseName:  "Transfer",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "Transfer",
			FullURL: sender.Host + "/api/1/transfers/" + token.Address + "/" + receiver.AccountAddress,
			Method:  http.MethodPost,
			Payload: string(p),
			Timeout: time.Second * 180,
		},
		TargetStatusCode: 200,
	}
	case1.Run()

	// 2. get transfer status
	case2 := &APITestCase{
		CaseName:  "GetSentTransferDetail",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "GetSentTransferDetail",
			FullURL: sender.Host + "/api/1/transferstatus/" + token.Address + "/" + lockSecretHash.String(),
			Method:  http.MethodGet,
			Timeout: time.Second * 180,
		},
		TargetStatusCode: 0,
	}
	case2.Run()

	// 3. receiver get unfinished transfers
	case3 := &APITestCase{
		CaseName:  "GetUnfinishedReceivedTransfer",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "GetUnfinishedReceivedTransfer",
			FullURL: receiver.Host + "/api/1/getunfinishedreceivedtransfer/" + token.Address + "/" + lockSecretHash.String(),
			Method:  http.MethodGet,
			Timeout: time.Second * 180,
		},
		TargetStatusCode: 200,
	}
	case3.Run()

	// 4. stop receiver
	case4 := &APITestCase{
		CaseName:  "Stop",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "Stop",
			FullURL: receiver.Host + "/api/1/stop",
			Method:  http.MethodGet,
			Payload: string(p),
			Timeout: time.Second * 180,
		},
		TargetStatusCode: 200,
	}
	case4.Run()
	// 5. sender allow reveal secret
	type AllowRevealSecretPayload struct {
		LockSecretHash string `json:"lock_secret_hash"`
		TokenAddress   string `json:"token_address"`
	}
	var p5 AllowRevealSecretPayload
	p5.TokenAddress = token.Address
	p5.LockSecretHash = lockSecretHash.String()
	p, err = json.Marshal(p5)
	if err != nil {
		panic(err)
	}
	case5 := &APITestCase{
		CaseName:  "CancelTransfer",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "CancelTransfer",
			FullURL: sender.Host + "/api/1/transfers/allowrevealsecret",
			Method:  http.MethodPost,
			Payload: string(p),
			Timeout: time.Second * 180,
		},
		TargetStatusCode: 200,
	}
	case5.Run()
	// 6. sender cancel transfer
	case6 := &APITestCase{
		CaseName:  "CancelTransfer",
		AllowFail: allowFail,
		Req: &models.Req{
			APIName: "CancelTransfer",
			FullURL: sender.Host + "/api/1/transfercancel/" + token.Address + "/" + lockSecretHash.String(),
			Method:  http.MethodPost,
			Payload: string(p),
			Timeout: time.Second * 180,
		},
		TargetStatusCode: 200,
	}
	case6.Run()
}
