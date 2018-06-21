package cases

import (
	"fmt"

	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/utils"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

// CrashCaseRecv02 场景二：ReceiveSecretRequestStateChange
// 收到Secretrequest后崩
// 节点1向节点6发送20个token,节点6向节点1发送secretrequest请求，节点1收到崩,
// 节点1、节点2、节点3各锁定20个token；重启节点1后，节点锁定token解锁，转账成功。
func (cm *CaseManager) CrashCaseRecv02() (err error) {
	env, err := models.NewTestEnv("./cases/CrashCaseRecv02.ENV")
	if err != nil {
		return
	}
	defer env.KillAllRaidenNodes()
	// 源数据
	var transAmount int32
	var msg string
	transAmount = 20
	tokenAddress := env.Tokens[0].Address
	N1, N2, N3, N6 := env.Nodes[0], env.Nodes[1], env.Nodes[2], env.Nodes[3]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点2,36
	N2.Start(env)
	N3.Start(env)
	N6.Start(env)
	// 启动节点1, ReceiveSecretRequestStateChange
	N1.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "ReceiveSecretRequestStateChange",
	})

	// 查询节点2，记录cd21数据
	cd21 := utils.GetChannelBetween(N2, N1, tokenAddress).PrintDataBeforeTransfer()
	// 查询节点2，记录cd24数据
	cd23 := utils.GetChannelBetween(N2, N3, tokenAddress).PrintDataBeforeTransfer()
	// 查询节点3，记录cd36数据
	cd36 := utils.GetChannelBetween(N3, N6, tokenAddress).PrintDataBeforeTransfer()

	// 节点1向节点6转账20token
	N1.SendTrans(tokenAddress, transAmount, N6.Address, false)
	time.Sleep(time.Second * 3)
	//  崩溃判断
	if N1.IsRunning() {
		panic("Node N1 should be exited,but it still running")
	}

	// 查询cd21，锁定对方20
	cd21middle := utils.GetChannelBetween(N2, N1, tokenAddress).PrintDataAfterCrash()
	if cd21middle.PartnerLockedAmount != transAmount {
		msg = fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", transAmount, cd21middle.PartnerLockedAmount)
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 查询cd23，锁定20
	cd23middle := utils.GetChannelBetween(N2, N3, tokenAddress).PrintDataAfterCrash()
	if cd23middle.LockedAmount != transAmount {
		msg = fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", transAmount, cd23middle.LockedAmount)
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}
	// 查询cd36，锁定20
	cd36middle := utils.GetChannelBetween(N3, N6, tokenAddress).PrintDataAfterCrash()
	if cd36middle.LockedAmount != transAmount {
		msg = fmt.Sprintf("Expect locked amount = %d,but got %d ,FAILED!!!", transAmount, cd36middle.LockedAmount)
		models.Logger.Println(msg)
		return fmt.Errorf(msg)
	}

	// 重启节点1，交易自动继续
	N1.ReStartWithoutConditionquit(env)

	// 查询cd21并校验
	cd21new := utils.GetChannelBetween(N2, N1, tokenAddress).PrintDataAfterRestart()
	if cd21new.Balance-cd21.Balance != transAmount {
		models.Logger.Println(env.CaseName + " END ====> FAILED")
		return fmt.Errorf("Case [%s] FAILED", env.CaseName)
	}
	// 查询cd23并校验
	cd23new := utils.GetChannelBetween(N2, N3, tokenAddress).PrintDataAfterRestart()
	if cd23new.PartnerBalance-cd23.PartnerBalance != transAmount {
		models.Logger.Println(env.CaseName + " END ====> FAILED")
		return fmt.Errorf("Case [%s] FAILED", env.CaseName)
	}
	// 查询cd36并校验
	cd36new := utils.GetChannelBetween(N3, N6, tokenAddress).PrintDataAfterRestart()
	if cd36new.PartnerBalance-cd36.PartnerBalance != transAmount {
		models.Logger.Println(env.CaseName + " END ====> FAILED")
		return fmt.Errorf("Case [%s] FAILED", env.CaseName)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
