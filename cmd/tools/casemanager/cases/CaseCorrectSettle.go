package cases

import (
	"fmt"
	"math/big"
	"time"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"

	"github.com/SmartMeshFoundation/Photon/utils"

	"github.com/SmartMeshFoundation/Photon/log"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

func init() {

}

// CaseCorrectSettle :
func (cm *CaseManager) CaseCorrectSettle() (err error) {
	if !cm.RunSlow {
		return
	}
	env, err := models.NewTestEnv("./cases/CaseCorrectSettle.ENV", cm.UseMatrix, cm.EthEndPoint)
	if err != nil {
		return
	}
	defer func() {
		log.Trace(fmt.Sprintf("CaseCorrectSettle err=%s", err))
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
	cm.startNodes(env, N1, N2)
	N0.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "ReceiveSecretRequestStateChange",
	})

	// 获取channel信息
	// get channel info
	c01 := N0.GetChannelWith(N1, tokenAddress).Println("before send tras")
	go N0.SendTrans(env.Tokens[0].TokenAddress.String(), 3, N2.Address, false)
	time.Sleep(3 * time.Second)
	if N0.IsRunning() {
		return cm.caseFailWithWrongChannelData(env.CaseName, "n0 should not running")
	}
	N0.ReStartWithoutConditionquit(env)
	err = N0.Close(c01.ChannelIdentifier)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, fmt.Sprintf("close failed %s", err))
	}

	var i = 0
	for i = 0; i < 100; i++ {
		var c channeltype.ChannelDataDetail
		time.Sleep(time.Second)
		c, err = N0.SpecifiedChannel(c01.ChannelIdentifier)
		if err != nil {
			return cm.caseFailWithWrongChannelData(env.CaseName, "specified channel")
		}
		if c.OurBalanceProof == nil {
			continue
		}
		if c.OurBalanceProof.ContractLocksRoot == utils.EmptyHash {
			continue
		}
		if c.OurBalanceProof.ContractTransferAmount != nil && c.OurBalanceProof.ContractTransferAmount.Cmp(big.NewInt(0)) != 0 {
			continue
		}
		break
	}
	if i == 100 {
		return cm.caseFailWithWrongChannelData(env.CaseName, "n0 spec err")
	}

	settleTime := c01.SettleTimeout + 3600/14
	err = cm.trySettleInSeconds(int(settleTime), N1, c01.ChannelIdentifier)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, c01.Name)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
