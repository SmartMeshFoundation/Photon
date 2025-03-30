package restful

import (
	"fmt"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

func writejson(w rest.ResponseWriter, result interface{}) {
	err := w.WriteJson(result)
	if err != nil {
		fmt.Println(fmt.Sprintf("writejson err %s", err))
	}
}

func HexToAddress(addr string) (address common.Address, err error) {
	address = common.HexToAddress(addr)
	if address.String() != addr {
		err = fmt.Errorf("address checksum error expect=%s,got=%s", address.String(), addr)
	}
	return
}

// EmptyHash all zero,invalid
var EmptyHash = common.Hash{}
