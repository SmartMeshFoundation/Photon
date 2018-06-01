package mobile

import (
	"encoding/json"
	"math/rand"

	"fmt"
	"time"

	"math/big"

	"errors"

	"github.com/SmartMeshFoundation/SmartRaiden"
	"github.com/SmartMeshFoundation/SmartRaiden/internal/rpanic"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/restful/v1"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

//API for export interface
type API struct {
	api *smartraiden.RaidenAPI
}

func marshal(v interface{}) (s string, err error) {
	d, err := json.Marshal(v)
	if err != nil {
		log.Error(err.Error())
		return
	}
	return string(d), nil
}

//GetChannelList GET /api/1/channels
func (a *API) GetChannelList() (channels string, err error) {
	chs, err := a.api.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	if err != nil {
		log.Error(err.Error())
		return
	}
	var datas []*v1.ChannelData
	for _, c := range chs {
		d := &v1.ChannelData{
			ChannelAddress:      c.ChannelAddress.String(),
			PartnerAddrses:      c.PartnerAddress.String(),
			Balance:             c.OurBalance,
			PartnerBalance:      c.PartnerBalance,
			LockedAmount:        c.OurAmountLocked,
			PartnerLockedAmount: c.PartnerAmountLocked,
			State:               c.State,
			TokenAddress:        c.TokenAddress.String(),
			SettleTimeout:       c.SettleTimeout,
			RevealTimeout:       c.RevealTimeout,
		}
		datas = append(datas, d)
	}
	channels, err = marshal(datas)
	return
}

//GetOneChannel GET /api/1/channels/0x2a65aca4d5fc5b5c859090a6c34d164135398226
func (a *API) GetOneChannel(channelAddress string) (channel string, err error) {
	chaddr := common.HexToAddress(channelAddress)
	c, err := a.api.GetChannel(chaddr)
	if err != nil {
		log.Error(err.Error())
		return
	}
	d := &v1.ChannelDataDetail{
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
	channel, err = marshal(d)
	return
}

//OpenChannel put request
func (a *API) OpenChannel(partnerAddress, tokenAddress string, settleTimeout int, balanceStr string) (channel string, err error) {
	partnerAddr := common.HexToAddress(partnerAddress)
	tokenAddr := common.HexToAddress(tokenAddress)
	balance, _ := new(big.Int).SetString(balanceStr, 0)
	c, err := a.api.Open(tokenAddr, partnerAddr, settleTimeout, params.DefaultRevealTimeout)
	if err != nil {
		log.Error(err.Error())
		return
	}
	d := &v1.ChannelData{
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
	if balance.Cmp(utils.BigInt0) > 0 {
		err = a.api.Deposit(tokenAddr, partnerAddr, balance, params.DefaultPollTimeout)
		if err == nil {
			d.Balance = c.OurBalance
		} else {
			log.Error(" RaidenAPI.Deposit error : ", err)
			return
		}
	}
	channel, err = marshal(d)
	return

}

//CloseChannel close a channel
func (a *API) CloseChannel(channelAddres string) (channel string, err error) {
	chAddr := common.HexToAddress(channelAddres)
	c, err := a.api.GetChannel(chAddr)
	if err != nil {
		log.Error(err.Error())
		return
	}
	c, err = a.api.Close(c.TokenAddress, c.PartnerAddress)
	if err != nil {
		log.Error(err.Error())
		return
	}
	d := &v1.ChannelData{
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
	channel, err = marshal(d)
	return
}

//SettleChannel settle a channel
func (a *API) SettleChannel(channelAddres string) (channel string, err error) {
	chAddr := common.HexToAddress(channelAddres)
	c, err := a.api.GetChannel(chAddr)
	if err != nil {
		log.Error(err.Error())
		return
	}
	c, err = a.api.Settle(c.TokenAddress, c.PartnerAddress)
	if err != nil {
		log.Error(err.Error())
		return
	}
	d := &v1.ChannelData{
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
	channel, err = marshal(d)
	return
}

//DepositChannel deposit balance to channel
func (a *API) DepositChannel(channelAddres string, balanceStr string) (channel string, err error) {
	chAddr := common.HexToAddress(channelAddres)
	balance, _ := new(big.Int).SetString(balanceStr, 0)
	c, err := a.api.GetChannel(chAddr)
	if err != nil {
		log.Error(fmt.Sprintf("GetChannel %s err %s", utils.APex(chAddr), err))
		return
	}
	err = a.api.Deposit(c.TokenAddress, c.PartnerAddress, balance, params.DefaultPollTimeout)
	if err != nil {
		log.Error(fmt.Sprintf("Deposit to %s:%s err %s", utils.APex(c.TokenAddress),
			utils.APex(c.PartnerAddress), err))
		return
	}

	d := &v1.ChannelData{
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
	channel, err = marshal(d)
	return
}

//NetworkEvent GET /api/<version>/events/network
func (a *API) NetworkEvent(fromBlock, toBlock int64) (eventsString string, err error) {
	events, err := a.api.GetNetworkEvents(fromBlock, toBlock)
	if err != nil {
		log.Error(err.Error())
		return
	}
	eventsString, err = marshal(events)
	return
}

//TokensEvent GET /api/1/events/tokens/0x61c808d82a3ac53231750dadc13c777b59310bd9
func (a *API) TokensEvent(fromBlock, toBlock int64, tokenAddress string) (eventsString string, err error) {
	token := common.HexToAddress(tokenAddress)
	events, err := a.api.GetTokenNetworkEvents(token, fromBlock, toBlock)
	if err != nil {
		log.Error(err.Error())
		return
	}
	eventsString, err = marshal(events)
	return
}

//ChannelsEvent GET /api/1/events/channels/0x2a65aca4d5fc5b5c859090a6c34d164135398226?from_block=1337
func (a *API) ChannelsEvent(fromBlock, toBlock int64, channelAddress string) (eventsString string, err error) {
	channel := common.HexToAddress(channelAddress)
	events, err := a.api.GetChannelEvents(channel, fromBlock, toBlock)
	if err != nil {
		log.Error(err.Error())
		return
	}
	eventsString, err = marshal(events)
	return
}

//Address GET /api/1/address
func (a *API) Address() (addr string) {
	return a.api.Address().String()
}

//Tokens GET /api/1/tokens
func (a *API) Tokens() (tokens string) {
	tokens, err := marshal(a.api.Tokens())
	if err != nil {
		log.Error(fmt.Sprintf("marshal tokens error %s", err))
	}
	return
}

type partnersData struct {
	PartnerAddress string `json:"partner_address"`
	Channel        string `json:"channel"`
}

//TokenPartners GET /api/1/tokens/0x61bb630d3b2e8eda0fc1d50f9f958ec02e3969f6/partners
func (a *API) TokenPartners(tokenAddress string) (channels string, err error) {
	tokenAddr := common.HexToAddress(tokenAddress)
	chs, err := a.api.GetChannelList(tokenAddr, utils.EmptyAddress)
	if err != nil {
		log.Error(err.Error())
		return
	}
	var datas []*partnersData
	for _, c := range chs {
		d := &partnersData{
			PartnerAddress: c.PartnerAddress.String(),
			Channel:        "api/1/channles/" + c.OurAddress.String(),
		}
		datas = append(datas, d)
	}
	channels, err = marshal(datas)
	return
}

//RegisterToken PUT /api/1/tokens/0xea674fdde714fd979de3edf0f56aa9716b898ec8 Registering a Token
func (a *API) RegisterToken(tokenAddress string) (managerAddress string, err error) {
	tokenAddr := common.HexToAddress(tokenAddress)
	mgr, err := a.api.RegisterToken(tokenAddr)
	if err != nil {
		log.Error(err.Error())
		return
	}
	return mgr.String(), err
}

/*
Transfers POST /api/1/transfers/0x2a65aca4d5fc5b5c859090a6c34d164135398226/0x61c808d82a3ac53231750dadc13c777b59310bd9
Initiating a Transfer
identifier:0 means random identifier generated by system
*/
func (a *API) Transfers(tokenAddress, targetAddress string, amountstr string, feestr string, id int64) (transfer string, err error) {
	tokenAddr := common.HexToAddress(tokenAddress)
	targetAddr := common.HexToAddress(targetAddress)
	amount, _ := new(big.Int).SetString(amountstr, 0)
	fee, _ := new(big.Int).SetString(feestr, 0)
	identifier := uint64(id)
	if identifier == 0 {
		identifier = rand.New(rand.NewSource(time.Now().UnixNano())).Uint64()
	}
	if amount.Cmp(utils.BigInt0) <= 0 {
		err = errors.New("amount should be positive")
		return
	}
	err = a.api.Transfer(tokenAddr, amount, fee, targetAddr, identifier, params.MaxRequestTimeout)
	if err != nil {
		log.Error(err.Error())
		return
	}
	req := &v1.TransferData{}
	req.Initiator = a.api.Raiden.NodeAddress.String()
	req.Target = targetAddress
	req.Token = tokenAddress
	req.Amount = amount
	req.Identifier = identifier
	req.Fee = fee
	return marshal(req)
}

/*
TokenSwap token swap for maker
role: "maker" or "taker"
*/
func (a *API) TokenSwap(role string, Identifier int64, SendingAmountStr, ReceivingAmountStr string, SendingToken, ReceivingToken, TargetAddress string) (err error) {
	type Req struct {
		Role            string   `json:"role"`
		SendingAmount   *big.Int `json:"sending_amount"`
		SendingToken    string   `json:"sending_token"`
		ReceivingAmount int64    `json:"receiving_amount"`
		ReceivingToken  *big.Int `json:"receiving_token"`
	}

	var target common.Address
	target = common.HexToAddress(TargetAddress)
	if Identifier <= 0 {
		err = errors.New("Identifier must be positive")
		return
	}
	SendingAmount, _ := new(big.Int).SetString(SendingAmountStr, 0)
	ReceivingAmount, _ := new(big.Int).SetString(ReceivingAmountStr, 0)
	if role == "maker" {
		err = a.api.TokenSwapAndWait(uint64(Identifier), common.HexToAddress(SendingToken), common.HexToAddress(ReceivingToken),
			a.api.Raiden.NodeAddress, target, SendingAmount, ReceivingAmount)
	} else if role == "taker" {
		err = a.api.ExpectTokenSwap(uint64(Identifier), common.HexToAddress(ReceivingToken), common.HexToAddress(SendingToken),
			target, a.api.Raiden.NodeAddress, ReceivingAmount, SendingAmount)
	} else {
		err = fmt.Errorf("Provided invalid token swap role %s", role)
	}
	return
}

//Stop stop raiden
func (a *API) Stop() {
	//test only
	a.api.Stop()
}

/*
QueryingConnections Querying connections details

GET /api/<version>/connections
*/
func (a *API) QueryingConnections() string {
	connections := a.api.GetConnectionManagersInfo()
	s, err := marshal(connections)
	if err != nil {
		log.Error(fmt.Sprintf("marshal connections error %s", err))
	}
	return s
}

/*
ConnectToTokenNetwork Connecting to a token network
PUT /api/1/connections/0x2a65aca4d5fc5b5c859090a6c34d164135398226
*/
//Connecting to a token network
func (a *API) ConnectToTokenNetwork(tokenAddress string, fundsStr string) (err error) {
	token := common.HexToAddress(tokenAddress)
	funds, _ := new(big.Int).SetString(fundsStr, 0)
	if funds.Cmp(utils.BigInt0) <= 0 {
		err = errors.New("funds <=0")
		return
	}
	err = a.api.ConnectTokenNetwork(token, funds, params.DefaultInitialChannelTarget, params.DefaultJoinableFundsTarget)
	return
}

/*
LeaveTokenNetwork leave a token network
*/
func (a *API) LeaveTokenNetwork(OnlyReceivingChannels bool, tokenAddress string) (channels string, err error) {
	token := common.HexToAddress(tokenAddress)

	chs, err := a.api.LeaveTokenNetwork(token, OnlyReceivingChannels)
	if err != nil {
		log.Error(err.Error())
		return
	}

	var addrs []string
	for _, c := range chs {
		addrs = append(addrs, c.OurAddress.String())
	}
	channels, err = marshal(addrs)
	return

}

/*
ChannelFor3rdParty generate info for 3rd party use,
for update transfer and withdraw.
*/
func (a *API) ChannelFor3rdParty(channelAddress, thirdPartyAddress string) (r string, err error) {
	channelAddr := common.HexToAddress(channelAddress)
	thirdPartyAddr := common.HexToAddress(thirdPartyAddress)
	if channelAddr == utils.EmptyAddress || thirdPartyAddr == utils.EmptyAddress {
		err = errors.New("invalid argument")
		return
	}
	result, err := a.api.ChannelInformationFor3rdParty(channelAddr, thirdPartyAddr)
	if err != nil {
		log.Error(err.Error())
		return
	}
	r, err = marshal(result)
	return
}

/*
SwitchToMesh  makes raiden switch to mesh mode,if nodestr is empty ,will switch back to normal mode.
*/
func (a *API) SwitchToMesh(nodesstr string) (err error) {
	var nodes []*network.NodeInfo
	err = json.Unmarshal([]byte(nodesstr), &nodes)
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = a.api.Raiden.Protocol.UpdateMeshNetworkNodes(nodes)
	if err != nil {
		log.Error(err.Error())
		return
	}
	return nil
}

/*
EthereumStatus  query the status between raiden and ethereum
todo fix it ,r is useless
*/
func (a *API) EthereumStatus() (r string, err error) {
	c := a.api.Raiden.Chain
	if c != nil && c.Client.Status == helper.ConnectionOk {
		return time.Now().String(), nil
	}
	return time.Now().String(), errors.New("connect failed")
}

/*
GetSentTransfers retuns list of sent transfer between `from_block` and `to_block`
*/
func (a *API) GetSentTransfers(from, to int64) (r string, err error) {
	log.Trace(fmt.Sprintf("from=%d,to=%d\n", from, to))
	trs, err := a.api.GetSentTransfers(from, to)
	if err != nil {
		log.Error(err.Error())
		return
	}
	r, err = marshal(trs)
	return
}

/*
GetReceivedTransfers retuns list of received transfer between `from_block` and `to_block`
it contains token swap
*/
func (a *API) GetReceivedTransfers(from, to int64) (r string, err error) {
	trs, err := a.api.GetReceivedTransfers(from, to)
	if err != nil {
		log.Error(err.Error())
		return
	}
	r, err = marshal(trs)
	return
}

// ErrorHandler is a client-side subscription callback to invoke on events and
// subscription failure.
type ErrorHandler interface {
	OnError(errCode int, failure string)
}

// SubscribeError subscribes to notifications about the current blockchain head
// on the given channel.
func (a *API) SubscribeError(handler ErrorHandler) (err error) {
	// Subscribe to the event internally
	go func() {
		rpanic.RegisterErrorNotifier("API SubscribeError")
		err = <-rpanic.GetNotify()
		handler.OnError(32, err.Error())
	}()
	return nil
}
