package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/utils"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

// CrashCaseSend06 场景六：EventSendRefundTransferAfter
// 发送refundtransfer交易崩溃
// 节点2发送45token给节点6 ，发送refundtransfer后节点3崩，节点2锁定45，其余节点无锁定;
// 重启节点3后，交易失败，通道2-3中节点2，3各锁定 45，通道2-7中节点2锁定45
func (cm *CaseManager) CrashCaseSend06() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCaseSend06.ENV")
	if err != nil {
		return
	}
	defer func() {
		if env.Debug == false {
			env.KillAllRaidenNodes()
		}
	}()
	// 源数据
	var transAmount int32
	var msg string
	transAmount = 45
	tokenAddress := env.Tokens[0].Address
	N2, N3, N6, N7 := env.Nodes[0], env.Nodes[1], env.Nodes[2], env.Nodes[3]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点2，6，7
	N2.Start(env)
	N6.Start(env)
	N7.Start(env)
	// 启动节点3, EventSendRefundTransferAfter
	N3.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "EventSendRefundTransferAfter",
	})

	// 节点2向节点6转账20token
	N2.SendTrans(tokenAddress, transAmount, N6.Address, false)
	time.Sleep(time.Second * 3)
	//  崩溃判断
	if N3.IsRunning() {
		msg = "Node " + N3.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 崩溃后数据校验
	utils.GetChannelBetween(N2, N3, tokenAddress).PrintDataAfterCrash()
	utils.GetChannelBetween(N6, N3, tokenAddress).PrintDataAfterCrash()
	utils.GetChannelBetween(N2, N7, tokenAddress).PrintDataAfterCrash()
	utils.GetChannelBetween(N7, N3, tokenAddress).PrintDataAfterCrash()

	// 重启节点3，自动发送之前中断的交易
	N3.ReStartWithoutConditionquit(env)
	time.Sleep(time.Second * 30)

	// 重启后校验
	// cd23,双锁定
	ar1 := utils.CheckChannelLockBoth(N2, N3, tokenAddress, transAmount, utils.PrintDataAfterRestart)
	// cd63，无锁定
	ar2 := utils.CheckChannelNoLock(N3, N6, tokenAddress, utils.PrintDataAfterRestart)
	// cd73，7锁定45
	ar3 := utils.CheckChannelLockPartner(N3, N7, tokenAddress, transAmount, utils.PrintDataAfterRestart)
	// cd27，2锁定45
	ar4 := utils.CheckChannelLockPartner(N7, N2, tokenAddress, transAmount, utils.PrintDataAfterRestart)

	isSuccess := ar1 && ar2 && ar3 && ar4
	if !isSuccess {
		models.Logger.Println(env.CaseName + " END ====> FAILED")
		return fmt.Errorf("Case [%s] FAILED", env.CaseName)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
