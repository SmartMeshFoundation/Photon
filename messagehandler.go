package smartraiden

import (
	"fmt"

	"math/big"

	"errors"

	"encoding/json"

	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/channel/channeltype"
	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/models"
	"github.com/SmartMeshFoundation/SmartRaiden/notify"
	"github.com/SmartMeshFoundation/SmartRaiden/rerr"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
 Class responsible to handle the protocol messages.

        This class is not intended to be used standalone, use RaidenService
        instead.
*/
type raidenMessageHandler struct {
	raiden        *RaidenService
	blockedTokens map[common.Address]bool
}

func newRaidenMessageHandler(raiden *RaidenService) *raidenMessageHandler {
	h := &raidenMessageHandler{
		raiden:        raiden,
		blockedTokens: make(map[common.Address]bool),
	}
	return h
}

/*
 Handles `message` and sends an ACK on success.
*/
func (mh *raidenMessageHandler) onMessage(msg encoding.SignedMessager, hash common.Hash) (err error) {
	msg.SetTag(&transfer.MessageTag{
		EchoHash: hash,
	})
	switch m2 := msg.(type) {
	case *encoding.SecretRequest:
		f := mh.raiden.SecretRequestPredictorMap[m2.LockSecretHash]
		if f != nil {
			ignore := (f)(m2)
			if ignore {
				return errors.New("ignore this secret request,because of SecretRequestPredictorMap ignores")
			}
		}
		err = mh.messageSecretRequest(m2)
	case *encoding.RevealSecret:
		f := mh.raiden.RevealSecretListenerMap[m2.LockSecretHash()]
		if f != nil {
			remove := (f)(m2)
			if remove {
				delete(mh.raiden.RevealSecretListenerMap, m2.LockSecretHash())
			}
		}
		err = mh.messageRevealSecret(m2) //has no relation with statemanager,duplicate message will be ok
	case *encoding.UnLock:
		err = mh.messageUnlock(m2)
	case *encoding.DirectTransfer:
		err = mh.messageDirectTransfer(m2)
	case *encoding.MediatedTransfer:
		err = mh.messageMediatedTransfer(m2)
		if err == nil {
			for f := range mh.raiden.ReceivedMediatedTrasnferListenerMap {
				remove := (*f)(m2)
				if remove {
					delete(mh.raiden.ReceivedMediatedTrasnferListenerMap, f)
				}
			}
		}
	case *encoding.AnnounceDisposed:
		err = mh.messageAnnounceDisposed(m2)
	case *encoding.AnnounceDisposedResponse:
		err = mh.messageAnnounceDisposedResponse(m2)
	case *encoding.RemoveExpiredHashlockTransfer:
		err = mh.messageRemoveExpiredHashlockTransfer(m2)
	case *encoding.SettleRequest:
		err = mh.messageSettleRequest(m2)
	case *encoding.SettleResponse:
		err = mh.messageSettleResponse(m2)
	case *encoding.WithdrawRequest:
		err = mh.messageWithdrawRequest(m2)
	case *encoding.WithdrawResponse:
		err = mh.messageWithdrawResponse(m2)
	default:
		log.Error(fmt.Sprintf("raidenMessageHandler unknown msg:%s", utils.StringInterface1(msg)))
		return fmt.Errorf("unhandled message cmdid:%d", msg.Cmd())
	}
	return err
}

func (mh *raidenMessageHandler) balanceProof(msg *encoding.UnLock) {
	blanceProof := transfer.NewBalanceProofStateFromEnvelopMessage(msg)
	balanceProof := &mediatedtransfer.ReceiveUnlockStateChange{
		LockSecretHash: msg.LockSecretHash(),
		NodeAddress:    msg.Sender,
		BalanceProof:   blanceProof,
		Message:        msg,
	}
	mh.raiden.StateMachineEventHandler.dispatchBySecretHash(balanceProof.LockSecretHash, balanceProof)
}

/*
 todo 收到密码,可能会影响到好多StateManager, 这些 StateManager 我如何做到原子保存呢?
*/
// todo receive secret may impact manay StateManager, how should I do atomic store for these StateManager?
func (mh *raidenMessageHandler) messageRevealSecret(msg *encoding.RevealSecret) error {
	secret := msg.LockSecret
	sender := msg.Sender
	mh.raiden.registerSecret(secret)
	stateChange := &mediatedtransfer.ReceiveSecretRevealStateChange{Secret: secret, Sender: sender, Message: msg}
	// save log to db
	channels := mh.raiden.findAllChannelsByLockSecretHash(msg.LockSecretHash())
	for _, c := range channels {
		mh.raiden.db.UpdateTransferStatusMessage(c.TokenAddress, msg.LockSecretHash(), fmt.Sprintf("收到 RevealSecret, from=%s", utils.APex2(msg.Sender)))
	}
	mh.raiden.StateMachineEventHandler.dispatchBySecretHash(msg.LockSecretHash(), stateChange)
	return nil
}

/*
引起的通道状态变化有 stateManager 触发保存机制
1. 找不到对应的 StateManager, 什么都不用做,忽略即可.
2. 找到了对应的 StateManager, 但是消息内容不对,比如 Amount 不对,忽略即可
3. 找到了对应的 StateManager, 并且消息内容正确,则肯定回发送 RevealSecret,这时保存消息收到并且更新状态
*/
/*
 *	messageSecretRequest : function to handle SecretRequest
 *
 *	Note that channel states should be triggered by stateManager to keep the record of these states.
 *	1. If there is no relevant StateManager, then do nothing.
 *	2. If we find relevant StateManager, but message content has faults, just neglect it.
 *	3. If we find relevant StateManager, and message is correct, then send RevealSecret, and store this message and switch channel state.
 */
func (mh *raidenMessageHandler) messageSecretRequest(msg *encoding.SecretRequest) error {
	stateChange := &mediatedtransfer.ReceiveSecretRequestStateChange{
		Amount:         new(big.Int).Set(msg.PaymentAmount),
		LockSecretHash: msg.LockSecretHash,
		Sender:         msg.Sender,
		Message:        msg,
	}
	// save log to db
	channels := mh.raiden.findAllChannelsByLockSecretHash(msg.LockSecretHash)
	for _, c := range channels {
		mh.raiden.db.UpdateTransferStatusMessage(c.TokenAddress, stateChange.LockSecretHash, fmt.Sprintf("收到 SecretRequest, from=%s", utils.APex2(msg.Sender)))
	}
	mh.raiden.StateMachineEventHandler.dispatchBySecretHash(stateChange.LockSecretHash, stateChange)
	return nil
}

/*
收到了 Unlock 消息,首先验证是否正确,如果不正确,则说明节点之间状态同步出错,通道只能关闭了
如果正确.
查找相应的StateManager, 肯定能找到,四中 StateManager 都有可能.
1. InititiatorStateManager ,则认为是错误的,忽略
2. MediatedStateManager 更新 State, 有可能认为这个交易结束(TransferPair 只有一个),也有可能继续(如果有多个 Transfer Pair).
3. TargetStateManager  更新 State, 认为交易彻底结束.
3. CrashStateManager 更新 State, 移除此锁
因此,适宜直接在此函数更新通道状态并保存ack
*/
/*
 *	messageUnlock : function to handle unlock.
 *
 *	Note that if received unlock message, first we need to verify if it is correct. If not, which means channel states are not synchronized, then close the channel.
 *	If correct, then just check relevant StateManager, there are four possible StateManager.
 *	1. InitiatorStateManager : assumed fault, neglect.
 *	2. MediatedStateManager : update channel state, may assume that this transfer completes (only one TransferPair), or this transfer goes on (multiple TransferPair).
 *	3. TargetStateManager : update channel state, assume this transfer completes.
 *	4. CrashStateManager : update channel state and remove the lock.
 *	So, it is reasonable that we update channel state and store ACK just in this function.
 */
func (mh *raidenMessageHandler) messageUnlock(msg *encoding.UnLock) error {
	lockSecretHash := msg.LockSecretHash()
	secret := msg.LockSecret
	mh.raiden.registerSecret(secret)
	var ch *channel.Channel
	var err error
	ch, err = mh.raiden.findChannelByIdentifier(msg.ChannelIdentifier)
	if err != nil {
		log.Info(fmt.Sprintf("Message for unknown channel: %s", err))
		return err
	}
	log.Trace(fmt.Sprintf("lockSecretHash=%s,nettingchannel=%s", utils.HPex(lockSecretHash), ch))
	/*
		收到unlock时,需要判断下通道的状态,如果该通道的状态已经不为open了,就不应该处理这笔unlock,否则有可能会损失钱.
		因为我已经提交过balance proof,如果不提交新的,我会损失钱,如果提交新的,那么之前在链上unlock过的锁,需要再unlock一遍,同样会损失gas
		所以应该拒绝该笔unlock,什么都不做
	*/
	/*
		When receive an unlock, I need to determine the state of the next channel.
		If the state of the channel is no longer open, I shouldn't process the unlock, otherwise I may lose money.
		Because I've already submitted balance proof, if I don't submit a new one, I'll lose money.
		If I submit a new one, then unlocked locks on the chain that were previously unlocked need to be unlocked again, and gas will also be lost.
		So we should abandon the unlock msg and do nothing.
	*/
	if !channeltype.CanDealUnlock[ch.State] {
		return errors.New("received unlock msg,but channel cannot deal unlock, do nothing")
	}
	err = ch.RegisterTransfer(mh.raiden.GetBlockNumber(), msg)
	if err != nil {
		log.Error(fmt.Sprintf("messageUnlock RegisterTransfer err=%s", err))
		return err
	}
	/*
		验证过消息是有效的,然后通知相应的 stateMana 该结束的结束,
	*/
	mh.balanceProof(msg)
	mh.raiden.updateChannelAndSaveAck(ch, msg.Tag())
	return nil
}

/*
如果消息错误,则说明节点之间状态同步出错,通道只能关闭了
相关的 StateManager 自己根据超时判断是否结束
适宜直接更新通道并保存 ack
*/
/*
 * messageRemoveExpiredHashlockTransfer : function to handle RemoveExpiredHashlock event.
 *
 *	Note that if message is faulty, which means channel states are not synchronized, this channel should be closed.
 *	Relevant StateManager should check this by settle_timeout
 *	Reasonable to update channel and store ACK.
 */
func (mh *raidenMessageHandler) messageRemoveExpiredHashlockTransfer(msg *encoding.RemoveExpiredHashlockTransfer) error {
	ch, err := mh.raiden.findChannelByIdentifier(msg.ChannelIdentifier)
	if err != nil {
		return fmt.Errorf("received  RemoveExpiredHashlockTransfer ,but relate channel cannot found %s", utils.StringInterface(msg, 7))
	}
	if !ch.CanContinueTransfer() {
		log.Warn(fmt.Sprintf("receive msg %s, but channel cannot continue transfer", msg))
		return nil
	}
	err = ch.RegisterRemoveExpiredHashlockTransfer(msg, mh.raiden.GetBlockNumber())
	if err != nil {
		log.Warn(fmt.Sprintf("RegisterRemoveExpiredHashlockTransfer err %s", err))
		/*
			这里不能直接丢弃掉消息,因为存在双方当前块不同步的情况,此时如果我丢弃了该条消息(本来是正确的),那么双方状态永远不会再同步了
			所以返回err让对方重发,如果是上诉情况,那么到时候会正常处理掉该消息,双方状态恢复正常.如果不是上诉情况,那么双方状态已经不同步了,
			一直重发也没什么
		*/
		return err
	}
	mh.raiden.updateChannelAndSaveAck(ch, msg.Tag())
	return nil
}

/*
收到 AnnounceDisposed 消息,如果错误,说明节点之间状态不同步了,通道只能关闭了
收到正常的的 AnnounceDisposed :
1. InitiatorStateManager 可能会选择其他节点进行路由,也可能会失败,但是肯定会发送 AnnounceDisposedResponse
2. MediatedStateManager 可能会选择其他节点进行路由,也可能会发送 AnnounceDisposed 表明交易无法继续,但是肯定会发送 AnnounceDisposedResponse
3. TargetStateManager  状态错误,不可能会出现
4. CrashStateManager 直接发送AnnounceDisposedResponse
因此通道状态更新以及保存消息收到,可以放在 EventSendAnnounceDisposedResponse
但是存在非原子更新的情况
1. 作为 InitiatorStateManager, 选择其他节点(EventSendMediatedTransfer)和发送EventSendAnnounceDisposedResponse无法原子处理
2. 作为MediatedStateManager 选择其他节点(EventSendMediatedTransfer)和发送EventSendAnnounceDisposedResponse无法原子处理
*/
/*
 * messageAnnounceDisposed : function to handle AnnounceDisposed message.
 *
 *	Note that when receiving AnnounceDisposed, if any fault occurs, which means channel states are not synchronized, this channel should be closed.
 *	When receiving normal AnnounceDisposed :
 *	1. InitiatorStateManager : Choose another node to route, and it may fail, but certainly it sends out AnnounceDisposedResponse.
 *	2. MediatedStateManager : Choose another node to route, may send AnnounceDisposed denoting that this transfer can't be furthered, but it must send AnnounceDisposedResponse.
 *	3. TargetStateManager : faulty channel state, impossible to occur.
 * 	4. CrashStateManager : immediately send AnnounceDisposedResponse
 *
 *	Hence, once channel states and stored messages are received, they can be put in EventSendAnnounceDisposedResponse.
 *	But there are cases of non-atomic update :
 *	1. As to InitiatorStateManager, there is no atomic operation between EventSendMediatedTransfer and EventSendAnnounceDisposedResponse.
 *	2. As to MediatedStateManager, there is no atomic operation between EventSendMediatedTransfer and EventSendAnnounceDisposedResponse.
 */
func (mh *raidenMessageHandler) messageAnnounceDisposed(msg *encoding.AnnounceDisposed) (err error) {
	graph := mh.raiden.getChannelGraph(msg.ChannelIdentifier)
	if graph == nil {
		return fmt.Errorf("unkonwn channel %s", msg.ChannelIdentifier.String())
	}
	if !graph.HasChannel(mh.raiden.NodeAddress, msg.Sender) {
		err = fmt.Errorf("direct transfer from node without an existing channel: %s", msg.Sender)
		return
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch == nil {
		return rerr.ChannelNotFound(fmt.Sprintf("channel:%s,partner:%s", utils.HPex(msg.ChannelIdentifier), utils.APex2(msg.Sender)))
	}
	err = ch.RegisterAnnouceDisposed(msg)
	if err != nil {
		return
	}
	punish := models.NewReceivedAnnounceDisposed(msg.Lock.Hash(), msg.ChannelIdentifier, msg.GetAdditionalHash(), msg.OpenBlockNumber, msg.Signature)
	err = mh.raiden.db.MarkLockHashCanPunish(punish)
	if err != nil {
		err = fmt.Errorf("MarkLockHashCanPunish %s err %s", utils.StringInterface(punish, 2), err)
		return
	}
	stateChange := &mediatedtransfer.ReceiveAnnounceDisposedStateChange{
		Sender:  msg.Sender,
		Token:   ch.TokenAddress,
		Lock:    msg.Lock,
		Message: msg,
	}
	mh.raiden.StateMachineEventHandler.dispatchBySecretHash(msg.Lock.LockSecretHash, stateChange)
	mh.raiden.db.UpdateTransferStatusMessage(ch.TokenAddress, msg.Lock.LockSecretHash, fmt.Sprintf("收到AnnounceDisposed from=%s", utils.APex2(msg.Sender)))
	return nil
}

/*
收到 AnnouceDisposedResponse,如果验证不通过,说明节点状态同步出了问题,通道只能关闭
收到正常的 AnnounceDisposedResponse:
相应的 StateManager 无需关心这个事件.
1. InitiatorStateManager 不可能收到,一定是个错误
2. MediatedStateManager 无需处理
3. TargetStateManager 不可能收到(目前是这样的,暂不允许接收方拒绝收款).
4. CrashStateManager 无需处理,等锁自动过期即可,因为这种情况,我是不会知道密码的.
因此适宜直接更新通道,并保存ack
*/
/*
 *	messageAnnounceDisposedResponse : function to handle AnnounceDisposedResponse event.
 *
 *	Note that when receiving AnnounceDisposed, if any fault occurs, which means channel states are not synchronized, this channel should be closed.
 *	When receiving normal AnnounceDisposedResponse :
 *	1. InitiatorStateManager : Cannot receive this event, faults occur.
 *	2. MediatedStateManager : No need to handle.
 *	3. TargetStateManager : impossible to receive this event.
 * 	4. CrashStateManager : No need to handle, just wait for expiration.
 *	Reasonable to update payment channel and store ACK.
 */
func (mh *raidenMessageHandler) messageAnnounceDisposedResponse(msg *encoding.AnnounceDisposedResponse) (err error) {
	graph := mh.raiden.getChannelGraph(msg.ChannelIdentifier)
	if graph == nil {
		return fmt.Errorf("unkonwn channel %s", msg.ChannelIdentifier.String())
	}
	if !graph.HasChannel(mh.raiden.NodeAddress, msg.Sender) {
		err = fmt.Errorf("direct transfer from node without an existing channel: %s", msg.Sender)
		return
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch == nil {
		return rerr.ChannelNotFound(fmt.Sprintf("channel:%s,partner:%s", utils.HPex(msg.ChannelIdentifier), utils.APex2(msg.Sender)))
	}
	/*
		必须验证我确实发送过这个Dispose
	*/
	// must check that I actually send this Dispose
	b := mh.raiden.db.IsLockSecretHashChannelIdentifierDisposed(msg.LockSecretHash, msg.ChannelIdentifier)
	if !b {
		return fmt.Errorf("maybe a attack, receive a announce disposed response,but i never send announce disposed,msg=%s", msg)
	}
	err = ch.RegisterTransfer(mh.raiden.GetBlockNumber(), msg)
	if err != nil {
		return
	}
	//保存通道状态即可.
	// Just store channel state.
	mh.raiden.updateChannelAndSaveAck(ch, msg.Tag())
	return nil
}

/*
如果验证错误,说明节点状态不同步,只能关闭通道
没有相关的 StateManager, 直接更新通道并保持 ack
*/
/*
 * messageDirectTransfer : function to handle directTransfer event.
 *
 *	Note that if verification is faulty, which means channel states are not synchronized, then we have to close this channel.
 *	There is no relevant StateManager, just update channel states and store ACK.
 */
func (mh *raidenMessageHandler) messageDirectTransfer(msg *encoding.DirectTransfer) error {
	// 用户调用了prepare-update,暂停接收新交易
	// halt new transfer because clients invoke prepare-update
	if mh.raiden.StopCreateNewTransfers {
		return rerr.ErrStopCreateNewTransfer
	}
	//mh.balanceProof(msg)
	graph := mh.raiden.getChannelGraph(msg.ChannelIdentifier)
	token := mh.raiden.getTokenForChannelIdentifier(msg.ChannelIdentifier)
	if graph == nil {
		return fmt.Errorf("unknown channel %s", utils.HPex(msg.ChannelIdentifier))
	}
	if _, ok := mh.blockedTokens[token]; ok {
		return rerr.ErrTransferUnwanted
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch == nil {
		return rerr.ChannelNotFound(fmt.Sprintf("token:%s,partner:%s", utils.APex2(token), utils.APex2(msg.Sender)))
	}
	if ch.State != channeltype.StateOpened {
		return rerr.TransferWhenClosed(ch.ChannelIdentifier.String())
	}
	var amount = new(big.Int)
	amount = amount.Sub(msg.TransferAmount, ch.PartnerState.TransferAmount())
	err := ch.RegisterTransfer(mh.raiden.GetBlockNumber(), msg)
	if err != nil {
		log.Error(fmt.Sprintf("RegisterTransfer error %s\n", msg))
		return err
	}
	receiveSuccess := &transfer.EventTransferReceivedSuccess{
		Amount:            amount,
		Initiator:         msg.Sender,
		ChannelIdentifier: msg.ChannelIdentifier,
	}
	mh.raiden.updateChannelAndSaveAck(ch, msg.Tag())
	err = mh.raiden.StateMachineEventHandler.OnEvent(receiveSuccess, nil)
	return err
}

/*
收到 MediatedTransfer, 如果验证不通过,说明节点之间状态不同步,通道只能关闭
验证通过:
1. 我是接收方, 创建 TargetStateManager, 并发送 SecretRequest, 适宜在 EventSendSecretRequest 更新通道状态并保存 ack
2. 我是中间节点
	2.1 第一次收到这个 LockSecretHash 创建 MediatedStateManager, 如果可以继续交易则发送 MediatedTransfer,
		否则发送 AnnounceDisposed,适宜在相关发送事件中更新通道状态并保存 ack
	2.2 第 n次收到这个 LockSecretHash, 可能的结果, 如果可以继续交易则发送MediatedTransfer,否则发送 AnnounceDisposed,适宜在相关发送事件中更新通道状态并保存 ack
3. 如果是 token swap
todo 需要设计如何保存 token swap 相关数据,并在崩溃恢复以后保证原子性.
*/
/*
 *	received MediatedTransfer. If verification fails, which means node states are not synchronous, and channel has to be closed.
 *
 *	Verification succeed :
 *		1. I am transfer initiator : Create TargetStateManager, send SecretRequest, and should update channel state in EventSendSecretRequest, and Store ACK.
 *		2. I am mediated node :
 *				2.1 First time receiving this LockSecretHash : Create a MediatedStateManager, if we can continue our transfer, then send MediatedTransfer,
 *				otherwise sending AnnounceDisposed, should update channel state and store ACK in relevant send event.
 *				2.2 n-th time receiving this LockSecretHash : if transfer can continue, then sending MediatedTransfer, otherwise, sending AnnounceDisposed,
 *				participants should update channel states and store ACK in relevant send event.
 *		3. if we use token swap.
 *		todo we should design how to store related data of token swap, and ensure atomicity after node crashes.
 */
func (mh *raidenMessageHandler) messageMediatedTransfer(msg *encoding.MediatedTransfer) error {
	// 用户调用了prepare-update,暂停接收新交易
	// Clients inovke prepare-update, stop receiving new transfers.
	if mh.raiden.StopCreateNewTransfers {
		return rerr.ErrStopCreateNewTransfer
	}
	token := mh.raiden.getTokenForChannelIdentifier(msg.ChannelIdentifier)
	if mh.raiden.Config.IgnoreMediatedNodeRequest && msg.Target != mh.raiden.NodeAddress {
		//todo what about return a AnnounceDisposed Message ?
		/*
			需要考虑恶意攻击的情况,比如发送一个我已经知道密码,但是尚未 unlock 的锁
		*/
		// We need to consider cases with potential attack risks, such as sending a lock that I know the secret but not yet unlock.
		return fmt.Errorf("ignored mh mediated transfer, because i don't want to route ")
	}
	if mh.raiden.Config.IsMeshNetwork {
		return fmt.Errorf("deny any mediated transfer when there is no internet connection")
	}
	if _, ok := mh.blockedTokens[token]; ok {
		return rerr.ErrTransferUnwanted
	}
	graph := mh.raiden.getToken2ChannelGraph(token)
	if graph == nil {
		return fmt.Errorf("received transfer on unkown token :%s", utils.APex2(token))
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch == nil {
		return rerr.ChannelNotFound(fmt.Sprintf("token:%s,partner:%s", utils.APex2(token), utils.APex2(msg.Sender)))
	}
	if !ch.CanTransfer() {
		return rerr.TransferWhenClosed(fmt.Sprintf("Mediated transfer received but the channel is  can not accept any transfer %s", ch.ChannelIdentifier.String()))
	}
	err := ch.RegisterTransfer(mh.raiden.GetBlockNumber(), msg)
	if err != nil {
		return err
	}
	// only for test
	dataForDebug := &struct {
		SearchKey           string
		TokenNetworkAddress string
		PartnerAddress      string
		TransferAmount      int64
		Expiration          int64
		Amount              int64
		LockSecretHash      string
		MerkleProof         []common.Hash
		Signature           string
	}{
		SearchKey:           "dataForDebug",
		TokenNetworkAddress: mh.raiden.Config.RegistryAddress.String(),
		PartnerAddress:      msg.Sender.String(),
		TransferAmount:      msg.TransferAmount.Int64(),
		Expiration:          msg.Expiration,
		Amount:              msg.PaymentAmount.Int64(),
		LockSecretHash:      msg.LockSecretHash.String(),
		MerkleProof:         channel.ComputeProofForLock(msg.GetLock(), ch.PartnerState.Tree).MerkleProof,
		Signature:           common.Bytes2Hex(msg.Signature),
	}

	buf, err := json.MarshalIndent(dataForDebug, "", "\t")
	log.Trace(string(buf))
	//mh.updateChannelAndSaveAck(ch, msg.Tag())
	if msg.Target == mh.raiden.NodeAddress {
		mh.raiden.targetMediatedTransfer(msg, ch)
	} else {
		mh.raiden.mediateMediatedTransfer(msg, ch)
	}
	/*
		start  taker's tokenswap ,only if receive a valid mediated transfer
	*/
	key := swapKey{
		LockSecretHash: msg.LockSecretHash,
		FromToken:      token,
		FromAmount:     msg.PaymentAmount.String(),
	}
	if tokenswap, ok := mh.raiden.SwapKey2TokenSwap[key]; ok {
		remove := mh.raiden.messageTokenSwapTaker(msg, tokenswap)
		if remove { //once the swap start,remove mh key immediately. otherwise,maker may repeat mh tokenswap operation.
			delete(mh.raiden.SwapKey2TokenSwap, key)
		}
		//return nil
	}
	return nil
}

/*
如果验证不通过,有可能是节点状态不同步,也有可能是其他链上异步事件导致,
无论那种情况,这个通道在对方看来都不能再继续使用,只能强制关闭.
直接更新通道状态并保存 ack 即可
*/
/*
 *	messageSettleRequest : function to handle SettleRequest event.
 *
 *	Note that if verification does not pass, the reason might be channel states are not synchronous,
 *	or maybe there are another on-chain asynchronous event.
 *  No matter which is our case, this channel are assumed to be not able to use, then enforce channel close.
 * 	Directly update channel states and store ACK.
 */
func (mh *raidenMessageHandler) messageSettleRequest(msg *encoding.SettleRequest) error {
	graph := mh.raiden.getChannelGraph(msg.ChannelIdentifier)
	token := mh.raiden.getTokenForChannelIdentifier(msg.ChannelIdentifier)
	if graph == nil {
		return fmt.Errorf("unknown channel %s", utils.HPex(msg.ChannelIdentifier))
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch == nil {
		return rerr.ChannelNotFound(fmt.Sprintf("token:%s,partner:%s", utils.APex2(token), utils.APex2(msg.Sender)))
	}
	if ch.State != channeltype.StateOpened {
		return fmt.Errorf("receive settle request but channel state is %s", ch.State)
	}
	err := ch.RegisterCooperativeSettleRequest(msg)
	if err != nil {
		log.Error(fmt.Sprintf("RegisterCooperativeSettleRequest error %s\n", err))
		return err
	}
	settleResponse, err := ch.CreateCooperativeSettleResponse(msg)
	if err != nil {
		//if err, channel can only be closed /settled
		log.Error(fmt.Sprintf("CreateCooperativeSettleResponse err %s", err))
		return err
	}
	// 如果这里有我发出的未解的锁,那么说明对方在老的balance_proof上withdraw,
	// 此时同意对我并没有坏处,所以正常返回response
	//if ch.HasAnyUnkonwnSecretTransferOnRoad() {
	//	//我自己理解 withdraw on channel就可以,防止上一笔交易额外损失
	//	result := ch.CooperativeSettleChannelOnRequest(msg.Participant1Signature, settleResponse)
	//	go func() {
	//		var err2 error
	//		err2 = <-result.Result
	//		if err2 != nil {
	//			log.Error(fmt.Sprintf("CooperativeSettleChannelOnRequest err %s", err2))
	//		} else {
	//			log.Info(fmt.Sprintf("CooperativeSettleChannelOnRequest success on channel %s", ch.ChannelIdentifier.String()))
	//		}
	//	}()
	//	return nil
	//}
	err = settleResponse.Sign(mh.raiden.PrivateKey, settleResponse)
	if err != nil {
		panic(fmt.Sprintf("sign message for settle response err %s", err))
	}
	err = mh.raiden.sendAsync(msg.Sender, settleResponse)
	if err != nil {
		log.Error(fmt.Sprintf("send message %s, to %s ,err %s", settleResponse, msg.Sender, err))
	}
	mh.raiden.updateChannelAndSaveAck(ch, msg.Tag())
	return nil
}
func (mh *raidenMessageHandler) messageSettleResponse(msg *encoding.SettleResponse) error {
	graph := mh.raiden.getChannelGraph(msg.ChannelIdentifier)
	token := mh.raiden.getTokenForChannelIdentifier(msg.ChannelIdentifier)
	if graph == nil {
		return fmt.Errorf("unknown channel %s", utils.HPex(msg.ChannelIdentifier))
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch == nil {
		return rerr.ChannelNotFound(fmt.Sprintf("token:%s,partner:%s", utils.APex2(token), utils.APex2(msg.Sender)))
	}
	/*
			会不会出现碰巧双方都发出了 settle request 这种情况?
			比如一方提出 settle, 另一方提出 withdraw,
			极小概率可能出现,如果出现,就退回到原始的 close/settle 模式
		也要防止另一方主动发出 settle response, 而我并没有发出 settle request 请求这种情况.
	*/
	/*
	 *	Is there any possibility that both participants send settle request?
	 *	Like one of them send settle, but the other send withdraw.
	 *	Not possible, if so, then revert to orignal close/settle mode.
	 *	Also, we need to prevent the other participant proactively send settle response, but I do not send settle request.
	 */
	if ch.State != channeltype.StateCooprativeSettle {
		return fmt.Errorf("receive settle response but channel state is %s", ch.State)
	}
	err := ch.RegisterCooperativeSettleResponse(msg)
	if err != nil {
		log.Error(fmt.Sprintf("RegisterCooperativeSettleResponse error %s\n", err))
		return err
	}
	mh.raiden.updateChannelAndSaveAck(ch, msg.Tag())
	result := ch.CooperativeSettleChannel(msg)
	go func() {
		err = <-result.Result
		if err != nil {
			log.Error(fmt.Sprintf("CooperativeSettleChannel %s failed, so we can only close/settle this channel, err = %s", utils.HPex(msg.ChannelIdentifier), err.Error()))
			mh.raiden.NotifyHandler.Notify(notify.LevelWarn, fmt.Sprintf("CooperateSettle通道失败,建议强制close/settle通道,ChannelIdentifier=%s", msg.ChannelIdentifier.String()))
		}
	}()
	return nil
}
func (mh *raidenMessageHandler) messageWithdrawRequest(msg *encoding.WithdrawRequest) error {
	graph := mh.raiden.getChannelGraph(msg.ChannelIdentifier)
	token := mh.raiden.getTokenForChannelIdentifier(msg.ChannelIdentifier)
	if graph == nil {
		return fmt.Errorf("unknown channel %s", utils.HPex(msg.ChannelIdentifier))
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch == nil {
		return rerr.ChannelNotFound(fmt.Sprintf("token:%s,partner:%s", utils.APex2(token), utils.APex2(msg.Sender)))
	}
	if ch.State != channeltype.StateOpened {
		return fmt.Errorf("receive settle request but channel state is %s", ch.State)
	}
	err := ch.RegisterWithdrawRequest(msg)
	if err != nil {
		log.Error(fmt.Sprintf("RegisterWithdrawRequest error %s\n", err))
		return err
	}
	// 现在只允许一方取现,直接构造response
	// Now we only allow one partcipant to withdraw, directly create response.
	withdrawResponse, err := ch.CreateWithdrawResponse(msg)
	if err != nil {
		//if err, channel can only be closed /settled
		log.Error(fmt.Sprintf("CreateWithdrawResponse err %s", err))
		return err
	}

	// 如果这里有我发出的未解的锁,那么说明对方在老的balance_proof上withdraw,
	// 此时同意对我并没有坏处,所以正常返回response
	//if ch.HasAnyUnkonwnSecretTransferOnRoad() {
	//	//我自己理解 withdraw on channel就可以,防止上一笔交易额外损失
	//	result := ch.WithdrawOnRequest(msg.Participant1Signature, withdrawResponse)
	//	go func() {
	//		var err2 error
	//		err2 = <-result.Result
	//		if err2 != nil {
	//			log.Error(fmt.Sprintf("WithdrawOnRequest err %s", err2))
	//		} else {
	//			log.Info(fmt.Sprintf("WithdrawOnRequest success on channel %s", ch.ChannelIdentifier.String()))
	//		}
	//	}()
	//	return nil
	//}
	err = withdrawResponse.Sign(mh.raiden.PrivateKey, withdrawResponse)
	if err != nil {
		panic(fmt.Sprintf("sign message for withdraw response err %s", err))
	}
	err = mh.raiden.sendAsync(msg.Sender, withdrawResponse)
	if err != nil {
		log.Error(fmt.Sprintf("send message %s, to %s ,err %s", withdrawResponse, msg.Sender, err))
	}
	mh.raiden.updateChannelAndSaveAck(ch, msg.Tag())
	return nil
}
func (mh *raidenMessageHandler) messageWithdrawResponse(msg *encoding.WithdrawResponse) error {
	graph := mh.raiden.getChannelGraph(msg.ChannelIdentifier)
	token := mh.raiden.getTokenForChannelIdentifier(msg.ChannelIdentifier)
	if graph == nil {
		return fmt.Errorf("unknown channel %s", utils.HPex(msg.ChannelIdentifier))
	}
	ch := graph.GetPartenerAddress2Channel(msg.Sender)
	if ch == nil {
		return rerr.ChannelNotFound(fmt.Sprintf("token:%s,partner:%s", utils.APex2(token), utils.APex2(msg.Sender)))
	}
	if ch.State != channeltype.StateWithdraw {
		return fmt.Errorf("receive settle request but channel state is %s", ch.State)
	}
	/*
		要先验证一下我发出去了 withdraw request,并且金额正确,然后才能注册
	*/
	// We need to verify withdraw request that is send and tokenAmount is correct, then register this request.
	err := ch.RegisterWithdrawResponse(msg)
	if err != nil {
		log.Error(fmt.Sprintf("RegisterTransfer error %s\n", msg))
		return err
	}
	mh.raiden.updateChannelAndSaveAck(ch, msg.Tag())
	//如果碰巧崩溃了,如果失败了,都只能回到 close/settle 这种老办法.
	// If crash happens, or register fails, we should revert to close/settle mode.
	result := ch.Withdraw(msg)
	go func() {
		err = <-result.Result
		if err != nil {
			log.Error(fmt.Sprintf("Withdraw %s failed, so we can only close/settle this channel", msg.ChannelIdentifier.String()))
		}
	}()
	return nil
}
