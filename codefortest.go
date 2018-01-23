package raiden_network

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"path"

	"time"

	"encoding/hex"

	"github.com/SmartMeshFoundation/raiden-network/network"
	"github.com/SmartMeshFoundation/raiden-network/network/helper"
	"github.com/SmartMeshFoundation/raiden-network/network/rpc"
	"github.com/SmartMeshFoundation/raiden-network/params"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
)

var curAccountIndex = 0

func newTestRaiden() *RaidenService {
	transport := network.MakeTestUDPTransport(50000 + curAccountIndex + 1)
	discover := network.NewHttpDiscovery() //share the same discovery ,so node can find each other
	bcs := newTestBlockChainService()
	config := params.DefaultConfig
	config.MyAddress = bcs.NodeAddress
	config.PrivateKey = bcs.PrivKey
	config.DataDir = path.Join(os.TempDir(), utils.RandomString(10))
	config.ExternIp = transport.Host
	config.ExternPort = transport.Port
	config.Host = transport.Host
	config.Port = transport.Port
	config.RevealTimeout = 4
	config.SettleTimeout = 20
	config.PrivateKeyHex = hex.EncodeToString(crypto.FromECDSA(config.PrivateKey))
	os.MkdirAll(config.DataDir, os.ModePerm)
	config.DataBasePath = path.Join(config.DataDir, "log.db")
	rd := NewRaidenService(bcs, bcs.PrivKey, transport, discover, &config)
	return rd
}
func newTestRaidenApi() *RaidenApi {
	api := NewRaidenApi(newTestRaiden())
	go func() {
		api.Raiden.Start()
	}()
	time.Sleep(time.Second * 3)
	return api
}

//maker sure these accounts are valid, and  engouh eths for test
func testGetnextValidAccount() (*ecdsa.PrivateKey, common.Address) {
	am := NewAccountManager("d:\\privnet\\keystore")
	privkey, err := am.GetPrivateKey(am.Accounts[curAccountIndex].Address, "123")
	if err != nil {
		fmt.Sprintf("testGetnextValidAccount err:", err)
		panic("")
	}
	curAccountIndex++
	return crypto.ToECDSAUnsafe(privkey), utils.PubkeyToAddress(privkey)
}
func newTestBlockChainService() *rpc.BlockChainService {
	conn, err := helper.NewSafeClient(node.DefaultIPCEndpoint("geth"))
	if err != nil {
		log.Error("Failed to connect to the Ethereum client: ", err)
	}
	privkey, _ := testGetnextValidAccount()
	if err != nil {
		log.Error("Failed to create authorized transactor: ", err)
	}
	return rpc.NewBlockChainService(privkey, params.ROPSTEN_REGISTRY_ADDRESS, conn)
}

func makeTestRaidens() (r1, r2, r3 *RaidenService) {
	r1 = newTestRaiden()
	r2 = newTestRaiden()
	r3 = newTestRaiden()
	go func() {
		r1.Start()
	}()
	go func() {
		r2.Start()
	}()
	go func() {
		r3.Start()
	}()
	time.Sleep(time.Second * 3)
	return
}
func newTestRaidenApiQuick() *RaidenApi {
	api := NewRaidenApi(newTestRaiden())
	go func() {
		api.Raiden.Start()
	}()
	return api
}
func makeTestRaidenApis() (rA, rB, rC, rD *RaidenApi) {
	rA = newTestRaidenApiQuick()
	rB = newTestRaidenApiQuick()
	rC = newTestRaidenApiQuick()
	rD = newTestRaidenApiQuick()
	go func() {
		rA.Raiden.Start()
	}()
	go func() {
		rB.Raiden.Start()
	}()
	go func() {
		rC.Raiden.Start()
	}()
	go func() {
		rD.Raiden.Start()
	}()
	time.Sleep(time.Second * 3)
	return
}
