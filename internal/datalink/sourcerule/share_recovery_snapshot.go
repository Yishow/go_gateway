package sourcerule

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"reflect"
	"slices"
	"strings"
	"time"

	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/schema"
)

// ErrShareRecoverySnapshot deliberately excludes raw stored data and SQL errors.
var ErrShareRecoverySnapshot = errors.New("source share recovery evidence is unavailable or changed")

// ShareRecoverySnapshot contains desired mappings and an opaque input witness.
// It does not authorize activation or replace the existing Share authority.
type ShareRecoverySnapshot struct {
	Desired       []modbusshare.DesiredMapping
	WitnessDigest string
}

type shareRecoveryReader struct{ tx *sql.Tx }

func (r shareRecoveryReader) List(ctx context.Context, f ListFilter) ([]*schema.SourceRule, error) {
	return readShareList(ctx, r.tx, f)
}
func (r shareRecoveryReader) ListLinks(ctx context.Context, id string) ([]*schema.SourceRuleLink, error) {
	return readShareListLinks(ctx, r.tx, id)
}
func (r shareRecoveryReader) ListCandidateSnapshots(ctx context.Context, id, rev string) ([]*schema.SourceRuleCandidateSnapshot, error) {
	return readShareListCandidateSnapshots(ctx, r.tx, id, rev)
}
func (r shareRecoveryReader) GetByID(ctx context.Context, id string) (*schema.Tag, error) {
	var found string
	if r.tx.QueryRowContext(ctx, `SELECT id FROM tags WHERE id=?`, id).Scan(&found) != nil {
		return nil, ErrShareRecoverySnapshot
	}
	return &schema.Tag{ID: found}, nil
}

// ReadShareRecoverySnapshot uses the caller's transaction and the original
// builder, never services that could open a second connection or refresh evidence.
func ReadShareRecoverySnapshot(ctx context.Context, tx *sql.Tx, workspaceID string, settings modbusshare.Settings) (*ShareRecoverySnapshot, error) {
	if tx == nil || strings.TrimSpace(workspaceID) == "" {
		return nil, ErrShareRecoverySnapshot
	}
	var workspaceRaw, settingsRaw string
	if tx.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key='studio_v2_workspace'`).Scan(&workspaceRaw) != nil || tx.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key='modbus_share'`).Scan(&settingsRaw) != nil {
		return nil, ErrShareRecoverySnapshot
	}
	var workspace struct {
		ID        string   `json:"id"`
		DeviceIDs []string `json:"ordered_device_ids"`
	}
	if json.Unmarshal([]byte(workspaceRaw), &workspace) != nil || workspace.ID != workspaceID || workspace.DeviceIDs == nil {
		return nil, ErrShareRecoverySnapshot
	}
	persisted := modbusshare.DefaultSettings()
	if json.Unmarshal([]byte(settingsRaw), &persisted) != nil {
		return nil, ErrShareRecoverySnapshot
	}
	expected := settings
	expected.UpdatedAt = time.Time{}
	persisted.UpdatedAt = time.Time{}
	if !reflect.DeepEqual(expected, persisted) {
		return nil, ErrShareRecoverySnapshot
	}
	reader := shareRecoveryReader{tx: tx}
	// Raw inputs are retained only inside this call, then represented by a digest.
	witness := []any{"share-recovery-input:v1", workspaceRaw, settingsRaw}
	tagIDs, mappingIDs := map[string]bool{}, map[string]bool{}
	for _, deviceID := range workspace.DeviceIDs {
		if strings.TrimSpace(deviceID) == "" {
			return nil, ErrShareRecoverySnapshot
		}
		rules, err := reader.List(ctx, ListFilter{DeviceID: &deviceID})
		if err != nil {
			return nil, ErrShareRecoverySnapshot
		}
		witness = append(witness, deviceID, rules)
		for _, rule := range rules {
			snapshots, err := reader.ListCandidateSnapshots(ctx, rule.ID, rule.RevisionID)
			if err != nil {
				return nil, ErrShareRecoverySnapshot
			}
			links, err := reader.ListLinks(ctx, rule.ID)
			if err != nil {
				return nil, ErrShareRecoverySnapshot
			}
			witness = append(witness, snapshots, links)
			current := currentLocalModbusSnapshot(snapshots)
			if settings.Enabled && rule.ShareEnabled && (strings.TrimSpace(rule.RevisionID) == "" || current == nil) {
				return nil, ErrShareRecoverySnapshot
			}
			for _, link := range links {
				var pointDevice string
				if tx.QueryRowContext(ctx, `SELECT device_id FROM points WHERE id=?`, link.PointID).Scan(&pointDevice) != nil || pointDevice != rule.DeviceID {
					return nil, ErrShareRecoverySnapshot
				}
				pointRows, err := shareRecoveryRows(ctx, tx, `SELECT id,device_id,address,function,data_type,data_format,format_mode,format_override,format_version,revision_id,enabled,mode,polling_group_id FROM points WHERE id=?`, link.PointID)
				if err != nil {
					return nil, ErrShareRecoverySnapshot
				}
				witness = append(witness, "point", link.PointID, pointRows)
				if link.TagID != nil {
					tagIDs[*link.TagID] = true
				}
				if link.MappingID != nil {
					mappingIDs[*link.MappingID] = true
				}
			}
			if current != nil {
				if settings.Enabled && rule.ShareEnabled {
					if current.Status != schema.SourceRuleCandidateStatusReady || strings.TrimSpace(current.Reason) != "" {
						return nil, ErrShareRecoverySnapshot
					}
					var envelope struct {
						Candidates json.RawMessage `json:"candidates"`
					}
					if json.Unmarshal([]byte(current.Payload), &envelope) != nil || len(envelope.Candidates) == 0 || string(envelope.Candidates) == "null" {
						return nil, ErrShareRecoverySnapshot
					}
				}
				candidates, err := decodeCandidatePayload[schema.SourceRuleLocalModbusOutputCandidate](current.Payload)
				if err != nil {
					return nil, ErrShareRecoverySnapshot
				}
				for _, candidate := range candidates {
					if settings.Enabled && rule.ShareEnabled {
						if candidate.Status != schema.SourceRuleLocalModbusOutputStatusReady || strings.TrimSpace(candidate.BlockingReason) != "" || candidate.TagID == nil || strings.TrimSpace(*candidate.TagID) == "" || candidate.Register == nil || candidate.PointID == "" {
							return nil, ErrShareRecoverySnapshot
						}
						if !slices.ContainsFunc(links, func(link *schema.SourceRuleLink) bool {
							return link.PointID == candidate.PointID && link.TagID != nil && *link.TagID == *candidate.TagID && (candidate.MappingID == nil || (link.MappingID != nil && *link.MappingID == *candidate.MappingID))
						}) {
							return nil, ErrShareRecoverySnapshot
						}
					}
					if candidate.TagID != nil {
						tagIDs[*candidate.TagID] = true
					}
					if candidate.MappingID != nil {
						mappingIDs[*candidate.MappingID] = true
					}
				}
			}
		}
	}
	for _, item := range []struct {
		table string
		ids   map[string]bool
	}{{"tags", tagIDs}, {"mappings", mappingIDs}} {
		ids := make([]string, 0, len(item.ids))
		for id := range item.ids {
			ids = append(ids, id)
		}
		slices.Sort(ids)
		for _, id := range ids {
			rows, err := shareRecoveryRows(ctx, tx, `SELECT * FROM `+item.table+` WHERE id=?`, id)
			if err != nil {
				return nil, ErrShareRecoverySnapshot
			}
			witness = append(witness, item.table, id, rows) // Empty row sets also form evidence.
		}
	}
	desired, err := buildDesiredShareMappingsFrom(ctx, workspaceID, settings, workspace.DeviceIDs, reader, reader)
	if err != nil {
		return nil, ErrShareRecoverySnapshot
	}
	for _, mapping := range desired {
		if mapping.MappingID == "" {
			return nil, ErrShareRecoverySnapshot
		}
		var tagID, pointID string
		if tx.QueryRowContext(ctx, `SELECT tag_id,point_id FROM mappings WHERE id=?`, mapping.MappingID).Scan(&tagID, &pointID) != nil || tagID != mapping.TagID {
			return nil, ErrShareRecoverySnapshot
		}
		var count int
		if tx.QueryRowContext(ctx, `SELECT count(*) FROM source_rule_links WHERE rule_id=? AND point_id=? AND tag_id=? AND mapping_id=?`, mapping.SourceRuleID, pointID, mapping.TagID, mapping.MappingID).Scan(&count) != nil || count != 1 {
			return nil, ErrShareRecoverySnapshot
		}
	}
	raw, err := json.Marshal(witness)
	if err != nil {
		return nil, ErrShareRecoverySnapshot
	}
	digest := sha256.Sum256(raw)
	return &ShareRecoverySnapshot{Desired: desired, WitnessDigest: hex.EncodeToString(digest[:])}, nil
}

// ValidateShareRecoverySnapshot checks the original witness without rebasing.
func ValidateShareRecoverySnapshot(ctx context.Context, tx *sql.Tx, workspaceID string, settings modbusshare.Settings, expectedDigest string) error {
	if len(expectedDigest) != sha256.Size*2 {
		return ErrShareRecoverySnapshot
	}
	snapshot, err := ReadShareRecoverySnapshot(ctx, tx, workspaceID, settings)
	if err != nil || snapshot.WitnessDigest != expectedDigest {
		return ErrShareRecoverySnapshot
	}
	return nil
}

func shareRecoveryRows(ctx context.Context, tx *sql.Tx, query string, args ...any) ([]any, error) {
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	result := []any{columns}
	for rows.Next() {
		values := make([]any, len(columns))
		dest := make([]any, len(columns))
		for i := range values {
			dest[i] = &values[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		for i, v := range values {
			if b, ok := v.([]byte); ok {
				values[i] = string(b)
			}
		}
		result = append(result, values)
	}
	return result, rows.Err()
}
