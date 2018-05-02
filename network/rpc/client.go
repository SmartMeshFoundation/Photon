package rpc

import (
	"crypto/ecdsa"

	"context"

	"math/big"

	"time"

	"errors"

	"fmt"

	"sync"

	"github.com/SmartMeshFoundation/SmartRaiden/abi/bind"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/rerr"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

//context for tx
func GetCallContext() context.Context {
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(params.Default_Tx_Timeout))
	return ctx
}

//context for query on chain
func GetQueryConext() context.Context {
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(params.DEFAULT_POLL_TIMEOUT))
	return ctx
}

type BlockChainService struct {
	PrivKey               *ecdsa.PrivateKey
	NodeAddress           common.Address
	RegistryAddress       common.Address
	Client                *helper.SafeEthClient
	PollTimeOut           time.Duration
	AddressToken          map[common.Address]*TokenProxy
	AddressDiscovery      map[common.Address]*EndpointRegistryProxy
	AddressChannelManager map[common.Address]*ChannelManagerContractProxy
	AddressChannel        map[common.Address]*NettingChannelContractProxy
	AddressRegistry       map[common.Address]*RegistryProxy
	Auth                  *bind.TransactOpts
	queryOpts             *bind.CallOpts
	Lock                  sync.RWMutex
}

func NewBlockChainService(privKey *ecdsa.PrivateKey, registryAddress common.Address, client *helper.SafeEthClient) *BlockChainService {
	bcs := &BlockChainService{
		PrivKey:               privKey,
		NodeAddress:           crypto.PubkeyToAddress(privKey.PublicKey),
		RegistryAddress:       registryAddress,
		Client:                client,
		PollTimeOut:           params.DEFAULT_POLL_TIMEOUT,
		AddressToken:          make(map[common.Address]*TokenProxy),
		AddressDiscovery:      make(map[common.Address]*EndpointRegistryProxy),
		AddressChannelManager: make(map[common.Address]*ChannelManagerContractProxy),
		AddressChannel:        make(map[common.Address]*NettingChannelContractProxy),
		AddressRegistry:       make(map[common.Address]*RegistryProxy),
		Auth:                  bind.NewKeyedTransactor(privKey),
	}
	bcs.queryOpts = &bind.CallOpts{
		Pending: false,
		From:    bcs.NodeAddress,
		Context: GetQueryConext(),
	}
	//It needs to be set up, otherwise, even the contract revert will not report wrong.
	bcs.Auth.GasLimit = uint64(params.GAS_LIMIT)
	bcs.Auth.GasPrice = big.NewInt(params.GAS_PRICE)
	return bcs
}
func (this *BlockChainService) QueryOpts() *bind.CallOpts {
	return &bind.CallOpts{
		Pending: false,
		From:    this.NodeAddress,
		Context: GetQueryConext(),
	}
}
func (this *BlockChainService) BlockNumber() (num *big.Int, err error) {
	h, err := this.Client.HeaderByNumber(context.Background(), nil)
	if err == nil {
		num = h.Number
	}
	return
}
func (this *BlockChainService) Nonce(account common.Address) (uint64, error) {
	return this.Client.PendingNonceAt(context.Background(), account)
}
func (this *BlockChainService) Balance(account common.Address) (*big.Int, error) {
	return this.Client.PendingBalanceAt(context.Background(), account)
}
func (this *BlockChainService) ContractExist(contractAddr common.Address) bool {
	code, err := this.Client.CodeAt(context.Background(), contractAddr, nil)
	//spew.Dump("code:", code)
	return err == nil && len(code) > 0
}
func (this *BlockChainService) GetBlockHeader(blockNumber *big.Int) (*types.Header, error) {
	return this.Client.HeaderByNumber(context.Background(), blockNumber)
}

func (this *BlockChainService) NextBlock() (currentBlock *big.Int, err error) {
	currentBlock, err = this.BlockNumber()
	if err != nil {
		return
	}
	targetBlockNumber := new(big.Int).Add(currentBlock, big.NewInt(1))
	for currentBlock.Cmp(targetBlockNumber) == -1 {
		time.Sleep(500 * time.Millisecond)
		currentBlock, err = this.BlockNumber()
		if err != nil {
			return
		}
	}
	return
}

// Return a proxy to interact with a token.
func (this *BlockChainService) Token(tokenAddress common.Address) (t *TokenProxy) {
	_, ok := this.AddressToken[tokenAddress]
	if !ok {
		token, _ := NewToken(tokenAddress, this.Client)
		this.AddressToken[tokenAddress] = &TokenProxy{
			Address: tokenAddress, bcs: this, Token: token}
	}
	return this.AddressToken[tokenAddress]
}

// Return a proxy to interact with the discovery.
func (this *BlockChainService) Discovery(discoverAddress common.Address) (t *EndpointRegistryProxy) {
	_, ok := this.AddressDiscovery[discoverAddress]
	if !ok {
		discover, _ := NewEndpointRegistry(discoverAddress, this.Client)
		this.AddressDiscovery[discoverAddress] = &EndpointRegistryProxy{Address: discoverAddress, Disovery: discover, bcs: this}
	}
	return this.AddressDiscovery[discoverAddress]
}

//  Return a proxy to interact with a NettingChannelContract.
func (this *BlockChainService) NettingChannel(address common.Address) (t *NettingChannelContractProxy, err error) {
	_, ok := this.AddressChannel[address]
	if !ok {
		ch, _ := NewNettingChannelContract(address, this.Client)
		if !this.ContractExist(address) {
			return nil, errors.New(fmt.Sprintf("no code at %s", address))
		}
		this.AddressChannel[address] = &NettingChannelContractProxy{Address: address, bcs: this, ch: ch}
	}
	return this.AddressChannel[address], nil
}

// Return a proxy to interact with a ChannelManagerContract.
func (this *BlockChainService) Manager(address common.Address) (t *ChannelManagerContractProxy) {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	_, ok := this.AddressChannelManager[address]
	if !ok {
		mgr, _ := NewChannelManagerContract(address, this.Client)
		this.AddressChannelManager[address] = &ChannelManagerContractProxy{Address: address, bcs: this, mgr: mgr}
	}
	return this.AddressChannelManager[address]
}

// Return a proxy to interact with Registry.
func (this *BlockChainService) Registry(address common.Address) (t *RegistryProxy) {
	_, ok := this.AddressRegistry[address]
	if !ok {
		reg, _ := NewRegistry(address, this.Client)
		this.AddressRegistry[address] = &RegistryProxy{address, this, reg}
	}
	return this.AddressRegistry[address]
}

/*
all channel managers
*/
func (this *BlockChainService) GetAllChannelManagers() (mgrs []*ChannelManagerContractProxy, err error) {
	reg := this.Registry(this.RegistryAddress)
	var mgrAddressess []common.Address
	mgrAddressess, err = reg.ChannelManagerAddresses()
	if err != nil {
		return
	}
	for _, mgrAddr := range mgrAddressess {
		mgr := this.Manager(mgrAddr)
		mgrs = append(mgrs, mgr)
	}
	return
}

type RegistryProxy struct {
	Address  common.Address //contract address
	bcs      *BlockChainService
	registry *Registry
}

/// @notice Get the ChannelManager address for a specific token
/// @param token_address The address of the given token
/// @return Address of channel manager
func (this *RegistryProxy) ChannelManagerByToken(tokenAddress common.Address) (mgr common.Address, err error) {
	mgr, err = this.registry.ChannelManagerByToken(this.bcs.QueryOpts(), tokenAddress)
	if err != nil && err.Error() == "abi: unmarshalling empty output" {
		err = rerr.NoTokenManager
	}
	return
}

/// @notice Get all registered tokens
/// @return addresses of all registered tokens
func (this *RegistryProxy) TokenAddresses() (tokens []common.Address, err error) {
	return this.registry.TokenAddresses(this.bcs.QueryOpts())
}

/// @notice Get the addresses of all channel managers for all registered tokens
/// @return addresses of all channel managers
func (this *RegistryProxy) ChannelManagerAddresses() (mgrs []common.Address, err error) {
	return this.registry.ChannelManagerAddresses(this.bcs.QueryOpts())
}
func (this *RegistryProxy) GetContract() *Registry {
	return this.registry
}
func (this *RegistryProxy) AddToken(tokenAddress common.Address) (mgrAddr common.Address, err error) {
	tx, err := this.registry.AddToken(this.bcs.Auth, tokenAddress)
	if err != nil {
		return
	}
	receipt, err := bind.WaitMined(GetCallContext(), this.bcs.Client, tx)
	if err != nil {
		return
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Info(fmt.Sprintf("AddToken failed %s,receipt=%s", utils.APex(this.Address), receipt))
		err = errors.New("AddToken tx execution failed")
		return
	} else {
		log.Info(fmt.Sprintf("AddToken success %s,token=%s", utils.APex(this.Address), tokenAddress.String()))
	}
	//receipt.Logs[0].Data
	//spew.Config.DisableMethods = true
	//spew.Dump("receipt:", receipt)
	//fmt.Printf("receipt=%s\n", receipt)
	//The return value of the contract can not be obtained directly
	return this.registry.ChannelManagerByToken(this.bcs.QueryOpts(), tokenAddress)
}

type ChannelManagerContractProxy struct {
	Address common.Address //contract address
	bcs     *BlockChainService
	mgr     *ChannelManagerContract
	lock    sync.Mutex //new channel protect
}

/// @notice Get all channels
/// @return All the open channels //settled,closed will give too.
func (this *ChannelManagerContractProxy) GetChannelsAddresses() (channels []common.Address, err error) {
	return this.mgr.GetChannelsAddresses(this.bcs.QueryOpts())
}

/// @notice Get all valid channels
/// @return All the channels that have not been settled
//func (this *ChannelManagerContractProxy) GetValidChannelsAddresses() (validChannels []common.Address, err error) {
//	channels, err := this.mgr.GetChannelsAddresses(this.bcs.QueryOpts())
//	if err != nil {
//		return
//	}
//	for _, ch := range channels {
//		_, err2 := this.bcs.NettingChannel(ch)
//		if err2 != nil {
//			continue
//		}
//		validChannels = append(validChannels, ch)
//	}
//	return
//}

/// @notice Get the address of the channel token
/// @return The token
func (this *ChannelManagerContractProxy) TokenAddress() (token common.Address, err error) {
	return this.mgr.TokenAddress(this.bcs.QueryOpts())
}

/// @notice Get the address of channel with a partner
/// @param partner The address of the partner
/// @return The address of the channel
func (this *ChannelManagerContractProxy) GetChannelWith(partenerAddress common.Address) (channel common.Address, err error) {
	return this.mgr.GetChannelWith(this.bcs.QueryOpts(), partenerAddress)
}

/// @notice Get all channels that an address participates in.
/// @param node_address The address of the node
/// @return The channel's addresses that node_address participates in.
func (this *ChannelManagerContractProxy) NettingContractsByAddress(node_address common.Address) (channels []common.Address, err error) {
	return this.mgr.NettingContractsByAddress(this.bcs.QueryOpts(), node_address)
}
func (this *ChannelManagerContractProxy) NettingChannelByAddress(node_address common.Address) (channels []*NettingChannelContractProxy, err error) {
	addrs, err := this.mgr.NettingContractsByAddress(this.bcs.QueryOpts(), node_address)
	if err != nil {
		return
	}
	for _, addr := range addrs {
		ch, err := this.bcs.NettingChannel(addr)
		if err != nil {
			continue //channel settled will be no code? todo make a test
		}
		channels = append(channels, ch)
	}
	return
}

/// @notice Get all participants of all channels
/// @return All participants in all channels
func (this *ChannelManagerContractProxy) GetChannelsParticipants() (nodes []common.Address, err error) {
	return this.mgr.GetChannelsParticipants(this.bcs.QueryOpts())
}
func (this *ChannelManagerContractProxy) GetContract() *ChannelManagerContract {
	return this.mgr
}

func (this *ChannelManagerContractProxy) NewChannel(partnerAddress common.Address, settleTimeout int) (chAddr common.Address, err error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	tx, err := this.mgr.NewChannel(this.bcs.Auth, partnerAddress, big.NewInt(int64(settleTimeout)))
	if err != nil {
		return
	}
	receipt, err := bind.WaitMined(GetCallContext(), this.bcs.Client, tx)
	if err != nil {
		return
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Info(fmt.Sprintf("NewChannel failed %s,receipt=%s", utils.APex(this.Address), receipt))
		err = errors.New("NewChannel tx execution failed")
		return
	} else {
		log.Info(fmt.Sprintf("NewChannel success %s, partnerAddress=%s", utils.APex(this.Address), utils.APex(partnerAddress)))
	}
	return this.GetChannelWith(partnerAddress)
}

type NettingChannelContractProxy struct {
	Address common.Address //this contract address
	bcs     *BlockChainService
	ch      *NettingChannelContract
}

/// @notice Get the address and balance of both partners in a channel.
/// @return The address and balance pairs.
func (this *NettingChannelContractProxy) AddressAndBalance() (addr1 common.Address, balance1 *big.Int, addr2 common.Address, balance2 *big.Int, err error) {
	result, err := this.ch.AddressAndBalance(this.bcs.QueryOpts())
	if err != nil {
		return
	}
	return result.Participant1, result.Balance1, result.Participant2, result.Balance2, err
}

/// @notice Returns the number of blocks until the settlement timeout.
/// @return The number of blocks until the settlement timeout.
func (this *NettingChannelContractProxy) SettleTimeout() (settleTimeout int, err error) {
	result, err := this.ch.SettleTimeout(this.bcs.QueryOpts())
	if err != nil {
		return
	}
	return int(result.Int64()), err
}

/// @notice Returns the address of the token.
/// @return The address of the token.
func (this *NettingChannelContractProxy) TokenAddress() (token common.Address, err error) {
	return this.ch.TokenAddress(this.bcs.QueryOpts())
}

/// @notice Returns the block number for when the channel was opened.
/// @return The block number for when the channel was opened.
func (this *NettingChannelContractProxy) Opened() (openBlockNumber int64, err error) {
	result, err := this.ch.Opened(this.bcs.QueryOpts())
	if err != nil {
		return
	}
	return result.Int64(), err
}

/// @notice Returns the block number for when the channel was closed.
/// @return The block number for when the channel was closed.
func (this *NettingChannelContractProxy) Closed() (closeBlockNumber int64, err error) {
	result, err := this.ch.Closed(this.bcs.QueryOpts())
	if err != nil {
		return
	}
	return result.Int64(), err
}

/// @notice Returns the address of the closing participant.
/// @return The address of the closing participant.
func (this *NettingChannelContractProxy) ClosingAddress() (closeAddr common.Address, err error) {
	return this.ch.ClosingAddress(this.bcs.QueryOpts())
}

func (this *NettingChannelContractProxy) GetContract() *NettingChannelContract {
	return this.ch
}

type EndpointRegistryProxy struct {
	Address  common.Address
	bcs      *BlockChainService
	Disovery *EndpointRegistry
}

/*
 * @notice Finds the socket if given an Ethereum Address
 * @dev Finds the socket if given an Ethereum Address
 * @param An eth_address which is a 20 byte Ethereum Address
 * @return A socket which the current Ethereum Address is using.
 */
func (this *EndpointRegistryProxy) FindEndpointByAddress(nodeAddress common.Address) (socket string, err error) {
	return this.Disovery.FindEndpointByAddress(this.bcs.QueryOpts(), nodeAddress)
}

/*
 * @notice Finds Ethreum Address if given an existing socket address
 * @dev Finds Ethreum Address if given an existing socket address
 * @param string of socket in this format "127.0.0.1:40001"
 * @return An ethereum address
 */
func (this *EndpointRegistryProxy) FindAddressByEndpoint(socket string) (nodeAddress common.Address, err error) {
	return this.Disovery.FindAddressByEndpoint(this.bcs.QueryOpts(), socket)
}

/*
 * @notice Registers the Ethereum Address to the Endpoint socket.
 * @dev Registers the Ethereum Address to the Endpoint socket.
 * @param string of socket in this format "127.0.0.1:40001"
 */
func (this *EndpointRegistryProxy) RegisterEndpoint(socket string) (err error) {
	tx, err := this.Disovery.RegisterEndpoint(this.bcs.Auth, socket)
	if err != nil {
		return err
	}
	receipt, err := bind.WaitMined(GetCallContext(), this.bcs.Client, tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Info(fmt.Sprintf("registerEndpoint failed %s,receipt=%s", utils.APex(this.Address), receipt))
		return errors.New("RegisterEndpoint tx execution failed")
	} else {
		log.Info(fmt.Sprint("RegisterEndpoint success %s,socket=%s", utils.APex(this.Address), socket))
	}
	return nil
}

type TokenProxy struct {
	Address common.Address
	bcs     *BlockChainService
	Token   *Token
	lock    sync.Mutex
}

/// @return total amount of tokens
func (this *TokenProxy) TotalSupply() (*big.Int, error) {
	return this.Token.TotalSupply(this.bcs.QueryOpts())
}

/// @param _owner The address from which the balance will be retrieved
/// @return The balance
func (this *TokenProxy) BalanceOf(addr common.Address) (*big.Int, error) {
	amount, err := this.Token.BalanceOf(this.bcs.QueryOpts(), addr)
	if err != nil {
		return nil, err
	}
	return amount, err
}

/// @param _owner The address of the account owning tokens
/// @param _spender The address of the account able to transfer the tokens
/// @return Amount of remaining tokens allowed to spent
func (this *TokenProxy) Allowance(owner, spender common.Address) (int64, error) {
	amount, err := this.Token.Allowance(this.bcs.QueryOpts(), owner, spender)
	if err != nil {
		return 0, err
	}
	return amount.Int64(), err //todo if amount larger than max int64?
}

/// @notice `msg.sender` approves `_spender` to spend `_value` tokens
/// @param _spender The address of the account able to transfer the tokens
/// @param _value The amount of wei to be approved for transfer
/// @return Whether the approval was successful or not
func (this *TokenProxy) Approve(spender common.Address, value *big.Int) (err error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	tx, err := this.Token.Approve(this.bcs.Auth, spender, value)
	if err != nil {
		return err
	}
	receipt, err := bind.WaitMined(GetCallContext(), this.bcs.Client, tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Info(fmt.Sprintf("Approve failed %s,receipt=%s", utils.APex(this.Address), receipt))
		return errors.New("Approve tx execution failed")
	} else {
		log.Info(fmt.Sprint("Approve success %s,spender=%s,value=%d", utils.APex(this.Address), utils.APex(spender), value))
	}
	return nil
}

/// @notice send `_value` token to `_to` from `msg.sender`
/// @param _to The address of the recipient
/// @param _value The amount of token to be transferred
/// @return Whether the transfer was successful or not
func (this *TokenProxy) Transfer(spender common.Address, value *big.Int) (err error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	tx, err := this.Token.Transfer(this.bcs.Auth, spender, value)
	if err != nil {
		return err
	}
	receipt, err := bind.WaitMined(GetCallContext(), this.bcs.Client, tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Info(fmt.Sprintf("Transfer failed %s,receipt=%s", utils.APex(this.Address), receipt))
		return errors.New("Transfer tx execution failed")
	} else {
		log.Info(fmt.Sprintf("Transfer success %s,spender=%s,value=%d", utils.APex(this.Address), utils.APex(spender), value))
	}
	return nil
}
