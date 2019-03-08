package v1

import (
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/Photon/pfsproxy"
	"github.com/SmartMeshFoundation/Photon/rerr"

	"github.com/SmartMeshFoundation/Photon/dto"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

/*
TokenSwap is the api of /api/1/tokenswap/:id
:id must be a unique identifier.
*/
func TokenSwap(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> TokenSwap ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	var err error
	/*
	   {
	       "role": "maker",
	       "sending_amount": 42,
	       "sending_token": "0xea674fdde714fd979de3edf0f56aa9716b898ec8",
	       "receiving_amount": 76,
	       "receiving_token": "0x2a65aca4d5fc5b5c859090a6c34d164135398226"
	   }
	*/
	// 用户调用了prepare-update,暂停接收新交易
	// client invokes prepare-update, halts receiving new transfers
	if API.Photon.StopCreateNewTransfers {
		resp = dto.NewExceptionAPIResponse(rerr.ErrStopCreateNewTransfer)
	}
	type Req struct {
		Role            string                      `json:"role"`
		SendingAmount   *big.Int                    `json:"sending_amount"`
		SendingToken    string                      `json:"sending_token"`
		ReceivingAmount *big.Int                    `json:"receiving_amount"`
		ReceivingToken  string                      `json:"receiving_token"`
		Secret          string                      `json:"secret"`     // taker无需填写,maker必填,且hash值需与url参数中的locksecrethash匹配,算法为SHA3
		RouteInfo       []pfsproxy.FindPathResponse `json:"route_info"` // 指定的路由信息
	}
	targetstr := r.PathParam("target")
	lockSecretHash := r.PathParam("locksecrethash")
	var target common.Address
	if len(targetstr) != len(target.String()) {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError)
		return
	}
	target, err = utils.HexToAddress(targetstr)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	if lockSecretHash == "" {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.Append("must provide a valid lockSecretHash "))
		return
	}
	req := &Req{}
	err = r.DecodeJsonPayload(req)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	makerToken, err := utils.HexToAddress(req.SendingToken)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	takerToken, err := utils.HexToAddress(req.ReceivingToken)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
		return
	}
	if req.Role == "maker" {
		// 校验secret和lockSecretHash是否匹配
		// check whether secret and lockSecretHash match.
		if req.Secret == "" || utils.ShaSecret(common.HexToHash(req.Secret).Bytes()) != common.HexToHash(lockSecretHash) {
			resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.Append("must provide a matching pair of secret and lockSecretHash"))
			return
		}
		err = API.TokenSwapAndWait(lockSecretHash, makerToken, takerToken,
			API.Photon.NodeAddress, target, req.SendingAmount, req.ReceivingAmount, req.Secret, req.RouteInfo)
	} else if req.Role == "taker" {
		err = API.ExpectTokenSwap(lockSecretHash, takerToken, makerToken,
			target, API.Photon.NodeAddress, req.ReceivingAmount, req.SendingAmount, req.RouteInfo)
	} else {
		err = rerr.ErrArgumentError.Errorf("Provided invalid token swap role %s", req.Role)
	}
	resp = dto.NewAPIResponse(err, nil)
}
