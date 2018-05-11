package v1

import (
	"fmt"
	"net/http"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

func RegisterToken(w rest.ResponseWriter, r *rest.Request) {
	token := r.PathParam("token")
	tokenAddr := common.HexToAddress(token)
	mgr, err := RaidenApi.RegisterToken(tokenAddr)
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

func Stop(w rest.ResponseWriter, r *rest.Request) {
	//test only
	RaidenApi.Stop()
	w.Header().Set("Content-Type", "text/plain")
	w.(http.ResponseWriter).Write([]byte("ok"))
}
func SwitchToMesh(w rest.ResponseWriter, r *rest.Request) {
	var nodes []*network.NodeInfo
	err := r.DecodeJsonPayload(&nodes)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(nodes) <= 0 {
		rest.Error(w, "no nodes", http.StatusBadRequest)
		return
	}
	err = RaidenApi.Raiden.Protocol.SwitchTransporterToMeshNetwork(nodes)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.(http.ResponseWriter).Write([]byte("ok"))
}
