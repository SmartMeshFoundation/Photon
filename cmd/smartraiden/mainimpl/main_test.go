package mainimpl

import (
	"testing"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

func TestPromptAccount(t *testing.T) {
	promptAccount(utils.EmptyAddress, `../../../testdata/keystore`, "")
}
