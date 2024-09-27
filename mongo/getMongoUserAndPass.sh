#!/bin/bash

if [ -z "$1" ]; then
    export MONGODB_ROOT_USER=$(minikube kubectl -- exec -it deployment.apps/mongo-mongodb -n mongo -- bash -c 'echo $MONGODB_ROOT_USER')
    export MONGODB_ROOT_PASSWORD=$(minikube kubectl -- exec -it deployment.apps/mongo-mongodb -n mongo -- bash -c 'echo $MONGODB_ROOT_PASSWORD')
else
    export MONGODB_ROOT_USER=$(minikube kubectl -- exec -it deployment.apps/mongo-mongodb -n $1 -- bash -c 'echo $MONGODB_ROOT_USER')
    export MONGODB_ROOT_PASSWORD=$(minikube kubectl -- exec -it deployment.apps/mongo-mongodb -n $1 -- bash -c 'echo $MONGODB_ROOT_PASSWORD')
fi

