package models

import (
	"encoding/gob"

	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
SentAnnounceDisposed 我发出了 AnnonuceDisposed, 那么就要标记这个 channel 上的这个锁我不要去链上兑现了,
如果对方发送过来 AnnounceDisposedResponse, 我要移除这个锁.
*/
/*
 *	SendAnnounceDisposed : a channel participant sends out an AnnounceDisposed message, then that lock in this channel should be tagged
 *	to inform that this participant should not register it on chain.
 *	If his partner sends AnnounceDisposedResponse, that lock has to be removed.
 */
type SentAnnounceDisposed struct {
	Key               []byte `storm:"id"`
	LockSecretHash    []byte `storm:"index"` //假设非恶意的情况下,锁肯定是不会重复的.但是我有可能在多个通道上发送 AnnounceDisposed,但是肯定不会在同一个通道上发送多次 announce disposed // Assume in honest case, locks are not repeated, but maybe I send AnnounceDisposed in multiple channels.
	ChannelIdentifier common.Hash
}

/*
ReceivedAnnounceDisposed 收到对方的 disposed, 主要是用来对方unlock 的时候,提交证据,惩罚对方
*/
/*
 *	ReceiveAnnounceDisposed : to receive AnnounceDisposed message from channel partner,
 *	mainly to submit proofs and punish fraudulent behaviors while partner submits unlock.
 */
type ReceivedAnnounceDisposed struct {
	Key               []byte `storm:"id"`
	LockHash          []byte `storm:"index"` //hash(expiration,locksecrethash,amount)
	ChannelIdentifier []byte `storm:"index"`
	OpenBlockNumber   int64
	AdditionalHash    common.Hash
	Signature         []byte
	IsSubmittedToPms  bool
}

func init() {
	gob.Register(&SentAnnounceDisposed{})
	gob.Register(&ReceivedAnnounceDisposed{})
}

//NewReceivedAnnounceDisposed create ReceivedAnnounceDisposed
func NewReceivedAnnounceDisposed(LockHash, ChannelIdentifier, additionalHash common.Hash, openBlockNumber int64, signature []byte) *ReceivedAnnounceDisposed {
	key := utils.Sha3(LockHash[:], ChannelIdentifier[:])
	return &ReceivedAnnounceDisposed{
		Key:               key[:],
		LockHash:          LockHash[:],
		ChannelIdentifier: ChannelIdentifier[:],
		OpenBlockNumber:   openBlockNumber,
		AdditionalHash:    additionalHash,
		Signature:         signature,
		IsSubmittedToPms:  false,
	}
}

func init() {
	gob.Register(&SentAnnounceDisposed{})
	gob.Register(&ReceivedAnnounceDisposed{})
}
