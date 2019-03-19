package codefortest

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"crypto/ecdsa"

	accountModule "github.com/SmartMeshFoundation/Photon/accounts"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/models/stormdb"
	"github.com/SmartMeshFoundation/Photon/network/helper"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// TestEthRPCEndPoint :
var TestEthRPCEndPoint = os.Getenv("ETHRPCENDPOINT")

// TestKeystorePath :
var TestKeystorePath = os.Getenv("KEYSTORE")

// TestPassword :
var TestPassword = "123"

// DeployRegistryContract :
func DeployRegistryContract() (registryAddress common.Address, registry *contracts.TokensNetwork, secretRegistryAddress common.Address, err error) {
	var tx *types.Transaction
	conn, err := GetEthClient()
	if err != nil {
		return
	}
	defer conn.Close()

	accounts, err := GetAccounts()
	if err != nil {
		return
	}
	key := accounts[0].PrivateKey
	auth := bind.NewKeyedTransactor(key)
	chainID, err := conn.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("failed to get network id %s", err)
	}
	registryAddress, tx, registry, err = contracts.DeployTokensNetwork(auth, conn, chainID)
	if err != nil {
		err = fmt.Errorf("failed to deploy TokenNetworkRegistry %s", err)
		return
	}
	ctx := context.Background()
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		err = fmt.Errorf("failed to deploy contact when mining :%v", err)
		return
	}
	fmt.Printf("deploy TokenNetworkRegistry complete...\n")
	fmt.Printf("TokenNetworkRegistryAddress=%s, SecretRgistryAddess=%s\n", registryAddress.String(), secretRegistryAddress.String())
	return
}

// GetEthClient :
func GetEthClient() (client *helper.SafeEthClient, err error) {
	return helper.NewSafeClient(TestEthRPCEndPoint)
}

// TestAccount :
type TestAccount struct {
	Address    common.Address
	PrivateKey *ecdsa.PrivateKey
}

// GetAccounts unlock all accounts on TestKeystorePath
func GetAccounts() (accounts []TestAccount, err error) {
	am := accountModule.NewAccountManager(TestKeystorePath)
	if len(am.Accounts) == 0 {
		err = fmt.Errorf("no ethereum accounts found in the directory [%s]", TestKeystorePath)
		return
	}
	for _, a := range am.Accounts {
		var keyBin []byte
		var key *ecdsa.PrivateKey
		keyBin, err = am.GetPrivateKey(a.Address, TestPassword)
		if err != nil {
			return
		}
		key, err = crypto.ToECDSA(keyBin)
		if err != nil {
			return
		}
		accounts = append(accounts, TestAccount{
			Address:    a.Address,
			PrivateKey: key,
		})
	}
	return
}

// GetAccountsByAddress : find and unlock  account by TestPassword
func GetAccountsByAddress(address common.Address) (account TestAccount, err error) {
	accounts, err := GetAccounts()
	if err != nil {
		return
	}
	for _, a := range accounts {
		if a.Address == address {
			account = a
			return
		}
	}
	err = errors.New("no account in keystore")
	return
}

// NewTestDB :
func NewTestDB(dbPath string) (dao models.Dao) {
	if dbPath == "" {
		dbPath = path.Join(os.TempDir(), "testxxxx.db")
		err := os.RemoveAll(dbPath)
		err = os.RemoveAll(dbPath + ".lock")
		if err != nil {
			fmt.Println(err)
		}
	}
	var err error
	//if os.Getenv("PHOTON_DB") == "gkv" {
	//	fmt.Println("use gkv db")
	//	dao, err = gkvdb.OpenDb(dbPath)
	//	if err != nil {
	//		panic(err)
	//	}
	//} else {
	fmt.Println("use storm db")
	dao, err = stormdb.OpenDb(dbPath)
	if err != nil {
		panic(err)
	}
	//}
	return
}
