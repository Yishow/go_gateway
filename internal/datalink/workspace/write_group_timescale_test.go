package workspace

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func timescaleGroup(mode WriteGroupStorageMode) *WriteGroup {
	return &WriteGroup{
		ID: "group-ts", WorkspaceID: "workspace-1", LayoutVersion: WriteGroupLayoutManagedLongV1, StorageMode: mode,
		RecordingProfile: &WriteGroupRecordingProfile{Kind: WriteGroupProfileInstantaneous, Interval: 15, Unit: WriteGroupIntervalUnitSecond},
		Destination: WriteGroupDestination{
			ConnectorID: "connector-1", Database: "gwtest", TableSchema: "app", StorageStrategy: WriteGroupStorageStrategyManaged,
		},
		Members: []WriteGroupMember{{DeviceID: "device", PointID: "point", TagID: "tag", TargetColumn: "v"}},
	}
}

// TestTimescaleStorageModeIsExplicit pins how the mode is selected and what it
// changes: it is accepted only on the long layout, the plain default is
// unchanged, and the storage scope (and so the table and its owner marker) is a
// different identity from plain, so a plain table is never reused or altered.
func TestTimescaleStorageModeIsExplicit(t *testing.T) {
	t.Run("the mode is accepted on a long group, and plain stays the default", func(t *testing.T) {
		explicit := timescaleGroup(WriteGroupStorageModeTimescale)
		require.NoError(t, normalizeWriteGroupStorageMode(explicit))
		require.Equal(t, WriteGroupStorageModeTimescale, explicit.StorageMode)

		defaulted := timescaleGroup("")
		require.NoError(t, normalizeWriteGroupStorageMode(defaulted))
		require.Equal(t, WriteGroupStorageModePlain, defaulted.StorageMode, "timescale is never the default")

		for name, group := range map[string]*WriteGroup{
			"an unknown mode": timescaleGroup(WriteGroupStorageMode("native")),
			"a wide group with a mode": func() *WriteGroup {
				g := timescaleGroup(WriteGroupStorageModeTimescale)
				g.LayoutVersion = ""
				return g
			}(),
		} {
			require.ErrorIs(t, normalizeWriteGroupStorageMode(group), ErrWriteGroupValidation, name)
		}
	})

	t.Run("the mode is part of the storage scope identity", func(t *testing.T) {
		plainTable, plainOwner, ok := LongStorageIdentity(timescaleGroup(WriteGroupStorageModePlain))
		require.True(t, ok)
		tsTable, tsOwner, ok := LongStorageIdentity(timescaleGroup(WriteGroupStorageModeTimescale))
		require.True(t, ok)
		require.NotEqual(t, plainTable, tsTable, "a timescale table is never the plain table")
		require.NotEqual(t, plainOwner, tsOwner, "and has its own ownership marker")
		again, _, _ := LongStorageIdentity(timescaleGroup(WriteGroupStorageModeTimescale))
		require.Equal(t, tsTable, again, "the identity is stable")
		other := timescaleGroup(WriteGroupStorageModeTimescale)
		other.ID = "another-group"
		sharedTable, _, _ := LongStorageIdentity(other)
		require.Equal(t, tsTable, sharedTable, "groups of one scope still share one table")
	})

	t.Run("the layout is the owned hypertable layout, and only PostgreSQL has one", func(t *testing.T) {
		tagTypes := map[string]schema.DataType{"tag": schema.DataTypeFloat64}
		layout, err := managedSchemaLayout(timescaleGroup(WriteGroupStorageModeTimescale), tagTypes, "postgres")
		require.NoError(t, err)
		require.Equal(t, &recordingplan.GroupSchemaPartition{
			Kind: recordingplan.PartitionKindTimescale, Column: "bucket_start", ChunkInterval: "7 days",
		}, layout.Partition)
		plain, err := managedSchemaLayout(timescaleGroup(WriteGroupStorageModePlain), tagTypes, "postgres")
		require.NoError(t, err)
		require.Nil(t, plain.Partition, "plain storage is untouched")

		_, err = managedSchemaLayout(timescaleGroup(WriteGroupStorageModeTimescale), tagTypes, "sqlite")
		require.ErrorIs(t, err, recordingplan.ErrTimescaleUnavailable, "SQLite rejects the selection")
		var unavailable *recordingplan.TimescaleUnavailableError
		require.ErrorAs(t, err, &unavailable)
		require.Equal(t, []string{"dialect-unsupported"}, unavailable.Reasons)
	})
}

// TestTimescaleHypertableJoinsTheSchemaDigest: the readiness digest a confirmation
// is bound to covers the hypertable evidence, so a changed interval, time column,
// compression or job invalidates it exactly like a changed column does, and a
// plain table never carries any.
func TestTimescaleHypertableJoinsTheSchemaDigest(t *testing.T) {
	group := timescaleGroup(WriteGroupStorageModeTimescale)
	group.Destination.TableName = "gw_long_x"
	inspect := func(hypertable *dbtarget.HypertableInfo) string {
		digest, err := digestWriteGroupInspection(group, &dbtarget.TableInspection{
			Status: dbtarget.TableInspectionExists, Hypertable: hypertable, HypertableInspected: true,
			Columns: []dbtarget.ColumnInfo{{Name: "bucket_start", DataType: "timestamp with time zone"}},
		})
		require.NoError(t, err)
		return digest
	}
	confirmed := dbtarget.HypertableInfo{TimeColumn: "bucket_start", ChunkInterval: 7 * 24 * time.Hour, Dimensions: 1}
	base := inspect(&confirmed)
	require.Equal(t, base, inspect(&confirmed), "the same catalog evidence gives the same digest")
	for name, changed := range map[string]dbtarget.HypertableInfo{
		"another interval":    {TimeColumn: "bucket_start", ChunkInterval: 24 * time.Hour, Dimensions: 1},
		"another time column": {TimeColumn: "bucket_end", ChunkInterval: 7 * 24 * time.Hour, Dimensions: 1},
		"a second dimension":  {TimeColumn: "bucket_start", ChunkInterval: 7 * 24 * time.Hour, Dimensions: 2},
		"compression":         {TimeColumn: "bucket_start", ChunkInterval: 7 * 24 * time.Hour, Dimensions: 1, CompressionEnabled: true},
		"a bound job":         {TimeColumn: "bucket_start", ChunkInterval: 7 * 24 * time.Hour, Dimensions: 1, Jobs: 1},
	} {
		require.NotEqual(t, base, inspect(&changed), name)
	}
	require.NotEqual(t, base, inspect(nil), "a table that is no longer a hypertable changes the digest")
}

// TestFrozenHypertableEvidenceIsExactlyTheConfirmedLayout: the runtime layout
// records a hypertable only when the catalog showed precisely the confirmed
// partition column and interval, in the catalog's own microseconds, so nothing
// coarser can freeze a differing hypertable as verified.
func TestFrozenHypertableEvidenceIsExactlyTheConfirmedLayout(t *testing.T) {
	confirmed := dbtarget.HypertableInfo{TimeColumn: "bucket_start", ChunkInterval: 7 * 24 * time.Hour, Dimensions: 1}
	require.Equal(t, &WriteGroupHypertableEvidence{TimeColumn: "bucket_start", ChunkIntervalMicros: 604800000000}, freezeHypertableEvidence(&confirmed))

	for name, change := range map[string]func(*dbtarget.HypertableInfo){
		"a microsecond more":        func(h *dbtarget.HypertableInfo) { h.ChunkInterval += time.Microsecond },
		"a microsecond less":        func(h *dbtarget.HypertableInfo) { h.ChunkInterval -= time.Microsecond },
		"0.4 seconds more":          func(h *dbtarget.HypertableInfo) { h.ChunkInterval += 400 * time.Millisecond },
		"another interval":          func(h *dbtarget.HypertableInfo) { h.ChunkInterval = 24 * time.Hour },
		"an unknown interval":       func(h *dbtarget.HypertableInfo) { h.ChunkInterval, h.ChunkIntervalUnknown = 0, true },
		"an unknown but seven days": func(h *dbtarget.HypertableInfo) { h.ChunkIntervalUnknown = true },
		"another time column":       func(h *dbtarget.HypertableInfo) { h.TimeColumn = "bucket_end" },
	} {
		changed := confirmed
		change(&changed)
		require.Nil(t, freezeHypertableEvidence(&changed), name)
	}
	require.Nil(t, freezeHypertableEvidence(nil), "no hypertable, no evidence")
}
