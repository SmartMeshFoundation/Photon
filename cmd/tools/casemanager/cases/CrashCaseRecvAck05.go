package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

// CrashCaseRecvAck05 场景五：RevealSecretRecevieAck
// 节点2向节点6发送20个token，节点2崩，节点2和节点3之间通道节点2锁定20 token，节点3和节点6之间转账完成；重启后，节点2锁定20解除，完成转账。
func (cm *CaseManager) CrashCaseRecvAck05() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCaseRecvAck05.ENV")
	if err != nil {
		return
	}
	defer func() {
		if env.Debug == false {
			env.KillAllPhotonNodes()
		}
	}()
	// 源数据
	var transAmount int32
	var msg string
	transAmount = 20
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N2, N3, N6 := env.Nodes[0], env.Nodes[1], env.Nodes[2]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 1. 启动
	// 启动节点3,6
	N3.Start(env)
	N6.Start(env)
	// 启动节点2, ReceiveRevealSecretAck
	N2.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "ReceiveRevealSecretAck",
	})
	// 初始数据记录
	N3.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	cd36 := N3.GetChannelWith(N6, tokenAddress).PrintDataBeforeTransfer()
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
	N3.GetChannelWith(N2, tokenAddress).PrintDataAfterCrash()
	cd36middle := N3.GetChannelWith(N6, tokenAddress).PrintDataAfterCrash()
	//// 校验cd32, 2锁定45
	//if !cd32middle.CheckLockPartner(transAmount) {
	//	return cm.caseFailWithWrongChannelData(env.CaseName, cd32middle.Name)
	//}
	// 校验cd36，交易成功
	if !cd36middle.CheckPartnerBalance(cd36.PartnerBalance + transAmount) {
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
	//// 校验cd32，2锁定45
	//if !cd32new.CheckLockPartner(transAmount) {
	//	return cm.caseFailWithWrongChannelData(env.CaseName, cd32new.Name)
	//}
	// 校验cd36，交易成功
	if !cd36new.CheckPartnerBalance(cd36.PartnerBalance + transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd36new.Name)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
