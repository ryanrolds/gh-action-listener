# K8s files

### Namespace

```
kubectl apply -f k8s/namespace.yaml
```

### Secrets

The service requires that an access token be configured.

```
export ACCESS_TOKEN=$(echo -n <token> | base64 -w 0)
envsubst < k8s/secrets.yaml | kubectl apply -f -
```

### Accounts, Roles, and Bindings

It's expected that the Screeps Server controller is installed and the role have been applied.

```
kubectl apply -f k8s/service_account.yaml
kubectl apply -f k8s/role_binding.yaml
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
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/ingress.yaml
```

