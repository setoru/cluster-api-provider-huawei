package network

import (
	"fmt"

	infrav1alpha1 "github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/api/v1alpha1"
	natMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/nat/v2/model"
	vpcMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2/model"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
)

func (s *Service) reconcileNatGateways() error {
	klog.Info("Reconciling Nat Gateways")

	if s.scope.VPC().Id == "" {
		klog.Infof("VPC ID is empty, skipping NAT gateway reconcile")
		return nil
	}

	if len(s.scope.Subnets()) == 0 {
		klog.Infof("No subnets available, skipping NAT gateways reconcile")
		return nil
	}

	existing, err := s.describeNatGatewaysBySubnet()
	if err != nil {
		return err
	}

	natGatewaysIds := make([]string, 0)
	subnetIds := make([]string, 0)

	for _, subnet := range s.scope.Subnets() {
		if ngw, ok := existing[subnet.Id]; ok {
			natGatewaysIds = append(natGatewaysIds, ngw)
			continue
		}
		subnetIds = append(subnetIds, subnet.Id)
	}

	natGatewaysIps, err := s.getNatGatewaysIps(natGatewaysIds)
	if err != nil {
		return err
	}

	s.scope.SetNatGatewaysIPs(natGatewaysIps)

	if len(subnetIds) > 0 {
		// set NatGatewayCreationStarted if the condition has never been set before
		if !conditions.Has(s.scope.InfraCluster(), infrav1alpha1.NatGatewaysReadyCondition) {
			conditions.MarkFalse(s.scope.InfraCluster(),
				infrav1alpha1.NatGatewaysReadyCondition,
				infrav1alpha1.NatGatewaysCreationStartedReason,
				clusterv1.ConditionSeverityInfo, "")
			if err := s.scope.PatchObject(); err != nil {
				return errors.Wrap(err, "failed to patch conditions")
			}
		}
		err := s.createNatGateways(subnetIds)
		if err != nil {
			return err
		}
		conditions.MarkTrue(s.scope.InfraCluster(), infrav1alpha1.NatGatewaysReadyCondition)
	}
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

func (s *Service) createNatGateways(subnetIds []string) (err error) {
	for _, subnetId := range subnetIds {
		showSubnetRequest := &vpcMdl.ShowSubnetRequest{SubnetId: subnetId}
		showSubnetResponse, err := s.vpcClient.ShowSubnet(showSubnetRequest)
		if err != nil {
			return errors.Wrap(err, "failed to find subnet")
		}
		if showSubnetResponse.Subnet.VpcId != s.scope.VPC().Id {
			continue
		}
		createNatGatewayRequest := &natMdl.CreateNatGatewayRequest{}
		createNatGatewayRequest.Body = &natMdl.CreateNatGatewayRequestBody{
			NatGateway: &natMdl.CreateNatGatewayOption{
				Name:              fmt.Sprintf("nat-%s", util.RandomString(4)),
				RouterId:          s.scope.VPC().Id,
				Spec:              natMdl.GetCreateNatGatewayOptionSpecEnum().E_1,
				InternalNetworkId: subnetId,
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
		if err = s.createSnatRule(createNatGatewayResponse.NatGateway.Id, publicIpId, subnetId); err != nil {
			return err
		}
		klog.Infof("Created Nat Gateway %s", createNatGatewayResponse.NatGateway.Id)
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

func (s *Service) describeNatGatewaysBySubnet() (map[string]string, error) {
	request := &natMdl.ListNatGatewaysRequest{
		RouterId: &s.scope.VPC().Id,
	}
	response, err := s.natClient.ListNatGateways(request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list nat gateways")
	}
	gatewaysIds := map[string]string{}
	for _, natGateway := range *response.NatGateways {
		gatewaysIds[natGateway.InternalNetworkId] = natGateway.Id
	}
	return gatewaysIds, nil
}

func (s *Service) getNatGatewaysIps(natGatewayIds []string) ([]string, error) {
	request := &natMdl.ListNatGatewaySnatRulesRequest{
		NatGatewayId: &natGatewayIds,
	}
	response, err := s.natClient.ListNatGatewaySnatRules(request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list nat gateway snat rule")
	}
	nateGatewaysIps := make([]string, 0)
	for _, snat := range *response.SnatRules {
		nateGatewaysIps = append(nateGatewaysIps, snat.FloatingIpAddress)
	}
	return nateGatewaysIps, nil
}
