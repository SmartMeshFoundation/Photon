package photon

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/channel"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/transfer"
	"github.com/SmartMeshFoundation/Photon/transfer/mediatedtransfer"
	"github.com/SmartMeshFoundation/Photon/transfer/mediatedtransfer/crashnode"
	"github.com/SmartMeshFoundation/Photon/transfer/mtree"
	"github.com/SmartMeshFoundation/Photon/utils"
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
func (rs *Service) restore() {
	//1. 处理未完成的锁
	// 1. handle incomplete locks
	rs.restoreLocks()
	//打印回复后的通道信息
	//log.Trace(fmt.Sprintf("tokengraph=%s", utils.StringInterface(rs.Token2ChannelGraph, 7)))
	//log.Trace(fmt.Sprintf("Transfer2StateManager=%s", utils.StringInterface(rs.Transfer2StateManager, 7)))
	//2. 为发送成功的 EnvelopMessage 继续发送
	// 2. keep sending EnvelopMessage that failed previously.
	rs.reSendEnvelopMessage()
}
func (rs *Service) reSendEnvelopMessage() {
	msgs := rs.dao.GetAllOrderedSentEnvelopMessager()
	for _, msg := range msgs {
		/*
			1. 可以立即发送消息,但是要在历史消息处理完毕以后再接收消息
			3. 已经过期的MediatedTransfer,对方还是会接受,但是需要我紧接着发送removeExpiredHashLock
			4. 通道已经 settle 的消息,会在protocol层被丢弃
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

func (rs *Service) restoreLocks() {
	token2ActionInitCrashRestartStateChange := make(map[common.Hash]*mediatedtransfer.ActionInitCrashRestartStateChange)
	var locks []*lockInfo
	//收集所有的锁,
	// collect all locks.
	//log.Trace(fmt.Sprintf("Token2TokenNetwork=%s", utils.StringInterface(rs.Token2TokenNetwork, 7)))
	//log.Trace(fmt.Sprintf("Token2ChannelGraph=%s", utils.StringInterface(rs.Token2ChannelGraph, 7)))
	for token := range rs.Token2TokenNetwork {
		g := rs.Token2ChannelGraph[token]
		//log.Trace(fmt.Sprintf("process token=%s", token.String()))
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
				//todo 密码已经链上注册的锁,需要跳过
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
				//todo 密码已经链上注册的锁,需要跳过
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
		//要注册密码,否则链上注册密码事件会找不到相关的通道.
		rs.registerChannelForHashlock(l.ch, l.l.LockSecretHash)
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
