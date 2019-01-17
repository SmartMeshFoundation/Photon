package mobile

import (
	"encoding/json"
	"fmt"
	"time"

	"math/big"

	"errors"

	"strings"

	photon "github.com/SmartMeshFoundation/Photon"
	"github.com/SmartMeshFoundation/Photon/internal/rpanic"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/network"
	"github.com/SmartMeshFoundation/Photon/network/netshare"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/restful/v1"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

// API for export interface
//
// should not export any member because of gomobile's protocol
type API struct {
	api           *photon.API
	callID2result map[string]Result
}

// Result :
type Result struct {
	Result string `json:"result"`
	Err    error  `json:"error"`
	Done   bool   `json:"done"`
}

func newResult() Result {
	return Result{
		Result: "",
		Err:    nil,
		Done:   false,
	}
}

func marshal(v interface{}) (s string, err error) {
	d, err := json.Marshal(v)
	if err != nil {
		log.Error(err.Error())
		return
	}
	return string(d), nil
}

/*
GetChannelList returns all the available channels

example returns:
[
    {
        "channel_address": "0xc502076485a3cff65f83c00095dc55e745f790eee4c259ea963969a343fc792a",
        "open_block_number": 5228715,
        "partner_address": "0x4B89Bff01009928784eB7e7d10Bf773e6D166066",
        "balance": 499490,
        "partner_balance": 1500506,
        "locked_amount": 0,
        "partner_locked_amount": 0,
        "token_address": "0x663495a1b8e9Be17083b37924cFE39e17858F9e8",
        "state": 1,
        "StateString": "opened",
        "settle_timeout": 100000,
        "reveal_timeout": 5000
    }
]
*/
func (a *API) GetChannelList() (channels string, err error) {
	defer func() {
		log.Trace(fmt.Sprintf("ApiCall GetChannelList channels=\n%s,err=%v", channels, err))
	}()
	chs, err := a.api.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	if err != nil {
		log.Error(err.Error())
		return
	}
	var datas []*v1.ChannelData
	for _, c := range chs {
		d := &v1.ChannelData{
			ChannelIdentifier:   common.BytesToHash(c.Key).String(),
			PartnerAddrses:      c.PartnerAddress().String(),
			Balance:             c.OurBalance(),
			PartnerBalance:      c.PartnerBalance(),
			LockedAmount:        c.OurAmountLocked(),
			PartnerLockedAmount: c.PartnerAmountLocked(),
			State:               c.State,
			StateString:         c.State.String(),
			OpenBlockNumber:     c.ChannelIdentifier.OpenBlockNumber,
			TokenAddress:        c.TokenAddress().String(),
			SettleTimeout:       c.SettleTimeout,
			RevealTimeout:       c.RevealTimeout,
		}
		datas = append(datas, d)
	}
	channels, err = marshal(datas)
	return
}

/*
GetOneChannel return one specified channel with more detail information

exmaple returns:
{
    "channel_identifier": "0xc502076485a3cff65f83c00095dc55e745f790eee4c259ea963969a343fc792a",
    "open_block_number": 5228715,
    "partner_address": "0x4B89Bff01009928784eB7e7d10Bf773e6D166066",
    "balance": 499490,
    "patner_balance": 1500506,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x663495a1b8e9Be17083b37924cFE39e17858F9e8",
    "state": 1,
    "StateString": "opened",
    "settle_timeout": 100000,
    "reveal_timeout": 0,
    "ClosedBlock": 0,
    "SettledBlock": 0,
    "OurUnkownSecretLocks": {},
    "OurKnownSecretLocks": {},
    "PartnerUnkownSecretLocks": {},
    "PartnerKnownSecretLocks": {},
    "OurLeaves": null,
    "PartnerLeaves": null,
    "OurBalanceProof": {
        "Nonce": 0,
        "TransferAmount": 0,
        "LocksRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "ChannelIdentifier": {
            "ChannelIdentifier": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "OpenBlockNumber": 0
        },
        "MessageHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "Signature": null,
        "ContractTransferAmount": 0,
        "ContractNonce": 0,
        "ContractLocksRoot": "0x0000000000000000000000000000000000000000000000000000000000000000"
    },
    "PartnerBalanceProof": {
        "Nonce": 0,
        "TransferAmount": 0,
        "LocksRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "ChannelIdentifier": {
            "ChannelIdentifier": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "OpenBlockNumber": 0
        },
        "MessageHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "Signature": null,
        "ContractTransferAmount": 0,
        "ContractNonce": 0,
        "ContractLocksRoot": "0x0000000000000000000000000000000000000000000000000000000000000000"
    },
    "Signature": null
}
*/
func (a *API) GetOneChannel(channelIdentifier string) (channel string, err error) {
	defer func() {
		log.Trace(fmt.Sprintf("Api GetOneChannel in channel address=%s,out channel=\n%s,err=%v", channelIdentifier, channel, err))
	}()
	channelIdentifierHash := common.HexToHash(channelIdentifier)
	c, err := a.api.GetChannel(channelIdentifierHash)
	if err != nil {
		log.Error(err.Error())
		return
	}
	d := &v1.ChannelDataDetail{
		ChannelIdentifier:        common.BytesToHash(c.Key).String(),
		OpenBlockNumber:          c.ChannelIdentifier.OpenBlockNumber,
		PartnerAddress:           c.PartnerAddress().String(),
		Balance:                  c.OurBalance(),
		PartnerBalance:           c.PartnerBalance(),
		State:                    c.State,
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
	channel, err = marshal(d)
	return
}

/*
Deposit try to open a new channel on contract with
`partnerAddress` . the `settleTimeout` is the settle time of
the new channel.  `balanceStr` is the token to deposit to this channel and it  must be positive
 if `NewChannel` is true,  a new channel must be created and if `settleTimeout` is zero then it will be set as default
settle timeout.
if `NewChannel` is false, `settleTimeout` must be zero.

	//如果NewChannel为true
	//  SettleTimeout表示新建通道的结算窗口,如果SettleTimeout为0,则用系统默认计算窗口
	//如果NewChannel为 false
	//  SettleTimeout 必须为0

example returns:
{
    "channel_identifier": "0x97f73562938f6d538a07780b29847330e97d40bb8d0f23845a798912e76970e1",
    "open_block_number": 2560271,
    "partner_address": "0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5",
    "balance": 50,
    "partner_balance": 0,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "state": 1,
    "StateString": "opened",
    "settle_timeout": 150,
    "reveal_timeout": 0
}
*/
func (a *API) Deposit(partnerAddress, tokenAddress string, settleTimeout int, balanceStr string, newChannel bool) (callID string, err error) {
	callID = utils.NewRandomHash().String()
	result := newResult()
	a.callID2result[callID] = result
	go func() {
		r, e := a.deposit(partnerAddress, tokenAddress, settleTimeout, balanceStr, newChannel)
		result.Result = r
		result.Err = e
		result.Done = true
		a.callID2result[callID] = result
	}()
	return
}

func (a *API) deposit(partnerAddress, tokenAddress string, settleTimeout int, balanceStr string, newcChannel bool) (channel string, err error) {
	defer func() {
		log.Trace(fmt.Sprintf("Api Deposit in partnerAddress=%s,tokenAddress=%s,settletTimeout=%d,balanceStr=%s\nout channel=\n%s,err=%v",
			partnerAddress, tokenAddress, settleTimeout, balanceStr, channel, err,
		))
	}()
	partnerAddr, err := utils.HexToAddressWithoutValidation(partnerAddress)
	if err != nil {
		return
	}
	tokenAddr, err := utils.HexToAddressWithoutValidation(tokenAddress)
	if err != nil {
		return
	}
	balance, _ := new(big.Int).SetString(balanceStr, 0)
	c, err := a.api.DepositAndOpenChannel(tokenAddr, partnerAddr, settleTimeout, a.api.Photon.Config.RevealTimeout, balance, newcChannel)
	if err != nil {
		log.Error(err.Error())
		return
	}
	d := &v1.ChannelData{
		ChannelIdentifier:   common.BytesToHash(c.Key).String(),
		PartnerAddrses:      c.PartnerAddress().String(),
		Balance:             c.OurBalance(),
		PartnerBalance:      c.PartnerBalance(),
		State:               c.State,
		StateString:         c.State.String(),
		OpenBlockNumber:     c.ChannelIdentifier.OpenBlockNumber,
		RevealTimeout:       c.RevealTimeout,
		SettleTimeout:       c.SettleTimeout,
		TokenAddress:        c.TokenAddress().String(),
		LockedAmount:        c.OurAmountLocked(),
		PartnerLockedAmount: c.PartnerAmountLocked(),
	}
	channel, err = marshal(d)
	return

}

/*
CloseChannel close the  channel
如果force 为false,则表示希望双方协商关闭通道,
如果force为true,则表示希望直接连上关闭通道,不需要对方同意.
example returns:
{
    "channel_identifier": "0x97f73562938f6d538a07780b29847330e97d40bb8d0f23845a798912e76970e1",
    "open_block_number": 2560271,
    "partner_address": "0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5",
    "balance": 50,
    "partner_balance": 0,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "state": 2,
    "StateString": "closed",
    "settle_timeout": 150,
    "reveal_timeout": 0
}
*/
func (a *API) CloseChannel(channelIdentifier string, force bool) (callID string, err error) {
	callID = utils.NewRandomHash().String()
	result := newResult()
	a.callID2result[callID] = result
	go func() {
		r, e := a.closeChannel(channelIdentifier, force)
		result.Result = r
		result.Err = e
		result.Done = true
		a.callID2result[callID] = result
	}()
	return

}
func (a *API) closeChannel(channelIdentifier string, force bool) (channel string, err error) {
	defer func() {
		log.Trace(fmt.Sprintf("Api CloseChannel in channelIdentifier=%s,out channel=\n%s,err=%v",
			channelIdentifier, channel, err,
		))
	}()
	channelIdentifierHash := common.HexToHash(channelIdentifier)
	c, err := a.api.GetChannel(channelIdentifierHash)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if force {
		c, err = a.api.Close(c.TokenAddress(), c.PartnerAddress())
		if err != nil {
			log.Error(err.Error())
			return
		}
	} else {
		c, err = a.api.CooperativeSettle(c.TokenAddress(), c.PartnerAddress())
		if err != nil {
			log.Error(err.Error())
			return
		}
	}
	d := &v1.ChannelData{
		ChannelIdentifier:   common.BytesToHash(c.Key).String(),
		PartnerAddrses:      c.PartnerAddress().String(),
		Balance:             c.OurBalance(),
		PartnerBalance:      c.PartnerBalance(),
		State:               c.State,
		SettleTimeout:       c.SettleTimeout,
		TokenAddress:        c.TokenAddress().String(),
		LockedAmount:        c.OurAmountLocked(),
		PartnerLockedAmount: c.PartnerAmountLocked(),
	}
	channel, err = marshal(d)
	return
}

/*
SettleChannel settle a channel
在通道已经关闭的情况下,过了结算窗口期以后,用户可以在合约上进行结算.
example returns:
{
    "channel_identifier": "0x97f73562938f6d538a07780b29847330e97d40bb8d0f23845a798912e76970e1",
    "open_block_number": 2560271,
    "partner_address": "0xf0f6E53d6bbB9Debf35Da6531eC9f1141cd549d5",
    "balance": 50,
    "partner_balance": 0,
    "locked_amount": 0,
    "partner_locked_amount": 0,
    "token_address": "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2",
    "state": 3,
    "StateString": "settled",
    "settle_timeout": 150,
    "reveal_timeout": 0
}
*/
func (a *API) SettleChannel(channelIdentifier string) (callID string, err error) {
	callID = utils.NewRandomHash().String()
	result := newResult()
	a.callID2result[callID] = result
	go func() {
		r, e := a.settleChannel(channelIdentifier)
		result.Result = r
		result.Err = e
		result.Done = true
		a.callID2result[callID] = result
	}()
	return
}
func (a *API) settleChannel(channelIdentifier string) (channel string, err error) {
	defer func() {
		log.Trace(fmt.Sprintf("Api SettleChannel in channelIdentifier=%s,out channel=\n%s,err=%v",
			channelIdentifier, channel, err,
		))
	}()

	channelIdentifierHash := common.HexToHash(channelIdentifier)
	c, err := a.api.GetChannel(channelIdentifierHash)
	if err != nil {
		log.Error(err.Error())
		return
	}
	c, err = a.api.Settle(c.TokenAddress(), c.PartnerAddress())
	if err != nil {
		log.Error(err.Error())
		return
	}
	d := &v1.ChannelData{
		ChannelIdentifier:   common.BytesToHash(c.Key).String(),
		PartnerAddrses:      c.PartnerAddress().String(),
		Balance:             c.OurBalance(),
		PartnerBalance:      c.PartnerBalance(),
		State:               c.State,
		SettleTimeout:       c.SettleTimeout,
		TokenAddress:        c.TokenAddress().String(),
		LockedAmount:        c.OurAmountLocked(),
		PartnerLockedAmount: c.PartnerAmountLocked(),
	}
	channel, err = marshal(d)
	return
}

// Deprecated
func (a *API) networkEvent(fromBlock, toBlock int64) (eventsString string, err error) {
	events, err := a.api.GetNetworkEvents(fromBlock, toBlock)
	if err != nil {
		log.Error(err.Error())
		return
	}
	eventsString, err = marshal(events)
	return
}

//Deprecated: TokensEvent GET /api/1/events/tokens/0x61c808d82a3ac53231750dadc13c777b59310bd9
func (a *API) tokensEvent(fromBlock, toBlock int64, tokenAddress string) (eventsString string, err error) {
	token, err := utils.HexToAddressWithoutValidation(tokenAddress)
	if err != nil {
		return
	}
	events, err := a.api.GetTokenNetworkEvents(token, fromBlock, toBlock)
	if err != nil {
		log.Error(err.Error())
		return
	}
	eventsString, err = marshal(events)
	return
}

//Deprecated: ChannelsEvent GET /api/1/events/channels/0x2a65aca4d5fc5b5c859090a6c34d164135398226?from_block=1337
func (a *API) channelsEvent(fromBlock, toBlock int64, channelIdentifier string) (eventsString string, err error) {
	channel := common.HexToHash(channelIdentifier)
	events, err := a.api.GetChannelEvents(channel, fromBlock, toBlock)
	if err != nil {
		log.Error(err.Error())
		return
	}
	eventsString, err = marshal(events)
	return
}

/*
Address returns node's checksum address
for example: returns "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2"
*/
func (a *API) Address() (addr string) {
	return a.api.Address().String()
}

/*
Tokens returns all the token have registered on Photon
for example:
[
    "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2"
]
*/
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

/*
TokenPartners  Get all the channel partners of this token.
获取我在`token`上与其他所有节点的通道.
for example:
[
    {
        "partner_address": "0x151E62a787d0d8d9EfFac182Eae06C559d1B68C2",
        "channel": "api/1/channles/0x79b789e88c3d2173af4048498f8c1ce66f019f33a6b8b06bedef51dde72bbbc1"
    },
    {
        "partner_address": "0x201B20123b3C489b47Fde27ce5b451a0fA55FD60",
        "channel": "api/1/channles/0xd971f803c7ea39ee050bf00ec9919269cf63ee5d0e968d5fe33a1a0f0004f73d"
    }
]
*/
func (a *API) TokenPartners(tokenAddress string) (channels string, err error) {
	tokenAddr, err := utils.HexToAddressWithoutValidation(tokenAddress)
	if err != nil {
		return
	}
	chs, err := a.api.GetChannelList(tokenAddr, utils.EmptyAddress)
	if err != nil {
		log.Error(err.Error())
		return
	}
	var datas []*partnersData
	for _, c := range chs {
		d := &partnersData{
			PartnerAddress: c.PartnerAddress().String(),
			Channel:        "api/1/channles/" + c.OurAddress.String(),
		}
		datas = append(datas, d)
	}
	channels, err = marshal(datas)
	return
}

/*
Transfers POST /api/1/transfers/0x2a65aca4d5fc5b5c859090a6c34d164135398226/0x61c808d82a3ac53231750dadc13c777b59310bd9
Initiating a Transfer
tokenAddress is  the token to transfer
targetAddress is address of the receipt of the transfer
amountstr is integer amount string
feestr is  always 0 now
isDirect is this should be True when no internet connection,otherwise false.
data: the info
example returns for a correct call:
transfer:
{
    "initiator_address": "0x292650fee408320D888e06ed89D938294Ea42f99",
    "target_address": "0x4B89Bff01009928784eB7e7d10Bf773e6D166066",
    "token_address": "0x663495a1b8e9Be17083b37924cFE39e17858F9e8",
    "amount": 1,
    "lockSecretHash": "0x5e86d58579cfbc77901a457d7f63e8ec6e47efc5848761f51e63729e7848a01d",
    "sync": true
}

the caller should call GetTransferStatus periodically to query this transfer's latest status.
*/
func (a *API) Transfers(tokenAddress, targetAddress string, amountstr string, feestr string, secretStr string, isDirect bool, data string) (transfer string, err error) {
	defer func() {
		log.Trace(fmt.Sprintf("Api Transfers tokenAddress=%s,targetAddress=%s,amountstr=%s,feestr=%s,secretStr=%s, isDirect=%v, data=%s \nout transfer=\n%s,err=%v",
			tokenAddress, targetAddress, amountstr, feestr, secretStr, isDirect, data, transfer, err,
		))
	}()
	tokenAddr, err := utils.HexToAddressWithoutValidation(tokenAddress)
	if err != nil {
		return
	}
	targetAddr, err := utils.HexToAddressWithoutValidation(targetAddress)
	if err != nil {
		return
	}
	if len(secretStr) != 0 && len(secretStr) != 64 && (strings.HasPrefix(secretStr, "0x") && len(secretStr) != 66) {
		err = errors.New("invalid secret")
		return
	}
	if len(data) > params.MaxTransferDataLen {
		err = errors.New("invalid data, data len must < 256")
		return
	}
	amount, _ := new(big.Int).SetString(amountstr, 0)
	fee, _ := new(big.Int).SetString(feestr, 0)
	secret := common.HexToHash(secretStr)
	if amount.Cmp(utils.BigInt0) <= 0 {
		err = errors.New("amount should be positive")
		return
	}
	result, err := a.api.TransferAsync(tokenAddr, amount, fee, targetAddr, secret, isDirect, data)
	if err != nil {
		log.Error(err.Error())
		return
	}
	req := &v1.TransferData{}
	req.LockSecretHash = result.LockSecretHash.String()
	req.Initiator = a.api.Photon.NodeAddress.String()
	req.Target = targetAddress
	req.Token = tokenAddress
	req.Amount = amount
	req.Secret = secretStr
	req.Fee = fee
	req.Data = data
	return marshal(req)
}

/*
TokenSwap token swap for maker for two Photon nodes
the role should only be  "maker" or "taker".
`role` only maker or taker, if i'm a taker ,I must call TokenSwap first,then maker call his TokenSwap
`lockSecretHash` if i'm taker,I only know lockSecretHash, I must specify a valid hash
`SecretStr` if i'm a maker, I know secret and also secret's hash, I must specify the `SecretStr` and can ignore `lockSecretHash`

*/
func (a *API) TokenSwap(role string, lockSecretHash string, SendingAmountStr, ReceivingAmountStr string, SendingToken, ReceivingToken, TargetAddress string, SecretStr string) (callID string, err error) {
	callID = utils.NewRandomHash().String()
	result := newResult()
	a.callID2result[callID] = result
	go func() {
		e := a.tokenSwap(role, lockSecretHash, SendingAmountStr, ReceivingAmountStr, SendingToken, ReceivingToken, TargetAddress, SecretStr)
		result.Err = e
		result.Done = true
		a.callID2result[callID] = result
	}()
	return
}
func (a *API) tokenSwap(role string, lockSecretHash string, SendingAmountStr, ReceivingAmountStr string, SendingToken, ReceivingToken, TargetAddress string, SecretStr string) (err error) {
	type Req struct {
		Role            string   `json:"role"`
		SendingAmount   *big.Int `json:"sending_amount"`
		SendingToken    string   `json:"sending_token"`
		ReceivingAmount int64    `json:"receiving_amount"`
		ReceivingToken  *big.Int `json:"receiving_token"`
	}

	var target common.Address
	target, err = utils.HexToAddressWithoutValidation(TargetAddress)
	if err != nil {
		return
	}
	if len(lockSecretHash) <= 0 {
		err = errors.New("LockSecretHash must not be empty")
		return
	}
	SendingAmount, _ := new(big.Int).SetString(SendingAmountStr, 0)
	ReceivingAmount, _ := new(big.Int).SetString(ReceivingAmountStr, 0)
	makerToken, err := utils.HexToAddressWithoutValidation(SendingToken)
	if err != nil {
		return
	}
	takerToken, err := utils.HexToAddressWithoutValidation(ReceivingToken)
	if err != nil {
		return
	}
	if role == "maker" {
		err = a.api.TokenSwapAndWait(lockSecretHash, makerToken, takerToken,
			a.api.Photon.NodeAddress, target, SendingAmount, ReceivingAmount, SecretStr)
	} else if role == "taker" {
		err = a.api.ExpectTokenSwap(lockSecretHash, takerToken, makerToken,
			target, a.api.Photon.NodeAddress, ReceivingAmount, SendingAmount)
	} else {
		err = fmt.Errorf("provided invalid token swap role %s", role)
	}
	return
}

//Stop stop Photon
func (a *API) Stop() {
	log.Trace("Api Stop")
	//test only
	a.api.Stop()
}

/*
ChannelFor3rdParty generate info for 3rd party use,
for update transfer and withdraw.

example returns:
{
    "channel_identifier": "0x029a853513e98050e670eb6d5f36217998a2c689ef2f1c65b5954051490d5965",
    "open_block_number": 2644876,
    "token_network_address": "0xa3b6481d1c6aa8ba538e8fa9d4d8b1dbadfd379c",
    "partner_address": "0x64d11d0cbb3f4f9bb3ee09709d4254f0899a6381",
    "update_transfer": {
        "nonce": 0,
        "transfer_amount": null,
        "locksroot": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "extra_hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "closing_signature": null,
        "non_closing_signature": null
    },
    "unlocks": null,
    "punishes": [
        {
            "lock_hash": "0xd4ec833949fa91e5f30b4e5e8b2e88cca10e8192a68e51bdb24d18220b3f519d",
            "additional_hash": "0xe800ff8e78b8e367fb165b76f6e0cd1f31d46e7fda640e02134eed4f5e983d53",
            "signature": "i24Lz6KVvDnlqsxhQzDu+IIx6jJKC4gdVyWg6NpkrfsEejzGV8F0CPB0oUUJjDZ2wmChKG6XjZQx24QkDmhsKhs="
        }
    ]
}
*/
func (a *API) ChannelFor3rdParty(channelIdentifier, thirdPartyAddress string) (r string, err error) {
	channelIdentifierHash := common.HexToHash(channelIdentifier)
	thirdPartyAddr, err := utils.HexToAddressWithoutValidation(thirdPartyAddress)
	if err != nil {
		return
	}
	if channelIdentifierHash == utils.EmptyHash || thirdPartyAddr == utils.EmptyAddress {
		err = errors.New("invalid argument")
		return
	}
	result, err := a.api.ChannelInformationFor3rdParty(channelIdentifierHash, thirdPartyAddr)
	if err != nil {
		log.Error(err.Error())
		return
	}
	r, err = marshal(result)
	return
}

/*
SwitchNetwork  switch between mesh and internet
*/
func (a *API) SwitchNetwork(isMesh bool) {
	log.Trace(fmt.Sprintf("Api SwitchNetwork isMesh=%v", isMesh))
	a.api.Photon.Config.IsMeshNetwork = isMesh
}

/*
UpdateMeshNetworkNodes updates all nodes in MeshNetwork.
Nodes within the same local network have higher priority.
*/
func (a *API) UpdateMeshNetworkNodes(nodesstr string) (err error) {
	defer func() {
		log.Trace(fmt.Sprintf("Api UpdateMeshNetworkNodes nodesstr=%s,out err=%v", nodesstr, err))
	}()
	var nodes []*network.NodeInfo
	err = json.Unmarshal([]byte(nodesstr), &nodes)
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = a.api.Photon.Protocol.UpdateMeshNetworkNodes(nodes)
	if err != nil {
		log.Error(err.Error())
		return
	}
	return nil
}

/*
EthereumStatus  query the status between Photon and ethereum
todo fix it,remove this deprecated api
*/
func (a *API) EthereumStatus() (r string, err error) {
	c := a.api.Photon.Chain
	if c != nil && c.Client.Status == netshare.Connected {
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

// Subscription represents an event subscription where events are
// delivered on a data channel.
type Subscription struct {
	quitChan chan struct{}
}

// Unsubscribe cancels the sending of events to the data channel
// and closes the error channel.
func (s *Subscription) Unsubscribe() {
	close(s.quitChan)
}

// NotifyHandler is a client-side subscription callback to invoke on events and
// subscription failure.
type NotifyHandler interface {
	//some unexpected error
	OnError(errCode int, failure string)
	//OnStatusChange server connection status change
	OnStatusChange(s string)
	//OnReceivedTransfer  receive a transfer
	OnReceivedTransfer(tr string)
	//OnSentTransfer a transfer sent success
	OnSentTransfer(tr string)
	// OnNotify get some important message Photon want to notify upper application
	OnNotify(level int, info string)
}

/*
Subscribe  As to Status Notification, we put these codebase into an individual package
 and use channel to communication.
 To avoid write block, we can write data through select.
 We should make effort to avoid start go routine.
 If there's need to create a new Photon instance, sub.Unsubscribe must be invoked to do that or memory leakage will occur.
*/
func (a *API) Subscribe(handler NotifyHandler) (sub *Subscription, err error) {
	sub = &Subscription{
		quitChan: make(chan struct{}),
	}
	cs := v1.ConnectionStatus{
		XMPPStatus: netshare.Disconnected,
		EthStatus:  netshare.Disconnected,
	}

	var xn <-chan netshare.Status
	switch t := a.api.Photon.Transport.(type) {
	case *network.MatrixMixTransport:
		xn, err = t.GetNotify()
		if err != nil {
			log.Error(fmt.Sprintf("matrix transport get nofity err %s", err))
			return
		}
	case *network.MixTransport:
		xn, err = t.GetNotify()
		if err != nil {
			log.Error(fmt.Sprintf("mix transport get notify err %s", err))
			return
		}
	default:
		xn = make(chan netshare.Status)
	}
	go func() {
		rpanic.RegisterErrorNotifier("API SubscribeNeighbour")
		for {
			var err error
			var d []byte
			select {
			case err = <-rpanic.GetNotify():
				handler.OnError(32, err.Error())
			case s := <-a.api.Photon.EthConnectionStatus:
				cs.EthStatus = s
				cs.LastBlockTime = a.api.Photon.GetDao().GetLastBlockNumberTime().Format(v1.BlockTimeFormat)
				d, err = json.Marshal(cs)
				handler.OnStatusChange(string(d))
			case s := <-xn:
				cs.XMPPStatus = s
				cs.LastBlockTime = a.api.Photon.GetDao().GetLastBlockNumberTime().Format(v1.BlockTimeFormat)
				d, err = json.Marshal(cs)
				handler.OnStatusChange(string(d))
			case t, ok := <-a.api.Photon.NotifyHandler.GetSentTransferChan():
				if ok {
					d, err = json.Marshal(t)
					handler.OnSentTransfer(string(d))
				}
			case t, ok := <-a.api.Photon.NotifyHandler.GetReceivedTransferChan():
				if ok {
					d, err = json.Marshal(t)
					handler.OnReceivedTransfer(string(d))
				}
			case n, ok := <-a.api.Photon.NotifyHandler.GetNoticeChan():
				if ok {
					handler.OnNotify(int(n.Level), n.Info)
				}
			case <-sub.quitChan:
				return
			}
			if err != nil {
				log.Error(fmt.Sprintf("err =%s", err))
			}
		}

	}()
	return
}

/*
GetTransferStatus return transfer result
status should be one the following
// TransferStatusInit init
TransferStatusInit = 0

// TransferStatusCanCancel transfer can cancel right now
TransferStatusCanCancel =1

// TransferStatusCanNotCancel transfer can not cancel
TransferStatusCanNotCancel =2

// TransferStatusSuccess transfer already success
TransferStatusSuccess =3

// TransferStatusCanceled transfer cancel by user request
TransferStatusCanceled =4

// TransferStatusFailed transfer already failed
TransferStatusFailed =5

example returns:
{
    "LockSecretHash": "0x2f6dbd44fa95d7edc840570d3bc847e24846a5422fffa324cdd9c5cab945857e",
    "Status": 2,
    "StatusMessage": "MediatedTransfer 正在发送 target=4b89\nMediatedTransfer 发送成功\n收到 SecretRequest, from=3af7\nRevealSecret 正在发送 target=3af7\nRevealSecret 发送成功\n收到 RevealSecret, from=4b89\nUnlock 正在发送 target=4b89\nUnLock 发送成功,交易成功.\n"
}
*/
func (a *API) GetTransferStatus(tokenAddressStr string, lockSecretHashStr string) (r string, err error) {
	defer func() {
		log.Trace(fmt.Sprintf("Api GetTransferStatus tokenAddressStr=%s,lockSecretHashStr=%s, err=%s\n",
			tokenAddressStr, lockSecretHashStr, err,
		))
	}()
	tokenAddress, err := utils.HexToAddress(tokenAddressStr)
	if err != nil {
		log.Error(err.Error())
		return
	}
	ts, err := a.api.Photon.GetDao().GetTransferStatus(tokenAddress, common.HexToHash(lockSecretHashStr))
	if err != nil {
		log.Error(fmt.Sprintf("err =%s", err))
		return
	}
	r, err = marshal(ts)
	return
}

// NotifyNetworkDown :
func (a *API) NotifyNetworkDown() error {
	return a.api.NotifyNetworkDown()
}

// GetCallResult :
func (a *API) GetCallResult(callID string) (r string, err error) {
	result, ok := a.callID2result[callID]
	if !ok {
		err = errors.New("not found")
		return
	}
	if result.Err != nil {
		err = result.Err
		return
	}
	r = result.Result
	done := result.Done
	if done {
		err = result.Err
		delete(a.callID2result, callID)
	} else {
		err = errors.New("dealing")
	}
	return
}

// Withdraw :
/*
	1. withdraw
	{ "amount":3333,}
	2. prepare for withdraw:
	{"op":"preparewithdraw",}
	3. cancel prepare:
	{"op": "cancelprepare"}
*/
func (a *API) Withdraw(channelIdentifierHashStr, amountstr, op string) (callID string, err error) {
	callID = utils.NewRandomHash().String()
	result := newResult()
	a.callID2result[callID] = result
	go func() {
		r, e := a.withdraw(channelIdentifierHashStr, amountstr, op)
		result.Result = r
		result.Err = e
		result.Done = true
		a.callID2result[callID] = result
	}()
	return
}

func (a *API) withdraw(channelIdentifierHashStr, amountStr, op string) (r string, err error) {
	defer func() {
		log.Trace(fmt.Sprintf("Api Withdraw channelIdentifierHashStr=%s,amountStr=%s,op=%s,response=%s,err=%v",
			channelIdentifierHashStr, amountStr, op, r, err))
	}()
	const OpPrepareWithdraw = "preparewithdraw"
	const OpCancelPrepare = "cancelprepare"
	channelIdentifier := common.HexToHash(channelIdentifierHashStr)
	amount, _ := new(big.Int).SetString(amountStr, 0)
	c, err := a.api.GetChannel(channelIdentifier)
	if err != nil {
		log.Error(fmt.Sprintf("GetChannel %s err %s", utils.HPex(channelIdentifier), err))
		return
	}
	if amount != nil && amount.Cmp(utils.BigInt0) > 0 { // withdraw
		c, err = a.api.Withdraw(c.TokenAddress(), c.PartnerAddress(), amount)
		if err != nil {
			log.Error(fmt.Sprintf("Withdraw %s err %s", utils.HPex(channelIdentifier), err))
			return
		}
	} else {
		if op == OpPrepareWithdraw {
			c, err = a.api.PrepareForWithdraw(c.TokenAddress(), c.PartnerAddress())
		} else if op == OpCancelPrepare {
			c, err = a.api.CancelPrepareForWithdraw(c.TokenAddress(), c.PartnerAddress())
		} else {
			err = fmt.Errorf("unkown operation %s", op)
		}
		if err != nil {
			log.Error(err.Error())
			return
		}
	}
	d := &v1.ChannelData{
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
	r, err = marshal(d)
	return
}

// OnResume :
// 手机从后台切换至前台时调用
func (a *API) OnResume() (err error) {
	// 1. 强制网络重连
	err = a.NotifyNetworkDown()
	// 2. 还需要作什么???
	return
}

// GetSystemStatus :
func (a *API) GetSystemStatus() (r string, err error) {
	resp := a.api.SystemStatus()
	return resp.ToString(), nil
}
