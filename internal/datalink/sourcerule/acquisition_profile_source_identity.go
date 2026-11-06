package sourcerule

import (
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
)

func sourceAcquisitionIdentity(device schema.Device) optimizer.AcquisitionIdentity {
	config, err := connector.ResolveAcquisitionConfig(device.Protocol, device.ConnectionConfig)
	if err == nil {
		return config.Identity(device.ID)
	}
	return optimizer.AcquisitionIdentity{Protocol: acquisitionNativeProtocol(device.Protocol)}
}

func sourceIdentityForProtocol(protocol schema.ProtocolType, identities []optimizer.AcquisitionIdentity) optimizer.AcquisitionIdentity {
	if len(identities) > 0 {
		return identities[0]
	}
	return optimizer.AcquisitionIdentity{Protocol: acquisitionNativeProtocol(protocol)}
}
