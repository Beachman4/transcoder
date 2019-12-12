#!/usr/bin/env bash

set -e

tag=$(openssl rand -base64 12)

docker build -t transcoding:${tag} .

gcr=gcr.io/engineering-sandbox-228018/transcoding:${tag}
gcrLatest=gcr.io/engineering-sandbox-228018/transcoding:latest

docker tag transcoding:${tag} ${gcr}
docker tag transcoding:${tag} ${gcrLatest}

docker push ${gcr}
docker push ${gcrLatest}

export TAG=${tag}

envsubst < k8/deployment.yml | kubectl apply --namespace ns-aylon-armstrong -f -