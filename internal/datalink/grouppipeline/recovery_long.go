package grouppipeline

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"
)

func longGroup(group *workspace.WriteGroup) bool {
	return group != nil && group.EffectiveLayoutVersion() == workspace.WriteGroupLayoutManagedLongV1
}

// validLongRuntimeProof is the depth check for a group on the shared typed long
// table. A recovered snapshot must prove, from its own frozen facts, that it
// belongs to the storage scope it names: plain storage, a profile that derives
// the row interval, the table and owner marker derived from the scope, the exact
// fixed column layout, an exact type for every source and a receipt table when
// the group dedupes by receipt. It never trusts a caller-supplied applied flag.
func validLongRuntimeProof(group *workspace.WriteGroup, layout *workspace.WriteGroupRuntimeLayout) bool {
	if group == nil || layout == nil || recordingReductionUnsupported(group, layout.TagTypes) {
		return false
	}
	table, owner, ok := workspace.LongStorageIdentity(group)
	if !ok || len(group.Members) == 0 || len(layout.Columns) == 0 ||
		layout.Dialect != dbtarget.SQLDialectSQLite && layout.Dialect != dbtarget.SQLDialectPostgres &&
			layout.Dialect != dbtarget.SQLDialectMySQL && layout.Dialect != dbtarget.SQLDialectSQLServer {
		return false
	}
	profile := group.RecordingProfile
	if profile == nil ||
		profile.Kind != workspace.WriteGroupProfileInstantaneous && profile.Kind != workspace.WriteGroupProfileCumulative ||
		profile.IntervalSeconds() != group.RowPolicy.IntervalSeconds || group.Destination.TableName != table {
		return false
	}
	if !validStorageModeProof(group, layout) {
		return false
	}
	if digest := strings.TrimSpace(layout.SchemaDigest); len(digest) != sha256.Size*2 {
		return false
	} else if _, err := hex.DecodeString(digest); err != nil {
		return false
	}
	ownerColumn := false
	for _, column := range layout.Columns {
		if strings.EqualFold(strings.TrimSpace(column.Name), owner) {
			ownerColumn = !column.Nullable && !column.PrimaryKey
		}
	}
	if !ownerColumn {
		return false
	}
	reduction := profile.EffectiveSampleReduction()
	members := make([]dbtarget.LongRowMember, 0, len(group.Members))
	for _, member := range group.Members {
		tagType, exact := layout.TagTypes[member.TagID]
		if !exact {
			return false
		}
		stored, ok := workspace.LongMemberStoredTypeForMember(reduction, member, tagType)
		if !ok {
			return false
		}
		members = append(members, dbtarget.LongRowMember{
			MemberKey: member.DeviceID + "/" + member.PointID + "/" + member.TagID,
			DeviceID:  member.DeviceID, PointID: member.PointID, TagID: member.TagID, Type: stored,
			Reduction: string(reduction), ValueRole: string(member.ValueRole), SourceValueType: member.SourceValueType,
		})
	}
	if _, issues := dbtarget.NewLongRowLayout(dbtarget.LongRowSpec{
		Dialect: layout.Dialect, Columns: layout.Columns, WorkspaceID: group.WorkspaceID, GroupID: group.ID,
		GroupRevision: group.Revision, ProfileKind: string(profile.Kind), ProfileInterval: profile.Interval,
		ProfileUnit: string(profile.Unit), Members: members,
	}); len(issues) > 0 {
		return false
	}
	if layout.Dialect == dbtarget.SQLDialectMySQL {
		return workspace.ValidateMySQLRuntimeLayout(group, layout) == nil
	}
	if layout.Dialect == dbtarget.SQLDialectSQLServer {
		return workspace.ValidateSQLServerRuntimeLayout(group, layout) == nil
	}
	return !strings.EqualFold(strings.TrimSpace(group.WritePolicy.DedupeCapability), groupdelivery.DedupeReceipt) || layout.ReceiptTableReady
}

// validStorageModeProof checks the storage mode against its frozen evidence.
// Plain storage never claims a hypertable. A timescale group must be on
// PostgreSQL and carry the catalog evidence Apply verified: the owned
// hypertable cut on bucket_start at the confirmed seven-day interval. The table
// and owner marker already differ by mode, because the mode is part of the
// storage scope that derives them.
func validStorageModeProof(group *workspace.WriteGroup, layout *workspace.WriteGroupRuntimeLayout) bool {
	switch group.StorageMode {
	case workspace.WriteGroupStorageModePlain:
		return layout.Hypertable == nil && !mysqlPartitionEvidencePresent(layout)
	case workspace.WriteGroupStorageModeNativeMonthly:
		if layout.Dialect == dbtarget.SQLDialectSQLServer {
			return workspace.ValidateSQLServerRuntimeLayout(group, layout) == nil
		}
		return mysqlNativeMonthlyProof(layout)
	case workspace.WriteGroupStorageModeTimescale:
		hypertable := layout.Hypertable
		return layout.Dialect == dbtarget.SQLDialectPostgres && hypertable != nil &&
			hypertable.TimeColumn == dbtarget.TimescalePartitionColumn &&
			hypertable.ChunkIntervalMicros == dbtarget.TimescaleChunkInterval.Microseconds()
	}
	return false
}

// mysqlPartitionEvidencePresent reports whether any frozen MySQL managed table
// in the layout carries native partition evidence. Plain storage, on every
// dialect, must never carry it.
func mysqlPartitionEvidencePresent(layout *workspace.WriteGroupRuntimeLayout) bool {
	if layout.MySQL == nil {
		return false
	}
	for _, physical := range []*dbtarget.MySQLManagedInspection{layout.MySQL.Group, layout.MySQL.Receipt} {
		if physical != nil && physical.Evidence != nil && physical.Evidence.Partitions != nil {
			return true
		}
	}
	return false
}

// mysqlNativeMonthlyProof checks the frozen native_monthly evidence: MySQL
// only, no hypertable claim, non-nil partition evidence on the group table and
// a receipt table that stays unpartitioned. The canonical workspace validator
// remains the authority; this is the pipeline's own depth check.
func mysqlNativeMonthlyProof(layout *workspace.WriteGroupRuntimeLayout) bool {
	if layout == nil || layout.Dialect != dbtarget.SQLDialectMySQL || layout.Hypertable != nil || layout.MySQL == nil {
		return false
	}
	var groupPartitions *dbtarget.MySQLPartitionEvidence
	if layout.MySQL.Group != nil && layout.MySQL.Group.Evidence != nil {
		groupPartitions = layout.MySQL.Group.Evidence.Partitions
	}
	if groupPartitions == nil {
		return false
	}
	return layout.MySQL.Receipt != nil && layout.MySQL.Receipt.Evidence != nil &&
		layout.MySQL.Receipt.Evidence.Partitions == nil
}
