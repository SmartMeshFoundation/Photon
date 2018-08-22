package blockchain

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// --------- tests for alarmtask.go
func TestNewAlarmTask(t *testing.T) {
	assert.NotEmpty(t, at)
	assert.NotEmpty(t, at.client)
	assert.EqualValues(t, -1, at.LastBlockNumber)
	// test chan
	go func() {
		s, ok := <-at.quitChan
		assert.True(t, ok)
		assert.NotEmpty(t, s)
	}()
	at.quitChan <- struct{}{}
	assert.False(t, at.stopped)
	assert.EqualValues(t, time.Second, at.waitTime)
	assert.EqualValues(t, 0, len(at.callback))
	at.lock.Lock()
	at.lock.Unlock()
}

var testCb AlarmCallback = func(blockNumber int64) error {
	return nil
}

func TestAlarmTask_RegisterCallback(t *testing.T) {
	num := len(at.callback)
	for i := 0; i < 10; i++ {
		go func() {
			at.RegisterCallback(&testCb)
		}()
	}
	time.Sleep(1 * time.Second)
	assert.EqualValues(t, num+10, len(at.callback))
}
func TestAlarmTask_RemoveCallback1(t *testing.T) {
	num := len(at.callback)
	at.RegisterCallback(&testCb)
	at.RemoveCallback(&testCb)
	assert.EqualValues(t, num, len(at.callback))
}

func TestAlarmTask_RemoveCallback2(t *testing.T) {
	for i := 0; i < 10; i++ {
		go at.RegisterCallback(&testCb)
	}
	time.Sleep(1 * time.Second)
	num := len(at.callback)
	for i := 0; i < 10; i++ {
		go at.RemoveCallback(&testCb)
	}
	time.Sleep(1 * time.Second)
	assert.EqualValues(t, num-10, len(at.callback))
}

func TestAlarmTask_StartAndStop(t *testing.T) {
	for _, value := range at.callback {
		at.RemoveCallback(value)
	}
	assert.Equal(t, 0, len(at.callback))
	at.RegisterCallback(&testCb)
	oldBlockNo := at.LastBlockNumber
	at.Start()
	// Start for 20s
	fmt.Println("let AlarmTask runs for 20s...")
	for i := 0; i < 20; i++ {
		time.Sleep(time.Second * 1)
	}
	newBlockNum := at.LastBlockNumber
	assert.NotEqual(t, oldBlockNo, newBlockNum)
	assert.False(t, at.stopped)
	at.Stop()
	assert.True(t, at.stopped)
	_, ok := <-at.quitChan
	assert.False(t, ok)
}
