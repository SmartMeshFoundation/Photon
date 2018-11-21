package models

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/stretchr/testify/assert"
)

func TestModelDB_NewReceivedTransfer(t *testing.T) {
	m := setupDb(t)
	taddr := utils.NewRandomAddress()
	caddr := utils.NewRandomHash()
	lockSecertHash := utils.NewRandomHash()
	m.NewReceivedTransfer(2, caddr, taddr, taddr, 3, big.NewInt(10), lockSecertHash)
	key := fmt.Sprintf("%s-%d", caddr.String(), 3)
	r, err := m.GetReceivedTransfer(key)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, r.FromAddress, taddr)
	assert.Equal(t, r.ChannelIdentifier, caddr)
	assert.EqualValues(t, r.Nonce, 3)
	assert.EqualValues(t, r.Amount, big.NewInt(10))
	m.NewReceivedTransfer(3, caddr, taddr, taddr, 4, big.NewInt(10), lockSecertHash)
	m.NewReceivedTransfer(5, caddr, taddr, taddr, 6, big.NewInt(10), lockSecertHash)

	trs, err := m.GetReceivedTransferInBlockRange(0, 3)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, len(trs), 2)
	trs, err = m.GetReceivedTransferInBlockRange(0, 5)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, len(trs), 3)

	trs, err = m.GetReceivedTransferInBlockRange(0, 1)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, len(trs), 0)
}

func TestModelDB_NewSentTransfer(t *testing.T) {
	m := setupDb(t)
	taddr := utils.NewRandomAddress()
	caddr := utils.NewRandomHash()
	lockSecertHash := utils.NewRandomHash()
	m.NewSentTransfer(2, caddr, taddr, taddr, 3, big.NewInt(10), lockSecertHash)
	key := fmt.Sprintf("%s-%d", caddr.String(), 3)
	r, err := m.GetSentTransfer(key)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, r.ToAddress, taddr)
	assert.Equal(t, r.ChannelIdentifier, caddr)
	assert.EqualValues(t, r.Nonce, 3)
	assert.EqualValues(t, r.Amount, big.NewInt(10))

	lockSecertHash = utils.NewRandomHash()
	m.NewSentTransfer(3, caddr, taddr, taddr, 4, big.NewInt(10), lockSecertHash)
	lockSecertHash = utils.NewRandomHash()
	m.NewSentTransfer(5, caddr, taddr, taddr, 6, big.NewInt(10), lockSecertHash)

	trs, err := m.GetSentTransferInBlockRange(0, 3)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, len(trs), 2)
	trs, err = m.GetSentTransferInBlockRange(0, 5)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, len(trs), 3)

	trs, err = m.GetSentTransferInBlockRange(0, 1)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, len(trs), 0)
}
