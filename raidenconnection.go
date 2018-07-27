package smartraiden

import (
	"fmt"
	"sync"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/rerr"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

func (rs *RaidenService) connectionManagerForToken(tokenAddress common.Address) (*ConnectionManager, error) {
	mgr, ok := rs.Tokens2ConnectionManager[tokenAddress]
	if ok {
		return mgr, nil
	}
	return nil, rerr.InvalidAddress(fmt.Sprintf("token %s is not registered", utils.APex(tokenAddress)))
}
func (rs *RaidenService) leaveAllTokenNetworksAsync() *utils.AsyncResult {
	var leaveResults []*utils.AsyncResult
	for token := range rs.Token2ChannelGraph {
		mgr, err := rs.connectionManagerForToken(token)
		if err != nil {
			log.Error(fmt.Sprintf("connectionManagerForToken %s err %s", utils.APex(token), err))
		}
		if mgr != nil {
			leaveResults = append(leaveResults, mgr.LeaveAsync())
		}
	}
	return waitGroupAsyncResult(leaveResults)
}
func waitGroupAsyncResult(results []*utils.AsyncResult) *utils.AsyncResult {
	totalResult := utils.NewAsyncResult()
	wg := sync.WaitGroup{}
	wg.Add(len(results))
	for i := range results {
		go func(i int) {
			<-results[i].Result
			wg.Done()
		}(i)
	}
	go func() {
		wg.Wait()
		totalResult.Result <- nil
		close(totalResult.Result)
	}()
	return totalResult
}
func (rs *RaidenService) closeAndSettle() {
	log.Info("raiden will close and settle all channels now")
	var Mgrs []*ConnectionManager
	for t := range rs.Token2ChannelGraph {
		mgr, err := rs.connectionManagerForToken(t)
		if err != nil {
			continue
		}
		Mgrs = append(Mgrs, mgr)
	}
	blocksToWait := func() int64 {
		var max int64
		for _, mgr := range Mgrs {
			if max < mgr.minSettleBlocks() {
				max = mgr.minSettleBlocks()
			}
		}
		return max
	}
	var AllChannels []*channel.Channel
	for _, mgr := range Mgrs {
		for _, c := range mgr.openChannels() {
			ch, err := rs.findChannelByAddress(c.ChannelIdentifier.ChannelIdentifier)
			if err != nil {
				panic(fmt.Sprintf("channel %s must exist", c.ChannelIdentifier))
			}
			AllChannels = append(AllChannels, ch)
		}
	}
	leavingResult := rs.leaveAllTokenNetworksAsync()
	//using the un-cached block number here
	lastBlock := rs.GetBlockNumber()
	earliestSettlement := lastBlock + blocksToWait()
	/*
			    TODO: estimate and set a `timeout` parameter in seconds
		     based on connection_manager.min_settle_blocks and an average
		     blocktime from the past
	*/
	currentBlock := lastBlock
	for currentBlock < earliestSettlement {
		time.Sleep(time.Second * 10)
		lastBlock := rs.GetBlockNumber()
		if lastBlock != currentBlock {
			currentBlock = lastBlock
			waitBlocksLeft := blocksToWait()
			notSettled := 0
			for _, c := range AllChannels {
				if c.State != channeltype.StateSettled {
					notSettled++
				}
			}
			if notSettled == 0 {
				log.Debug("nothing left to settle")
				break
			}
			log.Info(fmt.Sprintf("waiting at least %d more blocks for %d channels not yet settled", waitBlocksLeft, notSettled))
		}
		//why  leaving_greenlet.wait
		timeoutch := time.After(time.Second * time.Duration(blocksToWait()))
		select {
		case <-timeoutch:
		case <-leavingResult.Result:
		}
	}
	for _, c := range AllChannels {
		if c.State != channeltype.StateSettled {
			log.Error("channels were not settled:", c.ChannelIdentifier)
		}
	}
}
