package models

import (
	"testing"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/stretchr/testify/assert"
)

func TestModelDB_NewNonParticipantChannel(t *testing.T) {
	model := setupDb(t)
	defer func() {
		model.CloseDB()
	}()
	token := utils.NewRandomAddress()
	p1 := utils.NewRandomAddress()
	p2 := utils.NewRandomAddress()
	err := model.NewNonParticipantChannel(token, p1, p2)
	if err != nil {
		t.Error(err)
		return
	}
	p3 := utils.NewRandomAddress()
	err = model.NewNonParticipantChannel(token, p1, p3)
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
	err = model.RemoveNonParticipantChannel(token, p1, utils.NewRandomAddress())
	assert.EqualValues(t, err != nil, true)
	err = model.RemoveNonParticipantChannel(token, p1, p3)
	assert.EqualValues(t, err == nil, true)
	edges, err = model.GetAllNonParticipantChannel(token)
	assert.EqualValues(t, len(edges), 2)
}
