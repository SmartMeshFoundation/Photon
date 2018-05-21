package ice

import (
	"fmt"

	"errors"

	"sort"

	"time"

	"net"

	"sync"

	"encoding/hex"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/ice/attr"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/stun"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/turn"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

//SessionRole role of ICE
type SessionRole int

const (
	//SessionRoleUnkown error state
	SessionRoleUnkown SessionRole = iota
	//SessionRoleControlling role controlling
	SessionRoleControlling
	//SessionRoleControlled role controlled
	SessionRoleControlled
)

func (s SessionRole) String() string {
	switch s {
	case SessionRoleUnkown:
		return "unkown"
	case SessionRoleControlled:
		return "controlled"
	case SessionRoleControlling:
		return "controlling"
	}
	return "unknown"
}

type session struct {
	Name string //for debug
	/*
		controlling or controlled
	*/
	role SessionRole
	/**
	 * Specify whether to use aggressive nomination.
	 */
	aggresive bool

	/**
	 * For a controlled agent, specify how long it wants to wait (in
	 * milliseconds) for the controlling agent to complete sending
	 * connectivity check with nominated flag set to true for all components
	 * after the controlled agent has found that all connectivity checks in
	 * its checklist have been completed and there is at least one successful
	 * (but not nominated) check for every component.
	 *
	 * Default value for this option is
	 * ICE_CONTROLLED_AGENT_WAIT_NOMINATION_TIMEOUT. Specify -1 to disable
	 * this timer.
	 */
	controlledAgentWaitNomiatedTimeout time.Duration
	/* STUN credentials */
	txUserName       string /**< Uname for TX.	TxUserFrag:RxUserFrag    */
	txPassword       string /**< Remote password.   */
	rxUserFrag       string /**< Local ufrag.	    */
	rxUserName       string /**< Uname for RX	    */
	rxPassword       string /**< Local password.    */
	txCrendientials  stun.MessageIntegrity
	rxCrendientials  stun.MessageIntegrity
	sessionComponent sessionComponet
	localCandidates  []*Candidate
	remoteCandidates []*Candidate
	checkList        *sessionCheckList
	validCheckList   *sessionCheckList // check has been verified and is valid.
	transporter      stunTranporter    //获取 candidates 用的 stunclient, 可能只指定了一个 stun 服务器,而没有 turn 服务器,也可能两者都没有.
	/*
		探测的过程中,按照协议要求,必须从指定的 ip 地址和端口发送探测数据,因此,如果本机有多个 ip 地址,那么就会有多个 serverSocker
	*/
	serverSocks map[string]serverSocker
	/*
			按照现在的实现,连接到 stun/turn 服务器的那个需要特殊处理,
			只有他发送数据的时候,可能需要经过 turn server 中转.
		当然如果真的没有 turnserver, 也不影响,它会是 nil, 也不会从服务器发送中转数据
	*/
	turnServerSock *turnServerSock

	isNominating bool /* Nominating stage   */
	//write this chan to finish one check.
	checkMap map[string]chan error
	//todo refer state, etc.
	iceStreamTransport *StreamTransport
	/*
	   用于角色冲突的时候,自行进行角色切换 ICEROLECONTROLLING <--->ICEROLECONTROLLED
	*/
	tieBreaker     uint64
	earlyCheckList []*rxCheck
	msg2Check      map[stun.TransactionID]*sessionCheck
	mlock          sync.Mutex

	/*
		收到的stun message, 不要堵塞发送接收routine
	*/
	msgChan        chan *stunMessageWrapper
	dataChan       chan *stunDataWrapper
	tryFailChan    chan *checkFailedWrapper
	quitChan       chan struct{}         //close when stop
	hasStopped     bool                  //停止销毁相关资源时,标记.
	completeResult sessionCompleteResult //0,not complete ,1 complete success, 2 complete failure
	log            log.Logger
}
type sessionCompleteResult int

const (
	sessionNotComplete        sessionCompleteResult = iota
	sessionCheckComplete                            //wait for nomination.
	sessionCompleteSuccess                          //at least have one nomination
	sessionAllCompleteSuccess                       //所有的 check 都 finish 了.
	sessionCompleteFailure
)

type sessionComponet struct {
	/**
	 * Pointer to ICE check with highest priority which connectivity check
	 * has been successful. The value will be NULL if a no successful check
	 * has not been found for this component.
	 */
	validCheck *sessionCheck
	/**
	 * Pointer to ICE check with highest priority which connectivity check
	 * has been successful and it has been nominated. The value may be NULL
	 * if there is no such check yet.
	 */
	nominatedCheck *sessionCheck
	/*
		nominated server socker
		to send data to peer.
	*/
	nominatedServerSock serverSocker
}

/**
 * This structure represents an incoming check (an incoming Binding
 * request message), and is mainly used to keep early checks in the
 * list in the ICE session. An early check is a request received
 * from remote when we haven't received SDP answer yet, therefore we
 * can't perform triggered check. For such cases, keep the incoming
 * request in a list, and we'll do triggered checks (simultaneously)
 * as soon as we receive answer.
 */
type rxCheck struct {
	componentID   int
	remoteAddress string
	localAddress  string
	userCandidate bool
	priority      int
	role          SessionRole
}

func (r *rxCheck) String() string {
	return fmt.Sprintf("{remote=%s,local=%s,userCandidate=%v,priorit=%d,role=%d}", r.remoteAddress, r.localAddress, r.userCandidate, r.priority, r.role)
}

type stunMessageWrapper struct {
	localAddr  string
	remoteAddr string
	msg        *stun.Message
}
type stunDataWrapper struct {
	localAddr  string
	remoteAddr string
	data       []byte
}

/*
尝试多次失败以后,需要结束这次 check,
*/
type checkFailedWrapper struct {
	c   *sessionCheck
	err error
}

/*
ice session运行着四种协程
1.来自上层的调用
2.来自下层 socket 的消息协程
3.自身的 loop 协程
4.check 时候的大量协程,
*/
func newIceSession(name string, role SessionRole, localCandidates []*Candidate, transporter stunTranporter, ice *StreamTransport) *session {
	s := &session{
		Name:               name,
		role:               role,
		aggresive:          true,
		rxUserFrag:         utils.RandomString(8),
		rxPassword:         utils.RandomString(8),
		localCandidates:    localCandidates,
		transporter:        transporter,
		checkMap:           make(map[string]chan error),
		iceStreamTransport: ice,
		checkList:          new(sessionCheckList),
		validCheckList:     new(sessionCheckList),
		tieBreaker:         attr.RandUint64(),
		serverSocks:        make(map[string]serverSocker),
		msg2Check:          make(map[stun.TransactionID]*sessionCheck),
		msgChan:            make(chan *stunMessageWrapper, 10),
		dataChan:           make(chan *stunDataWrapper, 10),
		quitChan:           make(chan struct{}),
		tryFailChan:        make(chan *checkFailedWrapper, 10),
		log:                log.New("name", fmt.Sprintf("%s-icesession", name)),
		controlledAgentWaitNomiatedTimeout: time.Second * 10,
	}
	s.rxCrendientials = stun.NewShortTermIntegrity(s.rxPassword)
	//make sure the first candidates is used to communicate with stun/turn server

	return s
}

var errTooManyCandidates = errors.New("too many candidates")

func (s *session) addMsgCheck(id stun.TransactionID, check *sessionCheck) {
	s.mlock.Lock()
	s.msg2Check[id] = check
	s.mlock.Unlock()
}
func (s *session) getMsgCheck(id stun.TransactionID) *sessionCheck {
	s.mlock.Lock()
	defer s.mlock.Unlock()
	return s.msg2Check[id]
}
func (s *session) deleteMsgCheck(id stun.TransactionID) {
	s.mlock.Lock()
	delete(s.msg2Check, id)
	s.mlock.Unlock()
}
func (s *session) Stop() {
	s.hasStopped = true
	for _, srv := range s.serverSocks {
		srv.Close()
	}
	for _, c := range s.checkMap {
		close(c)
	}
	close(s.quitChan) //avoid send on close
}
func (s *session) createCheckList(sd *sessionDescription) error {
	if len(sd.candidates) > maxCandidates {
		return errTooManyCandidates
	}
	s.txUserName = fmt.Sprintf("%s:%s", sd.user, s.rxUserFrag)
	s.rxUserName = fmt.Sprintf("%s:%s", s.rxUserFrag, sd.user)
	s.txPassword = sd.password
	s.txCrendientials = stun.NewShortTermIntegrity(s.txPassword)
	for _, c := range sd.candidates {
		if c.ComponentID != 1 { //only one component,
			continue
		}
		s.remoteCandidates = append(s.remoteCandidates, c)
	}
	for _, l := range s.localCandidates {
		for _, r := range s.remoteCandidates {
			chk := &sessionCheck{
				localCandidate:  l,
				remoteCandidate: r,
				key:             fmt.Sprintf("%s-%s", l.addr, r.addr),
				state:           checkStateFrozen,
				priority:        calcPairPriority(s.role, l, r),
			}
			s.checkList.checks = append(s.checkList.checks, chk)
		}
	}
	if len(s.checkList.checks) == 0 {
		return errors.New("no matched candidate found")
	}
	//priority from high to low. not stable
	sort.Stable(s.checkList)
	s.pruneCheckList()
	return nil
}

/* Since an agent cannot sendData requests directly from a reflexive
 * candidate, but only from its base, the agent next goes through the
 * sorted list of candidate pairs.  For each pair where the local
 * candidate is server reflexive, the server reflexive candidate MUST be
 * replaced by its base.  Once this has been done, the agent MUST prune
 * the list.  This is done by removing a pair if its local and remote
 * candidates are identical to the local and remote candidates of a pair
 * higher up on the priority list.  The result is a sequence of ordered
 * candidate pairs, called the check list for that media stream.
 */
func (s *session) pruneCheckList() {
	m := make(map[string]bool)
	var checks []*sessionCheck
	for _, c := range s.checkList.checks {
		key := fmt.Sprintf("%d%s", c.localCandidate.Foundation, c.remoteCandidate.addr)
		if m[key] {
			continue
		}
		m[key] = true
		checks = append(checks, c)
	}
	s.checkList.checks = checks
}

/*
如何处理新的数据到来的通知.
来自其他人的 binding request,是必须包含短期认证的,否则可能出现错误
B.4. Importance of the STUN Username

ICE requires the usage of message integrity with STUN using its short-term credential functionality.
The actual short-term credential is formed by exchanging username fragments in the SDP offer/answer exchange.
The need for this mechanism goes beyond just security; it is actually required for correct operation
of ICE in the first place.
*/
func (s *session) StartServer() (err error) {
	defer func() {
		if err != nil {
			for _, srv := range s.serverSocks {
				srv.Close()
			}
		}
	}()
	s.transporter.Close() //首先要关闭这个连接,否则没法再次 Listen, 会提示被占用
	turnsock, hasRelay := s.transporter.(*turnSock)
	start := 0
	candidates := s.transporter.getListenCandidiates()
	if hasRelay {
		start = 1
		cfg := &turnServerSockConfig{
			user:         turnsock.user,
			password:     turnsock.password,
			nonce:        turnsock.nonce,
			realm:        turnsock.realm,
			credentials:  turnsock.credentials,
			relayAddress: turnsock.relayAddress,
			serverAddr:   turnsock.serverAddr,
			lifetime:     turnsock.lifetime,
		}
		s.turnServerSock, err = newTurnServerSockWrapper(candidates[0], s.Name, s, cfg)
		if err != nil {
			return err
		}
		s.serverSocks[candidates[0]] = s.turnServerSock
	}
	for ; start < len(candidates); start++ {
		var srv *stunServerSock
		srv, err = newStunServerSock(candidates[start], s, s.Name)
		if err != nil {
			return err
		}
		s.serverSocks[candidates[start]] = srv
	}
	go s.loop()
	return
}

/*
创建checklist 以后,如果本地有 relay 的候选地址,
那么需要在turn server 上专门设置,对方发送到 turn server 的数据才会中转给我.
否则会被 turn server 丢弃.
*/
func (s *session) createTurnPermissionIfNeeded() (err error) {
	var res *stun.Message
	if s.turnServerSock != nil {
		res, err = s.turnServerSock.createPermission(s.remoteCandidates)
		if err != nil {
			return
		}
		if res.Type != turn.CreatePermissionResponse {
			return errors.New("Create permission error")
		}
	}
	return nil
}

/*
check stage:
one check received a valid response
*/
func (s *session) handleCheckResponse(check *sessionCheck, from string, res *stun.Message) {
	var err error
	s.log.Trace(fmt.Sprintf("handle check response check=%s\nfrom=%s\n res=%s", check.String(), from, res.String()))
	if from != check.remoteCandidate.addr {
		s.log.Info(fmt.Sprintf("check received stun message not from expected address,got:%s,check is %s", from, check))
	}
	if s.completeResult >= sessionAllCompleteSuccess {
		/*
			收到前期bindingRequest 的 response, 但是已经协商完毕了,忽略即可.
		*/
		s.log.Info(fmt.Sprintf("%s received checkresponse ,but check is already finished.", check))
		return
	}
	if s.getMsgCheck(res.TransactionID) == nil {
		/*
			因为某些原因,认为这个 check 已经失败了,
			比如超时,直接丢弃即可.
		*/
		return
	}
	if res.Type.Class == stun.ClassErrorResponse {
		s.log.Info(fmt.Sprintf("%s received error response %s", check.key, res.Type))
		var code stun.ErrorCodeAttribute
		err = code.GetFrom(res)
		if err != nil || code.Code != stun.CodeRoleConflict {
			s.changeCheckState(check, checkStateFailed, fmt.Errorf("unkown error code %s", code))
			s.tryCompleteCheck(check)
			return
		}
		/* Role conclict response.
		 *
		 * 7.1.2.1.  Failure Cases:
		 *
		 * If the request had contained the ICE-CONTROLLED attribute,
		 * the agent MUST switch to the controlling role if it has not
		 * already done so.  If the request had contained the
		 * ICE-CONTROLLING attribute, the agent MUST switch to the
		 * controlled role if it has not already done so.  Once it has
		 * switched, the agent MUST immediately retry the request with
		 * the ICE-CONTROLLING or ICE-CONTROLLED attribute reflecting
		 * its new role.
		 */
		var newrole = SessionRoleUnkown
		_, err = res.Get(stun.AttrICEControlled)
		if err == nil {
			newrole = SessionRoleControlling
		} else if _, err = res.Get(stun.AttrICEControlling); err == nil {
			newrole = SessionRoleControlled
		}
		if newrole != s.role {
			s.log.Trace(fmt.Sprintf("change role from %s to %s", s.role, newrole))
			s.role = newrole
		}
		s.retryOneCheck(check)
		return
	}
	/*
			我作为 controlled 的时候,对方有可能先收到 binding request,然后才收到信令服务器过来的sdp, 从而导致
			对方发送的 bingding response 中的 Message Integrity是错的,我们不能把它当成错误处理.
		todo 当我作为 controlled 时候一旦收到对方的 bindingRequest ,应该明确知道以后的 response Message Integrity 必须是对的.
	*/
	if err = s.rxCrendientials.Check(res); err != nil {
		if s.role == SessionRoleControlling {
			err = fmt.Errorf("receive check response,but crendientials check failed %s", err)
			s.log.Error(err.Error())
			s.changeCheckState(check, checkStateFailed, err)
			s.tryCompleteCheck(check) //is this the last check?
		} else {
			s.log.Warn(fmt.Sprintf("receive check response,but crendientials check failed %s", err))
		}

	}
	/* 7.1.2.1.  Failure Cases
	 *
	 * The agent MUST check that the source IP address and port of the
	 * response equals the destination IP address and port that the Binding
	 * Request was sent to, and that the destination IP address and port of
	 * the response match the source IP address and port that the Binding
	 * Request was sent from.
	 */
	if check.remoteCandidate.addr != from {
		/*
			remote peer reflexive 只能从 binding request 中发现,那里有认证, response 是没有认证的.防止攻击
		*/
		err = fmt.Errorf("check %s got message from unkown address %s", check.key, from)
		s.log.Error(fmt.Sprintf("connectivity check failed,  check:%s remote address mismatch,err:%s", check, err))
		s.changeCheckState(check, checkStateFailed, err)
		s.tryCompleteCheck(check) //is this the last check?
		return
	}
	/* 7.1.2.2.  Success Cases
	 *
	 * A check is considered to be a success if all of the following are
	 * true:
	 *
	 * o  the STUN transaction generated a success response
	 *
	 * o  the source IP address and port of the response equals the
	 *    destination IP address and port that the Binding Request was sent
	 *    to
	 *
	 * o  the destination IP address and port of the response match the
	 *    source IP address and port that the Binding Request was sent from
	 */
	var xaddr stun.XORMappedAddress
	err = xaddr.GetFrom(res)
	if err != nil {
		s.changeCheckState(check, checkStateFailed, err)
		s.tryCompleteCheck(check)
		return
	}
	s.log.Trace(fmt.Sprintf("get xaddr =%s ", xaddr.String()))

	var lcand *Candidate
	for _, c := range s.localCandidates {
		if xaddr.String() == c.addr && c.baseAddr == check.localCandidate.baseAddr {
			lcand = c //在这里已经切换了 check. 根据实际会选择是 reflexive 的连接还是 host,还是 relay
			break
		}
	}
	if lcand == nil {
		/* 7.1.2.2.1.  Discovering Peer Reflexive Candidates
		 * If the transport address returned in XOR-MAPPED-ADDRESS does not match
		 * any of the local candidates that the agent knows about, the mapped
		 * address represents a new candidate - a peer reflexive candidate.
		 */
		foundation := calcFoundation(check.localCandidate.baseAddr)
		lcand = new(Candidate)
		lcand.Foundation = foundation
		lcand.baseAddr = check.localCandidate.baseAddr
		lcand.Type = CandidatePeerReflexive
		lcand.ComponentID = check.localCandidate.ComponentID
		lcand.addr = xaddr.String()
		lcand.transport = check.localCandidate.transport
		lcand.Priority = calcCandidatePriority(lcand.Type, defaultPreference, lcand.ComponentID)
		s.log.Trace(fmt.Sprintf("candidate add peer reflexive :%s", lcand))
		s.localCandidates = append(s.localCandidates, lcand)
	}
	/* 7.1.2.2.3.  Constructing a Valid Pair
	 * Next, the agent constructs a candidate pair whose local candidate
	 * equals the mapped address of the response, and whose remote candidate
	 * equals the destination address to which the request was sent.
	 */

	/* Add pair to valid list, if it's not there, otherwise just update
	 * nominated flag
	 */
	found := false
	var newcheck *sessionCheck
	for _, check2 := range s.validCheckList.checks {
		if check2.localCandidate == lcand && check2.remoteCandidate == check.remoteCandidate {
			found = true
			check2.nominated = check.nominated
			newcheck = check2
			break
		}
	}
	if !found {
		newcheck = &sessionCheck{
			localCandidate:  lcand,
			remoteCandidate: check.remoteCandidate,
			priority:        calcPairPriority(s.role, lcand, check.remoteCandidate),
			state:           checkStateSucced,
			nominated:       check.nominated,
			key:             fmt.Sprintf("%s-%s", lcand.addr, check.remoteCandidate.addr),
		}
		s.validCheckList.checks = append(s.validCheckList.checks, newcheck)
		sort.Sort(s.validCheckList) //todo 为什么要排序呢?看不出来有任何必要
	}
	//find valid check and nominated check
	s.markValidAndNonimated(newcheck)
	/* 7.1.2.2.2.  Updating Pair States
	 *
	 * The agent sets the state of the pair that generated the check to
	 * Succeeded.  The success of this check might also cause the state of
	 * other checks to change as well.
	 */
	s.changeCheckState(check, checkStateSucced, nil)
	/* Perform 7.1.2.2.2.  Updating Pair States.
	 * This may terminate ICE processing.
	 */
	s.tryCompleteCheck(check)
}
func (s *session) markValidAndNonimated(check *sessionCheck) {
	s.mlock.Lock()
	if s.sessionComponent.validCheck == nil || s.sessionComponent.validCheck.priority < check.priority {
		s.sessionComponent.validCheck = check
	}
	if check.nominated {
		if s.sessionComponent.nominatedCheck == nil || s.sessionComponent.nominatedCheck.priority < check.priority {
			s.log.Trace(fmt.Sprintf("old nominatedcheck=%s\n,new nominated=%s", s.sessionComponent.nominatedCheck, check))
			s.sessionComponent.nominatedCheck = check
		}
	}
	s.mlock.Unlock()
}

/*
if all check is failed or success, notify upper layer. return true, when this ice negotiation finished.
*/
func (s *session) tryCompleteCheck(check *sessionCheck) bool {
	/* 7.1.2.2.2.  Updating Pair States
	 *
	 * The agent sets the state of the pair that generated the check to
	 * Succeeded.  The success of this check might also cause the state of
	 * other checks to change as well.  The agent MUST perform the following
	 * two steps:
	 *
	 * 1.  The agent changes the states for all other Frozen pairs for the
	 *     same media stream and same foundation to Waiting.  Typically
	 *     these other pairs will have different component IDs but not
	 *     always.
	 */
	if check.err == nil {
		for _, c := range s.checkList.checks {
			if c.localCandidate.Foundation == check.localCandidate.Foundation && c.state == checkStateFrozen {
				s.changeCheckState(c, checkStateWaiting, nil)
			}
		}
		s.log.Trace(fmt.Sprintf("check  finished:%s", check.String()))
	}

	/* 8.2.  Updating States
	 *
	 * For both controlling and controlled agents, the state of ICE
	 * processing depends on the presence of nominated candidate pairs in
	 * the valid list and on the state of the check list:
	 *
	 * o  If there are no nominated pairs in the valid list for a media
	 *    stream and the state of the check list is Running, ICE processing
	 *    continues.
	 *
	 * o  If there is at least one nominated pair in the valid list:
	 *
	 *    - The agent MUST remove all Waiting and Frozen pairs in the check
	 *      list for the same component as the nominated pairs for that
	 *      media stream
	 *
	 *    - If an In-Progress pair in the check list is for the same
	 *      component as a nominated pair, the agent SHOULD cease
	 *      retransmissions for its check if its pair priority is lower
	 *      than the lowest priority nominated pair for that component
	 */
	if check.err == nil && check.nominated {
		for _, c := range s.checkList.checks {
			if c.state < checkStateInProgress {
				//just fail frozen/waiting check
				s.log.Trace(fmt.Sprintf("check %s to be failed because higher priority check finished.", c.key))
				s.cancelOneCheck(c)
			} else if c.state == checkStateInProgress && c.priority <= check.priority {
				/*
					c.priority<check.priority or <= todo if any error,change to <
						这种策略会尽快结束,但是存在问题,如果低优先级的先完成
						1. 对方可能会收到高优先级的 request, 进而以高优先级为准,如果只有一个 ip 地址,那没什么问题
						2. 应该被采用的高优先级被放弃.
				*/
				/* State is IN_PROGRESS, cancel transaction */
				s.cancelOneCheck(c)
			}
		}
	}
	/* Still in 8.2.  Updating States
	 *
	 * o  Once there is at least one nominated pair in the valid list for
	 *    every component of at least one media stream and the state of the
	 *    check list is Running:
	 *
	 *    *  The agent MUST change the state of processing for its check
	 *       list for that media stream to Completed.
	 *
	 *    *  The agent MUST continue to respond to any checks it may still
	 *       receive for that media stream, and MUST perform triggered
	 *       checks if required by the processing of Section 7.2.
	 *
	 *    *  The agent MAY begin transmitting media for this media stream as
	 *       described in Section 11.1
	 */
	/*
		only one component,so finish
	*/

	/* Note: this is the stuffs that we don't do in 7.1.2.2.2, since our
	 *       ICE session only supports one media stream for now:
	 *
	 * 7.1.2.2.2.  Updating Pair States
	 *
	 * 2.  If there is a pair in the valid list for every component of this
	 *     media stream (where this is the actual number of components being
	 *     used, in cases where the number of components signaled in the SDP
	 *     differs from offerer to answerer), the success of this check may
	 *     unfreeze checks for other media streams.
	 */

	/* 7.1.2.3.  Check List and Timer State Updates
	 * Regardless of whether the check was successful or failed, the
	 * completion of the transaction may require updating of check list and
	 * timer states.
	 *
	 * If all of the pairs in the check list are now either in the Failed or
	 * Succeeded state, and there is not a pair in the valid list for each
	 * component of the media stream, the state of the check list is set to
	 * Failed.
	 */

	/*
	 * See if all checks in the checklist have completed. If we do,
	 * then mark ICE processing as failed.
	 */
	hasNotFinished := false
	for _, c := range s.checkList.checks {
		if c.state < checkStateSucced {
			hasNotFinished = true
			break
		}
	}
	if s.sessionComponent.nominatedCheck != nil { //todo 非 aggressive 模式下,会不会出问题呢?
		s.iceComplete(nil, !hasNotFinished)
		return true
	}
	if !hasNotFinished {
		/* All checks have completed, but we don't have nominated pair.
		 * If agent's role is controlled, check if all components have
		 * valid pair. If it does, this means the controlled agent has
		 * finished the check list and it's waiting for controlling
		 * agent to sendData checks with USE-CANDIDATE flag set.
		 */
		if s.role == SessionRoleControlled {
			if s.sessionComponent.validCheck == nil {
				//todo notify ice failed.
				s.iceComplete(errors.New("no valid check"), true)
				return true
			}
			s.log.Trace(fmt.Sprintf("all checks completed. controlled agent now waits for nomination.."))
			s.changeCompleteResult(sessionCheckComplete)
			go func() {
				//start a timer,failed if there is no nomiated
				time.Sleep(s.controlledAgentWaitNomiatedTimeout) // time from pjnath
				//有可能这个连接已经因为其他原因已经被用户关闭了,要考虑这种可能性.
				if s.sessionComponent.nominatedCheck == nil && !s.hasStopped {
					s.iceComplete(errors.New("no nonimated"), true)
				}
			}()
			return false
		} else if s.isNominating { //如果我是 controlling, 那么总是采用 aggressive策略.
			s.iceComplete(fmt.Errorf("%s controlling no nominated ", s.Name), true)
			return true
		} else {
			/*
				如果我是 regular 模式,那么此时应该是再次发送 bingdingrequest, 并带上 usecandidate, 目前没必要.
			*/
			panic("not implemented")
			//return false
		}

	}
	/* If this connectivity check has been successful, scan all components
	 * and see if they have a valid pair, if we are controlling and we haven't
	 * started our nominated check yet.
	 */
	//目前只有一个 component, 另外只支持 aggressive 模式.
	return false
}
func (s *session) changeCompleteResult(r sessionCompleteResult) {
	if r <= s.completeResult {
		panic(fmt.Sprintf("%s  complete result must increase only, old=%d,new=%d", s.Name, s.completeResult, r))
	}
	s.log.Trace(fmt.Sprintf("change complete result from %d to %d", s.completeResult, r))
	s.completeResult = r
	return
}

/*
关闭除要使用的那个 serversock 以外其他所有的 sock, 因为只有一个是有效的,要使用的.
*/
func (s *session) closeUselessServerSock() {
	for k, srv2 := range s.serverSocks {
		if s.sessionComponent.nominatedServerSock != srv2 {
			delete(s.serverSocks, k)
			srv2.Close()
		}
	}
	if s.sessionComponent.nominatedServerSock != s.turnServerSock {
		s.turnServerSock = nil
	}
}
func (s *session) iceComplete(result error, allcomplete bool) {
	//应该继续允许处理 BindingRequest, 因为对方可能还没有结束.
	s.log.Debug(fmt.Sprintf("icesseion complete ,err:%v,allcomplete=%v", result, allcomplete))
	old := s.completeResult
	if result != nil {
		s.changeCompleteResult(sessionCompleteFailure)
	} else {
		if allcomplete {
			/*
				8.1.1.2. Aggressive Nomination
				With aggressive nomination, the controlling agent includes the USECANDIDATE attribute in every check it sends.
				Once the first check for a component succeeds, it will be added to the valid list and have its
				nominated flag set. When all components have a nominated pair in the valid list, media can begin
				to flow using the highest priority nominated pair. However, because the agent included the USECANDIDATE
				attribute in all of its checks, another check may yet complete, causing another valid pair to have its
				nominated flag set. ICE always selects the highest-priority nominated candidate pair from the valid list
				as the one used for media.
			*/
			/*
				 Consequently, the selected pair may actually change briefly as ICE checks
				complete, resulting in a set of transient selections until it stabilizes.
			*/
			/*
				到这里协商才算真正完毕,后续的 request 请求继续响应,但是我不再发送 request了
			*/
			s.changeCompleteResult(sessionAllCompleteSuccess)
			s.log.Debug(fmt.Sprintf("icesession allcomplete"))
			if len(s.checkMap) != 0 {
				panic("all check should finished")
			}
		} else {
			if s.completeResult < sessionCompleteSuccess {
				s.changeCompleteResult(sessionCompleteSuccess)
			}
		}
		s.log.Trace(fmt.Sprintf("valid check=%s\n nominated=%s\n", s.sessionComponent.validCheck, s.sessionComponent.nominatedCheck))
		srv, err := s.getSenderServerSock(s.sessionComponent.nominatedCheck.localCandidate.addr)
		if err != nil {
			panic(fmt.Sprintf("cannot found nominatedcheck corresponding serversock %s", err))
		}
		s.mlock.Lock()
		s.sessionComponent.nominatedServerSock = srv
		if allcomplete {
			s.closeUselessServerSock()
			s.mlock.Unlock()
			check := s.sessionComponent.nominatedCheck
			if check.localCandidate.Type == CandidateRelay {
				result = s.turnServerSock.channelBind(check.remoteCandidate.addr)
				if result != nil {
					/*
						失败了,不妨碍我继续使用sendIndication 来传输数据,继续这么做吧.
					*/
					s.log.Error(fmt.Sprintf("channel bind err:%s", result))
					//s.iceStreamTransport.State = TransportStateFailed
					//t.Stop()
					srv.FinishNegotiation(turnModeData)
					return
				}
				srv.FinishNegotiation(turnModeData)
			} else {
				srv.FinishNegotiation(stunModeData)
			}
		} else {
			s.mlock.Unlock()
		}
	}
	if old < sessionCompleteSuccess { //只通知上层一次,但是可能完成多次,不断更新状态.
		s.iceStreamTransport.onIceComplete(result)
	}
}

/*
cancel one started check
*/
func (s *session) cancelOneCheck(check *sessionCheck) {
	chr := s.checkMap[check.key]
	chr <- errors.New("canceled")
	s.changeCheckState(check, checkStateFailed, errors.New("canceled"))
}
func (s *session) finishOneCheck(check *sessionCheck) {
	chr := s.checkMap[check.key]
	delete(s.checkMap, check.key)
	close(chr)
}
func (s *session) retryOneCheck(check *sessionCheck) {
	if check.state != checkStateInProgress {
		s.log.Info(fmt.Sprintf("only can retry a check in progress, check=%s", check))
		return
	}
	chr := s.checkMap[check.key]
	chr <- errCheckRetry
}
func (s *session) startCheck() error {
	s.log.Trace(fmt.Sprintf("start ice check..."))
	if s.aggresive && s.role == SessionRoleControlling {
		s.isNominating = true
	}
	if s.checkList.checks[0].state != checkStateFrozen {
		return errors.New("already start another check")
	}
	s.allcheck(s.checkList.checks)
	return nil
}
func (s *session) changeCheckState(check *sessionCheck, newState SessionCheckState, err error) {
	s.log.Trace(fmt.Sprintf("check %s: state changed from %s to %s err:%s", check.key, check.state, newState, err))
	if check.state >= newState {
		s.log.Error(fmt.Sprintf("check state only can increase. newstate=%s,oldState=%s, check=%s", newState, check.state, check))
		return
	}
	check.state = newState
	check.err = err
	//停止探测
	if check.state >= checkStateSucced {
		s.finishOneCheck(check)
	}
}

//启动完毕以后立即返回,结果要从 ice complete中获取.
func (s *session) allcheck(checks []*sessionCheck) {
	const checkInterval = time.Millisecond * 20
	fmap := make(map[int]bool)
	for _, c := range checks {
		key := fmt.Sprintf("%s-%s", c.localCandidate.addr, c.remoteCandidate.addr)
		ch := make(chan error, 1)
		s.checkMap[key] = ch
	}
	/*
		only one compondent, all waiting...
	*/
	for _, c := range checks {
		if !fmap[c.localCandidate.Foundation] {
			fmap[c.localCandidate.Foundation] = true
			s.changeCheckState(c, checkStateWaiting, nil)
		}

	}
	for _, rc := range s.earlyCheckList {
		/*
			优先处理收到的请求,可能已经可以成功了.
		*/
		s.log.Trace(fmt.Sprintf("process early check list %s", rc))
		s.handleIncomingCheck(rc)
	}
	for _, c := range checks {
		ch := s.checkMap[c.key]
		//有可能还没有启动,其他 check 已经完毕,这个就没有必要了.
		if c.state == checkStateWaiting {
			s.changeCheckState(c, checkStateInProgress, nil)
			go s.onecheck(c, ch, s.isNominating)
			time.Sleep(checkInterval)
		}
	}
	/* If we don't have anything in Waiting state, perform check to
	 * highest priority pair that is in Frozen state.
	 */
	for _, c := range checks {
		ch := s.checkMap[c.key]
		//有可能还没有启动,其他 check 已经完毕,这个就没有必要了.
		if c.state == checkStateFrozen {
			s.changeCheckState(c, checkStateInProgress, nil)
			go s.onecheck(c, ch, s.isNominating)
			time.Sleep(checkInterval)
		}
	}
}
func (s *session) buildBindingRequest(c *sessionCheck) (req *stun.Message) {
	var (
		err      error
		priority attr.Priority
		control  stun.Setter
		prio     int
		setters  []stun.Setter
	)
	req = new(stun.Message)
	prio = calcCandidatePriority(CandidatePeerReflexive, defaultPreference, 1)
	priority = attr.Priority(prio)
	if s.role == SessionRoleControlling {
		control = attr.IceControlling(s.tieBreaker)
	} else {
		control = attr.IceControlled(s.tieBreaker)
	}
	setters = []stun.Setter{stun.BindingRequest,
		stun.TransactionIDSetter,
		priority, control,
		software,
		stun.Username(s.txUserName),
		s.txCrendientials,
		stun.Fingerprint}
	if c.nominated && s.role == SessionRoleControlling {
		//useCandidate 不能放在最后,
		setters = append([]stun.Setter{attr.UseCandidate}, setters...)
	}
	err = req.Build(setters...)
	if err != nil {
		panic("build error...")
	}
	return
}
func (s *session) getSenderServerSock(localAddr string) (ss serverSocker, err error) {
	srv, ok := s.serverSocks[localAddr]
	if ok {
		return srv, nil
	}
	for _, c := range s.localCandidates {
		if c.addr == localAddr && c.Type == CandidateRelay {
			return s.turnServerSock, nil
		} else if c.addr == localAddr && c.Type == CandidateServerReflexive {
			ss = s.serverSocks[c.baseAddr]
			return
		} else if c.addr == localAddr && c.Type == CandidatePeerReflexive {
			ss = s.serverSocks[c.baseAddr]
			return
		}
	}
	err = fmt.Errorf("%s localadr=%s,cannot found in maps=%#v", s.Name, localAddr, s.serverSocks)
	return nil, err
}

func calcRetransmitTimeout(count int, lastsleep time.Duration) time.Duration {
	if count == 0 {
		return defaultRTOValue
	} else if count < maxRetryBindingRequest-1 {
		return lastsleep * 2
	} else {
		return stunTimeoutValue
	}

}
func (s *session) onecheck(c *sessionCheck, chCheckResult chan error, nominate bool) {
	var (
		err        error
		req        *stun.Message
		sleep      time.Duration = defaultRTOValue
		serversock serverSocker
	)
	s.log.Trace(fmt.Sprintf("start check %s", c.key))
	defer func() {
		s.log.Trace(fmt.Sprintf("check complete %s", c.key))
		if req != nil {
			s.deleteMsgCheck(req.TransactionID)
		}
	}()
	serversock, err = s.getSenderServerSock(c.localCandidate.addr)
	if err != nil && s.completeResult < sessionAllCompleteSuccess {
		s.log.Error(err.Error())
		return
	}
	if nominate && s.role == SessionRoleControlling {
		c.nominated = true
	}
	//build req message
lblRestart:
	req = s.buildBindingRequest(c)
	for i := 0; i < maxRetryBindingRequest; i++ {
		s.log.Trace(fmt.Sprintf("%s sendData %d times,bindingrequestlength=%d", c.key, i+1, len(req.Raw)))
		sleep = calcRetransmitTimeout(i, sleep)
		s.addMsgCheck(req.TransactionID, c)

		err = serversock.sendStunMessageAsync(req, c.localCandidate.addr, c.remoteCandidate.addr)
		if err != nil {
			s.log.Debug(fmt.Sprintf("send binding request from %s to %s ,err %s", c.localCandidate.addr, c.remoteCandidate.addr, err))
		}
		select {
		case <-time.After(sleep):
			continue
		case err = <-chCheckResult:
			if err == errCheckRetry {
				goto lblRestart // 立即进行下一次探测.
			}
			return
		}
	}
	//探测了七次,没有任何结果,失败.
	s.tryFailChan <- &checkFailedWrapper{c, errTriedTooManyTimes}
	return
}

func (s *session) changeRole(newrole SessionRole) {
	s.log.Trace(fmt.Sprintf("role changed from %s to %s", s.role, newrole))
	s.role = newrole
}
func (s *session) sendResponse(localAddr, fromAddr string, req *stun.Message, code stun.ErrorCode) {
	var (
		err         error
		res         = new(stun.Message)
		fromUDPAddr *net.UDPAddr
	)
	fromUDPAddr = addrToUDPAddr(fromAddr)
	sc, err := s.getSenderServerSock(localAddr)
	if err != nil {
		/*
			ignore
			partner should have negotiation complete,
		*/
		if s.completeResult < sessionAllCompleteSuccess {
			s.log.Error(err.Error())
		}
		return
	}
	if code == 0 {
		err = res.Build(
			stun.NewTransactionIDSetter(req.TransactionID),
			stun.NewType(stun.MethodBinding, stun.ClassSuccessResponse),
			software,
			&stun.XORMappedAddress{
				IP:   fromUDPAddr.IP,
				Port: fromUDPAddr.Port,
			},
			s.txCrendientials,
			stun.Fingerprint,
		)
		if err != nil {
			panic(fmt.Sprintf("build res message error %s", err))
		}
		sc.sendStunMessageAsync(res, localAddr, fromAddr)
		return
	} else if code == stun.CodeRoleConflict {
		err = res.Build(
			stun.NewTransactionIDSetter(req.TransactionID),
			stun.NewType(stun.MethodBinding, stun.ClassErrorResponse),
			software,
			stun.CodeRoleConflict,
			&stun.XORMappedAddress{
				IP:   fromUDPAddr.IP,
				Port: fromUDPAddr.Port,
			},
			s.txCrendientials,
			stun.Fingerprint,
		)
		if err != nil {
			panic(fmt.Sprintf("build res message error %s", err))
		}
		sc.sendStunMessageAsync(res, localAddr, fromAddr)
		return
	} else if code == stun.CodeUnauthorised {
		res.Build(stun.NewTransactionIDSetter(req.TransactionID), stun.BindingError,
			stun.CodeUnauthorised, software, s.txCrendientials, stun.Fingerprint)
		sc.sendStunMessageAsync(res, localAddr, fromAddr)
	}
}

//binding request 和普通的 stun message 一样处理.
func (s *session) processBindingRequest(localAddr, fromAddr string, req *stun.Message) {
	var (
		err         error
		hasControll = false
		rcheck      = new(rxCheck)
		priority    attr.Priority
	)
	var userName stun.Username
	s.log.Trace(fmt.Sprintf("received binding request  %s<----------%s %s", localAddr, fromAddr, hex.EncodeToString(req.TransactionID[:])))
	err = priority.GetFrom(req)
	if err != nil {
		s.log.Info(fmt.Sprintf("stun bind request has no priority,ingored."))
		return
	}
	rcheck.priority = int(priority)
	err = userName.GetFrom(req)
	if err != nil {
		s.log.Warn(fmt.Sprintf("%s received bind request  with no username %s", localAddr, err))
		s.sendResponse(localAddr, fromAddr, req, stun.CodeUnauthorised)
		return
	}
	if len(s.rxUserName) > 0 {
		if userName.String() != s.rxUserName {
			s.log.Warn(fmt.Sprintf("%s received bind request ,but user name not match expect=%s,got=%s", localAddr, s.rxUserName, userName.String()))
			s.sendResponse(localAddr, fromAddr, req, stun.CodeUnauthorised)
			return
		}

	}
	/*
		必须进行权限检查,以防止收到错误的消息
	*/
	err = s.rxCrendientials.Check(req)
	if err != nil {
		s.log.Warn(fmt.Sprintf("%s received bind request ,crendientials check failed %s", localAddr, err))
		s.sendResponse(localAddr, fromAddr, req, stun.CodeUnauthorised)
		return
	}
	_, err = req.Get(stun.AttrICEControlling)
	if err == nil {
		hasControll = true
		rcheck.role = SessionRoleControlling
		if s.role != SessionRoleControlled {
			var peerTieBreaker attr.IceControlling
			peerTieBreaker.GetFrom(req)
			/*
				tiebreaker, 谁的大以谁的为准.
			*/
			if s.tieBreaker < uint64(peerTieBreaker) {
				s.changeRole(SessionRoleControlled)
			} else {
				s.sendResponse(localAddr, fromAddr, req, stun.CodeRoleConflict)
				return
			}

		}
	}
	_, err = req.Get(stun.AttrICEControlled)
	if err == nil {
		hasControll = true
		rcheck.role = SessionRoleControlled
		if s.role != SessionRoleControlling {
			var peerTieBreaker attr.IceControlled
			peerTieBreaker.GetFrom(req)
			if s.tieBreaker < uint64(peerTieBreaker) {
				s.changeRole(SessionRoleControlling)
			} else {
				s.sendResponse(localAddr, fromAddr, req, stun.CodeRoleConflict)
				return
			}

		}
	}
	if !hasControll {
		s.log.Info(fmt.Sprintf("received stun binding request,but no icecontrolling and icecontrolled"))
		s.sendResponse(localAddr, fromAddr, req, stun.CodeUnauthorised)
		return
	}
	/*
		如果是 earlycheck, 那么发送过去的 response 中 username 应该是错的,所以我们不能认为 username 不对就是错的.
	*/
	s.sendResponse(localAddr, fromAddr, req, 0)
	if s.completeResult >= sessionAllCompleteSuccess {
		return // 不应该继续处理了,因为negotiation 已经完成了.
	}
	//early check received.
	if len(s.checkMap) <= 0 && s.completeResult == sessionNotComplete {
		s.rxUserName = string(userName)
		s.log.Info(fmt.Sprintf("received early check from %s, username=%s", fromAddr, s.rxUserName))
	}
	/*
	 * Handling early check.
	 *
	 * It's possible that we receive this request before we receive SDP
	 * answer. In this case, we can't perform trigger check since we
	 * don't have checklist yet, so just save this check in a pending
	 * triggered check array to be acted upon later.
	 */
	//init check
	_, err = req.Get(stun.AttrUseCandidate)
	if err == nil {
		rcheck.userCandidate = true
	}
	rcheck.componentID = 1
	rcheck.remoteAddress = fromAddr
	rcheck.localAddress = localAddr
	if len(s.checkMap) <= 0 && s.completeResult == sessionNotComplete { //checkmap为空表示我还没开始协商,当然也可能是我已经把所有的 check 都检查完了.
		/*
			We don't have answer yet, so keep this request for later
		*/
		s.earlyCheckList = append(s.earlyCheckList, rcheck)
	} else {
		//其他阶段忽略,我已经选定了用于通信的 check
		s.handleIncomingCheck(rcheck)

	}
}

/* Handle incoming Binding request and perform triggered check.
 * This function may be called by processBindingRequest, or when
 * SDP answer is received and we have received early checks.
 */

func (s *session) handleIncomingCheck(rcheck *rxCheck) {
	var (
		lcand *Candidate
		rcand *Candidate
	)
	/* 7.2.1.3.  Learning Peer Reflexive Candidates
	 * If the source transport address of the request does not match any
	 * existing remote candidates, it represents a new peer reflexive remote
	 * candidate.
	 */
	for _, c := range s.remoteCandidates {
		if c.addr == rcheck.remoteAddress {
			rcand = c
			break
		}
	}
	if rcand == nil {
		if len(s.remoteCandidates) > maxCandidates {
			s.log.Warn(fmt.Sprintf("unable to add new peer reflexive candidate: too many candidates ."))
			return
		}
		rcand = new(Candidate)
		rcand.ComponentID = 1
		rcand.Type = CandidatePeerReflexive
		rcand.Priority = rcheck.priority
		rcand.addr = rcheck.remoteAddress
		rcand.Foundation = calcFoundation(rcheck.remoteAddress)
		s.remoteCandidates = append(s.remoteCandidates, rcand)
		s.log.Info(fmt.Sprintf("add new remote candidate from the request %s", rcand.addr))
	}
	/*
		寻找匹配这个 rcheck 的 localCandidates, 就找优先级最高的那个就可以了.
	*/
	for _, cand := range s.localCandidates {
		if cand.addr == rcheck.localAddress {
			lcand = cand
			break
		}
	}
	/*
	 * Create candidate pair for this request.
	 */

	/*
	 * 7.2.1.4.  Triggered Checks
	 *
	 * Now that we have local and remote candidate, check if we already
	 * have this pair in our checklist.
	 */
	var c *sessionCheck
	for _, chk := range s.checkList.checks {
		if chk.localCandidate == lcand && chk.remoteCandidate == rcand {
			c = chk
			break
		}
	}
	/* If the pair is already on the check list:
	 * - If the state of that pair is Waiting or Frozen, its state is
	 *   changed to In-Progress and a check for that pair is performed
	 *   immediately.  This is called a triggered check.
	 *
	 * - If the state of that pair is In-Progress, the agent SHOULD
	 *   generate an immediate retransmit of the Binding Request for the
	 *   check in progress.  This is to facilitate rapid completion of
	 *   ICE when both agents are behind NAT.
	 *
	 * - If the state of that pair is Failed or Succeeded, no triggered
	 *   check is sent.
	 */
	if c != nil {
		oldnominated := c.nominated
		c.nominated = rcheck.userCandidate || c.nominated
		s.log.Trace(fmt.Sprintf("change check %s nominated from %v to %v", c.key, oldnominated, c.nominated))
		if c.state == checkStateFrozen || c.state == checkStateWaiting {
			s.log.Trace(fmt.Sprintf("performing triggered check for %s", c.key))
			chResult, ok := s.checkMap[c.key]
			if !ok {
				panic("must ...")
			}
			s.changeCheckState(c, checkStateInProgress, nil)
			go s.onecheck(c, chResult, c.nominated || s.isNominating)
		} else if c.state == checkStateInProgress {
			//Should retransmit immediately
			s.log.Trace(fmt.Sprintf("triggered check for check %s not performed, because its in progress. Retransmitting", c.key))
			s.retryOneCheck(c)
		} else if c.state == checkStateSucced {
			if rcheck.userCandidate {
				for _, vc := range s.validCheckList.checks {
					if vc.remoteCandidate == c.remoteCandidate {
						vc.nominated = true
						s.markValidAndNonimated(vc)
						s.log.Trace(fmt.Sprintf("valid check %s is nominated", vc.key))
					}
				}
			}
			s.log.Trace(fmt.Sprintf("triggered check for check %s not performed because it's completed", c.key))
			complete := s.tryCompleteCheck(c)
			if complete {
				return
			}
		}
	} else {
		/* If the pair is not already on the check list:
		 * - The pair is inserted into the check list based on its priority.
		 * - Its state is set to In-Progress
		 * - A triggered check for that pair is performed immediately.
		 */
		/* Note: only do this if we don't have too many checks in checklist */
		c := &sessionCheck{
			localCandidate:  lcand,
			remoteCandidate: rcand,
			priority:        calcPairPriority(s.role, lcand, rcand),
			state:           checkStateInProgress,
			nominated:       rcheck.userCandidate,
			key:             fmt.Sprintf("%s-%s", lcand.addr, rcand.addr),
		}
		s.checkList.checks = append(s.checkList.checks, c)
		ch := make(chan error, 1)
		s.checkMap[c.key] = ch
		nominated := c.nominated || s.isNominating
		go s.onecheck(c, ch, nominated)
		s.log.Trace(fmt.Sprintf("New triggered check added:%s", c.key))
	}
}

/*
在localAddr 收到了来自remoteAddr 的 stun binding response
注意这里的 localAddr 并不是代表本机地址,
代表的是SessionCheck 中的 localCandidate
*/
func (s *session) processBindingResponse(localAddr, remoteAddr string, msg *stun.Message) {
	id := msg.TransactionID
	check := s.getMsgCheck(id)
	if check == nil {
		s.log.Info(fmt.Sprintf("receive bind response ,but has no related check %s", msg))
		return
	}
	if check.localCandidate.addr != localAddr {
		s.log.Warn(fmt.Sprintf("received bind response ,but local addr err ,expect %s,got %s", check.localCandidate.addr, localAddr))
		return
	}
	if check.state >= checkStateSucced {
		s.log.Info(fmt.Sprintf("check %s has been finished", check.key))
		return
	}
	s.handleCheckResponse(check, remoteAddr, msg)
}

/*
ice 协商的核心就是处理
1. 收到的 binding Request msgChan
2. 收到的 binding response  msgChan
3. 自己不断的发送binding Request 发送完毕对应的 tryFailChan
4. 协商找到可用连接以后,收发数据. 收数据用dataChan

*/
func (s *session) loop() {
	for {
		r := utils.RandomString(20)
		s.log.Trace(fmt.Sprintf("session loop %s start @%s", r, time.Now().Format("15:04:05.999")))
		select {
		case msg, ok := <-s.msgChan:
			if ok {
				s.processStunMessage(msg.localAddr, msg.remoteAddr, msg.msg)
			} else {
				return
			}
		case data, ok := <-s.dataChan:
			if ok {
				s.iceStreamTransport.onRxData(data.data, data.remoteAddr)
			} else {
				return
			}
		case c, ok := <-s.tryFailChan:
			if ok {
				if c.c.state < checkStateSucced {
					//可能已经成功了,
					s.changeCheckState(c.c, checkStateFailed, c.err)
					s.tryCompleteCheck(c.c)
				}
			} else {
				return
			}
		case <-s.quitChan:
			return
		}
		s.log.Trace(fmt.Sprintf("loop %s end @%s", r, time.Now().Format("15:04:05.999")))
	}
}

/*
ice 协商只应该收到 binding response 和 bindingRequest
其他消息都应该是某种错误,或者恶意攻击.
*/
func (s *session) processStunMessage(localAddr, remoteAddr string, msg *stun.Message) {
	if msg.Type == stun.BindingRequest {
		s.processBindingRequest(localAddr, remoteAddr, msg)
		return
	}
	//binding response?
	if msg.Type == stun.BindingError || msg.Type == stun.BindingSuccess {
		s.processBindingResponse(localAddr, remoteAddr, msg)
		return
	}
	s.log.Warn(fmt.Sprintf("%s receive unexpected stun message from  %s, msg:%s", localAddr, remoteAddr, msg.Type))
}

/*
message received from peer or stun server after negiotiation complete.
*/
func (s *session) RecieveStunMessage(localAddr, remoteAddr string, msg *stun.Message) {
	s.log.Trace(fmt.Sprintf("%s receive stun message from  %s, msg:%s", localAddr, remoteAddr, msg.Type))
	if s.hasStopped {
		return
	}
	//不要阻塞发送接收消息线程.
	s.msgChan <- &stunMessageWrapper{localAddr, remoteAddr, msg}
	return

}

/*
1. 在协商未完全结束之前就有可能收到数据,只要有一个可用的连接,对方就会发送数据,
2. 随着协商的完成,最终双方会确认一个一致的 check, 如果这时候是走的 relay, 那么才会启用 turn channel 模式.
*/
func (s *session) ReceiveData(localAddr, peerAddr string, data []byte) {
	if s.hasStopped {
		return
	}
	s.log.Trace(fmt.Sprintf("recevied data %s<-----%s l=%d", localAddr, peerAddr, len(data)))
	s.dataChan <- &stunDataWrapper{localAddr, peerAddr, data}
	return

}
func (s *session) SendData(data []byte) error {
	s.mlock.Lock()
	//nominiatedcheck可能会在可以发送数据以后变化.
	check := s.sessionComponent.nominatedCheck
	srv := s.sessionComponent.nominatedServerSock
	fromaddr := check.localCandidate.addr
	s.mlock.Unlock()
	if check == nil {
		return errors.New("no check")
	}
	if srv == nil {
		return errors.New("no stun transport")
	}
	s.log.Trace(fmt.Sprintf("send data from %s to %s datalen=%d", fromaddr, check.remoteCandidate.addr, len(data)))
	if check.localCandidate.Type == CandidateServerReflexive || check.localCandidate.Type == CandidatePeerReflexive {
		fromaddr = check.localCandidate.baseAddr
		s.log.Trace(fmt.Sprintf("accutally send data from %s to %s datalen=%d", fromaddr, check.remoteCandidate.addr, len(data)))
	}
	return srv.sendData(data, fromaddr, check.remoteCandidate.addr)
}

/*
pair priority = 2^32*MIN(G,D) + 2*MAX(G,D) + (G>D?1:0)
*/
func calcPairPriority(role SessionRole, l, r *Candidate) uint64 {
	var o, a uint32
	var min, max uint32
	if role == SessionRoleControlling {
		o = uint32(l.Priority)
		a = uint32(r.Priority)
	} else {
		o = uint32(r.Priority)
		a = uint32(l.Priority)
	}
	if o > a {
		min = a
		max = o
	} else {
		min = o
		max = a
	}
	var p uint64
	p = uint64(min) << 32
	max = max << 1
	p += uint64(max)
	if o > a {
		p++
	}
	return p
}
