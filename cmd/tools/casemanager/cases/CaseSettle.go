package cases

import (
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

// CaseSettle :
func (cm *CaseManager) CaseSettle() (err error) {
	env, err := models.NewTestEnv("./cases/CaseSettle.ENV", cm.UseMatrix)
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
	c01 := N0.GetChannelWith(N1, tokenAddress).Println("BeforeClose")
	N0.SendTrans(env.Tokens[0].TokenAddress.String(), 1, N1.Address, false)
	time.Sleep(3 * time.Second)
	// Close
	N0.Close(c01.ChannelIdentifier)
	N0.GetChannelWith(N1, tokenAddress).Println("AfterClose")
	// Settle
	time.Sleep(time.Duration(c01.SettleTimeout+257+10) * time.Second)
	N0.Settle(c01.ChannelIdentifier)
	time.Sleep(10 * time.Second)
	// 验证
	// verify
	c01new := N0.GetChannelWith(N1, tokenAddress).Println("AfterSettle")

	if c01new != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, c01new.Name)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
