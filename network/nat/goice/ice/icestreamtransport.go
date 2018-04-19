package ice

import (
	"bytes"
	"fmt"
	"net"

	"strings"

	"errors"

	"github.com/nkbai/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/sdp"
)

type TransportOperation int

const (
	/** Initialization (candidate gathering) */
	ICE_TRANSPORT_OP_INIT TransportOperation = iota
	/** Negotiation */
	ICE_TRANSPORT_OP_NEGOTIATION
	/** This operation is used to report failure in keep-alive operation.
	*  Currently it is only used to report TURN Refresh failure.
	 */
	ICE_TRANSPORT_OP_KEEP_ALIVE
	/** IP address change notification from STUN keep-alive operation.
	 */
	ICE_TRANSPORT_OP_ADDR_CHANGE
)

type StreamTransportCallbacker interface {
	/*
			This callback will be called when the ICE transport receives
		     * incoming packet from the sockets which is not related to ICE
		     * (for example, normal RTP/RTCP packet destined for application).
	*/
	OnReceiveData(data []byte, from net.Addr)
	/*
		Callback to report status of various ICE operations.
	*/
	OnIceComplete(result error)
}

type TransportConfig struct {
	Server          string
	StunSever       string            //maybe empty
	TurnSever       string            //maybe empty
	TurnUserName    string
	TurnPassword    string
	ComponentNumber int //must be 1,right now
}

func NewTransportConfigHostonly() *TransportConfig {
	return &TransportConfig{
		ComponentNumber: 1,
	}
}
func NewTransportConfigWithStun(stunServer string) *TransportConfig {
	return &TransportConfig{
		StunSever:       stunServer,
		ComponentNumber: 1,
	}
}
func NewTransportConfigWithTurn(turnServer, turnUser, turnPass string) *TransportConfig {
	return &TransportConfig{
		TurnSever:       turnServer,
		TurnUserName:    turnUser,
		TurnPassword:    turnPass,
		ComponentNumber: 1,
	}
}

type TransportState int

const (
	/**
	 * ICE stream transport initialization/candidate gathering process is
	 * complete, ICE session may be created on this stream transport.
	 */
	TransportStateReady TransportState = iota
	/**
	 * New session has been created and the session is ready.
	 */
	TransportStateSessionReady
	/**
	 * ICE negotiation is in progress.
	 */
	TransportStateNegotiation
	/**
	 * ICE negotiation has completed successfully and media is ready
	 * to be used.
	 */
	TransportStateRunning
	/**
	 * ICE negotiation has completed with failure.
	 */
	TransportStateFailed
	/**
	 * ICE negotiation has completed with failure.
	 */
	TransportStateStopped
)

func (s TransportState) String() string {
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

type IceStreamTransport struct {
	Name        string //debug info
	cfg         *TransportConfig
	transporter StunTranporter
	component   *TransportComponent
	State       TransportState
	session     *IceSession
	cb          StreamTransportCallbacker
}

func NewIceStreamTransport(cfg *TransportConfig, name string) (it *IceStreamTransport, err error) {
	it = &IceStreamTransport{
		cfg:   cfg,
		State: TransportStateReady,
		Name:  name,
	}
	if len(it.cfg.StunSever) > 0 {
		it.transporter, err = NewStunSocket(it.cfg.StunSever)
	} else if len(it.cfg.TurnSever) > 0 {
		it.transporter, err = NewTurnSock(it.cfg.TurnSever, cfg.TurnUserName, cfg.TurnPassword)
	} else {
		it.transporter = new(HostOnlySock)
	}
	if err != nil {
		return
	}
	it.component = NewTransportComponent(it.transporter, 1)
	_, err = it.component.GetCandidates()
	if err != nil {
		return
	}
	log.Trace("candidates=%s", log.StringInterface(it.component.candidates, 3))
	return
}
func (t *IceStreamTransport) InitIce(role SessionRole) error {
	s := NewIceSession(t.Name, role, t.component.candidates, t.transporter, t)
	t.session = s
	for i, c := range s.localCandidates {
		log.Trace("%s Candidate %d added componentId=%d type=%s foundation=%d,addr=%s,base=%s,priority=%d",
			t.Name, i, c.ComponentID, c.Type, c.Foundation, c.addr, c.baseAddr, c.Priority,
		)
	}
	err := t.session.StartServer()
	if err != nil {
		return err
	}
	t.State = TransportStateSessionReady
	return nil
}
func (t *IceStreamTransport) SetCallBack(cb StreamTransportCallbacker) {
	t.cb = cb
}
func (t *IceStreamTransport) StartNegotiation(remoteSDP string) error {
	var err error
	if t.session == nil || t.State != TransportStateSessionReady {
		return errors.New("no session")
	}
	//defer func() {
	//	if err != nil {
	//		t.session.Stop()
	//	}
	//}()
	log.Trace("%s received sdp \n%s\n", t.Name, remoteSDP)
	sd, err := DecodeSession(remoteSDP)
	if err != nil {
		return err
	}
	err = t.session.createCheckList(sd)
	if err != nil {
		return err
	}
	log.Trace("%s checklist created\n%s", t.Name, t.session.checkList)
	err = t.session.createTurnPermissionIfNeeded()
	if err != nil {
		return err
	}
	log.Trace("create permission success for all remote address")
	t.State = TransportStateNegotiation
	err = t.session.startCheck()
	if err != nil {
		return err
	}

	return nil
}
func (t *IceStreamTransport) EncodeSession() (s string, err error) {
	if t.session == nil {
		err = errors.New(fmt.Sprintf("no session and state =%d", t.State))
		return
	}
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "v=0\no=- 3414953978 3414953978 IN IP4 localhost\ns=ice\nt=0 0\n")
	fmt.Fprintf(buf, "a=ice-ufrag:%s\na=ice-pwd:%s\n", t.session.rxUserFrag, t.session.rxPassword)
	//only on component now....
	uaddr := addrToUdpAddr(t.component.defaultCandidate.addr)
	fmt.Fprintf(buf, "m=audio %d RTP/AVP 0\nc=IN IP4 %s\n", uaddr.Port, uaddr.IP.String())
	for _, c := range t.component.candidates {
		fmt.Fprintf(buf, "%s\n", c)
	}
	return string(buf.Bytes()), nil
}

//不支持复用,只能完全重新构建.
func (t *IceStreamTransport) Stop() {
	if t.State == TransportStateStopped {
		log.Error("%s has already stopped", t.Name)
		return
	}
	t.State = TransportStateFailed
	if t.session != nil {
		t.session.Stop()
	}
}
func (t *IceStreamTransport) SendData(data []byte) error {
	if t.State != TransportStateRunning {
		return errors.New("transport not running.")
	}
	check := t.session.sessionComponent.nominatedCheck
	srv := t.session.sessionComponent.nominatedServerSock
	if check == nil {
		return errors.New("no check.")
	}
	if srv == nil {
		return errors.New("no stun transport")
	}
	return srv.sendData(data, check.localCandidate.addr, check.remoteCandidate.addr)
}
func (t *IceStreamTransport) onIceComplete(result error) {

	if t.State != TransportStateNegotiation {
		log.Error("%s finish reulst %s", t.Name, result)
		panic("only finish once")
	}
	defer func() {
		if t.cb != nil {
			t.cb.OnIceComplete(result)
		}
		log.Debug("%s ice negotiation finished ,new state is %s", t.Name, t.State.String())
	}()
	if result != nil {
		log.Info("%s ice negotiation failed", t.Name)
		t.State = TransportStateFailed
		//t.Stop()
		return
	}
	srv := t.session.getSenderServerSock(t.session.sessionComponent.nominatedCheck.localCandidate.addr)
	t.session.sessionComponent.nominatedServerSock = srv
	check := t.session.sessionComponent.nominatedCheck
	if check.localCandidate.Type == CandidateRelay {
		result = t.session.turnServerSock.channelBind(check.remoteCandidate.addr)
		if result != nil {
			log.Error("%s channel bind err:%s", t.Name, result)
			t.State = TransportStateFailed
			//t.Stop()
			return
		}
		srv.FinishNegotiation(TurnModeData)
	} else {
		srv.FinishNegotiation(StunModeData)
	}
	for k, srv2 := range t.session.serverSocks {
		if srv != srv2 {
			delete(t.session.serverSocks, k)
			srv2.Close()
		}
	}
	if srv != t.session.turnServerSock {
		t.session.turnServerSock = nil
	}
	t.State = TransportStateRunning
}
func (t *IceStreamTransport) onRxData(data []byte, from string) {
	if t.cb != nil {
		t.cb.OnReceiveData(data, addrToUdpAddr(from))
	}
}
func DecodeSession(str string) (session *sessionDescription, err error) {
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
			fmt.Sscanf(v, "IN IP4 %s", &session.defaultIp)
		}
	}
	if len(session.user) <= 0 || len(session.password) == 0 ||
		len(session.defaultIp) == 0 ||
		len(session.candidates) == 0 {
		err = fmt.Errorf("remote session description error %s", str)
		return
	}
	s2 := fmt.Sprintf("%s:%d", session.defaultIp, session.defaultPort)
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

type sessionDescription struct {
	user            string
	password        string
	defaultPort     int
	defaultIp       string
	candidates      []*Candidate
	defautCandidate *Candidate
}
type TransportComponent struct {
	Name             string
	componentId      int
	candidates       []*Candidate
	defaultCandidate *Candidate
	candidateGetter  CandidateGetter
}

func NewTransportComponent(candidateGetter CandidateGetter, id int) *TransportComponent {
	return &TransportComponent{
		candidateGetter: candidateGetter,
		componentId:     id,
	}
}

func (t *TransportComponent) GetCandidates() (candidates []*Candidate, err error) {
	candidates, err = t.candidateGetter.GetCandidates()
	if err != nil {
		return
	}
	for _, c := range candidates {
		c.ComponentID = t.componentId
		c.Priority = calcCandidatePriority(c.Type, DefaultPreference, c.ComponentID)
		c.transport = TransportUDP
	}
	t.candidates = candidates
	t.defaultCandidate = candidates[len(candidates)-1]
	return
}
