package gkvdb

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/models/cb"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

//GetAllTokens returna all tokens on this registry contract
func (dao *GkvDB) GetAllTokens() (tokens models.AddressMap, err error) {
	err = dao.getKeyValueToBucket(models.BucketToken, models.KeyToken, &tokens)
	if err != nil {
		if err == ErrorNotFound {
			tokens = make(models.AddressMap)
		}
	}
	return
}

//AddToken add a new token to db,
func (dao *GkvDB) AddToken(token common.Address, tokenNetworkAddress common.Address) error {
	var m models.AddressMap
	err := dao.getKeyValueToBucket(models.BucketToken, models.KeyToken, &m)
	if err != nil {
		return err
	}
	if m[token] != utils.EmptyAddress {
		//startup ...
		log.Info("AddToken ,but already exists,should be ignored when startup...")
		return nil
	}
	m[token] = tokenNetworkAddress
	err = dao.saveKeyValueToBucket(models.BucketToken, models.KeyToken, m)
	dao.handleTokenCallback(dao.newTokenCallbacks, token)
	return err
}
func (dao *GkvDB) handleTokenCallback(m map[*cb.NewTokenCb]bool, token common.Address) {
	var cbs []*cb.NewTokenCb
	dao.mlock.Lock()
	for f := range m {
		remove := (*f)(token)
		if remove {
			cbs = append(cbs, f)
		}
	}
	for _, f := range cbs {
		delete(m, f)
	}
	dao.mlock.Unlock()
}

//UpdateTokenNodes update all nodes that open channel
func (dao *GkvDB) UpdateTokenNodes(token common.Address, nodes []common.Address) error {
	return dao.saveKeyValueToBucket(models.BucketTokenNodes, token[:], nodes)
}

//GetTokenNodes return all nodes has channel with me
func (dao *GkvDB) GetTokenNodes(token common.Address) (nodes []common.Address) {
	err := dao.getKeyValueToBucket(models.BucketTokenNodes, token[:], &nodes)
	if err != nil {
		log.Warn(fmt.Sprintf("GetTokenNodes for %s err=%s", token.String(), err))
	}
	return
}
