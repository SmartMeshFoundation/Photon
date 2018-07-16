package models

import (
	"fmt"
	"testing"

	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/stretchr/testify/assert"
)

func TestChannel(t *testing.T) {
	model := setupDb(t)
	defer func() {
		model.CloseDB()
	}()
	newchannelcb := func(c *channeltype.Serialization) bool {

		return true
	}
	updateContractBalancechannelcb := func(c *channeltype.Serialization) bool {

		return true
	}
	UpdateChannelStatecb := func(c *channeltype.Serialization) bool {

		return true
	}
	model.RegisterNewChannellCallback(newchannelcb)
	model.RegisterChannelDepositCallback(updateContractBalancechannelcb)
	model.RegisterChannelStateCallback(UpdateChannelStatecb)

	ch1 := &channeltype.Serialization{
		ChannelIdentifier: &contracts.ChannelUniqueID{
			ChannelIdentifier: utils.NewRandomHash(),
			OpenBlockNumber:   3,
		},
		Key:            utils.NewRandomHash(),
		TokenAddress:   utils.NewRandomAddress(),
		PartnerAddress: utils.NewRandomAddress(),
	}
	c := ch1
	ch2 := &channeltype.Serialization{
		ChannelIdentifier: &contracts.ChannelUniqueID{
			ChannelIdentifier: utils.NewRandomHash(),
			OpenBlockNumber:   3,
		},
		Key:            utils.NewRandomHash(),
		TokenAddress:   utils.NewRandomAddress(),
		PartnerAddress: utils.NewRandomAddress(),
	}
	err := model.NewChannel(ch1)
	if err != nil {
		t.Error(err)
		return
	}
	err = model.NewChannel(ch2)
	if err != nil {
		t.Error(err)
		return
	}

	chs, err := model.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	if err != nil || len(chs) != 2 {
		t.Error(err)
		t.Log(fmt.Sprintf("chs=%v", utils.StringInterface(chs, 5)))
		return
	}
	//log.Trace(fmt.Sprintf("ch1=%s,ch2=%s", utils.StringInterface(ch1, 3), utils.StringInterface(ch2, 3)))
	ch, err := model.GetChannel(c.TokenAddress, c.PartnerAddress)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, c, ch)
	err = model.UpdateChannelNoTx(c)
	if err != nil {
		t.Error(err)
		return
	}
	chs, err = model.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	if err != nil || len(chs) != 2 {
		t.Error(err)
		log.Error(fmt.Sprintf("chs=%s", utils.StringInterface(chs, 3)))
		return
	}
	//log.Trace(fmt.Sprintf("chs=%s", utils.StringInterface(chs, 3)))
	err = model.UpdateChannelContractBalance(c)
	if err != nil {
		t.Error(err)
		return
	}
	err = model.UpdateChannelContractBalance(c)
	if err != nil {
		t.Error(err)
		return
	}
	err = model.UpdateChannelState(c)
	if err != nil {
		t.Error(err)
		return
	}
	err = model.UpdateChannelState(c)
	if err != nil {
		t.Error(err)
		return
	}
}
func TestChannelTwice(t *testing.T) {
	TestChannel(t)
	TestChannel(t)
}
