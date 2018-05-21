package models

import (
	"os"
	"path"
	"testing"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/gob"
)

type Product struct {
	Pk                  int `storm:"id,increment"` // primary key with auto increment
	Name                string
	IntegerField        uint64 `storm:"increment"`
	IndexedIntegerField uint32 `storm:"index,increment"`
	UniqueIntegerField  int16  `storm:"unique,increment=100"` // the starting value can be set
}

func TestProduct(t *testing.T) {
	model := setupDb(t)
	defer model.CloseDB()
	model.db.Init(&Product{})
	p := &Product{
		Name:               "123",
		UniqueIntegerField: 123,
	}
	err := model.db.Save(p)
	if err != nil {
		t.Error(err)
	}
	var all []*Product
	model.db.All(&all)
	if len(all) != 1 {
		t.Error("number error")
	}
	p.Pk++
	err = model.db.Save(p)
	if err == nil {
		t.Error("should not save")
	}
}

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
	os.Remove(dbpath)
	testUniqueArray(t, dbpath)
	//testUniqueArray(t, dbpath)
	//testUniqueArray(t, dbpath)
}
