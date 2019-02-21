package debug

import "testing"

func TestSetup(t *testing.T) {
	t.Logf("mac=%s", mac())
}
