package cases

import (
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

// CrashCase006 : only for local test
func (cm *CaseManager) CrashCase006() (err error) {
	if !cm.RunThisCaseOnly {
		return
	}
	env, err := models.NewTestEnv("./cases/CrashCase006.ENV", cm.UseMatrix, cm.EthEndPoint)
	if err != nil {
		return
	}
	defer func() {
		if env.Debug == false {
			env.KillAllPhotonNodes()
		}
	}()
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	n0, n1, n2 := env.Nodes[0], env.Nodes[1], env.Nodes[2]
	transAmount := int32(150)
	tokenAddress := env.Tokens[0].TokenAddress.String()
	// 启动
	cm.startNodes(env, n2, n1)
	n0.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "EventSendAnnouncedDisposedResponseAfter",
	})
	// 初始数据记录
	n0.GetChannelWith(n1, tokenAddress).PrintDataBeforeTransfer()
	n1.GetChannelWith(n2, tokenAddress).PrintDataBeforeTransfer()
	// 转账
	go n0.SendTrans(tokenAddress, transAmount, n2.Address, false)
	time.Sleep(time.Second * 3)
	// 崩溃判断
	for i := 0; i < cm.HighMediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		if !n0.IsRunning() {
			break
		}
	}
	if n0.IsRunning() {
		msg := "Node " + n0.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 中间数据记录
	models.Logger.Println("------------ Data After Crash ------------")
	n1.GetChannelWith(n0, tokenAddress).PrintDataAfterCrash()
	n2.GetChannelWith(n1, tokenAddress).PrintDataAfterCrash()
	// 重启
	//time.Sleep(30 * time.Second)
	n0.ReStartWithoutConditionquit(env)
	for i := 0; i < cm.HighMediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		// 查询重启后数据
		models.Logger.Println("------------ Data After Restart ------------")
		c01new := n0.GetChannelWith(n1, tokenAddress).PrintDataAfterRestart()
		c12new := n1.GetChannelWith(n2, tokenAddress).PrintDataAfterRestart()
		// 校验对等
		models.Logger.Println("------------ Data After Fail ------------")
		if !c01new.CheckEqualByPartnerNode(env) || !c12new.CheckEqualByPartnerNode(env) {
			continue
		}
		if !c01new.CheckLockBoth(0) {
			continue
		}
		if !c12new.CheckLockBoth(0) {
			continue
		}
		models.Logger.Println(env.CaseName + " END ====> SUCCESS")
		return
	}
	return cm.caseFail(env.CaseName)
}
