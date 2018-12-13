package cases

import (
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

// CaseMatrixResend :
func (cm *CaseManager) CaseMatrixResend() (err error) {
	if cm.IsAutoRun {
		return
	}
	env, err := models.NewTestEnv("./cases/CaseMatrixResend.ENV")
	if err != nil {
		return
	}
	defer func() {
		if env.Debug == false {
			env.KillAllPhotonNodes()
		}
	}()
	// 源数据
	transAmount := int32(1)
	tokenAddress := env.Tokens[0].TokenAddress.String()
	//number := 1000
	N0, N1 := env.Nodes[0], env.Nodes[1]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点,让节点0在收到
	N0.Start(env)
	c01 := N0.GetChannelWith(N1, tokenAddress).PrintDataBeforeTransfer()
	go N0.SendTransWithData(tokenAddress, transAmount, N1.Address, true, "123")
	time.Sleep(60 * time.Second)
	N1.Start(env)
	c10new := N1.GetChannelWith(N0, tokenAddress).PrintDataAfterTransfer()
	if !c10new.CheckSelfBalance(c01.PartnerBalance + transAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, c10new.Name)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
