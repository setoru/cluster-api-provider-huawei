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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	capierrors "sigs.k8s.io/cluster-api/errors"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HuaweiClusterSpec defines the desired state of HuaweiCluster.
type HuaweiClusterSpec struct {
	// The ECS Region the cluster lives in.
	Region string `json:"region,omitempty"`

	// credentialsSecret is a local reference to a secret that contains the
	// credentials data to access HuaweiCloud PC client
	// +kubebuilder:validation:Required
	CredentialsSecret *corev1.LocalObjectReference `json:"credentialsSecret"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`
}

// HuaweiClusterStatus defines the observed state of HuaweiCluster.
type HuaweiClusterStatus struct {
	// Ready describes if the Huawei Cluster can be considered ready for machine creation.
	// +optional
	Ready bool `json:"ready,omitempty"`
	// Conditions defines current service state of the Huawei cluster.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`

	// failureReason will be set in the event that there is a terminal problem reconciling the Cluster
	// and will contain a succinct value suitable for machine interpretation.
	//
	// This field should not be set for transitive errors that can be fixed automatically or with manual intervention,
	// but instead indicate that something is fundamentally wrong with the FooCluster and that it cannot be recovered.
	// +optional
	FailureReason *capierrors.ClusterStatusError `json:"failureReason,omitempty"`

	// failureMessage will be set in the event that there is a terminal problem reconciling the Cluster
	// and will contain a more verbose string suitable for logging and human consumption.
	//
	// This field should not be set for transitive errors that can be fixed automatically or with manual intervention,
	// but instead indicate that something is fundamentally wrong with the FooCluster and that it cannot be recovered.
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// HuaweiCluster is the Schema for the huaweiclusters API.
type HuaweiCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HuaweiClusterSpec   `json:"spec,omitempty"`
	Status HuaweiClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HuaweiClusterList contains a list of HuaweiCluster.
type HuaweiClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HuaweiCluster `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &HuaweiCluster{}, &HuaweiClusterList{})
}
