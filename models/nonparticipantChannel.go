package models

import (
	"fmt"

	"bytes"

	"math/big"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
)

/*
NonParticipantChannel 所有的通道信息在本地的存储
因为合约不提供直接查询通道信息,只能通过事件获取,所以需要在本地保存一份,以便查询
*/
/*
 *	NonParticipantChannel : structure for back up of channel information at local storage.
 *	Because contract does not provide direct check for channel information, so we need to backup at local storage.
 */
type NonParticipantChannel struct {
	Participant1        common.Address
	Participant2        common.Address
	Participant1Balance *big.Int
	Participant2Balance *big.Int
}

func participant2bytes(p1, p2 common.Address) []byte {
	b := make([]byte, len(p1)*2)
	copy(b[0:len(p1)], p1[:])
	copy(b[len(p1):], p2[:])
	return b
}
func bytes2participant(data []byte) (p1, p2 common.Address) {
	if len(data) != len(p1)*2 {
		return
	}
	copy(p1[:], data[:len(p1)])
	copy(p2[:], data[len(p1):])
	return
}
func participantKey(p1, p2 common.Address) common.Address {
	t := utils.Sha3(p1[:], p2[:])
	return common.BytesToAddress(t[:])
}

/*
ChannelParticipantMap :
todo 如果这个 map 很大,怎么办?存储效率肯定会很低.
否则怎么遍历呢?
*/
type ChannelParticipantMap map[common.Hash][]byte

const bucketChannel = "bucketChannel"

//NewNonParticipantChannel 需要保存 channel identifier, 通道的事件都是与此有关系的
func (model *ModelDB) NewNonParticipantChannel(token common.Address, channel common.Hash, participant1, participant2 common.Address) error {
	var m ChannelParticipantMap
	log.Trace(fmt.Sprintf("NewNonParticipantChannel token=%s,participant1=%s,participant2=%s",
		utils.APex2(token),
		utils.APex2(participant1),
		utils.APex2(participant2),
	))
	err := model.db.Get(bucketChannel, token[:], &m)
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
	key := channel
	if m[key] != nil {
		//startup ...
		log.Warn(fmt.Sprintf("add channel ,but channel already exists, maybe duplicates channelnew events,participant1=%s,participant2=%s",
			utils.APex2(participant1), utils.APex2(participant2)))
		return nil
	}
	m[key] = participant2bytes(participant1, participant2)
	log.Trace(fmt.Sprintf("NewNonParticipantChannel token=%s,p1=%s,p2=%s,len(m)=%d", utils.APex2(token),
		utils.APex2(participant1), utils.APex2(participant2), len(m)))
	err = model.db.Set(bucketChannel, token[:], m)
	return err
}

//RemoveNonParticipantChannel a channel is settled
func (model *ModelDB) RemoveNonParticipantChannel(token common.Address, channel common.Hash) error {
	var m ChannelParticipantMap
	err := model.db.Get(bucketChannel, token[:], &m)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil
		}
		return err
	}
	if m[channel] == nil {
		//startup ...
		return fmt.Errorf("delete channel ,but channel don't exists")
	}
	delete(m, channel)
	log.Trace(fmt.Sprintf("RemoveNonParticipantChannel token=%s,channel=%s", utils.APex2(token),
		utils.HPex(channel)))
	err = model.db.Set(bucketChannel, token[:], m)
	return err
}

//GetAllNonParticipantChannel returna all channel on this `token`
func (model *ModelDB) GetAllNonParticipantChannel(token common.Address) (edges []common.Address, err error) {
	var m ChannelParticipantMap
	err = model.db.Get(bucketChannel, token[:], &m)
	log.Trace(fmt.Sprintf("GetAllNonParticipantChannel,token=%s,err=%v", utils.APex2(token), err))
	if err == storm.ErrNotFound {
		err = nil
		return
	}
	for _, data := range m {
		p1, p2 := bytes2participant(data)
		edges = append(edges, p1, p2)
	}
	return
}
