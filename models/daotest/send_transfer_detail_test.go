package daotest

import (
	"testing"

	"fmt"

	"math/big"

	"github.com/SmartMeshFoundation/Photon/codefortest"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/stretchr/testify/assert"
)

func TestModelDB_SentTransferDetail(t *testing.T) {
	dao := codefortest.NewTestDB("")
	defer dao.CloseDB()
	dao.SaveLatestBlockNumber(1)
	tokenAddress := utils.NewRandomAddress()
	target := utils.NewRandomAddress()
	amount := big.NewInt(1)
	data := "123"
	lockSecretHash := utils.NewRandomHash()
	dao.NewSentTransferDetail(tokenAddress, target, amount, data, false, lockSecretHash)

	std, err := dao.GetSentTransferDetail(tokenAddress, lockSecretHash)
	assert.Empty(t, err)
	assert.EqualValues(t, std.Status, models.TransferStatusInit)
	fmt.Println(utils.StringInterface(std, 0))

	dao.UpdateSentTransferDetailStatus(tokenAddress, lockSecretHash, models.TransferStatusSuccess, "msg1", nil)

	list, err := dao.GetSentTransferDetailList(utils.EmptyAddress, -1, -1, -1, -1)
	fmt.Println(utils.StringInterface(list, 0))
	assert.Empty(t, err)
	assert.EqualValues(t, 1, len(list))
	assert.EqualValues(t, list[0].Status, models.TransferStatusSuccess)

	lockSecretHash2 := utils.NewRandomHash()
	dao.NewSentTransferDetail(tokenAddress, target, amount, data, false, lockSecretHash2)

	list, err = dao.GetSentTransferDetailList(tokenAddress, -1, -1, -1, -1)
	fmt.Println(utils.StringInterface(list, 0))
	assert.Empty(t, err)
	assert.EqualValues(t, 2, len(list))
}
