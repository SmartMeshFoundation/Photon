package contracttest

import (
	"testing"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

// TestSecretRegistryRight : 正确调用测试
func TestSecretRegistryRight(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	a1 := env.getRandomAccountExcept(t)

	// version test
	_, err := env.SecretRegistry.ContractVersion(nil)
	assertSuccess(t, &count, err)

	// get blockNo
	blockNoWhenRegister := getLatestBlockNumber()
	// register right
	secret := utils.ShaSecret([]byte(utils.RandomString(9)))
	tx, err := env.SecretRegistry.RegisterSecret(a1.Auth, secret)
	assertTxSuccess(t, &count, tx, err)

	// get block height by secret hash
	secretHash := utils.ShaSecret(secret[:])
	blockNo, err := env.SecretRegistry.GetSecretRevealBlockHeight(nil, secretHash)
	assertSuccess(t, &count, err)
	assertEqual(t, &count, true, blockNo.Cmp(blockNoWhenRegister.Number) >= 0)

	// register repeat
	tx, err = env.SecretRegistry.RegisterSecret(a1.Auth, secret)
	assertTxFail(t, &count, tx, err)

	t.Log(endMsg("SecretRegistry 正确调用测试", count))
}

// TestSecretRegistryException : 异常调用测试
func TestSecretRegistryException(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	a1 := env.getRandomAccountExcept(t)

	// register with empty hash
	tx, err := env.SecretRegistry.RegisterSecret(a1.Auth, utils.EmptyHash)
	assertTxFail(t, &count, tx, err)

	t.Log(endMsg("SecretRegistry 异常调用测试", count))

}

// TestSecretRegistryEdge : 边界测试
func TestSecretRegistryEdge(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	t.Log(endMsg("SecretRegistry 边界测试", count))
}

// TestSecretRegistryAttack : 恶意调用测试
func TestSecretRegistryAttack(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("SecretRegistry 恶意调用测试", count))
}
