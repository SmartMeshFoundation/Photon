package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/utils"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

// CrashCase05 场景五：EventSendRefundTransferAfter
// 发送refundtransfer交易崩溃
// 节点2发送45token给节点6 ，发送refundtransfer后节点3崩，节点2锁定45，其余节点无锁定;
// 重启节点3后，节点2，3各锁定 45，节点2，4、节点4、5，节点5、6交易成功，转账成功，但节点2、3各锁定45token
func (cm *CaseManager) CrashCase05() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCase05.ENV")
	if err != nil {
		return
	}
	defer env.KillAllRaidenNodes()
	// 源数据
	var transAmount int32
	var msg string
	transAmount = 45
	tokenAddress := env.Tokens[0].Address
	N2, N3, N4, N5, N6 := env.Nodes[0], env.Nodes[1], env.Nodes[2], env.Nodes[3], env.Nodes[4]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点2，4,5,6
	N2.Start(env)
	N4.Start(env)
	N5.Start(env)
	N6.Start(env)
	// 启动节点3, EventSendRefundTransferAfter
	N3.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "EventSendRefundTransferAfter",
	})

	// 查询节点6，记录cd65数据
	cd63 := utils.GetChannelBetween(N6, N3, tokenAddress)
	cd63.Println("Channel data before transfer send, cd63:")
	// 查询节点2，记录cd24数据
	cd24 := utils.GetChannelBetween(N2, N4, tokenAddress)
	cd24.Println("Channel data before transfer send, cd24:")
	// 查询节点4，记录cd45数据
	cd45 := utils.GetChannelBetween(N4, N5, tokenAddress)
	cd45.Println("Channel data before transfer send, cd45:")
	// 查询节点5，记录cd56数据
	cd56 := utils.GetChannelBetween(N5, N6, tokenAddress)
	cd56.Println("Channel data before transfer send, cd56:")

	// 节点2向节点6转账20token
	N2.SendTrans(tokenAddress, transAmount, N6.Address, false)
	time.Sleep(time.Second * 3)
	//  崩溃判断
	if N3.IsRunning() {
		panic("Node N3 should be exited,but it still running")
	}

	// 查询cd23，锁定45
	cd23middle := utils.GetChannelBetween(N2, N3, tokenAddress)
	cd23middle.Println("Channel data after transfer send, cd23middle:")
	if cd23middle.LockedAmount != transAmount {
		msg = fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", transAmount, cd23middle.PartnerLockedAmount)
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 查询cd63,cd24,cd45,均无锁定
	cd63middle := utils.GetChannelBetween(N6, N3, tokenAddress)
	cd63middle.Println("Channel data after transfer send, cd63middle:")
	if cd63middle.LockedAmount != 0 || cd63middle.PartnerLockedAmount != 0 {
		msg = fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", transAmount, cd63middle.LockedAmount)
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	cd24middle := utils.GetChannelBetween(N2, N4, tokenAddress)
	cd24middle.Println("Channel data after transfer send, cd24middle:")
	if cd63middle.LockedAmount != 0 || cd24middle.PartnerLockedAmount != 0 {
		msg = fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", transAmount, cd24middle.LockedAmount)
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	cd45middle := utils.GetChannelBetween(N4, N5, tokenAddress)
	cd45middle.Println("Channel data after transfer send, cd45middle:")
	if cd45middle.LockedAmount != 0 || cd45middle.PartnerLockedAmount != 0 {
		msg = fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", transAmount, cd45middle.LockedAmount)
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}

	// 重启节点3，自动发送之前中断的交易
	N3.DebugCrash = false
	N3.ConditionQuit = nil
	N3.Name = "RestartNode"
	N3.Start(env)

	// 查询cd23并校验
	cd23new := utils.GetChannelBetween(N2, N3, tokenAddress)
	cd23new.Println("Channel data after transfer success, cd23new:")
	if cd23new.PartnerLockedAmount != transAmount || cd23new.LockedAmount != transAmount {
		models.Logger.Println(env.CaseName + " END ====> FAILED")
		return fmt.Errorf("Case [%s] FAILED", env.CaseName)
	}
	// 查询cd63并校验
	cd63new := utils.GetChannelBetween(N6, N3, tokenAddress)
	cd63new.Println("Channel data after transfer success, cd63new:")
	if cd63new.Balance-cd63.Balance != 0 {
		models.Logger.Println(env.CaseName + " END ====> FAILED")
		return fmt.Errorf("Case [%s] FAILED", env.CaseName)
	}
	// 查询cd56并校验
	cd56new := utils.GetChannelBetween(N5, N6, tokenAddress)
	cd56new.Println("Channel data after transfer success, cd56new:")
	if cd56new.Balance-cd56.Balance != transAmount {
		models.Logger.Println(env.CaseName + " END ====> FAILED")
		return fmt.Errorf("Case [%s] FAILED", env.CaseName)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
