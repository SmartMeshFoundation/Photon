package cases

import (
	"time"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/utils"
)

// CaseEthNonce :
func (cm *CaseManager) CaseEthNonce() (err error) {
	if !cm.RunSlow {
		return
	}
	env, err := models.NewTestEnv("./cases/CaseEthNonce.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	settleTimeout := int64(120)
	N0 := env.Nodes[0]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点2，3
	// start node 2, 3
	N0.Start(env)
	time.Sleep(time.Second * 3)
	// 获取channel信息
	for i := 0; i < 10; i++ {
		go func() {
			err2 := N0.OpenChannel(utils.NewRandomAddress().String(), tokenAddress, 1, settleTimeout)
			if err2 != nil {
				fmt.Printf("----------err : %s \n", err2.Error())
			}
		}()
	}
	for i := 0; i < cm.MediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		channels := N0.GetChannels(tokenAddress)
		if len(channels) >= 10 {
			//time.Sleep(time.Second * 5)
			models.Logger.Println(env.CaseName + " END ====> SUCCESS")
			return nil
		}
	}
	return cm.caseFail(env.CaseName)

}
