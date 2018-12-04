package daotest

import (
	"reflect"
	"testing"

	"github.com/SmartMeshFoundation/Photon/codefortest"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
)

func TestTX(t *testing.T) {
	dao := codefortest.NewTestDB("")
	defer dao.CloseDB()
	// tx commit
	echoHash := utils.NewRandomHash()
	tx := dao.StartTx(models.BucketAck)
	dao.SaveAck(echoHash, echoHash.Bytes(), tx)
	tx.Commit()
	r1 := dao.GetAck(echoHash)
	if !reflect.DeepEqual(r1, echoHash.Bytes()) {
		t.Error("not equal")
		return
	}
	// tx rollback
	tx2 := dao.StartTx(models.BucketAck)
	echoHash = utils.NewRandomHash()
	dao.SaveAck(echoHash, echoHash.Bytes(), tx2)
	tx.Rollback()
	r2 := dao.GetAck(echoHash)
	if r2 != nil {
		t.Error("should nil")
		return
	}
}
