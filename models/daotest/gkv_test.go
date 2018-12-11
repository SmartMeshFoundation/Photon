package daotest

import (
	"os"
	"path"
	"testing"
	"time"

	"fmt"

	"gitee.com/johng/gkvdb/gkvdb"
)

func TestGKV(t *testing.T) {
	if testing.Short() {
		return
	}
	dbPath := path.Join(os.TempDir(), "testxxxx.db")
	err := os.RemoveAll(dbPath)
	err = os.RemoveAll(dbPath + ".lock")
	if err != nil {
		panic(err)
	}
	db, err := gkvdb.New(dbPath)
	if err != nil {
		panic(err)
	}
	key := []byte("key")
	value := []byte("value")
	table, err := db.Table("TestTable")
	if err != nil {
		panic(err)
	}
	err = table.Set(key[:], value[:])
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			tx := db.Begin("TestTable")
			err2 := tx.Set(key[:], value[:])
			err2 = tx.Commit()
			err2 = table.Set(key[:], value[:]) // 不管是使用table.Set还是手动事务并在commit时添加sync=true参数,都无法避免下面的异常发生
			if err2 != nil {
				fmt.Println("errr", err2)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
	count := 0
	for {
		count++
		buf := table.Items(-1)
		//if buf == nil {
		//	fmt.Printf("异常, buf=nil count = %d err=%v\n", count, err)
		//}
		if len(buf) == 0 {
			fmt.Printf("异常, len(buf)=0 count = %d err=%v\n", count, err)
		}
	}
}

func TestGkv_Remove(t *testing.T) {
	dbPath := path.Join(os.TempDir(), "testxxxx.db")
	err := os.RemoveAll(dbPath)
	err = os.RemoveAll(dbPath + ".lock")
	if err != nil {
		panic(err)
	}
	db, err := gkvdb.New(dbPath)
	if err != nil {
		panic(err)
	}
	key := []byte("key1")
	value := []byte("value")
	table, err := db.Table("TestTable")
	if err != nil {
		panic(err)
	}
	err = table.Set(key[:], value[:])
	if err != nil {
		panic(err)
	}
	m := table.Items(-1)
	if len(m) != 1 {
		t.Error("should have one record")
		return
	}
	fmt.Println(m)
	fmt.Println("删除...")
	err = table.Remove(key[:])
	if err != nil {
		t.Error(err)
		return
	}
	//err = db.RemoveFrom(key[:], "TestTable")
	//time.Sleep(time.Second)
	m = table.Items(-1)
	fmt.Println(m)
	if len(m) != 0 {
		t.Error("must be empty")
		time.Sleep(time.Second)
		return
	}
}
