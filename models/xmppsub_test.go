package models

import (
	"testing"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

func TestModelDB_XMPPIsAddrSubed(t *testing.T) {
	db := setupDb(t)
	defer db.db.Close()
	addr := utils.NewRandomAddress()
	if db.XMPPIsAddrSubed(addr) {
		t.Error("should not marked")
		return
	}
	db.XMPPMarkAddrSubed(addr)
	if !db.XMPPIsAddrSubed(addr) {
		t.Error("should marked")
		return
	}
	db.XMPPUnMarkAddr(addr)
	if db.XMPPIsAddrSubed(addr) {
		t.Error("should not marked")
	}
}
