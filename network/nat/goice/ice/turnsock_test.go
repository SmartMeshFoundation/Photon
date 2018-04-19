package ice

import (
	"fmt"
	"testing"
)

func newTestTurnSock() (turn *TurnSock) {
	turn, err := NewTurnSock("182.254.155.208:3478", "bai", "bai")
	if err != nil {
		panic(err)
		return
	}
	return turn
}
func TestNewTurnSock(t *testing.T) {
	turn, err := NewTurnSock("182.254.155.208:3478", "bai", "bai")
	if err != nil {
		t.Error(err)
		return
	}
	cands, err := turn.GetCandidates()
	if err != nil {
		t.Error(err)
		return
	}
	for i, c := range cands {
		t.Log(fmt.Sprintf("cands[%d]=%s", i, c))
	}
}
