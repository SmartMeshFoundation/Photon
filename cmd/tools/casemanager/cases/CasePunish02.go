package cases

import (
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

/*
CasePunish02 : test for punish
#N0-N1-N2 交易,N1 AnnouceDisposed 给N0,但是N0在EventSendAnnouncedDisposedResponseBefore崩溃,
#然后N0重新启动,但是不收发任何消息,同时强制关闭通道,然后等结算窗口过期以后settle通道. N1不应该unlock,N0不应该
#有机会惩罚N0

# 测试N0不能自身诱导N1犯错获利
*/
func (cm *CaseManager) CasePunish02() (err error) {
	env, err := models.NewTestEnv("./cases/CasePunish02.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	//N1,N2不再和N0发生通信
	N1.UpdateMeshNetworkNodes(cm.nodesExcept(env.Nodes, N0)...)
	N2.UpdateMeshNetworkNodes(cm.nodesExcept(env.Nodes, N1)...)
	//N0启动以后无法和N1,N2通信
	N0.ReStartWithoutConditionquitAndNetwork(env)
	//N0注册密码
	err = N0.RegisterSecret(secret)
	if err != nil {
		return fmt.Errorf("register secret err %s", err)
	}
	//N1关闭通道,
	err = N1.Close(c10.ChannelIdentifier)
	if err != nil {
		return fmt.Errorf("close channel err %s", err)
	}
	/*
		这时候N1不应该解锁,因为他已经发出去了annoucedisposed,虽然他没有收到annoucedisposed response
	*/

	expectN1 := int32(n1value) + n1balance
	expectN0 := n0balance + int32(n0value)
	err = cm.trySettleInSeconds(cm.HighMediumWaitSeconds, N0, c10.ChannelIdentifier)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, c10.Name)
	}

	var n0NewValue, n1NewValue int

	n0NewValue, err = N0.TokenBalance(tokenAddress)
	n1NewValue, err = N1.TokenBalance(tokenAddress)
	if n1NewValue != int(expectN1) || n0NewValue != int(expectN0) {
		return cm.caseFailWithWrongChannelData(env.CaseName, fmt.Sprintf("check balance error n0=%d,n0expect=%d,n1=%d,n1expect=%d ", n0NewValue, expectN0, n1NewValue, expectN1))
	}

	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
