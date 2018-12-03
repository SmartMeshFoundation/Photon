package daotest

import (
	"fmt"
	"math/big"
	"testing"

	"time"

	"sync"

	"github.com/SmartMeshFoundation/Photon/codefortest"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/stretchr/testify/assert"
)

func TestModelDB_NewReceivedTransfer(t *testing.T) {
	dao := codefortest.NewTestDB("")
	defer dao.CloseDB()
	taddr := utils.NewRandomAddress()
	caddr := utils.NewRandomHash()
	lockSecertHash := utils.NewRandomHash()
	dao.NewReceivedTransfer(2, caddr, taddr, taddr, 3, big.NewInt(10), lockSecertHash, "123")
	key := fmt.Sprintf("%s-%d", caddr.String(), 3)
	r, err := dao.GetReceivedTransfer(key)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, r.FromAddress, taddr)
	assert.Equal(t, r.ChannelIdentifier, caddr)
	assert.EqualValues(t, r.Nonce, 3)
	assert.EqualValues(t, r.Amount, big.NewInt(10))
	dao.NewReceivedTransfer(3, caddr, taddr, taddr, 4, big.NewInt(10), lockSecertHash, "123")
	dao.NewReceivedTransfer(5, caddr, taddr, taddr, 6, big.NewInt(10), lockSecertHash, "123")

	trs, err := dao.GetReceivedTransferInBlockRange(0, 3)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, len(trs), 2)
	trs, err = dao.GetReceivedTransferInBlockRange(0, 5)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, len(trs), 3)

	trs, err = dao.GetReceivedTransferInBlockRange(0, 1)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, len(trs), 0)
}

func TestModelDB_NewSentTransfer(t *testing.T) {
	dao := codefortest.NewTestDB("")
	defer dao.CloseDB()
	taddr := utils.NewRandomAddress()
	caddr := utils.NewRandomHash()
	lockSecertHash := utils.NewRandomHash()
	dao.NewSentTransfer(2, caddr, taddr, taddr, 3, big.NewInt(10), lockSecertHash, "123")
	key := fmt.Sprintf("%s-%d", caddr.String(), 3)
	r, err := dao.GetSentTransfer(key)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, r.ToAddress, taddr)
	assert.Equal(t, r.ChannelIdentifier, caddr)
	assert.EqualValues(t, r.Nonce, 3)
	assert.EqualValues(t, r.Amount, big.NewInt(10))

	lockSecertHash = utils.NewRandomHash()
	dao.NewSentTransfer(3, caddr, taddr, taddr, 4, big.NewInt(10), lockSecertHash, "123")
	lockSecertHash = utils.NewRandomHash()
	dao.NewSentTransfer(5, caddr, taddr, taddr, 6, big.NewInt(10), lockSecertHash, "123")

	trs, err := dao.GetSentTransferInBlockRange(0, 3)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, len(trs), 2)
	trs, err = dao.GetSentTransferInBlockRange(0, 5)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, len(trs), 3)

	trs, err = dao.GetSentTransferInBlockRange(0, 1)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, len(trs), 0)
}

func TestBatchWriteDb(t *testing.T) {
	dao := codefortest.NewTestDB("")
	defer dao.CloseDB()
	caddr := utils.NewRandomHash()
	taddr := utils.NewRandomAddress()
	lockSecertHash := utils.NewRandomHash()
	dao.NewTransferStatus(taddr, lockSecertHash)
	number := float64(1000)
	wg := sync.WaitGroup{}
	wg.Add(int(number))
	begin := time.Now()
	for i := uint64(0); i < uint64(number); i++ {
		go func(index uint64) {
			//b := time.Now()
			//dao.SaveLatestBlockNumber(111)
			//dao.UpdateTransferStatusMessage(taddr, lockSecertHash, strconv.Itoa(int(index)))
			dao.NewSentTransfer(3, caddr, taddr, taddr, index, big.NewInt(10), lockSecertHash, "123")
			//fmt.Println("use ", time.Since(b).Seconds())
			wg.Done()
		}(i)
	}
	wg.Wait()
	total := time.Since(begin).Seconds()
	fmt.Println("total use seconds ", total)
	fmt.Println("avg use seconds ", total/number)
}
