package connector

import "context"

// AcquisitionConnectionToken identifies ephemeral capability-cache ownership.
// It is never source, point, sample or persisted configuration identity.
type AcquisitionConnectionToken struct {
	ManagerGeneration   uint64
	TransportGeneration uint64
}

// Valid excludes unowned legacy fixtures and adapters without a success token.
func (t AcquisitionConnectionToken) Valid() bool {
	return t.ManagerGeneration != 0 && t.TransportGeneration != 0
}

// CurrentAcquisitionToken reads the current incarnation under the existing gate.
func (mc *ManagedConnection) CurrentAcquisitionToken(ctx context.Context) (AcquisitionConnectionToken, error) {
	if err := mc.LockContext(ctx); err != nil {
		return AcquisitionConnectionToken{}, err
	}
	defer mc.Unlock()
	return mc.acquisitionTokenLocked(), nil
}

// acquisitionTokenLocked is also used at the actual-send boundary while the
// managed gate is already held. Implementations must read only atomic state;
// acquiring an adapter/transport gate here would deadlock the send path.
func (mc *ManagedConnection) acquisitionTokenLocked() AcquisitionConnectionToken {
	token := AcquisitionConnectionToken{ManagerGeneration: mc.managerGeneration}
	if provider, ok := mc.Protocol.(interface{ AcquisitionTransportGeneration() uint64 }); ok {
		token.TransportGeneration = provider.AcquisitionTransportGeneration()
	}
	return token
}
