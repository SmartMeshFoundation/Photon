package rpc

import (
	"math/big"

	"context"

	"github.com/SmartMeshFoundation/Photon/network/helper"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

//EventsGetInternal get events of history
func EventsGetInternal(ctx context.Context, contractsAddress []common.Address, fromBlock,
	toBlock int64, client *helper.SafeEthClient) ([]types.Log, error) {
	q, err := buildQueryBatch(contractsAddress, rpc.BlockNumber(fromBlock), rpc.BlockNumber(toBlock))
	if err != nil {
		return nil, err
	}
	ctx = ensureContext(ctx)
	return client.FilterLogs(ctx, *q)
}

func buildQueryBatch(contractsAddress []common.Address, fromBlock rpc.BlockNumber,
	toBlock rpc.BlockNumber) (q *ethereum.FilterQuery, err error) {
	q = &ethereum.FilterQuery{}
	q.FromBlock = big.NewInt(int64(fromBlock))
	if toBlock == rpc.LatestBlockNumber {
		q.ToBlock = nil
	} else {
		q.ToBlock = big.NewInt(int64(toBlock))
	}
	q.Addresses = contractsAddress
	return
}

func ensureContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.TODO()
	}
	return ctx
}
