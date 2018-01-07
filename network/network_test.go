package network

import (
	"fmt"
	"testing"

	"time"

	"github.com/SmartMeshFoundation/raiden-network/params"
	"github.com/prestonTao/upnp"
)

func TestUpnp(t *testing.T) {
	upnpMan := new(upnp.Upnp)
	err := upnpMan.SearchGateway()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("local ip address: ", upnpMan.LocalHost)
		fmt.Println("gateway ip address: ", upnpMan.Gateway.Host)

	}
	mapping := new(upnp.Upnp)
	if err := mapping.AddPortMapping(55789, 55789, "UDP"); err == nil {
		fmt.Println("success !")
		fmt.Println("extern ip:", upnpMan.GatewayOutsideIP)
		// remove Port mapping in gatway
		mapping.Reclaim()
	} else {
		fmt.Println("fail !")
	}
}

func TestStun(t *testing.T) {
	pms, err := StunMapping("0.0.0.0", params.INITIAL_PORT)
	if err != nil { //stunmap should always work
		t.Error(err)
	}
	t.Logf("pms=%#v", pms)
	pms.Conn.Close()
}

func TestSocketFactory(t *testing.T) {
	pms, err := SocketFactory("0.0.0.0", params.INITIAL_PORT, "stun")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("pms=%#v", pms)
	RelaseMappedSocket(pms)
	time.Sleep(2 * time.Second)
}
