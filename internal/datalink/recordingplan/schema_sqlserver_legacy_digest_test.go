package recordingplan

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLServerEvidencePreservesLegacyCatalogDigest(t *testing.T) {
	// This is the serialized inspection contract before the SQL Server field.
	// Existing MySQL operation manifests hash the exact bytes, including nulls.
	const legacy = `{"Status":"missing","Columns":null,"ColumnTypes":null,"ColumnNullable":null,"ColumnPrimaryKey":null,"Indexes":null,"IndexesInspected":false,"Timescale":null,"MySQL":null}`
	want := sha256.Sum256([]byte(legacy))
	got, err := mysqlCatalogDigest(TargetTableInspection{Status: inspectionMissing})
	require.NoError(t, err)
	require.Equal(t, hex.EncodeToString(want[:]), got, "a new provider must not change existing catalog fingerprints")
}
