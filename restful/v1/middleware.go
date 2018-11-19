package v1

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
)

// HTTPPasswordCheckMiddleWare check http password when api call
type HTTPPasswordCheckMiddleWare struct{}

// MiddlewareFunc makes AccessLogApacheMiddleware implement the Middleware interface.
func (mw *HTTPPasswordCheckMiddleWare) MiddlewareFunc(h rest.HandlerFunc) rest.HandlerFunc {

	return func(w rest.ResponseWriter, r *rest.Request) {

		if HTTPPassword != "" {
			hp := r.Header.Get("http-password")
			if hp != HTTPPassword {
				rest.Error(w, "wrong http password", http.StatusBadRequest)
				return
			}
		}
		// call the handler
		h(w, r)
	}
}
