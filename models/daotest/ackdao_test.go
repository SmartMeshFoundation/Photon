package daotest

import (
	"testing"

	"reflect"

	"github.com/SmartMeshFoundation/Photon/codefortest"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
)

func TestAckDao(t *testing.T) {
	dao := codefortest.NewTestDB("")
	defer dao.CloseDB()
	echoHash := utils.NewRandomHash()
	tx := dao.StartTx(models.BucketAck)
	// save with tx and get
	dao.SaveAck(echoHash, echoHash.Bytes(), tx)
	tx.Commit()
	r1 := dao.GetAck(echoHash)
	if !reflect.DeepEqual(r1, echoHash.Bytes()) {
		t.Error("not equal")
		return
	}
	// save and get
	echoHash = utils.NewRandomHash()
	dao.SaveAckNoTx(echoHash, echoHash.Bytes())
	r2 := dao.GetAck(echoHash)
	if !reflect.DeepEqual(r2, echoHash.Bytes()) {
		t.Error("not equal")
		return
	}
}
