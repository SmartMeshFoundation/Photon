package models

import (
	"encoding/gob"
	"time"

	"sort"

	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/ethereum/go-ethereum/common"
)

/*
SentEnvelopMessager is record of envelop message,that don't received a Ack
*/
type SentEnvelopMessager struct {
	Message  encoding.EnvelopMessager
	Receiver common.Address
	Time     time.Time
	EchoHash []byte `storm:"id"`
}

type envelopMessageSorter []*SentEnvelopMessager

func (c envelopMessageSorter) Len() int {
	return len(c)
}
func (c envelopMessageSorter) Less(i, j int) bool {
	m1 := c[i].Message.GetEnvelopMessage()
	m2 := c[j].Message.GetEnvelopMessage()
	return m1.Nonce < m2.Nonce
}
func (c envelopMessageSorter) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

/*
SortEnvelopMessager :
must be stable
对于 ChannelOpenedAndDeposit 事件,会产生两个 stateChange,
严格要求有先后顺序
*/
/*
 *	sortEnvelopMessager : function to sort arrays of sent messenger.
 *
 *	Note that for event of ChannelOpenedAndDeposit, two stateChange will be generated.
 *	And they must be in order.
 */
func SortEnvelopMessager(msgs []*SentEnvelopMessager) {
	sort.Stable(envelopMessageSorter(msgs))
}

func init() {
	gob.Register(&SentEnvelopMessager{})
}
