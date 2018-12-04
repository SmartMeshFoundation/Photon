package daotest

import (
	"os"
	"path"
	"testing"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/gob"
)

type UniqueArray struct {
	Name               string
	UniqueIntegerField [20]byte `storm:"id"` // the starting value can be set
}

func testUniqueArray(t *testing.T, dbPath string) {
	db, err := storm.Open(dbPath, storm.Codec(gob.Codec))
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	db.Init(&UniqueArray{})
	p := &UniqueArray{
		Name:               "123",
		UniqueIntegerField: [20]byte{1, 2, 3},
	}
	err = db.Save(p)
	if err != nil {
		t.Error(err)
	}
	var all []*UniqueArray
	db.All(&all)
	if len(all) != 1 {
		t.Error("number error")
	}
	p.Name = "1234"
	err = db.Save(p)
	if err == nil {
		//todo save is create?
		//t.Error("should not save")
	}
}
func TestTwice(t *testing.T) {
	dbpath := path.Join(os.TempDir(), "testxxxx.db")
	os.RemoveAll(dbpath)
	testUniqueArray(t, dbpath)
	//testUniqueArray(t, dbpath)
	//testUniqueArray(t, dbpath)
}
