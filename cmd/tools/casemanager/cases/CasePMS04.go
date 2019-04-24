package cases

import (
	"time"

	"fmt"
	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

// CasePMS04 :
func (cm *CaseManager) CasePMS04() (err error) {
	if !cm.RunSlow {
		return ErrorSkip
	}
	env, err := models.NewTestEnv("./cases/CasePMS04.ENV", cm.UseMatrix, cm.EthEndPoint, "CasePMS04")
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
	// 启动受托节点n1-pms
	env.StartPMS()
	cm.startNodes(env, N1.SetConditionQuit(&params.ConditionQuit{
		QuitEvent: "ReceiveAnnounceDisposedStateChange",
	}), N2.PMS(), N3.SetConditionQuit(&params.ConditionQuit{
		QuitEvent: "ReceiveAnnounceDisposedStateChange",
	}))

	transAmount := int32(20)
	// 获取channel信息
	c12 := N1.GetChannelWith(N2, tokenAddress).Println("before send trans")
	c32 := N3.GetChannelWith(N2, tokenAddress).Println("before send trans")
	n1TokenOld, err := N1.TokenBalance(tokenAddress)
	n2TokenOld, err := N2.TokenBalance(tokenAddress)
	n3TokenOld, err := N3.TokenBalance(tokenAddress)
	if err != nil {
		return
	}
	n1TokenOld += int(c12.Balance)
	n2TokenOld += int(c12.PartnerBalance)
	n2TokenOld += int(c32.PartnerBalance)
	n3TokenOld += int(c32.Balance)
	/*
		# 路由：a-b-c，b委托pms,
		# a和c同时发起mtr,c在收到reveal secret后，c和a同时掉线，然后b掉线
		# a和c上线到链上注册密码，然后a unlock通道a-b ,c unlock 通道c-b
		# 测试：两笔交易应均失败
	*/
	secret1, _, err := N1.GenerateSecret()
	if err != nil {
		return
	}
	secret3, _, err := N3.GenerateSecret()
	if err != nil {
		return
	}
	go N1.SendTransWithSecret(tokenAddress, transAmount, N3.Address, secret1)
	go N3.SendTransWithSecret(tokenAddress, transAmount, N1.Address, secret3)
	// 崩溃判断
	for i := 0; i < cm.MediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		if !N1.IsRunning() && !N3.IsRunning() {
			break
		}
	}
	if N1.IsRunning() || N3.IsRunning() {
		msg := "Node " + N1.Name + " " + N3.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	N2.Shutdown(env)
	cm.startNodes(env, N1.RestartName().SetConditionQuit(nil), N3.RestartName().SetConditionQuit(nil))
	time.Sleep(time.Second * 2)
	N1.GetChannelWith(N2, tokenAddress).Println("after send trans")
	N3.GetChannelWith(N2, tokenAddress).Println("after send trans")

	//N1注册密码
	err = N1.RegisterSecret(secret1)
	if err != nil {
		return fmt.Errorf("n1 register secret err %s", err)
	}
	models.Logger.Printf("n1 register secret on chain finished...")
	//N3注册密码
	err = N3.RegisterSecret(secret3)
	if err != nil {
		return fmt.Errorf("n3 register secret err %s", err)
	}
	models.Logger.Printf("n3 register secret on chain finished...")

	err = N1.Close(c12.ChannelIdentifier)
	err = N3.Close(c32.ChannelIdentifier)
	settleTime := c12.SettleTimeout + 257
	err = cm.trySettleInSeconds(int(settleTime), N1, c12.ChannelIdentifier)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, c12.Name)
	}
	models.Logger.Printf("n1 settle finished...")

	err = cm.trySettleInSeconds(int(settleTime), N3, c32.ChannelIdentifier)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, c32.Name)
	}
	models.Logger.Printf("n3 settle finished...")
	N2.RestartName().Start(env)
	time.Sleep(time.Second)
	n1NewValue, err := N1.TokenBalance(tokenAddress)
	n2NewValue, err := N2.TokenBalance(tokenAddress)
	n3NewValue, err := N3.TokenBalance(tokenAddress)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}

	if n1NewValue != n1TokenOld && n2NewValue != n2TokenOld && n3NewValue != n3TokenOld {
		return cm.caseFailWithWrongChannelData(env.CaseName, fmt.Sprintf("check balance on chain err,n1=%d,expect=%d,n2=%d,expect=%d,n3=%d,expect=%d",
			n1NewValue, n1TokenOld, n2NewValue, n2TokenOld, n3NewValue != n3TokenOld))
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
