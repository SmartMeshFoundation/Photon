package ice

import (
	"bytes"
	"fmt"
)

/**
 * This enumeration describes the state of ICE check.
 */
type SessionCheckState int

const (
	/**
	 * A check for this pair hasn't been performed, and it can't
	 * yet be performed until some other check succeeds, allowing this
	 * pair to unfreeze and move into the Waiting state.
	 */
	CheckStateFrozen SessionCheckState = iota
	/**
	 * A check has not been performed for this pair, and can be
	 * performed as soon as it is the highest priority Waiting pair on
	 * the check list.
	 */
	CheckStateWaiting
	/**
	 * A check has not been performed for this pair, and can be
	 * performed as soon as it is the highest priority Waiting pair on
	 * the check list.
	 */
	CheckStateInProgress
	/**
	 * A check has not been performed for this pair, and can be
	 * performed as soon as it is the highest priority Waiting pair on
	 * the check list.
	 */
	CheckStateSucced
	/**
	 * A check for this pair was already done and failed, either
	 * never producing any response or producing an unrecoverable failure
	 * response.
	 */
	CheckStateFailed
)

func (s SessionCheckState) String() string {
	switch s {
	case CheckStateFrozen:
		return "frozen"
	case CheckStateWaiting:
		return "waiting"
	case CheckStateInProgress:
		return "inprogress"
	case CheckStateSucced:
		return "success"
	case CheckStateFailed:
		return "failed"
	}
	return "unknown"
}

type SessionCheckListState int

const (
	/**
	* The checklist is not yet running.
	 */
	CHECKLIST_STATE_IDLE SessionCheckListState = iota
	/**
	 * In this state, ICE checks are still in progress for this
	 * media stream.
	 */
	CHECKLIST_STATE_RUNNING
	/**
	 * In this state, ICE checks have completed for this media stream,
	 * either successfully or with failure.
	 */
	CHECKLIST_STATE_COMPLETED
)

func (s SessionCheckListState) String() string {
	switch s {
	case CHECKLIST_STATE_IDLE:
		return "idle"
	case CHECKLIST_STATE_RUNNING:
		return "running"
	case CHECKLIST_STATE_COMPLETED:
		return "completed"
	}
	return "unkown"
}

/**
 * This structure describes an ICE connectivity check. An ICE check
 * contains a candidate pair, and will involve sending STUN Binding
 * Request transaction for the purposes of verifying connectivity.
 * A check is sent from the local candidate to the remote candidate
 * of a candidate pair.
 */
type SessionCheck struct {
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

func (s *SessionCheck) String() string {
	return fmt.Sprintf("{l=%s,r=%s,priorit=%x,state=%s,nominated=%v,err=%s}",
		s.localCandidate.addr, s.remoteCandidate.addr, s.priority, s.state, s.nominated, s.err)
}

type SessionCheckList struct {
	state  SessionCheckListState
	checks []*SessionCheck
}

func (s *SessionCheckList) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "{\nstate=%s\n", s.state)
	for i, v := range s.checks {
		fmt.Fprintf(w, "\t [%d]=%s\n", i, v)
	}
	fmt.Fprintf(w, "}")
	return w.String()
}
func (sc *SessionCheckList) Len() int {
	return len(sc.checks)
}
func (sc *SessionCheckList) Less(i, j int) bool {
	return sc.checks[i].priority > sc.checks[j].priority
}

func (sc *SessionCheckList) Swap(i, j int) {
	var t *SessionCheck
	t = sc.checks[i]
	sc.checks[i] = sc.checks[j]
	sc.checks[j] = t
}
