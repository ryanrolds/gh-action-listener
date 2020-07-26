package api

import (
	"github.com/ryanrolds/gh_action_listener/pkg/config"
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
