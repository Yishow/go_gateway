package runtime

import (
	"context"
	"fmt"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

func (s *Service) bootstrapFromWorkspaceProjection(ctx context.Context, projection *workspace.RuntimeProjection) error {
	if err := s.ApplyWorkspaceProjection(ctx, projection); err != nil {
		return err
	}
	if err := s.scheduler.Start(projection.PollingGroups); err != nil { //nolint:contextcheck // Scheduler lifecycle API does not accept context.
		return fmt.Errorf("啟動 scheduler 失敗: %w", err)
	}
	return nil
}

// ApplyWorkspaceProjection replaces runtime device, point, and mapping scope with one persisted workspace projection.
func (s *Service) ApplyWorkspaceProjection(_ context.Context, projection *workspace.RuntimeProjection) error {
	if s == nil || projection == nil {
		return fmt.Errorf("workspace runtime projection 不可為 nil")
	}

	activeDevices := make(map[string]struct{}, len(projection.Devices))
	for _, deviceRecord := range projection.Devices {
		if deviceRecord == nil || deviceRecord.Status != schema.DeviceStatusActive {
			continue
		}
		activeDevices[deviceRecord.ID] = struct{}{}
	}

	for deviceID := range s.scheduler.GetAllDeviceBreakerStates() {
		if _, ok := activeDevices[deviceID]; !ok {
			s.RemoveDevice(deviceID)
		}
	}
	for _, deviceRecord := range projection.Devices {
		if deviceRecord == nil {
			continue
		}
		if deviceRecord.Status != schema.DeviceStatusActive {
			s.RemoveDevice(deviceRecord.ID)
			continue
		}
		s.scheduler.AddDevice(deviceRecord)
	}

	desiredPoints := make(map[string]struct{}, len(projection.Points))
	for _, pointRecord := range projection.Points {
		if pointRecord == nil {
			continue
		}
		desiredPoints[pointRecord.ID] = struct{}{}
	}
	for pointID := range s.currentPointMetaIDs() {
		if _, ok := desiredPoints[pointID]; !ok {
			s.RemovePoint(pointID)
		}
	}
	for _, pointRecord := range projection.Points {
		if pointRecord == nil {
			continue
		}
		s.registerPointMeta(pointRecord.ID, pointMeta{
			DeviceID: pointRecord.DeviceID,
			Address:  pointRecord.Address,
		})
		if _, ok := activeDevices[pointRecord.DeviceID]; pointRecord.Enabled && ok {
			s.scheduler.AddPoint(pointRecord)
			continue
		}
		s.scheduler.RemovePoint(pointRecord.ID)
	}

	if s.IsRunning() {
		for _, group := range projection.PollingGroups {
			s.scheduler.AddPollingGroup(group) //nolint:contextcheck // Scheduler lifecycle API does not accept context.
		}
	}

	if err := s.applyWorkspaceProjectionMappings(projection); err != nil {
		return err
	}
	s.replaceProjectionAligned(projection)
	return nil
}

func (s *Service) currentPointMetaIDs() map[string]struct{} {
	s.pointMetaMu.RLock()
	defer s.pointMetaMu.RUnlock()

	ids := make(map[string]struct{}, len(s.pointMetaIndex))
	for pointID := range s.pointMetaIndex {
		ids[pointID] = struct{}{}
	}
	return ids
}

func (s *Service) currentPointMetaIDsForDevice(deviceID string) map[string]struct{} {
	s.pointMetaMu.RLock()
	defer s.pointMetaMu.RUnlock()

	ids := make(map[string]struct{})
	for pointID, meta := range s.pointMetaIndex {
		if meta.DeviceID == deviceID {
			ids[pointID] = struct{}{}
		}
	}
	return ids
}

func (s *Service) applyWorkspaceProjectionForDevice(
	ctx context.Context,
	projection *workspace.RuntimeProjection,
	deviceID string,
) error {
	if deviceID == "" {
		return s.ApplyWorkspaceProjection(ctx, projection)
	}
	if s == nil || projection == nil {
		return fmt.Errorf("workspace runtime projection 不可為 nil")
	}

	var targetDevice *schema.Device
	for _, deviceRecord := range projection.Devices {
		if deviceRecord != nil && deviceRecord.ID == deviceID {
			targetDevice = deviceRecord
			break
		}
	}
	active := targetDevice != nil && targetDevice.Status == schema.DeviceStatusActive
	if active {
		s.scheduler.AddDevice(targetDevice)
	} else {
		s.RemoveDevice(deviceID)
	}

	desiredPoints := make(map[string]struct{})
	for _, pointRecord := range projection.Points {
		if pointRecord == nil || pointRecord.DeviceID != deviceID {
			continue
		}
		desiredPoints[pointRecord.ID] = struct{}{}
	}
	for pointID := range s.currentPointMetaIDsForDevice(deviceID) {
		if _, ok := desiredPoints[pointID]; !ok {
			s.RemovePoint(pointID)
		}
	}
	for _, pointRecord := range projection.Points {
		if pointRecord == nil || pointRecord.DeviceID != deviceID {
			continue
		}
		s.registerPointMeta(pointRecord.ID, pointMeta{
			DeviceID: pointRecord.DeviceID,
			Address:  pointRecord.Address,
		})
		if active && pointRecord.Enabled {
			s.scheduler.AddPoint(pointRecord)
			continue
		}
		s.scheduler.RemovePoint(pointRecord.ID)
	}

	if s.IsRunning() {
		for _, group := range projection.PollingGroups {
			s.scheduler.AddPollingGroup(group) //nolint:contextcheck // Scheduler lifecycle API does not accept context.
		}
	}
	if err := s.applyWorkspaceProjectionMappingsForPoints(projection, desiredPoints); err != nil {
		return err
	}
	s.markDeviceProjectionAligned(deviceID, projection)
	return nil
}

func (s *Service) applyWorkspaceProjectionMappings(projection *workspace.RuntimeProjection) error {
	next, err := workspaceProjectionMappingBindings(projection, nil)
	if err != nil {
		return err
	}

	s.mappingMu.Lock()
	s.mappingIndex = next
	s.mappingMu.Unlock()
	return nil
}

func (s *Service) applyWorkspaceProjectionMappingsForPoints(
	projection *workspace.RuntimeProjection,
	pointIDs map[string]struct{},
) error {
	next, err := workspaceProjectionMappingBindings(projection, pointIDs)
	if err != nil {
		return err
	}

	s.mappingMu.Lock()
	for pointID := range pointIDs {
		delete(s.mappingIndex, pointID)
	}
	for pointID, bindings := range next {
		s.mappingIndex[pointID] = bindings
	}
	s.mappingMu.Unlock()
	return nil
}

func workspaceProjectionMappingBindings(
	projection *workspace.RuntimeProjection,
	pointIDs map[string]struct{},
) (map[string][]mappingBinding, error) {
	tagTypes := make(map[string]schema.DataType, len(projection.Tags))
	for _, tagRecord := range projection.Tags {
		if tagRecord == nil {
			continue
		}
		tagTypes[tagRecord.ID] = tagRecord.DataType
	}

	next := make(map[string][]mappingBinding)
	for _, mappingRecord := range projection.Mappings {
		if mappingRecord == nil || !mappingRecord.Enabled {
			continue
		}
		if pointIDs != nil {
			if _, ok := pointIDs[mappingRecord.PointID]; !ok {
				continue
			}
		}
		tagType, ok := tagTypes[mappingRecord.TagID]
		if !ok {
			return nil, fmt.Errorf("workspace projection mapping %s references missing tag %s", mappingRecord.ID, mappingRecord.TagID)
		}
		next[mappingRecord.PointID] = append(next[mappingRecord.PointID], mappingBinding{
			TagID:             mappingRecord.TagID,
			TagDataType:       tagType,
			TransformPipeline: mappingRecord.TransformPipeline,
		})
	}
	return next, nil
}
