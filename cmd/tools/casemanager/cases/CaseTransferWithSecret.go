package cases

import (
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

// CaseTransferWithSecret :
func (cm *CaseManager) CaseTransferWithSecret() (err error) {
	env, err := models.NewTestEnv("./cases/CaseTransferWithSecret.ENV", cm.UseMatrix, cm.EthEndPoint)
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

	secret, SecretHash, err := N0.GenerateSecret()
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	// 获取channel信息
	// get channel info
	c01 := N0.GetChannelWith(N1, tokenAddress).Println("BeforeSendTransWithSecret")
	go N0.SendTransWithSecret(tokenAddress, 1, N1.Address, secret)
	time.Sleep(3 * time.Second)
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
	//没有发送密码允许,对方肯定接收不到
	c01new := N0.GetChannelWith(N1, tokenAddress).Println("after send transfer with secret")
	if c01new.CheckSelfBalance(c01.Balance - 1) {
		return cm.caseFail(env.CaseName)
	}
	if c01new.CheckPartnerBalance(c01.PartnerBalance + 1) {
		return cm.caseFail(env.CaseName)
	}
	N0.AllowSecret(SecretHash, tokenAddress)

	for i := 0; i < cm.HighMediumWaitSeconds; i++ {
		c01new = N0.GetChannelWith(N1, tokenAddress).Println("after  allow reveal secret")
		time.Sleep(time.Second) //保证photon在十秒之内会尝试发送一次消息
		if !c01new.CheckSelfBalance(c01.Balance - 1) {
			continue
		}
		if !c01new.CheckPartnerBalance(c01.PartnerBalance + 1) {
			continue
		}
		models.Logger.Println(env.CaseName + " END ====> SUCCESS")
		return
	}
	return cm.caseFail(env.CaseName)
}
