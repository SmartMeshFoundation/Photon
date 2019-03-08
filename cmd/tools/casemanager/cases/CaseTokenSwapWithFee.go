package cases

import (
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/ethereum/go-ethereum/common"
)

// CaseTokenSwapWithFee :
func (cm *CaseManager) CaseTokenSwapWithFee() (err error) {
	if !cm.RunSlow {
		return
	}
	env, err := models.NewTestEnv("./cases/CaseTokenSwapWithFee.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	n0, n1, n2, n3, n4, n5 := env.Nodes[0], env.Nodes[1], env.Nodes[2], env.Nodes[3], env.Nodes[4], env.Nodes[5]
	models.Logger.Println(env.CaseName + " BEGIN ====>")

	env.StartPFS()
	cm.startNodesWithFee(env, n0, n1, n2, n3, n4, n5)
	time.Sleep(time.Second)
	models.Logger.Println("start  token swap.")
	secret2, secrethash2, err := n0.GenerateSecret()
	if err != nil {
		return
	}
	token2 := env.Tokens[1].TokenAddress.String()
	c01t0 := n0.GetChannelWith(n1, tokenAddress).Println("before token swap")
	c01t1 := n0.GetChannelWith(n1, token2).Println("before token swap")
	takerPath := n0.FindPath(n4, common.HexToAddress(tokenAddress), 1)
	err = n0.TokenSwap(n4.Address, secrethash2, tokenAddress, token2, "taker", "", 1, 3, takerPath)
	if err != nil {
		return fmt.Errorf(" token swap taker err=%s", err)
	}
	makerPath := n4.FindPath(n0, common.HexToAddress(token2), 3)
	err = n4.TokenSwap(n0.Address, secrethash2, token2, tokenAddress, "maker", secret2, 3, 1, makerPath)
	if err != nil {
		models.Logger.Println("token swap fail")
		return fmt.Errorf(" token sdwap maker err=%s", err)
	}
	time.Sleep(time.Second * 3) //必须多等一会儿,否则查到的信息不准确.
	c01t0new := n0.GetChannelWith(n1, tokenAddress).Println("after token swap")
	c01t1new := n0.GetChannelWith(n1, token2).Println("after token swap")
	if !c01t0new.CheckSelfBalance(c01t0.Balance - 1) {
		return fmt.Errorf(" token swap check sending banlance 0 err")
	}
	if !c01t1new.CheckSelfBalance(c01t1.Balance + 3) {
		return fmt.Errorf(" token swap check receiving banlance 0 err")
	}
	c12 := n1.GetChannelWith(n2, tokenAddress).Println("after token swap")
	if !c12.CheckLockBoth(0) {
		return fmt.Errorf("c12 check error")
	}
	c12t1 := n1.GetChannelWith(n2, token2).Println("after token swap")
	if !c12t1.CheckLockBoth(0) {
		return fmt.Errorf("c12t1 check error")
	}
	c23 := n2.GetChannelWith(n3, tokenAddress).Println("after token swap")
	c25 := n2.GetChannelWith(n5, token2).Println("")
	c34 := n3.GetChannelWith(n4, tokenAddress).Println("")
	c54 := n5.GetChannelWith(n4, token2).Println("")
	if !c23.CheckLockBoth(0) {
		return fmt.Errorf("c23 check error")
	}
	if !c25.CheckLockBoth(0) {
		return fmt.Errorf("c25 check error")
	}
	if !c34.CheckLockBoth(0) {
		return fmt.Errorf("c34 check error")
	}
	if !c54.CheckLockBoth(0) {
		return fmt.Errorf("c54 check error")
	}
	return nil
}
