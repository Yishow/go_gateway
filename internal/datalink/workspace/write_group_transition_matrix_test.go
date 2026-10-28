package workspace

import (
	"testing"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func transitionMatrixGroup(layout WriteGroupLayoutVersion) *WriteGroup {
	group := &WriteGroup{
		ID: "G1", WorkspaceID: "workspace-1", Revision: "R1", Status: WriteGroupStatusRunning,
		LayoutVersion: layout,
		RecordingProfile: &WriteGroupRecordingProfile{
			Kind: WriteGroupProfileInstantaneous, Interval: 15, Unit: WriteGroupIntervalUnitSecond,
			SampleReduction: WriteGroupSampleReductionLatest,
		},
		Members: []WriteGroupMember{{
			DeviceID: "device-1", PointID: "point-1", TagID: "tag-1", CircuitID: "circuit-1",
			SourceRevision: "source-1", MappingRevision: "mapping-1", TargetColumn: "state_value",
			ValueRole: WriteGroupValueRoleState, SourceValueType: schema.DataTypeBool, ResultValueType: measurement.ExactBool,
		}},
		Destination: WriteGroupDestination{
			ConnectorID: "A", ConnectorRevision: "A1", TableName: "history_g1",
			StorageStrategy: WriteGroupStorageStrategyManaged,
		},
	}
	if layout == WriteGroupLayoutManagedLongV1 {
		group.StorageMode = WriteGroupStorageModePlain
	}
	return group
}

func transitionMatrixPair(from, to WriteGroupLayoutVersion) (*WriteGroup, *WriteGroup, map[string]schema.DataType) {
	source := transitionMatrixGroup(from)
	target := transitionMatrixGroup(to)
	target.ID, target.Revision, target.Status = "G2", "R2", WriteGroupStatusDraft
	target.Destination.ConnectorID, target.Destination.ConnectorRevision, target.Destination.TableName = "B", "B1", "history_g2"
	return source, target, map[string]schema.DataType{"tag-1": schema.DataTypeBool}
}

func TestWriteGroupTransitionMatrixBidirectionalPlainState(t *testing.T) {
	// Spec example: SQLite G1 wide/plain instantaneous state → G2 long/plain,
	// latest meaning retained. This pure check cannot copy/delete SQL history.
	for _, engine := range []string{"sqlite", "postgres"} {
		for _, layouts := range [][2]WriteGroupLayoutVersion{
			{WriteGroupLayoutManagedWideProfileV1, WriteGroupLayoutManagedLongV1},
			{WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedWideProfileV1},
		} {
			t.Run(engine+"/"+string(layouts[0]), func(t *testing.T) {
				source, target, types := transitionMatrixPair(layouts[0], layouts[1])
				beforeSource, beforeTarget := cloneWriteGroup(source), cloneWriteGroup(target)
				require.NoError(t, ValidateWriteGroupTransitionStructure(source, target, engine, engine, types))
				require.Equal(t, beforeSource, source)
				require.Equal(t, beforeTarget, target)
				require.Equal(t, map[string]schema.DataType{"tag-1": schema.DataTypeBool}, types)
			})
		}
	}
}

func TestWriteGroupTransitionMatrixDestinationReplacement(t *testing.T) {
	for _, tc := range []struct {
		engine string
		mode   WriteGroupStorageMode
	}{
		{"sqlite", WriteGroupStorageModePlain}, {"postgres", WriteGroupStorageModePlain},
		{"mysql", WriteGroupStorageModePlain}, {"mariadb", WriteGroupStorageModePlain}, {"sqlserver", WriteGroupStorageModePlain},
		{"postgres", WriteGroupStorageModeTimescale}, {"mysql", WriteGroupStorageModeNativeMonthly},
		{"mariadb", WriteGroupStorageModeNativeMonthly}, {"sqlserver", WriteGroupStorageModeNativeMonthly},
	} {
		t.Run(tc.engine+"/"+string(tc.mode), func(t *testing.T) {
			source, target, types := transitionMatrixPair(WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedLongV1)
			source.StorageMode, target.StorageMode = tc.mode, tc.mode
			// MariaDB spec example retains A/A1 and B/B1 as distinct destinations.
			// Nil means only eligible to seek NEW target proof, not verified or ready.
			require.NoError(t, ValidateWriteGroupTransitionStructure(source, target, tc.engine, tc.engine, types))
			require.Empty(t, target.Destination.SchemaDigest)
			require.Equal(t, WriteGroupStatusDraft, target.Status)
		})
	}
	for _, engine := range []string{"sqlite", "postgres"} {
		source, target, types := transitionMatrixPair(WriteGroupLayoutManagedWideProfileV1, WriteGroupLayoutManagedWideProfileV1)
		require.NoError(t, ValidateWriteGroupTransitionStructure(source, target, engine, engine, types))
	}
}

func TestWriteGroupTransitionMatrixEffectiveMeaningAndPresentation(t *testing.T) {
	source, target, types := transitionMatrixPair(WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedWideProfileV1)
	source.StorageMode, source.RecordingProfile.SampleReduction = "", ""
	source.Status = WriteGroupStatusDeleted // Structural eligibility is not recreate authorization.
	target.Name, target.Members[0].CircuitName, target.Members[0].TargetColumn = "renamed", "renamed circuit", "new_column"
	require.NoError(t, ValidateWriteGroupTransitionStructure(source, target, "sqlite", "sqlite", types))
}

// Matrix eligibility deliberately does not authorize a changed runtime policy.
func TestWriteGroupTransitionMatrixPolicyValidationRemainsSeparate(t *testing.T) {
	source, target, types := transitionMatrixPair(WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedWideProfileV1)
	maxAge := 30
	target.Members[0].Required, target.Members[0].MaxAgeSeconds = true, &maxAge
	target.RowPolicy.AllowedLatenessSeconds, target.RowPolicy.IncompletePolicy = 5, "partial"
	target.RowPolicy.ValueColumn = "new_value_column"
	require.NoError(t, ValidateWriteGroupTransitionStructure(source, target, "sqlite", "sqlite", types))
	require.Equal(t, WriteGroupStatusDraft, target.Status)
	require.Empty(t, target.AppliedRevision)
}
