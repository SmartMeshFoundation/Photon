package v1

import (
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

//TransferData post for transfers
type TransferData struct {
	Initiator      string   `json:"initiator_address"`
	Target         string   `json:"target_address"`
	Token          string   `json:"token_address"`
	Amount         *big.Int `json:"amount"`
	Secret         string   `json:"secret"` // 当用户想使用自己指定的密码,而非随机密码时使用	// client can assign specific secret
	LockSecretHash string   `json:"lockSecretHash"`
	Fee            *big.Int `json:"fee"`
	IsDirect       bool     `json:"is_direct"`
}

/*
GetSentTransfers returns list of sent transfer between `from_block` and `to_block`
*/
func GetSentTransfers(w rest.ResponseWriter, r *rest.Request) {
	from, to := getFromTo(r)
	log.Trace(fmt.Sprintf("from=%d,to=%d\n", from, to))
	trs, err := RaidenAPI.GetSentTransfers(from, to)
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
	trs, err := RaidenAPI.GetReceivedTransfers(from, to)
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
	// 用户调用了prepare-update,暂停接收新交易
	// client invokes prepare-update, halts receiving new transfers.
	if RaidenAPI.Raiden.StopCreateNewTransfers {
		rest.Error(w, "Stop create new transfers, please restart smartraiden", http.StatusBadRequest)
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
	lockSecretHash, err := RaidenAPI.TransferAndNotWait(tokenAddr, req.Amount, req.Fee, targetAddr, common.HexToHash(req.Secret), params.MaxRequestTimeout, req.IsDirect)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusConflict)
		return
	}
	req.Initiator = RaidenAPI.Raiden.NodeAddress.String()
	req.Target = target
	req.Token = token
	req.LockSecretHash = lockSecretHash.String()
	err = w.WriteJson(req)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

// GetTransferStatus : query transfer status by lockSecretHash
func GetTransferStatus(w rest.ResponseWriter, r *rest.Request) {
	lockSecretHashStr := r.PathParam("locksecrethash")
	lockSecretHash := common.HexToHash(lockSecretHashStr)

	ts, err := RaidenAPI.Raiden.GetDb().GetTransferStatus(lockSecretHash)
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
	lockSecretHashStr := r.PathParam("locksecrethash")
	lockSecretHash := common.HexToHash(lockSecretHashStr)
	token := r.PathParam("token")
	tokenAddr, err := utils.HexToAddress(token)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = RaidenAPI.CancelTransfer(lockSecretHash, tokenAddr)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusConflict)
		return
	}
}
