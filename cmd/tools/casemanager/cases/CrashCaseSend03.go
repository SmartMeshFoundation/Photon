package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/utils"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

// CrashCaseSend03 场景三：EventSendBalanceProofAfter
// 发送余额证明后崩溃（发送方崩）
// 节点2向节点6转账20 token,发送balanceProof后，节点2崩，路由走2-3-6，查询节点3，节点6，节点3和6之间交易完成。
// 节点2、3交易未完成，节点2锁定20token。重启节点2后，节点2、3交易完成，实现转账继续。
func (cm *CaseManager) CrashCaseSend03() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCaseSend03.ENV")
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
	// 启动节点2, EventSendRevealSecretAfter
	N2.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "EventSendBalanceProofAfter",
	})
	// 启动节点3，6
	N3.Start(env)
	N6.Start(env)

	// 查询节点6，记录cd63数据
	cd63 := utils.GetChannelBetween(N6, N3, tokenAddress)
	cd63.Println("Channel data before transfer send, cd63:")
	// 节点2向节点6转账20token
	N2.SendTrans(tokenAddress, transAmount, N6.Address, false)
	time.Sleep(time.Second * 3)
	//  崩溃判断
	if N2.IsRunning() {
		panic("Node N2 should be exited,but it still running")
	}
	// 查询节点3，与节点2交易未完成，锁定节点2 20个token
	cd32 := utils.GetChannelBetween(N3, N2, tokenAddress)
	cd32.Println("Channel data after transfer send, cd32:")
	if cd32.PartnerLockedAmount != transAmount {
		msg = fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", transAmount, cd32.PartnerLockedAmount)
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 查询节点6，与节点3交易完成
	cd63new := utils.GetChannelBetween(N6, N3, tokenAddress)
	cd63new.Println("Channel data after transfer send, cd63new:")
	if cd63new.Balance-cd63.Balance != transAmount {
		msg = "Expect transfer on cd63 success,but got failed ,FAILED!!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 重启节点2，自动发送之前中断的交易
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
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
