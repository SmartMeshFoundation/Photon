package helper

import (
	"context"
	"math/big"
	"sync"

	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatedier/frp/src/utils/log"
)

//how to recover from a restart of geth
type SafeEthClient struct {
	*ethclient.Client
	lock      sync.Mutex
	url       string
	ReConnect map[string]chan struct{}
}

func NewSafeClient(rawurl string) (*SafeEthClient, error) {
	c := new(SafeEthClient)
	c.ReConnect = make(map[string]chan struct{})
	c.url = rawurl
	var err error
	c.Client, err = ethclient.Dial(rawurl)
	return c, err
}
func (this *SafeEthClient) RegisterReConnectNotify(name string) <-chan struct{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	c, ok := this.ReConnect[name]
	if ok {
		log.Warn("NeedReConnectNotify should only call once")
		return c
	}
	c = make(chan struct{}, 1)
	this.ReConnect[name] = c
	return c
}

//try to reconnect with geth after a restart of geth
func (this *SafeEthClient) RecoverDisconnect() {
	var err error
	var client *ethclient.Client
	for {
		log.Info("tyring to reconnect geth ...")
		client, err = ethclient.Dial(this.url)
		if err != nil {
			log.Info("reconnect to geth error:", err)
		} else {
			//reconnect ok
			this.lock.Lock()
			this.Client = client
			var keys []string
			for name, c := range this.ReConnect {
				keys = append(keys, name)
				c <- struct{}{}
				close(c)
			}
			for _, name := range keys {
				delete(this.ReConnect, name)
			}
			this.lock.Unlock()
			return
		}
		time.Sleep(time.Second * 3)
	}
}
func (this *SafeEthClient) BlockByHash(ctx context.Context, hash common.Hash) (r1 *types.Block, err error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	r1, err = this.Client.BlockByHash(ctx, hash)
	return
}

func (this *SafeEthClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.BlockByNumber(ctx, number)
}

// HeaderByHash returns the block header with the given hash.
func (this *SafeEthClient) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.HeaderByHash(ctx, hash)
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (this *SafeEthClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.HeaderByNumber(ctx, number)
}

func (this *SafeEthClient) TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.TransactionByHash(ctx, hash)
}

func (this *SafeEthClient) TransactionSender(ctx context.Context, tx *types.Transaction, block common.Hash, index uint) (common.Address, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.TransactionSender(ctx, tx, block, index)
}

// TransactionCount returns the total number of transactions in the given block.
func (this *SafeEthClient) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.TransactionCount(ctx, blockHash)
}

func (this *SafeEthClient) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.TransactionInBlock(ctx, blockHash, index)
}

func (this *SafeEthClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.TransactionReceipt(ctx, txHash)
}

func (this *SafeEthClient) SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.SyncProgress(ctx)
}

func (this *SafeEthClient) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.SubscribeNewHead(ctx, ch)
}

func (this *SafeEthClient) NetworkID(ctx context.Context) (*big.Int, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.NetworkID(ctx)
}

func (this *SafeEthClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.BalanceAt(ctx, account, blockNumber)
}

func (this *SafeEthClient) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.StorageAt(ctx, account, key, blockNumber)
}
func (this *SafeEthClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.CodeAt(ctx, account, blockNumber)
}

func (this *SafeEthClient) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.NonceAt(ctx, account, blockNumber)
}

func (this *SafeEthClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.FilterLogs(ctx, q)
}

func (this *SafeEthClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.SubscribeFilterLogs(ctx, q, ch)
}

func (this *SafeEthClient) PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.PendingBalanceAt(ctx, account)
}

func (this *SafeEthClient) PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.PendingStorageAt(ctx, account, key)
}

func (this *SafeEthClient) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.PendingCodeAt(ctx, account)
}

func (this *SafeEthClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.PendingNonceAt(ctx, account)
}

// PendingTransactionCount returns the total number of transactions in the pending state.
func (this *SafeEthClient) PendingTransactionCount(ctx context.Context) (uint, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.PendingTransactionCount(ctx)
}

func (this *SafeEthClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.CallContract(ctx, msg, blockNumber)
}

func (this *SafeEthClient) PendingCallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.PendingCallContract(ctx, msg)
}

func (this *SafeEthClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.SuggestGasPrice(ctx)
}

func (this *SafeEthClient) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.EstimateGas(ctx, msg)
}

func (this *SafeEthClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.Client.SendTransaction(ctx, tx)
}
