package models

import (
	"math/big"
	"testing"

	"math/rand"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/assert"
)

func TestModelDB_NewSentEnvelopMessager(t *testing.T) {
	m := setupDb(t)
	bp := &encoding.BalanceProof{
		Nonce:             11,
		ChannelIdentifier: utils.Sha3([]byte("123")),
		TransferAmount:    big.NewInt(12),
		OpenBlockNumber:   3,
		Locksroot:         utils.EmptyHash,
	}
	p := encoding.NewDirectTransfer(bp)
	receiverPrivKey, receiver := utils.MakePrivateKeyAddress()
	err := p.Sign(receiverPrivKey, p)
	if err != nil {
		t.Error(err)
	}
	m.NewSentEnvelopMessager(p, receiver)
	msgs := m.GetAllOrderedSentEnvelopMessager()
	assert.EqualValues(t, len(msgs), 1)
	echohash := utils.Sha3(p.Pack(), receiver[:])
	m.DeleteEnvelopMessager(echohash)
	msgs = m.GetAllOrderedSentEnvelopMessager()
	assert.EqualValues(t, len(msgs), 0)
}

func TestModelDB_NewSentEnvelopMessager2(t *testing.T) {
	m := setupDb(t)
	msgs := m.GetAllOrderedSentEnvelopMessager()
	assert.EqualValues(t, len(msgs), 0)
	m.DeleteEnvelopMessager(utils.NewRandomHash())
	msgs = m.GetAllOrderedSentEnvelopMessager()
	assert.EqualValues(t, len(msgs), 0)
}

func TestModelDB_NewSentEnvelopMessager3(t *testing.T) {
	m := setupDb(t)
	var msgs []*encoding.DirectTransfer
	total := 10
	var min int64 = math.MaxInt64
	for i := 0; i < total; i++ {
		bp := &encoding.BalanceProof{
			Nonce:             int64(rand.Int63()),
			ChannelIdentifier: utils.Sha3([]byte("123")),
			TransferAmount:    big.NewInt(12),
			OpenBlockNumber:   3,
			Locksroot:         utils.EmptyHash,
		}
		p := encoding.NewDirectTransfer(bp)
		msgs = append(msgs, p)
		receiverPrivKey, receiver := utils.MakePrivateKeyAddress()
		err := p.Sign(receiverPrivKey, p)
		if err != nil {
			t.Error(err)
		}
		m.NewSentEnvelopMessager(p, receiver)
		if bp.Nonce < min {
			min = bp.Nonce
		}
	}

	smsgs := m.GetAllOrderedSentEnvelopMessager()
	assert.EqualValues(t, len(smsgs), 10)
	for i := 0; i < len(smsgs); i++ {
		s := smsgs[i]
		if s.Message.GetEnvelopMessage().Nonce < min {
			panic("order error")
		}
		min = s.Message.GetEnvelopMessage().Nonce
	}
}
