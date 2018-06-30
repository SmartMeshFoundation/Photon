package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

// CrashCaseSend05 场景五：EventSendRefundTransferAfter
// 发送refundtransfer交易崩溃
// 节点2发送45token给节点6 ，发送refundtransfer后节点3崩，节点2锁定45，其余节点无锁定;
// 重启节点3后，节点2，3各锁定 45，节点2，4、节点4、5，节点5、6交易成功，转账成功，但节点2、3各锁定45token
func (cm *CaseManager) CrashCaseSend05() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCaseSend05.ENV")
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
	// 启动节点2，4,5,6
	N2.Start(env)
	N4.Start(env)
	N5.Start(env)
	N6.Start(env)
	// 启动节点3, EventSendRefundTransferAfter
	N3.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "EventSendRefundTransferAfter",
	})

	// 节点2向节点6转账20token
	N2.SendTrans(tokenAddress, transAmount, N6.Address, false)
	time.Sleep(time.Second * 3)
	//  崩溃判断
	if N3.IsRunning() {
		msg = "Node " + N3.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}

	// 查询cd23，锁定45,其余无锁定
	models.Logger.Println("------------ Data After Crash ------------")
	N2.GetChannelWith(N3, tokenAddress).PrintDataAfterCrash()
	cd24middle := N2.GetChannelWith(N4, tokenAddress).PrintDataAfterCrash()
	N6.GetChannelWith(N3, tokenAddress).PrintDataAfterCrash()
	cd45middle := N4.GetChannelWith(N5, tokenAddress).PrintDataAfterCrash()
	cd56middle := N5.GetChannelWith(N6, tokenAddress).PrintDataAfterCrash()

	// 重启节点3，自动发送之前中断的交易
	N3.ReStartWithoutConditionquit(env)
	time.Sleep(time.Second * 30)

	// 查询重启后数据
	models.Logger.Println("------------ Data After Restart ------------")
	cd23new := N2.GetChannelWith(N3, tokenAddress).PrintDataAfterRestart()
	cd24new := N2.GetChannelWith(N4, tokenAddress).PrintDataAfterRestart()
	cd63new := N6.GetChannelWith(N3, tokenAddress).PrintDataAfterRestart()
	cd45new := N4.GetChannelWith(N5, tokenAddress).PrintDataAfterRestart()
	cd56new := N5.GetChannelWith(N6, tokenAddress).PrintDataAfterRestart()

	// 校验对等
	models.Logger.Println("------------ Data After Fail ------------")
	if !cd23new.CheckEqualByPartnerNode(env) || !cd24new.CheckEqualByPartnerNode(env) ||
		!cd63new.CheckEqualByPartnerNode(env) || !cd45new.CheckEqualByPartnerNode(env) || !cd56new.CheckEqualByPartnerNode(env) {
		return cm.caseFail(env.CaseName)
	}
	// cd23,双锁定
	if !cd23new.CheckLockBoth(transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd23new.Name)
	}
	// cd63,无锁
	if !cd63new.CheckNoLock() {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd63new.Name)
	}
	// cd24,交易成功
	if !cd24new.CheckPartnerBalance(cd24middle.PartnerBalance + transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd24new.Name)
	}
	// cd45,交易成功
	if !cd45new.CheckPartnerBalance(cd45middle.PartnerBalance + transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd45new.Name)
	}
	// cd56,交易成功
	if !cd56new.CheckPartnerBalance(cd56middle.PartnerBalance + transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd56new.Name)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
