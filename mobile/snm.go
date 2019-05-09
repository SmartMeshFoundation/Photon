package mobile

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/SmartMeshFoundation/Photon/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/SmartMeshFoundation/Photon/dto"
	"github.com/SmartMeshFoundation/Photon/rerr"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/SmartMeshFoundation/Photon/network/helper"

	"github.com/SmartMeshFoundation/Photon/mobile/snm"

	"github.com/SmartMeshFoundation/Photon/accounts"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
)

var (
	// errInsufficientFundsForSub 撤出资金,金额过大
	errInsufficientFundsForSub        = rerr.NewError(10001, "insufficient value for subfunds")
	errAlreadyExistSubFunds           = rerr.NewError(10002, "ErrAlreadyExistSubFunds")
	errNotAllowSubFundsWhenFunding    = rerr.NewError(10003, "ErrNotAllowSubFundsWhenFunding")
	errValueCannotBeZero              = rerr.NewError(10004, "errValueCannotBeZero")
	errContractsAlreadyStopped        = rerr.NewError(10005, "errContractsAlreadyStopped")
	errSubFundsWithoutCallPreSubFunds = rerr.NewError(10006, "errSubFundsWithoutCallPreSubFunds")
	errSubFundsBeforeTimeout          = rerr.NewError(10007, "errSubFundsBeforeTimeout")
	errCannotStopContract             = rerr.NewError(10008, "errCannotStopContract")
	errCannotGetFundsNow              = rerr.NewError(10009, "errCannotGetFundsNow")
)

//SNM app与合约打交道
type SNM struct {
	addr           common.Address
	key            *ecdsa.PrivateKey
	contract       common.Address
	c              *helper.SafeEthClient
	m              *snm.Mortgage
	lockTime       *big.Int
	endTimeOfFunds *big.Int
	isRunning      bool
}

//NewSNM 创建管理接口
func NewSNM(address, keystorePath, ethRPCEndPoint, password, contract string) (s *SNM, err error) {
	addr := common.HexToAddress(address)
	_, keybin, err := accounts.PromptAccount(addr, keystorePath, password)
	if err != nil {
		return
	}
	key, err := crypto.ToECDSA(keybin)
	if err != nil {
		return
	}
	c, err := helper.NewSafeClient(ethRPCEndPoint)
	if err != nil {
		return
	}
	m, err := snm.NewMortgage(common.HexToAddress(contract), c)
	if err != nil {
		return
	}
	endTimeOfFunds, err := m.EndTimeOfFunds(nil)
	if err != nil {
		return
	}
	isRuning, err := m.IsRunning(nil)
	if err != nil {
		return
	}
	lockTime, err := m.LockTime(nil)
	if err != nil {
		return
	}
	return &SNM{
		addr:           addr,
		key:            key,
		contract:       common.HexToAddress(contract),
		m:              m,
		endTimeOfFunds: endTimeOfFunds,
		lockTime:       lockTime,
		isRunning:      isRuning,
	}, nil
}

//AddFunds 追加投资
func (s *SNM) AddFunds(value string) (result string) {
	if !s.isRunning {
		return dto.NewErrorMobileResponse(errContractsAlreadyStopped)
	}
	v, b := new(big.Int).SetString(value, 0)
	if !b {
		return dto.NewErrorMobileResponse(rerr.ErrArgumentError)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	av, err := s.c.BalanceAt(ctx, s.addr, nil)
	if err != nil {
		return dto.NewErrorMobileResponse(err)
	}
	if av.Cmp(v) < 0 {
		return dto.NewErrorMobileResponse(rerr.ErrInsufficientBalance)
	}
	opts := bind.NewKeyedTransactor(s.key)
	opts.Value = v
	tx, err := s.m.AddFunds(opts)
	return s.helpTx(tx, err)
}

//PreSubFunds 预备撤回投资
func (s *SNM) PreSubFunds(value string) (result string) {
	if !s.isRunning {
		return dto.NewErrorMobileResponse(errContractsAlreadyStopped)
	}
	v, b := new(big.Int).SetString(value, 0)
	if !b {
		return dto.NewErrorMobileResponse(rerr.ErrArgumentError)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	li, err := s.m.Locked(nil, s.addr)
	if err != nil {
		return dto.NewErrorMobileResponse(rerr.ErrContractQueryError.WithData(err))
	}
	if li.Value.Cmp(utils.BigInt0) > 0 {
		return dto.NewErrorMobileResponse(errAlreadyExistSubFunds)
	}
	header, err := s.c.HeaderByNumber(ctx, nil)
	if err != nil {
		return dto.NewErrorMobileResponse(rerr.ErrContractQueryError.WithData(err))
	}
	if header.Number.Cmp(s.endTimeOfFunds) <= 0 {
		return dto.NewErrorMobileResponse(errNotAllowSubFundsWhenFunding)
	}
	mortgage, err := s.m.Mortgage(nil, s.addr)
	if err != nil {
		return dto.NewErrorMobileResponse(rerr.ErrContractQueryError.WithData(err))
	}
	if mortgage.Cmp(v) < 0 {
		return dto.NewErrorMobileResponse(rerr.ErrInsufficientBalance)
	}
	opts := bind.NewKeyedTransactor(s.key)

	tx, err := s.m.PreSubFunds(opts, v)
	return s.helpTx(tx, err)
}

//SubFunds 锁定到期,撤回投资
func (s *SNM) SubFunds() (result string) {

	li, err := s.m.Locked(nil, s.addr)
	if err != nil {
		return dto.NewErrorMobileResponse(rerr.ErrContractQueryError)
	}
	if li.Value.Cmp(utils.BigInt0) <= 0 {
		return dto.NewErrorMobileResponse(errSubFundsWithoutCallPreSubFunds)
	}
	header, err := s.c.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return dto.NewErrorMobileResponse(rerr.ErrContractQueryError.WithData(err))
	}
	isRunning, err := s.m.IsRunning(nil)
	if err != nil {
		return dto.NewErrorMobileResponse(rerr.ErrContractQueryError)
	}
	//如果合约已经停止运行,随时都可以撤回资金
	if header.Number.Cmp(li.EndBlock) <= 0 && isRunning {
		return dto.NewErrorMobileResponse(errNotAllowSubFundsWhenFunding)
	}
	opts := bind.NewKeyedTransactor(s.key)

	tx, err := s.m.SubFunds(opts)
	return s.helpTx(tx, err)
}

//TryStopContract 当募集资金结合后,募集自己不够,任何人可以停止合约
func (s *SNM) TryStopContract() (result string) {
	header, err := s.c.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return dto.NewErrorMobileResponse(rerr.ErrContractQueryError.WithData(err))
	}
	if header.Number.Cmp(s.endTimeOfFunds) <= 0 {
		return dto.NewErrorMobileResponse(errCannotStopContract)
	}
	balance, err := s.c.BalanceAt(context.Background(), s.contract, nil)
	if err != nil {
		return dto.NewErrorMobileResponse(rerr.ErrContractQueryError)
	}
	minimumFunds, err := s.m.MinimumFunds(nil)
	if err != nil {
		return dto.NewErrorMobileResponse(rerr.ErrContractQueryError)
	}
	if minimumFunds.Cmp(balance) <= 0 {
		return dto.NewErrorMobileResponse(errCannotStopContract)
	}
	isRunning, err := s.m.IsRunning(nil)
	if err != nil {
		return dto.NewErrorMobileResponse(rerr.ErrContractQueryError)
	}
	if !isRunning {
		return dto.NewErrorMobileResponse(errContractsAlreadyStopped)
	}
	opts := bind.NewKeyedTransactor(s.key)
	tx, err := s.m.TryStopContract(opts)
	return s.helpTx(tx, err)
}

//GetFunds 当合约停止运行以后,投资人可以立即撤回投资
func (s *SNM) GetFunds() (result string) {
	isRunning, err := s.m.IsRunning(nil)
	if err != nil {
		return dto.NewErrorMobileResponse(rerr.ErrContractQueryError)
	}
	if isRunning {
		return dto.NewErrorMobileResponse(errCannotGetFundsNow)
	}
	opts := bind.NewKeyedTransactor(s.key)
	tx, err := s.m.GetFunds(opts)
	return s.helpTx(tx, err)
}
func (s *SNM) helpTx(tx *types.Transaction, err error) string {
	if err != nil {
		return dto.NewErrorMobileResponse(rerr.ErrTxWaitMined.WithData("failed tx"))
	}
	r, err := bind.WaitMined(context.Background(), s.c, tx)
	if err != nil {
		return dto.NewErrorMobileResponse(rerr.ErrTxWaitMined)
	}
	if r.Status != types.ReceiptStatusSuccessful {
		return dto.NewErrorMobileResponse(rerr.ErrTxReceiptStatus)
	}
	return dto.NewSuccessMobileResponse("")
}

type status struct {
	IsRunning     bool
	TatalInterest *big.Int //可以从合约走,但是每天利息收益呢?
	LockedValue   *big.Int
	LockEndBlock  *big.Int
	//正在进行中的投资,只能通过超节点查询
}

//Status 查询投资状态
func (s *SNM) Status() (result string) {
	isRunning, err := s.m.IsRunning(nil)
	if err != nil {
		return dto.NewErrorMobileResponse(rerr.ErrContractQueryError.WithData(err))
	}
	st := status{}
	st.IsRunning = isRunning
	return dto.NewSuccessMobileResponse(st)
}
