package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ryanrolds/gh_action_listener/internal/api"
	"github.com/ryanrolds/gh_action_listener/internal/config"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	configFilename = "listener.yaml"
)

func main() {
	cfg, err := config.GetConfig(configFilename)
	if err != nil {
		log.Fatal("problem reading backend.yaml")
	}

	initLogging()

	// creates the in-cluster config
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		logrus.WithError(err).Fatal("problem creating k8s in-cluster config")
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		logrus.WithError(err).Fatal("problem creating k8s client")
	}

	a := api.NewAPI(cfg, clientset)
	r := mux.NewRouter()

	// TODO add request logging middleware
	r.Use(a.MiddlewareCheckAccess)

	r.HandleFunc("/deploy", a.DeployHandler).Methods(http.MethodPut)
	r.HandleFunc("/resource/screeps/server", a.CreateUpdateScreepsServerResourceHandler).Methods(http.MethodPost)
	r.HandleFunc("/resource/screeps/server", a.DeleteScreepsServerResourceHandler).Methods(http.MethodDelete)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err = srv.ListenAndServe()
	logrus.Fatal(err)
}
