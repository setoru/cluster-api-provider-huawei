/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package scope

import (
	"github.com/go-logr/logr"
	infrav1 "github.com/setoru/cluster-api-provider-huawei/api/v1alpha1"
	hwclient "github.com/setoru/cluster-api-provider-huawei/internal/cloud/client"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ClusterScope struct {
	Logger        logr.Logger
	Client        client.Client
	Cluster       *clusterv1.Cluster
	HuaweiCluster *infrav1.HuaweiCluster
	HuaweiClient  hwclient.Client
}
