package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

// CrashCaseRecvAck03 场景三 MediatedTransferRecevieAck
// 节点2向节点6发送45个token，节点2崩，节点2，3各锁定45 token
// 重启后，节点2、3token解锁，成功转账节点6。
func (cm *CaseManager) CrashCaseRecvAck03() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCaseRecvAck03.ENV")
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
	N2, N3, N6 := env.Nodes[0], env.Nodes[1], env.Nodes[2]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 1. 启动
	// 启动节点3,6
	N3.Start(env)
	N6.Start(env)
	// 启动节点2, MediatedTransferRecevieAck
	N2.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "ReceiveMediatedTransferAck",
	})
	// 初始数据记录
	N3.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	N3.GetChannelWith(N6, tokenAddress).PrintDataBeforeTransfer()
	// 3. 节点2向节点6转账
	go N2.SendTrans(tokenAddress, transAmount, N6.Address, false)
	time.Sleep(time.Second * 3)
	// 4. 崩溃判断
	if N2.IsRunning() {
		msg = "Node " + N2.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 6. 中间数据记录
	models.Logger.Println("------------ Data After Crash ------------")
	cd32middle := N3.GetChannelWith(N2, tokenAddress).PrintDataAfterCrash()
	cd36middle := N3.GetChannelWith(N6, tokenAddress).PrintDataAfterCrash()
	// 校验cd32, 2锁定45
	if !cd32middle.CheckLockPartner(transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd32middle.Name)
	}
	// 校验cd36，3锁定45
	if !cd36middle.CheckLockSelf(transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd36middle.Name)
	}

	// 6. 重启节点2，交易自动继续
	N2.ReStartWithoutConditionquit(env)
	time.Sleep(time.Second * 15)

	// 查询重启后数据
	models.Logger.Println("------------ Data After Restart ------------")
	cd32new := N3.GetChannelWith(N2, tokenAddress).PrintDataAfterRestart()
	cd36new := N3.GetChannelWith(N6, tokenAddress).PrintDataAfterRestart()

	// 校验对等
	models.Logger.Println("------------ Data After Fail ------------")
	if !cd32new.CheckEqualByPartnerNode(env) || !cd36new.CheckEqualByPartnerNode(env) {
		return cm.caseFail(env.CaseName)
	}
	// 校验cd32，2锁定45
	if !cd32new.CheckLockPartner(transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd32new.Name)
	}
	// 校验cd36，3锁定45
	if !cd36new.CheckLockSelf(transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd36new.Name)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
