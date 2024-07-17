# NGINX Prometheus Exporter in Kubernetes

This example shows how to run NGINX Prometheus Exporter in a Kubernetes cluster.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
## Table of Contents

- [Prerequisites](#prerequisites)
- [Create a kind cluster](#create-a-kind-cluster)
- [Deploy the NGINX Hello application and NGINX Prometheus Exporter](#deploy-the-nginx-hello-application-and-nginx-prometheus-exporter)
- [Configure port forwarding](#configure-port-forwarding)
- [Verification](#verification)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Prerequisites

- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) installed.
- [kind](https://kind.sigs.k8s.io/) installed.

## Create a kind cluster

You can create a kind cluster with the following command:

```console
kind create cluster
```

For details, see the [kind documentation](https://kind.sigs.k8s.io/docs/user/quick-start/#creating-a-cluster).

## Deploy the NGINX Hello application and NGINX Prometheus Exporter

You can deploy the NGINX Hello application and NGINX Prometheus Exporter with the following command:

```console
kubectl apply -f nginx-hello.yaml
```

## Configure port forwarding

Port forwarding is used to access the NGINX Hello application and NGINX Prometheus Exporter from your local machine.

You can configure port forwarding with the following command:

```console
kubectl port-forward service/nginx-demo 8080:80 9113:9113
```

## Verification

You can access the NGINX Hello application at [http://localhost:8080](http://localhost:8080) and the
NGINX Prometheus Exporter at [http://localhost:9113](http://localhost:9113).
