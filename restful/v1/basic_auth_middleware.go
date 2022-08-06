package v1

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/SmartMeshFoundation/Photon/params"

	"github.com/ant0ine/go-json-rest/rest"
)

// AuthMap 定义了需要密码的url列表
var AuthMap = map[string]string{
	"/api/1/withdraw/":                        http.MethodPut,
	"/api/1/deposit":                          http.MethodPut,
	"/api/1/channels/preparecooperatesettle/": http.MethodPut,
	"/api/1/channels/cancelcooperatesettle/":  http.MethodPut,
	"/api/1/channels/":                        http.MethodPatch,
	"/api/1/transfers/":                       http.MethodPost,
	"/api/1/transfer-smt/":                    http.MethodPost,
	"/api/1/auth-check":                       http.MethodGet,
}

// AuthBasicMiddleware provides a simple AuthBasic implementation. On failure, a 401 HTTP response
//is returned. On success, the wrapped middleware is called, and the userId is made available as
// request.Env["REMOTE_USER"].(string)
type AuthBasicMiddleware struct {

	// Realm name to display to the user. Required.
	Realm string
}

// MiddlewareFunc makes AuthBasicMiddleware implement the Middleware interface.
func (mw *AuthBasicMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	if mw.Realm == "" {
		log.Fatal("Realm is required")
	}

	return func(writer rest.ResponseWriter, request *rest.Request) {
		needAuth := false
		for url, method := range AuthMap {
			if strings.HasPrefix(request.URL.Path, url) && method == request.Method {
				needAuth = true
				break
			}
		}
		if needAuth {
			authHeader := request.Header.Get("Authorization")
			if authHeader == "" {
				mw.unauthorized(writer)
				return
			}

			providedUserID, providedPassword, err := mw.decodeBasicAuthHeader(authHeader)
			if err != nil {
				rest.Error(writer, "Invalid authentication", http.StatusBadRequest)
				return
			}
			if !(providedUserID == params.Cfg.HTTPUsername && providedPassword == params.Cfg.HTTPPassword) {
				mw.unauthorized(writer)
				return
			}
		}
		handler(writer, request)
	}
}

func (mw *AuthBasicMiddleware) unauthorized(writer rest.ResponseWriter) {
	writer.Header().Set("WWW-Authenticate", "Basic realm="+mw.Realm)
	rest.Error(writer, "Not Authorized", http.StatusUnauthorized)
}

func (mw *AuthBasicMiddleware) decodeBasicAuthHeader(header string) (user string, password string, err error) {

	parts := strings.SplitN(header, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Basic") {
		return "", "", errors.New("Invalid authentication")
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", errors.New("Invalid base64")
	}

	creds := strings.SplitN(string(decoded), ":", 2)
	if len(creds) != 2 {
		return "", "", errors.New("Invalid authentication")
	}

	return creds[0], creds[1], nil
}
