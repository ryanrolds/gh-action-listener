package api

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func (a *API) DeployHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("deploy handler called")
	w.WriteHeader(http.StatusOK)
}
