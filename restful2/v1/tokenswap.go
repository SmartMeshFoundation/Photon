package v1

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

func TokenSwap(w rest.ResponseWriter, r *rest.Request) {
	/*
	   {
	       "role": "maker",
	       "sending_amount": 42,
	       "sending_token": "0xea674fdde714fd979de3edf0f56aa9716b898ec8",
	       "receiving_amount": 76,
	       "receiving_token": "0x2a65aca4d5fc5b5c859090a6c34d164135398226"
	   }
	*/
	type Req struct {
		Role            string   `json:"role"`
		SendingAmount   *big.Int `json:"sending_amount"`
		SendingToken    string   `json:"sending_token"`
		ReceivingAmount *big.Int `json:"receiving_amount"`
		ReceivingToken  string   `json:"receiving_token"`
	}
	targetstr := r.PathParam("target")
	idstr := r.PathParam("id")
	var target common.Address
	var id int
	if len(targetstr) != len(target.String()) {
		rest.Error(w, "target address error", http.StatusBadRequest)
		return
	}
	target = common.HexToAddress(targetstr)
	id, _ = strconv.Atoi(idstr)
	if id <= 0 {
		rest.Error(w, "must provide a valid id ", http.StatusBadRequest)
		return
	}
	req := &Req{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Role == "maker" {
		err = RaidenApi.TokenSwapAndWait(uint64(id), common.HexToAddress(req.SendingToken), common.HexToAddress(req.ReceivingToken),
			RaidenApi.Raiden.NodeAddress, target, req.SendingAmount, req.ReceivingAmount)
	} else if req.Role == "taker" {
		err = RaidenApi.ExpectTokenSwap(uint64(id), common.HexToAddress(req.ReceivingToken), common.HexToAddress(req.SendingToken),
			target, RaidenApi.Raiden.NodeAddress, req.ReceivingAmount, req.SendingAmount)
	} else {
		err = fmt.Errorf("Provided invalid token swap role %s", req.Role)
	}
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		rest.Error(w, "", http.StatusCreated)
	}
}
