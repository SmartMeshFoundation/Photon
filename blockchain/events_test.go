package blockchain

import (
	"os"
	"testing"

	"github.com/SmartMeshFoundation/Photon/log"

	"github.com/SmartMeshFoundation/Photon/network/rpc"

	"fmt"

	"time"

	"math/big"

	"github.com/SmartMeshFoundation/Photon/codefortest"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
}

type fakeRPCModule struct {
	RegistryAddress       common.Address
	SecretRegistryAddress common.Address
}

func (r *fakeRPCModule) GetRegistryAddress() common.Address {
	return r.RegistryAddress
}

func (r *fakeRPCModule) GetSecretRegistryAddress() common.Address {
	return r.SecretRegistryAddress
}

type fakeChainEventRecordDao struct{}

func (f *fakeChainEventRecordDao) NewDeliveredChainEvent(id models.ChainEventID, blockNumber uint64) {
	return
}
func (f *fakeChainEventRecordDao) CheckChainEventDelivered(id models.ChainEventID) (blockNumber uint64, delivered bool) {
	return
}
func (f *fakeChainEventRecordDao) ClearOldChainEventRecord(blockNumber uint64) {
	return
}
func (f *fakeChainEventRecordDao) MakeChainEventID(l *types.Log) models.ChainEventID {
	return ""
}

func TestNewBlockChainEvents(t *testing.T) {
	client, err := codefortest.GetEthClient()
	if err != nil {
		panic(err)
	}
	be := NewBlockChainEvents(client, &fakeRPCModule{}, &fakeChainEventRecordDao{})
	if be == nil {
		t.Error("NewBlockChainEvents failed")
	}
}

func TestEvents_Start(t *testing.T) {
	client, err := codefortest.GetEthClient()
	if err != nil {
		panic(err)
	}
	be := NewBlockChainEvents(client, &fakeRPCModule{
		RegistryAddress: rpc.TestGetTokenNetworkRegistryAddress(),
	}, &fakeChainEventRecordDao{})
	if be == nil {
		t.Error("NewBlockChainEvents failed")
	}
	params.ChainID = big.NewInt(8888)
	be.Start(-1)
	begin := time.Now()
	for {
		if time.Since(begin) > 10*time.Second {
			be.Stop()
			time.Sleep(5 * time.Second)
			return
		}
		select {
		case sc := <-be.StateChangeChannel:
			fmt.Println(utils.StringInterface(sc, 0))
			//BlockStateChange, ok := sc.(transfer.BlockStateChange)
			//if ok {
			//	fmt.Println(BlockStateChange.BlockNumber)
			//}
		case t := <- be.EffectiveChainChan:
			fmt.Println("from be.EffectiveChainChan:", t)
		}
	}
}

func TestEvents_QueryAllStateChanges(t *testing.T) {
	client, err := codefortest.GetEthClient()
	if err != nil {
		panic(err)
	}
	logs, err := rpc.EventsGetInternal(
		rpc.GetQueryConext(), nil, 50000, 40000, client)
	fmt.Println(logs, err)
	if err != nil {
		return
	}
}

func TestEvents_QueryAllStateChanges2(t *testing.T) {
	client, err := codefortest.GetEthClient()
	if err != nil {
		panic(err)
	}
	logs, err := rpc.EventsGetInternal(
		rpc.GetQueryConext(), []common.Address{common.HexToAddress("0x71849b4f2fd77146f17298a363c1a750a14fc2ba")}, 13362235, 13362235, client)
	log.Trace(fmt.Sprintf("logs=%s", utils.StringInterface(logs, 5)))
	if err != nil {
		return
	}
}

func TestEvents_Start2(t *testing.T) {
	client, err := codefortest.GetEthClient()
	if err != nil {
		panic(err)
	}
	be := NewBlockChainEvents(client, &fakeRPCModule{
		RegistryAddress: common.HexToAddress("0x71849b4f2fd77146f17298a363c1a750a14fc2ba"),
	}, &fakeChainEventRecordDao{})
	if be == nil {
		t.Error("NewBlockChainEvents failed")
	}
	params.ChainID = big.NewInt(8888)
	chs, err := be.queryAllStateChange(13362234, 13362238)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("chs=%s", utils.StringInterface(chs, 5))
}
