package cases

import (
	"math/big"

	"errors"
	"time"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/ethereum/go-ethereum/common"
)

// CaseTransferWithRouteInfo :
func (cm *CaseManager) CaseTransferWithRouteInfo() (err error) {
	if cm.IsAutoRun {
		return ErrorSkip
	}
	env, err := models.NewTestEnv("./cases/CaseTransferWithRouteInfo.ENV", cm.UseMatrix, cm.EthEndPoint)
	if err != nil {
		return
	}
	defer func() {
		if env.Debug == false {
			env.KillAllPhotonNodes()
		}
	}()
	// 源数据
	tokenAddressStr := env.Tokens[0].TokenAddress.String()
	tokenAddress := common.HexToAddress(tokenAddressStr)
	n0, n1, n2 := env.Nodes[0], env.Nodes[1], env.Nodes[2]
	// 0. 启动pfs
	env.StartPFS()
	// 1. 启动节点
	n0.StartWithFeeAndPFS(env)
	n1.StartWithFeeAndPFS(env)
	n2.StartWithFeeAndPFS(env)
	if cm.UseMatrix {
		time.Sleep(time.Second * 5)
	}
	// 2. 查询n0-n2,金额=10000的path
	transferAmount := big.NewInt(10000)
	pfsProxy := env.GetPfsProxy(env.GetPrivateKeyByNode(n0))
	route, err := pfsProxy.FindPath(n0.GetAddress(), n2.GetAddress(), tokenAddress, transferAmount, true)
	if err != nil {
		models.Logger.Println(err)
		return cm.caseFail(env.CaseName)
	}
	models.Logger.Printf("pfsProxy.FindPath response:\n%s\n", models.MarshalIndent(route))
	// 3. 发送交易
	c01 := n0.GetChannelWith(n1, tokenAddressStr).Println("before transfer")
	c12 := n1.GetChannelWith(n2, tokenAddressStr).Println("before transfer")
	n0.SendTransWithRouteInfo(n2, tokenAddressStr, int32(transferAmount.Int64()), route)
	// 4. 余额校验
	err = cm.tryInSeconds(15, func() error {
		c01new := n0.GetChannelWith(n1, tokenAddressStr).Println("after transfer")
		c12new := n1.GetChannelWith(n2, tokenAddressStr).Println("after transfer")
		if !c01new.CheckEqualByPartnerNode(env) {
			return errors.New("not equal with partner")
		}
		if !c01new.CheckSelfBalance(c01.Balance - int32(transferAmount.Int64()) - 1) {
			return errors.New("unexpected balance")
		}
		if !c12new.CheckEqualByPartnerNode(env) {
			return errors.New("not equal with partner")
		}
		if !c12new.CheckSelfBalance(c12.Balance - int32(transferAmount.Int64())) {
			return errors.New("unexpected balance")
		}
		return nil
	})
	if err != nil {
		models.Logger.Println(err)
		return cm.caseFail(env.CaseName)
	}
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
