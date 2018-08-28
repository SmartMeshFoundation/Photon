package main

import (
	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/smoketest/cases"
	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/smoketest/models"
)

var env *models.RaidenEnvReader
var allowFail = true

func main() {
	// 1. init log
	cases.InitCaseLogger("./log/smoketest.log")

	// 2. start raiden nodes
	StartRaidenNode("")

	// 3. init RaidenEnvReader
	hosts := []string{
		"http://127.0.0.1:6000",
		"http://127.0.0.1:6001",
		"http://127.0.0.1:6002",
		"http://127.0.0.1:6003",
		"http://127.0.0.1:6004",
		"http://127.0.0.1:6005",
	}
	env = models.NewRaidenEnvReader(hosts)

	// 4. save all data to before.data
	env.SaveToFile("./log/before.data")

	// 5. Run smoke test
	SmokeTest()

	// 6. save all data to after.data
	env.Refresh()
	env.SaveToFile("./log/after.data")
}
