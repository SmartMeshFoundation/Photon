package smartraiden

import (
	"testing"

	"sync"

	"time"

	"fmt"

	"math/big"

	"os"

	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fatedier/frp/src/utils/log"
)

var repeatCount = 10
var big1 = big.NewInt(1)

//a valid channel address onchain
func getAChannel(api *RaidenAPI) common.Address {
	for _, g := range api.Raiden.Token2ChannelGraph {
		for addr := range g.ChannelAddress2Channel {
			return addr
		}
	}
	panic("no channel")
}
func getAToken(api *RaidenAPI) common.Address {
	for t := range api.Raiden.Token2ChannelGraph {
		return t
	}
	panic("no token")
}
func TestSwapKeyAsMapKey(t *testing.T) {
	reinit()
	key1 := swapKey{
		LockSecretHash: 32,
		FromToken:      utils.NewRandomAddress(),
		FromAmount:     big.NewInt(300).String(),
	}
	key2 := key1
	m := make(map[swapKey]bool)
	m[key1] = true
	if m[key2] != true {
		t.Error("expect equal")
	}
	key2.LockSecretHash = 3
	if m[key2] == true {
		t.Error("should not equal")
	}
}

func testRaidenAPIGetNetworkEvents(t *testing.T, api *RaidenAPI) {
	events, err := api.GetNetworkEvents(-1, -1)
	if err != nil {
		t.Error(err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(repeatCount)
	for i := 0; i < repeatCount; i++ {
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

func testRaidenAPIGetChannelEvents(t *testing.T, api *RaidenAPI) {
	addr := getAChannel(api)
	events, err := api.GetChannelEvents(addr, -1, -1)
	if err != nil {
		t.Error(err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(repeatCount)
	for i := 0; i < repeatCount; i++ {
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
func testGetTokenNetworkEvents(t *testing.T, api *RaidenAPI) {
	addr := getAToken(api)
	events, err := api.GetTokenNetworkEvents(addr, -1, -1)
	if err != nil {
		t.Error(err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(repeatCount)
	for i := 0; i < repeatCount; i++ {
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
	reinit()
	api := newTestRaidenAPI()
	defer api.Stop()
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		testRaidenAPIGetChannelEvents(t, api)
		wg.Done()
	}()
	go func() {
		testRaidenAPIGetNetworkEvents(t, api)
		wg.Done()
	}()
	go func() {
		testGetTokenNetworkEvents(t, api)
		wg.Done()
	}()
	wg.Wait()
}

func testRaidenAPIClose(t *testing.T, api *RaidenAPI, ch *channel.Channel) {
	tokenAddr := ch.TokenAddress
	partnerAddr := ch.PartnerState.Address
	wg := sync.WaitGroup{}
	wg.Add(repeatCount)
	for i := 0; i < repeatCount; i++ {
		go func() {
			api.Close(tokenAddr, partnerAddr)
			wg.Done()
		}()
	}
	wg.Wait()
}
func raidenAPISettle(t *testing.T, api *RaidenAPI, ch *channel.Channel) {
	tokenAddr := ch.TokenAddress
	partnerAddr := ch.PartnerState.Address
	wg := sync.WaitGroup{}
	wg.Add(repeatCount)
	for i := 0; i < repeatCount; i++ {
		go func() {
			//settle error ignore
			api.Settle(tokenAddr, partnerAddr)
			wg.Done()
		}()
	}
	wg.Wait()
}
func TestRaidenAPICloseAndSettle(t *testing.T) {
	reinit()
	if true {
		return
	}
	api := newTestRaidenAPI()
	defer api.Stop()
	ch := api.Raiden.getChannelWithAddr(getAChannel(api))
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		testRaidenAPIClose(t, api, ch)
		wg.Done()
	}()
	go func() {
		raidenAPISettle(t, api, ch)
		wg.Done()
	}()
	wg.Wait()
}

func findAValidChannel(ra, rb *RaidenAPI) (addr common.Address, money *big.Int) {
	for _, g := range ra.Raiden.Token2ChannelGraph {
		c := g.GetPartenerAddress2Channel(rb.Raiden.NodeAddress)
		if c != nil && c.Balance().Cmp(big.NewInt(10)) > 0 && c.State() == transfer.ChannelStateOpened {
			return c.ChannelIdentifier, c.Balance()
		}
	}
	return
}

/*
to test transfer concurrency, todo, return only channel between a,b,c
*/
func findAllCanTransferChannel(ra, rb, rc *RaidenAPI) map[common.Address]common.Address {
	var allAddresses = map[common.Address]bool{
		ra.Raiden.NodeAddress: true,
		rb.Raiden.NodeAddress: true,
		rc.Raiden.NodeAddress: true,
	}
	m := make(map[common.Address]common.Address)
	for _, g := range ra.Raiden.Token2ChannelGraph {
		for addr, c := range g.ChannelAddress2Channel {
			if c.Balance().Cmp(utils.BigInt0) > 0 && c.State() == transfer.ChannelStateOpened && allAddresses[c.PartnerState.Address] {
				if m[addr] == utils.EmptyAddress {
					m[addr] = ra.Raiden.NodeAddress
				}
			}
		}
	}
	for _, g := range rb.Raiden.Token2ChannelGraph {
		for addr, c := range g.ChannelAddress2Channel {
			if c.Balance().Cmp(utils.BigInt0) > 0 && c.State() == transfer.ChannelStateOpened && allAddresses[c.PartnerState.Address] {
				if m[addr] == utils.EmptyAddress {
					m[addr] = rb.Raiden.NodeAddress
				}
			}
		}
	}
	return m
}
func TestTransfer(t *testing.T) {
	reinit()
	if true {
		return
	}
	ra, rb, rc, rd := makeTestRaidenAPIs()
	defer ra.Stop()
	defer rb.Stop()
	defer rc.Stop()
	defer rd.Stop()
	chm := findAllCanTransferChannel(ra, rb, rc)

	wgStart := sync.WaitGroup{}
	wgEnd := sync.WaitGroup{}
	log.Info("channels number ", len(chm))
	var i uint64
	values := make(map[*RaidenAPI]map[common.Address]*big.Int)
	for chaddr, nodeAddr := range chm {
		i++
		r := rb
		if nodeAddr == ra.Raiden.NodeAddress {
			r = ra
		}
		ch, err := r.GetChannel(chaddr)
		if err != nil {
			t.Error(err)
		}
		_, ok := values[r]
		if !ok {
			values[r] = make(map[common.Address]*big.Int)
		}
		values[r][chaddr] = ch.OurBalance
		go func(r *RaidenAPI, tokenAddr, partnerAddr common.Address, id uint64) {
			wgStart.Add(1)
			wgStart.Wait() //start at the same time
			err := r.Transfer(tokenAddr, big1, utils.BigInt0, partnerAddr, id, time.Minute*2, false)
			if err != nil {
				t.Error()
			}
			wgEnd.Done()
		}(r, ch.TokenAddress, ch.PartnerAddress, i)
	}
	time.Sleep(time.Second)
	wgEnd.Add(int(i))
	for j := 0; j < int(i); j++ {
		wgStart.Done()
	}
	wgEnd.Wait()
	time.Sleep(time.Second * 3)
	for r, m := range values {
		for addr, v := range m {
			c, err := r.GetChannel(addr)
			if err != nil {
				t.Error(err)
			}
			if c.OurBalance.Cmp(x.Sub(v, big.NewInt(1))) != 0 {
				log.Error(fmt.Sprintf("transfer amount misatch expect %d,get %d @%s", x.Sub(v, big.NewInt(1)), c.OurBalance, c.OurAddress.String()))
			}
		}
	}
	//there are unfinished transfers?
	time.Sleep(time.Second)
}

//test must must fail,because transfer must ordered between two partners.
//but should not panic
func TestTransferWithPython(t *testing.T) {
	reinit()
	ra := newTestRaidenAPI()
	defer ra.Stop()
	log.Info("node addr:=", ra.Address().String())
	c, _ := ra.GetChannel(common.HexToAddress(os.Getenv("CHANNEL")))
	wg := sync.WaitGroup{}
	//cnt := int(money) - 1
	cnt := 10
	wg.Add(cnt)
	for i := 1; i < cnt+1; i++ {
		go func(id int) {
			err := ra.Transfer(c.TokenAddress, big1, utils.BigInt0, c.PartnerAddress, uint64(id), time.Second*10, false)
			if err != nil {
				log.Error(fmt.Sprintf("err=%s", err))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	time.Sleep(time.Second * 3)
	//assert(t, c.OurBalance, x.Sub(originalBalance, big.NewInt(int64(cnt))))
}

func TestPairTransfer(t *testing.T) {
	reinit()
	ra, rb, rc, rd := makeTestRaidenAPIs()
	rc.Stop()
	rd.Stop()
	defer ra.Stop()
	defer rb.Stop()
	log.Info("nodes startup complete...")
	addr, _ := findAValidChannel(ra, rb)
	if addr == utils.EmptyAddress {
		t.Logf("no channel to transfer..\n")
		return
	}
	c := ra.Raiden.getChannelWithAddr(addr)
	amoney := ra.Raiden.getChannelWithAddr(addr).Balance()
	bmoney := rb.Raiden.getChannelWithAddr(addr).Balance()
	wg := sync.WaitGroup{}
	cnt := 5
	fmt.Printf("start transfer...\n")
	wg.Add(cnt * 2)
	for i := 1; i <= cnt; i++ {
		//wg.Add(2)
		go func(index int) {
			err := ra.Transfer(c.TokenAddress, big1, utils.BigInt0, rb.Raiden.NodeAddress, uint64(2*index), time.Minute*20, false)
			if err != nil {
				t.Error(err)
			}
			wg.Done()
		}(i)
		go func(index int) {
			err := rb.Transfer(c.TokenAddress, big1, utils.BigInt0, ra.Raiden.NodeAddress, uint64(2*index+1), time.Minute*20, false)
			if err != nil {
				t.Error(err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Print("end transfer...\n")
	time.Sleep(time.Second * 18) //let ra,rb update
	c1 := ra.Raiden.getChannelWithAddr(addr)
	if c1.Balance().Cmp(amoney) != 0 {
		t.Errorf("money not equal expect=%d,get =%d", c1.Balance(), amoney)
		//return
	}
	if c1.PartnerState.Balance(c1.OurState).Cmp(bmoney) != 0 {
		t.Errorf("money not equalexpect=%d,get =%d", c1.PartnerState.Balance(c1.OurState), bmoney)
		//return
	}
	c2 := rb.Raiden.getChannelWithAddr(addr)
	if c2.Balance().Cmp(bmoney) != 0 {
		t.Errorf("money not equal expect=%d get=%d", c2.Balance(), bmoney)
		//return
	}
	if c2.PartnerState.Balance(c2.OurState).Cmp(amoney) != 0 {
		t.Errorf("money not equal expect=%d get=%d", c2.PartnerState.Balance(c2.OurState), amoney)
		//return
	}
}
