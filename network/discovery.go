package network

import (
	"context"
	"fmt"

	"strconv"
	"strings"

	"errors"

	"github.com/SmartMeshFoundation/SmartRaiden/blockchain"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/astaxie/beego/httplib"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var hostportPrefix = []byte("hostport_")
var nodePrefix = []byte("node_")

//DiscoveryInterface all discovery should follow the interface
type DiscoveryInterface interface {
	Register(address common.Address, host string, port int) error
	Get(address common.Address) (host string, port int, err error)
	NodeIDByHostPort(host string, port int) (node common.Address, err error)
}

//HTTPDiscovery for test only,
type HTTPDiscovery struct {
	Path string
}

//NewHTTPDiscovery create httpDiscovery
func NewHTTPDiscovery() *HTTPDiscovery {
	return &HTTPDiscovery{
		Path: "http://182.254.155.208:50000/public/",
	}
}

//Register a node's ip and port
func (h *HTTPDiscovery) Register(address common.Address, host string, port int) error {
	_, err := httplib.Get(fmt.Sprintf(h.Path+"register?addr=%s&hostport=%s", strings.ToLower(address.String()), fmt.Sprintf("%s:%d", host, port))).String()
	return err
}

//Get return a node's ip and port
func (h *HTTPDiscovery) Get(address common.Address) (host string, port int, err error) {
	s, err := httplib.Get(fmt.Sprintf(h.Path+"gethostport?addr=%s", strings.ToLower(address.String()))).String()
	host, port = SplitHostPort(s)
	return
}

//NodeIDByHostPort return a node's address by ip and port
func (h *HTTPDiscovery) NodeIDByHostPort(host string, port int) (node common.Address, err error) {
	s, err := httplib.Get(fmt.Sprintf(h.Path+"getaddr?hostport=%s", fmt.Sprintf("%s:%d", host, port))).String()
	if err != nil {
		return
	}
	node = common.HexToAddress(s)
	return
}

//Discovery endpoint info saved in memory ,test only
type Discovery struct {
	NodeIDHostPortMap map[common.Address]string
}

//NewDiscovery create Discovery
func NewDiscovery() *Discovery {
	return &Discovery{
		NodeIDHostPortMap: make(map[common.Address]string),
	}
}

//Register a node
func (d *Discovery) Register(address common.Address, host string, port int) error {
	d.NodeIDHostPortMap[address] = fmt.Sprintf("%s:%d", host, port)
	return nil
}

//SplitHostPort is the same as net.SplitHostPort
//todo remove this function
func SplitHostPort(hostport string) (string, int) {
	ss := strings.Split(hostport, ":")
	port, _ := strconv.Atoi(ss[1])
	return ss[0], port
}

//Get a node's ip and port
func (d *Discovery) Get(address common.Address) (host string, port int, err error) {
	hostport, ok := d.NodeIDHostPortMap[address]
	if !ok {
		err = errors.New("no such address")
		return
	}
	host, port = SplitHostPort(hostport)
	return
}

//NodeIDByHostPort find a node by ip and port
func (d *Discovery) NodeIDByHostPort(host string, port int) (node common.Address, err error) {
	hostport := tohostport(host, port)
	for k, v := range d.NodeIDHostPortMap {
		if v == hostport {
			return k, nil
		}
	}
	return utils.EmptyAddress, errors.New("not found")
}

//ContractDiscovery Allows registering and looking up by endpoint (Host, Port) for node_address.
type ContractDiscovery struct {
	discovery   *contracts.EndpointRegistry
	client      *helper.SafeEthClient
	auth        *bind.TransactOpts
	node        common.Address
	cacheAddr   map[common.Address]string
	cacheSocket map[string]common.Address
	myaddress   common.Address
}

//NewContractDiscovery create ContractDiscovery
func NewContractDiscovery(mynode, myaddress common.Address, client *helper.SafeEthClient, auth *bind.TransactOpts) *ContractDiscovery {
	c := &ContractDiscovery{
		client:      client,
		auth:        auth,
		node:        mynode,
		cacheAddr:   make(map[common.Address]string),
		cacheSocket: make(map[string]common.Address),
		myaddress:   myaddress,
	}
	c.discovery, _ = contracts.NewEndpointRegistry(myaddress, client)
	ch := make(chan types.Log, 1)
	_, err := rpc.EventSubscribe(myaddress, "AddressRegistered", contracts.EndpointRegistryABI, client, ch)
	if err == nil {
		//c.db = db.GetDefaultDb()
		go func() {
			for { //monitor event on chain
				l := <-ch
				ev, err := blockchain.NewEventAddressRegistered(&l)
				if err != nil {
					continue
				}
				log.Debug(fmt.Sprintf("receive register node=%s, socket=%s", ev.EthAddress.String(), ev.Socket))
				c.put(ev.EthAddress, ev.Socket)
			}

		}()
	}
	return c
}

func (c *ContractDiscovery) put(addr common.Address, hostport string) {
	c.cacheAddr[addr] = hostport
	c.cacheSocket[hostport] = addr
	//c.db.Put(append(nodePrefix, addr[:]...), []byte(hostport))
	//c.db.Put(append(hostportPrefix, []byte(hostport)...), addr[:])
}

/*
Register a node
it may emit a Tx on block chain
*/
func (c *ContractDiscovery) Register(node common.Address, host string, port int) error {
	if node != c.node {
		log.Crit(fmt.Sprintf("register node to contract with unknown addr %s", utils.APex(node)))
	}
	log.Info(fmt.Sprintf("ContractDiscovery register %s %s:%d", utils.APex(node), host, port))
	h1, p1, err := c.Get(node)
	//my node's host and port donesn't change after restart
	if err == nil && h1 == host && p1 == port {
		return nil
	}
	hostport := tohostport(host, port)
	tx, err := c.discovery.RegisterEndpoint(c.auth, hostport)
	if err != nil {
		return fmt.Errorf("RegisterEndpoint %s", err)
	}
	//wait for completion ?
	_, err = bind.WaitMined(context.Background(), c.client, tx)
	if err != nil {
		return fmt.Errorf("WaitMined %s", err)
	}
	return nil
}

//Get a node's ip and port from blockchain
func (c *ContractDiscovery) Get(node common.Address) (host string, port int, err error) {
	//hostport, err := c.db.Get(append(nodePrefix, node[:]...))
	hostport, ok := c.cacheAddr[node]
	if ok {
		host, port = SplitHostPort(string(hostport))
		return
	}
	hostportstr, err := c.discovery.FindEndpointByAddress(nil, node)
	if err == nil && len(hostportstr) > 0 {
		c.put(node, hostportstr)
		host, port = SplitHostPort(hostportstr)
	}
	return
}

//NodeIDByHostPort find a node by ip and port
func (c *ContractDiscovery) NodeIDByHostPort(host string, port int) (node common.Address, err error) {
	hostport := tohostport(host, port)
	//addr, err := c.db.Get(append(hostportPrefix, []byte(hostport)...))
	addr, ok := c.cacheSocket[hostport]
	if ok {
		return addr, nil
	}
	node, err = c.discovery.FindAddressByEndpoint(nil, hostport)
	if err != nil && node != utils.EmptyAddress {
		c.put(node, hostport)
	}
	return
}
