package main

import (
	// "context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	// "k8s.io/apimachinery/pkg/api/errors"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/ryanrolds/gh_action_listener/pkg/api"
	"github.com/ryanrolds/gh_action_listener/pkg/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	configFilename = "listener.yaml"
	// waitBetweenTicks = 10 * time.Second
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
	r.Use(a.MiddlewareCheckAccess)
	r.HandleFunc("/deploy", a.DeployHandler).Methods(http.MethodPost)
	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:80",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err = srv.ListenAndServe()
	logrus.Fatal(err)

	/*
			for {
				// get pods in all the namespaces by omitting namespace
				// Or specify namespace to get pods in particular namespace
				pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
				if err != nil {
					logrus.WithError(err).Fatal("problem getting pods")
				}
				logrus.Printf("There are %d pods in the cluster\n", len(pods.Items))

				// Examples for error handling:
				// - Use helper functions e.g. errors.IsNotFound()
				// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
				_, err = clientset.CoreV1().Pods("default").Get(context.TODO(), "example-xxxxx", metav1.GetOptions{})
				if errors.IsNotFound(err) {
					logrus.Printf("Pod example-xxxxx not found in default namespace\n")
				} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
					logrus.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
				} else if err != nil {
					logrus.WithError(err).Fatal("problem getting pods")
				} else {
					logrus.Printf("Found example-xxxxx pod in default namespace\n")
				}

				time.Sleep(waitBetweenTicks)
		  }
	*/
}
