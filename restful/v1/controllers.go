package v1

import (
	"fmt"
	"net/http"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

func RegisterToken(w rest.ResponseWriter, r *rest.Request) {
	token := r.PathParam("token")
	tokenAddr := common.HexToAddress(token)
	mgr, err := RaidenApi.RegisterToken(tokenAddr)
	type Ret struct {
		Channel_manager_address string
	}
	if err != nil {
		log.Error(fmt.Sprintf("RegisterToken %s err:%s", tokenAddr.String(), err))
		rest.Error(w, err.Error(), http.StatusConflict)
	} else {
		ret := &Ret{Channel_manager_address: mgr.String()}
		w.WriteJson(ret)
	}
}

func Stop(w rest.ResponseWriter, r *rest.Request) {
	//test only
	RaidenApi.Stop()
	w.Header().Set("Content-Type", "text/plain")
	w.(http.ResponseWriter).Write([]byte("ok"))
}
