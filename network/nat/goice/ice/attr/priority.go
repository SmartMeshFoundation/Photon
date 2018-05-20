package attr

import (
	"strconv"

	"encoding/binary"

	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/stun"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/turn"
)

//Priority  https://trac.tools.ietf.org/html/rfc5245#section-19.1
type Priority uint32 // encoded as uint16

func (n Priority) String() string { return strconv.Itoa(int(n)) }

/*
It is a 32-bit unsigned integer, and has an attribute
value of 0x0024.
*/
const prioritySize = 4

// AddTo adds PRIORITY to message.
func (n Priority) AddTo(m *stun.Message) error {
	v := make([]byte, prioritySize)
	binary.BigEndian.PutUint32(v, uint32(n))
	// v[2:4] are zeroes (RFFU = 0)
	m.Add(stun.AttrPriority, v)
	return nil
}

// GetFrom decodes PRIORITY from message.
func (n *Priority) GetFrom(m *stun.Message) error {
	v, err := m.Get(stun.AttrPriority)
	if err != nil {
		return err
	}
	if len(v) != prioritySize {
		return &turn.BadAttrLength{
			Attr:     stun.AttrPriority,
			Got:      len(v),
			Expected: prioritySize,
		}
	}
	*n = Priority(binary.BigEndian.Uint32(v))
	return nil
}
