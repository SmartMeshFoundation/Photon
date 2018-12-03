package daotest

import (
	"testing"

	"github.com/SmartMeshFoundation/Photon/codefortest"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/stretchr/testify/assert"
)

func TestModelDB_MarkLockSecretHashDisposed(t *testing.T) {
	dao := codefortest.NewTestDB("")
	defer func() {
		dao.CloseDB()
	}()
	lock1 := utils.NewRandomHash()
	lock2 := utils.NewRandomHash()
	ch := utils.NewRandomHash()
	err := dao.MarkLockSecretHashDisposed(lock1, ch)
	if err != nil {
		t.Error(err)
	}
	r := dao.IsLockSecretHashDisposed(lock1)
	assert.EqualValues(t, r, true)
	r = dao.IsLockSecretHashDisposed(lock2)
	assert.EqualValues(t, r, false)
	r = dao.IsLockSecretHashChannelIdentifierDisposed(lock1, ch)
	assert.EqualValues(t, r, true)
	r = dao.IsLockSecretHashChannelIdentifierDisposed(lock2, ch)
	assert.EqualValues(t, r, false)
}

func TestNewReceivedAnnounceDisposed(t *testing.T) {
	lockHash := utils.NewRandomHash()
	channel := utils.NewRandomHash()
	r := models.NewReceivedAnnounceDisposed(lockHash, channel, utils.NewRandomHash(), 3, nil)
	dao := codefortest.NewTestDB("")
	defer func() {
		dao.CloseDB()
	}()
	err := dao.MarkLockHashCanPunish(r)
	if err != nil {
		t.Error(err)
		return
	}
	b := dao.IsLockHashCanPunish(lockHash, channel)
	assert.EqualValues(t, b, true)
	b = dao.IsLockHashCanPunish(utils.NewRandomHash(), channel)
	assert.EqualValues(t, b, false)
	b = dao.IsLockHashCanPunish(lockHash, utils.NewRandomHash())
	assert.EqualValues(t, b, false)
	r2 := dao.GetReceivedAnnounceDisposed(lockHash, channel)
	assert.EqualValues(t, r, r2)
	rs := dao.GetChannelAnnounceDisposed(channel)
	assert.EqualValues(t, len(rs), 1)
	rs = dao.GetChannelAnnounceDisposed(utils.NewRandomHash())
	assert.EqualValues(t, len(rs), 0)
	r2 = dao.GetReceivedAnnounceDisposed(lockHash, utils.NewRandomHash())
	if r2 != nil {
		t.Error("should be nil")
		return
	}
}
