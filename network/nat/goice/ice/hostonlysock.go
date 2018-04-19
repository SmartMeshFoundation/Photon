package ice

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

/*
没有指定 turn /stun server 也是允许的,
*/
type HostOnlySock struct {
	localCandidates []string
}

func (h *HostOnlySock) GetCandidates() (candidates []*Candidate, err error) {
	addrs, err := DefaultGatherer.Gather()
	if err != nil {
		return
	}
	if len(addrs) < 0 {
		//no ip
		err = errors.New("no network")
	}
	port := rand.NewSource(time.Now().UnixNano()).Int63() % 50000
	primaryAddress := fmt.Sprintf("%s:%d", addrs[0].IP.String(), port)
	candidates, err = GetLocalCandidates(primaryAddress)
	if err != nil {
		return
	}
	for _, c := range candidates {
		h.localCandidates = append(h.localCandidates, c.addr)
	}
	return
}
func (h *HostOnlySock) Close() {

}

/*
address need to listen for input stun binding request...
*/
func (h *HostOnlySock) GetListenCandidiates() []string {
	return h.localCandidates
}
