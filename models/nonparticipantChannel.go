package models

import (
	"fmt"

	"math/big"

	"bytes"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
)

type NonParticipantChannel struct {
	Participant1        common.Address
	Participant2        common.Address
	Participant1Balance *big.Int
	Participant2Balance *big.Int
}

/*
如果这个 map 很大,怎么办?存储效率肯定会很低.
否则怎么遍历呢?
*/
type ChannelParticipantMap map[common.Address]common.Address

const bucketChannel = "bucketChannel"

func (model *ModelDB) NewNonParticipantChannel(token, participant1, participant2 common.Address) error {
	var m ChannelParticipantMap
	err := model.db.Get(bucketChannel, token, &m)
	if err != nil {
		if err == storm.ErrNotFound {
			m = make(ChannelParticipantMap)
		} else {
			return err
		}

	}
	if participant1 == participant2 {
		panic(fmt.Sprintf("channel error, p1 andf p2 is the same,token=%s,participant=%s", token.String(), participant1.String()))
	}
	if bytes.Compare(participant1[:], participant2[:]) > 0 {
		participant1, participant2 = participant2, participant1
	}
	if m[participant1] != utils.EmptyAddress {
		//startup ...
		log.Warn("add channel ,but channel already exists, maybe duplicates channelnew events")
		return nil
	}
	m[participant1] = participant2
	err = model.db.Set(bucketToken, keyToken, m)
	return err
}
func (model *ModelDB) RemoveNonParticipantChannel(token, participant1, participant2 common.Address) error {
	var m ChannelParticipantMap
	err := model.db.Get(bucketChannel, token, &m)
	if err != nil {
		if err == storm.ErrNotFound {
			m = make(ChannelParticipantMap)
		} else {
			return err
		}

	}
	if participant1 == participant2 {
		panic(fmt.Sprintf("channel error, p1 andf p2 is the same,token=%s,participant=%s", token.String(), participant1.String()))
	}
	if bytes.Compare(participant1[:], participant2[:]) > 0 {
		participant1, participant2 = participant2, participant1
	}
	if m[participant1] == utils.EmptyAddress {
		//startup ...
		log.Warn("delete channel ,but channel don't exists")
		return nil
	}
	delete(m, participant1)
	err = model.db.Set(bucketToken, keyToken, m)
	return err
}

//GetAllTokens returna all tokens on this registry contract
func (model *ModelDB) GetAllNonParticipantChannel(token common.Address) (edges ChannelParticipantMap, err error) {
	err = model.db.Get(bucketChannel, token, &edges)
	if err == storm.ErrNotFound {
		err = nil
	}
	return
}
