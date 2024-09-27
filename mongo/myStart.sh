#!/bin/bash

if [ -z "$1" ]; then
  helm install mongo oci://registry-1.docker.io/bitnamicharts/mongodb -n mongo --create-namespace --set service.nodePorts.mongodb=30020,resources.requests.cpu=2,resources.requests.memory=8Gi,resources.limits.memory=8Gi,resources.requests.ephemeral-storage=1Gi,resources.limits.ephemeral-storage=100Gi
else
  helm install mongo oci://registry-1.docker.io/bitnamicharts/mongodb -n $1 --create-namespace --set service.nodePorts.mongodb=30021,resources.requests.cpu=2,resources.requests.memory=8Gi,resources.limits.memory=8Gi,resources.requests.ephemeral-storage=1Gi,resources.limits.ephemeral-storage=100Gi
fi

#// service.type="NodePort",