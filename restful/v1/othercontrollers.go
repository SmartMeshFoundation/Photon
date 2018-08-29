package v1

import (
	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
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
