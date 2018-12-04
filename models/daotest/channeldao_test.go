package daotest

import (
	"fmt"
	"testing"

	"time"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/codefortest"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestChannel1(t *testing.T) {
	dao := codefortest.NewTestDB("")
	defer dao.CloseDB()
	newchannelcb := func(c *channeltype.Serialization) bool {

		return true
	}
	updateContractBalancechannelcb := func(c *channeltype.Serialization) bool {

		return true
	}
	UpdateChannelStatecb := func(c *channeltype.Serialization) bool {

		return true
	}
	dao.RegisterNewChannelCallback(newchannelcb)
	dao.RegisterChannelDepositCallback(updateContractBalancechannelcb)
	dao.RegisterChannelStateCallback(UpdateChannelStatecb)
	h := utils.NewRandomHash()
	a1 := utils.NewRandomAddress()
	a2 := utils.NewRandomAddress()
	ch1 := &channeltype.Serialization{
		ChannelIdentifier: &contracts.ChannelUniqueID{
			ChannelIdentifier: h,
			OpenBlockNumber:   3,
		},
		Key:                 h[:],
		TokenAddressBytes:   a1[:],
		PartnerAddressBytes: a2[:],
	}
	c := ch1
	h2 := utils.NewRandomHash()
	a21 := utils.NewRandomAddress()
	a22 := utils.NewRandomAddress()
	ch2 := &channeltype.Serialization{
		ChannelIdentifier: &contracts.ChannelUniqueID{
			ChannelIdentifier: h2,
			OpenBlockNumber:   3,
		},
		Key:                 h2[:],
		TokenAddressBytes:   a21[:],
		PartnerAddressBytes: a22[:],
	}
	err := dao.NewChannel(ch1)
	if err != nil {
		t.Error(err)
		return
	}
	err = dao.NewChannel(ch2)
	if err != nil {
		t.Error(err)
		return
	}

	chs, err := dao.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	if err != nil || len(chs) != 2 {
		t.Error(err)
		t.Log(fmt.Sprintf("chs=%v", utils.StringInterface(chs, 2)))
		return
	}
	//log.Trace(fmt.Sprintf("ch1=%s,ch2=%s", utils.StringInterface(ch1, 3), utils.StringInterface(ch2, 3)))
	ch, err := dao.GetChannel(c.TokenAddress(), c.PartnerAddress())
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, c, ch)
	ch, err = dao.GetChannelByAddress(common.BytesToHash(c.Key))
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, c, ch)
	err = dao.UpdateChannelNoTx(c)
	if err != nil {
		t.Error(err)
		return
	}
	chs2, err := dao.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	assert.EqualValues(t, err == nil, true)
	assert.EqualValues(t, len(chs), len(chs2))
	err = dao.UpdateChannelContractBalance(c)
	if err != nil {
		t.Error(err)
		return
	}
	err = dao.UpdateChannelContractBalance(c)
	if err != nil {
		t.Error(err)
		return
	}
	err = dao.UpdateChannelState(c)
	if err != nil {
		t.Error(err)
		return
	}
	err = dao.UpdateChannelState(c)
	if err != nil {
		t.Error(err)
		return
	}
	c.State = channeltype.StateSettled
	err = dao.RemoveChannel(c)
	if err != nil {
		t.Error(err)
		return
	}
	chs, err = dao.GetChannelList(utils.EmptyAddress, utils.EmptyAddress)
	if err != nil || len(chs) != 1 {
		t.Error(err)
		t.Log(fmt.Sprintf("chs=%v", utils.StringInterface(chs, 2)))
		return
	}
}

func TestChannelTwice(t *testing.T) {
	TestChannel1(t)
	time.Sleep(time.Second)
	TestChannel1(t)
}

func TestModelDB_NewSettledChannel(t *testing.T) {
	dao := codefortest.NewTestDB("")
	defer dao.CloseDB()
	h := utils.NewRandomHash()
	a1 := utils.NewRandomAddress()
	a2 := utils.NewRandomAddress()
	ch1 := &channeltype.Serialization{
		ChannelIdentifier: &contracts.ChannelUniqueID{
			ChannelIdentifier: h,
			OpenBlockNumber:   3,
		},
		Key:                 h[:],
		TokenAddressBytes:   a1[:],
		PartnerAddressBytes: a2[:],
		State:               channeltype.StateSettled,
	}
	h2 := utils.NewRandomHash()
	a21 := utils.NewRandomAddress()
	a22 := utils.NewRandomAddress()
	ch2 := &channeltype.Serialization{
		ChannelIdentifier: &contracts.ChannelUniqueID{
			ChannelIdentifier: h2,
			OpenBlockNumber:   3,
		},
		Key:                 h2[:],
		TokenAddressBytes:   a21[:],
		PartnerAddressBytes: a22[:],
		State:               channeltype.StateSettled,
	}
	err := dao.NewSettledChannel(ch1)
	if err != nil {
		t.Error(err)
		return
	}
	ch, err := dao.GetSettledChannel(ch1.ChannelIdentifier.ChannelIdentifier, ch1.ChannelIdentifier.OpenBlockNumber)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, ch, ch1)
	ch, err = dao.GetSettledChannel(ch1.ChannelIdentifier.ChannelIdentifier, 32)
	assert.EqualValues(t, err != nil, true)
	chs, err := dao.GetAllSettledChannel()
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, chs[0], ch1)
	err = dao.NewSettledChannel(ch2)
	if err != nil {
		t.Error(err)
		return
	}
	chs, err = dao.GetAllSettledChannel()
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, len(chs), 2)
}
