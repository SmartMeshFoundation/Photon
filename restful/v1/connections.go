package v1

import (
	"net/http"

	"math/big"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

/*
GetConnections Get a dict whose keys are token addresses and whose values are
open channels, funds of last request, sum of deposits and number of channels
*/
func GetConnections(w rest.ResponseWriter, r *rest.Request) {
	err := w.WriteJson(RaidenAPI.GetConnectionManagersInfo())
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

/*
ConnectTokenNetwork open channels with existing addresses on this token network.
and deposit to the new  opened channel
*/
func ConnectTokenNetwork(w rest.ResponseWriter, r *rest.Request) {
	type Req struct {
		Funds *big.Int `json:"funds"`
	}
	tokenstr := r.PathParam("token")
	var token common.Address
	if len(tokenstr) != len(token.String()) {
		rest.Error(w, "argument error", http.StatusBadRequest)
		return
	}
	token = common.HexToAddress(tokenstr)
	req := &Req{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Trace(fmt.Sprintf("req=%#v", req))
	if req.Funds == nil || req.Funds.Cmp(utils.BigInt0) <= 0 {
		rest.Error(w, "funds error", http.StatusBadRequest)
		return
	}
	err = RaidenAPI.ConnectTokenNetwork(token, req.Funds, params.DefaultInitialChannelTarget, params.DefaultJoinableFundsTarget)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.(http.ResponseWriter).WriteHeader(http.StatusCreated)
	_, err = w.(http.ResponseWriter).Write(nil)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

/*
LeaveTokenNetwork may take very long time.
1.close all the channels on this token network
2. waiting for the settle time
3. settle all the channel.
*/
func LeaveTokenNetwork(w rest.ResponseWriter, r *rest.Request) {
	type Req struct {
		OnlyReceivingChannels bool `json:"only_receiving_channels"`
	}
	tokenstr := r.PathParam("token")
	var token common.Address
	if len(tokenstr) != len(token.String()) {
		rest.Error(w, "argument error", http.StatusBadRequest)
		return
	}
	token = common.HexToAddress(tokenstr)
	req := &Req{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		rest.Error(w, "argument error", http.StatusBadRequest)
		return
	}
	chs, err := RaidenAPI.LeaveTokenNetwork(token, req.OnlyReceivingChannels)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if len(chs) > 0 {
		var addrs []string
		for _, c := range chs {
			addrs = append(addrs, c.ChannelIdentifier.ChannelIdentifier.String())
		}
		err = w.WriteJson(addrs)
		if err != nil {
			log.Warn(fmt.Sprintf("writejson err %s", err))
		}
		return
	} else {
		w.(http.ResponseWriter).WriteHeader(http.StatusCreated)
		_, err = w.(http.ResponseWriter).Write(nil)
		if err != nil {
			log.Warn(fmt.Sprintf("writejson err %s", err))
		}
	}

}
