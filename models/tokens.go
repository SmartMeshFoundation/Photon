package models

import (
	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/ethereum/go-ethereum/log"
)

type AddressMap map[common.Address]common.Address

const bucketToken = "bucketToken"
const keyToken = "tokens"

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
		log.Error("AddToken ,but already exists,should be ignored when startup...")
		return nil
	}
	m[token] = manager
	err = model.db.Set(bucketToken, keyToken, m)
	model.handleTokenCallback(model.NewTokenCallbacks, token)
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
