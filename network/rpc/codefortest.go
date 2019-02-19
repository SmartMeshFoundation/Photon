package rpc

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/network/helper"
	"github.com/SmartMeshFoundation/Photon/notify"

	"os"

	"encoding/hex"

	"crypto/ecdsa"

	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

//PrivateRopstenRegistryAddress test registry address, todo use env
var PrivateRopstenRegistryAddress = common.HexToAddress(os.Getenv("TOKEN_NETWORK"))

//TestRPCEndpoint test eth rpc url, todo use env
var TestRPCEndpoint = os.Getenv("ETHRPCENDPOINT")

//TestPrivKey for test only
var TestPrivKey *ecdsa.PrivateKey

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
func (dao *FakeTXINfoDao) UpdateTXInfoStatus(txHash common.Hash, status models.TXInfoStatus, pendingBlockNumber int64) (err error) {
	return
}

// GetTXInfoList :
func (dao *FakeTXINfoDao) GetTXInfoList(channelIdentifier common.Hash, openBlockNumber int64, tokenAddress common.Address, txType models.TXInfoType, status models.TXInfoStatus) (list []*models.TXInfo, err error) {
	return
}

func init() {
	if encoding.IsTest {
		keybin, err := hex.DecodeString(os.Getenv("KEY1"))
		if err != nil {
			//启动错误不要用 log, 那时候 log 还没准备好
			// do not use log to print start error, it's not ready
			panic(fmt.Sprintf("err %s", err))
		}
		TestPrivKey, err = crypto.ToECDSA(keybin)
		if err != nil {
			panic(fmt.Sprintf("err %s", err))
		}
	}

}

//MakeTestBlockChainService creat test BlockChainService
func MakeTestBlockChainService() *BlockChainService {
	conn, err := helper.NewSafeClient(TestRPCEndpoint)
	//conn, err := ethclient.Dial("ws://" + node.DefaultWSEndpoint())
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the Ethereum client: %s\n", err))
	}
	bcs, err := NewBlockChainService(TestPrivKey, PrivateRopstenRegistryAddress, conn, notify.NewNotifyHandler(), &FakeTXINfoDao{})
	if err != nil {
		panic(err)
	}
	return bcs
}

//TestGetTokenNetworkAddress for test only
func TestGetTokenNetworkAddress() common.Address {
	addr := common.HexToAddress(os.Getenv("TOKEN_NETWORK"))
	if addr == utils.EmptyAddress {
		panic("TOKENNETWORK env error")
	}
	log.Trace(fmt.Sprintf("test TOKEN_NETWORK=%s ", addr.String()))
	return addr
}

//TestGetTokenNetworkRegistryAddress for test only
func TestGetTokenNetworkRegistryAddress() common.Address {
	addr := common.HexToAddress(os.Getenv("TOKEN_NETWORK"))
	if addr == utils.EmptyAddress {
		panic("REGISTRY env error")
	}
	log.Info(fmt.Sprintf("TOKEN_NETWORK=%s", addr.String()))
	return addr
}

//TestGetParticipant1 for test only
func TestGetParticipant1() (privKey *ecdsa.PrivateKey, addr common.Address) {
	keybin, err := hex.DecodeString(os.Getenv("KEY1"))
	if err != nil {
		panic("KEY1 ERRor")
	}
	return testGetParticipant(keybin)
}

//TestGetParticipant2 for test only
func TestGetParticipant2() (privKey *ecdsa.PrivateKey, addr common.Address) {
	keybin, err := hex.DecodeString(os.Getenv("KEY2"))
	if err != nil {
		panic("KEY1 ERRor")
	}
	return testGetParticipant(keybin)
}
func testGetParticipant(keybin []byte) (privKey *ecdsa.PrivateKey, addr common.Address) {
	privKey, err := crypto.ToECDSA(keybin)
	if err != nil {
		panic(fmt.Sprintf("toecda err %s,keybin=%s", err, hex.EncodeToString(keybin)))
	}
	addr = crypto.PubkeyToAddress(privKey.PublicKey)
	return
}

/*
CreateChannelBetweenAddress these two address doesn't need have token and ether
*/
//func CreateChannelBetweenAddress(client *ethclient.Client, addr1, addr2 common.Address, key1, key2 *ecdsa.PrivateKey) (err error) {
//	token := os.Getenv("TOKEN")
//	tokenNetwork := os.Getenv("TOKENNETWORK")
//	auth := bind.NewKeyedTransactor(TestPrivKey)
//	tokenNetworkAddress := common.HexToAddress(tokenNetwork)
//	tp, err := contracts.NewToken(common.HexToAddress(token), client)
//	if err != nil {
//		return
//	}
//	err = createchannel.TransferTo(client, TestPrivKey, addr1, big.NewInt(params.Ether))
//	if err != nil {
//		return
//	}
//	err = createchannel.TransferTo(client, TestPrivKey, addr2, big.NewInt(params.Ether))
//	if err != nil {
//		return
//	}
//	tx, err := tp.Transfer(auth, addr1, big.NewInt(500))
//	if err != nil {
//		return
//	}
//	_, err = bind.WaitMined(context.Background(), client, tx)
//	if err != nil {
//		return
//	}
//	tx, err = tp.Transfer(auth, addr2, big.NewInt(500))
//	if err != nil {
//		return
//	}
//	_, err = bind.WaitMined(context.Background(), client, tx)
//	if err != nil {
//		return
//	}
//	createchannel.CreatAChannelAndDeposit(addr1, addr2, key1, key2, 40, tokenNetworkAddress, tp, client)
//	return
//}
