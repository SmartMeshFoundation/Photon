package cases

import (
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/utils"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

func init() {

}

// CaseCannotUpdateBalanceProofAfterChannelClosed02 :
func (cm *CaseManager) CaseCannotUpdateBalanceProofAfterChannelClosed02() (err error) {
	if !cm.RunSlow {
		return
	}
	env, err := models.NewTestEnv("./cases/CaseCannotUpdateBalanceProofAfterChannelClosed02.ENV", cm.UseMatrix, cm.EthEndPoint)
	if err != nil {
		return
	}
	defer func() {
		log.Trace(fmt.Sprintf("CaseCannotUpdateBalanceProofAfterChannelClosed02 err=%s", err))
		if env.Debug == false {
			env.KillAllPhotonNodes()
		}
	}()
	// 源数据
	// original data
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N0, N1, N2 := env.Nodes[0], env.Nodes[1], env.Nodes[2]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点2，3
	// start node 2, 3
	cm.startNodes(env, N2)
	N0.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "ReceiveSecretRevealStateChange",
	})
	N1.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "EventSendRevealSecretAfter",
	})
	// 获取channel信息
	// get channel info
	if cm.UseMatrix{
		time.Sleep(time.Second*5)
	}
	c01 := N0.GetChannelWith(N1, tokenAddress).Println("before send tras")
	N1.GetChannelWith(N2, tokenAddress).Println("before send  trans")

	go N0.SendTrans(env.Tokens[0].TokenAddress.String(), 3, N2.Address, false)
	time.Sleep(3 * time.Second)
	// 崩溃判断
	// 验证n0 n1崩溃
	for i := 0; i < cm.MediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		if !N0.IsRunning() && !N1.IsRunning(){
			break
		}
	}
	if N0.IsRunning() {
		return cm.caseFailWithWrongChannelData(env.CaseName, "n0 should quit")
	}
	if N1.IsRunning() {
		return cm.caseFailWithWrongChannelData(env.CaseName, "n1 should quit")
	}
	N1.ReStartWithoutConditionquit(env)
	err = N1.Close(c01.ChannelIdentifier)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, fmt.Sprintf("close failed %s", err))
	}
	//N0务必启启动,尝试发送removeExpiredHashlock失败
	N0.ReStartWithoutConditionquit(env)
	if cm.UseMatrix{
		time.Sleep(time.Second*5)
	}
	var i = 0
	settleTime := c01.SettleTimeout + 3600/14
	for i = 0; i < int(settleTime); i++ {
		time.Sleep(time.Second)
		c, err := N1.SpecifiedChannel(c01.ChannelIdentifier)
		log.Trace(fmt.Sprintf("c=%s,err=%s", utils.StringInterface(c, 3), err))
		if err != nil {
			continue
		}
		if len(c.PartnerKnownSecretLocks) != 1 {
			continue
		}
		reg := false
		for _, s := range c.PartnerKnownSecretLocks {
			if s.IsRegisteredOnChain {
				reg = true
				break
			}
		}
		if !reg {
			continue
		}
		break
	}
	if i == int(settleTime) {
		return cm.caseFailWithWrongChannelData(env.CaseName, c01.Name)
	}

	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
