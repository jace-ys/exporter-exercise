# Exporter Exercise

We would like to measure your proficiency in software with this challenge. You may complete as much or as little as you like of it.

Please do not spend more than 2 hours on the challenge - if you have reach this time limit then please leave a note with any future enhancements you would like to make.

## Overview

We are looking for a bespoke Prometheus Metrics Exporter for Redis so that we can export Redis metrics and scrape them using Prometheus, that exports at least the following keys:

```
braze_redis_keys # a gauge of the total count of keys at a given time
braze_redis_build_version_info # Info of the current version of Redis running
```

The exporter can be written in any language, so long as it exports the above keys in the [Prometheus text-based format](https://prometheus.io/docs/instrumenting/exposition_formats/#text-based-format).

### Tasks

1. Fork the repository
2. Bootstrap the environment
    1. Verify that Redis is able to come up with seeded data.
3. Write the Prometheus exporter
    1. You may choose any language to write this in, or to use an existing exporter, forked to expose the required keys.
4. Dockerise the exporter
    1. Produce a Dockerfile that packages the exporter into a Docker image.
5. Deploy the exporter so that it can read Redis
    1. Provide screenshots of the exporter pulling the required keys from Redis
6. Documentation
    1. Please document your design choices, explanations, along with a step by step guide that the reviewer can follow to reproduce your results.

### Getting started

You will require the following tools, [minikube](https://minikube.sigs.k8s.io/docs/start/) and [helm](https://helm.sh/docs/helm/helm_install/). If you dont have the tools already installed below are the steps required to get these running on macOS.

```
# Download minikube (e.g. for macOS using brew) https://minikube.sigs.k8s.io/docs/start/
brew install minikube
 
# Create a cluster
minikube start
 
# Download helm (e.g. for macOS using brew) https://helm.sh/docs/intro/install/
brew install helm
```

1. Add the bitnami helm chart repo
`helm repo add bitnami https://charts.bitnami.com/bitnami`
2. Install the Redis app using the values file in this repo
`helm install redis bitnami/redis -f kubernetes/values.yaml`
3. You will now have 2 stateful sets deployed `redis-master` and `redis-slave` which have been pre-populated with keys.

### Hints

See tips on writing a [Prometheus exporter](https://prometheus.io/docs/instrumenting/writing_exporters/). Below are some links to Prometheus clients:
 * [Go](https://github.com/prometheus/client_golang)
 * [Python](https://github.com/prometheus/client_python)
 * [Ruby](https://github.com/prometheus/client_java)