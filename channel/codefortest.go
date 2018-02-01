package channel

import (
	"math/big"

	"fmt"

	"github.com/SmartMeshFoundation/raiden-network/network/helper"
	"github.com/SmartMeshFoundation/raiden-network/network/rpc"
	"github.com/SmartMeshFoundation/raiden-network/transfer"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gotips/log"
)

func newTestBlockChainService() *rpc.BlockChainService {
	conn, err := helper.NewSafeClient(rpc.TestRpcEndpoint)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to connect to the Ethereum client: %s", err))
	}
	privkey, _ := utils.MakePrivateKeyAddress()
	if err != nil {
		log.Error("Failed to create authorized transactor: ", err)
	}
	return rpc.NewBlockChainService(privkey, rpc.PRIVATE_ROPSTEN_REGISTRY_ADDRESS, conn)
}

func makeTestExternState() *ChannelExternalState {
	bcs := newTestBlockChainService()
	//must provide a valid netting channel address
	nettingChannel, _ := bcs.NettingChannel(common.HexToAddress("0x1309B5eAd287aD4A41CEB6df6A359529a561ECC6"))
	return NewChannelExternalState(func(channel *Channel, hashlock common.Hash) {}, nettingChannel, nettingChannel.Address, bcs)
}
func MakeTestPairChannel() (*Channel, *Channel) {
	tokenAddress := utils.NewRandomAddress()
	externState1 := makeTestExternState()
	externState2 := makeTestExternState()
	var balance1 = big.NewInt(330)
	var balance2 = big.NewInt(110)
	revealTimeout := 7
	settleTimeout := 30
	ourState := NewChannelEndState(externState1.bcs.NodeAddress, balance1, nil, transfer.EmptyMerkleTreeState)
	partnerState := NewChannelEndState(externState2.bcs.NodeAddress, balance2, nil, transfer.EmptyMerkleTreeState)

	testChannel, _ := NewChannel(ourState, partnerState, externState1, tokenAddress, externState1.ChannelAddress, externState1.bcs, revealTimeout, settleTimeout)

	ourState = NewChannelEndState(externState1.bcs.NodeAddress, balance1, nil, transfer.EmptyMerkleTreeState)
	partnerState = NewChannelEndState(externState2.bcs.NodeAddress, balance2, nil, transfer.EmptyMerkleTreeState)
	testChannel2, _ := NewChannel(partnerState, ourState, externState2, tokenAddress, externState2.ChannelAddress, externState2.bcs, revealTimeout, settleTimeout)
	return testChannel, testChannel2
}
