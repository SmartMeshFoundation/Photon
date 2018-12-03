package models

import (
	"math/big"
	"time"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/models/cb"
	"github.com/ethereum/go-ethereum/common"
)

// TX :
type TX interface {
	Set(table string, key interface{}, value interface{}) error
	Save(v interface{}) error
	Commit() error
	Rollback() error
}

// AckDao :
type AckDao interface {
	GetAck(echoHash common.Hash) []byte
	SaveAck(echoHash common.Hash, ack []byte, tx TX)
	SaveAckNoTx(echoHash common.Hash, ack []byte)
}

// BlockNumberDao :
type BlockNumberDao interface {
	GetLatestBlockNumber() int64
	SaveLatestBlockNumber(blockNumber int64)
	GetLastBlockNumberTime() time.Time
}

// ChainIDDao :
type ChainIDDao interface {
	GetChainID() int64
	SaveChainID(chainID int64)
}

// ChannelDao :
type ChannelDao interface {
	// add
	NewChannel(c *channeltype.Serialization) error
	// remove
	RemoveChannel(c *channeltype.Serialization) error
	// update
	UpdateChannel(c *channeltype.Serialization, tx TX) error
	UpdateChannelNoTx(c *channeltype.Serialization) error
	UpdateChannelState(c *channeltype.Serialization) error
	// mix update
	UpdateChannelAndSaveAck(c *channeltype.Serialization, echoHash common.Hash, ack []byte) (err error)
	UpdateChannelContractBalance(c *channeltype.Serialization) error
	//query
	GetChannel(token, partner common.Address) (c *channeltype.Serialization, err error)
	GetChannelByAddress(channelIdentifier common.Hash) (c *channeltype.Serialization, err error)
	GetChannelList(token, partner common.Address) (cs []*channeltype.Serialization, err error)
}

// UnlockDao :
type UnlockDao interface {
	IsThisLockHasUnlocked(channelIdentifier common.Hash, lockHash common.Hash) bool
	UnlockThisLock(channelIdentifier common.Hash, lockHash common.Hash)
}

// ExpiredLockDao :
type ExpiredLockDao interface {
	IsThisLockRemoved(channelIdentifier common.Hash, sender common.Address, lockHash common.Hash) bool
	RemoveLock(channelIdentifier common.Hash, sender common.Address, lockHash common.Hash)
}

// DbStatusDao :
type DbStatusDao interface {
	MarkDbOpenedStatus()
	IsDbCrashedLastTime() bool
}

// RegistryAddressDao :
type RegistryAddressDao interface {
	SaveRegistryAddress(registryAddress common.Address)
	GetRegistryAddress() common.Address
}

// SentEnvelopMessagerDao :
type SentEnvelopMessagerDao interface {
	NewSentEnvelopMessager(msg encoding.EnvelopMessager, receiver common.Address)
	DeleteEnvelopMessager(echohash common.Hash)
	GetAllOrderedSentEnvelopMessager() []*SentEnvelopMessager
}

// FeeChargeRecordDao :
type FeeChargeRecordDao interface {
	SaveFeeChargeRecord(r *FeeChargeRecord) (err error)
	GetAllFeeChargeRecord() (records []*FeeChargeRecord, err error)
	GetFeeChargeRecordByLockSecretHash(lockSecretHash common.Hash) (records []*FeeChargeRecord, err error)
}

// FeePolicyDao :
type FeePolicyDao interface {
	SaveFeePolicy(fp *FeePolicy) (err error)
	GetFeePolicy() (fp *FeePolicy)
}

// NonParticipantChannelDao :
type NonParticipantChannelDao interface {
	NewNonParticipantChannel(token common.Address, channelIdentifier common.Hash, participant1, participant2 common.Address) error
	RemoveNonParticipantChannel(token common.Address, channelIdentifier common.Hash) error
	GetAllNonParticipantChannel(token common.Address) (edges []common.Address, err error)
	GetParticipantAddressByTokenAndChannel(token common.Address, channel common.Hash) (p1, p2 common.Address)
}

// SentAnnounceDisposedDao :
type SentAnnounceDisposedDao interface {
	MarkLockSecretHashDisposed(lockSecretHash common.Hash, channelIdentifier common.Hash) error
	IsLockSecretHashDisposed(lockSecretHash common.Hash) bool
	IsLockSecretHashChannelIdentifierDisposed(lockSecretHash common.Hash, ChannelIdentifier common.Hash) bool
}

// ReceivedAnnounceDisposedDao :
type ReceivedAnnounceDisposedDao interface {
	MarkLockHashCanPunish(r *ReceivedAnnounceDisposed) error
	IsLockHashCanPunish(lockHash, channelIdentifier common.Hash) bool
	GetReceivedAnnounceDisposed(lockHash, channelIdentifier common.Hash) *ReceivedAnnounceDisposed
	GetChannelAnnounceDisposed(channelIdentifier common.Hash) []*ReceivedAnnounceDisposed
}

// SettledChannelDao :
type SettledChannelDao interface {
	NewSettledChannel(c *channeltype.Serialization) error
	GetAllSettledChannel() (chs []*channeltype.Serialization, err error)
	GetSettledChannel(channelIdentifier common.Hash, openBlockNumber int64) (c *channeltype.Serialization, err error)
}

// TokenDao :
type TokenDao interface {
	GetAllTokens() (tokens AddressMap, err error)
	AddToken(token common.Address, tokenNetworkAddress common.Address) error
}

// SentTransferDao :
type SentTransferDao interface {
	NewSentTransfer(blockNumber int64, channelIdentifier common.Hash, tokenAddr, toAddr common.Address, nonce uint64, amount *big.Int, lockSecretHash common.Hash, data string) *SentTransfer
	GetSentTransfer(key string) (*SentTransfer, error)
	GetSentTransferInBlockRange(fromBlock, toBlock int64) (transfers []*SentTransfer, err error)
}

// ReceivedTransferDao :
type ReceivedTransferDao interface {
	NewReceivedTransfer(blockNumber int64, channelIdentifier common.Hash, tokenAddr, fromAddr common.Address, nonce uint64, amount *big.Int, lockSecretHash common.Hash, data string) *ReceivedTransfer
	GetReceivedTransfer(key string) (*ReceivedTransfer, error)
	GetReceivedTransferInBlockRange(fromBlock, toBlock int64) (transfers []*ReceivedTransfer, err error)
}

// TransferStatusDao :
type TransferStatusDao interface {
	NewTransferStatus(tokenAddress common.Address, lockSecretHash common.Hash)
	UpdateTransferStatus(tokenAddress common.Address, lockSecretHash common.Hash, status TransferStatusCode, statusMessage string)
	UpdateTransferStatusMessage(tokenAddress common.Address, lockSecretHash common.Hash, statusMessage string)
	GetTransferStatus(tokenAddress common.Address, lockSecretHash common.Hash) (*TransferStatus, error)
}

// XMPPSubDao :
type XMPPSubDao interface {
	XMPPMarkAddrSubed(addr common.Address)
	XMPPIsAddrSubed(addr common.Address) bool
	XMPPUnMarkAddr(addr common.Address)
}

// Dao :
type Dao interface {
	AckDao
	BlockNumberDao
	ChainIDDao
	ChannelDao
	UnlockDao
	ExpiredLockDao
	DbStatusDao
	RegistryAddressDao
	SentEnvelopMessagerDao
	FeeChargeRecordDao
	FeePolicyDao
	NonParticipantChannelDao
	SentAnnounceDisposedDao
	ReceivedAnnounceDisposedDao
	SettledChannelDao
	TokenDao
	SentTransferDao
	ReceivedTransferDao
	TransferStatusDao
	XMPPSubDao

	StartTx() (tx TX)
	CloseDB()

	RegisterNewTokenCallback(f cb.NewTokenCb)
	RegisterNewChannelCallback(f cb.ChannelCb)
	RegisterChannelDepositCallback(f cb.ChannelCb)
	RegisterChannelStateCallback(f cb.ChannelCb)
	RegisterChannelSettleCallback(f cb.ChannelCb)
}
