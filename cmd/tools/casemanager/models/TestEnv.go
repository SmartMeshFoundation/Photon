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

	"io"

	"github.com/SmartMeshFoundation/SmartRaiden/accounts"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts/test/tokens/tokenerc223approve"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
	XMPPServer              string
	EthRPCEndpoint          string
	RegistryContractAddress string
	Verbosity               int
	Debug                   bool
	Nodes                   []*RaidenNode
	Tokens                  []*Token
	Channels                []*Channel
	Keys                    []*ecdsa.PrivateKey `json:"-"`
}

// Logger : global case logger
var Logger *log.Logger
var globalPassword = "123"

type logTee struct {
	w1 io.Writer
	w2 io.Writer
}

func (t *logTee) Write(p []byte) (n int, err error) {
	n, err = t.w1.Write(p)
	_, err = t.w2.Write(p)
	if err != nil {
		panic(err)
	}
	return
}

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
	Logger = log.New(&logTee{logFile, os.Stderr}, "", log.LstdFlags|log.Lshortfile)
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
	registryAddress, registry := loadTokenNetworkContract(c, conn, key)
	env.RegistryContractAddress = registryAddress.String()
	env.Nodes = loadNodes(c)
	env.Tokens = loadTokenAddrs(c, env, conn, key, registry)
	env.Channels = loadAndBuildChannels(c, env, conn)
	env.KillAllRaidenNodes()
	env.ClearHistoryData()
	env.Println(env.CaseName + " env:")
	Logger.Println("Env Prepare SUCCESS")
	return
}

func loadTokenNetworkContract(c *config.Config, conn *ethclient.Client, key *ecdsa.PrivateKey) (registryAddress common.Address, registry *contracts.TokenNetworkRegistry) {
	addr := c.RdString("COMMON", "token_network_address", "new")
	if addr == "new" {
		registryAddress, registry = deployRegistryContract(conn, key)
		Logger.Printf("New RegistryAddress : %s\n", registryAddress.String())
	} else {
		registryAddress = common.HexToAddress(addr)
	}
	Logger.Println("Load RegistryAddress SUCCESS")
	return
}
func deployRegistryContract(conn *ethclient.Client, key *ecdsa.PrivateKey) (registryAddress common.Address, registry *contracts.TokenNetworkRegistry) {
	auth := bind.NewKeyedTransactor(key)
	//Deploy Secret Registry
	secretRegistryAddress, tx, _, err := contracts.DeploySecretRegistry(auth, conn)
	if err != nil {
		log.Fatalf("Failed to deploy SecretRegistry contract: %v", err)
	}
	ctx := context.Background()
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to deploy contact when mining :%v", err)
	}
	fmt.Printf("Deploy SecretRegistry complete...\n")
	chainID, err := conn.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("failed to get network id %s", err)
	}
	registryAddress, tx, registry, err = contracts.DeployTokenNetworkRegistry(auth, conn, secretRegistryAddress, chainID)
	if err != nil {
		log.Fatalf("failed to deploy TokenNetworkRegistry %s", err)
	}
	ctx = context.Background()
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to deploy contact when mining :%v", err)
	}
	fmt.Printf("deploy TokenNetworkRegistry complete...\n")
	fmt.Printf("TokenNetworkRegistry=%s\n", registryAddress.String())
	return
}
func promptAccount(keystorePath string) (addr common.Address, key *ecdsa.PrivateKey) {
	am := accounts.NewAccountManager(keystorePath)
	if len(am.Accounts) == 0 {
		log.Fatal(fmt.Sprintf("No Ethereum accounts found in the directory %s", keystorePath))
		os.Exit(1)
	}
	addr = am.Accounts[0].Address
	for i := 0; i < 3; i++ {
		//retries three times
		if len(globalPassword) <= 0 {
			fmt.Printf("Enter the password to unlock")
			_, err := fmt.Scanln(&globalPassword)
			if err != nil {
				log.Fatal(err)
			}
		}
		keyBin, err := am.GetPrivateKey(addr, globalPassword)
		if err != nil && i == 3 {
			log.Fatal(fmt.Sprintf("Exhausted passphrase unlock attempts for %s. Aborting ...", addr))
			os.Exit(1)
		}
		if err != nil {
			log.Println(fmt.Sprintf("password incorrect\n Please try again or kill the process to quit.\nUsually Ctrl-c."))
			continue
		}
		key, err = crypto.ToECDSA(keyBin)
		if err != nil {
			log.Println(fmt.Sprintf("private key to bytes err %s", err))
		}
		break
	}
	return
}
func loadNodes(c *config.Config) (nodes []*RaidenNode) {
	options, err := c.Options("NODE")
	if err != nil {
		panic(err)
	}
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

func loadTokenAddrs(c *config.Config, env *TestEnv, conn *ethclient.Client, key *ecdsa.PrivateKey, registry *contracts.TokenNetworkRegistry) (tokens []*Token) {
	options, err := c.Options("TOKEN")
	if err != nil {
		panic(err)
	}
	sort.Strings(options)
	for _, option := range options {
		addr := c.RdString("TOKEN", option, "")
		if addr == "new" {
			tokenNetwork, tokenNetworkAddress, token, tokenAddress := deployNewToken(env, conn, key, registry)
			Logger.Printf("New TokenAddress %s : token=%s token_network=%s", option, tokenAddress.String(), tokenNetworkAddress.String())
			tokens = append(tokens, &Token{
				Name:                option,
				Token:               token,
				TokenAddress:        tokenAddress,
				TokenNetwork:        tokenNetwork,
				TokenNetworkAddress: tokenNetworkAddress,
			})
		} else {
			tokens = append(tokens, &Token{
				Name:         option,
				TokenAddress: common.HexToAddress(addr),
			})
		}
	}
	Logger.Println("Load Tokens SUCCESS")
	return
}
func deployNewToken(env *TestEnv, conn *ethclient.Client, key *ecdsa.PrivateKey, registry *contracts.TokenNetworkRegistry) (tokenNetwork *contracts.TokenNetwork, tokenNetworkAddress common.Address, token *contracts.Token, tokenAddress common.Address) {
	var err error
	tokenNetworkAddress, tokenAddress = newToken(key, conn, registry)
	tokenNetwork, err = contracts.NewTokenNetwork(tokenNetworkAddress, conn)
	if err != nil {
		panic(fmt.Sprintf("err for NewChannelManagerContract %s", err))
	}
	token, err = contracts.NewToken(tokenAddress, conn)
	if err != nil {
		panic(fmt.Sprintf("err for newtoken err %s", err))
	}
	am := accounts.NewAccountManager(env.KeystorePath)
	var accounts []common.Address
	for _, node := range env.Nodes {
		address := common.HexToAddress(node.Address)
		accounts = append(accounts, address)
		keyBin, err := am.GetPrivateKey(address, globalPassword)
		if err != nil {
			Logger.Fatalf("password error for %s", address.String())
		}
		keyTemp, err := crypto.ToECDSA(keyBin)
		if err != nil {
			Logger.Fatalf("ToECDSA err %s", err)
		}
		env.Keys = append(env.Keys, keyTemp)
	}
	transferMoneyForAccounts(key, conn, accounts, token)
	return
}
func newToken(key *ecdsa.PrivateKey, conn *ethclient.Client, tokenNetwork *contracts.TokenNetworkRegistry) (tokenNetworkAddr common.Address, tokenAddr common.Address) {
	auth := bind.NewKeyedTransactor(key)
	tokenAddr, tx, _, err := tokenerc223approve.DeployHumanERC223Token(auth, conn, big.NewInt(50000000000), "test symoble")
	if err != nil {
		log.Fatalf("Failed to DeployHumanStandardToken: %v", err)
	}
	fmt.Printf("token deploy tx=%s\n", tx.Hash().String())
	ctx := context.Background()
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to deploy contact when mining :%v", err)
	}
	fmt.Printf("DeployHumanStandardToken complete...\n")
	tx, err = tokenNetwork.CreateERC20TokenNetwork(auth, tokenAddr)
	if err != nil {
		log.Fatalf("Failed to AddToken: %v", err)
	}
	ctx = context.Background()
	_, err = bind.WaitMined(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to AddToken when mining :%v", err)
	}
	tokenNetworkAddr, err = tokenNetwork.TokenToTokenNetworks(nil, tokenAddr)
	fmt.Printf("DeployHumanStandardToken complete... %s,token_network_address=%s\n", tokenAddr.String(), tokenNetworkAddr.String())
	return
}

// TransferMoneyForAccounts :
func transferMoneyForAccounts(key *ecdsa.PrivateKey, conn *ethclient.Client, accounts []common.Address, token *contracts.Token) {
	wg := sync.WaitGroup{}
	wg.Add(len(accounts))
	auth := bind.NewKeyedTransactor(key)
	nonce, err := conn.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		panic(err)
	}
	for index, account := range accounts {
		go func(account common.Address, i int) {
			auth2 := bind.NewKeyedTransactor(key)
			auth2.Nonce = big.NewInt(int64(nonce) + int64(i))
			tx, err := token.Transfer(auth2, account, big.NewInt(5000000), nil)
			if tx == nil {
				panic("transfer should use approve and transfer from instead")
			}
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
		b, err := token.BalanceOf(nil, account)
		if err != nil {
			panic(err)
		}
		fmt.Printf("account %s has token %s\n", utils.APex(account), b)
	}
}
func loadAndBuildChannels(c *config.Config, env *TestEnv, conn *ethclient.Client) (channels []*Channel) {
	options, err := c.Options("CHANNEL")
	if err != nil {
		panic(err)
	}
	if options == nil || len(options) == 0 {
		return
	}
	for _, option := range options {
		s := strings.Split(c.RdString("CHANNEL", option, ""), ",")
		_, token := env.GetTokenByName(s[2])
		if token.Token == nil {
			fmt.Println("use old token , do not create channel...")
			return
		}
		index1, account1 := env.GetNodeAddressByName(s[0])
		key1 := env.Keys[index1]
		amount1, err := strconv.ParseInt(s[3], 10, 64)
		index2, account2 := env.GetNodeAddressByName(s[1])
		key2 := env.Keys[index2]
		amount2, err := strconv.ParseInt(s[4], 10, 64)
		settledTimeout, err := strconv.ParseUint(s[5], 10, 64)
		if err != nil {
			panic(err)
		}
		creatAChannelAndDeposit(account1, account2, key1, key2, big.NewInt(amount1), big.NewInt(amount2), settledTimeout, token, conn)
	}
	Logger.Println("Load and create channels SUCCESS")
	return nil
}

func creatAChannelAndDeposit(account1, account2 common.Address, key1, key2 *ecdsa.PrivateKey, amount1 *big.Int, amount2 *big.Int, settledTimeout uint64, token *Token, conn *ethclient.Client) {
	log.Printf("createchannel between %s-%s\n", utils.APex(account1), utils.APex(account2))
	var tx *types.Transaction
	var err error
	auth1 := bind.NewKeyedTransactor(key1)
	auth2 := bind.NewKeyedTransactor(key2)
	if amount1.Int64() > 0 {
		approveAccount(token.Token, auth1, token.TokenNetworkAddress, amount1, conn)
		tx, err = token.TokenNetwork.OpenChannelWithDeposit(auth1, account1, account2, settledTimeout, amount1)
	} else {
		tx, err = token.TokenNetwork.OpenChannel(auth1, account1, account2, settledTimeout)
	}
	if err != nil {
		panic(err)
	}
	_, err = bind.WaitMined(context.Background(), conn, tx)
	if err != nil {
		panic(err)
	}
	if amount2.Int64() > 0 {
		approveAccount(token.Token, auth2, token.TokenNetworkAddress, amount2, conn)
		tx, err = token.TokenNetwork.Deposit(auth2, account2, account1, amount2)
		if err != nil {
			panic(err)
		}
		_, err = bind.WaitMined(context.Background(), conn, tx)
		if err != nil {
			panic(err)
		}
	}
}

func approveAccount(token *contracts.Token, auth *bind.TransactOpts, tokenNetworkAddress common.Address, amount *big.Int, conn *ethclient.Client) {
	tx, err := token.Approve(auth, tokenNetworkAddress, amount)
	if err != nil {
		log.Fatalf("Failed to Approve: %v", err)
	}
	log.Printf("approve gas %s:%d\n", tx.Hash().String(), tx.Gas())
	ctx := context.Background()
	_, err = bind.WaitMined(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to Approve when mining :%v", err)
	}
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
		pstr2 = append(pstr2, "-9")
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
