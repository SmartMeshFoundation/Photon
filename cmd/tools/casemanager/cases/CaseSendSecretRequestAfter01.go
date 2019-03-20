package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

/*
CaseSendSecretRequestAfter01 :
#路由：1-2-3
#1-2-3，设置3崩溃条件为EventSendSecretRequestAfter，1向3转帐10 token，

#重启后： bug:
#1-2通道锁移除；2-3通道双方查询锁都没有移除；过期后2-3通道锁同样没有移除

#期望:
#1-2通道锁定；2-3通道锁定；过期后1-2,2-3通道锁移除
*/
func (cm *CaseManager) CaseSendSecretRequestAfter01() (err error) {
	if !cm.RunSlow {
		return ErrorSkip
	}
	env, err := models.NewTestEnv("./cases/CaseSendSecretRequestAfter01.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	// 启动节点1, EventSendSecretRequestAfter
	N3.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "EventSendSecretRequestAfter",
	})
	cm.startNodes(env, N1, N2)

	// 初始数据记录
	N1.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	N2.GetChannelWith(N3, tokenAddress).PrintDataBeforeTransfer()
	// 节点1向节点3转账
	go N1.SendTrans(tokenAddress, transAmount, N3.Address, false)
	time.Sleep(time.Second * 3)
	//  崩溃判断
	for i := 0; i < cm.HighMediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		if !N3.IsRunning() {
			break
		}
	}
	if N3.IsRunning() {
		msg := "Node " + N3.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 中间数据记录
	models.Logger.Println("------------ Data After Crash ------------")
	cd12middle := N1.GetChannelWith(N2, tokenAddress).PrintDataAfterCrash()
	cd23middle := N2.GetChannelWith(N3, tokenAddress).PrintDataAfterCrash()
	// c12,n1锁10,
	// c23,n2锁10
	if !cd12middle.CheckLockSelf(transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, "cd12middle check failed")
	}
	if !cd23middle.CheckLockSelf(transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, "cd23middle check failed")
	}
	// 重启节点3,所有节点锁定,过期以后都移除
	N3.ReStartWithoutConditionquit(env)
	//确保HighMediumWaitSeconds>通道settle timeout
	for i := 0; i < cm.HighMediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		cd12new := N1.GetChannelWith(N2, tokenAddress).PrintDataAfterCrash()
		cd23new := N2.GetChannelWith(N3, tokenAddress).PrintDataAfterCrash()
		if !cd23new.CheckLockSelf(0) {
			continue
		}
		if !cd12new.CheckLockSelf(0) {
			continue
		}
		models.Logger.Printf("%s transfer success END ====> SUCCESS", env.CaseName)
		return
	}
	models.Logger.Printf("case fail....")
	return cm.caseFail(env.CaseName)
}
