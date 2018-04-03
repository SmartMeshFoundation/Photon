package v1

import (
	"fmt"

	"net/http"

	"encoding/json"

	"github.com/SmartMeshFoundation/raiden-network/params"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"

	"math/rand"
	"time"

	"strconv"

	"math/big"

	"github.com/ethereum/go-ethereum/log"
)

type DataMap map[string]interface{}

type Controller struct {
	BaseController
}

func (this *Controller) Address() {
	data := make(DataMap)
	data["our_address"] = RaidenApi.Raiden.NodeAddress.String()
	this.Data["json"] = data
	this.ServeJSON()
}

func (this *Controller) Tokens() {
	this.Data["json"] = RaidenApi.GetTokenList()
	this.ServeJSON()
}

type partnersData struct {
	PartnerAddress string `json:"partner_address"`
	Channel        string `json:"channel"`
}

func (this *Controller) TokenPartners() {
	tokenAddr := common.HexToAddress(this.Ctx.Input.Param(":token"))
	log.Trace(fmt.Sprintf("TokenPartners tokenAddr=%s", utils.APex(tokenAddr)))
	chs, err := RaidenApi.GetChannelList(tokenAddr, utils.EmptyAddress)
	if err != nil {
		this.Abort(http.StatusInternalServerError)
		return
	}
	var datas []*partnersData
	for _, c := range chs {
		d := &partnersData{
			PartnerAddress: c.PartnerAddress.String(),
			Channel:        "api/1/channles/" + c.ChannelAddress.String(),
		}
		datas = append(datas, d)
	}
	this.Data["json"] = datas
	this.ServeJSON()
}

func (this *Controller) RegisterToken() {
	token := this.Ctx.Input.Param(":token")
	tokenAddr := common.HexToAddress(token)
	mgr, err := RaidenApi.RegisterToken(tokenAddr)
	type Ret struct {
		Channel_manager_address string
	}
	if err != nil {
		log.Error(fmt.Sprintf("RegisterToken %s err:%s", tokenAddr.String(), err))
		this.Abort(http.StatusConflict)
	} else {
		r := &Ret{Channel_manager_address: mgr.String()}
		this.Data["json"] = r
		this.ServeJSON()
	}
}

//post for transfers
type transferData struct {
	/*
			  "initiator_address": "0xea674fdde714fd979de3edf0f56aa9716b898ec8",
		    "target_address": "0x61c808d82a3ac53231750dadc13c777b59310bd9",
		    "token_address": "0x2a65aca4d5fc5b5c859090a6c34d164135398226",
		    "amount": 200,
		    "identifier": 42
	*/
	Initiator  string   `json:"initiator_address"`
	Target     string   `json:"target_address"`
	Token      string   `json:"token_address"`
	Amount     *big.Int `json:"amount"`
	Identifier uint64   `json:"identifier"`
	Fee        *big.Int `json:"fee"`
}

func (this *Controller) Transfers() {
	token := this.Ctx.Input.Param(":token")
	tokenAddr := common.HexToAddress(token)
	target := this.Ctx.Input.Param(":target")
	targetAddr := common.HexToAddress(target)
	log.Trace(fmt.Sprintf("request body:%s", this.Ctx.Input.RequestBody))
	req := &transferData{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, req)
	if err != nil {
		log.Error(err.Error())
		this.Abort(http.StatusBadRequest)
		return
	}
	if req.Identifier == 0 {
		req.Identifier = rand.New(rand.NewSource(time.Now().UnixNano())).Uint64()
	}
	if req.Amount.Cmp(utils.BigInt0) <= 0 {
		this.Abort(http.StatusBadRequest)
		return
	}
	if req.Fee == nil {
		req.Fee = utils.BigInt0
	}
	if req.Fee.Cmp(utils.BigInt0) < 0 {
		this.Abort(http.StatusBadRequest)
		return
	}
	err = RaidenApi.Transfer(tokenAddr, req.Amount, req.Fee, targetAddr, req.Identifier, params.MaxRequestTimeout)
	if err != nil {
		this.Abort(http.StatusConflict)
		return
	}
	req.Initiator = RaidenApi.Raiden.NodeAddress.String()
	req.Target = target
	req.Token = token
	this.Data["json"] = req
	this.ServeJSON()
}

func (this *Controller) TokenSwap() {
	/*
	   {
	       "role": "maker",
	       "sending_amount": 42,
	       "sending_token": "0xea674fdde714fd979de3edf0f56aa9716b898ec8",
	       "receiving_amount": 76,
	       "receiving_token": "0x2a65aca4d5fc5b5c859090a6c34d164135398226"
	   }
	*/
	type Req struct {
		Role            string   `json:"role"`
		SendingAmount   *big.Int `json:"sending_amount"`
		SendingToken    string   `json:"sending_token"`
		ReceivingAmount *big.Int `json:"receiving_amount"`
		ReceivingToken  string   `json:"receiving_token"`
	}
	targetstr := this.Ctx.Input.Param(":target")
	idstr := this.Ctx.Input.Param(":id")
	var target common.Address
	var id int
	if len(targetstr) != len(target.String()) {
		this.Abort(http.StatusBadRequest)
		return
	}
	target = common.HexToAddress(targetstr)
	id, _ = strconv.Atoi(idstr)
	if id <= 0 {
		this.Abort(http.StatusBadRequest)
		return
	}
	req := &Req{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, req)
	if err != nil {
		log.Error(err.Error())
		this.Abort(http.StatusBadRequest)
		return
	}
	if req.Role == "maker" {
		err = RaidenApi.TokenSwapAndWait(uint64(id), common.HexToAddress(req.SendingToken), common.HexToAddress(req.ReceivingToken),
			RaidenApi.Raiden.NodeAddress, target, req.SendingAmount, req.ReceivingAmount)
	} else if req.Role == "taker" {
		err = RaidenApi.ExpectTokenSwap(uint64(id), common.HexToAddress(req.ReceivingToken), common.HexToAddress(req.SendingToken),
			target, RaidenApi.Raiden.NodeAddress, req.ReceivingAmount, req.SendingAmount)
	} else {
		err = fmt.Errorf("Provided invalid token swap role %s", req.Role)
	}
	if err != nil {
		log.Error(err.Error())
		this.Abort(http.StatusBadRequest)
	} else {
		this.Abort(http.StatusCreated)
	}
}

func (this *Controller) Stop() {
	//test only
	RaidenApi.Stop()
	this.Ctx.WriteString("ok")
}
