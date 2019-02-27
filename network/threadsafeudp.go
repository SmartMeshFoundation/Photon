package network

import (
	"net"
)

//SafeUDPConnection a udp connection with lock
type SafeUDPConnection struct {
	*net.UDPConn
	//lock sync.Mutex
}

//NewSafeUDPConnection create udp connection
func NewSafeUDPConnection(protocol string, laddr *net.UDPAddr) (*SafeUDPConnection, error) {
	suc := new(SafeUDPConnection)
	var err error
	addr2 := *laddr
	addr2.IP = net.ParseIP("0.0.0.0") //确保listen的是0.0.0.0
	suc.UDPConn, err = net.ListenUDP(protocol, &addr2)
	return suc, err
}

//WriteTo only writeto needs protection
func (su *SafeUDPConnection) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	//su.lock.Lock()
	//defer su.lock.Unlock()
	return su.UDPConn.WriteTo(b, addr)
}
