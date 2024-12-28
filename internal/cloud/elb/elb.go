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

package elb

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3/model"
	"github.com/pkg/errors"
	"github.com/setoru/cluster-api-provider-huawei/internal/cloud/scope"
)

func CreateLoadBalancer(clusterScope *scope.ClusterScope) (loadbalancerId string, publicIp string, err error) {
	client := clusterScope.HuaweiClient
	name := clusterScope.HuaweiCluster.Spec.LoadBalancerSpec.Name
	subnetId := clusterScope.HuaweiCluster.Spec.LoadBalancerSpec.SubnetId
	per := model.GetCreateLoadBalancerBandwidthOptionShareTypeEnum().PER
	publicIpInfo := &model.CreateLoadBalancerPublicIpOption{Bandwidth: &model.CreateLoadBalancerBandwidthOption{ShareType: &per}}
	req := &model.CreateLoadBalancerRequest{Body: &model.CreateLoadBalancerRequestBody{Loadbalancer: &model.CreateLoadBalancerOption{Name: &name, VipSubnetCidrId: &subnetId, Publicip: publicIpInfo}}}
	response, err := client.CreateLoadBalancer(req)
	loadbalancerId = *response.LoadbalancerId
	publicIp = response.Loadbalancer.Publicips[0].PublicipAddress
	if err != nil {
		return "", "", errors.Wrap(err, "failed to create load balancer")
	}
	clusterScope.HuaweiCluster.Spec.LoadBalancerSpec.ID = loadbalancerId
	return loadbalancerId, publicIp, nil
}

func DeleteLoadBalancer(clusterScope *scope.ClusterScope) error {
	client := clusterScope.HuaweiClient
	if clusterScope.HuaweiCluster.Spec.LoadBalancerSpec == nil {
		clusterScope.Logger.Info("no load balancer")
		return nil
	}
	loadbalancerId := clusterScope.HuaweiCluster.Spec.LoadBalancerSpec.ID
	req := &model.DeleteLoadBalancerRequest{LoadbalancerId: loadbalancerId}
	_, err := client.DeleteLoadBalancer(req)
	return errors.Wrap(err, "failed to delete load balancer")
}
