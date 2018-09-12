package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

// CrashCaseRecv06 场景六：BeforeSendRevealSecret
// （发送secret之前）
// 节点1向节点6发送20 token,节点1崩。节点1锁定20 token, 节点2 成功，节点3成功，
// 节点6成功。重启节点1后，节点1锁定解锁，转账成功。
// 此种情况下，转账继续，不影响使用。
func (cm *CaseManager) CrashCaseRecv06() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCaseRecv06.ENV")
	if err != nil {
		return
	}
	defer func() {
		if env.Debug == false {
			env.KillAllRaidenNodes()
		}
	}()
	// 源数据
	var transAmount int32
	var msg string
	transAmount = 20
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N1, N2, N3, N6 := env.Nodes[0], env.Nodes[1], env.Nodes[2], env.Nodes[3]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 1. 启动
	// 启动节点2,3,6
	N2.Start(env)
	N3.Start(env)
	N6.Start(env)
	// 启动节点1, BeforeSendRevealSecret
	N1.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "EventSendRevealSecretBefore",
	})
	// 初始数据记录
	cd21 := N2.GetChannelWith(N1, tokenAddress).PrintDataBeforeTransfer()
	cd23 := N2.GetChannelWith(N3, tokenAddress).PrintDataBeforeTransfer()
	cd36 := N3.GetChannelWith(N6, tokenAddress).PrintDataBeforeTransfer()
	// 3. 节点1向节点6转账20token
	N1.SendTrans(tokenAddress, transAmount, N6.Address, false)
	time.Sleep(time.Second * 3)
	// 4. 崩溃判断
	if N1.IsRunning() {
		msg = "Node " + N1.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 6. 中间数据记录
	models.Logger.Println("------------ Data After Crash ------------")
	cd21middle := N2.GetChannelWith(N1, tokenAddress).PrintDataAfterCrash()
	cd23middle := N2.GetChannelWith(N3, tokenAddress).PrintDataAfterCrash()
	cd36middle := N3.GetChannelWith(N6, tokenAddress).PrintDataAfterCrash()
	// 校验cd12，锁定20
	if !cd21middle.CheckLockPartner(transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd21middle.Name)
	}
	// 校验cd23，交易成功
	if !cd23middle.CheckPartnerBalance(cd23.PartnerBalance + transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd23middle.Name)
	}
	// 校验cd36，交易成功
	if !cd36middle.CheckPartnerBalance(cd36.PartnerBalance + transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd36middle.Name)
	}

	// 6. 重启节点2，交易自动继续
	N1.ReStartWithoutConditionquit(env)
	time.Sleep(time.Second * 15)

	// 查询重启后数据
	models.Logger.Println("------------ Data After Restart ------------")
	cd21new := N2.GetChannelWith(N1, tokenAddress).PrintDataAfterRestart()
	N2.GetChannelWith(N3, tokenAddress).PrintDataAfterRestart()
	N3.GetChannelWith(N6, tokenAddress).PrintDataAfterRestart()

	// 校验对等
	models.Logger.Println("------------ Data After Fail ------------")
	if !cd21new.CheckEqualByPartnerNode(env) {
		return cm.caseFail(env.CaseName)
	}
	// 校验cd21, 交易成功
	if !cd21new.CheckSelfBalance(cd21.Balance + transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd21new.Name)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
