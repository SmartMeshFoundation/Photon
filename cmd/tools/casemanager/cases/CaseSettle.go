package cases

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	pmodels "github.com/SmartMeshFoundation/Photon/models"
)

// CaseSettle :
func (cm *CaseManager) CaseSettle() (err error) {
	if !cm.RunSlow {
		return ErrorSkip
	}
	env, err := models.NewTestEnv("./cases/CaseSettle.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	N0, N1 := env.Nodes[0], env.Nodes[1]
	models.Logger.Println(env.CaseName + " BEGIN ====>")
	// 启动节点2，3
	// start node 2, 3
	cm.startNodes(env, N0, N1)
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
	// 获取channel信息
	// get channel info
	c01 := N0.GetChannelWith(N1, tokenAddress).Println("BeforeClose")
	N0.SendTrans(env.Tokens[0].TokenAddress.String(), 1, N1.Address, false)
	//time.Sleep(3 * time.Second)
	// Close
	if cm.UseMatrix {
		time.Sleep(time.Second * 7)
	}
	err = N0.Close(c01.ChannelIdentifier)
	if err != nil {
		return
	}
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
	N0.GetChannelWith(N1, tokenAddress).Println("AfterClose")
	err = cm.trySettleInSeconds(int(c01.SettleTimeout)+257+10, N0, c01.ChannelIdentifier)
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, err.Error())
	}
	var i int
	for i = 0; i < cm.MediumWaitSeconds; i++ {
		time.Sleep(time.Second)
		// 验证
		// verify
		c01new := N0.GetChannelWith(N1, tokenAddress).Println("AfterSettle")
		if c01new == nil {
			break
		}
	}
	if i == cm.MediumWaitSeconds {
		return cm.caseFailWithWrongChannelData(env.CaseName, c01.Name)
	}
	err = N0.Transfer(tokenAddress, 1, N1.Address, false)
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
	if err == nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, "Transfer must failed after cooperate settle")
	}
	txs, err := N0.ContractCallTXQuery(c01.ChannelIdentifier, pmodels.TXInfoTypeSettle)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, "query tx failed")
	}
	if len(txs) != 1 {
		return cm.caseFailWithWrongChannelData(env.CaseName, fmt.Sprintf("txs length =%d", len(txs)))
	}
	tx := txs[0]
	settleParams := &pmodels.ChannelSettleTXParams{}
	err = json.Unmarshal([]byte(tx.TXParams), settleParams)
	if err != nil {
		return cm.caseFailWithWrongChannelData(env.CaseName, fmt.Sprintf("tx params error %s", err))
	}
	if settleParams.P1Balance.Int64() != 49 || settleParams.P2Balance.Int64() != 51 {
		return cm.caseFailWithWrongChannelData(env.CaseName, tx.TXParams)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
