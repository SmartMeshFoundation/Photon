package cases

import (
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

// CrashCase011 : only for local test
func (cm *CaseManager) CrashCase011() (err error) {
	if !cm.RunThisCaseOnly {
		return ErrorSkip
	}
	env, err := models.NewTestEnv("./cases/CrashCase011.ENV", cm.UseMatrix, cm.EthEndPoint)
	if err != nil {
		return
	}
	defer func() {
		if env.Debug == false {
			env.KillAllPhotonNodes()
		}
	}()
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	n0, n1, n2, n3, n4, n5, n6 := env.Nodes[0], env.Nodes[1], env.Nodes[2], env.Nodes[3], env.Nodes[4], env.Nodes[5], env.Nodes[6]
	transAmount := int32(150)
	tokenAddress := env.Tokens[0].TokenAddress.String()
	// 启动
	cm.startNodes(env, n0, n2, n3, n4, n5, n6)
	n1.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "MediatorReReceiveStateChange",
	})
	// 初始数据记录
	n0.GetChannelWith(n1, tokenAddress).PrintDataBeforeTransfer()
	n1.GetChannelWith(n2, tokenAddress).PrintDataBeforeTransfer()
	n2.GetChannelWith(n3, tokenAddress).PrintDataBeforeTransfer()
	n3.GetChannelWith(n1, tokenAddress).PrintDataBeforeTransfer()
	n2.GetChannelWith(n4, tokenAddress).PrintDataBeforeTransfer()
	n4.GetChannelWith(n6, tokenAddress).PrintDataBeforeTransfer()
	n6.GetChannelWith(n5, tokenAddress).PrintDataBeforeTransfer()
	n1.GetChannelWith(n5, tokenAddress).PrintDataBeforeTransfer()
	// 转账
	go n0.SendTrans(tokenAddress, transAmount, n4.Address, false)
	time.Sleep(time.Second * 3)
	// 崩溃判断
	for i := 0; i < cm.HighMediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		if !n1.IsRunning() {
			break
		}
	}
	if n1.IsRunning() {
		msg := "Node " + n1.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 中间数据记录
	models.Logger.Println("------------ Data After Crash ------------")
	n0.GetChannelWith(n1, tokenAddress).PrintDataBeforeTransfer()
	n2.GetChannelWith(n1, tokenAddress).PrintDataBeforeTransfer()
	n2.GetChannelWith(n3, tokenAddress).PrintDataBeforeTransfer()
	n3.GetChannelWith(n1, tokenAddress).PrintDataBeforeTransfer()
	n2.GetChannelWith(n4, tokenAddress).PrintDataBeforeTransfer()
	n4.GetChannelWith(n6, tokenAddress).PrintDataBeforeTransfer()
	n6.GetChannelWith(n5, tokenAddress).PrintDataBeforeTransfer()
	n5.GetChannelWith(n1, tokenAddress).PrintDataBeforeTransfer()
	// 重启
	//time.Sleep(30 * time.Second)
	n1.ReStartWithoutConditionquit(env)
	for i := 0; i < cm.HighMediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		// 查询重启后数据
		models.Logger.Println("------------ Data After Restart ------------")
		c01new := n0.GetChannelWith(n1, tokenAddress).PrintDataBeforeTransfer()
		c12new := n1.GetChannelWith(n2, tokenAddress).PrintDataBeforeTransfer()
		c23new := n2.GetChannelWith(n3, tokenAddress).PrintDataBeforeTransfer()
		c31new := n3.GetChannelWith(n1, tokenAddress).PrintDataBeforeTransfer()
		c24new := n2.GetChannelWith(n4, tokenAddress).PrintDataBeforeTransfer()
		c46new := n4.GetChannelWith(n6, tokenAddress).PrintDataBeforeTransfer()
		c65new := n6.GetChannelWith(n5, tokenAddress).PrintDataBeforeTransfer()
		c15new := n1.GetChannelWith(n5, tokenAddress).PrintDataBeforeTransfer()
		// 校验对等
		models.Logger.Println("------------ Data After Fail ------------")
		if !c01new.CheckEqualByPartnerNode(env) || !c12new.CheckEqualByPartnerNode(env) {
			continue
		}
		if !c01new.CheckLockBoth(0) {
			continue
		}
		if !c23new.CheckLockBoth(0) {
			continue
		}
		if !c31new.CheckLockBoth(0) {
			continue
		}
		if !c24new.CheckLockBoth(0) {
			continue
		}
		if !c46new.CheckLockBoth(0) {
			continue
		}
		if !c65new.CheckLockBoth(0) {
			continue
		}
		if !c15new.CheckLockBoth(0) {
			continue
		}
		models.Logger.Println(env.CaseName + " END ====> SUCCESS")
		return
	}
	return cm.caseFail(env.CaseName)
}
