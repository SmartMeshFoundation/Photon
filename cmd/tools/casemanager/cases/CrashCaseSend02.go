package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/utils"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

// CrashCaseSend02 场景二：EventSendRevealSecretAfter
// 节点2向节点6转账20token,发送revealsecret后，节点2崩，路由走2-3-6，查询节点6，节点3，交易未完成，锁定节点3 20个token,节点2 20个token，
// 重启后，节点3和节点6的交易完成，节点2和节点3交易未完成，继续锁定20token。再次发送转账交易后，两次交易都完成。
// 再次发转账会解锁第一笔交易，存在问题。
func (cm *CaseManager) CrashCaseSend02() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCaseSend02.ENV")
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

	// 初始数据记录
	utils.GetChannelBetween(N3, N2, tokenAddress).PrintDataBeforeTransfer()
	cd63 := utils.GetChannelBetween(N6, N3, tokenAddress).PrintDataBeforeTransfer()

	// 节点2向节点6转账20token
	N2.SendTrans(tokenAddress, transAmount, N6.Address, false)
	time.Sleep(time.Second * 3)
	//  崩溃判断
	if N2.IsRunning() {
		msg := "Node " + N2.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 中间数据记录
	utils.GetChannelBetween(N3, N2, tokenAddress).PrintDataAfterCrash()
	utils.GetChannelBetween(N6, N3, tokenAddress).PrintDataAfterCrash()

	// 重启节点2，自动发送之前中断的交易
	N2.ReStartWithoutConditionquit(env)

	// 查询结果并校验
	utils.GetChannelBetween(N3, N2, tokenAddress).PrintDataAfterRestart()
	cd63new := utils.GetChannelBetween(N6, N3, tokenAddress).PrintDataAfterRestart()
	if cd63new.Balance-cd63.Balance != transAmount {
		models.Logger.Println("Expect transfer on" + cd63new.Name + " success,but got failed ,FAILED!!!")
		models.Logger.Println(env.CaseName + " END ====> FAILED")
		return fmt.Errorf("Case [%s] FAILED", env.CaseName)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
