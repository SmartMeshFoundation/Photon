package turn

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/stun"
)

const MinChannelNumber = 0x4000
const MaxChannelNumber = 0x7fff

type ChannelData struct {
	ChannelNumber uint16 //must between 0x4000-0x7fff
	Data          []byte //data to send, can be empty.
}

func (c *ChannelData) String() string {
	return fmt.Sprintf("{channel number=%d,data len=%d}", c.ChannelNumber, len(c.Data))
}

// AddTo adds Channel Data to message.
func (c *ChannelData) AddTo(m *stun.Message) error {
	if m.Type.Method != stun.MethodChannelData {
		return fmt.Errorf("channel data works only on channelData method, now is %s", m.Type.String())
	}
	buf := new(bytes.Buffer)
	if c.ChannelNumber < MinChannelNumber || c.ChannelNumber > MaxChannelNumber {
		return fmt.Errorf("channel number is invalid :%d", c.ChannelNumber)
	}
	binary.Write(buf, binary.BigEndian, c.ChannelNumber)
	binary.Write(buf, binary.BigEndian, uint16(len(c.Data)))
	binary.Write(buf, binary.BigEndian, c.Data) //todo fix padding, when works on tcp... let 4 bytes align
	/*
			   Over TCP and TLS-over-TCP, the ChannelData message MUST be padded to
		   a multiple of four bytes in order to ensure the alignment of
		   subsequent messages.  The padding is not reflected in the length
		   field of the ChannelData message, so the actual size of a ChannelData
		   message (including padding) is (4 + Length) rounded up to the nearest
		   multiple of 4.  Over UDP, the padding is not required but MAY be
		   included.
	*/
	m.Raw = buf.Bytes()
	return nil
}

// GetFrom decodes Channel Data from message.
func (c *ChannelData) GetFrom(m *stun.Message) error {
	if m.Type.Method != stun.MethodChannelData {
		return fmt.Errorf("expect MethodChannelData,but got %s", m.String())
	}
	buf := bytes.NewBuffer(m.Raw)
	err := binary.Read(buf, binary.BigEndian, &c.ChannelNumber)
	if err != nil {
		return err
	}
	if c.ChannelNumber < MinChannelNumber || c.ChannelNumber > MaxChannelNumber {
		return fmt.Errorf("channel number is invalid :%d", c.ChannelNumber)
	}
	var length uint16
	err = binary.Read(buf, binary.BigEndian, &length)
	if err != nil {
		return err
	}
	c.Data = make([]byte, length)
	err = binary.Read(buf, binary.BigEndian, c.Data)
	if err != nil {
		return err
	}
	return nil
}
