package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
)

// StudioV2LocatorReaders exposes existing saved reads only, never bootstrap or repair.
type StudioV2LocatorReaders struct {
	Connectors interface {
		GetByID(context.Context, string) (*schema.DatabaseConnector, error)
	}
	Targets interface {
		List(context.Context, dbtarget.TargetMappingListFilter) ([]*schema.DatabaseTargetMapping, error)
	}
	Devices interface {
		GetByID(context.Context, string) (*schema.Device, error)
	}
	Rules interface {
		GetByID(context.Context, string) (*schema.SourceRule, error)
		ListByDeviceIDs(context.Context, []string) ([]*schema.SourceRule, error)
		ListLinks(context.Context, string) ([]*schema.SourceRuleLink, error)
	}
	Points interface {
		GetByID(context.Context, string) (*schema.Point, error)
	}
	Mappings interface {
		GetByID(context.Context, string) (*schema.Mapping, error)
		List(context.Context, mapping.ListFilter) ([]*schema.Mapping, error)
	}
	Tags interface {
		GetByID(context.Context, string) (*schema.Tag, error)
	}
	Groups interface {
		Get(context.Context, string) (*workspace.WriteGroupSaveResult, error)
	}
}

// WithLocatorEntryServices wires saved-read authorities for an explicit viewing entry.
func (h *StudioV2WorkspaceHandler) WithLocatorEntryServices(readers StudioV2LocatorReaders) *StudioV2WorkspaceHandler {
	h.locatorReaders = readers
	return h
}

type studioLocatorRequest struct {
	runtimeview.EntityLocator
	Editor                                           bool
	ExpectedWorkspaceRevision, ExpectedGroupRevision string
}

func parseStudioLocatorQuery(raw string) (*studioLocatorRequest, error) {
	if raw == "" {
		return nil, nil
	}
	q, err := url.ParseQuery(raw)
	allowed := map[string]bool{"entry": true, "workspace_id": true, "step": true, "device_id": true, "rule_id": true, "mapping_id": true, "group_id": true, "group_revision": true, "editor": true, "expected_workspace_revision": true, "expected_group_revision": true}
	if err != nil || q.Get("entry") != "locator" {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	for key, values := range q {
		if !allowed[key] || len(values) != 1 || (values[0] == "" && key != "expected_workspace_revision") || (key != "entry" && key != "step" && key != "editor" && !(key == "expected_workspace_revision" && values[0] == "") && !runtimeDeliveryIDs(values[0])) {
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
	}
	step, err := strconv.Atoi(q.Get("step"))
	if err != nil || q.Get("step") != strconv.Itoa(step) || !runtimeDeliveryIDs(q.Get("workspace_id")) {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	required := []string{"entry", "workspace_id", "step"}
	switch step {
	case 1:
		required = append(required, "device_id")
	case 2:
		required = append(required, "device_id", "rule_id")
	case 3:
		required = append(required, "device_id", "mapping_id")
	case 4:
		required = append(required, "group_id")
	default:
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	editor := q.Get("editor") == "1"
	if _, present := q["editor"]; present && (!editor || step != 4) {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	if editor {
		if _, present := q["expected_workspace_revision"]; !present {
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
		required = append(required, "editor", "expected_group_revision")
	}
	for _, key := range required {
		if q.Get(key) == "" {
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
	}
	for key := range q {
		if !slices.Contains(required, key) && !(step == 4 && (key == "device_id" || key == "group_revision" || (editor && key == "expected_workspace_revision"))) {
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
	}
	return &studioLocatorRequest{Editor: editor, ExpectedWorkspaceRevision: q.Get("expected_workspace_revision"), ExpectedGroupRevision: q.Get("expected_group_revision"), EntityLocator: runtimeview.EntityLocator{Step: step, WorkspaceID: q.Get("workspace_id"), DeviceID: runtimeEvidenceString(q.Get("device_id")), RuleID: runtimeEvidenceString(q.Get("rule_id")), MappingID: runtimeEvidenceString(q.Get("mapping_id")), GroupID: runtimeEvidenceString(q.Get("group_id")), GroupRevision: runtimeEvidenceString(q.Get("group_revision"))}}, nil
}

// locatorEntry handles every query-bearing request, so malformed scope cannot fall back to bootstrap.
func (h *StudioV2WorkspaceHandler) locatorEntry(c *gin.Context) bool {
	request, err := parseStudioLocatorQuery(c.Request.URL.RawQuery)
	if err != nil {
		renderRuntimeHandlerError(c, http.StatusBadRequest, string(runtimeview.ReasonUnavailable), false)
		return true
	}
	if request == nil {
		return false
	}
	record, err := h.svc.ReadSnapshot(c.Request.Context())
	if err != nil || record == nil {
		renderRuntimeEvidenceError(c, studioLocatorReadError(err))
		return true
	}
	if record.ID != request.WorkspaceID {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonForeign))
		return true
	}
	if !runtimeDeliveryIDs(record.ID) || (record.DatabaseSetupRevision != "" && !runtimeDeliveryIDs(record.DatabaseSetupRevision)) {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonUnavailable))
		return true
	}
	if request.Editor && request.ExpectedWorkspaceRevision != record.DatabaseSetupRevision {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonStale))
		return true
	}
	first, err := h.readStudioLocatorEntry(c.Request.Context(), record, *request)
	if err != nil {
		renderRuntimeEvidenceError(c, err)
		return true
	}
	entry, err := projectStudioLocator(record, first)
	if err != nil {
		renderRuntimeEvidenceError(c, err)
		return true
	}
	second, err := h.readStudioLocatorEntry(c.Request.Context(), record, *request)
	if err != nil {
		renderRuntimeEvidenceError(c, err)
		return true
	}
	if !sameStudioLocatorObservation(first, second) {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonStale))
		return true
	}
	after, err := h.svc.ReadSnapshot(c.Request.Context())
	if err != nil || after == nil {
		renderRuntimeEvidenceError(c, studioLocatorReadError(err))
		return true
	}
	if !sameStudioLocatorRecord(record, after) {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonStale))
		return true
	}
	// Heads come from the independent second read, not a cached native DTO.
	if second.Editor != nil {
		for _, source := range second.Editor.Sources {
			entry.RuleHeads[source.Rule.ID] = source.Rule.RevisionID
			entry.RevalidatedMappingIDs = append(entry.RevalidatedMappingIDs, source.Mapping.ID)
		}
	}
	if second.Rule != nil {
		entry.RuleHeads[second.Rule.ID] = second.Rule.RevisionID
	}
	if second.Group != nil {
		entry.GroupHeads[second.Group.Group.ID] = second.Group.Group.Revision
	}
	if second.Mapping != nil {
		entry.RevalidatedMappingIDs = append(entry.RevalidatedMappingIDs, second.Mapping.ID)
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: studioV2WorkspaceBootstrapResponse{ID: record.ID, Kind: record.Kind, Status: record.Status, OrderedDeviceIDs: append([]string{}, record.OrderedDeviceIDs...), CreatedAt: record.CreatedAt, UpdatedAt: record.UpdatedAt, LocatorEntry: entry}})
	return true
}

func studioLocatorReadError(err error) error {
	if errors.Is(err, workspace.ErrNotFound) || errors.Is(err, workspace.ErrWriteGroupNotFound) || errors.Is(err, sql.ErrNoRows) || errors.Is(err, sourcerule.ErrSourceRuleNotFound) || errors.Is(err, point.ErrPointNotFound) || errors.Is(err, mapping.ErrMappingNotFound) || errors.Is(err, tag.ErrTagNotFound) {
		return runtimeEvidenceError(runtimeview.ReasonStale)
	}
	return runtimeEvidenceError(runtimeview.ReasonUnavailable)
}

type studioLocatorObservation struct {
	Editor  *studioLocatorEditorObservation
	Device  *schema.Device
	Rule    *schema.SourceRule
	Link    *schema.SourceRuleLink
	Point   *schema.Point
	Mapping *schema.Mapping
	Tag     *schema.Tag
	Group   *workspace.WriteGroupSaveResult
}

func (h *StudioV2WorkspaceHandler) readStudioLocator(ctx context.Context, record *workspace.Record, request runtimeview.EntityLocator) (*studioLocatorObservation, error) {
	o := &studioLocatorObservation{}
	r := h.locatorReaders
	if request.DeviceID != nil {
		id := *request.DeviceID
		count := 0
		for _, owned := range record.OrderedDeviceIDs {
			if owned == id {
				count++
			}
		}
		if count != 1 {
			return nil, runtimeEvidenceError(runtimeview.ReasonForeign)
		}
		if studioLocatorReaderMissing(r.Devices) {
			return nil, studioLocatorReadError(nil)
		}
		d, err := r.Devices.GetByID(ctx, id)
		if err != nil || d == nil {
			return nil, studioLocatorReadError(err)
		}
		if d.ID != id {
			return nil, runtimeEvidenceError(runtimeview.ReasonForeign)
		}
		o.Device, err = cloneStudioLocatorValue(d)
		if err != nil {
			return nil, err
		}
	}
	switch request.Step {
	case 2:
		if studioLocatorReaderMissing(r.Rules) {
			return nil, studioLocatorReadError(nil)
		}
		rule, err := r.Rules.GetByID(ctx, *request.RuleID)
		if err != nil || rule == nil {
			return nil, studioLocatorReadError(err)
		}
		if rule.ID != *request.RuleID || rule.DeviceID != *request.DeviceID {
			return nil, runtimeEvidenceError(runtimeview.ReasonForeign)
		}
		if !runtimeDeliveryIDs(rule.ID, rule.RevisionID) {
			return nil, studioLocatorReadError(nil)
		}
		o.Rule, err = cloneStudioLocatorValue(rule)
		if err != nil {
			return nil, err
		}
	case 3:
		if err := h.readStudioLocatorMapping(ctx, *request.DeviceID, *request.MappingID, o); err != nil {
			if errors.Is(err, errStudioLocatorMissingBinding) {
				return nil, runtimeEvidenceError(runtimeview.ReasonStale)
			}
			return nil, err
		}
	case 4:
		if studioLocatorReaderMissing(r.Groups) {
			return nil, studioLocatorReadError(nil)
		}
		group, err := r.Groups.Get(ctx, *request.GroupID)
		if err != nil || group == nil || group.Group == nil {
			return nil, studioLocatorReadError(err)
		}
		if group.Group.ID != *request.GroupID || group.Group.WorkspaceID != record.ID {
			return nil, runtimeEvidenceError(runtimeview.ReasonForeign)
		}
		if group.WorkspaceRevision != record.DatabaseSetupRevision {
			return nil, runtimeEvidenceError(runtimeview.ReasonStale)
		}
		if !runtimeDeliveryIDs(group.Group.Revision) {
			return nil, studioLocatorReadError(nil)
		}
		if request.DeviceID != nil {
			owner := group.Group.BasicManagedDeviceID
			if (owner != "" && owner != *request.DeviceID) || (owner != *request.DeviceID && !runtimeGroupHasDevice(group.Group, *request.DeviceID)) {
				return nil, runtimeEvidenceError(runtimeview.ReasonForeign)
			}
		}
		// Deleted is retained as a canonical tombstone; the selection resolver refuses it.
		o.Group, err = cloneStudioLocatorValue(group)
		if err != nil {
			return nil, err
		}
	}
	return o, nil
}

func (h *StudioV2WorkspaceHandler) readStudioLocatorMapping(ctx context.Context, deviceID, mappingID string, o *studioLocatorObservation) error {
	r := h.locatorReaders
	if studioLocatorReaderMissing(r.Mappings) || studioLocatorReaderMissing(r.Points) || studioLocatorReaderMissing(r.Tags) || studioLocatorReaderMissing(r.Rules) {
		return studioLocatorReadError(nil)
	}
	m, err := r.Mappings.GetByID(ctx, mappingID)
	if err != nil || m == nil {
		return studioLocatorReadError(err)
	}
	if m.ID != mappingID || !runtimeDeliveryIDs(m.PointID, m.TagID) {
		return runtimeEvidenceError(runtimeview.ReasonStale)
	}
	o.Mapping, err = cloneStudioLocatorValue(m)
	if err != nil {
		return err
	}
	p, err := r.Points.GetByID(ctx, o.Mapping.PointID)
	if err != nil || p == nil {
		return studioLocatorReadError(err)
	}
	if p.ID != o.Mapping.PointID || p.DeviceID != deviceID {
		return runtimeEvidenceError(runtimeview.ReasonForeign)
	}
	o.Point, err = cloneStudioLocatorValue(p)
	if err != nil {
		return err
	}
	t, err := r.Tags.GetByID(ctx, o.Mapping.TagID)
	if err != nil || t == nil {
		return studioLocatorReadError(err)
	}
	if t.ID != o.Mapping.TagID {
		return runtimeEvidenceError(runtimeview.ReasonStale)
	}
	o.Tag, err = cloneStudioLocatorValue(t)
	if err != nil {
		return err
	}
	rules, err := r.Rules.ListByDeviceIDs(ctx, []string{deviceID})
	if err != nil {
		return studioLocatorReadError(err)
	}
	for _, candidate := range rules {
		if candidate == nil || !runtimeDeliveryIDs(candidate.ID, candidate.RevisionID) {
			return studioLocatorReadError(nil)
		}
		rule, err := cloneStudioLocatorValue(candidate)
		if err != nil {
			return err
		}
		links, err := r.Rules.ListLinks(ctx, rule.ID)
		if err != nil {
			return studioLocatorReadError(err)
		}
		for _, link := range links {
			if link == nil || link.PointID != o.Point.ID {
				continue
			}
			if o.Rule != nil || link.RuleID != rule.ID || link.Address != o.Point.Address || link.MappingID == nil || *link.MappingID != mappingID || link.TagID == nil || *link.TagID != o.Tag.ID {
				return runtimeEvidenceError(runtimeview.ReasonStale)
			}
			if rule.DeviceID != deviceID {
				return runtimeEvidenceError(runtimeview.ReasonForeign)
			}
			if !isRuleManagedWorkspaceTagOwnedBy(o.Tag, rule.ID, link.Address) && !isLegacyDirectWorkspaceBinding(link, o.Mapping, o.Tag) {
				return runtimeEvidenceError(runtimeview.ReasonStale)
			}
			o.Rule = rule
			o.Link, err = cloneStudioLocatorValue(link)
			if err != nil {
				return err
			}
		}
	}
	if o.Rule == nil {
		return errStudioLocatorMissingBinding
	}
	head, err := r.Rules.GetByID(ctx, o.Rule.ID)
	if err != nil || head == nil {
		return studioLocatorReadError(err)
	}
	if !sameStudioLocatorRecord(o.Rule, head) {
		return runtimeEvidenceError(runtimeview.ReasonStale)
	}
	return nil
}

func studioLocatorReaderMissing(value any) bool {
	return value == nil || (reflect.ValueOf(value).Kind() == reflect.Pointer && reflect.ValueOf(value).IsNil())
}
