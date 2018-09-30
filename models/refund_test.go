package models

import (
	"testing"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/stretchr/testify/assert"
)

func TestModelDB_MarkLockSecretHashDisposed(t *testing.T) {
	model := setupDb(t)
	defer func() {
		model.CloseDB()
	}()
	lock1 := utils.NewRandomHash()
	lock2 := utils.NewRandomHash()
	ch := utils.NewRandomHash()
	err := model.MarkLockSecretHashDisposed(lock1, ch)
	if err != nil {
		t.Error(err)
	}
	r := model.IsLockSecretHashDisposed(lock1)
	assert.EqualValues(t, r, true)
	r = model.IsLockSecretHashDisposed(lock2)
	assert.EqualValues(t, r, false)
	r = model.IsLockSecretHashChannelIdentifierDisposed(lock1, ch)
	assert.EqualValues(t, r, true)
	r = model.IsLockSecretHashChannelIdentifierDisposed(lock2, ch)
	assert.EqualValues(t, r, false)
}

func TestNewReceivedAnnounceDisposed(t *testing.T) {
	lockHash := utils.NewRandomHash()
	channel := utils.NewRandomHash()
	r := NewReceivedAnnounceDisposed(lockHash, channel, utils.NewRandomHash(), 3, nil)
	model := setupDb(t)
	defer func() {
		model.CloseDB()
	}()
	err := model.MarkLockHashCanPunish(r)
	if err != nil {
		t.Error(err)
		return
	}
	b := model.IsLockHashCanPunish(lockHash, channel)
	assert.EqualValues(t, b, true)
	b = model.IsLockHashCanPunish(utils.NewRandomHash(), channel)
	assert.EqualValues(t, b, false)
	b = model.IsLockHashCanPunish(lockHash, utils.NewRandomHash())
	assert.EqualValues(t, b, false)
	r2 := model.GetReceiviedAnnounceDisposed(lockHash, channel)
	assert.EqualValues(t, r, r2)
	rs := model.GetChannelAnnounceDisposed(channel)
	assert.EqualValues(t, len(rs), 1)
	rs = model.GetChannelAnnounceDisposed(utils.NewRandomHash())
	assert.EqualValues(t, len(rs), 0)
	r2 = model.GetReceiviedAnnounceDisposed(lockHash, utils.NewRandomHash())
	if r2 != nil {
		t.Error("should be nil")
		return
	}
}
