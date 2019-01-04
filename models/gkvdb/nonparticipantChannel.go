package gkvdb

import (
	"fmt"

	"github.com/kataras/go-errors"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

type nonParticipantChannel struct {
	ChannelIdentifierBytes []byte
	TokenAddressBytes      []byte
	Participant1Bytes      []byte
	Participant2Bytes      []byte
}

//NewNonParticipantChannel 需要保存 channel identifier, 通道的事件都是与此有关系的
func (dao *GkvDB) NewNonParticipantChannel(token common.Address, channel common.Hash, participant1, participant2 common.Address) error {
	var m nonParticipantChannel
	log.Trace(fmt.Sprintf("NewNonParticipantChannel token=%s,participant1=%s,participant2=%s",
		utils.APex2(token),
		utils.APex2(participant1),
		utils.APex2(participant2),
	))
	err := dao.getKeyValueToBucket(models.BucketChannel, channel[:], &m)
	if err == nil {
		return errors.New("duplicate key")
	}
	if participant1 == participant2 {
		panic(fmt.Sprintf("channel error, p1 andf p2 is the same,token=%s,participant=%s", token.String(), participant1.String()))
	}

	log.Trace(fmt.Sprintf("NewNonParticipantChannel token=%s,p1=%s,p2=%s", utils.APex2(token),
		utils.APex2(participant1), utils.APex2(participant2)))
	return dao.saveKeyValueToBucket(models.BucketChannel, channel[:], &nonParticipantChannel{
		TokenAddressBytes:      token[:],
		Participant1Bytes:      participant1[:],
		Participant2Bytes:      participant2[:],
		ChannelIdentifierBytes: channel[:],
	})
}

//RemoveNonParticipantChannel a channel is settled
func (dao *GkvDB) RemoveNonParticipantChannel(channel common.Hash) error {
	var m nonParticipantChannel
	err := dao.getKeyValueToBucket(models.BucketChannel, channel[:], &m)
	if err != nil {
		return err
	}
	return dao.removeKeyValueFromBucket(models.BucketChannel, channel[:])
}

//GetAllNonParticipantChannelByToken returna all channel on this `token`
func (dao *GkvDB) GetAllNonParticipantChannelByToken(token common.Address) (edges []common.Address, err error) {
	tb, err := dao.db.Table(models.BucketChannel)
	if err != nil {
		panic(err)
	}
	buf := tb.Values(-1)
	if buf == nil || len(buf) == 0 {
		return
	}
	for _, v := range buf {
		var m nonParticipantChannel
		gobDecode(v, &m)
		edges = append(edges, common.BytesToAddress(m.Participant1Bytes), common.BytesToAddress(m.Participant2Bytes))
	}
	return
}

// GetNonParticipantChannelByID  :
func (dao *GkvDB) GetNonParticipantChannelByID(channelIdentifierForQuery common.Hash) (
	tokenAddress common.Address, participant1, participant2 common.Address, err error) {
	var m nonParticipantChannel
	err = dao.getKeyValueToBucket(models.BucketChannel, channelIdentifierForQuery[:], &m)
	log.Trace(fmt.Sprintf("GetNonParticipantChannelByID,channel=%s,err=%v", utils.HPex(channelIdentifierForQuery), err))
	if err == ErrorNotFound {
		return
	}
	tokenAddress = common.BytesToAddress(m.TokenAddressBytes)
	participant1 = common.BytesToAddress(m.Participant1Bytes)
	participant2 = common.BytesToAddress(m.Participant2Bytes)
	return
}
