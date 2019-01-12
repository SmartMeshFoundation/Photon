package cases

import (
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

// CaseCooperateSettle :
func (cm *CaseManager) CaseCooperateSettle() (err error) {
	env, err := models.NewTestEnv("./cases/CaseCooperateSettle.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	N0.Start(env)
	N1.Start(env)

	// 获取channel信息
	// get channel info
	c01 := N0.GetChannelWith(N1, tokenAddress).Println("before cooperative settle")
	N0.SendTrans(env.Tokens[0].TokenAddress.String(), 1, N1.Address, false)

	// Cooperate settle
	N0.CooperateSettle(c01.ChannelIdentifier)
	var i = 0
	for i = 0; i < cm.MediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		// 验证
		// verify
		c01new := N0.GetChannelWith(N1, tokenAddress).Println("AfterSettle")

		if c01new == nil {
			break
		}
	}
	if i == cm.MediumWaitSeconds {
		return cm.caseFailWithWrongChannelData(env.CaseName, c01.Name)
	}
	err = N0.Transfer(tokenAddress, 1, N1.Address, false)
	if err == nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, "Transfer must failed after cooperate settle")
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
