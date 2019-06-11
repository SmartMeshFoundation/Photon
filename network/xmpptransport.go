package network

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"time"

	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/network/netshare"
	"github.com/SmartMeshFoundation/Photon/network/xmpptransport"
	"github.com/SmartMeshFoundation/Photon/network/xmpptransport/xmpppass"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var errXMPPConnectionNotReady = errors.New("xmpp connection not ready")

//XMPPTransport use XMPP to comminucate with other Photon nodes
type XMPPTransport struct {
	conn          *xmpptransport.XMPPConnection
	quitChan      chan struct{}
	stopped       bool
	stopReceiving bool
	log           log.Logger
	protocol      ProtocolReceiver
	NodeAddress   common.Address
	key           *ecdsa.PrivateKey
	statusChan    chan netshare.Status
}

/*
NewXMPPTransport create xmpp transporter,
if not success ,for example cannot connect to xmpp server, will try background
*/
func NewXMPPTransport(name, ServerURL string, key *ecdsa.PrivateKey, deviceType string) (x *XMPPTransport) {
	x = &XMPPTransport{
		quitChan:    make(chan struct{}),
		NodeAddress: crypto.PubkeyToAddress(key.PublicKey),
		key:         key,
		statusChan:  make(chan netshare.Status, 10),
	}
	addr := crypto.PubkeyToAddress(key.PublicKey)
	x.log = log.New("name", name)
	// 2019.06.10 启动时主线程不等待,优化无网情况下的启动速度
	// 就算如果matrix连接不上,而主线程正常启动完成开始发送消息,也会被Send方法拒绝,重要消息会进入重发阶段,没有影响
	//wg := sync.WaitGroup{}
	//wg.Add(1)
	go func() {
		wait := time.Millisecond
		var err error
		//only wait one time
		//var first bool
		for {
			select {
			case <-time.After(wait):
				x.conn, err = xmpptransport.NewConnection(ServerURL, addr, x, x, name, deviceType, x.statusChan)
				//if !first {
				//	first = true
				//	wg.Done()
				//}
				if err != nil {
					x.log.Error(fmt.Sprintf("cannot connect to xmpp server %s, retry in 5 seconds", ServerURL))
					time.Sleep(time.Second * 5)
				} else {
					return
				}
			case <-x.quitChan:
				return
			}

		}

	}()
	//only wait for one try
	//wg.Wait()
	return
}

//GetPassWord returns current login password
func (x *XMPPTransport) GetPassWord() string {
	pass, err := xmpppass.CreatePassword(x.key)
	if err != nil {
		log.Error(fmt.Sprintf("GetPassWord for %s err %s", utils.APex2(x.NodeAddress), err))
	}
	err = xmpppass.VerifySignature(x.NodeAddress.String(), pass)
	if err != nil {
		panic(err)
	}
	return pass
}

//Send a message
func (x *XMPPTransport) Send(receiver common.Address, data []byte) error {
	x.log.Trace(fmt.Sprintf("send to %s, message=%s", utils.APex2(receiver), encoding.MessageType(data[0])))
	if x.stopped || x.conn == nil {
		return errXMPPConnectionNotReady
	}
	return x.conn.SendData(receiver, data)
}

//DataHandler call back of xmpp connection
func (x *XMPPTransport) DataHandler(from common.Address, data []byte) {
	x.log.Trace(fmt.Sprintf("received from %s, message=%s", utils.APex2(from), encoding.MessageType(data[0])))
	if x.stopped || x.stopReceiving {
		return
	}
	if x.protocol != nil {
		x.protocol.receive(data)
	}
}

//Start ,ready for send and receive
func (x *XMPPTransport) Start() {

}

//Stop send and receive
func (x *XMPPTransport) Stop() {
	x.stopped = true
	close(x.quitChan)
	if x.conn != nil {
		x.conn.Close()
	}
}

//StopAccepting stops receiving
func (x *XMPPTransport) StopAccepting() {
	x.stopReceiving = true
}

//RegisterProtocol a receiver
func (x *XMPPTransport) RegisterProtocol(protcol ProtocolReceiver) {
	x.protocol = protcol
}

//NodeStatus get node's status and is online right now
func (x *XMPPTransport) NodeStatus(addr common.Address) (deviceType string, isOnline bool) {
	if x.conn == nil {
		return "", false
	}
	var err error
	deviceType, isOnline, err = x.conn.IsNodeOnline(addr)
	if err != nil {
		x.log.Error(fmt.Sprintf("IsNodeOnline query %s err %s", utils.APex2(addr), err))
	}
	return
}

// RegisterWakeUpChan impl wakeuphandler.IWakeUpHandler 由于xmpp的节点在线状态维护在连接层,所以在这里转发下
func (x *XMPPTransport) RegisterWakeUpChan(addr common.Address, c chan int) {
	x.conn.RegisterWakeUpChan(addr, c)
}

// UnRegisterWakeUpChan impl wakeuphandler.IWakeUpHandler 由于xmpp的节点在线状态维护在连接层,所以在这里转发下
func (x *XMPPTransport) UnRegisterWakeUpChan(addr common.Address) {
	x.conn.UnRegisterWakeUpChan(addr)
}

// WakeUp impl wakeuphandler.IWakeUpHandler, shouldn't call
func (x *XMPPTransport) WakeUp(addr common.Address) {
	panic("wrong call")
}
