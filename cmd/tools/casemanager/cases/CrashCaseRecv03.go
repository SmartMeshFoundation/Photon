package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/utils"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

// CrashCaseRecv03 场景三：ReceiveTransferRefundStateChange
//（收到refundtransfer崩）
// 节点1向节点6发送45个token，（提前进行两次转账，降低部分余额，新余额分配为节点3和节点6 余额：30， 320；节点3和节点7余额： 30 90），
// 因此，节点3要回退节点2，节点2崩；节点1锁定45，节点2，节点3锁定45，节点6未锁定；重启节点2后，重启转账成功，锁定token解锁。
func (cm *CaseManager) CrashCaseRecv03() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCaseRecv03.ENV")
	if err != nil {
		return
	}
	defer env.KillAllRaidenNodes()
	// 源数据
	var transAmount int32
	var msg string
	transAmount = 45
	tokenAddress := env.Tokens[0].Address
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
		QuitEvent: "ReceiveTransferRefundStateChange",
	})

	// 2. 记录所有通道历史数据
	cd12 := utils.GetChannelBetween(N1, N2, tokenAddress).PrintDataBeforeTransfer()
	cd32 := utils.GetChannelBetween(N3, N2, tokenAddress).PrintDataBeforeTransfer()
	cd42 := utils.GetChannelBetween(N4, N2, tokenAddress).PrintDataBeforeTransfer()
	cd36 := utils.GetChannelBetween(N3, N6, tokenAddress).PrintDataBeforeTransfer()
	cd45 := utils.GetChannelBetween(N4, N5, tokenAddress).PrintDataBeforeTransfer()
	cd56 := utils.GetChannelBetween(N5, N6, tokenAddress).PrintDataBeforeTransfer()

	// 3. 节点1向节点6转账45token
	N1.SendTrans(tokenAddress, transAmount, N6.Address, false)
	time.Sleep(time.Second * 3)
	// 4. 崩溃判断
	if N2.IsRunning() {
		msg = "Node " + N2.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 5. 崩溃后重启前数据校验
	cd12middle := utils.GetChannelBetween(N1, N2, tokenAddress).PrintDataAfterCrash()
	cd32middle := utils.GetChannelBetween(N3, N2, tokenAddress).PrintDataAfterCrash()
	cd42middle := utils.GetChannelBetween(N4, N2, tokenAddress).PrintDataAfterCrash()
	cd36middle := utils.GetChannelBetween(N3, N6, tokenAddress).PrintDataAfterCrash()
	cd45middle := utils.GetChannelBetween(N4, N5, tokenAddress).PrintDataAfterCrash()
	cd56middle := utils.GetChannelBetween(N5, N6, tokenAddress).PrintDataAfterCrash()
	// 校验cd12，锁定45
	if cd12middle.LockedAmount != transAmount {
		msg = fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", transAmount, cd12middle.LockedAmount)
		return cm.caseFail(env.CaseName, msg)
	}
	// 校验cd32，锁定对方45
	if cd32middle.PartnerLockedAmount != transAmount {
		msg = fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", transAmount, cd32middle.PartnerLockedAmount)
		return cm.caseFail(env.CaseName, msg)
	}
	// 校验cd42，无锁定
	if cd42middle.PartnerLockedAmount != 0 {
		msg = fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", 0, cd42middle.PartnerLockedAmount)
		return cm.caseFail(env.CaseName, msg)
	}
	// 校验cd36，无锁定
	if cd36middle.LockedAmount != 0 {
		msg = fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", 0, cd36middle.LockedAmount)
		return cm.caseFail(env.CaseName, msg)
	}
	// 校验cd45，无锁定
	if cd45middle.LockedAmount != 0 {
		msg = fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", 0, cd45middle.LockedAmount)
		return cm.caseFail(env.CaseName, msg)
	}
	// 校验cd56，无锁定
	if cd56middle.LockedAmount != 0 {
		msg = fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", 0, cd56middle.LockedAmount)
		return cm.caseFail(env.CaseName, msg)
	}

	// 6. 重启节点2，交易自动继续
	N2.ReStartWithoutConditionquit(env)

	time.Sleep(time.Second * 30)
	// 7. 重启后数据校验，交易成功，全通道无锁定
	cd12new := utils.GetChannelBetween(N1, N2, tokenAddress).PrintDataAfterRestart()
	cd32new := utils.GetChannelBetween(N3, N2, tokenAddress).PrintDataAfterRestart()
	cd42new := utils.GetChannelBetween(N4, N2, tokenAddress).PrintDataAfterRestart()
	cd36new := utils.GetChannelBetween(N3, N6, tokenAddress).PrintDataAfterRestart()
	cd45new := utils.GetChannelBetween(N4, N5, tokenAddress).PrintDataAfterRestart()
	cd56new := utils.GetChannelBetween(N5, N6, tokenAddress).PrintDataAfterRestart()
	// 校验cd12
	if cd12.Balance-cd12new.Balance != transAmount {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd12new.Name)
	}
	// 校验cd32,解锁
	if cd32new.Balance != cd32.Balance || cd32new.PartnerLockedAmount != 0 {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd32new.Name)
	}
	// 校验cd42
	if cd42.PartnerBalance-cd42new.PartnerBalance != transAmount {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd42new.Name)
	}
	// 校验cd36
	if cd36new.Balance != cd36.Balance || cd36new.LockedAmount != 0 {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd36new.Name)
	}
	// 校验cd45
	if cd45.Balance-cd45new.Balance != transAmount {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd45new.Name)
	}
	// 校验cd56
	if cd56.Balance-cd56new.Balance != transAmount {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd56new.Name)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
