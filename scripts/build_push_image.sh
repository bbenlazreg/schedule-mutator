#!/bin/bash -e

NAME=schedulemutator
ECR_PREFIX=eu.gcr.io
IMAGE_NAME=$(basename `pwd`)
IMAGE_VERSION=$(git log --abbrev-commit --format=%h -s | head -n 1)

echo $TOKEN | docker login -u oauth2accesstoken --password-stdin https://eu.gcr.io
docker build --no-cache -t $ECR_PREFIX/$TARGET_PROJECT/$IMAGE_NAME:$IMAGE_VERSION .
docker tag $ECR_PREFIX/$TARGET_PROJECT/$IMAGE_NAME:$IMAGE_VERSION $ECR_PREFIX/$TARGET_PROJECT/$IMAGE_NAME:$TAG
docker push $ECR_PREFIX/$TARGET_PROJECT/$IMAGE_NAME:$IMAGE_VERSION
docker push $ECR_PREFIX/$TARGET_PROJECT/$IMAGE_NAME:$TAG