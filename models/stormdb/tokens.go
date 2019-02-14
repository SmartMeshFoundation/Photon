package stormdb

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/rerr"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/models/cb"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
)

//GetAllTokens returna all tokens on this registry contract
func (model *StormDB) GetAllTokens() (tokens models.AddressMap, err error) {
	err = model.db.Get(models.BucketToken, models.KeyToken, &tokens)
	if err != nil {
		if err == storm.ErrNotFound {
			tokens = make(models.AddressMap)
		} else {
			err = rerr.ErrGeneralDBError.AppendError(err)
		}
	}
	return
}

//AddToken add a new token to db,
func (model *StormDB) AddToken(token common.Address, tokenNetworkAddress common.Address) error {
	var m models.AddressMap
	err := model.db.Get(models.BucketToken, models.KeyToken, &m)
	if err != nil {
		return models.GeneratDBError(err)
	}
	if m[token] != utils.EmptyAddress {
		//startup ...
		log.Info("AddToken ,but already exists,should be ignored when startup...")
		return nil
	}
	m[token] = tokenNetworkAddress
	err = model.db.Set(models.BucketToken, models.KeyToken, m)
	model.handleTokenCallback(model.newTokenCallbacks, token)
	return models.GeneratDBError(err)
}
func (model *StormDB) handleTokenCallback(m map[*cb.NewTokenCb]bool, token common.Address) {
	var cbs []*cb.NewTokenCb
	model.mlock.Lock()
	for f := range m {
		remove := (*f)(token)
		if remove {
			cbs = append(cbs, f)
		}
	}
	for _, f := range cbs {
		delete(m, f)
	}
	model.mlock.Unlock()
}

//UpdateTokenNodes update all nodes that open channel
func (model *StormDB) UpdateTokenNodes(token common.Address, nodes []common.Address) error {
	err := model.db.Set(models.BucketTokenNodes, token[:], nodes)
	return models.GeneratDBError(err)
}

//GetTokenNodes return all nodes has channel with me
func (model *StormDB) GetTokenNodes(token common.Address) (nodes []common.Address) {
	err := model.db.Get(models.BucketTokenNodes, token[:], &nodes)
	if err != nil {
		log.Warn(fmt.Sprintf("GetTokenNodes for %s err=%s", token.String(), err))
	}
	return
}
