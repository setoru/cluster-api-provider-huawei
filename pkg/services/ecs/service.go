package ecs

import (
	"k8s.io/klog/v2"

	"github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/pkg/scope"
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	ecsregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/region"
	eip "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/eip/v2"
	eipregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/eip/v2/region"
)

type Service struct {
	scope     *scope.ClusterScope
	ecsClient *ecs.EcsClient
	eipClient *eip.EipClient
}

func NewService(scope *scope.ClusterScope) (*Service, error) {
	ecsreg, err := ecsregion.SafeValueOf(scope.Region())
	if err != nil {
		klog.Errorf("Failed to get region: %v", err)
		return nil, err
	}

	ecsHCHttpCli, err := ecs.EcsClientBuilder().
		WithRegion(ecsreg).
		WithCredential(scope.Credentials).
		SafeBuild()
	if err != nil {
		klog.Errorf("Failed to create ELB client: %v", err)
		return nil, err
	}

	ecsCli := ecs.NewEcsClient(ecsHCHttpCli)

	eipReg, err := eipregion.SafeValueOf(scope.Region())
	if err != nil {
		klog.Errorf("Failed to get region: %v", err)
		return nil, err
	}

	eipHCHttpCli, err := eip.EipClientBuilder().
		WithRegion(eipReg).
		WithCredential(scope.Credentials).
		SafeBuild()
	if err != nil {
		klog.Errorf("Failed to create EIP client: %v", err)
		return nil, err
	}
	eipCli := eip.NewEipClient(eipHCHttpCli)

	return &Service{
		ecsClient: ecsCli,
		eipClient: eipCli,
		scope:     scope,
	}, nil
}
