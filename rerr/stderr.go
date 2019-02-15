package rerr

import (
	"encoding/json"
	"fmt"
)

//StandardError 标准错误，包含错误码和错误信息
type StandardError struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_message"`
}

// Error 实现了 Error接口
func (e StandardError) Error() string {
	return fmt.Sprintf("errorCode: %d, errorMsg %s", e.ErrorCode, e.ErrorMsg)
}

//Printf 附加错误描述信息
func (e StandardError) Printf(format string, a ...interface{}) StandardError {
	s := fmt.Sprintf(format, a...)
	err2 := e
	err2.ErrorMsg = fmt.Sprintf("%s:%s", e.ErrorMsg, s)
	return err2
}

//Errorf alias of printf
func (e StandardError) Errorf(format string, a ...interface{}) StandardError {
	return e.Printf(format, a...)
}

//Append 附加错误描述信息
func (e StandardError) Append(info string) StandardError {
	err2 := e
	err2.ErrorMsg = fmt.Sprintf("%s:%s", e.ErrorMsg, info)
	return err2
}

//AppendError 附加错误信息
func (e StandardError) AppendError(err error) StandardError {
	if err != nil {
		err2 := e
		err2.ErrorMsg = fmt.Sprintf("%s:%s", e.ErrorMsg, err.Error())
		return err2
	}
	return e
}

//StandardDataError 用于有结构化错误描述的场景
type StandardDataError struct {
	StandardError
	Data json.RawMessage `json:"data"`
}

//WithData 附加结构化错误信息
func (e StandardError) WithData(data interface{}) StandardDataError {
	err2 := StandardDataError{
		StandardError: e,
	}
	d, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	err2.Data = json.RawMessage(d)
	return err2
}

//ContractCallError 将合约调用错误转换为standarderror
func ContractCallError(err error) StandardError {
	//金额不够需要特殊处理
	if err.Error() == "insufficient balance to pay for gas" || err.Error() == "insufficient funds for gas * price + value" {
		return ErrInsufficientBalanceForGas

	}
	return ErrUnkownSpectrumRPCError.AppendError(err)
}

var (
	//ErrSuccess 成功
	ErrSuccess = StandardError{ErrorCode: 0, ErrorMsg: "success"}
	//ErrUnrecognized 未知错误,
	ErrUnrecognized = StandardError{ErrorCode: -1, ErrorMsg: "unknown error"}
	//ErrArgumentError 参数错误
	ErrArgumentError = StandardError{ErrorCode: 1, ErrorMsg: "ArgumentError"}
	//ErrPhotonAlreadyRunning 启动了多个photon实例
	ErrPhotonAlreadyRunning = StandardError{ErrorCode: 2, ErrorMsg: "PhotonAlreadyRunning"}
	//ErrHashLengthNot32  参数错误
	ErrHashLengthNot32 = StandardError{ErrorCode: 1000, ErrorMsg: "HashLengthNot32"}
	/*ErrNotFound Raised when something not found
	 */
	ErrNotFound = StandardError{ErrorCode: 1001, ErrorMsg: "Not found"}
	/*ErrInsufficientBalance Raised when the netting channel doesn't enough available capacity to
	  pay for the transfer.

	  Used for the validation of an *incoming* messages.
	*/
	ErrInsufficientBalance = StandardError{ErrorCode: 1002, ErrorMsg: "InsufficientBalance"}
	/*ErrInvalidAmount Raised when the user provided value is not a integer and cannot be used
	  to defined a transfer value
	*/
	ErrInvalidAmount = StandardError{ErrorCode: 1003, ErrorMsg: "InvalidAmount"}

	/*ErrNoPathError Raised when there is no path to the requested target address in the
	  payment network.

	  This exception is raised if there is not a single path in the network to
	  reach the target, it's not used if there is a path but the transfre failed
	  because of the lack of capacity or network problems.
	*/
	ErrNoPathError = StandardError{ErrorCode: 1005, ErrorMsg: "NoPathError"}
	/*ErrSamePeerAddress Raised when a user tries to create a channel where the address of both
	  peers is the same.
	*/
	ErrSamePeerAddress = StandardError{ErrorCode: 1006, ErrorMsg: "SamePeerAddress"}
	/*ErrInvalidState Raised when the user requested action cannot be done due to the current
	  state of the channel.
	*/
	ErrInvalidState = StandardError{ErrorCode: 1007, ErrorMsg: "InvalidState"}
	//ErrTransferWhenClosed Raised when a user tries to request a transfer is a closed channel.
	ErrTransferWhenClosed = StandardError{ErrorCode: 1008, ErrorMsg: "TransferWhenClosed"}
	/*ErrUnknownAddress Raised when the user provided address is valid but is not from a known
	  node.
	*/
	ErrUnknownAddress = StandardError{ErrorCode: 1009, ErrorMsg: "UnknownAddress"}
	/*ErrInvalidLocksRoot Raised when the received message has an invalid locksroot.

	  Used to reject a message when a pending lock is missing from the locksroot,
	  otherwise if the message is accepted there is a pontential loss of token.
	*/
	ErrInvalidLocksRoot = StandardError{ErrorCode: 1010, ErrorMsg: "Locksroot mismatch"}
	/*ErrInvalidNonce Raised when the received messages has an invalid value for the nonce.

	  The nonce field must change incrementally
	*/
	ErrInvalidNonce = StandardError{ErrorCode: 1011, ErrorMsg: "InvalidNonce"}

	/*ErrTransferUnwanted Raised when the node is not receiving new transfers.
	 */
	ErrTransferUnwanted = StandardError{ErrorCode: 1012, ErrorMsg: "TransferUnwanted"}
	// ErrStopCreateNewTransfer reject new transactions
	ErrStopCreateNewTransfer = StandardError{ErrorCode: 1013, ErrorMsg: "new transactions are not allowed"}

	//ErrNotAllowMediatedTransfer not allow mediated transfer when mesh
	ErrNotAllowMediatedTransfer = StandardError{ErrorCode: 1014, ErrorMsg: "no mediated transfer on mesh only network"}
	//ErrDuplicateTransfer token和secret相同的交易
	ErrDuplicateTransfer = StandardError{ErrorCode: 1015, ErrorMsg: "secret and token cannot duplicate"}
	//ErrNodeNotOnline 发送消息时,对方不在线
	ErrNodeNotOnline = StandardError{ErrorCode: 1016, ErrorMsg: "NodeOffline"}
	//ErrTransferCannotCancel 试图取消已经泄露秘密的交易
	ErrTransferCannotCancel = StandardError{ErrorCode: 1017, ErrorMsg: "TranasferCannotCancel"}
	/*
		DB error
	*/

	//ErrGeneralDBError 未归类数据库错误,需要进一步细化
	ErrGeneralDBError = StandardError{ErrorCode: 1018, ErrorMsg: "DBError"}
	//ErrDBDuplicateKey 重复的key
	ErrDBDuplicateKey = StandardError{ErrorCode: 1019, ErrorMsg: "duplicate key"}
	//ErrTransferTimeout 交易超时,不代表交易肯定会成功或者失败,只是在给定时间内交易没有成功而已
	ErrTransferTimeout = StandardError{ErrorCode: 1020, ErrorMsg: "ErrTransferTimeout"}
	//ErrUpdateButHaveTransfer 试图升级,发现还有交易在进行
	ErrUpdateButHaveTransfer = StandardError{ErrorCode: 1021, ErrorMsg: "ErrUpdateButHaveTransfer"}
	//ErrNotChargeFee 进行与收费相关的操作,但是没有启用收费
	ErrNotChargeFee = StandardError{ErrorCode: 1022, ErrorMsg: "ErrNotChargeFee"}
	/*
		以太坊报公链节点报的错误


	*/

	//ErrInsufficientBalanceForGas gas problem
	ErrInsufficientBalanceForGas = StandardError{ErrorCode: 2000, ErrorMsg: "insufficient balance to pay for gas"}

	/*
		Tx 相关
		链上操作相关
	*/

	//ErrCloseChannel 链上执行关闭通道时发生了错误
	ErrCloseChannel = StandardError{ErrorCode: 2001, ErrorMsg: "closeChannel"}
	//ErrRegisterSecret 链上注册密码的时候发生了错误
	ErrRegisterSecret = StandardError{ErrorCode: 2002, ErrorMsg: "RegisterSecret"}
	//ErrUnlock 链上unlock的时候发生了错误
	ErrUnlock = StandardError{ErrorCode: 2003, ErrorMsg: "Unlock"}
	//ErrUpdateBalanceProof 链上提交balance proof发生错误
	ErrUpdateBalanceProof = StandardError{ErrorCode: 2004, ErrorMsg: "UpdateBalanceProof"}
	//ErrPunish 链上执行punish的时候发生错误
	ErrPunish = StandardError{ErrorCode: 2005, ErrorMsg: "punish"}
	//ErrSettle 链上执行settle操作的时候发生错误
	ErrSettle = StandardError{ErrorCode: 2006, ErrorMsg: "settle"}
	//ErrDeposit 链上执行deposit发生错误
	ErrDeposit = StandardError{ErrorCode: 2007, ErrorMsg: "deposit"}
	//ErrSpectrumNotConnected 没有连接到公链.
	ErrSpectrumNotConnected = StandardError{ErrorCode: 2008, ErrorMsg: "ErrSpectrumNotConnected"}
	//ErrTxWaitMined waitMined return error
	ErrTxWaitMined = StandardError{ErrorCode: 2008, ErrorMsg: "ErrTxWaitMined"}
	//ErrTxReceiptStatus tx 被打包了,但是结果失败
	ErrTxReceiptStatus = StandardError{ErrorCode: 2009, ErrorMsg: "ErrTxReceiptStatus"}
	//ErrSecretAlreadyRegistered 尝试连上注册密码,但是密码已经注册了
	ErrSecretAlreadyRegistered = StandardError{ErrorCode: 2010, ErrorMsg: "ErrSecretAlreadyRegistered"}
	//ErrSpectrumSyncError 连接到的公链长时间不出块或者正在同步
	ErrSpectrumSyncError = StandardError{ErrorCode: 2011, ErrorMsg: "ErrSpectrumSyncError"}
	//ErrSpectrumBlockError 本地已处理的块数和公链汇报块数不一致,比如我本地已经处理到了50000块,但是公链节点报告现在只有3000块
	ErrSpectrumBlockError = StandardError{ErrorCode: 2012, ErrorMsg: "ErrSpectrumBlockError"}
	//ErrUnkownSpectrumRPCError 其他以太坊rpc错误
	ErrUnkownSpectrumRPCError = StandardError{ErrorCode: 2999, ErrorMsg: "unkown spectrum rpc error"}
	/*ErrTokenNotFound Raised when token not found
	 */
	ErrTokenNotFound = StandardError{ErrorCode: 3001, ErrorMsg: "TokenNotFound"}
	/*ErrChannelNotFound Raised when token not found
	 */
	ErrChannelNotFound = StandardError{ErrorCode: 3002, ErrorMsg: "ChannelNotFound"}
	//ErrNoAvailabeRoute no availabe route
	ErrNoAvailabeRoute = StandardError{ErrorCode: 3003, ErrorMsg: "NoAvailabeRoute"}
	//ErrTransferNotFound not found transfer
	ErrTransferNotFound = StandardError{ErrorCode: 3004, ErrorMsg: "TransferNotFound"}
	//ErrChannelAlreadExist 通道已存在
	ErrChannelAlreadExist = StandardError{ErrorCode: 3005, ErrorMsg: "ChannelAlreadExist"}
	/*ErrPFS PFS Error
	向PFS发起请求错误
	*/
	ErrPFS = StandardError{ErrorCode: 4000, ErrorMsg: "ErrorPFS"}

	/*
		Channel Error
	*/

	//ErrChannelNotAllowWithdraw 通道现在不能合作取现,比如有交易在进行
	ErrChannelNotAllowWithdraw = StandardError{ErrorCode: 5000, ErrorMsg: "CannotWithdarw"}
	//ErrChannelState 在不能执行相应操作的通道状态,试图执行某些交易,比如在关闭的通道上发起交易
	ErrChannelState = StandardError{ErrorCode: 5001, ErrorMsg: "ErrChannelState"}
	//ErrChannelSettleTimeout 没到settle时间尝试去settle
	ErrChannelSettleTimeout = StandardError{ErrorCode: 5002, ErrorMsg: "Channel only can settle after timeout"}
	//ErrChannelNotParticipant 给定地址不是通道的任何参与一方
	ErrChannelNotParticipant = StandardError{ErrorCode: 5003, ErrorMsg: "NotParticipant"}
	//ErrChannelLockSecretHashNotFound 通道中没有相应的锁
	ErrChannelLockSecretHashNotFound = StandardError{ErrorCode: 5004, ErrorMsg: "ChannelNoSuchLock"}
	//ErrChannelEndStateNoSuchLock 通道当前参与方中找不到相应的锁
	ErrChannelEndStateNoSuchLock = StandardError{ErrorCode: 5005, ErrorMsg: "ErrChannelEndStateNoSuchLock"}
	//ErrChannelLockAlreadyExpired 通道中锁已过期
	ErrChannelLockAlreadyExpired = StandardError{ErrorCode: 5006, ErrorMsg: "ErrChannelLockAlreadyExpired"}
	//ErrChannelBalanceDecrease 发生了降低通道balance(指的是合约中的balance)的行为
	ErrChannelBalanceDecrease = StandardError{ErrorCode: 5007, ErrorMsg: "ErrChannelBalanceDecrease"}
	//ErrChannelTransferAmountMismatch 收到的交易中transferamount不匹配
	ErrChannelTransferAmountMismatch = StandardError{ErrorCode: 5008, ErrorMsg: "ErrChannelTransferAmountMismatch"}
	//ErrChannelBalanceProofAlreadyRegisteredOnChain  已经提交过balanceproof以后试图修改本地balance proof
	ErrChannelBalanceProofAlreadyRegisteredOnChain = StandardError{ErrorCode: 5009, ErrorMsg: "ErrChannelBalanceProofAlreadyRegisteredOnChain"}
	//ErrChannelDuplicateLock 通道中已存在该密码的锁
	ErrChannelDuplicateLock = StandardError{ErrorCode: 5010, ErrorMsg: "ErrChannelDuplicateLock"}
	//ErrChannelTransferAmountDecrease 收到交易,TransferAmount变小了
	ErrChannelTransferAmountDecrease = StandardError{ErrorCode: 5011, ErrorMsg: "ErrChannelTransferAmountDecrease"}
	//ErrRemoveNotExpiredLock 试图移除没有过期的锁
	ErrRemoveNotExpiredLock = StandardError{ErrorCode: 5012, ErrorMsg: "ErrRemoveNotExpiredLock"}
	//ErrUpdateBalanceProofAfterClosed 试图在通道关闭以后还更新对方或者我的balance proof,基本意思和ErrChannelBalanceProofAlreadyRegisteredOnChain一样
	ErrUpdateBalanceProofAfterClosed = StandardError{ErrorCode: 5013, ErrorMsg: "ErrUpdateBalanceProofAfterClosed"}
	//ErrChannelIdentifierMismatch 通道id不匹配
	ErrChannelIdentifierMismatch = StandardError{ErrorCode: 5014, ErrorMsg: "ErrChannelIdentifierMismatch"}
	//ErrChannelInvalidSender 收到来自未知参与方的交易
	ErrChannelInvalidSender = StandardError{ErrorCode: 5015, ErrorMsg: "ErrChannelInvalidSender"}
	//ErrChannelBalanceNotMatch  合作关闭通道,取现时金额检查不匹配,
	ErrChannelBalanceNotMatch = StandardError{ErrorCode: 5016, ErrorMsg: "ErrChannelBalanceNotMatch"}
	//ErrChannelLockMisMatch 收到交易中指定的锁和本地不匹配
	ErrChannelLockMisMatch = StandardError{ErrorCode: 5017, ErrorMsg: "ErrChannelLockMisMatch"}
	//ErrChannelWithdrawAmount  合作取现的金额过大
	ErrChannelWithdrawAmount = StandardError{ErrorCode: 5018, ErrorMsg: "ErrChannelWithdrawAmount"}
	//ErrChannelLockExpirationTooLarge 收到交易,指定的过期时间太长了,这可能是一种攻击
	ErrChannelLockExpirationTooLarge = StandardError{ErrorCode: 5019, ErrorMsg: "ErrChannelLockExpirationTooLarge"}
	//ErrChannelRevealTimeout 指定的reveal timeout 非法
	ErrChannelRevealTimeout = StandardError{ErrorCode: 5020, ErrorMsg: "ErrChannelRevealTimeout"}
	//ErrChannelBalanceProofNil balanceproof为空
	ErrChannelBalanceProofNil = StandardError{ErrorCode: 5021, ErrorMsg: "ErrChannelBalanceProofNil"}
	//ErrChannelCloseClosedChannel 试图关闭已经关闭的通道
	ErrChannelCloseClosedChannel = StandardError{ErrorCode: 5022, ErrorMsg: "ErrChannelCloseClosedChannel"}
	//ErrChannelBackgroundTx 后台执行Tx发生错误
	ErrChannelBackgroundTx = StandardError{ErrorCode: 5023, ErrorMsg: "ErrChannelBackgroundTx"}

	/*ErrChannelWithdrawButHasLocks : we can't send a request for withdraw when there are locks.
	 */
	ErrChannelWithdrawButHasLocks = StandardError{ErrorCode: 5014, ErrorMsg: "ErrChannelWithdrawButHasLocks"}
	/*ErrChannelCooperativeSettleButHasLocks : we can't send a request for settle when there are locks.
	 */
	ErrChannelCooperativeSettleButHasLocks = StandardError{ErrorCode: 5015, ErrorMsg: "ErrChannelCooperativeSettleButHasLocks"}
	/*ErrChannelInvalidSttleTimeout Raised when the user provided timeout value is less than the minimum
	  settle timeout
	*/
	ErrChannelInvalidSttleTimeout = StandardError{ErrorCode: 5003, ErrorMsg: "ErrInvalidSettleTimeout"}
	/*
		Transport error
	*/

	//ErrTransportTypeUnknown  未知的transport层错误,
	ErrTransportTypeUnknown = StandardError{ErrorCode: 6000, ErrorMsg: "transport type error"}
	//ErrSubScribeNeighbor 订阅节点在线信息错误
	ErrSubScribeNeighbor = StandardError{ErrorCode: 6001, ErrorMsg: "ErrSubScribeNeighbor"}
)
