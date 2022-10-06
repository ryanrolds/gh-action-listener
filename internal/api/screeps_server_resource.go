package api

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"strings"

	screepsv1 "github.com/ryanrolds/screeps-server-controller/api/v1"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ResourceScreepsServer = "screeps-server"
	ParamBranch           = "branch"
)

func getClient() (client.Client, error) {
	// Setup client that understands the ScreepsServer resources
	scheme := runtime.NewScheme()
	err := screepsv1.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}

	kubeconfig := ctrl.GetConfigOrDie()
	controllerClient, err := client.New(kubeconfig, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	return controllerClient, nil
}

func (a *API) CreateUpdateScreepsServerResourceHandler(w http.ResponseWriter, r *http.Request) {
	branchValues, ok := r.URL.Query()[ParamBranch]
	if !ok {
		logrus.Error("missing branch query param")
		writeResponse(w, http.StatusBadRequest, "missing branch query param")
		return
	}

	branchName := branchValues[0]

	// munge branch name to be a valid k8s resource name
	branchName, err := mungeBranchName(branchName)
	if err != nil {
		logrus.Error(err)
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	tagValues, ok := r.URL.Query()[ParamTag]
	if !ok {
		logrus.Error("missing tag query param")
		writeResponse(w, http.StatusBadRequest, "missing tag query param")
		return
	}

	tag := tagValues[0]

	logrus.Infof("create/update screeps resource handler called %s %s", branchName, tag)

	resource, ok := a.config.Resources[ResourceScreepsServer]
	if !ok {
		logrus.Errorf("unknown resource: %s", ResourceScreepsServer)
		writeResponse(w, http.StatusNotFound, "unknown resource")
		return
	}

	logrus.Infof("resource found: %v", resource)

	controllerClient, err := getClient()
	if err != nil {
		logrus.WithError(err).Error("problem getting client")
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check for existing CR
	server := &screepsv1.ScreepsServer{}
	err = controllerClient.Get(context.TODO(), client.ObjectKey{
		Name:      branchName,
		Namespace: resource.Namespace,
	}, server)

	// if error, check if not found error and create - otherwise return error
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			// Error other than not found
			logrus.WithError(err).Error("problem getting screeps server")
			writeResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Create new CR
		server = &screepsv1.ScreepsServer{
			ObjectMeta: metav1.ObjectMeta{
				Name:      branchName,
				Namespace: resource.Namespace,
			},
			Spec: screepsv1.ScreepsServerSpec{
				Branch: branchName,
				Tag:    tag,
			},
		}

		err = controllerClient.Create(context.TODO(), server)
		if err != nil {
			logrus.WithError(err).Error("problem creating screeps server")
			writeResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		logrus.Infof("created screeps server: %v", server)
		w.WriteHeader(http.StatusOK)
		return
	}

	// If tag is different than spec, update
	if server.Spec.Tag != tag {
		// If exists, update
		server.Spec.Tag = tag
		err = controllerClient.Update(context.TODO(), server)
		if err != nil {
			logrus.WithError(err).Error("problem updating screeps server")
			writeResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		logrus.Infof("updated screeps server: %v", server)
		w.WriteHeader(http.StatusOK)
	}

	logrus.Info("screeps server already up to date")
	w.WriteHeader(http.StatusOK)
}

func (a *API) DeleteScreepsServerResourceHandler(w http.ResponseWriter, r *http.Request) {
	branchValues, ok := r.URL.Query()[ParamBranch]
	if !ok {
		logrus.Error("missing branch query param")
		writeResponse(w, http.StatusBadRequest, "missing branch query param")
		return
	}

	branchName := branchValues[0]

	// munge branch name to be a valid k8s resource name
	branchName, err := mungeBranchName(branchName)
	if err != nil {
		logrus.Error(err)
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	logrus.Infof("delete screeps resource handler called %s", branchName)

	resource, ok := a.config.Resources[ResourceScreepsServer]
	if !ok {
		logrus.Errorf("unknown resource: %s", ResourceScreepsServer)
		writeResponse(w, http.StatusNotFound, "unknown resource")
		return
	}

	logrus.Infof("resource found: %v", resource)

	controllerClient, err := getClient()
	if err != nil {
		logrus.WithError(err).Error("problem getting client")
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = controllerClient.Delete(context.TODO(), &screepsv1.ScreepsServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      branchName,
			Namespace: resource.Namespace,
		},
	})
	if err != nil {
		logrus.WithError(err).Error("problem deleting screeps server")
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Infof("deleted screeps server: %v", branchName)
	w.WriteHeader(http.StatusOK)
}

func mungeBranchName(branchName string) (string, error) {
	branchName = strings.ReplaceAll(branchName, "/", "-")
	branchName = strings.ReplaceAll(branchName, "_", "-")
	branchName = strings.ReplaceAll(branchName, ".", "-")
	branchName = strings.ReplaceAll(branchName, " ", "")
	branchName = strings.ReplaceAll(branchName, ":", "")
	branchName = strings.ReplaceAll(branchName, "@", "")
	branchName = strings.ReplaceAll(branchName, "#", "")

	if regexp.MustCompile("[^a-zA-Z0-9-]+").MatchString(branchName) {
		return "", errors.New("invalid branch name: " + branchName)
	}

	return branchName, nil
}
