package network

import (
	"net"

	"fmt"

	"time"

	"errors"

	"github.com/ethereum/go-ethereum/log"
	stun "github.com/nkbai/go-stun/stun"
	"github.com/prestonTao/upnp"
)

type PortMappedSocket struct {
	Conn         *SafeUdpConnection
	Ip           string
	Port         int
	ExternalIp   string
	ExternalPort int
	Method       string
}

func OpenBareSocket(ip string, port int) (*SafeUdpConnection, error) {
	return NewSafeUdpConnection("udp", &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	})
}

func UpnpMapping(ip string, port int) (pms *PortMappedSocket, err error) {
	upnpMan := new(upnp.Upnp)
	err = upnpMan.SearchGateway()
	if err != nil {
		return
	}
	if err = upnpMan.AddPortMapping(port, port, "UDP"); err == nil {
		externalIp := upnpMan.GatewayOutsideIP
		if externalIp == "" { //multi nat routers
			err = errors.New("no outside ip")
			return
		}
	}
	conn, err := OpenBareSocket(ip, port)
	if err != nil {
		return
	}
	pms = &PortMappedSocket{
		Conn:         conn,
		Ip:           ip,
		Port:         port,
		ExternalIp:   upnpMan.GatewayOutsideIP,
		ExternalPort: port,
		Method:       "upnp",
	}
	return
}

func StunMapping(ip string, port int) (pms *PortMappedSocket, err error) {
	conn, err := OpenBareSocket(ip, port)
	if err != nil {
		return
	}
	c := stun.NewClientWithConnection(conn)
	c.SetVerbose(false)
	c.SetVVerbose(false)
	nattype, host, err := c.Discover()
	if err != nil {
		return
	}
	//disable timeout
	conn.SetDeadline(time.Now().Add(time.Hour * 24 * 365 * 10))
	log.Info(fmt.Sprintf("stun nattype:%s", nattype.String()))
	pms = &PortMappedSocket{
		Conn:         conn,
		Ip:           ip,
		Port:         port,
		ExternalIp:   host.IP(),
		ExternalPort: int(host.Port()),
		Method:       "stun",
	}
	go func() {
		for {
			err := c.KeepaliveOnlySend()
			if err != nil {
				log.Info("stun keep alive error:", err)
			}
			time.Sleep(time.Second * 30)
		}
	}()
	return
}

func noneMaping(ip string, port int) (pms *PortMappedSocket, err error) {
	conn, err := OpenBareSocket(ip, port)
	if err != nil {
		return
	}
	pms = &PortMappedSocket{
		Conn:         conn,
		Ip:           ip,
		Port:         port,
		ExternalIp:   ip,
		ExternalPort: port,
		Method:       "none",
	}
	return
}
func SocketFactory(ip string, port int, strategy string) (pms *PortMappedSocket, err error) {
	switch strategy {
	case "upnp":
		pms, err = UpnpMapping(ip, port)
	case "stun":
		pms, err = StunMapping(ip, port)
	case "auto":
		pms, err = UpnpMapping(ip, port)
		if err != nil {
			log.Info(fmt.Sprintf("upnp mapping failure err:%s", err))
		}
		pms, err = StunMapping(ip, port)
	case "none":
		pms, err = noneMaping(ip, port)
	default:
		err = errors.New("unkown  strategy")
	}
	return
}
func RelaseMappedSocket(pms *PortMappedSocket) {
	if pms.Method == "upnp" {
		upnpMan := new(upnp.Upnp)
		upnpMan.DelPortMapping(pms.ExternalPort, "UDP")
	}
	pms.Conn.Close()
}
