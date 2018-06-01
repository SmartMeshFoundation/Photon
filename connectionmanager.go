package smartraiden

import (
	"errors"
	"sync"

	"fmt"

	"time"

	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/internal/rpanic"
	"github.com/SmartMeshFoundation/SmartRaiden/network"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fatedier/frp/src/utils/log"
)

//ConnectionManager for connection api
type ConnectionManager struct {
	BootstrapAddr       common.Address //class member
	raiden              *RaidenService
	api                 *RaidenAPI
	lock                sync.Mutex
	tokenAddress        common.Address
	funds               *big.Int
	initChannelTarget   int64
	joinableFundsTarget float64
}

/*
NewConnectionManager if crash,does connection manager need to restore?
*/
func NewConnectionManager(raiden *RaidenService, tokenAddress common.Address) *ConnectionManager {
	cm := &ConnectionManager{
		raiden:              raiden,
		api:                 NewRaidenAPI(raiden),
		tokenAddress:        tokenAddress,
		funds:               utils.BigInt0,
		initChannelTarget:   3,
		joinableFundsTarget: 0.4,
	}
	cm.BootstrapAddr = common.HexToAddress("0x0202020202020202020202020202020202020202")
	return cm
}

/*
Connect to the network.
        Use this to establish a connection with the token network.

        Subsequent calls to `connect` are allowed, but will only affect the spendable
        funds and the connection strategy parameters for the future. `connect` will not
        close any channels.

        Note: the ConnectionManager does not discriminate manually opened channels from
        automatically opened ones. If the user manually opened channels, those deposit
        amounts will affect the funding per channel and the number of new channels opened.

        Args:
            funds (int): the amount of tokens spendable for this
            ConnectionManager.
            initial_channel_target (int): number of channels to open immediately
            joinable_funds_target (float): amount of funds not initially assigned
*/
func (cm *ConnectionManager) Connect(funds *big.Int, initialChannelTarget int64, joinableFundsTarget float64) error {
	if funds.Cmp(utils.BigInt0) <= 0 {
		return errors.New("connecting needs a positive value for `funds`")
	}
	_, ok := cm.raiden.MessageHandler.blockedTokens[cm.tokenAddress]
	if ok { //first leave ,then connect to cm token network
		delete(cm.raiden.MessageHandler.blockedTokens, cm.tokenAddress)
	}
	cm.initChannelTarget = initialChannelTarget
	cm.joinableFundsTarget = joinableFundsTarget
	openChannels := cm.openChannels()
	if len(openChannels) > 0 {
		log.Debug(fmt.Sprintf("connect() called on an already joined token network tokenaddress=%s,openchannels=%d,sumdeposits=%d,funds=%d", utils.APex(cm.tokenAddress), len(openChannels), cm.sumDeposits(), cm.funds))
	}
	chs, err := cm.raiden.db.GetChannelList(cm.tokenAddress, utils.EmptyAddress)
	if err != nil {
		return err
	}
	if len(chs) == 0 {
		log.Debug("bootstrapping token network.")
		cm.lock.Lock()
		_, err2 := cm.api.Open(cm.tokenAddress, cm.BootstrapAddr, cm.raiden.Config.SettleTimeout, cm.raiden.Config.RevealTimeout)
		if err2 != nil {
			log.Error(fmt.Sprintf("open channel between %s and %s error:%s", utils.APex(cm.tokenAddress), utils.APex(cm.BootstrapAddr), err2))
		}
		cm.lock.Unlock()
	}
	cm.lock.Lock()
	cm.funds = funds
	err = cm.addNewPartners()
	cm.lock.Unlock()
	return err
}
func (cm *ConnectionManager) openChannels() []*channel.Serialization {
	chs, _ := cm.api.GetChannelList(cm.tokenAddress, utils.EmptyAddress)
	var chs2 []*channel.Serialization
	for _, c := range chs {
		if c.State == transfer.ChannelStateOpened {
			chs2 = append(chs2, c)
		}
	}
	return chs2
}

//"The calculated funding per partner depending on configuration and
//overall funding of the ConnectionManager.
func (cm *ConnectionManager) initialFundingPerPartner() *big.Int {
	if cm.initChannelTarget > 0 {
		f1 := new(big.Float).SetInt(cm.funds)
		f3 := big.NewFloat(1 - cm.joinableFundsTarget)
		f1.Mul(f1, f3)
		i1, _ := f1.Int(nil)
		return i1.Div(i1, big.NewInt(cm.initChannelTarget))
	}
	return utils.BigInt0
}

/*
WantsMoreChannels returns True, if funds available and the `initial_channel_target` was not yet
        reached.
*/
func (cm *ConnectionManager) WantsMoreChannels() bool {
	_, ok := cm.raiden.MessageHandler.blockedTokens[cm.tokenAddress]
	if ok {
		return false
	}
	return cm.fundsRemaining().Cmp(utils.BigInt0) > 0 && len(cm.openChannels()) < int(cm.initChannelTarget)
}

//The remaining funds after subtracting the already deposited amounts.
func (cm *ConnectionManager) fundsRemaining() *big.Int {
	if cm.funds.Cmp(utils.BigInt0) > 0 {
		remaining := new(big.Int)
		remaining.Sub(cm.funds, cm.sumDeposits())
		return remaining
	}
	return utils.BigInt0
}

//Shorthand for getting sum of all open channels deposited funds
func (cm *ConnectionManager) sumDeposits() *big.Int {
	chs := cm.openChannels()
	var sum = big.NewInt(0)
	for _, c := range chs {
		sum.Add(sum, c.OurContractBalance)
	}
	return sum
}

//Shorthand for getting channels that had received any transfers in this token network
func (cm *ConnectionManager) receivingChannels() (chs []*channel.Serialization) {
	for _, c := range cm.openChannels() {
		if c.PartnerBalanceProof != nil && c.PartnerBalanceProof.Nonce > 0 {
			chs = append(chs, c)
		}
	}
	return
}

//Returns the minimum necessary waiting time to settle all channels.
func (cm *ConnectionManager) minSettleBlocks() int64 {
	chs := cm.receivingChannels()
	var maxTimeout int64 = -1
	currentBlock := cm.raiden.GetBlockNumber()
	for _, c := range chs {
		var sinceClosed int64
		if c.State == transfer.ChannelStateClosed {
			//todo fix cm!
			log.Info(fmt.Sprintf("calc minSettleBlocks need fix:%d", currentBlock))
			// sinceClosed = currentBlock - c.ExternState.ClosedBlock
			sinceClosed = int64(c.SettleTimeout)
		} else if c.State == transfer.ChannelStateOpened {
			sinceClosed = -1
		} else {
			sinceClosed = 0
		}
		t := int64(c.SettleTimeout) - sinceClosed
		if maxTimeout < t {
			maxTimeout = t
		}
	}
	return maxTimeout
}
func (cm *ConnectionManager) leaveState() bool {
	_, ok := cm.raiden.MessageHandler.blockedTokens[cm.tokenAddress]
	return ok || cm.initChannelTarget < 1
}

/*
Close all channels in the token network.
        Note: By default we're just discarding all channels we haven't received anything.
        This potentially leaves deposits locked in channels after `closing`. This is "safe"
        from an accounting point of view (deposits can not be lost), but may still be
        undesirable from a liquidity point of view (deposits will only be freed after
        manually closing or after the partner closed the channel).

        If only_receiving is False then we close and settle all channels irrespective of them
        having received transfers or not.
*/
func (cm *ConnectionManager) closeAll(onlyReceiving bool) []*channel.Serialization {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.initChannelTarget = 0
	var channelsToClose []*channel.Serialization
	if onlyReceiving {
		channelsToClose = cm.receivingChannels()
	} else {
		channelsToClose = cm.openChannels()
	}
	for _, c := range channelsToClose {
		_, err := cm.api.Close(cm.tokenAddress, c.PartnerAddress)
		if err != nil {
			log.Error(fmt.Sprintf("close channel %s error:%s", utils.APex(c.ChannelAddress), err))
		}
	}
	return channelsToClose
}

//LeaveAsync leave raiden network
func (cm *ConnectionManager) LeaveAsync() *network.AsyncResult {
	result := network.NewAsyncResult()
	go func() {
		defer rpanic.PanicRecover("LeaveAsync")
		cm.Leave(true)
		result.Result <- nil
		close(result.Result)
	}()
	return result
}

/*
Leave the token network.
        This implies closing all channels and waiting for all channels to be settled.
*/
func (cm *ConnectionManager) Leave(onlyReceiving bool) []*channel.Serialization {
	cm.raiden.MessageHandler.blockedTokens[cm.tokenAddress] = true
	if cm.initChannelTarget > 0 {
		cm.initChannelTarget = 0
	}
	closedChannels := cm.closeAll(onlyReceiving)
	cm.WaitForSettle(closedChannels)
	return closedChannels
}

/*
WaitForSettle Wait for all closed channels of the token network to settle.
        Note, that this does not time out.
*/
func (cm *ConnectionManager) WaitForSettle(closedChannels []*channel.Serialization) bool {
	found := false
	for {
		found = false
		for _, c := range closedChannels {
			if c.State != transfer.ChannelStateSettled {
				found = true
				break
			}
		}
		if found {
			time.Sleep(time.Minute)
		} else {
			break
		}
	}
	return true
}

/*
Open a channel with `partner` and deposit `funding_amount` tokens.

        If the channel was already opened (a known race condition),
        this skips the opening and only deposits.
*/
func (cm *ConnectionManager) openAndDeposit(partner common.Address, fundingAmount *big.Int) error {
	_, err := cm.api.Open(cm.tokenAddress, partner, cm.raiden.Config.SettleTimeout, cm.raiden.Config.RevealTimeout)
	if err != nil {
		return err
	}
	ch, err := cm.raiden.db.GetChannel(cm.tokenAddress, partner)
	if err != nil {
		return err
	}
	if ch == nil {
		err = fmt.Errorf("Opening new channel failed; channel already opened,  but partner not in channelgraph ,partner=%s,tokenaddress=%s", utils.APex(partner), utils.APex(cm.tokenAddress))
		log.Error(err.Error())
		return err
	}
	err = cm.api.Deposit(cm.tokenAddress, partner, fundingAmount, params.DefaultPollTimeout)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

/*
This opens channels with a number of new partners according to the
        connection strategy parameter `self.initial_channel_target`.
        Each new channel will receive `self.initial_funding_per_partner` funding.
*/
func (cm *ConnectionManager) addNewPartners() error {
	newPartnerCount := int(cm.initChannelTarget) - len(cm.openChannels())
	if newPartnerCount <= 0 {
		return nil
	}
	for _, partner := range cm.findNewPartners(newPartnerCount) {
		err := cm.openAndDeposit(partner, cm.initialFundingPerPartner())
		if err != nil {
			log.Error(fmt.Sprintf("addNewPartners %s ,err:%s", utils.APex(partner), err))
			return err
		}
	}
	return nil
}

/*
RetryConnect Will be called when new channels in the token network are detected.
        If the minimum number of channels was not yet established, it will try
        to open new channels.

        If the connection manager has no funds, this is a noop.
*/
func (cm *ConnectionManager) RetryConnect() {
	if cm.funds.Cmp(utils.BigInt0) <= 0 {
		return
	}
	if cm.leaveState() {
		return
	}
	cm.lock.Lock()
	defer cm.lock.Unlock()
	if cm.fundsRemaining().Cmp(utils.BigInt0) <= 0 {
		return
	}
	if len(cm.openChannels()) >= int(cm.initChannelTarget) {
		return
	}
	//try to fullfill our connection goal
	cm.addNewPartners()
}

/*
JoinChannel Will be called, when we were selected as channel partner by another
        node. It will fund the channel with up to the partner's deposit, but
        not more than remaining funds or the initial funding per channel.

        If the connection manager has no funds, this is a noop.
*/
func (cm *ConnectionManager) JoinChannel(partnerAddress common.Address, partnerDepost *big.Int) {
	if cm.funds.Cmp(utils.BigInt0) <= 0 {
		return
	}
	if cm.leaveState() {
		return
	}
	cm.lock.Lock()
	defer cm.lock.Unlock()
	remaining := cm.fundsRemaining()
	initial := cm.initialFundingPerPartner()
	joiningFunds := partnerDepost
	if joiningFunds.Cmp(remaining) > 0 {
		joiningFunds = remaining
	}
	if joiningFunds.Cmp(initial) > 0 {
		joiningFunds = initial
	}
	if joiningFunds.Cmp(utils.BigInt0) <= 0 {
		return
	}
	err := cm.api.Deposit(cm.tokenAddress, partnerAddress, joiningFunds, params.DefaultPollTimeout)
	log.Debug("joined a channel funds=%d,me=%s,partner=%s err=%s", joiningFunds, utils.APex(cm.raiden.NodeAddress), utils.APex(partnerAddress), err)
	return
}

/*
Search the token network for potential channel partners.

        Args:
            number (int): number of partners to return
*/
func (cm *ConnectionManager) findNewPartners(number int) []common.Address {
	var known = make(map[common.Address]bool)
	for _, c := range cm.openChannels() {
		known[c.PartnerAddress] = true
	}
	known[cm.BootstrapAddr] = true
	known[cm.raiden.NodeAddress] = true
	channelAddresses := cm.raiden.db.GetTokenNodes(cm.tokenAddress)
	var availables []common.Address
	for _, addr := range channelAddresses {
		if !known[addr] {
			availables = append(availables, addr)
		}
	}
	log.Debug(fmt.Sprintf("found %d partners", len(availables)))
	if number < len(availables) {
		return availables[:number]
	}
	return availables
}
