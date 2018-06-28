package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/utils"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

// CrashCaseRecv04 场景四：ReceiveTransferRefundStateChange
//（收到refundtransfer崩）
// 节点1向节点6发送45个token，（提前进行两次转账，降低部分余额，节点3和节点7余额： 30 90），
// 因此，节点3要回退节点2，节点2崩；节点1锁定45，节点2，节点3锁定45，节点6未锁定；重启节点2后，重启转账失败，cd12,23,27全锁定，cd36无锁定
func (cm *CaseManager) CrashCaseRecv04() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCaseRecv04.ENV")
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
	N1, N2, N3, N6, N7 := env.Nodes[0], env.Nodes[1], env.Nodes[2], env.Nodes[3], env.Nodes[4]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 1. 启动
	// 启动节点1,3,6,7
	N1.Start(env)
	N3.Start(env)
	N6.Start(env)
	N7.Start(env)
	// 启动节点2, ReceiveTransferRefundStateChange
	N2.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "ReceiveTransferRefundStateChange",
	})

	// 2. 记录所有通道历史数据
	utils.GetChannelBetween(N1, N2, tokenAddress).PrintDataBeforeTransfer()
	utils.GetChannelBetween(N3, N2, tokenAddress).PrintDataBeforeTransfer()
	utils.GetChannelBetween(N3, N6, tokenAddress).PrintDataBeforeTransfer()
	utils.GetChannelBetween(N7, N2, tokenAddress).PrintDataBeforeTransfer()
	utils.GetChannelBetween(N7, N3, tokenAddress).PrintDataBeforeTransfer()

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
	cd36middle := utils.GetChannelBetween(N3, N6, tokenAddress).PrintDataAfterCrash()
	cd72middle := utils.GetChannelBetween(N7, N2, tokenAddress).PrintDataAfterCrash()
	cd73middle := utils.GetChannelBetween(N7, N3, tokenAddress).PrintDataAfterCrash()
	// 校验cd12，锁定45
	if cd12middle.LockedAmount != transAmount {
		msg = fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", transAmount, cd12middle.LockedAmount)
		return cm.caseFail(env.CaseName, msg)
	}
	// 校验cd32，双锁定45
	if cd32middle.PartnerLockedAmount != transAmount && cd32middle.LockedAmount != transAmount {
		msg = fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", transAmount, cd32middle.PartnerLockedAmount)
		return cm.caseFail(env.CaseName, msg)
	}
	// 校验cd36，无锁定
	if cd36middle.LockedAmount != 0 {
		msg = fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", 0, cd36middle.LockedAmount)
		return cm.caseFail(env.CaseName, msg)
	}
	// 校验cd72，无锁定
	if cd72middle.PartnerLockedAmount != 0 {
		msg = fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", 0, cd72middle.PartnerLockedAmount)
		return cm.caseFail(env.CaseName, msg)
	}
	// 校验cd73，无锁定
	if cd73middle.LockedAmount != 0 {
		msg = fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", 0, cd73middle.LockedAmount)
		return cm.caseFail(env.CaseName, msg)
	}

	// 6. 重启节点2，交易自动继续
	N2.ReStartWithoutConditionquit(env)

	time.Sleep(time.Second * 30)
	// 7. 重启后数据校验，转账失败
	cd12new := utils.GetChannelBetween(N1, N2, tokenAddress).PrintDataAfterRestart()
	cd32new := utils.GetChannelBetween(N3, N2, tokenAddress).PrintDataAfterRestart()
	cd36new := utils.GetChannelBetween(N3, N6, tokenAddress).PrintDataAfterRestart()
	cd72new := utils.GetChannelBetween(N7, N2, tokenAddress).PrintDataAfterRestart()
	cd73new := utils.GetChannelBetween(N7, N3, tokenAddress).PrintDataAfterRestart()
	// 校验cd12, 1锁定
	if cd12new.LockedAmount != transAmount {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd12new.Name)
	}
	// 校验cd32, 双锁定
	if cd32new.PartnerLockedAmount != transAmount && cd32new.LockedAmount != transAmount {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd32new.Name)
	}
	// 校验cd36，无锁定
	if cd36new.LockedAmount != 0 {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd36new.Name)
	}
	// 校验cd72,2锁定
	if cd72new.PartnerLockedAmount != transAmount {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd72new.Name)
	}
	// 校验cd72,7锁定45
	if cd73new.LockedAmount != transAmount {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd73new.Name)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
