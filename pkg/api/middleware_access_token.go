package api

import (
  "net/http"

  "github.com/sirupsen/logrus"
)

const AccessTokenHeader = "X-Access-Token"

func (a *API) MiddlewareCheckAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    token := r.Header.Get(AccessTokenHeader)

		if token != "" && token == a.config.AccessToken {
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Forbidden", http.StatusForbidden)
	})
}
