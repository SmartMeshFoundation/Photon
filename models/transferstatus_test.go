package models

import (
	"testing"

	"fmt"

	"time"

	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/stretchr/testify/assert"
)

func TestModelDB_TransferStatus(t *testing.T) {
	m := setupDb(t)
	lockSecretHash := utils.NewRandomHash()
	tokenAddress := utils.NewRandomAddress()
	m.NewTransferStatus(tokenAddress, lockSecretHash)
	msg1 := "1111"
	m.UpdateTransferStatus(tokenAddress, lockSecretHash, TransferStatusCanCancel, msg1)

	ts, err := m.GetTransferStatus(tokenAddress, lockSecretHash)
	assert.Empty(t, err)
	assert.EqualValues(t, lockSecretHash, ts.LockSecretHash)
	assert.EqualValues(t, TransferStatusCanCancel, ts.Status)
	assert.EqualValues(t, fmt.Sprintf("%s\n", msg1), ts.StatusMessage)

	msg2 := "2222"
	m.UpdateTransferStatus(tokenAddress, lockSecretHash, TransferStatusCanNotCancel, msg2)
	ts2, err := m.GetTransferStatus(tokenAddress, lockSecretHash)
	assert.Empty(t, err)
	assert.EqualValues(t, lockSecretHash, ts2.LockSecretHash)
	assert.EqualValues(t, TransferStatusCanNotCancel, ts2.Status)
	assert.EqualValues(t, fmt.Sprintf("%s\n%s\n", msg1, msg2), ts2.StatusMessage)
}

func TestModelDb_BatchTransferStatus(t *testing.T) {
	m := setupDb(t)
	lockSecretHash := utils.NewRandomHash()
	tokenAddress := utils.NewRandomAddress()
	m.NewTransferStatus(tokenAddress, lockSecretHash)
	msg1 := "1111"

	// write once
	start := time.Now()
	m.UpdateTransferStatusMessage(tokenAddress, lockSecretHash, msg1)
	fmt.Println("update once use ", time.Since(start))

	// write sync
	start = time.Now()
	i := 0
	for i < 1000 {
		m.UpdateTransferStatusMessage(tokenAddress, lockSecretHash, msg1)
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
	//		m.UpdateTransferStatusMessage(tokenAddress, lockSecretHash, msg1)
	//		wg.Done()
	//	}()
	//	i++
	//}
	//wg.Wait()
	//fmt.Println("update 100 times async use ", time.Since(start))
}
