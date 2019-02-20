package daotest

import (
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/SmartMeshFoundation/Photon/codefortest"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/stretchr/testify/assert"
)

func TestModelDB_NewReceivedTransfer(t *testing.T) {
	dao := codefortest.NewTestDB("")
	defer dao.CloseDB()
	taddr := utils.NewRandomAddress()
	caddr := utils.NewRandomHash()
	var openBlockNumber int64 = 3
	lockSecertHash := utils.NewRandomHash()
	dao.NewReceivedTransfer(2, caddr, openBlockNumber, taddr, taddr, 3, big.NewInt(10), lockSecertHash, "123")
	key := fmt.Sprintf("%s-%d-%d", caddr.String(), openBlockNumber, 3)
	r, err := dao.GetReceivedTransfer(key)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, r.FromAddress, taddr)
	assert.Equal(t, r.ChannelIdentifier, caddr)
	assert.EqualValues(t, r.Nonce, 3)
	assert.EqualValues(t, r.Amount, big.NewInt(10))
	dao.NewReceivedTransfer(3, caddr, openBlockNumber, taddr, taddr, 4, big.NewInt(10), lockSecertHash, "123")
	dao.NewReceivedTransfer(5, caddr, openBlockNumber, taddr, taddr, 6, big.NewInt(10), lockSecertHash, "123")

	trs, err := dao.GetReceivedTransferList(utils.EmptyAddress, 0, 3)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, len(trs), 2)
	trs, err = dao.GetReceivedTransferList(utils.EmptyAddress, 0, 5)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, len(trs), 3)

	trs, err = dao.GetReceivedTransferList(utils.EmptyAddress, 0, 1)
	if err != nil {
		t.Error(err)
		return
	}
	//assert.EqualValues(t, len(trs), 0)
	//from := time.Now().Add(0 - time.Minute)
	//to := time.Now().Add(time.Minute)
	//trs, err = dao.GetReceivedTransferList(utils.EmptyAddress, from, to)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//assert.EqualValues(t, len(trs), 3)
	//from = time.Now().Add(time.Second)
	//trs, err = dao.GetReceivedTransferInTimeRange(from, to)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//assert.EqualValues(t, len(trs), 0)
}

func TestBatchWriteDb(t *testing.T) {
	dao := codefortest.NewTestDB("")
	defer dao.CloseDB()
	//caddr := utils.NewRandomHash()
	//var openBlockNumber int64 = 3
	taddr := utils.NewRandomAddress()
	lockSecertHash := utils.NewRandomHash()
	number := float64(1000)
	wg := sync.WaitGroup{}
	wg.Add(int(number))
	begin := time.Now()
	for i := uint64(0); i < uint64(number); i++ {
		go func(index uint64) {
			//b := time.Now()
			//dao.SaveLatestBlockNumber(111)
			//dao.UpdateTransferStatusMessage(taddr, lockSecertHash, strconv.Itoa(int(index)))
			dao.NewSentTransferDetail(utils.NewRandomAddress(), taddr, big.NewInt(10), "123", true, lockSecertHash)
			//dao.NewSentTransfer(3, caddr, openBlockNumber, taddr, taddr, index, big.NewInt(10), lockSecertHash, "123")
			//fmt.Println("use ", time.Since(b).Seconds())
			wg.Done()
		}(i)
	}
	wg.Wait()
	total := time.Since(begin).Seconds()
	fmt.Println("total use seconds ", total)
	fmt.Println("avg use seconds ", total/number)
}
