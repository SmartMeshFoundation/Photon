package smartraiden

import (
	"time"

	"fmt"

	"math/big"

	"sync"

	"errors"

	"github.com/SmartMeshFoundation/SmartRaiden/blockchain"
	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/models"
	"github.com/SmartMeshFoundation/SmartRaiden/network"
	"github.com/SmartMeshFoundation/SmartRaiden/rerr"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

type RaidenApi struct {
	Raiden *RaidenService
}

//CLI interface.
func NewRaidenApi(raiden *RaidenService) *RaidenApi {
	return &RaidenApi{Raiden: raiden}
}

func (this *RaidenApi) Address() common.Address {
	return this.Raiden.NodeAddress
}

//Return a list of the tokens registered with the default registry.
func (this *RaidenApi) Tokens() []common.Address {
	addresses, _ := this.Raiden.Registry.TokenAddresses()
	return addresses
}

/*
Returns a list of channels associated with the optionally given
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
func (this *RaidenApi) GetChannelList(tokenAddress common.Address, partnerAddress common.Address) (cs []*channel.ChannelSerialization, err error) {
	return this.Raiden.db.GetChannelList(tokenAddress, partnerAddress)
}
func (this *RaidenApi) GetChannel(channelAddress common.Address) (c *channel.ChannelSerialization, err error) {
	return this.Raiden.db.GetChannelByAddress(channelAddress)
}

/*
	 If the token is registered then, return the channel manager address.
		   Also make sure that the channel manager is registered with the node.

		   Returns None otherwise.
*/
func (this *RaidenApi) ManagerAddressIfTokenRegistered(tokenAddress common.Address) (mgrAddr common.Address, err error) {
	mgrAddr, err = this.Raiden.Registry.ChannelManagerByToken(tokenAddress)
	if err != nil {
		return
	}
	return
}

/*
Will register the token at `token_address` with raiden. If it's already
    registered, will throw an exception.
*/
func (this *RaidenApi) RegisterToken(tokenAddress common.Address) (mgrAddr common.Address, err error) {
	mgrAddr, err = this.Raiden.Registry.ChannelManagerByToken(tokenAddress)
	if err == nil && mgrAddr != utils.EmptyAddress {
		err = errors.New("Token already registered")
		return
	}
	//for non exist tokenaddress, ChannelManagerByToken will return a error: `abi : unmarshalling empty output`
	if err == rerr.NoTokenManager {
		return this.Raiden.Registry.AddToken(tokenAddress)
	} else {
		return
	}

}

/*
Instruct the ConnectionManager to establish and maintain a connection to the token
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
func (this *RaidenApi) ConnectTokenNetwork(tokenAddress common.Address, funds *big.Int, initialChannelTarget int64, joinableFundsTarget float64) error {
	cm, err := this.Raiden.ConnectionManagerForToken(tokenAddress)
	if err != nil {
		return err
	}
	return cm.Connect(funds, initialChannelTarget, joinableFundsTarget)
}

/*
Instruct the ConnectionManager to close all channels and wait for
    settlement
*/
func (this *RaidenApi) LeaveTokenNetwork(tokenAddress common.Address, onlyReceiving bool) ([]*channel.ChannelSerialization, error) {
	cm, err := this.Raiden.ConnectionManagerForToken(tokenAddress)
	if err != nil {
		return nil, err
	}
	chs := cm.Leave(onlyReceiving)
	return chs, nil
}

/*
Get a dict whose keys are token addresses and whose values are
    open channels, funds of last request, sum of deposits and number of channels
*/
func (this *RaidenApi) GetConnectionManagersInfo() map[string]interface{} {
	type info struct {
		Funds       *big.Int `json:"funds"`
		SumDeposits *big.Int `json:"sum_deposits"`
		Channels    int      `json:"channels"`
	}
	infos := make(map[string]interface{})
	for _, t := range this.GetTokenList() {
		cm, err := this.Raiden.ConnectionManagerForToken(t)
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
func (this *RaidenApi) Open(tokenAddress, partnerAddress common.Address, settleTimeout, revealTimeout int) (ch *channel.ChannelSerialization, err error) {
	if revealTimeout <= 0 {
		revealTimeout = this.Raiden.Config.RevealTimeout
	}
	if settleTimeout <= 0 {
		settleTimeout = this.Raiden.Config.SettleTimeout
	}
	if settleTimeout <= revealTimeout {
		err = rerr.InvalidSettleTimeout
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	this.Raiden.db.RegisterNewChannellCallback(func(c *channel.ChannelSerialization) (remove bool) {
		if c.TokenAddress == tokenAddress && c.PartnerAddress == partnerAddress {
			wg.Done()
			return true
		}
		return false
	})
	result := this.Raiden.NewChannelClient(tokenAddress, partnerAddress, settleTimeout)
	err = <-result.Result
	if err != nil {
		return
	}
	//wait
	wg.Wait()
	return this.Raiden.db.GetChannel(tokenAddress, partnerAddress)
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
func (this *RaidenApi) Deposit(tokenAddress, partnerAddress common.Address, amount *big.Int, pollTimeout time.Duration) (err error) {
	c, err := this.Raiden.db.GetChannel(tokenAddress, partnerAddress)
	if err != nil {
		return
	}
	token := this.Raiden.Chain.Token(tokenAddress)
	balance, err := token.BalanceOf(this.Raiden.NodeAddress)
	if err != nil {
		return
	}
	/*
			 Checking the balance is not helpful since this requires multiple
		     transactions that can race, e.g. the deposit check succeed but the
		     user spent his balance before deposit.
	*/
	if balance.Cmp(amount) < 0 {
		err = fmt.Errorf("Not enough balance to deposit. %s Available=%d Tried=%d", tokenAddress.String(), balance, amount)
		log.Error(err.Error())
		return rerr.InsufficientFunds
	}
	err = token.Approve(c.ChannelAddress, amount)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	this.Raiden.db.RegisterChannelDepositCallback(func(c2 *channel.ChannelSerialization) (remove bool) {
		if c2.ChannelAddress == c.ChannelAddress {
			wg.Done()
			return true
		}
		return false
	})
	//deposit move ... todo
	result := this.Raiden.DepositChannelClient(c.ChannelAddress, amount)
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
Start an atomic swap operation by sending a MediatedTransfer with
    `maker_amount` of `maker_token` to `taker_address`. Only proceed when a
    new valid MediatedTransfer is received with `taker_amount` of
    `taker_token`.
*/
func (this *RaidenApi) TokenSwapAndWait(identifier uint64, makerToken, takerToken, makerAddress, takerAddress common.Address,
	makerAmount, takerAmount *big.Int) error {
	result, err := this.TokenSwapAsync(identifier, makerToken, takerToken, makerAddress, takerAddress,
		makerAmount, takerAmount)
	if err != nil {
		return err
	}
	err = <-result.Result
	return err
}

func (this *RaidenApi) TokenSwapAsync(identifier uint64, makerToken, takerToken, makerAddress, takerAddress common.Address,
	makerAmount, takerAmount *big.Int) (result *network.AsyncResult, err error) {
	chs, err := this.Raiden.db.GetChannelList(takerToken, utils.EmptyAddress)
	if err != nil || len(chs) == 0 {
		err = errors.New("unkown taker token")
		return
	}
	chs, err = this.Raiden.db.GetChannelList(makerToken, utils.EmptyAddress)
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
	result = this.Raiden.TokenSwapMakerClient(tokenSwap)
	return
}

/*
Register an expected transfer for this node.

    If a MediatedMessage is received for the `maker_asset` with
    `maker_amount` then proceed to send a MediatedTransfer to
    `maker_address` for `taker_asset` with `taker_amount`.
*/
func (this *RaidenApi) ExpectTokenSwap(identifier uint64, makerToken, takerToken, makerAddress, takerAddress common.Address,
	makerAmount, takerAmount *big.Int) (err error) {
	chs, err := this.Raiden.db.GetChannelList(takerToken, utils.EmptyAddress)
	if err != nil || len(chs) == 0 {
		err = errors.New("unkown taker token")
		return
	}
	chs, err = this.Raiden.db.GetChannelList(makerToken, utils.EmptyAddress)
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
	this.Raiden.TokenSwapTakerClient(tokenSwap)
	return nil
}

//Returns the currently network status of `node_address
func (this *RaidenApi) GetNodeNetworkState(nodeAddress common.Address) string {
	return this.Raiden.Protocol.GetNetworkStatus(nodeAddress)
}

//Returns the currently network status of `node_address`.
func (this *RaidenApi) StartHealthCheckFor(nodeAddress common.Address) string {
	this.Raiden.StartHealthCheckFor(nodeAddress)
	return this.GetNodeNetworkState(nodeAddress)
}

func (this *RaidenApi) GetTokenList() (tokens []common.Address) {
	tokensmap, _ := this.Raiden.db.GetAllTokens()
	for k, _ := range tokensmap {
		tokens = append(tokens, k)
	}
	return
}

//Do a transfer with `target` with the given `amount` of `token_address`.
func (this *RaidenApi) TransferAndWait(token common.Address, amount *big.Int, fee *big.Int, target common.Address, identifier uint64, timeout time.Duration) (err error) {
	result, err := this.TransferAsync(token, amount, fee, target, identifier)
	if err != nil {
		return err
	}
	if timeout > 0 {
		timeoutCh := time.After(timeout)
		select {
		case <-timeoutCh:
			err = rerr.TransferTimeout
		case err = <-result.Result:
		}
	} else {
		err = <-result.Result
	}
	return
}
func (this *RaidenApi) Transfer(token common.Address, amount *big.Int, fee *big.Int, target common.Address, identifier uint64, timeout time.Duration) error {
	return this.TransferAndWait(token, amount, fee, target, identifier, timeout)
}

func (this *RaidenApi) TransferAsync(tokenAddress common.Address, amount *big.Int, fee *big.Int, target common.Address, identifier uint64) (result *network.AsyncResult, err error) {
	tokens := this.Tokens()
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
	if amount.Cmp(utils.BigInt0) <= 0 {
		err = rerr.InvalidAmount
		return
	}
	log.Debug(fmt.Sprintf("initiating transfer initiator=%s target=%s token=%s amount=%d identifier=%d",
		this.Raiden.NodeAddress.String(), target.String(), tokenAddress.String(), amount, identifier))
	result = this.Raiden.MediatedTransferAsyncClient(tokenAddress, amount, fee, target, identifier)
	return
}

//Close a channel opened with `partner_address` for the given `token_address`.
//Close a channel opened with `partner_address` for the given `token_address`. return when state has been updated to database
func (this *RaidenApi) Close(tokenAddress, partnerAddress common.Address) (c *channel.ChannelSerialization, err error) {
	c, err = this.Raiden.db.GetChannel(tokenAddress, partnerAddress)
	if err != nil {
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	this.Raiden.db.RegisterChannelStateCallback(func(c2 *channel.ChannelSerialization) (remove bool) {
		if c2.ChannelAddress == c.ChannelAddress {
			wg.Done()
			return true
		}
		return false
	})
	//send close channel request
	result := this.Raiden.CloseChannelClient(c.ChannelAddress)
	err = <-result.Result
	if err != nil {
		return
	}
	wg.Wait()
	//reload data from database,
	return this.Raiden.db.GetChannelByAddress(c.ChannelAddress)
}

//Settle a closed channel with `partner_address` for the given `token_address`.return when state has been updated to database
func (this *RaidenApi) Settle(tokenAddress, partnerAddress common.Address) (ch *channel.ChannelSerialization, err error) {
	c, err := this.Raiden.db.GetChannel(tokenAddress, partnerAddress)
	if c.State == transfer.CHANNEL_STATE_OPENED {
		err = rerr.InvalidState("channel is still open")
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	this.Raiden.db.RegisterChannelStateCallback(func(c2 *channel.ChannelSerialization) (remove bool) {
		log.Trace(fmt.Sprintf("wait %s settled ,get channle %s update",
			utils.APex2(c.ChannelAddress), utils.APex2(c2.ChannelAddress)))
		if c2.ChannelAddress == c.ChannelAddress {
			wg.Done()
			return true
		}
		return false
	})
	//send settle request
	result := this.Raiden.SettleChannelClient(c.ChannelAddress)
	err = <-result.Result
	if err != nil {
		return
	}
	wg.Wait()
	//reload data from database,
	return this.Raiden.db.GetChannelByAddress(c.ChannelAddress)
}
func (this *RaidenApi) GetTokenNetworkEvents(tokenAddress common.Address, fromBlock, toBlock int64) (data []interface{}, err error) {
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
	tokens, err := this.Raiden.db.GetAllTokens()
	if err != nil {
		return
	}
	for t, manager := range tokens {
		if tokenAddress == utils.EmptyAddress || t == tokenAddress {
			events, err := this.Raiden.BlockChainEvents.GetAllChannelManagerEvents(manager, fromBlock, toBlock)
			if err != nil {
				return nil, err
			}
			for _, e := range events {
				e2 := e.(*blockchain.EventChannelNew)
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

func (this *RaidenApi) GetNetworkEvents(fromBlock, toBlock int64) ([]interface{}, error) {
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
	events, err := this.Raiden.BlockChainEvents.GetAllRegistryEvents(this.Raiden.RegistryAddress, fromBlock, toBlock)
	if err != nil {
		return nil, err
	}
	var data []interface{}
	for _, e := range events {
		e2 := e.(*blockchain.EventTokenAdded)
		ed := &eventData{
			EventType:             e2.EventName,
			TokenAddress:          e2.TokenAddress.String(),
			ChannelManagerAddress: e2.ChannelManagerAddress.String(),
		}
		data = append(data, ed)
	}
	return data, nil
}

func (this *RaidenApi) GetChannelEvents(channelAddress common.Address, fromBlock, toBlock int64) (data []transfer.Event, err error) {

	var events []transfer.Event
	events, err = this.Raiden.BlockChainEvents.GetAllNettingChannelEvents(channelAddress, fromBlock, toBlock)
	if err != nil {
		return
	}
	for _, e := range events {
		m := make(map[string]interface{})
		switch e2 := e.(type) {
		case *blockchain.EventChannelNewBalance:
			m["event_type"] = e2.EventName
			m["participant"] = e2.ParticipantAddress.String()
			m["balance"] = e2.Balance
			m["block_number"] = e2.BlockNumber
			data = append(data, m)
		case *blockchain.EventChannelClosed:
			m["event_type"] = e2.EventName
			m["netting_channel_address"] = e2.ContractAddress.String()
			m["closing_address"] = e2.ClosingAddress.String()
			data = append(data, m)
		case *blockchain.EventChannelSettled:
			m["event_type"] = e2.EventName
			m["netting_channel_address"] = e2.ContractAddress.String()
			m["block_number"] = e2.BlockNumber
			data = append(data, m)
		case *blockchain.EventChannelSecretRevealed:
			m["event_type"] = e2.EventName
			m["netting_channel_address"] = e2.ContractAddress.String()
			m["secret"] = e2.Secret.String()
			data = append(data, m)
			//case *blockchain.EventTransferUpdated:
			//	m["event_type"] = e2.EventName
			//	m["token_address"] = t.String()
			//	m["channel_manager_address"] = graph.ChannelManagerAddress.String()
		}

	}

	var raidenEvents []*models.InternalEvent
	raidenEvents, err = this.Raiden.db.GetEventsInBlockRange(fromBlock, toBlock)
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
func (this *RaidenApi) Stop() {
	log.Info("calling api stop..")
	this.Raiden.Stop()
	log.Info("stop successful..")
}

type EventTransferSentSuccessWrapper struct {
	transfer.EventTransferSentSuccess
	BlockNumber int64
	Name        string
}
type EventTransferSentFailedWrapper struct {
	transfer.EventTransferSentFailed
	BlockNumber int64
	Name        string
}
type EventEventTransferReceivedSuccessWrapper struct {
	transfer.EventTransferReceivedSuccess
	BlockNumber int64
	Name        string
}
