package main

import (
	"log"

	"context"

	"fmt"

	"os"

	"crypto/ecdsa"

	"github.com/SmartMeshFoundation/SmartRaiden/accounts"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethutils "github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/node"
	"github.com/slonzok/getpass"
	"gopkg.in/urfave/cli.v1"
)

func main() {
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
	}
	app.Action = mainctx
	app.Name = "raidendeploy"
	app.Version = "0.1"
	app.Run(os.Args)
}

func mainctx(ctx *cli.Context) error {
	// Create an IPC based RPC connection to a remote node and an authorized transactor
	conn, err := ethclient.Dial(ctx.String("eth-rpc-endpoint"))
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to connect to the Ethereum client: %v", err))
	}
	address := common.HexToAddress(ctx.String("address"))
	address, key := promptAccount(address, ctx.String("keystore-path"))
	fmt.Println("start to deploy ...")
	deployContract(key, conn)
	return nil
}
func deployContract(key *ecdsa.PrivateKey, conn *ethclient.Client) {
	auth := bind.NewKeyedTransactor(key)
	//Deploy Secret Registry
	secretRegistryAddress, tx, _, err := contracts.DeploySecretRegistry(auth, conn)
	if err != nil {
		log.Fatalf("Failed to deploy new token contract: %v", err)
	}
	ctx := context.Background()
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to deploy contact when mining :%v", err)
	}
	fmt.Printf("Deploy Secret Registry complete...\n")
	chainId, err := conn.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("failed to get network id %s", err)
	}
	registryAddress, tx, _, err := contracts.DeployTokenNetworkRegistry(auth, conn, secretRegistryAddress, chainId)
	if err != nil {
		log.Fatalf("failed to deploy registry %s", err)
	}
	ctx = context.Background()
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to deploy contact when mining :%v", err)
	}
	fmt.Printf("deploy registry complete...\n")
	fmt.Printf("RegistryAddress=%s\nSecretyRegistryAddress=%s\n", registryAddress.String(), secretRegistryAddress.String())
	//RegistryAddress=0x1026a4441921EcF88aaF13014d96aF90f735a02c
	//EndpointRegistryAddress=0xB85b8b57e2b701d5E918D7d9027A7330472a663a
}
func promptAccount(adviceAddress common.Address, keystorePath string) (addr common.Address, key *ecdsa.PrivateKey) {
	am := accounts.NewAccountManager(keystorePath)
	if len(am.Accounts) == 0 {
		log.Fatal(fmt.Sprintf("No Ethereum accounts found in the directory %s", keystorePath))
		os.Exit(1)
	}
	if !am.AddressInKeyStore(adviceAddress) {
		if adviceAddress != utils.EmptyAddress {
			log.Fatal(fmt.Sprintf("account %s could not be found on the sytstem. aborting...", adviceAddress))
			os.Exit(1)
		}
		shouldPromt := true
		fmt.Println("The following accounts were found in your machine:")
		for i := 0; i < len(am.Accounts); i++ {
			fmt.Printf("%3d -  %s\n", i, am.Accounts[i].Address.String())
		}
		fmt.Println("")
		for shouldPromt {
			fmt.Printf("Select one of them by index to continue:")
			idx := -1
			fmt.Scanf("%d\n", &idx)
			if idx >= 0 && idx < len(am.Accounts) {
				shouldPromt = false
				addr = am.Accounts[idx].Address
			} else {
				fmt.Printf("Error: Provided index %d is out of bounds", idx)
			}
		}
	} else {
		addr = adviceAddress
	}
	var password string
	for i := 0; i < 3; i++ {
		//retries three times
		password = getpass.Prompt("Enter the password to unlock")
		//fmt.Printf("\npassword is %s\n", password)
		keybin, err := am.GetPrivateKey(addr, password)
		if err != nil && i == 3 {
			log.Fatal(fmt.Sprintf("Exhausted passphrase unlock attempts for %s. Aborting ...", addr))
			os.Exit(1)
		}
		if err != nil {
			log.Println(fmt.Sprintf("password incorrect\n Please try again or kill the process to quit.\nUsually Ctrl-c."))
			continue
		}
		key, err = crypto.ToECDSA(keybin)
		if err != nil {
			log.Fatalf("private key to bytes err %s", err)
			os.Exit(1)
		}
		break
	}
	return
}
