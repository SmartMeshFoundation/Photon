package v1

import (
	"fmt"

	"net/http"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/ant0ine/go-json-rest/rest"
)

/*
Address is api of /api/1/address
*/
func Address(w rest.ResponseWriter, r *rest.Request) {
	data := make(map[string]interface{})
	data["our_address"] = RaidenAPI.Raiden.NodeAddress.String()
	err := w.WriteJson(data)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

/*
GetBalance : get account's balance and locked account on each token
*/
func GetBalance(w rest.ResponseWriter, r *rest.Request) {
	resp, err := RaidenAPI.GetBalance()
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = w.WriteJson(resp)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}
