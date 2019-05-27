package notify

import (
	"fmt"

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

// NewNotifyHandler :
func NewNotifyHandler() *Handler {
	return &Handler{
		receivedTransferChan: make(chan *models.ReceivedTransfer, 10),
		noticeChan:           make(chan *Notice, 10),
		stopped:              false,
	}
}

// Stop :
func (h *Handler) Stop() {
	h.stopped = true
	close(h.receivedTransferChan)
	close(h.noticeChan)
}

// GetNoticeChan :
// return read-only, keep chan private
func (h *Handler) GetNoticeChan() <-chan *Notice {
	return h.noticeChan
}

// GetReceivedTransferChan :
// keep chan private
func (h *Handler) GetReceivedTransferChan() <-chan *models.ReceivedTransfer {
	return h.receivedTransferChan
}

// Notify : 通知上层,不让阻塞,以免影响正常业务
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

// NotifyString : 通知上层,不让阻塞,以免影响正常业务
func (h *Handler) NotifyString(level Level, info string) {
	h.Notify(level, &InfoStruct{
		Type:    InfoTypeString,
		Message: info,
	})
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

// NotifyReceiveMediatedTransfer :通知收到了MediatedTransfer
func (h *Handler) NotifyReceiveMediatedTransfer(msg *encoding.MediatedTransfer, tokenAddress common.Address) {
	if h.stopped || msg == nil {
		return
	}
	info := fmt.Sprintf("收到token=%s,amount=%d,locksecrethash=%s的交易",
		utils.APex2(tokenAddress), msg.PaymentAmount, utils.HPex(msg.LockSecretHash))
	h.NotifyString(LevelInfo, info)
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
	h.Notify(LevelInfo, &InfoStruct{
		Type:    InfoTypeContractCallTXInfo,
		Message: txInfo,
	})
}

//NotifyInconsistentDatabase 通知在进行交易的时候发生了错误,因为交易双方的数据库不一致
func (h *Handler) NotifyInconsistentDatabase(channelIdentifier common.Hash, target common.Address) {
	log.Info(fmt.Sprintf("NotifyInconsistentDatabase on channel %s", channelIdentifier.String()))
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
