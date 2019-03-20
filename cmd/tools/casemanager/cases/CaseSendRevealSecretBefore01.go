package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

/*
CaseSendRevealSecretBefore01 :
#路由：1-2-3

#1-2-3，通道资金均是100，设置1崩溃条件为 EventSendRevealSecretBefore，1向3转帐10token，

#重启后，bug: 1-2；2-3锁钱；过期后，1-2锁移除；2-3锁没有移除

#期望:
#过期以后 1-2,2-3锁均移除

*/
func (cm *CaseManager) CaseSendRevealSecretBefore01() (err error) {
	if !cm.RunSlow {
		return ErrorSkip
	}
	env, err := models.NewTestEnv("./cases/CaseSendRevealSecretBefore01.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	transAmount = 10
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N1, N2, N3 := env.Nodes[0], env.Nodes[1], env.Nodes[2]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点2,
	N1.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "EventSendRevealSecretBefore",
	})
	// 启动节点3，6
	cm.startNodes(env, N2, N3)

	// 初始数据记录
	N1.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	N2.GetChannelWith(N3, tokenAddress).PrintDataBeforeTransfer()
	// 节点2向节点6转账20token
	go N1.SendTrans(tokenAddress, transAmount, N3.Address, false)
	time.Sleep(time.Second * 3)
	//  崩溃判断
	for i := 0; i < cm.HighMediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		if !N1.IsRunning() {
			break
		}
	}
	if N1.IsRunning() {
		msg := "Node " + N1.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 中间数据记录
	models.Logger.Println("------------ Data After Crash ------------")
	cd12middle := N2.GetChannelWith(N1, tokenAddress).PrintDataAfterCrash()
	cd23middle := N3.GetChannelWith(N2, tokenAddress).PrintDataAfterCrash()

	// c12,n1锁10,
	// c23,n2锁10
	if !cd12middle.CheckLockPartner(transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, "cd12middle check failed")
	}
	if !cd23middle.CheckLockPartner(transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, "cd23middle check failed")
	}
	// 重启节点1， 过期后所有锁解除
	N1.ReStartWithoutConditionquit(env)
	models.Logger.Println("------------ Data After Restart ------------")
	//确保HighMediumWaitSeconds>通道settle timeout
	for i := 0; i < cm.HighMediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		cd12new := N1.GetChannelWith(N2, tokenAddress).PrintDataAfterCrash()
		cd23new := N2.GetChannelWith(N3, tokenAddress).PrintDataAfterCrash()

		if !cd12new.CheckLockBoth(0) {
			continue
		}

		if cd23new.CheckLockBoth(0) {
			models.Logger.Printf("%s transfer success END ====> SUCCESS", env.CaseName)
			return
		}
	}
	models.Logger.Printf("case fail....")
	return cm.caseFail(env.CaseName)
}
