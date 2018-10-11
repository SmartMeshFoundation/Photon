package rpc

import (
	"context"

	"math/big"

	"time"

	"fmt"

	"crypto/ecdsa"

	"sync"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/models"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/SmartMeshFoundation/SmartRaiden/network/netshare"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
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
	RegistryAddress     common.Address
	SecretRegistryProxy *SecretRegistryProxy
	RegistryProxy       *RegistryProxy
	//Client if eth rpc client
	Client          *helper.SafeEthClient
	addressTokens   map[common.Address]*TokenProxy
	addressChannels map[common.Address]*TokenNetworkProxy
	//Auth needs by call on blockchain todo remove this
	Auth      *bind.TransactOpts
	queryOpts *bind.CallOpts
}

// 1. 尝试从数据库获取RegistryAddress
// 2. 如果db里有且与用户指定的值不同,返回err
// 3. 如果数据库没有且用户指定了值,使用指定值并保存到db
// 4. 如果数据没有且用户没指定,根据公链判断使用哪一个默认值并保存到db
func getRegistryAddress(userRegistryAddress common.Address, db *models.ModelDB, client *helper.SafeEthClient) (registryAddress common.Address, err error) {
	dbRegistryAddress := db.GetRegistryAddress()
	if dbRegistryAddress == utils.EmptyAddress {
		registryAddress = userRegistryAddress
	} else if userRegistryAddress != utils.EmptyAddress && dbRegistryAddress != userRegistryAddress {
		err = fmt.Errorf(fmt.Sprintf("db mismatch, db's registry=%s,now registry=%s",
			registryAddress.String(), userRegistryAddress.String()))
		return
	} else {
		registryAddress = dbRegistryAddress
	}
	if registryAddress == utils.EmptyAddress {
		if client.Status == netshare.Connected {
			var genesisBlockHash common.Hash
			genesisBlockHash, err = client.GenesisBlockHash(context.Background())
			if err != nil {
				log.Error(err.Error())
				return
			}
			registryAddress = params.GenesisBlockHashToDefaultRegistryAddress[genesisBlockHash]
		} else {
			err = fmt.Errorf(fmt.Sprintf("registry address is needed when no network"))
			return
		}
	}
	if registryAddress != dbRegistryAddress {
		db.SaveRegistryAddress(registryAddress)
	}
	log.Info(fmt.Sprintf("start with registry address %s", registryAddress.String()))
	return
}

//NewBlockChainService create BlockChainService
func NewBlockChainService(config *params.Config, db *models.ModelDB) (bcs *BlockChainService, err error) {
	client, err := helper.NewSafeClient(config.EthRPCEndPoint)
	if err != nil {
		err = fmt.Errorf("cannot connect to geth :%s err=%s", config.EthRPCEndPoint, err)
		err = nil
	}
	// 获取registryAddress
	registryAddress, err := getRegistryAddress(config.RegistryAddress, db, client)
	if err != nil {
		client.Close()
		return
	}
	bcs = &BlockChainService{
		PrivKey:         config.PrivateKey,
		NodeAddress:     crypto.PubkeyToAddress(config.PrivateKey.PublicKey),
		RegistryAddress: registryAddress,
		Client:          client,
		addressTokens:   make(map[common.Address]*TokenProxy),
		addressChannels: make(map[common.Address]*TokenNetworkProxy),
		Auth:            bind.NewKeyedTransactor(config.PrivateKey),
	}
	bcs.queryOpts = &bind.CallOpts{
		Pending: false,
		From:    bcs.NodeAddress,
		Context: GetQueryConext(),
	}
	// remove gas limit config and let it calculate automatically
	//bcs.Auth.GasLimit = uint64(params.GasLimit)
	bcs.Auth.GasPrice = big.NewInt(params.GasPrice)
	return bcs, nil
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
func (bcs *BlockChainService) Token(tokenAddress common.Address) (t *TokenProxy, err error) {
	_, ok := bcs.addressTokens[tokenAddress]
	if !ok {
		token, err := contracts.NewToken(tokenAddress, bcs.Client)
		if err != nil {
			log.Error(fmt.Sprintf("NewToken %s err %s", tokenAddress.String(), err))
			return nil, err
		}
		bcs.addressTokens[tokenAddress] = &TokenProxy{
			Address: tokenAddress, bcs: bcs, Token: token}
	}
	return bcs.addressTokens[tokenAddress], nil
}

//TokenNetwork return a proxy to interact with a NettingChannelContract.
func (bcs *BlockChainService) TokenNetwork(address common.Address) (t *TokenNetworkProxy, err error) {
	_, ok := bcs.addressChannels[address]
	if !ok {
		var tokenNetwork *contracts.TokenNetwork
		tokenNetwork, err = contracts.NewTokenNetwork(address, bcs.Client)
		if err != nil {
			log.Error(fmt.Sprintf("NewNettingChannelContract %s err %s", address.String(), err))
			return
		}
		if !bcs.contractExist(address) {
			return nil, fmt.Errorf("no code at %s", address)
		}
		bcs.addressChannels[address] = &TokenNetworkProxy{Address: address, bcs: bcs, ch: tokenNetwork}
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
			return
		}
		bcs.addressChannels[address] = &TokenNetworkProxy{Address: address, bcs: bcs, ch: ch}
	}
	return bcs.addressChannels[address], nil
}

// Registry Return a proxy to interact with Registry.
func (bcs *BlockChainService) Registry(address common.Address) (t *RegistryProxy) {
	if bcs.RegistryProxy != nil {
		return bcs.RegistryProxy
	}
	reg, err := contracts.NewTokenNetworkRegistry(address, bcs.Client)
	if err != nil {
		log.Error(fmt.Sprintf("NewRegistry %s err %s ", address.String(), err))
		return
	}
	r := &RegistryProxy{address, bcs, reg}
	secAddr, err := r.registry.SecretRegistryAddress(nil)
	if err != nil {
		log.Error(fmt.Sprintf("get Secret_registry_address %s", err))
		return
	}
	s, err := contracts.NewSecretRegistry(secAddr, bcs.Client)
	if err != nil {
		log.Error(fmt.Sprintf("NewSecretRegistry err %s", err))
		return
	}
	bcs.RegistryProxy = r
	bcs.SecretRegistryProxy = &SecretRegistryProxy{
		Address:          secAddr,
		bcs:              bcs,
		registry:         s,
		RegisteredSecret: make(map[common.Hash]*sync.Mutex),
	}
	return bcs.RegistryProxy
}
