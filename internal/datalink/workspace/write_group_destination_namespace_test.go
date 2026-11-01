package workspace

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWriteGroupDestinationNamespaceDefaultsPreserveResolverContract(t *testing.T) {
	for _, tc := range []struct{ kind, want string }{{"sqlite", "main"}, {"postgres", "public"}, {"sqlserver", "dbo"}} {
		t.Run(tc.kind, func(t *testing.T) {
			destination := WriteGroupDestination{}
			require.NoError(t, resolveWriteGroupDestinationScope(tc.kind, `{}`, &destination))
			require.Equal(t, tc.want, destination.TableSchema)
			destination = WriteGroupDestination{TableSchema: "  explicit  "}
			require.NoError(t, resolveWriteGroupDestinationScope(tc.kind, `{}`, &destination))
			require.Equal(t, "explicit", destination.TableSchema)
			require.NoError(t, resolveWriteGroupDestinationScope(tc.kind, `{"schema":" stored "}`, &destination))
			require.Equal(t, "stored", destination.TableSchema)
			require.Error(t, resolveWriteGroupDestinationScope(tc.kind, `{invalid`, &destination))
		})
	}
	destination := WriteGroupDestination{Database: "foreign", TableSchema: "foreign"}
	require.Error(t, resolveWriteGroupDestinationScope("mysql", `{"database":"owned"}`, &destination))
}
