package main

import (
	"fmt"

	"os"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/mobile"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	ethutils "github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/urfave/cli"
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
		},
		cli.StringFlag{
			Name: "eth-rpc-endpoint",
			Usage: `"host:port" address of ethereum JSON-RPC server.\n'
	           'Also accepts a protocol prefix (ws:// or ipc channel) with optional port',`,
		},
		cli.StringFlag{
			Name:  "registry-contract-address",
			Usage: `hex encoded address of the registry contract.`,
		},
		cli.StringFlag{
			Name:  "listen-address",
			Usage: `"host:port" for the raiden service to listen on.`,
			Value: fmt.Sprintf("0.0.0.0:%d", params.InitialPort),
		},
		cli.StringFlag{
			Name:  "api-address",
			Value: "127.0.0.1:5001",
		},
		ethutils.DirectoryFlag{
			Name:  "datadir",
			Usage: "Directory for storing raiden data.",
		},
		cli.StringFlag{
			Name:  "password-file",
			Usage: "Text file containing password for provided account",
		},
	}
	app.Action = mainCtx
	app.Name = "mobiletest"
	app.Version = "0.3"
	err := app.Run(os.Args)
	if err != nil {
		log.Crit(err.Error())
	}
}
func mainCtx(ctx *cli.Context) (err error) {
	fmt.Printf("Welcom to mobiletest,version %s\n", ctx.App.Version)
	address := ctx.String("address")
	keystorePath := ctx.String("keystore-path")
	ethRPCEndpoint := ctx.String("eth-rpc-endpoint")
	registryContractAddress := ctx.String("registry-contract-address")
	listenAddress := ctx.String("listen-address")
	dataDir := ctx.String("datadir")
	password := ctx.String("password-file")
	apiAddr := ctx.String("api-address")
	otherArgs := mobile.NewStrings(1)
	err = otherArgs.Set(0, fmt.Sprintf("--registry-contract-address=%s", registryContractAddress))
	if err != nil {
		return err
	}
	api, err := mobile.StartUp(address, keystorePath, ethRPCEndpoint, dataDir, password, apiAddr, listenAddress,
		"", params.SpectrumTestNetRegistryAddress.String(),
		otherArgs)
	if err != nil {
		log.Crit(fmt.Sprintf("start up err %s", err))
		return
	}
	sub, err := api.Subscribe(handler{})
	if err != nil {
		log.Crit(fmt.Sprintf("sub err %s", err))
	}
	time.Sleep(time.Hour)
	sub.Unsubscribe()
	return nil
}

type handler struct {
}

//some unexpected error
func (h handler) OnError(errCode int, failure string) {
	log.Error(fmt.Sprintf("receive err %d, %s", errCode, failure))
}

//OnStatusChange server connection status change
func (h handler) OnStatusChange(s string) {
	log.Error(fmt.Sprintf("receive status change %s", s))
}

//OnReceivedTransfer  receive a transfer
func (h handler) OnReceivedTransfer(tr string) {
	log.Error(fmt.Sprintf("receive transfer %s", tr))
}

//OnSentTransfer a transfer sent success
func (h handler) OnSentTransfer(tr string) {
	log.Error(fmt.Sprintf("sent transfer %s", tr))
}

func (h handler) OnNotify(level int, info string) {
	log.Info(fmt.Sprintf("Receive notice : level=%d info=%s\n", level, info))
}
