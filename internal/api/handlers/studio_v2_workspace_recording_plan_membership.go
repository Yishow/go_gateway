package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

// RecordingPlanMeasurementService reads only measurements in a workspace.
type RecordingPlanMeasurementService interface {
	ListByWorkspace(context.Context, string) ([]measurement.MeasurementDefinition, error)
}

// RecordingPlanPointService reads points using a device-scoped filter.
type RecordingPlanPointService interface {
	List(context.Context, point.ListFilter) ([]*schema.Point, error)
}

var (
	errRecordingPlanMembershipUnavailable = errors.New("recording plan membership services are unavailable")
	errRecordingPlanMembershipInvalid     = errors.New("recording plan membership is invalid")
)

const (
	recordingPlanMembersInvalidCode    = "RECORDING_PLAN_MEMBERS_INVALID"
	recordingPlanMembersInvalidMessage = "recording plan members must be saved measurements in the workspace"
)

type recordingPlanMembershipServices struct {
	measurement RecordingPlanMeasurementService
	point       RecordingPlanPointService
}

// SetRecordingMembershipServices wires the persisted membership readers used
// by recording-plan create and update validation. Missing readers fail closed.
func (h *StudioV2WorkspaceRecordingPlansHandler) SetRecordingMembershipServices(
	measurementSvc RecordingPlanMeasurementService,
	pointSvc RecordingPlanPointService,
) *StudioV2WorkspaceRecordingPlansHandler {
	h.membership = recordingPlanMembershipServices{
		measurement: measurementSvc,
		point:       pointSvc,
	}
	return h
}

// validateRecordingPlanMembership accepts only members backed by a saved
// measurement whose device belongs to the workspace and whose point is still
// enabled. Each member must keep the equipment identity of its measurement.
func (h *StudioV2WorkspaceRecordingPlansHandler) validateRecordingPlanMembership(ctx context.Context, record *workspace.Record, plan *recordingplan.RecordingPlan) error {
	if h.membership.measurement == nil || h.membership.point == nil {
		return errRecordingPlanMembershipUnavailable
	}
	if record == nil || plan == nil || len(plan.Members) == 0 {
		return errRecordingPlanMembershipInvalid
	}

	definitions, err := h.membership.measurement.ListByWorkspace(ctx, record.ID)
	if err != nil {
		return errRecordingPlanMembershipUnavailable
	}
	definitionsByID := make(map[string]measurement.MeasurementDefinition, len(definitions))
	for _, definition := range definitions {
		if definition.WorkspaceID == record.ID {
			definitionsByID[strings.TrimSpace(definition.ID)] = definition
		}
	}
	workspaceDevices := make(map[string]struct{}, len(record.OrderedDeviceIDs))
	for _, deviceID := range record.OrderedDeviceIDs {
		workspaceDevices[strings.TrimSpace(deviceID)] = struct{}{}
	}

	enabledPointsByDevice := make(map[string]map[string]struct{})
	equipmentByMeasurement := make(map[string]string, len(plan.Members))
	memberIDs := make(map[string]struct{}, len(plan.Members))
	for _, member := range plan.Members {
		memberID := strings.TrimSpace(member.MemberID)
		measurementID := strings.TrimSpace(member.MeasurementID)
		if memberID == "" || measurementID == "" {
			return errRecordingPlanMembershipInvalid
		}
		if _, duplicate := memberIDs[memberID]; duplicate {
			return errRecordingPlanMembershipInvalid
		}
		if _, duplicate := equipmentByMeasurement[measurementID]; duplicate {
			return errRecordingPlanMembershipInvalid
		}

		definition, exists := definitionsByID[measurementID]
		if !exists {
			return errRecordingPlanMembershipInvalid
		}
		deviceID := strings.TrimSpace(definition.DeviceID)
		if _, selected := workspaceDevices[deviceID]; !selected || deviceID == "" {
			return errRecordingPlanMembershipInvalid
		}
		equipmentID := strings.TrimSpace(definition.EquipmentID)
		if equipmentID == "" || strings.TrimSpace(member.EquipmentID) != equipmentID {
			return errRecordingPlanMembershipInvalid
		}
		enabledPoints, err := h.enabledRecordingDevicePoints(ctx, deviceID, enabledPointsByDevice)
		if err != nil {
			return err
		}
		if _, enabled := enabledPoints[strings.TrimSpace(definition.PointID)]; !enabled {
			return errRecordingPlanMembershipInvalid
		}

		memberIDs[memberID] = struct{}{}
		equipmentByMeasurement[measurementID] = equipmentID
	}

	for _, stream := range plan.Streams {
		equipmentID, exists := equipmentByMeasurement[strings.TrimSpace(stream.MeasurementID)]
		if !exists {
			return errRecordingPlanMembershipInvalid
		}
		if streamEquipment := strings.TrimSpace(stream.EquipmentID); streamEquipment != "" && streamEquipment != equipmentID {
			return errRecordingPlanMembershipInvalid
		}
	}
	return nil
}

// enabledRecordingDevicePoints loads each referenced device once; devices that
// no member uses are never read.
func (h *StudioV2WorkspaceRecordingPlansHandler) enabledRecordingDevicePoints(ctx context.Context, deviceID string, cache map[string]map[string]struct{}) (map[string]struct{}, error) {
	if points, loaded := cache[deviceID]; loaded {
		return points, nil
	}
	listed, err := h.membership.point.List(ctx, point.ListFilter{DeviceID: &deviceID})
	if err != nil {
		return nil, errRecordingPlanMembershipUnavailable
	}
	points := make(map[string]struct{}, len(listed))
	for _, candidate := range listed {
		if candidate != nil && candidate.DeviceID == deviceID && candidate.Enabled {
			points[strings.TrimSpace(candidate.ID)] = struct{}{}
		}
	}
	cache[deviceID] = points
	return points, nil
}

func renderRecordingPlanMembershipError(c *gin.Context, err error) {
	if errors.Is(err, errRecordingPlanMembershipInvalid) {
		renderRecordingPlanError(c, http.StatusUnprocessableEntity, recordingPlanMembersInvalidCode, recordingPlanMembersInvalidMessage, false, "save each measurement in the workspace, then reload the plan")
		return
	}
	renderRecordingPlanError(c, http.StatusServiceUnavailable, recordingPlanUnavailableCode, recordingPlanUnavailableMessage, true, "retry after saved measurements are available")
}
