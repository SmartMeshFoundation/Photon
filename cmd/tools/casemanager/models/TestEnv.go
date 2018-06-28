package models

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"fmt"

	"encoding/json"

	"time"

	"strconv"

	"github.com/SmartMeshFoundation/SmartRaiden"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/huamou/config"
)

// TestEnv env manager for test
type TestEnv struct {
	CaseName                string
	Main                    string
	DataDir                 string
	KeystorePath            string
	PasswordFile            string
	RegistryContractAddress string
	XMPPServer              string
	EthRPCEndpoint          string
	Verbosity               int
	Debug                   bool
	Nodes                   []*RaidenNode
	Tokens                  []*Token
	Channels                []*Channel
	Keys                    []*ecdsa.PrivateKey
}

// Logger : global case logger
var Logger *log.Logger
var globalPassword = "123"

// NewTestEnv default contractor
func NewTestEnv(configFilePath string) (env *TestEnv, err error) {
	c, err := config.ReadDefault(configFilePath)
	if err != nil {
		log.Println("Load config error:", err)
		return
	}
	env = new(TestEnv)
	env.CaseName = c.RdString("COMMON", "case_name", "DefaultName")
	// init logger
	logfile := "./log/" + env.CaseName + ".log"
	logFile, err := os.Create(logfile)
	if err != nil {
		log.Fatalln("Create log file error !", logfile)
	}
	Logger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
	Logger.Println("Start to prepare env for " + env.CaseName + "...")
	env.Main = c.RdString("COMMON", "main", "smartraiden")
	env.DataDir = c.RdString("COMMON", "data_dir", ".smartraiden")
	env.KeystorePath = c.RdString("COMMON", "keystore_path", "../../../testdata/casemanager-keystore")
	env.PasswordFile = c.RdString("COMMON", "password_file", "../../../testdata/casemanager-keystore/pass")
	env.XMPPServer = c.RdString("COMMON", "xmpp-server", "")
	env.EthRPCEndpoint = c.RdString("COMMON", "eth_rpc_endpoint", "ws://182.254.155.208:30306")
	env.Verbosity = c.RdInt("COMMON", "verbosity", 5)
	env.Debug = c.RdBool("COMMON", "debug", false)
	// Create an IPC based RPC connection to a remote node and an authorized transactor
	conn, err := ethclient.Dial(env.EthRPCEndpoint)
	if err != nil {
		Logger.Fatalf(fmt.Sprintf("Failed to connect to the Ethereum client: %v", err))
	}
	_, key := promptAccount(env.KeystorePath)
	registryAddress, registry := loadRegistryContract(c, conn, key)
	env.RegistryContractAddress = registryAddress
	env.Nodes = loadNodes(c)
	env.Tokens = loadTokenAddrs(c, env, conn, key, registry)
	env.Channels = loadAndBuildChannels(c, env, conn)
	env.KillAllRaidenNodes()
	env.ClearHistoryData()
	env.Println(env.CaseName + " env:")
	Logger.Println("Env Prepare SUCCESS")
	return
}

func loadRegistryContract(c *config.Config, conn *ethclient.Client, key *ecdsa.PrivateKey) (addr string, registry *contracts.Registry) {
	addr = c.RdString("COMMON", "registry_contract_address", "new")
	if addr == "new" {
		addr, registry = deployRegistryContract(conn, key)
		Logger.Printf("New RegistryContractAddress : %s\n", addr)
	}
	Logger.Println("Load RegistryContractAddress SUCCESS")
	return
}
func deployRegistryContract(conn *ethclient.Client, key *ecdsa.PrivateKey) (registryContractAddress string, registry *contracts.Registry) {
	auth := bind.NewKeyedTransactor(key)
	//DeployNettingChannelLibrary
	NettingChannelLibraryAddress, tx, _, err := contracts.DeployNettingChannelLibrary(auth, conn)
	if err != nil {
		Logger.Fatalf("Failed to DeployNettingChannelLibrary: %v", err)
	}
	ctx := context.Background()
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		Logger.Fatalf("failed to deploy contact when mining :%v", err)
	}
	//DeployChannelManagerLibrary link nettingchannle library before deploy
	contracts.ChannelManagerLibraryBin = strings.Replace(contracts.ChannelManagerLibraryBin, "__NettingChannelLibrary.sol:NettingCha__", NettingChannelLibraryAddress.String()[2:], -1)
	ChannelManagerLibraryAddress, tx, _, err := contracts.DeployChannelManagerLibrary(auth, conn)
	if err != nil {
		Logger.Fatalf("Failed to deploy new token contract: %v", err)
	}
	ctx = context.Background()
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		Logger.Fatalf("failed to deploy contact when mining :%v", err)
	}
	//DeployRegistry link channelmanagerlibrary before deploy
	contracts.RegistryBin = strings.Replace(contracts.RegistryBin, "__ChannelManagerLibrary.sol:ChannelMan__", ChannelManagerLibraryAddress.String()[2:], -1)
	RegistryContractAddress, tx, _, err := contracts.DeployRegistry(auth, conn)
	if err != nil {
		Logger.Fatalf("Failed to deploy new token contract: %v", err)
	}
	ctx = context.Background()
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		Logger.Fatalf("failed to deploy contact when mining :%v", err)
	}
	fmt.Printf("RegistryAddress=%s\n", RegistryContractAddress.String())
	registry, err = contracts.NewRegistry(RegistryContractAddress, conn)
	if err != nil {
		panic(err)
	}
	registryContractAddress = RegistryContractAddress.String()
	return
}
func promptAccount(keystorePath string) (addr common.Address, key *ecdsa.PrivateKey) {
	am := smartraiden.NewAccountManager(keystorePath)
	if len(am.Accounts) == 0 {
		Logger.Fatal(fmt.Sprintf("No Ethereum accounts found in the directory %s", keystorePath))
		os.Exit(1)
	}
	addr = am.Accounts[0].Address
	for i := 0; i < 3; i++ {
		//retries three times
		if len(globalPassword) <= 0 {
			fmt.Printf("Enter the password to unlock")
			fmt.Scanln(&globalPassword)
		}
		keybin, err := am.GetPrivateKey(addr, globalPassword)
		if err != nil && i == 3 {
			Logger.Fatal(fmt.Sprintf("Exhausted passphrase unlock attempts for %s. Aborting ...", addr))
			os.Exit(1)
		}
		if err != nil {
			Logger.Println(fmt.Sprintf("password incorrect\n Please try again or kill the process to quit.\nUsually Ctrl-c."))
			continue
		}
		key, err = crypto.ToECDSA(keybin)
		if err != nil {
			Logger.Println(fmt.Sprintf("private key to bytes err %s", err))
		}
		break
	}
	return
}
func loadNodes(c *config.Config) (nodes []*RaidenNode) {
	options, _ := c.Options("NODE")
	sort.Strings(options)
	for _, option := range options {
		s := strings.Split(c.RdString("NODE", option, ""), ",")
		nodes = append(nodes, &RaidenNode{
			Name:          option,
			Host:          "http://" + s[1],
			Address:       s[0],
			APIAddress:    s[1],
			ListenAddress: s[1] + "0",
			DebugCrash:    false,
		})
	}
	Logger.Println("Load Nodes SUCCESS")
	return
}

func loadTokenAddrs(c *config.Config, env *TestEnv, conn *ethclient.Client, key *ecdsa.PrivateKey, registry *contracts.Registry) (tokens []*Token) {
	options, _ := c.Options("TOKEN")
	sort.Strings(options)
	for _, option := range options {
		addr := c.RdString("TOKEN", option, "")
		if addr == "new" {
			manager, token, tokenAddress := deployNewToken(env, conn, key, registry)
			addr = tokenAddress.String()
			Logger.Printf("New TokenAddress %s : %s\n", option, addr)
			tokens = append(tokens, &Token{
				Name:    option,
				Address: addr,
				Manager: manager,
				Token:   token,
			})
		} else {
			tokens = append(tokens, &Token{
				Name:    option,
				Address: addr,
			})
		}
	}
	Logger.Println("Load Tokens SUCCESS")
	return
}
func deployNewToken(env *TestEnv, conn *ethclient.Client, key *ecdsa.PrivateKey, registry *contracts.Registry) (manager *contracts.ChannelManagerContract, token *contracts.Token, tokenAddress common.Address) {
	var err error
	mgrAddress, tokenAddress := newToken(key, conn, registry)
	manager, err = contracts.NewChannelManagerContract(mgrAddress, conn)
	if err != nil {
		panic(fmt.Sprintf("err for NewChannelManagerContract %s", err))
	}
	token, err = contracts.NewToken(tokenAddress, conn)
	if err != nil {
		panic(fmt.Sprintf("err for newtoken err %s", err))
	}
	am := smartraiden.NewAccountManager(env.KeystorePath)
	var accounts []common.Address
	for _, node := range env.Nodes {
		address := common.HexToAddress(node.Address)
		accounts = append(accounts, address)
		keybin, err := am.GetPrivateKey(address, globalPassword)
		if err != nil {
			Logger.Fatalf("password error for %s", address.String())
		}
		keytemp, err := crypto.ToECDSA(keybin)
		if err != nil {
			Logger.Fatalf("ToECDSA err %s", err)
		}
		env.Keys = append(env.Keys, keytemp)
	}
	transferMoneyForAccounts(key, conn, accounts, token)
	return manager, token, tokenAddress
}
func newToken(key *ecdsa.PrivateKey, conn *ethclient.Client, registry *contracts.Registry) (mgrAddress common.Address, tokenAddr common.Address) {
	auth := bind.NewKeyedTransactor(key)
	tokenAddr, tx, _, err := contracts.DeployHumanStandardToken(auth, conn, big.NewInt(50000000000), "test", 2, "test symoble")
	if err != nil {
		Logger.Fatalf("Failed to DeployHumanStandardToken: %v", err)
	}
	ctx := context.Background()
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		Logger.Fatalf("failed to deploy contact when mining :%v", err)
	}
	tx, err = registry.AddToken(auth, tokenAddr)
	if err != nil {
		Logger.Fatalf("Failed to AddToken: %v", err)
	}
	ctx = context.Background()
	_, err = bind.WaitMined(ctx, conn, tx)
	if err != nil {
		Logger.Fatalf("failed to AddToken when mining :%v", err)
	}
	mgrAddress, err = registry.ChannelManagerByToken(nil, tokenAddr)
	fmt.Printf("DeployHumanStandardToken complete... token=%s,mgr=%s\n", tokenAddr.String(), mgrAddress.String())
	return
}

// TransferMoneyForAccounts :
func transferMoneyForAccounts(key *ecdsa.PrivateKey, conn *ethclient.Client, accounts []common.Address, token *contracts.Token) {
	wg := sync.WaitGroup{}
	wg.Add(len(accounts))
	auth := bind.NewKeyedTransactor(key)
	nonce, _ := conn.PendingNonceAt(context.Background(), auth.From)
	for index, account := range accounts {
		go func(account common.Address, i int) {
			auth2 := bind.NewKeyedTransactor(key)
			auth2.Nonce = big.NewInt(int64(nonce) + int64(i))
			tx, err := token.Transfer(auth2, account, big.NewInt(5000000))
			if err != nil {
				Logger.Fatalf("Failed to Transfer: %v", err)
			}
			ctx := context.Background()
			_, err = bind.WaitMined(ctx, conn, tx)
			if err != nil {
				Logger.Fatalf("failed to Transfer when mining :%v", err)
			}
			wg.Done()
		}(account, index)
		time.Sleep(time.Millisecond * 100)
	}
	wg.Wait()
	for _, account := range accounts {
		b, _ := token.BalanceOf(nil, account)
		fmt.Printf("account %s has token %s\n", utils.APex(account), b)
	}
}
func loadAndBuildChannels(c *config.Config, env *TestEnv, conn *ethclient.Client) (channels []*Channel) {
	options, _ := c.Options("CHANNEL")
	if options == nil || len(options) == 0 {
		return
	}
	for _, option := range options {
		s := strings.Split(c.RdString("CHANNEL", option, ""), ",")
		_, token := env.GetTokenByName(s[2])
		if token == nil {
			fmt.Println("use old token , do not create channel...")
			return
		}
		index1, account1 := env.GetNodeAddressByName(s[0])
		key1 := env.Keys[index1]
		amount1, _ := strconv.ParseInt(s[3], 10, 64)
		index2, account2 := env.GetNodeAddressByName(s[1])
		key2 := env.Keys[index2]
		amount2, _ := strconv.ParseInt(s[4], 10, 64)
		settledTimeout, _ := strconv.ParseInt(s[5], 10, 64)
		creatAChannelAndDeposit(account1, account2, key1, key2, amount1, amount2, settledTimeout, token.Manager, token.Token, conn)
	}
	Logger.Println("Load and create channels SUCCESS")
	return nil
}

func creatAChannelAndDeposit(account1, account2 common.Address, key1, key2 *ecdsa.PrivateKey, amount1 int64, amount2 int64, settledTimeout int64, manager *contracts.ChannelManagerContract, token *contracts.Token, conn *ethclient.Client) {
	log.Printf("createchannel between %s-%s\n", utils.APex(account1), utils.APex(account2))
	auth1 := bind.NewKeyedTransactor(key1)
	auth1.GasLimit = uint64(params.GasLimit)
	auth1.GasPrice = big.NewInt(params.GasPrice)
	callAuth1 := &bind.CallOpts{
		Pending: false,
		From:    account1,
		Context: context.Background(),
	}
	auth2 := bind.NewKeyedTransactor(key2)
	auth2.GasLimit = uint64(params.GasLimit)
	auth2.GasPrice = big.NewInt(params.GasPrice)
	tx, err := manager.NewChannel(auth1, account2, big.NewInt(settledTimeout))
	if err != nil {
		log.Printf("Failed to NewChannel: %v,%s,%s", err, auth1.From.String(), account2.String())
		return
	}
	ctx := context.Background()
	_, err = bind.WaitMined(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to NewChannel when mining :%v", err)
	}
	//step 2 deopsit
	//step 2.1 aprove
	channelAddress, err := manager.GetChannelWith(callAuth1, account2)
	if err != nil {
		log.Fatalf("failed to get channel %s", err)
		return
	}
	Logger.Printf("New Channel : %s\n", channelAddress.String())
	channel, _ := contracts.NewNettingChannelContract(channelAddress, conn)
	wg2 := sync.WaitGroup{}
	go func() {
		wg2.Add(1)
		defer wg2.Done()
		tx, err := token.Approve(auth1, channelAddress, big.NewInt(amount1))
		if err != nil {
			log.Fatalf("Failed to Approve: %v", err)
		}
		log.Printf("approve gas %s:%d\n", tx.Hash().String(), tx.Gas())
		ctx = context.Background()
		_, err = bind.WaitMined(ctx, conn, tx)
		if err != nil {
			log.Fatalf("failed to Approve when mining :%v", err)
		}
		tx, err = channel.Deposit(auth1, big.NewInt(amount1))
		if err != nil {
			log.Fatalf("Failed to Deposit: %v", err)
		}
		ctx = context.Background()
		_, err = bind.WaitMined(ctx, conn, tx)
		if err != nil {
			log.Fatalf("failed to Deposit when mining :%v", err)
		}
		fmt.Printf("Deposit to account1 %d tokends complete...\n", amount1)
	}()
	go func() {
		wg2.Add(1)
		defer wg2.Done()
		tx, err := token.Approve(auth2, channelAddress, big.NewInt(amount2))
		if err != nil {
			log.Fatalf("Failed to Approve: %v", err)
		}
		ctx = context.Background()
		_, err = bind.WaitMined(ctx, conn, tx)
		if err != nil {
			log.Fatalf("failed to Approve when mining :%v", err)
		}
		tx, err = channel.Deposit(auth2, big.NewInt(amount2))
		if err != nil {
			log.Fatalf("Failed to Deposit: %v", err)
		}
		ctx = context.Background()
		_, err = bind.WaitMined(ctx, conn, tx)
		if err != nil {
			log.Fatalf("failed to Deposit when mining :%v", err)
		}
		fmt.Printf("Deposit to account2 %d tokens complete...\n", amount2)
	}()
	time.Sleep(time.Second * 10)
	wg2.Wait()
}

// KillAllRaidenNodes kill all raiden node
func (env *TestEnv) KillAllRaidenNodes() {
	var pstr2 []string
	//kill the old process
	if runtime.GOOS == "windows" {
		pstr2 = append(pstr2, "-F")
		pstr2 = append(pstr2, "-IM")
		pstr2 = append(pstr2, "smartraiden*")
		ExecShell("taskkill", pstr2, "./log/killall.log", true)
	} else {
		pstr2 = append(pstr2, "smartraiden")
		ExecShell("killall", pstr2, "./log/killall.log", true)
	}
	Logger.Println("Kill all raiden nodes SUCCESS")
}

// ClearHistoryData :
func (env *TestEnv) ClearHistoryData() {
	if env.DataDir == "" {
		return
	}
	filepath.Walk(env.DataDir, func(path string, fi os.FileInfo, err error) error {
		if nil == fi {
			return err
		}
		if !fi.IsDir() {
			return nil
		}
		name := fi.Name()

		if name == ".smartraiden" {
			err := os.RemoveAll(path)
			if err != nil {
				fmt.Println("delet dir error:", err)
			}
		}
		Logger.Println("Clear history data SUCCESS")
		return nil
	})
}

// GetTokenByName :
func (env *TestEnv) GetTokenByName(tokenName string) (index int, token *Token) {
	for index, token := range env.Tokens {
		if token.Name == tokenName {
			return index, token
		}
	}
	return
}

// GetNodeAddressByName :
func (env *TestEnv) GetNodeAddressByName(nodeName string) (index int, address common.Address) {
	for index, node := range env.Nodes {
		if node.Name == nodeName {
			return index, common.HexToAddress(node.Address)
		}
	}
	return
}

// GetNodeByAddress :
func (env *TestEnv) GetNodeByAddress(nodeAddress string) *RaidenNode {
	for _, node := range env.Nodes {
		if node.Address == nodeAddress {
			return node
		}
	}
	return nil
}

//Println print all
func (env *TestEnv) Println(header string) {
	Logger.Println(header)
	buf, err := json.MarshalIndent(env, "", "\t")
	if err != nil {
		panic(err)
	}
	Logger.Println(string(buf))
}
