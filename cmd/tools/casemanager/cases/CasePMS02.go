package cases

import (
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/params"
)

// CasePMS02 :
func (cm *CaseManager) CasePMS02() (err error) {
	if !cm.RunSlow {
		return ErrorSkip
	}
	env, err := models.NewTestEnv("./cases/CasePMS02.ENV", cm.UseMatrix, cm.EthEndPoint, "CasePMS02")
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

	N1, N2, N3 := env.Nodes[1], env.Nodes[2], env.Nodes[3]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动pms
	env.StartPMS()
	// 启动节点2、3
	cm.startNodes(env, N1, N2.PMS().SetConditionQuit(&params.ConditionQuit{
		QuitEvent: "ReceiveRevealSecretAck",
	}), N3)

	transAmount := int32(10)
	// 获取channel信息
	c12 := N1.GetChannelWith(N2, tokenAddress).Println("before send trans")
	c32 := N3.GetChannelWith(N2, tokenAddress).Println("before send trans")
	// N1 send trans to N3
	secret, _, err := N1.GenerateSecret()
	if err != nil {
		return
	}
	N1.SendTrans(tokenAddress, transAmount, N3.Address, false)
	// 崩溃判断
	for i := 0; i < cm.MediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		if !N2.IsRunning() {
			break
		}
	}
	if N2.IsRunning() {
		msg := "Node " + N2.Name + " should be exited,but it still running, FAILED !!!"
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}

	err = N1.Close(c12.ChannelIdentifier)
	if err != nil {
		return
	}
	err = N3.RegisterSecret(secret)
	if err != nil {
		return
	}
	time.Sleep(time.Second)
	N1.GetChannelWith(N2, tokenAddress).Println("after send trans")
	N3.GetChannelWith(N2, tokenAddress).Println("after send trans")
	err = cm.tryInSeconds(cm.LowWaitSeconds, func() error {
		c12new := N1.GetChannelWith(N2, tokenAddress).Println("check CD-N1-N2")
		c32new := N3.GetChannelWith(N2, tokenAddress).Println("check CD-N3-N2")
		if c12new.LockedAmount != int32(0) && c32new.PartnerLockedAmount != int32(0) {
			return cm.caseFailWithWrongChannelData(env.CaseName, fmt.Sprintf("check lock failed"))
		}
		if c32new.Balance != c32.Balance+transAmount && c12new.PartnerBalance != c12.PartnerBalance+transAmount {
			return cm.caseFailWithWrongChannelData(env.CaseName, "check balance failed")
		}
		return nil
	})
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, err.Error())
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
