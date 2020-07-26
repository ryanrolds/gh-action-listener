package main

import (
	"context"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	configFilename   = "listener.yaml"
	waitBetweenTicks = 10 * time.Second
)

func main() {
	config, err := GetConfig(configFilename)
	if err != nil {
		log.Fatal("problem reading backend.yaml")
	}

	initLogging(config)

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
}
