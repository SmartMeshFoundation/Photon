package main

import (
	"fmt"
	"log"
	"os"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/cases"
	"github.com/urfave/cli"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:  "caselist",
			Usage: "list all cases",
			Action: func(*cli.Context) error {
				cases.NewCaseManager()
				return nil
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "case",
			Usage: "The case number that you want to run. For example, --case=CrashCaseSend01 will run CrashCaseSend01. --case=all run all cases in this path",
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
	// init env
	caseName := ctx.String("case")
	fmt.Println(caseName)
	if caseName != "" {
		// load all cases
		caseManager := cases.NewCaseManager()
		fmt.Println("Start Crash Test...")
		// run case
		if caseName == "all" {
			caseManager.RunAll()
		} else {
			caseManager.RunOne(caseName)
		}
		return
	}
	cli.ShowAppHelp(ctx)
	return
}
