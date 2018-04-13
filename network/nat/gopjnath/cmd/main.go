package main

import (
	"time"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/gopjnath"
	"github.com/labstack/gommon/log"
)

func testCreateStream(role gopjnath.IceSessRole, neoch chan error) *gopjnath.IceStreamTransport {
	ch := make(chan error, 1)
	defer close(ch)
	stream, err := gopjnath.NewIceStreamTransport(fmt.Sprintf("test%d", role), func(u uint, bytes []byte, addr gopjnath.SockAddr) {
		log.Info(fmt.Sprintf("receive:%s", string(bytes)))
	}, func(op gopjnath.IceTransportOp, e error) {
		fmt.Printf("ice complete callback op=%d,err=%s\n", op, e)
		if op == gopjnath.IceTransportOpStateInit {
			ch <- e
		} else if op == gopjnath.IceTransportOpStateNegotiation {
			neoch <- e
		}
	})
	if err != nil {
		log.Error(err)
		return nil
	}
	err = <-ch
	if err != nil {
		log.Error("create instance error ", err)
		return nil
	}
	err = stream.InitIceSession(gopjnath.IceSessRoleControlled)
	if err != nil {
		log.Error(err)
		return nil
	}
	stream.ShowIceInfo()
	return stream
}
func testpair() {
	err := gopjnath.IceInit("182.254.155.208:3478", "182.254.155.208:3478", "bai", "bai")
	if err != nil {
		log.Error(err)
		return
	}
	stream1NeoChan := make(chan error, 1)
	stream2NeoChan := make(chan error, 1) //make sure not block sender(pjlib mainthread)
	stream := testCreateStream(gopjnath.IceSessRoleControlled, stream1NeoChan)
	if stream == nil {
		log.Error("create stream error")
		return
	}
	log.Info("start create stream2...")
	stream2 := testCreateStream(gopjnath.IceSessRoleControlling, stream2NeoChan)
	if stream2 == nil {
		log.Error("create stream2 error")
		return
	}

	sdp, err := stream.GetLocalSdp()
	if err != nil {
		log.Error("get local sdp ", err)
		return
	}

	log.Info("stream local sdp=%s\n", sdp)
	sdp2, err := stream2.GetLocalSdp()
	if err != nil {
		log.Error("get local sdp2 ", err)
	}

	err = stream.StartIce(sdp2)
	if err != nil {
		log.Error("set remote sdp err", err)
		return
	}
	err = stream2.StartIce(sdp)
	if err != nil {
		log.Error("set remote sdp err ", err)
		return
	}
	err = <-stream1NeoChan
	close(stream1NeoChan)
	log.Info("stream1NeoChan complete...")
	if err != nil {
		log.Error("stream1 neo fail ", err)
		return
	}
	err = <-stream2NeoChan
	close(stream2NeoChan)
	log.Info("stream2NeoChan complete...")
	if err != nil {
		log.Error("stream2 neo fail ", err)
		return
	}
	//time.Sleep(time.Second*10)

	err = stream.Send([]byte("data from stream1"))
	if err != nil {
		log.Error("stream send data err ", err)
		return
	}

	log.Info("stream send ...")
	err = stream2.Send([]byte("from stream2"))
	if err != nil {
		log.Error("stream2 send data err ", err)
		return
	}
	log.Info("stream2 send ...")
	time.Sleep(time.Millisecond * 10)
	//err=stream.StopIce()
	//if err!=nil{
	//	log.Error(err)
	//	return
	//}
	//err=stream.Send([]byte("after mine close"))
	//if err==nil{
	//	t.Error("should not send after stop..")
	//	return
	//}
	//err=stream2.Send([]byte("after remote close"))
	//if err!=nil{
	//	t.Error(err)
	//	return
	//}
	//time.Sleep(time.Millisecond*10)
	stream.Destroy()
	err = stream2.Send([]byte("after remote destroy destroy"))
	if err != nil {
		log.Error(err)
		return
	}
	time.Sleep(time.Millisecond * 10)
	stream2.Destroy()
}
func main() {
	log.Error("err log ")
	for i := 0; i < 1000; i++ {
		testpair()
	}
}
