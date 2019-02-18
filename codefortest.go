package photon

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"path"

	"time"

	"sync"

	"github.com/SmartMeshFoundation/Photon/accounts"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/models/stormdb"
	"github.com/SmartMeshFoundation/Photon/network"
	"github.com/SmartMeshFoundation/Photon/network/helper"
	"github.com/SmartMeshFoundation/Photon/network/rpc"
	"github.com/SmartMeshFoundation/Photon/network/rpc/fee"
	"github.com/SmartMeshFoundation/Photon/notify"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//reinit this variable before test case start
var curAccountIndex = 0

func reinit() {
	curAccountIndex = 0
}
func newTestPhoton() *Service {
	return newTestPhotonWithPolicy(&NoFeePolicy{})
}

func newTestPhotonWithPolicy(feePolicy fee.Charger) *Service {
	config := params.DefaultConfig
	config.DataDir = os.Getenv("DATADIR")
	if config.DataDir == "" {
		config.DataDir = path.Join(os.TempDir(), utils.RandomString(10))
	}
	config.DataBasePath = path.Join(config.DataDir, "log.dao")
	db, err := stormdb.OpenDb(config.DataBasePath)
	if err != nil {
		panic(err)
	}
	bcs := newTestBlockChainService(db)
	notifyHandler := notify.NewNotifyHandler()
	transport := network.MakeTestMixTransport(utils.APex2(bcs.NodeAddress), bcs.PrivKey)
	config.MyAddress = bcs.NodeAddress
	config.PrivateKey = bcs.PrivKey
	log.Info(fmt.Sprintf("DataDir=%s", config.DataDir))
	config.RevealTimeout = 10
	config.SettleTimeout = 600
	err = os.MkdirAll(config.DataDir, os.ModePerm)
	if err != nil {
		log.Error(err.Error())
	}
	config.NetworkMode = params.MixUDPXMPP
	rd, err := NewPhotonService(bcs, bcs.PrivKey, transport, &config, notifyHandler, db)
	if err != nil {
		log.Error(err.Error())
	}
	rd.FeePolicy = feePolicy
	return rd
}
func newTestPhotonAPI() *API {
	api := NewPhotonAPI(newTestPhoton())
	err := api.Photon.Start()
	if err != nil {
		panic(fmt.Sprintf("Photon start err %s", err))
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
func newTestBlockChainService(dao models.Dao) *rpc.BlockChainService {
	conn, err := helper.NewSafeClient(rpc.TestRPCEndpoint)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to connect to the Ethereum client: %s", err))
	}
	privkey, addr := testGetnextValidAccount()
	log.Trace(fmt.Sprintf("privkey=%s,addr=%s", privkey, addr.String()))
	bcs, err := rpc.NewBlockChainService(privkey, rpc.PrivateRopstenRegistryAddress, conn, notify.NewNotifyHandler(), nil)
	if err != nil {
		log.Error(err.Error())
	}
	return bcs
}

func makeTestPhotons() (r1, r2, r3 *Service) {
	r1 = newTestPhoton()
	r2 = newTestPhoton()
	r3 = newTestPhoton()
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
func newTestPhotonAPIQuick() *API {
	api := NewPhotonAPI(newTestPhoton())
	//go func() {
	//	/*#nosec*/
	//	api.Photon.Start()
	//}()
	return api
}

func makeTestPhotonAPIs() (rA, rB, rC, rD *API) {
	rA = newTestPhotonAPIQuick()
	rB = newTestPhotonAPIQuick()
	rC = newTestPhotonAPIQuick()
	rD = newTestPhotonAPIQuick()
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		/*#nosec*/
		rA.Photon.Start()
		wg.Done()
	}()
	go func() {
		/*#nosec*/
		rB.Photon.Start()
		wg.Done()
	}()
	go func() {
		/*#nosec*/
		rC.Photon.Start()
		wg.Done()
	}()
	go func() {
		/*#nosec*/
		rD.Photon.Start()
		wg.Done()
	}()
	wg.Wait()
	return
}

func makeTestPhotonAPIArrays(datadirs ...string) (apis []*API) {
	if datadirs == nil || len(datadirs) == 0 {
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(len(datadirs))
	for _, datadir := range datadirs {
		// #nosec
		os.Setenv("DATADIR", datadir)
		api := newTestPhotonAPIQuick()
		go func() {
			/*#nosec*/
			api.Photon.Start()
			wg.Done()
		}()
		apis = append(apis, api)
	}
	wg.Wait()
	return
}

func makeTestPhotonAPIsWithFee(policy fee.Charger) (rA, rB, rC, rD *API) {
	rA = NewPhotonAPI(newTestPhotonWithPolicy(policy))
	rB = NewPhotonAPI(newTestPhotonWithPolicy(policy))
	rC = NewPhotonAPI(newTestPhotonWithPolicy(policy))
	rD = NewPhotonAPI(newTestPhotonWithPolicy(policy))
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		/*#nosec*/
		rA.Photon.Start()
		wg.Done()
	}()
	go func() {
		/*#nosec*/
		rB.Photon.Start()
		wg.Done()
	}()
	go func() {
		/*#nosec*/
		rC.Photon.Start()
		wg.Done()
	}()
	go func() {
		/*#nosec*/
		rD.Photon.Start()
		wg.Done()
	}()
	wg.Wait()
	return
}
