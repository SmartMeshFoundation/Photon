package v1

import (
	"net/http"

	"math/big"

	"github.com/SmartMeshFoundation/raiden-network/params"
	"github.com/SmartMeshFoundation/raiden-network/transfer"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type channelData struct {
	ChannelAddress      string   `json:"channel_address"`
	PartnerAddrses      string   `json:"partner_address"`
	Balance             *big.Int `json:"balance"`
	PartnerBalance      *big.Int `json:"patner_balance"`
	LockedAmount        *big.Int `json:"locked_amount"`
	PartnerLockedAmount *big.Int `json:"partner_locked_amount"`
	TokenAddress        string   `json:"token_address"`
	State               string   `json:"state"`
	SettleTimeout       int      `json:"settle_timeout"`
	RevealTimeout       int      `json:"reveal_timeout"`
}

func GetChannelList(w rest.ResponseWriter, r *rest.Request) {
	chs, err := RaidenApi.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var datas []*channelData
	for _, c := range chs {
		d := &channelData{
			ChannelAddress:      c.ChannelAddress.String(),
			PartnerAddrses:      c.PartnerAddress.String(),
			Balance:             c.OurBalance,
			PartnerBalance:      c.PartnerBalance,
			State:               c.State,
			TokenAddress:        c.TokenAddress.String(),
			SettleTimeout:       c.SettleTimeout,
			RevealTimeout:       c.RevealTimeout,
			LockedAmount:        c.OurAmountLocked,
			PartnerLockedAmount: c.PartnerAmountLocked,
		}
		datas = append(datas, d)
	}
	w.WriteJson(datas)
}

//get  channel state
func SpecifiedChannel(w rest.ResponseWriter, r *rest.Request) {
	ch := r.PathParam("channel")
	chaddr := common.HexToAddress(ch)
	c, err := RaidenApi.GetChannel(chaddr)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	d := &channelData{
		ChannelAddress:      c.ChannelAddress.String(),
		PartnerAddrses:      c.PartnerAddress.String(),
		Balance:             c.OurBalance,
		PartnerBalance:      c.PartnerBalance,
		State:               c.State,
		SettleTimeout:       c.SettleTimeout,
		TokenAddress:        c.TokenAddress.String(),
		LockedAmount:        c.OurAmountLocked,
		PartnerLockedAmount: c.PartnerAmountLocked,
	}
	w.WriteJson(d)
}

//put request
func OpenChannel(w rest.ResponseWriter, r *rest.Request) {
	req := &channelData{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	partnerAddr := common.HexToAddress(req.PartnerAddrses)
	tokenAddr := common.HexToAddress(req.TokenAddress)
	if req.State == "" { //open channel
		c, err := RaidenApi.Open(tokenAddr, partnerAddr, req.SettleTimeout, params.DEFAULT_REVEAL_TIMEOUT)
		if err != nil {
			log.Error(err.Error())
			rest.Error(w, err.Error(), http.StatusConflict)
			return
		} else {
			d := &channelData{
				ChannelAddress:      c.ChannelAddress.String(),
				PartnerAddrses:      c.PartnerAddress.String(),
				Balance:             c.OurBalance,
				PartnerBalance:      c.PartnerBalance,
				State:               c.State,
				SettleTimeout:       c.SettleTimeout,
				TokenAddress:        c.TokenAddress.String(),
				LockedAmount:        c.OurAmountLocked,
				PartnerLockedAmount: c.PartnerAmountLocked,
			}
			if req.Balance.Cmp(utils.BigInt0) > 0 {
				err = RaidenApi.Deposit(tokenAddr, partnerAddr, req.Balance, params.DEFAULT_POLL_TIMEOUT)
				if err == nil {
					c, _ := RaidenApi.GetChannel(c.ChannelAddress)
					d.Balance = c.OurBalance
				} else {
					log.Error(" RaidenApi.Deposit error : ", err)
				}
			}
			data,err:= w.EncodeJson(d)
			if err!=nil{
				log.Error(err.Error())
				rest.Error(w,err.Error(),http.StatusConflict)
				return
			}
			rest.Error(w,string(data),http.StatusCreated)
			return
		}
	}
	rest.Error(w, "argument error", http.StatusBadRequest)
	return
}

func CloseSettleDepositChannel(w rest.ResponseWriter, r *rest.Request) {
	chstr := r.PathParam("channel")
	if len(chstr) != len(utils.EmptyAddress.String()) {
		rest.Error(w, "argument error", http.StatusBadRequest)
	}
	chAddr := common.HexToAddress(chstr)
	type Req struct {
		State   string
		Balance *big.Int
	}
	req := &Req{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c, err := RaidenApi.GetChannel(chAddr)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if req.Balance != nil && req.Balance.Cmp(utils.BigInt0) > 0 { //deposit
		err = RaidenApi.Deposit(c.TokenAddress, c.PartnerAddress, req.Balance, params.DEFAULT_POLL_TIMEOUT)
		if err != nil {
			rest.Error(w, err.Error(), http.StatusRequestTimeout)
			return
		}
	} else {
		//close or settle
		if req.State != transfer.CHANNEL_STATE_CLOSED && req.State != transfer.CHANNEL_STATE_SETTLED {
			rest.Error(w, "argument error", http.StatusBadRequest)
			return
		}
		if req.State == transfer.CHANNEL_STATE_CLOSED {
			c, err = RaidenApi.Close(c.TokenAddress, c.PartnerAddress)
			if err != nil {
				log.Error(err.Error())
				rest.Error(w, err.Error(), http.StatusConflict)
				return
			}
		} else {
			c, err = RaidenApi.Settle(c.TokenAddress, c.PartnerAddress)
			if err != nil {
				log.Error(err.Error())
				rest.Error(w, err.Error(), http.StatusConflict)
				return
			}
		}
	}
	//reload new data from database
	c, _ = RaidenApi.GetChannel(c.ChannelAddress)
	d := &channelData{
		ChannelAddress:      c.ChannelAddress.String(),
		PartnerAddrses:      c.PartnerAddress.String(),
		Balance:             c.OurBalance,
		PartnerBalance:      c.PartnerBalance,
		State:               c.State,
		SettleTimeout:       c.SettleTimeout,
		TokenAddress:        c.TokenAddress.String(),
		LockedAmount:        c.OurAmountLocked,
		PartnerLockedAmount: c.PartnerAmountLocked,
	}
	w.WriteJson(d)
}
