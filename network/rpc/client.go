package rpc

import (
	"context"

	"math/big"

	"time"

	"errors"

	"fmt"

	"sync"

	"crypto/ecdsa"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/rerr"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

//GetCallContext context for tx
func GetCallContext() context.Context {
	ctx, cf := context.WithDeadline(context.Background(), time.Now().Add(params.DefaultTxTimeout))
	if cf != nil {
	}
	return ctx
}

//GetQueryConext context for query on chain
func GetQueryConext() context.Context {
	ctx, cf := context.WithDeadline(context.Background(), time.Now().Add(params.DefaultPollTimeout))
	if cf != nil {
	}
	return ctx
}

/*
BlockChainService provides quering on blockchain.
*/
type BlockChainService struct {
	//PrivKey of this node, todo remove this
	PrivKey *ecdsa.PrivateKey
	//NodeAddress is address of this node
	NodeAddress common.Address
	//RegistryAddress registy contract address
	RegistryAddress common.Address
	//Client if eth rpc client
	Client            *helper.SafeEthClient
	addressTokens     map[common.Address]*TokenProxy
	addressChannels   map[common.Address]*TokenNetworkProxy
	addressRegistries map[common.Address]*RegistryProxy
	//Auth needs by call on blockchain todo remove this
	Auth      *bind.TransactOpts
	queryOpts *bind.CallOpts
}

//NewBlockChainService create BlockChainService
func NewBlockChainService(privKey *ecdsa.PrivateKey, registryAddress common.Address, client *helper.SafeEthClient) *BlockChainService {
	bcs := &BlockChainService{
		PrivKey:           privKey,
		NodeAddress:       crypto.PubkeyToAddress(privKey.PublicKey),
		RegistryAddress:   registryAddress,
		Client:            client,
		addressTokens:     make(map[common.Address]*TokenProxy),
		addressChannels:   make(map[common.Address]*TokenNetworkProxy),
		addressRegistries: make(map[common.Address]*RegistryProxy),
		Auth:              bind.NewKeyedTransactor(privKey),
	}
	bcs.queryOpts = &bind.CallOpts{
		Pending: false,
		From:    bcs.NodeAddress,
		Context: GetQueryConext(),
	}
	//It needs to be set up, otherwise, even the contract revert will not report wrong.
	bcs.Auth.GasLimit = uint64(params.GasLimit)
	bcs.Auth.GasPrice = big.NewInt(params.GasPrice)
	return bcs
}
func (bcs *BlockChainService) getQueryOpts() *bind.CallOpts {
	return &bind.CallOpts{
		Pending: false,
		From:    bcs.NodeAddress,
		Context: GetQueryConext(),
	}
}
func (bcs *BlockChainService) blockNumber() (num *big.Int, err error) {
	h, err := bcs.Client.HeaderByNumber(context.Background(), nil)
	if err == nil {
		num = h.Number
	}
	return
}
func (bcs *BlockChainService) nonce(account common.Address) (uint64, error) {
	return bcs.Client.PendingNonceAt(context.Background(), account)
}
func (bcs *BlockChainService) balance(account common.Address) (*big.Int, error) {
	return bcs.Client.PendingBalanceAt(context.Background(), account)
}
func (bcs *BlockChainService) contractExist(contractAddr common.Address) bool {
	code, err := bcs.Client.CodeAt(context.Background(), contractAddr, nil)
	//spew.Dump("code:", code)
	return err == nil && len(code) > 0
}
func (bcs *BlockChainService) getBlockHeader(blockNumber *big.Int) (*types.Header, error) {
	return bcs.Client.HeaderByNumber(context.Background(), blockNumber)
}

func (bcs *BlockChainService) nextBlock() (currentBlock *big.Int, err error) {
	currentBlock, err = bcs.blockNumber()
	if err != nil {
		return
	}
	targetBlockNumber := new(big.Int).Add(currentBlock, big.NewInt(1))
	for currentBlock.Cmp(targetBlockNumber) == -1 {
		time.Sleep(500 * time.Millisecond)
		currentBlock, err = bcs.blockNumber()
		if err != nil {
			return
		}
	}
	return
}

// TokenNetworkAddres return a proxy to interact with a token.
func (bcs *BlockChainService) Token(tokenAddress common.Address) (t *TokenProxy) {
	_, ok := bcs.addressTokens[tokenAddress]
	if !ok {
		token, err := contracts.NewToken(tokenAddress, bcs.Client)
		if err != nil {
			log.Error(fmt.Sprintf("NewToken %s err %s", tokenAddress.String(), err))
		}
		bcs.addressTokens[tokenAddress] = &TokenProxy{
			Address: tokenAddress, bcs: bcs, Token: token}
	}
	return bcs.addressTokens[tokenAddress]
}

//TokenNetwork return a proxy to interact with a NettingChannelContract.
func (bcs *BlockChainService) TokenNetwork(address common.Address) (t *TokenNetworkProxy, err error) {
	_, ok := bcs.addressChannels[address]
	if !ok {
		ch, err := contracts.NewTokenNetwork(address, bcs.Client)
		if err != nil {
			log.Error(fmt.Sprintf("NewNettingChannelContract %s err %s", address.String(), err))
		}
		if !bcs.contractExist(address) {
			return nil, fmt.Errorf("no code at %s", address)
		}
		bcs.addressChannels[address] = &TokenNetworkProxy{Address: address, bcs: bcs, ch: ch}
	}
	return bcs.addressChannels[address], nil
}

//TokenNetworkWithoutCheck return a proxy to interact with a NettingChannelContract,don't check this address is valid or not
func (bcs *BlockChainService) TokenNetworkWithoutCheck(address common.Address) (t *TokenNetworkProxy, err error) {
	_, ok := bcs.addressChannels[address]
	if !ok {
		var ch *contracts.TokenNetwork
		ch, err = contracts.NewTokenNetwork(address, bcs.Client)
		if err != nil {
			log.Error(fmt.Sprintf("NewNettingChannelContract %s err %s", address.String(), err))
		}
		bcs.addressChannels[address] = &TokenNetworkProxy{Address: address, bcs: bcs, ch: ch}
	}
	return bcs.addressChannels[address], nil
}

// Registry Return a proxy to interact with Registry.
func (bcs *BlockChainService) Registry(address common.Address) (t *RegistryProxy) {
	_, ok := bcs.addressRegistries[address]
	if !ok {
		reg, err := contracts.NewTokenNetworkRegistry(address, bcs.Client)
		if err != nil {
			log.Error(fmt.Sprintf("NewRegistry %s err %s ", address.String(), err))
		}
		bcs.addressRegistries[address] = &RegistryProxy{address, bcs, reg}
	}
	return bcs.addressRegistries[address]
}

//RegistryProxy proxy for registry contract
type RegistryProxy struct {
	Address  common.Address //contract address
	bcs      *BlockChainService
	registry *contracts.TokenNetworkRegistry
}

// TokenNetworkByToken Get the ChannelManager address for a specific token
// @param token_address The address of the given token
// @return Address of tokenNetwork
func (r *RegistryProxy) TokenNetworkByToken(tokenAddress common.Address) (tokenNetworkAddress common.Address, err error) {
	tokenNetworkAddress, err = r.registry.Token_to_token_networks(r.bcs.getQueryOpts(), tokenAddress)
	if err != nil && err.Error() == "abi: unmarshalling empty output" {
		err = rerr.ErrNoTokenManager
	}
	return
}

//GetContract return contract itself
func (r *RegistryProxy) GetContract() *contracts.TokenNetworkRegistry {
	return r.registry
}

//AddToken register a new token,this token must be a valid erc20
func (r *RegistryProxy) AddToken(tokenAddress common.Address) (tokenNetworkAddress common.Address, err error) {
	tx, err := r.registry.CreateERC20TokenNetwork(r.bcs.Auth, tokenAddress)
	if err != nil {
		return
	}
	receipt, err := bind.WaitMined(GetCallContext(), r.bcs.Client, tx)
	if err != nil {
		return
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Info(fmt.Sprintf("AddToken failed %s,receipt=%s", utils.APex(r.Address), receipt))
		err = errors.New("AddToken tx execution failed")
		return
	}
	log.Info(fmt.Sprintf("AddToken success %s,token=%s", utils.APex(r.Address), tokenAddress.String()))
	return r.registry.Token_to_token_networks(r.bcs.getQueryOpts(), tokenAddress)
}

//TokenNetworkProxy proxy of TokenNetwork Contract
type TokenNetworkProxy struct {
	Address common.Address //this contract address
	bcs     *BlockChainService
	ch      *contracts.TokenNetwork
}

//NewChannel create new channel ,block until a new channel create
func (t *TokenNetworkProxy) NewChannel(partnerAddress common.Address, settleTimeout int) (err error) {
	tx, err := t.ch.OpenChannel(t.bcs.Auth, t.bcs.NodeAddress, partnerAddress, uint64(settleTimeout))
	if err != nil {
		return
	}
	log.Info(fmt.Sprintf("NewChannel txhash=%s", tx.Hash().String()))
	receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	if err != nil {
		return
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Info(fmt.Sprintf("NewChannel failed %s,receipt=%s", utils.APex(t.Address), receipt))
		err = errors.New("NewChannel tx execution failed")
		return
	}
	log.Info(fmt.Sprintf("NewChannel success %s, partnerAddress=%s", utils.APex(t.Address), utils.APex(partnerAddress)))
	return
}

/*GetChannelInfo Returns the channel specific data.
@param participant1 Address of one of the channel participants.
@param participant2 Address of the other channel participant.
@return Channel state and settle_block_number.
if state is 1, settleBlockNumber is settle timeout, if state is 2,settleBlockNumber is the min block number ,settle can be called.
*/
func (t *TokenNetworkProxy) GetChannelInfo(participant1, participant2 common.Address) (channelID common.Hash, settleBlockNumber, openBlockNumber uint64, state uint8, settleTimeout uint64, err error) {
	return t.ch.GetChannelInfo(t.bcs.getQueryOpts(), participant1, participant2)
}

//GetChannelParticipantInfo Returns Info of this channel.
//@return The address of the token.
func (t *TokenNetworkProxy) GetChannelParticipantInfo(participant, partner common.Address) (deposit *big.Int, balanceHash common.Hash, nonce uint64, err error) {
	deposit, h, nonce, err := t.ch.GetChannelParticipantInfo(t.bcs.getQueryOpts(), participant, partner)
	balanceHash = common.BytesToHash(h[:])
	return
}

//GetContract return contract
func (t *TokenNetworkProxy) GetContract() *contracts.TokenNetwork {
	return t.ch
}

//TokenProxy proxy of ERC20 token
type TokenProxy struct {
	Address common.Address
	bcs     *BlockChainService
	Token   *contracts.Token
	lock    sync.Mutex
}

// TotalSupply total amount of tokens
func (t *TokenProxy) TotalSupply() (*big.Int, error) {
	return t.Token.TotalSupply(t.bcs.getQueryOpts())
}

// BalanceOf The balance
// @param _owner The address from which the balance will be retrieved
func (t *TokenProxy) BalanceOf(addr common.Address) (*big.Int, error) {
	amount, err := t.Token.BalanceOf(t.bcs.getQueryOpts(), addr)
	if err != nil {
		return nil, err
	}
	return amount, err
}

// Allowance Amount of remaining tokens allowed to spent
// @param _owner The address of the account owning tokens
// @param _spender The address of the account able to transfer the tokens
func (t *TokenProxy) Allowance(owner, spender common.Address) (int64, error) {
	amount, err := t.Token.Allowance(t.bcs.getQueryOpts(), owner, spender)
	if err != nil {
		return 0, err
	}
	return amount.Int64(), err //todo if amount larger than max int64?
}

// Approve Whether the approval was successful or not
// @notice `msg.sender` approves `_spender` to spend `_value` tokens
// @param _spender The address of the account able to transfer the tokens
// @param _value The amount of wei to be approved for transfer
func (t *TokenProxy) Approve(spender common.Address, value *big.Int) (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	tx, err := t.Token.Approve(t.bcs.Auth, spender, value)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("Approve %s, txhash=%s", utils.APex(spender), tx.Hash().String()))
	receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Info(fmt.Sprintf("Approve failed %s,receipt=%s", utils.APex(t.Address), receipt))
		return errors.New("Approve tx execution failed")
	}
	log.Info(fmt.Sprintf("Approve success %s,spender=%s,value=%d", utils.APex(t.Address), utils.APex(spender), value))
	return nil
}

// Transfer return Whether the transfer was successful or not
// @notice send `_value` token to `_to` from `msg.sender`
// @param _to The address of the recipient
// @param _value The amount of token to be transferred
func (t *TokenProxy) Transfer(spender common.Address, value *big.Int) (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	tx, err := t.Token.Transfer(t.bcs.Auth, spender, value)
	if err != nil {
		return err
	}
	receipt, err := bind.WaitMined(GetCallContext(), t.bcs.Client, tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Info(fmt.Sprintf("Transfer failed %s,receipt=%s", utils.APex(t.Address), receipt))
		return errors.New("Transfer tx execution failed")
	}
	log.Info(fmt.Sprintf("Transfer success %s,spender=%s,value=%d", utils.APex(t.Address), utils.APex(spender), value))
	return nil
}
