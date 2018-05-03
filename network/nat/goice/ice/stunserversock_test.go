package ice

import (
	"testing"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/stun"
)

type mockcb struct {
	s ServerSocker
}

/*
 收到一个 stun.Message, 可能是 Bind Request/Bind Response 等等.
*/
func (m *mockcb) RecieveStunMessage(localAddr, remoteAddr string, req *stun.Message) {
	if req.Type != stun.BindingRequest {
		return
	}
	log.Info(fmt.Sprintf("recevied binding request %s<----%s", localAddr, remoteAddr))
	var res *stun.Message = new(stun.Message)
	from := addrToUdpAddr(remoteAddr)

	err := res.Build(
		stun.NewTransactionIDSetter(req.TransactionID),
		stun.NewType(stun.MethodBinding, stun.ClassSuccessResponse),
		software,
		&stun.XORMappedAddress{
			IP:   from.IP,
			Port: from.Port,
		},
		stun.Fingerprint,
	)
	if err != nil {
		panic(fmt.Sprintf("build res message error %s", err))
	}
	m.s.sendStunMessageAsync(res, localAddr, remoteAddr)
	return
}

/*
	ICE 协商建立连接以后,收到了对方发过来的数据,可能是经过 turn server 中转的 channel data( 不接受 sendData data request),也可能直接是数据.
	如果是经过 turn server 中转的, channelNumber 一定介于0x4000-0x7fff 之间.否则一定为0
*/
func (m *mockcb) ReceiveData(localAddr, peerAddr string, data []byte) {

}

//binding request 和普通的 stun message 一样处理.
//func (s *StunServerSock) processBindingRequest(from net.Addr, req *stun.Message) {

//notauthrized:
//	res.Build(stun.NewTransactionIDSetter(req.TransactionID), stun.BindingError,
//		stun.CodeUnauthorised, software, stun.Fingerprint)
//	s.sendStunMessageAsync(res, from)
//}
func setupTestServerSock() (s1, s2 *StunServerSock) {
	var err error
	mybindaddr := "127.0.0.1:8700"
	peerbindaddr := "127.0.0.1:8800"
	m1 := new(mockcb)
	m2 := new(mockcb)
	s1, err = NewStunServerSock(mybindaddr, m1, "s1")
	if err != nil {
		log.Crit(fmt.Sprintf("create new sock error %s %s", mybindaddr, err))
	}
	s2, err = NewStunServerSock(peerbindaddr, m2, "s2")
	if err != nil {
		log.Crit(fmt.Sprintf("creat new sock error %s %s", peerbindaddr, err))
	}
	m1.s = s1
	m2.s = s2
	return s1, s2
}
func TestNewServerSock(t *testing.T) {
	s1, s2 := setupTestServerSock()
	req, _ := stun.Build(stun.TransactionIDSetter, stun.BindingRequest, software, stun.Fingerprint)
	res, err := s1.sendStunMessageSync(req, s1.Addr, s2.Addr)
	if err != nil {
		t.Error(err)
		return
	}
	if res.Type != stun.BindingSuccess {
		t.Error("should success")
		return
	}
	log.Trace(fmt.Sprintf("s1 received :%s", res.String()))

}
