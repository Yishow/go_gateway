package recordingplan

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func validTestSQLServerOptions() SQLServerGroupSchemaOptions {
	return SQLServerGroupSchemaOptions{
		Database:           "recording",
		Schema:             "dbo",
		Collation:          sqlserverCollation,
		ServerVersion:      "16.0.4180.1",
		ServerEdition:      "Enterprise Edition: Core-based Licensing (64-bit)",
		CapabilityRevision: "rev-2026",
		EngineEdition:      3,
		CompatibilityLevel: 160,
	}
}

func validTestSQLServerTableInspection(layout GroupSchemaLayout, options SQLServerGroupSchemaOptions) TargetTableInspection {
	cols := make([]string, len(layout.Columns))
	types := make(map[string]string, len(layout.Columns))
	nullable := make(map[string]bool, len(layout.Columns))
	pks := make(map[string]bool, len(layout.Columns))
	collations := make(map[string]string)
	defaults := make(map[string]*string)

	for i, col := range layout.Columns {
		cols[i] = col.Name
		types[col.Name] = col.SQLType
		nullable[col.Name] = col.Nullable
		pks[col.Name] = col.PrimaryKey
		if sqlserverCollatedSQLType(col.SQLType) {
			collations[col.Name] = options.Collation
		}
		if col.Name == layout.OwnerColumn {
			d := managedOwnerDefault
			defaults[col.Name] = &d
		}
	}

	return TargetTableInspection{
		Status:           inspectionExists,
		Columns:          cols,
		ColumnTypes:      types,
		ColumnNullable:   nullable,
		ColumnPrimaryKey: pks,
		Indexes: []TargetIndexInspection{{
			Name:    "idx_" + layout.TableName + "_src_time",
			Columns: []string{"workspace_id", "tag_id", "bucket_start"},
			Unique:  false,
		}},
		IndexesInspected: true,
		SQLServer: &TargetSQLServerTableEvidence{
			Database:              options.Database,
			Schema:                options.Schema,
			ServerVersion:         options.ServerVersion,
			ServerEdition:         options.ServerEdition,
			CapabilityRevision:    options.CapabilityRevision,
			EngineEdition:         options.EngineEdition,
			CompatibilityLevel:    options.CompatibilityLevel,
			MetadataVerified:      true,
			TableType:             "USER_TABLE",
			ObjectID:              1001,
			Plain:                 true,
			PrimaryKey:            []string{"bucket_start", "record_id"},
			ColumnCollations:      collations,
			ColumnDefaults:        defaults,
			ProvenanceJSONChecked: true,
		},
	}
}

func validTestSQLServerReceiptInspection(options SQLServerGroupSchemaOptions) TargetTableInspection {
	receiptCols := sqlserverReceiptColumns()
	cols := make([]string, len(receiptCols))
	types := make(map[string]string, len(receiptCols))
	nullable := make(map[string]bool, len(receiptCols))
	pks := make(map[string]bool, len(receiptCols))
	collations := make(map[string]string)
	defaults := make(map[string]*string)

	for i, col := range receiptCols {
		cols[i] = col.Name
		types[col.Name] = col.SQLType
		nullable[col.Name] = col.Nullable
		pks[col.Name] = col.PrimaryKey
		if sqlserverCollatedSQLType(col.SQLType) {
			collations[col.Name] = options.Collation
		}
		if col.Name == managedEffectReceiptMarker {
			d := managedOwnerDefault
			defaults[col.Name] = &d
		}
	}

	return TargetTableInspection{
		Status:           inspectionExists,
		Columns:          cols,
		ColumnTypes:      types,
		ColumnNullable:   nullable,
		ColumnPrimaryKey: pks,
		IndexesInspected: true,
		SQLServer: &TargetSQLServerTableEvidence{
			Database:           options.Database,
			Schema:             options.Schema,
			ServerVersion:      options.ServerVersion,
			ServerEdition:      options.ServerEdition,
			CapabilityRevision: options.CapabilityRevision,
			EngineEdition:      options.EngineEdition,
			CompatibilityLevel: options.CompatibilityLevel,
			MetadataVerified:   true,
			TableType:          "USER_TABLE",
			ObjectID:           1002,
			Plain:              true,
			PrimaryKey:         []string{"effect_key"},
			ColumnCollations:   collations,
			ColumnDefaults:     defaults,
		},
	}
}

func TestSQLServerForeignOwnerDenied(t *testing.T) {
	opts := validTestSQLServerOptions()
	layout, err := ManagedSQLServerLongGroupLayout("gw_long_test", "_gw_owner_0123456789abcdef01234567", opts)
	require.NoError(t, err)

	t.Run("owner_foreign_default_value", func(t *testing.T) {
		insp := validTestSQLServerTableInspection(layout, opts)
		foreign := "foreign_app"
		insp.SQLServer.ColumnDefaults[layout.OwnerColumn] = &foreign
		err := validateSQLServerCatalogTable(insp, layout.Columns, layout.OwnerColumn, opts)
		require.ErrorIs(t, err, ErrIncompatibleExistingTable)
	})

	t.Run("owner_missing_default", func(t *testing.T) {
		insp := validTestSQLServerTableInspection(layout, opts)
		delete(insp.SQLServer.ColumnDefaults, layout.OwnerColumn)
		err := validateSQLServerCatalogTable(insp, layout.Columns, layout.OwnerColumn, opts)
		require.ErrorIs(t, err, ErrIncompatibleExistingTable)
	})

	t.Run("unexpected_default_on_value_column", func(t *testing.T) {
		insp := validTestSQLServerTableInspection(layout, opts)
		unexpected := "0"
		insp.SQLServer.ColumnDefaults["value_int"] = &unexpected
		err := validateSQLServerCatalogTable(insp, layout.Columns, layout.OwnerColumn, opts)
		require.ErrorIs(t, err, ErrIncompatibleExistingTable)
	})
}

func TestSQLServerCatalogInspection(t *testing.T) {
	opts := validTestSQLServerOptions()
	layout, err := ManagedSQLServerLongGroupLayout("gw_long_test", "_gw_owner_0123456789abcdef01234567", opts)
	require.NoError(t, err)

	t.Run("missing_proof_cannot_become_create", func(t *testing.T) {
		// Missing table without SQLServer evidence
		group := TargetTableInspection{Status: inspectionMissing}
		receipt := validTestSQLServerReceiptInspection(opts)
		_, err := inspectSQLServerGroupSnapshot(layout, opts.Database, group, receipt)
		require.ErrorIs(t, err, ErrTargetInspectionUnconfirmed)

		// Missing table with MetadataVerified = false
		group = TargetTableInspection{
			Status: inspectionMissing,
			SQLServer: &TargetSQLServerTableEvidence{
				Database:         opts.Database,
				Schema:           opts.Schema,
				MetadataVerified: false,
			},
		}
		_, err = inspectSQLServerGroupSnapshot(layout, opts.Database, group, receipt)
		require.ErrorIs(t, err, ErrTargetInspectionUnconfirmed)
	})

	t.Run("foreign_collation_rejected", func(t *testing.T) {
		insp := validTestSQLServerTableInspection(layout, opts)
		insp.SQLServer.ColumnCollations["tag_id"] = "SQL_Latin1_General_CP1_CI_AS"
		err := validateSQLServerCatalogTable(insp, layout.Columns, layout.OwnerColumn, opts)
		require.ErrorIs(t, err, ErrIncompatibleExistingTable)
	})

	t.Run("foreign_pk_rejected", func(t *testing.T) {
		insp := validTestSQLServerTableInspection(layout, opts)
		insp.SQLServer.PrimaryKey = []string{"record_id"}
		err := validateSQLServerCatalogTable(insp, layout.Columns, layout.OwnerColumn, opts)
		require.ErrorIs(t, err, ErrIncompatibleExistingTable)

		insp.SQLServer.PrimaryKey = []string{"record_id", "bucket_start"} // wrong order
		err = validateSQLServerCatalogTable(insp, layout.Columns, layout.OwnerColumn, opts)
		require.ErrorIs(t, err, ErrIncompatibleExistingTable)
	})

	t.Run("foreign_index_rejected", func(t *testing.T) {
		// Missing index
		insp := validTestSQLServerTableInspection(layout, opts)
		insp.Indexes = nil
		err := validateSQLServerCatalogIndexes(insp, layout.Indexes)
		require.ErrorIs(t, err, ErrIncompatibleExistingTable)

		// Extra index
		insp = validTestSQLServerTableInspection(layout, opts)
		insp.Indexes = append(insp.Indexes, TargetIndexInspection{Name: "idx_extra", Columns: []string{"tag_id"}})
		err = validateSQLServerCatalogIndexes(insp, layout.Indexes)
		require.ErrorIs(t, err, ErrIncompatibleExistingTable)

		// Unique index
		insp = validTestSQLServerTableInspection(layout, opts)
		insp.Indexes[0].Unique = true
		err = validateSQLServerCatalogIndexes(insp, layout.Indexes)
		require.ErrorIs(t, err, ErrIncompatibleExistingTable)

		// Indexes not inspected
		insp = validTestSQLServerTableInspection(layout, opts)
		insp.IndexesInspected = false
		err = validateSQLServerCatalogIndexes(insp, layout.Indexes)
		require.ErrorIs(t, err, ErrTargetInspectionUnconfirmed)
	})

	t.Run("foreign_json_rejected", func(t *testing.T) {
		insp := validTestSQLServerTableInspection(layout, opts)
		insp.SQLServer.ProvenanceJSONChecked = false
		err := validateSQLServerCatalogTable(insp, layout.Columns, layout.OwnerColumn, opts)
		require.ErrorIs(t, err, ErrIncompatibleExistingTable)
	})

	t.Run("mismatched_server_facts_rejected", func(t *testing.T) {
		insp := validTestSQLServerTableInspection(layout, opts)
		insp.SQLServer.ServerVersion = "15.0.2000.5"
		err := validateSQLServerCatalogTable(insp, layout.Columns, layout.OwnerColumn, opts)
		require.ErrorIs(t, err, ErrIncompatibleExistingTable)

		insp = validTestSQLServerTableInspection(layout, opts)
		insp.SQLServer.CompatibilityLevel = 150
		err = validateSQLServerCatalogTable(insp, layout.Columns, layout.OwnerColumn, opts)
		require.ErrorIs(t, err, ErrIncompatibleExistingTable)

		insp = validTestSQLServerTableInspection(layout, opts)
		insp.SQLServer.EngineEdition = 1 // Personal
		err = validateSQLServerCatalogTable(insp, layout.Columns, layout.OwnerColumn, opts)
		require.ErrorIs(t, err, ErrIncompatibleExistingTable)
	})

	t.Run("complete_owned_layout_accepted", func(t *testing.T) {
		group := validTestSQLServerTableInspection(layout, opts)
		receipt := validTestSQLServerReceiptInspection(opts)
		found, err := inspectSQLServerGroupSnapshot(layout, opts.Database, group, receipt)
		require.NoError(t, err)
		require.False(t, found.pendingGroup)
		require.False(t, found.pendingReceipt)
		require.Len(t, found.tables, 2)
		require.Equal(t, SchemaTableActionUnchanged, found.tables[0].Action)
		require.NotEmpty(t, found.tables[0].CatalogDigest)
		require.Equal(t, SchemaTableActionUnchanged, found.tables[1].Action)
		require.NotEmpty(t, found.tables[1].CatalogDigest)

		generated, err := GenerateGroupSchemaDDL(dialectSQLServer, layout)
		require.NoError(t, err)
		stmts, err := sqlserverPendingSchemaStatements(layout, found.tables, generated)
		require.NoError(t, err)
		require.Empty(t, stmts, "complete existing layout produces zero statements")
	})

	t.Run("missing_group_and_receipt_accepted", func(t *testing.T) {
		missingEvidence := &TargetSQLServerTableEvidence{
			Database:           opts.Database,
			Schema:             opts.Schema,
			ServerVersion:      opts.ServerVersion,
			ServerEdition:      opts.ServerEdition,
			CapabilityRevision: opts.CapabilityRevision,
			EngineEdition:      opts.EngineEdition,
			CompatibilityLevel: opts.CompatibilityLevel,
			MetadataVerified:   true,
		}
		group := TargetTableInspection{Status: inspectionMissing, SQLServer: missingEvidence}
		receipt := TargetTableInspection{Status: inspectionMissing, SQLServer: missingEvidence}
		found, err := inspectSQLServerGroupSnapshot(layout, opts.Database, group, receipt)
		require.NoError(t, err)
		require.True(t, found.pendingGroup)
		require.True(t, found.pendingReceipt)
		require.Equal(t, SchemaTableActionCreate, found.tables[0].Action)
		require.Equal(t, SchemaTableActionCreate, found.tables[1].Action)

		generated, err := GenerateGroupSchemaDDL(dialectSQLServer, layout)
		require.NoError(t, err)
		stmts, err := sqlserverPendingSchemaStatements(layout, found.tables, generated)
		require.NoError(t, err)
		require.Len(t, stmts, 3, "creates main table, index, and receipt table")
	})

	t.Run("independently_missing_receipt_accepted", func(t *testing.T) {
		group := validTestSQLServerTableInspection(layout, opts)
		missingEvidence := &TargetSQLServerTableEvidence{
			Database:           opts.Database,
			Schema:             opts.Schema,
			ServerVersion:      opts.ServerVersion,
			ServerEdition:      opts.ServerEdition,
			CapabilityRevision: opts.CapabilityRevision,
			EngineEdition:      opts.EngineEdition,
			CompatibilityLevel: opts.CompatibilityLevel,
			MetadataVerified:   true,
		}
		receipt := TargetTableInspection{Status: inspectionMissing, SQLServer: missingEvidence}
		found, err := inspectSQLServerGroupSnapshot(layout, opts.Database, group, receipt)
		require.NoError(t, err)
		require.False(t, found.pendingGroup)
		require.True(t, found.pendingReceipt)
		require.Equal(t, SchemaTableActionUnchanged, found.tables[0].Action)
		require.Equal(t, SchemaTableActionCreate, found.tables[1].Action)

		generated, err := GenerateGroupSchemaDDL(dialectSQLServer, layout)
		require.NoError(t, err)
		stmts, err := sqlserverPendingSchemaStatements(layout, found.tables, generated)
		require.NoError(t, err)
		require.Len(t, stmts, 1, "only creates receipt table")
	})

	t.Run("receipt_indexes_not_inspected_rejected", func(t *testing.T) {
		group := validTestSQLServerTableInspection(layout, opts)
		receipt := validTestSQLServerReceiptInspection(opts)
		receipt.IndexesInspected = false
		_, err := inspectSQLServerGroupSnapshot(layout, opts.Database, group, receipt)
		require.ErrorIs(t, err, ErrTargetInspectionUnconfirmed)
	})

	t.Run("receipt_non_pk_indexes_rejected", func(t *testing.T) {
		group := validTestSQLServerTableInspection(layout, opts)
		receipt := validTestSQLServerReceiptInspection(opts)
		receipt.Indexes = []TargetIndexInspection{{Name: "idx_extra_receipt", Columns: []string{"effect_key"}}}
		_, err := inspectSQLServerGroupSnapshot(layout, opts.Database, group, receipt)
		require.ErrorIs(t, err, ErrIncompatibleExistingTable)
	})

	t.Run("snapshot_rejects_options_missing_facts", func(t *testing.T) {
		group := validTestSQLServerTableInspection(layout, opts)
		receipt := validTestSQLServerReceiptInspection(opts)
		noFacts := layout
		noFacts.SQLServer = &SQLServerGroupSchemaOptions{Database: opts.Database, Schema: opts.Schema}
		_, err := inspectSQLServerGroupSnapshot(noFacts, opts.Database, group, receipt)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidGroupSchemaLayout), "must reject missing facts in options")
	})

	t.Run("binding_rejects_evidence_missing_facts", func(t *testing.T) {
		group := validTestSQLServerTableInspection(layout, opts)
		group.SQLServer.ServerVersion = ""
		receipt := validTestSQLServerReceiptInspection(opts)
		_, err := inspectSQLServerGroupSnapshot(layout, opts.Database, group, receipt)
		require.ErrorIs(t, err, ErrSchemaTargetChanged)
	})

	t.Run("dashed_and_unicode_db_accepted", func(t *testing.T) {
		for _, dbName := range []string{"recording-prod", "生產記錄庫"} {
			customOpts := opts
			customOpts.Database = dbName
			customLayout, err := ManagedSQLServerLongGroupLayout("gw_long_test", "_gw_owner_0123456789abcdef01234567", customOpts)
			require.NoError(t, err)
			group := validTestSQLServerTableInspection(customLayout, customOpts)
			receipt := validTestSQLServerReceiptInspection(customOpts)
			found, err := inspectSQLServerGroupSnapshot(customLayout, customOpts.Database, group, receipt)
			require.NoError(t, err, "database %q must be accepted", dbName)
			require.False(t, found.pendingGroup)
			require.False(t, found.pendingReceipt)
		}
	})

	t.Run("permission_denied", func(t *testing.T) {
		missingEvidence := &TargetSQLServerTableEvidence{
			Database:         opts.Database,
			Schema:           opts.Schema,
			MetadataVerified: true,
		}
		group := TargetTableInspection{Status: inspectionForbidden, SQLServer: missingEvidence}
		receipt := TargetTableInspection{Status: inspectionMissing, SQLServer: missingEvidence}
		_, err := inspectSQLServerGroupSnapshot(layout, opts.Database, group, receipt)
		require.ErrorIs(t, err, ErrTargetPermissionDenied)
	})
}

func TestSQLServerPreviewScopeValidation(t *testing.T) {
	opts := validTestSQLServerOptions()

	t.Run("valid_scope_and_facts", func(t *testing.T) {
		err := validateSQLServerPreviewScope(opts.Database, opts.Schema, "sqlserver", opts)
		require.NoError(t, err)
	})

	t.Run("valid_dashed_database_scope", func(t *testing.T) {
		dashedOpts := opts
		dashedOpts.Database = "recording-prod"
		err := validateSQLServerPreviewScope(dashedOpts.Database, dashedOpts.Schema, "sqlserver", dashedOpts)
		require.NoError(t, err)
	})

	t.Run("valid_unicode_chinese_database_scope", func(t *testing.T) {
		unicodeOpts := opts
		unicodeOpts.Database = "生產記錄庫"
		err := validateSQLServerPreviewScope(unicodeOpts.Database, unicodeOpts.Schema, "sqlserver", unicodeOpts)
		require.NoError(t, err)
	})

	t.Run("database_with_nul_rejected", func(t *testing.T) {
		nulOpts := opts
		nulOpts.Database = "recording\x00prod"
		err := validateSQLServerPreviewScope(nulOpts.Database, nulOpts.Schema, "sqlserver", nulOpts)
		require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)
	})

	t.Run("schema_still_strictly_validates_identifier", func(t *testing.T) {
		badSchema := opts
		badSchema.Schema = "dbo-invalid"
		err := validateSQLServerPreviewScope(opts.Database, badSchema.Schema, "sqlserver", badSchema)
		require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)
	})

	t.Run("wrong_dialect", func(t *testing.T) {
		err := validateSQLServerPreviewScope(opts.Database, opts.Schema, "postgres", opts)
		require.ErrorIs(t, err, ErrUnsupportedSchemaDialect)
	})

	t.Run("missing_database", func(t *testing.T) {
		err := validateSQLServerPreviewScope("", opts.Schema, "sqlserver", opts)
		require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)
	})

	t.Run("missing_schema", func(t *testing.T) {
		err := validateSQLServerPreviewScope(opts.Database, "", "sqlserver", opts)
		require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)
	})

	t.Run("schema_mismatch", func(t *testing.T) {
		err := validateSQLServerPreviewScope(opts.Database, "other_schema", "sqlserver", opts)
		require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)
	})

	t.Run("database_mismatch", func(t *testing.T) {
		err := validateSQLServerPreviewScope("other_db", opts.Schema, "sqlserver", opts)
		require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)
	})

	t.Run("options_missing_database", func(t *testing.T) {
		noDb := opts
		noDb.Database = ""
		err := validateSQLServerPreviewScope("recording", opts.Schema, "sqlserver", noDb)
		require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)
	})

	t.Run("non_2022_version", func(t *testing.T) {
		badVer := opts
		badVer.ServerVersion = "15.0.2000"
		err := validateSQLServerPreviewScope(opts.Database, opts.Schema, "sqlserver", badVer)
		require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)
	})

	t.Run("missing_server_edition", func(t *testing.T) {
		bad := opts
		bad.ServerEdition = ""
		err := validateSQLServerPreviewScope(opts.Database, opts.Schema, "sqlserver", bad)
		require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)
	})

	t.Run("invalid_engine_edition", func(t *testing.T) {
		bad := opts
		bad.EngineEdition = 5 // SQL Database / Azure
		err := validateSQLServerPreviewScope(opts.Database, opts.Schema, "sqlserver", bad)
		require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)
	})

	t.Run("invalid_compatibility_level", func(t *testing.T) {
		bad := opts
		bad.CompatibilityLevel = 150
		err := validateSQLServerPreviewScope(opts.Database, opts.Schema, "sqlserver", bad)
		require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)
	})
}
