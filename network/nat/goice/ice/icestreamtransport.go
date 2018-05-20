package ice

import (
	"bytes"
	"fmt"
	"net"

	"strings"

	"errors"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/sdp"
)

//StreamTransportCallbacker callback of ICE
type StreamTransportCallbacker interface {
	/*
			OnReceiveData will be called when the ICE transport receives
		     * incoming packet from the sockets which is not related to ICE
		     * (for example, normal RTP/RTCP packet destined for application).
	*/
	OnReceiveData(data []byte, from net.Addr)
	/*
		OnIceComplete report status of various ICE operations.
	*/
	OnIceComplete(result error)
}

/*
TransportConfig is configuration of ICE
*/
type TransportConfig struct {
	Server          string
	StunSever       string //maybe empty
	TurnSever       string //maybe empty
	TurnUserName    string
	TurnPassword    string
	ComponentNumber int //must be 1,right now
}

//StreamTransport is a transport
type StreamTransport struct {
	Name        string //debug info
	cfg         *TransportConfig
	transporter stunTranporter
	component   *transportComponent
	State       transportState
	session     *session
	cb          StreamTransportCallbacker
	log         log.Logger
}

type sessionDescription struct {
	user            string
	password        string
	defaultPort     int
	defaultIP       string
	candidates      []*Candidate
	defautCandidate *Candidate
}
type transportComponent struct {
	Name             string
	componentID      int
	candidates       []*Candidate
	defaultCandidate *Candidate
	candidateGetter  candidateGetter
}

//NewTransportConfigHostonly return a  hostonly config
func NewTransportConfigHostonly() *TransportConfig {
	return &TransportConfig{
		ComponentNumber: 1,
	}
}

//NewTransportConfigWithStun return a stun config
func NewTransportConfigWithStun(stunServer string) *TransportConfig {
	return &TransportConfig{
		StunSever:       stunServer,
		ComponentNumber: 1,
	}
}

//NewTransportConfigWithTurn return a turn config
func NewTransportConfigWithTurn(turnServer, turnUser, turnPass string) *TransportConfig {
	return &TransportConfig{
		TurnSever:       turnServer,
		TurnUserName:    turnUser,
		TurnPassword:    turnPass,
		ComponentNumber: 1,
	}
}

type transportState int

const (
	/*TransportStateReady represents that ICE stream transport initialization/candidate gathering process is
	complete, ICE session may be created on this stream transport.
	*/
	TransportStateReady transportState = iota
	/*TransportStateSessionReady means that New session has been created and the session is ready.
	 */
	TransportStateSessionReady
	/*TransportStateNegotiation means that ICE negotiation is in progress.
	 */
	TransportStateNegotiation
	/*TransportStateRunning means that ICE negotiation has completed successfully and media is ready
	  to be used.
	*/
	TransportStateRunning
	/*TransportStateFailed means ICE negotiation has completed with failure.
	 */
	TransportStateFailed
	/*TransportStateStopped ICE negotiation has completed with failure.
	 */
	TransportStateStopped
)

func (s transportState) String() string {
	switch s {
	case TransportStateReady:
		return "Candidate Gathering Complete"
	case TransportStateSessionReady:
		return "Session Initialized"
	case TransportStateNegotiation:
		return "Negotiation In Progress"
	case TransportStateRunning:
		return "Negotiation Success"
	case TransportStateFailed:
		return "Negotiation Failed"
	case TransportStateStopped:
		return "stopped"
	}
	return "unkown"
}

//NewIceStreamTransport create streamTransport from configuration
//name is for debug
func NewIceStreamTransport(cfg *TransportConfig, name string) (it *StreamTransport, err error) {
	it = &StreamTransport{
		cfg:   cfg,
		State: TransportStateReady,
		Name:  name,
		log:   log.New("name", fmt.Sprintf("%s-StreamTransport", name)),
	}
	if len(it.cfg.StunSever) > 0 {
		it.transporter, err = newStunSocket(it.cfg.StunSever)
	} else if len(it.cfg.TurnSever) > 0 {
		it.transporter, err = newTurnSock(it.cfg.TurnSever, cfg.TurnUserName, cfg.TurnPassword)
	} else {
		it.transporter = new(HostOnlySock)
	}
	if err != nil {
		return
	}
	it.component = newTransportComponent(it.transporter, 1)
	_, err = it.component.GetCandidates()
	if err != nil {
		return
	}
	it.log.Trace(fmt.Sprintf("candidates=%#v", it.component.candidates))
	return
}

//InitIce set role of this transport
func (t *StreamTransport) InitIce(role SessionRole) error {
	s := newIceSession(t.Name, role, t.component.candidates, t.transporter, t)
	t.session = s
	for i, c := range s.localCandidates {
		t.log.Trace(fmt.Sprintf("%s Candidate %d added componentID=%d type=%s foundation=%d,addr=%s,base=%s,priority=%d",
			t.Name, i, c.ComponentID, c.Type, c.Foundation, c.addr, c.baseAddr, c.Priority,
		))
	}
	err := t.session.StartServer()
	if err != nil {
		return err
	}
	t.State = TransportStateSessionReady
	return nil
}

//SetCallBack set the callback
// TODO should move set to NewIceStreamTransport
func (t *StreamTransport) SetCallBack(cb StreamTransportCallbacker) {
	t.cb = cb
}

//StartNegotiation starts negotiation process. can be called only once.
func (t *StreamTransport) StartNegotiation(remoteSDP string) (err error) {
	defer func() {
		if err != nil {
			t.log.Error(fmt.Sprintf("StartNegotiation with remotesdp err =%s", err))
		}
	}()
	if t.session == nil || t.State != TransportStateSessionReady {
		err = errors.New("no session")
		return
	}
	t.log.Trace(fmt.Sprintf("%s received sdp \n%s\n", t.Name, remoteSDP))
	t.State = TransportStateNegotiation
	sd, err := decodeSession(remoteSDP)
	if err != nil {
		return
	}
	err = t.session.createCheckList(sd)
	if err != nil {
		return
	}
	t.log.Trace(fmt.Sprintf("checklist created\n%s", t.session.checkList))
	err = t.session.createTurnPermissionIfNeeded()
	if err != nil {
		return
	}
	t.log.Trace(fmt.Sprintf("create permission success for all remote address"))

	err = t.session.startCheck()
	if err != nil {
		return
	}
	return nil
}

//EncodeSession encoding ice info to sdp
func (t *StreamTransport) EncodeSession() (s string, err error) {
	if t.session == nil {
		err = fmt.Errorf("no session and state =%d", t.State)
		return
	}
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "v=0\no=- 3414953978 3414953978 IN IP4 localhost\ns=ice\nt=0 0\n")
	fmt.Fprintf(buf, "a=ice-ufrag:%s\na=ice-pwd:%s\n", t.session.rxUserFrag, t.session.rxPassword)
	//only on component now....
	uaddr := addrToUDPAddr(t.component.defaultCandidate.addr)
	fmt.Fprintf(buf, "m=audio %d RTP/AVP 0\nc=IN IP4 %s\n", uaddr.Port, uaddr.IP.String())
	for _, c := range t.component.candidates {
		fmt.Fprintf(buf, "%s\n", c)
	}
	return string(buf.Bytes()), nil
}

//Stop destroy this transport, and cannot be reused
func (t *StreamTransport) Stop() {
	if t.State == TransportStateStopped {
		t.log.Error(fmt.Sprintf("%s has already stopped", t.Name))
		return
	}
	t.State = TransportStateFailed
	if t.session != nil {
		t.session.Stop()
	}
}

//SendData send data to peer, peer's ip and port are select by ice
func (t *StreamTransport) SendData(data []byte) error {
	if t.State != TransportStateRunning {
		return errors.New("transport not running")
	}
	return t.session.SendData(data)
}

/*
保证只会被调用一次,表示已经找到了至少一个有效连接,可以发送数据了,
但是这个连接未必是最后确定的,可能会发生变化.
*/
func (t *StreamTransport) onIceComplete(result error) {

	if t.State != TransportStateNegotiation {
		t.log.Error(fmt.Sprintf("%s finish reulst %s,t.State=%d", t.Name, result, t.State))
		panic(fmt.Sprintf("%s only finish once", t.Name))
	}
	defer func() {
		if t.cb != nil {
			t.cb.OnIceComplete(result)
		}
		t.log.Debug(fmt.Sprintf("%s ice negotiation finished ,new state is %s", t.Name, t.State.String()))
	}()
	if result != nil {
		t.log.Info(fmt.Sprintf("%s ice negotiation failed", t.Name))
		t.State = TransportStateFailed
		return
	}
	t.State = TransportStateRunning
}

/*
收到数据,并不表示协商已经完毕,而是对方找到了一条有效连接.
*/
func (t *StreamTransport) onRxData(data []byte, from string) {
	if t.cb != nil {
		t.cb.OnReceiveData(data, addrToUDPAddr(from))
	}
}
func decodeSession(str string) (session *sessionDescription, err error) {
	var s sdp.Session
	s, err = sdp.DecodeSession([]byte(str), s)
	if err != nil {
		return
	}
	session = &sessionDescription{}
	for _, line := range s {
		v := string(line.Value)
		//log.Trace(v)
		switch line.Type {
		case sdp.TypeAttribute:
			ss := strings.Split(v, ":")
			if len(ss) != 2 {
				err = fmt.Errorf("attribute error :%s", v)
				return
			}
			switch ss[0] {
			case "ice-ufrag":
				session.user = ss[1]
			case "ice-pwd":
				session.password = ss[1]
			case "candidate":
				parser := candidateParser{
					buf: line.Value,
					c:   new(Candidate),
				}
				err = parser.parse()
				if err != nil {
					return
				}
				session.candidates = append(session.candidates, parser.c)
			}
		case sdp.TypeMediaDescription:
			fmt.Sscanf(v, "audio %d RTP/", &session.defaultPort)
		case sdp.TypeConnectionData:
			fmt.Sscanf(v, "IN IP4 %s", &session.defaultIP)
		}
	}
	if len(session.user) <= 0 || len(session.password) == 0 ||
		len(session.defaultIP) == 0 ||
		len(session.candidates) == 0 {
		err = fmt.Errorf("remote session description error %s", str)
		return
	}
	s2 := fmt.Sprintf("%s:%d", session.defaultIP, session.defaultPort)
	for _, c := range session.candidates {
		if c.addr == s2 {
			session.defautCandidate = c
			break
		}
	}
	if session.defautCandidate == nil {
		err = fmt.Errorf("no default candidate found %s", s2)
	}
	return
}

func newTransportComponent(candidateGetter candidateGetter, id int) *transportComponent {
	return &transportComponent{
		candidateGetter: candidateGetter,
		componentID:     id,
	}
}

func (t *transportComponent) GetCandidates() (candidates []*Candidate, err error) {
	candidates, err = t.candidateGetter.GetCandidates()
	if err != nil {
		return
	}
	for _, c := range candidates {
		c.ComponentID = t.componentID
		c.Priority = calcCandidatePriority(c.Type, defaultPreference, c.ComponentID)
		c.transport = TransportUDP
	}
	t.candidates = candidates
	t.defaultCandidate = candidates[len(candidates)-1]
	return
}
