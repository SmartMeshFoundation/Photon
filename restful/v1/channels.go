package v1

import (
	"net/http"

	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

/*
ChannelData export json data format
*/
type ChannelData struct {
	ChannelAddress      string   `json:"channel_address"`
	PartnerAddrses      string   `json:"partner_address"`
	Balance             *big.Int `json:"balance"`
	PartnerBalance      *big.Int `json:"partner_balance"`
	LockedAmount        *big.Int `json:"locked_amount"`
	PartnerLockedAmount *big.Int `json:"partner_locked_amount"`
	TokenAddress        string   `json:"token_address"`
	State               string   `json:"state"`
	SettleTimeout       int      `json:"settle_timeout"`
	RevealTimeout       int      `json:"reveal_timeout"`
}

type channelDataDetail struct {
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

	/*
		extended
	*/
	ClosedBlock              int64
	SettledBlock             int64
	OurUnkownSecretLocks     map[common.Hash]channel.PendingLock
	OurKnownSecretLocks      map[common.Hash]channel.UnlockPartialProof
	PartnerUnkownSecretLocks map[common.Hash]channel.PendingLock
	PartnerKnownSecretLocks  map[common.Hash]channel.UnlockPartialProof
	OurLeaves                []common.Hash
	PartnerLeaves            []common.Hash
	OurBalanceProof          *transfer.BalanceProofState
	PartnerBalanceProof      *transfer.BalanceProofState
	Signature                []byte //my signature of PartnerBalanceProof
}

/*
GetChannelList list all my channels
*/
func GetChannelList(w rest.ResponseWriter, r *rest.Request) {
	chs, err := RaidenAPI.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var datas []*ChannelData
	for _, c := range chs {
		d := &ChannelData{
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

/*
ChannelFor3rdParty generate info for 3rd party use,
for update transfer and withdraw.
*/
func ChannelFor3rdParty(w rest.ResponseWriter, r *rest.Request) {
	ch := r.PathParam("channel")
	thirdParty := r.PathParam("3rd")
	channelAddress := common.HexToAddress(ch)
	thirdAddress := common.HexToAddress(thirdParty)
	if channelAddress == utils.EmptyAddress || thirdAddress == utils.EmptyAddress {
		rest.Error(w, "argument error", http.StatusBadRequest)
		return
	}
	result, err := RaidenAPI.ChannelInformationFor3rdParty(channelAddress, thirdAddress)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(result)
}

/*
SpecifiedChannel get  a channel state
*/
func SpecifiedChannel(w rest.ResponseWriter, r *rest.Request) {
	ch := r.PathParam("channel")
	chaddr := common.HexToAddress(ch)
	c, err := RaidenAPI.GetChannel(chaddr)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	d := &channelDataDetail{
		ChannelAddress:           c.ChannelAddress.String(),
		PartnerAddrses:           c.PartnerAddress.String(),
		Balance:                  c.OurBalance,
		PartnerBalance:           c.PartnerBalance,
		State:                    c.State,
		SettleTimeout:            c.SettleTimeout,
		TokenAddress:             c.TokenAddress.String(),
		LockedAmount:             c.OurAmountLocked,
		PartnerLockedAmount:      c.PartnerAmountLocked,
		ClosedBlock:              c.ClosedBlock,
		SettledBlock:             c.SettledBlock,
		OurLeaves:                c.OurLeaves,
		PartnerLeaves:            c.PartnerLeaves,
		OurKnownSecretLocks:      c.OurLock2UnclaimedLocks,
		OurUnkownSecretLocks:     c.OurLock2PendingLocks,
		PartnerUnkownSecretLocks: c.PartnerLock2PendingLocks,
		PartnerKnownSecretLocks:  c.PartnerLock2UnclaimedLocks,
		OurBalanceProof:          c.OurBalanceProof,
		PartnerBalanceProof:      c.PartnerBalanceProof,
	}
	w.WriteJson(d)
}

/*
OpenChannel open a channel with partner.
token must exist
partner maybe an invalid address
*/
func OpenChannel(w rest.ResponseWriter, r *rest.Request) {
	req := &ChannelData{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	partnerAddr := common.HexToAddress(req.PartnerAddrses)
	tokenAddr := common.HexToAddress(req.TokenAddress)
	if req.State == "" { //open channel
		c, err := RaidenAPI.Open(tokenAddr, partnerAddr, req.SettleTimeout, params.DefaultRevealTimeout)
		if err != nil {
			log.Error(err.Error())
			rest.Error(w, err.Error(), http.StatusConflict)
			return
		}
		d := &ChannelData{
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
			err = RaidenAPI.Deposit(tokenAddr, partnerAddr, req.Balance, params.DefaultPollTimeout)
			if err == nil {
				c2, err2 := RaidenAPI.GetChannel(c.ChannelAddress)
				if err2 != nil {
					rest.Error(w, err2.Error(), http.StatusInternalServerError)
				}
				d.Balance = c2.OurBalance
			} else {
				log.Error(" RaidenAPI.Deposit error : ", err)
			}
		}
		w.WriteJson(d)
		return
	}
	rest.Error(w, "argument error", http.StatusBadRequest)
	return
}

/*
CloseSettleDepositChannel can do the following jobs:
close channel
settle channel
deposit to channel
*/
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
	c, err := RaidenAPI.GetChannel(chAddr)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if req.Balance != nil && req.Balance.Cmp(utils.BigInt0) > 0 { //deposit
		err = RaidenAPI.Deposit(c.TokenAddress, c.PartnerAddress, req.Balance, params.DefaultPollTimeout)
		if err != nil {
			rest.Error(w, err.Error(), http.StatusRequestTimeout)
			return
		}
	} else {
		//close or settle
		if req.State != transfer.ChannelStateClosed && req.State != transfer.ChannelStateSettled {
			rest.Error(w, "argument error", http.StatusBadRequest)
			return
		}
		if req.State == transfer.ChannelStateClosed {
			c, err = RaidenAPI.Close(c.TokenAddress, c.PartnerAddress)
			if err != nil {
				log.Error(err.Error())
				rest.Error(w, err.Error(), http.StatusConflict)
				return
			}
		} else {
			c, err = RaidenAPI.Settle(c.TokenAddress, c.PartnerAddress)
			if err != nil {
				log.Error(err.Error())
				rest.Error(w, err.Error(), http.StatusConflict)
				return
			}
		}
	}
	//reload new data from database
	c, _ = RaidenAPI.GetChannel(c.ChannelAddress)
	d := &ChannelData{
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
