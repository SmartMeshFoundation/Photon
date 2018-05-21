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
	Client                 *helper.SafeEthClient
	addressTokens          map[common.Address]*TokenProxy
	addressDiscoveries     map[common.Address]*EndpointRegistryProxy
	addressChannelManagers map[common.Address]*ChannelManagerContractProxy
	addressChannels        map[common.Address]*NettingChannelContractProxy
	addressRegistries      map[common.Address]*RegistryProxy
	//Auth needs by call on blockchain todo remove this
	Auth      *bind.TransactOpts
	queryOpts *bind.CallOpts
}

//NewBlockChainService create BlockChainService
func NewBlockChainService(privKey *ecdsa.PrivateKey, registryAddress common.Address, client *helper.SafeEthClient) *BlockChainService {
	bcs := &BlockChainService{
		PrivKey:                privKey,
		NodeAddress:            crypto.PubkeyToAddress(privKey.PublicKey),
		RegistryAddress:        registryAddress,
		Client:                 client,
		addressTokens:          make(map[common.Address]*TokenProxy),
		addressDiscoveries:     make(map[common.Address]*EndpointRegistryProxy),
		addressChannelManagers: make(map[common.Address]*ChannelManagerContractProxy),
		addressChannels:        make(map[common.Address]*NettingChannelContractProxy),
		addressRegistries:      make(map[common.Address]*RegistryProxy),
		Auth:                   bind.NewKeyedTransactor(privKey),
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

// Token return a proxy to interact with a token.
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

// Discovery return a proxy to interact with the discovery.
func (bcs *BlockChainService) Discovery(discoverAddress common.Address) (t *EndpointRegistryProxy) {
	_, ok := bcs.addressDiscoveries[discoverAddress]
	if !ok {
		discover, err := contracts.NewEndpointRegistry(discoverAddress, bcs.Client)
		if err != nil {
			log.Error(fmt.Sprintf("NewEndpointRegistry on %s err %s", discoverAddress.String(), err))
		}
		bcs.addressDiscoveries[discoverAddress] = &EndpointRegistryProxy{Address: discoverAddress, Disovery: discover, bcs: bcs}
	}
	return bcs.addressDiscoveries[discoverAddress]
}

//NettingChannel return a proxy to interact with a NettingChannelContract.
func (bcs *BlockChainService) NettingChannel(address common.Address) (t *NettingChannelContractProxy, err error) {
	_, ok := bcs.addressChannels[address]
	if !ok {
		ch, err := contracts.NewNettingChannelContract(address, bcs.Client)
		if err != nil {
			log.Error(fmt.Sprintf("NewNettingChannelContract %s err %s", address.String(), err))
		}
		if !bcs.contractExist(address) {
			return nil, fmt.Errorf("no code at %s", address)
		}
		bcs.addressChannels[address] = &NettingChannelContractProxy{Address: address, bcs: bcs, ch: ch}
	}
	return bcs.addressChannels[address], nil
}

// Manager return a proxy to interact with a ChannelManagerContract.
func (bcs *BlockChainService) Manager(address common.Address) (t *ChannelManagerContractProxy) {
	_, ok := bcs.addressChannelManagers[address]
	if !ok {
		mgr, err := contracts.NewChannelManagerContract(address, bcs.Client)
		if err != nil {
			log.Error(fmt.Sprintf("NewChannelManagerContract %s err %s", address.String(), err))
		}
		bcs.addressChannelManagers[address] = &ChannelManagerContractProxy{Address: address, bcs: bcs, mgr: mgr}
	}
	return bcs.addressChannelManagers[address]
}

// Registry Return a proxy to interact with Registry.
func (bcs *BlockChainService) Registry(address common.Address) (t *RegistryProxy) {
	_, ok := bcs.addressRegistries[address]
	if !ok {
		reg, err := contracts.NewRegistry(address, bcs.Client)
		if err != nil {
			log.Error(fmt.Sprintf("NewRegistry %s err %s ", address.String(), err))
		}
		bcs.addressRegistries[address] = &RegistryProxy{address, bcs, reg}
	}
	return bcs.addressRegistries[address]
}

/*
GetAllChannelManagers return all channel managers
*/
func (bcs *BlockChainService) GetAllChannelManagers() (mgrs []*ChannelManagerContractProxy, err error) {
	reg := bcs.Registry(bcs.RegistryAddress)
	var mgrAddressess []common.Address
	mgrAddressess, err = reg.ChannelManagerAddresses()
	if err != nil {
		return
	}
	for _, mgrAddr := range mgrAddressess {
		mgr := bcs.Manager(mgrAddr)
		mgrs = append(mgrs, mgr)
	}
	return
}

//RegistryProxy proxy for registry contract
type RegistryProxy struct {
	Address  common.Address //contract address
	bcs      *BlockChainService
	registry *contracts.Registry
}

// ChannelManagerByToken Get the ChannelManager address for a specific token
// @param token_address The address of the given token
// @return Address of channel manager
func (r *RegistryProxy) ChannelManagerByToken(tokenAddress common.Address) (mgr common.Address, err error) {
	mgr, err = r.registry.ChannelManagerByToken(r.bcs.getQueryOpts(), tokenAddress)
	if err != nil && err.Error() == "abi: unmarshalling empty output" {
		err = rerr.ErrNoTokenManager
	}
	return
}

// TokenAddresses Get all registered tokens
// @return addresses of all registered tokens
func (r *RegistryProxy) TokenAddresses() (tokens []common.Address, err error) {
	return r.registry.TokenAddresses(r.bcs.getQueryOpts())
}

// ChannelManagerAddresses Get the addresses of all channel managers for all registered tokens
// @return addresses of all channel managers
func (r *RegistryProxy) ChannelManagerAddresses() (mgrs []common.Address, err error) {
	return r.registry.ChannelManagerAddresses(r.bcs.getQueryOpts())
}

//GetContract return contract itself
func (r *RegistryProxy) GetContract() *contracts.Registry {
	return r.registry
}

//AddToken register a new token,this token must be a valid erc20
func (r *RegistryProxy) AddToken(tokenAddress common.Address) (mgrAddr common.Address, err error) {
	tx, err := r.registry.AddToken(r.bcs.Auth, tokenAddress)
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
	return r.registry.ChannelManagerByToken(r.bcs.getQueryOpts(), tokenAddress)
}

//ChannelManagerContractProxy is proxy of ChannelMananger Contract
type ChannelManagerContractProxy struct {
	Address common.Address //contract address
	bcs     *BlockChainService
	mgr     *contracts.ChannelManagerContract
	lock    sync.Mutex //new channel protect
}

// GetChannelsAddresses Get all channels
// @return All the open channels //settled,closed will give too.
func (cm *ChannelManagerContractProxy) GetChannelsAddresses() (channels []common.Address, err error) {
	return cm.mgr.GetChannelsAddresses(cm.bcs.getQueryOpts())
}

// @notice Get all valid channels
// @return All the channels that have not been settled
//func (this *ChannelManagerContractProxy) GetValidChannelsAddresses() (validChannels []common.Address, err error) {
//	channels, err := this.mgr.GetChannelsAddresses(this.bcs.getQueryOpts())
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

// TokenAddress Get the address of the channel token
// @return The token
func (cm *ChannelManagerContractProxy) TokenAddress() (token common.Address, err error) {
	return cm.mgr.TokenAddress(cm.bcs.getQueryOpts())
}

// GetChannelWith Get the address of channel with a partner
// @param partner The address of the partner
// @return The address of the channel
func (cm *ChannelManagerContractProxy) GetChannelWith(partenerAddress common.Address) (channel common.Address, err error) {
	return cm.mgr.GetChannelWith(cm.bcs.getQueryOpts(), partenerAddress)
}

// NettingContractsByAddress Get all channels that an address participates in.
// @param node_address The address of the node
// @return The channel's addresses that node_address participates in.
func (cm *ChannelManagerContractProxy) NettingContractsByAddress(nodeAddress common.Address) (channels []common.Address, err error) {
	return cm.mgr.NettingContractsByAddress(cm.bcs.getQueryOpts(), nodeAddress)
}

//NettingChannelByAddress get NettChannel by partner
func (cm *ChannelManagerContractProxy) NettingChannelByAddress(nodeAddress common.Address) (channels []*NettingChannelContractProxy, err error) {
	addrs, err := cm.mgr.NettingContractsByAddress(cm.bcs.getQueryOpts(), nodeAddress)
	if err != nil {
		return
	}
	for _, addr := range addrs {
		ch, err := cm.bcs.NettingChannel(addr)
		if err != nil {
			continue //channel settled will be no code? todo make a test
		}
		channels = append(channels, ch)
	}
	return
}

// GetChannelsParticipants Get all participants of all channels
// @return All participants in all channels
func (cm *ChannelManagerContractProxy) GetChannelsParticipants() (nodes []common.Address, err error) {
	return cm.mgr.GetChannelsParticipants(cm.bcs.getQueryOpts())
}

//GetContract returns contract
func (cm *ChannelManagerContractProxy) GetContract() *contracts.ChannelManagerContract {
	return cm.mgr
}

//NewChannel create new channel ,block until a new channel create
func (cm *ChannelManagerContractProxy) NewChannel(partnerAddress common.Address, settleTimeout int) (chAddr common.Address, err error) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	tx, err := cm.mgr.NewChannel(cm.bcs.Auth, partnerAddress, big.NewInt(int64(settleTimeout)))
	if err != nil {
		return
	}
	receipt, err := bind.WaitMined(GetCallContext(), cm.bcs.Client, tx)
	if err != nil {
		return
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Info(fmt.Sprintf("NewChannel failed %s,receipt=%s", utils.APex(cm.Address), receipt))
		err = errors.New("NewChannel tx execution failed")
		return
	}
	log.Info(fmt.Sprintf("NewChannel success %s, partnerAddress=%s", utils.APex(cm.Address), utils.APex(partnerAddress)))
	return cm.GetChannelWith(partnerAddress)
}

//NettingChannelContractProxy proxy of NettingChannel Contract
type NettingChannelContractProxy struct {
	Address common.Address //this contract address
	bcs     *BlockChainService
	ch      *contracts.NettingChannelContract
}

// AddressAndBalance Get the address and balance of both partners in a channel.
// @return The address and balance pairs.
func (c *NettingChannelContractProxy) AddressAndBalance() (addr1 common.Address, balance1 *big.Int, addr2 common.Address, balance2 *big.Int, err error) {
	result, err := c.ch.AddressAndBalance(c.bcs.getQueryOpts())
	if err != nil {
		return
	}
	return result.Participant1, result.Balance1, result.Participant2, result.Balance2, err
}

// SettleTimeout Returns the number of blocks until the settlement timeout.
// @return The number of blocks until the settlement timeout.
func (c *NettingChannelContractProxy) SettleTimeout() (settleTimeout int, err error) {
	result, err := c.ch.SettleTimeout(c.bcs.getQueryOpts())
	if err != nil {
		return
	}
	return int(result.Int64()), err
}

// TokenAddress Returns the address of the token.
// @return The address of the token.
func (c *NettingChannelContractProxy) TokenAddress() (token common.Address, err error) {
	return c.ch.TokenAddress(c.bcs.getQueryOpts())
}

// Opened Returns the block number for when the channel was opened.
// @return The block number for when the channel was opened.
func (c *NettingChannelContractProxy) Opened() (openBlockNumber int64, err error) {
	result, err := c.ch.Opened(c.bcs.getQueryOpts())
	if err != nil {
		return
	}
	return result.Int64(), err
}

// Closed Returns the block number for when the channel was closed.
// @return The block number for when the channel was closed.
func (c *NettingChannelContractProxy) Closed() (closeBlockNumber int64, err error) {
	result, err := c.ch.Closed(c.bcs.getQueryOpts())
	if err != nil {
		return
	}
	return result.Int64(), err
}

// ClosingAddress Returns the address of the closing participant.
// @return The address of the closing participant.
func (c *NettingChannelContractProxy) ClosingAddress() (closeAddr common.Address, err error) {
	return c.ch.ClosingAddress(c.bcs.getQueryOpts())
}

//GetContract return contract
func (c *NettingChannelContractProxy) GetContract() *contracts.NettingChannelContract {
	return c.ch
}

//EndpointRegistryProxy is proxy of contract endpoint discovery
type EndpointRegistryProxy struct {
	Address  common.Address
	bcs      *BlockChainService
	Disovery *contracts.EndpointRegistry
}

/*
FindEndpointByAddress Finds the socket if given an Ethereum Address
@dev Finds the socket if given an Ethereum Address
@param An eth_address which is a 20 byte Ethereum Address
@return A socket which the current Ethereum Address is using.
*/
func (er *EndpointRegistryProxy) FindEndpointByAddress(nodeAddress common.Address) (socket string, err error) {
	return er.Disovery.FindEndpointByAddress(er.bcs.getQueryOpts(), nodeAddress)
}

/*
FindAddressByEndpoint Finds Ethreum Address if given an existing socket address
@dev Finds Ethreum Address if given an existing socket address
@param string of socket in this format "127.0.0.1:40001"
@return An ethereum address
*/
func (er *EndpointRegistryProxy) FindAddressByEndpoint(socket string) (nodeAddress common.Address, err error) {
	return er.Disovery.FindAddressByEndpoint(er.bcs.getQueryOpts(), socket)
}

/*
RegisterEndpoint Registers the Ethereum Address to the Endpoint socket.
@dev Registers the Ethereum Address to the Endpoint socket.
@param string of socket in this format "127.0.0.1:40001"
*/
func (er *EndpointRegistryProxy) RegisterEndpoint(socket string) (err error) {
	tx, err := er.Disovery.RegisterEndpoint(er.bcs.Auth, socket)
	if err != nil {
		return err
	}
	receipt, err := bind.WaitMined(GetCallContext(), er.bcs.Client, tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Info(fmt.Sprintf("registerEndpoint failed %s,receipt=%s", utils.APex(er.Address), receipt))
		return errors.New("RegisterEndpoint tx execution failed")
	}
	log.Info(fmt.Sprintf("RegisterEndpoint success %s,socket=%s", utils.APex(er.Address), socket))
	return nil
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
