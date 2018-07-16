package models

import (
	"fmt"

	log "github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/models/cb"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

//AddressMap is token address to mananger address
type AddressMap map[common.Address]common.Address

const bucketToken = "bucketToken"
const keyToken = "tokens"
const bucketTokenNodes = "bucketTokenNodes"
const keyTokenNodes = "nodes"

//GetAllTokens returna all tokens on this registry contract
func (model *ModelDB) GetAllTokens() (tokens AddressMap, err error) {
	err = model.db.Get(bucketToken, keyToken, &tokens)
	return
}

//AddToken add a new token to db,
func (model *ModelDB) AddToken(token common.Address, tokenNetworkAddress common.Address) error {
	var m AddressMap
	err := model.db.Get(bucketToken, keyToken, &m)
	if err != nil {
		return err
	}
	if m[token] != utils.EmptyAddress {
		//startup ...
		log.Info("AddToken ,but already exists,should be ignored when startup...")
		return nil
	}
	m[token] = tokenNetworkAddress
	err = model.db.Set(bucketToken, keyToken, m)
	model.handleTokenCallback(model.newTokenCallbacks, token)
	return err
}
func (model *ModelDB) handleTokenCallback(m map[*cb.NewTokenCb]bool, token common.Address) {
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
func (model *ModelDB) UpdateTokenNodes(token common.Address, nodes []common.Address) error {
	return model.db.Set(bucketTokenNodes, token[:], nodes)
}

//GetTokenNodes return all nodes has channel with me
func (model *ModelDB) GetTokenNodes(token common.Address) (nodes []common.Address) {
	err := model.db.Get(bucketTokenNodes, token[:], &nodes)
	if err != nil {
		log.Warn(fmt.Sprintf("GetTokenNodes for %s err=%s", token.String(), err))
	}
	return
}
