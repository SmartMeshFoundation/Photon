package photon

import (
	"math/big"

	"sync"

	"errors"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/pfsproxy"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

//NoFeePolicy charge no fee
type NoFeePolicy struct {
}

//GetNodeChargeFee always return 0
func (n *NoFeePolicy) GetNodeChargeFee(nodeAddress, tokenAddress common.Address, amount *big.Int) *big.Int {
	return utils.BigInt0
}

// FeeModule :
type FeeModule struct {
	db        *models.ModelDB
	pfsProxy  pfsproxy.PfsProxy
	feePolicy *models.FeePolicy
	lock      sync.Mutex
}

// NewFeeModule :
func NewFeeModule(db *models.ModelDB, pfsProxy pfsproxy.PfsProxy) (fm *FeeModule, err error) {
	if db == nil {
		panic("need init db first")
	}
	fp := db.GetFeePolicy()
	fm = &FeeModule{
		db:       db,
		pfsProxy: pfsProxy,
	}
	err = fm.SetFeePolicy(fp)
	if err != nil {
		return
	}
	if fm.pfsProxy != nil {
		log.Info("init fee module with pfs success")
	} else {
		log.Info("init fee module without pfs success")
	}
	return
}

// SetFeePolicy :
func (fm *FeeModule) SetFeePolicy(fp *models.FeePolicy) (err error) {
	if fp == nil {
		return errors.New("can not set nil fee policy")
	}
	if fp.AccountFee == nil {
		return errors.New("AccountFee can not be nil")
	}
	if fp.TokenFeeMap == nil {
		return errors.New("TokenFeeMap can not be nil")
	}
	if fp.ChannelFeeMap == nil {
		return errors.New("ChannelFeeMap can not be nil")
	}
	fm.lock.Lock()
	defer fm.lock.Unlock()
	// set fee policy to pfs
	if fm.pfsProxy != nil {
		err = fm.pfsProxy.SetFeePolicy(fp)
		if err != nil {
			log.Error(fmt.Sprintf("commit fee policy to pfs failed, err = %s", err.Error()))
			return
		}
	}
	// set fee policy to db
	err = fm.db.SaveFeePolicy(fp)
	if err != nil {
		if fm.pfsProxy != nil {
			log.Error("save fee policy to db err,may cause different fee policy between photon and pfs")
		}
		return
	}
	fm.feePolicy = fp
	return
}

//GetNodeChargeFee : impl of FeeCharge
func (fm *FeeModule) GetNodeChargeFee(nodeAddress, tokenAddress common.Address, amount *big.Int) *big.Int {
	var feeSetting *models.FeeSetting
	var ok bool
	// 优先channel
	c, err := fm.db.GetChannel(tokenAddress, nodeAddress)
	if c != nil && err == nil {
		feeSetting, ok = fm.feePolicy.ChannelFeeMap[c.ChannelIdentifier.ChannelIdentifier]
		if ok {
			return calculateFee(feeSetting, amount)
		}
	}
	// 其次token
	feeSetting, ok = fm.feePolicy.TokenFeeMap[tokenAddress]
	if ok {
		return calculateFee(feeSetting, amount)
	}
	// 最后account
	return calculateFee(fm.feePolicy.AccountFee, amount)
}

func calculateFee(feeSetting *models.FeeSetting, amount *big.Int) *big.Int {
	fee := big.NewInt(0)
	if feeSetting.FeePercent > 0 {
		fee = fee.Div(amount, big.NewInt(feeSetting.FeePercent))
	}
	if feeSetting.FeeConstant.Cmp(big.NewInt(0)) > 0 {
		fee = fee.Add(fee, feeSetting.FeeConstant)
	}
	return fee
}
