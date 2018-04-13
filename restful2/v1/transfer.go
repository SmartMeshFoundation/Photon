package v1

import (
	"math/big"
	"math/rand"
	"net/http"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

//post for transfers
type transferData struct {
	Initiator  string   `json:"initiator_address"`
	Target     string   `json:"target_address"`
	Token      string   `json:"token_address"`
	Amount     *big.Int `json:"amount"`
	Identifier uint64   `json:"identifier"`
	Fee        *big.Int `json:"fee"`
}

func Transfers(w rest.ResponseWriter, r *rest.Request) {
	token := r.PathParam("token")
	tokenAddr := common.HexToAddress(token)
	target := r.PathParam("target")
	targetAddr := common.HexToAddress(target)
	req := &transferData{}
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
	err = RaidenApi.Transfer(tokenAddr, req.Amount, req.Fee, targetAddr, req.Identifier, params.MaxRequestTimeout)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusConflict)
		return
	}
	req.Initiator = RaidenApi.Raiden.NodeAddress.String()
	req.Target = target
	req.Token = token
	w.WriteJson(req)
}
