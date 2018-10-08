package mainimpl

import (
	"testing"

	"github.com/SmartMeshFoundation/SmartRaiden/accounts"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

func TestPromptAccount(t *testing.T) {
	accounts.PromptAccount(utils.EmptyAddress, `../../../testdata/keystore`, "")
}
func panicOnNullValue() {
	var c []int
	c[0] = 0
}

func TestPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			//t.Error(err)
		} else {
			t.Error("should panic")
		}
	}()
	panicOnNullValue()
}

type T struct {
	a int
}

func TestStruct(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			//t.Error(err)
		} else {
			t.Error("should panic")
		}
	}()
	var a *T
	t.Logf("a.a=%d", a.a)
}
