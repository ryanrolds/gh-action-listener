# K8s files


### Secrets

The service requires that an access token be configured.

```
export ACCESS_TOKEN=$(echo -n <token> | base64 -w 0)
envsubst < k8s/secrets.yaml | kubectl apply -f -
```

### Deploy

```
docker build .
export TAG_NAME=$(docker images --format='{{.ID}}' | head -1)
docker tag $TAG_NAME docker.pedanticorderliness.com/gh-action-listener:$TAG_NAME
docker push docker.pedanticorderliness.com/gh-action-listener:$TAG_NAME
envsubst < k8s/deployment.yaml | kubectl apply -f -
```

### Setup service and ingress

```
envsubst < k8s/service.yaml | kubectl apply -f -
kubectl apply -f ingress.yaml
```

### Permission

```
kubectl create clusterrolebinding default-edit --clusterrole=edit --serviceaccount=default:default
```
