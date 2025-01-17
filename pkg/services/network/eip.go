package network

import (
	"fmt"

	eipMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/eip/v2/model"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/cluster-api/util"
)

func (s *Service) allocatePublicIp() (string, error) {
	createPublicIpRequest := &eipMdl.CreatePublicipRequest{}
	publicIpBody := &eipMdl.CreatePublicipOption{
		Type: "5_bgp",
	}
	bandwidthBody := &eipMdl.CreatePublicipBandwidthOption{
		ChargeMode: ptr.To(eipMdl.GetCreatePublicipBandwidthOptionChargeModeEnum().TRAFFIC),
		Name:       ptr.To(fmt.Sprintf("eip-%s", util.RandomString(4))),
		ShareType:  eipMdl.GetCreatePublicipBandwidthOptionShareTypeEnum().PER,
		Size:       ptr.To[int32](100),
	}
	createPublicIpRequest.Body = &eipMdl.CreatePublicipRequestBody{
		Publicip:  publicIpBody,
		Bandwidth: bandwidthBody,
	}
	createPublicIpResponse, err := s.eipClient.CreatePublicip(createPublicIpRequest)
	if err != nil {
		return "", errors.Wrap(err, "failed to create public ip")
	}
	return *createPublicIpResponse.Publicip.Id, nil
}

func (s *Service) releasePublicIp(publicIpId string) error {
	delPubIpReq := &eipMdl.DeletePublicipRequest{
		PublicipId: publicIpId,
	}
	delPubIpRes, err := s.eipClient.DeletePublicip(delPubIpReq)
	if err != nil {
		return errors.Wrapf(err, "failed to delete public ip %s", publicIpId)
	}
	klog.Infof("Delete public ip response: %v", delPubIpRes)
	return nil
}
