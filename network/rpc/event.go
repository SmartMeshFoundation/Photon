package rpc

import (
	"math/big"
	"strings"

	"context"

	"github.com/SmartMeshFoundation/raiden-network/network/helper"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func buildQuery(contractAddress common.Address, fromBlock rpc.BlockNumber,
	toBlock rpc.BlockNumber, eventName string, abistr string) (q *ethereum.FilterQuery, err error) {
	parsed, err := abi.JSON(strings.NewReader(abistr))
	if err != nil {
		return
	}
	q = &ethereum.FilterQuery{}
	q.FromBlock = big.NewInt(int64(fromBlock))
	if toBlock == rpc.LatestBlockNumber {
		q.ToBlock = nil
	} else {
		q.ToBlock = big.NewInt(int64(toBlock))
	}
	q.Topics = [][]common.Hash{
		{parsed.Events[eventName].Id()}, //event signature
	}
	//spew.Dump("signature,", parsed.Events[eventName].Id())
	if contractAddress != utils.EmptyAddress {
		q.Addresses = []common.Address{contractAddress}
	}
	return
}
func ensureContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.TODO()
	}
	return ctx
}

//events of history
//if contractAddress is empty,it will query all contract
func EventGetInternal(ctx context.Context, contractAddress common.Address, fromBlock rpc.BlockNumber,
	toBlock rpc.BlockNumber, eventName string, abistr string, client *helper.SafeEthClient) ([]types.Log, error) {
	q, err := buildQuery(contractAddress, fromBlock, toBlock, eventName, abistr)
	if err != nil {
		return nil, err
	}
	ctx = ensureContext(ctx)
	return client.FilterLogs(ctx, *q)
}

func EventGet(contractAddress common.Address, eventName string, abistr string, client *helper.SafeEthClient) ([]types.Log, error) {
	return EventGetInternal(context.Background(), contractAddress, rpc.EarliestBlockNumber, rpc.LatestBlockNumber,
		eventName, abistr, client)
}

//events of future
//if contractAddress is empty,it will subscribe all contract
func EventSubscribeInternal(ctx context.Context, contractAddress common.Address, fromBlock rpc.BlockNumber,
	toBlock rpc.BlockNumber, eventName string, abistr string,
	client *ethclient.Client, ch chan types.Log) (sub ethereum.Subscription, err error) {
	//subcribe logs that will happen.
	q, err := buildQuery(contractAddress, fromBlock, toBlock, eventName, abistr)
	if err != nil {
		return nil, err
	}
	sub, err = client.SubscribeFilterLogs(ctx, *q, ch)
	if err != nil {
		return nil, err
	}
	return
	//node.DefaultIPCEndpoint("geth")
}

func EventSubscribe(contractAddress common.Address,
	eventName string, abistr string, client *helper.SafeEthClient, ch chan types.Log) (ethereum.Subscription, error) {
	return EventSubscribeInternal(context.Background(), contractAddress, rpc.EarliestBlockNumber, rpc.LatestBlockNumber,
		eventName, abistr, client.Client, ch)
}
