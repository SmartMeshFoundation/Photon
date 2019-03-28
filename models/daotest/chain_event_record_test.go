package daotest

import (
	"testing"

	"time"

	"github.com/SmartMeshFoundation/Photon/codefortest"
	"github.com/SmartMeshFoundation/Photon/models"
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
	assert.EqualValues(t, 0, blockNumber)

}

func Test1(t *testing.T) {
	dbPath := "./temp"
	dao := codefortest.NewTestDB(dbPath)
	defer dao.CloseDB()
	var idList []models.ChainEventID
	for i := uint(0); i < 100; i++ {
		l := new(types.Log)
		l.TxHash = utils.NewRandomHash()
		l.Index = i
		id := dao.MakeChainEventID(l)
		idList = append(idList, id)
		dao.NewDeliveredChainEvent(id, uint64(i)+1)
	}
	//fmt.Println("total==============", len(idList))
	//for _, id := range idList {
	//	b, err := dao.CheckChainEventDelivered(id)
	//	fmt.Println(common.Bytes2Hex(id[:]), b, err)
	//}
	//dao.ClearOldChainEventRecord(1000)

	//dao.CloseDB()
	//for _, id := range idList {
	//	b, err := dao.CheckChainEventDelivered(id)
	//	fmt.Println(b, err)
	//}
	//dao = codefortest.NewTestDB(dbPath)
	dao.ClearOldChainEventRecord(1000)
	dao.ClearOldChainEventRecord(1000)
	time.Sleep(5 * time.Second)
}
