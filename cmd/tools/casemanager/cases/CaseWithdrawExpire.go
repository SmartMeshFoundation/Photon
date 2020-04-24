package cases

import (
	"time"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/go-errors/errors"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

// CaseWithdrawExpire :
func (cm *CaseManager) CaseWithdrawExpire() (err error) {
	env, err := models.NewTestEnv("./cases/CaseWithdrawExpire.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	N0, N1 := env.Nodes[0], env.Nodes[1]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点2，3
	// start node 2, 3
	cm.startNodes(env, N0, N1)

	// 获取channel信息
	// get channel info
	c01 := N0.GetChannelWith(N1, tokenAddress).Println("before withdraw")
	N1.SendTrans(env.Tokens[0].TokenAddress.String(), 1, N0.Address, false)
	time.Sleep(time.Second * 2)
	// n1删除数据并重启
	N1.Shutdown(env)
	N1.ClearHistoryData(env.DataDir)
	N1.Start(env)

	// N0 withdraw
	N0.Withdraw(c01.ChannelIdentifier, 51)
	// 等待N0通道状态恢复
	err = cm.tryInSeconds(50, func() error {
		c01 = N0.GetChannelWith(N1, tokenAddress).Println("after withdraw expire")
		if c01.State == channeltype.StateWithdraw {
			return errors.New("waiting... ")
		} else if c01.State == channeltype.StateOpened {
			return nil
		}
		return errors.New("channel state error")
	})
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
