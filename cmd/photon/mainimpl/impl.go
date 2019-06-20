package mainimpl

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	debug2 "runtime/debug"

	"github.com/SmartMeshFoundation/Photon/rerr"

	"github.com/SmartMeshFoundation/Photon/network/netshare"

	"github.com/SmartMeshFoundation/Photon/network/mdns"

	"encoding/hex"

	"path"

	"path/filepath"

	"encoding/json"
	"os/signal"
	"time"

	"net"
	"strconv"

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
	"github.com/SmartMeshFoundation/Photon/notify"
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
			Usage: `hex encoded address of the registry contract.it's the token network contract address '`,
		},
		cli.StringFlag{
			Name:  "listen-address",
			Usage: `"host:port" for the photon service to listen on.`,
			Value: fmt.Sprintf("0.0.0.0:%d", params.InitialPort),
		},
		cli.StringFlag{
			Name:  "api-address",
			Usage: `host:port" for the RPC server to listen on.`,
			Value: "127.0.0.1:5001",
		},
		ethutils.DirectoryFlag{
			Name:  "datadir",
			Usage: "Directory for storing photon data.",
			Value: ethutils.DirectoryString{Value: params.DefaultDataDir()},
		},
		cli.StringFlag{
			Name:  "password-file",
			Usage: "Text file containing password for provided account",
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
		cli.StringFlag{
			Name:  "matrix-server",
			Usage: "use another matrix server,only domainname ,for example: transport01.smartmesh.cn",
			Value: "",
		},
		cli.BoolFlag{
			Name:  "matrix",
			Usage: "use matrix as transport,this is the default transport",
		},
		cli.IntFlag{
			Name:  "reveal-timeout",
			Usage: "channels' reveal timeout",
			Value: params.DefaultRevealTimeout,
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
			Name:  "db",
			Usage: "use --db=gkv when need photon run with gkvdb,default db is boltdb,photon doesn't support change db type once db is created.",
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
	var isFirstStartUp, hasConnectedChain bool
	//photon是否已经创建成功,成功以后,dao和client的所有权也将会移动到Service中,不能自己close了
	//否则会二次close,造成错误
	var photonServiceCreated bool
	// load config
	cfg, err := config(ctx)
	if err != nil {
		return
	}
	// connect to blockchain
	client, err := helper.NewSafeClient(cfg.EthRPCEndPoint)
	if err != nil {
		err = fmt.Errorf("cannot connect to geth :%s err=%s", cfg.EthRPCEndPoint, err)
		err = nil
	}
	defer func() {
		if client != nil && err != nil && !photonServiceCreated {
			client.Close()
		}
	}()
	// open db
	var dao models.Dao
	err = checkDbMeta(cfg.DataBasePath, "boltdb")
	if err != nil {
		err = rerr.ErrArgumentError.Printf("checkDbMeta err %s", err.Error())
		return
	}
	dao, err = stormdb.OpenDb(cfg.DataBasePath)
	if err != nil {
		err = rerr.ErrGeneralDBError.Printf("open db error %s", err.Error())
		return
	}
	defer func() {
		if err != nil && !photonServiceCreated {
			dao.CloseDB()
		}
	}()
	cfg.RegistryAddress, isFirstStartUp, hasConnectedChain, err = getRegistryAddress(cfg, dao, client)
	if err != nil {
		return
	}
	//没有pfs一样可以启动,只不过在收费模式下,交易会失败而已.
	if cfg.PfsHost == "" {
		cfg.PfsHost, err = getDefaultPFSByTokenNetworkAddress(cfg.RegistryAddress, cfg.NetworkMode == params.MixUDPMatrix)
		if err != nil {
			log.Error(fmt.Sprintf("getDefaultPFSByTokenNetworkAddress err %s", err))
			err = nil
		}
	}
	log.Info(fmt.Sprintf("pfs server=%s", cfg.PfsHost))
	// get ChainID
	if isFirstStartUp {
		if !hasConnectedChain {
			err = rerr.ErrFirstStartWithoutNetwork
			return
		}
		params.ChainID, err = client.NetworkID(context.Background())
		if err != nil {
			err = rerr.ErrUnkownSpectrumRPCError.Append(err.Error())
			return
		}
		dao.SaveChainID(params.ChainID.Int64())
	} else {
		params.ChainID = big.NewInt(dao.GetChainID())
	}
	//  init notify handler
	notifyHandler := notify.NewNotifyHandler()
	// init blockchain module
	bcs, err := rpc.NewBlockChainService(cfg.PrivateKey, cfg.RegistryAddress, client, notifyHandler, dao)
	if err != nil {
		return
	}
	if isFirstStartUp {
		var contractVersion string
		var secretRegisteryAddress common.Address
		var punishBlockNumber uint64
		var chainID *big.Int
		contractVersion, secretRegisteryAddress, punishBlockNumber, chainID, err = verifyContractCode(bcs)
		if err != nil {
			return
		}
		dao.SaveContractStatus(models.ContractStatus{
			RegistryAddress:       cfg.RegistryAddress,
			SecretRegistryAddress: secretRegisteryAddress,
			PunishBlockNumber:     int64(punishBlockNumber),
			ChainID:               chainID,
			ContractVersion:       contractVersion,
		})
	}
	//主网networkID是主网,无论在以太坊还是spectrum,都是如此
	params.IsMainNet = params.ChainID.Cmp(big.NewInt(params.MainNetNetworkID)) == 0
	if params.IsMainNet {
		cfg.SettleTimeout = params.MainNetChannelSettleTimeoutMin
	} else {
		cfg.SettleTimeout = params.TestNetChannelSettleTimeoutMin
	}
	/*
		由于数据库设计历史原因,chainID是单独保存的,为了保持兼容,暂时不做修改
	*/
	cs := dao.GetContractStatus()
	if cs.ChainID.Cmp(params.ChainID) != 0 {
		panic(fmt.Sprintf("chainid not equal ,there must be error, db status=%s,params=%s", utils.StringInterface(cs, 3), params.ChainID))
	}
	params.PunishBlockNumber = cs.PunishBlockNumber
	log.Info(fmt.Sprintf("punish block number=%d", params.PunishBlockNumber))
	transport, err := buildTransport(cfg, bcs, dao)
	if err != nil {
		err = rerr.ErrUnknown
		return
	}
	defer func() {
		if err != nil && !photonServiceCreated {
			transport.Stop()
		}
	}()
	service, err := photon.NewPhotonService(bcs, cfg.PrivateKey, transport, cfg, notifyHandler, dao)
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
func buildTransport(cfg *params.Config, bcs *rpc.BlockChainService, dao models.Dao) (transport network.Transporter, err error) {
	/*
		use ice and doesn't work as route node,means this node runs  on a mobile phone.
	*/
	switch cfg.NetworkMode {
	case params.NoNetwork:
		params.EnableMDNS = false
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
		if params.MobileMode {
			deviceType = network.DeviceTypeMobile
		}
		transport, err = network.NewMixTranspoter(bcs.NodeAddress.String(), cfg.XMPPServer, cfg.Host, cfg.Port, bcs.PrivKey, nil, policy, deviceType, dao)
	case params.MixUDPMatrix:
		log.Info(fmt.Sprintf("use mix matrix, server=%s ", params.MatrixServerConfig))
		policy := network.NewTokenBucket(10, 1, time.Now)
		deviceType := network.DeviceTypeOther
		if params.MobileMode {
			deviceType = network.DeviceTypeMobile
		}
		transport, err = network.NewMatrixMixTransporter(bcs.NodeAddress.String(), cfg.Host, cfg.Port, bcs.PrivKey, nil, policy, deviceType, dao)
	}
	return
}
func regQuitHandler(api *photon.API) {
	go func() {
		if params.MobileMode {
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
func config(ctx *cli.Context) (config *params.Config, err error) {
	config = &params.DefaultConfig
	config.EthRPCEndPoint = ctx.String("eth-rpc-endpoint")

	listenhost, listenport, err := net.SplitHostPort(ctx.String("listen-address"))
	if err != nil {
		err = rerr.ErrArgumentError.Append("--listen-address err")
		return
	}
	apihost, apiport, err := net.SplitHostPort(ctx.String("api-address"))
	if err != nil {
		err = rerr.ErrArgumentError.Append("--api-address err")
		return
	}
	config.Host = listenhost
	config.Port, err = strconv.Atoi(listenport)
	if err != nil {
		err = rerr.ErrArgumentError.Append("--listen-address err")
		return
	}
	config.UseConsole = ctx.Bool("console")
	config.APIHost = apihost
	config.APIPort, err = strconv.Atoi(apiport)
	if err != nil {
		err = rerr.ErrArgumentError.Append("--api-address err")
		return
	}
	config.PrivateKey, err = getPrivateKey(ctx)
	if err != nil {
		err = rerr.ErrArgumentError.Printf("private key err %s", err.Error())
		return
	}
	//log.Trace(fmt.Sprintf("privatekey=%s", hex.EncodeToString(crypto.FromECDSA(config.PrivateKey))))
	config.MyAddress = crypto.PubkeyToAddress(config.PrivateKey.PublicKey)
	log.Info(fmt.Sprintf("Start with account %s", config.MyAddress.String()))
	registAddrStr := ctx.String("registry-contract-address")
	if len(registAddrStr) > 0 {
		config.RegistryAddress = common.HexToAddress(registAddrStr)
	}
	dataDir := ctx.String("datadir")
	if len(dataDir) == 0 {
		dataDir = path.Join(utils.GetHomePath(), ".photon")
	}
	config.DataDir = dataDir
	if !utils.Exists(config.DataDir) {
		err = os.MkdirAll(config.DataDir, os.ModePerm)
		if err != nil {
			err = rerr.ErrArgumentError.Printf("datadir:%s doesn't exist and cannot create %v", config.DataDir, err)
			return
		}
	}
	userDbPath := hex.EncodeToString(config.MyAddress[:])
	userDbPath = userDbPath[:8]
	userDbPath = filepath.Join(config.DataDir, userDbPath)
	if !utils.Exists(userDbPath) {
		err = os.MkdirAll(userDbPath, os.ModePerm)
		if err != nil {
			err = rerr.ErrArgumentError.Printf("datadir:%s doesn't exist and cannot create %v", config.DataDir, err)
			return
		}
	}
	databasePath := filepath.Join(userDbPath, "log.db")
	config.Debug = ctx.Bool("debug")
	config.DataBasePath = databasePath
	if ctx.Bool("debugcrash") {
		config.DebugCrash = true
		conditionquit := ctx.String("conditionquit")
		err = json.Unmarshal([]byte(conditionquit), &config.ConditionQuit)
		if err != nil {
			err = rerr.ErrArgumentError.Printf("conditioquit parse error %s", err)
			return
		}
		log.Info(fmt.Sprintf("condition quit=%#v", config.ConditionQuit))
	}
	config.IgnoreMediatedNodeRequest = ctx.Bool("ignore-mediatednode-request")
	if ctx.Bool("debug-nonetwork") {
		config.NetworkMode = params.NoNetwork
	} else if ctx.Bool("debug-udp-only") {
		config.NetworkMode = params.UDPOnly
	} else if ctx.Bool("xmpp") {
		config.NetworkMode = params.MixUDPXMPP
	} else {
		config.NetworkMode = params.MixUDPMatrix //默认用matrix
	}
	config.EnableMediationFee = true
	if ctx.Bool("disable-fee") {
		config.EnableMediationFee = false
	}
	if ctx.Bool("enable-health-check") {
		config.EnableHealthCheck = true
	}
	config.XMPPServer = ctx.String("xmpp-server")
	if len(ctx.String("matrix-server")) > 0 {
		s := ctx.String("matrix-server")
		s = strings.TrimSpace(s)
		log.Info(fmt.Sprintf("use matrix server %s", s))
		for k := range params.MatrixServerConfig {
			delete(params.MatrixServerConfig, k)
		}
		params.MatrixServerConfig[s] = fmt.Sprintf("http://%s:8008", s)
	}

	if ctx.IsSet("reveal-timeout") {
		config.RevealTimeout = ctx.Int("reveal-timeout")
		if config.RevealTimeout <= 0 {
			log.Warn("reveal timeout should > 0")
		}
	}
	config.PfsHost = ctx.String("pfs")

	if ctx.Bool("enable-fork-confirm") {
		log.Info("fork-confirm enable...")
		params.EnableForkConfirm = true
	}
	if ctx.IsSet("http-username") && ctx.IsSet("http-password") {
		config.HTTPUsername = ctx.String("http-username")
		config.HTTPPassword = ctx.String("http-password")
	}
	mi := ctx.String("debug-mdns-interval")
	dur, err := time.ParseDuration(mi)
	if err != nil {
		err = rerr.ErrArgumentError.Printf("arg debug-mdns-interval err %s", err)
		return
	}
	params.DefaultMDNSQueryInterval = dur
	log.Info(fmt.Sprintf("mdns query interval=%s", params.DefaultMDNSQueryInterval))
	mo := ctx.String("debug-mdns-keepalive")
	dur, err = time.ParseDuration(mo)
	if err != nil {
		err = rerr.ErrArgumentError.Printf("arg debug-mdns-keepalive err %s", err)
		return
	}
	params.DefaultMDNSKeepalive = dur
	mdns.ServiceTag = ctx.String("debug-mdns-servicetag")
	config.PmsHost = ctx.String("pms")
	config.PmsAddress = common.HexToAddress(ctx.String("pms-address"))
	config.LogFilePath = ctx.String("logfile")
	return
}

/*
getPrivateKey: 如果是meshbox,则通过专用插件获取私钥,否则根据指定的keystore-path找相应的私钥
*/
func getPrivateKey(ctx *cli.Context) (privateKey *ecdsa.PrivateKey, err error) {
	if os.Getenv("IS_MESH_BOX") == "true" || os.Getenv("IS_MESH_BOX") == "TRUE" {
		// load photon_plugin.so
		var plug *plugin.Plugin
		var privateKeyGetter plugin.Symbol
		var privateKeyBytes []byte
		plug, err = plugin.Open("photon_plugin.so")
		if err != nil {
			err = fmt.Errorf("plugin open photo_plugin.so err %s", err)
			return
		}
		privateKeyGetter, err = plug.Lookup("GetPrivateKeyForMeshBox")
		if err != nil {
			err = fmt.Errorf("plugin lockup symbol err %s", err)
			return
		}

		privateKeyBytes, err = privateKeyGetter.(func() ([]byte, error))()
		if err != nil {
			err = fmt.Errorf("privateKeyGetter fail err %s", err)
			return
		}
		return crypto.ToECDSA(privateKeyBytes)
	}
	var keyBin []byte
	address := common.HexToAddress(ctx.String("address"))
	address, keyBin, err = accounts.PromptAccount(address, ctx.String("keystore-path"), ctx.String("password-file"))
	if err != nil {
		return
	}
	debug2.FreeOSMemory() //强制立即释放scrypt分配的256M内存
	return crypto.ToECDSA(keyBin)
}

/*
getRegistryAddress:
系统第一次初始化的时候必须保证有网,无网也没有历史数据,则没有任何意义.
是否是第一次启动的判断标准就是是否有历史数据库.
todo:存在问题,有可能第一次启动的时候连接到了一条无效的公链.
*/
func getRegistryAddress(config *params.Config, dao models.Dao, client *helper.SafeEthClient) (registryAddress common.Address, isFirstStartUp, hasConnectedChain bool, err error) {
	log.Info(fmt.Sprintf("contract status=%s", utils.StringInterface(dao.GetContractStatus(), 5)))
	dbRegistryAddress := dao.GetContractStatus().RegistryAddress
	isFirstStartUp = dbRegistryAddress == utils.EmptyAddress
	hasConnectedChain = client.Status == netshare.Connected
	if isFirstStartUp && !hasConnectedChain {
		err = rerr.ErrFirstStartWithoutNetwork
		return
	}
	if !isFirstStartUp && config.RegistryAddress != utils.EmptyAddress && dbRegistryAddress != config.RegistryAddress {
		err = rerr.ErrArgumentError.Printf(fmt.Sprintf("db mismatch, db's registry=%s,now registry=%s",
			dbRegistryAddress.String(), config.RegistryAddress.String()))
		return
	}
	if isFirstStartUp {
		if config.RegistryAddress == utils.EmptyAddress {
			registryAddress, err = getDefaultRegistryByEthClient(client)
			if err != nil {
				return
			}
			log.Info(fmt.Sprintf("start with TokenNetworkAddress default : %s", registryAddress.String()))
		} else {
			registryAddress = config.RegistryAddress
			log.Info(fmt.Sprintf("start with TokenNetworkAddress in param : %s", registryAddress.String()))
		}
		//等交验完合约没问题以后再存,否则合约有问题还需要重新来过
		//dao.SaveContractStatus(registryAddress)
	} else {
		registryAddress = dbRegistryAddress
		log.Info(fmt.Sprintf("start with TokenNetworkAddress in db : %s", registryAddress.String()))
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
	registryAddress = params.GenesisBlockHashToDefaultRegistryAddress[genesisBlockHash]
	return
}
func getDefaultPFSByTokenNetworkAddress(tokenNetworkAddress common.Address, isMatrix bool) (pfs string, err error) {
	if isMatrix {
		var ok bool
		pfs, ok = params.DefaultMatrixContractToPFS[tokenNetworkAddress]
		if !ok {
			err = rerr.ErrArgumentError.Printf("can not find default pfs host by TokenNetworkAddress[%s]", tokenNetworkAddress.String())
			return
		}
		return
	}
	pfs, ok := params.DefaultContractToPFS[tokenNetworkAddress]
	if !ok {
		err = rerr.ErrArgumentError.Printf("can not find default pfs host by TokenNetworkAddress[%s]", tokenNetworkAddress.String())
		return
	}
	return
}

/*
	校验链上的合约代码版本
*/
func verifyContractCode(bcs *rpc.BlockChainService) (contractVersion string, secretRegisteryAddress common.Address, punishBlockNumber uint64, chainID *big.Int, err error) {
	log.Info(fmt.Sprintf("registry address=%s", bcs.GetRegistryAddress().String()))
	contractVersion, err = bcs.RegistryProxy.GetContractVersion()
	if err != nil {
		return
	}
	if !strings.HasPrefix(contractVersion, params.ContractVersionPrefix) {
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
