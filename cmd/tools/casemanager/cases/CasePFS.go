package cases

import (
	"math/big"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/ethereum/go-ethereum/common"
)

// CasePFS :
func (cm *CaseManager) CasePFS() (err error) {
	if !cm.RunThisCaseOnly {
		return
	}
	env, err := models.NewTestEnv("./cases/CasePFS.ENV", cm.UseMatrix, cm.EthEndPoint)
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
	// 2. 查询n0-n2,金额=10000的path
	transferAmount := big.NewInt(10000)
	pfsProxy := env.GetPfsProxy(env.GetPrivateKeyByNode(n0))
	route, err := pfsProxy.FindPath(n0.GetAddress(), n1.GetAddress(), tokenAddress, transferAmount, true)
	if err != nil {
		models.Logger.Println(err)
		return cm.caseFail(env.CaseName)
	}
	// 3. 发送交易
	n0.SendTrans(tokenAddressStr, int32(transferAmount.Int64()), n2.Address, false)
	models.Logger.Printf("pfsProxy.FindPath response:\n%s\n", models.MarshalIndent(route))
	// 3. 查询
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}
