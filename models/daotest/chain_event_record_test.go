package daotest

import (
	"testing"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/codefortest"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestChainEventRecord(t *testing.T) {
	dao := codefortest.NewTestDB("")
	defer dao.CloseDB()
	l1 := new(types.Log)
	l1.TxHash = utils.NewRandomHash()
	l1.Index = 1

	l2 := new(types.Log)
	l2.TxHash = utils.NewRandomHash()
	l2.Index = 2

	id1 := dao.MakeChainEventID(l1)
	id2 := dao.MakeChainEventID(l2)
	fmt.Println(id1, len(id1))
	fmt.Println(id2, len(id2))

	dao.NewDeliveredChainEvent(id1, 1)

	blockNumber, delivered := dao.CheckChainEventDelivered(id1)
	assert.EqualValues(t, true, delivered)
	assert.EqualValues(t, 1, blockNumber)

	blockNumber, delivered = dao.CheckChainEventDelivered(id2)
	assert.EqualValues(t, false, delivered)
	assert.EqualValues(t, 0, blockNumber)

	dao.ClearOldChainEventRecord(0)

	blockNumber, delivered = dao.CheckChainEventDelivered(id1)
	assert.EqualValues(t, true, delivered)
	assert.EqualValues(t, 1, blockNumber)

	dao.ClearOldChainEventRecord(100)

	blockNumber, delivered = dao.CheckChainEventDelivered(id1)
	assert.EqualValues(t, false, delivered)

}
