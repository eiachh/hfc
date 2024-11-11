#!/bin/bash

#ns,port,pwd
if [ -n "$1" ] && [ -n "$2" ]; then
  helm install mongo oci://registry-1.docker.io/bitnamicharts/mongodb -n $1 --create-namespace --set service.nodePorts.mongodb=30020,auth.rootPassword=$2,resources.requests.cpu=2,resources.requests.memory=500Mi,resources.limits.memory=8Gi,resources.requests.ephemeral-storage=1Gi,resources.limits.ephemeral-storage=10Gi
else
  helm install mongo oci://registry-1.docker.io/bitnamicharts/mongodb -n mongo --create-namespace --set service.nodePorts.mongodb=30020,resources.requests.cpu=2,resources.requests.memory=8Gi,resources.limits.memory=8Gi,resources.requests.ephemeral-storage=1Gi,resources.limits.ephemeral-storage=100Gi
fi

#// service.type="NodePort",