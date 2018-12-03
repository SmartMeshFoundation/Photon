package daotest

import (
	"fmt"
	"testing"

	"github.com/SmartMeshFoundation/Photon/codefortest"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

var token common.Address

func TestModelDB_NewNonParticipantChannel(t *testing.T) {
	dao := codefortest.NewTestDB("")
	defer func() {
		dao.CloseDB()
	}()
	token = utils.NewRandomAddress()
	p1 := utils.NewRandomAddress()
	p2 := utils.NewRandomAddress()
	channel := utils.Sha3(p1[:], p2[:], token[:])
	err := dao.NewNonParticipantChannel(token, channel, p1, p2)
	if err != nil {
		t.Error(err)
		return
	}
	p3 := utils.NewRandomAddress()
	channel2 := utils.Sha3(p1[:], p3[:], token[:])
	err = dao.NewNonParticipantChannel(token, channel2, p1, p3)
	if err != nil {
		t.Error(err)
		return
	}
	edges, err := dao.GetAllNonParticipantChannel(token)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, len(edges), 4)
	err = dao.RemoveNonParticipantChannel(token, utils.NewRandomHash())
	assert.EqualValues(t, err != nil, true)
	err = dao.RemoveNonParticipantChannel(token, channel2)
	assert.EqualValues(t, err == nil, true)
	edges, err = dao.GetAllNonParticipantChannel(token)
	assert.EqualValues(t, len(edges), 2)
}
func TestReadDbAgain(t *testing.T) {
	TestModelDB_NewNonParticipantChannel(t)
	dao := codefortest.NewTestDB("")
	defer func() {
		dao.CloseDB()
	}()
	edges, err := dao.GetAllNonParticipantChannel(token)
	if err != nil {
		t.Error(err)
		return
	}
	log.Trace(fmt.Sprintf("len edges=%d", len(edges)))
}
