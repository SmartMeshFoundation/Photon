package notify

import (
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/Photon/rerr"

	"github.com/SmartMeshFoundation/Photon/log"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"

	"github.com/SmartMeshFoundation/Photon/encoding"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
Handler :
deal notice info for upper app
*/
type Handler struct {

	//receivedTransferChan  ReceivedTransfer notify, should never close
	receivedTransferChan chan *models.ReceivedTransfer
	//noticeChan should never close
	noticeChan chan *Notice
	// work status
	stopped bool
}

// NewNotifyHandler 创建一个photon给上层应用通知的handler
func NewNotifyHandler() *Handler {
	return &Handler{
		receivedTransferChan: make(chan *models.ReceivedTransfer, 10),
		noticeChan:           make(chan *Notice, 10),
		stopped:              false,
	}
}

// Stop 关闭,仅在photon停止时调用
func (h *Handler) Stop() {
	h.stopped = true
	close(h.receivedTransferChan)
	close(h.noticeChan)
}

// GetNoticeChan return read-only, keep chan private
func (h *Handler) GetNoticeChan() <-chan *Notice {
	return h.noticeChan
}

// GetReceivedTransferChan keep chan private
func (h *Handler) GetReceivedTransferChan() <-chan *models.ReceivedTransfer {
	return h.receivedTransferChan
}

// Notify 通知上层,不让阻塞,以免影响正常业务
func (h *Handler) Notify(level Level, info *InfoStruct) {
	if h.stopped || info == nil {
		return
	}
	select {
	case h.noticeChan <- newNotice(level, info):
	default:
		// never block
	}
}

// NotifySentTransferDetail : 通知上层,不让阻塞,以免影响正常业务
func (h *Handler) NotifySentTransferDetail(sentTransferDetail *models.SentTransferDetail) {
	h.Notify(LevelInfo, &InfoStruct{
		Type:    InfoTypeSentTransferDetail,
		Message: sentTransferDetail,
	})
}

const (
	//CallStatusFinishedSuccess 调用成功
	CallStatusFinishedSuccess = iota + 1
	//CallStatusError 失败
	CallStatusError
)

type channelCallIDResult struct {
	CallID       string      `json:"call_id"`
	Status       int         `json:"status"`
	ErrorMessage string      `json:"error_message"`
	Channel      interface{} `json:"channel"`
}

//NotifyChannelCallIDError 通知channel callid结果出错
func (h *Handler) NotifyChannelCallIDError(callID string, err error) {
	h.Notify(LevelInfo, &InfoStruct{
		Type: InfoTypeChannelCallID,
		Message: &channelCallIDResult{
			CallID:       callID,
			Status:       CallStatusError,
			ErrorMessage: err.Error(),
		},
	})
}

//NotifyChannelCallIDSuccess 通知channel callid成功
func (h *Handler) NotifyChannelCallIDSuccess(callID string, channel *channeltype.ChannelDataDetail) {
	h.Notify(LevelInfo, &InfoStruct{
		Type: InfoTypeChannelCallID,
		Message: &channelCallIDResult{
			CallID:       callID,
			Status:       CallStatusFinishedSuccess,
			ErrorMessage: "",
			Channel:      channel,
		},
	})
}

//NotifyChannelStatus 通知channel发生了变化,包括balance,locked_amount,state等等
func (h *Handler) NotifyChannelStatus(ch *channeltype.ChannelDataDetail) {
	//log.Trace(fmt.Sprintf("notify channel status changed:%s", utils.StringInterface(ch, 5)))
	h.Notify(LevelInfo, &InfoStruct{
		Type:    InfoTypeChannelStatus,
		Message: ch,
	})
}

type receivedTranser struct {
	Token      common.Address `json:"token"`
	From       common.Address `json:"from"`
	Amount     *big.Int       `json:"amount"`
	ID         string         `json:"id"`
	Expiration int64          `json:"expiration"`
}

// NotifyReceiveMediatedTransfer :通知收到了MediatedTransfer,和NotifyReceiveTransfer不一样,不能表示交易成功,
func (h *Handler) NotifyReceiveMediatedTransfer(msg *encoding.MediatedTransfer, tokenAddress common.Address) {
	if h.stopped || msg == nil {
		return
	}
	log.Info(fmt.Sprintf("NotifyReceiveMediatedTransfer token=%s,amount=%d,locksecrethash=%s的交易",
		utils.APex2(tokenAddress), msg.PaymentAmount, utils.HPex(msg.LockSecretHash)))
	h.Notify(LevelInfo, &InfoStruct{
		Type: InfoTypeReceivedMediatedTransfer,
		Message: &receivedTranser{
			Token:      tokenAddress,
			From:       msg.Initiator,
			Amount:     msg.PaymentAmount,
			ID:         msg.LockSecretHash.String(),
			Expiration: msg.Expiration,
		},
	})
}

// NotifyReceiveTransfer : 通知成功收到一笔token
func (h *Handler) NotifyReceiveTransfer(rt *models.ReceivedTransfer) {

	if h.stopped || rt == nil {
		return
	}
	select {
	case h.receivedTransferChan <- rt:
	default:
		// never block
	}
}

/*
NotifyContractCallTXInfo 当自己发起的合约调用tx被成功打包时,通知上层
*/
func (h *Handler) NotifyContractCallTXInfo(txInfo *models.TXInfo) {
	if h.stopped {
		return
	}
	h.Notify(LevelInfo, &InfoStruct{
		Type:    InfoTypeContractCallTXInfo,
		Message: txInfo,
	})
}

//NotifyInconsistentDatabase 通知在进行交易的时候发生了错误,因为交易双方的数据库不一致
func (h *Handler) NotifyInconsistentDatabase(channelIdentifier common.Hash, target common.Address) {
	log.Info(fmt.Sprintf("NotifyInconsistentDatabase on channel %s", channelIdentifier.String()))
	if h.stopped {
		return
	}
	type inconsistentDatabase struct {
		ChannelIdentifier common.Hash    `json:"channel_identifier"`
		Target            common.Address `json:"target"`
	}
	h.Notify(LevelInfo, &InfoStruct{
		Type: InfoTypeInconsistentDatabase,
		Message: inconsistentDatabase{
			ChannelIdentifier: channelIdentifier,
			Target:            target,
		},
	})
}

// NotifyPhotonBalanceNotEnough 通知上层账户余额不足
func (h *Handler) NotifyPhotonBalanceNotEnough(balance, needed *big.Int) {
	log.Warn(fmt.Sprintf("NotifyPhotonBalanceNotEnough balance=%d needed=%d", balance, needed))
	if h.stopped {
		return
	}
	type notEnough struct {
		Need *big.Int `json:"need"`
		Have *big.Int `json:"have"`
	}
	h.Notify(LevelError, &InfoStruct{
		Type: InfoTypeBalanceNotEnoughError,
		Message: &notEnough{
			Need: needed,
			Have: balance,
		},
	})
}

type failedCooperate struct {
	Channel   common.Hash `json:"channel"`
	ErrorCode int         `json:"error_code"`
	ErrorMsg  string      `json:"error_message"`
}

func newFailedCooperate(channel common.Hash, err error) *failedCooperate {
	errCode := rerr.ErrUnknown.ErrorCode
	if e2, ok := err.(*rerr.StandardError); ok {
		errCode = e2.ErrorCode
	}

	return &failedCooperate{
		Channel:   channel,
		ErrorCode: errCode,
		ErrorMsg:  err.Error(),
	}
}
func (h *Handler) notifyCooperateFailed(typ int, msg *failedCooperate) {
	if h.stopped {
		return
	}
	h.Notify(LevelError, &InfoStruct{
		Type:    typ,
		Message: msg,
	})
}

//NotifyCooperateSettleRefused 通知对方拒绝合作关闭通道
func (h *Handler) NotifyCooperateSettleRefused(channel common.Hash, err error) {
	log.Warn(fmt.Sprintf("NotifyCooperateSettleRefused on channel=%s,reason=%s", utils.HPex(channel), err))
	h.notifyCooperateFailed(InfoTypeCooperateSettleRefused, newFailedCooperate(channel, err))
}

//NotifyCooperateSettleFailed 通知Tx执行失败
func (h *Handler) NotifyCooperateSettleFailed(channel common.Hash, err error) {
	log.Warn(fmt.Sprintf("NotifyCooperateSettleFailed on channel=%s,err=%s", utils.HPex(channel), err))
	h.notifyCooperateFailed(InfoTypeCooperateSettleFailed, newFailedCooperate(channel, err))

}

//NotifyWithdrawRefused 通知对方拒绝合作取现
func (h *Handler) NotifyWithdrawRefused(channel common.Hash, err error) {
	log.Warn(fmt.Sprintf("NotifyWithdrawRefused on channel=%s,reason=%s", utils.HPex(channel), err))
	h.notifyCooperateFailed(InfoTypeWithdrawRefused, newFailedCooperate(channel, err))

}

//NotifyWithdrawFailed 通知Tx执行失败
func (h *Handler) NotifyWithdrawFailed(channel common.Hash, err error) {
	log.Warn(fmt.Sprintf("NotifyWithdrawFailed on channel=%s,err=%s", utils.HPex(channel), err))
	h.notifyCooperateFailed(InfoTypeWithdrawFailed, newFailedCooperate(channel, err))
}
