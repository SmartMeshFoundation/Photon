package models

import (
	"math/big"
	"time"

	"github.com/SmartMeshFoundation/Photon/rerr"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/models/cb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// KeyGetter :
type KeyGetter interface {
	GetKey() []byte
}

// TX :
type TX interface {
	Set(table string, key interface{}, value interface{}) error
	Save(v KeyGetter) error
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

//ChannelUpdateDao update channel status in db
type ChannelUpdateDao interface {
	// update
	UpdateChannel(c *channeltype.Serialization, tx TX) error
	UpdateChannelNoTx(c *channeltype.Serialization) error
	UpdateChannelState(c *channeltype.Serialization) error
	// mix update
	UpdateChannelAndSaveAck(c *channeltype.Serialization, echoHash common.Hash, ack []byte) (err error)
	UpdateChannelContractBalance(c *channeltype.Serialization) error
}

// ChannelDao :
type ChannelDao interface {
	ChannelUpdateDao
	// add
	NewChannel(c *channeltype.Serialization) error
	// remove
	RemoveChannel(c *channeltype.Serialization) error
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

// ContractStatusDao :
type ContractStatusDao interface {
	SaveContractStatus(contractStatus ContractStatus)
	GetContractStatus() ContractStatus
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
	RemoveNonParticipantChannel(channel common.Hash) error
	GetAllNonParticipantChannelByToken(token common.Address) (edges []common.Address, err error)
	GetNonParticipantChannelByID(channelIdentifierForQuery common.Hash) (
		tokenAddress common.Address, participant1, participant2 common.Address, err error)
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

//// SentTransferDao :
//type SentTransferDao interface {
//	NewSentTransfer(blockNumber int64, channelIdentifier common.Hash, openBlockNumber int64, tokenAddr, toAddr common.Address, nonce uint64, amount *big.Int, lockSecretHash common.Hash, data string) *SentTransfer
//	GetSentTransfer(key string) (*SentTransfer, error)
//	GetSentTransferInBlockRange(fromBlock, toBlock int64) (transfers []*SentTransfer, err error)
//	GetSentTransferInTimeRange(from, to time.Time) (transfers []*SentTransfer, err error)
//}

// ReceivedTransferDao :
type ReceivedTransferDao interface {
	NewReceivedTransfer(blockNumber int64, channelIdentifier common.Hash, openBlockNumber int64, tokenAddr, fromAddr common.Address, nonce uint64, amount *big.Int, lockSecretHash common.Hash, data string) *ReceivedTransfer
	GetReceivedTransfer(key string) (*ReceivedTransfer, error)
	GetReceivedTransferList(tokenAddress common.Address, fromBlock, toBlock int64) (transfers []*ReceivedTransfer, err error)
}

//// TransferStatusDao :
//type TransferStatusDao interface {
//	NewTransferStatus(tokenAddress common.Address, lockSecretHash common.Hash)
//	UpdateTransferStatus(tokenAddress common.Address, lockSecretHash common.Hash, status TransferStatusCode, statusMessage string)
//	UpdateTransferStatusMessage(tokenAddress common.Address, lockSecretHash common.Hash, statusMessage string)
//	GetTransferStatus(tokenAddress common.Address, lockSecretHash common.Hash) (*TransferStatus, error)
//}

// SentTransferDetailDao :
type SentTransferDetailDao interface {
	NewSentTransferDetail(tokenAddress, target common.Address, amount *big.Int, data string, isDirect bool, lockSecretHash common.Hash)
	UpdateSentTransferDetailStatus(tokenAddress common.Address, lockSecretHash common.Hash, status TransferStatusCode, statusMessage string, otherParams interface{})
	UpdateSentTransferDetailStatusMessage(tokenAddress common.Address, lockSecretHash common.Hash, statusMessage string)
	GetSentTransferDetail(tokenAddress common.Address, lockSecretHash common.Hash) (*SentTransferDetail, error)
	GetSentTransferDetailList(tokenAddress common.Address, fromTime, toTime int64, fromBlock, toBlock int64) (transfers []*SentTransferDetail, err error)
}

// XMPPSubDao :
type XMPPSubDao interface {
	XMPPMarkAddrSubed(addr common.Address)
	XMPPIsAddrSubed(addr common.Address) bool
	XMPPUnMarkAddr(addr common.Address)
}

// TXInfoDao :
type TXInfoDao interface {
	NewPendingTXInfo(tx *types.Transaction, txType TXInfoType, channelIdentifier common.Hash, openBlockNumber int64, txParams TXParams) (txInfo *TXInfo, err error)
	SaveEventToTXInfo(event interface{}) (txInfo *TXInfo, err error)
	UpdateTXInfoStatus(txHash common.Hash, status TXInfoStatus, pendingBlockNumber int64) (err error)
	GetTXInfoList(channelIdentifier common.Hash, openBlockNumber int64, tokenAddress common.Address, txType TXInfoType, status TXInfoStatus) (list []*TXInfo, err error)
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
	ContractStatusDao
	SentEnvelopMessagerDao
	FeeChargeRecordDao
	FeePolicyDao
	NonParticipantChannelDao
	SentAnnounceDisposedDao
	ReceivedAnnounceDisposedDao
	SettledChannelDao
	TokenDao
	ReceivedTransferDao
	XMPPSubDao
	TXInfoDao
	SentTransferDetailDao

	StartTx() (tx TX)
	CloseDB()

	RegisterNewTokenCallback(f cb.NewTokenCb)
	RegisterNewChannelCallback(f cb.ChannelCb)
	RegisterChannelDepositCallback(f cb.ChannelCb)
	RegisterChannelStateCallback(f cb.ChannelCb)
	RegisterChannelSettleCallback(f cb.ChannelCb)
}

//GeneratDBError helper function
func GeneratDBError(err error) error {
	if err != nil {
		e2, ok := err.(rerr.StandardError)
		if ok {
			return e2
		}
		return rerr.ErrGeneralDBError.Append(err.Error())
	}
	return nil
}
