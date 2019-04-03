package v1

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/rerr"

	"strconv"

	"github.com/SmartMeshFoundation/Photon/dto"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/network"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

/*
UpdateMeshNetworkNodes update nodes of this intranet
*/
func UpdateMeshNetworkNodes(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> UpdateMeshNetworkNodes ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	var err error
	var nodes []*network.NodeInfo
	err = r.DecodeJsonPayload(&nodes)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError)
		return
	}
	err = API.Photon.Protocol.UpdateMeshNetworkNodes(nodes)
	resp = dto.NewAPIResponse(err, "ok")
}

/*
SwitchNetwork  switch between mesh and internet
*/
func SwitchNetwork(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> SwitchNetwork ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	mesh := r.PathParam("mesh")
	isMesh, err := strconv.ParseBool(mesh)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError)
		return
	}
	log.Trace(fmt.Sprintf("Api SwitchNetwork isMesh=%v", isMesh))
	if isMesh {
		log.Trace("use NotifyNetworkDown")
		err := API.NotifyNetworkDown()
		if err != nil {
			panic("never happen")
		}
	}
	resp = dto.NewSuccessAPIResponse(nil)
}

// PrepareUpdate : 停止创建新的交易,返回当前是否可以升级
// PrepareUpdate : stop sending new transfers, return boolean that if we can update now?
func PrepareUpdate(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> PrepareUpdate ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	var err error
	// 这里没并发问题,直接操作即可
	// no concurrent issue, just do it.
	API.Photon.StopCreateNewTransfers = true
	num := len(API.Photon.Transfer2StateManager)
	if num > 0 {
		err = fmt.Errorf("%d transactions are still in progress. Please wait until all transactions are over", num)
		resp = dto.NewExceptionAPIResponse(rerr.ErrUpdateButHaveTransfer.AppendError(err))
		return
	}
	resp = dto.NewSuccessAPIResponse(nil)
}

// NotifyNetworkDown :
func NotifyNetworkDown(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> NotifyNetworkDown ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	err := API.NotifyNetworkDown()
	resp = dto.NewAPIResponse(err, "ok")
}

// GetFeePolicy :
func GetFeePolicy(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> GetFeePolicy ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	fp, err := API.GetFeePolicy()
	resp = dto.NewAPIResponse(err, fp)
}

// SetFeePolicy :
func SetFeePolicy(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> SetFeePolicy ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	req := &models.FeePolicy{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	err = API.SetFeePolicy(req)
	resp = dto.NewAPIResponse(err, "ok")
}

// FindPath :
func FindPath(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> FindPath ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	targetAddressStr := r.PathParam("target_address")
	targetAddress, err := utils.HexToAddress(targetAddressStr)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	tokenAddressStr := r.PathParam("token")
	tokenAddress, err := utils.HexToAddress(tokenAddressStr)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	amountStr := r.PathParam("amount")
	amount, ok := math.ParseBig256(amountStr)
	if !ok {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError)
		return
	}
	result, err := API.FindPath(targetAddress, tokenAddress, amount)
	resp = dto.NewAPIResponse(err, result)

}

// GetAllFeeChargeRecord :
func GetAllFeeChargeRecord(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> GetAllFeeChargeRecord ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	result, err := API.GetAllFeeChargeRecord()
	resp = dto.NewAPIResponse(err, result)
}

// GetSystemStatus :
func GetSystemStatus(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> GetSystemStatus ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	result, err := API.SystemStatus()
	resp = dto.NewAPIResponse(err, result)
}

// GetIncomeDetailsRequest :
type GetIncomeDetailsRequest struct {
	TokenAddress string `json:"token_address"`
	FromTime     int64  `json:"from_time"`
	ToTime       int64  `json:"to_time"`
	Limit        int    `json:"limit"`
}

// GetIncomeDetails :
func GetIncomeDetails(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> GetIncomeDetails ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	req := &GetIncomeDetailsRequest{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	tokenAddress := common.HexToAddress(req.TokenAddress)
	result, err := API.GetIncomeDetails(tokenAddress, req.FromTime, req.ToTime, req.Limit)
	resp = dto.NewAPIResponse(err, result)
}

// GetOneWeekIncomeRequest :
type GetOneWeekIncomeRequest struct {
	TokenAddress string `json:"token_address"`
	Days         int    `json:"days"`
}

// GetDaysIncome :
func GetDaysIncome(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> GetDaysIncome ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	req := &GetOneWeekIncomeRequest{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	tokenAddress := common.HexToAddress(req.TokenAddress)
	result, err := API.GetDaysIncome(tokenAddress, req.Days)
	resp = dto.NewAPIResponse(err, result)
}

//GetBuildInfo :
func GetBuildInfo(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> GetBuildInfo ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	resp = dto.NewSuccessAPIResponse(API.GetBuildInfo())
}
