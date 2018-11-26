package cases

import (
	"time"

	"sync"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

// LocalCase : only for local test
func (cm *CaseManager) LocalCase() (err error) {
	env, err := models.NewTestEnv("./cases/LocalCase.ENV")
	if err != nil {
		return
	}
	defer func() {
		if env.Debug == false {
			env.KillAllPhotonNodes()
		}
	}()
	// 源数据
	transAmount := int32(1)
	tokenAddress := env.Tokens[0].TokenAddress.String()
	number := 1000
	N0, N1, _ := env.Nodes[0], env.Nodes[1], env.Nodes[2]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点,让节点0在收到
	N0.Start(env)
	N1.Start(env)
	//N2.Start(env)
	//channel := N1.GetChannelWith(N0, tokenAddress)
	begin := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(number)
	wg2 := sync.WaitGroup{}
	wg2.Add(number)
	for i := 0; i < number; i++ {
		go func(index int) {
			wg.Done()
			wg.Wait()
			bt := time.Now()
			N0.SendTransSync(tokenAddress, transAmount, N1.Address, true)
			fmt.Println("transfer ", index, "use  ", time.Since(bt).Seconds())
			wg2.Done()
		}(i)
	}
	wg2.Wait()
	total := time.Since(begin).Seconds()
	fmt.Println("total=", total)
	fmt.Println("tps=", 1000/total)
	//time.Sleep(10 * time.Second)
	//N1.Close(channel.ChannelIdentifier)
	time.Sleep(1000 * time.Second)
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return
}
