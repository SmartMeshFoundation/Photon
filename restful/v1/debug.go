package v1

import (
	"math/big"
	"net/http"

	"context"

	"github.com/SmartMeshFoundation/SmartRaiden/network/xmpptransport"
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
	token := common.HexToAddress(tokenstr)
	addr := common.HexToAddress(addrstr)
	t := RaidenAPI.Raiden.Chain.Token(token)
	v, err := t.BalanceOf(addr)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.(http.ResponseWriter).Write([]byte(v.String()))
}

/*
TransferToken for test only
Transfer from this node to `addr` `value` tokens on token `token`
*/
func TransferToken(w rest.ResponseWriter, r *rest.Request) {
	tokenstr := r.PathParam("token")
	addrstr := r.PathParam("addr")
	valuestr := r.PathParam("value")
	token := common.HexToAddress(tokenstr)
	addr := common.HexToAddress(addrstr)
	v, b := new(big.Int).SetString(valuestr, 0)
	if !b {
		rest.Error(w, "arg error ", http.StatusBadRequest)
		return
	}
	t := RaidenAPI.Raiden.Chain.Token(token)
	err := t.Transfer(addr, v)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson("ok")
}

//EthBalance how many eth `addr` have.
func EthBalance(w rest.ResponseWriter, r *rest.Request) {
	addrstr := r.PathParam("addr")
	addr := common.HexToAddress(addrstr)
	v, err := RaidenAPI.Raiden.Chain.Client.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.(http.ResponseWriter).Write([]byte(v.String()))
}

//BlockTimeFormat  is time format of last block
const BlockTimeFormat = "01-02|15:04:05.999"

//ConnectionStatus status of network connection
type ConnectionStatus struct {
	XMPPStatus    xmpptransport.Status
	EthStatus     xmpptransport.Status
	LastBlockTime string
}

/*
EthereumStatus  query the status between raiden and ethereum
*/
func EthereumStatus(w rest.ResponseWriter, r *rest.Request) {
	c := RaidenAPI.Raiden.Chain
	cs := &ConnectionStatus{
		XMPPStatus:    xmpptransport.Disconnected,
		LastBlockTime: RaidenAPI.Raiden.GetDb().GetLastBlockNumberTime().Format(BlockTimeFormat),
	}
	if c != nil && c.Client.Status == xmpptransport.Connected {
		cs.EthStatus = xmpptransport.Connected
	} else {
		cs.EthStatus = xmpptransport.Disconnected
	}
	w.WriteJson(cs)
}
