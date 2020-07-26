# K8s files

Docker Build

```
docker build .
<get final tag (4ea6105f31a3) from output >
docker tag 4ea6105f31a3 docker.pedanticorderliness.com/gh_action_listener:4ea6105f31a3
docker push docker.pedanticorderliness.com/gh_action_listener:4ea6105f31a3
```

K8s

```
export ENV=prod
export TAG_NAME=4ea6105f31a3
envsubst < deployment.yaml | kubectl apply -f -
envsubst < service.yaml | kubectl apply -f -
kubectl apply -f ingress-prod.yaml
```
