package dbtarget

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLServerManagedRowProofClonePreservesFrozenEvidence(t *testing.T) {
	owner := "gateway"
	original := &SQLServerManagedRowProof{Group: &SQLServerManagedInspection{
		Inspection: &TableInspection{Columns: []ColumnInfo{{Name: "tag_id"}}, Indexes: []IndexInfo{{Name: "source_time", Columns: []string{"tag_id"}}}},
		Evidence: &SQLServerManagedTableEvidence{
			PrimaryKey: []string{"bucket_start", "record_id"}, ColumnCollations: map[string]string{"tag_id": "BIN2"},
			ColumnDefaults: map[string]*string{"owner": &owner},
		},
		Indexes: []IndexInfo{},
	}}
	cloned := CloneSQLServerManagedRowProof(original)
	cloned.Group.Inspection.Columns[0].Name = "changed"
	cloned.Group.Inspection.Indexes[0].Columns[0] = "changed"
	cloned.Group.Evidence.PrimaryKey[0] = "changed"
	cloned.Group.Evidence.ColumnCollations["tag_id"] = "CI"
	*cloned.Group.Evidence.ColumnDefaults["owner"] = "changed"
	require.Equal(t, "tag_id", original.Group.Inspection.Columns[0].Name)
	require.Equal(t, "tag_id", original.Group.Inspection.Indexes[0].Columns[0])
	require.Equal(t, "bucket_start", original.Group.Evidence.PrimaryKey[0])
	require.Equal(t, "BIN2", original.Group.Evidence.ColumnCollations["tag_id"])
	require.Equal(t, "gateway", *original.Group.Evidence.ColumnDefaults["owner"])
	require.NotNil(t, cloned.Group.Indexes, "verified empty index evidence must remain distinct from nil")
	require.Nil(t, cloned.Receipt)
	require.Nil(t, CloneSQLServerManagedRowProof(nil))
}
