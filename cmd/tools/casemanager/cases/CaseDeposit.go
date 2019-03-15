package cases

import (
	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

// CaseDeposit :
func (cm *CaseManager) CaseDeposit() (err error) {
	env, err := models.NewTestEnv("./cases/CaseDeposit.ENV", cm.UseMatrix, cm.EthEndPoint)
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

	// 获取channel信息
	// get channel info
	c01 := N0.GetChannelWith(N1, tokenAddress).Println("before deposit")

	// N0 deposit
	err = N0.Deposit(N1.Address, tokenAddress, 50)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, err.Error())
	}
	err = cm.tryInSeconds(cm.HighMediumWaitSeconds, func() error {
		// check
		c01new := N0.GetChannelWith(N1, tokenAddress).Println("after deposit")
		if !c01new.CheckEqualByPartnerNode(env) || !c01new.CheckSelfBalance(c01.Balance+50) {
			return cm.caseFailWithWrongChannelData(env.CaseName, c01new.Name)
		}
		return nil
	})
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, err.Error())
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
