package v1

import (
	"fmt"
	"net/http"

	"github.com/SmartMeshFoundation/SmartRaiden/network"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nkbai/log"
)

/*
RegisterToken register a new token to the raiden network.
this address must be a valid ERC20 token
*/
func RegisterToken(w rest.ResponseWriter, r *rest.Request) {
	token := r.PathParam("token")
	tokenAddr := common.HexToAddress(token)
	mgr, err := RaidenAPI.RegisterToken(tokenAddr)
	type Ret struct {
		ChannelManagerAddress string `json:"channel_manager_address"`
	}
	if err != nil {
		log.Error(fmt.Sprintf("RegisterToken %s err:%s", tokenAddr.String(), err))
		rest.Error(w, err.Error(), http.StatusConflict)
	} else {
		ret := &Ret{ChannelManagerAddress: mgr.String()}
		w.WriteJson(ret)
	}
}

/*
Stop for app user, call this api before quit.
*/
func Stop(w rest.ResponseWriter, r *rest.Request) {
	//test only
	RaidenAPI.Stop()
	w.Header().Set("Content-Type", "text/plain")
	w.(http.ResponseWriter).Write([]byte("ok"))
}

/*
SwitchToMesh for test only now
*/
func SwitchToMesh(w rest.ResponseWriter, r *rest.Request) {
	var nodes []*network.NodeInfo
	err := r.DecodeJsonPayload(&nodes)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = RaidenAPI.Raiden.Protocol.SwitchTransporterToMeshNetwork(nodes)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.(http.ResponseWriter).Write([]byte("ok"))
}
