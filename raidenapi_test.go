package raiden_network

import (
	"testing"

	"sync"

	"time"

	"fmt"

	"math/big"

	"github.com/SmartMeshFoundation/raiden-network/channel"
	"github.com/SmartMeshFoundation/raiden-network/transfer"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fatedier/frp/src/utils/log"
)

var RepeatCount int = 100
var big1 = big.NewInt(1)

//a valid channel address onchain
func getAChannel(api *RaidenApi) common.Address {
	for _, g := range api.Raiden.Token2ChannelGraph {
		for addr, _ := range g.ChannelAddress2Channel {
			return addr
		}
	}
	panic("no channel")
}
func getAToken(api *RaidenApi) common.Address {
	for t, _ := range api.Raiden.Token2ChannelGraph {
		return t
	}
	panic("no token")
}
func TestSwapKeyAsMapKey(t *testing.T) {
	key1 := SwapKey{
		Identifier: 32,
		FromToken:  utils.NewRandomAddress(),
		FromAmount: big.NewInt(300),
	}
	key2 := key1
	m := make(map[SwapKey]bool)
	m[key1] = true
	if m[key2] != true {
		t.Error("expect equal")
	}
	key2.Identifier = 3
	if m[key2] == true {
		t.Error("should not equal")
	}
}

func testRaidenApi_GetNetworkEvents(t *testing.T, api *RaidenApi) {
	events, err := api.GetNetworkEvents(-1, -1)
	if err != nil {
		t.Error(err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(RepeatCount)
	for i := 0; i < RepeatCount; i++ {
		go func() {
			ev2, err := api.GetNetworkEvents(-1, -1)
			if err != nil {
				t.Error(err)
			}
			if len(ev2) != len(events) {
				t.Error("not equal")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func testRaidenApi_GetChannelEvents(t *testing.T, api *RaidenApi) {
	addr := getAChannel(api)
	events, err := api.GetChannelEvents(addr, -1, -1)
	if err != nil {
		t.Error(err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(RepeatCount)
	for i := 0; i < RepeatCount; i++ {
		go func() {
			ev2, err := api.GetChannelEvents(addr, -1, -1)
			if err != nil {
				t.Error(err)
			}
			if len(ev2) != len(events) {
				t.Error("not equal")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
func testGetTokenNetworkEvents(t *testing.T, api *RaidenApi) {
	addr := getAToken(api)
	events, err := api.GetTokenNetworkEvents(addr, -1, -1)
	if err != nil {
		t.Error(err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(RepeatCount)
	for i := 0; i < RepeatCount; i++ {
		go func() {
			ev2, err := api.GetTokenNetworkEvents(addr, -1, -1)
			if err != nil {
				t.Error(err)
			}
			if len(ev2) != len(events) {
				t.Error("not equal")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
func TestEvents(t *testing.T) {
	api := newTestRaidenApi()
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		testRaidenApi_GetChannelEvents(t, api)
		wg.Done()
	}()
	go func() {
		testRaidenApi_GetNetworkEvents(t, api)
		wg.Done()
	}()
	go func() {
		testGetTokenNetworkEvents(t, api)
		wg.Done()
	}()
	wg.Wait()
}

func testRaidenApi_Close(t *testing.T, api *RaidenApi, ch *channel.Channel) {
	tokenAddr := ch.TokenAddress
	partnerAddr := ch.PartnerState.Address
	wg := sync.WaitGroup{}
	wg.Add(RepeatCount)
	for i := 0; i < RepeatCount; i++ {
		go func() {
			api.Close(tokenAddr, partnerAddr)
			wg.Done()
		}()
	}
	wg.Wait()
}
func testRaidenApi_Settle(t *testing.T, api *RaidenApi, ch *channel.Channel) {
	tokenAddr := ch.TokenAddress
	partnerAddr := ch.PartnerState.Address
	wg := sync.WaitGroup{}
	wg.Add(RepeatCount)
	for i := 0; i < RepeatCount; i++ {
		go func() {
			//settle error ignore
			api.Settle(tokenAddr, partnerAddr)
			wg.Done()
		}()
	}
	wg.Wait()
}
func TestRaidenApi_CloseAndSettle(t *testing.T) {
	api := newTestRaidenApi()
	ch := api.Raiden.GetChannelWithAddr(getAChannel(api))
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		testRaidenApi_Close(t, api, ch)
		wg.Done()
	}()
	go func() {
		testRaidenApi_Settle(t, api, ch)
		wg.Done()
	}()
	wg.Wait()
}

func findAValidChannel(ra, rb *RaidenApi) (addr common.Address, money *big.Int) {
	for _, g := range ra.Raiden.CloneToken2ChannelGraph() {
		c := g.GetPartenerAddress2Channel(rb.Raiden.NodeAddress)
		if c != nil && c.Balance().Cmp(big.NewInt(10)) > 0 && c.State() == transfer.CHANNEL_STATE_OPENED {
			return c.MyAddress, c.Balance()
		}
	}
	return
}

/*
to test transfer concurrency
*/
func findAllCanTransferChannel(ra, rb, rc *RaidenApi) map[common.Address]common.Address {
	m := make(map[common.Address]common.Address)
	for _, g := range ra.Raiden.CloneToken2ChannelGraph() {
		for addr, c := range g.ChannelAddress2Channel {
			if c.Balance().Cmp(utils.BigInt0) > 0 && c.State() == transfer.CHANNEL_STATE_OPENED {
				if m[addr] == utils.EmptyAddress {
					m[addr] = ra.Raiden.NodeAddress
				}
			}
		}
	}
	for _, g := range rb.Raiden.CloneToken2ChannelGraph() {
		for addr, c := range g.ChannelAddress2Channel {
			if c.Balance().Cmp(utils.BigInt0) > 0 && c.State() == transfer.CHANNEL_STATE_OPENED {
				if m[addr] == utils.EmptyAddress {
					m[addr] = rb.Raiden.NodeAddress
				}
			}
		}
	}
	return m
}
func TestTransfer(t *testing.T) {
	ra, rb, rc, _ := makeTestRaidenApis()
	chm := findAllCanTransferChannel(ra, rb, rc)

	wgStart := sync.WaitGroup{}
	wgEnd := sync.WaitGroup{}
	log.Info("channels number ", len(chm))
	var i uint64 = 0
	values := make(map[*channel.Channel]*big.Int)
	for chaddr, nodeAddr := range chm {
		i++
		r := rb
		if nodeAddr == ra.Raiden.NodeAddress {
			r = ra
		}
		ch, _ := r.GetChannel(chaddr)
		values[ch] = ch.Balance()
		go func(r *RaidenApi, tokenAddr, partnerAddr common.Address, id uint64) {
			wgStart.Add(1)
			wgStart.Wait() //同时开始
			err := r.Transfer(tokenAddr, big1, partnerAddr, id, time.Minute*2)
			if err != nil {
				t.Error()
			}
			wgEnd.Done()
		}(r, ch.TokenAddress, ch.PartnerState.Address, i)
	}
	time.Sleep(time.Second)
	wgEnd.Add(int(i))
	for j := 0; j < int(i); j++ {
		wgStart.Done()
	}
	wgEnd.Wait()
	time.Sleep(time.Second * 30)
	for c, v := range values {
		if c.Balance().Cmp(x.Sub(v, big.NewInt(1))) != 0 {
			log.Error(fmt.Sprintf("transfer amount misatch expect %d,get %d @%s", x.Sub(v, big.NewInt(1)), c.Balance(), c.MyAddress.String()))
		}
	}
	//there are unfinished transfers?
	time.Sleep(time.Second)
}

//test must must fail,because transfer must ordered between two partners.
//but should not panic
func TestTransferWithPython(t *testing.T) {
	ra := newTestRaidenApi()
	log.Info("node addr:=", ra.Address().String())
	c, _ := ra.GetChannel(common.HexToAddress("0x82e7281fc9d42a66e195ed66b6718bb706c1af7c"))
	money := c.Balance()
	originalBalance := money
	wg := sync.WaitGroup{}
	//cnt := int(money) - 1
	cnt := 10
	wg.Add(cnt)
	for i := 1; i < cnt+1; i++ {
		go func(id int) {
			err := ra.Transfer(c.TokenAddress, big1, c.PartnerState.Address, uint64(id), time.Second*50)
			if err != nil {
				t.Error(err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	time.Sleep(time.Second * 3)
	assert(t, c.Balance(), x.Sub(originalBalance, big.NewInt(int64(cnt))))
}

func TestPairTransfer(t *testing.T) {
	ra, rb, _, _ := makeTestRaidenApis()
	addr, _ := findAValidChannel(ra, rb)
	c := ra.Raiden.GetChannelWithAddr(addr)
	amoney := ra.Raiden.GetChannelWithAddr(addr).Balance()
	bmoney := rb.Raiden.GetChannelWithAddr(addr).Balance()
	for i := 1; i < 100; i++ {
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			err := ra.Transfer(c.TokenAddress, big1, rb.Raiden.NodeAddress, uint64(2*i), time.Minute)
			if err != nil {
				t.Error(err)
			}
			wg.Done()
		}()
		go func() {
			err := rb.Transfer(c.TokenAddress, big1, ra.Raiden.NodeAddress, uint64(2*i+1), time.Minute)
			if err != nil {
				t.Error(err)
			}
			wg.Done()
		}()
		wg.Wait()
		time.Sleep(time.Second * 5) //let ra,rb update
		c1 := ra.Raiden.GetChannelWithAddr(addr)
		if c1.Balance().Cmp(amoney) != 0 {
			t.Error(fmt.Sprintf("money not equal expect=%d,get =%d", c1.Balance(), amoney))
			return
		}
		if c1.PartnerState.Balance(c1.OurState) != bmoney {
			t.Error(fmt.Sprintf("money not equalexpect=%d,get =%d", c1.PartnerState.Balance(c1.OurState), bmoney))
			return
		}
		c2 := rb.Raiden.GetChannelWithAddr(addr)
		if c2.Balance().Cmp(bmoney) != 0 {
			t.Error(fmt.Sprintf("money not equal expect=%d get=%d", c2.Balance(), bmoney))
			return
		}
		if c2.PartnerState.Balance(c2.OurState) != amoney {
			t.Error(fmt.Sprintf("money not equal expect=%d get=%d", c2.PartnerState.Balance(c2.OurState), amoney))
			return
		}
	}
}
