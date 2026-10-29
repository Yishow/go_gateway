package handlers

import (
	"encoding/json"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"slices"
	"strings"
)

// MarshalJSON is the shared HTTP serialization boundary for both runtime status paths.
// Copies retain runtime evidence privately; only safe operator projections leave the handler.
func (response runtimeStatusResponse) MarshalJSON() ([]byte, error) {
	type wireResponse runtimeStatusResponse
	safe := wireResponse(response)
	if runtimeSafeEnum(safe.SnapshotState.State, "ready", "empty", "degraded", "unavailable", "stale") == "" {
		safe.SnapshotState = datalinkruntime.RuntimeUnavailableTruthState("runtime snapshot unavailable")
	}
	safe.SnapshotState.Reason = runtimeSafeDiagnostic(safe.SnapshotState.Reason)
	safe.Collectors = slices.Clone(response.Collectors)
	for i := range safe.Collectors {
		safe.Collectors[i] = safeRuntimeCollector(safe.Collectors[i])
	}
	safe.DatabaseDelivery = slices.Clone(response.DatabaseDelivery)
	for i := range safe.DatabaseDelivery {
		d := &safe.DatabaseDelivery[i]
		d.DeviceID = runtimeSafeID(d.DeviceID)
		d.PointID = runtimeSafeID(d.PointID)
		d.TagID = runtimeSafeID(d.TagID)
		d.Status = datalinkruntime.DatabaseDeliveryStatus(runtimeSafeEnum(string(d.Status), "succeeded", "failed"))
		d.FailedStage = safeRuntimeDatabaseStage(d.FailedStage)
		d.Stages = slices.Clone(d.Stages)
		for j := range d.Stages {
			d.Stages[j] = safeRuntimeDatabaseStage(d.Stages[j])
		}
		d.Error = runtimeSafeDiagnostic(d.Error)
		d.LastFailureReason = runtimeSafeDiagnostic(d.LastFailureReason)
	}
	safe.ModbusShareDelivery = slices.Clone(response.ModbusShareDelivery)
	for i := range safe.ModbusShareDelivery {
		d := &safe.ModbusShareDelivery[i]
		d.DeviceID = runtimeSafeID(d.DeviceID)
		d.PointID = runtimeSafeID(d.PointID)
		d.TagID = runtimeSafeID(d.TagID)
		d.Status = datalinkruntime.ModbusShareDeliveryStatus(runtimeSafeEnum(string(d.Status), "succeeded", "failed"))
		d.Stage = runtimeSafeID(d.Stage)
		d.Error = runtimeSafeDiagnostic(d.Error)
	}
	safe.Diagnostics = slices.Clone(response.Diagnostics)
	for i := range safe.Diagnostics {
		d := &safe.Diagnostics[i]
		d.Scope = runtimeSafeDiagnosticScope(d.Scope)
		d.DeviceID = runtimeSafeID(d.DeviceID)
		d.PointID = runtimeSafeID(d.PointID)
		d.TagID = runtimeSafeID(d.TagID)
		d.FailureStage = safeRuntimeFlowStage(d.FailureStage)
		d.LatestSuccessfulStage = safeRuntimeFlowStage(d.LatestSuccessfulStage)
		d.FailureReason = runtimeSafeDiagnostic(d.FailureReason)
		d.Stages = slices.Clone(d.Stages)
		for j := range d.Stages {
			d.Stages[j].Stage = safeRuntimeFlowStage(d.Stages[j].Stage)
			d.Stages[j].Status = datalinkruntime.RuntimeFlowStageStatus(runtimeSafeEnum(string(d.Stages[j].Status), "success", "failed", "unknown"))
			d.Stages[j].Reason = runtimeSafeDiagnostic(d.Stages[j].Reason)
		}
	}
	return json.Marshal(safe)
}

func safeRuntimeCollector(c runtimeCollectorResponse) runtimeCollectorResponse {
	c.DeviceID = runtimeSafeID(c.DeviceID)
	c.DeviceName = runtimeSafeLabel(c.DeviceName)
	c.Protocol = runtimeSafeStatusChoice(c.Protocol, "modbus_tcp", "modbus_rtu", "modbus_udp", "fatek_fbs", "mc_3e", "mqtt")
	c.Status = runtimeSafeStatusChoice(c.Status, "idle", "running", "warning", "error")
	c.AvailabilityStatus = runtimeSafeStatusChoice(c.AvailabilityStatus, "available", "unavailable")
	if c.Status == "unavailable" || c.AvailabilityStatus == "unavailable" || c.Protocol == "unavailable" || c.DeviceID == "" {
		c.Running = false
	}
	c.AvailabilityReason = runtimeSafeAvailabilityReason(c.AvailabilityReason)
	c.LastError = runtimeSafeDiagnosticPointer(c.LastError)
	c.BreakerState = runtimeSafeStatusChoice(c.BreakerState, "open", "half-open", "closed")
	c.ProjectionAlignment = runtimeSafeStatusChoice(c.ProjectionAlignment, "aligned", "deferred", "restart-required", "restart_required", "stale", "degraded")
	c.RuntimeProjectionVersion = runtimeSafeID(c.RuntimeProjectionVersion)
	c.WorkspaceProjectionVersion = runtimeSafeID(c.WorkspaceProjectionVersion)
	c.ProjectionMessage = runtimeSafeDiagnostic(c.ProjectionMessage)
	c.ProjectionCode = runtimeSafeID(c.ProjectionCode)
	return c
}

func safeRuntimeStatusEvent(event datalinkruntime.DeviceStatusEvent) runtimeCollectorResponse {
	return safeRuntimeCollector(mapRuntimeCollectorResponse(event))
}

func safeRuntimeDatabaseStage(value datalinkruntime.DatabaseDeliveryStage) datalinkruntime.DatabaseDeliveryStage {
	return datalinkruntime.DatabaseDeliveryStage(runtimeSafeEnum(string(value), "collected", "mapped", "db_write", "db_write_failed", "db_write_succeeded"))
}

func safeRuntimeFlowStage(value datalinkruntime.RuntimeFlowStage) datalinkruntime.RuntimeFlowStage {
	return datalinkruntime.RuntimeFlowStage(runtimeSafeEnum(string(value), "collector", "mapping", "runtime_projection", "database_delivery"))
}

func runtimeSafeStatusChoice(value string, allowed ...string) string {
	if safe := runtimeSafeEnum(value, allowed...); safe != "" {
		return safe
	}
	return "unavailable"
}

func runtimeSafeDiagnosticScope(value string) string {
	if value == "runtime" {
		return value
	}
	kind, id, found := strings.Cut(value, ":")
	if found && runtimeSafeEnum(kind, "device", "point", "tag") != "" && runtimeSafeID(id) != "" {
		return kind + ":" + id
	}
	return ""
}
