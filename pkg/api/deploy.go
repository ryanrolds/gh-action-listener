package api

import (
	"context"
	"fmt"
	"net/http"

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
		writeResponse(w, http.StatusBadRequest, "missing repo query param")
		return
	}

	repoKey := repoValues[0]

	tagValues, ok := r.URL.Query()[ParamTag]
	if !ok {
		writeResponse(w, http.StatusBadRequest, "missing tag query param")
		return
	}

	tag := tagValues[0]

	logrus.Infof("deploy handler called %s %s", repoKey, tag)

	repo, ok := a.config.Repos[repoKey]
	if !ok {
		writeResponse(w, http.StatusNotFound, "unknown repo")
		return
	}

	deploymentsClient := a.k8sClient.AppsV1().Deployments(repo.Namespace)

	deployment, err := deploymentsClient.Get(context.TODO(), repo.ID, metav1.GetOptions{})
	if err != nil {
		logrus.WithError(err).Error("problem getting deployment")
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	for idx := range deployment.Spec.Template.Spec.Containers {
		if deployment.Spec.Template.Spec.Containers[idx].Name == repo.Name {
			deployment.Spec.Template.Spec.Containers[idx].Image = fmt.Sprintf("%s:%s", repo.Image, tag)
		}
	}

	data, err := deployment.Marshal()
	if err != nil {
		logrus.WithError(err).Error("problem marshaling deployment")
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Info(repo.ID, string(data))

	_, err = deploymentsClient.Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		logrus.WithError(err).Error("problem updating deployment")
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
