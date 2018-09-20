package cases

import (
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
)

// LocalCase : only for local test
func (cm *CaseManager) LocalCase() (err error) {
	env, err := models.NewTestEnv("./cases/LocalCase.ENV")
	if err != nil {
		return
	}
	defer func() {
		if env.Debug == false {
			env.KillAllRaidenNodes()
		}
	}()
	// 源数据
	transAmount := int32(20)
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N0, N1, N2 := env.Nodes[0], env.Nodes[1], env.Nodes[2]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点,让节点0在收到
	N0.StartWithConditionQuit(env, &params.ConditionQuit{
		QuitEvent: "EventSendAnnouncedDisposedResponseBefore",
	})
	N1.Start(env)
	N2.Start(env)
	//channel := N1.GetChannelWith(N0, tokenAddress)
	go N0.SendTrans(tokenAddress, transAmount, N2.Address, false)
	//time.Sleep(10 * time.Second)
	//N1.Close(channel.ChannelAddress)
	time.Sleep(1000 * time.Second)
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
