package smartraiden

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"path"

	"time"

	"encoding/hex"

	"sync"

	"github.com/SmartMeshFoundation/SmartRaiden/accounts"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/fee"
	"github.com/SmartMeshFoundation/SmartRaiden/notify"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//reinit this variable before test case start
var curAccountIndex = 0

func reinit() {
	curAccountIndex = 0
}
func newTestRaiden() *RaidenService {
	return newTestRaidenWithPolicy(&NoFeePolicy{})
}

func newTestRaidenWithPolicy(feePolicy fee.Charger) *RaidenService {
	bcs := newTestBlockChainService()
	notifyHandler := notify.NewNotifyHandler()
	transport := network.MakeTestMixTransport(utils.APex2(bcs.NodeAddress), bcs.PrivKey)
	config := params.DefaultConfig
	config.MyAddress = bcs.NodeAddress
	config.PrivateKey = bcs.PrivKey
	config.DataDir = os.Getenv("DATADIR")
	if config.DataDir == "" {
		config.DataDir = path.Join(os.TempDir(), utils.RandomString(10))
	}
	log.Info(fmt.Sprintf("DataDir=%s", config.DataDir))
	config.RevealTimeout = 10
	config.SettleTimeout = 600
	config.PrivateKeyHex = hex.EncodeToString(crypto.FromECDSA(config.PrivateKey))
	err := os.MkdirAll(config.DataDir, os.ModePerm)
	if err != nil {
		log.Error(err.Error())
	}
	config.DataBasePath = path.Join(config.DataDir, "log.db")
	config.NetworkMode = params.MixUDPXMPP
	rd, err := NewRaidenService(bcs, bcs.PrivKey, transport, &config, notifyHandler)
	if err != nil {
		log.Error(err.Error())
	}
	rd.SetFeePolicy(feePolicy)
	return rd
}
func newTestRaidenAPI() *RaidenAPI {
	api := NewRaidenAPI(newTestRaiden())
	err := api.Raiden.Start()
	if err != nil {
		panic(fmt.Sprintf("raiden start err %s", err))
	}
	return api
}

//maker sure these accounts are valid, and  engouh eths for test
func testGetnextValidAccount() (*ecdsa.PrivateKey, common.Address) {
	am := accounts.NewAccountManager("testdata/keystore")
	privkeybin, err := am.GetPrivateKey(am.Accounts[curAccountIndex].Address, "123")
	if err != nil {
		log.Error(fmt.Sprintf("testGetnextValidAccount err: %s", err))
		panic("")
	}
	curAccountIndex++
	privkey, err := crypto.ToECDSA(privkeybin)
	if err != nil {
		log.Error(fmt.Sprintf("to privkey err %s", err))
		panic("")
	}
	return privkey, crypto.PubkeyToAddress(privkey.PublicKey)
}
func newTestBlockChainService() *rpc.BlockChainService {
	conn, err := helper.NewSafeClient(rpc.TestRPCEndpoint)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to connect to the Ethereum client: %s", err))
	}
	privkey, addr := testGetnextValidAccount()
	log.Trace(fmt.Sprintf("privkey=%s,addr=%s", privkey, addr.String()))
	return rpc.NewBlockChainService(privkey, rpc.PrivateRopstenRegistryAddress, conn)
}

func makeTestRaidens() (r1, r2, r3 *RaidenService) {
	r1 = newTestRaiden()
	r2 = newTestRaiden()
	r3 = newTestRaiden()
	go func() {
		/*#nosec*/
		r1.Start()
	}()
	go func() {
		/*#nosec*/
		r2.Start()
	}()
	go func() {
		/*#nosec*/
		r3.Start()
	}()
	time.Sleep(time.Second * 3)
	return
}
func newTestRaidenAPIQuick() *RaidenAPI {
	api := NewRaidenAPI(newTestRaiden())
	//go func() {
	//	/*#nosec*/
	//	api.Raiden.Start()
	//}()
	return api
}

func makeTestRaidenAPIs() (rA, rB, rC, rD *RaidenAPI) {
	rA = newTestRaidenAPIQuick()
	rB = newTestRaidenAPIQuick()
	rC = newTestRaidenAPIQuick()
	rD = newTestRaidenAPIQuick()
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		/*#nosec*/
		rA.Raiden.Start()
		wg.Done()
	}()
	go func() {
		/*#nosec*/
		rB.Raiden.Start()
		wg.Done()
	}()
	go func() {
		/*#nosec*/
		rC.Raiden.Start()
		wg.Done()
	}()
	go func() {
		/*#nosec*/
		rD.Raiden.Start()
		wg.Done()
	}()
	wg.Wait()
	return
}

func makeTestRaidenAPIArrays(datadirs ...string) (apis []*RaidenAPI) {
	if datadirs == nil || len(datadirs) == 0 {
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(len(datadirs))
	for _, datadir := range datadirs {
		// #nosec
		os.Setenv("DATADIR", datadir)
		api := newTestRaidenAPIQuick()
		go func() {
			/*#nosec*/
			api.Raiden.Start()
			wg.Done()
		}()
		apis = append(apis, api)
	}
	wg.Wait()
	return
}

func makeTestRaidenAPIsWithFee(policy fee.Charger) (rA, rB, rC, rD *RaidenAPI) {
	rA = NewRaidenAPI(newTestRaidenWithPolicy(policy))
	rB = NewRaidenAPI(newTestRaidenWithPolicy(policy))
	rC = NewRaidenAPI(newTestRaidenWithPolicy(policy))
	rD = NewRaidenAPI(newTestRaidenWithPolicy(policy))
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		/*#nosec*/
		rA.Raiden.Start()
		wg.Done()
	}()
	go func() {
		/*#nosec*/
		rB.Raiden.Start()
		wg.Done()
	}()
	go func() {
		/*#nosec*/
		rC.Raiden.Start()
		wg.Done()
	}()
	go func() {
		/*#nosec*/
		rD.Raiden.Start()
		wg.Done()
	}()
	wg.Wait()
	return
}
