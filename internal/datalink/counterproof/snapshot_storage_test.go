package counterproof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSnapshotStorageProjectionIgnoresOnlyCounterBindingsAndPreservesInput(t *testing.T) {
	const original = `{"workspace_id":"w","counter_bindings":{"future":{"version":77,"meaning":[null]}},"runtime_layout":{"tag_types":{"電錶":"uint64"}},"group":{"id":"g","destination":{"table_name":"old_table"}}}`
	const expected = `{"group":{"id":"g","destination":{"table_name":"old_table"}},"runtime_layout":{"tag_types":{"電錶":"uint64"}},"workspace_id":"w"}`
	data := []byte(original)
	projected, err := SnapshotStorageProjection(data)
	require.NoError(t, err)
	require.Equal(t, expected, string(projected))
	require.Equal(t, original, string(data), "read-only projection does not repair or rewrite evidence")
	legacy, err := SnapshotStorageProjection([]byte(expected))
	require.NoError(t, err)
	require.Equal(t, expected, string(legacy))
}

func TestSnapshotStorageProjectionKeepsStrictEnvelopeRules(t *testing.T) {
	for _, raw := range []string{
		`null`, `[]`, `{"workspace_id":"w"`,
		`{"workspace_id":"w","workspace_id":"foreign"}`,
		`{"workspace_id":"w","Workspace_ID":"foreign"}`,
		`{"group":{"destination":{"table_name":"old","Table_Name":"foreign"}}}`,
		`{"group":{"members":[{"source_revision":"one","ſource_revision":"two"}]}}`,
	} {
		t.Run(raw, func(t *testing.T) {
			projected, err := SnapshotStorageProjection([]byte(raw))
			require.ErrorIs(t, err, ErrCounterBindingProofConflict)
			require.Nil(t, projected)
		})
	}
}
