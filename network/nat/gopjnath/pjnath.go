package gopjnath

/*
#cgo pkg-config: libpjproject
#include <pjlib.h>
#include <pjlib-util.h>
#include <pjnath.h>
#include "helper.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"sync"
	"unsafe"

	"github.com/ethereum/go-ethereum/log"
)

type IceTransportOp int

const (
	IceTransportOpStateInit        = IceTransportOp(C.PJ_ICE_STRANS_OP_INIT)
	IceTransportOpStateNegotiation = IceTransportOp(C.PJ_ICE_STRANS_OP_NEGOTIATION)
	IceTransportOpStateKeepAlive   = IceTransportOp(C.PJ_ICE_STRANS_OP_KEEP_ALIVE)
)

type TransportState int

const (
	TransportStateNull      = TransportState(C.PJ_ICE_STRANS_STATE_NULL)
	TransportStateInit      = TransportState(C.PJ_ICE_STRANS_STATE_INIT)
	TransportStateReady     = TransportState(C.PJ_ICE_STRANS_STATE_READY)
	TransportStateSessReady = TransportState(C.PJ_ICE_STRANS_STATE_SESS_READY)
	TransportStateNego      = TransportState(C.PJ_ICE_STRANS_STATE_NEGO)
	TransportStateRunning   = TransportState(C.PJ_ICE_STRANS_STATE_RUNNING)
	TransportStateFailed    = TransportState(C.PJ_ICE_STRANS_STATE_FAILED)
)

type IceSessRole int

const (
	IceSessRoleUnknown     = IceSessRole(C.PJ_ICE_SESS_ROLE_UNKNOWN)
	IceSessRoleControlled  = IceSessRole(C.PJ_ICE_SESS_ROLE_CONTROLLED)
	IceSessRoleControlling = IceSessRole(C.PJ_ICE_SESS_ROLE_CONTROLLING)
)

type QosType int

const (
	QosTypeBestEffort = QosType(C.PJ_QOS_TYPE_BEST_EFFORT)
	QosTypeBackground = QosType(C.PJ_QOS_TYPE_BACKGROUND)
	QosTypeVideo      = QosType(C.PJ_QOS_TYPE_VIDEO)
	QosTypeVoice      = QosType(C.PJ_QOS_TYPE_VOICE)
	QosTypeControl    = QosType(C.PJ_QOS_TYPE_CONTROL)
)

var initErr = errors.New("ice init error")

func casterr(err C.pj_status_t) error {
	if err == C.PJ_SUCCESS {
		return nil
	}
	buf := unsafe.Pointer(C.calloc(80, 1))
	s := C.pj_strerror(err, (*C.char)(buf), 80)
	str := C.pj_strbuf(&s)
	defer C.free(unsafe.Pointer(str))
	return errors.New(C.GoString(str))
}

func toString(s C.pj_str_t) string {
	str := C.pj_strbuf(&s)
	ptr := uintptr(unsafe.Pointer(str))
	var str2 = unsafe.Pointer(C.calloc(C.size_t(C.pj_strlen(&s)+1), 1))
	ptr2 := uintptr(str2)
	for i := 0; i < int(C.pj_strlen(&s)); i++ {
		*(*C.char)(unsafe.Pointer(ptr2)) = *(*C.char)(unsafe.Pointer(ptr))
		ptr++
		ptr2++
	}
	defer C.free(str2)
	return C.GoString((*C.char)(str2))
}

func ptrToString(s unsafe.Pointer) string {
	str := C.pj_strbuf((*C.pj_str_t)(s))
	ptr := uintptr(unsafe.Pointer(str))
	var str2 = unsafe.Pointer(C.calloc(C.size_t(C.pj_strlen((*C.pj_str_t)(s))+1), 1))
	ptr2 := uintptr(str2)
	for i := 0; i < int(C.pj_strlen((*C.pj_str_t)(s))); i++ {
		*(*C.char)(unsafe.Pointer(ptr2)) = *(*C.char)(unsafe.Pointer(ptr))
		ptr++
		ptr2++
	}
	defer C.free(str2)
	return C.GoString((*C.char)(str2))
}

func destroyString(s C.pj_str_t) {
	str := C.pj_strbuf(&s)
	C.free(unsafe.Pointer(str))
}

var once sync.Once

func SetIceLogLevel(lvl log.Lvl) {
	switch lvl {
	case log.LvlCrit:
		fallthrough
	case log.LvlError:
		C.pj_log_set_level(1)
	case log.LvlWarn:
		C.pj_log_set_level(2)
	case log.LvlInfo:
		C.pj_log_set_level(3)
	case log.LvlDebug:
		C.pj_log_set_level(3)
	case log.LvlTrace:
		C.pj_log_set_level(3)
	}
}
func iceInitInternal(sturnServer, turnServer, turnUserName, turnPassword string) error {
	result := C.goice_init(C.CString(sturnServer), C.CString(turnServer), C.CString(turnUserName), C.CString(turnPassword))
	if result != C.PJ_SUCCESS {
		fmt.Println("result :", result)
		return initErr
	}
	return regThisThread()
}
func IceInit(sturnServer, turnServer, turnUserName, turnPassword string) error {
	fmt.Printf("turnserver=%s,stunserver=%s,user=%s,password=%s\n", turnServer, sturnServer, turnUserName, turnPassword)
	//make sure only init once.
	once.Do(func() {
		err := iceInitInternal(sturnServer, turnServer, turnUserName, turnPassword)
		if err != nil {
			log.Crit(fmt.Sprintf("IceInit init err %v", err))
		}
	})
	return nil
}
func regThisThread() error {
	status := C.regThisThread()
	if status != C.PJ_SUCCESS {
		return casterr(status)
	}
	return nil
}
