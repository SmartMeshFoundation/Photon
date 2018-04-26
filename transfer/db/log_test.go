package db

import (
	"bytes"
	"encoding/gob"
	"os"
	"testing"

	"path"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/davecgh/go-spew/spew"
)

/*
test gob encoding and decoding
*/

type struct1 struct {
	Name string
	A    int
	B    int
}
type struct2 struct {
	Name string
	C    int
	D    int
}
type structbase struct {
	Name          string
	Data2WriteXXX interface{}
}

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
}
func TestObjTypeName(t *testing.T) {
	if objTypeName(struct1{}) != "struct1" {
		t.Error("typename error")
	}
}
func TestGob(t *testing.T) {
	s1 := struct1{"", 3, 5}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(s1)
	if err != nil {
		t.Error(err)
		return
	}
	encodedData := buf.Bytes()
	dec := gob.NewDecoder(bytes.NewBuffer(encodedData))
	var sb structbase
	err = dec.Decode(&sb)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Name:", sb.Name)
	var r interface{}
	if sb.Name == "struct1" {
		r = new(struct1)
	}
	dec2 := gob.NewDecoder(bytes.NewBuffer(encodedData))
	err = dec2.Decode(r)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("r=%#v", r)

}

/*
Use gob to serialize storage and de serialization of objects.
*/
func TestGobWithWrapper(t *testing.T) {
	gob.Register(&struct1{})
	gob.Register(&struct2{})
	s1 := struct1{"", 3, 5}
	sb := structbase{Name: "base", Data2WriteXXX: &s1}
	var buf = new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	err := enc.Encode(sb)
	if err != nil {
		t.Error(err)
		return
	}
	encodedData := buf.Bytes()
	dec := gob.NewDecoder(bytes.NewBuffer(encodedData))
	var sb2 structbase
	err = dec.Decode(&sb2)
	t.Logf("sb before=%#v", sb2)
	if err != nil {
		t.Error(err)
		return
	}
	sb = structbase{"base2", &struct2{"s2", 5, 8}}
	buf = new(bytes.Buffer)
	enc = gob.NewEncoder(buf)
	err = enc.Encode(sb)
	if err != nil {
		t.Error(err)
		return
	}
	dec = gob.NewDecoder(bytes.NewBuffer(buf.Bytes()))
	var sb3 structbase
	err = dec.Decode(&sb3)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("sb3=%#v", sb3)
	if _, ok := sb3.Data2WriteXXX.(*struct2); !ok {
		t.Error("type convert error")
	}
}

func TestNewStateChangeLog(t *testing.T) {
	dbPath := path.Join(os.TempDir(), "test.db")
	os.Remove(dbPath)
	db := NewStateChangeLogBoltBackend(dbPath)
	t.Log(db)
	data := []byte{1, 2, 3}
	id, err := db.WriteStateChange(data)
	if err != nil {
		t.Error(err)
	}
	//if id != 1 {
	//	t.Error("id not equal 1, ", id)
	//}
	data2, err := db.GetStateChangeById(id)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(data, data2) {
		spew.Dump(data2)
		t.Error("data not equal")
	}
	s := &snapshotToWrite{
		StateChangeId: 1,
		State:         data,
	}
	_, err = db.WriteStateSnapshot(s)
	if err != nil {
		t.Error(err)
	}
	s2, err := db.GetStateSnapshot()
	if err != nil {
		t.Error(err)
	}
	if s2.StateChangeId != s.StateChangeId {
		t.Error("id not equal")
	}
	data3 := s2.State.([]byte)
	if !bytes.Equal(data3, data) {
		t.Error("data3 not equal")
	}
	number := utils.RandSrc.Int63()
	events := []*InternalEvent{{StateChangeId: id,
		BlockNumber: number,
		EventObject: data}}
	err = db.WriteStateEvents(id, events)
	if err != nil {
		t.Error(err)
	}
	events2, err := db.GetEventsInRange(number, number+1)
	if err != nil {
		t.Error(err)
	}
	if len(events2) != 1 {
		t.Error("events length error")
	}
	t.Log("events2=%#v", events2[0])
}
