package rpc

import (
	"fmt"

	"os"

	"encoding/hex"

	"crypto/ecdsa"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//PrivateRopstenRegistryAddress test registry address, todo use env
var PrivateRopstenRegistryAddress = common.HexToAddress(os.Getenv("REGISTRY"))

//TestRPCEndpoint test eth rpc url, todo use env
var TestRPCEndpoint = os.Getenv("ETHRPCENDPOINT")

//TestPrivKey for test only
var TestPrivKey *ecdsa.PrivateKey

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
		fmt.Printf("Failed to connect to the Ethereum client: %s\n", err)
	}
	return NewBlockChainService(TestPrivKey, PrivateRopstenRegistryAddress, conn)
}

//GetTestChannelUniqueID for test only,get from env
func GetTestChannelUniqueID() *contracts.ChannelUniqueID {

	cu := &contracts.ChannelUniqueID{
		OpenBlockNumber: 3,
	}
	b, err := hex.DecodeString(os.Getenv("CHANNEL"))
	if err != nil || len(b) != len(cu.ChannelIdentifier) {
		panic("CHANNEL env error")
	}
	copy(cu.ChannelIdentifier[:], b)
	return cu
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
	addr := common.HexToAddress(os.Getenv("TOKEN_NETWORK_REGISTRY"))
	if addr == utils.EmptyAddress {
		panic("REGISTRY env error")
	}
	log.Info(fmt.Sprintf("TOKEN_NETWORK_REGISTRY=%s", addr.String()))
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
