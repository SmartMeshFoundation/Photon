package helper

import (
	"context"
	"math/big"
	"sync"

	"fmt"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/network/xmpptransport"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatedier/frp/src/utils/log"
	"github.com/go-errors/errors"
)

var errNotConnectd = errors.New("eth not connected")

//SafeEthClient how to recover from a restart of geth
type SafeEthClient struct {
	*ethclient.Client
	lock       sync.Mutex
	url        string
	ReConnect  map[string]chan struct{}
	Status     xmpptransport.Status
	StatusChan chan xmpptransport.Status
	quitChan   chan struct{}
}

//NewSafeClient create safeclient
func NewSafeClient(rawurl string) (*SafeEthClient, error) {
	c := &SafeEthClient{
		ReConnect:  make(map[string]chan struct{}),
		url:        rawurl,
		StatusChan: make(chan xmpptransport.Status, 10),
		quitChan:   make(chan struct{}),
	}
	var err error
	c.Client, err = ethclient.Dial(rawurl)
	if err == nil {
		c.changeStatus(xmpptransport.Connected)
	} else {
		//c.changeStatus(xmpptransport.Disconnected)
		go c.RecoverDisconnect()
	}
	return c, nil
}

//Close connection when destroy raiden service
func (c *SafeEthClient) Close() {
	c.Client.Close()
	close(c.quitChan)
}

//IsConnected return true when connected to eth rpc server
func (c *SafeEthClient) IsConnected() bool {
	return c.Status == xmpptransport.Connected
}

//RegisterReConnectNotify register notify when reconnect
func (c *SafeEthClient) RegisterReConnectNotify(name string) <-chan struct{} {
	c.lock.Lock()
	defer c.lock.Unlock()
	ch, ok := c.ReConnect[name]
	if ok {
		log.Warn("NeedReConnectNotify should only call once")
		return ch
	}
	ch = make(chan struct{}, 1)
	c.ReConnect[name] = ch
	return ch
}
func (c *SafeEthClient) changeStatus(newStatus xmpptransport.Status) {
	log.Info("ethclient connection status changed from %d to %d", c.Status, newStatus)
	c.Status = newStatus
	select {
	case c.StatusChan <- c.Status:
	default:
		//never block
	}
}

//RecoverDisconnect try to reconnect with geth after a restart of geth
func (c *SafeEthClient) RecoverDisconnect() {
	var err error
	var client *ethclient.Client
	c.changeStatus(xmpptransport.Reconnecting)
	for {
		log.Info("tyring to reconnect geth ...")
		select {
		case <-c.quitChan:
			return
		default:
			//never block
		}
		client, err = ethclient.Dial(c.url)
		if err != nil {
			log.Info(fmt.Sprintf("reconnect to geth error: %s", err))
			time.Sleep(time.Second * 3)
		} else {
			//reconnect ok
			c.Client = client
			c.changeStatus(xmpptransport.Connected)
			c.lock.Lock()
			var keys []string
			for name, c := range c.ReConnect {
				keys = append(keys, name)
				c <- struct{}{}
				close(c)
			}
			for _, name := range keys {
				delete(c.ReConnect, name)
			}
			c.lock.Unlock()
			return
		}
	}
}

//BlockByHash wrapper of BlockByHash
func (c *SafeEthClient) BlockByHash(ctx context.Context, hash common.Hash) (r1 *types.Block, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	r1, err = c.Client.BlockByHash(ctx, hash)
	return
}

//BlockByNumber wrapper of BlockByNumber
func (c *SafeEthClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, errNotConnectd
	}
	return c.Client.BlockByNumber(ctx, number)
}

// HeaderByHash returns the block header with the given hash.
func (c *SafeEthClient) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, errNotConnectd
	}
	return c.Client.HeaderByHash(ctx, hash)
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (c *SafeEthClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, errNotConnectd
	}
	return c.Client.HeaderByNumber(ctx, number)
}

//TransactionByHash wrapper of TransactionByHash
func (c *SafeEthClient) TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, false, errNotConnectd
	}
	return c.Client.TransactionByHash(ctx, hash)
}

//TransactionSender wrapper of TransactionSender
func (c *SafeEthClient) TransactionSender(ctx context.Context, tx *types.Transaction, block common.Hash, index uint) (common.Address, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return common.Address{}, errNotConnectd
	}
	return c.Client.TransactionSender(ctx, tx, block, index)
}

// TransactionCount returns the total number of transactions in the given block.
func (c *SafeEthClient) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return 0, errNotConnectd
	}
	return c.Client.TransactionCount(ctx, blockHash)
}

//TransactionInBlock wrapper of TransactionInBlock
func (c *SafeEthClient) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, errNotConnectd
	}
	return c.Client.TransactionInBlock(ctx, blockHash, index)
}

//TransactionReceipt wrappper of TransactionReceipt
func (c *SafeEthClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, errNotConnectd
	}
	return c.Client.TransactionReceipt(ctx, txHash)
}

//SyncProgress wrapper of SyncProgress
func (c *SafeEthClient) SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, errNotConnectd
	}
	return c.Client.SyncProgress(ctx)
}

//SubscribeNewHead wrapper of SubscribeNewHead
func (c *SafeEthClient) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, errNotConnectd
	}
	return c.Client.SubscribeNewHead(ctx, ch)
}

//NetworkID wrapper of NetworkID
func (c *SafeEthClient) NetworkID(ctx context.Context) (*big.Int, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, errNotConnectd
	}
	return c.Client.NetworkID(ctx)
}

//BalanceAt wrapper of BalanceAt
func (c *SafeEthClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, errNotConnectd
	}
	return c.Client.BalanceAt(ctx, account, blockNumber)
}

//StorageAt wrapper of StorageAt
func (c *SafeEthClient) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, errNotConnectd
	}
	return c.Client.StorageAt(ctx, account, key, blockNumber)
}

//CodeAt wrapper of CodeAt
func (c *SafeEthClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, errNotConnectd
	}
	return c.Client.CodeAt(ctx, account, blockNumber)
}

//NonceAt wrapper of NonceAt
func (c *SafeEthClient) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return 0, errNotConnectd
	}
	return c.Client.NonceAt(ctx, account, blockNumber)
}

//FilterLogs wrapper of FilterLogs
func (c *SafeEthClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, errNotConnectd
	}
	return c.Client.FilterLogs(ctx, q)
}

//SubscribeFilterLogs wrapper of SubscribeFilterLogs
func (c *SafeEthClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, errNotConnectd
	}
	return c.Client.SubscribeFilterLogs(ctx, q, ch)
}

//PendingBalanceAt wrapper of PendingBalanceAt
func (c *SafeEthClient) PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, errNotConnectd
	}
	return c.Client.PendingBalanceAt(ctx, account)
}

//PendingStorageAt wrapper of PendingStorageAt
func (c *SafeEthClient) PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, errNotConnectd
	}
	return c.Client.PendingStorageAt(ctx, account, key)
}

//PendingCodeAt wrapper of PendingCodeAt
func (c *SafeEthClient) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, errNotConnectd
	}
	return c.Client.PendingCodeAt(ctx, account)
}

//PendingNonceAt wrapper of PendingNonceAt
func (c *SafeEthClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return 0, errNotConnectd
	}
	return c.Client.PendingNonceAt(ctx, account)
}

// PendingTransactionCount returns the total number of transactions in the pending state.
func (c *SafeEthClient) PendingTransactionCount(ctx context.Context) (uint, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return 0, errNotConnectd
	}
	return c.Client.PendingTransactionCount(ctx)
}

//CallContract wrapper of CallContract
func (c *SafeEthClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, errNotConnectd
	}
	return c.Client.CallContract(ctx, msg, blockNumber)
}

//PendingCallContract wrapper of PendingCallContract
func (c *SafeEthClient) PendingCallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, errNotConnectd
	}
	return c.Client.PendingCallContract(ctx, msg)
}

//SuggestGasPrice wrapper of SuggestGasPrice
func (c *SafeEthClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return nil, errNotConnectd
	}
	return c.Client.SuggestGasPrice(ctx)
}

//EstimateGas wrapper of EstimateGas
func (c *SafeEthClient) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return 0, errNotConnectd
	}
	return c.Client.EstimateGas(ctx, msg)
}

//SendTransaction wrapper of SendTransaction
func (c *SafeEthClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.Client == nil {
		return errNotConnectd
	}
	return c.Client.SendTransaction(ctx, tx)
}
