package cases

import (
	"time"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

// CaseSendTransferOnReOpenedChannel :
func (cm *CaseManager) CaseSendTransferOnReOpenedChannel() (err error) {
	env, err := models.NewTestEnv("./cases/CaseSendTransferOnReOpenedChannel.ENV", cm.UseMatrix)
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
	N0.SendTrans(env.Tokens[0].TokenAddress.String(), 1, N1.Address, true)
	//time.Sleep(3 * time.Second)
	// Cooperate settle
	N0.CooperateSettle(c01.ChannelIdentifier)
	time.Sleep(10 * time.Second)
	// 验证
	// verify
	c01new := N0.GetChannelWith(N1, tokenAddress).Println("AfterSettle")

	if c01new != nil && c01new.State != channeltype.StateCooprativeSettle {
		return cm.caseFailWithWrongChannelData(env.CaseName, c01new.Name)
	}

	//重新创建通道,并交易
	err = N0.OpenChannel(N1.Address, tokenAddress, int64(c01.Balance), int64(c01.SettleTimeout))
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, c01.Name)
	}
	N0.SendTrans(tokenAddress, 1, N1.Address, true)
	c01new = N0.GetChannelWith(N1, tokenAddress).Println("after Reopen channel")
	if !c01new.CheckSelfBalance(c01.Balance - 1) {
		return cm.caseFail(env.CaseName)
	}
	if !c01new.CheckPartnerBalance(1) {
		return cm.caseFail(env.CaseName)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
