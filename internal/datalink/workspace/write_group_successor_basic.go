package workspace

import (
	"context"
	"database/sql"
	"reflect"
	"strings"

	"go-gateway/internal/datalink/measurement"
)

// Basic recreation is a separate same-shape contract, not an addition to the
// managed-profile transition matrix. Both checks run inside the caller's tx;
// the second check also verifies generic normalization did not change meaning.
func (r *SQLWriteGroupRepository) validateRecreatedBasicInTx(ctx context.Context, tx *sql.Tx, predecessor, candidate *WriteGroup) error {
	source, err := r.recreateStructuralSourceInTx(ctx, tx, predecessor)
	if err != nil {
		return err
	}
	if predecessor.RecordingProfile != nil || source.RecordingProfile != nil || candidate.RecordingProfile != nil ||
		predecessor.EffectiveLayoutVersion() != WriteGroupLayoutManagedWideV1 || source.EffectiveLayoutVersion() != WriteGroupLayoutManagedWideV1 || candidate.EffectiveLayoutVersion() != WriteGroupLayoutManagedWideV1 ||
		source.Destination.StorageStrategy != WriteGroupStorageStrategyManaged || candidate.Destination.StorageStrategy != WriteGroupStorageStrategyManaged ||
		source.StorageMode != "" || candidate.StorageMode != "" {
		return transitionRejection("unsupported-shape")
	}
	if strings.TrimSpace(predecessor.BasicManagedDeviceID) == "" {
		return lineageCorrupt("Basic predecessor has no canonical device")
	}
	if err := r.validateRecreatedBasicAliasInTx(ctx, tx, predecessor, candidate); err != nil {
		return err
	}
	// A nil profile cannot silently inherit or create a profile alias.
	if _, found, err := r.loadProfileKey(ctx, tx, predecessor.WorkspaceID, predecessor.ID); err != nil {
		return err
	} else if found {
		return lineageCorrupt("Basic-only predecessor has a profile alias")
	}
	if source.Destination.ConnectorID != candidate.Destination.ConnectorID || source.Destination.ConnectorRevision != candidate.Destination.ConnectorRevision ||
		source.Destination.Database != candidate.Destination.Database || source.Destination.TableSchema != candidate.Destination.TableSchema {
		return transitionRejection("destination-change")
	}
	sourceEngine, err := r.recreateSourceEngineInTx(ctx, tx, source)
	if err != nil {
		return err
	}
	targetEngine, err := r.recreateConnectorEngineInTx(ctx, tx, candidate.Destination.ConnectorID)
	if err != nil {
		return err
	}
	if sourceEngine != targetEngine {
		return transitionRejection("cross-engine")
	}
	if sourceEngine != writeGroupConnectorKindSQLite && sourceEngine != writeGroupConnectorKindPostgres {
		return transitionRejection("unsupported-engine")
	}
	sourcePolicy, targetPolicy := source.RowPolicy, candidate.RowPolicy
	sourcePolicy.GroupKeyColumns = normalizeStringSetForComparison(sourcePolicy.GroupKeyColumns)
	sourcePolicy.UniqueKeyColumns = normalizeStringSetForComparison(sourcePolicy.UniqueKeyColumns)
	targetPolicy.GroupKeyColumns = normalizeStringSetForComparison(targetPolicy.GroupKeyColumns)
	targetPolicy.UniqueKeyColumns = normalizeStringSetForComparison(targetPolicy.UniqueKeyColumns)
	if !reflect.DeepEqual(sourcePolicy, targetPolicy) || !reflect.DeepEqual(source.WritePolicy, candidate.WritePolicy) {
		return transitionRejection("recording-meaning-change")
	}
	before, beforeOK := transitionMembers(source)
	after, afterOK := transitionMembers(candidate)
	if !beforeOK || !afterOK || len(before) != len(after) {
		return transitionRejection("partial-source-set")
	}
	types, err := loadWriteGroupTagTypes(ctx, tx, r.query, candidate)
	if err != nil {
		return err
	}
	for identity, member := range before {
		next, exists := after[identity]
		if !exists {
			return transitionRejection("partial-source-set")
		}
		if member.DeviceID != predecessor.BasicManagedDeviceID {
			return lineageCorrupt("Basic applied source disagrees with canonical device")
		}
		if member.SourceRevision == "" || member.MappingRevision == "" {
			return transitionRejection("source-definition-unavailable")
		}
		// Circuit names are display text; every other persisted member field owns
		// source, binding, column, requiredness, age or recording meaning.
		member.CircuitName, next.CircuitName = "", ""
		if !reflect.DeepEqual(member, next) {
			return transitionRejection("recording-meaning-change")
		}
		if _, known := measurement.ExactTypeForTag(types[member.TagID]); !known {
			return transitionRejection("unsupported-shape")
		}
	}
	return nil
}
