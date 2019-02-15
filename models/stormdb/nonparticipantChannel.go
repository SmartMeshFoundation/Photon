package stormdb

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/models"

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
	ChannelIdentifierBytes []byte `storm:"id"`
	TokenAddressBytes      []byte `storm:"index"`
	Participant1Bytes      []byte
	Participant2Bytes      []byte
}

//NewNonParticipantChannel 需要保存 channel identifier, 通道的事件都是与此有关系的
func (model *StormDB) NewNonParticipantChannel(token common.Address, channel common.Hash, participant1, participant2 common.Address) error {
	if participant1 == participant2 {
		return fmt.Errorf("channel error, p1 andf p2 is the same,token=%s,participant=%s", token.String(), participant1.String())
	}
	_, _, _, err := model.GetNonParticipantChannelByID(channel)
	if err == nil {
		return fmt.Errorf("channel %s already exists", channel.String())
	}
	err = model.db.Save(&NonParticipantChannel{
		ChannelIdentifierBytes: channel[:],
		TokenAddressBytes:      token[:],
		Participant1Bytes:      participant1[:],
		Participant2Bytes:      participant2[:],
	})
	return models.GeneratDBError(err)
}

//RemoveNonParticipantChannel a channel is settled
func (model *StormDB) RemoveNonParticipantChannel(channel common.Hash) error {
	err := model.db.DeleteStruct(&NonParticipantChannel{
		ChannelIdentifierBytes: channel[:],
	})
	return models.GeneratDBError(err)
}

//GetNonParticipantChannelByID return one channel's information
func (model *StormDB) GetNonParticipantChannelByID(channelIdentifierForQuery common.Hash) (
	tokenAddress common.Address, participant1, participant2 common.Address, err error) {
	var channel NonParticipantChannel
	err = model.db.One("ChannelIdentifierBytes", channelIdentifierForQuery[:], &channel)
	if err != nil {
		err = fmt.Errorf("GetNonParticipantChannelByID err %s", err)
		return
	}
	tokenAddress = common.BytesToAddress(channel.TokenAddressBytes)
	participant1 = common.BytesToAddress(channel.Participant1Bytes)
	participant1 = common.BytesToAddress(channel.Participant2Bytes)
	return
}

//GetAllNonParticipantChannelByToken returna all channel on this `token`
func (model *StormDB) GetAllNonParticipantChannelByToken(token common.Address) (edges []common.Address, err error) {
	var channels []*NonParticipantChannel
	err = model.db.Find("TokenAddressBytes", token[:], &channels)
	if err == storm.ErrNotFound {
		err = nil
		return
	}
	if err != nil {
		err = fmt.Errorf("GetAllNonParticipantChannelByToken err %s", err)
		err = models.GeneratDBError(err)
		return
	}
	for _, c := range channels {
		edges = append(edges, common.BytesToAddress(c.Participant1Bytes), common.BytesToAddress(c.Participant2Bytes))
	}
	return
}
