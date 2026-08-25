package sourcerule

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/schema"
)

// BuildDesiredShareMappingsWithSettings returns desired Modbus Share mappings for the given workspace,
// enforcing global settings (B2: disabled settings return empty slice).
func (s *Service) BuildDesiredShareMappingsWithSettings(
	ctx context.Context,
	workspaceID string,
	settings modbusshare.Settings,
) ([]modbusshare.DesiredMapping, error) {
	return s.buildDesiredShareMappings(ctx, workspaceID, settings, nil)
}

// BuildDesiredShareMappingsForDevices scopes candidate hydration to the
// durable device membership of a workspace. Foreign rules are never allowed
// to make this workspace's hydration fail.
func (s *Service) BuildDesiredShareMappingsForDevices(ctx context.Context, workspaceID string, settings modbusshare.Settings, deviceIDs []string) ([]modbusshare.DesiredMapping, error) {
	return s.buildDesiredShareMappings(ctx, workspaceID, settings, deviceIDs)
}

func (s *Service) buildDesiredShareMappings(
	ctx context.Context,
	workspaceID string,
	settings modbusshare.Settings,
	deviceIDs []string,
) ([]modbusshare.DesiredMapping, error) {
	if !settings.Enabled {
		return []modbusshare.DesiredMapping{}, nil
	}

	var rules []*schema.SourceRule
	var err error
	if deviceIDs == nil {
		rules, err = s.repo.List(ctx, ListFilter{})
	} else {
		rules, err = s.ListByDeviceIDs(ctx, deviceIDs)
	}
	if err != nil {
		return nil, fmt.Errorf("列出來源規則失敗: %w", err)
	}

	sort.Slice(rules, func(i, j int) bool {
		if rules[i] == nil {
			return false
		}
		if rules[j] == nil {
			return true
		}
		return rules[i].ID < rules[j].ID
	})

	result := make([]modbusshare.DesiredMapping, 0)
	for _, rule := range rules {
		if rule == nil || !rule.ShareEnabled || strings.TrimSpace(rule.RevisionID) == "" {
			continue
		}

		snapshots, err := s.repo.ListCandidateSnapshots(ctx, rule.ID, rule.RevisionID)
		if err != nil {
			return nil, fmt.Errorf("列出 local modbus 候選快照失敗: %w", err)
		}

		snapshot := currentLocalModbusSnapshot(snapshots)
		if snapshot == nil {
			continue
		}

		candidates, err := decodeCandidatePayload[schema.SourceRuleLocalModbusOutputCandidate](snapshot.Payload)
		if err != nil {
			return nil, fmt.Errorf("解析 local modbus 候選快照失敗: %w", err)
		}

		for _, c := range candidates {
			span := localModbusRegisterCount(c.DataType)
			stride := 1
			if rule.ShareStride != nil {
				stride = *rule.ShareStride
			}
			if stride < span {
				return nil, &modbusshare.Error{
					Code:      modbusshare.ErrCodeInvalidGeometry,
					Message:   fmt.Sprintf("source-rule share stride %d is smaller than datatype span %d", stride, span),
					Retryable: false,
					Action:    localModbusStrideAction,
				}
			}
			if c.TagID == nil || c.Register == nil || strings.TrimSpace(*c.TagID) == "" {
				continue
			}
			if c.Status == schema.SourceRuleLocalModbusOutputStatusOutOfSync || strings.TrimSpace(c.BlockingReason) != "" {
				continue
			}
			if c.MappingID == nil {
				links, linkErr := s.ListLinks(ctx, rule.ID)
				if linkErr != nil {
					return nil, linkErr
				}
				for _, link := range links {
					if link != nil && link.TagID != nil && strings.TrimSpace(*link.TagID) == strings.TrimSpace(*c.TagID) && link.MappingID != nil {
						c.MappingID = stringPtr(strings.TrimSpace(*link.MappingID))
						break
					}
				}
			}
			// Candidate payloads are not ownership proof by themselves. Require
			// both the persisted tag and the rule link before projecting it.
			if s.tagSvc != nil {
				if _, err := s.tagSvc.GetByID(ctx, strings.TrimSpace(*c.TagID)); err != nil {
					return nil, &modbusshare.Error{Code: modbusshare.ErrCodeWorkspaceScope, Message: "local Modbus candidate tag is not persisted", Retryable: false, Action: "Repair the source-rule candidate and apply it again"}
				}
				links, err := s.ListLinks(ctx, rule.ID)
				if err != nil {
					return nil, err
				}
				linked := false
				for _, link := range links {
					if link != nil && link.TagID != nil && strings.TrimSpace(*link.TagID) == strings.TrimSpace(*c.TagID) {
						linked = true
						break
					}
				}
				if !linked {
					return nil, &modbusshare.Error{Code: modbusshare.ErrCodeWorkspaceScope, Message: "local Modbus candidate has no persisted source-rule relationship", Retryable: false, Action: "Apply the source-rule candidate before activation"}
				}
			}

			// Candidate.Register is the canonical zero-based runtime coordinate.
			// configuredLocalModbusRegister has already applied the durable human
			// start register, so adding that base here would shift the mapping twice.
			zeroBasedInt := int(*c.Register)
			if zeroBasedInt > int(modbusshare.MaxRegisterIndex) {
				return nil, &modbusshare.Error{Code: modbusshare.ErrCodeCapacityExceeded, Message: "source-rule share register exceeds capacity", Retryable: false, Action: "Choose a register within the configured capacity"}
			}
			zeroBased := uint16(zeroBasedInt)
			desired := modbusshare.DesiredMapping{
				WorkspaceID:        workspaceID,
				SourceRuleID:       rule.ID,
				SourceRuleRevision: rule.RevisionID,
				TagID:              strings.TrimSpace(*c.TagID),
				MappingID:          derefOptional(c.MappingID),
				DataType:           c.DataType,
				ShareStartRegister: modbusshare.ZeroBasedToHuman(zeroBased),
				ZeroBasedRegister:  zeroBased,
				SpanRegisters:      span,
				StrideRegisters:    stride,
				CapacityRegisters:  settings.CapacityRegisters,
				TagKey:             c.TagKey,
				DisplayName:        c.DisplayName,
			}
			result = append(result, desired)
		}
	}

	return result, nil
}
