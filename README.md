# GitHub Action Listener for K8s

A in-cluster K8s service that listens for events over HTTP for a pre-configured list repo-deployemnt pairs and Custom Resources. Intended to be used with GitHub Action, but is not specifically designed for that.

## Running

Service expects to be run inside of a K8s cluster, the service will fail when it cannot discover the K8s API service.

See `k8s/`.
