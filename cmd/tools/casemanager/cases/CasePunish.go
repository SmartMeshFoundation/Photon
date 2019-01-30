package cases

import (
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/Photon/restful/v1"

	"github.com/SmartMeshFoundation/Photon/utils"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

/*
CasePunish : test for punish
# N0-N1-N2,N1-N2因为钱不够,N1 AnnouceDisposed,N0在EventSendAnnouncedDisposedResponseBefore崩溃,
# 这时候N1强制关闭通道,调用测试接口,forceUnlock,
# N0重启以后,应该punish N1,拿走N1通道中的所有钱
*/
func (cm *CaseManager) CasePunish() (err error) {
	env, err := models.NewTestEnv("./cases/CasePunish.ENV", cm.UseMatrix, cm.EthEndPoint)
	if err != nil {
		return
	}
	defer func() {
		if env.Debug == false {
			env.KillAllPhotonNodes()
		}
	}()
	// 源数据
	transAmount := int32(20)
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N0, N1, N2 := env.Nodes[0], env.Nodes[1], env.Nodes[2]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点,让节点0发送SendAnnounce
	N0.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "EventSendAnnouncedDisposedResponseBefore",
	})
	N1.Start(env)
	N2.Start(env)
	secret, _, err := N0.GenerateSecret()
	if err != nil {
		return
	}
	//获取交易前金额状况
	c10 := N1.GetChannelWith(N0, tokenAddress)
	n0balance := c10.PartnerBalance
	n1balance := c10.Balance
	n0value, err := N0.TokenBalance(tokenAddress)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, "query balance error")
	}
	n1value, err := N1.TokenBalance(tokenAddress)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, "query balance n1 error")
	}
	models.Logger.Printf("n0balance=%d,n1balance=%d,n0value=%d,n1value=%d", n0balance, n1balance, n0value, n1value)
	go N0.SendTransWithSecret(tokenAddress, transAmount, N2.Address, secret)
	time.Sleep(time.Second * 3)
	// N0 crash
	if N0.IsRunning() {
		return fmt.Errorf("n0 should shutdown")
	}
	err = N1.ForceUnlock(c10.ChannelIdentifier, secret)
	if err != nil {
		return fmt.Errorf("force unlock err %s", err)
	}

	N0.ReStartWithoutConditionquit(env)
	//N0 should punish N1
	//check balance
	expectN1 := n1value
	expectN0 := n0balance + n1balance + int32(n0value)
	var i = 0
	for i = 0; i < cm.MediumWaitSeconds; i++ {
		var c v1.ChannelDataDetail
		time.Sleep(time.Second)
		c, err = N0.SpecifiedChannel(c10.ChannelIdentifier)
		if err != nil {
			continue
		}
		if (c.OurBalanceProof.ContractTransferAmount != nil && c.OurBalanceProof.ContractTransferAmount.Uint64() != 0) ||
			c.OurBalanceProof.ContractLocksRoot != utils.EmptyHash ||
			c.OurBalanceProof.ContractNonce != 0xffffffffffffffff {
			models.Logger.Printf("c=%s", utils.StringInterface(c, 5))
			continue
		}
		break
	}
	if i == cm.MediumWaitSeconds {
		return cm.caseFailWithWrongChannelData(env.CaseName, "check balance proof error")
	}
	for i = 0; i < cm.HighMediumWaitSeconds; i++ {
		var n0NewValue, n1NewValue int
		time.Sleep(time.Second)
		err = N0.Settle(c10.ChannelIdentifier)
		if err != nil {
			continue
		}
		time.Sleep(time.Second)
		n0NewValue, err = N0.TokenBalance(tokenAddress)
		n1NewValue, err = N1.TokenBalance(tokenAddress)
		if n1NewValue != expectN1 || n0NewValue != int(expectN0) {
			return cm.caseFailWithWrongChannelData(env.CaseName, fmt.Sprintf("check balance error n0=%d,n0expect=%d,n1=%d,n1expect=%d ", n0NewValue, expectN0, n1NewValue, expectN1))
		}
		break
	}
	if i == cm.HighMediumWaitSeconds {
		return cm.caseFailWithWrongChannelData(env.CaseName, "settle   error")
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
