package main

import (
	"log"

	"context"

	"fmt"

	"os"

	"crypto/ecdsa"

	"github.com/SmartMeshFoundation/Photon/accounts"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethutils "github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
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
	deployContract(key, conn)
	return nil
}
func deployContract(key *ecdsa.PrivateKey, conn *ethclient.Client) {
	auth := bind.NewKeyedTransactor(key)
	chainID, err := conn.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("failed to get network id %s", err)
	}
	tokenNetworkAddress, tx, _, err := contracts.DeployTokensNetwork(auth, conn, chainID)
	if err != nil {
		log.Fatalf("failed to deploy registry %s", err)
	}
	ctx := context.Background()
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to deploy contact when mining :%v", err)
	}
	fmt.Printf("deploy registry complete...\n")
	fmt.Printf("RegistryAddress=%s\n \n", tokenNetworkAddress.String())
}
