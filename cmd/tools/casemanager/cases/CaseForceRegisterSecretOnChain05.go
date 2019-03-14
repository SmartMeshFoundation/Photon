package cases

import (
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/Photon/log"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

func init() {

}

// CaseForceRegisterSecretOnChain05 :
func (cm *CaseManager) CaseForceRegisterSecretOnChain05() (err error) {
	if !cm.RunSlow {
		return
	}
	env, err := models.NewTestEnv("./cases/CaseForceRegisterSecretOnChain05.ENV", cm.UseMatrix, cm.EthEndPoint)
	if err != nil {
		return
	}
	defer func() {
		log.Trace(fmt.Sprintf("CaseForceRegisterSecretOnChain05 err=%s", err))
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
	c01 := N0.GetChannelWith(N1, tokenAddress).Println("before send tras")
	N1.GetChannelWith(N2, tokenAddress).Println("before send  trans")
	n0value, err := N0.TokenBalance(tokenAddress)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, "query balance error")
	}
	n1value, err := N1.TokenBalance(tokenAddress)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, "query balance n1 error")
	}

	n0value += int(c01.Balance)
	n1value += int(c01.PartnerBalance)
	log.Trace(fmt.Sprintf("before transfer ,n0value=%d,n1value=%d", n0value, n1value))

	go N0.SendTrans(env.Tokens[0].TokenAddress.String(), 3, N2.Address, false)
	time.Sleep(3 * time.Second)
	for i := 0; i < cm.HighMediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		if !N0.IsRunning() && !N1.IsRunning() {
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
	settleTime := c01.SettleTimeout + 3600/14
	err = cm.trySettleInSeconds(int(settleTime), N1, c01.ChannelIdentifier)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, c01.Name)
	}
	n1valuenew, err := N1.TokenBalance(tokenAddress)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, "query balance n1 error")
	}
	if n1value != n1valuenew-3 {
		return cm.caseFailWithWrongChannelData(env.CaseName, fmt.Sprintf("n0=%d,n1=%d, n1new=%d", n0value, n1value, n1valuenew))
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
