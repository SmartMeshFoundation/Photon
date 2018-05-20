package ice

import (
	"testing"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/stun"
)

func setupTurnServerSock() (s1, s2 *turnServerSock) {
	t1 := newTestTurnSock()
	t2 := newTestTurnSock()
	err := t1.allocateAddress()
	if err != nil {
		panic(err)
	}
	err = t2.allocateAddress()
	if err != nil {
		panic(err)
	}
	t1.Close()
	t2.Close()
	cfg1 := &turnServerSockConfig{
		user:         t1.user,
		password:     t1.password,
		nonce:        t1.nonce,
		realm:        t1.realm,
		credentials:  t1.credentials,
		lifetime:     t1.lifetime,
		serverAddr:   t1.serverAddr,
		relayAddress: t1.relayAddress,
	}
	cfg2 := &turnServerSockConfig{
		user:         t2.user,
		password:     t2.password,
		nonce:        t2.nonce,
		realm:        t2.realm,
		credentials:  t2.credentials,
		lifetime:     t2.lifetime,
		serverAddr:   t2.serverAddr,
		relayAddress: t2.relayAddress,
	}
	m1 := new(mockcb)
	m2 := new(mockcb)
	s1, err = newTurnServerSockWrapper(t1.s.LocalAddr, "s1", m1, cfg1)
	if err != nil {
		panic(err)
	}
	s2, err = newTurnServerSockWrapper(t2.s.LocalAddr, "s2", m2, cfg2)
	if err != nil {
		panic(err)
	}
	candidates1, err := t1.GetCandidates()
	if err != nil {
		panic(err)
	}
	candidates2, err := t2.GetCandidates()
	if err != nil {
		panic(err)
	}
	m1.s = s1
	m2.s = s2
	_, err = s1.createPermission(candidates2)
	if err != nil {
		panic(err)
	}
	_, err = s2.createPermission(candidates1)
	if err != nil {
		panic(err)
	}
	log.Trace("------------------------------")
	return
}
func TestNewTurnServerSockWrapper(t *testing.T) {
	s1, s2 := setupTurnServerSock()
	req, _ := stun.Build(stun.TransactionIDSetter, stun.BindingRequest, software, stun.Fingerprint)
	res, err := s1.sendStunMessageSync(req, s1.cfg.relayAddress, s2.cfg.relayAddress)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}
