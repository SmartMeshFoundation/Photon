package contracttest

import "testing"

// TestXXXXRight : 正确调用测试
func TestXXXXRight(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("XXXX 正确调用测试", count))
}

// TestXXXXException : 异常调用测试
func TestXXXXException(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("XXXX 异常调用测试", count))

}

// TestXXXXEdge : 边界测试
func TestXXXXEdge(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("XXXX 边界测试", count))
}

// TestXXXXAttack : 恶意调用测试
func TestXXXXAttack(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("XXXX 恶意调用测试", count))
}