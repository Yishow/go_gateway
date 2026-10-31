package sourcerule

import (
	"context"
	"fmt"
	"go-gateway/internal/datalink/point"
	"strings"
)

// pointAddressPreflightError is emitted only before a mutation or cleanup.
// A validation error raised after a write must still use normal compensation.
type pointAddressPreflightError struct{ cause error }

func (e *pointAddressPreflightError) Error() string { return e.cause.Error() }
func (e *pointAddressPreflightError) Unwrap() error { return e.cause }

type pointAddressConflictError struct{ address string }

func (e *pointAddressConflictError) Error() string {
	return fmt.Sprintf("位址 %s 已存在既有點位", e.address)
}
func (e *pointAddressConflictError) Unwrap() error { return ErrValidation }

func (s *Service) preflightPointAddresses(ctx context.Context, deviceID string, keepPointIDs map[string]struct{}, addresses, skipped []string) error {
	active := make([]string, 0, len(addresses))
	for _, address := range addresses {
		if !containsAddress(skipped, address) {
			active = append(active, address)
		}
	}
	if len(active) == 0 {
		return nil
	}
	if err := s.validatePointAddresses(ctx, deviceID, keepPointIDs, active); err != nil {
		return &pointAddressPreflightError{cause: err}
	}
	return nil
}

func (s *Service) validatePointAddresses(ctx context.Context, deviceID string, keepPointIDs map[string]struct{}, addresses []string) error {
	filter := point.ListFilter{DeviceID: &deviceID}
	points, err := s.pointSvc.List(ctx, filter)
	if err != nil {
		return fmt.Errorf("取得既有點位失敗: %w", err)
	}

	existing := make(map[string]string, len(points))
	for _, item := range points {
		if _, ok := keepPointIDs[item.ID]; ok {
			continue
		}
		existing[strings.ToUpper(strings.TrimSpace(item.Address))] = item.ID
	}

	for _, address := range addresses {
		if existing[strings.ToUpper(strings.TrimSpace(address))] != "" {
			return &pointAddressConflictError{address: address}
		}
	}
	return nil
}
