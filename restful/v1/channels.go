package v1

import (
	"net/http"

	"math/big"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/transfer"
	"github.com/SmartMeshFoundation/Photon/transfer/mtree"
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

//ChannelDataDetail more info
type ChannelDataDetail struct {
	ChannelIdentifier   string            `json:"channel_identifier"`
	OpenBlockNumber     int64             `json:"open_block_number"`
	PartnerAddress      string            `json:"partner_address"`
	Balance             *big.Int          `json:"balance"`
	PartnerBalance      *big.Int          `json:"partner_balance"`
	LockedAmount        *big.Int          `json:"locked_amount"`
	PartnerLockedAmount *big.Int          `json:"partner_locked_amount"`
	TokenAddress        string            `json:"token_address"`
	State               channeltype.State `json:"state"`
	StateString         string
	SettleTimeout       int `json:"settle_timeout"`
	RevealTimeout       int `json:"reveal_timeout"`

	/*
		extended
	*/
	ClosedBlock              int64
	SettledBlock             int64
	OurUnkownSecretLocks     map[common.Hash]channeltype.PendingLock
	OurKnownSecretLocks      map[common.Hash]channeltype.UnlockPartialProof
	PartnerUnkownSecretLocks map[common.Hash]channeltype.PendingLock
	PartnerKnownSecretLocks  map[common.Hash]channeltype.UnlockPartialProof
	OurLeaves                []*mtree.Lock
	PartnerLeaves            []*mtree.Lock
	OurBalanceProof          *transfer.BalanceProofState
	PartnerBalanceProof      *transfer.BalanceProofState
	Signature                []byte //my signature of PartnerBalanceProof
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
	d := &ChannelDataDetail{
		ChannelIdentifier:        c.ChannelIdentifier.ChannelIdentifier.String(),
		OpenBlockNumber:          c.ChannelIdentifier.OpenBlockNumber,
		PartnerAddress:           c.PartnerAddress().String(),
		Balance:                  c.OurBalance(),
		PartnerBalance:           c.PartnerBalance(),
		State:                    c.State,
		StateString:              c.State.String(),
		SettleTimeout:            c.SettleTimeout,
		TokenAddress:             c.TokenAddress().String(),
		LockedAmount:             c.OurAmountLocked(),
		PartnerLockedAmount:      c.PartnerAmountLocked(),
		ClosedBlock:              c.ClosedBlock,
		SettledBlock:             c.SettledBlock,
		OurLeaves:                c.OurLeaves,
		PartnerLeaves:            c.PartnerLeaves,
		OurKnownSecretLocks:      c.OurLock2UnclaimedLocks(),
		OurUnkownSecretLocks:     c.OurLock2PendingLocks(),
		PartnerUnkownSecretLocks: c.PartnerLock2PendingLocks(),
		PartnerKnownSecretLocks:  c.PartnerLock2UnclaimedLocks(),
		OurBalanceProof:          c.OurBalanceProof,
		PartnerBalanceProof:      c.PartnerBalanceProof,
	}
	err = w.WriteJson(d)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

/*
Deposit open a channel with partner.
token must exist
partner maybe an invalid address
*/
func Deposit(w rest.ResponseWriter, r *rest.Request) {
	var err error
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> Deposit ,err=%v", err))
	}()
	req := &ChannelData{}
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
	if req.State == 0 { //open channel
		c, err := API.DepositAndOpenChannel(tokenAddr, partnerAddr, req.SettleTimeout, API.Photon.Config.RevealTimeout, req.Balance)
		if err != nil {
			log.Error(err.Error())
			rest.Error(w, err.Error(), http.StatusConflict)
			return
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
	var err error
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> CloseSettleDepositChannel ,err=%v", err))
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
