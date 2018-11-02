package v1

import (
	"fmt"
	"net/http"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

// AllowRevealSecret :
// 当用户发出了一笔指定密码的交易时使用,这种交易在调用本接口解锁之前,是不会接收方发来的SecretRequest的
// AllowRevealSecret : used when clients send a transfer with specific secrets.
// that secret will not receive SecretRequest before invoking this function to unlock.
func AllowRevealSecret(w rest.ResponseWriter, r *rest.Request) {
	var err error
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> RegisterSecret ,err=%v", err))
	}()
	type AllowRevealSecretPayload struct {
		LockSecretHash string `json:"lock_secret_hash"`
		TokenAddress   string `json:"token_address"`
	}
	var payload AllowRevealSecretPayload
	err = r.DecodeJsonPayload(&payload)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	lockSecretHash := common.HexToHash(payload.LockSecretHash)
	tokenAddress, err := utils.HexToAddress(payload.TokenAddress)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = API.AllowRevealSecret(lockSecretHash, tokenAddress)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetUnfinishedReceivedTransfer :根据lockSecretHash查询未完成的交易
// GetUnfinishedReceivedTransfer : check incomplete transfers according to lockSecretHash
func GetUnfinishedReceivedTransfer(w rest.ResponseWriter, r *rest.Request) {
	tokenAddressStr := r.PathParam("tokenaddress")
	tokenAddress, err := utils.HexToAddress(tokenAddressStr)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	lockSecretHashStr := r.PathParam("locksecrethash")
	lockSecretHash := common.HexToHash(lockSecretHashStr)
	if lockSecretHash == utils.EmptyHash {
		rest.Error(w, "Invalid lockSecretHash", http.StatusBadRequest)
		return
	}
	transferData := API.GetUnfinishedReceivedTransfer(lockSecretHash, tokenAddress)
	err = w.WriteJson(transferData)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

// RegisterSecret :
// 当从其他渠道知道一笔交易的密码时,可以注册密码到state manger中
// RegisterSecret : when knowing the secret of a transfer from other sources,
// we can register secret in statemanager
func RegisterSecret(w rest.ResponseWriter, r *rest.Request) {
	var err error
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> RegisterSecret ,err=%v", err))
	}()
	type RegisterSecretPayload struct {
		Secret       string `json:"secret"`
		TokenAddress string `json:"token_address"`
	}
	var payload RegisterSecretPayload
	err = r.DecodeJsonPayload(&payload)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	secret := common.HexToHash(payload.Secret)
	tokenAddress, err := utils.HexToAddress(payload.TokenAddress)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = API.RegisterSecret(secret, tokenAddress)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
