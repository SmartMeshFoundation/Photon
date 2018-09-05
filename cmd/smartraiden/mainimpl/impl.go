package mainimpl

import (
	"fmt"
	"os"

	"encoding/hex"

	"path"

	"path/filepath"

	"encoding/json"
	"os/signal"
	"time"

	"net"
	"strconv"

	"strings"

	"github.com/SmartMeshFoundation/SmartRaiden"
	"github.com/SmartMeshFoundation/SmartRaiden/accounts"
	"github.com/SmartMeshFoundation/SmartRaiden/internal/debug"
	"github.com/SmartMeshFoundation/SmartRaiden/internal/rpanic"
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
	//debug2.SetTraceback("crash")
}

var api *smartraiden.RaidenAPI

//StartMain entry point of raiden app
func StartMain() (*smartraiden.RaidenAPI, error) {
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
			Value: ethutils.DirectoryString{Value: params.DefaultKeyStoreDir()},
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
			Value: params.SpectrumTestNetRegistryAddress.String(),
		},
		cli.StringFlag{
			Name:  "listen-address",
			Usage: `"host:port" for the raiden service to listen on.`,
			Value: fmt.Sprintf("0.0.0.0:%d", params.InitialPort),
		},
		cli.StringFlag{
			Name:  "api-address",
			Usage: `host:port" for the RPC server to listen on.`,
			Value: "127.0.0.1:5001",
		},
		ethutils.DirectoryFlag{
			Name:  "datadir",
			Usage: "Directory for storing raiden data.",
			Value: ethutils.DirectoryString{Value: params.DefaultDataDir()},
		},
		cli.StringFlag{
			Name:  "password-file",
			Usage: "Text file containing password for provided account",
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
		cli.BoolFlag{
			Name:  "nonetwork",
			Usage: "disable network, for example ,when we want to settle all channels",
		},
		cli.BoolFlag{
			Name:  "fee",
			Usage: "enable mediation fee",
		},
		cli.StringFlag{
			Name:  "xmpp-server",
			Usage: "use another xmpp server ",
			Value: params.DefaultXMPPServer,
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
	app.Action = mainCtx
	app.Name = "smartraiden"
	app.Version = "0.8"
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
	err := app.Run(os.Args)
	return api, err
}

func mainCtx(ctx *cli.Context) (err error) {
	log.Info(fmt.Sprintf("Welcom to smartraiden,version %s\n", ctx.App.Version))
	log.Info(fmt.Sprintf("os.args=%q", os.Args))
	cfg, err := config(ctx)
	if err != nil {
		return
	}
	//log.Debug(fmt.Sprintf("Config:%s", utils.StringInterface(cfg, 2)))
	ethEndpoint := ctx.String("eth-rpc-endpoint")
	// 禁止使用http协议启动,smartraiden,因为会出现以有网状态启动,但始终无法获取到链上的事件的情况,这会带来风险
	if strings.HasPrefix(ethEndpoint, "http") {
		err = fmt.Errorf("cannot connect to geth :%s err= does not support http protocol,please use websocket instead", ethEndpoint)
		return
	}
	client, err := helper.NewSafeClient(ethEndpoint)
	if err != nil {
		err = fmt.Errorf("cannot connect to geth :%s err=%s", ethEndpoint, err)
		return
	}
	bcs := rpc.NewBlockChainService(cfg.PrivateKey, cfg.RegistryAddress, client)
	transport, err := buildTransport(cfg, bcs)
	if err != nil {
		return
	}
	raidenService, err := smartraiden.NewRaidenService(bcs, cfg.PrivateKey, transport, cfg)
	if err != nil {
		transport.Stop()
		return
	}
	if cfg.EnableMediationFee {
		//do nothing.
	} else {
		raidenService.SetFeePolicy(&smartraiden.NoFeePolicy{})
	}
	err = raidenService.Start()
	if err != nil {
		raidenService.Stop()
		return
	}
	api = smartraiden.NewRaidenAPI(raidenService)
	regQuitHandler(api)
	if params.MobileMode {
		if cfg.APIHost == "0.0.0.0" {
			log.Info("start http server for test only...")
			go restful.Start(api, cfg)
			time.Sleep(time.Millisecond * 100)
		}
	} else {
		restful.Start(api, cfg)
	}

	return nil
}
func buildTransport(cfg *params.Config, bcs *rpc.BlockChainService) (transport network.Transporter, err error) {
	/*
		use ice and doesn't work as route node,means this node runs  on a mobile phone.
	*/
	if params.MobileMode {
		cfg.NetworkMode = params.MixUDPXMPP
	}
	switch cfg.NetworkMode {
	case params.NoNetwork:
		policy := network.NewTokenBucket(10, 1, time.Now)
		transport, err = network.NewUDPTransport(utils.APex2(bcs.NodeAddress), "127.0.0.1", cfg.Port, nil, policy)
		return
	case params.UDPOnly:
		policy := network.NewTokenBucket(10, 1, time.Now)
		transport, err = network.NewUDPTransport(utils.APex2(bcs.NodeAddress), cfg.Host, cfg.Port, nil, policy)
	case params.XMPPOnly:
		transport = network.NewXMPPTransport(utils.APex2(bcs.NodeAddress), cfg.XMPPServer, bcs.PrivKey, network.DeviceTypeOther)
	case params.MixUDPXMPP:
		policy := network.NewTokenBucket(10, 1, time.Now)
		deviceType := network.DeviceTypeOther
		if params.MobileMode {
			deviceType = network.DeviceTypeMobile
		}
		transport, err = network.NewMixTranspoter(utils.APex2(bcs.NodeAddress), cfg.XMPPServer, cfg.Host, cfg.Port, bcs.PrivKey, nil, policy, deviceType)
	}
	return
}
func regQuitHandler(api *smartraiden.RaidenAPI) {
	go func() {
		defer rpanic.PanicRecover("regQuitHandler")
		quitSignal := make(chan os.Signal, 1)
		signal.Notify(quitSignal, os.Interrupt, os.Kill)
		<-quitSignal
		signal.Stop(quitSignal)
		api.Stop()
		utils.SystemExit(0)
	}()
}
func config(ctx *cli.Context) (config *params.Config, err error) {
	config = &params.DefaultConfig
	listenhost, listenport, err := net.SplitHostPort(ctx.String("listen-address"))
	if err != nil {
		return
	}
	apihost, apiport, err := net.SplitHostPort(ctx.String("api-address"))
	if err != nil {
		return
	}
	config.Host = listenhost
	config.Port, err = strconv.Atoi(listenport)
	if err != nil {
		return
	}
	config.UseConsole = ctx.Bool("console")
	config.APIHost = apihost
	config.APIPort, err = strconv.Atoi(apiport)
	if err != nil {
		return
	}
	address := common.HexToAddress(ctx.String("address"))
	address, privkeyBin, err := accounts.PromptAccount(address, ctx.String("keystore-path"), ctx.String("password-file"))
	if err != nil {
		return
	}
	config.PrivateKeyHex = hex.EncodeToString(privkeyBin)
	config.PrivateKey, err = crypto.ToECDSA(privkeyBin)
	config.MyAddress = address
	if err != nil {
		err = fmt.Errorf("privkey error: %s", err)
		return
	}
	registAddrStr := ctx.String("registry-contract-address")
	if len(registAddrStr) > 0 {
		config.RegistryAddress = common.HexToAddress(registAddrStr)
	}
	dataDir := ctx.String("datadir")
	if len(dataDir) == 0 {
		dataDir = path.Join(utils.GetHomePath(), ".smartraiden")
	}
	config.DataDir = dataDir
	if !utils.Exists(config.DataDir) {
		err = os.MkdirAll(config.DataDir, os.ModePerm)
		if err != nil {
			err = fmt.Errorf("datadir:%s doesn't exist and cannot create %v", config.DataDir, err)
			return
		}
	}
	userDbPath := hex.EncodeToString(config.MyAddress[:])
	userDbPath = userDbPath[:8]
	userDbPath = filepath.Join(config.DataDir, userDbPath)
	if !utils.Exists(userDbPath) {
		err = os.MkdirAll(userDbPath, os.ModePerm)
		if err != nil {
			err = fmt.Errorf("datadir:%s doesn't exist and cannot create %v", userDbPath, err)
			return
		}
	}
	databasePath := filepath.Join(userDbPath, "log.db")
	config.DataBasePath = databasePath
	if ctx.Bool("debugcrash") {
		config.DebugCrash = true
		conditionquit := ctx.String("conditionquit")
		err = json.Unmarshal([]byte(conditionquit), &config.ConditionQuit)
		if err != nil {
			err = fmt.Errorf("conditioquit parse error %s", err)
			return
		}
		log.Info(fmt.Sprintf("condition quit=%#v", config.ConditionQuit))
	}
	config.IgnoreMediatedNodeRequest = ctx.Bool("ignore-mediatednode-request")
	if ctx.Bool("nonetwork") {
		config.NetworkMode = params.NoNetwork
	} else {
		config.NetworkMode = params.MixUDPXMPP
	}
	if ctx.Bool("fee") {
		config.EnableMediationFee = true
	}
	if ctx.Bool("enable-health-check") {
		config.EnableHealthCheck = true
	}
	config.XMPPServer = ctx.String("xmpp-server")
	return
}
