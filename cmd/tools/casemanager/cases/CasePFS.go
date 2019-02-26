package cases

import (
	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

// CasePFS :
func (cm *CaseManager) CasePFS() (err error) {
	if !cm.RunThisCaseOnly {
		return
	}
	env, err := models.NewTestEnv("./cases/CasePFS.ENV", cm.UseMatrix, cm.EthEndPoint)
	if err != nil {
		return
	}
	defer func() {
		if env.Debug == false {
			env.KillAllPhotonNodes()
		}
	}()
	// 源数据
	tokenAddress := env.Tokens[0].TokenAddress.String()
	n0, n1, n2 := env.Nodes[0], env.Nodes[1], env.Nodes[2]
	// 0. 启动pfs
	env.StartPFS()
	// 1. 启动节点
	n0.StartWithFeeAndPFS(env)
	n1.StartWithFeeAndPFS(env)
	n2.StartWithFeeAndPFS(env)
	// 2. 交易
	c01 := n0.GetChannelWith(n1, tokenAddress).Println("before transfer")
	c12 := n1.GetChannelWith(n2, tokenAddress).Println("before transfer")
	transferAmount := int32(10000)
	n0.SendTrans(tokenAddress, transferAmount, n2.Address, false)
	// 3. 余额验证
	err = cm.tryInSeconds(cm.HighMediumWaitSeconds, func() error {
		c01new := n0.GetChannelWith(n1, tokenAddress).Println("after transfer")
		c12new := n1.GetChannelWith(n2, tokenAddress).Println("after transfer")
		if !c01new.CheckEqualByPartnerNode(env) {
			return cm.caseFailWithWrongChannelData(env.CaseName, c01new.Name)
		}
		if !c01new.CheckSelfBalance(c01.Balance - transferAmount - transferAmount/10000) {
			return cm.caseFailWithWrongChannelData(env.CaseName, c01new.Name)
		}
		if !c12new.CheckEqualByPartnerNode(env) {
			return cm.caseFailWithWrongChannelData(env.CaseName, c12new.Name)
		}
		if !c12new.CheckSelfBalance(c12.Balance - transferAmount) {
			return cm.caseFailWithWrongChannelData(env.CaseName, c12new.Name)
		}
		return nil
	})
	if err != nil {
		return
	}
	// 3. 查询
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
