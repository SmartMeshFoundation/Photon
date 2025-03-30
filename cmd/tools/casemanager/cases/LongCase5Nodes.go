package cases

import (
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/utils"
)

// LongCase5Nodes :
func (cm *CaseManager) LongCase5Nodes() (err error) {
	if !cm.RunSlow {
		return ErrorSkip
	}
	env, err := models.NewTestEnv("./cases/LongCase5Nodes.ENV", cm.UseMatrix, cm.EthEndPoint)
	if err != nil {
		return
	}
	defer func() {
		if env.Debug == false {
			env.KillAllPhotonNodes()
		} else {
			//time.Sleep(time.Minute * 1000)
		}
	}()
	// 源数据
	settleTimeout := int64(500)
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N0, N1, N2, N3, N4 := env.Nodes[0], env.Nodes[1], env.Nodes[2], env.Nodes[3], env.Nodes[4]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// step 1 : Start 5 Atmosphere nodes
	models.Logger.Println("step 1 ---->")
	cm.startNodes(env, N0, N1, N2, N3, N4)
	if cm.UseMatrix {
		time.Sleep(time.Second * 10)
	}
	// step 2 : Create the following channels: N0 - N1, N1 - N2, N2 - N3 with 100 deposit
	models.Logger.Println("step 2 ---->")
	depositAmount := int64(100)
	err = N0.OpenChannel(N1.Address, tokenAddress, depositAmount, settleTimeout)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
	C01 := N0.GetChannelWith(N1, tokenAddress)
	if C01 == nil {
		return cm.caseFail(env.CaseName)
	}
	err = N1.OpenChannel(N2.Address, tokenAddress, depositAmount, settleTimeout)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	C12 := N1.GetChannelWith(N2, tokenAddress)
	if C12 == nil {
		return cm.caseFail(env.CaseName)
	}
	err = N2.OpenChannel(N3.Address, tokenAddress, depositAmount, settleTimeout)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	C23 := N2.GetChannelWith(N3, tokenAddress)
	if C23 == nil {
		return cm.caseFail(env.CaseName)
	}

	// step 3 : N0 N1 N2 N3 make a deposit of 100 to their channels
	models.Logger.Println("step 3 ---->")
	err = N1.Deposit(N0.Address, tokenAddress, depositAmount)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
	err = N2.Deposit(N1.Address, tokenAddress, depositAmount)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	err = N3.Deposit(N2.Address, tokenAddress, depositAmount)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}

	// step 4 : N4 tries to make a deposit to a channel that does not exist (fail channel with N3 does not exist)
	models.Logger.Println("step 4 ---->")
	err = N4.Deposit(N0.Address, utils.NewRandomAddress().String(), depositAmount)
	if err == nil {
		return cm.caseFail(env.CaseName)
	}
	err = nil

	// step 5 : N1 makes a deposit of 50 tokens on both of his channels (skip)
	// step 6 : N1 tries to open a channel with N0, but it already has a channel (skip)
	// step 7 : N2 opens a channel with N4 (0 tokens)
	// step 8 : N2 makes a deposit of 100 tokens
	models.Logger.Println("step 8 ---->")
	err = N2.OpenChannel(N4.Address, tokenAddress, depositAmount, settleTimeout)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
	C24 := N2.GetChannelWith(N4, tokenAddress)
	if C24 == nil {
		return cm.caseFail(env.CaseName)
	}

	// step 9 : N2 tries to deposit 30 tokens (skip)
	// step 10 : N2 deposits 50 tokens to the channel (N2 - N4)
	models.Logger.Println("step 10 ---->")
	depositAmount = 50
	err = N2.Deposit(N4.Address, tokenAddress, depositAmount)
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
	if err != nil {
		return cm.caseFail(env.CaseName)
	}

	// step 12 : N0 opens a channel with N4 (initial deposit of 10)
	models.Logger.Println("step 12 ---->")
	depositAmount = 10
	err = N0.OpenChannel(N4.Address, tokenAddress, depositAmount, settleTimeout)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	C04 := N0.GetChannelWith(N4, tokenAddress)
	if C04 == nil {
		return cm.caseFail(env.CaseName)
	}

	// step 13 : N4 deposits 25 tokens to N0<->N4 channel
	models.Logger.Println("step 13 ---->")
	depositAmount = 25
	err = N4.Deposit(N0.Address, tokenAddress, depositAmount)
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
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
	time.Sleep(1 * time.Second)
	if cm.UseMatrix {
		time.Sleep(time.Second * 10)
	}
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
	time.Sleep(1 * time.Second)
	if cm.UseMatrix {
		time.Sleep(time.Second * 7)
	}
	C12new = N1.GetChannelWith(N2, tokenAddress).PrintDataAfterTransfer()
	if !C12new.CheckSelfBalance(C12.Balance + C12.PartnerBalance) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C12new.Name)
	}

	// step 16 : N2 tries to send another > 25 tokens payment to N1 (fail no route with enough capacity)
	models.Logger.Println("step 16 ---->")
	transferAmount = 30
	if cm.UseMatrix {
		time.Sleep(time.Second * 7)
	}
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
		if cm.UseMatrix {
			time.Sleep(time.Second * 3)
		}
	}
	time.Sleep(6 * time.Second)
	C24new = N2.GetChannelWith(N4, tokenAddress).PrintDataAfterTransfer()
	if !C24new.CheckPartnerBalance(C24.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C24new.Name)
	}
	C04new := N0.GetChannelWith(N4, tokenAddress).PrintDataAfterTransfer()
	if !C04new.CheckPartnerBalance(C04.PartnerBalance - transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C04new.Name)
	}
	C01new = N0.GetChannelWith(N1, tokenAddress).PrintDataAfterTransfer()
	if !C01new.CheckPartnerBalance(C01.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C01new.Name)
	}

	// step 18 : N1 shuts down
	models.Logger.Println("step 18 ---->")
	N1.Shutdown(env)
	if cm.UseMatrix {
		time.Sleep(cm.MDNSLifeTime + time.Second*7)
	} else {
		// 等待mdns检测下线
		time.Sleep(cm.MDNSLifeTime)
	}
	if N1.IsRunning() {
		return cm.caseFail(env.CaseName)
	}

	// step 19 : N0 sends 10 tokens to N2 (using the N0<-> N4 <-> N2 route)
	models.Logger.Println("step 19 ---->")
	transferAmount = 10
	C04 = N0.GetChannelWith(N4, tokenAddress).PrintDataBeforeTransfer()
	C24 = N2.GetChannelWith(N4, tokenAddress).PrintDataBeforeTransfer()
	err = N0.Transfer(tokenAddress, transferAmount, N2.Address, false)
	if cm.UseMatrix {
		time.Sleep(time.Second * 10)
	}
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	time.Sleep(1 * time.Second)
	C04new = N0.GetChannelWith(N4, tokenAddress).PrintDataAfterTransfer()
	if !C04new.CheckPartnerBalance(C04.PartnerBalance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C04new.Name)
	}
	C24new = N2.GetChannelWith(N4, tokenAddress).PrintDataAfterTransfer()
	if !C24new.CheckSelfBalance(C24.Balance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C24new.Name)
	}

	// step 20 : N0 tries to open a channel with N1 (fail - it already has a channel with N1)(skip)

	// step 21 : N0 tries to make a payment to N1 (Node is offline - fails)
	models.Logger.Println("step 21 ---->")
	err = N0.Transfer(tokenAddress, transferAmount, N1.Address, false)
	if cm.UseMatrix {
		time.Sleep(time.Second * 10)
	}
	if err == nil {
		return cm.caseFail(env.CaseName)
	}
	err = nil

	// step 22 : N1 is back online
	models.Logger.Println("step 22 ---->")
	N1.ReStartWithoutConditionquit(env)
	if cm.UseMatrix {
		time.Sleep(time.Second * 10)
	}

	// step 23 : N3 sends all 100 tokens to N2 on payments of 1 token/each.
	models.Logger.Println("step 23 ---->")
	transferAmount = 100
	C23 = N2.GetChannelWith(N3, tokenAddress).PrintDataBeforeTransfer()
	for i := int32(0); i < transferAmount/2; i++ {
		err = N3.Transfer(tokenAddress, 2, N2.Address, false)
		if err != nil {
			return cm.caseFailWithWrongChannelData(env.CaseName, fmt.Sprintf("mass transfer i=%d,err=%s", i, err.Error()))
		}
		if cm.UseMatrix {
			time.Sleep(time.Second * 3)
		}
	}
	//等30秒,确认100笔交易成功
	time.Sleep(30 * time.Second)
	C23new := N2.GetChannelWith(N3, tokenAddress).PrintDataAfterTransfer()
	if !C23new.CheckSelfBalance(C23.Balance + transferAmount) {
		return cm.caseFailWithWrongChannelData(env.CaseName, C23new.Name)
	}

	// step 24 : N0 deposits 160 tokens on his channel with N1
	models.Logger.Println("step 24 ---->")
	depositAmount = 160
	if cm.UseMatrix {
		time.Sleep(time.Second * 20)
	}
	err = N0.Deposit(N1.Address, tokenAddress, depositAmount)
	if err != nil {
		return cm.caseFail(env.CaseName)
	}

	// step 25 : Assert that route N0->N1->N2->N3 has enough capacity to send 200 tokens from N0 to N3
	models.Logger.Println("step 25 ---->")
	transferAmount = 190
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
	if cm.UseMatrix {
		time.Sleep(time.Second * 10)
	}
	if err != nil {
		return cm.caseFail(env.CaseName)
	}
	time.Sleep(2 * time.Second)
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
	err = cm.tryInSeconds(25, func() error {
		return N4.Close(C24.ChannelIdentifier)
	})
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, err.Error())
	}
	time.Sleep(2 * time.Second) // 等待节点接收close事件
	C24 = N2.GetChannelWith(N4, tokenAddress)
	if C24.State != int(channeltype.StateClosed) {
		return cm.caseFail(env.CaseName)
	}

	// step 29 : N2 tries to make a deposit in the channel that is being closed (fail 409)
	models.Logger.Println("step 29 ---->")
	err = N2.Deposit(N4.Address, tokenAddress, depositAmount)
	if cm.UseMatrix {
		time.Sleep(time.Second * 10)
	}
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
	time.Sleep(2 * time.Second)
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
	time.Sleep(2 * time.Second)
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
	time.Sleep(2 * time.Second)
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
	err = N0.CooperateSettle(C01.ChannelIdentifier)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, err.Error())
	}
	err = N0.CooperateSettle(C04.ChannelIdentifier)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, err.Error())
	}
	err = N1.CooperateSettle(C12.ChannelIdentifier)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, err.Error())
	}
	err = N2.CooperateSettle(C23.ChannelIdentifier)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, err.Error())
	}
	err = cm.trySettleInSeconds(int(settleTimeout+260), N2, C24.ChannelIdentifier)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, err.Error())
	}
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
