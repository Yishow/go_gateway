package sourcerule

import (
	"context"
	"fmt"
)

// sourceProjectionRollback retains only transient affected-key state from the
// existing runtime authority; it does not introduce a second running model.
type sourceProjectionRollback interface {
	SnapshotSourceRuleProjection(string, []string) any
	RestoreSourceRuleProjection(any, *RuntimeProjectionProof) error
}

func (s *Service) snapshotRuntimeProjection(snapshot *sourceRuleSnapshot) {
	if runtime, ok := s.runtimeSync.(sourceProjectionRollback); ok {
		snapshot.runtimeProjection = runtime.SnapshotSourceRuleProjection(snapshot.rule.DeviceID, snapshotPointIDs(snapshot))
	}
}

func (s *Service) restoreRuntimeProjection(ctx context.Context, snapshot *sourceRuleSnapshot) error {
	if snapshot.runtimeProjection == nil {
		return nil
	}
	return s.restoreCapturedRuntimeProjection(ctx, snapshot.runtimeProjection)
}

func (s *Service) snapshotEmptyRuntimeProjection(deviceID string) any {
	if runtime, ok := s.runtimeSync.(sourceProjectionRollback); ok {
		return runtime.SnapshotSourceRuleProjection(deviceID, nil)
	}
	return nil
}

func (s *Service) restoreCapturedRuntimeProjection(ctx context.Context, value any) error {
	if value == nil {
		return nil
	}
	if runtime, ok := s.runtimeSync.(sourceProjectionRollback); ok {
		return runtime.RestoreSourceRuleProjection(value, s.sourceProjectionProof(ctx))
	}
	return fmt.Errorf("source projection rollback authority unavailable")
}
