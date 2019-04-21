package cases

import (
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/utils"
)

// CasePMSPunish :
func (cm *CaseManager) CasePMSPunish() (err error) {
	if !cm.RunSlow {
		return ErrorSkip
	}
	env, err := models.NewTestEnv("./cases/CasePMSPunish.ENV", cm.UseMatrix, cm.EthEndPoint, "CasePMSPunish")
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

	N1, N2, N3 := env.Nodes[1], env.Nodes[2], env.Nodes[3]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动pms
	env.StartPMS()
	// 启动节点2、3
	cm.startNodes(env, N2, N3, N1.SetConditionQuit(&params.ConditionQuit{
		QuitEvent: "EventSendAnnouncedDisposedResponseBefore",
	}))
	transAmount := int32(30)
	secret, _, err := N1.GenerateSecret()
	if err != nil {
		return
	}
	// 获取channel信息
	c12 := N1.GetChannelWith(N2, tokenAddress).Println("before send trans")
	n1value, err := N1.TokenBalance(tokenAddress)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, "query balance n1 error")
	}
	n2value, err := N2.TokenBalance(tokenAddress)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, "query balance n2 error")
	}

	n1value += int(c12.Balance)
	n2value += int(c12.PartnerBalance)
	log.Trace(fmt.Sprintf("before transfer ,n1value=%d,n2value=%d", n1value, n2value))

	// N1 send trans to N3
	N1.SendTransWithSecret(tokenAddress, transAmount, N3.Address, secret)
	// 崩溃判断
	for i := 0; i < cm.MediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		if !N1.IsRunning() {
			break
		}
	}
	if N1.IsRunning() {
		msg := "Node " + N1.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	//check lock
	c12lock := N2.GetChannelWith(N1, tokenAddress).Println("check CD-N2-N1 should has a lock")
	if c12lock.PartnerLockedAmount != transAmount {
		return fmt.Errorf("check n1 lock failed,lockAmount= %d,expect =%d", c12lock.PartnerLockedAmount, transAmount)
	}

	// n2 shut down
	N2.Shutdown(env)
	models.Logger.Println("n2 shutdown")
	// N1 restart pms
	cm.startNodes(env, N1.RestartName().PMS())
	models.Logger.Println("n1 restart with pms")
	//time.Sleep(time.Second * 2) //wait for submit to pms
	c1, err := N1.SpecifiedChannel(c12.ChannelIdentifier)
	if err != nil {
		return
	}
	models.Logger.Println(fmt.Sprintf("n1 specifiledchannel %s:", utils.StringInterface(c1, 5)))
	if c1.DelegateStateString != "delegate success" {
		return fmt.Errorf("n1 delegate failed")
	}
	//N1 shut down
	N1.Shutdown(env)
	models.Logger.Println("n1 shutdown after delegate CD-N1-N2")

	// n2 restart
	N2.RestartName().Start(env)
	models.Logger.Println("n2 restart to force unlock")

	// N2 force unlock channel
	err = N2.ForceUnlock(c12.ChannelIdentifier, secret)
	if err != nil {
		return fmt.Errorf("n2 force unlock err %s", err)
	}

	//pms should punish n2
	settleTime := c12.SettleTimeout + 257
	err = cm.trySettleInSeconds(int(settleTime), N2, c12.ChannelIdentifier)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, c12.Name)
	}
	models.Logger.Printf("n2 settle finished...")

	N1.RestartName().ReStartWithoutConditionquit(env)
	var n1NewValue, n2NewValue, i int
	for i = 0; i < 5; i++ {
		time.Sleep(time.Second * 1)
		n1NewValue, err = N1.TokenBalance(tokenAddress)
		n2NewValue, err = N2.TokenBalance(tokenAddress)

		if n1NewValue != n1value+100 {
			continue
		}
		if n2NewValue != n2value-100 {
			continue
		}
		if N1.GetChannelWith(N2, tokenAddress) != nil {
			continue
		}
		models.Logger.Println("punish n2 success")
		models.Logger.Println(env.CaseName + " END ====> SUCCESS")
		return
	}
	return cm.caseFailWithWrongChannelData(env.CaseName, fmt.Sprintf("punish n2 failed, error n1=%d,n1expect=%d,n2=%d,n2expect=%d ", n1NewValue, n1value+100, n2NewValue, n2value-100))
}
