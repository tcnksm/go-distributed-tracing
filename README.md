# Distributed tracing for Golang

This repository contains sample k8s config and golang microservices apps to try distributed tracing with GCP [Stackdriver trace](https://cloud.google.com/trace/). After setup, you can see the following tracing. This contains tracing example both HTTP and [gRPC](http://www.grpc.io/) requests.

![](/trace.png)

## Usage

To create k8s cluster for this,

```bash
$ ./cluster.sh 
```

To build docker images on GCP Container Builder, 

```bash
$ ./cloudbuild.sh
```

To run apps in k8s, replace `$PROJECT_ID` with your GCP project ID in `kubernetes.yaml.tmpl`,

```bash
$ kubectl apply -f kubernetes.yaml
```

## Reference 

- [Google Cloud Platform Blog: Distributed tracing for Go](https://cloudplatform.googleblog.com/2017/04/distributed-tracing-for-Go.html)
- [Automatic Stackdriver Tracing for gRPC · Go, the unwritten parts](https://rakyll.org/grpc-trace/)
