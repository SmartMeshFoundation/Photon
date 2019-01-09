package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

// CrashCaseSend02 场景二：EventSendRevealSecretAfter
// 节点2向节点6转账20token,发送revealsecret后，节点2崩，路由走2-3-6，查询节点6，节点3，交易未完成，锁定节点3 20个token,节点2 20个token，
// 重启后，节点3和节点6的交易完成，节点2和节点3交易完成，交易成功
func (cm *CaseManager) CrashCaseSend02() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCaseSend02.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	transAmount = 20
	tokenAddress := env.Tokens[0].TokenAddress.String()
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
	N3.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	N6.GetChannelWith(N3, tokenAddress).PrintDataBeforeTransfer()

	// 节点2向节点6转账20token
	N2.SendTrans(tokenAddress, transAmount, N6.Address, false)
	//time.Sleep(time.Second * 3)
	//  崩溃判断
	if N2.IsRunning() {
		msg := "Node " + N2.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 中间数据记录
	models.Logger.Println("------------ Data After Crash ------------")
	N3.GetChannelWith(N2, tokenAddress).PrintDataAfterCrash()
	N6.GetChannelWith(N3, tokenAddress).PrintDataAfterCrash()

	// 重启节点2，自动发送之前中断的交易
	N2.ReStartWithoutConditionquit(env)
	for i := 0; i < 30; i++ {
		time.Sleep(time.Second)

		// 查询重启后数据
		models.Logger.Println("------------ Data After Restart ------------")
		cd32new := N3.GetChannelWith(N2, tokenAddress).PrintDataAfterRestart()
		cd63new := N6.GetChannelWith(N3, tokenAddress).PrintDataAfterRestart()
		// 校验对等
		models.Logger.Println("------------ Data After Fail ------------")
		if !cd32new.CheckEqualByPartnerNode(env) || !cd63new.CheckEqualByPartnerNode(env) {
			return cm.caseFail(env.CaseName)
		}
		models.Logger.Println(env.CaseName + " END ====> SUCCESS")
		return
	}
	return cm.caseFail(env.CaseName)
}
