package cases

import (
	"time"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/utils"
)

// CaseEthNonce :
func (cm *CaseManager) CaseEthNonce() (err error) {
	env, err := models.NewTestEnv("./cases/CaseEthNonce.ENV", cm.UseMatrix)
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

	// 获取channel信息
	for i := 0; i < 100; i++ {
		go func() {
			err = N0.OpenChannel(utils.NewRandomAddress().String(), tokenAddress, 0, settleTimeout)
			if err != nil {
				fmt.Printf("----------err : %s \n", err.Error())
			}
		}()
	}
	time.Sleep(100 * time.Second)
	channels := N0.GetChannels(tokenAddress)
	if len(channels) != 100 {
		fmt.Println("channel num :", len(channels))
		return cm.caseFail(env.CaseName)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
