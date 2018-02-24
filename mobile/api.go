package mobile

import (
	"encoding/json"
	"math/rand"

	"fmt"
	"time"

	"math/big"

	"errors"

	"github.com/SmartMeshFoundation/raiden-network"
	"github.com/SmartMeshFoundation/raiden-network/params"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type Api struct {
	api *raiden_network.RaidenApi
}

func marshal(v interface{}) (s string, err error) {
	d, err := json.Marshal(v)
	if err != nil {
		return
	}
	return string(d), nil
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

//GET /api/1/channels
func (this *Api) GetChannelList() (channels string, err error) {
	chs, _ := this.api.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	var datas []*channelData
	for _, c := range chs {
		d := &channelData{
			ChannelAddress: c.OurAddress.String(),
			PartnerAddrses: c.PartnerAddress.String(),
			Balance:        c.OurBalance,
			State:          c.State,
			TokenAddress:   c.TokenAddress.String(),
			SettleTimeout:  c.SettleTimeout,
			RevealTimeout:  c.RevealTimeout,
		}
		datas = append(datas, d)
	}
	channels, err = marshal(datas)
	return
}

//GET /api/1/channels/0x2a65aca4d5fc5b5c859090a6c34d164135398226
func (this *Api) GetOneChannel(channelAddress string) (channel string, err error) {
	chaddr := common.HexToAddress(channelAddress)
	c, err := this.api.GetChannel(chaddr)
	if err != nil {
		return
	}
	d := &channelData{
		ChannelAddress: c.OurAddress.String(),
		PartnerAddrses: c.PartnerAddress.String(),
		Balance:        c.OurBalance,
		State:          c.State,
		SettleTimeout:  c.SettleTimeout,
		TokenAddress:   c.TokenAddress.String(),
	}
	channel, err = marshal(d)
	return
}

//put request
func (this *Api) OpenChannel(partnerAddress, tokenAddress string, settleTimeout int, balanceStr string) (channel string, err error) {
	partnerAddr := common.HexToAddress(partnerAddress)
	tokenAddr := common.HexToAddress(tokenAddress)
	balance, _ := new(big.Int).SetString(balanceStr, 0)
	c, err := this.api.Open(tokenAddr, partnerAddr, settleTimeout, params.DEFAULT_REVEAL_TIMEOUT)
	if err != nil {
		return
	} else {
		d := &channelData{
			ChannelAddress: c.OurAddress.String(),
			PartnerAddrses: c.PartnerAddress.String(),
			Balance:        c.OurBalance,
			State:          c.State,
			SettleTimeout:  c.SettleTimeout,
			TokenAddress:   c.TokenAddress.String(),
		}
		if balance.Cmp(utils.BigInt0) > 0 {
			err = this.api.Deposit(tokenAddr, partnerAddr, balance, params.DEFAULT_POLL_TIMEOUT)
			if err == nil {
				d.Balance = c.OurBalance
			} else {
				log.Error(" RaidenApi.Deposit error : ", err)
				return
			}
		}
		channel, err = marshal(d)
		return
	}
	return
}

func (this *Api) CloseChannel(channelAddres string) (channel string, err error) {
	chAddr := common.HexToAddress(channelAddres)
	c, err := this.api.GetChannel(chAddr)
	if err != nil {
		log.Error(err.Error())
		return
	}
	c, err = this.api.Close(c.TokenAddress, c.PartnerAddress)
	if err != nil {
		log.Error(err.Error())
		return
	}
	d := &channelData{
		ChannelAddress: c.OurAddress.String(),
		PartnerAddrses: c.PartnerAddress.String(),
		Balance:        c.OurBalance,
		State:          c.State,
		SettleTimeout:  c.SettleTimeout,
		TokenAddress:   c.TokenAddress.String(),
	}
	channel, err = marshal(d)
	return
}
func (this *Api) SettleChannel(channelAddres string) (channel string, err error) {
	chAddr := common.HexToAddress(channelAddres)
	c, err := this.api.GetChannel(chAddr)
	if err != nil {
		log.Error(err.Error())
		return
	}
	c, err = this.api.Settle(c.TokenAddress, c.PartnerAddress)
	if err != nil {
		log.Error(err.Error())
		return
	}
	d := &channelData{
		ChannelAddress: c.OurAddress.String(),
		PartnerAddrses: c.PartnerAddress.String(),
		Balance:        c.OurBalance,
		State:          c.State,
		SettleTimeout:  c.SettleTimeout,
		TokenAddress:   c.TokenAddress.String(),
	}
	channel, err = marshal(d)
	return
}
func (this *Api) DepositChannel(channelAddres string, balanceStr string) (channel string, err error) {
	chAddr := common.HexToAddress(channelAddres)
	balance, _ := new(big.Int).SetString(balanceStr, 0)
	c, err := this.api.GetChannel(chAddr)
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = this.api.Deposit(c.TokenAddress, c.PartnerAddress, balance, params.DEFAULT_POLL_TIMEOUT)
	if err != nil {
		return
	}

	d := &channelData{
		ChannelAddress: c.OurAddress.String(),
		PartnerAddrses: c.PartnerAddress.String(),
		Balance:        c.OurBalance,
		State:          c.State,
		SettleTimeout:  c.SettleTimeout,
		TokenAddress:   c.TokenAddress.String(),
	}
	channel, err = marshal(d)
	return
}

//GET /api/<version>/events/network
func (this *Api) NetworkEvent(fromBlock, toBlock int64) (eventsString string, err error) {
	events, err := this.api.GetNetworkEvents(fromBlock, toBlock)
	if err != nil {
		return
	}
	eventsString, err = marshal(events)
	return
}

//GET /api/1/events/tokens/0x61c808d82a3ac53231750dadc13c777b59310bd9
func (this *Api) TokensEvent(fromBlock, toBlock int64, tokenAddress string) (eventsString string, err error) {
	token := common.HexToAddress(tokenAddress)
	events, err := this.api.GetTokenNetworkEvents(token, fromBlock, toBlock)
	if err != nil {
		return
	}
	eventsString, err = marshal(events)
	return
}

//GET /api/1/events/channels/0x2a65aca4d5fc5b5c859090a6c34d164135398226?from_block=1337
func (this *Api) ChannelsEvent(fromBlock, toBlock int64, channelAddress string) (eventsString string, err error) {
	channel := common.HexToAddress(channelAddress)
	events, err := this.api.GetChannelEvents(channel, fromBlock, toBlock)
	if err != nil {
		return
	}
	eventsString, err = marshal(events)
	return
}

//GET /api/1/address
func (this *Api) Address() (addr string) {
	return this.api.Address().String()
}

//GET /api/1/tokens
func (this *Api) Tokens() (tokens string) {
	tokens, _ = marshal(this.api.Tokens())
	return
}

type partnersData struct {
	PartnerAddress string `json:"partner_address"`
	Channel        string `json:"channel"`
}

//GET /api/1/tokens/0x61bb630d3b2e8eda0fc1d50f9f958ec02e3969f6/partners
func (this *Api) TokenPartners(tokenAddress string) (channels string, err error) {
	tokenAddr := common.HexToAddress(tokenAddress)
	chs, _ := this.api.GetChannelList(tokenAddr, utils.EmptyAddress)
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

//PUT /api/1/tokens/0xea674fdde714fd979de3edf0f56aa9716b898ec8 Registering a Token
func (this *Api) RegisterToken(tokenAddress string) (managerAddress string, err error) {
	tokenAddr := common.HexToAddress(tokenAddress)
	mgr, err := this.api.RegisterToken(tokenAddr)
	if err != nil {
		return
	} else {
		return mgr.String(), err
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
}

/*
POST /api/1/transfers/0x2a65aca4d5fc5b5c859090a6c34d164135398226/0x61c808d82a3ac53231750dadc13c777b59310bd9
Initiating a Transfer
identifier:0 means random identifier generated by system
*/
func (this *Api) Transfers(tokenAddress, targetAddress string, amountstr string, id int64) (transfer string, err error) {
	tokenAddr := common.HexToAddress(tokenAddress)
	targetAddr := common.HexToAddress(targetAddress)
	amount, _ := new(big.Int).SetString(amountstr, 0)
	identifier := uint64(id)
	if identifier == 0 {
		identifier = rand.New(rand.NewSource(time.Now().UnixNano())).Uint64()
	}
	if amount.Cmp(utils.BigInt0) <= 0 {
		err = errors.New("amount should be positive")
		return
	}
	err = this.api.Transfer(tokenAddr, amount, targetAddr, identifier, params.MaxRequestTimeout)
	if err != nil {
		return
	}
	req := &transferData{}
	req.Initiator = this.api.Raiden.NodeAddress.String()
	req.Target = targetAddress
	req.Token = tokenAddress
	req.Amount = amount
	req.Identifier = identifier
	return marshal(req)
}

/*
token swap for maker
role: "maker" or "taker"
*/
func (this *Api) TokenSwap(role string, Identifier int64, SendingAmountStr, ReceivingAmountStr string, SendingToken, ReceivingToken, TargetAddress string) (err error) {
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
		err = this.api.TokenSwapAndWait(uint64(Identifier), common.HexToAddress(SendingToken), common.HexToAddress(ReceivingToken),
			this.api.Raiden.NodeAddress, target, SendingAmount, ReceivingAmount)
	} else if role == "taker" {
		err = this.api.ExpectTokenSwap(uint64(Identifier), common.HexToAddress(ReceivingToken), common.HexToAddress(SendingToken),
			target, this.api.Raiden.NodeAddress, ReceivingAmount, SendingAmount)
	} else {
		err = fmt.Errorf("Provided invalid token swap role %s", role)
	}
	return
}

func (this *Api) Stop() {
	//test only
	this.api.Stop()
}

/*
Querying connections details

GET /api/<version>/connections
*/
func (this *Api) QueryingConnections() string {
	connections := this.api.GetConnectionManagersInfo()
	s, _ := marshal(connections)
	return s
}

/*
Connecting to a token network
PUT /api/1/connections/0x2a65aca4d5fc5b5c859090a6c34d164135398226
*/
//Connecting to a token network
func (this *Api) ConnectToTokenNetwork(tokenAddress string, fundsStr string) (err error) {
	token := common.HexToAddress(tokenAddress)
	funds, _ := new(big.Int).SetString(fundsStr, 0)
	if funds.Cmp(utils.BigInt0) <= 0 {
		err = errors.New("funds <=0")
		return
	}
	err = this.api.ConnectTokenNetwork(token, funds, params.DEFAULT_INITIAL_CHANNEL_TARGET, params.DEFAULT_JOINABLE_FUNDS_TARGET)
	return
}

/*
leave a token network
*/
func (this *Api) LeaveTokenNetwork(OnlyReceivingChannels bool, tokenAddress string) (channels string, err error) {
	token := common.HexToAddress(tokenAddress)

	chs, err := this.api.LeaveTokenNetwork(token, OnlyReceivingChannels)
	if err != nil {
		return
	}

	var addrs []string
	for _, c := range chs {
		addrs = append(addrs, c.OurAddress.String())
	}
	channels, err = marshal(addrs)
	return

}
