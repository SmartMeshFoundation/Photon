package network

import (
	"net"
	"sync"
)

//SafeUDPConnection a udp connection with lock
type SafeUDPConnection struct {
	*net.UDPConn
	lock sync.Mutex
}

//NewSafeUDPConnection create udp connection
func NewSafeUDPConnection(protocol string, laddr *net.UDPAddr) (*SafeUDPConnection, error) {
	suc := new(SafeUDPConnection)
	var err error
	suc.UDPConn, err = net.ListenUDP(protocol, laddr)
	return suc, err
}

//WriteTo only writeto needs protection
func (su *SafeUDPConnection) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	su.lock.Lock()
	defer su.lock.Unlock()
	return su.UDPConn.WriteTo(b, addr)
}
