package cases

import (
	"time"

	"fmt"
	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

// CasePMS03 :
func (cm *CaseManager) CasePMS03() (err error) {
	if !cm.RunSlow {
		return ErrorSkip
	}
	env, err := models.NewTestEnv("./cases/CasePMS03.ENV", cm.UseMatrix, cm.EthEndPoint, "CasePMS03")
	if err != nil {
		return
	}
	defer func() {
		if env.Debug == false {
			env.KillAllPhotonNodes()
		}
	}()

	// 源数据
	// original data
	tokenAddress := env.Tokens[0].TokenAddress.String()

	N1, N2 := env.Nodes[1], env.Nodes[2]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动pms
	env.StartPMS()
	// 启动节点1、2
	cm.startNodes(env, N1, N2.PMS())

	transAmount := int32(10)
	// 获取channel信息
	c12 := N1.GetChannelWith(N2, tokenAddress).Println("before send trans")
	n1TokenOld, err := N1.TokenBalance(tokenAddress)
	n2TokenOld, err := N2.TokenBalance(tokenAddress)
	if err != nil {
		return
	}
	n1TokenOld += int(c12.Balance)
	n2TokenOld += int(c12.PartnerBalance)
	// N1 send trans to N3
	N1.SendTrans(tokenAddress, transAmount, N2.Address, false)
	time.Sleep(time.Second)
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
	// step2: b无网重启，
	N1.Shutdown(env)
	N2.Shutdown(env)
	time.Sleep(time.Second)
	cm.startNodes(env, N1.RestartName().NoNetwork(), N2.RestartName().NoNetwork())
	N1.SendTrans(tokenAddress, transAmount, N2.Address, true)
	time.Sleep(time.Second)
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
	// tep3:a关闭通道
	c12 = N1.GetChannelWith(N2, tokenAddress).Println("check CD-N1-N2")
	err = N1.Close(c12.ChannelIdentifier)
	if err != nil {
		return
	}
	models.Logger.Printf("n1 close CD-N1-N2...")
	// N2 settle channel
	settleTime := c12.SettleTimeout + 257
	err = cm.trySettleInSeconds(int(settleTime), N1, c12.ChannelIdentifier)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, c12.Name)
	}
	models.Logger.Printf("n1 settle finished...")

	N1.Shutdown(env)
	N2.Shutdown(env)
	time.Sleep(time.Second)
	cm.startNodes(env, N1.RestartName().HaveNetwork(), N2.RestartName().HaveNetwork())

	err = cm.tryInSeconds(cm.LowWaitSeconds, func() error {
		n1Token, err := N1.TokenBalance(tokenAddress)
		n2Token, err := N2.TokenBalance(tokenAddress)
		if err != nil {
			return cm.caseFail(env.CaseName)
		}
		if n1Token != n1TokenOld-10 && n2Token != n2TokenOld+10 {
			return cm.caseFailWithWrongChannelData(env.CaseName, fmt.Sprintf("check balance onchain, n1=%d,n1expect=%d,n2=%d,n2expect=%d ", n1Token, n1TokenOld-10, n2Token, n2TokenOld+10))
		}
		return nil
	})
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, err.Error())
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
