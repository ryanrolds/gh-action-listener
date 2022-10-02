package api

import (
	"net/http"

	"github.com/ryanrolds/gh_action_listener/internal/config"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

type API struct {
	config    *config.Config
	k8sClient *kubernetes.Clientset
}

func NewAPI(cfg *config.Config, k8sClient *kubernetes.Clientset) *API {
	return &API{
		config:    cfg,
		k8sClient: k8sClient,
	}
}

func writeResponse(w http.ResponseWriter, status int, reason string) {
	w.WriteHeader(status)
	_, err := w.Write([]byte(reason))
	if err != nil {
		logrus.WithError(err).Warn("problem writing response")
	}
}
