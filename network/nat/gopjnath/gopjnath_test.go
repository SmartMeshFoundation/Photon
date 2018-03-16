package gopjnath

import (
	"fmt"
	"testing"
	"time"
)

func TestInitIce(t *testing.T) {
	err := IceInit("182.254.155.208:3478", "182.254.155.208:3478", "bai", "bai")
	if err != nil {
		t.Error(err)
	}
}
func testCreateStream(t *testing.T, role IceSessRole, neoch chan error) *IceStreamTransport {
	ch := make(chan error, 1)
	defer close(ch)
	var stream *IceStreamTransport
	stream, err := NewIceStreamTransport(fmt.Sprintf("test%d", role), func(u uint, bytes []byte, addr SockAddr) {
		t.Log(fmt.Sprintf("receive:%s", string(bytes)))
	}, func(op IceTransportOp, e error) {
		fmt.Printf("%s ice complete callback op=%d,err=%s\n", stream.name, op, e)
		if op == IceTransportOpStateInit {
			ch <- e
		} else if op == IceTransportOpStateNegotiation {
			neoch <- e
		}
	})
	if err != nil {
		t.Error(err)
		return nil
	}
	err = <-ch
	if err != nil {
		t.Error("create instance error ", err)
		return nil
	}
	err = stream.InitIceSession(IceSessRoleControlled)
	if err != nil {
		t.Error(err)
		return nil
	}
	stream.ShowIceInfo()
	return stream
}
func TestIcestream(t *testing.T) {
	err := IceInit("182.254.155.208:3478", "182.254.155.208:3478", "bai", "bai")
	if err != nil {
		t.Error(err)
		return
	}
	stream1NeoChan := make(chan error, 1)
	stream2NeoChan := make(chan error, 1) //make sure not block sender(pjlib mainthread)
	stream := testCreateStream(t, IceSessRoleControlled, stream1NeoChan)
	if stream == nil {
		t.Error("create stream error")
		return
	}
	t.Log("start create stream2...")
	stream2 := testCreateStream(t, IceSessRoleControlling, stream2NeoChan)
	if stream2 == nil {
		t.Error("create stream2 error")
		return
	}

	sdp, err := stream.GetLocalSdp()
	if err != nil {
		t.Error("get local sdp ", err)
		return
	}

	t.Log("stream local sdp=", sdp)
	sdp2, err := stream2.GetLocalSdp()
	if err != nil {
		t.Error("get local sdp2 ", err)
	}

	err = stream.StartIce(sdp2)
	if err != nil {
		t.Error("set remote sdp err", err)
		return
	}
	err = stream2.StartIce(sdp)
	if err != nil {
		t.Error("set remote sdp err ", err)
		return
	}
	err = <-stream1NeoChan
	close(stream1NeoChan)
	t.Log("stream1NeoChan complete...")
	if err != nil {
		t.Error("stream1 neo fail ", err)
		return
	}
	err = <-stream2NeoChan
	close(stream2NeoChan)
	t.Log("stream2NeoChan complete...")
	if err != nil {
		t.Error("stream2 neo fail ", err)
		return
	}
	//time.Sleep(time.Second*10)

	err = stream.Send([]byte("data from stream1"))
	if err != nil {
		t.Error("stream send data err ", err)
		return
	}

	t.Log("stream send ...")
	err = stream2.Send([]byte("from stream2"))
	if err != nil {
		t.Error("stream2 send data err ", err)
		return
	}
	t.Log("stream2 send ...")
	time.Sleep(time.Millisecond * 10)
	err = stream.StopIce()
	if err != nil {
		t.Error(err)
		return
	}
	err = stream.Send([]byte("after mine close"))
	if err == nil {
		t.Error("should not send after stop..")
		return
	}
	//err=stream2.Send([]byte("after remote close"))
	//if err!=nil{
	//	t.Error(err)
	//	return
	//}
	time.Sleep(time.Millisecond * 10)
	stream.Destroy()
	err = stream2.Send([]byte("after remote destroy destroy"))
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(time.Millisecond * 10)
	stream2.Destroy()
}
func TestCrash(t *testing.T) {

	for i := 0; i < 1000; i++ {
		fmt.Println("i=", i)
		TestIcestream(t)
	}
}
