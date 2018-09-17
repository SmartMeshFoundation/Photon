package v1

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

/*
TokenSwap is the api of /api/1/tokenswap/:id
:id must be a unique identifier.
*/
func TokenSwap(w rest.ResponseWriter, r *rest.Request) {
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
	if RaidenAPI.Raiden.StopCreateNewTransfers {
		rest.Error(w, "Stop create new transfers, please restart smartraiden", http.StatusBadRequest)
		return
	}
	type Req struct {
		Role            string   `json:"role"`
		SendingAmount   *big.Int `json:"sending_amount"`
		SendingToken    string   `json:"sending_token"`
		ReceivingAmount *big.Int `json:"receiving_amount"`
		ReceivingToken  string   `json:"receiving_token"`
		Secret          string   `json:"secret"` // taker无需填写,maker必填,且hash值需与url参数中的locksecrethash匹配,算法为SHA3
	}
	targetstr := r.PathParam("target")
	lockSecretHash := r.PathParam("locksecrethash")
	var target common.Address
	if len(targetstr) != len(target.String()) {
		rest.Error(w, "target address error", http.StatusBadRequest)
		return
	}
	target, err := utils.HexToAddress(targetstr)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if lockSecretHash == "" {
		rest.Error(w, "must provide a valid lockSecretHash ", http.StatusBadRequest)
		return
	}
	req := &Req{}
	err = r.DecodeJsonPayload(req)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	makerToken, err := utils.HexToAddress(req.SendingToken)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	takerToken, err := utils.HexToAddress(req.ReceivingToken)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Role == "maker" {
		// 校验secret和lockSecretHash是否匹配
		if req.Secret == "" || utils.ShaSecret(common.HexToHash(req.Secret).Bytes()) != common.HexToHash(lockSecretHash) {
			rest.Error(w, "must provide a matching pair of secret and lockSecretHash", http.StatusBadRequest)
			return
		}
		err = RaidenAPI.TokenSwapAndWait(lockSecretHash, makerToken, takerToken,
			RaidenAPI.Raiden.NodeAddress, target, req.SendingAmount, req.ReceivingAmount, req.Secret)
	} else if req.Role == "taker" {
		err = RaidenAPI.ExpectTokenSwap(lockSecretHash, takerToken, makerToken,
			target, RaidenAPI.Raiden.NodeAddress, req.ReceivingAmount, req.SendingAmount)
	} else {
		err = fmt.Errorf("Provided invalid token swap role %s", req.Role)
	}
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.(http.ResponseWriter).WriteHeader(http.StatusCreated)
		_, err = w.(http.ResponseWriter).Write(nil)
		if err != nil {
			log.Warn(fmt.Sprintf("writejson err %s", err))
		}
	}
}
