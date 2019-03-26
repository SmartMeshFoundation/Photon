package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/cases"
	"github.com/urfave/cli"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:  "caselist",
			Usage: "list all cases",
			Action: func(ctx *cli.Context) error {
				cases.NewCaseManager(false, false, ctx.String("eth-rpc-endpoint"), false)
				return nil
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "case",
			Usage: "The case number that you want to run. For example, --case=CrashCaseSend01 will run CrashCaseSend01. --case=all run all cases in this path",
		},
		cli.StringFlag{
			Name:  "skip",
			Usage: "true to skip failed cases,default false",
			Value: "false",
		},
		cli.BoolFlag{
			Name:  "auto",
			Usage: "true if auto run",
		},
		cli.BoolFlag{
			Name:  "matrix",
			Usage: "true if run with matrix",
		},
		cli.StringFlag{
			Name:  "eth-rpc-endpoint",
			Usage: "eth rpc end point like : http://127.0.0.1:8545",
			Value: "http://127.0.0.1:30307",
		},
		cli.BoolFlag{
			Name:  "slow",
			Usage: "for long test,run every case",
		},
	}
	app.Action = Main
	app.Name = "case-manager"
	err := app.Run(os.Args)
	if err != nil {
		log.Printf("run err %s\n", err)
	}
}

// Main crash test
func Main(ctx *cli.Context) (err error) {
	start := time.Now()
	// init env
	caseName := ctx.String("case")
	fmt.Println(caseName)
	if caseName != "" {
		// load all cases
		caseManager := cases.NewCaseManager(ctx.Bool("auto"), ctx.Bool("matrix"), ctx.String("eth-rpc-endpoint"), ctx.Bool("slow"))
		fmt.Println("Start Casemanager Test...")
		// run case
		if caseName == "all" {
			//caseManager.RunThisCaseOnly = true
			caseManager.RunAll(ctx.String("skip"))
		} else {
			caseManager.RunThisCaseOnly = true
			caseManager.RunSlow = true
			caseManager.RunOne(caseName)
		}
		end := time.Now()
		log.Printf("casemanager time use:%s", end.Sub(start))
		return
	}
	err = cli.ShowAppHelp(ctx)
	return
}
