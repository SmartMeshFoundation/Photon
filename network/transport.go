package network

import (
	"bytes"
	"context"
	"time"

	mdns2 "github.com/whyrusleeping/mdns"

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

	/*
		为了避免消息接收方不在线时不停尝试重发,添加挂起/唤醒机制
	*/
	iWakeUpHandler
}

//MixTranspoter support udp and others(xmpp or matrix)
type MixTranspoter interface {
	Transporter
	//UDPNodeStatus status of UDPTransport
	UDPNodeStatus(addr common.Address) (deviceType string, isOnline bool)
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
	conn                   *net.UDPConn
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
	monitorIP              chan struct{}
	mdnsLock               sync.Mutex
	*wakeupHandler
}

//NewUDPTransport create UDPTransport,name必须是完整的地址
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
		wakeupHandler:          newWakeupHandler("udp"),
	}
	//127.0.0.1 作为一个特殊地址来处理,作为不启用mdns的指示,但是127.1.0.1等其他本机ip地址都认为有效
	if params.EnableMDNS {
		ctx, cf := context.WithCancel(context.Background())
		t.msrv, err = mdns.NewMdnsService(ctx, port, name, params.DefaultMDNSQueryInterval)
		if err != nil {
			log.Error(fmt.Sprintf("NewMdnsService err %s", err))
			cf()
			err = nil //无论mdns是否启动完成,都不能影响photon启动
		} else {
			t.cf = cf
			t.msrv.RegisterNotifee(t)
		}
		t.monitorIP = make(chan struct{})
		go t.monitorIPChange()
	}
	return
}

//考虑到手机或者盒子在使用photon的过程中可能会发生连接热点切换的问题,从而导致ip地址变化
func (ut *UDPTransport) monitorIPChange() {
	var err error
	lastip := mdns2.GetLocalIP()
	for {
		select {
		case <-time.After(time.Second * 10):
			ut.removeExpiredNodes() //先移除过期的无效的地址信息
			newip := mdns2.GetLocalIP()
			changeip := func() {
				log.Info(fmt.Sprintf("dectecipchange,last=%s,now=%s", lastip, newip))
				ut.mdnsLock.Lock()
				defer ut.mdnsLock.Unlock()
				if ut.msrv != nil {
					ut.cf()
					ut.cf = nil
					err = ut.msrv.Close()
					if err != nil {
						log.Error(fmt.Sprintf("close msrv err %s", err))
					}
					ut.msrv = nil
				}
				ctx, cf := context.WithCancel(context.Background())
				ut.msrv, err = mdns.NewMdnsService(ctx, ut.UAddr.Port, ut.name, params.DefaultMDNSQueryInterval)
				if err != nil {
					log.Error(fmt.Sprintf("NewMdnsService err %s", err))
					cf()
					return
				}
				ut.cf = cf
				ut.msrv.RegisterNotifee(ut)
				lastip = newip

			}
			if len(lastip) != len(newip) {
				//do changeip
				changeip()
				break
			}
			for i := 0; i < len(lastip); i++ {
				if bytes.Compare(lastip[i], newip[i]) != 0 {
					changeip()
					break
				}
			}
		case <-ut.monitorIP:
			return
		}
	}
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
	ut.mdnsLock.Lock()
	if ut.cf != nil {
		ut.cf()
	}
	if ut.msrv != nil {
		err := ut.msrv.Close()
		if err != nil {
			log.Error(fmt.Sprintf("udp transport stop err %s", err))
		}
	}
	if ut.monitorIP != nil {
		close(ut.monitorIP)
	}
	ut.mdnsLock.Unlock()
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

func (ut *UDPTransport) removeExpiredNodes() {
	ut.lock.Lock()
	defer ut.lock.Unlock()
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
	}
	for _, idToDelete := range idsToDelete {
		previousIP := ut.intranetNodes[idToDelete]
		delete(ut.intranetNodes, idToDelete)
		delete(ut.intranetNodesTimestamp, idToDelete)
		log.Info(fmt.Sprintf("peer UDP offline id=%s,previous Ip=%s", idToDelete.String(), previousIP.String()))
	}
}

//HandlePeerFound notification  from mdns
func (ut *UDPTransport) HandlePeerFound(id string, addr *net.UDPAddr) {
	ut.removeExpiredNodes()
	ut.lock.Lock()
	defer ut.lock.Unlock()
	idFound := common.HexToAddress(id)
	_, alreadyFound := ut.intranetNodes[idFound]

	// 标记发现的除自己以外的节点
	if id != ut.name {
		if !alreadyFound {
			log.Info(fmt.Sprintf("peer UDP found id=%s,addr=%s", id, addr))
			// 唤醒
			ut.wakeUp(idFound)
		}
		ut.intranetNodes[idFound] = addr
		ut.intranetNodesTimestamp[idFound] = time.Now()
	}
}
