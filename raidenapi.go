package smartraiden

import (
	"context"
	"encoding/binary"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/channel"

	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mtree"

	"github.com/SmartMeshFoundation/SmartRaiden/params"

	"fmt"

	"math/big"

	"sync"

	"errors"

	"bytes"
	"crypto/ecdsa"

	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/models"
	"github.com/SmartMeshFoundation/SmartRaiden/rerr"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var errEthConnectionNotReady = errors.New("eth connection not ready")

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
	tokens, err := r.Raiden.db.GetAllTokens()
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
func (r *RaidenAPI) GetChannelList(tokenAddress common.Address, partnerAddress common.Address) (cs []*channeltype.Serialization, err error) {
	return r.Raiden.db.GetChannelList(tokenAddress, partnerAddress)
}

//GetChannel get channel by address
func (r *RaidenAPI) GetChannel(ChannelIdentifier common.Hash) (c *channeltype.Serialization, err error) {
	return r.Raiden.db.GetChannelByAddress(ChannelIdentifier)
}

/*
TokenAddressIfTokenRegistered return the channel manager address,If the token is registered then
Also make sure that the channel manager is registered with the node.
*/
func (r *RaidenAPI) TokenAddressIfTokenRegistered(tokenAddress common.Address) (mgrAddr common.Address, err error) {
	if r.Raiden.Registry == nil {
		err = errEthConnectionNotReady
		return
	}
	mgrAddr, err = r.Raiden.Registry.TokenNetworkByToken(tokenAddress)
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
	if r.Raiden.Registry == nil {
		err = errEthConnectionNotReady
		return
	}
	mgrAddr, err = r.Raiden.Registry.TokenNetworkByToken(tokenAddress)
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
Open a channel with the peer at `partner_address`
    with the given `token_address`.
*/
func (r *RaidenAPI) Open(tokenAddress, partnerAddress common.Address, settleTimeout, revealTimeout int, deposit *big.Int) (ch *channeltype.Serialization, err error) {
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
	r.Raiden.db.RegisterNewChannellCallback(func(c *channeltype.Serialization) (remove bool) {
		if c.TokenAddress() == tokenAddress && c.PartnerAddress() == partnerAddress {
			wg.Done()
			return true
		}
		return false
	})
	result := r.Raiden.newChannelClient(tokenAddress, partnerAddress, settleTimeout, deposit)
	err = <-result.Result
	if err != nil {
		return
	}
	//wait
	wg.Wait()
	ch, err = r.Raiden.db.GetChannel(tokenAddress, partnerAddress)
	if err == nil {
		//must be success, no need to wait event and register a callback
		if deposit != nil {
			ch.OurContractBalance = deposit
		} else {
			ch.OurContractBalance = big.NewInt(0)
		}
	}
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
func (r *RaidenAPI) Deposit(tokenAddress, partnerAddress common.Address, amount *big.Int, pollTimeout time.Duration) (c *channeltype.Serialization, err error) {
	c, err = r.Raiden.db.GetChannel(tokenAddress, partnerAddress)
	if err != nil {
		return
	}
	token, err := r.Raiden.Chain.Token(tokenAddress)
	if err != nil {
		return
	}
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
		err = fmt.Errorf("not enough balance to deposit. %s Available=%d Tried=%d", tokenAddress.String(), balance, amount)
		log.Error(err.Error())
		err = rerr.ErrInsufficientFunds
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	r.Raiden.db.RegisterChannelDepositCallback(func(c2 *channeltype.Serialization) (remove bool) {
		if bytes.Equal(c2.Key, c.Key) {
			wg.Done()
			return true
		}
		return false
	})
	//deposit move ... todo
	result := r.Raiden.depositChannelClient(c.ChannelIdentifier.ChannelIdentifier, amount)
	err = <-result.Result
	if err != nil {
		return
	}
	/*
	 Wait until the `ChannelNewBalance` event is processed.
	*/
	wg.Wait()
	//reload data from database,
	return r.Raiden.db.GetChannelByAddress(c.ChannelIdentifier.ChannelIdentifier)
}

/*
TokenSwapAndWait Start an atomic swap operation by sending a MediatedTransfer with
    `maker_amount` of `maker_token` to `taker_address`. Only proceed when a
    new valid MediatedTransfer is received with `taker_amount` of
    `taker_token`.
*/
func (r *RaidenAPI) TokenSwapAndWait(lockSecretHash string, makerToken, takerToken, makerAddress, takerAddress common.Address,
	makerAmount, takerAmount *big.Int, secret string) error {
	result, err := r.tokenSwapAsync(lockSecretHash, makerToken, takerToken, makerAddress, takerAddress,
		makerAmount, takerAmount, secret)
	if err != nil {
		return err
	}
	err = <-result.Result
	return err
}

func (r *RaidenAPI) tokenSwapAsync(lockSecretHash string, makerToken, takerToken, makerAddress, takerAddress common.Address,
	makerAmount, takerAmount *big.Int, secret string) (result *utils.AsyncResult, err error) {
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
		LockSecretHash:  common.HexToHash(lockSecretHash),
		Secret:          common.HexToHash(secret),
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
func (r *RaidenAPI) ExpectTokenSwap(lockSecretHash string, makerToken, takerToken, makerAddress, takerAddress common.Address,
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
		LockSecretHash:  common.HexToHash(lockSecretHash),
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
	tokensmap, err := r.Raiden.db.GetAllTokens()
	if err != nil {
		log.Error(fmt.Sprintf("GetAllTokens err %s", err))
	}
	for k := range tokensmap {
		tokens = append(tokens, k)
	}
	return
}

//GetTokenTokenNetorks return all tokens and token networks
func (r *RaidenAPI) GetTokenTokenNetorks() (tokens []string) {
	tokenMap, err := r.Raiden.db.GetAllTokens()
	if err != nil {
		log.Error(fmt.Sprintf("GetAllTokens err %s", err))
	}
	for k := range tokenMap {
		tokens = append(tokens, k.String())
	}
	return
}

//TransferAndWait Do a transfer with `target` with the given `amount` of `token_address`.
func (r *RaidenAPI) TransferAndWait(token common.Address, amount *big.Int, fee *big.Int, target common.Address, secret common.Hash, timeout time.Duration, isDirectTransfer bool) (err error) {
	result, err := r.transferAsync(token, amount, fee, target, secret, isDirectTransfer)
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
func (r *RaidenAPI) Transfer(token common.Address, amount *big.Int, fee *big.Int, target common.Address, secret common.Hash, timeout time.Duration, isDirectTransfer bool) error {
	return r.TransferAndWait(token, amount, fee, target, secret, timeout, isDirectTransfer)
}

//transferAsync
func (r *RaidenAPI) transferAsync(tokenAddress common.Address, amount *big.Int, fee *big.Int, target common.Address, secret common.Hash, isDirectTransfer bool) (result *utils.AsyncResult, err error) {
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
		var c *channeltype.Serialization
		c, err = r.Raiden.db.GetChannel(tokenAddress, target)
		if err != nil {
			err = fmt.Errorf("no direct channel token:%s,partner:%s", tokenAddress.String(), target.String())
			return
		}
		if c.State != channeltype.StateOpened {
			err = fmt.Errorf("channel %s not opened", c.ChannelIdentifier)
			return
		}
	}
	if amount.Cmp(utils.BigInt0) <= 0 {
		err = rerr.ErrInvalidAmount
		return
	}
	log.Debug(fmt.Sprintf("initiating transfer initiator=%s target=%s token=%s amount=%d secret=%s",
		r.Raiden.NodeAddress.String(), target.String(), tokenAddress.String(), amount, secret.String()))
	result = r.Raiden.transferAsyncClient(tokenAddress, amount, fee, target, secret, isDirectTransfer)
	return
}

// AllowRevealSecret :
// 1. find state manager by lockSecretHash and tokenAddress
// 2. check secret matches lockSecretHash or not
// 3. remove the predictor
func (r *RaidenAPI) AllowRevealSecret(lockSecretHash common.Hash, tokenAddress common.Address) (err error) {
	key := utils.Sha3(lockSecretHash[:], tokenAddress[:])
	manager := r.Raiden.Transfer2StateManager[key]
	if manager == nil {
		return rerr.InvalidState("can not find transfer by lock_secret_hash and token_address")
	}
	state, ok := manager.CurrentState.(*mediatedtransfer.InitiatorState)
	if !ok {
		return rerr.InvalidState("wrong state")
	}
	if lockSecretHash != state.LockSecretHash || lockSecretHash != utils.ShaSecret(state.Secret.Bytes()) {
		return rerr.InvalidState("wrong lock_secret_hash")
	}
	delete(r.Raiden.SecretRequestPredictorMap, lockSecretHash)
	log.Trace(fmt.Sprintf("Remove SecretRequestPredictor for lockSecretHash="))
	return
}

// RegisterSecret :
func (r *RaidenAPI) RegisterSecret(secret common.Hash, tokenAddress common.Address) (err error) {
	lockSecretHash := utils.ShaSecret(secret.Bytes())
	//在channel 中注册密码
	// register secret in channel
	r.Raiden.registerSecret(secret)

	key := utils.Sha3(lockSecretHash[:], tokenAddress[:])
	manager := r.Raiden.Transfer2StateManager[key]
	if manager == nil {
		return rerr.InvalidState("can not find transfer by lock_secret_hash and token_address")
	}
	state, ok := manager.CurrentState.(*mediatedtransfer.TargetState)
	if !ok {
		return rerr.InvalidState("wrong state")
	}
	if lockSecretHash != state.FromTransfer.LockSecretHash {
		return rerr.InvalidState("wrong secret")
	}
	// 在state manager中注册密码
	// register secret in state manager
	state.FromTransfer.Secret = secret
	state.Secret = secret
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
func (r *RaidenAPI) GetUnfinishedReceivedTransfer(lockSecretHash common.Hash, tokenAddress common.Address) (resp *TransferDataResponse) {

	if r.Raiden.SecretRequestPredictorMap[lockSecretHash] != nil {
		return
	}
	key := utils.Sha3(lockSecretHash[:], tokenAddress[:])
	manager := r.Raiden.Transfer2StateManager[key]
	if manager == nil {
		log.Warn(fmt.Sprintf("can not find transfer by lock_secret_hash[%s] and token_address[%s]", lockSecretHash.String(), tokenAddress.String()))
		return
	}
	state, ok := manager.CurrentState.(*mediatedtransfer.TargetState)
	if !ok {
		// 接收人不是自己
		// I'm not the recipient
		return
	}
	resp = new(TransferDataResponse)
	resp.Initiator = state.FromTransfer.Initiator.String()
	resp.Target = state.FromTransfer.Target.String()
	resp.Token = tokenAddress.String()
	resp.Amount = state.FromTransfer.Amount
	resp.LockSecretHash = state.FromTransfer.LockSecretHash.String()
	resp.Expiration = state.FromTransfer.Expiration - state.BlockNumber
	return
}

//Close a channel opened with `partner_address` for the given `token_address`. return when state has been updated to database
func (r *RaidenAPI) Close(tokenAddress, partnerAddress common.Address) (c *channeltype.Serialization, err error) {
	c, err = r.Raiden.db.GetChannel(tokenAddress, partnerAddress)
	if err != nil {
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	r.Raiden.db.RegisterChannelStateCallback(func(c2 *channeltype.Serialization) (remove bool) {
		log.Trace(fmt.Sprintf("wait %s closed ,get channle %s update",
			c.ChannelIdentifier, c2.ChannelIdentifier))
		if bytes.Equal(c2.Key, c.Key) {
			wg.Done()
			return true
		}
		return false
	})
	//send close channel request
	result := r.Raiden.closeChannelClient(c.ChannelIdentifier.ChannelIdentifier)
	err = <-result.Result
	if err != nil {
		return
	}
	wg.Wait()
	//reload data from database,
	return r.Raiden.db.GetChannelByAddress(c.ChannelIdentifier.ChannelIdentifier)
}

//Settle a closed channel with `partner_address` for the given `token_address`.return when state has been updated to database
func (r *RaidenAPI) Settle(tokenAddress, partnerAddress common.Address) (c *channeltype.Serialization, err error) {
	c, err = r.Raiden.db.GetChannel(tokenAddress, partnerAddress)
	if c.State == channeltype.StateOpened {
		err = rerr.InvalidState("channel is still open")
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	r.Raiden.db.RegisterChannelSettleCallback(func(c2 *channeltype.Serialization) (remove bool) {
		log.Trace(fmt.Sprintf("wait %s settled ,get channle %s update",
			c.ChannelIdentifier, c2.ChannelIdentifier))
		if bytes.Equal(c2.Key, c.Key) {
			wg.Done()
			return true
		}
		return false
	})
	//send settle request
	result := r.Raiden.settleChannelClient(c.ChannelIdentifier.ChannelIdentifier)
	err = <-result.Result
	log.Trace(fmt.Sprintf("%s settled finish , err %v", c.ChannelIdentifier, err))
	if err != nil {
		return
	}
	wg.Wait()
	//reload data from database, this channel has been removed.
	return r.Raiden.db.GetSettledChannel(c.ChannelIdentifier.ChannelIdentifier, c.ChannelIdentifier.OpenBlockNumber)
}

//CooperativeSettle a channel opened with `partner_address` for the given `token_address`. return when state has been updated to database
func (r *RaidenAPI) CooperativeSettle(tokenAddress, partnerAddress common.Address) (c *channeltype.Serialization, err error) {
	c, err = r.Raiden.db.GetChannel(tokenAddress, partnerAddress)
	if c.State != channeltype.StateOpened && c.State != channeltype.StatePrepareForCooperativeSettle {
		err = rerr.InvalidState("channel must be  open")
		return
	}
	//send settle request
	result := r.Raiden.cooperativeSettleChannelClient(c.ChannelIdentifier.ChannelIdentifier)
	err = <-result.Result
	log.Trace(fmt.Sprintf("%s settled finish , err %v", c.ChannelIdentifier, err))
	if err != nil {
		return
	}
	//reload data from database, this channel has been removed.
	return r.Raiden.db.GetChannelByAddress(c.ChannelIdentifier.ChannelIdentifier)
}

//PrepareForCooperativeSettle  mark a channel prepared for settle,  return when state has been updated to database
func (r *RaidenAPI) PrepareForCooperativeSettle(tokenAddress, partnerAddress common.Address) (c *channeltype.Serialization, err error) {
	c, err = r.Raiden.db.GetChannel(tokenAddress, partnerAddress)
	if c.State != channeltype.StateOpened {
		err = rerr.InvalidState("channel must be  open")
		return
	}
	//send settle request
	result := r.Raiden.markChannelForCooperativeSettleClient(c.ChannelIdentifier.ChannelIdentifier)
	err = <-result.Result
	log.Trace(fmt.Sprintf("%s PrepareForCooperativeSettle finish , err %v", c.ChannelIdentifier, err))
	if err != nil {
		return
	}
	//reload data from database, this channel has been removed.
	return r.Raiden.db.GetChannelByAddress(c.ChannelIdentifier.ChannelIdentifier)
}

//CancelPrepareForCooperativeSettle  cancel a mark. return when state has been updated to database
func (r *RaidenAPI) CancelPrepareForCooperativeSettle(tokenAddress, partnerAddress common.Address) (c *channeltype.Serialization, err error) {
	c, err = r.Raiden.db.GetChannel(tokenAddress, partnerAddress)
	if c.State != channeltype.StatePrepareForCooperativeSettle {
		err = rerr.InvalidState("channel must be  open")
		return
	}
	//send settle request
	result := r.Raiden.cancelMarkChannelForCooperativeSettleClient(c.ChannelIdentifier.ChannelIdentifier)
	err = <-result.Result
	log.Trace(fmt.Sprintf("%s CancelPrepareForCooperativeSettle finish , err %v", c.ChannelIdentifier, err))
	if err != nil {
		return
	}
	//reload data from database, this channel has been removed.
	return r.Raiden.db.GetChannelByAddress(c.ChannelIdentifier.ChannelIdentifier)
}

//Withdraw on a channel opened with `partner_address` for the given `token_address`. return when state has been updated to database
func (r *RaidenAPI) Withdraw(tokenAddress, partnerAddress common.Address, amount *big.Int) (c *channeltype.Serialization, err error) {
	c, err = r.Raiden.db.GetChannel(tokenAddress, partnerAddress)
	if c.State != channeltype.StateOpened && c.State != channeltype.StatePrepareForWithdraw {
		err = rerr.InvalidState("channel must be  open")
		return
	}
	if c.OurBalance().Cmp(amount) < 0 {
		err = fmt.Errorf("invalid withdraw amount, availabe=%s,want=%s", c.OurBalance(), amount)
		return
	}
	//send settle request
	result := r.Raiden.withdrawClient(c.ChannelIdentifier.ChannelIdentifier, amount)
	err = <-result.Result
	log.Trace(fmt.Sprintf("%s withdraw finish , err %v", c.ChannelIdentifier, err))
	if err != nil {
		return
	}
	//reload data from database, this channel has been removed.
	return r.Raiden.db.GetChannelByAddress(c.ChannelIdentifier.ChannelIdentifier)
}

//PrepareForWithdraw  mark a channel prepared for withdraw,  return when state has been updated to database
func (r *RaidenAPI) PrepareForWithdraw(tokenAddress, partnerAddress common.Address) (c *channeltype.Serialization, err error) {
	c, err = r.Raiden.db.GetChannel(tokenAddress, partnerAddress)
	if c.State != channeltype.StateOpened {
		err = rerr.InvalidState("channel must be  open")
		return
	}
	//send settle request
	result := r.Raiden.markWithdraw(c.ChannelIdentifier.ChannelIdentifier)
	err = <-result.Result
	log.Trace(fmt.Sprintf("%s PrepareForWithdraw finish , err %v", c.ChannelIdentifier, err))
	if err != nil {
		return
	}
	//reload data from database, this channel has been removed.
	return r.Raiden.db.GetChannelByAddress(c.ChannelIdentifier.ChannelIdentifier)
}

//CancelPrepareForWithdraw  cancel a mark. return when state has been updated to database
func (r *RaidenAPI) CancelPrepareForWithdraw(tokenAddress, partnerAddress common.Address) (c *channeltype.Serialization, err error) {
	c, err = r.Raiden.db.GetChannel(tokenAddress, partnerAddress)
	if c.State != channeltype.StatePrepareForWithdraw {
		err = rerr.InvalidState("channel must be  open")
		return
	}
	//send settle request
	result := r.Raiden.cancelMarkWithdraw(c.ChannelIdentifier.ChannelIdentifier)
	err = <-result.Result
	log.Trace(fmt.Sprintf("%s CancelPrepareForWithdraw finish , err %v", c.ChannelIdentifier, err))
	if err != nil {
		return
	}
	//reload data from database, this channel has been removed.
	return r.Raiden.db.GetChannelByAddress(c.ChannelIdentifier.ChannelIdentifier)
}

//GetTokenNetworkEvents return events about this token
func (r *RaidenAPI) GetTokenNetworkEvents(tokenAddress common.Address, fromBlock, toBlock int64) (data []interface{}, err error) {
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
	//tokens, err := r.Raiden.db.GetAllTokens()
	//if err != nil {
	//	return
	//}
	//for t, manager := range tokens {
	//	if tokenAddress == utils.EmptyAddress || t == tokenAddress {
	//		events, err := r.Raiden.BlockChainEvents.GetAllChannelManagerEvents(manager, fromBlock, toBlock)
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

//GetNetworkEvents all raiden events
func (r *RaidenAPI) GetNetworkEvents(fromBlock, toBlock int64) ([]interface{}, error) {
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
	//events, err := r.Raiden.BlockChainEvents.GetAllRegistryEvents(r.Raiden.RegistryAddress, fromBlock, toBlock)
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
func (r *RaidenAPI) GetChannelEvents(channelAddress common.Hash, fromBlock, toBlock int64) (data []transfer.Event, err error) {

	//var events []transfer.Event
	//events, err = r.Raiden.BlockChainEvents.GetAllNettingChannelEvents(channelAddress, fromBlock, toBlock)
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
	//var raidenEvents []*models.InternalEvent
	//raidenEvents, err = r.Raiden.db.GetEventsInBlockRange(fromBlock, toBlock)
	//if err != nil {
	//	return
	//}
	////Here choose which raiden internal events we want to expose to the end user
	//for _, ev := range raidenEvents {
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

type updateTransfer struct {
	Nonce               uint64      `json:"nonce"`
	TransferAmount      *big.Int    `json:"transfer_amount"`
	Locksroot           common.Hash `json:"locksroot"`
	ExtraHash           common.Hash `json:"extra_hash"`
	ClosingSignature    []byte      `json:"closing_signature"`
	NonClosingSignature []byte      `json:"non_closing_signature"`
}

//todo 需要第三方服务帮忙注册密码么?如果不需要,是否应该自己注册密码?
// todo do we need delegation service to help us register secret? If not, should we register secret in person?
type unlock struct {
	Lock        *mtree.Lock `json:"lock"`
	MerkleProof []byte      `json:"merkle_proof"`
	Secret      common.Hash `json:"secret"`
	Signature   []byte      `json:"signature"`
}

//需要委托给第三方的 punish证据
// punish proof that is delegated to third-party.
type punish struct {
	LockHash       common.Hash `json:"lock_hash"` //the whole lock's hash,not lock secret hash
	AdditionalHash common.Hash `json:"additional_hash"`
	Signature      []byte      `json:"signature"`
}

//ChannelFor3rd is for 3rd party to call update transfer
type ChannelFor3rd struct {
	ChannelIdentifier  common.Hash    `json:"channel_identifier"`
	OpenBlockNumber    int64          `json:"open_block_number"`
	TokenNetworkAddrss common.Address `json:"token_network_address"`
	PartnerAddress     common.Address `json:"partner_address"`
	UpdateTransfer     updateTransfer `json:"update_transfer"`
	Unlocks            []*unlock      `json:"unlocks"`
	Punishes           []*punish      `json:"punishes"`
}

/*
ChannelInformationFor3rdParty generate all information need by 3rd party
*/
func (r *RaidenAPI) ChannelInformationFor3rdParty(ChannelIdentifier common.Hash, thirdAddr common.Address) (result *ChannelFor3rd, err error) {
	var sig []byte
	c, err := r.GetChannel(ChannelIdentifier)
	if err != nil {
		return
	}
	c3 := new(ChannelFor3rd)
	c3.ChannelIdentifier = ChannelIdentifier
	c3.OpenBlockNumber = c.ChannelIdentifier.OpenBlockNumber
	c3.TokenNetworkAddrss = r.Raiden.Token2TokenNetwork[c.TokenAddress()]
	c3.PartnerAddress = c.PartnerAddress()
	if c.PartnerBalanceProof == nil {
		result = c3
		return
	}
	if c.PartnerBalanceProof.Nonce > 0 {
		c3.UpdateTransfer.Nonce = c.PartnerBalanceProof.Nonce
		c3.UpdateTransfer.TransferAmount = c.PartnerBalanceProof.TransferAmount
		c3.UpdateTransfer.Locksroot = c.PartnerBalanceProof.LocksRoot
		c3.UpdateTransfer.ExtraHash = c.PartnerBalanceProof.MessageHash
		c3.UpdateTransfer.ClosingSignature = c.PartnerBalanceProof.Signature
		sig, err = signBalanceProofFor3rd(c, r.Raiden.PrivateKey)
		if err != nil {
			return
		}
		c3.UpdateTransfer.NonClosingSignature = sig
	}

	tree := mtree.NewMerkleTree(c.PartnerLeaves)
	var ws []*unlock
	for _, l := range c.PartnerLock2UnclaimedLocks() {
		proof := channel.ComputeProofForLock(l.Lock, tree)
		w := &unlock{
			Lock:        l.Lock,
			Secret:      l.Secret,
			MerkleProof: mtree.Proof2Bytes(proof.MerkleProof),
		}
		w.Signature, err = signUnlockFor3rd(c, w, thirdAddr, r.Raiden.PrivateKey)
		log.Trace(fmt.Sprintf("prootf=%s", utils.StringInterface(proof, 3)))
		ws = append(ws, w)
	}
	c3.Unlocks = ws
	var ps []*punish
	for _, annouceDisposed := range r.Raiden.db.GetChannelAnnounceDisposed(c.ChannelIdentifier.ChannelIdentifier) {
		//跳过历史 channel
		// omit history channel
		if annouceDisposed.OpenBlockNumber != c.ChannelIdentifier.OpenBlockNumber {
			continue
		}
		p := &punish{
			LockHash:       common.BytesToHash(annouceDisposed.LockHash),
			AdditionalHash: annouceDisposed.AdditionalHash,
			Signature:      annouceDisposed.Signature,
		}
		ps = append(ps, p)
	}
	c3.Punishes = ps
	result = c3
	return
}

//make sure PartnerBalanceProof is not nil
func signBalanceProofFor3rd(c *channeltype.Serialization, privkey *ecdsa.PrivateKey) (sig []byte, err error) {
	if c.PartnerBalanceProof == nil {
		log.Error(fmt.Sprintf("PartnerBalanceProof is nil,must ber a error"))
		return nil, errors.New("empty PartnerBalanceProof")
	}
	buf := new(bytes.Buffer)
	_, err = buf.Write(params.ContractSignaturePrefix)
	_, err = buf.Write([]byte(params.ContractBalanceProofDelegateMessageLength))
	_, err = buf.Write(utils.BigIntTo32Bytes(c.PartnerBalanceProof.TransferAmount))
	_, err = buf.Write(c.PartnerBalanceProof.LocksRoot[:])
	err = binary.Write(buf, binary.BigEndian, c.PartnerBalanceProof.Nonce)
	_, err = buf.Write(c.ChannelIdentifier.ChannelIdentifier[:])
	err = binary.Write(buf, binary.BigEndian, c.ChannelIdentifier.OpenBlockNumber)
	_, err = buf.Write(utils.BigIntTo32Bytes(params.ChainID))
	if err != nil {
		log.Error(fmt.Sprintf("buf write error %s", err))
	}
	dataToSign := buf.Bytes()
	return utils.SignData(privkey, dataToSign)
}

func signUnlockFor3rd(c *channeltype.Serialization, u *unlock, thirdAddress common.Address, privkey *ecdsa.PrivateKey) (sig []byte, err error) {
	buf := new(bytes.Buffer)
	_, err = buf.Write(params.ContractSignaturePrefix)
	_, err = buf.Write([]byte(params.ContractUnlockDelegateProofMessageLength))
	_, err = buf.Write(utils.BigIntTo32Bytes(c.PartnerBalanceProof.TransferAmount))
	_, err = buf.Write(thirdAddress[:])
	_, err = buf.Write(utils.BigIntTo32Bytes(big.NewInt(u.Lock.Expiration)))
	_, err = buf.Write(utils.BigIntTo32Bytes(u.Lock.Amount))
	_, err = buf.Write(u.Lock.LockSecretHash[:])
	_, err = buf.Write(c.ChannelIdentifier.ChannelIdentifier[:])
	err = binary.Write(buf, binary.BigEndian, c.ChannelIdentifier.OpenBlockNumber)
	_, err = buf.Write(utils.BigIntTo32Bytes(params.ChainID))
	if err != nil {
		log.Error(fmt.Sprintf("buf write error %s", err))
		return
	}
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

// AccountTokenBalanceVo for api
type AccountTokenBalanceVo struct {
	TokenAddress string   `json:"token_address"`
	Balance      *big.Int `json:"balance"`
	LockedAmount *big.Int `json:"locked_amount"`
}

// GetBalanceByTokenAddress : get account's balance and locked account on token
func (r *RaidenAPI) GetBalanceByTokenAddress(tokenAddress common.Address) (balances []*AccountTokenBalanceVo, err error) {
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
		err = errors.New("token not registered")
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
func (r *RaidenAPI) getBalance() (balances []*AccountTokenBalanceVo, err error) {
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
func (r *RaidenAPI) ForceUnlock(channelIdentifier common.Hash, lockSecretHash common.Hash, secretHash common.Hash) (err error) {
	channel := r.Raiden.getChannelWithAddr(channelIdentifier)
	tokenNetwork, err := r.Raiden.Chain.TokenNetwork(channel.TokenAddress)
	if err != nil {
		return
	}
	auth := bind.NewKeyedTransactor(r.Raiden.PrivateKey)
	partnerAddress := channel.PartnerState.Address
	tr, err := channel.CreateUnlock(lockSecretHash)
	if err != nil {
		return
	}
	lock := channel.PartnerState.Lock2PendingLocks[lockSecretHash]
	expiration := big.NewInt(lock.Lock.Expiration)
	proof := channel.PartnerState.Tree.MakeProof(lock.LockHash)

	// unlock
	tx, err := tokenNetwork.GetContract().Unlock(auth, partnerAddress, tr.TransferAmount,
		expiration, lock.Lock.Amount, secretHash, mtree.Proof2Bytes(proof))
	if err != nil {
		return
	}
	log.Info(fmt.Sprintf("ForceUnlock  txhash=%s", tx.Hash().String()))
	receipt, err := bind.WaitMined(context.Background(), r.Raiden.Chain.Client, tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Info(fmt.Sprintf("ForceUnlock failed %s", receipt))
		return errors.New("ForceUnlock tx execution failed")
	}
	log.Info(fmt.Sprintf("ForceUnlock success %s ,partner=%s", lockSecretHash.String(), utils.APex(partnerAddress)))
	return nil
}
