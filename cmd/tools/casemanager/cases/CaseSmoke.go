package cases

import (
	"fmt"
	"time"

	models2 "github.com/SmartMeshFoundation/Photon/models"

	"github.com/SmartMeshFoundation/Photon/utils"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
)

// CaseSmoke :
func (cm *CaseManager) CaseSmoke() (err error) {
	env, err := models.NewTestEnv("./cases/CaseSmoke.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	n0, n1, n2, n3 := env.Nodes[0], env.Nodes[1], env.Nodes[2], env.Nodes[3]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点2，3
	// start node 2, 3
	cm.startNodes(env, n0, n1, n2, n3)

	//第一类 升级api
	models.Logger.Println("start api class1..")
	err = n3.PrepareUpdate()
	if err != nil {
		return
	}
	err = n3.Transfer(tokenAddress, 1, n2.Address, true)
	if err == nil {
		return fmt.Errorf("can not send any transfer when update")
	}

	//第二类 基本交易相关api
	models.Logger.Println("start api class2..")
	// get channel info
	c01 := n0.GetChannelWith(n1, tokenAddress).Println("before transfer")
	err = n0.Transfer(env.Tokens[0].TokenAddress.String(), 1, n1.Address, false)
	if err != nil {
		return err
	}
	c01new := n0.GetChannelWith(n1, tokenAddress).Println("after transfer")
	if !c01new.CheckLockBoth(0) {
		return fmt.Errorf("transfer check lock err ")
	}
	if !c01new.CheckSelfBalance(c01.Balance - 1) {
		return fmt.Errorf("check balance error")
	}
	if !c01new.CheckPartnerBalance(c01.PartnerBalance + 1) {
		return fmt.Errorf("check partner balance error")
	}

	trs, err := n0.GetSentTransfers()
	if err != nil {
		return fmt.Errorf("GetSentTransfers err %s", err)
	}
	if len(trs) != 1 || trs[0].Amount.Uint64() != 1 {
		return fmt.Errorf("GetSentTransfers err trs=%s", utils.StringInterface(trs, 3))
	}
	rrs, err := n1.GetReceivedTransfers()
	if err != nil {
		return fmt.Errorf("GetReceivedTransfers err %s", err)
	}
	if len(rrs) != 1 || rrs[0].Amount.Uint64() != 1 {
		return fmt.Errorf("GetReceivedTransfers err rrs=%s", utils.StringInterface(rrs, 3))
	}

	//2.2 带密码的交易
	secret, secrethash, err := n0.GenerateSecret()
	if err != nil {
		return err
	}
	go n0.SendTransWithSecret(tokenAddress, 1, n1.Address, secret)
	time.Sleep(time.Second)
	urt, err := n1.GetUnfinishedReceivedTransfer(tokenAddress, secrethash)
	if err != nil {
		return err
	}
	if urt == nil || urt.LockSecretHash != secrethash {
		return fmt.Errorf("GetUnfinishedReceivedTransfer secrethash=%s,urt=%s", secrethash, utils.StringInterface(urt, 3))
	}
	c01new = n0.GetChannelWith(n1, tokenAddress).Println("after transfer")
	if !c01new.CheckLockSelf(1) {
		return fmt.Errorf("SendTransWithSecret check lock err ")
	}
	n0.AllowSecret(secrethash, tokenAddress)
	var i = 0
	for i = 0; i < 10; i++ {
		time.Sleep(time.Second) //等待足够的时间,重发一次消息
		c01new = n0.GetChannelWith(n1, tokenAddress).Println("after transfer")
		if c01new.CheckLockBoth(0) {
			break
		}
	}
	if i == 10 {
		return fmt.Errorf("after allow secret, should no lock")
	}
	//2.3 带密码的交易,取消
	secret, secrethash, err = n0.GenerateSecret()
	if err != nil {
		return err
	}
	go n0.SendTransWithSecret(tokenAddress, 1, n1.Address, secret)
	time.Sleep(time.Second)
	err = n0.CancelTransfer(tokenAddress, secrethash)
	if err != nil {
		return err
	}
	c01new = n0.GetChannelWith(n1, tokenAddress).Println("after transfer")
	if !c01new.CheckLockSelf(1) {
		return fmt.Errorf("CancelTransfer check lock err ")
	}
	st, err := n0.GetTransferStatus(tokenAddress, secrethash)
	if err != nil {
		return err
	}
	if st == nil || st.Status != models2.TransferStatusCanceled {
		return fmt.Errorf("cancel transfer status err expect=%d,got=%s", models2.TransferStatusCanceled, utils.StringInterface(st, 3))
	}
	//这个锁只能等待过期才会自动消失.

	//第三类 基本查询
	models.Logger.Println("start api class3..")
	ts, err := n0.Tokens()
	if err != nil || len(ts) != 2 {
		return fmt.Errorf("Tokens err ts=%s,err=%s", utils.StringInterface(ts, 2), err)
	}
	ps, err := n0.TokenPartners(tokenAddress)
	if err != nil || len(ps) != 1 {
		return fmt.Errorf("TokenPartners err,ps=%s,err=%s", utils.StringInterface(ps, 2), err)
	}

	//第四类 tokenswap
	// 直接通道的tokenswap
	models.Logger.Println("start api class4..")
	models.Logger.Println("start direct token swap.")
	secret2, secrethash2, err := n0.GenerateSecret()
	if err != nil {
		return
	}
	token2 := env.Tokens[1].TokenAddress.String()
	c01t0 := n0.GetChannelWith(n1, tokenAddress)
	c01t1 := n0.GetChannelWith(n1, token2)
	err = n0.TokenSwap(n1.Address, secrethash2, tokenAddress, token2, "taker", "", 1, 3)
	if err != nil {
		return fmt.Errorf("direct token swap taker err=%s", err)
	}
	err = n1.TokenSwap(n0.Address, secrethash2, token2, tokenAddress, "maker", secret2, 3, 1)
	if err != nil {
		models.Logger.Println("direct token swap fail")
		return fmt.Errorf("direct token sdwap maker err=%s", err)
	}
	//time.Sleep(time.Second)
	c01t0new := n0.GetChannelWith(n1, tokenAddress)
	c01t1new := n0.GetChannelWith(n1, token2)
	if !c01t0new.CheckSelfBalance(c01t0.Balance - 1) {
		return fmt.Errorf("direct token swap check sending banlance 0 err")
	}
	if !c01t1new.CheckSelfBalance(c01t1.Balance + 3) {
		return fmt.Errorf("direct token swap check receiving banlance 0 err")
	}

	models.Logger.Println("start  token swap.")
	secret2, secrethash2, err = n0.GenerateSecret()
	if err != nil {
		return
	}
	token2 = env.Tokens[1].TokenAddress.String()
	c01t0 = n0.GetChannelWith(n1, tokenAddress)
	c01t1 = n0.GetChannelWith(n1, token2)
	err = n0.TokenSwap(n2.Address, secrethash2, tokenAddress, token2, "taker", "", 1, 3)
	if err != nil {
		return fmt.Errorf(" token swap taker err=%s", err)
	}
	err = n2.TokenSwap(n0.Address, secrethash2, token2, tokenAddress, "maker", secret2, 3, 1)
	if err != nil {
		models.Logger.Println("token swap fail")
		return fmt.Errorf(" token sdwap maker err=%s", err)
	}
	time.Sleep(time.Second * 3) //必须多等一会儿,否则查到的信息不准确.
	c01t0new = n0.GetChannelWith(n1, tokenAddress)
	c01t1new = n0.GetChannelWith(n1, token2)
	if !c01t0new.CheckSelfBalance(c01t0.Balance - 1) {
		return fmt.Errorf(" token swap check sending banlance 0 err")
	}
	if !c01t1new.CheckSelfBalance(c01t1.Balance + 3) {
		return fmt.Errorf(" token swap check receiving banlance 0 err")
	}

	return nil
}
