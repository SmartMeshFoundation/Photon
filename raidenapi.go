package smartraiden

import (
	"time"

	"fmt"

	"math/big"

	"sync"

	"errors"

	"bytes"
	"encoding/binary"

	"crypto/ecdsa"

	"github.com/SmartMeshFoundation/SmartRaiden/blockchainold"
	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/models"
	"github.com/SmartMeshFoundation/SmartRaiden/network"
	"github.com/SmartMeshFoundation/SmartRaiden/rerr"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

//RaidenAPI raiden for user
type RaidenAPI struct {
	Raiden *RaidenService
}

//NewRaidenAPI create CLI interface.
func NewRaidenAPI(raiden *RaidenService) *RaidenAPI {
	return &RaidenAPI{Raiden: raiden}
}

//Address return this node's address
func (r *RaidenAPI) Address() common.Address {
	return r.Raiden.NodeAddress
}

//Tokens Return a list of the tokens registered with the default registry.
func (r *RaidenAPI) Tokens() (addresses []common.Address) {
	addresses, err := r.Raiden.Registry.TokenAddresses()
	if err == nil {
		return
	}
	tm, err := r.Raiden.db.GetAllTokens()
	if err != nil {
		return
	}
	for t := range tm {
		addresses = append(addresses, t)
	}
	return addresses
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
func (r *RaidenAPI) GetChannelList(tokenAddress common.Address, partnerAddress common.Address) (cs []*channel.Serialization, err error) {
	return r.Raiden.db.GetChannelList(tokenAddress, partnerAddress)
}

//GetChannel get channel by address
func (r *RaidenAPI) GetChannel(channelAddress common.Address) (c *channel.Serialization, err error) {
	return r.Raiden.db.GetChannelByAddress(channelAddress)
}

/*
ManagerAddressIfTokenRegistered return the channel manager address,If the token is registered then
Also make sure that the channel manager is registered with the node.
*/
func (r *RaidenAPI) ManagerAddressIfTokenRegistered(tokenAddress common.Address) (mgrAddr common.Address, err error) {
	mgrAddr, err = r.Raiden.Registry.ChannelManagerByToken(tokenAddress)
	if err != nil {
		return
	}
	return
}

/*
RegisterToken Will register the token at `token_address` with raiden. If it's already
    registered, will throw an exception.
*/
func (r *RaidenAPI) RegisterToken(tokenAddress common.Address) (mgrAddr common.Address, err error) {
	mgrAddr, err = r.Raiden.Registry.ChannelManagerByToken(tokenAddress)
	if err == nil && mgrAddr != utils.EmptyAddress {
		err = errors.New("TokenNetworkAddres already registered")
		return
	}
	//for non exist tokenaddress, ChannelManagerByToken will return a error: `abi : unmarshalling empty output`
	if err == rerr.ErrNoTokenManager {
		return r.Raiden.Registry.AddToken(tokenAddress)
	}
	return
}

/*
ConnectTokenNetwork Instruct the ConnectionManager to establish and maintain a connection to the token
    network.

    If the `token_address` is not already part of the raiden network, this will also register
    the token. //

    Args:
        token_address (bin): the ERC20 token network to connect to.
        funds (int): the amount of funds that can be used by the ConnectionMananger.
        initial_channel_target (int): number of channels to open proactively.
        joinable_funds_target (float): fraction of the funds that will be used to join
            channels opened by other participants.
*/
func (r *RaidenAPI) ConnectTokenNetwork(tokenAddress common.Address, funds *big.Int, initialChannelTarget int64, joinableFundsTarget float64) error {
	cm, err := r.Raiden.connectionManagerForToken(tokenAddress)
	if err != nil {
		return err
	}
	return cm.Connect(funds, initialChannelTarget, joinableFundsTarget)
}

/*
LeaveTokenNetwork Instruct the ConnectionManager to close all channels and wait for
    settlement
*/
func (r *RaidenAPI) LeaveTokenNetwork(tokenAddress common.Address, onlyReceiving bool) ([]*channel.Serialization, error) {
	cm, err := r.Raiden.connectionManagerForToken(tokenAddress)
	if err != nil {
		return nil, err
	}
	chs, err := cm.Leave(onlyReceiving)
	return chs, err
}

/*
GetConnectionManagersInfo Get a dict whose keys are token addresses and whose values are
    open channels, funds of last request, sum of deposits and number of channels
*/
func (r *RaidenAPI) GetConnectionManagersInfo() map[string]interface{} {
	type info struct {
		Funds       *big.Int `json:"funds"`
		SumDeposits *big.Int `json:"sum_deposits"`
		Channels    int      `json:"channels"`
	}
	infos := make(map[string]interface{})
	for _, t := range r.GetTokenList() {
		cm, err := r.Raiden.connectionManagerForToken(t)
		if err != nil {
			continue
		}
		if len(cm.openChannels()) > 0 {
			info := &info{
				Funds:       cm.funds,
				SumDeposits: cm.sumDeposits(),
				Channels:    len(cm.openChannels()),
			}
			infos[cm.tokenAddress.String()] = info
		}
	}
	return infos
}

/*
Open a channel with the peer at `partner_address`
    with the given `token_address`.
*/
func (r *RaidenAPI) Open(tokenAddress, partnerAddress common.Address, settleTimeout, revealTimeout int) (ch *channel.Serialization, err error) {
	if revealTimeout <= 0 {
		revealTimeout = r.Raiden.Config.RevealTimeout
	}
	if settleTimeout <= 0 {
		settleTimeout = r.Raiden.Config.SettleTimeout
	}
	if settleTimeout <= revealTimeout {
		err = rerr.ErrInvalidSettleTimeout
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	r.Raiden.db.RegisterNewChannellCallback(func(c *channel.Serialization) (remove bool) {
		if c.TokenAddress == tokenAddress && c.PartnerAddress == partnerAddress {
			wg.Done()
			return true
		}
		return false
	})
	result := r.Raiden.newChannelClient(tokenAddress, partnerAddress, settleTimeout)
	err = <-result.Result
	if err != nil {
		return
	}
	//wait
	wg.Wait()
	return r.Raiden.db.GetChannel(tokenAddress, partnerAddress)
}

/*
Deposit `amount` in the channel with the peer at `partner_address` and the
    given `token_address` in order to be able to do transfers.

    Raises:
        InvalidAddress: If either token_address or partner_address is not
        20 bytes long.
        TransactionThrew: May happen for multiple reasons:
            - If the token approval fails, e.g. the token may validate if
              account has enough balance for the allowance.
            - The deposit failed, e.g. the allowance did not set the token
              aside for use and the user spent it before deposit was called.
            - The channel was closed/settled between the allowance call and
              the deposit call.
        AddressWithoutCode: The channel was settled during the deposit
        execution.
*/
func (r *RaidenAPI) Deposit(tokenAddress, partnerAddress common.Address, amount *big.Int, pollTimeout time.Duration) (err error) {
	c, err := r.Raiden.db.GetChannel(tokenAddress, partnerAddress)
	if err != nil {
		return
	}
	token := r.Raiden.Chain.Token(tokenAddress)
	balance, err := token.BalanceOf(r.Raiden.NodeAddress)
	if err != nil {
		return
	}
	/*
			 Checking the balance is not helpful since r requires multiple
		     transactions that can race, e.g. the deposit check succeed but the
		     user spent his balance before deposit.
	*/
	if balance.Cmp(amount) < 0 {
		err = fmt.Errorf("Not enough balance to deposit. %s Available=%d Tried=%d", tokenAddress.String(), balance, amount)
		log.Error(err.Error())
		return rerr.ErrInsufficientFunds
	}
	err = token.Approve(c.ChannelAddress, amount)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	r.Raiden.db.RegisterChannelDepositCallback(func(c2 *channel.Serialization) (remove bool) {
		if c2.ChannelAddress == c.ChannelAddress {
			wg.Done()
			return true
		}
		return false
	})
	//deposit move ... todo
	result := r.Raiden.depositChannelClient(c.ChannelAddress, amount)
	err = <-result.Result
	if err != nil {
		return
	}
	/*
	 Wait until the `ChannelNewBalance` event is processed.
	*/
	wg.Wait()
	return nil
}

/*
TokenSwapAndWait Start an atomic swap operation by sending a MediatedTransfer with
    `maker_amount` of `maker_token` to `taker_address`. Only proceed when a
    new valid MediatedTransfer is received with `taker_amount` of
    `taker_token`.
*/
func (r *RaidenAPI) TokenSwapAndWait(identifier uint64, makerToken, takerToken, makerAddress, takerAddress common.Address,
	makerAmount, takerAmount *big.Int) error {
	result, err := r.tokenSwapAsync(identifier, makerToken, takerToken, makerAddress, takerAddress,
		makerAmount, takerAmount)
	if err != nil {
		return err
	}
	err = <-result.Result
	return err
}

func (r *RaidenAPI) tokenSwapAsync(identifier uint64, makerToken, takerToken, makerAddress, takerAddress common.Address,
	makerAmount, takerAmount *big.Int) (result *network.AsyncResult, err error) {
	chs, err := r.Raiden.db.GetChannelList(takerToken, utils.EmptyAddress)
	if err != nil || len(chs) == 0 {
		err = errors.New("unkown taker token")
		return
	}
	chs, err = r.Raiden.db.GetChannelList(makerToken, utils.EmptyAddress)
	if err != nil || len(chs) == 0 {
		err = errors.New("unkown maker token")
		return
	}
	tokenSwap := &TokenSwap{
		Identifier:      identifier,
		FromToken:       makerToken,
		FromAmount:      new(big.Int).Set(makerAmount),
		FromNodeAddress: makerAddress,
		ToToken:         takerToken,
		ToAmount:        new(big.Int).Set(takerAmount),
		ToNodeAddress:   takerAddress,
	}
	result = r.Raiden.tokenSwapMakerClient(tokenSwap)
	return
}

/*
ExpectTokenSwap Register an expected transfer for this node.

    If a MediatedMessage is received for the `maker_asset` with
    `maker_amount` then proceed to send a MediatedTransfer to
    `maker_address` for `taker_asset` with `taker_amount`.
*/
func (r *RaidenAPI) ExpectTokenSwap(identifier uint64, makerToken, takerToken, makerAddress, takerAddress common.Address,
	makerAmount, takerAmount *big.Int) (err error) {
	chs, err := r.Raiden.db.GetChannelList(takerToken, utils.EmptyAddress)
	if err != nil || len(chs) == 0 {
		err = errors.New("unkown taker token")
		return
	}
	chs, err = r.Raiden.db.GetChannelList(makerToken, utils.EmptyAddress)
	if err != nil || len(chs) == 0 {
		err = errors.New("unkown maker token")
		return
	}
	tokenSwap := &TokenSwap{
		Identifier:      identifier,
		FromToken:       makerToken,
		FromAmount:      new(big.Int).Set(makerAmount),
		FromNodeAddress: makerAddress,
		ToToken:         takerToken,
		ToAmount:        new(big.Int).Set(takerAmount),
		ToNodeAddress:   takerAddress,
	}
	r.Raiden.tokenSwapTakerClient(tokenSwap)
	return nil
}

//GetNodeNetworkState Returns the currently network status of `node_address
func (r *RaidenAPI) GetNodeNetworkState(nodeAddress common.Address) (deviceType string, isOnline bool) {
	return r.Raiden.Protocol.GetNetworkStatus(nodeAddress)
}

//StartHealthCheckFor Returns the currently network status of `node_address`.
func (r *RaidenAPI) StartHealthCheckFor(nodeAddress common.Address) (deviceType string, isOnline bool) {
	r.Raiden.startHealthCheckFor(nodeAddress)
	return r.GetNodeNetworkState(nodeAddress)
}

//GetTokenList returns all available tokens
func (r *RaidenAPI) GetTokenList() (tokens []common.Address) {
	tokensmap, _ := r.Raiden.db.GetAllTokens()
	for k := range tokensmap {
		tokens = append(tokens, k)
	}
	return
}

//TransferAndWait Do a transfer with `target` with the given `amount` of `token_address`.
func (r *RaidenAPI) TransferAndWait(token common.Address, amount *big.Int, fee *big.Int, target common.Address, identifier uint64, timeout time.Duration, isDirectTransfer bool) (err error) {
	result, err := r.transferAsync(token, amount, fee, target, identifier, isDirectTransfer)
	if err != nil {
		return err
	}
	if timeout > 0 {
		timeoutCh := time.After(timeout)
		select {
		case <-timeoutCh:
			err = rerr.ErrTransferTimeout
		case err = <-result.Result:
		}
	} else {
		err = <-result.Result
	}
	return
}

//Transfer transfer and wait
func (r *RaidenAPI) Transfer(token common.Address, amount *big.Int, fee *big.Int, target common.Address, identifier uint64, timeout time.Duration, isDirectTransfer bool) error {
	return r.TransferAndWait(token, amount, fee, target, identifier, timeout, isDirectTransfer)
}

//transferAsync
func (r *RaidenAPI) transferAsync(tokenAddress common.Address, amount *big.Int, fee *big.Int, target common.Address, identifier uint64, isDirectTransfer bool) (result *network.AsyncResult, err error) {
	tokens := r.Tokens()
	found := false
	for _, t := range tokens {
		if t == tokenAddress {
			found = true
			break
		}
	}
	if !found {
		err = errors.New("token not exist")
		return
	}
	if isDirectTransfer {
		var c *channel.Serialization
		c, err = r.Raiden.db.GetChannel(tokenAddress, target)
		if err != nil {
			err = fmt.Errorf("no direct channel token:%s,partner:%s", tokenAddress.String(), target.String())
			return
		}
		if c.State != transfer.ChannelStateOpened {
			err = fmt.Errorf("channel %s not opened", c.ChannelAddress.String())
			return
		}
	}
	if amount.Cmp(utils.BigInt0) <= 0 {
		err = rerr.ErrInvalidAmount
		return
	}
	log.Debug(fmt.Sprintf("initiating transfer initiator=%s target=%s token=%s amount=%d identifier=%d",
		r.Raiden.NodeAddress.String(), target.String(), tokenAddress.String(), amount, identifier))
	result = r.Raiden.transferAsyncClient(tokenAddress, amount, fee, target, identifier, isDirectTransfer)
	return
}

//Close a channel opened with `partner_address` for the given `token_address`. return when state has been updated to database
func (r *RaidenAPI) Close(tokenAddress, partnerAddress common.Address) (c *channel.Serialization, err error) {
	c, err = r.Raiden.db.GetChannel(tokenAddress, partnerAddress)
	if err != nil {
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	r.Raiden.db.RegisterChannelStateCallback(func(c2 *channel.Serialization) (remove bool) {
		log.Trace(fmt.Sprintf("wait %s closed ,get channle %s update",
			utils.APex2(c.ChannelAddress), utils.APex2(c2.ChannelAddress)))
		if c2.ChannelAddress == c.ChannelAddress {
			wg.Done()
			return true
		}
		return false
	})
	//send close channel request
	result := r.Raiden.closeChannelClient(c.ChannelAddress)
	err = <-result.Result
	if err != nil {
		return
	}
	wg.Wait()
	//reload data from database,
	return r.Raiden.db.GetChannelByAddress(c.ChannelAddress)
}

//Settle a closed channel with `partner_address` for the given `token_address`.return when state has been updated to database
func (r *RaidenAPI) Settle(tokenAddress, partnerAddress common.Address) (ch *channel.Serialization, err error) {
	c, err := r.Raiden.db.GetChannel(tokenAddress, partnerAddress)
	if c.State == transfer.ChannelStateOpened {
		err = rerr.InvalidState("channel is still open")
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	r.Raiden.db.RegisterChannelStateCallback(func(c2 *channel.Serialization) (remove bool) {
		log.Trace(fmt.Sprintf("wait %s settled ,get channle %s update",
			utils.APex2(c.ChannelAddress), utils.APex2(c2.ChannelAddress)))
		if c2.ChannelAddress == c.ChannelAddress {
			wg.Done()
			return true
		}
		return false
	})
	//send settle request
	result := r.Raiden.settleChannelClient(c.ChannelAddress)
	err = <-result.Result
	log.Trace(fmt.Sprintf("%s settled finish , err %v", utils.APex(c.ChannelAddress), err))
	if err != nil {
		return
	}
	wg.Wait()
	//reload data from database,
	return r.Raiden.db.GetChannelByAddress(c.ChannelAddress)
}

//GetTokenNetworkEvents return events about this token
func (r *RaidenAPI) GetTokenNetworkEvents(tokenAddress common.Address, fromBlock, toBlock int64) (data []interface{}, err error) {
	type eventData struct {
		/*
					 {
			        "event_type": "ChannelNew",
			        "settle_timeout": 10,
			        "netting_channel": "0xc0ea08a2d404d3172d2add29a45be56da40e2949",
			        "participant1": "0x4894a542053248e0c504e3def2048c08f73e1ca6",
			        "participant2": "0x356857Cd22CBEFccDa4e96AF13b408623473237A"
			    }
		*/
		EventType      string `json:"event_type"`
		SettleTimeout  int    `json:"settle_timeout"`
		NettingChannel string `json:"netting_channel"`
		Participant1   string `json:"participant1"`
		Participant2   string `json:"participant2"`
		TokenAddress   string `json:"token_address"`
	}
	tokens, err := r.Raiden.db.GetAllTokens()
	if err != nil {
		return
	}
	for t, manager := range tokens {
		if tokenAddress == utils.EmptyAddress || t == tokenAddress {
			events, err := r.Raiden.BlockChainEvents.GetAllChannelManagerEvents(manager, fromBlock, toBlock)
			if err != nil {
				return nil, err
			}
			for _, e := range events {
				e2 := e.(*blockchainold.EventChannelOpen)
				ed := &eventData{
					EventType:      e2.EventName,
					SettleTimeout:  e2.SettleTimeout,
					NettingChannel: e2.NettingChannelAddress.String(),
					Participant1:   e2.Participant1.String(),
					Participant2:   e2.Participant2.String(),
					TokenAddress:   t.String(),
				}
				data = append(data, ed)
			}
		}
	}
	return
}

//GetNetworkEvents all raiden events
func (r *RaidenAPI) GetNetworkEvents(fromBlock, toBlock int64) ([]interface{}, error) {
	type eventData struct {
		/*
					 "event_type": "TokenAdded",
			        "token_address": "0xea674fdde714fd979de3edf0f56aa9716b898ec8",
			        "channel_manager_address": "0xc0ea08a2d404d3172d2add29a45be56da40e2949"
		*/
		EventType             string `json:"event_type"`
		TokenAddress          string `json:"token_address"`
		ChannelManagerAddress string `json:"channel_manager_address"`
	}
	events, err := r.Raiden.BlockChainEvents.GetAllRegistryEvents(r.Raiden.RegistryAddress, fromBlock, toBlock)
	if err != nil {
		return nil, err
	}
	var data []interface{}
	for _, e := range events {
		e2 := e.(*blockchainold.EventTokenNetworkCreated)
		ed := &eventData{
			EventType:             e2.EventName,
			TokenAddress:          e2.TokenAddress.String(),
			ChannelManagerAddress: e2.TokenNetworkAddress.String(),
		}
		data = append(data, ed)
	}
	return data, nil
}

//GetChannelEvents events of this channel
func (r *RaidenAPI) GetChannelEvents(channelAddress common.Address, fromBlock, toBlock int64) (data []transfer.Event, err error) {

	var events []transfer.Event
	events, err = r.Raiden.BlockChainEvents.GetAllNettingChannelEvents(channelAddress, fromBlock, toBlock)
	if err != nil {
		return
	}
	for _, e := range events {
		m := make(map[string]interface{})
		switch e2 := e.(type) {
		case *blockchainold.EventChannelNewBalance:
			m["event_type"] = e2.EventName
			m["participant"] = e2.ParticipantAddress.String()
			m["balance"] = e2.Balance
			m["block_number"] = e2.BlockNumber
			data = append(data, m)
		case *blockchainold.EventChannelClosed:
			m["event_type"] = e2.EventName
			m["netting_channel_address"] = e2.ContractAddress.String()
			m["closing_address"] = e2.ClosingAddress.String()
			data = append(data, m)
		case *blockchainold.EventChannelSettled:
			m["event_type"] = e2.EventName
			m["netting_channel_address"] = e2.ContractAddress.String()
			m["block_number"] = e2.BlockNumber
			data = append(data, m)
		case *blockchainold.EventSecretRevealed:
			m["event_type"] = e2.EventName
			m["netting_channel_address"] = e2.ContractAddress.String()
			m["secret"] = e2.Secret.String()
			data = append(data, m)
			//case *blockchain.EventNonClosingBalanceProofUpdated:
			//	m["event_type"] = e2.EventName
			//	m["token_address"] = t.String()
			//	m["channel_manager_address"] = graph.TokenAddress.String()
		}

	}

	var raidenEvents []*models.InternalEvent
	raidenEvents, err = r.Raiden.db.GetEventsInBlockRange(fromBlock, toBlock)
	if err != nil {
		return
	}
	//Here choose which raiden internal events we want to expose to the end user
	for _, ev := range raidenEvents {
		m := make(map[string]interface{})
		switch e2 := ev.EventObject.(type) {
		case *transfer.EventTransferSentSuccess:
			m["event_type"] = "EventTransferSentSuccess"
			m["identifier"] = e2.Identifier
			m["block_number"] = ev.BlockNumber
			m["amount"] = e2.Amount
			m["target"] = e2.Target
			data = append(data, m)
		case *transfer.EventTransferSentFailed:
			m["event_type"] = "EventTransferSentFailed"
			m["identifier"] = e2.Identifier
			m["block_number"] = ev.BlockNumber
			m["reason"] = e2.Reason
			data = append(data, m)
		case *transfer.EventTransferReceivedSuccess:
			m["event_type"] = "EventTransferReceivedSuccess"
			m["identifier"] = e2.Identifier
			m["block_number"] = ev.BlockNumber
			m["amount"] = e2.Amount
			m["initiator"] = e2.Initiator.String()
			data = append(data, m)
		}
	}
	return
}

/*
GetSentTransfers query sent transfers from db
*/
func (r *RaidenAPI) GetSentTransfers(from, to int64) ([]*models.SentTransfer, error) {
	return r.Raiden.db.GetSentTransferInBlockRange(from, to)
}

/*
GetReceivedTransfers query received transfers from db
*/
func (r *RaidenAPI) GetReceivedTransfers(from, to int64) ([]*models.ReceivedTransfer, error) {
	return r.Raiden.db.GetReceivedTransferInBlockRange(from, to)
}

//Stop stop for mobile app
func (r *RaidenAPI) Stop() {
	log.Info("calling api stop..")
	r.Raiden.Stop()
	log.Info("stop successful..")
}

/*
{
    "channel_address": "0x5B3F0E96E45e1e4351F6460feBfB6007af25FBB0",
 "update_transfer":{
        "nonce": 32,
        "transferred_amount": 1800000000000000,
        "locksroot": " 0x447b478a024ade59c5c18e348c357aae6a4ec6e30131213f8cf6444214c57e89",
        "extra_hash": " 0x557b478a024ade59c5c18e348c357aae6a4ec6e30131213f8cf6444214c57e89",
        "closing_signature": " 0x557b478a024ade59c5c18e348c357aae6a4ec6e30131213f8cf6444214c57e89557b478a024ade59c5c18e348c357aae6a4ec6e30131213f8cf6444214c57e8927",
        "non_closing_signature": " 0x557b478a024ade59c5c18e348c357aae6a4ec6e30131213f8cf6444214c57e89557b478a024ade59c5c18e348c357aae6a4ec6e30131213f8cf6444214c57e8927"
 },
 "withdraws":[
     {
        "locked_encoded": "0x00000033333333333333333333333333333333333333333",
        "merkle_proof": "0x3333333333333333333333333333",
        "secret": "0x333333333333333333333333333333333333333",
     },
      {
        "locked_encoded": "0x00000033333333333333333333333333333333333333333",
        "merkle_proof": "0x3333333333333333333333333333",
        "secret": "0x333333333333333333333333333333333333333",
     },
 ],
}
*/
type updateTransfer struct {
	Nonce               int64    `json:"nonce"`
	TransferAmount      *big.Int `json:"transfer_amount"`
	Locksroot           string   `json:"locksroot"`
	ExtraHash           string   `json:"extra_hash"`
	ClosingSignature    string   `json:"closing_signature"`
	NonClosingSignature string   `json:"non_closing_signature"`
}
type withdraw struct {
	LockedEncoded string `json:"locked_encoded"`
	MerkleProof   string `json:"merkle_proof"`
	Secret        string `json:"secret"`
}

//ChannelFor3rd is for 3rd party to call update transfer
type ChannelFor3rd struct {
	ChannelAddress string         `json:"channel_address"`
	UpdateTransfer updateTransfer `json:"update_transfer"`
	Withdraws      []*withdraw    `json:"withdraws"`
}

/*
ChannelInformationFor3rdParty generate all information need by 3rd party
*/
func (r *RaidenAPI) ChannelInformationFor3rdParty(channelAddr, thirdAddr common.Address) (result *ChannelFor3rd, err error) {
	var sig []byte
	c, err := r.GetChannel(channelAddr)
	if err != nil {
		return
	}
	c3 := new(ChannelFor3rd)
	c3.ChannelAddress = channelAddr.String()
	if c.PartnerBalanceProof == nil {
		result = c3
		return
	}
	c3.UpdateTransfer.Nonce = c.PartnerBalanceProof.Nonce
	c3.UpdateTransfer.TransferAmount = c.PartnerBalanceProof.TransferAmount
	c3.UpdateTransfer.Locksroot = c.PartnerBalanceProof.LocksRoot.String()
	c3.UpdateTransfer.ExtraHash = c.PartnerBalanceProof.MessageHash.String()
	c3.UpdateTransfer.ClosingSignature = common.Bytes2Hex(c.PartnerBalanceProof.Signature)
	sig, err = signFor3rd(c, thirdAddr, r.Raiden.PrivateKey)
	if err != nil {
		return
	}
	c3.UpdateTransfer.NonClosingSignature = common.Bytes2Hex(sig)
	tree, err := transfer.NewMerkleTree(c.PartnerLeaves)
	if err != nil {
		return
	}
	var ws []*withdraw
	for _, l := range c.PartnerLock2UnclaimedLocks {
		proof := channel.ComputeProofForLock(l.Secret, l.Lock, tree)
		w := &withdraw{
			LockedEncoded: common.Bytes2Hex(proof.LockEncoded),
			Secret:        l.Secret.String(),
			MerkleProof:   common.Bytes2Hex(transfer.Proof2Bytes(proof.MerkleProof)),
		}
		log.Trace(fmt.Sprintf("prootf=%s", utils.StringInterface(proof, 3)))
		ws = append(ws, w)
	}
	c3.Withdraws = ws

	result = c3
	return
}

//make sure PartnerBalanceProof is not nil
func signFor3rd(c *channel.Serialization, thirdAddr common.Address, privkey *ecdsa.PrivateKey) (sig []byte, err error) {
	if c.PartnerBalanceProof == nil {
		log.Error(fmt.Sprintf("PartnerBalanceProof is nil,must ber a error"))
		return nil, errors.New("empty PartnerBalanceProof")
	}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, c.PartnerBalanceProof.Nonce)
	buf.Write(utils.BigIntTo32Bytes(c.PartnerBalanceProof.TransferAmount))
	buf.Write(c.PartnerBalanceProof.LocksRoot[:])
	buf.Write(c.ChannelAddress[:])
	buf.Write(c.PartnerBalanceProof.MessageHash[:])
	buf.Write(c.PartnerBalanceProof.Signature)
	buf.Write(thirdAddr[:])
	dataToSign := buf.Bytes()
	return utils.SignData(privkey, dataToSign)
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
