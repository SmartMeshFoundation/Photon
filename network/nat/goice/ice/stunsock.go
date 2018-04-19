package ice

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/stun"
)

const DefaultReadDeadLine = time.Second * 10

type StunSocket struct {
	ServerAddr   string
	MappedAddr   net.UDPAddr
	LocalAddr    string // local addr used to  connect server
	Client       *stun.Client
	ReadDeadline time.Duration
	localAddrs   []string //for listen
}

func NewStunSocket(serverAddr string) (s *StunSocket, err error) {
	s = &StunSocket{
		ServerAddr:   serverAddr,
		ReadDeadline: DefaultReadDeadLine,
	}
	conn, err := net.Dial("udp", serverAddr)
	if err != nil {
		log.Fatalln("failed to dial:", err)
	}
	client, err := stun.NewClient(stun.ClientOptions{
		Connection: conn,
	})
	if err != nil {
		return
	}
	s.Client = client
	s.LocalAddr = conn.(*net.UDPConn).LocalAddr().String()
	return
}

//get mapped address from server
func (s *StunSocket) mapAddress() error {
	deadline := time.Now().Add(s.ReadDeadline)
	var err error
	wg := sync.WaitGroup{}
	wg.Add(1)
	err = s.Client.Do(stun.MustBuild(stun.TransactionIDSetter, stun.BindingRequest), deadline, func(res stun.Event) {
		defer wg.Done()
		if res.Error != nil {
			err = res.Error
			return
		}
		var xorAddr stun.XORMappedAddress
		if err = xorAddr.GetFrom(res.Message); err != nil {
			var addr stun.MappedAddress
			err = addr.GetFrom(res.Message)
			if err != nil {
				return
			}
			s.MappedAddr = net.UDPAddr{IP: addr.IP, Port: addr.Port}
		} else {
			s.MappedAddr = net.UDPAddr{IP: xorAddr.IP, Port: xorAddr.Port}
		}
	})
	wg.Wait()
	//keep alive todo
	return err
}

/*
获取有一部分信息的candidiate.第一个是本机主要地址,最后一个是缺省 Candidate
*/
func (s *StunSocket) GetCandidates() (candidates []*Candidate, err error) {
	err = s.mapAddress()
	if err != nil {
		return
	}
	c := new(Candidate)
	c.baseAddr = s.LocalAddr
	c.Type = CandidateServerReflexive
	c.addr = s.MappedAddr.String()
	c.Foundation = calcFoundation(c.baseAddr)
	candidates, err = GetLocalCandidates(c.baseAddr)
	if err != nil {
		return
	}
	for _, c := range candidates {
		s.localAddrs = append(s.localAddrs, c.addr)
	}
	if c.addr != c.baseAddr { //we have a public ip
		candidates = append(candidates, c)
	}
	return
}

func (s *StunSocket) Close() {
	if s.Client != nil {
		s.Client.Close()
	}
}

/*
address need to listen for input stun binding request...
*/
func (s *StunSocket) GetListenCandidiates() []string {
	return s.localAddrs
}
