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

	"runtime/debug"

	"github.com/SmartMeshFoundation/SmartRaiden"
	"github.com/SmartMeshFoundation/SmartRaiden/cmd/smartraiden/mainimpl"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/slonzok/getpass"
)

var (
	argAddress                  string
	argKeyStorePath             string
	argEthRPCEndpoint           string
	argRegistryContractAddress  = params.RopstenRegistryAddress.String()
	argDiscoveryContractAddress = params.RopstenDiscoveryAddress.String()
	argListenAddress            = "0.0.0.0:40001"
	argAPIAddress               = "0.0.0.0:5001"
	argDataDir                  string
	argPasswordFile             string
	argNat                      = "stun"
	argLogging                  = "trace"
	argLogfile                  = ""
)

func init() {
	debug.SetTraceback("crash")
}

func panicOnNullValue() {
	var c []int
	c[0] = 0
}

/*
StartUp is entry point for mobile raiden
address :Node address,such as 0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa
keystorePath:The address of the private key,  geth keystore directory . eg ~/.geth/keystore
ethRpcEndPoint:URL connected to geth ,such as:ws://10.0.0.2:8546
dataDir:The working directory of a node, such as ~/.smartraiden
passwordfile: file to storage password eg ~/.geth/pass.txt
apiAddr: 127.0.0.1:5001 for product,0.0.0.1:5001 for test
*/
func StartUp(address, keystorePath, ethRPCEndPoint, dataDir, passwordfile, apiAddr string) {
	argAddress = address
	argKeyStorePath = keystorePath
	argEthRPCEndpoint = ethRPCEndPoint
	argDataDir = dataDir
	argPasswordFile = passwordfile
	os.Args = make([]string, 0, 20)
	os.Args = append(os.Args, "smartraidenmobile")
	os.Args = append(os.Args, fmt.Sprintf("--address=%s", address))
	os.Args = append(os.Args, fmt.Sprintf("--keystore-path=%s", keystorePath))
	os.Args = append(os.Args, fmt.Sprintf("--eth-rpc-endpoint=%s", ethRPCEndPoint))
	os.Args = append(os.Args, fmt.Sprintf("--datadir=%s", dataDir))
	os.Args = append(os.Args, fmt.Sprintf("--password-file=%s", passwordfile))
	os.Args = append(os.Args, fmt.Sprintf("--nat=ice"))
	os.Args = append(os.Args, fmt.Sprintf("--ignore-mediatednode-request"))
	os.Args = append(os.Args, fmt.Sprintf("--api-address=%s", apiAddr))
	os.Args = append(os.Args, fmt.Sprintf("--verbosity=5"))
	os.Args = append(os.Args, fmt.Sprintf("--debug"))
	os.Args = append(os.Args, fmt.Sprintf("--enable-health-check"))
	//panicOnNullValue()
	mainimpl.StartMain()
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
	log.Root().SetHandler(log.LvlFilterHandler(lvl, utils.MyStreamHandler(writer)))
}
func mobileMain() (api *API, err error) {
	fmt.Printf("Welcom to smartraiden,version %f\n", 0.1)
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
	//log.Trace(fmt.Sprintf("cfg=", spew.Sdump(cfg)))
	//spew.Dump("Config:", cfg)
	ethEndpoint := argEthRPCEndpoint
	client, err := helper.NewSafeClient(ethEndpoint)
	if err != nil {
		log.Error(fmt.Sprintf("cannot connect to geth :%s err=%s", ethEndpoint, err))
		utils.SystemExit(1)
	}
	bcs := rpc.NewBlockChainService(cfg.PrivateKey, cfg.RegistryAddress, client)
	discovery := network.NewContractDiscovery(bcs.NodeAddress, common.HexToAddress(argDiscoveryContractAddress), bcs.Client, bcs.Auth)
	policy := network.NewTokenBucket(10, 1, time.Now)
	transport := network.NewUDPTransport(host, port, pms.Conn, nil, policy)
	raidenService := smartraiden.NewRaidenService(bcs, cfg.PrivateKey, transport, discovery, cfg)
	//startup may take long time
	raidenService.Start()
	api = &API{smartraiden.NewRaidenApi(raidenService)}
	regQuitHandler(api.api)
	return api, nil
}
func regQuitHandler(api *smartraiden.RaidenApi) {
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
	am := smartraiden.NewAccountManager(keystorePath)
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
		var data []byte
		data, err = ioutil.ReadFile(passwordfile)
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
	apihost, apiport := network.SplitHostPort(argAPIAddress)
	config.Host = listenhost
	config.Port = listenport
	config.UseConsole = false
	config.UseRpc = false
	config.ApiHost = apihost
	config.ApiPort = apiport
	config.ExternIp = pms.ExternalIp
	config.ExternPort = pms.ExternalPort
	maxUnresponsiveTime := int64(time.Minute)
	config.Protocol.NatKeepAliveTimeout = maxUnresponsiveTime / params.DefaultKeepAliveReties
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
		dataDir = path.Join(utils.GetHomePath(), ".smartraiden")
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
