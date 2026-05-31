package workspace

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
)

var ErrValidation = errors.New("studio v2 workspace validation failed")

func workspaceValidationError(message string) error {
	return fmt.Errorf("%w: %s", ErrValidation, message)
}

func (s *Service) AttachDevice(ctx context.Context, deviceID string) (*Record, error) {
	deviceID = strings.TrimSpace(deviceID)
	if deviceID == "" {
		return nil, workspaceValidationError("device id 不能為空")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	record, err := s.getOrCreate(ctx)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(record.OrderedDeviceIDs, deviceID) {
		record.OrderedDeviceIDs = append(record.OrderedDeviceIDs, deviceID)
		record.Status = WorkspaceStatusReady
		record.UpdatedAt = s.now()
		if err := s.repo.Save(ctx, record); err != nil {
			return nil, fmt.Errorf("attach workspace device: %w", err)
		}
	}

	return cloneRecord(record), nil
}

func (s *Service) DetachDevice(ctx context.Context, deviceID string) (*Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, err := s.getOrCreate(ctx)
	if err != nil {
		return nil, err
	}

	filtered := make([]string, 0, len(record.OrderedDeviceIDs))
	for _, existingID := range record.OrderedDeviceIDs {
		if existingID != deviceID {
			filtered = append(filtered, existingID)
		}
	}

	if len(filtered) == len(record.OrderedDeviceIDs) {
		// deviceID 不在工作區中，無需變更
		return cloneRecord(record), nil
	}

	record.OrderedDeviceIDs = filtered
	record.Status = workspaceStatusForDevices(filtered)
	record.UpdatedAt = s.now()
	if err := s.repo.Save(ctx, record); err != nil {
		return nil, fmt.Errorf("detach workspace device: %w", err)
	}

	return cloneRecord(record), nil
}

func (s *Service) ReplaceDeviceOrder(ctx context.Context, orderedDeviceIDs []string) (*Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, err := s.getOrCreate(ctx)
	if err != nil {
		return nil, err
	}

	expected := make(map[string]struct{}, len(record.OrderedDeviceIDs))
	for _, deviceID := range record.OrderedDeviceIDs {
		expected[deviceID] = struct{}{}
	}

	if len(orderedDeviceIDs) != len(record.OrderedDeviceIDs) {
		return nil, workspaceValidationError("ordered_device_ids 必須完整覆蓋工作區內的設備")
	}

	seen := make(map[string]struct{}, len(orderedDeviceIDs))
	for _, deviceID := range orderedDeviceIDs {
		deviceID = strings.TrimSpace(deviceID)
		if deviceID == "" {
			return nil, workspaceValidationError("ordered_device_ids 不能包含空值")
		}
		if _, exists := expected[deviceID]; !exists {
			return nil, workspaceValidationError("ordered_device_ids 包含不屬於工作區的設備")
		}
		if _, duplicated := seen[deviceID]; duplicated {
			return nil, workspaceValidationError("ordered_device_ids 不能包含重複設備")
		}
		seen[deviceID] = struct{}{}
	}

	record.OrderedDeviceIDs = append([]string{}, orderedDeviceIDs...)
	record.Status = workspaceStatusForDevices(record.OrderedDeviceIDs)
	record.UpdatedAt = s.now()
	if err := s.repo.Save(ctx, record); err != nil {
		return nil, fmt.Errorf("replace workspace device order: %w", err)
	}

	return cloneRecord(record), nil
}

func workspaceStatusForDevices(deviceIDs []string) Status {
	if len(deviceIDs) == 0 {
		return WorkspaceStatusEmpty
	}

	return WorkspaceStatusReady
}
