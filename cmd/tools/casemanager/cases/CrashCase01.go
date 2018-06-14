package cases

import (
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/utils"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

// CrashCase01 场景一：EventSendMediatedTransferAfter
// 发送中转转账后崩溃
// 节点1向节点2发送MTR后，节点1崩溃，此时，节点2默认收到MTR，但由于没有ACK确认，没发生转账，余额不变。节点2没收到转账token.
// 重启节点1后，继续转账，转账成功。
func (cm *CaseManager) CrashCase01() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCase01.ENV")
	if err != nil {
		return
	}
	// 源数据
	var transAmount int32
	transAmount = 5
	tokenAddress := env.Tokens[0].Address
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 1. /启动节点0 with EventSendMediatedTransferAfter
	N0 := env.Nodes[0]
	N0.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "EventSendMediatedTransferAfter",
	})

	// 2. 启动节点1
	N1 := env.Nodes[1]
	N1.Start(env)

	// 3. 获取0-1之间的channel，记录channel数据d1
	cd1 := utils.GetChannelBetween(N1, N0, tokenAddress)
	if cd1 == nil {
		return fmt.Errorf("Can not find channel between %s and %s on token %s", N0.Name, N1.Name, tokenAddress)
	}
	cd1.Println("Channel data before transfer send, cd1:")

	// 4. 从节点0发起到节点1的转账
	N0.SendTrans(tokenAddress, transAmount, N1.Address, false)

	// 5. 记录channel数据d2并与d1比对，assert(d1==d2)
	cd2 := utils.GetChannelBetween(N1, N0, tokenAddress)
	cd2.Println("Channel data after transfer send, cd2:")
	if !utils.IsEqualChannelData(cd1, cd2) {
		models.Logger.Println("Expect cd1 == cd2 but got cd1 != cd2, FAILED !!!")
		return fmt.Errorf("Expect cd1 == cd2 but got cd1 != cd2")
	}

	// 6. 重启节点1，自动发送之前中断的交易
	if N0.IsRunning() {
		models.Logger.Println("Expect cd1 == cd2 but got cd1 != cd2, FAILED !!!")
		return fmt.Errorf("Node N0 should be exited,but it still running")
	}
	N0.DebugCrash = false
	N0.ConditionQuit = nil
	N0.Name = "RestartNode"
	N0.Start(env)

	// 7. 记录channel数据d3, assert(余额校验)
	cd3 := utils.GetChannelBetween(N0, N1, tokenAddress)
	if cd3.SelfAddress != cd2.SelfAddress {
		utils.SwitchChannel(cd3)
	}
	cd3.Println("Channel data after N0 restart, cd3:")
	if cd3.Balance-cd2.Balance == transAmount {
		models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	} else {
		models.Logger.Println(env.CaseName + " END ====> FAILED")
		return fmt.Errorf("Case [%s] FAILED", env.CaseName)
	}
	return
}
