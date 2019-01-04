package cases

import (
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

// CaseFeeWithPFS : test for fee module with pfs
func (cm *CaseManager) CaseFeeWithPFS() (err error) {
	if cm.IsAutoRun {
		return
	}
	env, err := models.NewTestEnv("./cases/CaseFeeWithPFS.ENV", cm.UseMatrix)
	if err != nil {
		return
	}
	defer func() {
		if env.Debug == false {
			env.KillAllPhotonNodes()
		}
	}()
	// 源数据
	params := []string{
		"--fee", "--pfs=http://transport01.smartmesh.cn:7002",
		//"--fee", "--pfs=http://192.168.124.9:7000",
	}
	var transferAmount int32
	var fee int64
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N0, N1, N2, N3, N4, N5 := env.Nodes[0], env.Nodes[1], env.Nodes[2], env.Nodes[3], env.Nodes[4], env.Nodes[5]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点
	N0.StartWithParams(env, params...)
	N1.StartWithParams(env, params...)
	N2.StartWithParams(env, params...)
	N3.StartWithParams(env, params...)
	N4.StartWithParams(env, params...)
	N5.StartWithParams(env, params...)

	cm.logSeparatorLine("记录初始数据")
	C01 := N0.GetChannelWith(N1, tokenAddress) // C01 50000,50000
	C02 := N0.GetChannelWith(N2, tokenAddress) // C02 50000,50000
	C15 := N1.GetChannelWith(N5, tokenAddress) // C15 50000,50000
	C23 := N2.GetChannelWith(N3, tokenAddress) // C23 50000,50000
	C24 := N2.GetChannelWith(N4, tokenAddress) // c01 50000,50000
	C35 := N3.GetChannelWith(N5, tokenAddress) // C35 50000,50000
	C45 := N4.GetChannelWith(N5, tokenAddress) // C45 50000,50000
	C01.Println("初始")
	C02.Println("初始")
	C15.Println("初始")
	C23.Println("初始")
	C24.Println("初始")
	C35.Println("初始")
	C45.Println("初始")
	time.Sleep(30 * time.Second)

	cm.logSeparatorLine("Test 1 : transfer on route 0-1-5 with fee 1, should SUCCESS")
	transferAmount = 10000
	fee = 1
	N0.SendTransSyncWithFee(tokenAddress, transferAmount, N5.Address, false, fee)
	C01new := N0.GetChannelWith(N1, tokenAddress).PrintDataAfterTransfer() // C01 39999,60001
	C15new := N1.GetChannelWith(N5, tokenAddress).PrintDataAfterTransfer() // C15 40000,60000
	if !C01new.CheckPartnerBalance(C01.PartnerBalance + transferAmount + int32(fee)) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C01new.Name)
	}
	if !C15new.CheckPartnerBalance(C15.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C15new.Name)
	}
	time.Sleep(1000 * time.Second)
	//cm.logSeparatorLine("Test 1 : transfer with fee 0, should FAIL")
	//transferAmount = 10000
	//fee = 0
	//C01 := N0.GetChannelWith(N1, tokenAddress).PrintDataBeforeTransfer()
	//C12 := N1.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	//N0.SendTransSyncWithFee(tokenAddress, transferAmount, N2.Address, false, fee)
	//C01new := N0.GetChannelWith(N1, tokenAddress).PrintDataAfterTransfer()
	//C12new := N1.GetChannelWith(N2, tokenAddress).PrintDataAfterTransfer()
	//if !C01new.CheckPartnerBalance(C01.PartnerBalance) {
	//	return cm.caseFailWithWrongChannelData(env.CaseName, C01new.Name)
	//}
	//if !C12new.CheckPartnerBalance(C12.PartnerBalance) {
	//	return cm.caseFailWithWrongChannelData(env.CaseName, C12new.Name)
	//}
	//
	//cm.logSeparatorLine("Test 2 : transfer with fee 1, should SUCCESS")
	//transferAmount = 10000
	//fee = 1
	//C01 = C01new
	//C12 = C12new
	//N0.SendTransSyncWithFee(tokenAddress, transferAmount, N2.Address, false, fee)
	//C01new = N0.GetChannelWith(N1, tokenAddress).PrintDataAfterTransfer()
	//C12new = N1.GetChannelWith(N2, tokenAddress).PrintDataAfterTransfer()
	//if !C01new.CheckPartnerBalance(C01.PartnerBalance + transferAmount + int32(fee)) {
	//	return cm.caseFailWithWrongChannelData(env.CaseName, C01new.Name)
	//}
	//if !C12new.CheckPartnerBalance(C12.PartnerBalance + transferAmount) {
	//	return cm.caseFailWithWrongChannelData(env.CaseName, C12new.Name)
	//}

	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
