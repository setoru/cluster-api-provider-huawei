# Developer Guide

## Initial setup for development environment

### Install prerequisite

- git
- make
- [Go v1.22](https://golang.org/doc/install)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Kustomize](https://github.com/kubernetes-sigs/kustomize)
- [docker](https://docs.docker.com/install/)
- [Kind](https://kind.sigs.k8s.io/)
- [clusterctl](https://github.com/kubernetes-sigs/cluster-api/releases)
- [Tilt](https://docs.tilt.dev/install.html)
- [envsubst](https://github.com/drone/envsubst)

### Setup

Get the source:

- Fork the [Cluster API Repo](https://github.com/kubernetes-sigs/cluster-api) and clone the forked Cluster API repo into your GOPATH (i.e. `~/go/src/sigs.k8s.io/cluster-api`)

- Fork the [cluster-api-provider-huawei repo](https://github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei) and clone your new repo into the GOPATH (i.e. `~/go/src/github.com/myname/cluster-api-provider-huawei`)

## Create a kind cluster

make sure you have a kind cluster and that your `KUBECONFIG` is set up correctly:

```bash
kind create cluster --name=capi-test
```

This local cluster will be running all the cluster api controllers and become the management cluster which then can be used to spin up workload clusters on Huawei.

## Running local management cluster for development

Before the next steps, make sure [initial setup for development environment](https://github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/blob/master-dev/docs/development.md#initial-setup-for-development-environment) steps are complete.

build huawei manager from local cluster-api-provider-huawei source and run it in local kind cluster:

### Setting up Development Environment with Tilt

[Tilt](https://tilt.dev/) is a tool for quickly building, pushing, and reloading Docker containers as part of a Kubernetes deployment. Many of the Cluster API engineers use it for quick iteration. Please see our [Tilt instructions](https://github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/blob/master-dev/docs/tilt.md) to get started.

## Create a workload cluster

After your kind management cluster is up and running with Tilt, you can deploy a workload cluster in the Tilt web UI based on YAML templates from the directories specified in the `template_dirs` field from the tilt-settings.yaml file.

Open the Tilt web UI and navigate to the CAPH.templates section, Click **Create default cluster** icon to deploy a new workload cluster.

## Accessing the workload cluster

The cluster will now start provisioning. You can check status with:

```bash
kubectl get cluster
```

and see an output similar to this:

```bash
NAME              PHASE         AGE   VERSION
default-8847   Provisioned   8s    v1.32.0
```

To verify the status of cluster,ensure cluster is ready:

```bash
clusterctl describe cluster <clustername> --show-conditions all
```

After the first control plane node is up and running, we can retrieve the workload cluster Kubeconfig.

```bash
clusterctl get kubeconfig <clustername> > ~/kubeconfig
```

you can access the cluster via `kubectl`. E.g.

```bash
kubectl --kubeconfig ~/kubeconfig get nodes
```

## Deploy a CNI solution

Calico is used here as an example.

```bash
kubectl --kubeconfig ~/kubeconfig apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/calico.yaml
```

## Clean up

Before deleting the kind cluster, make sure you delete all the workload clusters.

```bash
kubectl delete cluster <clustername>
tilt up (ctrl-c)
kind delete cluster
```

## Submitting PRs and testing

Pull requests and issues are highly encouraged! If you're interested in submitting PRs to the project, please be sure to run some initial checks prior to submission:

```
make lint # Runs a suite of quick scripts to check code structure
make test # Runs tests on the Go code
```

### Executing unit tests

`make test` executes the project's unit tests. These tests do not stand up a Kubernetes cluster, nor do they have external dependencies.

## Troubleshooting

**Error: Failed to Create Load Balancers - "Available Clusters  could not be Found"**

Switch to a region with sufficient cluster capacity.

