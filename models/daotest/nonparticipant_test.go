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
	defer dao.CloseDB()
	token = utils.NewRandomAddress()
	p1 := utils.NewRandomAddress()
	p2 := utils.NewRandomAddress()
	channel := utils.Sha3(p1[:], p2[:], token[:])
	err := dao.NewNonParticipantChannel(token, channel, p1, p2)
	if err != nil {
		t.Error(err)
		return
	}
	err = dao.NewNonParticipantChannel(token, channel, p1, p2)
	if err == nil {
		t.Error("must report duplicate")
		return
	}
	p3 := utils.NewRandomAddress()
	channel2 := utils.Sha3(p1[:], p3[:], token[:])
	err = dao.NewNonParticipantChannel(token, channel2, p1, p3)
	if err != nil {
		t.Error(err)
		return
	}
	edges, err := dao.GetAllNonParticipantChannelByToken(token)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, len(edges), 4)
	err = dao.RemoveNonParticipantChannel(utils.NewRandomHash())
	assert.EqualValues(t, err != nil, true)
	err = dao.RemoveNonParticipantChannel(channel2)
	assert.EqualValues(t, err, nil)
	edges, err = dao.GetAllNonParticipantChannelByToken(token)
	assert.EqualValues(t, 2, len(edges))
}
func TestReadDbAgain(t *testing.T) {
	TestModelDB_NewNonParticipantChannel(t)
	dao := codefortest.NewTestDB("")
	defer dao.CloseDB()
	edges, err := dao.GetAllNonParticipantChannelByToken(token)
	if err != nil {
		t.Error(err)
		return
	}
	log.Trace(fmt.Sprintf("len edges=%d", len(edges)))
}
