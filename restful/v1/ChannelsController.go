package v1

import (
	"encoding/json"
	"net/http"

	"math/big"

	"github.com/SmartMeshFoundation/raiden-network/params"
	"github.com/SmartMeshFoundation/raiden-network/transfer"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type ChannelsController struct {
	BaseController
}

type channelData struct {
	ChannelAddress string   `json:"channel_address"`
	PartnerAddrses string   `json:"partner_address"`
	Balance        *big.Int `json:"balance"`
	TokenAddress   string   `json:"token_address"`
	State          string   `json:"state"`
	SettleTimeout  int      `json:"settle_timeout"`
	RevealTimeout  int      `json:"reveal_timeout"`
}

func (this *ChannelsController) Get() {
	if len(this.Ctx.Input.Param(":channel")) > 0 {
		this.SpecifiedChannel()
		return
	}
	chs, err := RaidenApi.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	if err != nil {
		this.Abort(http.StatusInternalServerError)
	}
	var datas []*channelData
	for _, c := range chs {
		d := &channelData{
			ChannelAddress: c.ChannelAddress.String(),
			PartnerAddrses: c.PartnerAddress.String(),
			Balance:        c.OurBalance,
			State:          c.State,
			TokenAddress:   c.TokenAddress.String(),
			SettleTimeout:  c.SettleTimeout,
			RevealTimeout:  c.RevealTimeout,
		}
		datas = append(datas, d)
	}
	this.Data["json"] = datas
	this.ServeJSON()
}
func (this *ChannelsController) Put() {
	this.OpenChannel()
}

func (this *ChannelsController) Patch() {
	this.CloseSettleDepositChannel()
}

//get  channel state
func (this *ChannelsController) SpecifiedChannel() {
	ch := this.Ctx.Input.Param(":channel")
	chaddr := common.HexToAddress(ch)
	c, err := RaidenApi.GetChannel(chaddr)
	if err != nil {
		this.Abort(http.StatusNotFound)
		return
	}
	d := &channelData{
		ChannelAddress: c.ChannelAddress.String(),
		PartnerAddrses: c.PartnerAddress.String(),
		Balance:        c.OurBalance,
		State:          c.State,
		SettleTimeout:  c.SettleTimeout,
		TokenAddress:   c.TokenAddress.String(),
	}
	this.Data["json"] = d
	this.ServeJSON()
}

//put request
func (this *ChannelsController) OpenChannel() {
	req := &channelData{}
	log.Trace("request body:", this.Ctx.Input.RequestBody)
	err := json.Unmarshal(this.Ctx.Input.RequestBody, req)
	if err != nil {
		log.Error(err.Error())
		this.Abort(http.StatusBadRequest)
		return
	}
	partnerAddr := common.HexToAddress(req.PartnerAddrses)
	tokenAddr := common.HexToAddress(req.TokenAddress)
	if req.State == "" { //open channel
		c, err := RaidenApi.Open(tokenAddr, partnerAddr, req.SettleTimeout, params.DEFAULT_REVEAL_TIMEOUT)
		if err != nil {
			log.Error(err.Error())
			this.Abort(http.StatusConflict)
			return
		} else {
			d := &channelData{
				ChannelAddress: c.ChannelAddress.String(),
				PartnerAddrses: c.PartnerAddress.String(),
				Balance:        c.OurBalance,
				State:          c.State,
				SettleTimeout:  c.SettleTimeout,
				TokenAddress:   c.TokenAddress.String(),
			}
			if req.Balance.Cmp(utils.BigInt0) > 0 {
				err = RaidenApi.Deposit(tokenAddr, partnerAddr, req.Balance, params.DEFAULT_POLL_TIMEOUT)
				if err == nil {
					d.Balance = c.OurBalance
				} else {
					log.Error(" RaidenApi.Deposit error : ", err)
				}
			}
			this.Data["json"] = d
			this.ServeJSON()
			return
		}
	}
	this.Abort(http.StatusBadRequest)
	return
}

func (this *ChannelsController) CloseSettleDepositChannel() {
	chstr := this.Ctx.Input.Param(":channel")
	if len(chstr) != common.AddressLength*2+2 {
		this.Abort(http.StatusConflict)
	}
	chAddr := common.HexToAddress(chstr)
	type req struct {
		State   string
		Balance *big.Int
	}
	r := &req{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, r)
	if err != nil {
		this.Abort(http.StatusBadRequest)
		return
	}
	c, err := RaidenApi.GetChannel(chAddr)
	if err != nil {
		log.Error(err.Error())
		this.Abort(http.StatusConflict)
		return
	}
	if r.Balance.Cmp(utils.BigInt0) > 0 { //deposit
		err = RaidenApi.Deposit(c.TokenAddress, c.PartnerAddress, r.Balance, params.DEFAULT_POLL_TIMEOUT)
		if err != nil {
			this.Abort(http.StatusRequestTimeout)
			return
		}
	} else {
		//close or settle
		if r.State != transfer.CHANNEL_STATE_CLOSED && r.State != transfer.CHANNEL_STATE_SETTLED {
			this.Abort(http.StatusBadRequest)
			return
		}
		if r.State == transfer.CHANNEL_STATE_CLOSED {
			c, err = RaidenApi.Close(c.TokenAddress, c.PartnerAddress)
			if err != nil {
				log.Error(err.Error())
				this.Abort(http.StatusInternalServerError)
				return
			}
		} else {
			c, err = RaidenApi.Settle(c.TokenAddress, c.PartnerAddress)
			if err != nil {
				log.Error(err.Error())
				this.Abort(http.StatusConflict)
				return
			}
		}
	}
	//reload new data from database
	c, _ = RaidenApi.GetChannel(c.ChannelAddress)
	d := &channelData{
		ChannelAddress: c.ChannelAddress.String(),
		PartnerAddrses: c.PartnerAddress.String(),
		Balance:        c.OurBalance,
		State:          c.State,
		SettleTimeout:  c.SettleTimeout,
		TokenAddress:   c.TokenAddress.String(),
	}
	this.Data["json"] = d
	this.ServeJSON()
}
