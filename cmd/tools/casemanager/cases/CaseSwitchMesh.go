package cases

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

// CaseSwitchMesh :
func (cm *CaseManager) CaseSwitchMesh() (err error) {
	env, err := models.NewTestEnv("./cases/CaseSwitchMesh.ENV", cm.UseMatrix, cm.EthEndPoint)
	if err != nil {
		return
	}
	defer func() {
		if env.Debug == false {
			env.KillAllPhotonNodes()
		}
	}()
	// 源数据
	// original data
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N0, N1 := env.Nodes[0], env.Nodes[1]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点2，3
	// start node 2, 3
	cm.startNodes(env, N0, N1)

	err = N0.Transfer(env.Tokens[0].TokenAddress.String(), 1, N1.Address, false)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, err.Error())
	}
	err = N0.SwitchNetwork("true")
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, fmt.Sprintf("switch mesh err %s", err))
	}
	err = N0.Transfer(tokenAddress, 1, N1.Address, false)
	if err == nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, "must fail when mesh is true")
	}
	err = N0.Transfer(tokenAddress, 1, N1.Address, true)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, fmt.Sprintf("fail when direct transfer on mesh %s", err))
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
