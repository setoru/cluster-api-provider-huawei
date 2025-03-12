# Kubernetes Cluster API Provider Huawei Cloud (CAPHW)

Kubernetes-native declarative infrastructure for [Huawei Cloud](https://www.huaweicloud.com/).

<!-- go doc / reference card -->
<a href="https://godoc.org/github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei">
<img src="https://godoc.org/github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei?status.svg"></a>
<!-- goreportcard badge -->
<a href="https://goreportcard.com/report/HuaweiCloudDeveloper/cluster-api-provider-huawei">
<img src="https://goreportcard.com/badge/HuaweiCloudDeveloper/cluster-api-provider-huawei"></a>

## What is the Cluster API Provider Huawei Cloud?

The Cluster API Provider Huawei Cloud is a Kubernetes project to bring declarative,
Kubernetes-style APIs to cluster creation, configuration, and management.
It provides optional, additive functionality on top of [Cluster API](https://github.com/kubernetes-sigs/cluster-api)
to deploy and manage Kubernetes clusters on Huawei Cloud.

## Pre-built Kubernetes HMIs

New HMIs are built on a best effort basis when a new Kubernetes version is released for each supported OS distribution
and then published to supported regions.

### Supported Os Distribution

- Ubuntu (ubuntu-22.04)

### Supported Images On Huawei Regions

supported Ubuntu-22.04 with Kubernetes v1.32.0:

| **Regions**    | **Image ID**                         |
|----------------|--------------------------------------|
| af-north-1     | 8fb791b0-9570-4f7c-b036-cc4bbc2853f0 |
| af-south-1     | 73545729-a2be-474f-adb0-e00c0ec2a1e5 |
| ap-southeast-1 | 4e98ff86-1c31-4ede-997c-44c39e618fd3 |
| ap-southeast-2 | 12896087-f3b9-40b3-9ffd-f68484d477b9 |
| ap-southeast-3 | af5d578d-6e8e-44bf-9bc3-895c2ff37449 |
| ap-southeast-4 | 939cdd00-19a8-4add-a9c5-615578f48ee2 |
| ap-southeast-5 | 1c290290-ed5d-4b4a-ae74-0aeb27e13bf0 |
| cn-east-3      | 30546bb8-2366-4468-9fdd-0293f56c405a |
| cn-north-1     | 67843a95-92f2-450a-a767-b9b9fd1a772e |
| cn-north-2     | a361110c-d91d-404d-89d2-b3e0b97b42a9 |
| cn-north-4     | 6b201604-c4a7-4abe-b6dc-cb074f744804 |
| cn-north-9     | f6f47e80-1aa3-4081-ae81-0600a7f45f31 |
| cn-north-11    | 7e39ea5c-4f69-43c7-8009-d69273312d0c |
| cn-south-1     | 9dbc85e2-eaa8-4b38-b52d-8b67e1f45127 |
| cn-south-2     | 829ee143-8c3f-4b52-996d-f0a99dbb10ba |
| cn-south-4     | e2263275-0045-4290-853d-c2c83b33d480 |
| cn-southwest-2 | 3b95a02e-3335-4388-8d12-a69afe60a2e6 |
| la-south-2     | dbc9f49a-2c1c-4a6b-abd9-0714114007bf |
| la-north-2     | 8fb791b0-9570-4f7c-b036-cc4bbc2853f0 |
| me-east-1      | fb79afff-7d52-4aa9-bbe7-eee2b6534828 |
| na-mexico-1    | 31bbf5ad-e8c6-4fdc-92b7-b0d8a07654c5 |
| sa-brazil-1    | 76c87bd3-a001-4311-ac7b-35e9fbf4353c |
| tr-west-1      | aec0fce2-d059-4e56-8532-d1fb0d0f5f6f |