package daotest

import (
	"testing"

	"fmt"

	"time"

	"github.com/SmartMeshFoundation/Photon/codefortest"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/stretchr/testify/assert"
)

func TestModelDB_TransferStatus(t *testing.T) {
	dao := codefortest.NewTestDB("")
	defer dao.CloseDB()
	lockSecretHash := utils.NewRandomHash()
	tokenAddress := utils.NewRandomAddress()
	dao.NewTransferStatus(tokenAddress, lockSecretHash)
	msg1 := "1111"
	dao.UpdateTransferStatus(tokenAddress, lockSecretHash, models.TransferStatusCanCancel, msg1)

	ts, err := dao.GetTransferStatus(tokenAddress, lockSecretHash)
	assert.Empty(t, err)
	assert.EqualValues(t, lockSecretHash, ts.LockSecretHash)
	assert.EqualValues(t, models.TransferStatusCanCancel, ts.Status)
	assert.EqualValues(t, fmt.Sprintf("%s\n", msg1), ts.StatusMessage)

	msg2 := "2222"
	dao.UpdateTransferStatus(tokenAddress, lockSecretHash, models.TransferStatusCanNotCancel, msg2)
	ts2, err := dao.GetTransferStatus(tokenAddress, lockSecretHash)
	assert.Empty(t, err)
	assert.EqualValues(t, lockSecretHash, ts2.LockSecretHash)
	assert.EqualValues(t, models.TransferStatusCanNotCancel, ts2.Status)
	assert.EqualValues(t, fmt.Sprintf("%s\n%s\n", msg1, msg2), ts2.StatusMessage)
}

func TestModelDb_BatchTransferStatus(t *testing.T) {
	dao := codefortest.NewTestDB("")
	defer dao.CloseDB()
	lockSecretHash := utils.NewRandomHash()
	tokenAddress := utils.NewRandomAddress()
	dao.NewTransferStatus(tokenAddress, lockSecretHash)
	msg1 := "1111"

	// write once
	start := time.Now()
	dao.UpdateTransferStatusMessage(tokenAddress, lockSecretHash, msg1)
	fmt.Println("update once use ", time.Since(start))

	// write sync
	start = time.Now()
	i := 0
	for i < 100 {
		dao.UpdateTransferStatusMessage(tokenAddress, lockSecretHash, msg1)
		i++
	}
	fmt.Println("update 100 times sync use ", time.Since(start))

	//// write async
	//start = time.Now()
	//i = 0
	//wg := sync.WaitGroup{}
	//wg.Add(1000)
	//for i < 1000 {
	//	go func() {
	//		dao.UpdateTransferStatusMessage(tokenAddress, lockSecretHash, msg1)
	//		wg.Done()
	//	}()
	//	i++
	//}
	//wg.Wait()
	//fmt.Println("update 100 times async use ", time.Since(start))
}
