package models

import (
	"testing"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/stretchr/testify/assert"
)

func TestModelDB_TransferStatus(t *testing.T) {
	m := setupDb(t)
	lockSecretHash := utils.NewRandomHash()
	msg1 := "1111"
	m.UpdateTransferStatus(lockSecretHash, TransferStatusCanCancel, msg1)

	ts, err := m.GetTransferStatus(lockSecretHash)
	assert.Empty(t, err)
	assert.EqualValues(t, lockSecretHash, ts.LockSecretHash)
	assert.EqualValues(t, TransferStatusCanCancel, ts.Status)
	assert.EqualValues(t, fmt.Sprintf("%s\n", msg1), ts.StatusMessage)

	msg2 := "2222"
	m.UpdateTransferStatus(lockSecretHash, TransferStatusCanNotCancel, msg2)
	ts2, err := m.GetTransferStatus(lockSecretHash)
	assert.Empty(t, err)
	assert.EqualValues(t, lockSecretHash, ts2.LockSecretHash)
	assert.EqualValues(t, TransferStatusCanNotCancel, ts2.Status)
	assert.EqualValues(t, fmt.Sprintf("%s\n%s\n", msg1, msg2), ts2.StatusMessage)
}
