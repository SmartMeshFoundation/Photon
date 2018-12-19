package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

// CrashCaseRecvAck01 场景一：SecretRequestRecevieAck
// 节点2向节点6发送20个token，发送成功，节点6崩。
// 此种情况下，转账成功，崩溃不影响交易。
// 继续转账，转账成功。
func (cm *CaseManager) CrashCaseRecvAck01() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCaseRecvAck01.ENV")
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
	// 启动节点2,3,6
	N2.Start(env)
	N3.Start(env)
	// 启动节点6, SecretRequestRecevieAck
	N6.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "ReceiveSecretRequestAck",
	})
	// 初始数据记录
	N2.GetChannelWith(N3, tokenAddress).PrintDataBeforeTransfer()
	N3.GetChannelWith(N6, tokenAddress).PrintDataBeforeTransfer()
	// 3. 节点2向节点6转账20token
	go N2.SendTrans(tokenAddress, transAmount, N6.Address, false)
	time.Sleep(time.Second * 10)
	// 4. 崩溃判断
	if N6.IsRunning() {
		msg = "Node " + N6.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 6. 中间数据记录
	models.Logger.Println("------------ Data After Crash ------------")
	cd23middle := N2.GetChannelWith(N3, tokenAddress).PrintDataAfterCrash()
	cd36middle := N3.GetChannelWith(N6, tokenAddress).PrintDataAfterCrash()
	// 校验cd23，无锁
	if !cd23middle.CheckNoLock() {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd23middle.Name)
	}
	// 校验cd36，无锁
	if !cd36middle.CheckNoLock() {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd36middle.Name)
	}

	// 6. 重启节点2，交易自动继续
	N6.ReStartWithoutConditionquit(env)
	time.Sleep(time.Second * 15)

	// 查询重启后数据
	models.Logger.Println("------------ Data After Restart ------------")
	cd23new := N2.GetChannelWith(N3, tokenAddress).PrintDataAfterRestart()
	cd36new := N3.GetChannelWith(N6, tokenAddress).PrintDataAfterRestart()

	// 校验对等
	models.Logger.Println("------------ Data After Fail ------------")
	if !cd23new.CheckEqualByPartnerNode(env) || !cd36new.CheckEqualByPartnerNode(env) {
		return cm.caseFail(env.CaseName)
	}
	// 校验cd23，无锁
	if !cd23new.CheckNoLock() {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd23new.Name)
	}
	// 校验cd36，无锁
	if !cd36new.CheckNoLock() {
		return cm.caseFailWithWrongChannelData(env.CaseName, cd36new.Name)
	}

	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
