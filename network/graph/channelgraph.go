package graph

import (
	"errors"

	"fmt"

	"sort"

	"strings"

	"math/big"

	"github.com/SmartMeshFoundation/Photon/channel"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/network/dijkstra"
	"github.com/SmartMeshFoundation/Photon/network/rpc/fee"
	"github.com/SmartMeshFoundation/Photon/network/xmpptransport"
	"github.com/SmartMeshFoundation/Photon/transfer/route"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

//EmptyExlude 为了调用 GetBestRoutes 方便一点
// EmptyExclude : to make GetBestRoutes easy to invoke.
var EmptyExlude map[common.Address]bool

func init() {
	EmptyExlude = make(map[common.Address]bool)
}

//MakeExclude 为了调用 GetBestRoutes 方便一点
// MakeExclude : to make GetBestRoutes easy to invoke.
func MakeExclude(addrs ...common.Address) map[common.Address]bool {
	m := make(map[common.Address]bool)
	for _, addr := range addrs {
		m[addr] = true
	}
	return m
}

//NodesStatusGetter for route service
type NodesStatusGetter interface {
	//GetNetworkStatus return addr's status
	GetNetworkStatus(addr common.Address) (deviceType string, isOnline bool)
}

//ChannelGraph is a Graph based on the channels and can find path between participants.
//整个 ChannelGraph 只能单线程访问
// The whole ChannelGraph can only be accessed by a single process.
type ChannelGraph struct {
	g                         *dijkstra.Graph
	OurAddress                common.Address
	TokenAddress              common.Address
	PartenerAddress2Channel   map[common.Address]*channel.Channel
	ChannelIdentifier2Channel map[common.Hash]*channel.Channel
	address2index             map[common.Address]int
	index2address             map[int]common.Address
}

/*
NewChannelGraph create ChannelGraph,one token one channelGraph
*/
func NewChannelGraph(ourAddress, tokenAddress common.Address, edges []common.Address) *ChannelGraph {
	cg := &ChannelGraph{
		OurAddress:                ourAddress,
		TokenAddress:              tokenAddress,
		PartenerAddress2Channel:   make(map[common.Address]*channel.Channel),
		ChannelIdentifier2Channel: make(map[common.Hash]*channel.Channel),
		address2index:             make(map[common.Address]int),
		index2address:             make(map[int]common.Address),
		g:                         dijkstra.NewGraph(),
	}
	cg.makeGraph(edges)
	//cg.printGraph()
	return cg
}
func (cg *ChannelGraph) printGraph() {
	fmt.Printf("%s channel graph", utils.APex2(cg.TokenAddress))
	rowheader := fmt.Sprintf("%s", strings.Repeat(" ", 14))
	for i := 0; i < len(cg.index2address); i++ {
		rowheader += fmt.Sprintf("     %s:%2d", utils.APex2(cg.index2address[i]), i)
	}
	fmt.Println(rowheader)
	for i := 0; i < len(cg.index2address); i++ {
		fmt.Printf("       %s:%2d", utils.APex2(cg.index2address[i]), i)
		for j := 0; j < len(cg.index2address); j++ {
			v, err := cg.g.GetVertex(i)
			if err != nil {
				log.Crit(fmt.Sprintf("addr %s:%d not exist", utils.APex(cg.index2address[i]), i))
			}

			if _, ok := v.GetArc(j); ok {
				fmt.Printf("%12d", 1)
			} else {
				fmt.Printf("%12d", 0)
			}
		}
		fmt.Println("")
	}
}

func (cg *ChannelGraph) makeGraph(edges []common.Address) {
	for i := 0; i < len(edges); i += 2 {
		cg.AddPath(edges[i], edges[i+1])
	}
}

//AddPath Add a new edge into the network.
func (cg *ChannelGraph) AddPath(source, target common.Address) {
	addr1 := source
	addr2 := target
	if index1, ok := cg.address2index[addr1]; !ok {
		index1 = len(cg.address2index)
		cg.address2index[addr1] = index1
		cg.index2address[index1] = addr1
	}
	if index2, ok := cg.address2index[addr2]; !ok {
		index2 = len(cg.address2index)
		cg.address2index[addr2] = index2
		cg.index2address[index2] = addr2
	}
	var index1, index2 int
	index1 = cg.address2index[addr1]
	index2 = cg.address2index[addr2]
	_, err := cg.g.GetVertex(index1)
	if err != nil {
		cg.g.AddVertex(index1)
	}
	_, err = cg.g.GetVertex(index2)
	if err != nil {
		cg.g.AddVertex(index2)
	}
	//todo int64 cannot store too much tokens only about 18 tokens. we should divide 1000 or 1000000,
	err = cg.g.AddArc(index1, index2, 1)
	if err != nil {
		log.Error(fmt.Sprintf("add path err%s", err))
	}
	err = cg.g.AddArc(index2, index1, 1) //now our graph is the least fee first.
	if err != nil {
		log.Error(fmt.Sprintf("add path err%s", err))
	}
	//cg.printGraph()
}

/*
AddChannel Instantiate a channel this node participates and add to the graph.

        If the channel is already registered do nothing.
*/
func (cg *ChannelGraph) AddChannel(ch *channel.Channel) error {
	if ch.OurState.Address != cg.OurAddress {
		return errors.New("Address mismatch, our_address doesn't match the channel details")
	}
	if ch2 := cg.ChannelIdentifier2Channel[ch.ChannelIdentifier.ChannelIdentifier]; ch2 == nil {

		cg.PartenerAddress2Channel[ch.PartnerState.Address] = ch
		cg.ChannelIdentifier2Channel[ch.ChannelIdentifier.ChannelIdentifier] = ch
		cg.AddPath(ch.OurState.Address, ch.PartnerState.Address)
	} else {
		log.Info(fmt.Sprintf("add channel %s,but channel already exist", ch.ChannelIdentifier.String()))
	}
	return nil
}

/*
HasChannel return  True if there is a connecting path regardless of the number of hops.
*/
func (cg *ChannelGraph) HasChannel(source, target common.Address) bool {
	sourceIndex, ok := cg.address2index[source]
	if !ok {
		return false
	}
	targetIndex, ok := cg.address2index[target]
	if !ok {
		return false
	}
	_, err := cg.g.Shortest(sourceIndex, targetIndex)
	return err == nil
}

var errAddressNotFoundInGraph = errors.New("address not found in channelgraph")

/*
ShortestPath returns the shortestpath weight from source to target.  make sure only be called in one thread.
*/
func (cg *ChannelGraph) ShortestPath(source, target common.Address, amount *big.Int, feeCharger fee.Charger) (totalWeight int64, err error) {
	sourceIndex, ok := cg.address2index[source]
	if !ok {
		err = errAddressNotFoundInGraph
		return
	}
	targetIndex, ok := cg.address2index[target]
	if !ok {
		err = errAddressNotFoundInGraph
		return
	}
	if sourceIndex == targetIndex {
		return 0, nil
	}
	var g2 *dijkstra.Graph
	if false { //make sure only be called in one thread.
		g2 = cg.g.CloneGraph()
	} else {
		g2 = cg.g
	}
	for _, v := range g2.Verticies {
		w := feeCharger.GetNodeChargeFee(cg.index2address[v.ID], cg.TokenAddress, amount).Int64()
		if w > 0 { //for no fee policy, all nodes charge 0 ,so use the shortest path first.
			v.SetWeight(w) // from v's fee is w.
		}
	}
	path, err := g2.Shortest(sourceIndex, targetIndex)
	if err != nil {
		return
	}
	return path.Distance, nil
}

//RemoveChannel remove a channel from graph,and i'm a participant of this channel
func (cg *ChannelGraph) RemoveChannel(ch *channel.Channel) {
	delete(cg.ChannelIdentifier2Channel, ch.ChannelIdentifier.ChannelIdentifier)
	delete(cg.PartenerAddress2Channel, ch.PartnerState.Address)
	cg.RemovePath(ch.OurState.Address, ch.PartnerState.Address)
}

//RemovePath Remove an edge from the network.  this edge may  not exist
func (cg *ChannelGraph) RemovePath(source, target common.Address) {
	sourceIndex, ok := cg.address2index[source]
	if !ok {
		return
	}
	targetIndex, ok := cg.address2index[target]
	if !ok {
		return
	}
	err := cg.g.DeleteArc(sourceIndex, targetIndex)
	if err != nil {
		log.Error(fmt.Sprintf("remove arc %d-%d err %s", sourceIndex, targetIndex, err))
	}
	err = cg.g.DeleteArc(targetIndex, sourceIndex)
	if err != nil {
		log.Error(fmt.Sprintf("remove arc %d-%d err %s", sourceIndex, targetIndex, err))
	}
}

/*
ChannelCanTransfer returns  True if the channel with `partner_address` is open and has spendable funds. """
        TODO: check if the partner's network is alive
*/
func (cg *ChannelGraph) ChannelCanTransfer(partenerAddress common.Address) bool {
	return cg.GetPartenerAddress2Channel(partenerAddress).CanTransfer()
}

//getNeighbours Get all neighbours adjacent to self.our_address. g is not thread safe
func (cg *ChannelGraph) getNeighbours() []common.Address {
	neighboursIndex, err := cg.g.GetAllNeighbors(cg.address2index[cg.OurAddress])
	if err != nil {
		return nil
	}
	var neighbours []common.Address
	for _, i := range neighboursIndex {
		neighbours = append(neighbours, cg.index2address[i])
	}
	return neighbours
}

type neighborWeight struct {
	neighbor common.Address
	weight   int64 //nerghbor to target's hops
}
type neighborWeightList []*neighborWeight

func (nw neighborWeightList) Len() int {
	return len(nw)
}
func (nw neighborWeightList) Less(i, j int) bool {
	return nw[i].weight < nw[j].weight
}
func (nw neighborWeightList) Swap(i, j int) {
	var temp *neighborWeight
	temp = nw[i]
	nw[i] = nw[j]
	nw[j] = temp
}

/*
all the neighbors that can reach target
they are ordered by hops to the target
*/
func (cg *ChannelGraph) orderedNeighbours(ourAddress, targetAddress common.Address, amount *big.Int, charger fee.Charger) neighborWeightList {

	neighbors := cg.getNeighbours()
	var nws neighborWeightList
	for _, n := range neighbors {
		w, err := cg.ShortestPath(n, targetAddress, amount, charger)
		if err != nil {
			continue
		}
		nws = append(nws, &neighborWeight{n, w})
	}
	sort.Sort(nws)
	return nws
}

/*
GetBestRoutes returns all neighbor nodes order by weight from it to target.
我们现在的路由算法应该是有历史记忆的最短路径/最小费用算法.
跳过所有已经走过的路径.
*/
/*
 *	GetBestRoutes :function to return all neighbor nodes order by weight from it to target.
 *
 *	Note that the routing algorithm we currently use should be the shortest-path/minimized-fee algorithm with history record,
 *	which circumvents all routes that have been iterated.
 */
func (cg *ChannelGraph) GetBestRoutes(nodesStatus NodesStatusGetter, ourAddress common.Address,
	targetAdress common.Address, amount *big.Int, targetAmount *big.Int, excludeAddresses map[common.Address]bool, feeCharger fee.Charger) (onlineNodes []*route.State) {
	/*

	   XXX: consider using multiple channels for a single transfer. Useful
	   for cases were the `amount` is larger than what is available
	   individually in any of the channels.

	   One possible approach is to _not_ filter these channels based on the
	   distributable amount, but to sort them based on available balance and
	   let the task use as many as required to finish the transfer.

	*/
	nws := cg.orderedNeighbours(ourAddress, targetAdress, amount, feeCharger)
	if len(nws) == 0 {
		log.Warn(fmt.Sprintf("no routes avaiable from %s to %s", utils.APex(ourAddress), utils.APex(targetAdress)))
		return
	}
	for _, nw := range nws {
		c := cg.GetPartenerAddress2Channel(nw.neighbor)
		//don't send the message backwards
		if excludeAddresses[nw.neighbor] {
			continue
		}
		if !c.CanTransfer() {
			log.Debug(fmt.Sprintf("channel %s-%s cannot transfer ,ignoring ..", utils.APex(ourAddress), utils.APex(nw.neighbor)))
			continue
		}
		if amount.Cmp(c.Distributable()) > 0 {
			log.Debug(fmt.Sprintf("channel %s-%s doesn't have enough funds[%d],ignoring...", utils.APex(ourAddress), utils.APex(nw.neighbor), amount))
			continue
		}
		deviceType, isOnline := nodesStatus.GetNetworkStatus(nw.neighbor)
		if !isOnline || (deviceType == xmpptransport.TypeMobile && nw.neighbor != targetAdress) {
			log.Debug(fmt.Sprintf("partener %s network ignored.. isOnline:%v,deviceType:%s", utils.APex(nw.neighbor), isOnline, deviceType))
			continue
		}
		routeState := Channel2RouteState(c, nw.neighbor, targetAmount, feeCharger)
		if routeState.Fee.Cmp(utils.BigInt0) > 0 {
			routeState.TotalFee = big.NewInt(int64(nw.weight))
		} else { //no fee policy,
			routeState.TotalFee = utils.BigInt0
		}

		onlineNodes = append(onlineNodes, routeState)
	}
	return
}
func (cg *ChannelGraph) haveNodes() bool {
	return len(cg.g.Verticies) > 0
}

//AllNodes returns all neighbor nodes
func (cg *ChannelGraph) AllNodes() (nodes []common.Address) {
	for n := range cg.address2index {
		nodes = append(nodes, n)
	}
	return nodes
}

//GetPartenerAddress2Channel returns a channel between me and address
func (cg *ChannelGraph) GetPartenerAddress2Channel(address common.Address) (c *channel.Channel) {
	c = cg.PartenerAddress2Channel[address]
	if c == nil {
		log.Trace(fmt.Sprintf("no channel with %s on token %s", utils.APex(address), utils.APex(cg.TokenAddress)))
	}
	return
}

//Channel2RouteState create a routeState from a channel
func Channel2RouteState(c *channel.Channel, partenerAddress common.Address, amount *big.Int, charger fee.Charger) *route.State {
	rs := route.NewState(c)
	rs.Fee = charger.GetNodeChargeFee(partenerAddress, c.TokenAddress, amount)
	return rs
}
