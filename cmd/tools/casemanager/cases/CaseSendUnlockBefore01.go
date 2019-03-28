package cases

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/params"

	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

/*
CaseSendUnlockBefore01 :
#路由：1-2-3

#1-2-3，通道资金均是100，设置2崩溃条件为EventSendUnlockBefore，1向3转帐10token，

#重启后，bug状况: 1-2转账成功，2-3通道锁钱；过期后，2-3移除过期锁(交易失败)，1-2转账成功，导致1丢钱

#重启后希望结果,12成功,但是2-3之间要等待3链上注册密码,然后才能成功.
*/
func (cm *CaseManager) CaseSendUnlockBefore01() (err error) {
	if !cm.RunSlow {
		return ErrorSkip
	}
	env, err := models.NewTestEnv("./cases/CaseSendUnlockBefore01.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	// 启动节点1,2,3
	cm.startNodes(env, N1, N3, N2.SetConditionQuit(&params.ConditionQuit{
		QuitEvent: "EventSendUnlockBefore",
	}))
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
	// 初始数据记录
	N1.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	N2.GetChannelWith(N3, tokenAddress).PrintDataBeforeTransfer()
	// 节点2向节点6转账20token
	go N1.SendTrans(tokenAddress, transAmount, N3.Address, false)
	time.Sleep(time.Second * 3)
	//  崩溃判断
	for i := 0; i < cm.HighMediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		if !N2.IsRunning() {
			break
		}
	}
	if N2.IsRunning() {
		msg := "Node " + N2.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 中间数据记录
	models.Logger.Println("------------ Data After Crash ------------")
	cd12middle := N1.GetChannelWith(N2, tokenAddress).PrintDataAfterCrash()
	cd23middle := N3.GetChannelWith(N2, tokenAddress).PrintDataAfterCrash()
	_ = cd12middle
	// c12,n1锁10,
	// c23,n2锁10
	//1可能这边已经解锁了.
	//if !cd12middle.CheckLockPartner(transAmount) {
	//	return cm.caseFailWithWrongChannelData(env.CaseName, "cd12middle check failed")
	//}
	//_ = cd23middle
	//不能检查23,有可能两者交易已经成功了
	if !cd23middle.CheckLockPartner(transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, "cd23middle check failed")
	}
	// 重启节点2，2,3之间等待超时注册密码,然后交易都成功
	N2.ReStartWithoutConditionquit(env)
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
