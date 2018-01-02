package transfer

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"strings"

	"fmt"

	"encoding/hex"

	"database/sql"

	"sync"

	"github.com/SmartMeshFoundation/raiden-network/utils"
	_ "github.com/cznic/ql/driver"
	"github.com/ethereum/go-ethereum/log"
)

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
	WriteStateSnapshot(stateChangeId string, data []byte) (int64, error)
}

type StateChangeLogQlBackend struct {
	db   *sql.DB
	lock sync.RWMutex
}

func createQlDb(dbPath string) *StateChangeLogQlBackend {
	db, err := sql.Open("ql", dbPath)
	if err != nil {
		log.Crit(fmt.Sprintf("cannot create or open db:%s,makesure you have write permission err:%v", dbPath, err))
	}
	sqlStmt := `
CREATE TABLE  if not EXISTS state_changes (
    id   string,
    data blob
);
CREATE TABLE  if not EXISTS   state_snapshot (
    identifier     int,
    statechange_id string,
    data           blob
);
CREATE TABLE  if not EXISTS  state_events (
	identifier string,
    source_statechange_id string,
    block_number          int,
    data                  blob
);
`
	tx, _ := db.Begin()
	_, err = tx.Exec(sqlStmt)
	if err != nil {
		log.Error(fmt.Sprintf("create db table error:%v", err))
	}
	tx.Commit()
	sb := &StateChangeLogQlBackend{db: db}
	return sb
}
func NewStateChangeLogQlBackend(dbPath string) *StateChangeLogQlBackend {
	return createQlDb(dbPath)
}

func (sb *StateChangeLogQlBackend) WriteStateChange(data []byte) (string, error) {
	sb.lock.Lock()
	defer sb.lock.Unlock()
	//插入数据
	tx, err := sb.db.Begin()
	if err != nil {
		return "", err
	}
	key := utils.RandomString(32)
	utils.NewRandomAddress()
	_, err = tx.Exec("INSERT INTO state_changes(id,data) VALUES($1,$2)", key, data)
	if err != nil {
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		return "", err
	}
	return key, err
}

func (sb *StateChangeLogQlBackend) WriteStateSnapshot(stateChangeId string, data []byte) (int64, error) {
	sb.lock.Lock()
	defer sb.lock.Unlock()
	tx, _ := sb.db.Begin()
	defer tx.Commit()
	rows, err := tx.Query("select * from state_changes where id=1")
	if err != nil {
		return 0, err
	}
	if rows.Next() {
		_, err = tx.Exec("update  state_snapshot set  VALUES($1,$2,$3)", 1, stateChangeId, data)
	} else {
		_, err = tx.Exec("INSERT INTO state_snapshot(identifier, statechange_id, data) VALUES($1,$2,$3)", 1, stateChangeId, data)
	}
	rows.Close()
	if err != nil {
		return 0, err
	}
	return 1, err
}

type eventWriterStruct struct {
	Identifier    string
	StateChangeId string
	BlockNumber   int64
	EventData     []byte
}
type InternalEvent struct {
	Identifier    string
	StateChangeId string
	BlockNumber   int64
	EventObject   Event
}

/*
 """Do an 'execute_many' write of state events. `events_data` should be a
        list of tuples of the form:
        (None, source_statechange_id, block_number, serialized_event_data)
        """
*/
func (sb *StateChangeLogQlBackend) WriteStateEvents(stateChangeId string, events []*eventWriterStruct) error {
	sb.lock.Lock()
	defer sb.lock.Unlock()
	tx, err := sb.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO state_events(identifier,source_statechange_id, block_number, data) VALUES($1,$2,$3,$4)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, ev := range events {
		_, err = stmt.Exec(utils.RandomString(32), ev.StateChangeId, ev.BlockNumber, ev.EventData)
		if err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}

//Return the last state snapshot
func (sb *StateChangeLogQlBackend) GetStateSnapshot() (string, []byte, error) {
	rows, err := sb.db.Query("SELECT statechange_id,data from state_snapshot")
	if err != nil {
		return "", nil, err
	}
	defer rows.Close()
	//only one row
	for rows.Next() {
		var id string
		var data []byte
		err = rows.Scan(&id, &data)
		if err != nil {
			return "", nil, err
		}
		return id, data, nil
	}
	return "", nil, err
}

func (sb *StateChangeLogQlBackend) GetStateChangeById(id string) (data []byte, err error) {
	rows, err := sb.db.Query("SELECT data from state_changes where id=$1", id)
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&data)
		if err != nil {
			break
		}
		break
	}
	rows.Close()
	return
}

func (sb *StateChangeLogQlBackend) GetEventsInRange(fromblock, toblock int64) (events []*eventWriterStruct, err error) {
	var stmt *sql.Stmt
	tx, _ := sb.db.Begin()
	defer tx.Commit()
	if toblock > 0 {
		stmt, err = tx.Prepare("SELECT * from  state_events WHERE block_number BETWEEN $1 AND $2 ")
	} else {
		stmt, err = tx.Prepare("SELECT * from  state_events WHERE block_number >= ?")
	}
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(fromblock, toblock)
	if err != nil {
		return
	}
	for rows.Next() {
		ev := &eventWriterStruct{}
		err = rows.Scan(&ev.Identifier, &ev.StateChangeId, &ev.BlockNumber, &ev.EventData)
		if err != nil {
			return
		}
		events = append(events, ev)
	}
	return
}

func (sb *StateChangeLogQlBackend) Close() {
	sb.db.Close()
}

type StateChangeLog struct {
	Storage    *StateChangeLogQlBackend
	Serializer StateChangeLogSerializer
}

func NewStateChangeLog(dbPath string) *StateChangeLog {
	return &StateChangeLog{
		Storage:    NewStateChangeLogQlBackend(dbPath),
		Serializer: new(GobSerialize),
	}
}

//Log a state change and return its identifier
func (this *StateChangeLog) Log(stateChange StateChange) (string, error) {
	/*
			 # TODO: Issue 587
		        # Implement a queue of state changes for batch writting
	*/
	data := this.Serializer.Serialize(stateChange)
	return this.Storage.WriteStateChange(data)
}

// Log the events that were generated by `state_change_id` into the write ahead Log
func (this *StateChangeLog) LogEvents(stateChangeId string, events []Event, currentBlockNumber int64) error {
	var eventsWriter []*eventWriterStruct
	for _, e := range events {
		eventsWriter = append(eventsWriter, &eventWriterStruct{
			StateChangeId: stateChangeId,
			BlockNumber:   currentBlockNumber,
			EventData:     this.Serializer.Serialize(e),
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
	es, err := this.Storage.GetEventsInRange(fromBlock, toBlock)
	if err != nil {
		return
	}
	for _, e := range es {
		events = append(events, &InternalEvent{
			Identifier:    e.Identifier,
			StateChangeId: e.StateChangeId,
			BlockNumber:   e.BlockNumber,
			EventObject:   Event(this.Serializer.DeSerialize(e.EventData)),
		})
	}
	return
}

func (this *StateChangeLog) GetStateChangeById(id string) (st StateChange, err error) {
	data, err := this.Storage.GetStateChangeById(id)
	if err != nil {
		return
	}
	st = StateChange(this.Serializer.DeSerialize(data))
	return
}
func (this *StateChangeLog) Snapshort(stateChangeId string, state State) (int64, error) {
	data := this.Serializer.Serialize(state)
	return this.Storage.WriteStateSnapshot(stateChangeId, data)
}
