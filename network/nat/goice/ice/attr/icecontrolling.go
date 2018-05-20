package attr

import (
	"encoding/binary"
	"math/rand"
	"strconv"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/stun"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/turn"
)

//IceControlling role of ice
type IceControlling uint64

func (n IceControlling) String() string { return strconv.Itoa(int(n)) }

const iceControllingSize = 8

// AddTo adds ICE-CONTROLLING  to message.
func (n IceControlling) AddTo(m *stun.Message) error {
	v := make([]byte, iceControllingSize)
	binary.BigEndian.PutUint64(v, uint64(n))
	m.Add(stun.AttrICEControlling, v)
	return nil
}

// GetFrom decodes ICE-CONTROLLING from message.
func (n *IceControlling) GetFrom(m *stun.Message) error {
	v, err := m.Get(stun.AttrICEControlling)
	if err != nil {
		return err
	}
	if len(v) != iceControllingSize {
		return &turn.BadAttrLength{
			Attr:     stun.AttrICEControlling,
			Got:      len(v),
			Expected: iceControllingSize,
		}
	}
	*n = IceControlling(binary.BigEndian.Uint64(v))
	return nil
}

//RandUint64 random uint64
func RandUint64() uint64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Uint64()
}
