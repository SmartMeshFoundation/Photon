package rpc

import (
	"context"

	"math/big"

	"time"

	"fmt"

	"crypto/ecdsa"

	"sync"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/network/helper"
	"github.com/SmartMeshFoundation/Photon/network/netshare"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum"
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
	NodeAddress         common.Address
	tokenNetworkAddress common.Address
	SecretRegistryProxy *SecretRegistryProxy
	//Client if eth rpc client
	Client          *helper.SafeEthClient
	addressTokens   map[common.Address]*TokenProxy
	addressChannels map[common.Address]*TokenNetworkProxy
	RegistryProxy   *RegistryProxy
	//Auth needs by call on blockchain todo remove this
	Auth *bind.TransactOpts
}

//NewBlockChainService create BlockChainService
func NewBlockChainService(privateKey *ecdsa.PrivateKey, registryAddress common.Address, client *helper.SafeEthClient) (bcs *BlockChainService, err error) {
	bcs = &BlockChainService{
		PrivKey:             privateKey,
		NodeAddress:         crypto.PubkeyToAddress(privateKey.PublicKey),
		Client:              client,
		addressTokens:       make(map[common.Address]*TokenProxy),
		addressChannels:     make(map[common.Address]*TokenNetworkProxy),
		Auth:                bind.NewKeyedTransactor(privateKey),
		tokenNetworkAddress: registryAddress,
	}
	// remove gas limit config and let it calculate automatically
	//bcs.Auth.GasLimit = uint64(params.GasLimit)
	bcs.Auth.GasPrice = big.NewInt(params.DefaultGasPrice)

	bcs.Registry(registryAddress, client.Status == netshare.Connected)
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
func (bcs *BlockChainService) TokenNetwork(tokenAddress common.Address) (t *TokenNetworkProxy, err error) {
	_, ok := bcs.addressChannels[tokenAddress]
	if !ok {
		bcs.addressChannels[tokenAddress] = &TokenNetworkProxy{bcs.RegistryProxy, bcs, tokenAddress}
	}
	return bcs.addressChannels[tokenAddress], nil
}

//TokenNetworkWithoutCheck return a proxy to interact with a NettingChannelContract,don't check this address is valid or not
func (bcs *BlockChainService) TokenNetworkWithoutCheck(tokenAddress common.Address) (t *TokenNetworkProxy, err error) {
	_, ok := bcs.addressChannels[tokenAddress]
	if !ok {
		bcs.addressChannels[tokenAddress] = &TokenNetworkProxy{bcs.RegistryProxy, bcs, tokenAddress}
	}
	return bcs.addressChannels[tokenAddress], nil
}

// Registry Return a proxy to interact with Registry.
func (bcs *BlockChainService) Registry(address common.Address, hasConnectChain bool) (t *RegistryProxy) {
	if bcs.RegistryProxy != nil && bcs.RegistryProxy.ch != nil {
		return bcs.RegistryProxy
	}
	r := &RegistryProxy{
		Address: address,
	}
	if hasConnectChain {
		reg, err := contracts.NewTokensNetwork(address, bcs.Client)
		if err != nil {
			log.Error(fmt.Sprintf("NewRegistry %s err %s ", address.String(), err))
			return
		}
		r.ch = reg
		secAddr, err := r.ch.SecretRegistry(nil)
		if err != nil {
			log.Error(fmt.Sprintf("get Secret_registry_address %s", err))
			return
		}
		s, err := contracts.NewSecretRegistry(secAddr, bcs.Client)
		if err != nil {
			log.Error(fmt.Sprintf("NewSecretRegistry err %s", err))
			return
		}
		bcs.SecretRegistryProxy = &SecretRegistryProxy{
			Address:          secAddr,
			bcs:              bcs,
			registry:         s,
			RegisteredSecret: make(map[common.Hash]*sync.Mutex),
		}
	}
	bcs.RegistryProxy = r
	return bcs.RegistryProxy
}

// GetRegistryAddress :
func (bcs *BlockChainService) GetRegistryAddress() common.Address {
	if bcs.RegistryProxy != nil {
		return bcs.RegistryProxy.Address
	}
	return utils.EmptyAddress
}

// GetSecretRegistryAddress :
func (bcs *BlockChainService) GetSecretRegistryAddress() common.Address {
	if bcs.SecretRegistryProxy != nil {
		return bcs.SecretRegistryProxy.Address
	}
	return utils.EmptyAddress
}

// SyncProgress :
func (bcs *BlockChainService) SyncProgress() (sp *ethereum.SyncProgress, err error) {
	return bcs.Client.SyncProgress(context.Background())
}
