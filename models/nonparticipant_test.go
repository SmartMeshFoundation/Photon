package models

import (
	"testing"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

var token common.Address

func TestModelDB_NewNonParticipantChannel(t *testing.T) {
	model := setupDb(t)
	defer func() {
		model.CloseDB()
	}()
	token = utils.NewRandomAddress()
	p1 := utils.NewRandomAddress()
	p2 := utils.NewRandomAddress()
	channel := utils.Sha3(p1[:], p2[:], token[:])
	err := model.NewNonParticipantChannel(token, channel, p1, p2)
	if err != nil {
		t.Error(err)
		return
	}
	p3 := utils.NewRandomAddress()
	channel2 := utils.Sha3(p1[:], p3[:], token[:])
	err = model.NewNonParticipantChannel(token, channel2, p1, p3)
	if err != nil {
		t.Error(err)
		return
	}
	edges, err := model.GetAllNonParticipantChannel(token)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, len(edges), 4)
	err = model.RemoveNonParticipantChannel(token, utils.NewRandomHash())
	assert.EqualValues(t, err != nil, true)
	err = model.RemoveNonParticipantChannel(token, channel2)
	assert.EqualValues(t, err == nil, true)
	edges, err = model.GetAllNonParticipantChannel(token)
	assert.EqualValues(t, len(edges), 2)
}
func TestReadDbAgain(t *testing.T) {
	TestModelDB_NewNonParticipantChannel(t)
	model, err := OpenDb(dbPath)
	if err != nil {
		t.Error(err)
		return
	}
	defer model.CloseDB()
	edges, err := model.GetAllNonParticipantChannel(token)
	if err != nil {
		t.Error(err)
		return
	}
	log.Trace(fmt.Sprintf("edges=%s", utils.StringInterface(edges, 3)))
}
