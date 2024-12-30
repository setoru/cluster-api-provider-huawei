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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	capierrors "sigs.k8s.io/cluster-api/errors"
)

const (
	// MachineFinalizer allows cleaning up resources associated with
	// DockerMachine before removing it from the API Server.
	MachineFinalizer = "huaweimachine.infrastructure.cluster.x-k8s.io"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HuaweiMachineSpec defines the desired state of HuaweiMachine.
type HuaweiMachineSpec struct {
	// ProviderID is the unique identifier as specified by the cloud provider.
	ProviderID *string `json:"providerID,omitempty"`

	//The flavor of the instance.
	Flavor string `json:"flavor"`

	// The ID of the vpc
	VpcID string `json:"vpcId"`

	// The ID of the subnet
	SubnetId     string `json:"subnetId"`
	SubnetIpv4Id string `json:"subnetIpv4Id"`

	// The ID of the region in which to create the instance. You can call the DescribeRegions operation to query the most recent region list.
	RegionID string `json:"regionId"`

	// The ID of the zone in which to create the instance. You can call the DescribeZones operation to query the most recent region list.
	ZoneID string `json:"zoneId"`

	// The ID of the image used to create the instance.
	ImageID string `json:"imageId"`

	// PublicIP specifies whether the instance should get a public IP.
	PublicIP bool `json:"publicIp,omitempty"`

	// RootVolume holds the properties regarding the system disk for the instance
	// +optional
	RootVolume RootVolumeProperties `json:"rootVolume,omitempty"`

	// DataVolumes holds information regarding the extra disks attached to the instance
	// +optional
	DataVolumes []DataVolumeProperties `json:"dataVolumes,omitempty"`

	ElbMembers []HuaweiElbMembers `json:"elbMembers,omitempty"`

	Charging Charging `json:"charging,omitempty"`

	AvailabilityZone string `json:"availabilityZone,omitempty"`

	BatchCreateInMultiAz bool `json:"batchCreateInMultiAz,omitempty"`

	ServerSchedulerHints ServerSchedulerHints `json:"serverSchedulerHints,omitempty"`

	SecurityGroups []SecurityGroup `json:"securityGroups,omitempty"`
}

// HuaweiMachineStatus defines the observed state of HuaweiMachine.
type HuaweiMachineStatus struct {
	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready,omitempty"`

	// Conditions defines current service state of the Huawei cluster.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`

	// Addresses contains the Huawei instance associated addresses.
	Addresses []clusterv1.MachineAddress `json:"addresses,omitempty"`

	// failureReason will be set in the event that there is a terminal problem reconciling the Machine
	// and will contain a succinct value suitable for machine interpretation.
	//
	// This field should not be set for transitive errors that can be fixed automatically or with manual intervention,
	// but instead indicate that something is fundamentally wrong with the FooMachine and that it cannot be recovered.
	// +optional
	FailureReason *capierrors.MachineStatusError `json:"failureReason,omitempty"`

	// failureMessage will be set in the event that there is a terminal problem reconciling the FooMachine
	// and will contain a more verbose string suitable for logging and human consumption.
	//
	// This field should not be set for transitive errors that can be fixed automatically or with manual intervention,
	// but instead indicate that something is fundamentally wrong with the FooMachine and that it cannot be recovered.
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`

	// See other rules for more details about mandatory/optional fields in InfraMachine status.
	// Other fields SHOULD be added based on the needs of your provider.
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// HuaweiMachine is the Schema for the huaweimachines API.
type HuaweiMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HuaweiMachineSpec   `json:"spec,omitempty"`
	Status HuaweiMachineStatus `json:"status,omitempty"`
}

// GetConditions returns the observations of the operational state of the HuaweiMachine resource.
func (h *HuaweiMachine) GetConditions() clusterv1.Conditions {
	return h.Status.Conditions
}

// SetConditions sets the underlying service state of the HuaweiMachine to the predescribed clusterv1.Conditions.
func (h *HuaweiMachine) SetConditions(conditions clusterv1.Conditions) {
	h.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// HuaweiMachineList contains a list of HuaweiMachine.
type HuaweiMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HuaweiMachine `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &HuaweiMachine{}, &HuaweiMachineList{})
}
