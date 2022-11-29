package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ParamRepo = "repo"
	ParamTag  = "tag"
)

func (a *API) DeployHandler(w http.ResponseWriter, r *http.Request) {
	repoValues, ok := r.URL.Query()[ParamRepo]
	if !ok {
		logrus.Errorf("missing repo parameter")
		writeResponse(w, http.StatusBadRequest, "missing repo query param")
		return
	}

	// get first string and lower case it
	repoKey := strings.ToLower(repoValues[0])

	tagValues, ok := r.URL.Query()[ParamTag]
	if !ok {
		logrus.Errorf("missing tag query param")
		writeResponse(w, http.StatusBadRequest, "missing tag query param")
		return
	}

	// get first tag and lower case it
	tag := strings.ToLower(tagValues[0])

	logrus.Infof("deploy handler called %s %s", repoKey, tag)

	repo, ok := a.config.Deployments[repoKey]
	if !ok {
		logrus.Errorf("repo %s not found", repoKey)
		writeResponse(w, http.StatusNotFound, "unknown repo")
		return
	}

	deploymentsClient := a.k8sClient.AppsV1().Deployments(repo.Namespace)

	deployment, err := deploymentsClient.Get(context.TODO(), repo.DeploymentName, metav1.GetOptions{})
	if err != nil {
		logrus.WithError(err).Error("problem getting deployment")
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	for idx := range deployment.Spec.Template.Spec.Containers {
		if deployment.Spec.Template.Spec.Containers[idx].Name == repo.ContainerName {
			deployment.Spec.Template.Spec.Containers[idx].Image = fmt.Sprintf("%s:%s", repo.Image, tag)
		}
	}

	_, err = deploymentsClient.Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		logrus.WithError(err).Error("problem updating deployment")
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Infof("deployment updated: %s", deployment.ObjectMeta.Name)

	w.WriteHeader(http.StatusOK)
}
