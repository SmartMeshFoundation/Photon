package v1

import (
	"fmt"

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
