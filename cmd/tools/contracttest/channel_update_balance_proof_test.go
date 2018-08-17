package contracttest

import "testing"

// TestUpdateBalanceProofRight : 正确调用测试
func TestUpdateBalanceProofRight(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("UpdateBalanceProof 正确调用测试", count))
}

// TestUpdateBalanceProofException : 异常调用测试
func TestUpdateBalanceProofException(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("UpdateBalanceProof 异常调用测试", count))

}

// TestUpdateBalanceProofEdge : 边界测试
func TestUpdateBalanceProofEdge(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("UpdateBalanceProof 边界测试", count))
}

// TestUpdateBalanceProofAttack : 恶意调用测试
func TestUpdateBalanceProofAttack(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("UpdateBalanceProof 恶意调用测试", count))
}
