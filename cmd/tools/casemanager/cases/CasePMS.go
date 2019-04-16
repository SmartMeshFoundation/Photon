package cases

import (
	"fmt"
	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/log"
	"time"
)

// CasePMS :
func (cm *CaseManager) CasePMS() (err error) {
	if !cm.RunSlow {
		return ErrorSkip
	}
	env, err := models.NewTestEnv("./cases/CasePMS.ENV", cm.UseMatrix, cm.EthEndPoint, "CasePMS")
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
	// 启动节点2
	cm.startNodes(env, N2)
	// 启动委托方节点N1
	cm.startNodesWithPMS(env, N1)
	// 获取channel信息
	c12 := N1.GetChannelWith(N2, tokenAddress).Println("before send trans")
	n1value, err := N1.TokenBalance(tokenAddress)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, "query balance n1 error")
	}
	n2value, err := N2.TokenBalance(tokenAddress)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, "query balance n2 error")
	}

	n1value += int(c12.Balance)
	n2value += int(c12.PartnerBalance)
	log.Trace(fmt.Sprintf("before transfer ,n1value=%d,n2value=%d", n1value, n2value))

	// N2 send trans
	N2.SendTrans(tokenAddress, 10, N1.Address, true)
	time.Sleep(time.Second) //wait for Delegate
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
	// n1 shut down
	N1.Shutdown(env)
	models.Logger.Println("n1 shutdown")

	// N2 close channel
	err = N2.Close(c12.ChannelIdentifier)
	if err != nil {
		return
	}
	// N2 settle channel
	settleTime := c12.SettleTimeout + 257
	err = cm.trySettleInSeconds(int(settleTime), N2, c12.ChannelIdentifier)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, c12.Name)
	}
	models.Logger.Printf("n2 settle finished...")

	N1.RestartName().StartWithPMS(env)
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
	var n1NewValue, n2NewValue, i int
	for i = 0; i < 5; i++ {
		time.Sleep(time.Second * 1)
		n1NewValue, err = N1.TokenBalance(tokenAddress)
		n2NewValue, err = N2.TokenBalance(tokenAddress)
		if n1NewValue != n1value+10 {
			continue
		}
		if n2NewValue != n2value-10 {
			continue
		}
		if N1.GetChannelWith(N2, tokenAddress) != nil {
			continue
		}
		models.Logger.Println(env.CaseName + " END ====> SUCCESS")
		return
	}
	return cm.caseFailWithWrongChannelData(env.CaseName, fmt.Sprintf("check TokenBalance error n1=%d,n1expect=%d,n2=%d,n2expect=%d ", n1NewValue, n1value+10, n2NewValue, n2value-10))
}
