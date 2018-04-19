// Package turn implements RFC 5766 Traversal Using Relays around NAT.
package turn

import (
	"encoding/binary"
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/stun"
)

// bin is shorthand for binary.BigEndian.
var bin = binary.BigEndian

// BadAttrLength means that length for attribute is invalid.
type BadAttrLength struct {
	Attr     stun.AttrType
	Got      int
	Expected int
}

func (e BadAttrLength) Error() string {
	return fmt.Sprintf("incorrect length for %s: got %d, expected %d",
		e.Attr,
		e.Got,
		e.Expected,
	)
}

// Default ports for TURN from RFC 5766 Section 4.
const (
	// DefaultPort for TURN is same as STUN.
	DefaultPort = stun.DefaultPort
	// DefaultTLSPort is for TURN over TLS and is same as STUN.
	DefaultTLSPort = stun.DefaultTLSPort
)

var (
	// AllocateRequest is shorthand for allocation request message type.
	AllocateRequest = stun.NewType(stun.MethodAllocate, stun.ClassRequest)
	// CreatePermissionRequest is shorthand for create permission request type.
	CreatePermissionRequest = stun.NewType(stun.MethodCreatePermission, stun.ClassRequest)
	// CreatePermissionResponse is shorthand for create permission response type
	CreatePermissionResponse = stun.NewType(stun.MethodCreatePermission, stun.ClassSuccessResponse)
	// SendIndication is shorthand for send indication message type to turn server.
	SendIndication = stun.NewType(stun.MethodSend, stun.ClassIndication)
	//DataIndication is shorthand for data indication message type from turn server.
	DataIndication = stun.NewType(stun.MethodData, stun.ClassIndication)
	//ChannelBind Request
	ChannelBindRequest = stun.NewType(stun.MethodChannelBind, stun.ClassRequest)
	// RefreshRequest is shorthand for refresh request message type.
	RefreshRequest = stun.NewType(stun.MethodRefresh, stun.ClassRequest)
	// ChannelData Request //request is fake
	ChannelDataRequest = stun.NewType(stun.MethodChannelData, stun.ClassRequest)
	// Refresh success
	RefreshResponse = stun.NewType(stun.MethodRefresh, stun.ClassSuccessResponse)
)
