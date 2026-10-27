package dbtarget

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

const (
	sqlServerManagedReceiptOwnerColumn     = "_gw_effect_receipts_v1"
	sqlServerManagedReceiptEffectKeyColumn = "effect_key"
	sqlServerManagedOwnerMarkerValue       = "gateway"
	sqlServerManagedUserTableType          = "USER_TABLE"
	sqlServerManagedCollation              = "Latin1_General_100_BIN2"
)

func validateSQLServerManagedRowProof(req GroupInsertRequest) error {
	unsupported := func(reason string) error {
		return fmt.Errorf("%w: %s", ErrGroupInsertUnsupported, reason)
	}
	proof := req.SQLServer
	if proof == nil || proof.Group == nil || proof.Receipt == nil {
		return unsupported("sqlserver managed row proof is required")
	}
	if strings.TrimSpace(proof.ConnectorID) == "" || strings.TrimSpace(proof.ConnectorRevision) == "" {
		return unsupported("sqlserver connector binding is required")
	}
	if req.Kind != schema.DatabaseConnectorKindSQLServer || req.Strategy != GroupEffectReceipt {
		return unsupported("sqlserver managed insert requires the receipt strategy")
	}
	if strings.TrimSpace(req.SchemaName) == "" || strings.TrimSpace(req.TableName) == "" || req.PayloadDigest == "" ||
		!sqlServerStableHashID(req.PayloadDigest) || !sqlServerStableHashID(req.Row.RecordID) ||
		!sqlServerStableHashID(req.Row.EffectKey) {
		return unsupported("sqlserver managed row identity is incomplete")
	}
	if _, err := parseSQLServerCommittedAt(req.CommittedAt); err != nil {
		return unsupported("committed_at must represent UTC microseconds")
	}

	group, receipt := proof.Group, proof.Receipt
	if group.Inspection == nil || group.Evidence == nil || receipt.Inspection == nil || receipt.Evidence == nil {
		return unsupported("sqlserver managed catalog proof is incomplete")
	}
	if !sqlServerManagedEvidenceBindingMatches(group.Evidence, receipt.Evidence) ||
		group.Inspection.Schema != req.SchemaName || receipt.Inspection.Schema != req.SchemaName ||
		group.Evidence.Schema != req.SchemaName || receipt.Evidence.Schema != req.SchemaName ||
		group.Inspection.Table != req.TableName || group.Inspection.Table == EffectReceiptTable ||
		receipt.Inspection.Table != EffectReceiptTable || group.Evidence.ObjectID == receipt.Evidence.ObjectID {
		return unsupported("sqlserver schema, table, or evidence bindings differ")
	}
	if !sqlServerManagedInspectionReady(group) || !sqlServerManagedInspectionReady(receipt) {
		return unsupported("sqlserver group and receipt catalog inspections must be complete")
	}
	groupOwner, groupOK := sqlServerManagedGroupShapeMatches(group)
	if !groupOK || !sqlServerManagedReceiptShapeMatches(receipt) {
		return unsupported("sqlserver group or receipt catalog shape is not canonical")
	}
	if err := validateSQLServerManagedRowCells(req, groupOwner); err != nil {
		return unsupported("sqlserver row cells do not match the frozen managed layout")
	}
	return nil
}

func sqlServerManagedEvidenceBindingMatches(group, receipt *SQLServerManagedTableEvidence) bool {
	return group.MetadataVerified && receipt.MetadataVerified &&
		strings.TrimSpace(group.Database) != "" && group.Database == receipt.Database &&
		strings.TrimSpace(group.Schema) != "" && group.Schema == receipt.Schema &&
		strings.TrimSpace(group.ServerVersion) != "" && group.ServerVersion == receipt.ServerVersion &&
		strings.TrimSpace(group.ServerEdition) != "" && group.ServerEdition == receipt.ServerEdition &&
		group.EngineEdition > 0 && group.EngineEdition == receipt.EngineEdition &&
		group.CompatibilityLevel > 0 && group.CompatibilityLevel == receipt.CompatibilityLevel &&
		group.CapabilityRevision != "" && group.CapabilityRevision == receipt.CapabilityRevision
}

func sqlServerManagedInspectionReady(inspection *SQLServerManagedInspection) bool {
	if inspection == nil || inspection.Inspection == nil || inspection.Evidence == nil {
		return false
	}
	physical, evidence := inspection.Inspection, inspection.Evidence
	return physical.Status == TableInspectionExists && physical.Schema != "" && physical.Table != "" &&
		len(physical.Columns) > 0 && inspection.Indexes != nil && evidence.MetadataVerified &&
		evidence.TableType == sqlServerManagedUserTableType && evidence.ObjectID > 0 &&
		strings.TrimSpace(evidence.Database) != "" && strings.TrimSpace(evidence.Schema) != "" &&
		strings.TrimSpace(evidence.ServerVersion) != "" && strings.TrimSpace(evidence.ServerEdition) != "" &&
		strings.TrimSpace(evidence.CapabilityRevision) != "" &&
		evidence.ColumnCollations != nil && evidence.ColumnDefaults != nil
}

func sqlServerManagedGroupShapeMatches(inspection *SQLServerManagedInspection) (string, bool) {
	columns := inspection.Inspection.Columns
	if len(columns) != len(longColumnRoles)+1 {
		return "", false
	}
	for index, role := range longColumnRoles {
		column := columns[index]
		primaryKey := role.name == longColBucketStart || role.name == longColRecordID
		if column.Name != role.name || column.DataType != sqlServerManagedLongType(role.name, role.role) ||
			column.Nullable != !role.required || column.HasDefault || column.PrimaryKey != primaryKey || column.Unique != primaryKey {
			return "", false
		}
	}
	owner := columns[len(longColumnRoles)]
	if !longOwnerPattern.MatchString(owner.Name) || owner.DataType != "NVARCHAR(16)" || owner.Nullable ||
		owner.PrimaryKey || owner.Unique || !owner.HasDefault {
		return "", false
	}
	nativeMonthly := inspection.Evidence.NativeMonthly != nil
	if !sqlServerManagedEvidenceShapeMatches(inspection, []string{longColBucketStart, longColRecordID}, owner.Name, true, nativeMonthly) {
		return "", false
	}
	wantIndexes := []IndexInfo{{
		Name:    "idx_" + inspection.Inspection.Table + "_src_time",
		Columns: []string{longColWorkspaceID, longColTagID, longColBucketStart},
	}}
	if !sqlServerManagedIndexesEqual(inspection.Indexes, wantIndexes) {
		return "", false
	}
	return owner.Name, true
}

func sqlServerManagedReceiptShapeMatches(inspection *SQLServerManagedInspection) bool {
	columns := inspection.Inspection.Columns
	if inspection.Evidence.RepairOperationID != "" {
		return false
	}
	if len(columns) != 4 {
		return false
	}
	want := []ColumnInfo{
		{Name: sqlServerManagedReceiptEffectKeyColumn, DataType: sqlServerHashSQLType, PrimaryKey: true, Unique: true},
		{Name: "payload_digest", DataType: sqlServerHashSQLType},
		{Name: "committed_at", DataType: sqlServerDateTime2SQLType},
		{Name: sqlServerManagedReceiptOwnerColumn, DataType: "NVARCHAR(16)", HasDefault: true},
	}
	for index := range want {
		if columns[index] != want[index] {
			return false
		}
	}
	if !sqlServerManagedEvidenceShapeMatches(inspection, []string{sqlServerManagedReceiptEffectKeyColumn}, sqlServerManagedReceiptOwnerColumn, false, false) {
		return false
	}
	return len(inspection.Indexes) == 0
}

func sqlServerManagedEvidenceShapeMatches(
	inspection *SQLServerManagedInspection,
	primaryKey []string,
	owner string,
	provenanceJSONChecked bool,
	nativeMonthly bool,
) bool {
	evidence := inspection.Evidence
	if evidence.TableType != sqlServerManagedUserTableType || evidence.ObjectID <= 0 ||
		evidence.ProvenanceJSONChecked != provenanceJSONChecked || !slices.Equal(evidence.PrimaryKey, primaryKey) {
		return false
	}
	if nativeMonthly {
		if evidence.Plain || !evidence.PartitionInspected || !sqlServerManagedNativeMonthlyEvidenceMatches(inspection) {
			return false
		}
	} else if !evidence.Plain || evidence.NativeMonthly != nil {
		return false
	}
	if len(evidence.ColumnDefaults) != len(inspection.Inspection.Columns) ||
		len(evidence.ColumnCollations) == 0 {
		return false
	}
	wantCollations := make(map[string]string)
	for _, column := range inspection.Inspection.Columns {
		if strings.HasPrefix(column.DataType, "NVARCHAR(") || strings.HasPrefix(column.DataType, "CHAR(") {
			wantCollations[column.Name] = sqlServerManagedCollation
		}
		value, exists := evidence.ColumnDefaults[column.Name]
		if !exists {
			return false
		}
		if column.Name == owner {
			if value == nil || *value != sqlServerManagedOwnerMarkerValue {
				return false
			}
		} else if value != nil {
			return false
		}
	}
	return maps.Equal(evidence.ColumnCollations, wantCollations)
}

func sqlServerManagedNativeMonthlyEvidenceMatches(inspection *SQLServerManagedInspection) bool {
	if inspection == nil || inspection.Inspection == nil || inspection.Evidence == nil {
		return false
	}
	return sqlServerManagedNativeMonthlyEvidenceMatchesForTable(inspection.Inspection.Table, inspection.Evidence.NativeMonthly)
}

func sqlServerManagedLongType(name string, role longColumnRole) string {
	switch role {
	case longTime:
		return sqlServerDateTime2SQLType
	case longText:
		if name == longColRecordID {
			return sqlServerHashSQLType
		}
		return "NVARCHAR(191)"
	case longProvenanceRole, longSlotText:
		return sqlServerNvarcharMaxSQLType
	case longSlotInt:
		return sqlServerBigintSQLType
	case longSlotDecimal, longSlotDelta:
		return "DECIMAL(38,18)"
	case longSlotFloat:
		return "FLOAT(53)"
	case longSlotBool:
		return sqlServerBitSQLType
	default:
		return ""
	}
}

func validateSQLServerManagedRowCells(req GroupInsertRequest, owner string) error {
	if !sqlServerUTCMicrosecond(req.Row.BucketStart) || len(req.Row.Cells) != len(longColumnRoles)+1 {
		return fmt.Errorf("invalid SQL Server row shape")
	}
	roles := make(map[string]struct {
		role     longColumnRole
		required bool
	}, len(longColumnRoles))
	for _, column := range longColumnRoles {
		roles[column.name] = struct {
			role     longColumnRole
			required bool
		}{column.role, column.required}
	}
	cells := make(map[string]any, len(req.Row.Cells))
	for _, cell := range req.Row.Cells {
		if cell.Column == "" {
			return fmt.Errorf("empty SQL Server row column")
		}
		if _, duplicate := cells[cell.Column]; duplicate {
			return fmt.Errorf("duplicate SQL Server row column")
		}
		cells[cell.Column] = cell.Value
	}
	for name, role := range roles {
		value, exists := cells[name]
		if !exists || (value == nil && role.required) {
			return fmt.Errorf("incomplete SQL Server row column")
		}
		column := ColumnInfo{Name: name, DataType: sqlServerManagedLongType(name, role.role)}
		if !sqlServerManagedRowCellValue(name, role.role, value, column) {
			return fmt.Errorf("invalid SQL Server row value")
		}
	}
	ownerValue, exists := cells[owner]
	if !exists || ownerValue != sqlServerManagedOwnerMarkerValue {
		return fmt.Errorf("invalid SQL Server row owner")
	}
	for name := range cells {
		if _, known := roles[name]; !known && name != owner {
			return fmt.Errorf("unknown SQL Server row column")
		}
	}
	if cells[longColRecordID] != req.Row.RecordID {
		return fmt.Errorf("SQL Server row record ID differs")
	}
	bucketStart, ok := cells[longColBucketStart].(time.Time)
	if !ok || !bucketStart.Equal(req.Row.BucketStart) {
		return fmt.Errorf("SQL Server row bucket start differs")
	}
	bucketEnd, ok := cells[longColBucketEnd].(time.Time)
	if !ok || !bucketEnd.After(bucketStart) {
		return fmt.Errorf("invalid SQL Server row bucket bounds")
	}
	if !sqlServerManagedProvenanceMatches(cells) {
		return fmt.Errorf("invalid SQL Server row provenance")
	}
	valueSlotCount := 0
	for _, name := range []string{longColValueInt, longColValueDecimal, longColValueFloat, longColValueBool, longColValueText} {
		if cells[name] != nil {
			valueSlotCount++
		}
	}
	if valueSlotCount > 1 {
		return fmt.Errorf("multiple SQL Server value slots are populated")
	}
	valueType, typeOK := cells[longColValueType].(string)
	if !typeOK || !sqlServerManagedValueSlotMatchesType(valueType, cells) {
		return fmt.Errorf("SQL Server value slot does not match the exact type")
	}
	return nil
}

func sqlServerManagedValueSlotMatchesType(valueType string, cells map[string]any) bool {
	var expected string
	switch measurement.ExactType(valueType) {
	case measurement.ExactInt64:
		expected = longColValueInt
	case measurement.ExactUint64, measurement.ExactDecimal:
		expected = longColValueDecimal
	case measurement.ExactFloat64:
		expected = longColValueFloat
	case measurement.ExactBool:
		expected = longColValueBool
	case measurement.ExactText:
		expected = longColValueText
	default:
		return false
	}
	for _, name := range []string{longColValueInt, longColValueDecimal, longColValueFloat, longColValueBool, longColValueText} {
		if name != expected && cells[name] != nil {
			return false
		}
	}
	return true
}

func sqlServerManagedRowCellValue(name string, role longColumnRole, value any, column ColumnInfo) bool {
	if value == nil {
		return true
	}
	switch role {
	case longTime:
		at, ok := value.(time.Time)
		return ok && sqlServerUTCMicrosecond(at)
	case longText, longSlotText, longProvenanceRole:
		text, ok := value.(string)
		if !ok || (text == "" && role != longSlotText) || !utf8.ValidString(text) {
			return false
		}
		if sqlServerManagedIdentityCell(name) && strings.HasSuffix(text, " ") {
			return false
		}
		if name == longColProvenance && !json.Valid([]byte(text)) {
			return false
		}
		encoded, err := EncodeExactValue(SQLDialectSQLServer, column.DataType, measurement.NewText(text))
		return err == nil && encoded == text
	case longSlotInt:
		integer, ok := value.(int64)
		if !ok {
			return false
		}
		encoded, err := EncodeExactValue(SQLDialectSQLServer, column.DataType, measurement.NewInt64(integer))
		return err == nil && encoded == integer
	case longSlotDecimal, longSlotDelta:
		digits, ok := value.(string)
		if !ok {
			return false
		}
		decimal, err := measurement.NewDecimal(digits)
		if err != nil {
			return false
		}
		encoded, err := EncodeExactValue(SQLDialectSQLServer, column.DataType, decimal)
		return err == nil && encoded == digits
	case longSlotFloat:
		floating, ok := value.(float64)
		if !ok {
			return false
		}
		exact, err := measurement.NewFloat64(floating)
		if err != nil {
			return false
		}
		encoded, err := EncodeExactValue(SQLDialectSQLServer, column.DataType, exact)
		return err == nil && encoded == floating
	case longSlotBool:
		boolean, ok := value.(bool)
		if !ok {
			return false
		}
		encoded, err := EncodeExactValue(SQLDialectSQLServer, column.DataType, measurement.NewBool(boolean))
		return err == nil && encoded == boolean
	default:
		return false
	}
}

func sqlServerManagedIdentityCell(name string) bool {
	switch name {
	case longColWorkspaceID, longColGroupID, longColGroupRev, longColCircuitID,
		longColDeviceID, longColPointID, longColTagID, longColMeasurement,
		longColSeriesEpoch, longColProfileKind, longColUnit:
		return true
	default:
		return false
	}
}

func sqlServerManagedProvenanceMatches(cells map[string]any) bool {
	text, ok := cells[longColProvenance].(string)
	if !ok {
		return false
	}
	decoder := json.NewDecoder(strings.NewReader(text))
	decoder.DisallowUnknownFields()
	var provenance longProvenance
	if decoder.Decode(&provenance) != nil || decoder.Decode(new(any)) != io.EOF {
		return false
	}
	valueType, typeOK := cells[longColValueType].(string)
	profileKind, profileOK := cells[longColProfileKind].(string)
	return typeOK && profileOK && provenance.Type == valueType && provenance.Member != "" &&
		provenance.Profile.Kind == profileKind
}

func sqlServerManagedIndexesEqual(left, right []IndexInfo) bool {
	if len(left) != len(right) {
		return false
	}
	return slices.EqualFunc(left, right, func(left, right IndexInfo) bool {
		return left.Name == right.Name && left.Unique == right.Unique && slices.Equal(left.Columns, right.Columns)
	})
}

func sqlServerUTCMicrosecond(value time.Time) bool {
	return value.Location() == time.UTC && value.Nanosecond()%1000 == 0
}

func parseSQLServerCommittedAt(value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil || parsed.Nanosecond()%1000 != 0 {
		return time.Time{}, fmt.Errorf("invalid committed_at")
	}
	return parsed.UTC(), nil
}

func sqlServerManagedPhysicalMismatch() error {
	return fmt.Errorf("%w: %w", ErrSQLServerManagedPhysicalMismatch, ErrGroupInsertUnsupported)
}
