#!/bin/bash

if [ -z "${PROJECT_ID}" ]; then
   echo "[ERROR] Set GCP project ID by PROJECT_ID env var"
   exit 1
fi

for dir in `find services -type d -depth 1` ; do
    pushd $dir
    echo "gcloud container builds submit --config cloudbuild.yaml . --project $PROJECT_ID"
    popd
done
