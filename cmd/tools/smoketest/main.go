package main

import (
	"os"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/smoketest/cases"
	"github.com/SmartMeshFoundation/Photon/cmd/tools/smoketest/models"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/utils"
)

var env *models.PhotonEnvReader
var allowFail = true

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
}
func main() {
	// 1. init log
	cases.InitCaseLogger("./log/smoketest.log")

	// 2. start photon nodes
	StartPhotonNode("")

	// 3. init PhotonEnvReader
	hosts := []string{
		"http://127.0.0.1:6000",
		"http://127.0.0.1:6001",
		"http://127.0.0.1:6002",
		"http://127.0.0.1:6003",
		"http://127.0.0.1:6004",
		"http://127.0.0.1:6005",
	}
	env = models.NewPhotonEnvReader(hosts)

	// 4. save all data to before.data
	env.SaveToFile("./log/before.data")

	// 5. Run smoke test
	SmokeTest()
}
