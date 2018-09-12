package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

// CrashCaseSend03 场景三：EventSendBalanceProofAfter
// 发送余额证明后崩溃（发送方崩）
// 节点2向节点6转账20 token,发送balanceProof后，节点2崩，路由走2-3-6，查询节点3，节点6，节点3和6之间交易完成。
// 节点2、3交易未完成，节点2锁定20token。重启节点2后，节点2、3交易完成，实现转账继续。
func (cm *CaseManager) CrashCaseSend03() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCaseSend03.ENV")
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
	transAmount = 20
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N2, N3, N6 := env.Nodes[0], env.Nodes[1], env.Nodes[2]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点2, EventSendRevealSecretAfter
	N2.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "EventSendUnlockAfter",
	})
	// 启动节点3，6
	N3.Start(env)
	N6.Start(env)

	// 初始数据记录
	cd32 := N3.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	cd63 := N6.GetChannelWith(N3, tokenAddress).PrintDataBeforeTransfer()
	// 节点2向节点6转账20token
	N2.SendTrans(tokenAddress, transAmount, N6.Address, false)
	time.Sleep(time.Second * 3)
	//  崩溃判断
	if N2.IsRunning() {
		msg := "Node " + N2.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 中间数据记录
	models.Logger.Println("------------ Data After Crash ------------")
	N3.GetChannelWith(N2, tokenAddress).PrintDataAfterCrash()
	cd63middle := N6.GetChannelWith(N3, tokenAddress).PrintDataAfterCrash()
	// cd63，交易成功
	if !cd63middle.CheckSelfBalance(cd63.Balance + transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd63middle.Name)
	}
	// 重启节点2，自动发送之前中断的交易
	N2.ReStartWithoutConditionquit(env)
	time.Sleep(time.Second * 30)

	// 查询重启后数据
	models.Logger.Println("------------ Data After Restart ------------")
	cd32new := N3.GetChannelWith(N2, tokenAddress).PrintDataAfterRestart()
	cd63new := N6.GetChannelWith(N3, tokenAddress).PrintDataAfterRestart()

	// 校验对等
	models.Logger.Println("------------ Data After Fail ------------")
	if !cd32new.CheckEqualByPartnerNode(env) || !cd63new.CheckEqualByPartnerNode(env) {
		return cm.caseFail(env.CaseName)
	}
	// cd32, 交易成功
	if !cd32new.CheckSelfBalance(cd32.Balance + transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd32new.Name)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
