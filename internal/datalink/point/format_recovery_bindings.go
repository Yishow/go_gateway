package point

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"time"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
)

// FormatRecoveryTargetInputs is internal detached input for the existing typed
// binding builder. Never serialize credentials or use this as a new resolver.
type FormatRecoveryTargetInputs struct {
	Device        schema.Device          `json:"-"`
	Mappings      []*schema.Mapping      `json:"-"`
	Tags          map[string]*schema.Tag `json:"-"`
	WitnessDigest string                 `json:"-"`
}

// ReadFormatRecoveryTargetInputs reads the complete target mapping set and all
// referenced tags in the caller's snapshot. Telemetry and timestamps are not
// configuration identity; insertion/deletion and disabled mappings are.
func ReadFormatRecoveryTargetInputs(ctx context.Context, tx *sql.Tx, p *schema.Point) (*FormatRecoveryTargetInputs, error) {
	if tx == nil || p == nil || p.ID == "" || p.DeviceID == "" {
		return nil, ErrFormatInvalid
	}
	result := &FormatRecoveryTargetInputs{Tags: make(map[string]*schema.Tag)}
	if tx.QueryRowContext(ctx, `SELECT id,name,COALESCE(description,''),protocol,status,connection_config FROM devices WHERE id=? AND EXISTS (SELECT 1 FROM points WHERE id=? AND device_id=? AND revision_id=?)`, p.DeviceID, p.ID, p.DeviceID, p.RevisionID).Scan(&result.Device.ID, &result.Device.Name, &result.Device.Description, &result.Device.Protocol, &result.Device.Status, &result.Device.ConnectionConfig) != nil {
		return nil, ErrFormatUnavailable
	}
	var err error
	result.Mappings, err = mapping.ReadByPointIDTx(ctx, tx, p.ID)
	if err != nil {
		return nil, ErrFormatUnavailable
	}
	// The witness has its own private encoding because public schema tags omit
	// KeyLower, and the exported evidence deliberately serializes to {}.
	type tagIdentity struct {
		Tag      schema.Tag
		KeyLower string
	}
	tags := make([]tagIdentity, 0, len(result.Mappings))
	mappings := make([]schema.Mapping, 0, len(result.Mappings))
	for _, m := range result.Mappings {
		if m == nil || m.PointID != p.ID || m.TagID == "" {
			return nil, ErrFormatUnavailable
		}
		identity := *m
		identity.CreatedAt = time.Time{}
		identity.UpdatedAt = time.Time{}
		mappings = append(mappings, identity)
		if _, exists := result.Tags[m.TagID]; exists {
			continue
		}
		target, e := tag.ReadByIDTx(ctx, tx, m.TagID)
		if e != nil {
			return nil, ErrFormatUnavailable
		}
		identityTag := *target
		identityTag.CreatedAt = time.Time{}
		identityTag.UpdatedAt = time.Time{}
		result.Tags[m.TagID] = target
		tags = append(tags, tagIdentity{Tag: identityTag, KeyLower: target.KeyLower})
	}
	raw, err := json.Marshal(struct {
		Version           string
		PointID, DeviceID string
		Device            schema.Device
		Mappings          []schema.Mapping
		Tags              []tagIdentity
	}{"format-recovery-target-input:v1", p.ID, p.DeviceID, result.Device, mappings, tags})
	if err != nil {
		return nil, ErrFormatUnavailable
	}
	digest := sha256.Sum256(raw)
	result.WitnessDigest = hex.EncodeToString(digest[:])
	return result, nil
}
