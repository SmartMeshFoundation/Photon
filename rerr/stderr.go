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

func (e StandardError) AppendError(err error) StandardError {
	if err != nil {
		err2 := e
		err2.ErrorMsg = fmt.Sprintf("%s:%s", e.ErrorMsg, err.Error())
		return err2
	}
	return e
}

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
	ErrSuccess              = StandardError{ErrorCode: 0, ErrorMsg: "success"}
	ErrUnrecognized         = StandardError{ErrorCode: -1, ErrorMsg: "unknown error"}
	ErrArgumentError        = StandardError{ErrorCode: 1, ErrorMsg: "ArgumentError"}
	ErrPhotonAlreadyRunning = StandardError{ErrorCode: 2, ErrorMsg: "PhotonAlreadyRunning"}
	ErrHashLengthNot32      = StandardError{ErrorCode: 1000, ErrorMsg: "HashLengthNot32"}
	/*
		ErrNotFound Raised when something not found
	*/
	ErrNotFound = StandardError{ErrorCode: 1001, ErrorMsg: "Not found"}
	/*
		ErrInsufficientBalance Raised when the netting channel doesn't enough available capacity to
		    pay for the transfer.

		    Used for the validation of an *incoming* messages.
	*/
	ErrInsufficientBalance = StandardError{ErrorCode: 1002, ErrorMsg: "InsufficientBalance"}
	/*
			ErrInvalidAmount Raised when the user provided value is not a integer and cannot be used
		    to defined a transfer value
	*/
	ErrInvalidAmount = StandardError{ErrorCode: 1003, ErrorMsg: "InvalidAmount"}

	//todo 添一个...
	/*
	   ErrNoPathError Raised when there is no path to the requested target address in the
	   payment network.

	   This exception is raised if there is not a single path in the network to
	   reach the target, it's not used if there is a path but the transfre failed
	   because of the lack of capacity or network problems.
	*/
	ErrNoPathError = StandardError{ErrorCode: 1005, ErrorMsg: "NoPathError"}
	/*
	   ErrSamePeerAddress Raised when a user tries to create a channel where the address of both
	       peers is the same.
	*/
	ErrSamePeerAddress = StandardError{ErrorCode: 1006, ErrorMsg: "SamePeerAddress"}
	/*
	   InvalidState Raised when the user requested action cannot be done due to the current
	   state of the channel.
	*/
	ErrInvalidState = StandardError{ErrorCode: 1007, ErrorMsg: "InvalidState"}
	//TransferWhenClosed Raised when a user tries to request a transfer is a closed channel.
	ErrTransferWhenClosed = StandardError{ErrorCode: 1008, ErrorMsg: "TransferWhenClosed"}
	/*
		UnknownAddress Raised when the user provided address is valid but is not from a known
		    node.
	*/
	ErrUnknownAddress = StandardError{ErrorCode: 1009, ErrorMsg: "UnknownAddress"}
	/*
		ErrInvalidLocksRoot Raised when the received message has an invalid locksroot.

		    Used to reject a message when a pending lock is missing from the locksroot,
		    otherwise if the message is accepted there is a pontential loss of token.
	*/
	ErrInvalidLocksRoot = StandardError{ErrorCode: 1010, ErrorMsg: "Locksroot mismatch"}
	/*
		ErrInvalidNonce Raised when the received messages has an invalid value for the nonce.

		    The nonce field must change incrementally
	*/
	ErrInvalidNonce = StandardError{ErrorCode: 1011, ErrorMsg: "InvalidNonce"}

	/*
	   ErrTransferUnwanted Raised when the node is not receiving new transfers.
	*/
	ErrTransferUnwanted = StandardError{ErrorCode: 1012, ErrorMsg: "TransferUnwanted"}
	// ErrStopCreateNewTransfer reject new transactions
	ErrStopCreateNewTransfer = StandardError{ErrorCode: 1013, ErrorMsg: "new transactions are not allowed"}

	//ErrNotAllowMediatedTransfer not allow mediated transfer when mesh
	ErrNotAllowMediatedTransfer = StandardError{ErrorCode: 1014, ErrorMsg: "no mediated transfer on mesh only network"}
	//ErrDuplicateTransfer token和secret相同的交易
	ErrDuplicateTranser     = StandardError{ErrorCode: 1015, ErrorMsg: "secret and token cannot duplicate"}
	ErrNodeNotOnline        = StandardError{ErrorCode: 1016, ErrorMsg: "NodeOffline"}
	ErrTransferCannotCancel = StandardError{ErrorCode: 1017, ErrorMsg: "TranasferCannotCancel"}
	/*
		DB error
	*/
	ErrGeneralDBError  = StandardError{ErrorCode: 1018, ErrorMsg: "DBError"}
	ErrDBDuplicateKey  = StandardError{ErrorCode: 1019, ErrorMsg: "duplicate key"}
	ErrTransferTimeout = StandardError{ErrorCode: 1020, ErrorMsg: "ErrTransferTimeout"}
	/*
		以太坊报公链节点报的错误



	*/
	//ErrInsufficientBalanceForGas gas problem
	ErrInsufficientBalanceForGas = StandardError{ErrorCode: 2000, ErrorMsg: "insufficient balance to pay for gas"}

	/*
		Tx 相关
		链上操作相关
	*/
	ErrCloseChannel         = StandardError{ErrorCode: 2001, ErrorMsg: "closeChannel"}
	ErrRegisterSecret       = StandardError{ErrorCode: 2002, ErrorMsg: "RegisterSecret"}
	ErrUnlock               = StandardError{ErrorCode: 2003, ErrorMsg: "Unlock"}
	ErrUpdateBalanceProof   = StandardError{ErrorCode: 2004, ErrorMsg: "UpdateBalanceProof"}
	ErrPunish               = StandardError{ErrorCode: 2005, ErrorMsg: "punish"}
	ErrSettle               = StandardError{ErrorCode: 2006, ErrorMsg: "settle"}
	ErrDeposit              = StandardError{ErrorCode: 2007, ErrorMsg: "deposit"}
	ErrSpectrumNotConnected = StandardError{ErrorCode: 2008, ErrorMsg: "ErrSpectrumNotConnected"}
	//ErrTxWaitMined waitMined return error
	ErrTxWaitMined = StandardError{ErrorCode: 2008, ErrorMsg: "ErrTxWaitMined"}
	//ErrTxReceiptStatus tx 被打包了,但是结果失败
	ErrTxReceiptStatus         = StandardError{ErrorCode: 2009, ErrorMsg: "ErrTxReceiptStatus"}
	ErrSecretAlreadyRegistered = StandardError{ErrorCode: 2010, ErrorMsg: "ErrSecretAlreadyRegistered"}
	ErrSpectrumSyncError       = StandardError{ErrorCode: 2011, ErrorMsg: "ErrSpectrumSyncError"}
	ErrSpectrumBlockError      = StandardError{ErrorCode: 2012, ErrorMsg: "ErrSpectrumBlockError"}
	//ErrUnkownSpectrumRPCError 其他以太坊rpc错误
	ErrUnkownSpectrumRPCError = StandardError{ErrorCode: 2999, ErrorMsg: "unkown spectrum rpc error"}
	/*
		ErrTokenNotFound Raised when token not found
	*/
	ErrTokenNotFound = StandardError{ErrorCode: 3001, ErrorMsg: "TokenNotFound"}
	/*
		ErrChannelNotFound Raised when token not found
	*/
	ErrChannelNotFound = StandardError{ErrorCode: 3002, ErrorMsg: "ChannelNotFound"}
	//ErrNoAvailabeRoute no availabe route
	ErrNoAvailabeRoute = StandardError{ErrorCode: 3003, ErrorMsg: "NoAvailabeRoute"}
	//ErrTransferNotFound not found transfer
	ErrTransferNotFound   = StandardError{ErrorCode: 3004, ErrorMsg: "TransferNotFound"}
	ErrChannelAlreadExist = StandardError{ErrorCode: 3005, ErrorMsg: "ChannelAlreadExist"}
	/*
		PFS Error
		向PFS发起请求错误
	*/
	ErrPFS = StandardError{ErrorCode: 4000, ErrorMsg: "ErrorPFS"}

	/*
		Channel Error
	*/
	ErrChannelNotAllowWithdraw                     = StandardError{ErrorCode: 5000, ErrorMsg: "CannotWithdarw"}
	ErrChannelState                                = StandardError{ErrorCode: 5001, ErrorMsg: "ErrChannelState"}
	ErrChannelSettleTimeout                        = StandardError{ErrorCode: 5002, ErrorMsg: "Channel only can settle after timeout"}
	ErrChannelNotParticipant                       = StandardError{ErrorCode: 5003, ErrorMsg: "NotParticipant"}
	ErrChannelLockSecretHashNotFound               = StandardError{ErrorCode: 5004, ErrorMsg: "ChannelNoSuchLock"}
	ErrChannelEndStateNoSuchLock                   = StandardError{ErrorCode: 5005, ErrorMsg: "ErrChannelEndStateNoSuchLock"}
	ErrChannelLockAlreadyExpired                   = StandardError{ErrorCode: 5006, ErrorMsg: "ErrChannelLockAlreadyExpired"}
	ErrChannelBalanceDecrease                      = StandardError{ErrorCode: 5007, ErrorMsg: "ErrChannelBalanceDecrease"}
	ErrChannelTransferAmountMismatch               = StandardError{ErrorCode: 5008, ErrorMsg: "ErrChannelTransferAmountMismatch"}
	ErrChannelBalanceProofAlreadyRegisteredOnChain = StandardError{ErrorCode: 5009, ErrorMsg: "ErrChannelBalanceProofAlreadyRegisteredOnChain"}
	ErrChannelDuplicateLock                        = StandardError{ErrorCode: 5010, ErrorMsg: "ErrChannelDuplicateLock"}
	ErrChannelTransferAmountDecrease               = StandardError{ErrorCode: 5011, ErrorMsg: "ErrChannelTransferAmountDecrease"}
	ErrRemoveNotExpiredLock                        = StandardError{ErrorCode: 5012, ErrorMsg: "ErrRemoveNotExpiredLock"}
	ErrUpdateBalanceProofAfterClosed               = StandardError{ErrorCode: 5013, ErrorMsg: "ErrUpdateBalanceProofAfterClosed"}
	ErrChannelIdentifierMismatch                   = StandardError{ErrorCode: 5014, ErrorMsg: "ErrChannelIdentifierMismatch"}
	ErrChannelInvalidSender                        = StandardError{ErrorCode: 5015, ErrorMsg: "ErrChannelInvalidSender"}
	ErrChannelBalanceNotMatch                      = StandardError{ErrorCode: 5016, ErrorMsg: "ErrChannelBalanceNotMatch"}
	ErrChannelLockMisMatch                         = StandardError{ErrorCode: 5017, ErrorMsg: "ErrChannelLockMisMatch"}
	ErrChannelWithdrawAmount                       = StandardError{ErrorCode: 5018, ErrorMsg: "ErrChannelWithdrawAmount"}
	ErrChannelLockExpirationTooLarge               = StandardError{ErrorCode: 5019, ErrorMsg: "ErrChannelLockExpirationTooLarge"}
	ErrChannelRevealTimeout                        = StandardError{ErrorCode: 5020, ErrorMsg: "ErrChannelRevealTimeout"}
	ErrChannelBalanceProofNil                      = StandardError{ErrorCode: 5021, ErrorMsg: "ErrChannelBalanceProofNil"}
	ErrChannelCloseClosedChannel                   = StandardError{ErrorCode: 5022, ErrorMsg: "ErrChannelCloseClosedChannel"}
	ErrChannelBackgroundTx                         = StandardError{ErrorCode: 5023, ErrorMsg: "ErrChannelBackgroundTx"}

	/*
	 *	ErrChannelWithdrawButHasLocks : we can't send a request for withdraw when there are locks.
	 */
	ErrChannelWithdrawButHasLocks = StandardError{ErrorCode: 5014, ErrorMsg: "ErrChannelWithdrawButHasLocks"}
	/*
	 *	ErrChannelCooperativeSettleButHasLocks : we can't send a request for settle when there are locks.
	 */
	ErrChannelCooperativeSettleButHasLocks = StandardError{ErrorCode: 5015, ErrorMsg: "ErrChannelCooperativeSettleButHasLocks"}
	/*
	   ErrChannelInvalidSttleTimeout Raised when the user provided timeout value is less than the minimum
	   settle timeout
	*/
	ErrChannelInvalidSttleTimeout = StandardError{ErrorCode: 5003, ErrorMsg: "ErrInvalidSettleTimeout"}
	/*
		Transport error
	*/
	ErrTransportTypeUnknown = StandardError{ErrorCode: 6000, ErrorMsg: "transport type error"}
	ErrSubScribeNeighbor    = StandardError{ErrorCode: 6001, ErrorMsg: "ErrSubScribeNeighbor"}
)
