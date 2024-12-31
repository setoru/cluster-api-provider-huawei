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

package ecs

import (
	"fmt"
	"time"

	ecsMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	"github.com/pkg/errors"
	"github.com/setoru/cluster-api-provider-huawei/internal/cloud/scope"
	"k8s.io/klog/v2"
)

const (
	// DefaultWaitForInterval default interval
	DefaultWaitForInterval   = 5
	ECSInstanceStatusRunning = "ACTIVE"
	// InstanceDefaultTimeout default timeout
	InstanceDefaultTimeout = 900
)

func CreateInstance(machineScope *scope.MachineScope) (*ecsMdl.ServerDetail, error) {
	client := machineScope.HuaweiClient
	request := &ecsMdl.CreateServersRequest{Body: &ecsMdl.CreateServersRequestBody{Server: &ecsMdl.PrePaidServer{}}}

	// ImageID
	imageId := machineScope.HuaweiMachine.Spec.ImageID
	request.Body.Server.ImageRef = imageId

	request.Body.Server.FlavorRef = machineScope.HuaweiMachine.Spec.Flavor
	request.Body.Server.Vpcid = machineScope.HuaweiMachine.Spec.VpcID

	nics := make([]ecsMdl.PrePaidServerNic, 0)
	nics = append(nics, ecsMdl.PrePaidServerNic{SubnetId: machineScope.HuaweiMachine.Spec.SubnetId})
	request.Body.Server.Nics = nics

	securityGroups := make([]ecsMdl.PrePaidServerSecurityGroup, 0)
	securityGroupIDs, err := getSecurityGroupIDs(machineScope)
	if err != nil {
		return nil, errors.Wrap(err, "error getting security groups ID")
	}
	for _, id := range *securityGroupIDs {
		securityGroups = append(securityGroups, ecsMdl.PrePaidServerSecurityGroup{Id: &id})
	}
	request.Body.Server.SecurityGroups = &securityGroups

	if machineScope.HuaweiMachine.Spec.PublicIP {
		var size int32 = 100
		mode := "traffic"
		publicip := ecsMdl.PrePaidServerPublicip{Eip: &ecsMdl.PrePaidServerEip{}}
		publicip.Eip.Iptype = "5_bgp"
		publicip.Eip.Bandwidth = &ecsMdl.PrePaidServerEipBandwidth{Size: &size, Sharetype: ecsMdl.GetPrePaidServerEipBandwidthSharetypeEnum().PER, Chargemode: &mode}

		request.Body.Server.Publicip = &publicip
	}

	request.Body.Server.Name = machineScope.HuaweiMachine.GetName()

	userData := ""
	if userData != "" {
		request.Body.Server.UserData = &userData
	}

	request.Body.Server.RootVolume = getRootVolumeProperties(machineScope)

	request.Body.Server.DataVolumes = getDataVolumeProperties(machineScope)

	request.Body.Server.AvailabilityZone = &machineScope.HuaweiMachine.Spec.AvailabilityZone

	request.Body.Server.BatchCreateInMultiAz = &machineScope.HuaweiMachine.Spec.BatchCreateInMultiAz

	metedata := make(map[string]string)
	metedata["virtual_env_type"] = "IsoImage"
	request.Body.Server.Metadata = metedata

	request.Body.Server.Extendparam = getCharging(machineScope)

	diskPrior := "true"
	request.Body.Server.Extendparam.DiskPrior = &diskPrior
	request.Body.Server.Extendparam.RegionID = &machineScope.HuaweiMachine.Spec.RegionID

	request.Body.Server.OsschedulerHints = getServerSchedulerHints(machineScope)

	response, err := client.CreateServers(request)
	klog.V(4).Infof("CreateServers response:%v", response)
	if err != nil {
		klog.Errorf("Error creating ECS instance: %v", err)
		return nil, errors.Wrap(err, "error creating ECS instance")
	}

	if response == nil || len(*response.ServerIds) != 1 {
		klog.Errorf("Unexpected reservation creating instances: %v", response)
		return nil, errors.Wrap(err, "unexpected reservation creating instance")
	}

	// Sleep
	time.Sleep(5 * time.Second)

	instance, err := waitForInstancesStatus(machineScope, *response.ServerIds, ECSInstanceStatusRunning, InstanceDefaultTimeout)

	if err != nil {
		klog.Errorf("Error waiting ECS instance to Running: %v", err)
		return nil, errors.Wrap(err, "error waiting ECS instance to Running")
	}

	if instance == nil || len(instance) < 1 {
		return nil, errors.Wrap(err, "ECS instance not found")
	}

	return instance[0], nil
}

func getSecurityGroupIDs(machineScope *scope.MachineScope) (*[]string, error) {
	var securityGroupIDs []string

	// If SecurityGroupID is assigned, use it directly
	if len(machineScope.HuaweiMachine.Spec.SecurityGroups) == 0 {
		return &[]string{}, nil
	}

	for _, sg := range machineScope.HuaweiMachine.Spec.SecurityGroups {
		securityGroupIDs = append(securityGroupIDs, *sg.ID)
	}
	if len(securityGroupIDs) == 0 {
		return nil, errors.New("no securitygroup IDs found from configuration")
	}
	return &securityGroupIDs, nil
}

func getRootVolumeProperties(machineScope *scope.MachineScope) *ecsMdl.PrePaidServerRootVolume {
	rootVolume := ecsMdl.PrePaidServerRootVolume{}
	switch machineScope.HuaweiMachine.Spec.RootVolume.VolumeType {
	case "SSD":
		rootVolume.Volumetype = ecsMdl.GetPrePaidServerRootVolumeVolumetypeEnum().SSD
	case "GPSSD":
		rootVolume.Volumetype = ecsMdl.GetPrePaidServerRootVolumeVolumetypeEnum().GPSSD
	case "SATA":
		rootVolume.Volumetype = ecsMdl.GetPrePaidServerRootVolumeVolumetypeEnum().SATA
	case "SAS":
		rootVolume.Volumetype = ecsMdl.GetPrePaidServerRootVolumeVolumetypeEnum().SAS
	case "GPSSD2":
		rootVolume.Volumetype = ecsMdl.GetPrePaidServerRootVolumeVolumetypeEnum().GPSSD2
		rootVolume.Iops = &machineScope.HuaweiMachine.Spec.RootVolume.Iops
		rootVolume.Throughput = &machineScope.HuaweiMachine.Spec.RootVolume.Throughput
	case "ESSD":
		rootVolume.Volumetype = ecsMdl.GetPrePaidServerRootVolumeVolumetypeEnum().ESSD
	case "ESSD2":
		rootVolume.Volumetype = ecsMdl.GetPrePaidServerRootVolumeVolumetypeEnum().ESSD2
		rootVolume.Iops = &machineScope.HuaweiMachine.Spec.RootVolume.Iops
	default:
		rootVolume.Volumetype = ecsMdl.GetPrePaidServerRootVolumeVolumetypeEnum().GPSSD
	}
	rootVolume.Extendparam = &ecsMdl.PrePaidServerRootVolumeExtendParam{SnapshotId: &machineScope.HuaweiMachine.Spec.RootVolume.SnapshotID}
	return &rootVolume
}

func getDataVolumeProperties(machineScope *scope.MachineScope) *[]ecsMdl.PrePaidServerDataVolume {
	dataVolumes := make([]ecsMdl.PrePaidServerDataVolume, 0)
	if len(machineScope.HuaweiMachine.Spec.DataVolumes) < 0 {
		dataVolumes = append(dataVolumes, ecsMdl.PrePaidServerDataVolume{Volumetype: ecsMdl.GetPrePaidServerDataVolumeVolumetypeEnum().GPSSD, Size: 40})
		return &dataVolumes
	}
	for _, volume := range machineScope.HuaweiMachine.Spec.DataVolumes {
		dataVolume := ecsMdl.PrePaidServerDataVolume{}
		switch volume.VolumeType {
		case "SSD":
			dataVolume.Volumetype = ecsMdl.GetPrePaidServerDataVolumeVolumetypeEnum().SSD
		case "GPSSD":
			dataVolume.Volumetype = ecsMdl.GetPrePaidServerDataVolumeVolumetypeEnum().GPSSD
		case "SATA":
			dataVolume.Volumetype = ecsMdl.GetPrePaidServerDataVolumeVolumetypeEnum().SATA
		case "SAS":
			dataVolume.Volumetype = ecsMdl.GetPrePaidServerDataVolumeVolumetypeEnum().SAS
		case "GPSSD2":
			dataVolume.Volumetype = ecsMdl.GetPrePaidServerDataVolumeVolumetypeEnum().GPSSD2
			dataVolume.Iops = &volume.Iops
			dataVolume.Throughput = &volume.Throughput
		case "ESSD":
			dataVolume.Volumetype = ecsMdl.GetPrePaidServerDataVolumeVolumetypeEnum().ESSD
		case "ESSD2":
			dataVolume.Volumetype = ecsMdl.GetPrePaidServerDataVolumeVolumetypeEnum().ESSD2
			dataVolume.Iops = &volume.Iops
		default:
			dataVolume.Volumetype = ecsMdl.GetPrePaidServerDataVolumeVolumetypeEnum().GPSSD
		}
		dataVolume.Extendparam = &ecsMdl.PrePaidServerDataVolumeExtendParam{SnapshotId: &volume.SnapshotID}
		dataVolume.DataImageId = &volume.DataImageId
		if volume.ClusterID != "" {
			dataVolume.ClusterId = &volume.ClusterID
			clusterType := ecsMdl.GetPrePaidServerDataVolumeClusterTypeEnum().DSS
			dataVolume.ClusterType = &clusterType
		}
		dataVolume.Multiattach = &volume.Multiattach
		dataVolume.Hwpassthrough = &volume.Passthrough
		dataVolumes = append(dataVolumes, dataVolume)
	}
	return &dataVolumes
}

func getCharging(machineScope *scope.MachineScope) *ecsMdl.PrePaidServerExtendParam {
	var charging ecsMdl.PrePaidServerExtendParam
	var chargingMode ecsMdl.PrePaidServerExtendParamChargingMode
	var isAutoPay ecsMdl.PrePaidServerExtendParamIsAutoPay
	var isAutoRenew ecsMdl.PrePaidServerExtendParamIsAutoRenew

	switch machineScope.HuaweiMachine.Spec.Charging.ChargingMode {
	case "prePaid":
		chargingMode = ecsMdl.GetPrePaidServerExtendParamChargingModeEnum().PRE_PAID
		periodType := machineScope.HuaweiMachine.Spec.Charging.PeriodType
		periodNum := machineScope.HuaweiMachine.Spec.Charging.PeriodNum
		switch periodType {
		case "month":
			month := ecsMdl.GetPrePaidServerExtendParamPeriodTypeEnum().MONTH
			charging.PeriodType = &month
		case "year":
			year := ecsMdl.GetPrePaidServerExtendParamPeriodTypeEnum().YEAR
			charging.PeriodType = &year
		}
		charging.PeriodNum = &periodNum
		switch machineScope.HuaweiMachine.Spec.Charging.IsAutoPay {
		case true:
			isAutoPay = ecsMdl.GetPrePaidServerExtendParamIsAutoPayEnum().TRUE
		case false:
			isAutoPay = ecsMdl.GetPrePaidServerExtendParamIsAutoPayEnum().FALSE
		}

		switch machineScope.HuaweiMachine.Spec.Charging.IsAutoRenew {
		case true:
			isAutoRenew = ecsMdl.GetPrePaidServerExtendParamIsAutoRenewEnum().TRUE
		case false:
			isAutoRenew = ecsMdl.GetPrePaidServerExtendParamIsAutoRenewEnum().FALSE
		}
	case "PostPaid":
		chargingMode = ecsMdl.GetPrePaidServerExtendParamChargingModeEnum().POST_PAID
	default:
		chargingMode = ecsMdl.GetPrePaidServerExtendParamChargingModeEnum().POST_PAID
	}

	charging.ChargingMode = &chargingMode
	charging.IsAutoPay = &isAutoPay
	charging.IsAutoRenew = &isAutoRenew
	charging.EnterpriseProjectId = &machineScope.HuaweiMachine.Spec.Charging.EnterpriseProjectId
	return &charging
}

func getServerSchedulerHints(machineScope *scope.MachineScope) *ecsMdl.PrePaidServerSchedulerHints {
	if machineScope.HuaweiMachine.Spec.ServerSchedulerHints.Group == "" {
		return nil
	}
	var serverSchedulerHints ecsMdl.PrePaidServerSchedulerHints
	switch machineScope.HuaweiMachine.Spec.ServerSchedulerHints.Tenancy {
	case "shared":
		shared := ecsMdl.GetPrePaidServerSchedulerHintsTenancyEnum().SHARED
		serverSchedulerHints.Tenancy = &shared
	case "dedicated":
		dedicated := ecsMdl.GetPrePaidServerSchedulerHintsTenancyEnum().DEDICATED
		serverSchedulerHints.Tenancy = &dedicated
	}
	serverSchedulerHints.Group = &machineScope.HuaweiMachine.Spec.ServerSchedulerHints.Group
	serverSchedulerHints.DedicatedHostId = &machineScope.HuaweiMachine.Spec.ServerSchedulerHints.DedicatedHostId
	return &serverSchedulerHints
}

// waitForInstancesStatus waits for instances to given status when instance.NotFound wait until timeout
func waitForInstancesStatus(machineScope *scope.MachineScope, instanceIds []string, instanceStatus string, timeout int) ([]*ecsMdl.ServerDetail, error) {
	client := machineScope.HuaweiClient
	publicIP := machineScope.HuaweiMachine.Spec.PublicIP
	result, err := WaitForResult(fmt.Sprintf("Wait for the instances %v state to change to %s ", instanceIds, instanceStatus), func() (stop bool, result interface{}, err error) {
		showServerRequest := ecsMdl.ListServersDetailsRequest{}
		//todo
		var serverIds string
		for _, id := range instanceIds {
			if serverIds == "" {
				serverIds = id
				continue
			}
			serverIds += ","
			serverIds += id
		}
		showServerRequest.ServerId = &serverIds
		listServersDetailsResponse, err := client.ListServersDetails(&showServerRequest)
		klog.V(3).Infof("get instance status resonpse： %v", listServersDetailsResponse)
		if err != nil {
			return false, nil, err
		}

		if len(*listServersDetailsResponse.Servers) <= 0 {
			return true, nil, fmt.Errorf("the instances %v not found. ", instanceIds)
		}

		idsLen := len(instanceIds)
		servers := make([]*ecsMdl.ServerDetail, 0)

		for _, server := range *listServersDetailsResponse.Servers {
			if server.Status == instanceStatus {
				servers = append(servers, &server)
			}
			if publicIP {
				wait := true
				for _, addresses := range server.Addresses {
					for _, address := range addresses {
						if *address.OSEXTIPStype == ecsMdl.GetServerAddressOSEXTIPStypeEnum().FLOATING {
							klog.V(3).Infof("get public ip： %v", address.Addr)
							wait = false
						}
					}
				}
				if wait {
					return false, nil, fmt.Errorf("wait for public ip ")
				}
			}
		}

		if len(servers) == idsLen {
			return true, servers, nil
		}

		return false, nil, fmt.Errorf("the instances  %v state are not  the expected state  %s ", instanceIds, instanceStatus)
	}, false, DefaultWaitForInterval, timeout)

	if err != nil {
		klog.Errorf("Wait for the instances %v state change to %v occur error %v", instanceIds, instanceStatus, err)
		return nil, err
	}
	klog.V(4).Infof(" Wait for the instances complete,result:%v", result)
	if result == nil {
		return nil, nil
	}
	return result.([]*ecsMdl.ServerDetail), nil
}

func GetExistingInstance(machineScope *scope.MachineScope) (*ecsMdl.ServerDetail, error) {
	providerId := machineScope.HuaweiMachine.Spec.ProviderID
	client := machineScope.HuaweiClient
	request := ecsMdl.ListServersDetailsRequest{}
	request.ServerId = providerId
	response, err := client.ListServersDetails(&request)
	klog.V(4).Infof("ListServersDetails response:%v,id:%v", response, providerId)
	if err != nil {
		return nil, err
	}
	instances := *response.Servers
	instance := instances[0]
	return &instance, err
}

// WaitForResult wait func
func WaitForResult(name string, predicate func() (bool, interface{}, error), returnWhenError bool, delay int, timeout int) (interface{}, error) {
	endTime := time.Now().Add(time.Duration(timeout) * time.Second)
	delaySecond := time.Duration(delay) * time.Second
	for {
		// Execute the function
		satisfied, result, err := predicate()
		if err != nil {
			klog.Errorf("%s Invoke func %++s error %++v", name, "predicate func() (bool, error)", err)
			if returnWhenError {
				return result, err
			}
		}
		if satisfied {
			return result, nil
		}
		// Sleep
		time.Sleep(delaySecond)
		// If a timeout is set, and that's been exceeded, shut it down
		if timeout >= 0 && time.Now().After(endTime) {
			return nil, fmt.Errorf("wait for %s timeout", name)
		}
	}
}

func DeleteInstance(machineScope *scope.MachineScope, instanceID string) error {
	client := machineScope.HuaweiClient
	request := ecsMdl.DeleteServersRequest{Body: &ecsMdl.DeleteServersRequestBody{}}
	ids := make([]ecsMdl.ServerId, 0)
	ids = append(ids, ecsMdl.ServerId{Id: instanceID})
	deletePublicip := true
	deleteVolume := true
	request.Body.Servers = ids
	request.Body.DeletePublicip = &deletePublicip
	request.Body.DeleteVolume = &deleteVolume
	response, err := client.DeleteServers(&request)
	if err != nil {
		klog.Errorf("failed to delete instances %v error %v", response, err)
		return fmt.Errorf("failed to delete instaces: %v", err)
	}
	return nil
}
