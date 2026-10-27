package grouppipeline

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

var errRuntimeSnapshotInvalid = errors.New("invalid runtime snapshot")

func groupKey(snap *workspace.WriteGroupAppliedSnapshot) groupdelivery.GroupKey {
	return groupdelivery.GroupKey{WorkspaceID: snap.WorkspaceID, GroupID: snap.GroupID, GroupRevision: snap.AppliedRevision}
}

func sameAppliedSnapshot(a, b *workspace.WriteGroupAppliedSnapshot) bool {
	return a.WorkspaceID == b.WorkspaceID && a.GroupID == b.GroupID &&
		a.AppliedRevision == b.AppliedRevision && a.EffectiveAt.Equal(b.EffectiveAt) && reflect.DeepEqual(a.Group, b.Group)
}

func (p *Pipeline) saveRuntimeSnapshot(ctx context.Context, snap *workspace.WriteGroupAppliedSnapshot) error {
	payload, err := json.Marshal(snap)
	if err != nil {
		return err
	}
	return p.deps.Store.SaveRuntimeVersion(ctx, groupKey(snap), payload)
}

func (p *Pipeline) loadRuntimeSnapshot(ctx context.Context, key groupdelivery.GroupKey) (*workspace.WriteGroupAppliedSnapshot, error) {
	return loadFrozenSnapshot(ctx, p.deps.Store, key)
}

// loadFrozenSnapshot reads and verifies the descriptor frozen when a group
// revision was applied. Nil without error means a legacy revision with none.
func loadFrozenSnapshot(ctx context.Context, store *groupdelivery.Store, key groupdelivery.GroupKey) (*workspace.WriteGroupAppliedSnapshot, error) {
	payload, err := store.RuntimeVersion(ctx, key)
	if err != nil || payload == nil {
		return nil, err
	}
	if _, err := counterproof.ResolveSnapshot(payload, counterproof.SnapshotScope{WorkspaceID: key.WorkspaceID, GroupID: key.GroupID, AppliedRevision: key.GroupRevision}); err != nil {
		return nil, err
	}
	return decodeFrozenStorageSnapshot(payload, key)
}

func decodeFrozenStorageSnapshot(payload []byte, key groupdelivery.GroupKey) (*workspace.WriteGroupAppliedSnapshot, error) {
	var snap workspace.WriteGroupAppliedSnapshot
	if err := json.Unmarshal(payload, &snap); err != nil {
		return nil, errRuntimeSnapshotInvalid
	}
	if groupKey(&snap) != key || snap.Group == nil || snap.RuntimeLayout == nil ||
		snap.Group.ID != key.GroupID || snap.Group.WorkspaceID != key.WorkspaceID ||
		snap.Group.Revision != key.GroupRevision || snap.Group.AppliedRevision != key.GroupRevision || snap.EffectiveAt.IsZero() {
		return nil, errRuntimeSnapshotInvalid
	}
	// A long group also looks like a canonical managed group, but its columns are
	// the fixed shared layout, not one column per source: it is judged by its own
	// proof and never by the wide one.
	switch {
	case longGroup(snap.Group):
		if !validLongRuntimeProof(snap.Group, snap.RuntimeLayout) {
			return nil, errRuntimeSnapshotInvalid
		}
	case canonicalManagedGroup(snap.Group) && !validCanonicalManagedRuntimeProof(snap.Group, snap.RuntimeLayout):
		return nil, errRuntimeSnapshotInvalid
	}
	return &snap, nil
}

func canonicalManagedGroup(group *workspace.WriteGroup) bool {
	return group != nil && group.Destination.StorageStrategy == workspace.WriteGroupStorageStrategyManaged &&
		strings.TrimSpace(group.RowPolicy.RecordKeyColumn) != ""
}

// recordingReductionUnsupported reports a frozen group whose recording shape
// no runtime writer can serve yet: an explicit mean, the new wide profile
// layout, or member recording metadata. Such a descriptor must never restore
// a writer: recovery keeps it blocked instead of degrading it to the legacy
// latest behavior.
func recordingReductionUnsupported(group *workspace.WriteGroup, tagTypes map[string]schema.DataType) bool {
	if group == nil {
		return false
	}
	switch group.EffectiveLayoutVersion() {
	case workspace.WriteGroupLayoutManagedWideProfileV1:
		// Served only when every frozen fact is complete and verified; a cumulative
		// wide profile and any partial or contradictory descriptor stay blocked.
		return !workspace.WideProfileShapeSupported(group, tagTypes)
	case workspace.WriteGroupLayoutManagedLongV1:
		return !workspace.LongRecordingShapeSupported(group, tagTypes)
	}
	if group.RecordingProfile != nil && group.RecordingProfile.EffectiveSampleReduction() == workspace.WriteGroupSampleReductionMean {
		return true
	}
	for i := range group.Members {
		if group.Members[i].ValueRole != "" || group.Members[i].SourceValueType != "" || group.Members[i].ResultValueType != "" {
			return true
		}
	}
	return false
}

func validCanonicalManagedRuntimeProof(group *workspace.WriteGroup, layout *workspace.WriteGroupRuntimeLayout) bool {
	if recordingReductionUnsupported(group, layout.TagTypes) {
		return false
	}
	// Only the legacy wide layout and the verified wide profile layout have a
	// proven wide writer; any unknown layout must never start as one.
	version := group.EffectiveLayoutVersion()
	if version != workspace.WriteGroupLayoutManagedWideV1 && version != workspace.WriteGroupLayoutManagedWideProfileV1 {
		return false
	}
	wideProfile := version == workspace.WriteGroupLayoutManagedWideProfileV1
	if !canonicalManagedGroup(group) || layout == nil || len(group.Members) == 0 || len(layout.Columns) == 0 ||
		layout.Dialect != dbtarget.SQLDialectSQLite && layout.Dialect != dbtarget.SQLDialectPostgres {
		return false
	}
	digest := strings.TrimSpace(layout.SchemaDigest)
	if len(digest) != sha256.Size*2 {
		return false
	}
	if _, err := hex.DecodeString(digest); err != nil {
		return false
	}
	columns := make(map[string]dbtarget.ColumnInfo, len(layout.Columns))
	for _, column := range layout.Columns {
		name := strings.TrimSpace(column.Name)
		if name == "" || strings.TrimSpace(column.DataType) == "" {
			return false
		}
		key := strings.ToLower(name)
		if _, exists := columns[key]; exists {
			return false
		}
		column.Name = name
		columns[key] = column
	}
	requiredColumns := []string{
		group.RowPolicy.RecordKeyColumn,
		group.RowPolicy.GroupIDColumn,
		group.RowPolicy.BucketStartColumn,
		group.RowPolicy.ProvenanceColumn,
	}
	if strings.TrimSpace(group.RowPolicy.DeviceIDColumn) != "" {
		requiredColumns = append(requiredColumns, group.RowPolicy.DeviceIDColumn)
	}
	if strings.TrimSpace(group.RowPolicy.EntityKeyColumn) != "" {
		requiredColumns = append(requiredColumns, group.RowPolicy.EntityKeyColumn)
	}
	for _, name := range requiredColumns {
		name = strings.TrimSpace(name)
		if name == "" {
			return false
		}
		column, exists := columns[strings.ToLower(name)]
		if !exists || column.Nullable {
			return false
		}
	}
	record, exists := columns[strings.ToLower(strings.TrimSpace(group.RowPolicy.RecordKeyColumn))]
	if !exists || !record.PrimaryKey {
		return false
	}
	owner, exists := columns[strings.ToLower(canonicalManagedOwnerColumn(group))]
	if !exists || owner.Nullable || owner.PrimaryKey || !strings.EqualFold(strings.TrimSpace(owner.DataType), "TEXT") {
		return false
	}
	for _, member := range group.Members {
		if strings.TrimSpace(member.TargetColumn) == "" {
			return false
		}
		memberColumn, exists := columns[strings.ToLower(strings.TrimSpace(member.TargetColumn))]
		if !exists || (wideProfile && !memberColumn.Nullable) {
			// A wide profile window may carry NULL for a member, so its column must
			// be nullable; the frozen catalog columns prove it.
			return false
		}
		tagType, exists := layout.TagTypes[member.TagID]
		if !exists || !tagType.IsValid() {
			return false
		}
	}
	return !strings.EqualFold(strings.TrimSpace(group.WritePolicy.DedupeCapability), groupdelivery.DedupeReceipt) || layout.ReceiptTableReady
}

// The owner marker is part of the persisted proof. Keep this derivation in
// lockstep with workspace.managedSchemaLayout; recovery never creates a new one.
func canonicalManagedOwnerColumn(group *workspace.WriteGroup) string {
	digest := sha256.Sum256([]byte(group.WorkspaceID + "\x00" + group.ID))
	return "_gw_owner_" + hex.EncodeToString(digest[:12])
}

func (p *Pipeline) restoreOpenJournals(ctx context.Context) error {
	journals, err := p.deps.Store.OpenJournals(ctx)
	if err != nil {
		return err
	}
	var recoveryErrors []error
	for _, journal := range journals {
		key := boundaryKey(journal.Key.GroupID, journal.Key.GroupRevision)
		p.mu.RLock()
		_, exists := p.boundaries[key]
		p.mu.RUnlock()
		if exists {
			continue
		}
		snap, err := p.loadRuntimeSnapshot(ctx, journal.Key)
		if err != nil {
			p.setBlocked(&workspace.WriteGroup{WorkspaceID: journal.Key.WorkspaceID, ID: journal.Key.GroupID, AppliedRevision: journal.Key.GroupRevision}, frozenFailureReason(err))
			if !isCounterProofReason(frozenFailureReason(err)) && !errors.Is(err, errRuntimeSnapshotInvalid) {
				recoveryErrors = append(recoveryErrors, err)
			}
			continue
		}
		if snap == nil {
			// Legacy journals may be recovered only through their original immutable
			// revision and verified destination, never through the current draft.
			snap, err = p.deps.Groups.ResolveAppliedAt(ctx, journal.Key.GroupID, journal.FirstBucket)
			if err != nil || snap == nil || groupKey(snap) != journal.Key {
				p.setBlocked(&workspace.WriteGroup{WorkspaceID: journal.Key.WorkspaceID, ID: journal.Key.GroupID, AppliedRevision: journal.Key.GroupRevision}, reasonSnapshotUnavailable)
				continue
			}
			_, reason := p.build(ctx, snap)
			if reason != "" {
				p.setBlocked(snap.Group, reason)
				continue
			}
			frozen, loadErr := p.loadRuntimeSnapshot(ctx, journal.Key)
			if loadErr != nil || frozen == nil {
				p.setBlocked(snap.Group, reasonSnapshotUnavailable)
				continue
			}
			snap = frozen
		}
		if source, ok := p.deps.Groups.(interface {
			RevisionEnd(context.Context, string, string) (time.Time, error)
		}); ok {
			if ends, err := source.RevisionEnd(ctx, journal.Key.GroupID, journal.Key.GroupRevision); err == nil {
				snap.EndsAt = ends
			}
		}
		until := journal.LastBucket.Add(time.Duration(snap.Group.RowPolicy.IntervalSeconds) * time.Second)
		m, reason, buildErr := p.buildFrozenWithCause(ctx, snap, journal.FirstBucket, until)
		if reason != "" {
			p.setBlocked(snap.Group, reason)
			if buildErr != nil && !isCounterProofReason(reason) && !errors.Is(buildErr, groupdelivery.ErrBaselineCorrupt) {
				recoveryErrors = append(recoveryErrors, buildErr)
			}
			continue
		}
		m.drainOnly = true
		p.installManaged(m)
	}
	return errors.Join(recoveryErrors...)
}

func frozenFailureReason(err error) string {
	var proof *counterproof.CounterBindingProofError
	if errors.As(err, &proof) {
		return proof.Reason.Error()
	}
	return reasonSnapshotUnavailable
}
