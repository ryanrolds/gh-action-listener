package api

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	ParamRepo = "repo"
	ParamTag  = "tag"
)

func (a *API) DeployHandler(w http.ResponseWriter, r *http.Request) {
	repo, ok := r.URL.Query()[ParamRepo]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing repo query param"))
		return
	}

	tag, ok := r.URL.Query()[ParamTag]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing tag query param"))
		return
	}

	logrus.Infof("deploy handler called %s %s", repo, tag)
	w.WriteHeader(http.StatusOK)
}
