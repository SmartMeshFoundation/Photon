package cases

import (
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/Photon/utils"

	"github.com/SmartMeshFoundation/Photon/log"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

// CaseZeroSecret :发送方指定密码为全0来进行交易,是否会造成意外.
func (cm *CaseManager) CaseZeroSecret() (err error) {
	if !cm.RunThisCaseOnly {
		return
	}
	env, err := models.NewTestEnv("./cases/CaseZeroSecret.ENV", cm.UseMatrix, cm.EthEndPoint)
	if err != nil {
		return
	}
	defer func() {
		time.Sleep(time.Minute * 100)
		if env.Debug == false {
			env.KillAllPhotonNodes()
		}
	}()
	// 源数据
	// original data
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N0, N1, N2 := env.Nodes[0], env.Nodes[1], env.Nodes[2]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	cm.startNodes(env, N0, N1, N2)

	// 获取channel信息
	// get channel info
	N0.GetChannelWith(N1, tokenAddress).Println("before transfer")
	secret := utils.EmptyHash
	secretHash := utils.ShaSecret(secret[:])
	N0.SendTransWithSecret(env.Tokens[0].TokenAddress.String(), 1, N2.Address, secret.String())
	if err != nil {
		log.Error(fmt.Sprintf("Transfer err %s", err))
		return
	}
	time.Sleep(time.Second * 2)
	log.Trace("allow reveal secret")
	N0.AllowSecret(secretHash.String(), tokenAddress)
	time.Sleep(time.Second)
	N0.GetChannelWith(N1, tokenAddress).Println("after transfer")
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
