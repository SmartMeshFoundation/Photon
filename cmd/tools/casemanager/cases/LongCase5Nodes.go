package cases

import (
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/network/netshare"
	"github.com/SmartMeshFoundation/Photon/utils"
)

// LongCase5Nodes :
func (cm *CaseManager) LongCase5Nodes() (err error) {
	env, err := models.NewTestEnv("./cases/LongCase5Nodes.ENV")
	if err != nil {
		return
	}
	defer func() {
		if env.Debug == false {
			env.KillAllPhotonNodes()
		}
	}()
	// 源数据
	settleTimeout := int64(120)
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N0, N1, N2, N3, N4 := env.Nodes[0], env.Nodes[1], env.Nodes[2], env.Nodes[3], env.Nodes[4]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// step 1 : Start 5 Photon nodes
	models.Logger.Println("step 1 ---->")
	N0.Start(env)
	N1.Start(env)
	N2.Start(env)
	N3.Start(env)
	N4.Start(env)

	// step 2 : Create the following channels: N0 - N1, N1 - N2, N2 - N3 with 0 deposit
	models.Logger.Println("step 2 ---->")
	err = N0.OpenChannel(N1.Address, tokenAddress, 0, settleTimeout)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	C01 := N0.GetChannelWith(N1, tokenAddress)
	if C01 == nil {
		return cm.caseFail(env.CaseName)
	}
	err = N1.OpenChannel(N2.Address, tokenAddress, 0, settleTimeout)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	C12 := N1.GetChannelWith(N2, tokenAddress)
	if C12 == nil {
		return cm.caseFail(env.CaseName)
	}
	err = N2.OpenChannel(N3.Address, tokenAddress, 0, settleTimeout)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	C23 := N2.GetChannelWith(N3, tokenAddress)
	if C23 == nil {
		return cm.caseFail(env.CaseName)
	}

	// step 3 : N0 N1 N2 N3 make a deposit of 100 to their channels
	models.Logger.Println("step 3 ---->")
	depositAmount := int64(100)
	err = N0.Deposit(C01.ChannelIdentifier, depositAmount)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	err = N1.Deposit(C01.ChannelIdentifier, depositAmount)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	err = N1.Deposit(C12.ChannelIdentifier, depositAmount)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	err = N2.Deposit(C12.ChannelIdentifier, depositAmount)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	err = N2.Deposit(C23.ChannelIdentifier, depositAmount)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	err = N3.Deposit(C23.ChannelIdentifier, depositAmount)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}

	// step 4 : N4 tries to make a deposit to a channel that does not exist (fail channel with N3 does not exist)
	models.Logger.Println("step 4 ---->")
	err = N4.Deposit(utils.NewRandomHash().String(), depositAmount)
	if err == nil {
		return cm.caseFail(env.CaseName)
	}
	err = nil

	// step 5 : N1 makes a deposit of 50 tokens on both of his channels (skip)
	// step 6 : N1 tries to open a channel with N0, but it already has a channel (fail)
	models.Logger.Println("step 5 ---->")
	err = N1.OpenChannel(N0.Address, tokenAddress, 0, settleTimeout)
	if err == nil {
		return cm.caseFail(env.CaseName)
	}
	err = nil

	// step 7 : N2 opens a channel with N4 (0 tokens)
	models.Logger.Println("step 7 ---->")
	err = N2.OpenChannel(N4.Address, tokenAddress, 0, settleTimeout)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	C24 := N2.GetChannelWith(N4, tokenAddress)
	if C24 == nil {
		return cm.caseFail(env.CaseName)
	}

	// step 8 : N2 makes a deposit of 100 tokens
	models.Logger.Println("step 8 ---->")
	err = N2.Deposit(C24.ChannelIdentifier, depositAmount)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}

	// step 9 : N2 tries to deposit 30 tokens (skip)
	// step 10 : N2 deposits 50 tokens to the channel (N2 - N4)
	models.Logger.Println("step 10 ---->")
	depositAmount = 50
	err = N2.Deposit(C24.ChannelIdentifier, depositAmount)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}

	// step 11 : N0 tries to open a channel with an initial deposit that is bigger then the Red Eyes Limit (skip)
	// step 12 : N0 opens a channel with N4 (initial deposit of 10)
	models.Logger.Println("step 12 ---->")
	err = N0.OpenChannel(N4.Address, tokenAddress, 0, settleTimeout)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	C04 := N0.GetChannelWith(N4, tokenAddress)
	if C04 == nil {
		return cm.caseFail(env.CaseName)
	}
	time.Sleep(5 * time.Second)

	// step 13 : N4 deposits 25 tokens to N0<->N4 channel
	models.Logger.Println("step 13 ---->")
	depositAmount = 25
	err = N4.Deposit(C04.ChannelIdentifier, depositAmount)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}

	// step 14 : N0 performs a payment to N4 of 50 tokens (path N0<->N1<->N2<->N4)
	models.Logger.Println("step 14 ---->")
	transferAmount := int32(50)
	C01 = N0.GetChannelWith(N1, tokenAddress).PrintDataBeforeTransfer()
	C12 = N1.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	C24 = N2.GetChannelWith(N4, tokenAddress).PrintDataBeforeTransfer()
	err = N0.Transfer(tokenAddress, transferAmount, N4.Address, false)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	time.Sleep(6 * time.Second)
	C01new := N0.GetChannelWith(N1, tokenAddress).PrintDataAfterTransfer()
	if !C01new.CheckPartnerBalance(C01.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C01new.Name)
	}
	C12new := N1.GetChannelWith(N2, tokenAddress).PrintDataAfterTransfer()
	if !C12new.CheckPartnerBalance(C12.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C12new.Name)
	}
	C24new := N2.GetChannelWith(N4, tokenAddress).PrintDataAfterTransfer()
	if !C24new.CheckPartnerBalance(C24.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C24new.Name)
	}

	// step 15 : N2 sends all of its tokens to N1 (one transfer)
	models.Logger.Println("step 15 ---->")
	C12 = N1.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	err = N2.Transfer(tokenAddress, C12.PartnerBalance, N1.Address, false)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	time.Sleep(6 * time.Second)
	C12new = N1.GetChannelWith(N2, tokenAddress).PrintDataAfterTransfer()
	if !C12new.CheckSelfBalance(C12.Balance + C12.PartnerBalance) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C12new.Name)
	}

	// step 16 : N2 tries to send another > 25 tokens payment to N1 (fail no route with enough capacity)
	models.Logger.Println("step 16 ---->")
	transferAmount = 30
	err = N2.Transfer(tokenAddress, transferAmount, N1.Address, false)
	if err == nil {
		return cm.caseFail(env.CaseName)
	}
	err = nil

	// step 17 : N2 sends 10 payments of 1 token to N1 by using the N2 <-> N4 <-> N0 <-> N1 route
	models.Logger.Println("step 17 ---->")
	transferAmount = 10
	C24 = N2.GetChannelWith(N4, tokenAddress).PrintDataBeforeTransfer()
	C04 = N0.GetChannelWith(N4, tokenAddress).PrintDataBeforeTransfer()
	C01 = N0.GetChannelWith(N1, tokenAddress).PrintDataBeforeTransfer()
	for i := int32(0); i < transferAmount; i++ {
		err = N2.Transfer(tokenAddress, 1, N1.Address, false)
		if err != nil {
			return cm.caseFail(env.CaseName)
		}
	}
	time.Sleep(60 * time.Second)
	C24new = N2.GetChannelWith(N4, tokenAddress).PrintDataBeforeTransfer()
	if !C24new.CheckPartnerBalance(C24.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C24new.Name)
	}
	C04new := N0.GetChannelWith(N4, tokenAddress).PrintDataBeforeTransfer()
	if !C04new.CheckPartnerBalance(C04.PartnerBalance - transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C04new.Name)
	}
	C01new = N0.GetChannelWith(N1, tokenAddress).PrintDataBeforeTransfer()
	if !C01new.CheckPartnerBalance(C01.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C01new.Name)
	}

	// step 18 : N1 shuts down
	models.Logger.Println("step 18 ---->")
	N1.Shutdown()
	if N1.IsRunning() {
		return cm.caseFail(env.CaseName)
	}

	// step 19 : N0 sends 10 tokens to N2 (using the N0<-> N4 <-> N2 route)
	models.Logger.Println("step 19 ---->")
	transferAmount = 10
	C04 = N0.GetChannelWith(N4, tokenAddress).PrintDataBeforeTransfer()
	C24 = N2.GetChannelWith(N4, tokenAddress).PrintDataBeforeTransfer()
	err = N0.Transfer(tokenAddress, transferAmount, N2.Address, false)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	time.Sleep(6 * time.Second)
	C04new = N0.GetChannelWith(N4, tokenAddress).PrintDataAfterTransfer()
	if !C04new.CheckPartnerBalance(C04.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C04new.Name)
	}
	C24new = N2.GetChannelWith(N4, tokenAddress).PrintDataAfterTransfer()
	if !C24new.CheckSelfBalance(C24.Balance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C24new.Name)
	}

	// step 20 : N0 tries to open a channel with N1 (fail - it already has a channel with N1)
	models.Logger.Println("step 20 ---->")
	err = N0.OpenChannel(N1.Address, tokenAddress, 0, settleTimeout)
	if err == nil {
		return cm.caseFail(env.CaseName)
	}

	// step 21 : N0 tries to make a payment to N1 (Node is offline - fails)
	models.Logger.Println("step 21 ---->")
	err = N0.Transfer(tokenAddress, transferAmount, N1.Address, false)
	if err == nil {
		return cm.caseFail(env.CaseName)
	}
	err = nil

	// step 22 : N1 is back online
	models.Logger.Println("step 22 ---->")
	N1.ReStartWithoutConditionquit(env)

	// step 23 : N3 sends all 100 tokens to N2 on payments of 1 token/each.
	models.Logger.Println("step 23 ---->")
	transferAmount = 100
	C23 = N2.GetChannelWith(N3, tokenAddress).PrintDataBeforeTransfer()
	for i := int32(0); i < transferAmount; i++ {
		err = N3.Transfer(tokenAddress, 1, N2.Address, false)
		if err != nil {
			return cm.caseFail(env.CaseName)
		}
	}
	time.Sleep(300 * time.Second)
	C23new := N2.GetChannelWith(N3, tokenAddress).PrintDataAfterTransfer()
	if !C23new.CheckSelfBalance(C23.Balance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C23new.Name)
	}

	// step 24 : N0 deposits 160 tokens on his channel with N1
	models.Logger.Println("step 24 ---->")
	depositAmount = 160
	err = N0.Deposit(C01.ChannelIdentifier, depositAmount)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}

	// step 25 : Assert that route N0->N1->N2->N3 has enough capacity to send 200 tokens from N0 to N3
	models.Logger.Println("step 25 ---->")
	transferAmount = 200
	C01 = N0.GetChannelWith(N1, tokenAddress).PrintDataBeforeTransfer()
	if C01.Balance < transferAmount {
		return cm.caseFailWithWrongChannelData(env.CaseName, C01.Name)
	}
	C12 = N1.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	if C12.Balance < transferAmount {
		return cm.caseFailWithWrongChannelData(env.CaseName, C12.Name)
	}
	C23 = N2.GetChannelWith(N3, tokenAddress).PrintDataBeforeTransfer()
	if C23.Balance < transferAmount {
		return cm.caseFailWithWrongChannelData(env.CaseName, C23.Name)
	}

	// step 26 : Perform 200 payments From N0 to N3
	// step 27 : assert
	models.Logger.Println("step 27 ---->")
	err = N0.Transfer(tokenAddress, transferAmount, N3.Address, false)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	time.Sleep(20 * time.Second)
	C01new = N0.GetChannelWith(N1, tokenAddress).PrintDataBeforeTransfer()
	if !C01new.CheckPartnerBalance(C01.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C01new.Name)
	}
	C12new = N1.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	if !C12new.CheckPartnerBalance(C12.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C12new.Name)
	}
	C23new = N2.GetChannelWith(N3, tokenAddress).PrintDataBeforeTransfer()
	if !C23new.CheckPartnerBalance(C12.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C23new.Name)
	}

	// step 28 : N4 closes his channel with N2
	models.Logger.Println("step 28 ---->")
	N4.Close(C24.ChannelIdentifier)
	time.Sleep(20 * time.Second)
	C24 = N2.GetChannelWith(N4, tokenAddress)
	if C24.State != int(netshare.Closed) {
		return cm.caseFail(env.CaseName)
	}

	// step 29 : N2 tries to make a deposit in the channel that is being closed (fail 409)
	models.Logger.Println("step 29 ---->")
	err = N2.Deposit(C24.ChannelIdentifier, depositAmount)
	if err == nil {
		return cm.caseFail(env.CaseName)
	}

	// step 30 : N2 sends 10 tokens to N1
	transferAmount = 10
	models.Logger.Println("step 30 ---->")
	C12 = N1.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	err = N2.Transfer(tokenAddress, transferAmount, N1.Address, false)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	time.Sleep(20 * time.Second)
	C12new = N1.GetChannelWith(N2, tokenAddress).PrintDataAfterTransfer()
	if !C12new.CheckPartnerBalance(C12.PartnerBalance - transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C12new.Name)
	}

	// step 31 : N1 sends 10 tokens to N0
	transferAmount = 10
	models.Logger.Println("step 31 ---->")
	C01 = N0.GetChannelWith(N1, tokenAddress).PrintDataBeforeTransfer()
	err = N1.Transfer(tokenAddress, transferAmount, N0.Address, false)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	time.Sleep(20 * time.Second)
	C01new = N0.GetChannelWith(N1, tokenAddress).PrintDataAfterTransfer()
	if !C01new.CheckPartnerBalance(C01.PartnerBalance - transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C01new.Name)
	}

	// step 32 : N4 sends 10 tokens to N2 (N0 -> N1-> N2)
	transferAmount = 10
	models.Logger.Println("step 32 ---->")
	C04 = N0.GetChannelWith(N4, tokenAddress).PrintDataBeforeTransfer()
	C01 = N0.GetChannelWith(N1, tokenAddress).PrintDataBeforeTransfer()
	C12 = N1.GetChannelWith(N2, tokenAddress).PrintDataBeforeTransfer()
	err = N4.Transfer(tokenAddress, transferAmount, N2.Address, false)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	time.Sleep(10 * time.Second)
	C04new = N0.GetChannelWith(N4, tokenAddress).PrintDataAfterTransfer()
	if !C04new.CheckSelfBalance(C04.Balance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C04new.Name)
	}
	C01new = N0.GetChannelWith(N1, tokenAddress).PrintDataAfterTransfer()
	if !C01new.CheckPartnerBalance(C01.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C01new.Name)
	}
	C12new = N1.GetChannelWith(N2, tokenAddress).PrintDataAfterTransfer()
	if !C12new.CheckPartnerBalance(C12.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C12new.Name)
	}

	// step 33 : settle all channel
	models.Logger.Println("step 33 ---->")
	N0.CooperateSettle(C01.ChannelIdentifier)
	N0.CooperateSettle(C04.ChannelIdentifier)
	N1.CooperateSettle(C12.ChannelIdentifier)
	N2.CooperateSettle(C23.ChannelIdentifier)
	time.Sleep(time.Duration(settleTimeout) * time.Second) // wait to settle C24
	N2.Settle(C24.ChannelIdentifier)
	time.Sleep(100 * time.Second) // wait to settle C24

	C01 = N0.GetChannelWith(N1, tokenAddress)
	if C01 != nil {
		C01.Println("")
		return cm.caseFail(env.CaseName)
	}
	C04 = N0.GetChannelWith(N4, tokenAddress)
	if C04 != nil {
		C04.Println("")
		return cm.caseFail(env.CaseName)
	}
	C12 = N1.GetChannelWith(N2, tokenAddress)
	if C12 != nil {
		C12.Println("")
		return cm.caseFail(env.CaseName)
	}
	C23 = N2.GetChannelWith(N3, tokenAddress)
	if C23 != nil {
		C23.Println("")
		return cm.caseFail(env.CaseName)
	}
	C24 = N2.GetChannelWith(N4, tokenAddress)
	if C24 != nil {
		C24.Println("")
		return cm.caseFail(env.CaseName)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
