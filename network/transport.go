package network

import (
	"context"
	"time"

	"github.com/SmartMeshFoundation/Photon/network/mdns"

	"fmt"

	"net"

	"errors"
	"sync"

	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/internal/rpanic"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/network/xmpptransport"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

//Policier to control the sending speed of transporter
type Policier interface {
	//Consume tokens.
	//Args:
	//tokens (float): number of transport tokens to consume
	//Returns:
	//wait_time (float): waiting time for the consumer
	Consume(tokens float64) time.Duration
}

//DeviceTypeMobile if you are a Photon running on a mobile phone
var DeviceTypeMobile = xmpptransport.TypeMobile

//DeviceTypeMeshBox if you are a Photon running on a meshbox
var DeviceTypeMeshBox = xmpptransport.TypeMeshBox

//DeviceTypeOther if you don't known the type,and is not a mobile phone, then other
var DeviceTypeOther = xmpptransport.TypeOtherDevice

//Transporter denotes a communication transport used by protocol
type Transporter interface {
	//Send a message to receiver
	Send(receiver common.Address, data []byte) error
	//Start ,ready for send and receive
	Start()
	//Stop send and receive
	Stop()
	//StopAccepting stops receiving
	StopAccepting()
	//RegisterProtocol a receiver
	RegisterProtocol(protcol ProtocolReceiver)
	//NodeStatus get node's status and is online right now
	NodeStatus(addr common.Address) (deviceType string, isOnline bool)
}

type dummyPolicy struct {
}

//Consume mocker
func (dp *dummyPolicy) Consume(tokens float64) time.Duration {
	time.Now()
	return 0
}

type timeFunc func() time.Time

//TokenBucket Implementation of the token bucket throttling algorithm.
type TokenBucket struct {
	Capacity  float64
	FillRate  float64
	Tokens    float64
	timeFunc  timeFunc
	Timestamp time.Time
}

//NewTokenBucket create a TokenBucket
func NewTokenBucket(capacity, fillRate float64, timeFunc ...timeFunc) *TokenBucket {
	tb := &TokenBucket{
		Capacity: capacity,
		FillRate: fillRate,
		Tokens:   capacity,
	}
	if len(timeFunc) == 1 {
		tb.timeFunc = timeFunc[0]
	} else {
		tb.timeFunc = time.Now
	}
	tb.Timestamp = tb.timeFunc()
	return tb
}

//Consume calc wait time.
func (tb *TokenBucket) Consume(tokens float64) time.Duration {
	waitTime := 0.0
	tb.Tokens -= tokens
	if tb.Tokens < 0 {
		tb.getTokens()
	}
	if tb.Tokens < 0 {
		waitTime = -tb.Tokens / tb.FillRate
	}
	return time.Duration(waitTime * float64(time.Second))
}
func (tb *TokenBucket) getTokens() {
	now := tb.timeFunc()
	fill := float64(now.Sub(tb.Timestamp)) / float64(time.Second)
	tb.Tokens += tb.FillRate * fill
	if tb.Tokens > tb.Capacity {
		tb.Tokens = tb.Capacity
	}
	tb.Timestamp = tb.timeFunc()
}

//ProtocolReceiver receive
type ProtocolReceiver interface {
	receive(data []byte)
}

//
/*
UDPTransport represents a UDP server
but how to handle listen error?
we need stop listen when switch to background
restart listen when switch foreground
*/
type UDPTransport struct {
	protocol               ProtocolReceiver
	conn                   *SafeUDPConnection
	UAddr                  *net.UDPAddr
	policy                 Policier
	stopped                bool
	stopReceiving          bool
	intranetNodes          map[common.Address]*net.UDPAddr
	intranetNodesTimestamp map[common.Address]time.Time
	lock                   sync.RWMutex
	name                   string
	log                    log.Logger
	msrv                   mdns.Service
	cf                     context.CancelFunc
}

//NewUDPTransport create UDPTransport
func NewUDPTransport(name, host string, port int, protocol ProtocolReceiver, policy Policier) (t *UDPTransport, err error) {
	t = &UDPTransport{
		UAddr: &net.UDPAddr{
			IP:   net.ParseIP(host),
			Port: port,
		},
		name:                   name,
		protocol:               protocol,
		policy:                 policy,
		log:                    log.New("name", name),
		intranetNodes:          make(map[common.Address]*net.UDPAddr),
		intranetNodesTimestamp: make(map[common.Address]time.Time),
	}
	//127.0.0.1 作为一个特殊地址来处理,作为不启用mdns的指示,但是127.1.0.1等其他本机ip地址都认为有效
	if host != "127.0.0.1" {
		ctx, cf := context.WithCancel(context.Background())
		t.msrv, err = mdns.NewMdnsService(ctx, port, name, params.DefaultMDNSQueryInterval)
		if err != nil {
			cf()
			return
		}
		t.cf = cf
		t.msrv.RegisterNotifee(t)
	}
	return
}

//Start udp listening
func (ut *UDPTransport) Start() {
	go func() {
		data := make([]byte, 4096)
		defer rpanic.PanicRecover("udptransport Start")
		for {
			conn, err := NewSafeUDPConnection("udp", ut.UAddr)
			if err != nil {
				log.Error(fmt.Sprintf("listen udp %s error %v", ut.UAddr.String(), err))
				time.Sleep(time.Second)
				continue
			}
			log.Info(fmt.Sprintf("udp server listening on %s", ut.UAddr.String()))
			ut.conn = conn
			ut.log.Info(fmt.Sprintf(" listen udp on %s", ut.UAddr))
			for {
				if ut.stopReceiving {
					return
				}
				read, remoteAddr, err := ut.conn.ReadFromUDP(data)
				if err != nil {
					if !ut.stopped {
						ut.log.Error(fmt.Sprintf("udp read data failure! %s", err))
						err = ut.conn.Close()
						break
					} else {
						return
					}

				}
				ut.log.Trace(fmt.Sprintf("receive from %s ,message=%s,hash=%s", remoteAddr,
					encoding.MessageType(data[0]), utils.HPex(utils.Sha3(data[:read]))))
				err = ut.Receive(data[:read])
			}
		}

	}()
	time.Sleep(time.Millisecond)
}

//Receive a message
func (ut *UDPTransport) Receive(data []byte) error {
	//ut.log.Trace(fmt.Sprintf("recevied data\n%s", hex.Dump(data)))
	if ut.stopReceiving {
		return errors.New("stop receive")
	}
	if ut.protocol != nil { //receive data before register a protocol
		ut.protocol.receive(data)
	}
	return nil
}

/*
Send `bytes_` to `host_port`.
Args:
    sender (address): The address of the running node.
    host_port (Tuple[(str, int)]): Tuple with the Host name and Port number.
    bytes_ (bytes): The bytes that are going to be sent through the wire.
*/
func (ut *UDPTransport) Send(receiver common.Address, data []byte) error {
	if ut.stopped {
		return fmt.Errorf("%s closed", ut.name)
	}
	ua, err := ut.getHostPort(receiver)
	if err != nil {
		return err
	}
	ut.log.Trace(fmt.Sprintf("%s send to %s %s:%d, message=%s,response hash=%s", ut.name,
		utils.APex2(receiver), ua.IP, ua.Port, encoding.MessageType(data[0]),
		utils.HPex(utils.Sha3(data, receiver[:]))))
	//ut.log.Trace(fmt.Sprintf("send data  \n%s", hex.Dump(data)))
	//only comment this line,if you want to test.
	//time.Sleep(ut.policy.Consume(1)) //force to wait,
	_, err = ut.conn.WriteToUDP(data, ua)
	return err
}

func (ut *UDPTransport) getHostPort(addr common.Address) (ua *net.UDPAddr, err error) {
	ut.lock.RLock()
	defer ut.lock.RUnlock()
	ua, ok := ut.intranetNodes[addr]
	if ok {
		return
	}
	err = fmt.Errorf("%s host port not found", utils.APex(addr))
	return
}
func (ut *UDPTransport) setHostPort(nodes map[common.Address]*net.UDPAddr) {
	ut.lock.Lock()
	defer ut.lock.Unlock()
	for k, v := range nodes {
		ut.intranetNodes[k] = v
	}
}

//RegisterProtocol register receiver
func (ut *UDPTransport) RegisterProtocol(proto ProtocolReceiver) {
	ut.protocol = proto
}

//Stop UDP connection
func (ut *UDPTransport) Stop() {
	if ut.cf != nil {
		ut.cf()
	}
	if ut.msrv != nil {
		err := ut.msrv.Close()
		if err != nil {
			log.Error(fmt.Sprintf("udp transport stop err %s", err))
		}
	}
	ut.stopReceiving = true
	ut.stopped = true
	ut.intranetNodes = make(map[common.Address]*net.UDPAddr)
	if ut.conn != nil {
		err := ut.conn.Close()
		if err != nil {
			log.Warn(fmt.Sprintf("close err %s ", err))
		}
	}
}

//StopAccepting stop receiving
func (ut *UDPTransport) StopAccepting() {
	ut.stopReceiving = true
}

//NodeStatus always mark the node offline
func (ut *UDPTransport) NodeStatus(addr common.Address) (deviceType string, isOnline bool) {
	ut.lock.RLock()
	defer ut.lock.RUnlock()
	if _, ok := ut.intranetNodes[addr]; ok {
		return DeviceTypeOther, true
	}
	return DeviceTypeOther, false
}

//HandlePeerFound notification  from mdns
func (ut *UDPTransport) HandlePeerFound(id string, addr *net.UDPAddr) {
	ut.lock.Lock()
	defer ut.lock.Unlock()
	idFound := common.HexToAddress(id)
	alreadyFound := false
	// 清除过期数据,即标志下线
	now := time.Now()
	var idsToDelete []common.Address
	for idTemp := range ut.intranetNodes {
		saveTime, ok := ut.intranetNodesTimestamp[idTemp]
		if !ok {
			// 不处理非自己发现的节点
			continue
		}
		if now.Sub(saveTime) > params.DefaultMDNSKeepalive {
			idsToDelete = append(idsToDelete, idTemp)
		}
		if idTemp == idFound {
			alreadyFound = true
		}
	}
	for _, idToDelete := range idsToDelete {
		delete(ut.intranetNodes, idToDelete)
		delete(ut.intranetNodesTimestamp, idToDelete)
		log.Info(fmt.Sprintf("peer UDP offline id=%s", idToDelete.String()))
	}
	// 标记发现的除自己以外的节点
	if id != ut.name {
		if !alreadyFound {
			log.Info(fmt.Sprintf("peer UDP found id=%s,addr=%s", id, addr))
		}
		ut.intranetNodes[idFound] = addr
		ut.intranetNodesTimestamp[idFound] = now
	}
}
