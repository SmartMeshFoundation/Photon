package v1

import (
	"net/http"

	"math/big"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

/*
ChannelData export json data format
*/
type ChannelData struct {
	ChannelIdentifier   string            `json:"channel_identifier"`
	OpenBlockNumber     int64             `json:"open_block_number"`
	PartnerAddrses      string            `json:"partner_address"`
	Balance             *big.Int          `json:"balance"`
	PartnerBalance      *big.Int          `json:"partner_balance"`
	LockedAmount        *big.Int          `json:"locked_amount"`
	PartnerLockedAmount *big.Int          `json:"partner_locked_amount"`
	TokenAddress        string            `json:"token_address"`
	State               channeltype.State `json:"state"`
	StateString         string            `json:"state_string"`
	SettleTimeout       int               `json:"settle_timeout"`
	RevealTimeout       int               `json:"reveal_timeout"`
}

/*
GetChannelList list all my channels
*/
func GetChannelList(w rest.ResponseWriter, r *rest.Request) {
	chs, err := API.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var datas []*ChannelData
	for _, c := range chs {
		d := &ChannelData{
			ChannelIdentifier:   c.ChannelIdentifier.ChannelIdentifier.String(),
			OpenBlockNumber:     c.ChannelIdentifier.OpenBlockNumber,
			PartnerAddrses:      c.PartnerAddress().String(),
			Balance:             c.OurBalance(),
			PartnerBalance:      c.PartnerBalance(),
			State:               c.State,
			StateString:         c.State.String(),
			TokenAddress:        c.TokenAddress().String(),
			SettleTimeout:       c.SettleTimeout,
			RevealTimeout:       c.RevealTimeout,
			LockedAmount:        c.OurAmountLocked(),
			PartnerLockedAmount: c.PartnerAmountLocked(),
		}
		datas = append(datas, d)
	}
	err = w.WriteJson(datas)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

/*
ChannelFor3rdParty generate info for 3rd party use,
for update transfer and withdraw.
/api/1/thirdparty/:channel/:3rd
channel:  channel identifier
3rd: the account that MS used to send transaction on blockchain
*/
func ChannelFor3rdParty(w rest.ResponseWriter, r *rest.Request) {
	ch := r.PathParam("channel")
	thirdParty := r.PathParam("3rd")
	channelIdentifier := common.HexToHash(ch)
	thirdAddress, err := utils.HexToAddress(thirdParty)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if channelIdentifier == utils.EmptyHash || thirdAddress == utils.EmptyAddress {
		rest.Error(w, "argument error", http.StatusBadRequest)
		return
	}
	result, err := API.ChannelInformationFor3rdParty(channelIdentifier, thirdAddress)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = w.WriteJson(result)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

/*
SpecifiedChannel get  a channel state
*/
func SpecifiedChannel(w rest.ResponseWriter, r *rest.Request) {
	ch := r.PathParam("channel")
	channelIdentifier := common.HexToHash(ch)
	c, err := API.GetChannel(channelIdentifier)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	d := channeltype.ChannelSerialization2ChannelDataDetail(c)
	err = w.WriteJson(d)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

/*
depositReq 用户存款请求
*/
type depositReq struct {
	PartnerAddrses string   `json:"partner_address"` //通道对方地址
	TokenAddress   string   `json:"token_address"`   //哪种token
	Balance        *big.Int `json:"balance"`         //存入金额,一定大于0
	//如果NewChannel为true
	//  SettleTimeout表示新建通道的结算窗口,如果SettleTimeout为0,则用系统默认计算窗口
	//如果NewChannel为 false
	//  SettleTimeout 必须为0,表示只是存款,一定不要创建通道
	SettleTimeout int  `json:"settle_timeout"`
	NewChannel    bool `json:"new_channel"` //此次行为是创建通道并存款还是只存款
}

/*
Deposit open a channel with partner if channel not exists  and deposit `balance` token to this channel.
token must exist
partner maybe an invalid address
Balance must be positive
*/
func Deposit(w rest.ResponseWriter, r *rest.Request) {
	var err error
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> Deposit ,err=%v", err))
	}()
	req := &depositReq{}
	err = r.DecodeJsonPayload(req)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	partnerAddr, err := utils.HexToAddress(req.PartnerAddrses)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tokenAddr, err := utils.HexToAddress(req.TokenAddress)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err := API.DepositAndOpenChannel(tokenAddr, partnerAddr, req.SettleTimeout, API.Photon.Config.RevealTimeout, req.Balance, req.NewChannel)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusConflict)
		return
	}
	var d *ChannelData
	if !req.NewChannel {
		d = &ChannelData{
			ChannelIdentifier:   c.ChannelIdentifier.ChannelIdentifier.String(),
			OpenBlockNumber:     c.ChannelIdentifier.OpenBlockNumber,
			PartnerAddrses:      c.PartnerAddress().String(),
			Balance:             c.OurBalance(),
			PartnerBalance:      c.PartnerBalance(),
			State:               c.State,
			StateString:         c.State.String(),
			SettleTimeout:       c.SettleTimeout,
			TokenAddress:        c.TokenAddress().String(),
			LockedAmount:        c.OurAmountLocked(),
			PartnerLockedAmount: c.PartnerAmountLocked(),
			RevealTimeout:       c.RevealTimeout,
		}

	}
	err = w.WriteJson(d)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
	return
}

/*
CloseSettleChannel can do the following jobs:
close channel
settle channel
deposit to channel
*/
func CloseSettleChannel(w rest.ResponseWriter, r *rest.Request) {
	var err error
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> CloseSettleChannel ,err=%v", err))
	}()
	chstr := r.PathParam("channel")
	if len(chstr) != len(utils.EmptyHash.String()) {
		rest.Error(w, "argument error", http.StatusBadRequest)
		return
	}
	channelIdentifier := common.HexToHash(chstr)
	type Req struct {
		State    string
		StateInt channeltype.State
		Balance  *big.Int
		Force    bool
	}
	req := &Req{}
	err = r.DecodeJsonPayload(req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c, err := API.GetChannel(channelIdentifier)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusConflict)
		return
	}

	if req.State == "closed" {
		req.StateInt = channeltype.StateClosed
	} else if req.State == "settled" {
		req.StateInt = channeltype.StateSettled
	} else {
		req.StateInt = channeltype.StateError
	}
	//close or settle
	if req.StateInt != channeltype.StateClosed && req.StateInt != channeltype.StateSettled {
		rest.Error(w, "argument error", http.StatusBadRequest)
		return
	}
	if req.StateInt == channeltype.StateClosed {
		if req.Force {
			c, err = API.Close(c.TokenAddress(), c.PartnerAddress())
			if err != nil {
				log.Error(err.Error())
				rest.Error(w, err.Error(), http.StatusConflict)
				return
			}
		} else {
			//cooperative settle channel
			c, err = API.CooperativeSettle(c.TokenAddress(), c.PartnerAddress())
			if err != nil {
				log.Error(err.Error())
				rest.Error(w, err.Error(), http.StatusConflict)
				return
			}
		}

	} else if req.StateInt == channeltype.StateSettled {
		c, err = API.Settle(c.TokenAddress(), c.PartnerAddress())
		if err != nil {
			log.Error(err.Error())
			rest.Error(w, err.Error(), http.StatusConflict)
			return
		}
	}
	d := &ChannelData{
		ChannelIdentifier:   c.ChannelIdentifier.ChannelIdentifier.String(),
		OpenBlockNumber:     c.ChannelIdentifier.OpenBlockNumber,
		PartnerAddrses:      c.PartnerAddress().String(),
		Balance:             c.OurBalance(),
		PartnerBalance:      c.PartnerBalance(),
		State:               c.State,
		StateString:         c.State.String(),
		SettleTimeout:       c.SettleTimeout,
		TokenAddress:        c.TokenAddress().String(),
		LockedAmount:        c.OurAmountLocked(),
		PartnerLockedAmount: c.PartnerAmountLocked(),
		RevealTimeout:       c.RevealTimeout,
	}
	err = w.WriteJson(d)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

func withdraw(w rest.ResponseWriter, r *rest.Request) {
	chstr := r.PathParam("channel")
	if len(chstr) != len(utils.EmptyHash.String()) {
		rest.Error(w, "argument error", http.StatusBadRequest)
		return
	}
	channelIdentifier := common.HexToHash(chstr)
	type Req struct {
		Amount *big.Int
		Op     string
	}
	const OpPrepareWithdraw = "preparewithdraw"
	const OpCancelPrepare = "cancelprepare"

	req := &Req{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c, err := API.GetChannel(channelIdentifier)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if req.Amount != nil && req.Amount.Cmp(utils.BigInt0) > 0 { //deposit
		c, err = API.Withdraw(c.TokenAddress(), c.PartnerAddress(), req.Amount)
		if err != nil {
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		if req.Op == OpPrepareWithdraw {
			c, err = API.PrepareForWithdraw(c.TokenAddress(), c.PartnerAddress())
		} else if req.Op == OpCancelPrepare {
			c, err = API.CancelPrepareForWithdraw(c.TokenAddress(), c.PartnerAddress())
		} else {
			err = fmt.Errorf("unkown operation %s", req.Op)
		}
		if err != nil {
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	d := &ChannelData{
		ChannelIdentifier:   c.ChannelIdentifier.ChannelIdentifier.String(),
		OpenBlockNumber:     c.ChannelIdentifier.OpenBlockNumber,
		PartnerAddrses:      c.PartnerAddress().String(),
		Balance:             c.OurBalance(),
		PartnerBalance:      c.PartnerBalance(),
		State:               c.State,
		StateString:         c.State.String(),
		SettleTimeout:       c.SettleTimeout,
		TokenAddress:        c.TokenAddress().String(),
		LockedAmount:        c.OurAmountLocked(),
		PartnerLockedAmount: c.PartnerAmountLocked(),
		RevealTimeout:       c.RevealTimeout,
	}
	err = w.WriteJson(d)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

//BalanceUpdateForPFS for path finding service, test only
func BalanceUpdateForPFS(w rest.ResponseWriter, r *rest.Request) {
	ch := r.PathParam("channel")
	channelIdentifier := common.HexToHash(ch)
	if channelIdentifier == utils.EmptyHash {
		rest.Error(w, "argument error", http.StatusBadRequest)
		return
	}
	result, err := API.BalanceProofForPFS(channelIdentifier)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = w.WriteJson(result)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}
