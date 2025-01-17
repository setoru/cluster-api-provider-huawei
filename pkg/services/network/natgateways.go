package network

import (
	"fmt"

	infrav1alpha1 "github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/api/v1alpha1"
	natMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/nat/v2/model"
	vpcMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2/model"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
)

func (s *Service) reconcileNatGateways() error {
	klog.Info("Reconciling Nat Gateways")

	if s.scope.VPC().Id == "" {
		klog.Infof("VPC ID is empty, skipping NAT gateway reconcile")
		return nil
	}

	listSubnetsRequest := &vpcMdl.ListSubnetsRequest{
		VpcId: &s.scope.VPC().Id,
	}
	listSubnetsResponse, err := s.vpcClient.ListSubnets(listSubnetsRequest)
	if err != nil {
		return errors.Wrap(err, "failed to list subnets")
	}

	if len(*listSubnetsResponse.Subnets) == 0 {
		klog.Infof("No subnets available, skipping NAT gateways reconcile")
		return nil
	}

	for _, subnet := range *listSubnetsResponse.Subnets {
		createNatGatewayRequest := &natMdl.CreateNatGatewayRequest{}
		createNatGatewayRequest.Body = &natMdl.CreateNatGatewayRequestBody{
			NatGateway: &natMdl.CreateNatGatewayOption{
				Name:              fmt.Sprintf("nat-%s", util.RandomString(4)),
				RouterId:          s.scope.VPC().Id,
				Spec:              natMdl.GetCreateNatGatewayOptionSpecEnum().E_1,
				InternalNetworkId: subnet.Id,
			},
		}
		createNatGatewayResponse, err := s.natClient.CreateNatGateway(createNatGatewayRequest)
		if err != nil {
			return errors.Wrap(err, "failed to create nat gateway")
		}
		// allocate EIP to Nat Gateway
		publicIpId, err := s.allocatePublicIp()
		if err != nil {
			return err
		}
		// create SNAT rules to access the Internet
		err = s.createSnatRule(createNatGatewayResponse.NatGateway.Id, publicIpId, subnet.Id)
		if err != nil {
			return err
		}
		klog.Infof("Created Nat Gateway %s", createNatGatewayResponse.NatGateway.Id)
	}
	conditions.MarkTrue(s.scope.InfraCluster(), infrav1alpha1.NatGatewaysReadyCondition)
	return nil
}

func (s *Service) deleteNatGateways() error {
	if s.scope.VPC().Id == "" {
		klog.Infof("VPC ID is empty")
		return nil
	}
	listNatGatewaysRequest := &natMdl.ListNatGatewaysRequest{
		RouterId: &s.scope.VPC().Id,
	}
	listNatGatewaysResponse, err := s.natClient.ListNatGateways(listNatGatewaysRequest)
	if err != nil {
		return errors.Wrap(err, "failed to list nat gateways")
	}

	for _, natGateway := range *listNatGatewaysResponse.NatGateways {
		if err := s.deleteNatGatewaysExistingRule(natGateway.Id); err != nil {
			return err
		}
		deleteNatGatewayRequest := &natMdl.DeleteNatGatewayRequest{
			NatGatewayId: natGateway.Id,
		}
		_, err = s.natClient.DeleteNatGateway(deleteNatGatewayRequest)
		if err != nil {
			return errors.Wrap(err, "failed to delete nat gateways")
		}
		klog.Infof("Delete Nat Gateway %s", natGateway.Id)
	}
	return nil
}

func (s *Service) createSnatRule(natGatewayId, publicIpId, subnetId string) error {
	snatRequest := &natMdl.CreateNatGatewaySnatRuleRequest{}
	snatRequest.Body = &natMdl.CreateNatGatewaySnatRuleRequestOption{
		SnatRule: &natMdl.CreateNatGatewaySnatRuleOption{
			NatGatewayId: natGatewayId,
			FloatingIpId: publicIpId,
			NetworkId:    &subnetId,
		},
	}
	_, err := s.natClient.CreateNatGatewaySnatRule(snatRequest)
	if err != nil {
		return errors.Wrap(err, "failed to create nat gateway snat rule")
	}
	return nil
}

func (s *Service) deleteNatGatewaysExistingRule(natGatewayId string) error {
	listNatGatewaySnatRulesRequest := &natMdl.ListNatGatewaySnatRulesRequest{
		NatGatewayId: ptr.To([]string{natGatewayId}),
	}
	listNatGatewaySnatRulesResponse, err := s.natClient.ListNatGatewaySnatRules(listNatGatewaySnatRulesRequest)
	if err != nil {
		return errors.Wrap(err, "failed to list nat gateway snat rule")
	}

	for _, snatRule := range *listNatGatewaySnatRulesResponse.SnatRules {
		deleteNatGatewaySnatRuleRequest := &natMdl.DeleteNatGatewaySnatRuleRequest{
			NatGatewayId: natGatewayId,
			SnatRuleId:   snatRule.Id,
		}
		_, err = s.natClient.DeleteNatGatewaySnatRule(deleteNatGatewaySnatRuleRequest)
		if err != nil {
			return errors.Wrap(err, "failed to delete nat gateway snat rule")
		}
		klog.Infof("Deleted NatGateway SnatRule %s", snatRule.Id)

		if err := s.releasePublicIp(snatRule.FloatingIpId); err != nil {
			return err
		}
	}

	listNatGatewayDnatRulesRequest := &natMdl.ListNatGatewayDnatRulesRequest{
		NatGatewayId: ptr.To([]string{natGatewayId}),
	}
	listNatGatewayDnatRulesResponse, err := s.natClient.ListNatGatewayDnatRules(listNatGatewayDnatRulesRequest)
	if err != nil {
		return errors.Wrap(err, "failed to list nat gateway snat rule")
	}

	for _, dnatRule := range *listNatGatewayDnatRulesResponse.DnatRules {
		deleteNatGatewayDnatRuleRequest := &natMdl.DeleteNatGatewayDnatRuleRequest{
			NatGatewayId: natGatewayId,
			DnatRuleId:   dnatRule.Id,
		}
		_, err = s.natClient.DeleteNatGatewayDnatRule(deleteNatGatewayDnatRuleRequest)
		if err != nil {
			return errors.Wrap(err, "failed to delete nat gateway dnat rule")
		}
		klog.Infof("Deleted NatGateway DnatRule %s", dnatRule.Id)

		if err := s.releasePublicIp(dnatRule.FloatingIpId); err != nil {
			return err
		}
	}
	return nil
}
