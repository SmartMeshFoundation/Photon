package cases

import (
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/Photon/utils"

	"github.com/SmartMeshFoundation/Photon/log"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

/*CaseTransferWithSameSecret :
# 连续两次交易使用相同的密码,第二笔在第一笔交易完成以后开始
# 交易成败与否不不关键,不能发生崩溃,中间节点不能丢钱.
*/
func (cm *CaseManager) CaseTransferWithSameSecret() (err error) {
	env, err := models.NewTestEnv("./cases/CaseTransferWithSameSecret.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	N0, N1, N2 := env.Nodes[0], env.Nodes[1], env.Nodes[2]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	cm.startNodes(env, N0, N1, N2)
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
	// 获取channel信息
	// get channel info
	N0.GetChannelWith(N1, tokenAddress).Println("before transfer")
	secret := utils.NewRandomHash()
	secretHash := utils.ShaSecret(secret[:])
	go N0.SendTransWithSecret(env.Tokens[0].TokenAddress.String(), 1, N2.Address, secret.String())

	time.Sleep(time.Second * 5)
	log.Trace("allow reveal secret")
	N0.AllowSecret(secretHash.String(), tokenAddress)
	time.Sleep(time.Second * 5)
	N0.GetChannelWith(N1, tokenAddress).Println("after transfer")
	c01 := N0.GetChannelWith(N1, tokenAddress).Println("before next transfer")
	c12 := N1.GetChannelWith(N2, tokenAddress).Println("before next transfer")
	go N0.SendTransWithSecret(tokenAddress, 1, N2.Address, secret.String())
	time.Sleep(time.Second * 3)
	N0.AllowSecret(secretHash.String(), tokenAddress)
	time.Sleep(time.Second * 3)

	/*
		检测结果:
		不能发生崩溃
		交易要么都失败,要么都成功.
	*/
	err = cm.tryInSeconds(cm.HighMediumWaitSeconds, func() error {
		if !N0.IsRunning() {
			return fmt.Errorf("n0 should not crash")
		}
		if !N1.IsRunning() {
			return fmt.Errorf("n1 should not crash")
		}
		if !N2.IsRunning() {
			return fmt.Errorf("n2 should not crash")
		}
		c01new := N0.GetChannelWith(N1, tokenAddress).Println("after next transfer")
		c12new := N1.GetChannelWith(N2, tokenAddress).Println("after next transfer")
		//无论哪种情况,有锁的时候都不检测.
		if !c01new.CheckLockBoth(0) || !c12new.CheckLockBoth(0) {
			return fmt.Errorf("have lock")
		}
		if c01new.PartnerBalance+c12new.Balance != c01.PartnerBalance+c12.Balance {
			return fmt.Errorf("n1 balance changed") //无论成败,这两者之和必须相等,否则就是出错了.
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("second transfer dectect err %s", err)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
