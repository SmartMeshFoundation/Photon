package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

// CrashCaseRecv02 场景二：ReceiveSecretRequestStateChange
// 收到Secretrequest后崩
// 节点1向节点6发送20个token,节点6向节点1发送secretrequest请求，节点1收到崩,
// 节点1、节点2、节点3各锁定20个token；重启节点1后，节点锁定token解锁，转账成功。
func (cm *CaseManager) CrashCaseRecv02() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCaseRecv02.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	N1, N2, N3, N6 := env.Nodes[0], env.Nodes[1], env.Nodes[2], env.Nodes[3]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点2,36
	N2.Start(env)
	N3.Start(env)
	N6.Start(env)
	// 启动节点1, ReceiveSecretRequestStateChange
	N1.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "ReceiveSecretRequestStateChange",
	})

	// 记录初始数据
	N2.GetChannelWith(N1, tokenAddress).PrintDataBeforeTransfer()
	N2.GetChannelWith(N3, tokenAddress).PrintDataBeforeTransfer()
	N3.GetChannelWith(N6, tokenAddress).PrintDataBeforeTransfer()

	// 节点1向节点6转账20token
	go N1.SendTrans(tokenAddress, transAmount, N6.Address, false)
	time.Sleep(time.Second * 3)
	//  崩溃判断
	if N1.IsRunning() {
		msg = "Node " + N1.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 中间数据记录
	models.Logger.Println("------------ Data After Crash ------------")
	cd21middle := N2.GetChannelWith(N1, tokenAddress).PrintDataAfterCrash()
	cd23middle := N2.GetChannelWith(N3, tokenAddress).PrintDataAfterCrash()
	cd36middle := N3.GetChannelWith(N6, tokenAddress).PrintDataAfterCrash()
	// 查询cd21，锁定对方20
	if !cd21middle.CheckLockPartner(transAmount) {
		return cm.caseFail(env.CaseName)
	}
	// 查询cd23，锁定20
	if !cd23middle.CheckLockSelf(transAmount) {
		return cm.caseFail(env.CaseName)
	}
	// 查询cd36，锁定20
	if !cd36middle.CheckLockSelf(transAmount) {
		return cm.caseFail(env.CaseName)
	}

	// 重启节点1，交易失败,各自锁定
	time.Sleep(30 * time.Second)
	N1.ReStartWithoutConditionquit(env)
	//for i := 0; i < cm.MediumWaitSeconds; i++ {
	//	time.Sleep(time.Second * 1)
	//
	//	// 查询重启后数据
	//	models.Logger.Println("------------ Data After Restart ------------")
	//	cd21new := N2.GetChannelWith(N1, tokenAddress).PrintDataAfterRestart()
	//	cd23new := N2.GetChannelWith(N3, tokenAddress).PrintDataAfterRestart()
	//	cd36new := N3.GetChannelWith(N6, tokenAddress).PrintDataAfterRestart()
	//
	//	// 校验对等
	//	models.Logger.Println("------------ Data After Fail ------------")
	//	if !cd21new.CheckEqualByPartnerNode(env) || !cd23new.CheckEqualByPartnerNode(env) || !cd36new.CheckEqualByPartnerNode(env) {
	//		continue
	//	}
	//
	//	// 查询cd21，锁定对方20
	//	if !cd21new.CheckLockPartner(transAmount) {
	//		continue
	//	}
	//	// 查询cd23，锁定20
	//	if !cd23new.CheckLockSelf(transAmount) {
	//		continue
	//	}
	//	// 查询cd36，锁定20
	//	if !cd36new.CheckLockSelf(transAmount) {
	//		continue
	//	}
	//	break
	//	//models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	//	//return
	//}
	// 等待超时remove
	for i := 0; i < cm.MediumWaitSeconds; i++ {
		time.Sleep(time.Second * 1)

		// 查询重启后数据
		models.Logger.Println("------------ Data After Restart ------------")
		cd21new := N2.GetChannelWith(N1, tokenAddress).PrintDataAfterRestart()
		cd23new := N2.GetChannelWith(N3, tokenAddress).PrintDataAfterRestart()
		cd36new := N3.GetChannelWith(N6, tokenAddress).PrintDataAfterRestart()

		// 校验对等
		models.Logger.Println("------------ Data After Fail ------------")
		if !cd21new.CheckEqualByPartnerNode(env) || !cd23new.CheckEqualByPartnerNode(env) || !cd36new.CheckEqualByPartnerNode(env) {
			continue
		}

		// 查询cd21，无锁
		if !cd21new.CheckLockBoth(0) {
			continue
		}
		// 查询cd23，无锁
		if !cd23new.CheckLockBoth(0) {
			continue
		}
		// 查询cd36，无锁
		if !cd36new.CheckLockBoth(0) {
			continue
		}
		models.Logger.Println(env.CaseName + " END ====> SUCCESS")
		return
	}
	return cm.caseFail(env.CaseName)
}
