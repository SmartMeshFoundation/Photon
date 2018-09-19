package cases

import (
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
)

// CaseWithdraw :
func (cm *CaseManager) CaseWithdraw() (err error) {
	env, err := models.NewTestEnv("./cases/CaseWithdraw.ENV")
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
	var withdrawAmount int32
	withdrawAmount = 1
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N0, N1 := env.Nodes[0], env.Nodes[1]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点2，3
	// start node 2, 3
	N0.Start(env)
	N1.Start(env)

	// 获取channel信息
	// get channel info
	c01 := N0.GetChannelWith(N1, tokenAddress).Println("BeforeWithdraw")

	// withdraw
	N0.Withdraw(c01.ChannelAddress, withdrawAmount)
	time.Sleep(10 * time.Second)
	// 验证
	// verify
	c01new := N0.GetChannelWith(N1, tokenAddress).Println("AfterWithdraw")

	if !c01new.CheckSelfBalance(c01.Balance - withdrawAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, c01new.Name)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
