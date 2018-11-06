package v1

import (
	"fmt"
	"net/http"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ant0ine/go-json-rest/rest"
)

// GetRandomSecret : create a secret and lockSecretHash with sha3
func GetRandomSecret(w rest.ResponseWriter, r *rest.Request) {

	type SecretPair struct {
		LockSecretHash string `json:"lock_secret_hash"`
		Secret         string `json:"secret"`
	}
	pair := new(SecretPair)
	seed := utils.ShaSecret(utils.NewRandomHash().Bytes())
	pair.Secret = seed.String()
	pair.LockSecretHash = utils.ShaSecret(seed.Bytes()).String()
	err := w.WriteJson(pair)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

// NotifyNetworkDown :
func NotifyNetworkDown(w rest.ResponseWriter, r *rest.Request) {
	err := API.NotifyNetworkDown()
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = w.(http.ResponseWriter).Write([]byte("ok"))
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

// GetFeePolicy :
func GetFeePolicy(w rest.ResponseWriter, r *rest.Request) {
	fp, err := API.GetFeePolicy()
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = w.WriteJson(fp)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

// SetFeePolicy :
func SetFeePolicy(w rest.ResponseWriter, r *rest.Request) {
	req := &models.FeePolicy{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = API.SetFeePolicy(req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = w.(http.ResponseWriter).Write([]byte("ok"))
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}
