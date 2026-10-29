package workspace

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"strings"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/snapshot"
)

// WriteGroupLifecycleHistory is a local read snapshot, never transition authority.
// Legacy groups have Canonical=false and no invented lineage or generation.
type WriteGroupLifecycleHistory struct {
	Canonical       bool                            `json:"canonical"`
	LineageID       string                          `json:"lineage_id,omitempty"`
	LineageRevision int64                           `json:"lineage_revision,omitempty"`
	CurrentGroupID  string                          `json:"current_group_id"`
	ReservedGroupID string                          `json:"reserved_group_id,omitempty"`
	Generations     []WriteGroupLifecycleGeneration `json:"generations"`
}

// WriteGroupLifecycleDestination is an allowlisted identity. DestinationScope is
// an opaque hash; database paths, schema names and connector config stay private.
type WriteGroupLifecycleDestination struct {
	ConnectorID       string `json:"connector_id"`
	ConnectorRevision string `json:"connector_revision"`
	DestinationScope  string `json:"destination_scope"`
}

// WriteGroupLifecycleGeneration separates saved configuration from immutable
// applied revisions. Reserved does not imply schema confirmation or intake.
type WriteGroupLifecycleGeneration struct {
	GroupID            string                         `json:"group_id"`
	Generation         int64                          `json:"generation,omitempty"`
	PredecessorGroupID string                         `json:"predecessor_group_id,omitempty"`
	Current            bool                           `json:"current"`
	Reserved           bool                           `json:"reserved"`
	Status             WriteGroupStatus               `json:"status"`
	GroupRevision      string                         `json:"group_revision"`
	AppliedRevision    string                         `json:"applied_revision"`
	SavedDestination   WriteGroupLifecycleDestination `json:"saved_destination"`
	Revisions          []WriteGroupLifecycleRevision  `json:"revisions"`
}

// WriteGroupLifecycleRevision describes the original immutable owner and units.
// An absent EndsAt is unknown/open-ended, not a fabricated transition boundary.
type WriteGroupLifecycleRevision struct {
	GroupRevision string                         `json:"group_revision"`
	EffectiveAt   time.Time                      `json:"effective_at"`
	EndsAt        *time.Time                     `json:"ends_at,omitempty"`
	Sources       []CutoverMemberRef             `json:"sources"`
	Destination   WriteGroupLifecycleDestination `json:"destination"`
	Pending       groupdelivery.LifecyclePending `json:"pending"`
}

func (r *SQLWriteGroupRepository) lifecycleHistoryInTx(ctx context.Context, tx *sql.Tx, workspaceID, groupID string) (*WriteGroupLifecycleHistory, error) {
	if r == nil || tx == nil {
		return nil, ErrWriteGroupValidation
	}
	metadata, canonical, err := r.loadWriteGroupLineage(ctx, tx, workspaceID, groupID)
	if err != nil {
		return nil, err
	}
	result := &WriteGroupLifecycleHistory{Canonical: canonical, CurrentGroupID: groupID, Generations: []WriteGroupLifecycleGeneration{}}
	ids := []string{groupID}
	if canonical {
		result.LineageID, result.LineageRevision = metadata.LineageID, metadata.Revision
		result.CurrentGroupID, result.ReservedGroupID = metadata.CurrentGroupID, metadata.ReservedGroupID.String
		ids, err = r.lifecycleGenerationIDs(ctx, tx, metadata.LineageID)
		if err != nil {
			return nil, err
		}
	}
	for _, id := range ids {
		generation, err := r.lifecycleGenerationInTx(ctx, tx, workspaceID, id, result)
		if err != nil {
			return nil, err
		}
		result.Generations = append(result.Generations, generation)
	}
	return result, nil
}

func (r *SQLWriteGroupRepository) lifecycleGenerationInTx(ctx context.Context, tx *sql.Tx, workspaceID, id string, history *WriteGroupLifecycleHistory) (WriteGroupLifecycleGeneration, error) {
	result := WriteGroupLifecycleGeneration{GroupID: id, Current: id == history.CurrentGroupID, Reserved: id == history.ReservedGroupID, Revisions: []WriteGroupLifecycleRevision{}}
	group, err := r.GetInTx(ctx, tx, workspaceID, id)
	if err != nil {
		return result, err
	}
	if history.Canonical {
		metadata, found, err := r.loadWriteGroupLineage(ctx, tx, workspaceID, id)
		if err != nil {
			return result, err
		}
		if !found || metadata.LineageID != history.LineageID {
			return result, ErrWriteGroupLineageCorrupt
		}
		result.Generation, result.PredecessorGroupID = metadata.Generation, metadata.PredecessorGroupID.String
	}
	if result.Reserved && (group.Status != WriteGroupStatusDraft || group.AppliedRevision != "") {
		return result, ErrWriteGroupLineageCorrupt
	}
	result.Status, result.GroupRevision, result.AppliedRevision = group.Status, group.Revision, group.AppliedRevision
	result.SavedDestination, _, err = lifecycleDestination(group)
	if err != nil {
		return result, err
	}
	revisions, err := r.lifecycleRevisionsInTx(ctx, tx, group)
	if err != nil {
		return result, err
	}
	for _, revision := range revisions {
		item, err := r.lifecycleRevisionInTx(ctx, tx, group, revision)
		if err != nil {
			return result, err
		}
		result.Revisions = append(result.Revisions, item)
	}
	sort.Slice(result.Revisions, func(i, j int) bool {
		a, b := result.Revisions[i], result.Revisions[j]
		if a.EffectiveAt.Equal(b.EffectiveAt) {
			return a.GroupRevision < b.GroupRevision
		}
		return a.EffectiveAt.Before(b.EffectiveAt)
	})
	return result, nil
}

func (r *SQLWriteGroupRepository) lifecycleRevisionInTx(ctx context.Context, tx *sql.Tx, group *WriteGroup, revision string) (WriteGroupLifecycleRevision, error) {
	result := WriteGroupLifecycleRevision{GroupRevision: revision, Sources: []CutoverMemberRef{}}
	version, err := r.VersionByRevisionInTx(ctx, tx, group.WorkspaceID, group.ID, revision)
	if err != nil {
		return result, err
	}
	if version.EffectiveAt.IsZero() {
		return result, ErrWriteGroupValidation
	}
	result.EffectiveAt = version.EffectiveAt
	var frozen groupdelivery.FrozenDestination
	result.Destination, frozen, err = lifecycleDestination(version.Group)
	if err != nil {
		return result, err
	}
	for _, member := range version.Group.Members {
		result.Sources = append(result.Sources, memberRefOf(member))
	}
	sortMemberRefs(result.Sources)
	// revisionEndInTx is shared with cutover. Check workspace before using its
	// group/revision lookup, so foreign persisted metadata is not silently trusted.
	var endWorkspace string
	err = tx.QueryRowContext(ctx, r.query(`SELECT workspace_id FROM write_group_revision_ends WHERE group_id=$1 AND group_revision=$2`), group.ID, revision).Scan(&endWorkspace)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return result, err
	}
	if err == nil && endWorkspace != group.WorkspaceID {
		return result, ErrWriteGroupValidation
	}
	end, err := r.revisionEndInTx(ctx, tx, group.ID, revision)
	if err != nil {
		return result, err
	}
	if !end.IsZero() {
		if end.Before(version.EffectiveAt) {
			return result, ErrWriteGroupValidation
		}
		result.EndsAt = &end
	}
	frozen.RecordKeyColumn, err = r.lifecycleRecordKeyInTx(ctx, tx, version.Group)
	if err != nil {
		return result, err
	}
	result.Pending, err = groupdelivery.ReadLifecyclePending(ctx, tx, groupdelivery.GroupKey{WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: revision}, frozen)
	return result, err
}

func lifecycleDestination(group *WriteGroup) (WriteGroupLifecycleDestination, groupdelivery.FrozenDestination, error) {
	d := group.Destination
	scope, err := snapshot.DestinationScope(d.ConnectorID, d.ConnectorRevision, d.Database, d.TableSchema, d.TableName)
	if err != nil {
		return WriteGroupLifecycleDestination{}, groupdelivery.FrozenDestination{}, ErrWriteGroupValidation
	}
	dedupe := strings.ToLower(strings.TrimSpace(group.WritePolicy.DedupeCapability))
	if dedupe == "" {
		dedupe = groupdelivery.DedupeNone
	}
	return WriteGroupLifecycleDestination{ConnectorID: d.ConnectorID, ConnectorRevision: d.ConnectorRevision, DestinationScope: scope}, groupdelivery.FrozenDestination{
		Scope: scope, ConnectorID: d.ConnectorID, ConnectorRevision: d.ConnectorRevision, Database: d.Database, TableSchema: d.TableSchema, TableName: d.TableName, DedupeCapability: dedupe,
	}, nil
}
