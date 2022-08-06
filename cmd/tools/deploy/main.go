package main

import (
	"log"

	"github.com/SmartMeshFoundation/Photon/utils"

	"context"

	"fmt"

	"os"

	"crypto/ecdsa"

	"github.com/SmartMeshFoundation/Photon/accounts"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts/test/tokens/smttoken"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethutils "github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/node"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "address",
			Usage: "The ethereum address you would like Photon to use and for which a keystore file exists in your local system.",
		},
		ethutils.DirectoryFlag{
			Name:  "keystore-path",
			Usage: "If you have a non-standard path for the ethereum keystore directory provide it using this argument. ",
			Value: ethutils.DirectoryString{Value: utils.DefaultKeyStoreDir()},
		},
		cli.StringFlag{
			Name: "eth-rpc-endpoint",
			Usage: `"host:port" address of ethereum JSON-RPC server.\n'
	           'Also accepts a protocol prefix (ws:// or ipc channel) with optional port',`,
			Value: node.DefaultIPCEndpoint("geth"),
		},
		cli.StringFlag{
			Name:  "token-network-address",
			Usage: "only deploy SMTToken with this token-network-address",
		},
	}
	app.Action = mainctx
	app.Name = "photondeploy"
	app.Version = "0.1"
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func mainctx(ctx *cli.Context) error {
	// Create an IPC based RPC connection to a remote node and an authorized transactor
	conn, err := ethclient.Dial(ctx.String("eth-rpc-endpoint"))
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to connect to the Ethereum client: %v", err))
	}
	address := common.HexToAddress(ctx.String("address"))
	address, keybin, err := accounts.PromptAccount(address, ctx.String("keystore-path"), "")
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to unlock account %s", err))
	}
	fmt.Println("start to deploy ...")
	key, err := crypto.ToECDSA(keybin)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to parse priv key %s", err))
	}
	deployContract(key, conn, ctx.String("token-network-address"))
	return nil
}
func deployContract(key *ecdsa.PrivateKey, conn *ethclient.Client, tokenNetworkAddressStr string) {
	auth := bind.NewKeyedTransactor(key)
	ctx := context.Background()
	var tokenNetworkAddress common.Address
	var tx *types.Transaction
	if tokenNetworkAddressStr != "" {
		tokenNetworkAddress = common.HexToAddress(tokenNetworkAddressStr)
	} else {
		// 1. deploy token network
		chainID, err := conn.NetworkID(context.Background())
		if err != nil {
			log.Fatalf("failed to get network id %s", err)
		}
		tokenNetworkAddress, tx, _, err = contracts.DeployTokensNetwork(auth, conn, chainID)
		if err != nil {
			log.Fatalf("failed to deploy registry %s", err)
		}
		_, err = bind.WaitDeployed(ctx, conn, tx)
		if err != nil {
			log.Fatalf("failed to deploy contact when mining :%v", err)
		}
		fmt.Printf("deploy registry complete... RegistryAddress=%s\n", tokenNetworkAddress.String())
	}
	// 2. deploy SMTToken
	tokenAddress, tx, _, err := smttoken.DeploySMTToken(auth, conn, "", tokenNetworkAddress)
	if err != nil {
		log.Fatalf("Failed to DeploySMTToken: %v", err)
	}
	fmt.Printf("SMTToken deploy tx=%s\n", tx.Hash().String())
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to deploy contact when mining :%v", err)
	}
	fmt.Printf("DeploySMTToken complete... tokenAddress=%s\n", tokenAddress.String())
	return
}
