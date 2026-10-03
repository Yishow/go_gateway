package runtime

import (
	"context"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

// ApplyWorkspaceProjectionForDevices updates runtime state for selected active devices.
// It intentionally reuses the existing per-device projection helper so unrelated
// scheduler devices, point metadata, mappings, and alignment state remain intact.
func (s *Service) ApplyWorkspaceProjectionForDevices(
	ctx context.Context,
	projection *workspace.RuntimeProjection,
	deviceIDs []string,
) error {
	return s.applyWorkspaceProjectionForDevicesAndPoints(ctx, projection, deviceIDs, nil)
}

// ApplyWorkspaceProjectionForDevicesAndPoints updates only selected points for
// the selected devices. Existing points and polling groups outside that group
// member scope remain untouched.
func (s *Service) ApplyWorkspaceProjectionForDevicesAndPoints(
	ctx context.Context,
	projection *workspace.RuntimeProjection,
	deviceIDs []string,
	pointIDs []string,
) error {
	selectedPoints := make(map[string]bool, len(pointIDs))
	for _, rawID := range pointIDs {
		id := strings.TrimSpace(rawID)
		if id == "" {
			return fmt.Errorf("scoped runtime projection point id 不可為空")
		}
		selectedPoints[id] = true
	}
	if len(selectedPoints) == 0 {
		return fmt.Errorf("scoped runtime projection 至少需要一個點位")
	}
	return s.applyWorkspaceProjectionForDevicesAndPoints(ctx, projection, deviceIDs, selectedPoints)
}

func (s *Service) applyWorkspaceProjectionForDevicesAndPoints(
	ctx context.Context,
	projection *workspace.RuntimeProjection,
	deviceIDs []string,
	selectedPoints map[string]bool,
) error {
	if s == nil || s.scheduler == nil {
		return fmt.Errorf("runtime service 不可為 nil")
	}
	if projection == nil {
		return fmt.Errorf("workspace runtime projection 不可為 nil")
	}
	selectedIDs := make([]string, 0, len(deviceIDs))
	selected := make(map[string]bool, len(deviceIDs))
	for _, rawID := range deviceIDs {
		id := strings.TrimSpace(rawID)
		if id == "" {
			return fmt.Errorf("scoped runtime projection device id 不可為空")
		}
		if selected[id] {
			continue
		}
		selected[id] = true
		selectedIDs = append(selectedIDs, id)
	}
	if len(selectedIDs) == 0 {
		return fmt.Errorf("scoped runtime projection 至少需要一台設備")
	}

	projectionIDs := make(map[string]bool, len(projection.DeviceIDs))
	for _, id := range projection.DeviceIDs {
		projectionIDs[strings.TrimSpace(id)] = true
	}
	projectionDevices := make(map[string]*schema.Device, len(projection.Devices))
	for _, deviceRecord := range projection.Devices {
		if deviceRecord == nil {
			continue
		}
		projectionDevices[deviceRecord.ID] = deviceRecord
	}
	for _, id := range selectedIDs {
		if !projectionIDs[id] {
			return fmt.Errorf("device %s is outside workspace runtime projection", id)
		}
		deviceRecord := projectionDevices[id]
		if deviceRecord == nil || deviceRecord.Status != schema.DeviceStatusActive {
			return fmt.Errorf("device %s is not active in workspace runtime projection", id)
		}
	}
	if err := s.validateSharedPollingGroups(projection, selected, selectedPoints); err != nil {
		return err
	}
	if len(selectedPoints) > 0 && !selectedProjectionPointsPresent(projection, selected, selectedPoints) {
		return fmt.Errorf("selected runtime projection point is outside the selected device scope")
	}
	for _, id := range selectedIDs {
		scopedProjection := scopedWorkspaceProjectionForDevice(projection, id, selected, selectedPoints)
		if err := s.applyWorkspaceProjectionForDevice(ctx, scopedProjection, id, len(selectedPoints) > 0); err != nil {
			return fmt.Errorf("apply workspace runtime projection for device %s: %w", id, err)
		}
	}
	return nil
}

func selectedProjectionPointsPresent(projection *workspace.RuntimeProjection, selected, selectedPoints map[string]bool) bool {
	present := make(map[string]bool, len(selectedPoints))
	for _, pointRecord := range projection.Points {
		if pointRecord == nil || !selected[pointRecord.DeviceID] {
			continue
		}
		if selectedPoints[pointRecord.ID] {
			present[pointRecord.ID] = true
		}
	}
	for pointID := range selectedPoints {
		if !present[pointID] {
			return false
		}
	}
	return true
}

func (s *Service) validateSharedPollingGroups(projection *workspace.RuntimeProjection, selected, selectedPoints map[string]bool) error {
	users := make(map[string][2]bool)
	for _, pointRecord := range projection.Points {
		if pointRecord == nil || pointRecord.PollingGroupID == nil {
			continue
		}
		groupID := strings.TrimSpace(*pointRecord.PollingGroupID)
		if groupID == "" {
			continue
		}
		usage := users[groupID]
		isSelected := selected[pointRecord.DeviceID]
		if len(selectedPoints) > 0 {
			isSelected = isSelected && selectedPoints[pointRecord.ID]
		}
		if isSelected {
			usage[0] = true
		} else {
			usage[1] = true
		}
		users[groupID] = usage
	}
	if len(users) == 0 {
		return nil
	}
	desiredGroups := make(map[string]*schema.PollingGroup, len(projection.PollingGroups))
	for _, group := range projection.PollingGroups {
		if group != nil {
			desiredGroups[group.ID] = group
		}
	}
	for groupID, usage := range users {
		if !usage[0] || !usage[1] {
			continue
		}
		desired := desiredGroups[groupID]
		current, ok := s.scheduler.PollingGroup(groupID)
		if !ok || desired == nil || !pollingGroupsEquivalent(current, desired) {
			return fmt.Errorf("shared polling group %s is not safely aligned", groupID)
		}
	}
	return nil
}

func pollingGroupsEquivalent(left, right *schema.PollingGroup) bool {
	if left == nil || right == nil {
		return left == right
	}
	return left.ID == right.ID &&
		left.Name == right.Name &&
		left.Description == right.Description &&
		left.IntervalMs == right.IntervalMs &&
		left.Priority == right.Priority &&
		left.Enabled == right.Enabled
}

func scopedWorkspaceProjectionForDevice(
	projection *workspace.RuntimeProjection,
	deviceID string,
	selected map[string]bool,
	selectedPoints map[string]bool,
) *workspace.RuntimeProjection {
	scoped := *projection
	scoped.Devices = make([]*schema.Device, 0, 1)
	for _, deviceRecord := range projection.Devices {
		if deviceRecord != nil && deviceRecord.ID == deviceID {
			scoped.Devices = append(scoped.Devices, deviceRecord)
			break
		}
	}

	pointIDs := make(map[string]bool)
	groupIDs := make(map[string]bool)
	scoped.Points = make([]*schema.Point, 0)
	for _, pointRecord := range projection.Points {
		if pointRecord == nil || pointRecord.DeviceID != deviceID {
			continue
		}
		if len(selectedPoints) > 0 && !selectedPoints[pointRecord.ID] {
			continue
		}
		scoped.Points = append(scoped.Points, pointRecord)
		pointIDs[pointRecord.ID] = true
		if pointRecord.PollingGroupID != nil && strings.TrimSpace(*pointRecord.PollingGroupID) != "" {
			groupIDs[*pointRecord.PollingGroupID] = true
		}
	}
	for _, pointRecord := range projection.Points {
		if pointRecord == nil || pointRecord.PollingGroupID == nil {
			continue
		}
		if selected[pointRecord.DeviceID] &&
			(len(selectedPoints) == 0 || selectedPoints[pointRecord.ID]) {
			continue
		}
		delete(groupIDs, *pointRecord.PollingGroupID)
	}

	scoped.Mappings = make([]*schema.Mapping, 0)
	tagIDs := make(map[string]bool)
	for _, mappingRecord := range projection.Mappings {
		if mappingRecord == nil || !pointIDs[mappingRecord.PointID] {
			continue
		}
		scoped.Mappings = append(scoped.Mappings, mappingRecord)
		tagIDs[mappingRecord.TagID] = true
	}
	scoped.Tags = make([]*schema.Tag, 0)
	for _, tagRecord := range projection.Tags {
		if tagRecord != nil && tagIDs[tagRecord.ID] {
			scoped.Tags = append(scoped.Tags, tagRecord)
		}
	}
	scoped.PollingGroups = make([]*schema.PollingGroup, 0)
	for _, group := range projection.PollingGroups {
		if group != nil && groupIDs[group.ID] {
			scoped.PollingGroups = append(scoped.PollingGroups, group)
		}
	}
	return &scoped
}
