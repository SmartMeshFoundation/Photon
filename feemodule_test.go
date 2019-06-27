package photon

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SmartMeshFoundation/Photon/channel"

	"math/big"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/codefortest"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/models/stormdb"
	"github.com/SmartMeshFoundation/Photon/pfsproxy"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

func TestFeeModule_Local(t *testing.T) {
	db, err := newTestStormDb()
	if err != nil {
		t.Error(err.Error())
		return
	}
	fm := NewFeeModule(db, nil)
	fakeAddress := utils.NewRandomAddress()
	var amount, fee *big.Int

	// default fee
	amount = big.NewInt(10000)
	fee = fm.GetNodeChargeFee(fakeAddress, fakeAddress, amount)
	fmt.Println(fee)
	if fee == nil || fee.Int64() != 1 {
		t.Error("fee wrong")
		return
	}

	fm.feePolicy.AccountFee.FeeConstant = big.NewInt(5)
	fm.SetFeePolicy(fm.feePolicy)
	fee = fm.GetNodeChargeFee(fakeAddress, fakeAddress, amount)
	fmt.Println(fee)
	if fee == nil || fee.Int64() != 6 {
		t.Error("fee wrong")
		return
	}

	fm.feePolicy.TokenFeeMap[fakeAddress] = &models.FeeSetting{
		FeeConstant: big.NewInt(10),
		FeePercent:  100,
	}
	fm.SetFeePolicy(fm.feePolicy)
	fee = fm.GetNodeChargeFee(fakeAddress, fakeAddress, amount)
	fmt.Println(fee)
	if fee == nil || fee.Int64() != 110 {
		t.Error("fee wrong")
		return
	}
}

func TestFeeModule_WithPFS(t *testing.T) {
	if testing.Short() {
		return
	}
	// dao
	db, err := newTestStormDb()
	if err != nil {
		t.Error(err.Error())
		return
	}
	// pfs proxy
	alice, err := codefortest.GetAccountsByAddress(common.HexToAddress("0x10b256b3C83904D524210958FA4E7F9cAFFB76c6"))
	if err != nil {
		t.Error(err.Error())
		return
	}
	pfsProxy := pfsproxy.NewPfsProxy("http://192.168.124.9:7000", alice.PrivateKey)
	// fee module
	fm := NewFeeModule(db, pfsProxy)
	fakeAddress := utils.NewRandomAddress()
	var amount, fee *big.Int

	// default fee
	amount = big.NewInt(10000)
	fee = fm.GetNodeChargeFee(fakeAddress, fakeAddress, amount)
	fmt.Println(fee)
	if fee == nil || fee.Int64() != 1 {
		t.Error("fee wrong")
		return
	}

	fm.feePolicy.AccountFee.FeeConstant = big.NewInt(5)
	fm.SetFeePolicy(fm.feePolicy)
	fee = fm.GetNodeChargeFee(fakeAddress, fakeAddress, amount)
	fmt.Println(fee)
	if fee == nil || fee.Int64() != 6 {
		t.Error("fee wrong")
		return
	}

	fm.feePolicy.TokenFeeMap[fakeAddress] = &models.FeeSetting{
		FeeConstant: big.NewInt(10),
		FeePercent:  100,
	}
	fm.SetFeePolicy(fm.feePolicy)
	fee = fm.GetNodeChargeFee(fakeAddress, fakeAddress, amount)
	fmt.Println(fee)
	if fee == nil || fee.Int64() != 110 {
		t.Error("fee wrong")
		return
	}
}

// newTestStormDb :
func newTestStormDb() (dao models.Dao, err error) {
	dbPath := path.Join(os.TempDir(), "testxxxx.dao")
	err = os.Remove(dbPath)
	err = os.Remove(dbPath + ".lock")
	return stormdb.OpenDb(dbPath)
}

func TestUninitMapMap(t *testing.T) {
	var ts map[common.Address]map[common.Hash][]*channel.Channel
	ts = make(map[common.Address]map[common.Hash][]*channel.Channel)
	channelsRegistered := ts[utils.NewRandomAddress()][utils.NewRandomHash()]
	assert.EqualValues(t, channelsRegistered == nil, true)
}

func TestSlice(t *testing.T) {
	s := []int{1}
	s = append(s[:0], s[1:]...)
	fmt.Println(s)
}
