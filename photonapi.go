package photon

import (
	"compress/gzip"
	"encoding/binary"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"

	"github.com/SmartMeshFoundation/Photon/params"

	"fmt"

	"math/big"

	"bytes"

	"sort"

	"context"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/network"
	"github.com/SmartMeshFoundation/Photon/network/netshare"
	"github.com/SmartMeshFoundation/Photon/pfsproxy"
	"github.com/SmartMeshFoundation/Photon/pmsproxy"
	"github.com/SmartMeshFoundation/Photon/rerr"
	"github.com/SmartMeshFoundation/Photon/transfer"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

//API photon for user
/* #nolint */
type API struct {
	Photon *Service
}

//NewPhotonAPI create CLI interface.
func NewPhotonAPI(photon *Service) *API {
	return &API{Photon: photon}
}

//Address return this node's address
func (r *API) Address() common.Address {
	return r.Photon.NodeAddress
}

//Tokens Return a list of the tokens registered with the default registry.
func (r *API) Tokens() (addresses []common.Address) {
	tokens, err := r.Photon.dao.GetAllTokens()
	if err != nil {
		log.Error(fmt.Sprintf("GetAllTokens err %s", err))
		return
	}
	for t := range tokens {
		addresses = append(addresses, t)
	}
	return
}

/*
GetChannelList Returns a list of channels associated with the optionally given
           `token_address` and/or `partner_address
Args:
            token_address (bin): an optionally provided token address
            partner_address (bin): an optionally provided partner address

        Return:
            A list containing all channels the node participates. Optionally
            filtered by a token address and/or partner address.

        Raises:
            KeyError: An error occurred when the token address is unknown to the node.
*/
func (r *API) GetChannelList(tokenAddress common.Address, partnerAddress common.Address) (cs []*channeltype.Serialization, err error) {
	return r.Photon.dao.GetChannelList(tokenAddress, partnerAddress)
}

//GetChannel get channel by address
func (r *API) GetChannel(ChannelIdentifier common.Hash) (c *channeltype.Serialization, err error) {
	return r.Photon.dao.GetChannelByAddress(ChannelIdentifier)
}

/*
DepositAndOpenChannel a channel with the peer at `partner_address`
    with the given `token_address`.
deposit必须大于0
settleTimeout: 如果为0 表示已经知道通道存在,只是为了存款,如果大于0,表示希望完全创建通道.
此接口并不等待交易打包才返回,因此如果是新创建通道,就算是成功了ch也会是nil
如果是单纯deposit,那么err为nil时,ch一定有效
*/
func (r *API) DepositAndOpenChannel(tokenAddress, partnerAddress common.Address, settleTimeout, revealTimeout int, deposit *big.Int, newChannel bool) (ch *channeltype.Serialization, err error) {
	if revealTimeout <= 0 {
		revealTimeout = r.Photon.Config.RevealTimeout
	}
	if newChannel {
		if settleTimeout <= 0 {
			settleTimeout = r.Photon.Config.SettleTimeout
		}
		if settleTimeout <= revealTimeout {
			err = rerr.ErrChannelInvalidSettleTimeout
			return
		}
		if bytes.Equal(partnerAddress[:], r.Photon.NodeAddress[:]) {
			err = rerr.ErrOpenChannelWithSelf
			return
		}
	} else {
		settleTimeout = 0
	}
	if deposit.Cmp(utils.BigInt0) <= 0 {
		err = rerr.ErrArgumentError.Append("invalid amount")
		return
	}
	if err = r.checkSmcStatus(); err != nil {
		return
	}
	if newChannel {
		_, err = r.Photon.dao.GetChannel(tokenAddress, partnerAddress)
		if err == nil {
			err = rerr.ErrChannelAlreadExist
			return
		}
		ch = channeltype.NewEmptySerialization()
		ch.ChannelIdentifier.ChannelIdentifier = utils.CalcChannelID(tokenAddress, r.Photon.Chain.GetRegistryAddress(), partnerAddress, r.Photon.NodeAddress)
		ch.TokenAddressBytes = tokenAddress[:]
		ch.OurAddress = r.Photon.NodeAddress
		ch.PartnerAddressBytes = partnerAddress[:]
		ch.State = channeltype.StateInValid
		ch.SettleTimeout = settleTimeout
		ch.RevealTimeout = revealTimeout
	} else {
		ch, err = r.Photon.dao.GetChannel(tokenAddress, partnerAddress)
		if err != nil {
			err = rerr.ErrChannelNotFound
			return
		}
		if ch.State != channeltype.StateOpened {
			err = rerr.ErrChannelState.Append(fmt.Sprintf("can not deposit to %s channel", ch.State))
			return
		}
	}
	result := r.Photon.depositAndOpenChannelClient(tokenAddress, partnerAddress, settleTimeout, deposit, newChannel)
	err = <-result.Result
	return
}

/*
TokenSwapAndWait Start an atomic swap operation by sending a MediatedTransfer with
    `maker_amount` of `maker_token` to `taker_address`. Only proceed when a
    new valid MediatedTransfer is received with `taker_amount` of
    `taker_token`.
*/
func (r *API) TokenSwapAndWait(lockSecretHash string, makerToken, takerToken, makerAddress, takerAddress common.Address,
	makerAmount, takerAmount *big.Int, secret string, routeInfo []pfsproxy.FindPathResponse) error {
	result, err := r.tokenSwapAsync(lockSecretHash, makerToken, takerToken, makerAddress, takerAddress,
		makerAmount, takerAmount, secret, routeInfo)
	if err != nil {
		return err
	}
	err = <-result.Result
	return err
}

func (r *API) tokenSwapAsync(lockSecretHash string, makerToken, takerToken, makerAddress, takerAddress common.Address,
	makerAmount, takerAmount *big.Int, secret string, routeInfo []pfsproxy.FindPathResponse) (result *utils.AsyncResult, err error) {
	chs, err := r.Photon.dao.GetChannelList(takerToken, utils.EmptyAddress)
	if err != nil || len(chs) == 0 {
		err = rerr.ErrTokenNotFound
		return
	}
	chs, err = r.Photon.dao.GetChannelList(makerToken, utils.EmptyAddress)
	if err != nil || len(chs) == 0 {
		err = rerr.ErrTokenNotFound
		return
	}

	tokenSwap := &TokenSwap{
		LockSecretHash:  common.HexToHash(lockSecretHash),
		Secret:          common.HexToHash(secret),
		FromToken:       makerToken,
		FromAmount:      new(big.Int).Set(makerAmount),
		FromNodeAddress: makerAddress,
		ToToken:         takerToken,
		ToAmount:        new(big.Int).Set(takerAmount),
		ToNodeAddress:   takerAddress,
		RouteInfo:       routeInfo,
	}
	result = r.Photon.tokenSwapMakerClient(tokenSwap)
	return
}

/*
ExpectTokenSwap Register an expected transfer for this node.

    If a MediatedMessage is received for the `maker_asset` with
    `maker_amount` then proceed to send a MediatedTransfer to
    `maker_address` for `taker_asset` with `taker_amount`.
*/
func (r *API) ExpectTokenSwap(lockSecretHash string, makerToken, takerToken, makerAddress, takerAddress common.Address,
	makerAmount, takerAmount *big.Int, routeInfo []pfsproxy.FindPathResponse) (err error) {
	chs, err := r.Photon.dao.GetChannelList(takerToken, utils.EmptyAddress)
	if err != nil || len(chs) == 0 {
		err = rerr.ErrTokenNotFound
		return
	}
	chs, err = r.Photon.dao.GetChannelList(makerToken, utils.EmptyAddress)
	if err != nil || len(chs) == 0 {
		err = rerr.ErrTokenNotFound
		return
	}
	tokenSwap := &TokenSwap{
		LockSecretHash:  common.HexToHash(lockSecretHash),
		FromToken:       makerToken,
		FromAmount:      new(big.Int).Set(makerAmount),
		FromNodeAddress: makerAddress,
		ToToken:         takerToken,
		ToAmount:        new(big.Int).Set(takerAmount),
		ToNodeAddress:   takerAddress,
		RouteInfo:       routeInfo,
	}
	r.Photon.tokenSwapTakerClient(tokenSwap)
	return nil
}

//GetNodeNetworkState Returns the currently network status of `node_address
func (r *API) GetNodeNetworkState(nodeAddress common.Address) (deviceType string, isOnline bool) {
	return r.Photon.Protocol.GetNetworkStatus(nodeAddress)
}

//StartHealthCheckFor Returns the currently network status of `node_address`.
func (r *API) StartHealthCheckFor(nodeAddress common.Address) (deviceType string, isOnline bool) {
	r.Photon.startHealthCheckFor(nodeAddress)
	return r.GetNodeNetworkState(nodeAddress)
}

//GetTokenList returns all available tokens
func (r *API) GetTokenList() (tokens []common.Address) {
	tokensmap, err := r.Photon.dao.GetAllTokens()
	if err != nil {
		log.Error(fmt.Sprintf("GetAllTokens err %s", err))
	}
	for k := range tokensmap {
		tokens = append(tokens, k)
	}
	return
}

//GetTokenTokenNetorks return all tokens and token networks
func (r *API) GetTokenTokenNetorks() (tokens []string) {
	tokenMap, err := r.Photon.dao.GetAllTokens()
	if err != nil {
		log.Error(fmt.Sprintf("GetAllTokens err %s", err))
	}
	for k := range tokenMap {
		tokens = append(tokens, k.String())
	}
	return
}

//Transfer transfer and wait
func (r *API) Transfer(token common.Address, amount *big.Int, target common.Address, secret common.Hash, timeout time.Duration, isDirectTransfer bool, data string, routeInfo []pfsproxy.FindPathResponse) (result *utils.AsyncResult, err error) {
	result, err = r.TransferInternal(token, amount, target, secret, isDirectTransfer, data, routeInfo)
	if err != nil {
		return
	}
	if timeout > 0 {
		timeoutCh := time.After(timeout)
		select {
		case <-timeoutCh:
			return result, rerr.ErrTransferTimeout
		case err = <-result.Result:
		}
	} else {
		err = <-result.Result
	}
	return result, err
}

// TransferAsync :
func (r *API) TransferAsync(tokenAddress common.Address, amount *big.Int, target common.Address, secret common.Hash, isDirectTransfer bool, data string, routeInfo []pfsproxy.FindPathResponse) (result *utils.AsyncResult, err error) {
	result, err = r.TransferInternal(tokenAddress, amount, target, secret, isDirectTransfer, data, routeInfo)
	if err != nil {
		return
	}
	timeoutCh := time.After(300 * time.Millisecond)
	select {
	case <-timeoutCh:
		return result, nil
	case err = <-result.Result:
	}
	return result, err
}

//TransferInternal :
func (r *API) TransferInternal(tokenAddress common.Address, amount *big.Int, target common.Address, secret common.Hash, isDirectTransfer bool, data string, routeInfo []pfsproxy.FindPathResponse) (result *utils.AsyncResult, err error) {
	log.Debug(fmt.Sprintf("initiating transfer initiator=%s target=%s token=%s amount=%d secret=%s,currentblock=%d",
		r.Photon.NodeAddress.String(), target.String(), tokenAddress.String(), amount, secret.String(), r.Photon.GetBlockNumber()))
	result = r.Photon.transferAsyncClient(tokenAddress, amount, target, secret, isDirectTransfer, data, routeInfo)
	return
}

// AllowRevealSecret :
// 1. find state manager by lockSecretHash and tokenAddress
// 2. check secret matches lockSecretHash or not
// 3. remove the predictor
func (r *API) AllowRevealSecret(lockSecretHash common.Hash, tokenAddress common.Address) (err error) {
	result := r.Photon.allowRevealSecretClient(lockSecretHash, tokenAddress)
	err = <-result.Result
	return
}

// RegisterSecret :
func (r *API) RegisterSecret(secret common.Hash, tokenAddress common.Address) (err error) {
	result := r.Photon.registerSecretClient(secret, tokenAddress)
	err = <-result.Result
	return
}

// TransferDataResponse :
type TransferDataResponse struct {
	Initiator      string   `json:"initiator_address"`
	Target         string   `json:"target_address"`
	Token          string   `json:"token_address"`
	Amount         *big.Int `json:"amount"`
	Secret         string   `json:"secret"`
	LockSecretHash string   `json:"lock_secret_hash"`
	Expiration     int64    `json:"expiration"`
	Fee            *big.Int `json:"fee"`
	IsDirect       bool     `json:"is_direct"`
}

// GetUnfinishedReceivedTransfer :
func (r *API) GetUnfinishedReceivedTransfer(lockSecretHash common.Hash, tokenAddress common.Address) (resp *TransferDataResponse) {
	result := r.Photon.getUnfinishedReceivedTransferClient(lockSecretHash, tokenAddress)
	err := <-result.Result
	if err != nil {
		log.Error(fmt.Sprintf("GetUnfinishedReceivedTransfer err %s", err))
		return nil
	}
	return result.Tag.(*TransferDataResponse)
}

//Close a channel opened with `partner_address` for the given `token_address`. return when state has been +d to database
func (r *API) Close(tokenAddress, partnerAddress common.Address) (c *channeltype.Serialization, err error) {
	if err = r.checkSmcStatus(); err != nil {
		return
	}
	c, err = r.Photon.dao.GetChannel(tokenAddress, partnerAddress)
	if err != nil {
		return
	}
	//send close channel request
	result := r.Photon.closeChannelClient(c.ChannelIdentifier.ChannelIdentifier)
	err = <-result.Result
	if err != nil {
		return
	}
	//reload data from database,
	return r.Photon.dao.GetChannelByAddress(c.ChannelIdentifier.ChannelIdentifier)
}

//Settle a closed channel with `partner_address` for the given `token_address`.return when state has been updated to database
func (r *API) Settle(tokenAddress, partnerAddress common.Address) (c *channeltype.Serialization, err error) {
	if err = r.checkSmcStatus(); err != nil {
		return
	}
	c, err = r.Photon.dao.GetChannel(tokenAddress, partnerAddress)
	if c.State == channeltype.StateOpened {
		err = rerr.InvalidState("channel is still open")
		return
	}
	//send settle request
	result := r.Photon.settleChannelClient(c.ChannelIdentifier.ChannelIdentifier)
	err = <-result.Result
	log.Info(fmt.Sprintf("%s settled finish , err %v", c.ChannelIdentifier, err))
	if err != nil {
		return
	}
	//reload data from database, this channel has been removed. 这时候channel应该是settling状态
	return r.Photon.dao.GetChannelByAddress(c.ChannelIdentifier.ChannelIdentifier)
}

//CooperativeSettle a channel opened with `partner_address` for the given `token_address`. return when state has been updated to database
func (r *API) CooperativeSettle(tokenAddress, partnerAddress common.Address) (c *channeltype.Serialization, err error) {
	if err = r.checkSmcStatus(); err != nil {
		return
	}
	c, err = r.Photon.dao.GetChannel(tokenAddress, partnerAddress)
	if c.State != channeltype.StateOpened && c.State != channeltype.StatePrepareForCooperativeSettle {
		err = rerr.InvalidState("channel must be  open")
		return
	}
	//send settle request
	result := r.Photon.cooperativeSettleChannelClient(c.ChannelIdentifier.ChannelIdentifier)
	err = <-result.Result
	log.Info(fmt.Sprintf("%s CooperativeSettle finish , err %v", c.ChannelIdentifier, err))
	if err != nil {
		return
	}
	//reload data from database, this channel has been removed.
	return r.Photon.dao.GetChannelByAddress(c.ChannelIdentifier.ChannelIdentifier)
}

//PrepareForCooperativeSettle  mark a channel prepared for settle,  return when state has been updated to database
func (r *API) PrepareForCooperativeSettle(channelIdentifier common.Hash) (c *channeltype.Serialization, err error) {
	c, err = r.Photon.dao.GetChannelByAddress(channelIdentifier)
	if err != nil {
		err = rerr.ChannelNotFound(channelIdentifier.String())
		return
	}
	if c.State != channeltype.StateOpened {
		err = rerr.InvalidState("channel must be  open")
		return
	}
	//send settle request
	result := r.Photon.markChannelForCooperativeSettleClient(channelIdentifier)
	err = <-result.Result
	log.Info(fmt.Sprintf("%s PrepareForCooperativeSettle finish , err %v", c.ChannelIdentifier, err))
	if err != nil {
		return
	}
	//reload data from database, this channel has been removed.
	return r.Photon.dao.GetChannelByAddress(channelIdentifier)
}

//CancelPrepareForCooperativeSettle  cancel a mark. return when state has been updated to database
func (r *API) CancelPrepareForCooperativeSettle(channelIdentifier common.Hash) (c *channeltype.Serialization, err error) {
	c, err = r.Photon.dao.GetChannelByAddress(channelIdentifier)
	if err != nil {
		err = rerr.ChannelNotFound(channelIdentifier.String())
		return
	}
	if c.State != channeltype.StatePrepareForCooperativeSettle {
		err = rerr.InvalidState("channel must be  prepareForCooperativeSettle")
		return
	}
	//send settle request
	result := r.Photon.cancelMarkChannelForCooperativeSettleClient(channelIdentifier)
	err = <-result.Result
	log.Info(fmt.Sprintf("%s CancelPrepareForCooperativeSettle finish , err %v", c.ChannelIdentifier, err))
	if err != nil {
		return
	}
	//reload data from database, this channel has been removed.
	return r.Photon.dao.GetChannelByAddress(channelIdentifier)
}

//Withdraw on a channel opened with `partner_address` for the given `token_address`. return when state has been updated to database
func (r *API) Withdraw(tokenAddress, partnerAddress common.Address, amount *big.Int) (c *channeltype.Serialization, err error) {
	if err = r.checkSmcStatus(); err != nil {
		return
	}
	c, err = r.Photon.dao.GetChannel(tokenAddress, partnerAddress)
	if c.State != channeltype.StateOpened && c.State != channeltype.StatePrepareForWithdraw {
		err = rerr.InvalidState("channel must be  open")
		return
	}
	if c.OurBalance().Cmp(amount) < 0 {
		err = rerr.ErrArgumentError.Printf("invalid withdraw amount, availabe=%s,want=%s", c.OurBalance(), amount)
		return
	}
	//send settle request
	result := r.Photon.withdrawClient(c.ChannelIdentifier.ChannelIdentifier, amount)
	err = <-result.Result
	log.Info(fmt.Sprintf("%s withdraw finish , err %v", c.ChannelIdentifier, err))
	if err != nil {
		return
	}
	//reload data from database, this channel has been removed.
	return r.Photon.dao.GetChannelByAddress(c.ChannelIdentifier.ChannelIdentifier)
}

//PrepareForWithdraw  mark a channel prepared for withdraw,  return when state has been updated to database
func (r *API) PrepareForWithdraw(tokenAddress, partnerAddress common.Address) (c *channeltype.Serialization, err error) {
	c, err = r.Photon.dao.GetChannel(tokenAddress, partnerAddress)
	if c.State != channeltype.StateOpened {
		err = rerr.InvalidState("channel must be  open")
		return
	}
	//send settle request
	result := r.Photon.markWithdrawClient(c.ChannelIdentifier.ChannelIdentifier)
	err = <-result.Result
	log.Info(fmt.Sprintf("%s PrepareForWithdraw finish , err %v", c.ChannelIdentifier, err))
	if err != nil {
		return
	}
	//reload data from database, this channel has been removed.
	return r.Photon.dao.GetChannelByAddress(c.ChannelIdentifier.ChannelIdentifier)
}

//CancelPrepareForWithdraw  cancel a mark. return when state has been updated to database
func (r *API) CancelPrepareForWithdraw(tokenAddress, partnerAddress common.Address) (c *channeltype.Serialization, err error) {
	c, err = r.Photon.dao.GetChannel(tokenAddress, partnerAddress)
	if c.State != channeltype.StatePrepareForWithdraw {
		err = rerr.InvalidState("channel must be  open")
		return
	}
	//send settle request
	result := r.Photon.cancelMarkWithdrawClient(c.ChannelIdentifier.ChannelIdentifier)
	err = <-result.Result
	log.Info(fmt.Sprintf("%s CancelPrepareForWithdraw finish , err %v", c.ChannelIdentifier, err))
	if err != nil {
		return
	}
	//reload data from database, this channel has been removed.
	return r.Photon.dao.GetChannelByAddress(c.ChannelIdentifier.ChannelIdentifier)
}

//GetTokenNetworkEvents return events about this token
func (r *API) GetTokenNetworkEvents(tokenAddress common.Address, fromBlock, toBlock int64) (data []interface{}, err error) {
	//type eventData struct {
	//	/*
	//				 {
	//		        "event_type": "ChannelNew",
	//		        "settle_timeout": 10,
	//		        "netting_channel": "0xc0ea08a2d404d3172d2add29a45be56da40e2949",
	//		        "participant1": "0x4894a542053248e0c504e3def2048c08f73e1ca6",
	//		        "participant2": "0x356857Cd22CBEFccDa4e96AF13b408623473237A"
	//		    }
	//	*/
	//	EventType      string `json:"event_type"`
	//	SettleTimeout  int    `json:"settle_timeout"`
	//	NettingChannel string `json:"netting_channel"`
	//	Participant1   string `json:"participant1"`
	//	Participant2   string `json:"participant2"`
	//	TokenAddress   string `json:"token_address"`
	//}
	//tokens, err := r.Photon.dao.GetAllTokens()
	//if err != nil {
	//	return
	//}
	//for t, manager := range tokens {
	//	if tokenAddress == utils.EmptyAddress || t == tokenAddress {
	//		events, err := r.Photon.BlockChainEvents.GetAllChannelManagerEvents(manager, fromBlock, toBlock)
	//		if err != nil {
	//			return nil, err
	//		}
	//		for _, e := range events {
	//			e2 := e.(*blockchain.EventChannelOpen)
	//			ed := &eventData{
	//				EventType:      e2.EventName,
	//				SettleTimeout:  e2.SettleTimeout,
	//				NettingChannel: e2.NettingChannelAddress.String(),
	//				Participant1:   e2.Participant1.String(),
	//				Participant2:   e2.Participant2.String(),
	//				TokenAddress:   t.String(),
	//			}
	//			data = append(data, ed)
	//		}
	//	}
	//}
	return
}

//GetNetworkEvents all photon events
func (r *API) GetNetworkEvents(fromBlock, toBlock int64) ([]interface{}, error) {
	//type eventData struct {
	//	/*
	//				 "event_type": "TokenAdded",
	//		        "token_address": "0xea674fdde714fd979de3edf0f56aa9716b898ec8",
	//		        "channel_manager_address": "0xc0ea08a2d404d3172d2add29a45be56da40e2949"
	//	*/
	//	EventType             string `json:"event_type"`
	//	TokenAddress          string `json:"token_address"`
	//	ChannelManagerAddress string `json:"channel_manager_address"`
	//}
	//events, err := r.Photon.BlockChainEvents.GetAllRegistryEvents(r.Photon.RegistryAddress, fromBlock, toBlock)
	//if err != nil {
	//	return nil, err
	//}
	//var data []interface{}
	//for _, e := range events {
	//	e2 := e.(*blockchain.EventTokenNetworkCreated)
	//	ed := &eventData{
	//		EventType:             e2.EventName,
	//		TokenAddress:          e2.TokenAddress.String(),
	//		ChannelManagerAddress: e2.TokenNetworkAddress.String(),
	//	}
	//	data = append(data, ed)
	//}
	return nil, nil
}

//GetChannelEvents events of this channel
func (r *API) GetChannelEvents(channelIdentifier common.Hash, fromBlock, toBlock int64) (data []transfer.Event, err error) {

	//var events []transfer.Event
	//events, err = r.Photon.BlockChainEvents.GetAllNettingChannelEvents(channelIdentifier, fromBlock, toBlock)
	//if err != nil {
	//	return
	//}
	//for _, e := range events {
	//	m := make(map[string]interface{})
	//	switch e2 := e.(type) {
	//	case *blockchain.EventChannelNewBalance:
	//		m["event_type"] = e2.EventName
	//		m["participant"] = e2.ParticipantAddress.String()
	//		m["balance"] = e2.Balance
	//		m["block_number"] = e2.BlockNumber
	//		data = append(data, m)
	//	case *blockchain.EventChannelClosed:
	//		m["event_type"] = e2.EventName
	//		m["netting_channel_address"] = e2.ContractAddress.String()
	//		m["closing_address"] = e2.ClosingAddress.String()
	//		data = append(data, m)
	//	case *blockchain.EventChannelSettled:
	//		m["event_type"] = e2.EventName
	//		m["netting_channel_address"] = e2.ContractAddress.String()
	//		m["block_number"] = e2.BlockNumber
	//		data = append(data, m)
	//	case *blockchain.EventSecretRevealed:
	//		m["event_type"] = e2.EventName
	//		m["netting_channel_address"] = e2.ContractAddress.String()
	//		m["secret"] = e2.Secret.String()
	//		data = append(data, m)
	//		//case *blockchain.EventNonClosingBalanceProofUpdated:
	//		//	m["event_type"] = e2.EventName
	//		//	m["token_address"] = t.String()
	//		//	m["channel_manager_address"] = graph.TokenAddress.String()
	//	}
	//
	//}
	//
	//var photonEvents []*models.InternalEvent
	//photonEvents, err = r.Photon.dao.GetEventsInBlockRange(fromBlock, toBlock)
	//if err != nil {
	//	return
	//}
	////Here choose which photon internal events we want to expose to the end user
	//for _, ev := range photonEvents {
	//	m := make(map[string]interface{})
	//	switch e2 := ev.EventObject.(type) {
	//	case *transfer.EventTransferSentSuccess:
	//		m["event_type"] = "EventTransferSentSuccess"
	//		m["identifier"] = e2.LockSecretHash
	//		m["block_number"] = ev.BlockNumber
	//		m["amount"] = e2.Amount
	//		m["target"] = e2.Target
	//		data = append(data, m)
	//	case *transfer.EventTransferSentFailed:
	//		m["event_type"] = "EventTransferSentFailed"
	//		m["identifier"] = e2.LockSecretHash
	//		m["block_number"] = ev.BlockNumber
	//		m["reason"] = e2.Reason
	//		data = append(data, m)
	//	case *transfer.EventTransferReceivedSuccess:
	//		m["event_type"] = "EventTransferReceivedSuccess"
	//		m["identifier"] = e2.LockSecretHash
	//		m["block_number"] = ev.BlockNumber
	//		m["amount"] = e2.Amount
	//		m["initiator"] = e2.Initiator.String()
	//		data = append(data, m)
	//	}
	//}
	return
}

/*
GetSentTransferDetails query sent transfers from dao
*/
func (r *API) GetSentTransferDetails(tokenAddress common.Address, from, to int64) ([]*models.SentTransferDetail, error) {
	return r.Photon.dao.GetSentTransferDetailList(tokenAddress, -1, -1, from, to)
}

/*
GetReceivedTransfers query received transfers from dao
*/
func (r *API) GetReceivedTransfers(tokenAddress common.Address, fromBlock, toBlock, fromTime, toTime int64) ([]*models.ReceivedTransfer, error) {
	return r.Photon.dao.GetReceivedTransferList(tokenAddress, fromBlock, toBlock, fromTime, toTime)
}

//Stop stop for mobile app
func (r *API) Stop() {
	log.Info("calling api stop..")
	r.Photon.Stop()
	log.Info("stop successful..")
}

/*
ChannelInformationFor3rdParty generate all information need by 3rd party
*/
func (r *API) ChannelInformationFor3rdParty(channelIdentifier common.Hash, thirdAddr common.Address) (result *pmsproxy.DelegateForPms, err error) {
	c, err := r.Photon.dao.GetChannelByAddress(channelIdentifier)
	if err != nil {
		return
	}
	return r.Photon.GetDelegateForPms(c, thirdAddr)
}

//EventTransferSentSuccessWrapper wrapper
type EventTransferSentSuccessWrapper struct {
	transfer.EventTransferSentSuccess
	BlockNumber int64
	Name        string
}

//EventTransferSentFailedWrapper wrapper
type EventTransferSentFailedWrapper struct {
	transfer.EventTransferSentFailed
	BlockNumber int64
	Name        string
}

//EventEventTransferReceivedSuccessWrapper wrapper
type EventEventTransferReceivedSuccessWrapper struct {
	transfer.EventTransferReceivedSuccess
	BlockNumber int64
	Name        string
}

// AccountTokenBalanceVo for api
type AccountTokenBalanceVo struct {
	TokenAddress string   `json:"token_address"`
	Balance      *big.Int `json:"balance"`
	LockedAmount *big.Int `json:"locked_amount"`
}

// GetBalanceByTokenAddress : get account's balance and locked account on token
func (r *API) GetBalanceByTokenAddress(tokenAddress common.Address) (balances []*AccountTokenBalanceVo, err error) {
	if tokenAddress == utils.EmptyAddress {
		return r.getBalance()
	}
	tokens := r.GetTokenList()
	hasRegistered := false
	for _, token := range tokens {
		if token == tokenAddress {
			hasRegistered = true
		}
	}
	if !hasRegistered {
		err = rerr.ErrTokenNotFound
		return
	}
	channels, err := r.GetChannelList(tokenAddress, utils.EmptyAddress)
	if err != nil {
		return
	}
	balance := new(AccountTokenBalanceVo)
	balance.TokenAddress = tokenAddress.String()
	balance.Balance = big.NewInt(0)
	balance.LockedAmount = big.NewInt(0)
	for _, channel := range channels {
		balance.Balance.Add(balance.Balance, channel.OurBalance())
		balance.LockedAmount.Add(balance.LockedAmount, channel.OurAmountLocked())
	}
	return []*AccountTokenBalanceVo{balance}, err
}

// getBalance : get account's balance and locked account on each token
func (r *API) getBalance() (balances []*AccountTokenBalanceVo, err error) {
	channels, err := r.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	if err != nil {
		return
	}
	token2ChannelMap := make(map[common.Address][]*channeltype.Serialization)
	for _, channel := range channels {
		token2ChannelMap[channel.TokenAddress()] = append(token2ChannelMap[channel.TokenAddress()], channel)
	}
	for tokenAddress, channels := range token2ChannelMap {
		balance := &AccountTokenBalanceVo{
			TokenAddress: tokenAddress.String(),
			Balance:      big.NewInt(0),
			LockedAmount: big.NewInt(0),
		}
		for _, channel := range channels {
			balance.Balance.Add(balance.Balance, channel.OurBalance())
			balance.LockedAmount.Add(balance.LockedAmount, channel.OurAmountLocked())
		}
		balances = append(balances, balance)
	}
	return
}

// ForceUnlock : only for debug
func (r *API) ForceUnlock(channelIdentifier, secret common.Hash) (err error) {
	result := r.Photon.forceUnlockClient(secret, channelIdentifier)
	err = <-result.Result
	return
}

// RegisterSecretOnChain : only for debug
func (r *API) RegisterSecretOnChain(secret common.Hash) (err error) {
	result := r.Photon.registerSecretOnChainClient(secret)
	err = <-result.Result
	return
}

// CancelTransfer : cancel a transfer when haven't send secret
func (r *API) CancelTransfer(lockSecretHash common.Hash, tokenAddress common.Address) error {
	result := r.Photon.cancelTransferClient(lockSecretHash, tokenAddress)
	return <-result.Result
}

type balanceProof struct {
	Nonce             uint64      `json:"nonce"`
	TransferAmount    *big.Int    `json:"transfer_amount"`
	LocksRoot         common.Hash `json:"locks_root"`
	ChannelIdentifier common.Hash `json:"channel_identifier"`
	OpenBlockNumber   int64       `json:"open_block_number"`
	MessageHash       common.Hash `json:"addition_hash"`
	//signature is nonce + transferred_amount + locksroot + channel_identifier + message_hash
	Signature []byte `json:"signature"`
}

//ProofForPFS proof for path finding service, test only
type ProofForPFS struct {
	BalanceProof balanceProof `json:"balance_proof"`
	Signature    []byte       `json:"balance_signature"`
	LockAmount   *big.Int     `json:"lock_amount"`
}

//BalanceProofForPFS proof for path finding service ,test only
func (r *API) BalanceProofForPFS(channelIdentifier common.Hash) (proof *ProofForPFS, err error) {
	ch, err := r.GetChannel(channelIdentifier)
	if err != nil {
		return
	}
	proof = &ProofForPFS{
		BalanceProof: balanceProof{
			Nonce:             ch.PartnerBalanceProof.Nonce,
			TransferAmount:    ch.PartnerBalanceProof.TransferAmount,
			LocksRoot:         ch.PartnerBalanceProof.LocksRoot,
			ChannelIdentifier: ch.ChannelIdentifier.ChannelIdentifier,
			OpenBlockNumber:   ch.ChannelIdentifier.OpenBlockNumber,
			MessageHash:       ch.PartnerBalanceProof.MessageHash,
			Signature:         ch.PartnerBalanceProof.Signature,
		},
		LockAmount: ch.PartnerAmountLocked(),
	}
	bpf := &proof.BalanceProof
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, bpf.Nonce)
	_, err = buf.Write(utils.BigIntTo32Bytes(bpf.TransferAmount))
	_, err = buf.Write(bpf.LocksRoot[:])
	_, err = buf.Write(bpf.ChannelIdentifier[:])
	err = binary.Write(buf, binary.BigEndian, bpf.OpenBlockNumber)
	_, err = buf.Write(bpf.MessageHash[:])
	_, err = buf.Write(bpf.Signature)
	_, err = buf.Write(utils.BigIntTo32Bytes(proof.LockAmount))
	dataToSign := buf.Bytes()
	proof.Signature, err = utils.SignData(r.Photon.PrivateKey, dataToSign)
	return
}

// NotifyNetworkDown :
func (r *API) NotifyNetworkDown() error {
	log.Info(fmt.Sprintf("NotifyNetworkDown from user"))
	// smc client
	client := r.Photon.Chain.Client
	if client.IsConnected() {
		//r.Photon.BlockChainEvents.Stop()
		client.Client.Close()
	}

	// xmpp client
	if t, ok := r.Photon.Protocol.Transport.(*network.MixTransport); ok {
		t.Reconnect()
	}
	// 置为无效公链状态
	//r.Photon.IsChainEffective = false
	//log.Info("photon works without effective chain now because user call NotifyNetworkDown...")
	return nil
}

// GetFeePolicy 如果没有启动收费会返回错误,否则返回当前账户设置的收费信息
func (r *API) GetFeePolicy() (fp *models.FeePolicy, err error) {
	feeModule, ok := r.Photon.FeePolicy.(*FeeModule)
	if !ok {
		err = rerr.ErrNotChargeFee
		return
	}
	return feeModule.feePolicy, nil
}

// SetFeePolicy 更新该账户所有的收费信息,不保留历史记录
func (r *API) SetFeePolicy(fp *models.FeePolicy) error {
	feeModule, ok := r.Photon.FeePolicy.(*FeeModule)
	if !ok {
		return rerr.ErrNotChargeFee.Append("photon start without param '--fee', can not set fee policy")
	}
	return feeModule.SetFeePolicy(fp)
}

// FindPath 向PFS询问路由,要求启用收费
func (r *API) FindPath(targetAddress, tokenAddress common.Address, amount *big.Int) (routes []pfsproxy.FindPathResponse, err error) {
	if r.Photon.PfsProxy == nil {
		err = rerr.ErrNotChargeFee.Append("photon start without param '--pfs', can not calculate total fee")
		return
	}
	routes, err = r.Photon.PfsProxy.FindPath(r.Photon.NodeAddress, targetAddress, tokenAddress, amount, true)
	if err != nil {
		return
	}
	return
}

// GetAllFeeChargeRecord :
func (r *API) GetAllFeeChargeRecord() (resp interface{}, err error) {
	type responce struct {
		TotalFee map[common.Address]*big.Int `json:"total_fee"`
		Details  []*models.FeeChargeRecord   `json:"details"`
	}
	var data responce

	data.Details, err = r.Photon.dao.GetAllFeeChargeRecord(utils.EmptyAddress, -1, -1)
	if err != nil {
		return
	}
	data.TotalFee = make(map[common.Address]*big.Int)
	for _, record := range data.Details {
		totalFee := data.TotalFee[record.TokenAddress]
		if totalFee == nil {
			totalFee = big.NewInt(0)
		}
		data.TotalFee[record.TokenAddress] = totalFee.Add(totalFee, record.Fee)
	}
	resp = data
	return
}

// SystemStatus :
func (r *API) SystemStatus() (resp interface{}, err error) {
	type transfers struct {
		SendNum    int `json:"send_num"`
		ReceiveNum int `json:"receive_num"`
		DealingNum int `json:"dealing_num"`
	}
	type systemStatus struct {
		EthRPCEndpoint      string                            `json:"eth_rpc_endpoint"`
		EthRPCStatus        string                            `json:"eth_rpc_status"` // disconnected, connected, closed, reconnecting
		NodeAddress         string                            `json:"node_address"`
		RegistryAddress     string                            `json:"registry_address"`
		TokenToTokenNetwork map[common.Address]common.Address `json:"token_to_token_network"`
		LastBlockNumber     int64                             `json:"block_number"`
		LastBlockNumberTime time.Time                         `json:"last_block_number_time"`
		IsMobileMode        bool                              `json:"is_mobile_mode"`
		NetworkType         string                            `json:"network_type"` // xmpp, xmpp-udp, matrix, matrix-udp,udp
		FeePolicy           *models.FeePolicy                 `json:"fee_policy"`
		ChannelNum          int                               `json:"channel_num"`
		Transfers           *transfers                        `json:"transfers,omitempty"`
		SyncProcess         *ethereum.SyncProgress            `json:"sync_process"`
	}
	var data systemStatus
	data.EthRPCEndpoint = r.Photon.Config.EthRPCEndPoint
	// EthRPCStatus
	// 这里只向外暴露两个状态:有效公链/无效公链
	if r.Photon.IsChainEffective {
		data.EthRPCStatus = "valid"
	} else {
		data.EthRPCStatus = "invalid"
	}
	//switch r.Photon.Chain.Client.Status {
	//case netshare.Disconnected:
	//	data.EthRPCStatus = "disconnected"
	//case netshare.Connected:
	//	data.EthRPCStatus = "connected"
	//	if !r.Photon.IsChainEffective {
	//		data.EthRPCStatus = "connected_but_invalid"
	//	}
	//case netshare.Closed:
	//	data.EthRPCStatus = "closed"
	//case netshare.Reconnecting:
	//	data.EthRPCStatus = "reconnecting"
	//}
	data.NodeAddress = r.Photon.NodeAddress.String()
	data.RegistryAddress = r.Photon.Chain.GetRegistryAddress().String()
	// TokenToTokenNetwork
	data.TokenToTokenNetwork = r.Photon.Token2TokenNetwork
	data.LastBlockNumber = r.Photon.dao.GetLatestBlockNumber()
	data.LastBlockNumberTime = r.Photon.dao.GetLastBlockNumberTime()
	data.IsMobileMode = params.MobileMode
	// network type
	switch r.Photon.Transport.(type) {
	case *network.XMPPTransport:
		data.NetworkType = "xmpp"
	case *network.MixTransport:
		data.NetworkType = "xmpp-udp"
	case *network.MatrixTransport:
		data.NetworkType = "matrix"
	case *network.MatrixMixTransport:
		data.NetworkType = "matrix-udp"
	case *network.UDPTransport:
		data.NetworkType = "udp"
	}
	// FeePolicy
	if r.Photon.Config.EnableMediationFee {
		data.FeePolicy = r.Photon.dao.GetFeePolicy()
	} else {
		data.FeePolicy = nil
	}
	// channel num
	cs, err := r.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	if err != nil {
		return
	}
	data.ChannelNum = len(cs)
	// Transfers
	sts, err := r.GetSentTransferDetails(utils.EmptyAddress, -1, -1)
	if err != nil {
		return
	}
	rts, err := r.GetReceivedTransfers(utils.EmptyAddress, -1, -1, -1, -1)
	if err != nil {
		return
	}
	data.Transfers = &transfers{
		SendNum:    len(sts),
		ReceiveNum: len(rts),
		DealingNum: len(r.Photon.Transfer2StateManager),
	}
	//只有确定公链链接有效才会查询,否则不查询
	if r.Photon.Chain.Client.Status == netshare.Connected {
		var sync *ethereum.SyncProgress
		sync, err = r.Photon.Chain.SyncProgress()
		if err != nil {
			err = rerr.ErrUnkownSpectrumRPCError.Append(err.Error())
			return
		}
		data.SyncProcess = sync
	}
	resp = data
	return
}

func (r *API) checkSmcStatus() error {
	var err error
	// 1. 校验最新块的时间
	lastBlockNumberTime := r.Photon.dao.GetLastBlockNumberTime()
	if time.Since(lastBlockNumberTime) > 60*time.Second {
		err = rerr.ErrSpectrumSyncError.Errorf("has't receive new block from smc since %s, maybe something wrong with smc", lastBlockNumberTime.String())
		log.Error(err.Error())
		return err
	}
	// 2. 校验smc节点同步情况
	sp, err := r.Photon.Chain.SyncProgress()
	if err != nil {
		err = rerr.ErrSpectrumSyncError.Errorf("call smc SyncProgress err %s", err)
		log.Error(err.Error())
		return err
	}
	var defaultSyncBlock uint64 = 3
	if params.ChainID.Uint64() == 8888 { //测试私链一秒一块,时间太短是不行的
		defaultSyncBlock = 30
	}
	if sp != nil && sp.HighestBlock-sp.CurrentBlock > defaultSyncBlock {
		err = rerr.ErrSpectrumBlockError.Errorf("smc block number error : HighestBlock=%d but CurrentBlock=%d", sp.HighestBlock, sp.CurrentBlock)
		log.Error(err.Error())
		return err
	}
	return nil
}

// ContractCallTXQueryParams 请求参数
type ContractCallTXQueryParams struct {
	ChannelIdentifier string              `json:"channel_identifier"`
	OpenBlockNumber   int64               `json:"open_block_number"`
	TokenAddress      string              `json:"token_address"`
	TXType            models.TXInfoType   `json:"tx_type"`
	TXStatus          models.TXInfoStatus `json:"tx_status"`
}

// ContractCallTXQuery 根据条件查询所有合约调用的信息
func (r *API) ContractCallTXQuery(req *ContractCallTXQueryParams) (list []*models.TXInfo, err error) {
	channelIdentifier := utils.EmptyHash
	openBlockNumber := int64(0)
	var txType models.TXInfoType
	var txStatus models.TXInfoStatus
	tokenAddress := utils.EmptyAddress
	if req.ChannelIdentifier != "" {
		channelIdentifier = common.HexToHash(req.ChannelIdentifier)
	}
	if req.OpenBlockNumber > 0 {
		openBlockNumber = req.OpenBlockNumber
	}
	if req.TokenAddress != "" {
		tokenAddress = common.HexToAddress(req.TokenAddress)
	}
	if req.TXType != "" {
		txType = req.TXType
	}
	if req.TXStatus != "" {
		txStatus = req.TXStatus
	}
	list, err = r.Photon.dao.GetTXInfoList(channelIdentifier, openBlockNumber, tokenAddress, txType, txStatus)
	return
}

// 手续类型常量
const incomeTypeTransfer = "0" // 转账收益
const incomeTypeFee = "1"      // 手续费收益

// IncomeDetail 收益接口返回的收益明细信息
type IncomeDetail struct {
	Amount    *big.Int `json:"amount"`
	Data      string   `json:"data"`
	Type      string   `json:"type"` // 0=转账收益 1-手续费收益
	TimeStamp int64    `json:"time_stamp"`
}

type incomeDetailSorter []*IncomeDetail

func (s incomeDetailSorter) Len() int {
	return len(s)
}

func (s incomeDetailSorter) Less(i, j int) bool {
	return s[i].TimeStamp < s[j].TimeStamp
}
func (s incomeDetailSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// GetIncomeDetails 收益明细查询接口,这里的收益包含手续费收益和收到的data不为""的交易
func (r *API) GetIncomeDetails(tokenAddress common.Address, fromTime, toTime int64, limit int) (list []*IncomeDetail, err error) {
	if tokenAddress == utils.EmptyAddress {
		err = rerr.ErrGeneralDBError.Append("token address can not be empty")
		return
	}
	receivedTransfers, err := r.Photon.dao.GetReceivedTransferList(tokenAddress, -1, -1, fromTime, toTime)
	if err != nil {
		err = rerr.ErrGeneralDBError.Append(err.Error())
		return
	}
	feeRecords, err := r.Photon.dao.GetAllFeeChargeRecord(tokenAddress, fromTime, toTime)
	if err != nil {
		err = rerr.ErrGeneralDBError.Append(err.Error())
		return
	}
	for _, rt := range receivedTransfers {
		if rt.Data == "" {
			// 过滤data为""的
			continue
		}
		list = append(list, &IncomeDetail{
			Amount:    rt.Amount,
			Data:      rt.Data,
			Type:      incomeTypeTransfer,
			TimeStamp: rt.TimeStamp,
		})
	}
	for _, f := range feeRecords {
		list = append(list, &IncomeDetail{
			Amount:    f.Fee,
			Data:      f.Data,
			Type:      incomeTypeFee,
			TimeStamp: f.Timestamp,
		})
	}
	// 排序,并裁剪至limit
	sort.Stable(incomeDetailSorter(list))
	if limit > 0 && len(list) > limit {
		list = list[:limit]
	}
	return
}

// OneDayIncome 一天的收益统计
type OneDayIncome struct {
	Amount    *big.Int `json:"amount"`
	TimeStamp int64    `json:"time_stamp"`
}

// DaysIncome 一周的收益统计
type DaysIncome struct {
	TokenAddress common.Address  `json:"token_address"`
	TotalAmount  *big.Int        `json:"total_amount"`
	Days         int             `json:"days"`
	Details      []*OneDayIncome `json:"details"`
}

// GetDaysIncome 获取过去n天的收益统计
func (r *API) GetDaysIncome(tokenAddress common.Address, n int) (resp []*DaysIncome, err error) {
	if n <= 0 {
		// 默认7天
		n = 7
	}
	var tokenList []common.Address
	if tokenAddress != utils.EmptyAddress {
		tokenList = append(tokenList, tokenAddress)
	} else {
		tokenList = r.GetTokenList()
	}
	if len(tokenList) == 0 {
		return
	}
	// 起始时间
	now := time.Now()
	t := now.AddDate(0, 0, 0-n)
	fromTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	toTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	/*
		构造每个token对应的DaysIncome
	*/
	for _, token := range tokenList {
		incomeDetailList, err2 := r.GetIncomeDetails(token, -1, -1, -1)
		if err2 != nil {
			err = err2
			return
		}
		r := &DaysIncome{
			TokenAddress: token,
			TotalAmount:  big.NewInt(0),
			Days:         n,
		}
		// 统计总收益
		var detailList []*IncomeDetail
		for _, incomeDetail := range incomeDetailList {
			r.TotalAmount = r.TotalAmount.Add(r.TotalAmount, incomeDetail.Amount)
			if incomeDetail.TimeStamp >= fromTime.Unix() && incomeDetail.TimeStamp < toTime.Unix() {
				detailList = append(detailList, incomeDetail)
			}
		}
		// 统计过去N天中每天的收益,没收益的天补0
		for i := 0; i < n; i++ {
			begin := fromTime.AddDate(0, 0, i)
			end := begin.AddDate(0, 0, 1)
			day := &OneDayIncome{
				Amount:    big.NewInt(0),
				TimeStamp: begin.Unix(),
			}
			for _, detail := range detailList {
				if detail.TimeStamp >= begin.Unix() && detail.TimeStamp < end.Unix() {
					day.Amount = day.Amount.Add(day.Amount, detail.Amount)
				}
			}
			r.Details = append(r.Details, day)
		}
		resp = append(resp, r)
	}
	return
}

// GetBuildInfo 获取当前版本信息
func (r *API) GetBuildInfo() *BuildInfo {
	return r.Photon.BuildInfo
}

// GetChannelSettleBlockResponse :
type GetChannelSettleBlockResponse struct {
	BlockNumberNow              int64 `json:"block_number_now"`
	BlockNumberChannelCanSettle int64 `json:"block_number_channel_can_settle,omitempty"`
}

// GetChannelSettleBlock :
func (r *API) GetChannelSettleBlock(channelIdentifier common.Hash) *GetChannelSettleBlockResponse {
	resp := &GetChannelSettleBlockResponse{}
	// 1. 获取当前块
	resp.BlockNumberNow = r.Photon.GetBlockNumber()
	// 2. 获取通道SettleBlock
	c := r.Photon.getChannelWithAddr(channelIdentifier)
	if c == nil || c.State != channeltype.StateClosed {
		return resp
	}
	resp.BlockNumberChannelCanSettle = c.ExternState.ClosedBlock + int64(c.SettleTimeout) + int64(params.ContractPunishBlockNumber)
	return resp
}

// GetTokenBalance 获取账户在token上的余额
func (r *API) GetTokenBalance(account, token common.Address) (*big.Int, error) {
	t, err := r.Photon.Chain.Token(token)
	if err != nil {
		return nil, rerr.ErrArgumentError.AppendError(err)
	}
	name, err := t.Token.Name(nil)
	if err != nil {
		log.Error(err.Error())
	}
	if name == params.SMTTokenName {
		v, err2 := r.Photon.Chain.Client.BalanceAt(context.Background(), account, big.NewInt(r.Photon.GetBlockNumber()))
		if err2 != nil {
			return nil, rerr.ErrArgumentError.AppendError(err2)
		}
		return v, nil
	}
	v, err := t.BalanceOf(account)
	if err != nil {
		return nil, rerr.ErrArgumentError.AppendError(err)
	}
	return v, nil
}

// GetAssetsOnTokenResponseDetail :
type GetAssetsOnTokenResponseDetail struct {
	TokenAddress    string   `json:"token_address"`
	BalanceOnChain  *big.Int `json:"balance_on_chain"`
	BalanceInPhoton *big.Int `json:"balance_in_photon"`
}

// GetAssetsOnToken :
func (r *API) GetAssetsOnToken(tokenList []common.Address) (resp []*GetAssetsOnTokenResponseDetail, err error) {
	if len(tokenList) == 0 {
		tokenList = r.GetTokenList()
	}
	for _, token := range tokenList {
		d := &GetAssetsOnTokenResponseDetail{
			TokenAddress:    token.String(),
			BalanceOnChain:  big.NewInt(0),
			BalanceInPhoton: big.NewInt(0),
		}
		// 1. 获取用户在链上的该token余额
		t, err2 := r.GetTokenBalance(r.Photon.NodeAddress, token)
		if err2 != nil {
			log.Error(err2.Error())
			err = err2
			return
		}
		d.BalanceOnChain = t
		// 2. 获取用户在photon中的该token余额
		balance, err2 := r.GetBalanceByTokenAddress(token)
		if err2 != nil {
			log.Error(err2.Error())
			err = err2
			return
		}
		if len(balance) == 1 {
			d.BalanceInPhoton = balance[0].Balance
		}
		if d.BalanceOnChain.Cmp(utils.BigInt0) > 0 || d.BalanceInPhoton.Cmp(utils.BigInt0) > 0 {
			resp = append(resp, d)
		}
	}
	return
}

// UploadLogFile 上传photon日志到日志server
func (r *API) UploadLogFile() error {
	return uploadLogFile(r.Photon.NodeAddress, r.Photon.Config.LogFilePath)
}

func uploadLogFile(address common.Address, logFilePath string) (err error) {
	now := time.Now()
	logFileName := fmt.Sprintf("photon-log-%s-%s.log", address.String(), now.Format("2006-01-02-15-04-05"))
	fmt.Println(logFileName)
	// 1. 读取文件并压缩 #nosec
	logFile, err := os.Open(logFilePath)
	if err != nil {
		return
	}
	defer logFile.Close()
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	fileStat, err := logFile.Stat()
	if err != nil {
		return
	}
	zw.Name = logFileName
	zw.ModTime = fileStat.ModTime()
	_, err = io.Copy(zw, logFile)
	if err != nil {
		return
	}
	err = zw.Flush()
	err = zw.Close()
	if err != nil {
		return
	}
	// 2. 发送
	var buf2 bytes.Buffer
	w := multipart.NewWriter(&buf2)
	fw, err := w.CreateFormFile("uploadfile", logFileName)
	if err != nil {
		return
	}
	_, err = io.Copy(fw, &buf)
	if err != nil {
		return
	}
	err = w.Close()
	if err != nil {
		return
	}
	res, err := http.Post(params.TestLogServer+"/logsrv/1/upload", w.FormDataContentType(), &buf2)
	if err != nil {
		return
	}
	defer res.Body.Close()
	message, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(err.Error())
	}
	log.Info(fmt.Sprintf("upload logfile and got response : %d %s", res.StatusCode, string(message)))
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("upload logfile and got response : %d %s", res.StatusCode, string(message))
	}
	return
}
