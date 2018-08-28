package v1

import (
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ant0ine/go-json-rest/rest"
)

func GetRandomSecret(w rest.ResponseWriter, r *rest.Request) {

	type SecretPair struct {
		LockSecretHash string
		Secret         string
	}
	pair := new(SecretPair)
	seed := utils.Sha3(utils.NewRandomHash().Bytes())
	pair.Secret = seed.String()
	pair.LockSecretHash = utils.Sha3(seed.Bytes()).String()
	w.WriteJson(pair)
}
