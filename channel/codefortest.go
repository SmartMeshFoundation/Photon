package channel

import (
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"

	"os"

	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mtree"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

func newTestBlockChainService() *rpc.BlockChainService {
	conn, err := helper.NewSafeClient(rpc.TestRPCEndpoint)
	if err != nil {
		log.Crit(fmt.Sprintf("Failed to connect to the Ethereum client: %s", err))
	}
	privkey, _ := utils.MakePrivateKeyAddress()
	if err != nil {
		log.Crit("Failed to create authorized transactor: ", err)
	}
	bcs, err := rpc.NewBlockChainService(privkey, rpc.PrivateRopstenRegistryAddress, conn)
	if err != nil {
		panic(err)
	}
	return bcs
}

var testFuncRegisterChannelForHashlock = func(channel *Channel, hashlock common.Hash) {}

func makeTestExternState() *ExternalState {
	bcs := newTestBlockChainService()
	//must provide a valid netting channel address
	tokenNetwork, err := bcs.TokenNetwork(common.HexToAddress(os.Getenv("TOKEN_NETWORK")))
	if err != nil {
		panic(err)
	}
	channelID := common.HexToHash(os.Getenv("CHANNEL"))
	channelIdentifer := &contracts.ChannelUniqueID{
		ChannelIdentifier: channelID,
		OpenBlockNumber:   3,
	}
	return NewChannelExternalState(testFuncRegisterChannelForHashlock,
		tokenNetwork, channelIdentifer,
		bcs.PrivKey, bcs.Client,
		nil, 0,
		bcs.NodeAddress, utils.NewRandomAddress(),
	)
}

//MakeTestPairChannel for test
func MakeTestPairChannel() (*Channel, *Channel) {
	externState1 := makeTestExternState()
	externState2 := makeTestExternState()
	var balance1 = big.NewInt(330)
	var balance2 = big.NewInt(110)
	tokenAddr := utils.NewRandomAddress()
	revealTimeout := 7
	settleTimeout := 30
	ourState := NewChannelEndState(externState1.MyAddress, balance1, nil, mtree.EmptyTree)
	partnerState := NewChannelEndState(externState2.MyAddress, balance2, nil, mtree.EmptyTree)
	//#nosec
	testChannel, _ := NewChannel(ourState, partnerState, externState1, tokenAddr, &externState1.ChannelIdentifier, revealTimeout, settleTimeout)

	ourState = NewChannelEndState(externState1.MyAddress, balance1, nil, mtree.EmptyTree)
	partnerState = NewChannelEndState(externState2.MyAddress, balance2, nil, mtree.EmptyTree)
	//#nosec
	testChannel2, _ := NewChannel(partnerState, ourState, externState2, tokenAddr, &externState2.ChannelIdentifier, revealTimeout, settleTimeout)
	return testChannel, testChannel2
}
