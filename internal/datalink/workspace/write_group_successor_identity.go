package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func (r *SQLWriteGroupRepository) validateRecreatedDraftIdentityInTx(ctx context.Context, tx *sql.Tx, predecessor, candidate *WriteGroup) error {
	source, err := r.recreatedProfileIdentityInTx(ctx, tx, predecessor, candidate)
	if err != nil {
		return err
	}
	if predecessor.RecordingProfile == nil {
		return nil
	}
	sourceEngine, err := r.recreateSourceEngineInTx(ctx, tx, source)
	if err != nil {
		return err
	}
	targetEngine, err := r.recreateConnectorEngineInTx(ctx, tx, candidate.Destination.ConnectorID)
	if err != nil {
		return err
	}
	tagTypes, err := loadWriteGroupTagTypes(ctx, tx, r.query, candidate)
	if err != nil {
		return err
	}
	return ValidateWriteGroupTransitionStructure(source, candidate, sourceEngine, targetEngine, tagTypes)
}

// recreatedProfileIdentityInTx checks aliases and immutable source meaning, not target family.
func (r *SQLWriteGroupRepository) recreatedProfileIdentityInTx(ctx context.Context, tx *sql.Tx, predecessor, candidate *WriteGroup) (*WriteGroup, error) {
	if predecessor.RecordingProfile == nil {
		return nil, r.validateRecreatedBasicInTx(ctx, tx, predecessor, candidate)
	}
	if candidate.RecordingProfile == nil {
		return nil, transitionRejection("unsupported-shape")
	}
	source, err := r.recreateStructuralSourceInTx(ctx, tx, predecessor)
	if err != nil {
		return nil, err
	}
	stored, found, err := r.loadProfileKey(ctx, tx, predecessor.WorkspaceID, predecessor.ID)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, lineageCorrupt("profile predecessor has no genesis profile alias")
	}
	original, err := circuitProfileKeyOf(predecessor)
	if err != nil || original != stored {
		return nil, lineageCorrupt("predecessor profile identity disagrees with genesis")
	}
	sourceKey, err := circuitProfileKeyOf(source)
	if err != nil || sourceKey != stored {
		return nil, lineageCorrupt("applied source identity disagrees with genesis")
	}
	derived, err := circuitProfileKeyOf(candidate)
	if err != nil {
		return nil, err
	}
	if derived != stored {
		return nil, transitionRejection("canonical-identity-change")
	}
	if err := r.validateRecreatedBasicAliasInTx(ctx, tx, predecessor, candidate); err != nil {
		return nil, err
	}
	return source, nil
}

func (r *SQLWriteGroupRepository) recreateConnectorEngineInTx(ctx context.Context, tx *sql.Tx, connectorID string) (string, error) {
	var kind string
	if err := tx.QueryRowContext(ctx, r.query(`SELECT kind FROM database_connectors WHERE id=$1`), connectorID).Scan(&kind); err != nil {
		return "", fmt.Errorf("read recreation engine evidence: %w", err)
	}
	kind = strings.ToLower(strings.TrimSpace(kind))
	// A MySQL driver is not proof of MySQL rather than MariaDB (or another
	// compatible server). This local primitive does not probe external servers.
	if kind == writeGroupConnectorKindMySQL || !transitionEngineSupported(kind) {
		return "", transitionRejection("engine-evidence-unavailable")
	}
	return kind, nil
}

// Existing dual aliases remain bound to the original device and canonical role;
// the compatibility field alone cannot detect a corrupted non-basic role.
func (r *SQLWriteGroupRepository) validateRecreatedBasicAliasInTx(ctx context.Context, tx *sql.Tx, predecessor, candidate *WriteGroup) error {
	var device, role string
	err := tx.QueryRowContext(ctx, r.query(`SELECT k.device_id,k.canonical_role
  FROM write_group_basic_keys k JOIN write_group_generations g ON g.lineage_id=k.lineage_id
  WHERE g.group_id=$1`), predecessor.ID).Scan(&device, &role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return lineageCorrupt("read recreation basic alias: %v", err)
	}
	if strings.TrimSpace(device) == "" || role != basicManagedCanonicalRole || device != predecessor.BasicManagedDeviceID {
		return lineageCorrupt("invalid recreation basic alias")
	}
	for _, member := range predecessor.Members {
		if member.DeviceID != device {
			return lineageCorrupt("predecessor source disagrees with basic genesis identity")
		}
	}
	for _, member := range candidate.Members {
		if member.DeviceID != device {
			return transitionRejection("canonical-identity-change")
		}
	}
	return nil
}

// Source engine identity belongs to its applied revision, not to an editable
// live connector. A present malformed descriptor never falls back to live state.
func (r *SQLWriteGroupRepository) recreateSourceEngineInTx(ctx context.Context, tx *sql.Tx, predecessor *WriteGroup) (string, error) {
	unavailable := func() (string, error) { return "", transitionRejection("source-engine-evidence-unavailable") }
	if predecessor.AppliedRevision != "" {
		// A globally unique group/revision cannot legitimately have runtime
		// descriptors in two workspaces; a foreign row is invalid, not absent.
		var descriptors int
		if err := tx.QueryRowContext(ctx, r.query(`SELECT COUNT(*) FROM wg_runtime_versions
            WHERE group_id=$1 AND group_revision=$2`), predecessor.ID, predecessor.AppliedRevision).Scan(&descriptors); err != nil || descriptors > 1 {
			return unavailable()
		}
		var payload string
		err := tx.QueryRowContext(ctx, r.query(`SELECT payload FROM wg_runtime_versions
   WHERE workspace_id=$1 AND group_id=$2 AND group_revision=$3`), predecessor.WorkspaceID, predecessor.ID, predecessor.AppliedRevision).Scan(&payload)
		if err != nil && (!errors.Is(err, sql.ErrNoRows) || descriptors != 0) {
			return unavailable()
		}
		if err == nil {
			var frozen WriteGroupAppliedSnapshot
			if json.Unmarshal([]byte(payload), &frozen) != nil || frozen.Group == nil || frozen.RuntimeLayout == nil ||
				frozen.WorkspaceID != predecessor.WorkspaceID || frozen.GroupID != predecessor.ID ||
				frozen.GroupRevision != predecessor.AppliedRevision || frozen.AppliedRevision != predecessor.AppliedRevision ||
				frozen.Group.WorkspaceID != predecessor.WorkspaceID || frozen.Group.ID != predecessor.ID ||
				frozen.Group.Revision != predecessor.AppliedRevision || frozen.Group.AppliedRevision != predecessor.AppliedRevision ||
				!sameRecreateDestinationIdentity(frozen.Group.Destination, predecessor.Destination) {
				return unavailable()
			}
			engine := string(frozen.RuntimeLayout.Dialect)
			if engine == writeGroupConnectorKindMySQL {
				if ValidateMySQLRuntimeLayout(frozen.Group, frozen.RuntimeLayout) != nil {
					return unavailable()
				}
				engine = frozen.RuntimeLayout.MySQL.Group.Evidence.ServerFamily
			}
			if engine == writeGroupConnectorKindSQLServer && (!sameWriteGroupAppliedSemantics(frozen.Group, predecessor) ||
				ValidateSQLServerRuntimeLayout(predecessor, frozen.RuntimeLayout) != nil) {
				return unavailable()
			}
			if !transitionEngineSupported(engine) {
				return unavailable()
			}
			return engine, nil
		}
	}
	// Legacy/draft fallback is evidence only while the live identity remains
	// exactly the predecessor's saved connector revision. No network is touched.
	var kind, revision string
	if err := tx.QueryRowContext(ctx, r.query(`SELECT kind,identity_revision FROM database_connectors WHERE id=$1`), predecessor.Destination.ConnectorID).Scan(&kind, &revision); err != nil {
		return unavailable()
	}
	kind = strings.ToLower(strings.TrimSpace(kind))
	if revision != predecessor.Destination.ConnectorRevision || kind == writeGroupConnectorKindMySQL || !transitionEngineSupported(kind) {
		return unavailable()
	}
	return kind, nil
}

func sameRecreateDestinationIdentity(before, after WriteGroupDestination) bool {
	return before.ConnectorID == after.ConnectorID && before.ConnectorRevision == after.ConnectorRevision &&
		before.Database == after.Database && before.TableSchema == after.TableSchema && before.TableName == after.TableName &&
		before.StorageStrategy == after.StorageStrategy
}

// Request CAS belongs to the current saved tombstone, but the last applied
// definition owns the recorded meaning. An unapplied edit is not a conversion.
func (r *SQLWriteGroupRepository) recreateStructuralSourceInTx(ctx context.Context, tx *sql.Tx, predecessor *WriteGroup) (*WriteGroup, error) {
	if predecessor.AppliedRevision == "" {
		return cloneWriteGroup(predecessor), nil
	}
	source, err := r.VersionByRevisionInTx(ctx, tx, predecessor.WorkspaceID, predecessor.ID, predecessor.AppliedRevision)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", transitionRejection("source-definition-unavailable"), err)
	}
	return cloneWriteGroup(source.Group), nil
}
