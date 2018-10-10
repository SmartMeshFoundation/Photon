package cases

import (
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
)

// CaseCooperateSettle :
func (cm *CaseManager) CaseCooperateSettle() (err error) {
	env, err := models.NewTestEnv("./cases/CaseCooperateSettle.ENV")
	if err != nil {
		return
	}
	defer func() {
		if env.Debug == false {
			env.KillAllRaidenNodes()
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
	// Cooperate settle
	N0.CooperateSettle(c01.ChannelIdentifier)
	time.Sleep(10 * time.Second)
	// 验证
	// verify
	c01new := N0.GetChannelWith(N1, tokenAddress).Println("AfterSettle")

	if c01new != nil && c01new.State != channeltype.StateCooprativeSettle {
		return cm.caseFailWithWrongChannelData(env.CaseName, c01new.Name)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
