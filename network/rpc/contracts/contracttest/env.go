package contracttest

import (
	"crypto/ecdsa"
	"log"

	"testing"

	"context"
	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/accounts"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/huamou/config"
)

// Env :
type Env struct {
	KeystorePath          string
	EthRPCEndpoint        string
	Token                 *contracts.Token
	TokenAddress          common.Address
	TokenNetworkAddress   common.Address
	Client                *ethclient.Client
	TokenNetworkRegistry  *contracts.TokenNetworkRegistry
	TokenNetwork          *contracts.TokenNetwork
	SecretRegistryAddress common.Address
	SecretRegistry        *contracts.SecretRegistry
	Accounts              []*Account
	isFirst               bool
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
	//  get the client
	env.Client, err = ethclient.Dial(env.EthRPCEndpoint)
	if err != nil {
		panic(err)
	}
	t.Logf("Geth client = %s", env.EthRPCEndpoint)
	// get secret registry
	env.SecretRegistryAddress = common.HexToAddress(c.RdString("COMMON", "secret_registry_address", "new"))
	env.SecretRegistry, err = contracts.NewSecretRegistry(env.SecretRegistryAddress, env.Client)
	if err != nil {
		t.Error(err)
		return
	}
	// get token
	env.TokenAddress = common.HexToAddress(c.RdString("COMMON", "token_address", "new"))
	env.Token, err = contracts.NewToken(env.TokenAddress, env.Client)
	if err != nil {
		panic(err)
	}
	t.Logf("Token = %s", env.TokenAddress.String())
	// get token_network
	tokenNetworkAddress := c.RdString("COMMON", "token_network_address", "")
	if tokenNetworkAddress == "new" || tokenNetworkAddress == "" {
		// Deploy a new token_network contract
	} else {
		env.TokenNetworkAddress = common.HexToAddress(tokenNetworkAddress)
		env.TokenNetwork, err = contracts.NewTokenNetwork(env.TokenNetworkAddress, env.Client)
		if err != nil {
			panic(err)
		}
	}
	t.Logf("TokenNetwork = %s", tokenNetworkAddress)
	// get token network registry
	tokenNetworkRegistryAddress := common.HexToAddress(c.RdString("COMMON", "token_network_registry_address", "new"))
	env.TokenNetworkRegistry, err = contracts.NewTokenNetworkRegistry(tokenNetworkRegistryAddress, env.Client)
	if err != nil {
		panic(err)
	}
	// init accounts, keys and auths
	initAccounts(t, env)
	t.Log("=======================================> env init done, test BEGIN ...")
	return
}

func initAccounts(t *testing.T, env *Env) {
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
		tx, err := env.Token.Approve(envAccount.Auth, env.TokenNetworkAddress, big.NewInt(50000000))
		if err != nil {
			t.Error(err)
			return
		}
		r, err := bind.WaitMined(context.Background(), env.Client, tx)
		if err != nil {
			t.Error(err)
			return
		}
		if r.Status != types.ReceiptStatusSuccessful {
			t.Error("receipt status error")
			return
		}
		env.Accounts = append(env.Accounts, envAccount)
	}
	t.Logf("load [%d] accouts from [%s] done ...", len(env.Accounts), env.KeystorePath)
}
