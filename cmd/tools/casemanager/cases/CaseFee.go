package cases

import (
	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

// CaseFee : test for fee module
func (cm *CaseManager) CaseFee() (err error) {
	env, err := models.NewTestEnv("./cases/CaseFee.ENV", cm.UseMatrix)
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
	N0, N1, N2, N3 := env.Nodes[0], env.Nodes[1], env.Nodes[2], env.Nodes[3]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点
	N0.Start(env)
	N1.StartWithFee(env)
	N2.StartWithFee(env)
	N3.Start(env)

	cm.logSeparatorLine("Test 1 : transfer with fee 0, should FAIL")
	transferAmount = 10000
	fee = 0
	C01 := N0.GetChannelWith(N1, tokenAddress).PrintDataBeforeTransfer()
	C12 := N1.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	C23 := N2.GetChannelWith(N3, tokenAddress).PrintDataBeforeTransfer()
	N0.SendTransSyncWithFee(tokenAddress, transferAmount, N3.Address, false, fee)
	C01new := N0.GetChannelWith(N1, tokenAddress).PrintDataAfterTransfer()
	C12new := N1.GetChannelWith(N2, tokenAddress).PrintDataAfterTransfer()
	C23new := N2.GetChannelWith(N3, tokenAddress).PrintDataAfterTransfer()
	if !C01new.CheckPartnerBalance(C01.PartnerBalance) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C01new.Name)
	}
	if !C12new.CheckPartnerBalance(C12.PartnerBalance) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C12new.Name)
	}
	if !C23new.CheckPartnerBalance(C23.PartnerBalance) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C23new.Name)
	}

	cm.logSeparatorLine("Test 2 : transfer with fee 2, should SUCCESS")
	transferAmount = 10000
	fee = 2
	C01 = C01new
	C12 = C12new
	C23 = C23new
	N0.SendTransSyncWithFee(tokenAddress, transferAmount, N3.Address, false, fee)
	C01new = N0.GetChannelWith(N1, tokenAddress).PrintDataAfterTransfer()
	C12new = N1.GetChannelWith(N2, tokenAddress).PrintDataAfterTransfer()
	C23new = N2.GetChannelWith(N3, tokenAddress).PrintDataAfterTransfer()
	if !C01new.CheckPartnerBalance(C01.PartnerBalance + transferAmount + 2) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C01new.Name)
	}
	if !C12new.CheckPartnerBalance(C12.PartnerBalance + transferAmount + 1) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C12new.Name)
	}
	if !C23new.CheckPartnerBalance(C23.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C23new.Name)
	}

	cm.logSeparatorLine("Test 3 : transfer with fee 3, should SUCCESS")
	transferAmount = 10000
	fee = 3
	C01 = C01new
	C12 = C12new
	C23 = C23new
	N0.SendTransSyncWithFee(tokenAddress, transferAmount, N3.Address, false, fee)
	C01new = N0.GetChannelWith(N1, tokenAddress).PrintDataAfterTransfer()
	C12new = N1.GetChannelWith(N2, tokenAddress).PrintDataAfterTransfer()
	C23new = N2.GetChannelWith(N3, tokenAddress).PrintDataAfterTransfer()
	if !C01new.CheckPartnerBalance(C01.PartnerBalance + transferAmount + 3) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C01new.Name)
	}
	if !C12new.CheckPartnerBalance(C12.PartnerBalance + transferAmount + 2) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C12new.Name)
	}
	if !C23new.CheckPartnerBalance(C23.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C23new.Name)
	}

	cm.logSeparatorLine("Test 4 : transfer with fee 5 and transferAmount < 10000, should SUCCESS")
	transferAmount = 5000
	fee = 5
	C01 = C01new
	C12 = C12new
	C23 = C23new
	N0.SendTransSyncWithFee(tokenAddress, transferAmount, N3.Address, false, fee)
	C01new = N0.GetChannelWith(N1, tokenAddress).PrintDataAfterTransfer()
	C12new = N1.GetChannelWith(N2, tokenAddress).PrintDataAfterTransfer()
	C23new = N2.GetChannelWith(N3, tokenAddress).PrintDataAfterTransfer()
	if !C01new.CheckPartnerBalance(C01.PartnerBalance + transferAmount + 5) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C01new.Name)
	}
	if !C12new.CheckPartnerBalance(C12.PartnerBalance + transferAmount + 5) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C12new.Name)
	}
	if !C23new.CheckPartnerBalance(C23.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C23new.Name)
	}

	cm.logSeparatorLine("Test 5 : transfer with fee 0 and transferAmount < 10000, should SUCCESS")
	transferAmount = 5000
	fee = 0
	C01 = C01new
	C12 = C12new
	C23 = C23new
	N0.SendTransSyncWithFee(tokenAddress, transferAmount, N3.Address, false, fee)
	C01new = N0.GetChannelWith(N1, tokenAddress).PrintDataAfterTransfer()
	C12new = N1.GetChannelWith(N2, tokenAddress).PrintDataAfterTransfer()
	C23new = N2.GetChannelWith(N3, tokenAddress).PrintDataAfterTransfer()
	if !C01new.CheckPartnerBalance(C01.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C01new.Name)
	}
	if !C12new.CheckPartnerBalance(C12.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C12new.Name)
	}
	if !C23new.CheckPartnerBalance(C23.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C23new.Name)
	}

	cm.logSeparatorLine("Test 6 : transfer with fee 20000 and transferAmount = 10000, should SUCCESS")
	transferAmount = 10000
	fee = 20000
	C01 = C01new
	C12 = C12new
	C23 = C23new
	N0.SendTransSyncWithFee(tokenAddress, transferAmount, N3.Address, false, fee)
	C01new = N0.GetChannelWith(N1, tokenAddress).PrintDataAfterTransfer()
	C12new = N1.GetChannelWith(N2, tokenAddress).PrintDataAfterTransfer()
	C23new = N2.GetChannelWith(N3, tokenAddress).PrintDataAfterTransfer()
	if !C01new.CheckPartnerBalance(C01.PartnerBalance + transferAmount + 20000) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C01new.Name)
	}
	if !C12new.CheckPartnerBalance(C12.PartnerBalance + transferAmount + 19999) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C12new.Name)
	}
	if !C23new.CheckPartnerBalance(C23.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C23new.Name)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
