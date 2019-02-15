package daotest

import (
	"math/big"
	"testing"

	"github.com/SmartMeshFoundation/Photon/codefortest"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestModelDB_ContractCallTX(t *testing.T) {
	dao := codefortest.NewTestDB("")
	defer dao.CloseDB()
	channelIdentifier := utils.NewRandomHash()
	openBlockNumber := int64(5)

	// 0. find
	list, err := dao.GetTXInfoList(utils.EmptyHash, 0, "", models.TXInfoStatusPending)
	assert.Empty(t, err)
	assert.EqualValues(t, 0, len(list))

	// 1. new
	tx := types.NewTransaction(1, utils.NewRandomAddress(), big.NewInt(1), 0, nil, nil)
	_, err = dao.NewPendingTXInfo(tx, models.TXInfoTypeDeposit, channelIdentifier, openBlockNumber, "")
	assert.Empty(t, err)

	// 2. find all
	list, err = dao.GetTXInfoList(utils.EmptyHash, 0, "", "")
	assert.Empty(t, err)
	assert.EqualValues(t, 1, len(list))

	// 3. find by channelIdentifier
	list, err = dao.GetTXInfoList(channelIdentifier, 0, "", "")
	assert.Empty(t, err)
	assert.EqualValues(t, 1, len(list))

	list, err = dao.GetTXInfoList(utils.NewRandomHash(), 0, "", "")
	assert.Empty(t, err)
	assert.EqualValues(t, 0, len(list))

	// 4. find by OpenBlockNumber
	list, err = dao.GetTXInfoList(utils.EmptyHash, openBlockNumber, "", "")
	assert.Empty(t, err)
	assert.EqualValues(t, 1, len(list))

	list, err = dao.GetTXInfoList(utils.EmptyHash, 2, "", "")
	assert.Empty(t, err)
	assert.EqualValues(t, 0, len(list))

	// 5. find by channelIdentifier && OpenBlockNumber
	list, err = dao.GetTXInfoList(channelIdentifier, openBlockNumber, "", "")
	assert.Empty(t, err)
	assert.EqualValues(t, 1, len(list))

	list, err = dao.GetTXInfoList(channelIdentifier, 2, "", "")
	assert.Empty(t, err)
	assert.EqualValues(t, 0, len(list))

	// 8. find by type
	list, err = dao.GetTXInfoList(utils.EmptyHash, 0, models.TXInfoTypeDeposit, "")
	assert.Empty(t, err)
	assert.EqualValues(t, 1, len(list))

	list, err = dao.GetTXInfoList(utils.EmptyHash, 0, models.TXInfoTypeClose, "")
	assert.Empty(t, err)
	assert.EqualValues(t, 0, len(list))

	// 7. find by status
	list, err = dao.GetTXInfoList(utils.EmptyHash, 0, "", models.TXInfoStatusPending)
	assert.Empty(t, err)
	assert.EqualValues(t, 1, len(list))

	list, err = dao.GetTXInfoList(utils.EmptyHash, 0, "", models.TXInfoStatusSuccess)
	assert.Empty(t, err)
	assert.EqualValues(t, 0, len(list))

	// 8. update
	err = dao.UpdateTXInfoStatus(tx.Hash(), models.TXInfoStatusSuccess, 2)
	assert.Empty(t, err)

	list, err = dao.GetTXInfoList(utils.EmptyHash, 0, "", "")
	assert.Empty(t, err)
	assert.EqualValues(t, 1, len(list))
	assert.EqualValues(t, models.TXInfoStatusSuccess, list[0].Status)
	assert.EqualValues(t, 2, list[0].PendingBlockNumber)
}
