package mainimpl

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

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
	"github.com/SmartMeshFoundation/Photon/network/netshare"
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
			return nil, fmt.Errorf("photon must build use makefile")
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
			Usage: "use xmpp as transport,default is xmpp, if two nodes use different transport,they cannot send message to each other",
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
			Usage: "use another matrix server",
			Value: "",
		},
		cli.BoolFlag{
			Name:  "matrix",
			Usage: "use matrix as transport,default is xmpp",
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
	// open db
	var dao models.Dao
	err = checkDbMeta(cfg.DataBasePath, "boltdb")
	if err != nil {
		return
	}
	dao, err = stormdb.OpenDb(cfg.DataBasePath)
	//}
	if err != nil {
		err = fmt.Errorf("open db error %s", err)
		client.Close()
		return
	}
	cfg.RegistryAddress, isFirstStartUp, hasConnectedChain, err = getRegistryAddress(cfg, dao, client)
	if err != nil {
		client.Close()
		dao.CloseDB()
		return
	}
	//没有pfs一样可以启动,只不过在收费模式下,交易会失败而已.
	if cfg.PfsHost == "" {
		cfg.PfsHost, err = getDefaultPFSByTokenNetworkAddress(cfg.RegistryAddress)
		if err != nil {
			log.Error(fmt.Sprintf("getDefaultPFSByTokenNetworkAddress err %s", err))
			//client.Close()
			//dao.CloseDB()
			//return
		}
	}
	// get ChainID
	if isFirstStartUp {
		if !hasConnectedChain {
			err = fmt.Errorf("first startup without ethereum rpc connection")
			dao.CloseDB()
			client.Close()
			return
		}
		params.ChainID, err = client.NetworkID(context.Background())
		if err != nil {
			dao.CloseDB()
			client.Close()
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
		dao.CloseDB()
		client.Close()
		return
	}
	if isFirstStartUp {
		var contractVersion string
		var secretRegisteryAddress common.Address
		var punishBlockNumber uint64
		var chainID *big.Int
		contractVersion, secretRegisteryAddress, punishBlockNumber, chainID, err = verifyContractCode(bcs)
		if err != nil {
			//dao.SaveContractStatus(models.ContractStatus{}) // return to first start up
			dao.CloseDB()
			client.Close()
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
	/*
		由于数据库设计历史原因,chainID是单独保存的,为了保持兼容,暂时不做修改
	*/
	cs := dao.GetContractStatus()
	if cs.ChainID.Cmp(params.ChainID) != 0 {
		panic(fmt.Sprintf("chainid not equal ,there must be error, db status=%s,params=%s", utils.StringInterface(cs, 3), params.ChainID))
	}
	params.PunishBlockNumber = cs.PunishBlockNumber
	log.Info(fmt.Sprintf("punish block number=%d", params.PunishBlockNumber))
	transport, err := buildTransport(cfg, bcs)
	if err != nil {
		dao.CloseDB()
		client.Close()
		return
	}
	service, err := photon.NewPhotonService(bcs, cfg.PrivateKey, transport, cfg, notifyHandler, dao)
	if err != nil {
		dao.CloseDB()
		client.Close()
		transport.Stop()
		return
	}
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
func buildTransport(cfg *params.Config, bcs *rpc.BlockChainService) (transport network.Transporter, err error) {
	/*
		use ice and doesn't work as route node,means this node runs  on a mobile phone.
	*/
	if params.MobileMode {
		cfg.NetworkMode = params.MixUDPXMPP
	}
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
		transport = network.NewXMPPTransport(bcs.NodeAddress.String(), cfg.XMPPServer, bcs.PrivKey, network.DeviceTypeOther)
	case params.MixUDPXMPP:
		policy := network.NewTokenBucket(10, 1, time.Now)
		deviceType := network.DeviceTypeOther
		if params.MobileMode {
			deviceType = network.DeviceTypeMobile
		}
		transport, err = network.NewMixTranspoter(bcs.NodeAddress.String(), cfg.XMPPServer, cfg.Host, cfg.Port, bcs.PrivKey, nil, policy, deviceType)
	case params.MixUDPMatrix:
		log.Trace(fmt.Sprintf("use mix matrix, server=%s ", params.MatrixServerConfig))
		policy := network.NewTokenBucket(10, 1, time.Now)
		deviceType := network.DeviceTypeOther
		if params.MobileMode {
			deviceType = network.DeviceTypeMobile
		}
		transport, err = network.NewMatrixMixTransporter(bcs.NodeAddress.String(), cfg.Host, cfg.Port, bcs.PrivKey, nil, policy, deviceType)
	}
	return
}
func regQuitHandler(api *photon.API) {
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
	config.EthRPCEndPoint = ctx.String("eth-rpc-endpoint")

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
	config.PrivateKey, err = getPrivateKey(ctx)
	if err != nil {
		err = fmt.Errorf("privkey error: %s", err)
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
	config.Debug = ctx.Bool("debug")
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
	if ctx.Bool("debug-nonetwork") {
		config.NetworkMode = params.NoNetwork
	} else if ctx.Bool("debug-udp-only") {
		config.NetworkMode = params.UDPOnly
	} else if ctx.Bool("matrix") {
		config.NetworkMode = params.MixUDPMatrix
	} else {
		config.NetworkMode = params.MixUDPXMPP //默认用xmpp做通信,matrix不太稳定
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
		log.Info(fmt.Sprintf("use matrix server %s", s))
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
		err = fmt.Errorf("arg debug-mdns-interval err %s", err)
		return
	}
	params.DefaultMDNSQueryInterval = dur
	log.Info(fmt.Sprintf("mdns query interval=%s", params.DefaultMDNSQueryInterval))
	mo := ctx.String("debug-mdns-keepalive")
	dur, err = time.ParseDuration(mo)
	if err != nil {
		err = fmt.Errorf("arg debug-mdns-keepalive err %s", err)
		return
	}
	params.DefaultMDNSKeepalive = dur
	mdns.ServiceTag = ctx.String("debug-mdns-servicetag")
	return
}

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
	return crypto.ToECDSA(keyBin)
}

func getRegistryAddress(config *params.Config, dao models.Dao, client *helper.SafeEthClient) (registryAddress common.Address, isFirstStartUp, hasConnectedChain bool, err error) {
	dbRegistryAddress := dao.GetContractStatus().RegistryAddress
	isFirstStartUp = dbRegistryAddress == utils.EmptyAddress
	hasConnectedChain = client.Status == netshare.Connected
	if isFirstStartUp && !hasConnectedChain {
		err = fmt.Errorf("first startup without ethereum rpc connection")
		return
	}
	if !isFirstStartUp && config.RegistryAddress != utils.EmptyAddress && dbRegistryAddress != config.RegistryAddress {
		err = fmt.Errorf(fmt.Sprintf("db mismatch, db's registry=%s,now registry=%s",
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
		log.Error(err.Error())
		return
	}
	// 根据创世区块hash设置主网标记
	if genesisBlockHash == params.MainNetGenesisBlockHash {
		params.IsMainNet = true
	}
	registryAddress = params.GenesisBlockHashToDefaultRegistryAddress[genesisBlockHash]
	return
}
func getDefaultPFSByTokenNetworkAddress(tokenNetworkAddress common.Address) (pfs string, err error) {
	pfs, ok := params.GenesisBlockHashToPFS[tokenNetworkAddress]
	if !ok {
		err = fmt.Errorf("can not find default pfs host by TokenNetworkAddress[%s]", tokenNetworkAddress.String())
		return
	}
	return
}

/*
	校验链上的合约代码版本
*/
func verifyContractCode(bcs *rpc.BlockChainService) (contractVersion string, secretRegisteryAddress common.Address, punishBlockNumber uint64, chainID *big.Int, err error) {
	log.Trace(fmt.Sprintf("registry address=%s", bcs.GetRegistryAddress().String()))
	contractVersion, err = bcs.RegistryProxy.GetContractVersion()
	if err != nil {
		return
	}
	if !strings.HasPrefix(contractVersion, params.ContractVersionPrefix) {
		err = fmt.Errorf("contract version on chain %s is incompatible with this photon version", contractVersion)
	}
	secretRegisteryAddress, err = bcs.RegistryProxy.GetContract().SecretRegistry(nil)
	if err != nil {
		err = fmt.Errorf("get SecretRegistry address err %s", err)
		return
	}
	punishBlockNumber, err = bcs.RegistryProxy.GetContract().PunishBlockNumber(nil)
	if err != nil {
		err = fmt.Errorf("get punish block number err %s", err)
	}
	chainID, err = bcs.RegistryProxy.GetContract().ChainId(nil)
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
