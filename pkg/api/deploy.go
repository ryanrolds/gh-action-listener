package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	ParamRepo = "repo"
	ParamTag  = "tag"
)

func (a *API) DeployHandler(w http.ResponseWriter, r *http.Request) {
	repoValues, ok := r.URL.Query()[ParamRepo]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("missing repo query param"))
		if err != nil {
			logrus.WithError(err).Warn("problem writing response")
		}

		return
	}

	repo := repoValues[0]

	tagValues, ok := r.URL.Query()[ParamTag]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("missing tag query param"))
		if err != nil {
			logrus.WithError(err).Warn("problem writing response")
		}

		return
	}

	tag := tagValues[0]

	logrus.Infof("deploy handler called %s %s", repo, tag)

	deployment, ok := a.config.Repos[repo]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("unknown repo"))
		if err != nil {
			logrus.WithError(err).Warn("problem writing response")
		}

		return
	}

	deploymentsClient := a.k8sClient.AppsV1().Deployments(apiv1.NamespaceDefault)
	patch := fmt.Sprintf(`{"spec":{"template":{"spec":{"containers":[{"name":"%s","image":"%s:%s"}]}}}}`,
		deployment.Name, deployment.Image, tag)
	_, err := deploymentsClient.Patch(context.TODO(), deployment.ID, types.MergePatchType, []byte(patch), metav1.PatchOptions{})
	if err != nil {
		logrus.WithError(err).Error("problem patching deployment")

		w.WriteHeader(http.StatusInternalServerError)
		_, writeErr := w.Write([]byte(err.Error()))
		if err != nil {
			logrus.WithError(writeErr).Warn("problem writing response")
		}

		return
	}

	w.WriteHeader(http.StatusOK)
}
