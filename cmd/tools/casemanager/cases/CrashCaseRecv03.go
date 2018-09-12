package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

// CrashCaseRecv03 场景三：ReceiveTransferRefundStateChange
//（收到refundtransfer崩）
// 节点1向节点6发送45个token，（提前进行两次转账，降低部分余额，新余额分配为节点3和节点6 余额：30， 320），
// 因此，节点3要回退节点2，节点2崩；节点1锁定45，节点2，节点3锁定45，节点6未锁定；重启节点2后，重启转账成功，锁定token解锁。
func (cm *CaseManager) CrashCaseRecv03() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCaseRecv03.ENV")
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
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N1, N2, N3, N4, N5, N6 := env.Nodes[0], env.Nodes[1], env.Nodes[2], env.Nodes[3], env.Nodes[4], env.Nodes[5]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 1. 启动
	// 启动节点1,3,4,5,6,7
	N1.Start(env)
	N3.Start(env)
	N4.Start(env)
	N5.Start(env)
	N6.Start(env)
	// 启动节点2, ReceiveTransferRefundStateChange
	N2.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "ReceiveAnnounceDisposedStateChange",
	})

	// 2. 记录所有通道历史数据
	N1.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	N3.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	N4.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	N3.GetChannelWith(N6, tokenAddress).PrintDataBeforeTransfer()
	N4.GetChannelWith(N5, tokenAddress).PrintDataBeforeTransfer()
	N5.GetChannelWith(N6, tokenAddress).PrintDataBeforeTransfer()

	// 3. 节点1向节点6转账45token
	go N1.SendTrans(tokenAddress, transAmount, N6.Address, false)
	time.Sleep(time.Second * 3)
	// 4. 崩溃判断
	if N2.IsRunning() {
		msg = "Node " + N2.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 中间数据记录
	models.Logger.Println("------------ Data After Crash ------------")
	cd12middle := N1.GetChannelWith(N2, tokenAddress).PrintDataAfterCrash()
	cd32middle := N3.GetChannelWith(N2, tokenAddress).PrintDataAfterCrash()
	cd42middle := N4.GetChannelWith(N2, tokenAddress).PrintDataAfterCrash()
	cd36middle := N3.GetChannelWith(N6, tokenAddress).PrintDataAfterCrash()
	cd45middle := N4.GetChannelWith(N5, tokenAddress).PrintDataAfterCrash()
	cd56middle := N5.GetChannelWith(N6, tokenAddress).PrintDataAfterCrash()
	// 校验cd12，锁定45
	if !cd12middle.CheckLockSelf(transAmount) {
		return cm.caseFail(env.CaseName)
	}
	// 校验cd32，锁定对方45
	if !cd32middle.CheckLockPartner(transAmount) {
		return cm.caseFail(env.CaseName)
	}
	// 校验cd42，无锁定
	if !cd42middle.CheckNoLock() {
		return cm.caseFail(env.CaseName)
	}
	// 校验cd36，无锁定
	if !cd36middle.CheckNoLock() {
		return cm.caseFail(env.CaseName)
	}
	// 校验cd45，无锁定
	if !cd45middle.CheckNoLock() {
		return cm.caseFail(env.CaseName)
	}
	// 校验cd56，无锁定
	if !cd56middle.CheckNoLock() {
		return cm.caseFail(env.CaseName)
	}

	// 6. 重启节点2，交易失败
	N2.ReStartWithoutConditionquit(env)
	time.Sleep(time.Second * 10)

	// 查询重启后数据
	models.Logger.Println("------------ Data After Restart ------------")
	cd12new := N1.GetChannelWith(N2, tokenAddress).PrintDataAfterRestart()
	cd32new := N3.GetChannelWith(N2, tokenAddress).PrintDataAfterRestart()
	cd42new := N4.GetChannelWith(N2, tokenAddress).PrintDataAfterRestart()
	cd36new := N3.GetChannelWith(N6, tokenAddress).PrintDataAfterRestart()
	cd45new := N4.GetChannelWith(N5, tokenAddress).PrintDataAfterRestart()
	cd56new := N5.GetChannelWith(N6, tokenAddress).PrintDataAfterRestart()

	// 校验对等
	models.Logger.Println("------------ Data After Fail ------------")
	if !cd12new.CheckEqualByPartnerNode(env) || !cd32new.CheckEqualByPartnerNode(env) ||
		!cd42new.CheckEqualByPartnerNode(env) || !cd36new.CheckEqualByPartnerNode(env) ||
		!cd45new.CheckEqualByPartnerNode(env) || !cd56new.CheckEqualByPartnerNode(env) {
		return cm.caseFail(env.CaseName)
	}
	// 校验cd12, 锁定45
	if !cd12new.CheckLockSelf(transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd12new.Name)
	}
	// 校验cd32,解锁
	if !cd32new.CheckNoLock() {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd32new.Name)
	}
	// 校验cd42,交易成功
	if !cd42new.CheckNoLock() {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd42new.Name)
	}
	// 校验cd36,解锁
	if !cd36new.CheckNoLock() {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd36new.Name)
	}
	// 校验cd45, 交易成功
	if !cd45new.CheckNoLock() {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd45new.Name)
	}
	// 校验cd56, 交易成功
	if !cd56new.CheckNoLock() {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd56new.Name)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
