package network

import (
	"net"
	"sync"
)

type SafeUdpConnection struct {
	*net.UDPConn
	lock sync.Mutex
}

func NewSafeUdpConnection(protocol string, laddr *net.UDPAddr) (*SafeUdpConnection, error) {
	suc := new(SafeUdpConnection)
	var err error
	suc.UDPConn, err = net.ListenUDP(protocol, laddr)
	return suc, err
}

//only writeto needs protection
func (this *SafeUdpConnection) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.UDPConn.WriteTo(b, addr)
}
