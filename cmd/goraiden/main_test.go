package main

import (
	"testing"

	"github.com/SmartMeshFoundation/raiden-network/utils"
)

func TestPromptAccount(t *testing.T) {
	promptAccount(utils.EmptyAddress, `D:\privnet\keystore\`, "")
}
