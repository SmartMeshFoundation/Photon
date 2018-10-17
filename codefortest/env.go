package codefortest

import (
	"context"
	"fmt"
	"log"
	"os"

	"crypto/ecdsa"

	accountModule "github.com/SmartMeshFoundation/SmartRaiden/accounts"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// TestEthRPCEndPoint :
var TestEthRPCEndPoint = os.Getenv("ETHRPCENDPOINT")

// TestKeystorePath :
var TestKeystorePath = os.Getenv("KEYSTORE")

// TestPassword :
var TestPassword = "123"

// DeployRegistryContract :
func DeployRegistryContract() (registryAddress common.Address, registry *contracts.TokenNetworkRegistry, err error) {
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

	//Deploy Secret Registry
	secretRegistryAddress, tx, _, err := contracts.DeploySecretRegistry(auth, conn)
	if err != nil {
		err = fmt.Errorf("failed to deploy SecretRegistry contract: %v", err)
		return
	}
	ctx := context.Background()
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		err = fmt.Errorf("failed to deploy contact when mining :%v", err)
		return
	}
	fmt.Printf("Deploy SecretRegistry complete...\n")
	chainID, err := conn.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("failed to get network id %s", err)
	}
	registryAddress, tx, registry, err = contracts.DeployTokenNetworkRegistry(auth, conn, secretRegistryAddress, chainID)
	if err != nil {
		err = fmt.Errorf("failed to deploy TokenNetworkRegistry %s", err)
		return
	}
	ctx = context.Background()
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

// GetAccounts :
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
