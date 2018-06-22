package network

import (
	"bytes"
	"testing"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

func TestXMPPTransport(t *testing.T) {
	key1, _ := utils.MakePrivateKeyAddress()
	key2, _ := utils.MakePrivateKeyAddress()
	x1 := MakeTestXMPPTransport("x1", key1)
	x2 := MakeTestXMPPTransport("x2", key2)
	d1 := newDummyProtocol("x1")
	d2 := newDummyProtocol("x2")
	x1.RegisterProtocol(d1)
	x2.RegisterProtocol(d2)
	x1.Start()
	x2.Start()
	defer x1.Stop()
	defer x2.Stop()
	err := x1.conn.SubscribeNeighbour(x2.NodeAddress)
	if err != nil {
		t.Error(err)
		return
	}
	deviceType, isOnline := x1.NodeStatus(x2.NodeAddress)
	if deviceType != DeviceTypeOther || !isOnline {
		t.Errorf("node status error deviceType=%s,isonline=%v", deviceType, isOnline)
		return
	}
	deviceType, isOnline = x1.NodeStatus(utils.NewRandomAddress())
	if isOnline {
		t.Error("should unkown")
		return
	}
	err = x1.Send(x2.NodeAddress, []byte("abcdefg"))
	if err != nil {
		t.Error(err)
		return
	}
	select {
	case <-time.After(time.Second):
		t.Error("time out")
	case data2 := <-d2.data:
		if !bytes.Equal(data2, []byte("abcdefg")) {
			t.Error("not equal")
		}
	}
}
