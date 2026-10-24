package grouppipeline

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// timescaleProof builds the group and frozen layout an Apply of a timescale
// long group leaves behind, with the catalog evidence it verified.
func timescaleProof(t *testing.T, mutate func(*workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout)) (*workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout) {
	t.Helper()
	return longProof(t, workspace.WriteGroupStorageModeTimescale, mutate)
}

// longProof builds the group and frozen layout an Apply of a long group in the
// given storage mode leaves behind: the mode's own table and owner marker, and
// for timescale the catalog evidence it verified.
func longProof(t *testing.T, mode workspace.WriteGroupStorageMode, mutate func(*workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout)) (*workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout) {
	t.Helper()
	group := &workspace.WriteGroup{
		ID: "group-ts", WorkspaceID: "workspace-A", Revision: "rev-1", AppliedRevision: "rev-1",
		LayoutVersion: workspace.WriteGroupLayoutManagedLongV1, StorageMode: mode,
		RecordingProfile: &workspace.WriteGroupRecordingProfile{Kind: workspace.WriteGroupProfileInstantaneous, Interval: 15, Unit: workspace.WriteGroupIntervalUnitSecond},
		Members:          []workspace.WriteGroupMember{{DeviceID: "device-1", PointID: "point-1", TagID: "tag-1", TargetColumn: "v"}},
		Destination: workspace.WriteGroupDestination{
			ConnectorID: "connector-1", ConnectorRevision: "crev-1", Database: "gwtest", TableSchema: "app", StorageStrategy: workspace.WriteGroupStorageStrategyManaged,
		},
		RowPolicy:   workspace.WriteGroupRowPolicy{IntervalSeconds: 15, RecordKeyColumn: "record_id"},
		WritePolicy: workspace.WriteGroupWritePolicy{Mode: "append", DedupeCapability: groupdelivery.DedupeReceipt},
	}
	table, owner, ok := workspace.LongStorageIdentity(group)
	require.True(t, ok)
	group.Destination.TableName = table
	layout, err := recordingplan.ManagedLongGroupLayout(string(dbtarget.SQLDialectPostgres), table, owner)
	if mode == workspace.WriteGroupStorageModeTimescale {
		layout, err = recordingplan.ManagedLongTimescaleLayout(table, owner)
	}
	require.NoError(t, err)
	columns := make([]dbtarget.ColumnInfo, 0, len(layout.Columns))
	for _, column := range layout.Columns {
		columns = append(columns, dbtarget.ColumnInfo{Name: column.Name, DataType: column.SQLType, Nullable: column.Nullable, PrimaryKey: column.PrimaryKey, Unique: column.PrimaryKey})
	}
	digest := sha256.Sum256([]byte(table))
	frozen := &workspace.WriteGroupRuntimeLayout{
		Dialect: dbtarget.SQLDialectPostgres, Columns: columns,
		TagTypes:          map[string]schema.DataType{"tag-1": schema.DataTypeFloat64},
		ReceiptTableReady: true, SchemaDigest: hex.EncodeToString(digest[:]),
	}
	if mode == workspace.WriteGroupStorageModeTimescale {
		frozen.Hypertable = &workspace.WriteGroupHypertableEvidence{TimeColumn: "bucket_start", ChunkIntervalMicros: 7 * 24 * 3600 * 1000000}
	}
	if mutate != nil {
		mutate(group, frozen)
	}
	return group, frozen
}

// TestTimescaleRuntimeProof pins what a recovered timescale group must prove
// before the pipeline writes to it: the same plain-table checks, and the frozen
// catalog evidence of the owned hypertable on the confirmed time dimension and
// interval. A plain proof never claims a hypertable.
func TestTimescaleRuntimeProof(t *testing.T) {
	group, layout := timescaleProof(t, nil)
	require.True(t, validLongRuntimeProof(group, layout), "the proof Apply leaves behind is accepted")

	for name, mutate := range map[string]func(*workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout){
		"no hypertable evidence": func(_ *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) { l.Hypertable = nil },
		"another chunk interval": func(_ *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
			l.Hypertable.ChunkIntervalMicros = 24 * 3600 * 1000000
		},
		"a microsecond more than the confirmed interval": func(_ *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
			l.Hypertable.ChunkIntervalMicros++
		},
		"a microsecond less than the confirmed interval": func(_ *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
			l.Hypertable.ChunkIntervalMicros--
		},
		"another time column": func(_ *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
			l.Hypertable.TimeColumn = "bucket_end"
		},
		"a SQLite dialect": func(_ *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
			l.Dialect = dbtarget.SQLDialectSQLite
		},
		"the plain table's name": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			plain := *g
			plain.StorageMode = workspace.WriteGroupStorageModePlain
			g.Destination.TableName, _, _ = workspace.LongStorageIdentity(&plain)
		},
	} {
		group, layout := timescaleProof(t, mutate)
		require.False(t, validLongRuntimeProof(group, layout), name)
	}

	// The plain mode's own proof is accepted, and never with a hypertable claimed.
	plain, plainLayout := longProof(t, workspace.WriteGroupStorageModePlain, nil)
	require.True(t, validLongRuntimeProof(plain, plainLayout))
	claimed, claimedLayout := longProof(t, workspace.WriteGroupStorageModePlain, func(_ *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
		l.Hypertable = &workspace.WriteGroupHypertableEvidence{TimeColumn: "bucket_start", ChunkIntervalMicros: 7 * 24 * 3600 * 1000000}
	})
	require.False(t, validLongRuntimeProof(claimed, claimedLayout), "a plain proof never claims a hypertable")

	// Timescale columns carry the timescale owner marker, so they never satisfy the plain mode.
	crossed, crossedLayout := timescaleProof(t, func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
		g.StorageMode = workspace.WriteGroupStorageModePlain
	})
	require.False(t, validLongRuntimeProof(crossed, crossedLayout), "the modes never own each other's table")
}
