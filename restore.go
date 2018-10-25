package smartraiden

import (
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/channel"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mediatedtransfer/crashnode"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mtree"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

/*
重启完毕以后,根据数据库中保存的数据,恢复操作
1. 未发送成功的 EnvelopMessage 继续发送
2. 持有的锁,建立对应的 StateManager, 对这些未完成的交易进行简单维护处理
*/
/*
 *	restore : function to restore data.
 *
 *	Note that
 *		1. unsuccessful EnvelopMessages resume to be sent.
 *		2. to create related StateManager as to those locks withholden by a particpant.
 */
func (rs *RaidenService) restore() {
	//1. 处理未完成的锁
	// 1. handle incomplete locks
	rs.restoreLocks()
	//2. 为发送成功的 EnvelopMessage 继续发送
	// 2. keep sending EnvelopMessage that failed previously.
	rs.reSendEnvelopMessage()
}
func (rs *RaidenService) reSendEnvelopMessage() {
	msgs := rs.db.GetAllOrderedSentEnvelopMessager()
	for _, msg := range msgs {
		/*
			todo 存在问题:
			1.应该等待历史消息处理完毕以后再发送消息
			2. 有些消息是不需要发送的,比如 unlock 消息,如果对在链上注册了密码,那么会在其他地方触发发送 unlock 消息
			3. 已经过期的MediatedTransfer?
			4. 通道已经 settle 的消息?
		*/
		/*
		 *	todo have problem here :
		 *	1. we should wait after handling history message then continue.
		 *	2. Some messages are no need to send, like unlock, if pairs register their secret on chain, then other pairs will be triggered to send unlock.
		 *	3. How to deal with expired MediatedTransfer?
		 *	4. How to deal with messages after channel settle?
		 */
		err := rs.sendAsync(msg.Receiver, msg.Message)
		if err != nil {
			log.Error(fmt.Sprintf("reSendEnvelopMessage %s to %s err %s", msg.Message, msg.Receiver, err))
		}
	}
}

type lockInfo struct {
	l      *mtree.Lock
	isSent bool
	token  common.Address
	ch     *channel.Channel
}

func (rs *RaidenService) restoreLocks() {
	token2ActionInitCrashRestartStateChange := make(map[common.Hash]*mediatedtransfer.ActionInitCrashRestartStateChange)
	var locks []*lockInfo
	//收集所有的锁,
	// collect all locks.
	for token := range rs.Token2TokenNetwork {
		g := rs.Token2ChannelGraph[token]
		for _, ch := range g.ChannelIdentifier2Channel {
			for _, l := range ch.OurState.Lock2PendingLocks {
				locks = append(locks, &lockInfo{
					l:      l.Lock,
					isSent: true,
					token:  token,
					ch:     ch,
				})
			}
			for _, l := range ch.OurState.Lock2UnclaimedLocks {
				locks = append(locks, &lockInfo{
					l:      l.Lock,
					isSent: true,
					token:  token,
					ch:     ch,
				})
			}
			for _, l := range ch.PartnerState.Lock2PendingLocks {
				locks = append(locks, &lockInfo{
					l:      l.Lock,
					isSent: false,
					token:  token,
					ch:     ch,
				})
			}
			for _, l := range ch.PartnerState.Lock2UnclaimedLocks {
				locks = append(locks, &lockInfo{
					l:      l.Lock,
					isSent: false,
					token:  token,
					ch:     ch,
				})
			}
		}
	}
	//log.Trace(fmt.Sprintf("after restart current locks %s", utils.StringInterface(locks, 4)))
	//将 lock 转换为ActionInitCrashRestartStateChange
	// switch lock to ActionInitCrashRestartStateChange
	for _, l := range locks {
		key := utils.Sha3(l.l.LockSecretHash[:], l.token[:])
		aicr := token2ActionInitCrashRestartStateChange[key]
		if aicr == nil {
			aicr = &mediatedtransfer.ActionInitCrashRestartStateChange{
				OurAddress:     rs.NodeAddress,
				Token:          l.token,
				LockSecretHash: l.l.LockSecretHash,
			}
		}
		if l.isSent {
			aicr.SentLocks = append(aicr.SentLocks, &mediatedtransfer.LockAndChannel{
				Lock:    l.l,
				Channel: l.ch,
			})
		} else {
			aicr.ReceivedLocks = append(aicr.ReceivedLocks, &mediatedtransfer.LockAndChannel{
				Lock:    l.l,
				Channel: l.ch,
			})
		}
		token2ActionInitCrashRestartStateChange[key] = aicr
	}
	//log.Trace(fmt.Sprintf("after restart ActionInitCrashRestartStateChanges=%s", utils.StringInterface(token2ActionInitCrashRestartStateChange, 5)))
	//根据ActionInitCrashRestartStateChange,创建对应的 stateManager
	// Create corresponding stateManager, according to ActionInitCrashRestartStateChange.
	for k, st := range token2ActionInitCrashRestartStateChange {
		stateManager := transfer.NewStateManager(crashnode.StateTransition, nil, crashnode.NameCrashNodeTransition, st.LockSecretHash, st.Token)
		rs.Transfer2StateManager[k] = stateManager
		rs.StateMachineEventHandler.dispatch(stateManager, st)
	}
}
