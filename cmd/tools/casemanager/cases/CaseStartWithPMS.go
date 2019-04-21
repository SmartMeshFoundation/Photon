package cases

import (
	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/utils"
)

// CaseStartWithPMS :
func (cm *CaseManager) CaseStartWithPMS() (err error) {
	if !cm.RunSlow {
		return ErrorSkip
	}
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
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N0, N1 := env.Nodes[0], env.Nodes[1]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动pms
	env.StartPMS()
	// 启动节点0,1
	cm.startNodes(env, N0.PMS(), N1.PMS())
	err = N0.OpenChannel(utils.NewRandomAddress().String(), tokenAddress, 1, 58)
	if err == nil {
		return cm.caseFail(env.CaseName)
	}
	err = N0.OpenChannel(utils.NewRandomAddress().String(), tokenAddress, 1, -1)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	//cm.waitForPostman()
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
