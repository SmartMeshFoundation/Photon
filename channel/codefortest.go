package channel

import (
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/network/helper"

	"os"

	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/network/rpc"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
	"github.com/SmartMeshFoundation/Photon/notify"
	"github.com/SmartMeshFoundation/Photon/transfer/mtree"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// FakeTXINfoDao only for test
type FakeTXINfoDao struct{}

// NewPendingTXInfo :
func (dao *FakeTXINfoDao) NewPendingTXInfo(tx *types.Transaction, txType models.TXInfoType, channelIdentifier common.Hash, openBlockNumber int64, txParams models.TXParams) (txInfo *models.TXInfo, err error) {
	return
}

// SaveEventToTXInfo :
func (dao *FakeTXINfoDao) SaveEventToTXInfo(event interface{}) (txInfo *models.TXInfo, err error) {
	return
}

// UpdateTXInfoStatus :
func (dao *FakeTXINfoDao) UpdateTXInfoStatus(txHash common.Hash, status models.TXInfoStatus, pendingBlockNumber int64, gasUsed uint64) (txInfo *models.TXInfo, err error) {
	return
}

// GetTXInfoList :
func (dao *FakeTXINfoDao) GetTXInfoList(channelIdentifier common.Hash, openBlockNumber int64, tokenAddress common.Address, txType models.TXInfoType, status models.TXInfoStatus) (list []*models.TXInfo, err error) {
	return
}

func newTestBlockChainService() *rpc.BlockChainService {
	conn, err := helper.NewSafeClient(rpc.TestRPCEndpoint)
	if err != nil {
		log.Crit(fmt.Sprintf("Failed to connect to the Ethereum client: %s", err))
	}
	privkey, _ := utils.MakePrivateKeyAddress()
	if err != nil {
		log.Crit("Failed to create authorized transactor: ", err)
	}
	bcs, err := rpc.NewBlockChainService(privkey, rpc.PrivateRopstenRegistryAddress, conn, notify.NewNotifyHandler(), &FakeTXINfoDao{})
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
	channelID := utils.NewRandomHash()
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
