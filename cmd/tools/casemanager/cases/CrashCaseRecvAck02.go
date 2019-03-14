package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

// CrashCaseRecvAck02 场景二 SecretRecevieAck
// 节点2向节点6发送45个token，发送成功，节点2崩。
// 转账成功，没有锁定token,重启后，节点2扣钱。
// 此种情况下，崩溃不影响交易。
func (cm *CaseManager) CrashCaseRecvAck02() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCaseRecvAck02.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	transAmount = 45
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N2, N3, N6 := env.Nodes[0], env.Nodes[1], env.Nodes[2]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 1. 启动
	// 启动节点3,6
	cm.startNodes(env, N3, N6)

	// 启动节点2, ReceiveSecretRequestAck
	N2.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "ReceiveRevealSecretAck",
	})
	if cm.UseMatrix{
		time.Sleep(time.Second*10)
	}
	// 初始数据记录
	N3.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	cd36 := N3.GetChannelWith(N6, tokenAddress).PrintDataBeforeTransfer()
	// 3. 节点2向节点6转账
	N2.SendTrans(tokenAddress, transAmount, N6.Address, false)
	time.Sleep(time.Second * 3)
	// 4. 崩溃判断
	for i := 0; i < cm.HighMediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		if !N2.IsRunning() {
			break
		}
	}
	if N2.IsRunning() {
		msg = "Node " + N2.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 6. 中间数据记录
	models.Logger.Println("------------ Data After Crash ------------")
	N3.GetChannelWith(N2, tokenAddress).PrintDataAfterCrash()
	cd36middle := N3.GetChannelWith(N6, tokenAddress).PrintDataAfterCrash()
	//// 校验cd32，交易成功
	//if !cd32middle.CheckSelfBalance(cd32.Balance + transAmount) {
	//	return cm.caseFailWithWrongChannelData(env.CaseName, cd32middle.Name)
	//}
	// 校验cd36，交易成功
	if !cd36middle.CheckPartnerBalance(cd36.PartnerBalance + transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd36middle.Name)
	}

	// 6. 重启节点2，交易自动继续
	N2.ReStartWithoutConditionquit(env)
	if cm.UseMatrix{
		time.Sleep(time.Second*5)
	}
	for i := 0; i < cm.HighMediumWaitSeconds; i++ {
		time.Sleep(time.Second)

		// 查询重启后数据
		models.Logger.Println("------------ Data After Restart ------------")
		cd32new := N3.GetChannelWith(N2, tokenAddress).PrintDataAfterRestart()
		cd36new := N3.GetChannelWith(N6, tokenAddress).PrintDataAfterRestart()

		// 校验对等
		models.Logger.Println("------------ Data After Fail ------------")
		if !cd32new.CheckEqualByPartnerNode(env) || !cd36new.CheckEqualByPartnerNode(env) {
			continue
		}
		//// 校验cd32，交易成功
		//if !cd32new.CheckSelfBalance(cd32.Balance + transAmount) {
		//	return cm.caseFailWithWrongChannelData(env.CaseName, cd32new.Name)
		//}
		// 校验cd36，交易成功
		if !cd36new.CheckPartnerBalance(cd36.PartnerBalance + transAmount) {
			continue
		}

		models.Logger.Println(env.CaseName + " END ====> SUCCESS")
		return
	}
	return cm.caseFail(env.CaseName)
}
