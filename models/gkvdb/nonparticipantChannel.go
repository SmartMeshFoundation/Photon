package gkvdb

import (
	"fmt"

	"bytes"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

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

//NewNonParticipantChannel 需要保存 channel identifier, 通道的事件都是与此有关系的
func (dao *GkvDB) NewNonParticipantChannel(token common.Address, channel common.Hash, participant1, participant2 common.Address) error {
	var m models.ChannelParticipantMap
	log.Trace(fmt.Sprintf("NewNonParticipantChannel token=%s,participant1=%s,participant2=%s",
		utils.APex2(token),
		utils.APex2(participant1),
		utils.APex2(participant2),
	))
	err := dao.getKeyValueToBucket(models.BucketChannel, token[:], &m)
	if err != nil {
		if err == ErrorNotFound {
			m = make(models.ChannelParticipantMap)
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
	return dao.saveKeyValueToBucket(models.BucketChannel, token[:], m)
}

//RemoveNonParticipantChannel a channel is settled
func (dao *GkvDB) RemoveNonParticipantChannel(token common.Address, channel common.Hash) error {
	var m models.ChannelParticipantMap
	err := dao.getKeyValueToBucket(models.BucketChannel, token[:], &m)
	if err != nil {
		if err == ErrorNotFound {
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
	return dao.saveKeyValueToBucket(models.BucketChannel, token[:], m)
}

//GetAllNonParticipantChannel returna all channel on this `token`
func (dao *GkvDB) GetAllNonParticipantChannel(token common.Address) (edges []common.Address, err error) {
	var m models.ChannelParticipantMap
	err = dao.getKeyValueToBucket(models.BucketChannel, token[:], &m)
	log.Trace(fmt.Sprintf("GetAllNonParticipantChannel,token=%s,err=%v", utils.APex2(token), err))
	if err == ErrorNotFound {
		err = nil
		return
	}
	for _, data := range m {
		p1, p2 := bytes2participant(data)
		edges = append(edges, p1, p2)
	}
	return
}

// GetParticipantAddressByTokenAndChannel :
func (dao *GkvDB) GetParticipantAddressByTokenAndChannel(token common.Address, channel common.Hash) (p1, p2 common.Address) {
	var m models.ChannelParticipantMap
	err := dao.getKeyValueToBucket(models.BucketChannel, token[:], &m)
	log.Trace(fmt.Sprintf("GetAllNonParticipantChannel,token=%s,err=%v", utils.APex2(token), err))
	if err == ErrorNotFound {
		return
	}
	for key, data := range m {
		if key == channel {
			return bytes2participant(data)
		}
	}
	return
}
