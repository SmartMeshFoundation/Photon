package network

import (
	"errors"

	"fmt"

	"sort"

	"strings"

	"sync"

	"math/big"

	"github.com/SmartMeshFoundation/raiden-network/channel"
	"github.com/SmartMeshFoundation/raiden-network/network/rpc"
	"github.com/SmartMeshFoundation/raiden-network/transfer"
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/nkbai/dijkstra"
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
		g:                       dijkstra.NewEmptyGraph(),
	}
	cg.makeGraph(edgeList)
	for _, d := range channelDetails {
		err := cg.AddChannel(d)
		if err != nil {
			log.Warn(fmt.Sprintf("'Error at registering opened channel contract. Perhaps contract is invalid? err=%s, channeladdress=%s",
				err, utils.APex(d.ChannelAddress)))
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
			if this.g.HasEdge(i, j) {
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
	this.g.AddEdge(index1, index2, 1)
	this.g.AddEdge(index2, index1, 1)
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
	this.Lock.Lock()
	defer this.Lock.Unlock()
	sourceIndex, ok := this.address2index[source]
	if !ok {
		err = errors.New("source address is unkown")
		return
	}
	targetIndex, ok := this.address2index[target]
	if !ok {
		err = errors.New("target address is unkown")
		return
	}
	indexPaths := dijkstra.NewGraph(this.g.GetAllVertices()).AllShortestPath(sourceIndex, targetIndex)
	for _, ip := range indexPaths {
		var p []common.Address
		for _, i := range ip {
			p = append(p, this.index2address[i])
		}
		paths = append(paths, p)
	}
	return
}

/*
 True if there is a connecting path regardless of the number of hops.
*/
func (this *ChannelGraph) HasPath(source, target common.Address) bool {
	return this.ShortestPath(source, target) != utils.MaxInt
}
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
	return this.g.HasEdge(sourceIndex, targetIndex)
}

/*
   def has_channel(self, source_address, target_address):
       """ True if there is a channel connecting both addresses. """
       return self.graph.has_edge(source_address, target_address)
*/
func (this *ChannelGraph) ShortestPath(source, target common.Address) int {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	sourceIndex, ok := this.address2index[source]
	if !ok {
		return utils.MaxInt
	}
	targetIndex, ok := this.address2index[target]
	if !ok {
		return utils.MaxInt
	}
	if sourceIndex == targetIndex {
		return 0
	}
	return dijkstra.NewGraph(this.g.GetAllVertices()).ShortestPath(sourceIndex, targetIndex)
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
	this.g.RemoveEdge(sourceIndex, targetIndex)
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
	neighboursIndex := this.g.GetAllNeighbours(this.address2index[this.OurAddress])
	var neighbours []common.Address
	for _, i := range neighboursIndex {
		neighbours = append(neighbours, this.index2address[i])
	}
	return neighbours
}

type neighborWeight struct {
	neighbor common.Address
	weight   int //nerghbor to target's hops
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
func (this *ChannelGraph) orderedNeighbours(ourAddress, targetAddress common.Address) []common.Address {

	neighbors := this.GetNeighbours()
	var nws neighborWeightList
	for _, n := range neighbors {
		w := this.ShortestPath(n, targetAddress)
		if n == targetAddress {
			log.Debug("equal")
		}
		nws = append(nws, &neighborWeight{n, w})
	}
	sort.Sort(nws)
	log.Trace(fmt.Sprintf("nws=%s\n", utils.StringInterface(nws, 3)))
	neighbors = []common.Address{}
	for _, nw := range nws {
		if nw.weight != utils.MaxInt {
			neighbors = append(neighbors, nw.neighbor)
		}
	}
	return neighbors
}

//todo other datatypes
/*
Yield a two-tuple (path, channel) that can be used to mediate the
    transfer. The result is ordered from the best to worst path.
*/
func (this *ChannelGraph) GetBestRoutes(nodesStatus NodesStatusGeter, ourAddress common.Address,
	targetAdress common.Address, amount *big.Int, previousAddress common.Address) (onlineNodes []*transfer.RouteState) {
	///*
	//	for direct transfer
	//*/
	//c := this.GetPartenerAddress2Channel(targetAdress)
	//if c != nil {
	//	onlineNodes = append(onlineNodes, Channel2RouteState(c, targetAdress))
	//	return
	//}
	/*

	   XXX: consider using multiple channels for a single transfer. Useful
	   for cases were the `amount` is larger than what is available
	   individually in any of the channels.

	   One possible approach is to _not_ filter these channels based on the
	   distributable amount, but to sort them based on available balance and
	   let the task use as many as required to finish the transfer.

	*/
	neighbors := this.orderedNeighbours(ourAddress, targetAdress)
	if len(neighbors) == 0 {
		log.Warn(fmt.Sprint("no routes avaiable from %s to %s", utils.APex(ourAddress), utils.APex(targetAdress)))
		return
	}
	for _, partnerAddress := range neighbors {
		c := this.GetPartenerAddress2Channel(partnerAddress)
		//don't send the message backwards
		if partnerAddress == previousAddress {
			continue
		}
		if c.State() != transfer.CHANNEL_STATE_OPENED {
			log.Debug(fmt.Sprintf("channel %s-%s is not opened ,ignoring ..", utils.APex(ourAddress), utils.APex(partnerAddress)))
			continue
		}
		if amount.Cmp(c.Distributable()) > 0 {
			log.Debug(fmt.Sprintf("channel %s-%s doesn't have enough funds[%d],ignoring...", utils.APex(ourAddress), utils.APex(partnerAddress), amount))
			continue
		}
		status := nodesStatus.GetNetworkStatus(partnerAddress)
		if status == NODE_NETWORK_UNREACHABLE {
			log.Debug(fmt.Sprintf("partener %s network ignored.. for:%s", utils.APex(partnerAddress), status))
			continue
		}
		routeState := Channel2RouteState(c, partnerAddress)
		onlineNodes = append(onlineNodes, routeState)
	}
	return
}
func (this *ChannelGraph) HaveNodes() bool {
	return this.g.Len() > 0
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
func Channel2RouteState(c *channel.Channel, partenerAddress common.Address) *transfer.RouteState {
	return &transfer.RouteState{
		State:          c.State(),
		HopNode:        partenerAddress,
		ChannelAddress: c.MyAddress,
		AvaibleBalance: c.Distributable(),
		SettleTimeout:  c.SettleTimeout,
		RevealTimeout:  c.RevealTimeout,
		ClosedBlock:    c.ExternState.ClosedBlock, //default is 0
	}
}
