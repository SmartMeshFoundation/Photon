package testfilter

import (
	"testing"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

func TestStruct(t *testing.T) {
	tt := contracts.TokenNetworkRegistryTokenNetworkCreated{
		TokenAddress: utils.NewRandomAddress(),
	}
	t2 := tt
	t2.TokenAddress = utils.NewRandomAddress()
	t.Logf(fmt.Sprintf("tt=%s,t2=%s", tt.TokenAddress.String(), t2.TokenAddress.String()))
}

func TestFilter1(t *testing.T) {
	//client, err := ethclient.Dial(rpc.TestRPCEndpoint)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//registryAddr := common.HexToAddress(os.Getenv("TOKEN_NETWORK_REGISTRY"))
	////registry, err := contracts.NewTokenNetworkRegistry(registryAddr, client)
	////if err != nil {
	////	t.Error(err)
	////	return
	////}
	//f, err := contracts.NewTokenNetworkRegistryFilterer(registryAddr, client)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//it, err := f.FilterTokenNetworkCreated(nil, nil, nil)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//for it.Next() {
	//	fmt.Printf("event=%s", utils.StringInterface(it.Event, 3))
	//}
	////ch := make(chan *contracts.TokenNetworkRegistryTokenNetworkCreated, 10)
	////var start uint64
	////sub, err := f.WatchTokenNetworkCreated(&bind.WatchOpts{
	////	Start: &start,
	////}, ch, nil, nil)
	////if err != nil {
	////	t.Error(err)
	////	return
	////}
	////for {
	////	select {
	////	case <-time.After(1 * time.Second):
	////		sub.Unsubscribe()
	////		return
	////	case e := <-ch:
	////		fmt.Printf("sub event=%s", utils.StringInterface(e, 3))
	////	}
	////}
}

func TestFilter2(t *testing.T) {
	//client, err := ethclient.Dial(rpc.TestRPCEndpoint)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//tokenNetworkAddr := common.HexToAddress(os.Getenv("TOKEN_NETWORK"))
	//channelIdentifier := "0x53d73bbd584d3154f49db7f5674c7118a08079a0b285016fffa1dc34e4e89fd9"
	//t.Logf("TOKEN_NETWORK = %s\n", tokenNetworkAddr.String())
	//t.Logf("CHANNEL_IDENTIFIER = %s\n", channelIdentifier)
	//f, err := contracts.NewTokenNetworkFilterer(tokenNetworkAddr, client)
	//ch := make(chan *contracts.TokenNetworkChannelNewDeposit, 10)
	//var start uint64
	//sub, err := f.WatchChannelNewDeposit(&bind.WatchOpts{
	//	Start: &start,
	//}, ch, [][32]byte{common.HexToHash(channelIdentifier)})
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//for {
	//	select {
	//	case <-time.After(50 * time.Second):
	//		sub.Unsubscribe()
	//		fmt.Println("timeout")
	//		return
	//	case e := <-ch:
	//		fmt.Printf("sub event=%s", utils.StringInterface(e, 3))
	//	}
	//}
}
