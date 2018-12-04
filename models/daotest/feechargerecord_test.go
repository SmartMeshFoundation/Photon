package daotest

import (
	"testing"

	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/Photon/codefortest"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
)

func TestModelDB_FeeChargeRecord(t *testing.T) {
	dao := codefortest.NewTestDB("")
	defer dao.CloseDB()
	lockSecretHash1 := utils.NewRandomHash()
	lockSecretHash2 := utils.NewRandomHash()
	r1 := &models.FeeChargeRecord{
		Fee:            big.NewInt(1),
		LockSecretHash: lockSecretHash1,
	}

	all, err := dao.GetAllFeeChargeRecord()
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	err = dao.SaveFeeChargeRecord(r1)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	all, err = dao.GetAllFeeChargeRecord()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if len(all) != 1 {
		t.Error("wrong data")
		return
	}
	for _, r := range all {
		fmt.Println(r.ToString())
	}

	r1.Key = utils.EmptyHash
	r1.Fee = big.NewInt(2)
	r1.LockSecretHash = lockSecretHash2

	err = dao.SaveFeeChargeRecord(r1)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	all, err = dao.GetAllFeeChargeRecord()
	fmt.Println(len(all))
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if len(all) != 2 {
		t.Error("wrong data")
		return
	}
	for _, r := range all {
		fmt.Println(r.ToString())
	}

	all, err = dao.GetFeeChargeRecordByLockSecretHash(lockSecretHash2)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	fmt.Println(len(all))
	if len(all) != 1 {
		t.Error("wrong data")
		return
	}
	for _, r := range all {
		fmt.Println(r.ToString())
	}
}
