package mobile

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"net/http"
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
	addr           common.Address        //投资人账户地址
	key            *ecdsa.PrivateKey     //投资人私钥
	contract       common.Address        //抵押合约地址
	c              *helper.SafeEthClient //公链rpc
	m              *snm.Mortgage         //抵押合约代理
	lockTime       *big.Int              //资金退出锁定时间
	endTimeOfFunds *big.Int              //募集资金结束时间
	isRunning      bool                  //正在运行?
	snmService     string                //超级节点地址,形如127.0.01:5003
}

//NewSNM 创建管理接口
func NewSNM(address, keystorePath, ethRPCEndPoint, password, contract, snmService string) (s *SNM, err error) {
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
		snmService:     snmService,
		c:              c,
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
	if !s.isRunning {
		return dto.NewErrorMobileResponse(errContractsAlreadyStopped)
	}
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
	//等到锁定期过了才能撤资
	if header.Number.Cmp(li.EndBlock) <= 0 {
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

/*
IncomeDetail 用户收益明细 ,it's reference of super-node-managermanaement
*/
type IncomeDetail struct {
	Key         string `json:"key" gorm:"primary_key"`
	Address     string `json:"address"`      // 获得该笔收入的用户
	Income      string `json:"income"`       // 本笔的收入金额
	BlockNumber int64  `json:"block_number"` // 该笔收入结算块号
	Date        string `json:"date"`         // 该笔收入结算日期,yyyy-mm-dd
	Timestamp   int64  `json:"timestamp"`    // 该笔收入结算时间戳
	TotalFund   string `json:"total_fund"`   // 该笔收入发生时的资金池总额
	Fund        string `json:"fund"`         // 该笔收入发生时的用户资金
	Proportion  string `json:"proportion"`   // 该笔收入发生时的用户资金占比
}

/*
GetMortgageInfoResponse :
*/
type getMortgageInfoResponse struct {
	/*
		超级节点部分信息
	*/
	SnmAccount            common.Address `json:"snm_account"`              // snm账户地址
	SnmTotalFund          *big.Int       `json:"snm_total_fund"`           // snm收到的投资(抵押)总额
	SnmAvailableTotalFund *big.Int       `json:"snm_available_total_fund"` // snm有效投资总额
	SnmBalance            *big.Int       `json:"snm_balance"`              // snm账户SMT余额
	SnmChannelNum         int            `json:"snm_channel_num"`          // snm在SMT的通道数量
	SnmAmountInChannel    *big.Int       `json:"snm_amount_in_channel"`    // snm通道中SMT总额
	SnmHistoryIncome      *big.Int       `json:"snm_history_income"`       // snm历史总收入

	/*
		用户部分信息
	*/
	UserFund                    *big.Int        `json:"user_fund,omitempty"`                      // 用户投资(抵押)总金额
	UserAvailableFund           *big.Int        `json:"user_available_fund,omitempty"`            // 用户有效资金总额(T+1)
	UserAvailableFundProportion string          `json:"user_available_fund_proportion,omitempty"` // 用户资金有效占比
	UserLockFund                *big.Int        `json:"user_lock_fund,omitempty"`                 // 用户锁定资金(正在撤出的资金)
	UserHistoryIncome           *big.Int        `json:"user_history_income,omitempty"`            // 用户历史收入总额
	UserIncomeDetailList        []*IncomeDetail `json:"user_income_detail_list,omitempty"`        // 用户历史收入明细
}

type status struct {
	IsRunning bool                    `json:"is_running"`
	SNM       getMortgageInfoResponse `json:"snm"`
}

//Status 查询投资状态
func (s *SNM) Status() (result string) {
	isRunning, err := s.m.IsRunning(nil)
	if err != nil {
		return dto.NewErrorMobileResponse(rerr.ErrContractQueryError.WithData(err))
	}
	st := status{}
	st.IsRunning = isRunning
	/*
	   访问snm获取信息
	*/
	FullURL := fmt.Sprintf("http://%s/api/1/mortgage-info/%s", s.snmService, s.addr.String())
	req := &utils.Req{
		FullURL: FullURL,
		Method:  http.MethodGet,
		Timeout: time.Second * 10,
	}
	statusCode, body, err := req.Invoke()
	if err != nil {
		err = fmt.Errorf("get snm info err :%s", err)
		return dto.NewErrorMobileResponse(rerr.ErrUnknown.WithData(err))
	}
	if statusCode != http.StatusOK {
		err = fmt.Errorf("get snm info  statusCode=%d", statusCode)
		return dto.NewErrorMobileResponse(rerr.ErrUnknown.WithData(err))
	}
	err = dto.ParseResult(string(body), &st.SNM)
	if err != nil {
		return dto.NewErrorMobileResponse(rerr.ErrUnknown.WithData(fmt.Sprintf("parse result err %s", err)))
	}
	return dto.NewSuccessMobileResponse(st)
}
