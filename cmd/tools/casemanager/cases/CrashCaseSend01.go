package cases

import (
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

// CrashCaseSend01 场景一：EventSendMediatedTransferAfter
// 发送中转转账后崩溃
// 节点1向节点2发送MTR后，节点1崩溃，此时，节点2默认收到MTR，但由于没有ACK确认，没发生转账，余额不变。节点2没收到转账token.
// 重启节点1后，继续转账，转账成功。
func (cm *CaseManager) CrashCaseSend01() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCaseSend01.ENV")
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
	transAmount = 5
	tokenAddress := env.Tokens[0].Address
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	N1, N2 := env.Nodes[0], env.Nodes[1]

	// 1. /启动节点 EventSendMediatedTransferAfter
	N1.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "EventSendMediatedTransferAfter",
	})
	// 2. 启动节点2
	N2.Start(env)
	// 3. 初始数据记录
	cd21 := N2.GetChannelWith(N1, tokenAddress).PrintDataBeforeTransfer()
	// 4. 从节点0发起到节点1的转账
	N1.SendTrans(tokenAddress, transAmount, N2.Address, false)
	time.Sleep(time.Second * 3)
	// 5. 崩溃判断
	if N1.IsRunning() {
		msg := "Node " + N1.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 6. 中间数据记录
	models.Logger.Println("------------ Data After Crash ------------")
	N2.GetChannelWith(N1, tokenAddress).PrintDataAfterCrash()
	// 6. 重启节点1，自动发送之前中断的交易
	N1.ReStartWithoutConditionquit(env)
	time.Sleep(time.Second * 3)

	// 查询重启后数据
	models.Logger.Println("------------ Data After Restart ------------")
	cd21new := N2.GetChannelWith(N1, tokenAddress).PrintDataAfterRestart()
	// 校验对等
	models.Logger.Println("------------ Data After Fail ------------")
	if !cd21new.CheckEqualByPartnerNode(env) {
		return cm.caseFail(env.CaseName)
	}
	models.Logger.Println("------------ Data After Restart ------------")
	// cd21,交易成功
	if !cd21new.CheckSelfBalance(cd21.Balance + transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd21new.Name)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
