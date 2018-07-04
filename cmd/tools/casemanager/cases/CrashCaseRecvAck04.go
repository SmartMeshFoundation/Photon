package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

// CrashCaseRecvAck04 场景四：RefundTransferRecevieAck
//节点2向节点6发送45个token,节点3崩。节点2、节点3各锁定45，走路由2，4，5，6成功；
//转账成功;重启后，2，3节点锁定45未解除。未完成。正常。
func (cm *CaseManager) CrashCaseRecvAck04() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCaseRecvAck04.ENV")
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
	transAmount = 45
	tokenAddress := env.Tokens[0].Address
	N2, N3, N4, N5, N6 := env.Nodes[0], env.Nodes[1], env.Nodes[2], env.Nodes[3], env.Nodes[4]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 1. 启动
	// 启动节点2,4,5,6
	N2.Start(env)
	N4.Start(env)
	N5.Start(env)
	N6.Start(env)
	// 启动节点3, RefundTransferRecevieAck
	N3.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "RefundTransferRecevieAck",
	})
	// 初始数据记录
	N2.GetChannelWith(N3, tokenAddress).PrintDataBeforeTransfer()
	N6.GetChannelWith(N3, tokenAddress).PrintDataBeforeTransfer()
	cd24 := N2.GetChannelWith(N4, tokenAddress).PrintDataBeforeTransfer()
	cd45 := N4.GetChannelWith(N5, tokenAddress).PrintDataBeforeTransfer()
	cd56 := N5.GetChannelWith(N6, tokenAddress).PrintDataBeforeTransfer()
	// 3. 节点2向节点6转账
	N2.SendTrans(tokenAddress, transAmount, N6.Address, false)
	time.Sleep(time.Second * 5)
	// 4. 崩溃判断
	if N3.IsRunning() {
		msg = "Node " + N3.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 6. 中间数据记录
	models.Logger.Println("------------ Data After Crash ------------")
	cd23middle := N2.GetChannelWith(N3, tokenAddress).PrintDataAfterCrash()
	cd63middle := N6.GetChannelWith(N3, tokenAddress).PrintDataAfterCrash()
	cd24middle := N2.GetChannelWith(N4, tokenAddress).PrintDataAfterCrash()
	cd45middle := N4.GetChannelWith(N5, tokenAddress).PrintDataAfterCrash()
	cd56middle := N5.GetChannelWith(N6, tokenAddress).PrintDataAfterCrash()

	// 校验cd23, 双锁定
	if !cd23middle.CheckLockBoth(transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd23middle.Name)
	}
	// 校验cd36，无锁定
	if !cd63middle.CheckNoLock() {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd63middle.Name)
	}
	// 校验cd24，交易成功
	if !cd24middle.CheckPartnerBalance(cd24.PartnerBalance + transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd24middle.Name)
	}
	// 校验cd45，交易成功
	if !cd45middle.CheckPartnerBalance(cd45.PartnerBalance + transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd45middle.Name)
	}
	// 校验cd56，交易成功
	if !cd56middle.CheckPartnerBalance(cd56.PartnerBalance + transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd56middle.Name)
	}

	// 6. 重启节点3，交易自动继续
	N3.ReStartWithoutConditionquit(env)
	time.Sleep(time.Second * 15)

	// 查询重启后数据
	models.Logger.Println("------------ Data After Restart ------------")
	cd23new := N2.GetChannelWith(N3, tokenAddress).PrintDataAfterRestart()
	cd63new := N6.GetChannelWith(N3, tokenAddress).PrintDataAfterRestart()
	cd24new := N2.GetChannelWith(N4, tokenAddress).PrintDataAfterRestart()
	cd45new := N4.GetChannelWith(N5, tokenAddress).PrintDataAfterRestart()
	cd56new := N5.GetChannelWith(N6, tokenAddress).PrintDataAfterRestart()

	// 校验对等
	models.Logger.Println("------------ Data After Fail ------------")
	if !cd23new.CheckEqualByPartnerNode(env) || !cd63new.CheckEqualByPartnerNode(env) || !cd24new.CheckEqualByPartnerNode(env) ||
		!cd45new.CheckEqualByPartnerNode(env) || !cd56new.CheckEqualByPartnerNode(env) {
		return cm.caseFail(env.CaseName)
	}

	// 校验cd23, 双锁定
	if !cd23new.CheckLockBoth(transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd23new.Name)
	}
	// 校验cd36，无锁定
	if !cd63new.CheckNoLock() {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd63new.Name)
	}
	// 校验cd24，交易成功
	if !cd24new.CheckPartnerBalance(cd24.PartnerBalance + transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd24new.Name)
	}
	// 校验cd45，交易成功
	if !cd45new.CheckPartnerBalance(cd45.PartnerBalance + transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd45new.Name)
	}
	// 校验cd56，交易成功
	if !cd56new.CheckPartnerBalance(cd56.PartnerBalance + transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd56new.Name)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
