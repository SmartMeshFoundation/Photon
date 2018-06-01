package network

import "github.com/ethereum/go-ethereum/common"

//XMPPTransporter use XMPP to comminucate with other raiden nodes
type XMPPTransporter struct {
}

//NewXMPPTransporter create xmpp transporter
func NewXMPPTransporter() (x *XMPPTransporter, err error) {
	return &XMPPTransporter{}, nil
}

//Send a message
func (x *XMPPTransporter) Send(receiver common.Address, host string, port int, data []byte) error {
	return nil
}

//receive a message
func (x *XMPPTransporter) receive(data []byte, host string, port int) error {
	return nil
}

//Start ,ready for send and receive
func (x *XMPPTransporter) Start() {

}

//Stop send and receive
func (x *XMPPTransporter) Stop() {

}

//StopAccepting stops receiving
func (x *XMPPTransporter) StopAccepting() {

}

//RegisterProtocol a receiver
func (x *XMPPTransporter) RegisterProtocol(protcol ProtocolReceiver) {

}
