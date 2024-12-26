/*
Copyright 2021 The Kubernetes Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=huaweimachinetemplates,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion

// HuaweiMachineTemplate is the Schema for the huaweimachinetemplates API.
type HuaweiMachineTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec HuaweiMachineTemplateSpec `json:"spec,omitempty"`
}

type HuaweiMachineTemplateSpec struct {
	Template HuaweiMachineTemplateResource `json:"template"`
}

type HuaweiMachineTemplateResource struct {
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	ObjectMeta clusterv1.ObjectMeta `json:"metadata,omitempty"`
	Spec       HuaweiMachineSpec    `json:"spec"`
}

//+kubebuilder:object:root=true

// HuaweiMachineTemplateList contains a list of HuaweiMachineTemplate.
type HuaweiMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HuaweiMachineTemplate `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &HuaweiMachineTemplate{}, &HuaweiMachineTemplateList{})
}
