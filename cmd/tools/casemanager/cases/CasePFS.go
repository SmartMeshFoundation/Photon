package cases

import (
	"math/big"

	"strings"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/pfsproxy"
	"github.com/ethereum/go-ethereum/common"
)

// CasePFS :
func (cm *CaseManager) CasePFS() (err error) {
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
	// 2. 查询n0-n2,金额=10000的path,expect 0-1-2
	transferAmount := big.NewInt(10000)
	pfsProxy := env.GetPfsProxy(env.GetPrivateKeyByNode(n0))
	err = testFindPath(pfsProxy, tokenAddress, n0, n2, transferAmount, []*models.PhotonNode{n1, n2})
	if err != nil {
		models.Logger.Println(err)
		return cm.caseFail(env.CaseName)
	}
	// 3. 查询n0-n2,金额=10000的path,expect 2-1-0
	transferAmount = big.NewInt(10000)
	pfsProxy = env.GetPfsProxy(env.GetPrivateKeyByNode(n2))
	err = testFindPath(pfsProxy, tokenAddress, n2, n0, transferAmount, []*models.PhotonNode{n1, n0})
	if err != nil {
		models.Logger.Println(err)
		return cm.caseFail(env.CaseName)
	}
	// 4. 查询n0-n1,金额=10000的path,expect 0-1
	transferAmount = big.NewInt(10000)
	pfsProxy = env.GetPfsProxy(env.GetPrivateKeyByNode(n0))
	err = testFindPath(pfsProxy, tokenAddress, n0, n1, transferAmount, []*models.PhotonNode{n1})
	if err != nil {
		models.Logger.Println(err)
		return cm.caseFail(env.CaseName)
	}
	// 3. 查询
	models.Logger.Println(env.CaseName + " END ====> SUCCESS")
	return nil
}

func testFindPath(pfsProxy pfsproxy.PfsProxy, token common.Address, from, to *models.PhotonNode, transferAmount *big.Int, expect ...[]*models.PhotonNode) (err error) {
	resp, err := pfsProxy.FindPath(from.GetAddress(), to.GetAddress(), token, transferAmount, true)
	if err != nil {
		return
	}
	if len(expect) == 0 {
		return
	}
	if len(resp) != len(expect) {
		err = fmt.Errorf("routes num wrong ,expect %d but got %d", len(expect), len(resp))
		return
	}
	for i := 0; i < len(resp); i++ {
		routeGot := resp[i].Result
		var routeExpect []string
		for _, n := range expect[i] {
			routeExpect = append(routeExpect, n.Address)
		}
		if len(routeGot) != len(routeExpect) {
			err = fmt.Errorf("path id=%d length wrong ,expect %d but got %d", resp[i].PathID, len(routeExpect), len(routeGot))
			return
		}
		for j := 0; j < len(routeGot); j++ {
			if !isEqualAddress(routeGot[j], routeExpect[j]) {
				err = fmt.Errorf("path id=%d wrong, expect %s but got %s", resp[i].PathID, routeExpect, routeGot)
				return
			}
		}
	}
	return
}

func isEqualAddress(a1, a2 string) bool {
	if strings.Compare(strings.ToLower(a1), strings.ToLower(a2)) != 0 {
		return false
	}
	return true
}
