package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

/*
CaseSendUnlockBefore02 ##构建1->2->3->4,2->4直接通道金额不够，设置2的奔溃事件EventSendUnlockBefore
#2崩溃后2-3锁钱，重启后锁保留，锁超时后,3连上注册密码,最后转账成功
*/
func (cm *CaseManager) CaseSendUnlockBefore02() (err error) {
	if !cm.RunSlow {
		return
	}
	env, err := models.NewTestEnv("./cases/CaseSendUnlockBefore02.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	transAmount = 50
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N1, N2, N3, N4 := env.Nodes[1], env.Nodes[2], env.Nodes[3], env.Nodes[4]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点2,
	N2.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "EventSendUnlockBefore",
	})
	// 启动节点1,3,4
	cm.startNodes(env, N1, N3, N4)
	if cm.UseMatrix {
		time.Sleep(time.Second * 10)
	}
	// 初始数据记录
	N3.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	// 节点2向节点6转账20token
	N1.SendTrans(tokenAddress, transAmount, N4.Address, false)
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
	c32new := N3.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	if !c32new.CheckLockPartner(transAmount) {
		return fmt.Errorf("CheckLockPartner 2 err %s", err)
	}
	waitForTimeout := cm.MediumWaitSeconds
	if cm.UseMatrix {
		waitForTimeout = cm.MediumWaitSeconds + 250
	}
	err = cm.tryInSeconds(waitForTimeout, func() error {
		models.Logger.Println("check...")
		var c channeltype.ChannelDataDetail
		c, err = N3.SpecifiedChannel(c32new.ChannelIdentifier)
		if len(c.PartnerKnownSecretLocks) <= 0 {
			return fmt.Errorf("CheckLockPartner after restart err %s", err)
		}
		found := false
		for _, l := range c.PartnerKnownSecretLocks {
			if l.IsRegisteredOnChain {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("secret not registered")
		}
		return nil
	})
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, "secret register error")
	}
	// 重启节点2，
	N2.ReStartWithoutConditionquit(env)
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
	err = cm.tryInSeconds(cm.MediumWaitSeconds, func() error {
		models.Logger.Println("check...")
		c32new = N3.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
		if !c32new.CheckLockPartner(0) {
			return fmt.Errorf("CheckLockPartner ReStartWithoutConditionquit err %s", err)
		}
		return nil
	})
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, "ReStartWithoutConditionquit lock not released")
	}
	models.Logger.Println(env.CaseName, "===>SUCCESS")
	return nil
}
