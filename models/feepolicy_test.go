package models

import (
	"testing"

	"math/big"
)

func TestModelDB_FeePolicy(t *testing.T) {
	model := setupDb(t)

	defaultFp := model.GetFeePolicy()
	if defaultFp.AccountFee.FeeConstant.Int64() != 0 {
		t.Error("wrong fee constant")
		return
	}
	if defaultFp.AccountFee.FeeRate.Int64() != 10000 {
		t.Error("wrong fee rate")
		return
	}

	defaultFp.AccountFee.FeeConstant = big.NewInt(5)
	defaultFp.AccountFee.FeeRate = big.NewInt(50000)

	err := model.SaveFeePolicy(defaultFp)
	if err != nil {
		t.Error(err)
		return
	}

	if defaultFp.AccountFee.FeeConstant.Int64() != 5 {
		t.Error("wrong fee constant")
		return
	}
	if defaultFp.AccountFee.FeeRate.Int64() != 50000 {
		t.Error("wrong fee rate")
		return
	}
}
