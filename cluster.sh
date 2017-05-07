#!/bin/bash

if [ -z "${PROJECT_ID}" ]; then
   echo "[ERROR] Set GCP project ID by PROJECT_ID env var"
   exit 1
fi

gcloud container clusters create test-trace-cluster \
       --project $PROJECT \
       --zone us-west1-b \
       --scopes default,bigquery,cloud-platform,compute-rw,datastore,logging-write,monitoring-write
