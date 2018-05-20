package ice

import (
	"fmt"
	"net"
	"strconv"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
)

func addrToUDPAddr(addr string) *net.UDPAddr {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		log.Error(fmt.Sprintf("SplitHostPort %s err %s", addr, err))
	}
	porti, err := strconv.Atoi(port)
	if err != nil {
		log.Error(fmt.Sprintf("port %s not int ,err %s", port, err))
	}
	return &net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: porti,
	}
}
func udpAddrToAddr(udpAddr net.Addr) string {
	addr := udpAddr.(*net.UDPAddr)
	return fmt.Sprintf("%s:%d", addr.IP.String(), addr.Port)
}
