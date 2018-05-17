package mainimpl

import (
	"fmt"
	"os"

	"io/ioutil"

	"encoding/hex"

	"path"

	"path/filepath"

	"encoding/json"
	"os/signal"
	debug2 "runtime/debug"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden"
	"github.com/SmartMeshFoundation/SmartRaiden/internal/debug"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/restful"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	ethutils "github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/node"
	"gopkg.in/urfave/cli.v1"
)

func init() {
	debug2.SetTraceback("crash")
}

func panicOnNullValue() {
	var c []int
	c[0] = 0
}

func StartMain() {
	os.Args[0] = "smartraiden"
	fmt.Printf("os.args=%q\n", os.Args)
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "address",
			Usage: "The ethereum address you would like raiden to use and for which a keystore file exists in your local system.",
		},
		ethutils.DirectoryFlag{
			Name:  "keystore-path",
			Usage: "If you have a non-standard path for the ethereum keystore directory provide it using this argument. ",
			Value: ethutils.DirectoryString{params.DefaultKeyStoreDir()},
		},
		cli.StringFlag{
			Name: "eth-rpc-endpoint",
			Usage: `"host:port" address of ethereum JSON-RPC server.\n'
	           'Also accepts a protocol prefix (ws:// or ipc channel) with optional port',`,
			Value: node.DefaultIPCEndpoint("geth"),
		},
		cli.StringFlag{
			Name:  "registry-contract-address",
			Usage: `hex encoded address of the registry contract.`,
			Value: params.RopstenRegistryAddress.String(),
		},
		cli.StringFlag{
			Name:  "discovery-contract-address",
			Usage: `hex encoded address of the discovery contract.`,
			Value: params.RopstenDiscoveryAddress.String(),
		},
		cli.StringFlag{
			Name:  "listen-address",
			Usage: `"host:port" for the raiden service to listen on.`,
			Value: fmt.Sprintf("0.0.0.0:%d", params.InitialPort),
		},
		cli.StringFlag{
			Name: "rpccorsdomain",
			Usage: `Comma separated list of domains to accept cross origin requests.
				(localhost enabled by default)`,
			Value: "http://localhost:* /*",
		},
		cli.IntFlag{Name: "max-unresponsive-time",
			Usage: `Max time in seconds for which an address can send no packets and
	               still be considered healthy.`,
			Value: 120,
		},
		cli.IntFlag{Name: "send-ping-time",
			Usage: `Time in seconds after which if we have received no message from a
	               node we have a connection with, we are going to send a PING message`,
			Value: 60,
		},
		cli.BoolTFlag{Name: "rpc",
			Usage: `Start with or without the RPC server. Default is to start
	               the RPC server`,
		},
		cli.StringFlag{
			Name:  "api-address",
			Usage: `host:port" for the RPC server to listen on.`,
			Value: "127.0.0.1:5001",
		},
		ethutils.DirectoryFlag{
			Name:  "datadir",
			Usage: "Directory for storing raiden data.",
			Value: ethutils.DirectoryString{params.DefaultDataDir()},
		},
		cli.StringFlag{
			Name:  "password-file",
			Usage: "Text file containing password for provided account",
		},
		cli.StringFlag{
			Name: "nat",
			Usage: `
				[auto|upnp|stun|none]
				Manually specify method to use for
				determining public IP / NAT traversal.
				Available methods:
				"auto" - Try UPnP, then
				STUN, fallback to none
				"upnp" - Try UPnP,
				fallback to none
				"stun" - Try STUN, fallback
				to none
				"none" - Use the local interface,only for test
				address (this will likely cause connectivity
				issues)
				"ice"- Use ice framework for nat punching
				[default: ice]`,
			Value: "ice",
		},
		cli.BoolFlag{
			Name:  "debugcrash",
			Usage: "enable debug crash feature",
		},
		cli.StringFlag{
			Name:  "conditionquit",
			Usage: "quit at specified point for test",
			Value: "",
		},
		cli.StringFlag{
			Name:  "turn-server",
			Usage: "tur server for ice",
			Value: params.DefaultTurnServer,
		},
		cli.StringFlag{
			Name:  "turn-user",
			Usage: "turn username for turn server",
			Value: params.DefaultTurnUserName,
		},
		cli.StringFlag{
			Name:  "turn-pass",
			Usage: "turn password for turn server",
			Value: params.DefaultTurnPassword,
		},
		cli.BoolFlag{
			Name:  "nonetwork",
			Usage: "disable network, for example ,when we want to settle all channels",
		},
		cli.BoolFlag{
			Name:  "fee",
			Usage: "enable mediation fee",
		},
		cli.StringFlag{
			Name:  "signal-server",
			Usage: "use another signal server ",
			Value: params.DefaultSignalServer,
		},
		cli.BoolFlag{
			Name:  "ignore-mediatednode-request",
			Usage: "this node doesn't work as a mediated node, only work as sender or receiver",
		},
		cli.BoolFlag{
			Name:  "enable-health-check",
			Usage: "enable health check ",
		},
	}
	app.Flags = append(app.Flags, debug.Flags...)
	app.Action = MainCtx
	app.Name = "smartraiden"
	app.Version = "0.2"
	app.Before = func(ctx *cli.Context) error {
		if err := debug.Setup(ctx); err != nil {
			return err
		}
		return nil
	}

	app.After = func(ctx *cli.Context) error {
		debug.Exit()
		return nil
	}
	app.Run(os.Args)
}
func MainCtx(ctx *cli.Context) error {
	var pms *network.PortMappedSocket
	var err error
	fmt.Printf("Welcom to smartraiden,version %s\n", ctx.App.Version)
	if ctx.String("nat") != "ice" {
		host, port := network.SplitHostPort(ctx.String("listen-address"))
		pms, err = network.SocketFactory(host, port, ctx.String("nat"))
		if err != nil {
			log.Crit(fmt.Sprintf("SocketFactory err=%s", err))
			return err
		}
		log.Trace(fmt.Sprintf("pms=%s", utils.StringInterface1(pms)))
	} else {
		host, port := network.SplitHostPort(ctx.String("listen-address"))
		pms = &network.PortMappedSocket{
			Ip:   host,
			Port: port,
		}
	}
	if err != nil {
		log.Error(fmt.Sprintf("start server on %s error:%s", ctx.String("listen-address"), err))
		utils.SystemExit(1)
	}
	cfg := config(ctx, pms)
	//log.Debug(fmt.Sprintf("Config:%s", utils.StringInterface(cfg, 2)))
	ethEndpoint := ctx.String("eth-rpc-endpoint")
	client, err := helper.NewSafeClient(ethEndpoint)
	if err != nil {
		log.Error(fmt.Sprintf("cannot connect to geth :%s err=%s", ethEndpoint, err))
		utils.SystemExit(1)
	}
	bcs := rpc.NewBlockChainService(cfg.PrivateKey, cfg.RegistryAddress, client)
	log.Trace(fmt.Sprintf("bcs=%#v", bcs))
	transport, discovery := buildTransportAndDiscovery(cfg, pms, bcs)
	raidenService := smartraiden.NewRaidenService(bcs, cfg.PrivateKey, transport, discovery, cfg)
	if cfg.EnableMediationFee {
		//do nothing.
	} else {
		raidenService.SetFeePolicy(&smartraiden.NoFeePolicy{})
	}
	raidenService.Start()
	api := smartraiden.NewRaidenApi(raidenService)
	regQuitHandler(api)
	restful.Start(api, cfg)
	return nil
}
func buildTransportAndDiscovery(cfg *params.Config, pms *network.PortMappedSocket, bcs *rpc.BlockChainService) (transport network.Transporter, discovery network.DiscoveryInterface) {
	var err error
	/*
		use ice and doesn't work as route node,means this node runs  on a mobile phone.
	*/
	if cfg.NetworkMode == params.ICEOnly && cfg.IgnoreMediatedNodeRequest {
		cfg.NetworkMode = params.MixUDPICE
	}
	switch cfg.NetworkMode {
	case params.NoNetwork:
		discovery = network.NewDiscovery()
		policy := network.NewTokenBucket(10, 1, time.Now)
		transport = network.NewDummyTransport(pms.Ip, pms.Port, nil, policy)
		return
	case params.UDPOnly:
		discovery = network.NewContractDiscovery(bcs.NodeAddress, cfg.DiscoveryAddress, bcs.Client, bcs.Auth)
		policy := network.NewTokenBucket(10, 1, time.Now)
		transport = network.NewUDPTransport(pms.Ip, pms.Port, pms.Conn, nil, policy)
	case params.ICEOnly:
		network.InitIceTransporter(cfg.Ice.TurnServer, cfg.Ice.TurnUser, cfg.Ice.TurnPassword, cfg.Ice.SignalServer)
		transport, err = network.NewIceTransporter(bcs.PrivKey, utils.APex2(bcs.NodeAddress))
		if err != nil {
			panic(err)
		}
		discovery = network.NewIceHelperDiscovery()
	case params.MixUDPICE:
		network.InitIceTransporter(cfg.Ice.TurnServer, cfg.Ice.TurnUser, cfg.Ice.TurnPassword, cfg.Ice.SignalServer)
		policy := network.NewTokenBucket(10, 1, time.Now)
		transport, discovery = network.NewMixTranspoter(bcs.PrivKey, utils.APex2(bcs.NodeAddress), pms.Ip, pms.Port, pms.Conn, nil, policy)
	}
	return
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
	if len(passwordfile) > 0 {
		data, err := ioutil.ReadFile(passwordfile)
		if err != nil {
			//pass, err := utils.PasswordDecrypt(passwordfile)
			//if err != nil {
			//	panic("decrypt pass err " + err.Error())
			//}
			//data = []byte(pass)
			data = []byte(passwordfile)
		}
		password := string(data)
		log.Trace(fmt.Sprintf("password is %s", password))
		keybin, err = am.GetPrivateKey(addr, password)
		if err != nil {
			log.Error(fmt.Sprintf("Incorrect password for %s in file. Aborting ... %s", addr.String(), err))
			utils.SystemExit(1)
		}
	} else {
		//for i := 0; i < 3; i++ {
		//	//retries three times
		//	password = getpass.Prompt("Enter the password to unlock:")
		//	keybin, err = am.GetPrivateKey(addr, password)
		//	if err != nil && i == 3 {
		//		log.Error(fmt.Sprintf("Exhausted passphrase unlock attempts for %s. Aborting ...", addr))
		//		utils.SystemExit(1)
		//	}
		//	if err != nil {
		//		log.Error(fmt.Sprintf("password incorrect\n Please try again or kill the process to quit.\nUsually Ctrl-c."))
		//		continue
		//	}
		//	break
		//}
		panic("must specified password")
	}
	return
}
func config(ctx *cli.Context, pms *network.PortMappedSocket) *params.Config {
	var err error
	config := params.DefaultConfig
	listenhost, listenport := network.SplitHostPort(ctx.String("listen-address"))
	apihost, apiport := network.SplitHostPort(ctx.String("api-address"))
	config.Host = listenhost
	config.Port = listenport
	config.UseConsole = ctx.Bool("console")
	config.UseRpc = ctx.Bool("rpc")
	config.ApiHost = apihost
	config.ApiPort = apiport
	config.ExternIp = pms.ExternalIp
	config.ExternPort = pms.ExternalPort
	maxUnresponsiveTime := ctx.Int64("max-unresponsive-time")
	config.Protocol.NatKeepAliveTimeout = maxUnresponsiveTime / params.DefaultKeepAliveReties
	address := common.HexToAddress(ctx.String("address"))
	address, privkeyBin := promptAccount(address, ctx.String("keystore-path"), ctx.String("password-file"))
	config.PrivateKeyHex = hex.EncodeToString(privkeyBin)
	config.PrivateKey, err = crypto.ToECDSA(privkeyBin)
	config.MyAddress = address
	if err != nil {
		log.Error("privkey error:", err)
		utils.SystemExit(1)
	}
	registAddrStr := ctx.String("registry-contract-address")
	if len(registAddrStr) > 0 {
		config.RegistryAddress = common.HexToAddress(registAddrStr)
	}
	discoverAddr := ctx.String("discovery-contract-address")
	if len(discoverAddr) > 0 {
		config.DiscoveryAddress = common.HexToAddress(discoverAddr)
	}
	dataDir := ctx.String("datadir")
	if len(dataDir) == 0 {
		dataDir = path.Join(utils.GetHomePath(), ".smartraiden")
	}
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
	if !utils.Exists(userDbPath) {
		err = os.MkdirAll(userDbPath, os.ModePerm)
		if err != nil {
			log.Error(fmt.Sprintf("Datadir:%s doesn't exist and cannot create %v", userDbPath, err))
			utils.SystemExit(1)
		}
	}
	databasePath := filepath.Join(userDbPath, "log.db")
	config.DataBasePath = databasePath
	if ctx.Bool("debugcrash") {
		config.DebugCrash = true
		conditionquit := ctx.String("conditionquit")
		err := json.Unmarshal([]byte(conditionquit), &config.ConditionQuit)
		if err != nil {
			panic(fmt.Sprintf("conditioquit parse error %s", err))
		}
		log.Info(fmt.Sprintf("condition quit=%#v", config.ConditionQuit))
	}
	config.Ice.StunServer = ctx.String("turn-server")
	config.Ice.TurnServer = ctx.String("turn-server")
	config.Ice.TurnUser = ctx.String("turn-user")
	config.Ice.TurnPassword = ctx.String("turn-pass")
	config.IgnoreMediatedNodeRequest = ctx.Bool("ignore-mediatednode-request")
	if ctx.String("nat") == "ice" {
		config.NetworkMode = params.ICEOnly
	} else if ctx.Bool("nonetwork") {
		config.NetworkMode = params.NoNetwork
	} else {
		config.NetworkMode = params.UDPOnly
	}
	if ctx.Bool("fee") {
		config.EnableMediationFee = true
	}
	config.Ice.SignalServer = ctx.String("signal-server")
	log.Trace(fmt.Sprintf("signal server=%s", config.Ice.SignalServer))
	if ctx.Bool("enable-health-check") {
		config.EnableHealthCheck = true
	}
	return &config
}
func init() {
	//many race condtions don't resolve
	setNativeThreadNumber()
}
func setNativeThreadNumber() {
	//runtime.GOMAXPROCS(1)
}
