package transfer

import (
	"encoding/gob"

	"math/big"

	"bytes"
	"encoding/binary"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//BalanceProofState is   proof need by contract
type BalanceProofState struct {
	Nonce             uint64
	TransferAmount    *big.Int
	LocksRoot         common.Hash
	ChannelIdentifier contracts.ChannelUniqueID
	MessageHash       common.Hash
	//signature is nonce + transferred_amount + locksroot + channel_identifier + message_hash
	Signature []byte

	/*
		由于合约上并没有存储transferamount 和 locksroot,
		而用户 unlock 的时候会改变对方的 TransferAmount, 虽然说这个没有对方的签名,但是必须凭此在合约上settle 以及 unlock
	*/
	// Because contract does not cache transferAmount and locksroot
	// and my partner's transferAmount will be changed, even I have no signature of my partner, but we should settle and unlock via it.
	ContractTransferAmount *big.Int
	ContractNonce          uint64
	ContractLocksRoot      common.Hash
}

//NewEmptyBalanceProofState init BalanceProof with proper state
func NewEmptyBalanceProofState() *BalanceProofState {
	return &BalanceProofState{
		TransferAmount:         new(big.Int),
		ContractTransferAmount: new(big.Int),
	}
}

//NewBalanceProofState create BalanceProofState
func NewBalanceProofState(nonce uint64, transferAmount *big.Int, locksRoot common.Hash,
	channelIdentifier contracts.ChannelUniqueID, messageHash common.Hash, signature []byte) *BalanceProofState {
	s := &BalanceProofState{
		Nonce:                  nonce,
		TransferAmount:         new(big.Int).Set(transferAmount),
		LocksRoot:              locksRoot,
		ChannelIdentifier:      channelIdentifier,
		MessageHash:            messageHash,
		Signature:              signature,
		ContractTransferAmount: new(big.Int),
		ContractNonce:          nonce,
		//ContractLocksRoot:      locksRoot,
	}
	if s.TransferAmount == nil {
		s.TransferAmount = new(big.Int)
	}
	return s
}

//NewBalanceProofStateFromEnvelopMessage from locked transfer
func NewBalanceProofStateFromEnvelopMessage(msg encoding.EnvelopMessager) *BalanceProofState {
	envmsg := msg.GetEnvelopMessage()
	msgHash := encoding.HashMessageWithoutSignature(msg)
	return NewBalanceProofState(envmsg.Nonce, envmsg.TransferAmount,
		envmsg.Locksroot,
		contracts.ChannelUniqueID{
			ChannelIdentifier: envmsg.ChannelIdentifier,
			OpenBlockNumber:   envmsg.OpenBlockNumber,
		},
		msgHash, envmsg.Signature)
}

//IsBalanceProofValid true if valid
func (bpf *BalanceProofState) IsBalanceProofValid() bool {
	var err error
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, bpf.Nonce)
	_, err = buf.Write(utils.BigIntTo32Bytes(bpf.TransferAmount))
	_, err = buf.Write(bpf.LocksRoot[:])
	_, err = buf.Write(bpf.ChannelIdentifier.ChannelIdentifier[:])
	_, err = buf.Write(bpf.MessageHash[:])
	dataToSign := buf.Bytes()

	hash := utils.Sha3(dataToSign)
	signature := make([]byte, len(bpf.Signature))
	copy(signature, bpf.Signature)
	signature[len(signature)-1] -= 27 //why?
	pubkey, err := crypto.Ecrecover(hash[:], signature)
	//log.Trace(fmt.Sprintf("signer =%s",utils.APex(utils.PubkeyToAddress(pubkey))))
	return err == nil && utils.PubkeyToAddress(pubkey) != utils.EmptyAddress
}

//StateName name of state
func (bpf *BalanceProofState) StateName() string {
	return "BalanceProofState"
}

func init() {
	gob.Register(&BalanceProofState{})
}
