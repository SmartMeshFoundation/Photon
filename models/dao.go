package models

import (
	"math/big"
	"time"

	"github.com/asdine/storm"

	"github.com/SmartMeshFoundation/Photon/rerr"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/models/cb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// KeyGetter for StormTx
type KeyGetter interface {
	GetKey() []byte
}

// TX for StormTx
type TX interface {
	Set(table string, key interface{}, value interface{}) error
	Save(v KeyGetter) error
	Commit() error
	Rollback() error
}

// AckDao 管理节点间通信ack的存储
type AckDao interface {
	GetAck(echoHash common.Hash) []byte
	SaveAck(echoHash common.Hash, ack []byte, tx TX)
	SaveAckNoTx(echoHash common.Hash, ack []byte)
}

// BlockNumberDao 管理当前最新块的存储
type BlockNumberDao interface {
	GetLatestBlockNumber() int64
	SaveLatestBlockNumber(blockNumber int64)
	GetLastBlockNumberTime() time.Time
}

// ChainIDDao 管理ChainID的存储,用作启动时校验
type ChainIDDao interface {
	GetChainID() int64
	SaveChainID(chainID int64)
}

//ChannelUpdateDao 该接口定义了所有关于通道更新操作
type ChannelUpdateDao interface {
	// update
	UpdateChannel(c *channeltype.Serialization, tx TX) error
	UpdateChannelNoTx(c *channeltype.Serialization) error
	UpdateChannelState(c *channeltype.Serialization) error
	// mix update
	UpdateChannelAndSaveAck(c *channeltype.Serialization, echoHash common.Hash, ack []byte) (err error)
	UpdateChannelContractBalance(c *channeltype.Serialization) error
}

// ChannelDao 该接口定义了通道内容存储及查询相关方法
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

// UnlockDao 保存自己在链上unlock成功的锁
type UnlockDao interface {
	IsThisLockHasUnlocked(channelIdentifier common.Hash, lockHash common.Hash) bool
	UnlockThisLock(channelIdentifier common.Hash, lockHash common.Hash)
}

// ExpiredLockDao 保存自己移除过的锁
type ExpiredLockDao interface {
	IsThisLockRemoved(channelIdentifier common.Hash, sender common.Address, lockHash common.Hash) bool
	RemoveLock(channelIdentifier common.Hash, sender common.Address, lockHash common.Hash)
}

// DbStatusDao 启动时在db中标记状态
type DbStatusDao interface {
	MarkDbOpenedStatus()
	IsDbCrashedLastTime() bool
}

// ContractStatusDao 保存合约相关常量,启动时使用或校验使用
type ContractStatusDao interface {
	SaveContractStatus(contractStatus ContractStatus)
	GetContractStatus() ContractStatus
}

// SentEnvelopMessagerDao 保存正在发送给其他节点的的带有BalanceProof的消息
type SentEnvelopMessagerDao interface {
	NewSentEnvelopMessager(msg encoding.EnvelopMessager, receiver common.Address)
	DeleteEnvelopMessager(echohash common.Hash)
	GetAllOrderedSentEnvelopMessager() []*SentEnvelopMessager
}

// FeeChargeRecordDao 保存节点手续费收取的详细流水
type FeeChargeRecordDao interface {
	SaveFeeChargeRecord(r *FeeChargeRecord) (err error)
	GetAllFeeChargeRecord(tokenAddress common.Address, fromTime, toTime int64) (records []*FeeChargeRecord, err error)
	GetFeeChargeRecordByLockSecretHash(lockSecretHash common.Hash) (records []*FeeChargeRecord, err error)
}

// FeePolicyDao 保存节点的手续费收取策略
type FeePolicyDao interface {
	SaveFeePolicy(fp *FeePolicy) (err error)
	GetFeePolicy() (fp *FeePolicy)
}

// NonParticipantChannelDao 保存所有通道的参与者信息,与Channel表不同,这里存储的信息比较简单
type NonParticipantChannelDao interface {
	NewNonParticipantChannel(token common.Address, channelIdentifier common.Hash, participant1, participant2 common.Address) error
	RemoveNonParticipantChannel(channel common.Hash) error
	GetAllNonParticipantChannelByToken(token common.Address) (edges []common.Address, err error)
	GetNonParticipantChannelByID(channelIdentifierForQuery common.Hash) (
		tokenAddress common.Address, participant1, participant2 common.Address, err error)
}

// SentAnnounceDisposedDao 记录自己声明放弃过的锁,防止自己误unlock
type SentAnnounceDisposedDao interface {
	MarkLockSecretHashDisposed(lockSecretHash common.Hash, channelIdentifier common.Hash) error
	IsLockSecretHashDisposed(lockSecretHash common.Hash) bool
	IsLockSecretHashChannelIdentifierDisposed(lockSecretHash common.Hash, ChannelIdentifier common.Hash) bool
	GetSendAnnounceDisposeByChannel(channelIdentifier common.Hash, isSubmitToPms bool) (list []*SentAnnounceDisposed)
	MarkSendAnnounceDisposeSubmittedByChannel(channelIdentifier common.Hash)
}

// ReceivedAnnounceDisposedDao 记录收到的通道partner声明放弃过的锁,punish时使用
type ReceivedAnnounceDisposedDao interface {
	MarkLockHashCanPunish(r *ReceivedAnnounceDisposed) error
	IsLockHashCanPunish(lockHash, channelIdentifier common.Hash) bool
	GetReceivedAnnounceDisposed(lockHash, channelIdentifier common.Hash) *ReceivedAnnounceDisposed
	GetChannelAnnounceDisposed(channelIdentifier common.Hash) []*ReceivedAnnounceDisposed
	MarkLockHashCanPunishSubmittedByChannel(channelIdentifier common.Hash)
}

// SettledChannelDao 记录历史上打开过,但已经settle掉的channel部分信息
type SettledChannelDao interface {
	NewSettledChannel(c *channeltype.Serialization) error
	GetAllSettledChannel() (chs []*channeltype.Serialization, err error)
	GetSettledChannel(channelIdentifier common.Hash, openBlockNumber int64) (c *channeltype.Serialization, err error)
}

// TokenDao 记录当前photon网络中,存在通道的token列表
type TokenDao interface {
	GetAllTokens() (tokens AddressMap, err error)
	AddToken(token common.Address, tokenNetworkAddress common.Address) error
}

// ReceivedTransferDao 记录他人通过photon发送给我的交易
type ReceivedTransferDao interface {
	NewReceivedTransfer(blockNumber int64, channelIdentifier common.Hash, openBlockNumber int64, tokenAddr, fromAddr common.Address, nonce uint64, amount *big.Int, lockSecretHash common.Hash, data string) *ReceivedTransfer
	GetReceivedTransfer(key string) (*ReceivedTransfer, error)
	GetReceivedTransferList(tokenAddress common.Address, fromBlock, toBlock, fromTime, toTime int64) (transfers []*ReceivedTransfer, err error)
}

// SentTransferDetailDao 记录自己发送的交易
type SentTransferDetailDao interface {
	NewSentTransferDetail(tokenAddress, target common.Address, amount *big.Int, data string, isDirect bool, lockSecretHash common.Hash)
	UpdateSentTransferDetailStatus(tokenAddress common.Address, lockSecretHash common.Hash, status TransferStatusCode, statusMessage string, otherParams interface{}) (transfer *SentTransferDetail)
	UpdateSentTransferDetailStatusMessage(tokenAddress common.Address, lockSecretHash common.Hash, statusMessage string) (transfer *SentTransferDetail)
	GetSentTransferDetail(tokenAddress common.Address, lockSecretHash common.Hash) (*SentTransferDetail, error)
	GetSentTransferDetailList(tokenAddress common.Address, fromTime, toTime int64, fromBlock, toBlock int64) (transfers []*SentTransferDetail, err error)
}

// XMPPSubDao 保存xmpp节点状态订阅情况
type XMPPSubDao interface {
	XMPPMarkAddrSubed(addr common.Address)
	XMPPIsAddrSubed(addr common.Address) bool
	XMPPUnMarkAddr(addr common.Address)
}

// TXInfoDao 保存photon节点发起的所有合约操作流水
type TXInfoDao interface {
	NewPendingTXInfo(tx *types.Transaction, txType TXInfoType, channelIdentifier common.Hash, openBlockNumber int64, txParams TXParams, isFake ...bool) (txInfo *TXInfo, err error)
	SaveEventToTXInfo(event interface{}) (txInfo *TXInfo, err error)
	UpdateTXInfoStatus(txHash common.Hash, status TXInfoStatus, pendingBlockNumber int64, gasUsed uint64) (txInfo *TXInfo, err error)
	GetTXInfoList(channelIdentifier common.Hash, openBlockNumber int64, tokenAddress common.Address, txType TXInfoType, status TXInfoStatus) (list []*TXInfo, err error)
}

// ChainEventRecordDao 记录链上事件处理流水,暂时没有使用
type ChainEventRecordDao interface {
	NewDeliveredChainEvent(id ChainEventID, blockNumber uint64)
	CheckChainEventDelivered(id ChainEventID) (blockNumber uint64, delivered bool)
	ClearOldChainEventRecord(blockNumber uint64)
	MakeChainEventID(l *types.Log) ChainEventID
}

// UnlockToSendDao 用于暂存无网状态下需要发送的unlock消息,等待恢复到有网之后再进行发送
type UnlockToSendDao interface {
	NewUnlockToSend(lockSecretHash common.Hash, tokenAddress, receiver common.Address, blockNumber int64) *UnlockToSend
	GetAllUnlockToSend() (list []*UnlockToSend)
	RemoveUnlockToSend(key []byte)
}

// Dao 汇总所有dao层接口
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
	ChainEventRecordDao
	UnlockToSendDao

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
		if err == storm.ErrNotFound {
			return rerr.ErrNotFound
		}
		return rerr.ErrGeneralDBError.Append(err.Error())
	}
	return nil
}
