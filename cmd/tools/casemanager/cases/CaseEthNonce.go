package cases

import (
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
)

// CaseEthNonce :
func (cm *CaseManager) CaseEthNonce() (err error) {
	env, err := models.NewTestEnv("./cases/CaseEthNonce.ENV")
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
	N0, N1, N2, N3, N4 := env.Nodes[0], env.Nodes[1], env.Nodes[2], env.Nodes[3], env.Nodes[4]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点2，3
	// start node 2, 3
	N0.Start(env)
	N1.Start(env)
	N2.Start(env)
	N3.Start(env)
	N4.Start(env)

	// 获取channel信息
	// get channel info
	c01 := N0.GetChannelWith(N1, tokenAddress).Println("BeforeCooperateSettle")
	c12 := N1.GetChannelWith(N2, tokenAddress).Println("BeforeCooperateSettle")
	c13 := N1.GetChannelWith(N3, tokenAddress).Println("BeforeCooperateSettle")
	c14 := N1.GetChannelWith(N4, tokenAddress).Println("BeforeCooperateSettle")
	// Cooperate settle
	N1.CooperateSettle(c01.ChannelIdentifier)
	N1.CooperateSettle(c12.ChannelIdentifier)
	N1.CooperateSettle(c13.ChannelIdentifier)
	N1.CooperateSettle(c14.ChannelIdentifier)
	time.Sleep(20 * time.Second)
	// 验证
	// verify
	c01 = N0.GetChannelWith(N1, tokenAddress)
	c12 = N1.GetChannelWith(N2, tokenAddress)
	c13 = N1.GetChannelWith(N3, tokenAddress)
	c14 = N1.GetChannelWith(N4, tokenAddress)
	if c01 != nil {
		c01.Println("Wrong, should be nil")
		return cm.caseFailWithWrongChannelData(env.CaseName, c01.Name)
	}
	if c12 != nil {
		c12.Println("Wrong, should be nil")
		return cm.caseFailWithWrongChannelData(env.CaseName, c12.Name)
	}
	if c13 != nil {
		c13.Println("Wrong, should be nil")
		return cm.caseFailWithWrongChannelData(env.CaseName, c13.Name)
	}
	if c14 != nil {
		c14.Println("Wrong, should be nil")
		return cm.caseFailWithWrongChannelData(env.CaseName, c14.Name)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
