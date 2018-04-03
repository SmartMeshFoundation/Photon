package network

import (
	"errors"

	"fmt"

	"sort"

	"strings"

	"sync"

	"math/big"

	"github.com/SmartMeshFoundation/raiden-network/channel"
	"github.com/SmartMeshFoundation/raiden-network/network/dijkstra"
	"github.com/SmartMeshFoundation/raiden-network/network/rpc"
	"github.com/SmartMeshFoundation/raiden-network/network/rpc/fee"
	"github.com/SmartMeshFoundation/raiden-network/transfer"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type ChannelDetails struct {
	ChannelAddress    common.Address
	OurState          *channel.ChannelEndState
	PartenerState     *channel.ChannelEndState
	ExternState       *channel.ChannelExternalState
	BlockChainService *rpc.BlockChainService
	RevealTimeout     int
	SettleTimeout     int
}

//Has Graph based on the channels and can find path between participants.
type ChannelGraph struct {
	g                       *dijkstra.Graph
	OurAddress              common.Address
	TokenAddress            common.Address
	ChannelManagerAddress   common.Address
	PartenerAddress2Channel map[common.Address]*channel.Channel //多线程
	ChannelAddress2Channel  map[common.Address]*channel.Channel //多线程
	address2index           map[common.Address]int
	index2address           map[int]common.Address
	Lock                    sync.Mutex
}

func NewChannelGraph(ourAddress, channelManagerAddress, tokenAddress common.Address, edgeList []common.Address, channelDetails []*ChannelDetails) *ChannelGraph {
	cg := &ChannelGraph{
		OurAddress:              ourAddress,
		TokenAddress:            tokenAddress,
		ChannelManagerAddress:   channelManagerAddress,
		PartenerAddress2Channel: make(map[common.Address]*channel.Channel),
		ChannelAddress2Channel:  make(map[common.Address]*channel.Channel),
		address2index:           make(map[common.Address]int),
		index2address:           make(map[int]common.Address),
		g:                       dijkstra.NewGraph(),
	}
	cg.makeGraph(edgeList)
	for _, d := range channelDetails {
		err := cg.AddChannel(d)
		if err != nil {
			log.Error(fmt.Sprintf("'Error at registering opened channel contract. Perhaps contract is invalid? err=%s, channeladdress=%s",
				err, utils.APex(d.ChannelAddress)))
			cg.RemovePath(d.OurState.Address, d.PartenerState.Address)
		}
	}
	cg.PrintGraph()
	return cg
}
func (this *ChannelGraph) PrintGraph() {
	rowheader := fmt.Sprintf("%s", strings.Repeat(" ", 14))
	for i := 0; i < len(this.index2address); i++ {
		rowheader += fmt.Sprintf("     %s:%2d", utils.APex2(this.index2address[i]), i)
	}
	fmt.Println(rowheader)
	for i := 0; i < len(this.index2address); i++ {
		fmt.Printf("       %s:%2d", utils.APex2(this.index2address[i]), i)
		for j := 0; j < len(this.index2address); j++ {
			v, err := this.g.GetVertex(i)
			if err != nil {
				log.Crit(fmt.Sprintf("addr %s:%d not exist", utils.APex(this.index2address[i]), i))
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

/*
Return a graph that represents the connections among the netting
    contracts.

    Args:
        edge_list (List[(address1, address2)]): All the channels that compose
            the graph.

    Returns:
        Graph A networkx.Graph instance were the graph nodes are nodes in the
            network and the edges are nodes that have a channel between them.
*/
func (this *ChannelGraph) makeGraph(edgeList []common.Address) {
	if len(edgeList)%2 != 0 {
		log.Crit("edgelist must be even")
	}
	for i := 0; i < len(edgeList); i += 2 {
		addr1 := edgeList[i]
		addr2 := edgeList[i+1]
		this.AddPath(addr1, addr2)
	}
}

//Add a new edge into the network.
func (this *ChannelGraph) AddPath(source, target common.Address) {
	addr1 := source
	addr2 := target
	var index1, index2 int
	if index1, ok := this.address2index[addr1]; !ok {
		index1 = len(this.address2index)
		this.address2index[addr1] = index1
		this.index2address[index1] = addr1
	}
	if index2, ok := this.address2index[addr2]; !ok {
		index2 = len(this.address2index)
		this.address2index[addr2] = index2
		this.index2address[index2] = addr2
	}
	index1 = this.address2index[addr1]
	index2 = this.address2index[addr2]
	_, err := this.g.GetVertex(index1)
	if err != nil {
		this.g.AddVertex(index1)
	}
	_, err = this.g.GetVertex(index2)
	if err != nil {
		this.g.AddVertex(index2)
	}
	//todo int64 cannot store too much tokens only about 18 tokens. we should divide 1000 or 1000000,
	err = this.g.AddArc(index1, index2, 1)
	if err != nil {
		log.Error(fmt.Sprintf("add path err%s", err))
	}
	err = this.g.AddArc(index2, index1, 1) //now our graph is the least fee first.
	if err != nil {
		log.Error(fmt.Sprintf("add path err%s", err))
	}
}

/*
Instantiate a channel this node participates and add to the graph.

        If the channel is already registered do nothing.
*/
func (this *ChannelGraph) AddChannel(details *ChannelDetails) error {
	if details.OurState.Address != this.OurAddress {
		return errors.New("Address mismatch, our_address doesn't match the channel details")
	}
	channelAddress := details.ChannelAddress
	if ch := this.GetChannelAddress2Channel(channelAddress); ch == nil {
		ch, err := channel.NewChannel(details.OurState, details.PartenerState,
			details.ExternState, this.TokenAddress, channelAddress, details.BlockChainService,
			details.RevealTimeout, details.SettleTimeout)
		if err != nil {
			return err
		}
		this.Lock.Lock()
		this.PartenerAddress2Channel[details.PartenerState.Address] = ch
		this.ChannelAddress2Channel[channelAddress] = ch
		this.Lock.Unlock()
		this.AddPath(details.OurState.Address, details.PartenerState.Address) //no need to do this? or makegraph is useless?
	}
	return nil
}

/*
Compute all shortest paths in the graph.

        Returns:
            all paths between source and     target.
*/
func (this *ChannelGraph) GetShortestPaths(source, target common.Address) (paths [][]common.Address, err error) {
	panic("not implement")
	return nil, nil
	//this.Lock.Lock()
	//defer this.Lock.Unlock()
	//sourceIndex, ok := this.address2index[source]
	//if !ok {
	//	err = errors.New("source address is unkown")
	//	return
	//}
	//targetIndex, ok := this.address2index[target]
	//if !ok {
	//	err = errors.New("target address is unkown")
	//	return
	//}
	//indexPaths := dijkstra.NewGraph(this.g.GetAllVertices()).AllShortestPath(sourceIndex, targetIndex)
	//for _, ip := range indexPaths {
	//	var p []common.Address
	//	for _, i := range ip {
	//		p = append(p, this.index2address[i])
	//	}
	//	paths = append(paths, p)
	//}
	//return
}

/*
 True if there is a connecting path regardless of the number of hops.
*/
//func (this *ChannelGraph) HasPath(source, target common.Address) bool {
//	return this.ShortestPath(source, target) != utils.MaxInt
//}
func (this *ChannelGraph) HashChannel(source, target common.Address) bool {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	sourceIndex, ok := this.address2index[source]
	if !ok {
		return false
	}
	targetIndex, ok := this.address2index[target]
	if !ok {
		return false
	}
	_, err := this.g.Shortest(sourceIndex, targetIndex)
	return err == nil
}

/*
   def has_channel(self, source_address, target_address):
       """ True if there is a channel connecting both addresses. """
       return self.graph.has_edge(source_address, target_address)
*/
var errAddressNotFoundInGraph = errors.New("address not found in channelgraph")

/*
make sure only be called in one thread.
*/
func (this *ChannelGraph) ShortestPath(source, target common.Address, amount *big.Int, feeCharger fee.FeeCharger) (totalWeight int64, err error) {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	sourceIndex, ok := this.address2index[source]
	if !ok {
		err = errAddressNotFoundInGraph
		return
	}
	targetIndex, ok := this.address2index[target]
	if !ok {
		err = errAddressNotFoundInGraph
		return
	}
	if sourceIndex == targetIndex {
		return 0, nil
	}
	var g2 *dijkstra.Graph
	if false { //make sure only be called in one thread.
		g2 = this.g.CloneGraph()
	} else {
		g2 = this.g
	}
	for _, v := range g2.Verticies {
		w := feeCharger.GetNodeChargeFee(this.index2address[v.ID], this.TokenAddress, amount).Int64()
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

//Remove an edge from the network.  this edge may  not exist
func (this *ChannelGraph) RemovePath(source, target common.Address) {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	sourceIndex, ok := this.address2index[source]
	if !ok {
		return
	}
	targetIndex, ok := this.address2index[target]
	if !ok {
		return
	}
	this.g.DeleteArc(sourceIndex, targetIndex)
	this.g.DeleteArc(targetIndex, sourceIndex)
}

/*
 """ True if the channel with `partner_address` is open and has spendable funds. """
        TODO: check if the partner's network is alive
*/
func (this *ChannelGraph) ChannelCanTransfer(partenerAddress common.Address) bool {
	return this.GetPartenerAddress2Channel(partenerAddress).CanTransfer()
}

//Get all neighbours adjacent to self.our_address. g is not thread safe
func (this *ChannelGraph) GetNeighbours() []common.Address {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	neighboursIndex, err := this.g.GetAllNeighbors(this.address2index[this.OurAddress])
	if err != nil {
		return nil
	}
	var neighbours []common.Address
	for _, i := range neighboursIndex {
		neighbours = append(neighbours, this.index2address[i])
	}
	return neighbours
}

type neighborWeight struct {
	neighbor common.Address
	weight   int64 //nerghbor to target's hops
}
type neighborWeightList []*neighborWeight

func (this neighborWeightList) Len() int {
	return len(this)
}
func (this neighborWeightList) Less(i, j int) bool {
	return this[i].weight < this[j].weight
}
func (this neighborWeightList) Swap(i, j int) {
	var temp *neighborWeight
	temp = this[i]
	this[i] = this[j]
	this[j] = temp
}

/*
all the neighbors that can reach target
they are ordered by hops to the target
*/
func (this *ChannelGraph) orderedNeighbours(ourAddress, targetAddress common.Address, amount *big.Int, charger fee.FeeCharger) neighborWeightList {

	neighbors := this.GetNeighbours()
	var nws neighborWeightList
	for _, n := range neighbors {
		w, err := this.ShortestPath(n, targetAddress, amount, charger)
		if err != nil {
			continue
		}
		nws = append(nws, &neighborWeight{n, w})
	}
	sort.Sort(nws)
	return nws
}

//todo other datatypes
/*
Yield a two-tuple (path, channel) that can be used to mediate the
    transfer. The result is ordered from the best to worst path.
*/
func (this *ChannelGraph) GetBestRoutes(nodesStatus NodesStatusGetter, ourAddress common.Address,
	targetAdress common.Address, amount *big.Int, previousAddress common.Address, feeCharger fee.FeeCharger) (onlineNodes []*transfer.RouteState) {
	/*

	   XXX: consider using multiple channels for a single transfer. Useful
	   for cases were the `amount` is larger than what is available
	   individually in any of the channels.

	   One possible approach is to _not_ filter these channels based on the
	   distributable amount, but to sort them based on available balance and
	   let the task use as many as required to finish the transfer.

	*/
	nws := this.orderedNeighbours(ourAddress, targetAdress, amount, feeCharger)
	if len(nws) == 0 {
		log.Warn(fmt.Sprintf("no routes avaiable from %s to %s", utils.APex(ourAddress), utils.APex(targetAdress)))
		return
	}
	for _, nw := range nws {
		c := this.GetPartenerAddress2Channel(nw.neighbor)
		//don't send the message backwards
		if nw.neighbor == previousAddress {
			continue
		}
		if c.State() != transfer.CHANNEL_STATE_OPENED {
			log.Debug(fmt.Sprintf("channel %s-%s is not opened ,ignoring ..", utils.APex(ourAddress), utils.APex(nw.neighbor)))
			continue
		}
		if amount.Cmp(c.Distributable()) > 0 {
			log.Debug(fmt.Sprintf("channel %s-%s doesn't have enough funds[%d],ignoring...", utils.APex(ourAddress), utils.APex(nw.neighbor), amount))
			continue
		}
		status := nodesStatus.GetNetworkStatus(nw.neighbor)
		if status == NODE_NETWORK_UNREACHABLE {
			log.Debug(fmt.Sprintf("partener %s network ignored.. for:%s", utils.APex(nw.neighbor), status))
			continue
		}
		routeState := Channel2RouteState(c, nw.neighbor, amount, feeCharger)
		if routeState.Fee.Cmp(utils.BigInt0) > 0 {
			routeState.TotalFee = big.NewInt(int64(nw.weight))
		} else { //no fee policy,
			routeState.TotalFee = utils.BigInt0
		}

		onlineNodes = append(onlineNodes, routeState)
	}
	return
}
func (this *ChannelGraph) HaveNodes() bool {
	return len(this.g.Verticies) > 0
}
func (this *ChannelGraph) AllNodes() (nodes []common.Address) {
	for n, _ := range this.address2index {
		nodes = append(nodes, n)
	}
	return nodes
}
func (this *ChannelGraph) GetPartenerAddress2Channel(address common.Address) (c *channel.Channel) {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	c = this.PartenerAddress2Channel[address]
	return
}
func (this *ChannelGraph) GetChannelAddress2Channel(address common.Address) (c *channel.Channel) {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	c = this.ChannelAddress2Channel[address]
	return
}
func Channel2RouteState(c *channel.Channel, partenerAddress common.Address, amount *big.Int, charger fee.FeeCharger) *transfer.RouteState {
	return &transfer.RouteState{
		State:          c.State(),
		HopNode:        partenerAddress,
		ChannelAddress: c.MyAddress,
		AvaibleBalance: c.Distributable(),
		SettleTimeout:  c.SettleTimeout,
		RevealTimeout:  c.RevealTimeout,
		ClosedBlock:    c.ExternState.ClosedBlock, //default is 0
		Fee:            charger.GetNodeChargeFee(partenerAddress, c.TokenAddress, amount),
	}
}
