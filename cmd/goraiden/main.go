package main

import (
	"os"

	"fmt"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/node"
	"github.com/nkbai/go-raiden/params"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	var language string

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "address",
			Usage:       "The ethereum address you would like raiden to use and for which a keystore file exists in your local system.",
			Destination: &language,
		},
		utils.DirectoryFlag{
			Name:  "keystore-path",
			Usage: "If you have a non-standard path for the ethereum keystore directory provide it using this argument. ",
			Value: utils.DirectoryString{params.DefaultKeyStoreDir()},
		},
		cli.StringFlag{
			Name: "eth-rpc-endpoint",
			Usage: `"host:port" address of ethereum JSON-RPC server.\n'
	           'Also accepts a protocol prefix (ws:// or ipc channel) with optional port',`,
			Value: node.DefaultIPCEndpoint("geth"),
		},
		cli.StringFlag{
			Name:  "registry-contract-address",
			Usage: `hex encoded address of the registry contract.`,
			Value: params.ROPSTEN_REGISTRY_ADDRESS.String(),
		},
		cli.StringFlag{
			Name:  "discovery-contract-address",
			Usage: `hex encoded address of the discovery contract.`,
			Value: params.ROPSTEN_DISCOVERY_ADDRESS.String(),
		},
		cli.StringFlag{
			Name:  "listen-address",
			Usage: `"host:port" for the raiden service to listen on.`,
			Value: fmt.Sprintf("0.0.0.0:%d", params.INITIAL_PORT),
		},
		cli.StringFlag{
			Name: "rpccorsdomain",
			Usage: `Comma separated list of domains to accept cross origin requests.
				(localhost enabled by default)`,
			Value: "http://localhost:* /*",
		},
		cli.StringFlag{Name: "logging",
			Usage: `ethereum.slogging config-string (\'<logger1>:<level>,<logger2>:<level>\')'`,
			Value: ":TRACE"},
		cli.StringFlag{
			Name:  "logfile",
			Usage: "file path for logging to file",
			Value: "",
		},
		cli.IntFlag{Name: "max-unresponsive-time",
			Usage: `Max time in seconds for which an address can send no packets and
	               still be considered healthy.`,
			Value: 120,
		},
		cli.IntFlag{Name: "send-ping-time",
			Usage: `Time in seconds after which if we have received no message from a
	               node we have a connection with, we are going to send a PING message`,
			Value: 60,
		},
		cli.BoolTFlag{Name: "rpc",
			Usage: `Start with or without the RPC server. Default is to start
	               the RPC server`,
		},
		cli.StringFlag{
			Name:  "api-address",
			Usage: `host:port" for the RPC server to listen on.`,
			Value: "0.0.0.0:5001",
		},
		utils.DirectoryFlag{
			Name:  "datadir",
			Usage: "Directory for storing raiden data.",
			Value: utils.DirectoryString{params.DefaultDataDir()},
		},
		cli.StringFlag{
			Name:  "password-file",
			Usage: "Text file containing password for provided account",
		},
	}
	app.Action = func(ctx *cli.Context) error {
		return nil
	}
	app.Name = "raiden"
	app.Version = "0.2"
	app.Run(os.Args)
}
