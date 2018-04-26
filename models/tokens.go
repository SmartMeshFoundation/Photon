package models

import (
	"fmt"

	log "github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
)

type AddressMap map[common.Address]common.Address

const bucketToken = "bucketToken"
const keyToken = "tokens"
const bucketTokenNodes = "bucketTokenNodes"
const keyTokenNodes = "nodes"

func (model *ModelDB) GetAllTokens() (tokens AddressMap, err error) {
	err = model.db.Get(bucketToken, keyToken, &tokens)
	return
}

//notify when new token add?
func (model *ModelDB) AddToken(token common.Address, manager common.Address) error {
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
	m[token] = manager
	err = model.db.Set(bucketToken, keyToken, m)
	model.handleTokenCallback(model.newTokenCallbacks, token)
	return err
}
func (model *ModelDB) handleTokenCallback(m map[*NewTokenCb]bool, token common.Address) {
	var cbs []*NewTokenCb
	model.mlock.Lock()
	for f, _ := range m {
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

//all nodes that open channel
func (model *ModelDB) UpdateTokenNodes(token common.Address, nodes []common.Address) error {
	return model.db.Set(bucketTokenNodes, token[:], nodes)
}
func (model *ModelDB) GetTokenNodes(token common.Address) (nodes []common.Address) {
	err := model.db.Get(bucketTokenNodes, token[:], &nodes)
	if err != nil {
		log.Warn(fmt.Sprintf("GetTokenNodes for %s err=%s", token.String(), err))
	}
	return
}
