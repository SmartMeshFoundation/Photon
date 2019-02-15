package v1

import (
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/Photon/rerr"

	"github.com/SmartMeshFoundation/Photon/dto"

	"context"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/network/netshare"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

/*
Balance for test only
query `addr`'s balance on `token`
*/
func Balance(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> Balance ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	tokenstr := r.PathParam("token")
	addrstr := r.PathParam("addr")
	token, err := utils.HexToAddress(tokenstr)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	addr, err := utils.HexToAddress(addrstr)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	t, err := API.Photon.Chain.Token(token)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	v, err := t.BalanceOf(addr)
	resp = dto.NewAPIResponse(err, v)
}

/*
TransferToken for test only
Transfer from this node to `addr` `value` tokens on token `token`
*/
func TransferToken(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> TransferToken ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	tokenstr := r.PathParam("token")
	addrstr := r.PathParam("addr")
	valuestr := r.PathParam("value")
	token, err := utils.HexToAddress(tokenstr)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	addr, err := utils.HexToAddress(addrstr)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	v, b := new(big.Int).SetString(valuestr, 0)
	if !b {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError)
		return
	}
	t, err := API.Photon.Chain.Token(token)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	err = t.Transfer(addr, v)
	resp = dto.NewAPIResponse(err, nil)
}

//EthBalance how many eth `addr` have.
func EthBalance(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> EthBalance ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	addrstr := r.PathParam("addr")
	addr, err := utils.HexToAddress(addrstr)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	v, err := API.Photon.Chain.Client.BalanceAt(context.Background(), addr, nil)
	resp = dto.NewAPIResponse(err, v)
}

//BlockTimeFormat  is time format of last block
const BlockTimeFormat = "01-02|15:04:05.999"

//ConnectionStatus status of network connection
type ConnectionStatus struct {
	XMPPStatus    netshare.Status `json:"xmpp_status"`
	EthStatus     netshare.Status `json:"eth_status"`
	LastBlockTime string          `json:"last_block_time"`
}

/*
EthereumStatus  query the status between Photon and ethereum
*/
func EthereumStatus(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> EthereumStatus ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	c := API.Photon.Chain
	cs := &ConnectionStatus{
		XMPPStatus:    netshare.Disconnected,
		LastBlockTime: API.Photon.GetDao().GetLastBlockNumberTime().Format(BlockTimeFormat),
	}
	if c != nil && c.Client.Status == netshare.Connected {
		cs.EthStatus = netshare.Connected
	} else {
		cs.EthStatus = netshare.Disconnected
	}
	resp = dto.NewAPIResponse(nil, cs)
}

/*
ForceUnlock force unlock by locksecrethash
*/
func ForceUnlock(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> ForceUnlock ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	channelIdentifierStr := r.PathParam("channel")
	channelIdentifier := common.HexToHash(channelIdentifierStr)
	secretStr := r.PathParam("secret")
	secret := common.HexToHash(secretStr)
	err := API.ForceUnlock(channelIdentifier, secret)
	resp = dto.NewAPIResponse(err, "ok")
}

/*
RegisterSecretOnChain register secret to contract
*/
func RegisterSecretOnChain(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> RegisterSecretOnChain ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	secretStr := r.PathParam("secret")
	secret := common.HexToHash(secretStr)
	err := API.RegisterSecretOnChain(secret)
	resp = dto.NewAPIResponse(err, "ok")
}
