package models

import (
	"os"
	"testing"

	"path"

	"reflect"

	"fmt"

	"bytes"
	"encoding/gob"

	"encoding/hex"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
}
func setupDb(t *testing.T) (model *ModelDB) {
	dbPath := path.Join(os.TempDir(), "testxxxx.db")
	os.Remove(dbPath)
	os.Remove(dbPath + ".lock")
	model, err := OpenDb(dbPath)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(model.db)
	return
}

func TestToken(t *testing.T) {
	model := setupDb(t)
	defer func() {
		model.CloseDB()
	}()
	var cbtokens []common.Address
	funcb := func(token common.Address) bool {
		cbtokens = append(cbtokens, token)
		return true
	}
	ts, err := model.GetAllTokens()
	if len(ts) > 0 {
		t.Error("should not found")
	}
	if len(ts) != 0 {
		t.Error("should be empty")
	}
	t1 := utils.NewRandomAddress()
	model.RegisterNewTokenCallback(funcb)
	err = model.AddToken(t1, utils.NewRandomAddress())
	if err != nil {
		t.Error(err)
	}
	t2 := utils.NewRandomAddress()
	err = model.AddToken(t2, utils.NewRandomAddress())
	if err != nil {
		t.Error(err)
	}
	if len(cbtokens) != 1 && cbtokens[0] != t1 {
		t.Error("add token error")
	}
}

func TestGob(t *testing.T) {
	s1 := params.RopstenRegistryAddress
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(s1)
	if err != nil {
		t.Error(err)
		return
	}
	encodedData := buf.Bytes()
	fmt.Printf("first\n%s", hex.Dump(encodedData))
	dec := gob.NewDecoder(bytes.NewBuffer(encodedData))
	var sb common.Address
	err = dec.Decode(&sb)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(s1, sb) {
		t.Error("not equal")
	}
	var buf2 bytes.Buffer
	enc2 := gob.NewEncoder(&buf2)
	enc2.Encode(&sb)
	encodedData2 := buf2.Bytes()
	fmt.Printf("second\n%s", hex.Dump(encodedData2))
	if !reflect.DeepEqual(encodedData, encodedData2) {
		t.Error("not equal")
	}

}

func TestWithdraw(t *testing.T) {

	model := setupDb(t)
	defer func() {
		model.CloseDB()
	}()
	channel := utils.NewRandomHash()
	secret := utils.Sha3(channel[:])
	r := model.IsThisLockHasWithdraw(channel, secret)
	if r == true {
		t.Error("should be false")
		return
	}
	model.WithdrawThisLock(channel, secret)
	r = model.IsThisLockHasWithdraw(channel, secret)
	if r == false {
		t.Error("should be true")
		return
	}
	r = model.IsThisLockHasWithdraw(utils.NewRandomHash(), secret)
	if r == true {
		t.Error("shoulde be false")
		return
	}
}

func TestModelDB_IsThisLockRemoved(t *testing.T) {
	model := setupDb(t)
	defer func() {
		model.CloseDB()
	}()
	channel := utils.NewRandomHash()
	secret := utils.Sha3(channel[:])
	sender := utils.NewRandomAddress()
	r := model.IsThisLockRemoved(channel, sender, secret)
	if r {
		t.Error("should be false")
		return
	}
	model.RemoveLock(channel, sender, secret)
	r = model.IsThisLockRemoved(channel, sender, secret)
	if !r {
		t.Error("should be true")
		return
	}
}
