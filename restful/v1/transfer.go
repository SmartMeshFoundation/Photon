package v1

import (
	"math/big"
	"math/rand"
	"net/http"
	"time"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

//TransferData post for transfers
type TransferData struct {
	Initiator  string   `json:"initiator_address"`
	Target     string   `json:"target_address"`
	Token      string   `json:"token_address"`
	Amount     *big.Int `json:"amount"`
	Identifier uint64   `json:"identifier"`
	Fee        *big.Int `json:"fee"`
	IsDirect   bool     `json:"is_direct"`
}

/*
Transfers is the api of /transfer/:token/:partner
*/
func Transfers(w rest.ResponseWriter, r *rest.Request) {
	token := r.PathParam("token")
	tokenAddr := common.HexToAddress(token)
	target := r.PathParam("target")
	targetAddr := common.HexToAddress(target)
	req := &TransferData{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Identifier == 0 {
		req.Identifier = rand.New(rand.NewSource(time.Now().UnixNano())).Uint64()
	}
	if req.Amount.Cmp(utils.BigInt0) <= 0 {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Fee == nil {
		req.Fee = utils.BigInt0
	}
	if req.Fee.Cmp(utils.BigInt0) < 0 {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = RaidenAPI.Transfer(tokenAddr, req.Amount, req.Fee, targetAddr, req.Identifier, params.MaxRequestTimeout, req.IsDirect)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusConflict)
		return
	}
	req.Initiator = RaidenAPI.Raiden.NodeAddress.String()
	req.Target = target
	req.Token = token
	w.WriteJson(req)
}

/*
GetSentTransfers retuns list of sent transfer between `from_block` and `to_block`
*/
func GetSentTransfers(w rest.ResponseWriter, r *rest.Request) {
	from, to := getFromTo(r)
	log.Trace(fmt.Sprintf("from=%d,to=%d\n", from, to))
	trs, err := RaidenAPI.GetSentTransfers(from, to)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(trs)
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
	w.WriteJson(trs)
}
