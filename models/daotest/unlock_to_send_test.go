package daotest

import (
	"testing"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/codefortest"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/stretchr/testify/assert"
)

func TestModelDB_UnlockToSend(t *testing.T) {
	dao := codefortest.NewTestDB("")
	defer dao.CloseDB()

	lockSecretHash := utils.NewRandomHash()
	token := utils.NewRandomAddress()
	receiver := utils.NewRandomAddress()
	dao.NewUnlockToSend(lockSecretHash, token, receiver, 5)
	list := dao.GetAllUnlockToSend()
	fmt.Println(list)
	assert.EqualValues(t, 1, len(list))
	key := utils.Sha3(lockSecretHash[:], token[:], receiver[:]).Bytes()
	dao.RemoveUnlockToSend(key)
	list = dao.GetAllUnlockToSend()
	fmt.Println(list)
	assert.EqualValues(t, 0, len(list))
}
