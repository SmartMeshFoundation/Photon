package cases

import (
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/Photon/utils"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

// CaseSendTransferOnReOpenedChannel :
func (cm *CaseManager) CaseSendTransferOnReOpenedChannel() (err error) {
	env, err := models.NewTestEnv("./cases/CaseSendTransferOnReOpenedChannel.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	trs, err := N0.GetSentTransfers()
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, fmt.Sprintf("GetSentTransfers err %s", err))
	}
	if len(trs) != 1 {
		return cm.caseFailWithWrongChannelData(env.CaseName, "trs 1 err")
	}
	//time.Sleep(3 * time.Second)
	// Cooperate settle
	N0.CooperateSettle(c01.ChannelIdentifier)
	i := 0
	for i = 0; i < cm.MediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		// 验证
		// verify
		c01new := N0.GetChannelWith(N1, tokenAddress).Println("AfterSettle")
		if c01new == nil {
			break
		}
	}
	if i == 10 {
		return cm.caseFailWithWrongChannelData(env.CaseName, c01.Name)
	}
	//重新创建通道,并交易
	err = N0.OpenChannel(N1.Address, tokenAddress, int64(c01.Balance), int64(c01.SettleTimeout))
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, c01.Name)
	}
	N0.SendTrans(tokenAddress, 1, N1.Address, true)
	c01new := N0.GetChannelWith(N1, tokenAddress).Println("after Reopen channel")
	if !c01new.CheckSelfBalance(c01.Balance - 1) {
		return cm.caseFail(env.CaseName)
	}
	if !c01new.CheckPartnerBalance(1) {
		return cm.caseFail(env.CaseName)
	}
	trs, err = N0.GetSentTransfers()
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, fmt.Sprintf("GetSentTransfers 2 err %s", err))
	}
	models.Logger.Printf("trs=%s", utils.StringInterface(trs, 3))
	if len(trs) != 2 {
		return cm.caseFailWithWrongChannelData(env.CaseName, "trs 2 err")
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
