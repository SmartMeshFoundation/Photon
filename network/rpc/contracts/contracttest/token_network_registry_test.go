package contracttest

import (
	"math/big"
	"testing"

	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts/test/tokens/tokenerc223"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

// TestTokenNetworkRegistryRight : 正确调用测试
func TestTokenNetworkRegistryRight(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	a1 := env.getRandomAccountExcept(t)

	// deploy new token
	tokenAddress, tx, _, err := tokenerc223.DeployHumanERC223Token(a1.Auth, env.Client, big.NewInt(500000000000000000), "test erc223", 0)
	assertTxSuccess(t, nil, tx, err)
	// create token network
	tx, err = env.TokenNetworkRegistry.CreateERC20TokenNetwork(a1.Auth, tokenAddress)
	assertTxSuccess(t, &count, tx, err)

	// get the new token network address
	tokenNetworkAddress, err := env.TokenNetworkRegistry.TokenToTokenNetworks(nil, tokenAddress)
	assertSuccess(t, &count, err)

	// check token network contract exists, and work right
	newTokenNetwork, err := contracts.NewTokenNetwork(tokenNetworkAddress, env.Client)
	assertSuccess(t, &count, err)

	version, err := newTokenNetwork.ContractVersion(nil)
	assertSuccess(t, &count, err)
	assertEqual(t, &count, true, version != "")
	t.Log(endMsg("TokenNetworkRegistry 正确调用测试", count))
}

// TestTokenNetworkRegistryException : 异常调用测试
func TestTokenNetworkRegistryException(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	a1 := env.getRandomAccountExcept(t)

	// create token network with wrong token address
	tx, err := env.TokenNetworkRegistry.CreateERC20TokenNetwork(a1.Auth, EmptyAccountAddress)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetworkRegistry.CreateERC20TokenNetwork(a1.Auth, common.HexToAddress(utils.RandomString(9)))
	assertTxFail(t, &count, tx, err)

	t.Log(endMsg("TokenNetworkRegistry 异常调用测试", count))

}

// TestTokenNetworkRegistryEdge : 边界测试
func TestTokenNetworkRegistryEdge(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("TokenNetworkRegistry 边界测试", count))
}

// TestTokenNetworkRegistryAttack : 恶意调用测试
func TestTokenNetworkRegistryAttack(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("TokenNetworkRegistry 恶意调用测试", count))
}
