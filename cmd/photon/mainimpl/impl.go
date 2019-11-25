package mainimpl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	debug2 "runtime/debug"
	"strconv"

	"github.com/SmartMeshFoundation/Photon/notify"
	"github.com/SmartMeshFoundation/Photon/rerr"

	"github.com/SmartMeshFoundation/Photon/network/netshare"

	"github.com/SmartMeshFoundation/Photon/network/mdns"

	"os/signal"
	"time"

	"strings"

	"crypto/ecdsa"

	"plugin"

	photon "github.com/SmartMeshFoundation/Photon"
	"github.com/SmartMeshFoundation/Photon/accounts"
	"github.com/SmartMeshFoundation/Photon/internal/debug"
	"github.com/SmartMeshFoundation/Photon/internal/rpanic"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/models/stormdb"
	"github.com/SmartMeshFoundation/Photon/network"
	"github.com/SmartMeshFoundation/Photon/network/helper"
	"github.com/SmartMeshFoundation/Photon/network/rpc"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/restful"
	"github.com/SmartMeshFoundation/Photon/utils"
	ethutils "github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/node"
	"gopkg.in/urfave/cli.v1"
)

func init() {
	//debug2.SetTraceback("crash")
}

var api *photon.API

//GoVersion genegate at build time
var GoVersion string

//GitCommit genegate at build time
var GitCommit string

//BuildDate genegate at build time
var BuildDate string

//Version version of this build
var Version string

//StartMain entry point of photon app
func StartMain() (*photon.API, error) {
	fmt.Printf("GoVersion=%s\nGitCommit=%s\nbuilddate=%sVersion=%s\n", GoVersion, GitCommit, BuildDate, Version)
	fmt.Printf("os.args=%q\n", os.Args)

	if len(GitCommit) != len(utils.EmptyAddress)*2 {
		if os.Getenv("ISTEST") == "" {
			return nil, rerr.ErrUnrecognized.Append("photon must build use makefile")
		}
	}
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "address",
			Usage: "The ethereum address you would like photon to use and for which a keystore file exists in your local system.",
		},
		cli.StringFlag{
			Name:  "private-key-file",
			Usage: "private key hex for run photon,only used by mesh box",
		},
		ethutils.DirectoryFlag{
			Name:  "keystore-path",
			Usage: "If you have a non-standard path for the ethereum keystore directory provide it using this argument. ",
			Value: ethutils.DirectoryString{Value: utils.DefaultKeyStoreDir()},
		},
		cli.StringFlag{
			Name:  "password-file",
			Usage: "Text file containing password for provided account",
		},
		cli.StringFlag{
			Name: "eth-rpc-endpoint",
			Usage: `"host:port" address of ethereum JSON-RPC server.\n'
	           'Also accepts a protocol prefix (ws:// or ipc channel) with optional port',`,
			Value: node.DefaultIPCEndpoint("geth"),
		},
		cli.StringFlag{
			Name:  "registry-contract-address",
			Usage: `hex encoded address of the registry contract.it's the token network contract address '`,
		},
		cli.StringFlag{
			Name:  "listen-address",
			Usage: `"host:port" for the photon service to listen on.`,
			Value: fmt.Sprintf("%s:%d", params.DefaultMainNetCfg.Host, params.DefaultMainNetCfg.Port),
		},
		cli.StringFlag{
			Name:  "api-address",
			Usage: `host:port" for the RPC server to listen on.`,
			Value: fmt.Sprintf("%s:%d", params.DefaultMainNetCfg.RestAPIHost, params.DefaultMainNetCfg.RestAPIPort),
		},
		ethutils.DirectoryFlag{
			Name:  "datadir",
			Usage: "Directory for storing photon data.",
			Value: ethutils.DirectoryString{Value: params.DefaultMainNetCfg.DataDir},
		},
		cli.BoolFlag{
			Name:  "debugcrash",
			Usage: "enable debug crash feature,only for test",
		},
		cli.StringFlag{
			Name:  "conditionquit",
			Usage: "quit at specified point for test",
			Value: "",
		},
		cli.BoolFlag{
			Name:  "debug-nonetwork",
			Usage: "disable network, for example ,when we want to settle all channels,only for test, should not be used in production",
		},
		cli.BoolFlag{
			Name:  "disable-fee",
			Usage: "disable mediation fee,default charge fee is 0.01%",
		},
		cli.BoolFlag{
			Name:  "xmpp",
			Usage: "use xmpp as transport,default is matrix, if two nodes use different transport,they cannot send message to each other",
		},
		cli.BoolFlag{
			Name:  "ignore-mediatednode-request",
			Usage: "this node doesn't work as a mediated node, only work as sender or receiver",
		},
		cli.BoolFlag{
			Name:  "enable-health-check",
			Usage: "enable health check ",
		},
		cli.BoolFlag{
			Name:  "matrix",
			Usage: "use matrix as transport,this is the default transport",
		},
		cli.IntFlag{
			Name:  "reveal-timeout",
			Usage: "channels' reveal timeout",
			Value: params.DefaultMainNetCfg.RevealTimeout,
		},
		cli.StringFlag{
			Name:  "pfs",
			Usage: "pathfinder service host,example http://transport01.smartmesh.cn:7000,default ",
		},
		cli.StringFlag{
			Name:  "pms",
			Usage: "photon-monitoring service host,example http://transport01.smartmesh.cn:8000",
		},
		cli.StringFlag{
			Name:  "pms-address",
			Usage: "account address of photon-monitoring",
		},
		cli.BoolFlag{
			Name:  "enable-fork-confirm",
			Usage: "enable fork confirm when receive events from chain,default is false,default is disabled",
		},
		cli.StringFlag{
			Name:  "http-username",
			Usage: "the username needed when call http api,only work with http-password",
		},
		cli.StringFlag{
			Name:  "http-password",
			Usage: "the password needed when call http api,only work with http-username",
		},
		cli.StringFlag{
			Name:  "http-auth-file",
			Usage: "path of http auth file",
		},
		cli.StringFlag{
			Name:  "debug-mdns-interval",
			Usage: "for test only",
			Value: "1s",
		},
		cli.StringFlag{
			Name:  "debug-mdns-keepalive", //mdns多久不响应就认为下线
			Usage: "for test only",
			Value: "20s",
		},
		cli.StringFlag{
			Name:  "debug-mdns-servicetag",
			Usage: "for test only",
			Value: mdns.ServiceTag,
		},
		cli.BoolFlag{
			Name:  "debug-udp-only",
			Usage: "for test only",
		},
		cli.BoolFlag{
			Name:  "mobile",
			Usage: "run photon in mobile mode,only used by mobile",
		},
		cli.StringFlag{
			Name:  "mobile-private-key-hex",
			Usage: "private key hex for run photon,only used by mobile",
		},
	}
	app.Flags = append(app.Flags, debug.Flags...)
	app.Action = mainCtx
	app.Name = "photon"
	app.Version = Version
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
	log.Info(fmt.Sprintf("Welcome to photon,version %s\n", ctx.App.Version))
	log.Info(fmt.Sprintf("os.args=%q", os.Args))
	log.Info(fmt.Sprintf("GoVersion=%s\nGitCommit=%s\nbuilddate=%sVersion=%s\n", GoVersion, GitCommit, BuildDate, Version))
	//photon是否已经创建成功,成功以后,dao和client的所有权也将会移动到Service中,不能自己close了
	//否则会二次close,造成错误
	var photonServiceCreated bool
	// 1. load config
	dao, client, isFirstStartUp, _, err := config(ctx)
	if err != nil {
		return
	}
	defer func() {
		if client != nil && err != nil && !photonServiceCreated {
			client.Close()
		}
		if dao != nil && err != nil && !photonServiceCreated {
			dao.CloseDB()
		}
	}()

	/*
		2.初始化notifyHandler
	*/
	notifyHandler := notify.NewNotifyHandler()

	/*
		3. 构造bcs
	*/
	bcs, err := rpc.NewBlockChainService(params.Cfg.PrivateKey, params.Cfg.RegistryAddress, client, notifyHandler, dao)
	if err != nil {
		return
	}

	/*
		4. 如果是第一次启动,校验合约信息,更新配置项,并保存ContractStatus到数据库
		这里不再使用ChainID表
	*/
	if isFirstStartUp {
		contractVersion, secretRegistryAddress, punishBlockNumber, chainID, err2 := verifyContractCode(bcs, params.Cfg.ContractVersionPrefix)
		if err2 != nil {
			err = err2
			return
		}
		dao.SaveContractStatus(models.ContractStatus{
			RegistryAddress:       params.Cfg.RegistryAddress,
			SecretRegistryAddress: secretRegistryAddress,
			PunishBlockNumber:     int64(punishBlockNumber),
			ChainID:               chainID,
			ContractVersion:       contractVersion,
		})
		params.Cfg.ChainID = chainID
		params.Cfg.PunishBlockNumber = int64(punishBlockNumber)
	}
	log.Info(fmt.Sprintf("Start Photon with ENV :\n%s", utils.MarshalIndent(params.Cfg)))

	/*
		5. 构造transport
	*/
	transport, err := buildTransport(bcs, dao)
	if err != nil {
		err = rerr.ErrUnknown
		return
	}
	defer func() {
		if err != nil && !photonServiceCreated {
			transport.Stop()
		}
	}()
	service, err := photon.NewPhotonService(bcs, transport, notifyHandler, dao)
	if err != nil {
		return
	}
	//transport,dao,client 所有权已经转移到PhotonService中了,如果以后失败,就有photonService自己管理
	photonServiceCreated = true
	// 保存构建信息
	service.SetBuildInfo(GoVersion, GitCommit, BuildDate, Version)
	err = service.Start()
	if err != nil {
		log.Error(fmt.Sprintf("photon service start error %s", err))
		service.Stop()
		return
	}
	api = photon.NewPhotonAPI(service)
	regQuitHandler(api)
	if params.Cfg.IsMobile {
		if params.Cfg.RestAPIHost == "0.0.0.0" {
			log.Info("start http server for test only...")
			go restful.Start(api)
			time.Sleep(time.Millisecond * 100)
		}
	} else {
		restful.Start(api)
	}

	return nil
}

func buildTransport(bcs *rpc.BlockChainService, dao models.Dao) (transport network.Transporter, err error) {
	/*
		use ice and doesn't work as route node,means this node runs  on a mobile phone.
	*/
	cfg := params.Cfg
	switch cfg.NetworkMode {
	case params.NoNetwork:
		params.Cfg.EnableMDNS = false
		policy := network.NewTokenBucket(10, 1, time.Now)
		transport, err = network.NewUDPTransport(bcs.NodeAddress.String(), "127.0.0.1", cfg.Port, nil, policy)
		return
	case params.UDPOnly:
		policy := network.NewTokenBucket(10, 1, time.Now)
		transport, err = network.NewUDPTransport(bcs.NodeAddress.String(), cfg.Host, cfg.Port, nil, policy)
	case params.XMPPOnly:
		transport = network.NewXMPPTransport(bcs.NodeAddress.String(), cfg.XMPPServer, bcs.PrivKey, network.DeviceTypeOther, dao)
	case params.MixUDPXMPP:
		policy := network.NewTokenBucket(10, 1, time.Now)
		deviceType := network.DeviceTypeOther
		if params.Cfg.IsMobile {
			deviceType = network.DeviceTypeMobile
		}
		transport, err = network.NewMixTranspoter(bcs.NodeAddress.String(), cfg.XMPPServer, cfg.Host, cfg.Port, bcs.PrivKey, nil, policy, deviceType, dao)
	case params.MixUDPMatrix:
		log.Info(fmt.Sprintf("use mix matrix, server=%s ", params.TrustMatrixServers))
		policy := network.NewTokenBucket(10, 1, time.Now)
		deviceType := network.DeviceTypeOther
		if params.Cfg.IsMobile {
			deviceType = network.DeviceTypeMobile
		}
		transport, err = network.NewMatrixMixTransporter(bcs.NodeAddress.String(), cfg.Host, cfg.Port, bcs.PrivKey, nil, policy, deviceType, dao)
	}
	return
}
func regQuitHandler(api *photon.API) {
	go func() {
		if params.Cfg.IsMobile {
			return
		}
		defer rpanic.PanicRecover("regQuitHandler")
		quitSignal := make(chan os.Signal, 1)
		signal.Notify(quitSignal, os.Interrupt, os.Kill)
		<-quitSignal
		signal.Stop(quitSignal)
		api.Stop()
		utils.SystemExit(0)
	}()
}

func config(ctx *cli.Context) (dao models.Dao, client *helper.SafeEthClient, isFirstStartUp, hasConnectedChain bool, err error) {
	/*
		1. 加载私钥
	*/
	privateKey, meshBoxPlugin, err := getPrivateKey(ctx)
	if err != nil {
		err = rerr.ErrArgumentError.Printf("private key err %s", err.Error())
		return
	}
	myAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	/*
		2.根据私钥构造dbpath,并打开数据库
	*/
	dataDir := ctx.String("datadir")
	dbPath, err := utils.BuildPhotonDbPath(dataDir, myAddress)
	if err != nil {
		err = rerr.ErrArgumentError.Append(err.Error())
		return
	}
	err = checkDbMeta(dbPath, "boltdb")
	if err != nil {
		err = rerr.ErrArgumentError.Printf("checkDbMeta err %s", err.Error())
		return
	}
	dao, err = stormdb.OpenDb(dbPath)
	if err != nil {
		err = rerr.ErrGeneralDBError.Printf("open db error %s", err.Error())
		return
	}
	contractStatus := dao.GetContractStatus()

	/*
		3. 连接公链节点
	*/
	ethRPCEndpoint := ctx.String("eth-rpc-endpoint")
	client, err = helper.NewSafeClient(ethRPCEndpoint)
	if err != nil {
		err = fmt.Errorf("cannot connect to geth :%s err=%s", ethRPCEndpoint, err)
		err = nil
	}

	/*
		4. 结合用户传入的参数,数据库及公链节点,确认以下内容:
			a. 当前所属的环境
			b. photon主合约地址
			c. 是否第一次启动
			d. 是否无网启动
	*/
	var paramRegistryAddress common.Address
	paramRegistryAddressStr := ctx.String("registry-contract-address")
	if len(paramRegistryAddressStr) > 0 {
		paramRegistryAddress = common.HexToAddress(paramRegistryAddressStr)
	}

	dbRegistryAddress := contractStatus.RegistryAddress
	env, registryAddress, isFirstStartUp, hasConnectedChain, err := checkEnvAndGetRegistryAddress(paramRegistryAddress, dbRegistryAddress, client)
	if err != nil {
		return
	}

	/*
		5. 使用环境对应的默认值初始化params.Cfg,之后根据用户参数调整其中内容
	*/
	params.InitDefaultCfg(env)
	if ctx.IsSet("mobile") {
		params.Cfg.IsMobile = ctx.Bool("mobile")
	}

	/*
		controlConfig
	*/
	params.Cfg.DataBasePath = dbPath
	params.Cfg.DataDir = dataDir
	params.Cfg.Debug = ctx.Bool("debug")
	if ctx.Bool("debugcrash") {
		params.Cfg.DebugCrash = true
		conditionquit := ctx.String("conditionquit")
		err = json.Unmarshal([]byte(conditionquit), &params.Cfg.ConditionQuit)
		if err != nil {
			err = rerr.ErrArgumentError.Printf("conditioquit parse error %s", err)
			return
		}
		log.Info(fmt.Sprintf("condition quit=%#v", params.Cfg.ConditionQuit))
	}
	params.Cfg.LogFilePath = ctx.String("logfile")
	params.Cfg.EnableHealthCheck = ctx.Bool("enable-health-check")

	/*
		nodeConfig
	*/
	params.Cfg.PrivateKey = privateKey
	params.Cfg.MyAddress = myAddress

	/*
		chainConfig,合约相关配置,如果第一次启动,上链获取,否则从数据库中获取
	*/
	params.Cfg.EthRPCEndPoint = ethRPCEndpoint
	params.Cfg.RegistryAddress = registryAddress
	if isFirstStartUp {
		// 因为需要调用合约,这里暂不填写,使用默认值,等后续初始化BlockChainService之后填写
	} else {
		params.Cfg.ChainID = contractStatus.ChainID
		params.Cfg.PunishBlockNumber = contractStatus.PunishBlockNumber
	}
	if ctx.IsSet("enable-fork-confirm") {
		params.Cfg.EnableForkConfirm = ctx.Bool("enable-fork-confirm")
	}

	/*
		rpcConfig
	*/
	listenHost, listenPort, err := net.SplitHostPort(ctx.String("listen-address"))
	if err != nil {
		err = rerr.ErrArgumentError.Append("--listen-address err")
		return
	}
	params.Cfg.Host = listenHost
	params.Cfg.Port, err = strconv.Atoi(listenPort)
	if err != nil {
		err = rerr.ErrArgumentError.Append("--listen-address err")
		return
	}
	if ctx.Bool("debug-nonetwork") {
		params.Cfg.NetworkMode = params.NoNetwork
	} else if ctx.Bool("debug-udp-only") {
		params.Cfg.NetworkMode = params.UDPOnly
	} else if ctx.Bool("xmpp") {
		params.Cfg.NetworkMode = params.MixUDPXMPP
	} else {
		params.Cfg.NetworkMode = params.MixUDPMatrix //默认用matrix
	}
	// mdns相关:
	if ctx.IsSet("debug-mdns-interval") {
		mi := ctx.String("debug-mdns-interval")
		dur, err2 := time.ParseDuration(mi)
		if err2 != nil {
			err = rerr.ErrArgumentError.Printf("arg debug-mdns-interval err %s", err2)
			return
		}
		params.Cfg.MDNSQueryInterval = dur
	}
	if ctx.IsSet("debug-mdns-keepalive") {
		mo := ctx.String("debug-mdns-keepalive")
		dur, err2 := time.ParseDuration(mo)
		if err2 != nil {
			err = rerr.ErrArgumentError.Printf("arg debug-mdns-keepalive err %s", err2)
			return
		}
		params.Cfg.MDNSKeepalive = dur

	}
	mdns.ServiceTag = ctx.String("debug-mdns-servicetag")
	/*
		channelConfig
	*/
	if ctx.IsSet("reveal-timeout") {
		params.Cfg.RevealTimeout = ctx.Int("reveal-timeout")
		if params.Cfg.RevealTimeout <= 0 {
			err = rerr.ErrArgumentError.Append("reveal timeout should > 0")
			return
		}
	}
	if ctx.Bool("disable-fee") {
		params.Cfg.EnableMediationFee = false
	}
	params.Cfg.IgnoreMediatedNodeRequest = ctx.Bool("ignore-mediatednode-request")

	/*
		pfsConfig
	*/
	if ctx.IsSet("pfs") {
		params.Cfg.PfsHost = ctx.String("pfs")
	}

	/*
		pmsConfig
	*/
	if ctx.IsSet("pms") {
		params.Cfg.PmsHost = ctx.String("pms")
	}
	if ctx.IsSet("pms-address") {
		params.Cfg.PmsAddress = common.HexToAddress(ctx.String("pms-address"))
	}

	/*
		apiConfig
	*/
	apiHost, apiPort, err := net.SplitHostPort(ctx.String("api-address"))
	if err != nil {
		err = rerr.ErrArgumentError.Append("--api-address err")
		return
	}
	params.Cfg.RestAPIHost = apiHost
	params.Cfg.RestAPIPort, err = strconv.Atoi(apiPort)
	if err != nil {
		err = rerr.ErrArgumentError.Append("--api-address err")
		return
	}
	params.Cfg.HTTPUsername, params.Cfg.HTTPPassword, err = getHTTPAuth(ctx, meshBoxPlugin)
	if err != nil {
		err = rerr.ErrArgumentError.Append("getHTTPAuth err")
		return
	}
	return
}

/*
Meshbox环境下从插件加载,其余情况从参数获取
*/
func getHTTPAuth(ctx *cli.Context, meshBoxPlugin *plugin.Plugin) (username, password string, err error) {
	if meshBoxPlugin != nil {
		// 这里不重复加载plugin,复用直接加载私钥时创建的对象
		httpAuthGetter, err2 := meshBoxPlugin.Lookup("GetHTTPAuthFromFile")
		if err2 != nil {
			err = fmt.Errorf("plugin lockup symbol err %s", err2)
			return
		}

		return httpAuthGetter.(func(string) (string, string, error))(ctx.String("http-auth-file"))
	}
	return ctx.String("http-username"), ctx.String("http-password"), nil
}

/*
getPrivateKey: 如果是meshbox,则通过专用插件获取私钥,否则根据指定的keystore-path找相应的私钥
*/
func getPrivateKey(ctx *cli.Context) (privateKey *ecdsa.PrivateKey, meshBoxPlugin *plugin.Plugin, err error) {
	if ctx.IsSet("mobile") && ctx.IsSet("mobile-private-key-hex") && ctx.Bool("mobile") {
		// 手机直接传递私钥的二进制字符串
		privateKeyBinHex := ctx.String("mobile-private-key-hex")
		privateKeyBytes := common.FromHex(privateKeyBinHex)
		privateKey, err = crypto.ToECDSA(privateKeyBytes)
		return
	}
	if os.Getenv("IS_MESH_BOX") == "true" || os.Getenv("IS_MESH_BOX") == "TRUE" {
		// load photon_plugin.so
		var privateKeyGetter plugin.Symbol
		var privateKeyBytes []byte
		meshBoxPlugin, err = plugin.Open("photon_plugin.so")
		if err != nil {
			err = fmt.Errorf("plugin open photo_plugin.so err %s", err)
			return
		}
		privateKeyGetter, err = meshBoxPlugin.Lookup("GetPrivateKeyForMeshBox")
		if err != nil {
			err = fmt.Errorf("plugin lockup symbol err %s", err)
			return
		}

		privateKeyBytes, err = privateKeyGetter.(func(string) ([]byte, error))(ctx.String("private-key-file"))
		if err != nil {
			err = fmt.Errorf("privateKeyGetter fail err %s", err)
			return
		}
		privateKey, err = crypto.ToECDSA(privateKeyBytes)
		return
	}
	var keyBin []byte
	address := common.HexToAddress(ctx.String("address"))
	address, keyBin, err = accounts.PromptAccount(address, ctx.String("keystore-path"), ctx.String("password-file"))
	if err != nil {
		return
	}
	debug2.FreeOSMemory() //强制立即释放scrypt分配的256M内存
	privateKey, err = crypto.ToECDSA(keyBin)
	return
}

/*
checkEnvAndGetRegistryAddress 结合用户参数,数据库及公链节点的状态,来确定当前的状态
*/
func checkEnvAndGetRegistryAddress(paramRegistryAddress common.Address, dbRegistryAddress common.Address, client *helper.SafeEthClient) (env params.ENV, registryAddress common.Address, isFirstStartUp, hasConnectedChain bool, err error) {
	isFirstStartUp = dbRegistryAddress == utils.EmptyAddress
	hasConnectedChain = client.Status == netshare.Connected
	if isFirstStartUp && !hasConnectedChain {
		err = rerr.ErrFirstStartWithoutNetwork
		return
	}
	if !isFirstStartUp && paramRegistryAddress != utils.EmptyAddress && dbRegistryAddress != paramRegistryAddress {
		err = rerr.ErrArgumentError.Printf(fmt.Sprintf("db mismatch, db's registry=%s,now registry=%s",
			dbRegistryAddress.String(), paramRegistryAddress.String()))
		return
	}
	if isFirstStartUp {
		if paramRegistryAddress == utils.EmptyAddress {
			registryAddress, err = getDefaultRegistryByEthClient(client)
			if err != nil {
				return
			}
			log.Info(fmt.Sprintf("start with TokenNetworkAddress default : %s", registryAddress.String()))
		} else {
			registryAddress = paramRegistryAddress
			log.Info(fmt.Sprintf("start with TokenNetworkAddress in param : %s", registryAddress.String()))
		}
	} else {
		registryAddress = dbRegistryAddress
		log.Info(fmt.Sprintf("start with TokenNetworkAddress in db : %s", registryAddress.String()))
	}
	if registryAddress == params.DefaultMainNetCfg.RegistryAddress {
		env = params.DefaultMainNetCfg.Env
	} else if registryAddress == params.DefaultTestNetCfg.RegistryAddress {
		env = params.DefaultTestNetCfg.Env
	} else {
		env = params.DefaultDevCfg.Env
	}
	return
}

func getDefaultRegistryByEthClient(client *helper.SafeEthClient) (registryAddress common.Address, err error) {
	var genesisBlockHash common.Hash
	genesisBlockHash, err = client.GenesisBlockHash(context.Background())
	if err != nil {
		err = rerr.ErrUnkownSpectrumRPCError.Append(err.Error())
		return
	}
	if genesisBlockHash == params.DefaultMainNetCfg.GenesisBlockHash {
		registryAddress = params.DefaultMainNetCfg.RegistryAddress
	} else if genesisBlockHash == params.DefaultTestNetCfg.GenesisBlockHash {
		registryAddress = params.DefaultTestNetCfg.RegistryAddress
	}
	return
}

/*
	校验链上的合约代码版本
*/
func verifyContractCode(bcs *rpc.BlockChainService, contractVersionPrefix string) (contractVersion string, secretRegisteryAddress common.Address, punishBlockNumber uint64, chainID *big.Int, err error) {
	log.Info(fmt.Sprintf("registry address=%s", bcs.GetRegistryAddress().String()))
	contractVersion, err = bcs.RegistryProxy.GetContractVersion()
	if err != nil {
		return
	}
	if !strings.HasPrefix(contractVersion, contractVersionPrefix) {
		err = rerr.ErrArgumentError.Printf("contract version on chain %s is incompatible with this photon version", contractVersion)
	}
	ch, err := bcs.RegistryProxy.GetContract()
	if err != nil {
		return
	}
	secretRegisteryAddress, err = ch.SecretRegistry(nil)
	if err != nil {
		err = rerr.ErrUnkownSpectrumRPCError.Printf("get SecretRegistry address err %s", err)
		return
	}
	punishBlockNumber, err = ch.PunishBlockNumber(nil)
	if err != nil {
		err = rerr.ErrUnkownSpectrumRPCError.Printf("get punish block number err %s", err)
	}
	chainID, err = ch.ChainId(nil)
	if err != nil {
		err = rerr.ErrUnkownSpectrumRPCError.Printf("get chain ID from register contract err %s", err)
	}
	return
}
func checkDbMeta(dbPath, dbType string) (err error) {
	//make sure db type not change since first start .
	dbInfo := fmt.Sprintf("%s.%s", dbPath, "info")
	if !common.FileExist(dbInfo) {
		err = ioutil.WriteFile(dbInfo, []byte(dbType), os.ModePerm)
		if err != nil {
			return
		}
	} else {
		var info []byte
		//#nosec#
		info, err = ioutil.ReadFile(dbInfo)
		if err != nil {
			return
		}
		if string(info) != dbType {
			err = errors.New("doesn't support switch db type right now")
			return
		}
	}
	return nil
}
