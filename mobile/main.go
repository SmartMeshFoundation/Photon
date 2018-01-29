package mobile

import (
	"os"

	"fmt"

	"io/ioutil"

	"encoding/hex"

	"path"

	"path/filepath"

	"os/signal"
	"runtime"
	"time"

	"github.com/SmartMeshFoundation/raiden-network"
	"github.com/SmartMeshFoundation/raiden-network/network"
	"github.com/SmartMeshFoundation/raiden-network/network/helper"
	"github.com/SmartMeshFoundation/raiden-network/network/rpc"
	"github.com/SmartMeshFoundation/raiden-network/params"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/slonzok/getpass"
)

var (
	argAddress                  string
	argKeyStorePath             string
	argEthRpcEndpoint           string
	argRegistryContractAddress  string = params.ROPSTEN_REGISTRY_ADDRESS.String()
	argDiscoveryContractAddress string = params.ROPSTEN_DISCOVERY_ADDRESS.String()
	argListenAddress            string = "0.0.0.0:40001"
	argApiAddress               string = "0.0.0.0:5001"
	argDataDir                  string
	argPasswordFile             string
	argNat                      string = "stun"
	argLogging                  string = "trace"
	argLogfile                  string = ""
)

/*
address :节点地址,比如0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa
keystorePath:存放私钥的目录地址 geth的keystore目录
ethRpcEndPoint: 连接geth的url 比如:ws://10.0.0.2:8546
dataDir:节点工作目录,存放数据用,要求必须可写
passwordfile:一个只包含私钥密码的文本文件
*/
func MobileStartUp(address, keystorePath, ethRpcEndPoint, dataDir, passwordfile string) (api *Api, err error) {
	argAddress = address
	argKeyStorePath = keystorePath
	argEthRpcEndpoint = ethRpcEndPoint
	argDataDir = dataDir
	argPasswordFile = passwordfile
	return mobileMain()
}
func setupLog() {
	loglevel := argLogging
	writer := os.Stderr
	lvl := log.LvlTrace
	switch loglevel {
	case "trace":
		lvl = log.LvlTrace
	case "debug":
		lvl = log.LvlDebug
	case "info":
		lvl = log.LvlInfo
	case "warn":
		lvl = log.LvlWarn
	case "error":
		lvl = log.LvlError
	case "critical":
		lvl = log.LvlCrit
	}
	logfilename := argLogfile
	if len(logfilename) > 0 {
		file, err := os.Create(logfilename)
		if err != nil {
			fmt.Printf("open logfile %s error:%s\n", logfilename, err)
			utils.SystemExit(1)
		}
		writer = file
	}
	fmt.Println("loglevel:", lvl.String())
	log.Root().SetHandler(log.LvlFilterHandler(lvl, log.StreamHandler(writer, log.TerminalFormat(true))))
}
func mobileMain() (api *Api, err error) {
	fmt.Printf("Welcom to GoRaiden,version %f\n", 0.1)
	//promptAccount(utils.EmptyAddress, `D:\privnet\keystore\`, "")
	setupLog()
	/*
	  TODO:
	        - Ask for confirmation to quit if there are any locked transfers that did
	        not timeout.
	*/
	host, port := network.SplitHostPort(argListenAddress)
	pms, err := network.SocketFactory(host, port, argNat)
	log.Trace(fmt.Sprintf("pms=%s", utils.StringInterface1(pms)))
	if err != nil {
		log.Error(fmt.Sprintf("start server on %s error:%s", argListenAddress, err))
		utils.SystemExit(1)
	}
	cfg := config(pms)
	log.Trace(fmt.Sprintf("cfg=", spew.Sdump(cfg)))
	//spew.Dump("Config:", cfg)
	ethEndpoint := argEthRpcEndpoint
	client, err := helper.NewSafeClient(ethEndpoint)
	if err != nil {
		log.Error(fmt.Sprintf("cannot connect to geth :%s err=%s", ethEndpoint, err))
		utils.SystemExit(1)
	}
	bcs := rpc.NewBlockChainService(cfg.PrivateKey, cfg.RegistryAddress, client)
	log.Trace(fmt.Sprintf("bcs=%#v", bcs))
	discovery := network.NewContractDiscovery(bcs.NodeAddress, bcs.Client, bcs.Auth)
	//discovery := network.NewHttpDiscovery()
	policy := network.NewTokenBucket(10, 1, time.Now)
	transport := network.NewUDPTransport(host, port, pms.Conn, nil, policy)
	raidenService := raiden_network.NewRaidenService(bcs, cfg.PrivateKey, transport, discovery, cfg)
	go func() {
		//startup may take long time and then enter loop
		raidenService.Start()
	}()
	api = &Api{raiden_network.NewRaidenApi(raidenService)}
	regQuitHandler(api.api)
	time.Sleep(time.Millisecond * 10) //let go func run
	raidenService.StartWg.Wait()
	return api, nil
}
func regQuitHandler(api *raiden_network.RaidenApi) {
	go func() {
		quitSignal := make(chan os.Signal, 1)
		signal.Notify(quitSignal, os.Interrupt, os.Kill)
		<-quitSignal
		signal.Stop(quitSignal)
		api.Stop()
		utils.SystemExit(0)
	}()
}
func promptAccount(adviceAddress common.Address, keystorePath, passwordfile string) (addr common.Address, keybin []byte) {
	am := raiden_network.NewAccountManager(keystorePath)
	if len(am.Accounts) == 0 {
		log.Error(fmt.Sprintf("No Ethereum accounts found in the directory %s", keystorePath))
		utils.SystemExit(1)
	}
	if !am.AddressInKeyStore(adviceAddress) {
		if adviceAddress != utils.EmptyAddress {
			log.Error(fmt.Sprintf("account %s could not be found on the sytstem. aborting...", adviceAddress))
			utils.SystemExit(1)
		}
		shouldPromt := true
		fmt.Println("The following accounts were found in your machine:")
		for i := 0; i < len(am.Accounts); i++ {
			fmt.Printf("%3d -  %s\n", i, am.Accounts[i].Address.String())
		}
		fmt.Println("")
		for shouldPromt {
			fmt.Printf("Select one of them by index to continue:\n")
			idx := -1
			fmt.Scanf("%d", &idx)
			if idx >= 0 && idx < len(am.Accounts) {
				shouldPromt = false
				addr = am.Accounts[idx].Address
			} else {
				fmt.Printf("Error: Provided index %d is out of bounds", idx)
			}
		}
	} else {
		addr = adviceAddress
	}
	var password string
	var err error
	if len(passwordfile) > 0 {
		data, err := ioutil.ReadFile(passwordfile)
		if err != nil {
			log.Error(fmt.Sprintf("password_file error:%s", err))
			utils.SystemExit(1)
		}
		password = string(data)
		log.Trace(fmt.Sprintf("password is %s", password))
		keybin, err = am.GetPrivateKey(addr, password)
		if err != nil {
			log.Error(fmt.Sprintf("Incorrect password for %s in file. Aborting ... %s", addr.String(), err))
			utils.SystemExit(1)
		}
	} else {
		for i := 0; i < 3; i++ {
			//retries three times
			password = getpass.Prompt("Enter the password to unlock:")
			keybin, err = am.GetPrivateKey(addr, password)
			if err != nil && i == 3 {
				log.Error(fmt.Sprintf("Exhausted passphrase unlock attempts for %s. Aborting ...", addr))
				utils.SystemExit(1)
			}
			if err != nil {
				log.Error(fmt.Sprintf("password incorrect\n Please try again or kill the process to quit.\nUsually Ctrl-c."))
				continue
			}
			break
		}
	}
	return
}
func config(pms *network.PortMappedSocket) *params.Config {
	var err error
	config := params.DefaultConfig
	listenhost, listenport := network.SplitHostPort(argListenAddress)
	apihost, apiport := network.SplitHostPort(argApiAddress)
	config.Host = listenhost
	config.Port = listenport
	config.UseConsole = false
	config.UseRpc = false
	config.ApiHost = apihost
	config.ApiPort = apiport
	config.ExternIp = pms.ExternalIp
	config.ExternPort = pms.ExternalPort
	max_unresponsive_time := int64(time.Minute)
	config.Protocol.NatKeepAliveTimeout = max_unresponsive_time / params.DEFAULT_NAT_KEEPALIVE_RETRIES
	address := common.HexToAddress(argAddress)
	address, privkeyBin := promptAccount(address, argKeyStorePath, argPasswordFile)
	config.PrivateKeyHex = hex.EncodeToString(privkeyBin)
	config.PrivateKey, err = crypto.ToECDSA(privkeyBin)
	log.Trace("private key")
	config.MyAddress = address
	if err != nil {
		log.Error(fmt.Sprintf("privkey error:%s", err))
		utils.SystemExit(1)
	}
	registAddrStr := argRegistryContractAddress
	if len(registAddrStr) > 0 {
		config.RegistryAddress = common.HexToAddress(registAddrStr)
	}
	discoverAddr := argDiscoveryContractAddress
	if len(discoverAddr) > 0 {
		config.DiscoveryAddress = common.HexToAddress(discoverAddr)
	}
	dataDir := argDataDir
	if len(dataDir) == 0 {
		dataDir = path.Join(utils.GetHomePath(), ".goraiden")
	}
	log.Trace("start dir...")
	config.DataDir = dataDir
	if !utils.Exists(config.DataDir) {
		err = os.MkdirAll(config.DataDir, os.ModePerm)
		if err != nil {
			log.Error(fmt.Sprintf("Datadir:%s doesn't exist and cannot create %v", config.DataDir, err))
			utils.SystemExit(1)
		}
	}
	userDbPath := hex.EncodeToString(config.MyAddress[:])
	userDbPath = userDbPath[:8]
	userDbPath = filepath.Join(config.DataDir, userDbPath)
	log.Trace("db dir")
	if !utils.Exists(userDbPath) {
		err = os.MkdirAll(userDbPath, os.ModePerm)
		if err != nil {
			log.Error(fmt.Sprintf("Datadir:%s doesn't exist and cannot create %v", userDbPath, err))
			utils.SystemExit(1)
		}
	}
	databasePath := filepath.Join(userDbPath, "log.db")
	config.DataBasePath = databasePath
	return &config
}
func init() {
	//many race condtions don't resolve
	setNativeThreadNumber()
}
func setNativeThreadNumber() {
	runtime.GOMAXPROCS(1)
}
