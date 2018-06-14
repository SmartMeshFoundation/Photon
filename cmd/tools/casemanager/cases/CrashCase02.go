package cases

import (
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/utils"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

// CrashCase02 场景二：EventSendRevealSecretAfter
// 节点2向节点6转账20token,发送revealsecret后，节点2崩，路由走2-3-6，查询节点6，节点3，交易未完成，锁定节点3 20个token,节点2 20个token，
// 重启后，节点3和节点6的交易完成，节点2和节点3交易未完成，继续锁定20token。再次发送转账交易后，两次交易都完成。
// 再次发转账会解锁第一笔交易，存在问题。
func (cm *CaseManager) CrashCase02() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCase02.ENV")
	if err != nil {
		return
	}
	// 源数据
	var transAmount int32
	transAmount = 20
	tokenAddress := env.Tokens[0].Address
	N2, N3, N6 := env.Nodes[0], env.Nodes[1], env.Nodes[2]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点2, EventSendRevealSecretAfter
	N2.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "EventSendRevealSecretAfter",
	})
	// 启动节点3，6
	N3.Start(env)
	N6.Start(env)

	// 节点2向节点6转账20token
	N2.SendTrans(tokenAddress, transAmount, N6.Address, false)

	// 查询节点3，交易未完成，锁定节点2 20个token
	cd32 := utils.GetChannelBetween(N3, N2, tokenAddress)
	cd32.Println("Channel data after transfer send, cd32:")
	if cd32.PartnerLockedAmount != transAmount {
		msg := fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", transAmount, cd32.PartnerLockedAmount)
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 查询节点6，交易未完成，锁定节点3 20个token
	cd63 := utils.GetChannelBetween(N6, N3, tokenAddress)
	cd63.Println("Channel data after transfer send, cd23:")
	if cd63.PartnerLockedAmount != transAmount {
		msg := fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", transAmount, cd63.PartnerLockedAmount)
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 重启节点2，自动发送之前中断的交易
	if N2.IsRunning() {
		panic("Node N2 should be exited,but it still running")
	}
	N2.DebugCrash = false
	N2.ConditionQuit = nil
	N2.Name = "RestartNode"
	N2.Start(env)
	// 查询结果并校验
	cd32new := utils.GetChannelBetween(N3, N2, tokenAddress)
	cd32new.Println("Channel data after transfer success, cd32new:")
	if cd32new.Balance-cd32.Balance != transAmount {
		models.Logger.Println(env.CaseName + " END ====> FAILED")
		return fmt.Errorf("Case [%s] FAILED", env.CaseName)
	}
	cd63new := utils.GetChannelBetween(N6, N3, tokenAddress)
	cd63new.Println("Channel data after transfer success, cd63new:")
	if cd63new.Balance-cd63.Balance != transAmount {
		models.Logger.Println(env.CaseName + " END ====> FAILED")
		return fmt.Errorf("Case [%s] FAILED", env.CaseName)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
