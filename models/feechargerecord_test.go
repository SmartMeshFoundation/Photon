package models

import (
	"testing"

	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/Photon/utils"
)

func TestModelDB_FeeChargeRecord(t *testing.T) {
	model, err := newTestDb()
	if err != nil {
		t.Error(err.Error())
		return
	}
	lockSecretHash1 := utils.NewRandomHash()
	lockSecretHash2 := utils.NewRandomHash()
	r1 := &FeeChargeRecord{
		Fee:            big.NewInt(1),
		LockSecretHash: lockSecretHash1,
	}

	all, err := model.GetAllFeeChargeRecord()
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	err = model.SaveFeeChargeRecord(r1.ToSerialized())
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	all, err = model.GetAllFeeChargeRecord()
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

	err = model.SaveFeeChargeRecord(r1.ToSerialized())
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	all, err = model.GetAllFeeChargeRecord()
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

	all, err = model.GetFeeChargeRecordByLockSecretHash(lockSecretHash2)
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
