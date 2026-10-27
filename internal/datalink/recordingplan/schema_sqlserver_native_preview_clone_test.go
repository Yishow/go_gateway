package recordingplan

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLServerNativePreviewCannotMutateStoredPartition(t *testing.T) {
	newToken := func() *SchemaPreviewToken {
		return &SchemaPreviewToken{
			Token: "native-preview",
			Tables: []SchemaPreviewTable{{
				Name: "gw_long_fixture",
				SQLServerPartition: &SQLServerPartitionPreview{
					Action: "create", FunctionName: "pf_gw_long_fixture",
					SchemeName: "ps_gw_long_fixture", Filegroup: "PRIMARY",
					Partitions: []SQLServerPartitionBoundary{{
						Name: "p_past", LessThan: "2026-10-01T00:00:00.000000", Ordinal: 1,
					}},
				},
			}},
		}
	}
	mutate := func(token *SchemaPreviewToken) {
		token.Tables[0].SQLServerPartition.FunctionName = "foreign_function"
		token.Tables[0].SQLServerPartition.Partitions[0].LessThan = "2027-10-01T00:00:00.000000"
	}
	assertFrozen := func(t *testing.T, repo *MemoryRepository) {
		t.Helper()
		stored, err := repo.GetPreviewToken(t.Context(), "native-preview")
		require.NoError(t, err)
		require.Equal(t, "pf_gw_long_fixture", stored.Tables[0].SQLServerPartition.FunctionName)
		require.Equal(t, "2026-10-01T00:00:00.000000", stored.Tables[0].SQLServerPartition.Partitions[0].LessThan)
	}
	t.Run("save freezes caller owned preview", func(t *testing.T) {
		repo, token := NewMemoryRepository(), newToken()
		require.NoError(t, repo.SavePreviewToken(t.Context(), token))
		mutate(token)
		assertFrozen(t, repo)
	})
	t.Run("read does not expose stored confirmation", func(t *testing.T) {
		repo := NewMemoryRepository()
		require.NoError(t, repo.SavePreviewToken(t.Context(), newToken()))
		returned, err := repo.GetPreviewToken(t.Context(), "native-preview")
		require.NoError(t, err)
		mutate(returned)
		assertFrozen(t, repo)
	})
}
