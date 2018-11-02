package v1

import (
	"fmt"
	"net/http"

	"strconv"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/network"
	"github.com/ant0ine/go-json-rest/rest"
)

/*
UpdateMeshNetworkNodes update nodes of this intranet
*/
func UpdateMeshNetworkNodes(w rest.ResponseWriter, r *rest.Request) {
	var err error
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> UpdateMeshNetworkNodes ,err=%v", err))
	}()
	var nodes []*network.NodeInfo
	err = r.DecodeJsonPayload(&nodes)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = API.Photon.Protocol.UpdateMeshNetworkNodes(nodes)
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
	API.Photon.Config.IsMeshNetwork = isMesh
	_, err = w.(http.ResponseWriter).Write([]byte("ok"))
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

// PrepareUpdate : 停止创建新的交易,返回当前是否可以升级
// PrepareUpdate : stop sending new transfers, return boolean that if we can update now?
func PrepareUpdate(w rest.ResponseWriter, r *rest.Request) {
	var err error
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> PrepareUpdate ,err=%v", err))
	}()
	// 这里没并发问题,直接操作即可
	// no concurrent issue, just do it.
	API.Photon.StopCreateNewTransfers = true
	num := len(API.Photon.Transfer2StateManager)
	if num > 0 {
		err = fmt.Errorf("%d transactions are still in progress. Please wait until all transactions are over", num)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = w.(http.ResponseWriter).Write([]byte("ok"))
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}
