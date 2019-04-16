package cases

import (
	"fmt"
	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/params"
	"time"
)

// CasePMSNoPunish :
func (cm *CaseManager) CasePMSNoPunish() (err error) {
	env, err := models.NewTestEnv("./cases/CasePMSNoPunish.ENV", cm.UseMatrix, cm.EthEndPoint, "CasePMSNoPunish")
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
	cm.startNodesWithPMS(env, N2)
	cm.startNodes(env, N3)
	// 启动委托节点N1
	cm.startNodes(env, N1.SetConditionQuit(&params.ConditionQuit{
		QuitEvent: "EventSendAnnouncedDisposedResponseBefore",
	}))

	transAmount := int32(30)
	secret, _, err := N1.GenerateSecret()
	if err != nil {
		return
	}
	// 获取channel信息
	c12 := N1.GetChannelWith(N2, tokenAddress).Println("before send trans")
	//c23 := N2.GetChannelWith(N3, tokenAddress).Println("before send trans")
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

	//n2 exit
	N2.Shutdown(env)
	models.Logger.Println("n2 shutdown")

	// N1 restart with pms
	N1.StartWithPMS(env)
	if err != nil {
		return
	}
	time.Sleep(time.Second)
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
	//N1注册密码
	err = N1.RegisterSecret(secret)
	if err != nil {
		return fmt.Errorf("n1 register secret err %s", err)
	}
	models.Logger.Printf("n1 register secret on chain finished...")
	// N1 close channel
	err = N1.Close(c12.ChannelIdentifier)
	if err != nil {
		return
	}
	models.Logger.Printf("n1 close channel finished...")
	// N1 settle
	settleTime := c12.SettleTimeout + 257
	err = cm.trySettleInSeconds(int(settleTime), N1, c12.ChannelIdentifier)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, c12.Name)
	}
	models.Logger.Printf("n1 settle finished...")
	//check n2's tokenBalance
	N2.ReStartWithoutConditionquit(env)
	time.Sleep(time.Second)
	var n2NewValue, i int
	for i = 0; i < 10; i++ {
		time.Sleep(time.Second * 1)
		n2NewValue, err = N2.TokenBalance(tokenAddress)
		if n2NewValue != n2value {
			continue
		}
		models.Logger.Println(fmt.Sprintf("n2NewValue=%dn2Value=%d", n2NewValue, n2value))
		models.Logger.Println(env.CaseName + " END ====> SUCCESS")
		return
	}
	return cm.caseFailWithWrongChannelData(env.CaseName, fmt.Sprintf("check TokenBalance error n2=%d,n2expect=%d ", n2NewValue, n2value))
}
