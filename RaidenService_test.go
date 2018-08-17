package smartraiden

import (
	"os"
	"testing"

	"time"

	"fmt"

	"runtime"

	_ "net/http/pprof"

	"net/http"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
}
func TestPing(t *testing.T) {
	reinit()
	r1 := newTestRaiden()
	r2 := newTestRaiden()
	r1.Start()
	r2.Start()
	defer r1.Stop()
	defer r2.Stop()
	ping := encoding.NewPing(32)
	ping.Sign(r1.PrivateKey, ping)
	err := r1.SendAndWait(r2.NodeAddress, ping, time.Second*10)
	if err != nil {
		t.Error(err)
	}
}

func TestRestart(t *testing.T) {
	var err error
	for i := 0; i < 10; i++ {
		reinit()
		r1 := newTestRaiden()
		log.Info(fmt.Sprintf("%d r1 create success", i))
		r2 := newTestRaiden()
		log.Info(fmt.Sprintf("%d r2 create success", i))
		err = r1.Start()
		if err != nil {
			t.Error(err)
			return
		}
		log.Info(fmt.Sprintf("%d r1 start success", i))
		err = r2.Start()
		if err != nil {
			t.Error(err)
			return
		}
		log.Info(fmt.Sprintf("%d r2 create success", i))
		ping := encoding.NewPing(32)
		ping.Sign(r1.PrivateKey, ping)
		err := r1.SendAndWait(r2.NodeAddress, ping, time.Second*10)
		if err != nil {
			t.Error(err)
			return
		}
		r1.Stop()
		log.Info(fmt.Sprintf("%d r1 create success", i))
		r2.Stop()
		log.Info(fmt.Sprintf("%d r2 create success", i))
		time.Sleep(time.Second * 5)
		runtime.GC()
	}

}
func TestRestart2(t *testing.T) {
	go http.ListenAndServe("0.0.0.0:6060", nil)
	for i := 0; i < 10; i++ {
		var err error
		reinit()
		r1 := newTestRaiden()
		log.Info(fmt.Sprintf("%d r1 create success", i))
		err = r1.Start()
		if err != nil {
			t.Error(err)
			return
		}
		log.Info(fmt.Sprintf("%d r1 start success", i))
		r1.Stop()
		log.Info(fmt.Sprintf("%d r1 stop success", i))
		runtime.GC()
		time.Sleep(time.Second * 10)
	}

}

func TestRestart3(t *testing.T) {
	go http.ListenAndServe("0.0.0.0:6060", nil)
	for i := 0; i < 10; i++ {
		var err error
		reinit()
		r1 := newTestRaiden()
		r2 := newTestRaiden()
		log.Info(fmt.Sprintf("%d r1 create success", i))
		err = r1.Start()
		if err != nil {
			t.Error(err)
			return
		}
		log.Info(fmt.Sprintf("%d r1 start success", i))
		err = r2.Start()
		if err != nil {
			t.Error(err)
			return
		}
		log.Info(fmt.Sprintf("%d r2 start success", i))
		ping := encoding.NewPing(int64(i + 1))
		err = ping.Sign(r1.PrivateKey, ping)
		if err != nil {
			t.Error(err)
			return
		}
		err = r1.SendAndWait(r2.NodeAddress, ping, time.Second*10)
		if err != nil {
			t.Error(err)
			return
		}
		log.Info(fmt.Sprintf("%d send ping message success", i))
		r1.Stop()
		log.Info(fmt.Sprintf("%d r1 stop success", i))
		r2.Stop()
		log.Info(fmt.Sprintf("%d r2 stop success", i))
		runtime.GC()
		time.Sleep(time.Second * 10)
	}

}

func TestTimer(t *testing.T) {
	tm := time.NewTimer(time.Second)
	//<-tm.C
	fmt.Println("fired")
	fmt.Printf("called stop %v\n", tm.Stop())
	fmt.Printf("called stop2 %v\n", tm.Stop())
	<-tm.C
	fmt.Println("fired after stop")
}
