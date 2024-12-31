/*
Copyright 2024.

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
	"context"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	infrav1 "github.com/setoru/cluster-api-provider-huawei/api/v1alpha1"
	hwclient "github.com/setoru/cluster-api-provider-huawei/internal/cloud/client"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MachineScope struct {
	Logger        logr.Logger
	Client        client.Client
	Cluster       *clusterv1.Cluster
	Machine       *clusterv1.Machine
	HuaweiMachine *infrav1.HuaweiMachine
	HuaweiClient  hwclient.Client
}

// GetProviderID returns the HuaweiMachine providerID from the spec.
func (m *MachineScope) GetProviderID() string {
	if m.HuaweiMachine.Spec.ProviderID != nil {
		return *m.HuaweiMachine.Spec.ProviderID
	}
	return ""
}

// GetRawBootstrapDataWithFormat returns the bootstrap data from the secret in the Machine's bootstrap.dataSecretName.
func (m *MachineScope) GetRawBootstrapDataWithFormat() ([]byte, string, error) {
	if m.Machine.Spec.Bootstrap.DataSecretName == nil {
		return nil, "", errors.New("error retrieving bootstrap data: linked Machine's bootstrap.dataSecretName is nil")
	}

	secret := &corev1.Secret{}
	key := types.NamespacedName{Namespace: m.Machine.GetNamespace(), Name: *m.Machine.Spec.Bootstrap.DataSecretName}
	if err := m.Client.Get(context.TODO(), key, secret); err != nil {
		return nil, "", errors.Wrapf(err, "failed to retrieve bootstrap data secret for AWSMachine %s/%s", m.Machine.GetNamespace(), m.Machine.GetName())
	}

	value, ok := secret.Data["value"]
	if !ok {
		return nil, "", errors.New("error retrieving bootstrap data: secret value key is missing")
	}

	return value, string(secret.Data["format"]), nil
}
