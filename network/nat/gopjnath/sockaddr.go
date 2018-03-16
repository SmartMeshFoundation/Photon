package gopjnath

/*
#include <pjnath.h>
#include <pjlib-util.h>
#include <pjlib.h>
*/
import "C"

import (
	//"bytes"
	"encoding/binary"
	"net"
	"unsafe"
)

type AddressFamily uint16

var (
	AfUnspec = AddressFamily(C.PJ_AF_UNSPEC)
	AfUnix   = AddressFamily(C.PJ_AF_UNIX)
	AfIP     = AddressFamily(C.PJ_AF_INET)
	AfIPv6   = AddressFamily(C.PJ_AF_INET6)
	AfPacket = AddressFamily(C.PJ_AF_PACKET)
	AfIRDA   = AddressFamily(C.PJ_AF_IRDA)
)

// SockAddr describes a generic socket address.
type SockAddr struct {
	s *C.pj_sockaddr
}

func NewSockAddr(af AddressFamily, adr string, port uint16) (*SockAddr, error) {
	s := &SockAddr{}
	s.s = (*C.pj_sockaddr)(C.malloc(224))
	addrCh := C.CString(adr)
	defer C.free(unsafe.Pointer(addrCh))
	addr := C.pj_str(addrCh)
	//defer destroyString(addr)
	err := C.pj_sockaddr_init(C.int(af), s.s, &addr, C.pj_uint16_t(port))
	if err != C.PJ_SUCCESS {
		return s, casterr(err)
	}
	return s, nil
}

func (s *SockAddr) Destroy() {
	C.free(unsafe.Pointer(s.s))
}

func (s *SockAddr) Af() AddressFamily {
	return AddressFamily(binary.LittleEndian.Uint16(s.s[:1]))
}

func (s *SockAddr) IP() net.IP {
	addrPtr := unsafe.Pointer(C.pj_sockaddr_get_addr(unsafe.Pointer(s.s)))
	addrLen := C.int(C.pj_sockaddr_get_addr_len(unsafe.Pointer(s.s)))
	addr := C.GoBytes(addrPtr, addrLen)
	return net.IP(addr)
}

func (s *SockAddr) SetIP(ip net.IP) error {
	// get address family
	var af AddressFamily
	if ip.To4() == nil {
		af = AfIPv6
	} else {
		af = AfIP
	}
	addrCh := C.CString(ip.String())
	defer C.free(unsafe.Pointer(addrCh))
	addr := C.pj_str(addrCh)
	//defer destroyString(addr)
	err := C.pj_sockaddr_set_str_addr(C.int(af), s.s, &addr)
	return casterr(err)
}

func (s *SockAddr) Port() uint16 {
	// return binary.LittleEndian.Uint16(s.s[2:4])
	return uint16(C.pj_sockaddr_get_port(unsafe.Pointer(s.s)))
}

func (s *SockAddr) SetPort(i uint16) {
	C.pj_sockaddr_set_port(s.s, C.pj_uint16_t(i))
}
