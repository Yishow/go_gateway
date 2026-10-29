package workspace

import (
	"testing"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func requireTransitionReason(t *testing.T, err error, reason string) {
	t.Helper()
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	var rejection *WriteGroupTransitionError
	require.ErrorAs(t, err, &rejection)
	require.Equal(t, reason, rejection.ReasonCode)
	require.Equal(t, "recording transition structure rejected", err.Error())
}

func TestWriteGroupTransitionMatrixRejectsEngineLayoutAndMode(t *testing.T) {
	for _, tc := range []struct {
		name, fromEngine, toEngine, reason string
		from, to                           WriteGroupLayoutVersion
		fromMode, toMode                   WriteGroupStorageMode
	}{
		{"SQL Server wide example", "sqlserver", "sqlserver", "unsupported-layout", WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedWideProfileV1, "plain", ""},
		{"MySQL wide target", "mysql", "mysql", "unsupported-layout", WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedWideProfileV1, "plain", ""},
		{"MariaDB wide target", "mariadb", "mariadb", "unsupported-layout", WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedWideProfileV1, "plain", ""},
		{"MySQL wide source", "mysql", "mysql", "unsupported-layout", WriteGroupLayoutManagedWideProfileV1, WriteGroupLayoutManagedLongV1, "", "plain"},
		{"cross engine", "sqlite", "postgres", "cross-engine", WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedLongV1, "plain", "plain"},
		{"MySQL is not MariaDB", "mysql", "mariadb", "cross-engine", WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedLongV1, "plain", "plain"},
		{"MariaDB is not MySQL", "mariadb", "mysql", "cross-engine", WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedLongV1, "plain", "plain"},
		{"unknown engine", "unknown", "unknown", "unsupported-engine", WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedLongV1, "plain", "plain"},
		{"missing engine", "", "sqlite", "unsupported-engine", WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedLongV1, "plain", "plain"},
		{"unknown target engine", "sqlite", "private://secret", "unsupported-engine", WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedLongV1, "plain", "plain"},
		{"legacy wide source", "sqlite", "sqlite", "unsupported-layout", WriteGroupLayoutManagedWideV1, WriteGroupLayoutManagedLongV1, "", "plain"},
		{"legacy absent layout", "sqlite", "sqlite", "unsupported-layout", "", WriteGroupLayoutManagedLongV1, "", "plain"},
		{"unknown target layout", "sqlite", "sqlite", "unsupported-layout", WriteGroupLayoutManagedLongV1, "private://secret", "plain", ""},
		{"plain to Timescale", "postgres", "postgres", "storage-mode-change", WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedLongV1, "plain", "timescale"},
		{"Timescale to plain", "postgres", "postgres", "storage-mode-change", WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedLongV1, "timescale", "plain"},
		{"plain to native", "mysql", "mysql", "storage-mode-change", WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedLongV1, "plain", "native_monthly"},
		{"native to plain", "sqlserver", "sqlserver", "storage-mode-change", WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedLongV1, "native_monthly", "plain"},
		{"Timescale to wide", "postgres", "postgres", "storage-mode-change", WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedWideProfileV1, "timescale", ""},
		{"wide must store implicit plain", "sqlite", "sqlite", "unsupported-storage-mode", WriteGroupLayoutManagedWideProfileV1, WriteGroupLayoutManagedLongV1, "plain", "plain"},
		{"unknown mode", "sqlite", "sqlite", "unsupported-storage-mode", WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedLongV1, "private://secret", "private://secret"},
		{"SQLite Timescale", "sqlite", "sqlite", "unsupported-storage-mode", WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedLongV1, "timescale", "timescale"},
		{"MySQL Timescale", "mysql", "mysql", "unsupported-storage-mode", WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedLongV1, "timescale", "timescale"},
		{"Postgres native", "postgres", "postgres", "unsupported-storage-mode", WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedLongV1, "native_monthly", "native_monthly"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			source, target, types := transitionMatrixPair(tc.from, tc.to)
			source.StorageMode, target.StorageMode = tc.fromMode, tc.toMode
			// Predecessor-looking schema fields cannot make an unsupported row legal.
			source.Destination.SchemaDigest, target.Destination.SchemaDigest = "old-proof", "old-proof"
			source.Destination.SchemaRevision, target.Destination.SchemaRevision = "old-schema", "old-schema"
			requireTransitionReason(t, ValidateWriteGroupTransitionStructure(source, target, tc.fromEngine, tc.toEngine, types), tc.reason)
		})
	}
}

func TestWriteGroupTransitionMatrixRejectsIdentityAndMeaningChanges(t *testing.T) {
	for _, tc := range []struct {
		name, reason string
		mutate       func(*WriteGroup, *WriteGroup)
	}{
		{"custom target", "unsupported-layout", func(_, g *WriteGroup) { g.Destination.StorageStrategy = WriteGroupStorageStrategyCustom }},
		{"custom source", "unsupported-layout", func(g, _ *WriteGroup) { g.Destination.StorageStrategy = WriteGroupStorageStrategyCustom }},
		{"missing strategy", "unsupported-layout", func(_, g *WriteGroup) { g.Destination.StorageStrategy = "" }},
		{"foreign workspace", "invalid-transition", func(_, g *WriteGroup) { g.WorkspaceID = "foreign" }},
		{"missing workspace", "invalid-transition", func(g, h *WriteGroup) { g.WorkspaceID, h.WorkspaceID = "", "" }},
		{"missing profile", "unsupported-shape", func(_, g *WriteGroup) { g.RecordingProfile = nil }},
		{"kind change", "profile-kind-change", func(_, g *WriteGroup) { g.RecordingProfile.Kind = WriteGroupProfileCumulative }},
		{"reduction change", "recording-meaning-change", func(_, g *WriteGroup) { g.RecordingProfile.SampleReduction = WriteGroupSampleReductionMean }},
		{"interval change", "recording-meaning-change", func(_, g *WriteGroup) { g.RecordingProfile.Interval = 30 }},
		{"unit change", "recording-meaning-change", func(_, g *WriteGroup) { g.RecordingProfile.Unit = WriteGroupIntervalUnitMinute }},
		{"partial source", "partial-source-set", func(_, g *WriteGroup) { g.Members = nil }},
		{"extra source", "partial-source-set", func(_, g *WriteGroup) { g.Members = append(g.Members, g.Members[0]) }},
		{"empty source", "partial-source-set", func(g, h *WriteGroup) { g.Members, h.Members = nil, nil }},
		{"duplicate sources", "partial-source-set", func(g, h *WriteGroup) {
			g.Members = append(g.Members, g.Members[0])
			h.Members = append(h.Members, h.Members[0])
		}},
		{"device identity", "partial-source-set", func(_, g *WriteGroup) { g.Members[0].DeviceID = "other" }},
		{"point identity", "partial-source-set", func(_, g *WriteGroup) { g.Members[0].PointID = "other" }},
		{"tag identity", "partial-source-set", func(_, g *WriteGroup) { g.Members[0].TagID = "other" }},
		{"empty member identity", "partial-source-set", func(g, h *WriteGroup) { g.Members[0].DeviceID, h.Members[0].DeviceID = "", "" }},
		{"circuit identity", "recording-meaning-change", func(_, g *WriteGroup) { g.Members[0].CircuitID = "other" }},
		{"entity identity", "recording-meaning-change", func(_, g *WriteGroup) { g.Members[0].EntityKey = "other" }},
		{"source revision", "recording-meaning-change", func(_, g *WriteGroup) { g.Members[0].SourceRevision = "other" }},
		{"mapping revision", "recording-meaning-change", func(_, g *WriteGroup) { g.Members[0].MappingRevision = "other" }},
		{"role change", "recording-meaning-change", func(_, g *WriteGroup) { g.Members[0].ValueRole = WriteGroupValueRoleGauge }},
		{"measurement change", "recording-meaning-change", func(_, g *WriteGroup) { id := "other"; g.Members[0].MeasurementID = &id }},
		{"source type change", "recording-meaning-change", func(_, g *WriteGroup) { g.Members[0].SourceValueType = schema.DataTypeInt64 }},
		{"result type coercion", "recording-meaning-change", func(_, g *WriteGroup) { g.Members[0].ResultValueType = measurement.ExactInt64 }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			source, target, types := transitionMatrixPair(WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedWideProfileV1)
			tc.mutate(source, target)
			beforeSource, beforeTarget := cloneWriteGroup(source), cloneWriteGroup(target)
			requireTransitionReason(t, ValidateWriteGroupTransitionStructure(source, target, "sqlite", "sqlite", types), tc.reason)
			require.Equal(t, beforeSource, source)
			require.Equal(t, beforeTarget, target)
		})
	}
	g := transitionMatrixGroup(WriteGroupLayoutManagedLongV1)
	requireTransitionReason(t, ValidateWriteGroupTransitionStructure(nil, g, "sqlite", "sqlite", nil), "invalid-transition")
	requireTransitionReason(t, ValidateWriteGroupTransitionStructure(g, nil, "sqlite", "sqlite", nil), "invalid-transition")
}
