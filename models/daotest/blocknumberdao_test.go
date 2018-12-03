package daotest

import (
	"testing"

	"time"

	"github.com/SmartMeshFoundation/Photon/codefortest"
)

func TestBlockNumberDao(t *testing.T) {
	dao := codefortest.NewTestDB("")
	bn := int64(500)
	dao.SaveLatestBlockNumber(bn)
	bn1 := dao.GetLatestBlockNumber()
	if bn1 != bn {
		t.Error("not equal")
		return
	}
	saveTime := dao.GetLastBlockNumberTime()
	if time.Since(saveTime) > 2*time.Second {
		t.Error("wrong time")
		return
	}
}
