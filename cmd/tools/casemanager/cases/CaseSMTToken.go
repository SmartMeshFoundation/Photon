package cases

import (
	"context"

	"math/big"

	"fmt"

	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/ethereum/go-ethereum/common"
)

// CaseSMTToken :
// 测试SMTToken
func (cm *CaseManager) CaseSMTToken() (err error) {
	if !cm.RunSlow {
		return ErrorSkip
	}
	env, err := models.NewTestEnv("./cases/CaseSMTToken.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	settleTimeout := int64(100)
	tokenAddress := env.Tokens[0].TokenAddress.String()
	N1, N2 := env.Nodes[0], env.Nodes[1]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点1，2
	// start node 2, 3
	cm.startNodes(env, N1, N2)
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
	// 1. 打印N1,N2余额
	showBalance(env, "begin", N1, N2)
	// 2. N1创建通道
	mainChainBalanceN1 := N1.Runtime.MainChainBalance
	depositAmount1 := int64(100)
	err = N1.OpenChannel(N2.Address, tokenAddress, depositAmount1, settleTimeout)
	if err != nil {
		models.Logger.Println(err)
		return cm.caseFail(env.CaseName)
	}
	// 3. 查询通道数据,并校验对等
	c12 := N1.GetChannelWith(N2, tokenAddress).Println("after N1 openAndDeposit 100")
	if !c12.CheckEqualByPartnerNode(env) {
		return cm.caseFailWithWrongChannelData(env.CaseName, c12.Name)
	}
	if !c12.CheckSelfBalance(int32(depositAmount1)) {
		return cm.caseFailWithWrongChannelData(env.CaseName, c12.Name)
	}
	// 4. 打印N1,N2余额
	showBalance(env, "after N1 openAndDeposit 100", N1, N2)
	// 4.5 校验N1主链余额
	fmt.Println(new(big.Int).Sub(mainChainBalanceN1, N1.Runtime.MainChainBalance).Int64())
	fmt.Println(depositAmount1)
	if new(big.Int).Sub(mainChainBalanceN1, N1.Runtime.MainChainBalance).Int64() <= depositAmount1 {
		models.Logger.Println("N1 mainChainBalance err")
		return cm.caseFail(env.CaseName)
	}
	// 5. N2 deposit
	mainChainBalanceN2 := N2.Runtime.MainChainBalance
	depositAmount2 := int64(50)
	err = N2.Deposit(N1.Address, tokenAddress, depositAmount2)
	if err != nil {
		models.Logger.Println(err)
		return cm.caseFail(env.CaseName)
	}
	// 6. 查询通道数据,并校验对等
	c21 := N2.GetChannelWith(N1, tokenAddress).Println("after N2 Deposit 50")
	if !c21.CheckEqualByPartnerNode(env) {
		return cm.caseFailWithWrongChannelData(env.CaseName, c21.Name)
	}
	if !c21.CheckSelfBalance(int32(depositAmount2)) {
		return cm.caseFailWithWrongChannelData(env.CaseName, c21.Name)
	}
	// 5. 打印N1,N2余额
	showBalance(env, "after N2 Deposit 50", N1, N2)
	// 5.5 校验N2主链余额
	if new(big.Int).Sub(mainChainBalanceN2, N2.Runtime.MainChainBalance).Int64() <= depositAmount2 {
		models.Logger.Println("N2 mainChainBalance err")
		return cm.caseFail(env.CaseName)
	}
	// 6. 结算
	err = N1.CooperateSettle(c12.ChannelIdentifier)
	if err != nil {
		models.Logger.Println(err)
		return cm.caseFail(env.CaseName)
	}
	// 7. 查询通道
	c12 = N1.GetChannelWith(N1, tokenAddress)
	if c12 != nil {
		models.Logger.Println("channel should not exist")
		return cm.caseFailWithWrongChannelData(env.CaseName, c12.Name)
	}
	// 8. 打印N1,N2余额
	showBalance(env, "after settle", N1, N2)

	// 9. N1 重新open
	depositAmount1 = int64(200)
	err = N1.OpenChannel(N2.Address, tokenAddress, depositAmount1, settleTimeout)
	if err != nil {
		models.Logger.Println(err)
		return cm.caseFail(env.CaseName)
	}
	// 10. 查询通道数据,并校验对等
	c12 = N1.GetChannelWith(N2, tokenAddress).Println("after N1 openAndDeposit 200")
	if !c12.CheckEqualByPartnerNode(env) {
		return cm.caseFailWithWrongChannelData(env.CaseName, c12.Name)
	}
	if !c12.CheckSelfBalance(int32(depositAmount1)) {
		return cm.caseFailWithWrongChannelData(env.CaseName, c12.Name)
	}
	// 11. N1 withdraw
	withdrawAmount := int32(100)
	N1.Withdraw(c12.ChannelIdentifier, withdrawAmount)

	// 12. 查询通道数据,并校验对等
	i := 0
	for i = 0; i < cm.HighMediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		c12 = N1.GetChannelWith(N2, tokenAddress).Println("after N1 openAndDeposit 200")
		if !c12.CheckEqualByPartnerNode(env) {
			continue
		}
		if !c12.CheckSelfBalance(int32(depositAmount1) - withdrawAmount) {
			continue
		}
		break
	}
	if i == cm.HighMediumWaitSeconds {
		return cm.caseFailWithWrongChannelData(env.CaseName, c12.Name)
	}
	// 12.5 N2 deposit
	err = N2.Deposit(N1.Address, tokenAddress, depositAmount2)
	if err != nil {
		models.Logger.Println(err)
		return cm.caseFail(env.CaseName)
	}
	// 13. N1 Close
	err = N1.Close(c12.ChannelIdentifier)
	if err != nil {
		models.Logger.Println(err)
		return cm.caseFail(env.CaseName)
	}
	// 14. N1 Settle
	err = cm.trySettleInSeconds(int(c12.SettleTimeout)+257+10, N1, c12.ChannelIdentifier)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, err.Error())
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}

func showBalance(env *models.TestEnv, prefix string, node ...*models.PhotonNode) {
	if len(node) == 0 {
		return
	}
	if prefix != "" {
		models.Logger.Println("-->", prefix)
	}
	models.Logger.Println("--> Balance of MainChain : ")
	for _, n := range node {
		b, err := env.Conn.BalanceAt(context.Background(), common.HexToAddress(n.Address), nil)
		if err != nil {
			return
		}
		n.Runtime.MainChainBalance = b
		models.Logger.Printf("\t%s = %d\n", n.Name, b.Uint64())
	}
	models.Logger.Println("--> Balance of MainChain Token :")
	for _, n := range node {
		b, err := env.Tokens[0].Token.BalanceOf(nil, common.HexToAddress(n.Address))
		if err != nil {
			return
		}
		models.Logger.Printf("\t%s = %d\n", n.Name, b.Uint64())
	}
	return
}
