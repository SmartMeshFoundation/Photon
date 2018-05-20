package ice

import (
	"fmt"
	"time"

	"errors"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/stun"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/turn"
)

/*
用于有 turn server 的情形下,收集本地候选地址列表.
*/
type turnSock struct {
	Client       *stun.Client
	s            *stunSocket
	user         string
	password     string
	nonce        string
	realm        string
	credentials  stun.MessageIntegrity //long term
	lifetime     turn.Lifetime         //how long is this allocate address valid
	localAddrs   []string
	mapAddress   string
	relayAddress string
	serverAddr   string
}

func newTurnSock(serverAddr, user, password string) (t *turnSock, err error) {
	var s *stunSocket
	s, err = newStunSocket(serverAddr)
	if err != nil {
		return
	}
	t = &turnSock{
		Client:     s.Client,
		s:          s,
		user:       user,
		password:   password,
		serverAddr: serverAddr,
	}
	return
}

func (t *turnSock) allocateAddress() error {
	deadline := time.Now().Add(t.s.ReadDeadline)
	var err error
	t.s.Client.Do(stun.MustBuild(stun.TransactionIDSetter, turn.AllocateRequest, turn.RequestedTransportUDP), deadline, func(res stun.Event) {
		if res.Error != nil {
			err = res.Error
			return
		}
		var (
			code  stun.ErrorCodeAttribute
			nonce stun.Nonce
			realm stun.Realm
		)
		err = code.GetFrom(res.Message)
		if err != nil {
			return
		}
		if code.Code != stun.CodeUnauthorised {
			log.Error(fmt.Sprintf("turn first allocate should faile, but code is %s", code))
			err = fmt.Errorf("unexpected turn code of error :%s", code)
			return
		}
		err = nonce.GetFrom(res.Message)
		if err != nil {
			return
		}
		err = realm.GetFrom(res.Message)
		if err != nil {
			return
		}
		log.Trace(fmt.Sprintf("get credentials nonce:%s,realm:%s,lieftime:%s", nonce, realm, t.lifetime.Duration))
		t.nonce = nonce.String()
		t.realm = realm.String()
		t.credentials = stun.NewLongTermIntegrity(t.user, t.realm, t.password)

	})
	if err != nil {
		return err
	}
	t.s.Client.Do(stun.MustBuild(stun.TransactionIDSetter, turn.AllocateRequest,
		turn.RequestedTransportUDP, stun.Realm(t.realm),
		stun.NewUsername(t.user), stun.Nonce(t.nonce), t.credentials), deadline, func(res stun.Event) {
		if res.Error != nil {
			err = res.Error
			return
		}
		var (
			code          stun.ErrorCodeAttribute
			RelayAddress  turn.RelayedAddress
			MappedAddress stun.XORMappedAddress
		)
		if res.Message.Type.Class == stun.ClassErrorResponse {
			code.GetFrom(res.Message)
			err = fmt.Errorf("got error response %s", code)
			log.Error(err.Error())
			return
		}
		err = MappedAddress.GetFrom(res.Message)
		if err != nil { //不考虑兼容rfc3489,肯定要有
			return
		}
		err = RelayAddress.GetFrom(res.Message)
		if err != nil {
			return
		}
		err = t.lifetime.GetFrom(res.Message)
		if err != nil {
			return
		}
		t.mapAddress = fmt.Sprintf("%s:%d", MappedAddress.IP, MappedAddress.Port)
		t.relayAddress = fmt.Sprintf("%s:%d", RelayAddress.IP, RelayAddress.Port)
	})
	if err != nil {
		return err
	}
	log.Trace(fmt.Sprintf("mappedaddr=%s,relay=%s", t.mapAddress, t.relayAddress))
	if len(t.mapAddress) == 0 || len(t.relayAddress) == 0 {
		return errors.New("can not get relay address")
	}
	//keep alive todo
	return nil
}

/*
第一个候选地址,必须是连接 turn server 的那个.
*/
func (t *turnSock) GetCandidates() (candidates []*Candidate, err error) {
	err = t.allocateAddress()
	if err != nil {
		return
	}
	c := new(Candidate)
	c.baseAddr = t.s.LocalAddr
	c.Type = CandidateServerReflexive
	c.addr = t.mapAddress
	c.Foundation = calcFoundation(c.baseAddr)
	c2 := new(Candidate)
	c2.Type = CandidateRelay
	c2.baseAddr = t.relayAddress
	c2.addr = t.relayAddress
	c2.Foundation = calcFoundation(c2.baseAddr)
	candidates, err = getLocalCandidates(c.baseAddr)
	if err != nil {
		return
	}
	for _, c := range candidates {
		t.localAddrs = append(t.localAddrs, c.addr)
	}
	if c.baseAddr != c.addr {
		candidates = append(candidates, c)
	}
	if c2.addr != c.baseAddr {
		candidates = append(candidates, c2)
	}
	return
}
func (t *turnSock) Close() {
	t.s.Close()
}

/*
address need to listen for input stun binding request...
*/
func (t *turnSock) getListenCandidiates() []string {
	return t.localAddrs
}
