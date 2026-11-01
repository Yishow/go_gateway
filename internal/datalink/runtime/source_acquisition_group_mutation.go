package runtime

import (
	"context"
	"fmt"
	"sync"
)

// TryBeginAcquisitionConfigurationMutations is the native atomic polling-group
// deletion hook. Its caller holds the canonical SQL writer reservation, so this
// method never waits for a runtime lock, reads SQL, or invokes a provider.
// Finalization must occur after that reservation ends and never grants gaps.
func (s *Service) TryBeginAcquisitionConfigurationMutations(ctx context.Context, deviceIDs []string) (func(), error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if s == nil || s.scheduler == nil {
		return func() {}, nil
	}
	// Reuse only the exact active synchronous owner, never release its token.
	if owner, ok := ctx.Value(sourceAcquisitionMutationKey{}).(*sourceAcquisitionMutation); ok && owner != nil && owner.owner == s && owner.active.Load() {
		filtered := make([]string, 0, len(deviceIDs))
		for _, id := range deviceIDs {
			if id != owner.deviceID {
				filtered = append(filtered, id)
			}
		}
		deviceIDs = filtered
	}
	tokens, ok := s.scheduler.TryBeginAcquisitionGeometryMutations(ctx, deviceIDs)
	if !ok {
		return nil, fmt.Errorf("source acquisition mutation busy or invalid")
	}
	var once sync.Once
	return func() {
		once.Do(func() {
			for _, token := range tokens {
				s.scheduler.FinishAcquisitionGeometryMutation(token)
			}
		})
	}, nil
}
