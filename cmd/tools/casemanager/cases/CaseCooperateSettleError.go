package cases

import (
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

// CaseCooperateSettleError :
func (cm *CaseManager) CaseCooperateSettleError() (err error) {
	env, err := models.NewTestEnv("./cases/CaseCooperateSettleError.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	c01 := N0.GetChannelWith(N1, tokenAddress).Println("before cooperative settle")
	N0.SendTrans(env.Tokens[0].TokenAddress.String(), 1, N1.Address, false)
	time.Sleep(time.Second)
	// n1删除数据并重启
	N1.Shutdown(env)
	N1.ClearHistoryData(env.DataDir)
	N1.Start(env)

	if cm.UseMatrix{
		time.Sleep(time.Second *5)
	}
	// Cooperate settle
	err = N0.CooperateSettle(c01.ChannelIdentifier, -1)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, err.Error())
	}
	if cm.UseMatrix{
		time.Sleep(time.Second *10)
	}
	// 等待N0接收消息
	time.Sleep(time.Second * 2)
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
