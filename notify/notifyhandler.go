package notify

import (
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/models"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

/*
Handler :
deal notice info for upper app
*/
type Handler struct {

	//sentTransferChan SentTransfer notify ,should never close
	sentTransferChan chan *models.SentTransfer
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
		sentTransferChan:     make(chan *models.SentTransfer),
		receivedTransferChan: make(chan *models.ReceivedTransfer),
		noticeChan:           make(chan *Notice),
		stopped:              false,
	}
}

// Stop :
func (h *Handler) Stop() {
	h.stopped = true
	close(h.sentTransferChan)
	close(h.receivedTransferChan)
	close(h.noticeChan)
}

// GetNoticeChan :
// return read-only, keep chan private
func (h *Handler) GetNoticeChan() <-chan *Notice {
	return h.noticeChan
}

// GetSentTransferChan :
// keep chan private
func (h *Handler) GetSentTransferChan() <-chan *models.SentTransfer {
	return h.sentTransferChan
}

// GetReceivedTransferChan :
// keep chan private
func (h *Handler) GetReceivedTransferChan() <-chan *models.ReceivedTransfer {
	return h.receivedTransferChan
}

// Notify : 通知上层,不让阻塞,以免影响正常业务
func (h *Handler) Notify(level Level, info interface{}) {
	if h.stopped || info == nil || info == "" {
		return
	}
	select {
	case h.noticeChan <- newNotice(level, info):
	default:
		// never block
	}
}

// NotifyReceiveMediatedTransfer :
func (h *Handler) NotifyReceiveMediatedTransfer(msg *encoding.MediatedTransfer, ch *channel.Channel) {
	if h.stopped || msg == nil {
		return
	}
	info := fmt.Sprintf("收到token=%s,amount=%d,locksecrethash=%s的交易",
		utils.APex2(ch.TokenAddress), msg.PaymentAmount, utils.HPex(msg.LockSecretHash))
	select {
	case h.noticeChan <- newNotice(LevelInfo, info):
	default:
		// never block
	}
}

// NotifySentTransfer :
func (h *Handler) NotifySentTransfer(st *models.SentTransfer) {
	if h.stopped || st == nil {
		return
	}
	select {
	case h.sentTransferChan <- st:
	default:
		// never block
	}
}

// NotifyReceiveTransfer :
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
