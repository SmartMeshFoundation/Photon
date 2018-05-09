package network

import (
	"time"

	"fmt"

	"net"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

type Policier interface {
	//"""Consume tokens.
	//Args:
	//tokens (float): number of transport tokens to consume
	//Returns:
	//wait_time (float): waiting time for the consumer
	Consume(tokens float64) time.Duration
}
type Transporter interface {
	Send(receiver common.Address, host string, port int, data []byte) error
	Receive(data []byte, host string, port int) error
	Start()
	Stop()
	StopAccepting()                            //stop receiving data
	RegisterProtocol(protcol ProtocolReceiver) //register transporter to protocol
}
type MessageCallBack func(sender common.Address, hostport string, msg []byte)

func tohostport(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

type DummyPolicy struct {
}

func (this *DummyPolicy) Consume(tokens float64) time.Duration {
	time.Now()
	return 0
}

type TimeFunc func() time.Time

//Implementation of the token bucket throttling algorithm.
type TokenBucket struct {
	Capacity  float64
	FillRate  float64
	Tokens    float64
	timeFunc  TimeFunc
	Timestamp time.Time
}

func NewTokenBucket(capacity, fillRate float64, timeFunc ...TimeFunc) *TokenBucket {
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

func (this *TokenBucket) Consume(tokens float64) time.Duration {
	waitTime := 0.0
	this.Tokens -= tokens
	if this.Tokens < 0 {
		this.getTokens()
	}
	if this.Tokens < 0 {
		waitTime = -this.Tokens / this.FillRate
	}
	return time.Duration(waitTime * float64(time.Second))
}
func (this *TokenBucket) getTokens() {
	now := this.timeFunc()
	fill := float64(now.Sub(this.Timestamp)) / float64(time.Second)
	this.Tokens += this.FillRate * fill
	if this.Tokens > this.Capacity {
		this.Tokens = this.Capacity
	}
	this.Timestamp = this.timeFunc()
}

//Store global state for an in process network, this won't use a real
//network protocol
type DummyNetwork struct {
	Transports              map[string]Transporter
	Counter                 int
	MessageSendCallbacks    []MessageCallBack
	MessageReceiveCallbacks []MessageCallBack
}

func NewDummyNetwork() *DummyNetwork {
	return &DummyNetwork{
		Transports: make(map[string]Transporter),
		Counter:    0,
	}
}

var dummyNetwork = NewDummyNetwork()

func RegisterSendCallback(cb MessageCallBack) {
	dummyNetwork.MessageSendCallbacks = append(dummyNetwork.MessageSendCallbacks, cb)
}
func RegisterReceiveCallback(cb MessageCallBack) {
	dummyNetwork.MessageReceiveCallbacks = append(dummyNetwork.MessageReceiveCallbacks, cb)
}

//Register a new node in the dummy network.
func (this *DummyNetwork) Register(transpoter Transporter, host string, port int) {
	hostport := fmt.Sprintf("%s:%d", host, port)
	this.Transports[hostport] = transpoter
}

//Register an attempt to send a packet. This method should be called
//everytime send() is used.
func (this *DummyNetwork) TrackSend(receiver common.Address, host string, port int, data []byte) error {
	this.Counter += 1
	for _, cb := range this.MessageSendCallbacks {
		cb(receiver, tohostport(host, port), data)
	}
	return nil
}

func (this *DummyNetwork) TrackReceive(receiver common.Address, host string, port int, data []byte) {
	for _, cb := range this.MessageReceiveCallbacks {
		cb(receiver, tohostport(host, port), data)
	}
}
func (this *DummyNetwork) Send(sender common.Address, host string, port int, data []byte) error {
	this.TrackSend(sender, host, port, data)
	hostport := tohostport(host, port)
	time.AfterFunc(time.Nanosecond, func() {
		this.Transports[hostport].Receive(data, host, port)
	})
	return nil
}

type ProtocolReceiver interface {
	Receive(data []byte, host string, port int)
}
type UDPTransport struct {
	protocol      ProtocolReceiver
	conn          *SafeUdpConnection
	Host          string
	Port          int
	policy        Policier
	isClosed      bool
	stopReceiving bool //todo use atomic to replace
}

func NewUDPTransport(host string, port int, conn *SafeUdpConnection, protocol ProtocolReceiver, policy Policier) *UDPTransport {
	t := &UDPTransport{
		Host:          host,
		Port:          port,
		protocol:      protocol,
		policy:        policy,
		isClosed:      false,
		stopReceiving: false,
	}
	addr := &net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: port}
	var err error
	if conn == nil {
		conn, err = NewSafeUdpConnection("udp", addr)
		if err != nil {
			log.Crit(fmt.Sprintf("listen udp %s:%d error %v", host, port, err))
		}
	}
	t.conn = conn
	log.Trace(fmt.Sprintf("listen udp on %s:%d", host, port))
	return t
}
func NewUDPTransportWithHostPort(host string, port int, protocol ProtocolReceiver, policy Policier) *UDPTransport {
	return NewUDPTransport(host, port, nil, protocol, policy)
}

func (this *UDPTransport) Start() {
	t := this
	go func() {
		data := make([]byte, 4096)
		for {
			if t.stopReceiving {
				break
			}
			read, remoteAddr, err := t.conn.ReadFromUDP(data)
			//log.Trace("receive data:")
			if err != nil {
				fmt.Println("udp read data failure!", err)
				if !t.isClosed {
					continue
				} else {
					return
				}

			}
			log.Trace(fmt.Sprintf("%d receive from %s:%d,message=%s,hash=%s\n", t.Port, remoteAddr.IP.String(),
				remoteAddr.Port, encoding.MessageType(data[0]), utils.HPex(utils.Sha3(data[:read]))))
			t.Receive(data[:read], remoteAddr.IP.String(), remoteAddr.Port)
		}

	}()
}
func (this *UDPTransport) Receive(data []byte, host string, port int) error {
	//todo fix get raiden address, my node address
	dummyNetwork.TrackReceive(common.Address{}, host, port, data)
	if this.protocol != nil { //receive data before register a protocol
		this.protocol.Receive(data, host, port)
	}

	return nil
}
func udpAddrFromHostport(host string, port int) *net.UDPAddr {
	//ss := strings.Split(hostport, ":")
	//Host := ss[0]
	//Port, _ := strconv.Atoi(ss[1])
	return &net.UDPAddr{IP: net.ParseIP(host), Port: port}
}

/*
Send `bytes_` to `host_port`.
Args:
    sender (address): The address of the running node.
    host_port (Tuple[(str, int)]): Tuple with the Host name and Port number.
    bytes_ (bytes): The bytes that are going to be sent through the wire.
*/
func (this *UDPTransport) Send(receiver common.Address, host string, port int, data []byte) error {
	dummyNetwork.TrackSend(receiver, host, port, data)
	log.Trace(fmt.Sprintf("%d send to %s %s:%d, message=%s,hash=%s\n", this.Port, utils.APex(receiver), host, port, encoding.MessageType(data[0]), utils.HPex(utils.Sha3(data, receiver[:]))))
	time.Sleep(this.policy.Consume(1))
	//todo need one lock for write?
	_, err := this.conn.WriteToUDP(data, udpAddrFromHostport(host, port))
	if err != nil {
		return err
	}
	return nil
}

func (this *UDPTransport) RegisterProtocol(proto ProtocolReceiver) {
	this.protocol = proto
}
func (this *UDPTransport) Stop() {
	this.isClosed = true
	this.conn.Close()
}
func (this *UDPTransport) StopAccepting() {
	this.stopReceiving = true
}

// Communication between inter-process nodes.
type DummyTransport struct {
	protocol ProtocolReceiver
	host     string
	port     int
	policy   Policier
}

func NewDummyTransport(host string, port int, protocol ProtocolReceiver, policy Policier) *DummyTransport {
	t := &DummyTransport{
		protocol: protocol,
		host:     host,
		port:     port,
		policy:   policy,
	}
	dummyNetwork.Register(t, host, port)
	return t
}
func (this *DummyTransport) Send(receiver common.Address, host string, port int, data []byte) error {
	time.Sleep(this.policy.Consume(1))
	return dummyNetwork.Send(receiver, host, port, data)
}
func (this *DummyTransport) Receive(data []byte, host string, port int) error {
	dummyNetwork.TrackReceive(common.Address{}, host, port, data)
	this.protocol.Receive(data, host, port)
	return nil
}
func (this *DummyTransport) RegisterProtocol(protcol ProtocolReceiver) {
	this.protocol = protcol
}
func (this *DummyTransport) Start() {

}
func (this *DummyTransport) Stop() {

}
func (this *DummyTransport) StopAccepting() {

}

type UnreliableTransport struct {
	DummyTransport
	DropRate int
}

func NewUnreliableTransport(t *DummyTransport) *UnreliableTransport {
	return &UnreliableTransport{DummyTransport: *t, DropRate: 2}
}
func (this *UnreliableTransport) Send(sender common.Address, host string, port int, data []byte) error {
	time.Sleep(this.policy.Consume(1))
	drop := dummyNetwork.Counter%this.DropRate == 0
	if !drop {
		return dummyNetwork.Send(sender, host, port, data)
	} else {
		dummyNetwork.TrackSend(sender, host, port, data)
		log.Debug("dropped packet ", dummyNetwork.Counter, utils.Pex(data))
	}
	return nil
}
