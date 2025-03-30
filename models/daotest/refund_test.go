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
	defer dao.CloseDB()
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

	l := dao.GetSendAnnounceDisposeByChannel(ch, false)
	assert.EqualValues(t, 1, len(l))
	assert.EqualValues(t, false, l[0].IsSubmitToPms)

	dao.MarkSendAnnounceDisposeSubmittedByChannel(ch)

	l = dao.GetSendAnnounceDisposeByChannel(ch, true)
	assert.EqualValues(t, 1, len(l))
	assert.EqualValues(t, true, l[0].IsSubmitToPms)
}

func TestNewReceivedAnnounceDisposed(t *testing.T) {
	lockHash := utils.NewRandomHash()
	channel := utils.NewRandomHash()
	r := models.NewReceivedAnnounceDisposed(lockHash, channel, utils.NewRandomHash(), 3, nil)
	dao := codefortest.NewTestDB("")
	defer dao.CloseDB()
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

	list := dao.GetChannelAnnounceDisposed(channel)
	for _, l := range list {
		assert.EqualValues(t, false, l.IsSubmittedToPms)
	}
	dao.MarkLockHashCanPunishSubmittedByChannel(channel)
	list = dao.GetChannelAnnounceDisposed(channel)
	for _, l := range list {
		assert.EqualValues(t, true, l.IsSubmittedToPms)
	}
}
