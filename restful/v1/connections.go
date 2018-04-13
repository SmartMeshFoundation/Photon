package v1

import (
	"net/http"

	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

func GetConnections(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(RaidenApi.GetConnectionManagersInfo())
}

//Connecting to a token network
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
	if req.Funds.Cmp(utils.BigInt0) <= 0 {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = RaidenApi.ConnectTokenNetwork(token, req.Funds, params.DEFAULT_INITIAL_CHANNEL_TARGET, params.DEFAULT_JOINABLE_FUNDS_TARGET)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.(http.ResponseWriter).WriteHeader(http.StatusCreated)
	w.(http.ResponseWriter).Write(nil)
}

//leave a token network
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
	chs, err := RaidenApi.LeaveTokenNetwork(token, req.OnlyReceivingChannels)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if len(chs) > 0 {
		var addrs []string
		for _, c := range chs {
			addrs = append(addrs, c.ChannelAddress.String())
		}
		w.WriteJson(addrs)
		return
	} else {
		w.(http.ResponseWriter).WriteHeader(http.StatusCreated)
		w.(http.ResponseWriter).Write(nil)
	}

}
