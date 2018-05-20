package ice

import (
	"bytes"
	"fmt"
)

/*
SessionCheckState describes the state of ICE check.
*/
type SessionCheckState int

const (
	/**
	 * A check for this pair hasn't been performed, and it can't
	 * yet be performed until some other check succeeds, allowing this
	 * pair to unfreeze and move into the Waiting state.
	 */
	checkStateFrozen SessionCheckState = iota
	/**
	 * A check has not been performed for this pair, and can be
	 * performed as soon as it is the highest priority Waiting pair on
	 * the check list.
	 */
	checkStateWaiting
	/**
	 * A check has not been performed for this pair, and can be
	 * performed as soon as it is the highest priority Waiting pair on
	 * the check list.
	 */
	checkStateInProgress
	/**
	 * A check has not been performed for this pair, and can be
	 * performed as soon as it is the highest priority Waiting pair on
	 * the check list.
	 */
	checkStateSucced
	/**
	 * A check for this pair was already done and failed, either
	 * never producing any response or producing an unrecoverable failure
	 * response.
	 */
	checkStateFailed
)

func (s SessionCheckState) String() string {
	switch s {
	case checkStateFrozen:
		return "frozen"
	case checkStateWaiting:
		return "waiting"
	case checkStateInProgress:
		return "inprogress"
	case checkStateSucced:
		return "success"
	case checkStateFailed:
		return "failed"
	}
	return "unknown"
}

/**
 * This structure describes an ICE connectivity check. An ICE check
 * contains a candidate pair, and will involve sending STUN Binding
 * Request transaction for the purposes of verifying connectivity.
 * A check is sent from the local candidate to the remote candidate
 * of a candidate pair.
 */
type sessionCheck struct {
	localCandidate  *Candidate
	remoteCandidate *Candidate
	key             string //简单与其他 check 区分,更多用于调试.
	priority        uint64
	state           SessionCheckState
	/**
	 * Flag to indicate whether this check is nominated. A nominated check
	 * contains USE-CANDIDATE attribute in its STUN Binding request.
	 */
	nominated bool
	/*
		what error
	*/
	err error
}

func (s *sessionCheck) String() string {
	return fmt.Sprintf("{l=%s,r=%s,priorit=%x,state=%s,nominated=%v,err=%s}",
		s.localCandidate.addr, s.remoteCandidate.addr, s.priority, s.state, s.nominated, s.err)
}

type sessionCheckList struct {
	checks []*sessionCheck
}

func (sc *sessionCheckList) String() string {
	w := new(bytes.Buffer)
	for i, v := range sc.checks {
		fmt.Fprintf(w, "\t [%d]=%s\n", i, v)
	}
	fmt.Fprintf(w, "}")
	return w.String()
}
func (sc *sessionCheckList) Len() int {
	return len(sc.checks)
}
func (sc *sessionCheckList) Less(i, j int) bool {
	return sc.checks[i].priority > sc.checks[j].priority
}

func (sc *sessionCheckList) Swap(i, j int) {
	var t *sessionCheck
	t = sc.checks[i]
	sc.checks[i] = sc.checks[j]
	sc.checks[j] = t
}
