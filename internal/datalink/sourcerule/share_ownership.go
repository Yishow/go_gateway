package sourcerule

import (
	"context"
	"reflect"
	"strings"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/schema"
)

type shareTagReader interface {
	GetByID(context.Context, string) (*schema.Tag, error)
}

type sharePointReader interface {
	GetByID(context.Context, string) (*schema.Point, error)
}

// NewPreviewMappingScopeChecker proves that a persisted mapping belongs to the
// requested workspace and returns its owning device for a scoped runtime
// subscription. A mapping id alone is not sufficient to authorize preview.
func NewPreviewMappingScopeChecker(workspaceReader interface {
	ShareOwnershipSnapshot(context.Context) (string, []string, error)
}, rules *Service, points sharePointReader) func(context.Context, string, *schema.Mapping) (string, error) {
	return func(ctx context.Context, workspaceID string, record *schema.Mapping) (string, error) {
		if workspaceReader == nil || rules == nil || points == nil || record == nil || workspaceID == "" {
			return "", modbusshare.NewError(modbusshare.ErrCodeWorkspaceScope, "preview mapping ownership could not be verified", true)
		}
		resolvedWorkspaceID, deviceIDs, err := workspaceReader.ShareOwnershipSnapshot(ctx)
		if err != nil || resolvedWorkspaceID != workspaceID {
			return "", modbusshare.NewError(modbusshare.ErrCodeWorkspaceScope, "preview mapping is outside the requested workspace", false)
		}
		pointRecord, err := points.GetByID(ctx, record.PointID)
		if err != nil || pointRecord == nil || !containsString(deviceIDs, pointRecord.DeviceID) {
			return "", modbusshare.NewError(modbusshare.ErrCodeWorkspaceScope, "preview point is outside the requested workspace", false)
		}
		ruleRecords, err := rules.ListByDeviceIDs(ctx, []string{pointRecord.DeviceID})
		if err != nil {
			return "", modbusshare.NewError(modbusshare.ErrCodeWorkspaceScope, "preview mapping ownership could not be verified", true)
		}
		for _, rule := range ruleRecords {
			if rule == nil {
				continue
			}
			links, err := rules.ListLinks(ctx, rule.ID)
			if err != nil {
				return "", modbusshare.NewError(modbusshare.ErrCodeWorkspaceScope, "preview mapping ownership could not be verified", true)
			}
			for _, link := range links {
				if link == nil || link.PointID != record.PointID {
					continue
				}
				if link.MappingID != nil && *link.MappingID == record.ID {
					return pointRecord.DeviceID, nil
				}
			}
		}
		return "", modbusshare.NewError(modbusshare.ErrCodeWorkspaceScope, "preview mapping is not durably owned by this workspace", false)
	}
}

func containsString(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

// NewShareDesiredMappingOwnershipChecker proves the complete persisted
// candidate identity, not merely that a tag is linked to some rule.
func NewShareDesiredMappingOwnershipChecker(workspaceReader interface {
	ShareOwnershipSnapshot(context.Context) (string, []string, error)
}, rules *Service, tags shareTagReader, settings func() modbusshare.Settings, mappingSvc *mapping.Service) func(context.Context, modbusshare.DesiredMapping) error {
	return func(ctx context.Context, desired modbusshare.DesiredMapping) error {
		if workspaceReader == nil || rules == nil || tags == nil || settings == nil || mappingSvc == nil {
			return modbusshare.NewError(modbusshare.ErrCodeWorkspaceScope, "workspace ownership could not be verified", true)
		}
		if strings.TrimSpace(desired.MappingID) == "" {
			return modbusshare.NewError(modbusshare.ErrCodeWorkspaceScope, "mapping ownership is not persisted", false)
		}
		mappingRecord, mappingErr := mappingSvc.GetByID(ctx, desired.MappingID)
		if mappingErr != nil || mappingRecord == nil || mappingRecord.TagID != desired.TagID {
			return modbusshare.NewError(modbusshare.ErrCodeWorkspaceScope, "mapping relationship is not persisted", false)
		}
		workspaceID, deviceIDs, err := workspaceReader.ShareOwnershipSnapshot(ctx)
		if err != nil || workspaceID != desired.WorkspaceID {
			return modbusshare.NewError(modbusshare.ErrCodeWorkspaceScope, "workspace ownership could not be verified", true)
		}
		if _, err := tags.GetByID(ctx, desired.TagID); err != nil {
			return modbusshare.NewError(modbusshare.ErrCodeWorkspaceScope, "mapping tag is not persisted", false)
		}
		candidates, err := rules.BuildDesiredShareMappingsForDevices(ctx, desired.WorkspaceID, settings(), deviceIDs)
		if err != nil {
			return modbusshare.NewError(modbusshare.ErrCodeWorkspaceScope, "workspace mapping ownership could not be verified", true)
		}
		for _, candidate := range candidates {
			candidate.OwnershipProof = nil
			persisted := desired
			persisted.OwnershipProof = nil
			if reflect.DeepEqual(candidate, persisted) {
				return nil
			}
		}
		return modbusshare.NewError(modbusshare.ErrCodeWorkspaceScope, "mapping identity is not durably owned by this workspace", false)
	}
}

// NewShareOwnershipChecker returns a durable workspace/rule/tag ownership
// predicate for Modbus Share operations.
func NewShareOwnershipChecker(workspaceReader interface {
	ShareOwnershipSnapshot(context.Context) (string, []string, error)
}, rules *Service, tags shareTagReader) func(context.Context, string, string) bool {
	return func(ctx context.Context, workspaceID, tagID string) bool {
		if workspaceReader == nil || rules == nil || tags == nil {
			return false
		}
		resolvedWorkspaceID, deviceIDs, err := workspaceReader.ShareOwnershipSnapshot(ctx)
		if err != nil || resolvedWorkspaceID != workspaceID {
			return false
		}
		ruleRecords, err := rules.ListByDeviceIDs(ctx, deviceIDs)
		if err != nil {
			return false
		}
		for _, rule := range ruleRecords {
			if rule == nil {
				continue
			}
			links, err := rules.ListLinks(ctx, rule.ID)
			if err != nil {
				return false
			}
			for _, link := range links {
				if link != nil && link.TagID != nil && *link.TagID == tagID {
					_, tagErr := tags.GetByID(ctx, tagID)
					return tagErr == nil
				}
			}
		}
		return false
	}
}
