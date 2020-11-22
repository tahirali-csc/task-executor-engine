#!/bin/sh

env GOOS=linux GOARCH=amd64 go build -o runner main.go

docker rmi -f runner:latest
docker build -t runner:latest .

kubectl delete -f pod.yaml
kubectl apply -f pod.yaml