package db

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"strings"

	"fmt"

	"encoding/hex"

	"sync"

	"time"

	"encoding/binary"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	bolt "github.com/coreos/bbolt"
	"github.com/ethereum/go-ethereum/common"
)

var bucketEvents = []byte("events")
var bucketEventsBlock = []byte("eventsBlock")
var bucketStateChange = []byte("statechange")
var bucketSnapshot = []byte("snapshot")
var bucketMeta = []byte("meta")

const dbVersion = 1

type StateChangeLogSerializer interface {
	Serialize(object interface{}) []byte
	DeSerialize(data []byte) interface{}
}

/*
A simple transaction serializer using pickle
todo buf,enc,dec is reusable?
*/
type GobSerialize struct {
}
type gobHelper struct {
	Name string
	Data interface{}
}

func objTypeName(object interface{}) string {
	t := reflect.ValueOf(object)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	names := strings.Split(t.Type().String(), ".")
	return names[len(names)-1]
}

type snapshotToWrite struct {
	StateChangeId int
	State         interface{}
}

/*
all the type must gob.Register
*/
func (g *GobSerialize) Serialize(object interface{}) []byte {
	gh := &gobHelper{Name: objTypeName(object), Data: object}
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(gh)
	if err != nil {
		//must be a bug
		log.Crit(fmt.Sprintf("cannot gob encode this object %#v", object))
	}
	return buf.Bytes()
}

func (g *GobSerialize) DeSerialize(data []byte) interface{} {
	gh := &gobHelper{}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(gh)
	if err != nil {
		log.Crit(fmt.Sprintf("cannot gob decode data:\n%s", hex.Dump(data)))
	}
	return gh.Data
}

/*
An abstract class defining the storage backend for the transaction log.
        Allows for pluggable storage backends.
*/
type StateChangeLogStorageBackender interface {
	WriteStateChange(data []byte) (int64, error) //return last insert id
	WriteStateSnapshot(stateChangeId int, data []byte) (int64, error)
}

type StateChangeLogBoltBackend struct {
	db         *bolt.DB
	lock       sync.RWMutex
	Serializer *GobSerialize
}

func createDb(dbPath string) *StateChangeLogBoltBackend {
	needUpdate := !common.FileExist(dbPath)
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Crit(fmt.Sprintf("cannot create or open db:%s,makesure you have write permission err:%v", dbPath, err))
	}
	if needUpdate {
		err = db.Update(func(tx *bolt.Tx) error {
			if _, err := tx.CreateBucketIfNotExists(bucketEvents); err != nil {
				return err
			}

			if _, err := tx.CreateBucketIfNotExists(bucketSnapshot); err != nil {
				return err
			}

			if _, err := tx.CreateBucketIfNotExists(bucketStateChange); err != nil {
				return err
			}
			if _, err := tx.CreateBucketIfNotExists(bucketEventsBlock); err != nil {
				return err
			}
			if _, err := tx.CreateBucketIfNotExists(bucketMeta); err != nil {
				return err
			}
			tx.Bucket(bucketMeta).Put([]byte("version"), itob(dbVersion))
			return nil
		})
		if err != nil {
			log.Crit(fmt.Sprintf("unable to create db "))
			return nil
		}
	}
	sb := &StateChangeLogBoltBackend{db: db, Serializer: new(GobSerialize)}
	return sb
}
func NewStateChangeLogBoltBackend(dbPath string) *StateChangeLogBoltBackend {
	return createDb(dbPath)
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
func i64tob(v int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
func u64tob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func (sb *StateChangeLogBoltBackend) WriteStateChange(data []byte) (id int, err error) {
	err = sb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketStateChange)
		id2, err := b.NextSequence()
		if err != nil {
			return err
		}
		id = int(id2)
		return b.Put(itob(id), data)
	})
	return
}

func (sb *StateChangeLogBoltBackend) WriteStateSnapshot(s *snapshotToWrite) (id int, err error) {
	id = 1
	sb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketSnapshot)
		return b.Put(itob(1), sb.Serializer.Serialize(s))
	})
	return
}

type InternalEvent struct {
	Identifier    int
	StateChangeId int
	BlockNumber   int64
	EventObject   transfer.Event
}

func blockKey(blockNumber int64, identifier int) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, blockNumber)
	binary.Write(buf, binary.BigEndian, identifier)
	return buf.Bytes()
}

/*
 """Do an 'execute_many' write of state events. `events_data` should be a
        list of tuples of the form:
        (None, source_statechange_id, block_number, serialized_event_data)
        """
*/
func (sb *StateChangeLogBoltBackend) WriteStateEvents(stateChangeId int, events []*InternalEvent) error {
	err := sb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketEvents)
		b2 := tx.Bucket(bucketEventsBlock)
		for _, e := range events {
			id, err := b.NextSequence()
			if err != nil {
				return err
			}
			e.Identifier = int(id)
			e.StateChangeId = stateChangeId
			err = b.Put(u64tob(id), sb.Serializer.Serialize(e))
			if err != nil {
				return err
			}
			err = b2.Put(blockKey(e.BlockNumber, e.Identifier), sb.Serializer.Serialize(e))
			if err != nil {
				return err
			}
			return nil
		}
		return nil
	})
	return err
}

//Return the last state snapshot
func (sb *StateChangeLogBoltBackend) GetStateSnapshot() (s *snapshotToWrite, err error) {
	err = sb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketSnapshot)
		d := b.Get(itob(1))
		if d == nil {
			return nil
		}
		s = new(GobSerialize).DeSerialize(d).(*snapshotToWrite)
		return nil
	})
	return
}

func (sb *StateChangeLogBoltBackend) GetStateChangeById(id int) (data []byte, err error) {
	err = sb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketStateChange)
		data = b.Get(itob(id))
		return nil
	})
	return
}

func (sb *StateChangeLogBoltBackend) GetEventsInRange(fromblock, toblock int64) (events []*InternalEvent, err error) {
	if fromblock < 0 {
		fromblock = 0
	}
	err = sb.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketEventsBlock).Cursor()
		if toblock <= 0 {
			k, v := c.Last()
			if k == nil {
				return nil
			}
			e := sb.Serializer.DeSerialize(v).(*InternalEvent)
			toblock = e.BlockNumber
		}
		min := i64tob(fromblock)
		max := i64tob(toblock)
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			e := sb.Serializer.DeSerialize(v).(*InternalEvent)
			if e.BlockNumber >= fromblock && e.BlockNumber <= toblock {
				events = append(events, e)
			}
		}
		return nil
	})
	return
}

func (sb *StateChangeLogBoltBackend) Close() {
	sb.db.Close()
}

type StateChangeLog struct {
	Storage    *StateChangeLogBoltBackend
	Serializer StateChangeLogSerializer
}

func NewStateChangeLog(dbPath string) *StateChangeLog {
	return &StateChangeLog{
		Storage:    NewStateChangeLogBoltBackend(dbPath),
		Serializer: new(GobSerialize),
	}
}

//Log a state change and return its identifier
func (this *StateChangeLog) Log(stateChange transfer.StateChange) (int, error) {
	/*
			  TODO: Issue 587
		         Implement a queue of state changes for batch writting
	*/
	data := this.Serializer.Serialize(stateChange)
	return this.Storage.WriteStateChange(data)
}

// Log the events that were generated by `state_change_id` into the write ahead Log
func (this *StateChangeLog) LogEvents(stateChangeId int, events []transfer.Event, currentBlockNumber int64) error {
	var eventsWriter []*InternalEvent
	for _, e := range events {
		eventsWriter = append(eventsWriter, &InternalEvent{
			StateChangeId: stateChangeId,
			BlockNumber:   currentBlockNumber,
			EventObject:   e,
		})
	}
	return this.Storage.WriteStateEvents(stateChangeId, eventsWriter)
}

/*
Get the raiden events in the period (inclusive) ranging from
        `from_block` to `to_block`.

        This function returns a list of tuples of the form:
        (identifier, generated_statechange_id, block_number, event_object)
*/
func (this *StateChangeLog) GetEventsInBlockRange(fromBlock, toBlock int64) (events []*InternalEvent, err error) {
	events, err = this.Storage.GetEventsInRange(fromBlock, toBlock)
	return
}

func (this *StateChangeLog) GetStateChangeById(id int) (st transfer.StateChange, err error) {
	data, err := this.Storage.GetStateChangeById(id)
	if err != nil {
		return
	}
	st = transfer.StateChange(this.Serializer.DeSerialize(data))
	return
}
func (this *StateChangeLog) Snapshot(stateChangeId int, state interface{}) (int, error) {
	s := &snapshotToWrite{
		StateChangeId: stateChangeId,
		State:         state,
	}
	return this.Storage.WriteStateSnapshot(s)
}

func (this *StateChangeLog) LoadSnapshot() (interface{}, error) {
	s, err := this.Storage.GetStateSnapshot()
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, nil
	}
	return s.State, nil
}
func init() {
	gob.Register(&InternalEvent{})
	gob.Register(&snapshotToWrite{})
}
