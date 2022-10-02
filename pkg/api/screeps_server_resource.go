package api

import (
	"context"
	"net/http"

	screepsv1 "github.com/ryanrolds/screeps-server-controller/api/v1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ResourceScreepsServer = "screeps-server"
	ParamBranch           = "branch"
)

func (a *API) ScreepsServerResourceHandler(w http.ResponseWriter, r *http.Request) {
	branchValues, ok := r.URL.Query()[ParamBranch]
	if !ok {
		writeResponse(w, http.StatusBadRequest, "missing branch query param")
		return
	}

	branchName := branchValues[0]

	tagValues, ok := r.URL.Query()[ParamTag]
	if !ok {
		writeResponse(w, http.StatusBadRequest, "missing tag query param")
		return
	}

	tag := tagValues[0]

	logrus.Infof("screeps resource handler called %s %s", branchName, tag)

	resource, ok := a.config.Resources[ResourceScreepsServer]
	if !ok {
		writeResponse(w, http.StatusNotFound, "unknown resource")
		return
	}

	logrus.Infof("resource found: %v", resource)

	// Create update CRD in K8s

	// Options
	// REST API call to abs url
	// unstructured client
	// controller-runtime client

	// additional info
	// https://github.com/kubernetes-sigs/kubebuilder/blob/master/docs/using_an_external_type.md
	// github.com/ryanrolds/screeps-server-controller/api/v1
	// screeps.pedanticorderliness.com
	// v1
	// screeps-server

	// Setup client that understands the ScreepsServer resources
	scheme := runtime.NewScheme()
	err := screepsv1.AddToScheme(scheme)
	if err != nil {
		logrus.WithError(err).Error("problem adding to scheme")
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	kubeconfig := ctrl.GetConfigOrDie()
	controllerClient, err := client.New(kubeconfig, client.Options{Scheme: scheme})
	if err != nil {
		logrus.WithError(err).Error("problem creating controller client")
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Get list of screeps servers
	list := &screepsv1.ScreepsServerList{}
	err = controllerClient.List(context.TODO(), list, &client.ListOptions{Namespace: resource.Namespace})
	if err != nil {
		logrus.WithError(err).Error("problem getting screeps server list")
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Infof("found server: %v", list)

	w.WriteHeader(http.StatusOK)
}
