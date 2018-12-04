package v1

import (
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

//TransferData post for transfers
type TransferData struct {
	Initiator      string   `json:"initiator_address"`
	Target         string   `json:"target_address"`
	Token          string   `json:"token_address"`
	Amount         *big.Int `json:"amount"`
	Secret         string   `json:"secret,omitempty"` // 当用户想使用自己指定的密码,而非随机密码时使用	// client can assign specific secret
	LockSecretHash string   `json:"lockSecretHash"`
	Fee            *big.Int `json:"fee,omitempty"`
	IsDirect       bool     `json:"is_direct,omitempty"`
	Sync           bool     `json:"sync,omitempty"` //是否同步
	Data           string   `json:"data"`           // 交易附加信息,长度不超过256
}

/*
GetSentTransfers returns list of sent transfer between `from_block` and `to_block`
*/
func GetSentTransfers(w rest.ResponseWriter, r *rest.Request) {
	from, to := getFromTo(r)
	log.Trace(fmt.Sprintf("from=%d,to=%d\n", from, to))
	trs, err := API.GetSentTransfers(from, to)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = w.WriteJson(trs)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

/*
GetReceivedTransfers retuns list of received transfer between `from_block` and `to_block`
it contains token swap
*/
func GetReceivedTransfers(w rest.ResponseWriter, r *rest.Request) {
	from, to := getFromTo(r)
	trs, err := API.GetReceivedTransfers(from, to)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = w.WriteJson(trs)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

/*
Transfers is the api of /transfer/:token/:partner
*/
func Transfers(w rest.ResponseWriter, r *rest.Request) {
	var err error
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> Transfers ,err=%v", err))
	}()
	// 用户调用了prepare-update,暂停接收新交易
	// client invokes prepare-update, halts receiving new transfers.
	if API.Photon.StopCreateNewTransfers {
		rest.Error(w, "Stop create new transfers, please restart photon", http.StatusBadRequest)
		return
	}
	token := r.PathParam("token")
	tokenAddr, err := utils.HexToAddress(token)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	target := r.PathParam("target")
	targetAddr, err := utils.HexToAddress(target)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req := &TransferData{}
	err = r.DecodeJsonPayload(req)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Amount.Cmp(utils.BigInt0) <= 0 {
		rest.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}
	if req.Fee == nil {
		req.Fee = utils.BigInt0
	}
	if req.Fee.Cmp(utils.BigInt0) < 0 {
		rest.Error(w, "Invalid fee", http.StatusBadRequest)
		return
	}
	if len(req.Secret) != 0 && len(req.Secret) != 64 && (strings.HasPrefix(req.Secret, "0x") && len(req.Secret) != 66) {
		rest.Error(w, "Invalid secret", http.StatusBadRequest)
		return
	}
	if len(req.Data) > params.MaxTransferDataLen {
		rest.Error(w, "Invalid data, length must < 256", http.StatusBadRequest)
		return
	}
	var result *utils.AsyncResult
	if req.Sync {
		result, err = API.Transfer(tokenAddr, req.Amount, req.Fee, targetAddr, common.HexToHash(req.Secret), params.MaxRequestTimeout, req.IsDirect, req.Data)
	} else {
		result, err = API.TransferAsync(tokenAddr, req.Amount, req.Fee, targetAddr, common.HexToHash(req.Secret), req.IsDirect, req.Data)
	}
	if err != nil {
		rest.Error(w, err.Error(), http.StatusConflict)
		return
	}
	//todo fix me support charging fee
	if req.Fee.Cmp(utils.BigInt0) == 0 {
		req.Fee = nil
	}
	req.Initiator = API.Photon.NodeAddress.String()
	req.Target = target
	req.Token = token
	req.LockSecretHash = result.LockSecretHash.String()
	err = w.WriteJson(req)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

// GetTransferStatus : query transfer status by lockSecretHash
func GetTransferStatus(w rest.ResponseWriter, r *rest.Request) {
	lockSecretHashStr := r.PathParam("locksecrethash")
	lockSecretHash := common.HexToHash(lockSecretHashStr)
	token := r.PathParam("token")
	tokenAddr, err := utils.HexToAddress(token)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ts, err := API.Photon.GetDao().GetTransferStatus(tokenAddr, lockSecretHash)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusConflict)
		return
	}
	err = w.WriteJson(ts)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

// CancelTransfer : cancel a transfer when haven't send secret
func CancelTransfer(w rest.ResponseWriter, r *rest.Request) {
	var err error
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> CancelTransfer ,err=%v", err))
	}()
	lockSecretHashStr := r.PathParam("locksecrethash")
	lockSecretHash := common.HexToHash(lockSecretHashStr)
	token := r.PathParam("token")
	tokenAddr, err := utils.HexToAddress(token)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = API.CancelTransfer(lockSecretHash, tokenAddr)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusConflict)
		return
	}
}
