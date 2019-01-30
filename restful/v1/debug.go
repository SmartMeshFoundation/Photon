package v1

import (
	"fmt"
	"math/big"
	"net/http"

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
	tokenstr := r.PathParam("token")
	addrstr := r.PathParam("addr")
	token, err := utils.HexToAddress(tokenstr)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	addr, err := utils.HexToAddress(addrstr)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t, err := API.Photon.Chain.Token(token)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusConflict)
		return
	}
	v, err := t.BalanceOf(addr)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusConflict)
		return
	}
	_, err = w.(http.ResponseWriter).Write([]byte(v.String()))
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

/*
TransferToken for test only
Transfer from this node to `addr` `value` tokens on token `token`
*/
func TransferToken(w rest.ResponseWriter, r *rest.Request) {
	tokenstr := r.PathParam("token")
	addrstr := r.PathParam("addr")
	valuestr := r.PathParam("value")
	token, err := utils.HexToAddress(tokenstr)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	addr, err := utils.HexToAddress(addrstr)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	v, b := new(big.Int).SetString(valuestr, 0)
	if !b {
		rest.Error(w, "arg error ", http.StatusBadRequest)
		return
	}
	t, err := API.Photon.Chain.Token(token)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Transfer(addr, v)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = w.WriteJson("ok")
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

//EthBalance how many eth `addr` have.
func EthBalance(w rest.ResponseWriter, r *rest.Request) {
	addrstr := r.PathParam("addr")
	addr, err := utils.HexToAddress(addrstr)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	v, err := API.Photon.Chain.Client.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusConflict)
		return
	}
	_, err = w.(http.ResponseWriter).Write([]byte(v.String()))
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
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
	err := w.WriteJson(cs)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

/*
ForceUnlock force unlock by locksecrethash
*/
func ForceUnlock(w rest.ResponseWriter, r *rest.Request) {
	channelIdentifierStr := r.PathParam("channel")
	channelIdentifier := common.HexToHash(channelIdentifierStr)
	secretStr := r.PathParam("secret")
	secret := common.HexToHash(secretStr)
	err := API.ForceUnlock(channelIdentifier, secret)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = w.WriteJson("ok")
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

/*
RegisterSecretOnChain register secret to contract
*/
func RegisterSecretOnChain(w rest.ResponseWriter, r *rest.Request) {
	secretStr := r.PathParam("secret")
	secret := common.HexToHash(secretStr)
	err := API.RegisterSecretOnChain(secret)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = w.WriteJson("ok")
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}
