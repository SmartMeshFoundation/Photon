package v1

import (
	"fmt"
	"net/http"

	"strconv"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/network"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common/math"
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

// NotifyNetworkDown :
func NotifyNetworkDown(w rest.ResponseWriter, r *rest.Request) {
	err := API.NotifyNetworkDown()
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = w.(http.ResponseWriter).Write([]byte("ok"))
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

// GetFeePolicy :
func GetFeePolicy(w rest.ResponseWriter, r *rest.Request) {
	fp, err := API.GetFeePolicy()
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = w.WriteJson(fp)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

// SetFeePolicy :
func SetFeePolicy(w rest.ResponseWriter, r *rest.Request) {
	req := &models.FeePolicy{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = API.SetFeePolicy(req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = w.(http.ResponseWriter).Write([]byte("ok"))
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

// FindPath :
func FindPath(w rest.ResponseWriter, r *rest.Request) {
	targetAddressStr := r.PathParam("target_address")
	targetAddress, err := utils.HexToAddress(targetAddressStr)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tokenAddressStr := r.PathParam("token")
	tokenAddress, err := utils.HexToAddress(tokenAddressStr)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	amountStr := r.PathParam("amount")
	amount, ok := math.ParseBig256(amountStr)
	if !ok {
		rest.Error(w, "wrong amount", http.StatusBadRequest)
		return
	}
	resp, err := API.FindPath(targetAddress, tokenAddress, amount)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = w.WriteJson(resp)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

// GetAllFeeChargeRecord :
func GetAllFeeChargeRecord(w rest.ResponseWriter, r *rest.Request) {
	err := w.WriteJson(API.GetAllFeeChargeRecord())
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

// GetSystemStatus :
func GetSystemStatus(w rest.ResponseWriter, r *rest.Request) {
	err := w.WriteJson(API.SystemStatus())
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

// GenerateSecret : generate debug secret for test
func GenerateSecret(w rest.ResponseWriter, r *rest.Request) {
	type resp struct {
		Secret     string
		SecretHash string
	}

	secret := utils.NewRandomHash()
	rs := resp{
		Secret:     secret.String(),
		SecretHash: utils.ShaSecret(secret[:]).String(),
	}
	err := w.WriteJson(rs)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}
