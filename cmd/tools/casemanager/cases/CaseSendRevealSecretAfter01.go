package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

/*
CaseSendRevealSecretAfter01 :
#路由：1-2-3

#1-2-3，设置1崩溃条件为EventSendRevealSecretAfter，1向3转帐10token，

#1重启后，有两种情况, 一种3没有收到密码,所有交易失败;一种是3收到了密码,2选择链上注册密码,从而所有节点成功
*/
func (cm *CaseManager) CaseSendRevealSecretAfter01() (err error) {
	if !cm.RunSlow {
		return ErrorSkip
	}
	env, err := models.NewTestEnv("./cases/CaseSendRevealSecretAfter01.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	// 启动节点1, EventSendRevealSecretAfter
	N1.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "EventSendRevealSecretAfter",
	})
	// 启动节点3，6
	cm.startNodes(env, N2, N3)

	// 初始数据记录
	N1.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	N2.GetChannelWith(N3, tokenAddress).PrintDataBeforeTransfer()
	// 节点1向节点3转账
	go N1.SendTrans(tokenAddress, transAmount, N3.Address, false)
	//time.Sleep(time.Second * 3)
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
	cd23middle := N2.GetChannelWith(N3, tokenAddress).PrintDataAfterCrash()
	// c12,n1锁10,
	// c23,n2锁10
	if !cd12middle.CheckLockPartner(transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, "cd12middle check failed")
	}
	_ = cd23middle
	//不能检查23,有可能两者交易已经成功了
	//if !cd23middle.CheckLockSelf(transAmount) {
	//	return cm.caseFailWithWrongChannelData(env.CaseName, "cd23middle check failed")
	//}
	// 重启节点1,要么所有交易都成功,要么都失败
	N1.ReStartWithoutConditionquit(env)
	//等待3秒,如果成功就会出现1,2锁定,但是2,3没有锁定
	//如果失败,就会是所有都锁定.
	time.Sleep(time.Second * 3)
	//确保HighMediumWaitSeconds>通道settle timeout
	for i := 0; i < cm.HighMediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		cd12new := N1.GetChannelWith(N2, tokenAddress).PrintDataAfterCrash()
		cd23new := N2.GetChannelWith(N3, tokenAddress).PrintDataAfterCrash()

		failure := false
		//失败的情况
		if cd12new.CheckLockSelf(transAmount) && cd23new.CheckLockSelf(transAmount) {
			failure = true
		}
		success := false
		//成功的情况
		if cd12new.CheckLockBoth(0) && cd23new.CheckLockBoth(0) {
			success = true
		}
		if success {
			models.Logger.Printf("%s transfer success END ====> SUCCESS", env.CaseName)
			return
		}
		if failure {
			models.Logger.Printf("%s transfer fail  END ====> SUCCESS", env.CaseName)
			return
		}
	}
	models.Logger.Printf("case fail....")
	return cm.caseFail(env.CaseName)
}
