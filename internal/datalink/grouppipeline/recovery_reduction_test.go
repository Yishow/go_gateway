package grouppipeline

import (
	"crypto/sha256"
	"encoding/json"
	"strings"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// canonicalProof builds the group and frozen layout an Apply of a legacy
// managed wide group leaves behind, minimal but complete for its proof.
func canonicalProof(t *testing.T) (*workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout) {
	t.Helper()
	group := groupAt("rev-1")
	group.Destination.StorageStrategy = workspace.WriteGroupStorageStrategyManaged
	group.Destination.TableName = "managed_readings"
	group.RowPolicy.RecordKeyColumn = "record_id"
	group.RowPolicy.GroupIDColumn = "group_id"
	group.RowPolicy.BucketStartColumn = "bucket_start"
	group.RowPolicy.ProvenanceColumn = "provenance"
	group.RecordingProfile = &workspace.WriteGroupRecordingProfile{Kind: workspace.WriteGroupProfileInstantaneous, Interval: 15, Unit: workspace.WriteGroupIntervalUnitSecond}
	group.RowPolicy.IntervalSeconds = 15
	owner := managedOwnerColumnForTest(group)
	layout := &workspace.WriteGroupRuntimeLayout{
		Dialect: dbtarget.SQLDialectSQLite,
		Columns: []dbtarget.ColumnInfo{
			{Name: "record_id", DataType: "TEXT", Nullable: false, PrimaryKey: true},
			{Name: "group_id", DataType: "TEXT", Nullable: false},
			{Name: "bucket_start", DataType: "TEXT", Nullable: false},
			{Name: "provenance", DataType: "TEXT", Nullable: false},
			{Name: "temperature", DataType: "REAL", Nullable: true},
			{Name: "pressure", DataType: "INTEGER", Nullable: true},
			{Name: owner, DataType: "TEXT", Nullable: false},
		},
		TagTypes: map[string]schema.DataType{
			"tag-t": schema.DataTypeFloat64,
			"tag-p": schema.DataTypeInt64,
		},
		ReceiptTableReady: true,
		SchemaDigest:      strings.Repeat("a", sha256.Size*2),
	}
	return group, layout
}

// injectMeanReduction decodes an explicit mean reduction into a group's saved
// profile, the same path every persisted payload takes.
func injectMeanReduction(t *testing.T, group *workspace.WriteGroup) {
	t.Helper()
	require.NotNil(t, group.RecordingProfile)
	require.NoError(t, json.Unmarshal(
		[]byte(`{"kind":"instantaneous","interval":15,"unit":"second","sample_reduction":"mean"}`),
		group.RecordingProfile))
}

// TestRecordingReductionCompatibility_RestoreProofRejectsMean pins the restore
// gate: a frozen descriptor whose group carries an explicit mean never proves a
// writer, because no runtime serves averages yet. A legacy group without a
// reduction keeps restoring.
func TestRecordingReductionCompatibility_RestoreProofRejectsMean(t *testing.T) {
	t.Run("canonical managed proof", func(t *testing.T) {
		group, layout := canonicalProof(t)
		require.True(t, validCanonicalManagedRuntimeProof(group, layout))
		injectMeanReduction(t, group)
		require.False(t, validCanonicalManagedRuntimeProof(group, layout),
			"a mean descriptor must fail closed instead of restoring a latest writer")
	})

	t.Run("long proof", func(t *testing.T) {
		group, layout := longProof(t, workspace.WriteGroupStorageModePlain, nil)
		require.True(t, validLongRuntimeProof(group, layout))
		injectMeanReduction(t, group)
		require.False(t, validLongRuntimeProof(group, layout),
			"a mean descriptor must fail closed instead of restoring a raw long writer")
	})
}

// TestRecordingReductionCompatibility_NewShapeProofsReject pins that a frozen
// descriptor of an unproven new shape never restores a writer: the new wide
// profile layout is not a legacy wide writer, and member metadata without a
// runtime proof blocks recovery.
func TestRecordingReductionCompatibility_NewShapeProofsReject(t *testing.T) {
	const newWideProfileLayout = workspace.WriteGroupLayoutVersion("managed_wide_profile_v1")

	t.Run("the new wide profile layout is not a legacy wide proof", func(t *testing.T) {
		group, layout := canonicalProof(t)
		require.Equal(t, workspace.WriteGroupLayoutManagedWideV1, group.EffectiveLayoutVersion())
		require.True(t, validCanonicalManagedRuntimeProof(group, layout))
		group.LayoutVersion = newWideProfileLayout
		require.False(t, validCanonicalManagedRuntimeProof(group, layout),
			"an unproven new wide descriptor must not start as a legacy wide writer")
	})

	t.Run("member metadata without proof is invalid", func(t *testing.T) {
		group, layout := canonicalProof(t)
		require.True(t, validCanonicalManagedRuntimeProof(group, layout))
		group.Members[0] = injectMemberRole(t, group.Members[0], "gauge")
		require.False(t, validCanonicalManagedRuntimeProof(group, layout),
			"member metadata has no runtime proof and must not restore a writer")
	})

	t.Run("a long descriptor with member metadata is invalid", func(t *testing.T) {
		group, layout := longProof(t, workspace.WriteGroupStorageModePlain, nil)
		require.True(t, validLongRuntimeProof(group, layout))
		group.Members[0] = injectMemberRole(t, group.Members[0], "gauge")
		require.False(t, validLongRuntimeProof(group, layout))
	})
}

// injectMemberRole decodes an explicit role into a member's wire form.
func injectMemberRole(t *testing.T, member workspace.WriteGroupMember, role string) workspace.WriteGroupMember {
	t.Helper()
	raw, err := json.Marshal(member)
	require.NoError(t, err)
	var payload map[string]any
	require.NoError(t, json.Unmarshal(raw, &payload))
	payload["value_role"] = role
	encoded, err := json.Marshal(payload)
	require.NoError(t, err)
	var decoded workspace.WriteGroupMember
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	return decoded
}
