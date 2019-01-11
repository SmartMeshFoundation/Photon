package utils

import (
	"fmt"
	"sync"

	"github.com/SmartMeshFoundation/Photon/log"
)

//DebugLock a tool for dead lock detect
type DebugLock struct {
	sync.Mutex
	name string
}

//Lock with name for dead lock detect
func (m *DebugLock) Lock() {
	s := RandomString(10)
	log.Info(fmt.Sprintf("try DebugLock lock %s", s))
	m.Mutex.Lock()
	m.name = s
	log.Info(fmt.Sprintf("%s DebugLock locked", s))
}

//Unlock with name for dead lock detect
func (m *DebugLock) Unlock() {
	log.Info(fmt.Sprintf("%s DebugLock unlocked", m.name))
	m.Mutex.Unlock()
}
