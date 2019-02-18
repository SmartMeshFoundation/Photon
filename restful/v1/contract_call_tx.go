package v1

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon"
	"github.com/SmartMeshFoundation/Photon/dto"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/rerr"
	"github.com/ant0ine/go-json-rest/rest"
)

/*
ContractCallTXQuery query tx info of contract call
*/
func ContractCallTXQuery(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> ContractCallTXQuery ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	var err error
	var req photon.ContractCallTXQueryParams
	err = r.DecodeJsonPayload(&req)
	if err == rest.ErrJsonPayloadEmpty {
		// 允许空条件的查询
		err = nil
	}
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError)
		return
	}
	list, err := API.ContractCallTXQuery(&req)
	resp = dto.NewAPIResponse(err, list)
}
