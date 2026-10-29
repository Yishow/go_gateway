package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
	"reflect"
	"slices"
	"sort"
	"time"
)

var errStudioLocatorMissingBinding = errors.New("locator binding missing")

type studioLocatorEditorObservation struct {
	TargetGroup *workspace.WriteGroup
	Connector   *schema.DatabaseConnector
	Targets     []*schema.DatabaseTargetMapping
	Devices     []*schema.Device
	Sources     []*studioLocatorObservation
	Missing     map[string]string
}
type studioLocatorEditorProjection struct {
	EditorHydrated  bool                                      `json:"editor_hydrated"`
	DatabaseConfig  *studioV2WorkspaceDatabaseConfigResponse  `json:"database_config"`
	DatabaseTargets []studioV2WorkspaceDatabaseTargetResponse `json:"database_targets"`
}

func (h *StudioV2WorkspaceHandler) readStudioLocatorEntry(ctx context.Context, record *workspace.Record, request studioLocatorRequest) (*studioLocatorObservation, error) {
	o, err := h.readStudioLocator(ctx, record, request.EntityLocator)
	if err != nil {
		return nil, err
	}
	if !request.Editor {
		return o, nil
	}
	if o.Group == nil || o.Group.Group.Revision != request.ExpectedGroupRevision || record.DatabaseSetupRevision != request.ExpectedWorkspaceRevision || o.Group.Group.Status == workspace.WriteGroupStatusDeleted {
		return nil, runtimeEvidenceError(runtimeview.ReasonStale)
	}
	o.Editor, err = h.readStudioLocatorEditor(ctx, record, o.Group.Group)
	return o, err
}
func (h *StudioV2WorkspaceHandler) readStudioLocatorEditor(ctx context.Context, record *workspace.Record, group *workspace.WriteGroup) (*studioLocatorEditorObservation, error) {
	r := h.locatorReaders
	out := &studioLocatorEditorObservation{TargetGroup: group, Targets: []*schema.DatabaseTargetMapping{}, Devices: []*schema.Device{}, Sources: []*studioLocatorObservation{}, Missing: map[string]string{}}
	if record.DatabaseConnectorID != "" {
		if studioLocatorReaderMissing(r.Connectors) || studioLocatorReaderMissing(r.Targets) {
			return nil, studioLocatorReadError(nil)
		}
		connector, err := r.Connectors.GetByID(ctx, record.DatabaseConnectorID)
		if err != nil || connector == nil {
			return nil, studioLocatorReadError(err)
		}
		if connector.ID != record.DatabaseConnectorID || !runtimeDeliveryIDs(connector.IdentityRevision) {
			return nil, runtimeEvidenceError(runtimeview.ReasonStale)
		}
		out.Connector, err = cloneStudioLocatorValue(connector)
		if err != nil {
			return nil, err
		}
	}
	devices := map[string]bool{}
	sources := map[string]bool{}
	for _, member := range group.Members {
		key := member.DeviceID + "/" + member.PointID + "/" + member.TagID
		if !runtimeDeliveryIDs(member.DeviceID, member.PointID, member.TagID) {
			return nil, studioLocatorReadError(nil)
		}
		if !slices.Contains(record.OrderedDeviceIDs, member.DeviceID) {
			out.Missing[key] = "unowned"
			continue
		}
		if studioLocatorReaderMissing(r.Devices) || studioLocatorReaderMissing(r.Points) || studioLocatorReaderMissing(r.Tags) || studioLocatorReaderMissing(r.Mappings) || studioLocatorReaderMissing(r.Rules) {
			return nil, studioLocatorReadError(nil)
		}
		d, err := r.Devices.GetByID(ctx, member.DeviceID)
		if errors.Is(err, sql.ErrNoRows) {
			out.Missing[key] = "device"
			continue
		}
		if err != nil || d == nil {
			return nil, studioLocatorReadError(err)
		}
		if d.ID != member.DeviceID {
			return nil, runtimeEvidenceError(runtimeview.ReasonForeign)
		}
		p, err := r.Points.GetByID(ctx, member.PointID)
		if errors.Is(err, point.ErrPointNotFound) {
			out.Missing[key] = "point"
			continue
		}
		if err != nil || p == nil {
			return nil, studioLocatorReadError(err)
		}
		if p.ID != member.PointID || p.DeviceID != member.DeviceID {
			return nil, runtimeEvidenceError(runtimeview.ReasonStale)
		}
		t, err := r.Tags.GetByID(ctx, member.TagID)
		if errors.Is(err, tag.ErrTagNotFound) {
			out.Missing[key] = "tag"
			continue
		}
		if err != nil || t == nil {
			return nil, studioLocatorReadError(err)
		}
		if t.ID != member.TagID {
			return nil, runtimeEvidenceError(runtimeview.ReasonStale)
		}
		rows, err := r.Mappings.List(ctx, mapping.ListFilter{PointID: &member.PointID, TagID: &member.TagID})
		if err != nil {
			return nil, studioLocatorReadError(err)
		}
		if len(rows) == 0 {
			out.Missing[key] = "mapping"
			continue
		}
		if len(rows) != 1 || rows[0] == nil || rows[0].PointID != member.PointID || rows[0].TagID != member.TagID {
			return nil, runtimeEvidenceError(runtimeview.ReasonStale)
		}
		source := &studioLocatorObservation{}
		err = h.readStudioLocatorMapping(ctx, member.DeviceID, rows[0].ID, source)
		if errors.Is(err, errStudioLocatorMissingBinding) {
			out.Missing[key] = "binding"
			continue
		}
		if err != nil {
			return nil, err
		}
		if source.Point.ID != member.PointID || source.Tag.ID != member.TagID {
			return nil, runtimeEvidenceError(runtimeview.ReasonStale)
		}
		if !devices[d.ID] {
			copy, err := cloneStudioLocatorValue(d)
			if err != nil {
				return nil, err
			}
			out.Devices = append(out.Devices, copy)
			devices[d.ID] = true
		}
		if !sources[source.Mapping.ID] {
			out.Sources = append(out.Sources, source)
			sources[source.Mapping.ID] = true
		}
	}
	sort.Slice(out.Devices, func(i, j int) bool { return out.Devices[i].ID < out.Devices[j].ID })
	sort.Slice(out.Sources, func(i, j int) bool { return out.Sources[i].Mapping.ID < out.Sources[j].Mapping.ID })
	if out.Connector != nil && group.Destination.ConnectorID == record.DatabaseConnectorID && len(group.Migration.SourceIDs) > 0 {
		rows, err := r.Targets.List(ctx, dbtarget.TargetMappingListFilter{ConnectorID: &record.DatabaseConnectorID})
		if err != nil {
			return nil, studioLocatorReadError(err)
		}
		canonicalSources := map[string]*workspace.WriteGroup{}
		for _, id := range group.Migration.SourceIDs {
			if !runtimeDeliveryIDs(id) || canonicalSources[id] != nil {
				return nil, studioLocatorReadError(nil)
			}
			canonicalSources[id] = group
		}
		if err := validateMigrationReadRows(rows, canonicalSources, dbtarget.TargetMappingListFilter{ConnectorID: &record.DatabaseConnectorID}); err != nil {
			return nil, studioLocatorReadError(err)
		}
		tags := map[string]string{}
		for _, source := range out.Sources {
			tags[source.Tag.ID] = source.Point.ID
		}
		for _, row := range rows {
			if row == nil {
				return nil, studioLocatorReadError(nil)
			}
			if row.ConnectorID != record.DatabaseConnectorID {
				return nil, runtimeEvidenceError(runtimeview.ReasonForeign)
			}
			if !slices.Contains(group.Migration.SourceIDs, row.ID) {
				continue
			}
			projected, err := projectMigrationRead(row, map[string]*workspace.WriteGroup{row.ID: group})
			if err != nil {
				return nil, studioLocatorReadError(err)
			}
			if projected.DatabaseTargetMapping.ConnectorID != record.DatabaseConnectorID || tags[projected.DatabaseTargetMapping.TagID] != projected.PointID {
				continue
			}
			copy, err := cloneStudioLocatorValue(row)
			if err != nil {
				return nil, err
			}
			out.Targets = append(out.Targets, copy)
		}
		sort.Slice(out.Targets, func(i, j int) bool { return out.Targets[i].ID < out.Targets[j].ID })
	}
	return out, nil
}
func projectStudioLocatorEditor(record *workspace.Record, o *studioLocatorEditorObservation, entry *studioLocatorEntry) error {
	entry.studioLocatorEditorProjection = &studioLocatorEditorProjection{EditorHydrated: true, DatabaseTargets: []studioV2WorkspaceDatabaseTargetResponse{}}
	if o.Connector != nil {
		c := buildWorkspaceDatabaseConfigResponse(record.ID, o.Connector, record.DatabaseRowGroups)
		c.SetupRevision = record.DatabaseSetupRevision
		c.LastSchemaEnsureError, c.LastWriteError, c.LastFlushError = "", "", ""
		entry.DatabaseConfig = &c
	}
	for _, d := range o.Devices {
		if slices.Contains(entry.HydrationScope.DeviceIDs, d.ID) {
			continue
		}
		item := mapStudioV2WorkspaceDeviceResponse(d, "", "")
		item.LastTestError, item.AvailabilityReason = "", ""
		entry.Devices = append(entry.Devices, item)
		entry.HydrationScope.DeviceIDs = append(entry.HydrationScope.DeviceIDs, d.ID)
	}

	for _, source := range o.Sources {
		projected, err := projectStudioLocator(record, source)
		if err != nil {
			return err
		}
		if !slices.Contains(entry.HydrationScope.RuleIDs, source.Rule.ID) {
			entry.SourceRules = append(entry.SourceRules, projected.SourceRules...)
			entry.HydrationScope.RuleIDs = append(entry.HydrationScope.RuleIDs, source.Rule.ID)
		}
		entry.Mappings = append(entry.Mappings, projected.Mappings...)
		entry.HydrationScope.MappingIDs = append(entry.HydrationScope.MappingIDs, source.Mapping.ID)

	}
	for _, raw := range o.Targets {
		projected, err := projectMigrationRead(raw, map[string]*workspace.WriteGroup{raw.ID: o.TargetGroup})
		if err != nil {
			return studioLocatorReadError(err)
		}
		row := projected.DatabaseTargetMapping
		entry.DatabaseTargets = append(entry.DatabaseTargets, studioV2WorkspaceDatabaseTargetResponse{ID: row.ID, WorkspaceID: record.ID, PointID: projected.PointID, TagID: row.TagID, ColumnName: row.ColumnName, Enabled: row.Enabled, RowGroupID: projected.RowGroupID, SaveState: workspaceSaveStateSaved, SetupRevision: record.DatabaseSetupRevision, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt})
	}
	return nil
}

// Empty editor collections are observed partial collections, not omitted unknowns.
func (e studioLocatorEntry) MarshalJSON() ([]byte, error) {
	type alias studioLocatorEntry
	if e.studioLocatorEditorProjection == nil {
		return json.Marshal(alias(e))
	}
	devices, rules, mappings := e.Devices, e.SourceRules, e.Mappings
	if devices == nil {
		devices = []studioV2WorkspaceDeviceResponse{}
	}
	if rules == nil {
		rules = []sourceRuleResponse{}
	}
	if mappings == nil {
		mappings = []studioV2WorkspaceMappingResponse{}
	}
	return json.Marshal(struct {
		alias
		Devices     []studioV2WorkspaceDeviceResponse  `json:"devices"`
		SourceRules []sourceRuleResponse               `json:"source_rules"`
		Mappings    []studioV2WorkspaceMappingResponse `json:"mappings"`
	}{alias(e), devices, rules, mappings})
}
func sameStudioLocatorEditor(a, b *studioLocatorEditorObservation) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	if len(a.Sources) != len(b.Sources) || len(a.Devices) != len(b.Devices) || !reflect.DeepEqual(a.Missing, b.Missing) || !reflect.DeepEqual(a.Targets, b.Targets) || !reflect.DeepEqual(a.TargetGroup, b.TargetGroup) {
		return false
	}
	for i := range a.Sources {
		if !sameStudioLocatorObservation(a.Sources[i], b.Sources[i]) {
			return false
		}
	}
	for i := range a.Devices {
		if !sameStudioLocatorObservation(&studioLocatorObservation{Device: a.Devices[i]}, &studioLocatorObservation{Device: b.Devices[i]}) {
			return false
		}
	}
	normalize := func(c *schema.DatabaseConnector) *schema.DatabaseConnector {
		if c == nil {
			return nil
		}
		v := *c
		v.Status = ""
		v.LastCheckAt = nil
		v.LastCheckError = ""
		v.LastSchemaEnsureAt = nil
		v.LastSchemaEnsureStatus = ""
		v.LastSchemaEnsureError = ""
		v.LastWriteAt = nil
		v.LastWriteStatus = ""
		v.LastWriteError = ""
		v.LastFlushAt = nil
		v.LastFlushStatus = ""
		v.LastFlushError = ""
		v.UpdatedAt = time.Time{}
		return &v
	}
	return reflect.DeepEqual(normalize(a.Connector), normalize(b.Connector))
}
