package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/utils"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

// CrashCase04 场景四：EventSendSecretRequestAfter
// 发送Secretrequest后崩溃
// 节点2向节点6转账20 token,节点6发送Secretrequest后，节点6崩。
// 查询节点2，节点3，节点2锁定20 token,节点3锁定20token,交易未完成。重启节点6后，交易完成，实现转账继续。
func (cm *CaseManager) CrashCase04() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCase04.ENV")
	if err != nil {
		return
	}
	defer env.KillAllRaidenNodes()
	// 源数据
	var transAmount int32
	var msg string
	transAmount = 20
	tokenAddress := env.Tokens[0].Address
	N2, N3, N6 := env.Nodes[0], env.Nodes[1], env.Nodes[2]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点2，3
	N2.Start(env)
	N3.Start(env)
	// 启动节点6, EventSendSecretRequestAfter
	N6.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "EventSendSecretRequestAfter",
	})
	// 查询节点6，记录cd63数据
	cd63 := utils.GetChannelBetween(N6, N3, tokenAddress)
	cd63.Println("Channel data before transfer send, cd63:")
	// 节点2向节点6转账20token
	N2.SendTrans(tokenAddress, transAmount, N6.Address, false)
	time.Sleep(time.Second * 3)
	//  崩溃判断
	if N6.IsRunning() {
		panic("Node N6 should be exited,but it still running")
	}
	// 查询节点3，与节点2交易未完成，锁定节点2 20个token
	cd32 := utils.GetChannelBetween(N3, N2, tokenAddress)
	cd32.Println("Channel data after transfer send, cd32:")
	if cd32.PartnerLockedAmount != transAmount {
		msg = fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", transAmount, cd32.PartnerLockedAmount)
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 查询节点3，与节点交易未完成，锁定节点3 20个token
	cd36 := utils.GetChannelBetween(N3, N6, tokenAddress)
	cd36.Println("Channel data after transfer send, cd36:")
	if cd36.LockedAmount != transAmount {
		msg = fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", transAmount, cd36.LockedAmount)
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 重启节点6，自动发送之前中断的交易
	N6.DebugCrash = false
	N6.ConditionQuit = nil
	N6.Name = "RestartNode"
	N6.Start(env)
	// 查询cd32并校验
	cd32new := utils.GetChannelBetween(N3, N2, tokenAddress)
	cd32new.Println("Channel data after transfer success, cd32new:")
	if cd32new.Balance-cd32.Balance != transAmount {
		models.Logger.Println(env.CaseName + " END ====> FAILED")
		return fmt.Errorf("Case [%s] FAILED", env.CaseName)
	}
	// 查询cd36并校验
	cd36new := utils.GetChannelBetween(N3, N6, tokenAddress)
	cd36new.Println("Channel data after transfer success, cd36new:")
	if cd36.Balance-cd36new.Balance != transAmount {
		models.Logger.Println(env.CaseName + " END ====> FAILED")
		return fmt.Errorf("Case [%s] FAILED", env.CaseName)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
