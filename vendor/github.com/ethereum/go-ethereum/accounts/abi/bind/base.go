// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package bind

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/SmartMeshFoundation/Photon/log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
)

// SignerFn is a signer function callback when a contract requires a method to
// sign the transaction before submission.
type SignerFn func(types.Signer, common.Address, *types.Transaction) (*types.Transaction, error)

// CallOpts is the collection of options to fine tune a contract call request.
type CallOpts struct {
	Pending bool           // Whether to operate on the pending state or the last known one
	From    common.Address // Optional the sender address, otherwise the first account is used

	Context context.Context // Network context to support cancellation and timeouts (nil = no timeout)
}

// TransactOpts is the collection of authorization data required to create a
// valid Ethereum transaction.
type TransactOpts struct {
	From   common.Address // Ethereum account to send the transaction from
	Nonce  *big.Int       // Nonce to use for the transaction execution (nil = use pending state)
	Signer SignerFn       // Method to use for signing the transaction (mandatory)

	Value    *big.Int // Funds to transfer along along the transaction (nil = 0 = no funds)
	GasPrice *big.Int // Gas price to use for the transaction execution (nil = gas price oracle)
	GasLimit uint64   // Gas limit to set for the transaction execution (0 = estimate)

	Context context.Context // Network context to support cancellation and timeouts (nil = no timeout)
}

// FilterOpts is the collection of options to fine tune filtering for events
// within a bound contract.
type FilterOpts struct {
	Start uint64  // Start of the queried range
	End   *uint64 // End of the range (nil = latest)

	Context context.Context // Network context to support cancellation and timeouts (nil = no timeout)
}

// WatchOpts is the collection of options to fine tune subscribing for events
// within a bound contract.
type WatchOpts struct {
	Start   *uint64         // Start of the queried range (nil = latest)
	Context context.Context // Network context to support cancellation and timeouts (nil = no timeout)
}

// BoundContract is the base wrapper object that reflects a contract on the
// Ethereum network. It contains a collection of methods that are used by the
// higher level contract bindings to operate.
type BoundContract struct {
	address    common.Address     // Deployment address of the contract on the Ethereum blockchain
	abi        abi.ABI            // Reflect based ABI to access the correct Ethereum methods
	caller     ContractCaller     // Read interface to interact with the blockchain
	transactor ContractTransactor // Write interface to interact with the blockchain
	filterer   ContractFilterer   // Event filtering to interact with the blockchain
}

// NewBoundContract creates a low level contract interface through which calls
// and transactions may be made through.
func NewBoundContract(address common.Address, abi abi.ABI, caller ContractCaller, transactor ContractTransactor, filterer ContractFilterer) *BoundContract {
	return &BoundContract{
		address:    address,
		abi:        abi,
		caller:     caller,
		transactor: transactor,
		filterer:   filterer,
	}
}

// DeployContract deploys a contract onto the Ethereum blockchain and binds the
// deployment address with a Go wrapper.
func DeployContract(opts *TransactOpts, abi abi.ABI, bytecode []byte, backend ContractBackend, params ...interface{}) (common.Address, *types.Transaction, *BoundContract, error) {
	// Otherwise try to deploy the contract
	c := NewBoundContract(common.Address{}, abi, backend, backend, backend)

	input, err := c.abi.Pack("", params...)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	tx, err := c.transact(opts, nil, append(bytecode, input...))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	c.address = crypto.CreateAddress(opts.From, tx.Nonce())
	return c.address, tx, c, nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (c *BoundContract) Call(opts *CallOpts, result interface{}, method string, params ...interface{}) error {
	// Don't crash on a lazy user
	if opts == nil {
		opts = new(CallOpts)
	}
	// Pack the input, call and unpack the results
	input, err := c.abi.Pack(method, params...)
	if err != nil {
		return err
	}
	var (
		msg    = ethereum.CallMsg{From: opts.From, To: &c.address, Data: input}
		ctx    = ensureContext(opts.Context)
		code   []byte
		output []byte
	)
	if opts.Pending {
		pb, ok := c.caller.(PendingContractCaller)
		if !ok {
			return ErrNoPendingState
		}
		output, err = pb.PendingCallContract(ctx, msg)
		if err == nil && len(output) == 0 {
			// Make sure we have a contract to operate on, and bail out otherwise.
			if code, err = pb.PendingCodeAt(ctx, c.address); err != nil {
				return err
			} else if len(code) == 0 {
				return ErrNoCode
			}
		}
	} else {
		output, err = c.caller.CallContract(ctx, msg, nil)
		if err == nil && len(output) == 0 {
			// Make sure we have a contract to operate on, and bail out otherwise.
			if code, err = c.caller.CodeAt(ctx, c.address, nil); err != nil {
				return err
			} else if len(code) == 0 {
				return ErrNoCode
			}
		}
	}
	if err != nil {
		return err
	}
	return c.abi.Unpack(result, method, output)
}

// Transact invokes the (paid) contract method with params as input values.
func (c *BoundContract) Transact(opts *TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	// Otherwise pack up the parameters and invoke the contract
	input, err := c.abi.Pack(method, params...)
	if err != nil {
		return nil, err
	}
	return c.transact(opts, &c.address, input)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (c *BoundContract) Transfer(opts *TransactOpts) (*types.Transaction, error) {
	return c.transact(opts, &c.address, nil)
}

var lock = sync.Mutex{}

//防止出现并发transaction调用,nonce重用现象. 假定前提是此账户在运行过程中只会在这一个地方使用.
var nonceMap = make(map[common.Address]uint64)

//ReInitNonceMap todo 以后如果有更好的办法,应该弃用. 只提供给casemanager使用,他和photon会共用一套账号,分别发起Tx
func ReInitNonceMap() {
	lock.Lock()
	nonceMap = make(map[common.Address]uint64)
	lock.Unlock()
}

// GetValidNonceAndLock 获取可用Nonce,本地内存维护,必须与ReturnNonce配对使用,否则会出现死锁
func GetValidNonceAndLock(conn ContractTransactor, address common.Address, ctx context.Context) (nonce uint64, err error){
	lock.Lock()
	if _, exist := nonceMap[address]; !exist {
		nonce, err =conn.PendingNonceAt(ensureContext(ctx), address)
		if err != nil {
			lock.Unlock()
			err = fmt.Errorf("GetValidNonce failed : %v", err)
			return
		}
		nonceMap[address] = nonce - 1
	}
	nonce = nonceMap[address]
	nonce++
	nonceMap[address] = nonce
	log.Trace(fmt.Sprintf("GetValidNonce success : address=%s nonce=%d ", address.String(), nonce))
	return
}
// ConfirmNonceAndUnlock :
func ConfirmNonceAndUnlock(address common.Address, nonce uint64, used bool) {
	if used {
	} else {
		nonceMap[address] = nonce -1
	}
	lock.Unlock()
}

// transact executes an actual transaction invocation, first deriving any missing
// authorization fields, and then scheduling the transaction for execution.
func (c *BoundContract) transact(opts *TransactOpts, contract *common.Address, input []byte) (*types.Transaction, error) {
	var err error
	lock.Lock()
	defer lock.Unlock()
	// Ensure a valid value field and resolve the account nonce
	value := opts.Value
	if value == nil {
		value = new(big.Int)
	}
	var nonce uint64
	if opts.Nonce == nil {
		if _, exist := nonceMap[opts.From]; !exist {
			nonce, err = c.transactor.PendingNonceAt(ensureContext(opts.Context), opts.From)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
			}
			log.Trace(fmt.Sprintf("first tx after  start up for %s,nonce=%d,this account should be only used by this application only ", opts.From.String(), nonce))
			nonceMap[opts.From] = nonce - 1
		}
		nonce = nonceMap[opts.From]
		nonce++
		defer func() {
			if err != nil {
				//发生了错误
				log.Info(fmt.Sprintf("%s nonce=%d Tx fail %s", opts.From.String(), nonce, err))
				nonce--
			}
			nonceMap[opts.From] = nonce
		}()
	} else {
		nonce = opts.Nonce.Uint64()
	}
	log.Trace(fmt.Sprintf("start tx addr=%s,nonce=%d", opts.From.String(), nonce))
	// Figure out the gas allowance and gas price values
	gasPrice := opts.GasPrice
	if gasPrice == nil {
		gasPrice, err = c.transactor.SuggestGasPrice(ensureContext(opts.Context))
		if err != nil {
			return nil, fmt.Errorf("failed to suggest gas price: %v", err)
		}
	}
	gasLimit := opts.GasLimit
	if gasLimit == 0 {
		// Gas estimation cannot succeed without code for method invocations
		if contract != nil {
			var code []byte
			if code, err = c.transactor.PendingCodeAt(ensureContext(opts.Context), c.address); err != nil {
				return nil, err
			} else if len(code) == 0 {
				return nil, ErrNoCode
			}
		}
		// If the contract surely has code (or code is not needed), estimate the transaction
		msg := ethereum.CallMsg{From: opts.From, To: contract, Value: value, Data: input}
		gasLimit, err = c.transactor.EstimateGas(ensureContext(opts.Context), msg)
		if err != nil {
			return nil, fmt.Errorf("failed to estimate gas needed: %v", err)
		}
	}
	//fmt.Printf("send Transaction with nonce=%d,fromaddress=%s\n", nonce, opts.From.String())
	// Create the transaction, sign it and schedule it for execution
	var rawTx *types.Transaction
	if contract == nil {
		rawTx = types.NewContractCreation(nonce, value, gasLimit, gasPrice, input)
	} else {
		rawTx = types.NewTransaction(nonce, c.address, value, gasLimit, gasPrice, input)
	}
	if opts.Signer == nil {
		return nil, errors.New("no signer to authorize the transaction with")
	}
	var networkID *big.Int
	networkID, err = c.transactor.NetworkID(ensureContext(opts.Context))
	if err != nil {
		return nil, err
	}

	// special dealing for specrum main net
	var head *types.Header
	head, err = c.transactor.HeaderByNumber(ensureContext(opts.Context), big.NewInt(1))
	if err != nil {
		return nil, err
	}
	if head != nil && head.Hash() == common.HexToHash("0x57e682b80257aad73c4f3ad98d20435b4e1644d8762ef1ea1ff2806c27a5fa3d") {
		fmt.Printf("change chainid to 20180430\n")
		networkID = big.NewInt(20180430)
	}
	var signedTx *types.Transaction
	signedTx, err = opts.Signer(types.NewEIP155Signer(networkID), opts.From, rawTx)
	if err != nil {
		return nil, err
	}
	if err = c.transactor.SendTransaction(ensureContext(opts.Context), signedTx); err != nil {
		return nil, err
	}
	//time.Sleep(100 * time.Millisecond)
	return signedTx, nil
}

// FilterLogs filters contract logs for past blocks, returning the necessary
// channels to construct a strongly typed bound iterator on top of them.
func (c *BoundContract) FilterLogs(opts *FilterOpts, name string, query ...[]interface{}) (chan types.Log, event.Subscription, error) {
	// Don't crash on a lazy user
	if opts == nil {
		opts = new(FilterOpts)
	}
	// Append the event selector to the query parameters and construct the topic set
	query = append([][]interface{}{{c.abi.Events[name].Id()}}, query...)

	topics, err := makeTopics(query...)
	if err != nil {
		return nil, nil, err
	}
	// Start the background filtering
	logs := make(chan types.Log, 128)

	config := ethereum.FilterQuery{
		Addresses: []common.Address{c.address},
		Topics:    topics,
		FromBlock: new(big.Int).SetUint64(opts.Start),
	}
	if opts.End != nil {
		config.ToBlock = new(big.Int).SetUint64(*opts.End)
	}
	/* TODO(karalabe): Replace the rest of the method below with this when supported
	sub, err := c.filterer.SubscribeFilterLogs(ensureContext(opts.Context), config, logs)
	*/
	buff, err := c.filterer.FilterLogs(ensureContext(opts.Context), config)
	if err != nil {
		return nil, nil, err
	}
	sub, err := event.NewSubscription(func(quit <-chan struct{}) error {
		for _, log := range buff {
			select {
			case logs <- log:
			case <-quit:
				return nil
			}
		}
		return nil
	}), nil

	if err != nil {
		return nil, nil, err
	}
	return logs, sub, nil
}

// WatchLogs filters subscribes to contract logs for future blocks, returning a
// subscription object that can be used to tear down the watcher.
func (c *BoundContract) WatchLogs(opts *WatchOpts, name string, query ...[]interface{}) (chan types.Log, event.Subscription, error) {
	// Don't crash on a lazy user
	if opts == nil {
		opts = new(WatchOpts)
	}
	// Append the event selector to the query parameters and construct the topic set
	query = append([][]interface{}{{c.abi.Events[name].Id()}}, query...)

	topics, err := makeTopics(query...)
	if err != nil {
		return nil, nil, err
	}
	// Start the background filtering
	logs := make(chan types.Log, 128)

	config := ethereum.FilterQuery{
		Addresses: []common.Address{c.address},
		Topics:    topics,
	}
	if opts.Start != nil {
		config.FromBlock = new(big.Int).SetUint64(*opts.Start)
	}
	sub, err := c.filterer.SubscribeFilterLogs(ensureContext(opts.Context), config, logs)
	if err != nil {
		return nil, nil, err
	}
	return logs, sub, nil
}

// UnpackLog unpacks a retrieved log into the provided output structure.
func (c *BoundContract) UnpackLog(out interface{}, event string, log types.Log) error {
	if len(log.Data) > 0 {
		if err := c.abi.Unpack(out, event, log.Data); err != nil {
			return err
		}
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	return parseTopics(out, indexed, log.Topics[1:])
}

// ensureContext is a helper method to ensure a context is not nil, even if the
// user specified it as such.
func ensureContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.TODO()
	}
	return ctx
}
