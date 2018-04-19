package ice

import (
	"fmt"
	"net"
	"strconv"
)

func addrToUdpAddr(addr string) *net.UDPAddr {
	host, port, _ := net.SplitHostPort(addr)
	porti, _ := strconv.Atoi(port)
	return &net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: porti,
	}
}
func udpAddrToAddr(udpAddr net.Addr) string {
	addr := udpAddr.(*net.UDPAddr)
	return fmt.Sprintf("%s:%d", addr.IP.String(), addr.Port)
}
