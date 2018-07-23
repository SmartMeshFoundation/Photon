package contracttest

import (
	"crypto/ecdsa"
	"log"

	"testing"

	"github.com/SmartMeshFoundation/SmartRaiden/accounts"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/huamou/config"
)

// Env :
type Env struct {
	KeystorePath        string
	EthRPCEndpoint      string
	TokenNetworkAddress common.Address
	Client              *ethclient.Client
	TokenNetwork        *contracts.TokenNetwork
	Accounts            []*Account
	isFirst             bool
}

// Account :
type Account struct {
	Address common.Address
	Key     *ecdsa.PrivateKey
	Auth    *bind.TransactOpts
}

var env *Env
var globalPassword = "123"

// InitEnv :
func InitEnv(t *testing.T, configFilePath string) {
	if env != nil {
		env.isFirst = false
		return
	}
	// load config
	c, err := config.ReadDefault(configFilePath)
	if err != nil {
		log.Println("Load config error:", err)
		return
	}
	env = new(Env)
	env.isFirst = true
	env.KeystorePath = c.RdString("COMMON", "keystore_path", "../../../testdata/casemanager-keystore")
	env.EthRPCEndpoint = c.RdString("COMMON", "eth_rpc_endpoint", "ws://182.254.155.208:30306")
	tokenNetworkAddress := c.RdString("COMMON", "token_network_address", "")
	//  get the client
	env.Client, err = ethclient.Dial(env.EthRPCEndpoint)
	if err != nil {
		panic(err)
	}
	// get token_network
	if tokenNetworkAddress == "new" || tokenNetworkAddress == "" {
		// Deploy a new token_network contract
	} else {
		env.TokenNetworkAddress = common.StringToAddress(tokenNetworkAddress)
		env.TokenNetwork, err = contracts.NewTokenNetwork(env.TokenNetworkAddress, env.Client)
		if err != nil {
			panic(err)
		}
	}
	// init accounts, keys and auths
	initAccounts(env)
	t.Log("env init done")
	return
}

func initAccounts(env *Env) {
	am := accounts.NewAccountManager(env.KeystorePath)
	for _, account := range am.Accounts {
		keyBin, err := am.GetPrivateKey(account.Address, globalPassword)
		if err != nil {
			log.Fatalf("password error for %s,err=%s", utils.APex2(account.Address), err)
		}
		keyTemp, err := crypto.ToECDSA(keyBin)
		if err != nil {
			log.Fatalf("toecdsa err %s", err)
		}
		envAccount := new(Account)
		envAccount.Address = account.Address
		envAccount.Key = keyTemp
		envAccount.Auth = bind.NewKeyedTransactor(keyTemp)
		env.Accounts = append(env.Accounts, envAccount)
	}
}
