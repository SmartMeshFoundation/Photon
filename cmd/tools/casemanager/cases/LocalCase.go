package cases

import (
	"strconv"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

// LocalCase : only for local test
func (cm *CaseManager) LocalCase() (err error) {
	env, err := models.NewTestEnv("./cases/LocalCase.ENV")
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
	times := 0
	name0 := N0.Name
	name1 := N1.Name
	for {
		if times > 100 {
			break
		}
		N0.Name = name0 + "-" + strconv.Itoa(times+1)
		N1.Name = name1 + "-" + strconv.Itoa(times+1)
		// 启动节点,让节点0在收到
		N0.Start(env)
		N1.Start(env)
		c01 := N0.GetChannelWith(N1, tokenAddress)
		N0.SendTransWithData(tokenAddress, transAmount, N1.Address, true, "123")
		c01new := N0.GetChannelWith(N1, tokenAddress)
		if !c01new.CheckPartnerBalance(c01.PartnerBalance + transAmount) {
			return cm.caseFailWithWrongChannelData(env.CaseName, c01new.Name)
		}
		N0.Shutdown(env)
		N1.Shutdown(env)
		c01new.PrintDataAfterTransfer()
		models.Logger.Printf("============== time %d done\n", times+1)
		times++
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
