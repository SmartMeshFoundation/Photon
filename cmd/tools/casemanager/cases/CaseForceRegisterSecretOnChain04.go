package cases

import (
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/utils"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

// CaseForceRegisterSecretOnChain04 :
func (cm *CaseManager) CaseForceRegisterSecretOnChain04() (err error) {
	if !cm.RunSlow {
		return ErrorSkip
	}
	env, err := models.NewTestEnv("./cases/CaseForceRegisterSecretOnChain04.ENV", cm.UseMatrix, cm.EthEndPoint)
	if err != nil {
		return
	}
	defer func() {
		log.Trace(fmt.Sprintf("CaseForceRegisterSecretOnChain04 err=%s", err))
		if env.Debug == false {
			env.KillAllPhotonNodes()
		}
	}()
	// 源数据
	// original data
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N0, N1 := env.Nodes[0], env.Nodes[1]
	models.Logger.Println(env.CaseName + " BEGIN ====>")

	cm.startNodes(env, N1)
	N0.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "ReceiveSecretRevealStateChange",
	})

	// 获取channel信息
	// get channel info
	c01 := N0.GetChannelWith(N1, tokenAddress).Println("before send tras")

	go N0.SendTrans(env.Tokens[0].TokenAddress.String(), 3, N1.Address, false)
	time.Sleep(3 * time.Second)
	err = N1.Close(c01.ChannelIdentifier)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, fmt.Sprintf("close failed %s", err))
	}
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
