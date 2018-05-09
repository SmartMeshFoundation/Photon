package network

import (
	"net"

	"fmt"

	"time"

	"errors"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/stun"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
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
	var err2 error
	laddr := &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	}
	rhost, rport := SplitHostPort(params.DefaultTurnServer)
	raddr := &net.UDPAddr{
		IP:   net.ParseIP(rhost),
		Port: rport,
	}
	conn, err := net.DialUDP("udp", laddr, raddr)
	c, err := stun.NewClient(stun.ClientOptions{
		Connection: conn,
	})
	if err != nil {
		return
	}
	pms = &PortMappedSocket{
		Ip:     ip,
		Port:   port,
		Method: "stun",
	}
	deadline := time.Now().Add(time.Second * 25)
	err = c.Do(stun.MustBuild(stun.TransactionIDSetter, stun.BindingRequest), deadline, func(res stun.Event) {
		if res.Error != nil {
			err2 = fmt.Errorf("res %s", res)
			return
		}
		var xorAddr stun.XORMappedAddress
		if err2 = xorAddr.GetFrom(res.Message); err2 != nil {
			var addr stun.MappedAddress
			err2 = addr.GetFrom(res.Message)
			if err2 != nil {
				return
			}
			log.Info(fmt.Sprintf("addr=%s", addr))
			pms.ExternalIp = addr.IP.String()
			pms.ExternalPort = addr.Port
		} else {
			pms.ExternalIp = xorAddr.IP.String()
			pms.ExternalPort = xorAddr.Port
			log.Info(fmt.Sprintf("xoraddr=%s", xorAddr))
		}
	})
	if err != nil {
		log.Crit(fmt.Sprintf("StunMapping do: %s", err))
	}
	if err2 != nil {
		log.Error(fmt.Sprintf("get external ip err %s", err2))
		err = err2
		return
	}
	err = c.Close()
	if err != nil {
		log.Crit(err.Error())
	}
	/*
		get our extern ip,port ,then restart
	*/
	conn2, err := OpenBareSocket(ip, port)
	if err != nil {
		return
	}
	pms.Conn = conn2
	go func() {

		for {
			req, _ := stun.Build(stun.TransactionIDSetter, stun.BindingIndication)
			_, err := conn.WriteTo(req.Raw, raddr)
			if err != nil {
				log.Info("stun keep alive err %s", err)
			}
			time.Sleep(time.Second * 30)
		}
	}()
	return
}

//func StunMapping2(ip string, port int) (pms *PortMappedSocket, err error) {
//	conn, err := OpenBareSocket(ip, port)
//	if err != nil {
//		return
//	}
//	c := stun.NewClientWithConnection(conn)
//	c.SetVerbose(false)
//	c.SetVVerbose(false)
//	//c.SetServerHost("stunserver.org", 3478)
//	//c.SetServerHost("182.254.155.208", 3478)
//	nattype, host, err := c.Discover()
//	if err != nil {
//		return
//	}
//	//disable timeout
//	conn.SetDeadline(time.Time{})
//	log.Info(fmt.Sprintf("stun nattype:%s", nattype.String()))
//	pms = &PortMappedSocket{
//		Conn:         conn,
//		Ip:           ip,
//		Port:         port,
//		ExternalIp:   host.IP(),
//		ExternalPort: int(host.Port()),
//		Method:       "stun",
//	}
//	go func() {
//		for {
//			err := c.KeepaliveOnlySend()
//			if err != nil {
//				log.Info("stun keep alive error:", err)
//			}
//			time.Sleep(time.Second * 30)
//		}
//	}()
//	return
//}

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
