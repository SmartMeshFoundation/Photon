package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

// CrashCaseRecv01 场景一：ActionInitTargetStateChange
// 收到mtr后崩,它是接收方
// 从节点2向节点6发送45个token，节点6崩后，节点2 锁定45token，节点3锁定45token，转帐失败；重启后，转账继续。
func (cm *CaseManager) CrashCaseRecv01() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCaseRecv01.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	// 启动节点2,3
	cm.startNodes(env, N2, N3)

	// 启动节点6, ActionInitTargetStateChange
	N6.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "ActionInitTargetStateChange",
	})

	// 记录初始数据
	N2.GetChannelWith(N3, tokenAddress).PrintDataBeforeTransfer()
	N3.GetChannelWith(N6, tokenAddress).PrintDataBeforeTransfer()

	// 节点2向节点6转账20token
	go N2.SendTrans(tokenAddress, transAmount, N6.Address, false)
	time.Sleep(time.Second * 3)
	//  崩溃判断
	for i := 0; i < cm.HighMediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		if !N6.IsRunning() {
			break
		}
	}
	if N6.IsRunning() {
		msg = "Node " + N6.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 中间数据记录
	models.Logger.Println("------------ Data After Crash ------------")
	cd23middle := N2.GetChannelWith(N3, tokenAddress).PrintDataAfterCrash()
	cd36middle := N3.GetChannelWith(N6, tokenAddress).PrintDataAfterCrash()
	// cd23，锁定45
	if !cd23middle.CheckLockSelf(transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd23middle.Name)
	}
	// cd36，锁定45
	if !cd36middle.CheckLockSelf(transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd36middle.Name)
	}

	// 重启节点6，交易自动继续
	N6.ReStartWithoutConditionquit(env)
	for i := 0; i < cm.MediumWaitSeconds; i++ {
		time.Sleep(time.Second * 1)

		// 查询重启后数据
		models.Logger.Println("------------ Data After Restart ------------")
		cd23new := N2.GetChannelWith(N3, tokenAddress).PrintDataAfterRestart()
		cd36new := N3.GetChannelWith(N6, tokenAddress).PrintDataAfterRestart()

		// 校验对等
		models.Logger.Println("------------ Data After Fail ------------")
		if !cd23new.CheckEqualByPartnerNode(env) || !cd36new.CheckEqualByPartnerNode(env) {
			continue
		}

		// cd23, 交易成功
		if !cd23new.CheckPartnerBalance(cd23middle.PartnerBalance + transAmount) {
			continue
		}
		// cd36，成功
		if !cd36new.CheckPartnerBalance(cd36middle.PartnerBalance + transAmount) {
			continue
		}
		models.Logger.Println(env.CaseName + " END ====> SUCCESS")
		return
	}
	return cm.caseFail(env.CaseName)
}
