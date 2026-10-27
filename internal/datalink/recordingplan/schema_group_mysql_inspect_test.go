package recordingplan

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func mysqlCatalogFixture(t *testing.T) (GroupSchemaLayout, *SchemaPreviewToken, TargetTableInspection, []string) {
	t.Helper()
	layout, err := ManagedMySQLLongGroupLayout("gw_catalog_test", "_gw_owner_0123456789abcdef01234567", MySQLGroupSchemaOptions{
		ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "cap-1",
		Engine: "InnoDB", Charset: "utf8mb4", Collation: "utf8mb4_0900_bin",
	})
	require.NoError(t, err)
	statements, err := GenerateGroupSchemaDDL("mysql", layout)
	require.NoError(t, err)
	state := TargetTableInspection{
		Status: inspectionExists, ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{},
		ColumnPrimaryKey: map[string]bool{}, IndexesInspected: true,
		MySQL: &TargetMySQLTableEvidence{
			ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "cap-1", Database: "recordings",
			Engine: "InnoDB", Charset: "utf8mb4", Collation: "utf8mb4_0900_bin", TableType: "BASE TABLE",
			PrimaryKey: []string{"bucket_start", "record_id"}, ColumnCollations: map[string]string{},
			ColumnDefaults: map[string]*string{}, JSONColumns: map[string]bool{"provenance": true},
		},
	}
	for _, column := range layout.Columns {
		state.Columns = append(state.Columns, column.Name)
		state.ColumnTypes[column.Name] = strings.ToLower(column.SQLType)
		state.ColumnNullable[column.Name] = column.Nullable
		state.ColumnPrimaryKey[column.Name] = column.PrimaryKey
		if strings.HasPrefix(column.SQLType, "VARCHAR(") || strings.HasPrefix(column.SQLType, "CHAR(") || column.SQLType == "LONGTEXT" {
			state.MySQL.ColumnCollations[column.Name] = "utf8mb4_0900_bin"
		}
	}
	state.ColumnTypes["value_bool"] = "tinyint(1)"
	ownerDefault := "gateway"
	state.MySQL.ColumnDefaults[layout.OwnerColumn] = &ownerDefault
	for _, index := range layout.Indexes {
		state.Indexes = append(state.Indexes, TargetIndexInspection{Name: index.Name, Columns: index.Columns})
	}
	token := &SchemaPreviewToken{Dialect: "mysql", Database: "recordings", Schema: "recordings", GroupLayout: &layout}
	return layout, token, state, statements
}

func mysqlStepReference(kind, name, sql string) mysqlApplyStepRef {
	return mysqlApplyStepRef{Kind: kind, Name: name, SQLHash: fmt.Sprintf("%x", sha256.Sum256([]byte(sql)))}
}

func TestMySQLCatalogStepRequiresExactEvidence(t *testing.T) {
	layout, token, state, statements := mysqlCatalogFixture(t)
	step := mysqlStepReference("group_table", layout.TableName, statements[0])
	inspect := func(context.Context, string) (TargetTableInspection, error) { return state, nil }
	evidence, err := inspectMySQLApplyStep(t.Context(), token, layout, step, inspect)
	require.NoError(t, err)
	require.Equal(t, "exact", evidence.State)
	require.NotEmpty(t, evidence.Fingerprint)

	mutations := map[string]func(*TargetTableInspection){
		"wrong engine":       func(s *TargetTableInspection) { s.MySQL.Engine = "MyISAM" },
		"case folding":       func(s *TargetTableInspection) { s.MySQL.ColumnCollations["tag_id"] = "utf8mb4_0900_ai_ci" },
		"lost scale":         func(s *TargetTableInspection) { s.ColumnTypes["value_decimal"] = "decimal(65,29)" },
		"unsigned int":       func(s *TargetTableInspection) { s.ColumnTypes["value_int"] = "bigint unsigned" },
		"owner default":      func(s *TargetTableInspection) { delete(s.MySQL.ColumnDefaults, layout.OwnerColumn) },
		"no JSON constraint": func(s *TargetTableInspection) { s.MySQL.JSONColumns["provenance"] = false },
		"reversed key":       func(s *TargetTableInspection) { s.MySQL.PrimaryKey = []string{"record_id", "bucket_start"} },
		"view":               func(s *TargetTableInspection) { s.MySQL.TableType = "VIEW" },
		"family drift":       func(s *TargetTableInspection) { s.MySQL.ServerFamily = "mariadb" },
		"capability drift":   func(s *TargetTableInspection) { s.MySQL.CapabilityRevision = "cap-2" },
		"missing proof":      func(s *TargetTableInspection) { s.MySQL = nil },
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			_, _, bad, _ := mysqlCatalogFixture(t)
			mutate(&bad)
			_, err := inspectMySQLApplyStep(t.Context(), token, layout, step, func(context.Context, string) (TargetTableInspection, error) { return bad, nil })
			require.Error(t, err)
		})
	}
}

func TestMySQLCatalogIndexCanBePendingWithoutRecreatingTable(t *testing.T) {
	layout, token, state, statements := mysqlCatalogFixture(t)
	state.Indexes = nil
	inspect := func(context.Context, string) (TargetTableInspection, error) { return state, nil }
	table, err := inspectMySQLApplyStep(t.Context(), token, layout, mysqlStepReference("group_table", layout.TableName, statements[0]), inspect)
	require.NoError(t, err)
	require.Equal(t, "exact", table.State)
	step := mysqlStepReference("index", layout.Indexes[0].Name, statements[1])
	index, err := inspectMySQLApplyStep(t.Context(), token, layout, step, inspect)
	require.NoError(t, err)
	require.Equal(t, "missing", index.State)
	state.Indexes = []TargetIndexInspection{{Name: layout.Indexes[0].Name, Columns: layout.Indexes[0].Columns}}
	after, err := inspectMySQLApplyStep(t.Context(), token, layout, mysqlStepReference("group_table", layout.TableName, statements[0]), inspect)
	require.NoError(t, err)
	require.Equal(t, table.Fingerprint, after.Fingerprint, "later index step does not alter the verified table step")
	state.Indexes = []TargetIndexInspection{{Name: layout.Indexes[0].Name, Columns: []string{"tag_id"}}}
	_, err = inspectMySQLApplyStep(t.Context(), token, layout, step, inspect)
	require.ErrorIs(t, err, ErrIncompatibleExistingTable)
	step.SQLHash = "wrong"
	_, err = inspectMySQLApplyStep(t.Context(), token, layout, step, inspect)
	require.ErrorIs(t, err, ErrPreviewTokenStale)
}

func TestMySQLCatalogPreviewResumesOnlyMissingIndexAndReceipt(t *testing.T) {
	layout, token, state, statements := mysqlCatalogFixture(t)
	state.Indexes = nil
	inspect := func(_ context.Context, name string) (TargetTableInspection, error) {
		if name == layout.TableName {
			return state, nil
		}
		return TargetTableInspection{Status: inspectionMissing, MySQL: state.MySQL}, nil
	}
	actual, tables, _, err := prepareGroupSchemaPreview(t.Context(), SchemaPreviewScope{Dialect: "mysql", Database: token.Database, Schema: token.Schema, GroupLayout: &layout}, inspect)
	require.NoError(t, err)
	require.Equal(t, statements[1:], actual, "never recreate the already verified table")
	require.Equal(t, "create_indexes", tables[0].Action)
}

func TestMySQLCatalogForbiddenIsNeverMissing(t *testing.T) {
	layout, token, state, statements := mysqlCatalogFixture(t)
	step := mysqlStepReference("group_table", layout.TableName, statements[0])
	for _, status := range []string{inspectionForbidden, "failed", ""} {
		state.Status = status
		_, err := inspectMySQLApplyStep(t.Context(), token, layout, step, func(context.Context, string) (TargetTableInspection, error) { return state, nil })
		require.ErrorIs(t, err, ErrTargetInspectionUnconfirmed)
	}
	state.Status = inspectionMissing
	evidence, err := inspectMySQLApplyStep(t.Context(), token, layout, step, func(context.Context, string) (TargetTableInspection, error) { return state, nil })
	require.NoError(t, err)
	require.Equal(t, "missing", evidence.State)
}

func TestMySQLPendingIndexPreviewCannotMutateStoredConfirmation(t *testing.T) {
	repo := NewMemoryRepository()
	token := &SchemaPreviewToken{Token: "stored-index-preview", Tables: []SchemaPreviewTable{{
		Name: "gw_table", Action: "create_indexes", Indexes: []GroupSchemaIndex{{Name: "idx_src", Columns: []string{"workspace_id", "tag_id", "bucket_start"}}},
	}}}
	require.NoError(t, repo.SavePreviewToken(t.Context(), token))
	token.Tables[0].Indexes[0].Columns[0] = "caller-input"
	loaded, err := repo.GetPreviewToken(t.Context(), token.Token)
	require.NoError(t, err)
	require.Equal(t, "workspace_id", loaded.Tables[0].Indexes[0].Columns[0])
	loaded.Tables[0].Indexes[0].Name = "caller-readback"
	again, err := repo.GetPreviewToken(t.Context(), token.Token)
	require.NoError(t, err)
	require.Equal(t, "idx_src", again.Tables[0].Indexes[0].Name)
}

func TestMySQLCatalogMariaDBIntegerDisplayWidthKeepsSignedRange(t *testing.T) {
	layout, token, state, _ := mysqlCatalogFixture(t)
	options := *layout.MySQL
	options.ServerFamily, options.ServerVersion, options.Collation = "mariadb", "11.4.13-MariaDB-ubu2404", "utf8mb4_nopad_bin"
	layout, err := ManagedMySQLLongGroupLayout(layout.TableName, layout.OwnerColumn, options)
	require.NoError(t, err)
	token.GroupLayout = &layout
	state.MySQL.ServerFamily, state.MySQL.ServerVersion, state.MySQL.Collation = options.ServerFamily, options.ServerVersion, options.Collation
	for name := range state.MySQL.ColumnCollations {
		state.MySQL.ColumnCollations[name] = options.Collation
	}
	state.ColumnTypes["value_int"] = "bigint(20)"
	state.ColumnTypes["provenance"] = "longtext"
	state.MySQL.ColumnCollations["provenance"] = "utf8mb4_bin"
	statements, err := GenerateGroupSchemaDDL("mysql", layout)
	require.NoError(t, err)
	step := mysqlStepReference("group_table", layout.TableName, statements[0])
	inspect := func(context.Context, string) (TargetTableInspection, error) { return state, nil }
	actual, err := inspectMySQLApplyStep(t.Context(), token, layout, step, inspect)
	require.NoError(t, err)
	require.Equal(t, "exact", actual.State)
	state.ColumnTypes["value_int"] = "bigint(20) unsigned"
	_, err = inspectMySQLApplyStep(t.Context(), token, layout, step, inspect)
	require.ErrorIs(t, err, ErrIncompatibleExistingTable)
	state.ColumnTypes["value_int"] = "bigint(20)"
	state.MySQL.JSONColumns["provenance"] = false
	_, err = inspectMySQLApplyStep(t.Context(), token, layout, step, inspect)
	require.ErrorIs(t, err, ErrIncompatibleExistingTable)
}

func TestMySQLSchemaLeaseStillConflictsWhenCapabilityFactsChange(t *testing.T) {
	layout, _, _, _ := mysqlCatalogFixture(t)
	first := mysqlTestToken(layout)
	first.Token, first.OperationID = "scope-first-token", "scope-first-op"
	second := first
	second.Token, second.OperationID = "scope-second-token", "scope-second-op"
	changed := layout
	options := *layout.MySQL
	options.CapabilityRevision = "cap-refreshed"
	changed.MySQL = &options
	second.GroupLayout = &changed
	var err error
	first.Digest, err = previewDigest(&first)
	require.NoError(t, err)
	second.Digest, err = previewDigest(&second)
	require.NoError(t, err)
	require.Equal(t, first.ScopeKey(), second.ScopeKey(), "the physical DDL scope cannot change with a capability projection")
	service := NewService(NewMemoryRepository())
	_, outcome, err := service.ClaimSchemaApply(t.Context(), &first)
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)
	op, outcome, err := service.ClaimSchemaApply(t.Context(), &second)
	require.NoError(t, err)
	require.Equal(t, ClaimScopeBusy, outcome)
	require.Equal(t, first.OperationID, op.OperationID)
}
