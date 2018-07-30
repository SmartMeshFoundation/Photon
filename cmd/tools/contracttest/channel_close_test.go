package contracttest

import "testing"

// TestChannelCloseRight : 正确调用测试
func TestChannelCloseRight(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("ChannelClose 正确调用测试", count))
}

// TestChannelCloseException : 异常调用测试
func TestChannelCloseException(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("ChannelClose 异常调用测试", count))

}

// TestChannelCloseEdge : 边界测试
func TestChannelCloseEdge(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("ChannelClose 边界测试", count))
}

// TestChannelCloseAttack : 恶意调用测试
func TestChannelCloseAttack(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	t.Log(endMsg("ChannelClose 恶意调用测试", count))
}
