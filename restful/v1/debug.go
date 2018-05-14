package v1

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
	"net/http"
	"math/big"
)

impithub.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
	"net/http"
	"math/big"
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
	t := RaidenApi.Raiden.Chain.Token(token)
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
	t := RaidenApi.Raiden.Chain.Token(token)
	err := t.Transfer(addr, v)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.(http.ResponseWriter).Write([]byte("ok"))

}
