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
	at.lock.Lock()
	at.lock.Unlock()
}

func TestAlarmTask_StartAndStop(t *testing.T) {
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
