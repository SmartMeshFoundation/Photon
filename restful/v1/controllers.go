package v1

import (
	"fmt"
	"net/http"

	"strconv"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
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
		err = w.WriteJson(ret)
		if err != nil {
			log.Warn(fmt.Sprintf("writejson err %s", err))
		}
	}
}

/*
Stop for app user, call this api before quit.
*/
func Stop(w rest.ResponseWriter, r *rest.Request) {
	//test only
	RaidenAPI.Stop()
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.(http.ResponseWriter).Write([]byte("ok"))
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

/*
UpdateMeshNetworkNodes update nodes of this intranet
*/
func UpdateMeshNetworkNodes(w rest.ResponseWriter, r *rest.Request) {
	var nodes []*network.NodeInfo
	err := r.DecodeJsonPayload(&nodes)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = RaidenAPI.Raiden.Protocol.UpdateMeshNetworkNodes(nodes)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.(http.ResponseWriter).Write([]byte("ok"))
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

/*
SwitchNetwork  switch between mesh and internet
*/
func SwitchNetwork(w rest.ResponseWriter, r *rest.Request) {
	mesh := r.PathParam("mesh")
	isMesh, err := strconv.ParseBool(mesh)
	if err != nil {
		rest.Error(w, "arg error", http.StatusBadRequest)
		return
	}
	RaidenAPI.Raiden.Config.IsMeshNetwork = isMesh
	_, err = w.(http.ResponseWriter).Write([]byte("ok"))
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}
