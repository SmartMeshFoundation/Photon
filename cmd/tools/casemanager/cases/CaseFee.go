package cases

import (
	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

// CaseFee : test for fee module
func (cm *CaseManager) CaseFee() (err error) {
	env, err := models.NewTestEnv("./cases/CaseFee.ENV")
	if err != nil {
		return
	}
	defer func() {
		if env.Debug == false {
			env.KillAllPhotonNodes()
		}
	}()
	// 源数据
	var transferAmount int32
	var fee int64
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N0, N1, N2 := env.Nodes[0], env.Nodes[1], env.Nodes[2]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点
	N0.Start(env)
	N1.StartWithFee(env)
	N2.Start(env)

	cm.logSeparatorLine("Test 1 : transfer with fee 0, should FAIL")
	transferAmount = 10000
	fee = 0
	C01 := N0.GetChannelWith(N1, tokenAddress).PrintDataBeforeTransfer()
	C12 := N1.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	N0.SendTransSyncWithFee(tokenAddress, transferAmount, N2.Address, false, fee)
	C01new := N0.GetChannelWith(N1, tokenAddress).PrintDataAfterTransfer()
	C12new := N1.GetChannelWith(N2, tokenAddress).PrintDataAfterTransfer()
	if !C01new.CheckPartnerBalance(C01.PartnerBalance) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C01new.Name)
	}
	if !C12new.CheckPartnerBalance(C12.PartnerBalance) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C12new.Name)
	}

	cm.logSeparatorLine("Test 2 : transfer with fee 1, should SUCCESS")
	transferAmount = 10000
	fee = 1
	C01 = C01new
	C12 = C12new
	N0.SendTransSyncWithFee(tokenAddress, transferAmount, N2.Address, false, fee)
	C01new = N0.GetChannelWith(N1, tokenAddress).PrintDataAfterTransfer()
	C12new = N1.GetChannelWith(N2, tokenAddress).PrintDataAfterTransfer()
	if !C01new.CheckPartnerBalance(C01.PartnerBalance + transferAmount + int32(fee)) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C01new.Name)
	}
	if !C12new.CheckPartnerBalance(C12.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C12new.Name)
	}

	cm.logSeparatorLine("Test 3 : transfer with fee 2, should SUCCESS")
	transferAmount = 10000
	fee = 2
	C01 = C01new
	C12 = C12new
	N0.SendTransSyncWithFee(tokenAddress, transferAmount, N2.Address, false, fee)
	C01new = N0.GetChannelWith(N1, tokenAddress).PrintDataAfterTransfer()
	C12new = N1.GetChannelWith(N2, tokenAddress).PrintDataAfterTransfer()
	if !C01new.CheckPartnerBalance(C01.PartnerBalance + transferAmount + int32(fee)) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C01new.Name)
	}
	if !C12new.CheckPartnerBalance(C12.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C12new.Name)
	}

	cm.logSeparatorLine("Test 4 : transfer with fee 2 and transferAmount < 10000, should SUCCESS")
	transferAmount = 5000
	fee = 2
	C01 = C01new
	C12 = C12new
	N0.SendTransSyncWithFee(tokenAddress, transferAmount, N2.Address, false, fee)
	C01new = N0.GetChannelWith(N1, tokenAddress).PrintDataAfterTransfer()
	C12new = N1.GetChannelWith(N2, tokenAddress).PrintDataAfterTransfer()
	if !C01new.CheckPartnerBalance(C01.PartnerBalance + transferAmount + int32(fee)) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C01new.Name)
	}
	if !C12new.CheckPartnerBalance(C12.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C12new.Name)
	}

	cm.logSeparatorLine("Test 5 : transfer with fee 0 and transferAmount < 10000, should SUCCESS")
	transferAmount = 5000
	fee = 0
	C01 = C01new
	C12 = C12new
	N0.SendTransSyncWithFee(tokenAddress, transferAmount, N2.Address, false, fee)
	C01new = N0.GetChannelWith(N1, tokenAddress).PrintDataAfterTransfer()
	C12new = N1.GetChannelWith(N2, tokenAddress).PrintDataAfterTransfer()
	if !C01new.CheckPartnerBalance(C01.PartnerBalance + transferAmount + int32(fee)) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C01new.Name)
	}
	if !C12new.CheckPartnerBalance(C12.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C12new.Name)
	}

	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
