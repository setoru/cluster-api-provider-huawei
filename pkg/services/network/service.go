package network

import (
	natReg "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/nat/v2/region"
	"k8s.io/klog/v2"

	"github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/pkg/scope"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	eip "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/eip/v2"
	nat "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/nat/v2"
	vpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2"
	vpcReg "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2/region"
)

type Service struct {
	scope     *scope.ClusterScope
	vpcClient *vpc.VpcClient
	eipClient *eip.EipClient
	natClient *nat.NatClient
}

func NewService(scope *scope.ClusterScope) (*Service, error) {
	region, err := vpcReg.SafeValueOf(scope.Region())
	if err != nil {
		klog.Errorf("Failed to get region: %v", err)
		return nil, err
	}

	vpcHCHttpCli, err := vpc.VpcClientBuilder().
		WithRegion(region).
		WithCredential(scope.Credentials).
		WithHttpConfig(config.DefaultHttpConfig()).
		SafeBuild()
	if err != nil {
		klog.Errorf("Failed to create VPC client: %v", err)
		return nil, err
	}
	vpcCli := vpc.NewVpcClient(vpcHCHttpCli)

	eipHCHttpCli, err := eip.EipClientBuilder().
		WithRegion(region).
		WithCredential(scope.Credentials).
		SafeBuild()
	if err != nil {
		klog.Errorf("Failed to create EIP client: %v", err)
		return nil, err
	}
	eipCli := eip.NewEipClient(eipHCHttpCli)

	natRegion, err := natReg.SafeValueOf(scope.Region())
	if err != nil {
		klog.Errorf("Failed to get region: %v", err)
		return nil, err
	}
	natHCHttpCli, err := nat.NatClientBuilder().
		WithRegion(natRegion).
		WithCredential(scope.Credentials).
		SafeBuild()
	if err != nil {
		klog.Errorf("Failed to create NAT client: %v", err)
		return nil, err
	}
	natCli := nat.NewNatClient(natHCHttpCli)

	return &Service{
		scope:     scope,
		vpcClient: vpcCli,
		eipClient: eipCli,
		natClient: natCli,
	}, nil
}
