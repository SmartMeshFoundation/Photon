package rpc

import (
	"fmt"

	"os"

	"encoding/hex"

	"crypto/ecdsa"

	"context"
	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/newtestenv/createchannel"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

var privkeybin = "2ddd679cb0f0754d0e20ef8206ea2210af3b51f159f55cfffbd8550f58daf779"

//PrivateRopstenRegistryAddress test registry address, todo use env
var PrivateRopstenRegistryAddress = common.HexToAddress(os.Getenv("REGISTRY")) // params.ROPSTEN_REGISTRY_ADDRESS
//TestRPCEndpoint test eth rpc url, todo use env
var TestRPCEndpoint = os.Getenv("ETHRPCENDPOINT")

//TestPrivKey for test only
var TestPrivKey *ecdsa.PrivateKey

func init() {
	keybin, err := hex.DecodeString(privkeybin)
	if err != nil {
		log.Crit("err %s", err)
	}
	TestPrivKey, err = crypto.ToECDSA(keybin)
	if err != nil {
		log.Crit("err %s", err)
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

/*
CreateChannelBetweenAddress these two address doesn't need have token and ether
*/
func CreateChannelBetweenAddress(client *ethclient.Client, addr1, addr2 common.Address, key1, key2 *ecdsa.PrivateKey) (channelAddress common.Address, err error) {
	token := os.Getenv("TOKEN")
	manager := os.Getenv("MANAGER")
	auth := bind.NewKeyedTransactor(TestPrivKey)

	mp, err := contracts.NewChannelManagerContract(common.HexToAddress(manager), client)
	if err != nil {
		return
	}
	tp, err := contracts.NewToken(common.HexToAddress(token), client)
	if err != nil {
		return
	}
	err = createchannel.TransferTo(client, TestPrivKey, addr1, big.NewInt(params.Ether))
	if err != nil {
		return
	}
	err = createchannel.TransferTo(client, TestPrivKey, addr2, big.NewInt(params.Ether))
	if err != nil {
		return
	}
	tx, err := tp.Transfer(auth, addr1, big.NewInt(50))
	if err != nil {
		return
	}
	_, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		return
	}
	tx, err = tp.Transfer(auth, addr2, big.NewInt(50))
	if err != nil {
		return
	}
	_, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		return
	}
	channelAddress = createchannel.CreatAChannelAndDeposit(addr1, addr2, key1, key2, 40, mp, tp, client)
	return
}
