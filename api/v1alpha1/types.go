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

type HuaweiLoadBalancerSpec struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SubnetId string `json:"subnetId"`
}

type HuaweiElbMembers struct {
	ID   string `json:"id"`
	Port int32  `json:"port"`
}

// RootVolumeProperties contains the information regarding the system disk including performance, size, name, and category
type RootVolumeProperties struct {
	VolumeType string `json:"volumeType,omitempty"`

	Size int32 `json:"size,omitempty"`

	Iops int32 `json:"iops,omitempty"`

	Throughput int32 `json:"throughput,omitempty"`

	SnapshotID string `json:"snapshotId,omitempty"`
}

type Charging struct {
	ChargingMode        string `json:"chargingMode,omitempty"`
	PeriodType          string `json:"periodType,omitempty"`
	PeriodNum           int32  `json:"periodNum,omitempty"`
	IsAutoPay           bool   `json:"isAutoPay,omitempty"`
	IsAutoRenew         bool   `json:"isAutoRenew,omitempty"`
	EnterpriseProjectId string `json:"enterpriseProjectId,omitempty"`
}

type ServerSchedulerHints struct {
	Group           string `json:"group,omitempty"`
	Tenancy         string `json:"tenancy,omitempty"`
	DedicatedHostId string `json:"dedicatedHostId,omitempty"`
}

// DataVolumeProperties contains the information regarding the datadisk attached to an instance
type DataVolumeProperties struct {
	VolumeType string `json:"volumeType,omitempty"`

	Size int32 `json:"size,omitempty"`

	Iops int32 `json:"iops,omitempty"`

	Throughput int32 `json:"throughput,omitempty"`

	SnapshotID string `json:"snapshotId,omitempty"`

	Multiattach bool `json:"multiattach,omitempty"`

	Passthrough bool `json:"passthrough,omitempty"`

	ClusterID string `json:"clusterId,omitempty"`

	ClusterType string `json:"clusterType,omitempty"`

	DataImageId string `json:"dataImageId,omitempty"`
}

type SecurityGroup struct {
	// ID of resource
	// +optional
	ID *string `json:"id,omitempty"`
}

// HuaweiCloudTag is the name/value pair for a tag
type HuaweiCloudTag struct {
	// Name of the tag
	Name string `json:"name"`
	// Value of the tag
	Value string `json:"value"`
}

// InstanceState describes the state of an Huawei instance.
type InstanceState string

var (
	// InstanceStateBuild is the string representing an instance in a build state.
	InstanceStateBuild = InstanceState("build")

	// InstanceStateActive is the string representing an instance in a active state.
	InstanceStateActive = InstanceState("active")

	// InstanceStateShutOff is the string representing an instance in a shutoff state.
	InstanceStateShutOff = InstanceState("shutoff")

	// InstanceStateUnknown is the string representing an instance in a unknown state.
	InstanceStateUnknown = InstanceState("unknown")

	// InstanceStateStopping is the string representing an instance
	// that is in the process of being stopped and can be restarted.
	InstanceStateStopping = InstanceState("stopping")

	// InstanceStateStopped is the string representing an instance
	// that has been stopped and can be restarted.
	InstanceStateStopped = InstanceState("stopped")
)
