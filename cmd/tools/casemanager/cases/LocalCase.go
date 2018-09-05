package cases

import (
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

// LocalCase : only for local test
func (cm *CaseManager) LocalCase() (err error) {
	env, err := models.NewTestEnv("./cases/LocalCase.ENV")
	if err != nil {
		return
	}
	defer func() {
		if env.Debug == false {
			env.KillAllRaidenNodes()
		}
	}()
	// 源数据
	transAmount := int32(20)
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N0, N1, N2, N3 := env.Nodes[0], env.Nodes[1], env.Nodes[2], env.Nodes[3]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点2，3
	N0.Start(env)
	N1.Start(env)
	N2.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "EventSendRevealSecretBefore",
	})
	N3.Start(env)
	go N0.SendTrans(tokenAddress, transAmount, N3.Address, false)
	time.Sleep(180 * time.Second)
	N2.ReStartWithoutConditionquit(env)
	time.Sleep(1000 * time.Second)
	//time.Sleep(time.Second * 30)
	// 节点2向节点3转账20token,带密码
	//secretSeed := "123"
	//N2.SendTransWithSecret(tokenAddress, transAmount, N3.Address, secretSeed)
	//// 源数据
	//var transAmount int32
	//var msg string
	//transAmount = 45
	//tokenAddress := env.Tokens[0].TokenAddress.String()
	//N2, N3, N6, N7 := env.Nodes[0], env.Nodes[1], env.Nodes[2], env.Nodes[3]
	//models.Logger.Println(env.CaseName + " BEGIN ====>")
	//// 启动节点2，6，7
	//N2.Start(env)
	//N6.Start(env)
	//N7.Start(env)
	//// 启动节点3, EventSendRefundTransferAfter
	//N3.StartWithConditionQuit(env, &params.ConditionQuit{
	//	QuitEvent: "EventSendRefundTransferAfter",
	//})
	//
	//// 节点2向节点6转账20token
	//N2.SendTrans(tokenAddress, transAmount, N6.Address, false)
	//time.Sleep(time.Second * 3)
	////  崩溃判断
	//if N3.IsRunning() {
	//	msg = "Node " + N3.Name + " should be exited,but it still running, FAILED !!!"
	//	models.Logger.Println(msg)
	//	return fmt.Errorf(msg)
	//}
	//// 崩溃后数据校验
	//models.Logger.Println("------------ Data After Crash ------------")
	//N2.GetChannelWith(N3, tokenAddress).PrintDataAfterCrash()
	//N6.GetChannelWith(N3, tokenAddress).PrintDataAfterCrash()
	//N2.GetChannelWith(N7, tokenAddress).PrintDataAfterCrash()
	//N7.GetChannelWith(N3, tokenAddress).PrintDataAfterCrash()
	//
	//// 重启节点3，自动发送之前中断的交易
	//N3.ReStartWithoutConditionquit(env)
	//time.Sleep(time.Second * 30)
	//
	//// 重启后校验
	//models.Logger.Println("------------ Data After Restart ------------")
	//cd23new := N2.GetChannelWith(N3, tokenAddress).PrintDataAfterRestart()
	//cd63new := N6.GetChannelWith(N3, tokenAddress).PrintDataAfterRestart()
	//cd27new := N2.GetChannelWith(N7, tokenAddress).PrintDataAfterRestart()
	//cd73new := N7.GetChannelWith(N3, tokenAddress).PrintDataAfterRestart()
	//
	//models.Logger.Println("------------ Data After Fail ------------")
	//// 校验通道双方相等
	//if !cd23new.CheckEqualByPartnerNode(env) || !cd63new.CheckEqualByPartnerNode(env) ||
	//	!cd27new.CheckEqualByPartnerNode(env) || !cd73new.CheckEqualByPartnerNode(env) {
	//	return cm.caseFail(env.CaseName)
	//}
	//// cd23,双锁定
	//if !cd23new.CheckLockBoth(transAmount) {
	//	return cm.caseFailWithWrongChannelData(env.CaseName, cd23new.Name)
	//}
	//// cd63，无锁定
	//if !cd63new.CheckNoLock() {
	//	return cm.caseFailWithWrongChannelData(env.CaseName, cd63new.Name)
	//}
	//// cd27，互锁
	//if !cd27new.CheckLockBoth(transAmount) {
	//	return cm.caseFailWithWrongChannelData(env.CaseName, cd27new.Name)
	//}
	//// cd37，互锁
	//if !cd73new.CheckLockBoth(transAmount) {
	//	return cm.caseFailWithWrongChannelData(env.CaseName, cd73new.Name)
	//}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
