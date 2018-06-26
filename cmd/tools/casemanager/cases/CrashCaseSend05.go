package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/utils"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

// CrashCaseSend05 场景五：EventSendRefundTransferAfter
// 发送refundtransfer交易崩溃
// 节点2发送45token给节点6 ，发送refundtransfer后节点3崩，节点2锁定45，其余节点无锁定;
// 重启节点3后，节点2，3各锁定 45，节点2，4、节点4、5，节点5、6交易成功，转账成功，但节点2、3各锁定45token
func (cm *CaseManager) CrashCaseSend05() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCaseSend05.ENV")
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
	cd63 := utils.GetChannelBetween(N6, N3, tokenAddress).PrintDataBeforeTransfer()
	utils.GetChannelBetween(N2, N4, tokenAddress).PrintDataBeforeTransfer()
	utils.GetChannelBetween(N4, N5, tokenAddress).PrintDataBeforeTransfer()
	cd56 := utils.GetChannelBetween(N5, N6, tokenAddress).PrintDataBeforeTransfer()

	// 节点2向节点6转账20token
	N2.SendTrans(tokenAddress, transAmount, N6.Address, false)
	time.Sleep(time.Second * 3)
	//  崩溃判断
	if N3.IsRunning() {
		msg = "Node " + N3.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}

	// 查询cd23，锁定45
	utils.GetChannelBetween(N2, N3, tokenAddress).PrintDataAfterCrash()
	// 查询cd63,cd24,cd45,均无锁定
	utils.GetChannelBetween(N6, N3, tokenAddress).PrintDataAfterCrash()
	utils.GetChannelBetween(N2, N4, tokenAddress).PrintDataAfterCrash()
	utils.GetChannelBetween(N4, N5, tokenAddress).PrintDataAfterCrash()
	// 重启节点3，自动发送之前中断的交易
	N3.ReStartWithoutConditionquit(env)
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
	if cd56.Balance-cd56new.Balance != transAmount {
		models.Logger.Println(env.CaseName + " END ====> FAILED")
		return fmt.Errorf("Case [%s] FAILED", env.CaseName)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
