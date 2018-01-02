package raiden_network

import (
	"time"

	"fmt"

	"github.com/SmartMeshFoundation/raiden-network/channel"
	"github.com/SmartMeshFoundation/raiden-network/network"
	"github.com/SmartMeshFoundation/raiden-network/rerr"
	"github.com/SmartMeshFoundation/raiden-network/transfer"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fatedier/frp/src/utils/log"
	"github.com/kataras/go-errors"
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
func (this *RaidenApi) GetChannelList(tokenAddress common.Address, partnerAddress common.Address) (channels []*channel.Channel) {
	if tokenAddress != utils.EmptyAddress && partnerAddress != utils.EmptyAddress {
		graph, ok := this.Raiden.Token2ChannelGraph[tokenAddress]
		if !ok {
			return
		}
		ch, ok := graph.PartenerAddress2Channel[partnerAddress]
		if !ok {
			return
		}
		channels = []*channel.Channel{ch}
		return
	} else if tokenAddress != utils.EmptyAddress {
		graph, ok := this.Raiden.Token2ChannelGraph[tokenAddress]
		if !ok {
			return
		}
		for _, c := range graph.ChannelAddres2Channel {
			channels = append(channels, c)
		}
	} else if partnerAddress != utils.EmptyAddress {
		for _, g := range this.Raiden.Token2ChannelGraph {
			for p, c := range g.PartenerAddress2Channel {
				if p == partnerAddress {
					channels = append(channels, c)
				}
			}
		}
	} else {
		for _, g := range this.Raiden.Token2ChannelGraph {
			for _, c := range g.PartenerAddress2Channel {
				channels = append(channels, c)
			}
		}
	}
	return
}

func (this *RaidenApi) GetChannel(channelAddress common.Address) (ch *channel.Channel, err error) {
	channels := this.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	for _, c := range channels {
		if c.MyAddress == channelAddress {
			return c, nil
		}
	}
	return nil, rerr.ChannelNotFound(channelAddress.String())
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
	/*
		为什么无关自己的mgr也要注册呢?没必要啊 todo fix it
	*/
	if !this.Raiden.ChannelManagerIsRegistered(mgrAddr) {
		this.Raiden.RegisterChannelManager(mgrAddr)
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
	if err != nil && err.Error() == "abi: unmarshalling empty output" {
		return this.Raiden.Registry.AddToken(tokenAddress)
	} else {
		return
	}

}

/*
def connect_token_network(
        self,
        token_address,
        funds,
        initial_channel_target=3,
        joinable_funds_target=.4):
    """Instruct the ConnectionManager to establish and maintain a connection to the token
    network.

    If the `token_address` is not already part of the raiden network, this will also register
    the token.

    Args:
        token_address (bin): the ERC20 token network to connect to.
        funds (int): the amount of funds that can be used by the ConnectionMananger.
        initial_channel_target (int): number of channels to open proactively.
        joinable_funds_target (float): fraction of the funds that will be used to join
            channels opened by other participants.
    """
    if not isaddress(token_address):
        raise InvalidAddress('token_address must be a valid address in binary')

    try:
        connection_manager = self.raiden.connection_manager_for_token(token_address)
    except InvalidAddress:
        # token is not yet registered
        self.raiden.default_registry.add_token(token_address)

        # wait for registration
        while token_address not in self.raiden.tokens_to_connectionmanagers:
            gevent.sleep(self.raiden.alarm.wait_time)
        connection_manager = self.raiden.connection_manager_for_token(token_address)

    connection_manager.connect(
        funds,
        initial_channel_target=initial_channel_target,
        joinable_funds_target=joinable_funds_target
    )

def leave_token_network(self, token_address, only_receiving=True):
    """Instruct the ConnectionManager to close all channels and wait for
    settlement.
    """
    connection_manager = self.raiden.connection_manager_for_token(token_address)
    return connection_manager.leave(only_receiving)

def get_connection_managers_info(self):
    """Get a dict whose keys are token addresses and whose values are
    open channels, funds of last request, sum of deposits and number of channels"""
    connection_managers = dict()

    for token in self.get_tokens_list():
        try:
            connection_manager = self.raiden.connection_manager_for_token(token)
        except InvalidAddress:
            connection_manager = None
        if connection_manager is not None and connection_manager.open_channels:
            connection_managers[connection_manager.token_address] = {
                'funds': connection_manager.funds,
                'sum_deposits': connection_manager.sum_deposits,
                'channels': len(connection_manager.open_channels),
            }

    return connection_managers
*/
/*
Open a channel with the peer at `partner_address`
    with the given `token_address`.
*/
func (this *RaidenApi) Open(tokenAddress, partnerAddress common.Address, settleTimeout, revealTimeout int) (ch *channel.Channel, err error) {
	if revealTimeout <= 0 {
		revealTimeout = this.Raiden.Config.RevealTimeout
	}
	if settleTimeout <= 0 {
		settleTimeout = this.Raiden.Config.SettleTimeout
	}
	if settleTimeout <= revealTimeout {
		err = rerr.InvalidSettleTimeout
	}
	chMgrAddr, err := this.Raiden.Registry.ChannelManagerByToken(tokenAddress)
	if err != nil {
		return
	}
	if _, ok := this.Raiden.Token2ChannelGraph[tokenAddress]; !ok {
		log.Error("open not exists token channel")
		err = rerr.UnknownTokenAddress("when open channel")
		return
	}
	chMgr := this.Raiden.Chain.Manager(chMgrAddr)
	_, err = chMgr.NewChannel(partnerAddress, settleTimeout)
	if err != nil {
		return
	}
	//wait channelNew Event
	for {
		g, ok := this.Raiden.Token2ChannelGraph[tokenAddress]
		if !ok {
			time.Sleep(time.Second)
			continue
		}
		_, ok = g.PartenerAddress2Channel[partnerAddress]
		if !ok {
			time.Sleep(time.Second)
			continue
		}
		break
	}
	g := this.Raiden.Token2ChannelGraph[tokenAddress]
	ch = g.PartenerAddress2Channel[partnerAddress]
	return
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
func (this *RaidenApi) Deposit(tokenAddress, partnerAddress common.Address, amount int64, pollTimeout time.Duration) (err error) {
	graph, ok := this.Raiden.Token2ChannelGraph[tokenAddress]
	if !ok {
		return rerr.InvalidAddress("Unknown token address")
	}
	ch, ok := graph.PartenerAddress2Channel[partnerAddress]
	if !ok {
		return rerr.InvalidAddress("No channel with partner_address for the given token")
	}
	if ch.TokenAddress != tokenAddress { //impossible
		return rerr.InvalidAddress("No channel with partner_address for the given token")
	}
	token := this.Raiden.Chain.Token(tokenAddress)
	balance, err := token.BalanceOf(this.Raiden.NodeAddress)
	if err != nil {
		return
	}
	/*
			# Checking the balance is not helpful since this requires multiple
		    # transactions that can race, e.g. the deposit check succeed but the
		    # user spent his balance before deposit.
	*/
	if balance < amount {
		err = fmt.Errorf("Not enough balance to deposit. %s Available=%d Tried=%d", tokenAddress.String(), balance, amount)
		log.Error(err.Error())
		return rerr.InsufficientFunds
	}
	oldBalance := ch.ContractBalance()
	err = ch.ExternState.Deposit(amount)
	if err != nil {
		return
	}
	/*
			# Wait until the `ChannelNewBalance` event is processed.
		        #
		        # Usually a single sleep is sufficient, since the `deposit` waits for
		        # the transaction to be polled.
	*/
	//为什么收不到我自己的newbalance事件呢?
	var totalSleep time.Duration = 0
	for {
		//too many race conditons! bai
		if oldBalance != ch.ContractBalance() {
			break
		}
		time.Sleep(time.Second)
		totalSleep += time.Second
		if totalSleep > pollTimeout {
			return errors.New("timeout")
		}
	}
	return
}

/*
def token_swap_and_wait(
        self,
        identifier,
        maker_token,
        maker_amount,
        maker_address,
        taker_token,
        taker_amount,
        taker_address):
    """ Start an atomic swap operation by sending a MediatedTransfer with
    `maker_amount` of `maker_token` to `taker_address`. Only proceed when a
    new valid MediatedTransfer is received with `taker_amount` of
    `taker_token`.
    """

    async_result = self.token_swap_async(
        identifier,
        maker_token,
        maker_amount,
        maker_address,
        taker_token,
        taker_amount,
        taker_address,
    )
    async_result.wait()

def token_swap_async(
        self,
        identifier,
        maker_token,
        maker_amount,
        maker_address,
        taker_token,
        taker_amount,
        taker_address):
    """ Start a token swap operation by sending a MediatedTransfer with
    `maker_amount` of `maker_token` to `taker_address`. Only proceed when a
    new valid MediatedTransfer is received with `taker_amount` of
    `taker_token`.
    """
    if not isaddress(maker_token):
        raise InvalidAddress(
            'Address for maker token is not in expected binary format in token swap'
        )
    if not isaddress(maker_address):
        raise InvalidAddress(
            'Address for maker is not in expected binary format in token swap'
        )

    if not isaddress(taker_token):
        raise InvalidAddress(
            'Address for taker token is not in expected binary format in token swap'
        )
    if not isaddress(taker_address):
        raise InvalidAddress(
            'Address for taker is not in expected binary format in token swap'
        )

    channelgraphs = self.raiden.token_to_channelgraph

    if taker_token not in channelgraphs:
        log.error('Unknown token {}'.format(pex(taker_token)))
        return

    if maker_token not in channelgraphs:
        log.error('Unknown token {}'.format(pex(maker_token)))
        return

    token_swap = TokenSwap(
        identifier,
        maker_token,
        maker_amount,
        maker_address,
        taker_token,
        taker_amount,
        taker_address,
    )

    async_result = AsyncResult()
    task = MakerTokenSwapTask(
        self.raiden,
        token_swap,
        async_result,
    )
    task.start()

    # the maker is expecting the taker transfer
    key = SwapKey(
        identifier,
        taker_token,
        taker_amount,
    )
    self.raiden.swapkey_to_greenlettask[key] = task
    self.raiden.swapkey_to_tokenswap[key] = token_swap

    return async_result

def expect_token_swap(
        self,
        identifier,
        maker_token,
        maker_amount,
        maker_address,
        taker_token,
        taker_amount,
        taker_address):
    """ Register an expected transfer for this node.

    If a MediatedMessage is received for the `maker_asset` with
    `maker_amount` then proceed to send a MediatedTransfer to
    `maker_address` for `taker_asset` with `taker_amount`.
    """

    if not isaddress(maker_token):
        raise InvalidAddress(
            'Address for maker token is not in expected binary format in expect_token_swap'
        )
    if not isaddress(maker_address):
        raise InvalidAddress(
            'Address for maker is not in expected binary format in expect_token_swap'
        )

    if not isaddress(taker_token):
        raise InvalidAddress(
            'Address for taker token is not in expected binary format in expect_token_swap'
        )
    if not isaddress(taker_address):
        raise InvalidAddress(
            'Address for taker is not in expected binary format in expect_token_swap'
        )

    channelgraphs = self.raiden.token_to_channelgraph

    if taker_token not in channelgraphs:
        log.error('Unknown token {}'.format(pex(taker_token)))
        return

    if maker_token not in channelgraphs:
        log.error('Unknown token {}'.format(pex(maker_token)))
        return

    # the taker is expecting the maker transfer
    key = SwapKey(
        identifier,
        maker_token,
        maker_amount,
    )

    token_swap = TokenSwap(
        identifier,
        maker_token,
        maker_amount,
        maker_address,
        taker_token,
        taker_amount,
        taker_address,
    )

    self.raiden.swapkey_to_tokenswap[key] = token_swap
# expose a synchronous interface to the user
token_swap = token_swap_and_wait
transfer = transfer_and_wait  # expose a synchronous interface to the user

*/
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
	for k, _ := range this.Raiden.Token2ChannelGraph {
		tokens = append(tokens, k)
	}
	return
}

/*
def transfer_and_wait(
        self,
        token_address,
        amount,
        target,
        identifier=None,
        timeout=None):
    """ Do a transfer with `target` with the given `amount` of `token_address`. """
    # pylint: disable=too-many-arguments

    async_result = self.transfer_async(
        token_address,
        amount,
        target,
        identifier,
    )
    return async_result.wait(timeout=timeout)
*/
//Do a transfer with `target` with the given `amount` of `token_address`.
func (this *RaidenApi) TransferAndWait(token common.Address, amount int64, target common.Address, identifier uint64, timeout time.Duration) (err error) {
	result, err := this.TransferAsync(token, amount, target, identifier)
	if err != nil {
		return err
	}
	if timeout > 0 {
		timeoutCh := time.After(timeout)
		select {
		case <-timeoutCh:
			//?无需关闭? close(timeoutCh)
			err = rerr.TransferTimeout
		case err = <-result.Result:
		}
	} else {
		err = <-result.Result
	}
	if result.Sub != nil { //for mediatedTranfer ,no way to cancel transfer
		result.Sub.Unsubscribe()
	}
	return
}
func (this *RaidenApi) Transfer(token common.Address, amount int64, target common.Address, identifier uint64, timeout time.Duration) error {
	return this.TransferAndWait(token, amount, target, identifier, timeout)
}

func (this *RaidenApi) TransferAsync(tokenAddress common.Address, amount int64, target common.Address, identifier uint64) (result *network.AsyncResult, err error) {
	if amount <= 0 {
		err = rerr.InvalidAmount
		return
	}
	graph, ok := this.Raiden.Token2ChannelGraph[tokenAddress]
	if !ok || graph.HasPath(this.Raiden.NodeAddress, target) {
		err = rerr.NoPathError
		return
	}
	log.Debug(fmt.Sprintf("initiating transfer initiator=%s target=%s token=%s amount=%d identifier=%d",
		this.Raiden.NodeAddress.Str(), target.String(), tokenAddress.String(), amount, identifier))
	result = this.Raiden.MediatedTransferAsync(tokenAddress, amount, target, identifier)
	return
}

//Close a channel opened with `partner_address` for the given `token_address`.
func (this *RaidenApi) Close(tokenAddress, partnerAddress common.Address) (ch *channel.Channel, err error) {
	graph, ok := this.Raiden.Token2ChannelGraph[tokenAddress]
	if !ok {
		err = rerr.InvalidAddress("token address is not valid.")
		return
	}
	ch, ok = graph.PartenerAddress2Channel[partnerAddress]
	if !ok {
		err = rerr.InvalidAddress("no such channel")
		return
	}
	err = ch.ExternState.Close(ch.PartnerState.BalanceProofState)
	return
}

//Settle a closed channel with `partner_address` for the given `token_address`.
func (this *RaidenApi) Settle(tokenAddress, partnerAddress common.Address) (ch *channel.Channel, err error) {
	graph, ok := this.Raiden.Token2ChannelGraph[tokenAddress]
	if !ok {
		err = rerr.InvalidAddress("token address is not valid.")
		return
	}
	ch, ok = graph.PartenerAddress2Channel[partnerAddress]
	if !ok {
		err = rerr.InvalidAddress("no such channel")
		return
	}
	if ch.CanTransfer() {
		err = rerr.InvalidState("channel is still open")
		return
	}
	currentblock := this.Raiden.GetBlockNumber()
	settleTimeout, err := ch.ExternState.NettingChannel.SettleTimeout()
	if err != nil {
		return
	}
	settleExpiration := ch.ExternState.ClosedBlock + settleTimeout
	if currentblock <= settleExpiration {
		err = rerr.InvalidState("settlement period is not yet over.")
		return
	}
	err = ch.ExternState.Settle()
	return
}
func (this *RaidenApi) GetTokenNetworkEvents(tokenAddress common.Address, fromBlock, toBlock int64) (events []transfer.Event, err error) {
	graph, ok := this.Raiden.Token2ChannelGraph[tokenAddress]
	if !ok {
		err = rerr.InvalidAddress("no such token")
	}
	return this.Raiden.BlockChainEvents.GetAllChannelManagerEvents(graph.ChannelManagerAddress, fromBlock, toBlock)
}

func (this *RaidenApi) GetNetworkEvents(fromBlock, toBlock int64) ([]transfer.Event, error) {
	return this.Raiden.BlockChainEvents.GetAllRegistryEvents(this.Raiden.RegistryAddress, fromBlock, toBlock)
}

func (this *RaidenApi) GetChannelEvents(channelAddress common.Address, fromBlock, toBlock int64) (events []transfer.Event, err error) {
	events, err = this.Raiden.BlockChainEvents.GetAllNettingChannelEvents(channelAddress, fromBlock, toBlock)
	if err != nil {
		return
	}
	raidenEvents, err := this.Raiden.TransactionLog.GetEventsInBlockRange(fromBlock, toBlock)
	if err != nil {
		return
	}
	//Here choose which raiden internal events we want to expose to the end user
	for _, ev := range raidenEvents {
		switch ev2 := ev.EventObject.(type) {
		case *transfer.EventTransferSentSuccess:
			e := &EventTransferSentSuccessWrapper{*ev2, ev.BlockNumber, "EventTransferSentSuccess"}
			events = append(events, e)
		case *transfer.EventTransferSentFailed:
			e := &EventTransferSentFailedWrapper{*ev2, ev.BlockNumber, "EventTransferSentFailed"}
			events = append(events, e)
		case *transfer.EventTransferReceivedSuccess:
			e := &EventEventTransferReceivedSuccessWrapper{*ev2, ev.BlockNumber, "EventTransferReceivedSuccess"}
			events = append(events, e)
		}
	}
	return
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
