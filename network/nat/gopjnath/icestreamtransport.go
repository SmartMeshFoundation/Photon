package gopjnath

/*
#cgo pkg-config: libpjproject
#include "helper.h"
*/
import "C"

import (
	//"log"
	"unsafe"

	"fmt"

	"sync"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/go-errors/errors"
)

type RxDataCb func(uint, []byte, SockAddr)
type IceCompleteCb func(IceTransportOp, error)
type IceStreamTransport struct {
	i             *C.IceInstance
	cb            *C.pj_ice_strans_cb
	OnRxData      RxDataCb
	OnIceComplete IceCompleteCb
	key           string //key for index
	name          string
}

var streamMap map[string]*IceStreamTransport
var streamlock sync.Mutex

func init() {
	streamMap = make(map[string]*IceStreamTransport)
}

// pj_status_t pj_ice_strans_create (const char *name, const pj_ice_strans_cfg *cfg, unsigned comp_cnt, void *user_data, const pj_ice_strans_cb *cb, pj_ice_strans **p_ice_st)
/*
dataCallback,iceCallback must not block .
*/
func NewIceStreamTransport(name string, dataCallback RxDataCb, iceCallback IceCompleteCb) (*IceStreamTransport, error) {
	n := C.CString(name)
	var err error
	defer C.free(unsafe.Pointer(n))
	regThisThread()
	if dataCallback == nil {
		dataCallback = defaultDataCallback
	}
	if iceCallback == nil {
		iceCallback = defaultIceCallback
	}
	stream := &IceStreamTransport{
		OnRxData:      dataCallback,
		OnIceComplete: iceCallback,
		name:          fmt.Sprintf("%s-%d", name, utils.NewRandomInt(10)),
	}
	stream.i = C.gopjnath_create_iceinstance(n) //must remember free.
	key := utils.RandomString(10)
	ckey := C.CString(key) //remember free,free this instance when destroy
	streamlock.Lock()
	streamMap[key] = stream
	streamlock.Unlock()
	stream.key = key
	C.gopjnath_set_user_data(stream.i, unsafe.Pointer(ckey))
	status := C.gopjnath_create_instance(stream.i, n)
	if status != C.PJ_SUCCESS {
		err = errors.New(fmt.Sprintf("%s NewIceStreamTransport error %d", name, casterr(status)))
		goto onerror
	}
	return stream, nil
onerror:
	stream.Destroy()
	return nil, err
}

// pj_ice_strans_state pj_ice_strans_get_state (pj_ice_strans *ice_st)
func (i *IceStreamTransport) State() TransportState {
	return TransportState(C.pj_ice_strans_get_state(C.gopjnath_get_icest(i.i)))
}

// const char * pj_ice_strans_state_name (pj_ice_strans_state state)
func TransportStateName(t TransportState) string {
	str := C.pj_ice_strans_state_name(C.pj_ice_strans_state(t))
	return C.GoString(str)
}

// void * pj_ice_strans_get_user_data (pj_ice_strans *ice_st)
// not implementing right now :)

// pj_grp_lock_t * pj_ice_strans_get_grp_lock (pj_ice_strans *ice_st)

// pj_status_t pj_ice_strans_init_ice (pj_ice_strans *ice_st, pj_ice_sess_role role, const pj_str_t *local_ufrag, const pj_str_t *local_passwd)
func (i *IceStreamTransport) InitIceSession(r IceSessRole) error {
	err := regThisThread()
	if err != nil {
		return err
	}
	status := C.gopjnath_init_session(i.i, C.pj_ice_sess_role(r))
	return casterr(status)
}

// pj_bool_t pj_ice_strans_has_sess (pj_ice_strans *ice_st)
func (i *IceStreamTransport) HasSess() bool {
	return int(C.pj_ice_strans_has_sess(C.gopjnath_get_icest(i.i))) != 0
}

// pj_bool_t pj_ice_strans_sess_is_running (pj_ice_strans *ice_st)
func (i *IceStreamTransport) SessIsRunning() bool {
	return int(C.pj_ice_strans_sess_is_running(C.gopjnath_get_icest(i.i))) != 0
}

// pj_bool_t pj_ice_strans_sess_is_complete (pj_ice_strans *ice_st)
func (i *IceStreamTransport) SessIsComplete() bool {
	return int(C.pj_ice_strans_sess_is_complete(C.gopjnath_get_icest(i.i))) != 0
}

// unsigned pj_ice_strans_get_running_comp_cnt (pj_ice_strans *ice_st)
func (i *IceStreamTransport) RunningCompCount() uint {
	return uint(C.pj_ice_strans_get_running_comp_cnt(C.gopjnath_get_icest(i.i)))
}

// pj_status_t pj_ice_strans_get_ufrag_pwd (pj_ice_strans *ice_st, pj_str_t *loc_ufrag, pj_str_t *loc_pwd, pj_str_t *rem_ufrag, pj_str_t *rem_pwd)
func (i *IceStreamTransport) UfragPwd() (string, string, string, string, error) {
	var locUfrag, locPwd, remUfrag, remPwd C.pj_str_t
	status := C.pj_ice_strans_get_ufrag_pwd(C.gopjnath_get_icest(i.i), &locUfrag, &locPwd, &remUfrag, &remPwd)
	var lu, lp, ru, rp string

	lu = toString(locUfrag)
	lp = toString(locPwd)
	ru = toString(remUfrag)
	rp = toString(remPwd)

	if status != C.PJ_SUCCESS {
		return lu, lp, ru, rp, casterr(status)
	}
	return lu, lp, ru, rp, nil
}

// unsigned pj_ice_strans_get_cands_count (pj_ice_strans *ice_st, unsigned comp_id)
func (i *IceStreamTransport) CandsCount(compId uint) uint {
	return uint(C.pj_ice_strans_get_cands_count(C.gopjnath_get_icest(i.i), C.uint(compId)))
}

// pj_ice_sess_role pj_ice_strans_get_role (pj_ice_strans *ice_st)
func (i *IceStreamTransport) Role() IceSessRole {
	return IceSessRole(C.pj_ice_strans_get_role(C.gopjnath_get_icest(i.i)))
}

// pj_status_t pj_ice_strans_change_role (pj_ice_strans *ice_st, pj_ice_sess_role new_role)
func (i *IceStreamTransport) ChangeRole(r IceSessRole) error {
	status := C.pj_ice_strans_change_role(C.gopjnath_get_icest(i.i), C.pj_ice_sess_role(C.int(r)))
	if status != C.PJ_SUCCESS {
		return casterr(status)
	}
	return nil
}

// pj_status_t pj_ice_strans_start_ice (pj_ice_strans *ice_st, const pj_str_t *rem_ufrag, const pj_str_t *rem_passwd, unsigned rcand_cnt, const pj_ice_sess_cand rcand[])
func (i *IceStreamTransport) StartIce(sdp string) error {
	regThisThread()
	err := i.SetRemoteSdp(sdp)
	if err != nil {
		return err
	}
	status := C.gopjnath_start_nego(i.i)
	if status != C.PJ_SUCCESS {
		return casterr(status)
	}
	return nil
}

//// const pj_ice_sess_check * pj_ice_strans_get_valid_pair (const pj_ice_strans *ice_st, unsigned comp_id)
//func (i *IceStreamTransport) ValidPair(compId uint) IceSessCheck {
//	return IceSessCheck{C.pj_ice_strans_get_valid_pair(i.i, C.uint(compId))}
//}

/*
we may use IceInstance again after initIceSession and StartIce
after stop ice, i cannot send, but can receive data from partner.
*/
func (i *IceStreamTransport) StopIce() error {
	err := regThisThread()
	if err != nil {
		return err
	}
	status := C.gopjnath_stop_session(i.i)
	if status != C.PJ_SUCCESS {
		return casterr(status)
	}
	return nil
}

/*
destroy ice instance ,may not reuse again
*/
func (i *IceStreamTransport) Destroy() {
	regThisThread()
	C.gopjnath_destroy_instance(i.i)
}

// pj_status_t pj_ice_strans_sendto (pj_ice_strans *ice_st, unsigned comp_id, const void *data, pj_size_t data_len, const pj_sockaddr_t *dst_addr, int dst_addr_len)
func (i *IceStreamTransport) Send(data []byte) error {
	cdata := C.CBytes(data)
	defer C.free(cdata)
	err := regThisThread()
	if err != nil {
		log.Error(fmt.Sprintf("%s send regthisthread fail %s", i.name, err))
	}
	log.Trace(fmt.Sprintf("icestreamtransport %s send data %v", i.name, data[0]))
	status := C.gopjnath_send_data(i.i, cdata, C.pj_size_t(len(data)))
	if status != C.PJ_SUCCESS {
		return casterr(status)
	}
	return nil
}
func (i *IceStreamTransport) ShowIceInfo() {
	C.gopjnath_show_ice(i.i)
}
func (i *IceStreamTransport) GetLocalSdp() (string, error) {
	regThisThread()
	d := C.malloc(C.size_t(1024))
	//defer C.free(d)
	r := 1024
	var l C.int
	for {
		l = C.gopjnath_encode_session(i.i, (*C.char)(d), C.int(r))
		if l < 0 && int(l) == int(C.PJ_ETOOSMALL) {
			r *= 2
			d = C.realloc(d, C.size_t(r))
			continue
		}
		if l < 0 {
			return "", casterr((C.pj_status_t)(-l))
		}
		break
	}
	return C.GoStringN((*C.char)(d), l), nil
}
func (i *IceStreamTransport) SetRemoteSdp(sdp string) error {
	csdp := C.CString(sdp)
	defer C.free(unsafe.Pointer(csdp))
	err := regThisThread()
	if err != nil {
		return err
	}
	status := C.gopjnath_input_remote(i.i, csdp)
	if status != C.PJ_SUCCESS {
		return casterr(status)
	}
	return nil
}

//export go_ice_callback
func go_ice_callback(i *C.pj_ice_strans, o C.pj_ice_strans_op, s C.pj_status_t) {
	key := C.GoString((*C.char)(C.pj_ice_strans_get_user_data(i)))
	strans := streamMap[key]
	if strans == nil {
		panic(fmt.Sprintf("key error=%s", key))
	}
	strans.OnIceComplete(IceTransportOp(o), casterr(s))
}

//export go_data_callback
func go_data_callback(i *C.pj_ice_strans, comp_id C.unsigned, pkt unsafe.Pointer, size C.pj_size_t, src_addr *C.pj_sockaddr_t, src_addr_len C.unsigned) {
	data := C.GoBytes(pkt, C.int(size))
	//TODO make ral SockAddr
	s := SockAddr{}
	key := C.GoString((*C.char)(C.pj_ice_strans_get_user_data(i)))
	strans := streamMap[key]
	if strans == nil {
		panic(fmt.Sprintf("key error=%s", key)) //碰到一次panic,为什么呢,
	}
	log.Trace(fmt.Sprintf("icedata callback %s received data %v", strans.name, data[0]))
	strans.OnRxData(uint(comp_id), data, s)
}

func defaultDataCallback(compId uint, data []byte, addr SockAddr) {
}

func defaultIceCallback(op IceTransportOp, err error) {
}
