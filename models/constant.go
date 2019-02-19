package models

import "github.com/ethereum/go-ethereum/common"

// DbVersion :
const DbVersion = 1

// ChannelParticipantMap : used by BucketChannel
type ChannelParticipantMap map[common.Hash][]byte

/*
 #no-golint
*/
const (
	BucketMeta        = "meta"
	BucketAck         = "ack"
	BucketBlockNumber = "bucketBlockNumber"
	BucketChainID     = "bucketChainID"
	/*
		保存channel完整信息
	*/
	BucketChannelSerialization = "bucketChannelSerialization"
	/*
		保存所有通道的ChannelParticipantMap
	*/
	BucketChannel = "bucketChannel"
	/*
	   保留 settle 的通道信息,供查询需要
	*/
	BucketSettledChannel = "settled_channel"

	BucketToken      = "bucketToken"
	BucketTokenNodes = "bucketTokenNodes"
	BucketXMPP       = "bucketxmpp"

	/*
		保存已经解锁的锁
	*/
	BucketWithDraw = "bucketWithdraw"
	/*
		保存已经过期的锁
	*/
	BucketExpiredHashlock          = "expiredHashlock"
	BucketEnvelopMessager          = "EnvelopMessager"
	BucketFeeChargeRecord          = "FeeChargeRecord"
	BucketFeePolicy                = "FeePolicy"
	BucketSentAnnounceDisposed     = "SentAnnounceDisposed"
	BucketReceivedAnnounceDisposed = "ReceivedAnnounceDisposed"
	BucketSentTransfer             = "SentTransfer"
	BucketReceivedTransfer         = "ReceivedTransfer"
	BucketTransferStatus           = "TransferStatus"
	BucketTXInfo                   = "TXInfo"
	BucketSentTransferDetail       = "SentTransferDetail"
)

/*
 #no-golint
*/
const (
	// keys of BucketMeta
	KeyVersion        = "version"
	KeyCloseFlag      = "close"
	KeyRegistry       = "registry"
	KeySecretRegistry = "secretregistry"

	// keys of BucketBlockNumber
	KeyBlockNumber     = "blocknumber"
	KeyBlockNumberTime = "blockTime"

	// keys of BucketChainID
	KeyChainID = "chainID"

	// keys of BucketFeePolicy
	KeyFeePolicy string = "feePolicy"
	// keys of BucketToken
	KeyToken = "tokens"
)
