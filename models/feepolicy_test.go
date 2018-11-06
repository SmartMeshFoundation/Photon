package models

import (
	"testing"

	"math/big"

	"github.com/SmartMeshFoundation/Photon/codefortest"
)

func TestModelDB_FeePolicy(t *testing.T) {
	model, err := codefortest.GetDb()
	if err != nil {
		t.Error(err.Error())
		return
	}

	defaultFp := model.GetFeePolicy()
	if defaultFp.AccountFee.FeeConstant.Int64() != 0 {
		t.Error("wrong fee constant")
		return
	}
	if defaultFp.AccountFee.FeePercent != 10000 {
		t.Error("wrong fee rate")
		return
	}

	defaultFp.AccountFee.FeeConstant = big.NewInt(5)
	defaultFp.AccountFee.FeePercent = 50000

	err = model.SaveFeePolicy(defaultFp)
	if err != nil {
		t.Error(err)
		return
	}

	if defaultFp.AccountFee.FeeConstant.Int64() != 5 {
		t.Error("wrong fee constant")
		return
	}
	if defaultFp.AccountFee.FeePercent != 50000 {
		t.Error("wrong fee rate")
		return
	}
}
