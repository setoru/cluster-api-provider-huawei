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

package client

import (
	"context"
	"fmt"
	"sync"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	reg "github.com/huaweicloud/huaweicloud-sdk-go-v3/core/region"
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	ecsMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	elb "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3"
	elbMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3/model"
	ims "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2"
	imsMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2/model"
	vpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3"
	vpcMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3/model"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apimachineryerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var mutex sync.Mutex

const (
	accessKey = "accessKey"
	secretKey = "secretKey"
)

// HuaweiCloudClientBuilderFuncType is function type for building huaweicloud client
type HuaweiCloudClientBuilderFuncType func(client client.Client, secretName, namespace, region string) (Client, error)

// Client is a wrapper object for actual HuaweiCloud SDK clients to allow for easier testing.
type Client interface {
	//ECS
	CreateServers(request *ecsMdl.CreateServersRequest) (*ecsMdl.CreateServersResponse, error)
	ShowServer(request *ecsMdl.ShowServerRequest) (*ecsMdl.ShowServerResponse, error)
	ListServersDetails(request *ecsMdl.ListServersDetailsRequest) (*ecsMdl.ListServersDetailsResponse, error)
	DeleteServers(request *ecsMdl.DeleteServersRequest) (*ecsMdl.DeleteServersResponse, error)
	ShowServerBlockDevice(request *ecsMdl.ShowServerBlockDeviceRequest) (*ecsMdl.ShowServerBlockDeviceResponse, error)
	CreateServerGroup(request *ecsMdl.CreateServerGroupRequest) (*ecsMdl.CreateServerGroupResponse, error)
	DeleteServerGroup(request *ecsMdl.DeleteServerGroupRequest) (*ecsMdl.DeleteServerGroupResponse, error)
	BatchCreateServerTags(request *ecsMdl.BatchCreateServerTagsRequest) (*ecsMdl.BatchCreateServerTagsResponse, error)
	NovaAssociateSecurityGroup(request *ecsMdl.NovaAssociateSecurityGroupRequest) (*ecsMdl.NovaAssociateSecurityGroupResponse, error)
	NovaDisassociateSecurityGroup(request *ecsMdl.NovaDisassociateSecurityGroupRequest) (*ecsMdl.NovaDisassociateSecurityGroupResponse, error)
	NovaListAvailabilityZones(request *ecsMdl.NovaListAvailabilityZonesRequest) (*ecsMdl.NovaListAvailabilityZonesResponse, error)
	NovaListServerSecurityGroups(request *ecsMdl.NovaListServerSecurityGroupsRequest) (*ecsMdl.NovaListServerSecurityGroupsResponse, error)
	BatchStopServers(request *ecsMdl.BatchStopServersRequest) (*ecsMdl.BatchStopServersResponse, error)
	ListFlavors(request *ecsMdl.ListFlavorsRequest) (*ecsMdl.ListFlavorsResponse, error)
	//VPC
	CreateVpc(request *vpcMdl.CreateVpcRequest) (*vpcMdl.CreateVpcResponse, error)
	DeleteVpc(request *vpcMdl.DeleteVpcRequest) (*vpcMdl.DeleteVpcResponse, error)
	ShowVpc(request *vpcMdl.ShowVpcRequest) (*vpcMdl.ShowVpcResponse, error)
	ListSecurityGroups(request *vpcMdl.ListSecurityGroupsRequest) (*vpcMdl.ListSecurityGroupsResponse, error)
	//ELB
	ShowLoadBalancer(request *elbMdl.ShowLoadBalancerRequest) (*elbMdl.ShowLoadBalancerResponse, error)
	CreateLoadBalancer(request *elbMdl.CreateLoadBalancerRequest) (*elbMdl.CreateLoadBalancerResponse, error)
	DeleteLoadBalancer(request *elbMdl.DeleteLoadBalancerRequest) (*elbMdl.DeleteLoadBalancerResponse, error)
	BatchCreateMembers(request *elbMdl.BatchCreateMembersRequest) (*elbMdl.BatchCreateMembersResponse, error)
	BatchDeleteMembers(request *elbMdl.BatchDeleteMembersRequest) (*elbMdl.BatchDeleteMembersResponse, error)
	//IMS
	ListImages(request *imsMdl.ListImagesRequest) (*imsMdl.ListImagesResponse, error)
}

type huaweicloudClient struct {
	ecsClient *ecs.EcsClient
	vpcClient *vpc.VpcClient
	elbClient *elb.ElbClient
	imsClient *ims.ImsClient
}

func (client *huaweicloudClient) ShowServerBlockDevice(request *ecsMdl.ShowServerBlockDeviceRequest) (*ecsMdl.ShowServerBlockDeviceResponse, error) {
	return client.ecsClient.ShowServerBlockDevice(request)
}

func (client *huaweicloudClient) CreateServerGroup(request *ecsMdl.CreateServerGroupRequest) (*ecsMdl.CreateServerGroupResponse, error) {
	return client.ecsClient.CreateServerGroup(request)
}

func (client *huaweicloudClient) DeleteServerGroup(request *ecsMdl.DeleteServerGroupRequest) (*ecsMdl.DeleteServerGroupResponse, error) {
	return client.ecsClient.DeleteServerGroup(request)
}

func (client *huaweicloudClient) BatchCreateServerTags(request *ecsMdl.BatchCreateServerTagsRequest) (*ecsMdl.BatchCreateServerTagsResponse, error) {
	return client.ecsClient.BatchCreateServerTags(request)
}

func (client *huaweicloudClient) NovaAssociateSecurityGroup(request *ecsMdl.NovaAssociateSecurityGroupRequest) (*ecsMdl.NovaAssociateSecurityGroupResponse, error) {
	return client.ecsClient.NovaAssociateSecurityGroup(request)
}

func (client *huaweicloudClient) NovaDisassociateSecurityGroup(request *ecsMdl.NovaDisassociateSecurityGroupRequest) (*ecsMdl.NovaDisassociateSecurityGroupResponse, error) {
	return client.ecsClient.NovaDisassociateSecurityGroup(request)
}

func (client *huaweicloudClient) NovaListAvailabilityZones(request *ecsMdl.NovaListAvailabilityZonesRequest) (*ecsMdl.NovaListAvailabilityZonesResponse, error) {
	return client.ecsClient.NovaListAvailabilityZones(request)
}

func (client *huaweicloudClient) NovaListServerSecurityGroups(request *ecsMdl.NovaListServerSecurityGroupsRequest) (*ecsMdl.NovaListServerSecurityGroupsResponse, error) {
	return client.ecsClient.NovaListServerSecurityGroups(request)
}

func (client *huaweicloudClient) CreateVpc(request *vpcMdl.CreateVpcRequest) (*vpcMdl.CreateVpcResponse, error) {
	return client.vpcClient.CreateVpc(request)
}

func (client *huaweicloudClient) DeleteVpc(request *vpcMdl.DeleteVpcRequest) (*vpcMdl.DeleteVpcResponse, error) {
	return client.vpcClient.DeleteVpc(request)
}

func (client *huaweicloudClient) ShowVpc(request *vpcMdl.ShowVpcRequest) (*vpcMdl.ShowVpcResponse, error) {
	return client.vpcClient.ShowVpc(request)
}

func (client *huaweicloudClient) ShowLoadBalancer(request *elbMdl.ShowLoadBalancerRequest) (*elbMdl.ShowLoadBalancerResponse, error) {
	return client.elbClient.ShowLoadBalancer(request)
}

func (client *huaweicloudClient) CreateLoadBalancer(request *elbMdl.CreateLoadBalancerRequest) (*elbMdl.CreateLoadBalancerResponse, error) {
	return client.elbClient.CreateLoadBalancer(request)
}

func (client *huaweicloudClient) DeleteLoadBalancer(request *elbMdl.DeleteLoadBalancerRequest) (*elbMdl.DeleteLoadBalancerResponse, error) {
	return client.elbClient.DeleteLoadBalancer(request)
}

func (client *huaweicloudClient) BatchCreateMembers(request *elbMdl.BatchCreateMembersRequest) (*elbMdl.BatchCreateMembersResponse, error) {
	return client.elbClient.BatchCreateMembers(request)
}

func (client huaweicloudClient) BatchDeleteMembers(request *elbMdl.BatchDeleteMembersRequest) (*elbMdl.BatchDeleteMembersResponse, error) {
	return client.elbClient.BatchDeleteMembers(request)
}

func (client *huaweicloudClient) CreateServers(request *ecsMdl.CreateServersRequest) (*ecsMdl.CreateServersResponse, error) {
	return client.ecsClient.CreateServers(request)
}

func (client *huaweicloudClient) ShowServer(request *ecsMdl.ShowServerRequest) (*ecsMdl.ShowServerResponse, error) {
	return client.ecsClient.ShowServer(request)
}

func (client *huaweicloudClient) ListServersDetails(request *ecsMdl.ListServersDetailsRequest) (*ecsMdl.ListServersDetailsResponse, error) {
	return client.ecsClient.ListServersDetails(request)
}

func (client *huaweicloudClient) DeleteServers(request *ecsMdl.DeleteServersRequest) (*ecsMdl.DeleteServersResponse, error) {
	return client.ecsClient.DeleteServers(request)
}

func (client *huaweicloudClient) ListImages(request *imsMdl.ListImagesRequest) (*imsMdl.ListImagesResponse, error) {
	return client.imsClient.ListImages(request)
}

func (client *huaweicloudClient) ListSecurityGroups(request *vpcMdl.ListSecurityGroupsRequest) (*vpcMdl.ListSecurityGroupsResponse, error) {
	return client.vpcClient.ListSecurityGroups(request)
}

func (client *huaweicloudClient) BatchStopServers(request *ecsMdl.BatchStopServersRequest) (*ecsMdl.BatchStopServersResponse, error) {
	return client.ecsClient.BatchStopServers(request)
}

func (client *huaweicloudClient) ListFlavors(request *ecsMdl.ListFlavorsRequest) (*ecsMdl.ListFlavorsResponse, error) {
	return client.ecsClient.ListFlavors(request)
}

// NewClient creates our client wrapper object for the actual HuaweiCloud clients we use.
func NewClient(ctrlRuntimeClient client.Client, secretName, namespace, region string) (Client, error) {
	credential, err := getCredentialFromSecret(ctrlRuntimeClient, secretName, namespace)
	if err != nil {
		klog.Errorf("fail to get credential %v ", err)
		return nil, err
	}
	klog.V(4).Infof("test func get client")
	if err != nil {
		return nil, err
	}
	ecsRegion := reg.NewRegion(region, fmt.Sprintf("https://ecs.%s.myhuaweicloud.com", region))
	ecsBuild, err := ecs.EcsClientBuilder().
		WithRegion(ecsRegion).
		WithCredential(credential).
		WithHttpConfig(config.DefaultHttpConfig()).
		SafeBuild()
	if err != nil {
		return nil, err
	}
	ecsClient := ecs.NewEcsClient(ecsBuild)

	vpcRegion := reg.NewRegion(region, fmt.Sprintf("https://vpc.%s.myhuaweicloud.com", region))
	vpcBuild, err := vpc.VpcClientBuilder().
		WithRegion(vpcRegion).
		WithCredential(credential).
		WithHttpConfig(config.DefaultHttpConfig()).
		SafeBuild()
	vpcClient := vpc.NewVpcClient(vpcBuild)

	elbRegion := reg.NewRegion(region, fmt.Sprintf("https://elb.%s.myhuaweicloud.com", region))
	elbBuild, err := elb.ElbClientBuilder().
		WithRegion(elbRegion).
		WithCredential(credential).
		WithHttpConfig(config.DefaultHttpConfig()).
		SafeBuild()
	elbClient := elb.NewElbClient(elbBuild)

	imsRegion := reg.NewRegion(region, fmt.Sprintf("https://ims.%s.myhuaweicloud.com", region))
	imsBuild, err := ims.ImsClientBuilder().
		WithRegion(imsRegion).
		WithCredential(credential).
		WithHttpConfig(config.DefaultHttpConfig()).
		SafeBuild()
	imsClient := ims.NewImsClient(imsBuild)

	return &huaweicloudClient{
		ecsClient: ecsClient,
		vpcClient: vpcClient,
		elbClient: elbClient,
		imsClient: imsClient,
	}, nil
}

func getCredentialFromSecret(ctrlRuntimeClient client.Client, secretName, namespace string) (auth.ICredential, error) {
	if secretName == "" {
		return nil, fmt.Errorf("secret name is empty")
	}
	var secret corev1.Secret
	if err := ctrlRuntimeClient.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: secretName}, &secret); err != nil {
		if apimachineryerrors.IsNotFound(err) {
			return nil, errors.Wrapf(err, "failed to patch secret name:%s namespace:%s", namespace, secretName)
		}
		return nil, err
	}
	return fetchCredentialsIniFromSecret(&secret)
}

func fetchCredentialsIniFromSecret(secret *corev1.Secret) (auth.ICredential, error) {
	ak, ok := secret.Data[accessKey]
	sk, ok := secret.Data[secretKey]
	if !ok {
		return nil, fmt.Errorf("failed to fetch key 'credentials' in secret data")
	}
	credential, err := basic.NewCredentialsBuilder().
		WithAk(string(ak)).
		WithSk(string(sk)).
		SafeBuild()
	if err != nil {
		return nil, err
	}
	return credential, nil
}
