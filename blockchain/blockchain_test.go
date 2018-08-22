package blockchain

import (
	"fmt"
	"os"

	"testing"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

var client *helper.SafeEthClient
var secretRegistryAddress common.Address
var TokenNetworkRegistryAddress common.Address
var be *Events
var at *AlarmTask

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
	setup()
}

func setup() {
	var err error
	client, err = helper.NewSafeClient(rpc.TestRPCEndpoint)
	if err != nil {
		panic(err)
	}
	TokenNetworkRegistryAddress = rpc.TestGetTokenNetworkRegistryAddress()
	tokenNetworkRegistry, err := contracts.NewTokenNetworkRegistry(TokenNetworkRegistryAddress, client)
	if err != nil {
		panic(err)
	}
	secretRegistryAddress, err = tokenNetworkRegistry.SecretRegistryAddress(nil)
	if err != nil {
		panic(err)
	}
	be = NewBlockChainEvents(client, TokenNetworkRegistryAddress, secretRegistryAddress, nil)
	tokens, err := be.GetAllTokenNetworks(0)
	if err != nil {
		panic(fmt.Sprintf("cannot get all token networks err %s", err))
	}
	if len(tokens) == 0 {
		panic(fmt.Sprintf("empty registyr network"))
	}
	at = NewAlarmTask(client)
}

func TestNewAlarmTask(t *testing.T) {
	assert.NotEmpty(t, at)
	assert.NotEmpty(t, at.client)
	assert.EqualValues(t, -1, at.LastBlockNumber)
	// test chan
	go func() {
		s, ok := <-at.quitChan
		assert.True(t, ok)
		assert.NotEmpty(t, s)
	}()
	at.quitChan <- struct{}{}
	assert.False(t, at.stopped)
	assert.EqualValues(t, time.Second, at.waitTime)
	assert.EqualValues(t, 0, len(at.callback))
	at.lock.Lock()
	at.lock.Unlock()
}

var testCb AlarmCallback = func(blockNumber int64) error {
	return nil
}

func TestAlarmTask_RegisterCallback(t *testing.T) {
	num := len(at.callback)
	for i := 0; i < 10; i++ {
		go func() {
			at.RegisterCallback(&testCb)
		}()
	}
	time.Sleep(1 * time.Second)
	assert.EqualValues(t, num+10, len(at.callback))
}
func TestAlarmTask_RemoveCallback1(t *testing.T) {
	num := len(at.callback)
	at.RegisterCallback(&testCb)
	at.RemoveCallback(&testCb)
	assert.EqualValues(t, num, len(at.callback))
}

func TestAlarmTask_RemoveCallback2(t *testing.T) {
	for i := 0; i < 10; i++ {
		go at.RegisterCallback(&testCb)
	}
	time.Sleep(1 * time.Second)
	num := len(at.callback)
	for i := 0; i < 10; i++ {
		go at.RemoveCallback(&testCb)
	}
	time.Sleep(1 * time.Second)
	assert.EqualValues(t, num-10, len(at.callback))
}

func TestAlarmTask_StartAndStop(t *testing.T) {
	for _, value := range at.callback {
		at.RemoveCallback(value)
	}
	assert.Equal(t, 0, len(at.callback))
	at.RegisterCallback(&testCb)
	oldBlockNo := at.LastBlockNumber
	at.Start()
	// Start for 20s
	fmt.Println("let AlarmTask runs for 20s...")
	for i := 0; i < 20; i++ {
		time.Sleep(time.Second * 1)
	}
	newBlockNum := at.LastBlockNumber
	assert.NotEqual(t, oldBlockNo, newBlockNum)
	assert.False(t, at.stopped)
	at.Stop()
	assert.True(t, at.stopped)
	_, ok := <-at.quitChan
	assert.False(t, ok)
}

func TestGetTokenNetworkCreated(t *testing.T) {
	//NewBlockChainEvents create BlockChainEvents
	tokens, err := be.GetAllTokenNetworks(0)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("tokens=%#v", tokens)
}

func TestEvents_GetAllChannels(t *testing.T) {
	channels, err := be.GetChannelNew(0, rpc.TestGetTokenNetworkAddress())
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("channels=%#v", channels)
}

func TestEvents_GetAllChannelClosed(t *testing.T) {
	events, err := be.GetChannelClosed(0, rpc.TestGetTokenNetworkAddress())
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("events=\n%s", utils.StringInterface(events, 3))
}

func TestEvents_GetAllChannelSettled(t *testing.T) {
	events, err := be.GetChannelSettled(0, rpc.TestGetTokenNetworkAddress())
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("events=\n%s", utils.StringInterface(events, 3))
}

func TestEvents_GetAllSecretRevealed(t *testing.T) {
	events, err := be.GetAllSecretRevealed(0)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("events=\n%s", utils.StringInterface(events, 3))
}

func TestEvents_GetChannelNewAndDeposit(t *testing.T) {
	events, err := be.GetChannelNewAndDeposit(0, utils.EmptyAddress)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("events=\n%s", utils.StringInterface(events, 3))
}
