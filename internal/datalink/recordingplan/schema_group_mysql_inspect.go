package recordingplan

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

const (
	mysqlStepKindGroupTable   = "group_table"
	mysqlStepKindIndex        = "index"
	mysqlStepKindReceiptTable = "receipt_table"
	mysqlStepStateMissing     = "missing"
	mysqlStepStateExact       = "exact"
)

// TargetMySQLTableEvidence is the saved destination's actual catalog evidence.
// Missing tables still require family, version, database and capability proof.
type TargetMySQLTableEvidence struct {
	ServerFamily       string
	ServerVersion      string
	CapabilityRevision string
	Database           string
	Engine             string
	Charset            string
	Collation          string
	TableType          string
	PrimaryKey         []string
	ColumnCollations   map[string]string
	ColumnDefaults     map[string]*string
	JSONColumns        map[string]bool
	Partitions         *MySQLPartitionEvidence `json:"partitions,omitempty"`
}

type mysqlApplyStepRef struct {
	Kind    string
	Name    string
	SQLHash string
}

type mysqlApplyStepEvidence struct {
	State       string
	Fingerprint string
}

func inspectMySQLApplyStep(ctx context.Context, token *SchemaPreviewToken, layout GroupSchemaLayout, step mysqlApplyStepRef, inspect TargetInspector) (mysqlApplyStepEvidence, error) {
	table, err := mysqlApplyStepTable(token, layout, step)
	if err != nil {
		return mysqlApplyStepEvidence{}, err
	}
	if inspect == nil {
		return mysqlApplyStepEvidence{}, ErrTargetInspectionUnconfirmed
	}
	state, err := inspect(ctx, table)
	if err != nil {
		return mysqlApplyStepEvidence{}, fmt.Errorf("%w: %w", ErrTargetInspectionUnconfirmed, err)
	}
	return inspectMySQLApplyStepState(token, layout, step, state)
}

func mysqlApplyStepTable(token *SchemaPreviewToken, layout GroupSchemaLayout, step mysqlApplyStepRef) (string, error) {
	table, statement, err := mysqlStepStatement(layout, step)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256([]byte(statement))
	if step.SQLHash != hex.EncodeToString(hash[:]) || token == nil || token.Dialect != dialectMySQL ||
		token.GroupLayout == nil || !groupSchemaLayoutsEqual(layout, *token.GroupLayout) {
		return "", ErrPreviewTokenStale
	}
	return table, nil
}

func inspectMySQLApplyStepState(token *SchemaPreviewToken, layout GroupSchemaLayout, step mysqlApplyStepRef, state TargetTableInspection) (mysqlApplyStepEvidence, error) {
	_, err := mysqlApplyStepTable(token, layout, step)
	if err != nil {
		return mysqlApplyStepEvidence{}, err
	}
	if err := checkMySQLCatalogBinding(token, layout, state.MySQL); err != nil {
		return mysqlApplyStepEvidence{}, err
	}
	switch state.Status {
	case inspectionMissing:
		if state.MySQL.Partitions != nil {
			return mysqlApplyStepEvidence{}, ErrIncompatibleExistingTable
		}
		return mysqlApplyStepEvidence{State: mysqlStepStateMissing}, nil
	case inspectionForbidden:
		return mysqlApplyStepEvidence{}, fmt.Errorf("%w: %w", ErrTargetInspectionUnconfirmed, ErrTargetPermissionDenied)
	case inspectionExists:
	default:
		return mysqlApplyStepEvidence{}, ErrTargetInspectionUnconfirmed
	}
	if err := validateMySQLCatalogPartitions(step.Kind, state.MySQL, *layout.MySQL); err != nil {
		return mysqlApplyStepEvidence{}, err
	}
	columns, owner := layout.Columns, layout.OwnerColumn
	indexes := layout.Indexes
	if step.Kind == mysqlStepKindReceiptTable {
		columns, owner, indexes = mysqlReceiptColumns(), managedEffectReceiptMarker, nil
	}
	if err := validateMySQLCatalogTable(state, columns, owner, *layout.MySQL); err != nil {
		return mysqlApplyStepEvidence{}, err
	}
	missing, err := missingMySQLIndexes(state, indexes)
	if err != nil {
		return mysqlApplyStepEvidence{}, err
	}
	if step.Kind == mysqlStepKindIndex && slices.Contains(missing, step.Name) {
		return mysqlApplyStepEvidence{State: mysqlStepStateMissing}, nil
	}
	if step.Kind == mysqlStepKindGroupTable {
		// Index creation is a later authorized step and cannot change the table
		// step's column/ownership fingerprint.
		state.Indexes, state.IndexesInspected = nil, false
	}
	payload, err := json.Marshal(struct {
		Step  mysqlApplyStepRef
		Table TargetTableInspection
	}{step, state})
	if err != nil {
		return mysqlApplyStepEvidence{}, fmt.Errorf("catalog fingerprint: %w", err)
	}
	digest := sha256.Sum256(payload)
	return mysqlApplyStepEvidence{State: mysqlStepStateExact, Fingerprint: hex.EncodeToString(digest[:])}, nil
}

func mysqlStepStatement(layout GroupSchemaLayout, step mysqlApplyStepRef) (table, statement string, err error) {
	if layout.MySQL == nil {
		return "", "", ErrSchemaAdapterUnverified
	}
	fixed, err := ManagedMySQLLongGroupLayout(layout.TableName, layout.OwnerColumn, *layout.MySQL)
	if err != nil || !groupSchemaLayoutsEqual(fixed, layout) {
		return "", "", ErrInvalidGroupSchemaLayout
	}
	statements, err := GenerateGroupSchemaDDL(dialectMySQL, layout)
	if err != nil {
		return "", "", err
	}
	switch step.Kind {
	case mysqlStepKindGroupTable:
		if step.Name == layout.TableName {
			return layout.TableName, statements[0], nil
		}
	case mysqlStepKindReceiptTable:
		if step.Name == managedEffectReceiptTable {
			return managedEffectReceiptTable, statements[len(statements)-1], nil
		}
	case mysqlStepKindIndex:
		for i, index := range layout.Indexes {
			if step.Name == index.Name {
				return layout.TableName, statements[i+1], nil
			}
		}
	}
	return "", "", ErrPreviewTokenStale
}

func checkMySQLCatalogBinding(token *SchemaPreviewToken, layout GroupSchemaLayout, evidence *TargetMySQLTableEvidence) error {
	if evidence == nil {
		return ErrTargetInspectionUnconfirmed
	}
	options := layout.MySQL
	if token.Database == "" || token.Schema != token.Database || evidence.Database != token.Database ||
		evidence.ServerFamily != options.ServerFamily || evidence.ServerVersion != options.ServerVersion ||
		evidence.CapabilityRevision != options.CapabilityRevision {
		return fmt.Errorf("%w: %w", ErrPreviewTokenStale, ErrSchemaTargetChanged)
	}
	return nil
}

func validateMySQLCatalogPartitions(kind string, evidence *TargetMySQLTableEvidence, options MySQLGroupSchemaOptions) error {
	if evidence == nil {
		return ErrTargetInspectionUnconfirmed
	}
	if kind == mysqlStepKindReceiptTable {
		if evidence.Partitions != nil {
			return ErrIncompatibleExistingTable
		}
		return nil
	}
	if kind != mysqlStepKindGroupTable && kind != mysqlStepKindIndex {
		return ErrPreviewTokenStale
	}
	if !sameMySQLPartitionEvidence(evidence.Partitions, options.NativeMonthly) {
		return ErrIncompatibleExistingTable
	}
	return nil
}

func mysqlReceiptColumns() []GroupSchemaColumn {
	return []GroupSchemaColumn{
		{Name: "effect_key", SQLType: mysqlIDSQLType, PrimaryKey: true},
		{Name: "payload_digest", SQLType: mysqlHashSQLType},
		{Name: "committed_at", SQLType: mysqlTimeSQLType},
		{Name: managedEffectReceiptMarker, SQLType: mysqlOwnerSQLType},
	}
}

func validateMySQLCatalogTable(state TargetTableInspection, columns []GroupSchemaColumn, owner string, options MySQLGroupSchemaOptions) error {
	evidence := state.MySQL
	if evidence == nil || state.ColumnTypes == nil || state.ColumnNullable == nil || state.ColumnPrimaryKey == nil {
		return ErrTargetInspectionUnconfirmed
	}
	if evidence.TableType != "BASE TABLE" || !strings.EqualFold(evidence.Engine, mysqlEngineInnoDB) ||
		evidence.Charset != options.Charset || evidence.Collation != options.Collation || len(state.Columns) != len(columns) {
		return ErrIncompatibleExistingTable
	}
	var primary []string
	for i, column := range columns {
		actualType, hasType := state.ColumnTypes[column.Name]
		nullable, hasNullable := state.ColumnNullable[column.Name]
		primaryKey, hasKey := state.ColumnPrimaryKey[column.Name]
		if state.Columns[i] != column.Name || !hasType || !hasNullable || !hasKey ||
			nullable != column.Nullable || primaryKey != column.PrimaryKey || !mysqlCatalogTypeMatches(column.SQLType, actualType, options.ServerFamily) {
			return ErrIncompatibleExistingTable
		}
		if column.PrimaryKey {
			primary = append(primary, column.Name)
		}
		if mysqlCollatedSQLType(column.SQLType) && evidence.ColumnCollations[column.Name] != options.Collation {
			return ErrIncompatibleExistingTable
		}
		if column.SQLType == mysqlJSONSQLType && !evidence.JSONColumns[column.Name] {
			return ErrIncompatibleExistingTable
		}
	}
	if !slices.Equal(primary, evidence.PrimaryKey) {
		return ErrIncompatibleExistingTable
	}
	ownerDefault := evidence.ColumnDefaults[owner]
	if ownerDefault == nil || *ownerDefault != managedOwnerDefault {
		return ErrIncompatibleExistingTable
	}
	return nil
}

func mysqlCatalogTypeMatches(expected, actual, family string) bool {
	actual = strings.ToUpper(strings.TrimSpace(actual))
	if expected == sqlTypeBoolean {
		return actual == "TINYINT(1)" || actual == sqlTypeBoolean
	}
	if expected == mysqlJSONSQLType && family == mysqlFamilyMariaDB {
		return actual == mysqlLongTextSQLType
	}
	if expected == sqlTypeBigint && family == mysqlFamilyMariaDB {
		// MariaDB reports the signed BIGINT display width; it is not precision
		// or an unsigned range and never permits the UNSIGNED suffix.
		return actual == sqlTypeBigint || actual == "BIGINT(20)"
	}
	return expected == actual
}

func missingMySQLIndexes(state TargetTableInspection, expected []GroupSchemaIndex) ([]string, error) {
	if !state.IndexesInspected {
		return nil, ErrTargetInspectionUnconfirmed
	}
	seen := make(map[string]bool, len(state.Indexes))
	for _, actual := range state.Indexes {
		if seen[actual.Name] {
			return nil, ErrIncompatibleExistingTable
		}
		seen[actual.Name] = true
		i := slices.IndexFunc(expected, func(index GroupSchemaIndex) bool { return index.Name == actual.Name })
		if i < 0 || actual.Unique || !slices.Equal(actual.Columns, expected[i].Columns) {
			return nil, ErrIncompatibleExistingTable
		}
	}
	var missing []string
	for _, index := range expected {
		if !seen[index.Name] {
			missing = append(missing, index.Name)
		}
	}
	return missing, nil
}
