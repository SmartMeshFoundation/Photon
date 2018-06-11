package smartraiden

import (
	"math/big"
	"testing"

	"bytes"
	"encoding/gob"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer/initiator"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

func TestStateManager(t *testing.T) {
	stateManager := transfer.NewStateManager(initiator.StateTransition, nil, initiator.NameInitiatorTransition, 1, utils.NewRandomAddress())
	lock := &encoding.Lock{
		Amount:     big.NewInt(34),
		Expiration: 4589895, //expiration block number
		HashLock:   utils.Sha3([]byte("hashlock")),
	}
	m1 := encoding.NewMediatedTransfer(11, 32, utils.NewRandomAddress(), utils.NewRandomAddress(), big.NewInt(33), utils.NewRandomAddress(),
		utils.Sha3([]byte("ddd")), lock, utils.NewRandomAddress(), utils.NewRandomAddress(), big.NewInt(33))
	tag := &transfer.MessageTag{
		MessageID:         utils.RandomString(10),
		EchoHash:          utils.Sha3(m1.Pack(), m1.Target[:]),
		IsASendingMessage: true,
		Receiver:          m1.Target,
	}
	m1.SetTag(tag)
	stateManager.LastSendMessage = m1
	buf := new(bytes.Buffer)
	err := gob.NewEncoder(buf).Encode(stateManager)
	if err != nil {
		t.Error(err)
		return
	}
	var sm2 *transfer.StateManager
	err = gob.NewDecoder(buf).Decode(&sm2)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("sm1=%s", utils.StringInterface(stateManager, 7))
	t.Logf("sm2=%s", utils.StringInterface(sm2, 7))

	var m2 *encoding.MediatedTransfer
	buf2 := new(bytes.Buffer)

	err = gob.NewEncoder(buf2).Encode(m1)
	gob.NewDecoder(buf2).Decode(&m2)
	t.Logf("m1=%s", utils.StringInterface(m1, 7))
	t.Logf("m2=%s", utils.StringInterface(m2, 7))
	if m2.Tag() == nil {
		t.Errorf("tag not saved")
		return
	}
}
