package cases

import (
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/Photon/utils"

	"github.com/SmartMeshFoundation/Photon/log"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

// CaseForceRegisterSecretOnChain03 :
func (cm *CaseManager) CaseForceRegisterSecretOnChain03() (err error) {
	if !cm.RunSlow {
		return ErrorSkip
	}
	env, err := models.NewTestEnv("./cases/CaseForceRegisterSecretOnChain03.ENV", cm.UseMatrix, cm.EthEndPoint)
	if err != nil {
		return
	}
	defer func() {
		log.Trace(fmt.Sprintf("CaseForceRegisterSecretOnChain03 err=%s", err))
		if env.Debug == false {
			env.KillAllPhotonNodes()
		}
	}()
	// 源数据
	// original data
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N0, N1, N2 := env.Nodes[0], env.Nodes[1], env.Nodes[2]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点2，3
	// start node 2, 3
	cm.startNodes(env, N1, N2,
		N0.SetConditionQuit(&params.ConditionQuit{
			QuitEvent: "ReceiveSecretRevealStateChange",
		}))

	// 获取channel信息
	// get channel info
	c01 := N0.GetChannelWith(N1, tokenAddress).Println("before send tras")
	N1.GetChannelWith(N2, tokenAddress).Println("before send  trans")

	go N0.SendTrans(env.Tokens[0].TokenAddress.String(), 3, N2.Address, false)
	time.Sleep(3 * time.Second)
	var i = 0
	settleTime := c01.SettleTimeout + 3600/14
	for i = 0; i < int(settleTime); i++ {
		time.Sleep(time.Second)
		c, err := N1.SpecifiedChannel(c01.ChannelIdentifier)
		log.Trace(fmt.Sprintf("c=%s,err=%s", utils.StringInterface(c, 3), err))
		if err != nil {
			continue
		}
		if len(c.PartnerKnownSecretLocks) != 1 {
			continue
		}
		reg := false
		for _, s := range c.PartnerKnownSecretLocks {
			if s.IsRegisteredOnChain {
				reg = true
				break
			}
		}
		if !reg {
			continue
		}
		break
	}
	if i == int(settleTime) {
		return cm.caseFailWithWrongChannelData(env.CaseName, c01.Name)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
