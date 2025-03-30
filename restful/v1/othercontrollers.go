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
		API.NotifyNetworkDown()
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

// NotifyNetworkDown 上层应用通知photon网络断开,强制photon对所有远程连接进行重连尝试
func NotifyNetworkDown(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> NotifyNetworkDown ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	API.NotifyNetworkDown()
	resp = dto.NewSuccessAPIResponse(nil)
}

// GetFeePolicy 查询节点手续费收费策略
func GetFeePolicy(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> GetFeePolicy ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	fp, err := API.GetFeePolicy()
	resp = dto.NewAPIResponse(err, fp)
}

// SetFeePolicy 全量更新节点手续费收费策略
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

// FindPath 通过photon调用pfs的路由查询服务
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

// GetAllFeeChargeRecord 查询节点收取手续费的记录
func GetAllFeeChargeRecord(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> GetAllFeeChargeRecord ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	result, err := API.GetAllFeeChargeRecord()
	resp = dto.NewAPIResponse(err, result)
}

// GetSystemStatus 查询节点汇总信息
func GetSystemStatus(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> GetSystemStatus ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	result, err := API.SystemStatus()
	resp = dto.NewAPIResponse(err, result)
}

// GetIncomeDetailsRequest GetIncomeDetails接口的返回结构
type GetIncomeDetailsRequest struct {
	TokenAddress string `json:"token_address"`
	FromTime     int64  `json:"from_time"`
	ToTime       int64  `json:"to_time"`
	Limit        int    `json:"limit"`
}

// GetIncomeDetails 查询photon节点收益信息明细
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

// GetOneWeekIncomeRequest 过去N天收益信息报表查询参数
type GetOneWeekIncomeRequest struct {
	TokenAddress string `json:"token_address"`
	Days         int    `json:"days"`
}

// GetDaysIncome 过去N天收益信息报表查询
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

//GetBuildInfo 查询版本构建信息
func GetBuildInfo(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> GetBuildInfo ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	resp = dto.NewSuccessAPIResponse(API.GetBuildInfo())
}

// GetAssetsOnTokenRequest GetAssetsOnToken接口查询参数
type GetAssetsOnTokenRequest struct {
	TokenList []string `json:"token_list"`
}

//GetAssetsOnToken 查询链上及photon资产信息
func GetAssetsOnToken(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> GetAssetsOnToken ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	req := &GetAssetsOnTokenRequest{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	var tokenList []common.Address
	for _, token := range req.TokenList {
		tokenList = append(tokenList, common.HexToAddress(token))
	}
	data, err := API.GetAssetsOnToken(tokenList)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(err)
	} else {
		resp = dto.NewSuccessAPIResponse(data)
	}
}
