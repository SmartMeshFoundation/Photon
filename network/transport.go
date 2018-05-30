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

//Policier to control the sending speed of transporter
type Policier interface {
	//Consume tokens.
	//Args:
	//tokens (float): number of transport tokens to consume
	//Returns:
	//wait_time (float): waiting time for the consumer
	Consume(tokens float64) time.Duration
}

//Transporter denotes a communication transport used by protocol
type Transporter interface {
	//Send a message
	Send(receiver common.Address, host string, port int, data []byte) error
	//receive a message
	receive(data []byte, host string, port int) error
	//Start ,ready for send and receive
	Start()
	//Stop send and receive
	Stop()
	//StopAccepting stops receiving
	StopAccepting()
	//RegisterProtocol a receiver
	RegisterProtocol(protcol ProtocolReceiver)
}
type messageCallBack func(sender common.Address, hostport string, msg []byte)

func tohostport(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
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

//dummyNetwork Store global state for an in process network, this won't use a real
//network protocol
type dummyNetwork struct {
	Transports              map[string]Transporter
	Counter                 int
	MessageSendCallbacks    []messageCallBack
	MessageReceiveCallbacks []messageCallBack
}

func newDummyNetwork() *dummyNetwork {
	return &dummyNetwork{
		Transports: make(map[string]Transporter),
		Counter:    0,
	}
}

var dummyNet = newDummyNetwork()

//RegisterSendCallback register callback
func RegisterSendCallback(cb messageCallBack) {
	dummyNet.MessageSendCallbacks = append(dummyNet.MessageSendCallbacks, cb)
}

//RegisterReceiveCallback register callback
func RegisterReceiveCallback(cb messageCallBack) {
	dummyNet.MessageReceiveCallbacks = append(dummyNet.MessageReceiveCallbacks, cb)
}

//Register a new node in the dummy network.
func (dn *dummyNetwork) Register(transpoter Transporter, host string, port int) {
	hostport := fmt.Sprintf("%s:%d", host, port)
	dn.Transports[hostport] = transpoter
}

//Register an attempt to send a packet. This method should be called
//everytime send() is used.
func (dn *dummyNetwork) trackSend(receiver common.Address, host string, port int, data []byte) error {
	dn.Counter++
	for _, cb := range dn.MessageSendCallbacks {
		cb(receiver, tohostport(host, port), data)
	}
	return nil
}

func (dn *dummyNetwork) trackReceive(receiver common.Address, host string, port int, data []byte) {
	for _, cb := range dn.MessageReceiveCallbacks {
		cb(receiver, tohostport(host, port), data)
	}
}

//Send a message
func (dn *dummyNetwork) Send(sender common.Address, host string, port int, data []byte) error {
	dn.trackSend(sender, host, port, data)
	hostport := tohostport(host, port)
	time.AfterFunc(time.Nanosecond, func() {
		dn.Transports[hostport].receive(data, host, port)
	})
	return nil
}

//ProtocolReceiver receive
type ProtocolReceiver interface {
	receive(data []byte, host string, port int)
}

//UDPTransport represents a UDP connection
type UDPTransport struct {
	protocol      ProtocolReceiver
	conn          *SafeUDPConnection
	Host          string
	Port          int
	policy        Policier
	isClosed      bool
	stopReceiving bool //todo use atomic to replace
}

//NewUDPTransport create UDPTransport
func NewUDPTransport(host string, port int, conn *SafeUDPConnection, protocol ProtocolReceiver, policy Policier) (t *UDPTransport, err error) {
	t = &UDPTransport{
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
	if conn == nil {
		conn, err = NewSafeUDPConnection("udp", addr)
		if err != nil {
			err = fmt.Errorf("listen udp %s:%d error %v", host, port, err)
			return
		}
	}
	t.conn = conn
	log.Trace(fmt.Sprintf("listen udp on %s:%d", host, port))
	return
}
func newUDPTransportWithHostPort(host string, port int, protocol ProtocolReceiver, policy Policier) *UDPTransport {
	t, err := NewUDPTransport(host, port, nil, protocol, policy)
	if err != nil {
		log.Error(err.Error())
	}
	return t
}

//Start udp listening
func (ut *UDPTransport) Start() {
	t := ut
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
			log.Trace(fmt.Sprintf("%d receive from %s:%d,message=%s,hash=%s", t.Port, remoteAddr.IP.String(),
				remoteAddr.Port, encoding.MessageType(data[0]), utils.HPex(utils.Sha3(data[:read]))))
			t.receive(data[:read], remoteAddr.IP.String(), remoteAddr.Port)
		}

	}()
}
func (ut *UDPTransport) receive(data []byte, host string, port int) error {
	//todo fix get raiden address, my node address
	dummyNet.trackReceive(common.Address{}, host, port, data)
	if ut.protocol != nil { //receive data before register a protocol
		ut.protocol.receive(data, host, port)
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
func (ut *UDPTransport) Send(receiver common.Address, host string, port int, data []byte) error {
	dummyNet.trackSend(receiver, host, port, data)
	log.Trace(fmt.Sprintf("%d send to %s %s:%d, message=%s,response hash=%s", ut.Port, utils.APex2(receiver), host, port, encoding.MessageType(data[0]), utils.HPex(utils.Sha3(data, receiver[:]))))
	//only comment this line,if you want to test.
	//time.Sleep(ut.policy.Consume(1)) //force to wait,
	//todo need one lock for write?
	_, err := ut.conn.WriteToUDP(data, udpAddrFromHostport(host, port))
	if err != nil {
		return err
	}
	return nil
}

//RegisterProtocol register receiver
func (ut *UDPTransport) RegisterProtocol(proto ProtocolReceiver) {
	ut.protocol = proto
}

//Stop UDP connection
func (ut *UDPTransport) Stop() {
	ut.isClosed = true
	ut.conn.Close()
}

//StopAccepting stop receiving
func (ut *UDPTransport) StopAccepting() {
	ut.stopReceiving = true
}

// Communication between inter-process nodes.
type dummyTransport struct {
	protocol ProtocolReceiver
	host     string
	port     int
	policy   Policier
}

//NewDummyTransport create a dummy transporter
func NewDummyTransport(host string, port int, protocol ProtocolReceiver, policy Policier) Transporter {
	t := &dummyTransport{
		protocol: protocol,
		host:     host,
		port:     port,
		policy:   policy,
	}
	dummyNet.Register(t, host, port)
	return t
}

//Send a message
func (dt *dummyTransport) Send(receiver common.Address, host string, port int, data []byte) error {
	time.Sleep(dt.policy.Consume(1))
	return dummyNet.Send(receiver, host, port, data)
}
func (dt *dummyTransport) receive(data []byte, host string, port int) error {
	dummyNet.trackReceive(common.Address{}, host, port, data)
	dt.protocol.receive(data, host, port)
	return nil
}

//RegisterProtocol a callback
func (dt *dummyTransport) RegisterProtocol(protcol ProtocolReceiver) {
	dt.protocol = protcol
}

//Start dummy
func (dt *dummyTransport) Start() {

}

//Stop dummy
func (dt *dummyTransport) Stop() {

}

//StopAccepting dummy
func (dt *dummyTransport) StopAccepting() {

}

type unreliableTransport struct {
	dummyTransport
	DropRate int
}

func newUnreliableTransport(t *dummyTransport) *unreliableTransport {
	return &unreliableTransport{dummyTransport: *t, DropRate: 2}
}

//Send a message ,it drops message randomly.
func (ut *unreliableTransport) Send(sender common.Address, host string, port int, data []byte) error {
	time.Sleep(ut.policy.Consume(1))
	drop := dummyNet.Counter%ut.DropRate == 0
	if !drop {
		return dummyNet.Send(sender, host, port, data)
	}
	dummyNet.trackSend(sender, host, port, data)
	log.Debug("dropped packet ", dummyNet.Counter, utils.Pex(data))
	return nil
}
