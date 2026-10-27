package recordingplan

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLServerPreviewClonePreservesFrozenScope(t *testing.T) {
	original := SchemaPreviewToken{GroupLayout: &GroupSchemaLayout{
		SQLServer: &SQLServerGroupSchemaOptions{Schema: "dbo", Collation: "Latin1_General_100_BIN2"},
	}}
	cloned := clonePreviewToken(original)
	cloned.GroupLayout.SQLServer.Schema = "other_schema"
	require.Equal(t, "dbo", original.GroupLayout.SQLServer.Schema, "editing a returned preview must not mutate its frozen schema")
}

func TestSQLServerLayoutEqualityBindsSchema(t *testing.T) {
	left := GroupSchemaLayout{TableName: "test_long", OwnerColumn: "_gw_owner_0123456789abcdef01234567",
		Columns:   []GroupSchemaColumn{{Name: "_gw_owner_0123456789abcdef01234567", SQLType: "NVARCHAR(16)"}},
		SQLServer: &SQLServerGroupSchemaOptions{Schema: "dbo", Collation: "Latin1_General_100_BIN2"},
	}
	right := left
	options := *left.SQLServer
	options.Schema = "another_schema"
	right.SQLServer = &options
	require.False(t, groupSchemaLayoutsEqual(left, right), "different SQL Server schemas cannot share confirmation layout identity")
}
