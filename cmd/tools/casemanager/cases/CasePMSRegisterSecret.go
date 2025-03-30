package cases

import (
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

// CasePMSRegisterSecret :
func (cm *CaseManager) CasePMSRegisterSecret() (err error) {
	if cm.IsAutoRun {
		return ErrorSkip
	}
	env, err := models.NewTestEnv("./cases/CasePMSRegisterSecret.ENV", cm.UseMatrix, cm.EthEndPoint, "CasePMSRegisterSecret")
	if err != nil {
		return
	}
	defer func() {
		if env.Debug == false {
			env.KillAllPhotonNodes()
		}
	}()

	// 源数据
	// original data
	tokenAddress := env.Tokens[0].TokenAddress.String()

	N0, N1, N2 := env.Nodes[0], env.Nodes[1], env.Nodes[2]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动受托节点n1-pms
	env.StartPMS()
	cm.startNodes(env, N0.SetConditionQuit(&params.ConditionQuit{
		QuitEvent: "ReceiveSecretRevealStateChange",
	}), N1.PMS(), N2.PMS())

	transAmount := int32(20)
	// 获取channel信息
	N1.GetChannelWith(N0, tokenAddress).Println("before send trans")
	N2.GetChannelWith(N1, tokenAddress).Println("before send trans")

	// 转账
	go N0.SendTrans(tokenAddress, transAmount, N2.Address, false)

	// 崩溃判断
	for i := 0; i < cm.MediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		if !N0.IsRunning() {
			break
		}
	}

	N1.GetChannelWith(N0, tokenAddress).Println("after send trans")
	N2.GetChannelWith(N1, tokenAddress).Println("after send trans")

	cm.waitForPostman()
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
