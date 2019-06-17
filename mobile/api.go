package mobile

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/Photon/pfsproxy"
	"github.com/SmartMeshFoundation/Photon/rerr"

	"github.com/SmartMeshFoundation/Photon/dto"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"

	"math/big"

	"errors"

	"strings"

	photon "github.com/SmartMeshFoundation/Photon"
	"github.com/SmartMeshFoundation/Photon/internal/rpanic"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/network"
	"github.com/SmartMeshFoundation/Photon/network/netshare"
	"github.com/SmartMeshFoundation/Photon/params"
	v1 "github.com/SmartMeshFoundation/Photon/restful/v1"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

// API for export interface
//
// should not export any member because of gomobile's protocol
type API struct {
	startTime time.Time
	api       *photon.API
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
func (a *API) GetChannelList() (result string) {
	defer func() {
		log.Trace(fmt.Sprintf("ApiCall GetChannelList result=%s", result))
	}()
	chs, err := a.api.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	if err != nil {
		log.Error(err.Error())
		result = dto.NewErrorMobileResponse(err)
		return
	}
	var datas []*v1.ChannelData
	for _, c := range chs {
		d := &v1.ChannelData{
			ChannelIdentifier:   c.ChannelIdentifier.ChannelIdentifier.String(),
			OpenBlockNumber:     c.ChannelIdentifier.OpenBlockNumber,
			PartnerAddrses:      c.PartnerAddress().String(),
			Balance:             c.OurBalance(),
			PartnerBalance:      c.PartnerBalance(),
			State:               c.State,
			StateString:         c.State.String(),
			DelegateState:       c.DelegateState,
			DelegateStateString: c.DelegateState.String(),
			TokenAddress:        c.TokenAddress().String(),
			SettleTimeout:       c.SettleTimeout,
			RevealTimeout:       c.RevealTimeout,
			LockedAmount:        c.OurAmountLocked(),
			PartnerLockedAmount: c.PartnerAmountLocked(),
		}
		if c.State == channeltype.StateClosed {
			res := a.api.GetChannelSettleBlock(c.ChannelIdentifier.ChannelIdentifier)
			d.BlockNumberNow = res.BlockNumberNow
			d.BlockNumberChannelCanSettle = res.BlockNumberChannelCanSettle
		}
		datas = append(datas, d)
	}
	result = dto.NewSuccessMobileResponse(datas)
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
    "OurUnknownSecretLocks": {},
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
func (a *API) GetOneChannel(channelIdentifier string) (result string) {
	defer func() {
		log.Trace(fmt.Sprintf("Api GetOneChannel in channel address=%s,out result=\n%s", channelIdentifier, result))
	}()
	channelIdentifierHash := common.HexToHash(channelIdentifier)
	c, err := a.api.GetChannel(channelIdentifierHash)
	if err != nil {
		log.Error(err.Error())
		result = dto.NewErrorMobileResponse(err)
		return
	}
	result = dto.NewSuccessMobileResponse(channeltype.ChannelSerialization2ChannelDataDetail(c))
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
func (a *API) Deposit(partnerAddress, tokenAddress string, settleTimeout int, balanceStr string, newChannel bool) (result string) {
	r, err := a.deposit(partnerAddress, tokenAddress, settleTimeout, balanceStr, newChannel)
	if err != nil {
		result = dto.NewErrorMobileResponse(err)
		return
	}
	result = dto.NewSuccessMobileResponse(r)
	return
}

func (a *API) deposit(partnerAddress, tokenAddress string, settleTimeout int, balanceStr string, newcChannel bool) (channel *channeltype.ChannelDataDetail, err error) {
	defer func() {
		log.Trace(fmt.Sprintf("Api Deposit in partnerAddress=%s,tokenAddress=%s,settletTimeout=%d,balanceStr=%s\nout channel=\n%s,err=%v",
			partnerAddress, tokenAddress, settleTimeout, balanceStr, utils.StringInterface(channel, 5), err,
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
	if c != nil {
		channel = channeltype.ChannelSerialization2ChannelDataDetail(c)
	}
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
func (a *API) CloseChannel(channelIdentifier string, force bool) (result string) {
	c, err := a.closeChannel(channelIdentifier, force)
	if err != nil {
		result = dto.NewErrorMobileResponse(err)
		return
	}
	result = dto.NewSuccessMobileResponse(c)
	return

}
func (a *API) closeChannel(channelIdentifier string, force bool) (channel *channeltype.ChannelDataDetail, err error) {
	defer func() {
		log.Trace(fmt.Sprintf("Api CloseChannel in channelIdentifier=%s,out channel=\n%s,err=%v",
			channelIdentifier, utils.StringInterface(channel, 5), err,
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
	channel = channeltype.ChannelSerialization2ChannelDataDetail(c)
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
func (a *API) SettleChannel(channelIdentifier string) (result string) {
	c, err := a.settleChannel(channelIdentifier)
	if err != nil {
		result = dto.NewErrorMobileResponse(err)
		return
	}
	result = dto.NewSuccessMobileResponse(c)
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
func (a *API) Withdraw(channelIdentifierHashStr, amountstr, op string) (result string) {
	c, err := a.withdraw(channelIdentifierHashStr, amountstr, op)
	if err != nil {
		result = dto.NewErrorMobileResponse(err)
		return
	}
	result = dto.NewSuccessMobileResponse(c)
	return
}

func (a *API) settleChannel(channelIdentifier string) (channel *channeltype.ChannelDataDetail, err error) {
	defer func() {
		log.Trace(fmt.Sprintf("Api SettleChannel in channelIdentifier=%s,out channel=\n%s,err=%v",
			channelIdentifier, utils.StringInterface(channel, 5), err,
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
	channel = channeltype.ChannelSerialization2ChannelDataDetail(c)
	return
}

// Deprecated
func (a *API) networkEvent(fromBlock, toBlock int64) (result string) {
	defer func() {
		log.Trace(fmt.Sprintf("ApiCall networkEvent result=%s", result))
	}()
	events, err := a.api.GetNetworkEvents(fromBlock, toBlock)
	if err != nil {
		log.Error(err.Error())
		return dto.NewErrorMobileResponse(err)
	}
	return dto.NewSuccessMobileResponse(events)
}

//Deprecated: ChannelsEvent GET /api/1/events/channels/0x2a65aca4d5fc5b5c859090a6c34d164135398226?from_block=1337
func (a *API) channelsEvent(fromBlock, toBlock int64, channelIdentifier string) (result string) {
	defer func() {
		log.Trace(fmt.Sprintf("ApiCall channelsEvent result=%s", result))
	}()
	channel := common.HexToHash(channelIdentifier)
	events, err := a.api.GetChannelEvents(channel, fromBlock, toBlock)
	if err != nil {
		log.Error(err.Error())
		return dto.NewErrorMobileResponse(err)
	}
	return dto.NewSuccessMobileResponse(events)
}

/*
Address returns node's checksum address
for example: returns "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2"
*/
func (a *API) Address() (result string) {
	defer func() {
		log.Trace(fmt.Sprintf("ApiCall Address result=%s", result))
	}()
	return dto.NewSuccessMobileResponse(a.api.Address().String())
}

/*
Tokens returns all the token have registered on Photon
for example:
[
    "0x7B874444681F7AEF18D48f330a0Ba093d3d0fDD2"
]
*/
func (a *API) Tokens() (result string) {
	defer func() {
		log.Trace(fmt.Sprintf("ApiCall Tokens result=%s", result))
	}()
	tokens := a.api.Tokens()
	return dto.NewSuccessMobileResponse(tokens)
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
func (a *API) TokenPartners(tokenAddress string) (result string) {
	defer func() {
		log.Trace(fmt.Sprintf("ApiCall TokenPartners result=%s", result))
	}()
	tokenAddr, err := utils.HexToAddressWithoutValidation(tokenAddress)
	if err != nil {
		return
	}
	chs, err := a.api.GetChannelList(tokenAddr, utils.EmptyAddress)
	if err != nil {
		log.Error(err.Error())
		return dto.NewErrorMobileResponse(err)
	}
	var datas []*partnersData
	for _, c := range chs {
		d := &partnersData{
			PartnerAddress: c.PartnerAddress().String(),
			Channel:        "api/1/channles/" + c.OurAddress.String(),
		}
		datas = append(datas, d)
	}
	return dto.NewSuccessMobileResponse(datas)
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

the caller should call GetSentTransferDetail periodically to query this transfer's latest status.
*/
func (a *API) Transfers(tokenAddress, targetAddress string, amountstr string, secretStr string, isDirect bool, data string, routeInfoStr string) (result string) {
	defer func() {
		log.Trace(fmt.Sprintf("Api Transfers tokenAddress=%s,targetAddress=%s,amountstr=%s,secretStr=%s,isDirect=%v, data=%s routeInfo=%s\nout transfer=\n%s ",
			tokenAddress, targetAddress, amountstr, secretStr, isDirect, data, result, routeInfoStr,
		))
	}()
	tokenAddr, err := utils.HexToAddressWithoutValidation(tokenAddress)
	if err != nil {
		err = rerr.ErrArgumentError.AppendError(err)
		return dto.NewErrorMobileResponse(err)
	}
	targetAddr, err := utils.HexToAddressWithoutValidation(targetAddress)
	if err != nil {
		err = rerr.ErrArgumentError.AppendError(err)
		return dto.NewErrorMobileResponse(err)
	}
	if len(secretStr) != 0 && len(secretStr) != 64 && (strings.HasPrefix(secretStr, "0x") && len(secretStr) != 66) {
		err = errors.New("invalid secret")
		err = rerr.ErrArgumentError.AppendError(err)
		return dto.NewErrorMobileResponse(err)
	}
	if len(data) > params.MaxTransferDataLen {
		err = errors.New("invalid data, data len must < 256")
		err = rerr.ErrArgumentError.AppendError(err)
		return dto.NewErrorMobileResponse(err)
	}
	amount, _ := new(big.Int).SetString(amountstr, 0)
	secret := common.HexToHash(secretStr)
	if amount.Cmp(utils.BigInt0) <= 0 {
		err = errors.New("amount should be positive")
		err = rerr.ErrArgumentError.AppendError(err)
		return dto.NewErrorMobileResponse(err)
	}
	// 解析指定的路由info
	var routeInfo []pfsproxy.FindPathResponse
	if routeInfoStr == "" {
		routeInfo = nil
	} else {
		err = json.Unmarshal([]byte(routeInfoStr), &routeInfo)
		if err != nil {
			err = fmt.Errorf("parse route info err=%s", err.Error())
			err = rerr.ErrArgumentError.AppendError(err)
			return dto.NewErrorMobileResponse(err)
		}
	}
	tr, err := a.api.TransferAsync(tokenAddr, amount, targetAddr, secret, isDirect, data, routeInfo)
	if err != nil {
		log.Error(err.Error())
		return dto.NewErrorMobileResponse(err)
	}
	req := &v1.TransferData{}
	req.LockSecretHash = tr.LockSecretHash.String()
	req.Initiator = a.api.Photon.NodeAddress.String()
	req.Target = targetAddress
	req.Token = tokenAddress
	req.Amount = amount
	req.Secret = secretStr
	req.Data = data
	return dto.NewSuccessMobileResponse(req)
}

/*
TokenSwap token swap for maker for two Photon nodes
the role should only be  "maker" or "taker".
`role` only maker or taker, if i'm a taker ,I must call TokenSwap first,then maker call his TokenSwap
`lockSecretHash` if i'm taker,I only know lockSecretHash, I must specify a valid hash
`SecretStr` if i'm a maker, I know secret and also secret's hash, I must specify the `SecretStr` and can ignore `lockSecretHash`

*/
//todo 手机接口暂时不导出tokenswap,
//func (a *API) TokenSwap(role string, lockSecretHash string, SendingAmountStr, ReceivingAmountStr string, SendingToken, ReceivingToken, TargetAddress string, SecretStr string) (callID string, err error) {
//	callID = utils.NewRandomHash().String()
//	result := newResult()
//	a.callID2result[callID] = result
//	go func() {
//		e := a.tokenSwap(role, lockSecretHash, SendingAmountStr, ReceivingAmountStr, SendingToken, ReceivingToken, TargetAddress, SecretStr)
//		result.Err = e
//		result.Done = true
//		a.callID2result[callID] = result
//	}()
//	return
//}
//func (a *API) tokenSwap(role string, lockSecretHash string, SendingAmountStr, ReceivingAmountStr string, SendingToken, ReceivingToken, TargetAddress string, SecretStr string) (err error) {
//	type Req struct {
//		Role            string   `json:"role"`
//		SendingAmount   *big.Int `json:"sending_amount"`
//		SendingToken    string   `json:"sending_token"`
//		ReceivingAmount int64    `json:"receiving_amount"`
//		ReceivingToken  *big.Int `json:"receiving_token"`
//	}
//
//	var target common.Address
//	target, err = utils.HexToAddressWithoutValidation(TargetAddress)
//	if err != nil {
//		return
//	}
//	if len(lockSecretHash) <= 0 {
//		err = errors.New("LockSecretHash must not be empty")
//		return
//	}
//	SendingAmount, _ := new(big.Int).SetString(SendingAmountStr, 0)
//	ReceivingAmount, _ := new(big.Int).SetString(ReceivingAmountStr, 0)
//	makerToken, err := utils.HexToAddressWithoutValidation(SendingToken)
//	if err != nil {
//		return
//	}
//	takerToken, err := utils.HexToAddressWithoutValidation(ReceivingToken)
//	if err != nil {
//		return
//	}
//	if role == "maker" {
//		err = a.api.TokenSwapAndWait(lockSecretHash, makerToken, takerToken,
//			a.api.Photon.NodeAddress, target, SendingAmount, ReceivingAmount, SecretStr)
//	} else if role == "taker" {
//		err = a.api.ExpectTokenSwap(lockSecretHash, takerToken, makerToken,
//			target, a.api.Photon.NodeAddress, ReceivingAmount, SendingAmount)
//	} else {
//		err = fmt.Errorf("provided invalid token swap role %s", role)
//	}
//	return
//}

//Stop stop Photon
func (a *API) Stop() {
	log.Info("Api Stop")
	if v1.QuitChan != nil {
		close(v1.QuitChan)
	}
	//test only
	a.api.Stop()
	//保证stop完成以后,再删除,也就是说上一个没有完全stop之前,是不能启动新的photon实例的
	delete(apiMonitor, a)
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
func (a *API) ChannelFor3rdParty(channelIdentifier, thirdPartyAddress string) (result string) {
	defer func() {
		log.Trace(fmt.Sprintf("ApiCall ChannelFor3rdParty result=%s", result))
	}()
	channelIdentifierHash := common.HexToHash(channelIdentifier)
	thirdPartyAddr, err := utils.HexToAddressWithoutValidation(thirdPartyAddress)
	if err != nil {
		err = rerr.ErrArgumentError.AppendError(err)
		return dto.NewErrorMobileResponse(err)
	}
	if channelIdentifierHash == utils.EmptyHash || thirdPartyAddr == utils.EmptyAddress {
		err = errors.New("invalid argument")
		err = rerr.ErrArgumentError.AppendError(err)
		return dto.NewErrorMobileResponse(err)
	}
	resp, err := a.api.ChannelInformationFor3rdParty(channelIdentifierHash, thirdPartyAddr)
	if err != nil {
		log.Error(err.Error())
		err = rerr.ErrArgumentError.AppendError(err)
		return dto.NewErrorMobileResponse(err)
	}
	return dto.NewSuccessMobileResponse(resp)
}

/*
SwitchNetwork  switch between mesh and internet
*/
func (a *API) SwitchNetwork(isMesh bool) {
	log.Trace(fmt.Sprintf("Api SwitchNetwork isMesh=%v", isMesh))
	if isMesh {
		log.Trace("use NotifyNetworkDown")
		err := a.api.NotifyNetworkDown()
		if err != nil {
			panic("never happen")
		}
	}
}

/*
UpdateMeshNetworkNodes updates all nodes in MeshNetwork.
Nodes within the same local network have higher priority.
*/
func (a *API) UpdateMeshNetworkNodes(nodesstr string) (result string) {
	defer func() {
		log.Trace(fmt.Sprintf("Api UpdateMeshNetworkNodes nodesstr=%s,out result=%v", nodesstr, result))
	}()
	var nodes []*network.NodeInfo
	err := json.Unmarshal([]byte(nodesstr), &nodes)
	if err != nil {
		log.Error(err.Error())
		err = rerr.ErrArgumentError.AppendError(err)
		return dto.NewErrorMobileResponse(err)
	}
	err = a.api.Photon.Protocol.UpdateMeshNetworkNodes(nodes)
	if err != nil {
		log.Error(err.Error())
		err = rerr.ErrArgumentError.AppendError(err)
		return dto.NewErrorMobileResponse(err)
	}
	return dto.NewSuccessMobileResponse(nil)
}

/*
GetSentTransfers retuns list of sent transfer between `from_block` and `to_block`
*/
func (a *API) GetSentTransfers(tokenAddressStr string, from, to int64) (result string) {
	defer func() {
		log.Trace(fmt.Sprintf("ApiCall GetSentTransferDetails result=%s", result))
	}()
	log.Trace(fmt.Sprintf("from=%d,to=%d\n", from, to))
	tokenAddress := utils.EmptyAddress
	if tokenAddressStr != "" {
		tokenAddress = common.HexToAddress(tokenAddressStr)
	}
	trs, err := a.api.GetSentTransferDetails(tokenAddress, from, to)
	if err != nil {
		log.Error(err.Error())
		return dto.NewErrorMobileResponse(err)
	}
	return dto.NewSuccessMobileResponse(trs)
}

/*
GetReceivedTransfers retuns list of received transfer between `from_block` and `to_block`
it contains token swap
*/
func (a *API) GetReceivedTransfers(tokenAddressStr string, from, to int64) (result string) {
	defer func() {
		log.Trace(fmt.Sprintf("ApiCall GetReceivedTransfers result=%s", result))
	}()
	tokenAddress := utils.EmptyAddress
	if tokenAddressStr != "" {
		tokenAddress = common.HexToAddress(tokenAddressStr)
	}
	trs, err := a.api.GetReceivedTransfers(tokenAddress, from, to, -1, -1)
	if err != nil {
		log.Error(err.Error())
		return dto.NewErrorMobileResponse(err)
	}
	return dto.NewSuccessMobileResponse(trs)
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
	/* OnNotify get some important message Photon want to notify upper application
	level: 0:info,1:warn,2:error
	info: type InfoStruct struct {
		Type    int
		Message interface{}
		}
	当info.Type=0 表示Message是一个string,1表示Message是TransferStatus
	*/
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
				log.Error(fmt.Sprintf("photon panic because of unkown err %s", err))
				handler.OnError(32, err.Error())
			case s := <-a.api.Photon.EthConnectionStatus:
				cs.EthStatus = s
				cs.LastBlockTime = a.api.Photon.GetDao().GetLastBlockNumberTime().Format(v1.BlockTimeFormat)
				d, err = json.Marshal(cs)
				log.Info(fmt.Sprintf("notify OnStatusChange=%s", d))
				handler.OnStatusChange(string(d))
			case s := <-xn:
				cs.XMPPStatus = s
				log.Info(fmt.Sprintf("status change to %d", cs.XMPPStatus))
				cs.LastBlockTime = a.api.Photon.GetDao().GetLastBlockNumberTime().Format(v1.BlockTimeFormat)
				d, err = json.Marshal(cs)
				log.Info(fmt.Sprintf("notify OnStatusChange=%s", d))
				handler.OnStatusChange(string(d))
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
func (a *API) GetTransferStatus(tokenAddressStr string, lockSecretHashStr string) (result string) {
	defer func() {
		log.Trace(fmt.Sprintf("Api GetSentTransferDetail tokenAddressStr=%s,lockSecretHashStr=%s, result=%s\n",
			tokenAddressStr, lockSecretHashStr, result,
		))
	}()
	tokenAddress, err := utils.HexToAddress(tokenAddressStr)
	if err != nil {
		log.Error(err.Error())
		err = rerr.ErrArgumentError.AppendError(err)
		return dto.NewErrorMobileResponse(err)
	}
	ts, err := a.api.Photon.GetDao().GetSentTransferDetail(tokenAddress, common.HexToHash(lockSecretHashStr))
	if err != nil {
		log.Error(fmt.Sprintf("err =%s", err))
		return dto.NewErrorMobileResponse(err)
	}
	return dto.NewSuccessMobileResponse(ts)
}

// NotifyNetworkDown :
func (a *API) NotifyNetworkDown() (result string) {
	defer func() {
		log.Trace(fmt.Sprintf("ApiCall NotifyNetworkDown result=%s", result))
	}()
	err := a.api.NotifyNetworkDown()
	return dto.NewErrorMobileResponse(err)
}

func (a *API) withdraw(channelIdentifierHashStr, amountStr, op string) (channel *channeltype.ChannelDataDetail, err error) {
	defer func() {
		log.Trace(fmt.Sprintf("Api Withdraw channelIdentifierHashStr=%s,amountStr=%s,op=%s,response=%s,err=%v",
			channelIdentifierHashStr, amountStr, op, utils.StringInterface(channel, 3), err))
	}()
	const OpPrepareWithdraw = "preparewithdraw"
	const OpCancelPrepare = "cancelprepare"
	channelIdentifier := common.HexToHash(channelIdentifierHashStr)
	amount, _ := new(big.Int).SetString(amountStr, 0)
	c, err := a.api.GetChannel(channelIdentifier)
	if err != nil {
		err = rerr.ErrChannelNotFound.WithData(channelIdentifier) //不要暴露底层错误
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
			err = rerr.ErrArgumentError.Errorf("unkown operation %s", op)
		}
		if err != nil {
			log.Error(err.Error())
			return
		}
	}
	channel = channeltype.ChannelSerialization2ChannelDataDetail(c)
	return
}

// OnResume 手机从后台切换至前台时调用
func (a *API) OnResume() string {
	// 1. 强制网络重连
	result := a.NotifyNetworkDown()
	// 2. 还需要作什么???
	return result
}

// GetSystemStatus 查询系统状态,
func (a *API) GetSystemStatus() (result string) {
	defer func() {
		log.Trace(fmt.Sprintf("ApiCall GetSystemStatus result=%s", result))
	}()
	resp, err := a.api.SystemStatus()
	return dto.NewMobileResponse(err, resp)
}

/*
FindPath 查询所有从我到target的最低费用路径,该调用总是找pfs问路
example:
{
        "path_id": 0,
        "path_hop": 2,
        "fee": 10000000000,
        "result": [
            "0x3bc7726c489e617571792ac0cd8b70df8a5d0e22",
            "0x8a32108d269c11f8db859ca7fac8199ca87a2722",
            "0xefb2e46724f675381ce0b3f70ea66383061924e9"
        ]
    }
*/
func (a *API) FindPath(targetStr, tokenStr, amountStr string) (result string) {
	defer func() {
		log.Trace(fmt.Sprintf("ApiCall FindPath result=%s", result))
	}()
	target := common.HexToAddress(targetStr)
	token := common.HexToAddress(tokenStr)
	amount, isSuccess := new(big.Int).SetString(amountStr, 0)
	if !isSuccess {
		err := rerr.ErrArgumentError.Errorf("arg amount err %s", amountStr)
		return dto.NewErrorMobileResponse(err)
	}
	routes, err := a.api.FindPath(target, token, amount)
	if err != nil {
		return dto.NewErrorMobileResponse(err)
	}
	return dto.NewSuccessMobileResponse(routes)
}

/*
ContractCallTXQuery 合约调用TX查询接口,4个参数均可传空值,空值即为不限制,4个参数对应的查询条件关系为and
channelIdentifierStr 有值时按通道ID查询
openBlockNumber 有值时按通道OpenBlockNumber查询,一般配合channelIdentifierStr参数一起使用,以精确定位到某一个通道
txTypeStr 有值时按tx类型查询,取值:
	TXInfoTypeDeposit            = "ChannelDeposit"
	TXInfoTypeClose              = "ChannelClose"
	TXInfoTypeSettle             = "ChannelSettle"
	TXInfoTypeCooperateSettle    = "CooperateSettle"
	TXInfoTypeUpdateBalanceProof = "UpdateBalanceProof"
	TXInfoTypeUnlock             = "Unlock"
	TXInfoTypePunish             = "Punish"
	TXInfoTypeWithdraw           = "Withdraw"
	TXInfoTypeApproveDeposit     = "ApproveDeposit"
	TXInfoTypeRegisterSecret     = "RegisterSecret"
txStatusStr 有值时按tx状态查询,取值:
	TXInfoStatusPending = "pending"
	TXInfoStatusSuccess = "success"
	TXInfoStatusFailed  = "failed"
*/
func (a *API) ContractCallTXQuery(channelIdentifierStr string, openBlockNumber int, tokenAddressStr, txTypeStr, txStatusStr string) (result string) {
	defer func() {
		log.Trace(fmt.Sprintf("ApiCall ContractCallTXQuery result=%s", result))
	}()
	req := &photon.ContractCallTXQueryParams{
		ChannelIdentifier: channelIdentifierStr,
		OpenBlockNumber:   int64(openBlockNumber),
		TokenAddress:      tokenAddressStr,
		TXType:            models.TXInfoType(txTypeStr),
		TXStatus:          models.TXInfoStatus(txStatusStr),
	}
	list, err := a.api.ContractCallTXQuery(req)
	if err != nil {
		return dto.NewErrorMobileResponse(err)
	}
	return dto.NewSuccessMobileResponse(list)
}

// Version 获取版本信息
func (a *API) Version() string {
	return dto.NewSuccessMobileResponse(a.api.GetBuildInfo())
}

// GetAssetsOnToken 参数逗号分隔
func (a *API) GetAssetsOnToken(tokenListStr string) (result string) {
	defer func() {
		log.Trace(fmt.Sprintf("ApiCall GetAssetsOnToken result=%s", result))
	}()
	var tokenList []common.Address
	if tokenListStr != "" {
		ss := strings.Split(tokenListStr, ",")
		for _, s := range ss {
			tokenList = append(tokenList, common.HexToAddress(s))
		}
	}
	resp, err := a.api.GetAssetsOnToken(tokenList)
	if err != nil {
		return dto.NewErrorMobileResponse(err)
	}
	return dto.NewSuccessMobileResponse(resp)
}

// DebugUploadLogfile 上传photon日志到logserver
func (a *API) DebugUploadLogfile() string {
	err := a.api.UploadLogFile()
	if err != nil {
		return dto.NewErrorMobileResponse(err)
	}
	return dto.NewSuccessMobileResponse(nil)
}
