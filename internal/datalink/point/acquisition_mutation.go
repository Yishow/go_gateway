package point

import "context"

// AcquisitionMutationLifecycle invalidates derived gap authorization before a
// configuration write. It cannot grant a new source graph or alter sample IDs.
// Implementations reuse an active outer source operation in the same context.
type AcquisitionMutationLifecycle interface {
	BeginAcquisitionConfigurationMutation(context.Context, string) func()
}

// SetAcquisitionMutationLifecycle configures the native runtime seam at startup.
// Telemetry writes deliberately never invoke this lifecycle.
func (s *Service) SetAcquisitionMutationLifecycle(lifecycle AcquisitionMutationLifecycle) {
	s.acquisitionMutation = lifecycle
}

func (s *Service) beginAcquisitionMutation(ctx context.Context, deviceID string) func() {
	if s.acquisitionMutation != nil {
		return s.acquisitionMutation.BeginAcquisitionConfigurationMutation(ctx, deviceID)
	}
	return func() {}
}
