package ice

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestNewStunSocket(t *testing.T) {
	stun, err := newStunSocket("182.254.155.208:3478")
	if err != nil {
		t.Error(err)
		return
	}
	cands, err := stun.GetCandidates()
	if err != nil {
		t.Error(err)
		return
	}
	spew.Dump("cands", cands)
}
