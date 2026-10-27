package recordingplan

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"slices"
)

// VerifyManagedMySQLCatalog checks a previously read group/receipt snapshot
// against the confirmed canonical layout and returns the same catalog digest
// schema confirmation records. It performs no database calls.
func VerifyManagedMySQLCatalog(layout GroupSchemaLayout, database string, group, receipt TargetTableInspection) (string, error) {
	found, err := inspectMySQLGroupSnapshot(layout, database, group, receipt)
	if err != nil {
		return "", err
	}
	if found.pendingGroup || found.pendingReceipt {
		return "", ErrTargetInspectionUnconfirmed
	}
	return schemaTablesDigest(found.tables), nil
}

func mysqlCatalogStepReference(kind, name, statement string) mysqlApplyStepRef {
	sum := sha256.Sum256([]byte(statement))
	return mysqlApplyStepRef{Kind: kind, Name: name, SQLHash: hex.EncodeToString(sum[:])}
}

// inspectMySQLGroupTables treats a table, its indexes and the shared receipt as
// separate catalog facts. It cannot turn a missing index into a missing table.
func inspectMySQLGroupTables(ctx context.Context, layout GroupSchemaLayout, database string, inspect TargetInspector) (groupInspection, error) {
	var found groupInspection
	if inspect == nil || database == "" {
		return found, ErrTargetInspectionUnconfirmed
	}
	if _, err := GenerateGroupSchemaDDL(dialectMySQL, layout); err != nil {
		return found, err
	}
	group, err := inspect(ctx, layout.TableName)
	if err != nil {
		return found, err
	}
	receipt, err := inspect(ctx, managedEffectReceiptTable)
	if err != nil {
		return found, err
	}
	return inspectMySQLGroupSnapshot(layout, database, group, receipt)
}

func inspectMySQLGroupSnapshot(layout GroupSchemaLayout, database string, group, receipt TargetTableInspection) (groupInspection, error) {
	var found groupInspection
	if database == "" {
		return found, ErrTargetInspectionUnconfirmed
	}
	statements, err := GenerateGroupSchemaDDL(dialectMySQL, layout)
	if err != nil {
		return found, err
	}
	token := &SchemaPreviewToken{Dialect: dialectMySQL, Database: database, Schema: database, GroupLayout: &layout}
	tableStep := mysqlCatalogStepReference(mysqlStepKindGroupTable, layout.TableName, statements[0])
	tableEvidence, err := inspectMySQLApplyStepState(token, layout, tableStep, group)
	if err != nil {
		return found, err
	}
	receiptStep := mysqlCatalogStepReference(mysqlStepKindReceiptTable, managedEffectReceiptTable, statements[len(statements)-1])
	receiptEvidence, err := inspectMySQLApplyStepState(token, layout, receiptStep, receipt)
	if err != nil {
		return found, err
	}
	groupTable := SchemaPreviewTable{Name: layout.TableName, Action: SchemaTableActionUnchanged}
	for i, index := range layout.Indexes {
		step := mysqlCatalogStepReference(mysqlStepKindIndex, index.Name, statements[i+1])
		actual, err := inspectMySQLApplyStepState(token, layout, step, group)
		if err != nil {
			return found, err
		}
		if actual.State == mysqlStepStateMissing {
			groupTable.Indexes = append(groupTable.Indexes, index)
		}
	}
	if tableEvidence.State == mysqlStepStateMissing {
		groupTable.Action, groupTable.Columns = SchemaTableActionCreate, groupColumnNames(layout.Columns)
	} else {
		groupTable.CatalogDigest, err = mysqlCatalogDigest(group)
		if err != nil {
			return found, err
		}
		if len(groupTable.Indexes) > 0 {
			groupTable.Action = SchemaTableActionCreateIndexes
		}
	}
	receiptTable := SchemaPreviewTable{Name: managedEffectReceiptTable, Action: SchemaTableActionUnchanged}
	if receiptEvidence.State == mysqlStepStateMissing {
		receiptTable.Action, receiptTable.Columns = SchemaTableActionCreate, groupColumnNames(mysqlReceiptColumns())
		found.pendingReceipt = true
	} else {
		receiptTable.CatalogDigest, err = mysqlCatalogDigest(receipt)
		if err != nil {
			return found, err
		}
	}
	found.pendingGroup = groupTable.Action != SchemaTableActionUnchanged
	found.tables = []SchemaPreviewTable{groupTable, receiptTable}
	return found, nil
}

func mysqlCatalogDigest(state TargetTableInspection) (string, error) {
	payload, err := json.Marshal(state)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

func mysqlPendingSchemaStatements(layout GroupSchemaLayout, tables []SchemaPreviewTable, generated []string) ([]string, error) {
	var statements []string
	for _, table := range tables {
		switch {
		case table.Name == layout.TableName && table.Action == SchemaTableActionCreate:
			statements = append(statements, generated[:len(generated)-1]...)
		case table.Name == layout.TableName && table.Action == SchemaTableActionCreateIndexes:
			for _, pending := range table.Indexes {
				i := slices.IndexFunc(layout.Indexes, func(index GroupSchemaIndex) bool {
					return index.Name == pending.Name && slices.Equal(index.Columns, pending.Columns)
				})
				if i < 0 {
					return nil, ErrInvalidGroupSchemaLayout
				}
				statements = append(statements, generated[i+1])
			}
		case table.Name == managedEffectReceiptTable && table.Action == SchemaTableActionCreate:
			statements = append(statements, generated[len(generated)-1])
		case table.Action == SchemaTableActionUnchanged:
		default:
			return nil, ErrInvalidGroupSchemaLayout
		}
	}
	return statements, nil
}
