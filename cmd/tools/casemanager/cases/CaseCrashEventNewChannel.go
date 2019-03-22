package cases

import (
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/params"
)

// CaseCrashEventNewChannel n0找n1打开通道并存款,n0崩溃条件EventNewChannelFromChainBeforeDeal, 恢复后验证双方状态
func (cm *CaseManager) CaseCrashEventNewChannel() (err error) {
	if !cm.RunSlow {
		return ErrorSkip
	}
	env, err := models.NewTestEnv("./cases/CaseCrashEventNewChannel.ENV", cm.UseMatrix, cm.EthEndPoint)
	if err != nil {
		return
	}
	defer func() {
		log.Trace(fmt.Sprintf("CaseCrashEventNewChannel err=%s", err))
		if env.Debug == false {
			env.KillAllPhotonNodes()
		}
	}()
	// 源数据
	// original data
	tokenAddress := env.Tokens[0].TokenAddress.String()
	n0, n1 := env.Nodes[0], env.Nodes[1]
	depositAmount := int64(100)
	settleTimeout := int64(100)
	models.Logger.Println(env.CaseName + " BEGIN ====>")

	// 0. start
	n0.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "EventNewChannelFromChainBeforeDeal",
	})
	n1.Start(env)
	// 1. open
	err = n0.OpenChannel(n1.Address, tokenAddress, depositAmount, settleTimeout, 0)
	if err == nil {
		return cm.caseFail(env.CaseName)
	}
	// 2. 验证n0崩溃
	i := 0
	for i = 0; i < cm.HighMediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		if !n0.IsRunning() {
			break
		}
	}
	if i == cm.HighMediumWaitSeconds {
		return cm.caseFail(env.CaseName)
	}
	// 3.重启n0
	n0.ReStartWithoutConditionquit(env)
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
	// 4. 校验双方通道状态
	c01 := n0.GetChannelWith(n1, tokenAddress)
	if c01 == nil {
		return cm.caseFail(env.CaseName)
	}
	if !c01.CheckEqualByPartnerNode(env) {
		return cm.caseFailWithWrongChannelData(env.CaseName, c01.Name)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
