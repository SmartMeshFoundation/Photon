package v1

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/SmartMeshFoundation/Photon/rerr"

	"github.com/SmartMeshFoundation/Photon/dto"
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
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> GetSentTransfers ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	from, to := getFromTo(r)
	log.Trace(fmt.Sprintf("from=%d,to=%d\n", from, to))
	trs, err := API.GetSentTransfers(utils.EmptyAddress, from, to)
	resp = dto.NewAPIResponse(err, trs)
}

/*
GetReceivedTransfers retuns list of received transfer between `from_block` and `to_block`
it contains token swap
*/
func GetReceivedTransfers(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> GetReceivedTransfers ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	from, to := getFromTo(r)
	trs, err := API.GetReceivedTransfers(utils.EmptyAddress, from, to)
	resp = dto.NewAPIResponse(err, trs)
}

/*
Transfers is the api of /transfer/:token/:partner
*/
func Transfers(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> Transfers ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	var err error
	// 用户调用了prepare-update,暂停接收新交易
	// client invokes prepare-update, halts receiving new transfers.
	if API.Photon.StopCreateNewTransfers {
		resp = dto.NewExceptionAPIResponse(rerr.ErrStopCreateNewTransfer)
		return
	}
	token := r.PathParam("token")
	tokenAddr, err := utils.HexToAddress(token)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	target := r.PathParam("target")
	targetAddr, err := utils.HexToAddress(target)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	req := &TransferData{}
	err = r.DecodeJsonPayload(req)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	if req.Amount.Cmp(utils.BigInt0) <= 0 {
		resp = dto.NewExceptionAPIResponse(rerr.ErrInvalidAmount.Append("invalid amount"))
		return
	}
	if req.Fee == nil {
		req.Fee = utils.BigInt0
	}
	if req.Fee.Cmp(utils.BigInt0) < 0 {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.Append("invalid amount"))
		return
	}
	if len(req.Secret) != 0 && len(req.Secret) != 64 && (strings.HasPrefix(req.Secret, "0x") && len(req.Secret) != 66) {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.Append("invalid secret"))
		return
	}
	if len(req.Data) > params.MaxTransferDataLen {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.Append("Invalid data, length must < 256"))
		return
	}
	var result *utils.AsyncResult
	if req.Sync {
		result, err = API.Transfer(tokenAddr, req.Amount, req.Fee, targetAddr, common.HexToHash(req.Secret), params.MaxRequestTimeout, req.IsDirect, req.Data)
	} else {
		result, err = API.TransferAsync(tokenAddr, req.Amount, req.Fee, targetAddr, common.HexToHash(req.Secret), req.IsDirect, req.Data)
	}
	if err != nil {
		resp = dto.NewExceptionAPIResponse(err)
		return
	}
	if req.Fee.Cmp(utils.BigInt0) == 0 {
		req.Fee = nil
	}
	req.Initiator = API.Photon.NodeAddress.String()
	req.Target = target
	req.Token = token
	req.LockSecretHash = result.LockSecretHash.String()
	resp = dto.NewSuccessAPIResponse(req)
}

// GetTransferStatus : query transfer status by lockSecretHash
func GetTransferStatus(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> GetTransferStatus ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	lockSecretHashStr := r.PathParam("locksecrethash")
	lockSecretHash := common.HexToHash(lockSecretHashStr)
	token := r.PathParam("token")
	tokenAddr, err := utils.HexToAddress(token)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	ts, err := API.Photon.GetDao().GetSentTransferDetail(tokenAddr, lockSecretHash)
	resp = dto.NewAPIResponse(err, ts)
}

// CancelTransfer : cancel a transfer when haven't send secret
func CancelTransfer(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> CancelTransfer ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	var err error
	lockSecretHashStr := r.PathParam("locksecrethash")
	lockSecretHash := common.HexToHash(lockSecretHashStr)
	token := r.PathParam("token")
	tokenAddr, err := utils.HexToAddress(token)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	err = API.CancelTransfer(lockSecretHash, tokenAddr)
	resp = dto.NewAPIResponse(err, nil)
}
