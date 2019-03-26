package cases

import (
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/params"
)

// CaseCrashEventClose n0找n1打开通道并存款,n0崩溃条件EventNewChannelFromChainBeforeDeal, 恢复后验证双方状态
func (cm *CaseManager) CaseCrashEventClose() (err error) {
	if !cm.RunSlow {
		return ErrorSkip
	}
	env, err := models.NewTestEnv("./cases/CaseCrashEventClose.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	cm.startNodes(env, n1,
		n0.SetConditionQuit(&params.ConditionQuit{
			QuitEvent: "EventChannelCloseFromChainBeforeDeal",
		}),
	)
	// 1. open
	err = n0.OpenChannel(n1.Address, tokenAddress, depositAmount, settleTimeout)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	// 2. 校验双方通道状态
	c01 := n0.GetChannelWith(n1, tokenAddress)
	if c01 == nil {
		return cm.caseFail(env.CaseName)
	}
	if !c01.CheckEqualByPartnerNode(env) {
		return cm.caseFailWithWrongChannelData(env.CaseName, c01.Name)
	}
	// 3. N0 Close
	err = n0.Close(c01.ChannelIdentifier, 0)
	if err == nil {
		return cm.caseFail(env.CaseName)
	}
	// 3. 验证n0崩溃
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
	// 4.重启n0
	n0.ReStartWithoutConditionquit(env)
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
	// 5. 校验双方通道状态
	c01 = n0.GetChannelWith(n1, tokenAddress).Println("after restart")
	if c01 == nil {
		return cm.caseFail(env.CaseName)
	}
	if !c01.CheckEqualByPartnerNode(env) {
		return cm.caseFailWithWrongChannelData(env.CaseName, c01.Name)
	}
	if !c01.CheckState(channeltype.StateClosed) {
		return cm.caseFailWithWrongChannelData(env.CaseName, c01.Name)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
