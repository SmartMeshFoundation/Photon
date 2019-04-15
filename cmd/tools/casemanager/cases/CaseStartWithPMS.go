package cases

import (
	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

// CaseStartWithPMS :
func (cm *CaseManager) CaseStartWithPMS() (err error) {
	env, err := models.NewTestEnv("./cases/2Nodes.ENV", cm.UseMatrix, cm.EthEndPoint, "CaseStartWithPMS")
	if err != nil {
		return
	}
	defer func() {
		if env.Debug == false {
			env.KillAllPhotonNodes()
		}
	}()
	// 源数据
	//tokenAddress := env.Tokens[0].TokenAddress.String()
	N0, N1 := env.Nodes[0], env.Nodes[1]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动pms
	env.StartPMS()
	// 启动节点0,1
	cm.startNodesWithPMS(env, N0, N1)
	cm.waitForPostman()
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
