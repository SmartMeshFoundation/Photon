package ice

import (
	"testing"

	"net"

	"time"

	"bytes"

	"fmt"

	"github.com/nkbai/log"
)

const (
	typHost = 1
	typStun = 2
	typTurn = 3
)

type icecb struct {
	data      chan []byte
	iceresult chan error
	name      string
}

func Newicecb(name string) *icecb {
	return &icecb{
		name:      name,
		data:      make(chan []byte, 1),
		iceresult: make(chan error, 1),
	}
}
func (c *icecb) OnReceiveData(data []byte, from net.Addr) {
	c.data <- data
}

/*
	Callback to report status of various ICE operations.
*/
func (c *icecb) OnIceComplete(result error) {
	c.iceresult <- result
	log.Trace("%s negotiation complete", c.name)
}
func setupTestIceStreamTransport(typ int) (s1, s2 *IceStreamTransport, err error) {
	var cfg *TransportConfig
	switch typ {
	case typHost:
		cfg = NewTransportConfigHostonly()
	case typStun:
		cfg = NewTransportConfigWithStun("182.254.155.208:3478")
	case typTurn:
		cfg = NewTransportConfigWithTurn("182.254.155.208:3478", "bai", "bai")
	}
	s1, err = NewIceStreamTransport(cfg, "s1")
	if err != nil {
		return
	}
	s2, err = NewIceStreamTransport(cfg, "s2")
	log.Trace("-----------------------------------------")
	return
}
func TestNewIceStreamTransport(t *testing.T) {
	cfg := NewTransportConfigHostonly()
	trans, err := NewIceStreamTransport(cfg, "hostonly")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("candidates host only:", log.StringInterface(trans.component.candidates, 3))
	cfg = NewTransportConfigWithStun("182.254.155.208:3478")
	trans, err = NewIceStreamTransport(cfg, "stun")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("candidates stun: ", log.StringInterface(trans.component.candidates, 3))
	cfg = NewTransportConfigWithTurn("182.254.155.208:3478", "bai", "bai")
	trans, err = NewIceStreamTransport(cfg, "turn")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("candidates turn:", log.StringInterface(trans.component.candidates, 3))
	trans.InitIce(SessionRoleControlling)
	s, err := trans.EncodeSession()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(s)
}

func TestIceStreamDecodeSession(t *testing.T) {
	s := `
v=0
o=- 3414953978 3414953978 IN IP4 localhost
s=ice
t=0 0
a=ice-ufrag:088e4954
a=ice-pwd:35702e2f
m=audio 52628 RTP/AVP 0
c=IN IP4 182.254.155.208
a=candidate:Sac140a06 1 UDP 1694498815 123.147.248.122 56465 typ srflx
a=candidate:Hac140a06 1 UDP 2130706431 172.20.10.6 59951 typ host
a=candidate:Ha000011 1 UDP 2130706431 10.0.0.17 59951 typ host
a=candidate:Had33702 1 UDP 2130706431 10.211.55.2 59951 typ host
a=candidate:Ha258102 1 UDP 2130706431 10.37.129.2 59951 typ host
a=candidate:Rb6fe9bd0 1 UDP 16777215 182.254.155.208 52628 typ relay
`
	session, err := DecodeSession(s)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("session=%s", log.StringInterface(session, 3))
}

func TestIceStreamTransport_StartNegotiation(t *testing.T) {
	s1, s2, err := setupTestIceStreamTransport(typHost)
	if err != nil {
		t.Error(err)
		return
	}
	cb1 := Newicecb("s1")
	cb2 := Newicecb("s2")
	s1.cb = cb1
	s2.cb = cb2
	err = s1.InitIce(SessionRoleControlling)
	if err != nil {
		t.Error(err)
		return
	}
	err = s2.InitIce(SessionRoleControlled)
	if err != nil {
		t.Error(err)
		return
	}
	lsdp, err := s1.EncodeSession()
	log.Trace("lsdp=%s", lsdp)
	if err != nil {
		t.Error(err)
		return
	}
	log.Trace("sdp length=%s", len(lsdp))
	rsdp, err := s2.EncodeSession()
	if err != nil {
		t.Error(err)
		return
	}
	log.Trace("rsdp=%s", rsdp)

	err = s2.StartNegotiation(lsdp)
	if err != nil {
		t.Error(err)
		return
	}

	err = s1.StartNegotiation(rsdp)
	if err != nil {
		t.Error(err)
		return
	}
	select {
	case <-time.After(20 * time.Second):
		t.Error("s2 negotiation timeout")
		return
	case err = <-cb2.iceresult:
		if err != nil {
			t.Error("s2 negotiation failed", err)
			return
		}
	}
	//return
	select {
	case <-time.After(20 * time.Second):
		t.Error("s1 negotiation timeout")
		return
	case err = <-cb1.iceresult:
		if err != nil {
			t.Error("s1 negotiation failed ", err)
			return
		}
	}
	s1data := []byte("hello,s2")
	s2data := []byte("hello,s1")
	err = s1.SendData(s1data)
	if err != nil {
		t.Error(err)
		return
	}
	err = s2.SendData(s2data)
	if err != nil {
		t.Error(err)
		return
	}
	select {
	case <-time.After(10 * time.Second):
		t.Error("s2 recevied timeout")
		return
	case data := <-cb2.data:
		if !bytes.Equal(data, s1data) {
			t.Error("s2 recevied error ,got ", string(data))
			return
		}
	}
	select {
	case <-time.After(10 * time.Second):
		t.Error("s1 recevied timeout")
		return
	case data := <-cb1.data:
		if !bytes.Equal(data, s2data) {
			t.Error("s1 recevied error ,got ", string(data))
			return
		}
	}
	return
}
func encodeSessionExclude(t *IceStreamTransport, excludes ...CandidateType) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "v=0\no=- 3414953978 3414953978 IN IP4 localhost\ns=ice\nt=0 0\n")
	fmt.Fprintf(buf, "a=ice-ufrag:%s\na=ice-pwd:%s\n", t.session.rxUserFrag, t.session.rxPassword)
	//only on component now....
	uaddr := addrToUdpAddr(t.component.defaultCandidate.addr)
	fmt.Fprintf(buf, "m=audio %d RTP/AVP 0\nc=IN IP4 %s\n", uaddr.Port, uaddr.IP.String())
	for _, c := range t.component.candidates {
		found := false
		for _, t := range excludes {
			if c.Type == t {
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintf(buf, "%s\n", c)
		}
	}
	return string(buf.Bytes())
}
func TestIceStreamTransport_StartNegotiationOnlyRelay(t *testing.T) {
	s1, s2, err := setupTestIceStreamTransport(typTurn)
	if err != nil {
		t.Error(err)
		return
	}
	cb1 := Newicecb("s1")
	cb2 := Newicecb("s2")
	s1.cb = cb1
	s2.cb = cb2
	err = s1.InitIce(SessionRoleControlling)
	if err != nil {
		t.Error(err)
		return
	}
	err = s2.InitIce(SessionRoleControlled)
	if err != nil {
		t.Error(err)
		return
	}
	rsdp := encodeSessionExclude(s2, CandidateHost, CandidateServerReflexive)
	err = s1.StartNegotiation(rsdp)
	if err != nil {
		t.Error(err)
		return
	}
	lsdp := encodeSessionExclude(s1, CandidateHost, CandidateServerReflexive)
	err = s2.StartNegotiation(lsdp)
	if err != nil {
		t.Error(err)
		return
	}
	select {
	case <-time.After(50 * time.Second):
		t.Error("s1 negotiation timeout")
		return
	case err = <-cb1.iceresult:
		if err != nil {
			t.Error("s1 negotiation failed ", err)
			//return
		}
	}
	log.Info("s1 negotiation success.")
	select {
	case <-time.After(50 * time.Second):
		t.Error("s2 negotiation timeout")
		return
	case err = <-cb2.iceresult:
		if err != nil {
			t.Error("s2 negotiation failed", err)
			return
		}
	}
	log.Info("s2 negotiation success")

}

func TestIceStreamTransport_StartNegotiationNoHost(t *testing.T) {
	s1, s2, err := setupTestIceStreamTransport(typTurn)
	if err != nil {
		t.Error(err)
		return
	}
	cb1 := Newicecb("s1")
	cb2 := Newicecb("s2")
	s1.cb = cb1
	s2.cb = cb2
	err = s1.InitIce(SessionRoleControlling)
	if err != nil {
		t.Error(err)
		return
	}
	err = s2.InitIce(SessionRoleControlled)
	if err != nil {
		t.Error(err)
		return
	}
	lsdp := encodeSessionExclude(s1, CandidateHost)
	err = s2.StartNegotiation(lsdp)
	if err != nil {
		t.Error(err)
		return
	}
	log.Trace("s2 StartNegotiation returned")
	rsdp := encodeSessionExclude(s2, CandidateHost)
	err = s1.StartNegotiation(rsdp)
	if err != nil {
		t.Error(err)
		return
	}
	log.Trace("s1 StartNegotiation returned")
	select {
	case <-time.After(50 * time.Second):
		t.Error("s1 negotiation timeout")
		return
	case err = <-cb1.iceresult:
		if err != nil {
			t.Error("s1 negotiation failed ", err)
			//return
		}
	}
	log.Info("s1 negotiation success.")
	select {
	case <-time.After(50 * time.Second):
		t.Error("s2 negotiation timeout")
		return
	case err = <-cb2.iceresult:
		if err != nil {
			t.Error("s2 negotiation failed", err)
			return
		}
	}
	log.Info("s2 negotiation success")

}

func BenchmarkIceStreamTransport_StartNegotiation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s1, s2, err := setupTestIceStreamTransport(typTurn)
		if err != nil {
			log.Error(err.Error())
			return
		}
		cb1 := Newicecb("s1")
		cb2 := Newicecb("s2")
		s1.cb = cb1
		s2.cb = cb2
		err = s1.InitIce(SessionRoleControlling)
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = s2.InitIce(SessionRoleControlled)
		if err != nil {
			log.Error(err.Error())
			return
		}
		lsdp := encodeSessionExclude(s1, CandidateHost)
		err = s2.StartNegotiation(lsdp)
		if err != nil {
			log.Error(err.Error())
			return
		}
		log.Trace("s2 StartNegotiation returned")
		rsdp := encodeSessionExclude(s2, CandidateHost)
		err = s1.StartNegotiation(rsdp)
		if err != nil {
			log.Error(err.Error())
			return
		}
		log.Trace("s1 StartNegotiation returned")
		select {
		case <-time.After(50 * time.Second):
			log.Error("s1 negotiation timeout")
			return
		case err = <-cb1.iceresult:
			if err != nil {
				log.Error("s1 negotiation failed ", err)
				return
			}
		}
		log.Info("s1 negotiation success.")
		select {
		case <-time.After(50 * time.Second):
			log.Error("s2 negotiation timeout")
			return
		case err = <-cb2.iceresult:
			if err != nil {
				log.Error("s2 negotiation failed", err)
				return
			}
		}
		log.Info("s2 negotiation success")
	}

}

func BenchmarkIceStreamTransport_StartNegotiationOnlyRelay(b *testing.B) {
	for i := 0; i < b.N; i++ {

		s1, s2, err := setupTestIceStreamTransport(typTurn)
		if err != nil {
			log.Error(err.Error())
			return
		}
		cb1 := Newicecb("s1")
		cb2 := Newicecb("s2")
		s1.cb = cb1
		s2.cb = cb2
		err = s1.InitIce(SessionRoleControlling)
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = s2.InitIce(SessionRoleControlled)
		if err != nil {
			log.Error(err.Error())
			return
		}
		rsdp := encodeSessionExclude(s2, CandidateHost, CandidateServerReflexive)
		err = s1.StartNegotiation(rsdp)
		if err != nil {
			log.Error(err.Error())
			return
		}
		lsdp := encodeSessionExclude(s1, CandidateHost, CandidateServerReflexive)
		err = s2.StartNegotiation(lsdp)
		if err != nil {
			log.Error(err.Error())
			return
		}
		select {
		case <-time.After(50 * time.Second):
			log.Error("s1 negotiation timeout")
			return
		case err = <-cb1.iceresult:
			if err != nil {
				log.Error("s1 negotiation failed ", err)
				//return
			}
		}
		log.Info("s1 negotiation success.")
		select {
		case <-time.After(50 * time.Second):
			log.Error("s2 negotiation timeout")
			return
		case err = <-cb2.iceresult:
			if err != nil {
				log.Error("s2 negotiation failed", err)
				return
			}
		}
		log.Info("s2 negotiation success")
	}

}

func BenchmarkIceStreamTransport_StartNegotiationNoHost(b *testing.B) {
	for i := 0; i < b.N; i++ {

		s1, s2, err := setupTestIceStreamTransport(typTurn)
		if err != nil {
			log.Error(err.Error())
			return
		}
		cb1 := Newicecb("s1")
		cb2 := Newicecb("s2")
		s1.cb = cb1
		s2.cb = cb2
		err = s1.InitIce(SessionRoleControlling)
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = s2.InitIce(SessionRoleControlled)
		if err != nil {
			log.Error(err.Error())
			return
		}
		lsdp := encodeSessionExclude(s1, CandidateHost)
		err = s2.StartNegotiation(lsdp)
		if err != nil {
			log.Error(err.Error())
			return
		}
		log.Trace("s2 StartNegotiation returned")
		rsdp := encodeSessionExclude(s2, CandidateHost)
		err = s1.StartNegotiation(rsdp)
		if err != nil {
			log.Error(err.Error())
			return
		}
		log.Trace("s1 StartNegotiation returned")
		select {
		case <-time.After(50 * time.Second):
			log.Error("s1 negotiation timeout")
			return
		case err = <-cb1.iceresult:
			if err != nil {
				log.Error("s1 negotiation failed ", err)
				return
			}
		}
		log.Info("s1 negotiation success.")
		select {
		case <-time.After(50 * time.Second):
			log.Error("s2 negotiation timeout")
			return
		case err = <-cb2.iceresult:
			if err != nil {
				log.Error("s2 negotiation failed", err)
				return
			}
		}
		log.Info("s2 negotiation success")
	}
}
