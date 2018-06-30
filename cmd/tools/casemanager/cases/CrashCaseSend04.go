package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

// CrashCaseSend04 场景四：EventSendSecretRequestAfter
// 发送Secretrequest后崩溃
// 节点2向节点6转账20 token,节点6发送Secretrequest后，节点6崩。
// 查询节点2，节点3，节点2锁定20 token,节点3锁定20token,交易未完成。重启节点6后，交易完成，实现转账继续。
func (cm *CaseManager) CrashCaseSend04() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCaseSend04.ENV")
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
	tokenAddress := env.Tokens[0].Address
	N2, N3, N6 := env.Nodes[0], env.Nodes[1], env.Nodes[2]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点2，3
	N2.Start(env)
	N3.Start(env)
	// 启动节点6, EventSendSecretRequestAfter
	N6.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "EventSendSecretRequestAfter",
	})
	// 记录初始数据
	cd23 := N2.GetChannelWith(N3, tokenAddress).PrintDataBeforeTransfer()
	cd36 := N3.GetChannelWith(N6, tokenAddress).PrintDataBeforeTransfer()
	// 节点2向节点6转账20token
	N2.SendTrans(tokenAddress, transAmount, N6.Address, false)
	time.Sleep(time.Second * 3)
	//  崩溃判断
	if N6.IsRunning() {
		msg = "Node " + N6.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 记录中间数据,锁定节点2节点3 20个token
	models.Logger.Println("------------ Data After Crash ------------")
	N2.GetChannelWith(N3, tokenAddress).PrintDataAfterCrash()
	N3.GetChannelWith(N6, tokenAddress).PrintDataAfterCrash()

	// 重启节点6，自动发送之前中断的交易
	N6.ReStartWithoutConditionquit(env)
	time.Sleep(time.Second * 3)

	// 查询重启后数据
	models.Logger.Println("------------ Data After Restart ------------")
	cd23new := N2.GetChannelWith(N3, tokenAddress).PrintDataAfterRestart()
	cd36new := N3.GetChannelWith(N6, tokenAddress).PrintDataAfterRestart()

	// 校验对等
	models.Logger.Println("------------ Data After Fail ------------")
	if !cd23new.CheckEqualByPartnerNode(env) || !cd36new.CheckEqualByPartnerNode(env) {
		return cm.caseFail(env.CaseName)
	}
	// cd23, 交易成功
	if !cd23new.CheckPartnerBalance(cd23.PartnerBalance + transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd23new.Name)
	}
	// 查询cd36并校验
	if !cd36new.CheckPartnerBalance(cd36.PartnerBalance + transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd36new.Name)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
