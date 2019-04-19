package cases

import (
	"fmt"
	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/utils"
	"time"
)

// CasePMS01 :
func (cm *CaseManager) CasePMS01() (err error) {
	if !cm.RunSlow {
		return ErrorSkip
	}
	env, err := models.NewTestEnv("./cases/CasePMS01.ENV", cm.UseMatrix, cm.EthEndPoint, "CasePMS01")
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
	cm.startNodes(env, N2, N3)
	// 启动委托方节点N1
	cm.startNodesWithPMS(env, N1.SetConditionQuit(&params.ConditionQuit{
		QuitEvent: "EventSendAnnouncedDisposedResponseBefore",
	}))

	transAmount1 := int32(11)
	transAmount2 := int32(1)
	transAmount3 := int32(11)
	// 获取channel信息
	c12 := N1.GetChannelWith(N2, tokenAddress).Println("before send trans")
	c23 := N2.GetChannelWith(N3, tokenAddress).Println("before send trans")
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

	// step1:N1 send trans to N3(30token,failed)
	secret, _, err := N1.GenerateSecret()
	if err != nil {
		return
	}
	models.Logger.Println("================step1")

	N1.SendTransWithSecret(tokenAddress, transAmount1, N3.Address, secret)
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
	c21First := N2.GetChannelWith(N1, tokenAddress).Println("step1: check CD-N2-N1")
	if c21First.PartnerLockedAmount != transAmount1 {
		return fmt.Errorf("step1: check n1 LockedAmount failed,LockedAmount= %d,expect =%d", c21First.PartnerLockedAmount, transAmount1)
	}
	// N1 restart pms
	//N1.ReStartWithoutConditionquit(env)
	cm.startNodesWithPMS(env, N1.RestartName().SetConditionQuit(nil))
	time.Sleep(time.Second)
	models.Logger.Println("n1 restart")
	// step2:N1 send trans to N3(10token,success)
	models.Logger.Println("================step2")
	for i := 0; i < 3; i++ {
		N1.SendTransWithSecret(tokenAddress, 1, N3.Address, string(i)+string(time.Now().Nanosecond()))
	}
	time.Sleep(time.Second)
	var waitTime int
	for waitTime = 0; waitTime < 45; waitTime++ {
		time.Sleep(time.Second)
		c23Second := N2.GetChannelWith(N3, tokenAddress).Println("step2: check CD-N2-N3,n3 should receive 3 token")
		if c23Second.PartnerBalance == c23.PartnerBalance+transAmount2*3 {
			break
		}
	}
	if waitTime == 45 {
		return fmt.Errorf("step2: check CD-N2-N3 failed,n3Balance<>expect =%d", c23.PartnerBalance+transAmount2*3)
	}
	// step3:N1 send trans to N3(40token,failed)
	models.Logger.Println("================step3")
	N1.Shutdown(env)
	time.Sleep(time.Second)
	cm.startNodesWithPMS(env, N1.RestartName().SetConditionQuit(&params.ConditionQuit{
		QuitEvent: "EventSendAnnouncedDisposedResponseBefore",
	}))
	time.Sleep(time.Second)
	secret3, _, err := N1.GenerateSecret()
	if err != nil {
		return
	}
	models.Logger.Println("start transfer3 ....")
	N1.SendTransWithSecret(tokenAddress, transAmount3, N3.Address, secret3)
	models.Logger.Println("transfer3 finish....")
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
	time.Sleep(time.Second)
	c21Third := N2.GetChannelWith(N1, tokenAddress).Println("step3: check CD-N2-N1 ")
	if c21Third.PartnerLockedAmount != transAmount3 {
		return fmt.Errorf("step3: check n1 LockedAmount failed,LockedAmount= %d,expect =%d", c21Third.PartnerLockedAmount, transAmount3)
	}

	models.Logger.Println("================n2 force unlock")
	////N2.UpdateMeshNetworkNodes(cm.nodesExcept(env.Nodes, N1)...)
	//cm.startNodesWithPMS(env, N1.RestartName().SetConditionQuit(nil).SetNoNetwork())
	///*//n1注册两个mtr的密码
	//err=N1.RegisterSecret(secret)
	//if err!=nil{
	//	return
	//}
	//err=N1.RegisterSecret(secret3)
	//if err!=nil{
	//	return
	//}*/
	//// N2 force unlock
	//time.Sleep(time.Second)
	//N1.Shutdown(env)
	c1, err := N2.SpecifiedChannel(c12.ChannelIdentifier)
	if err != nil {
		return
	}
	models.Logger.Println(fmt.Sprintf("n2 specifiledchannel %s:", utils.StringInterface(c1, 5)))

	err = N2.ForceUnlock(c12.ChannelIdentifier, secret3)
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

	N1.ReStartWithoutConditionquit(env)
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
