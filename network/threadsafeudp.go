package network

import (
	"net"
)

//NewSafeUDPConnection create udp connection
func NewSafeUDPConnection(protocol string, laddr *net.UDPAddr) (*net.UDPConn, error) {

	var err error
	addr2 := *laddr
	addr2.IP = net.ParseIP("0.0.0.0") //确保listen的是0.0.0.0
	conn, err := net.ListenUDP(protocol, &addr2)
	return conn, err
}
