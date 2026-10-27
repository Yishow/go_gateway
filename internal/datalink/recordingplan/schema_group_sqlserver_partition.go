package recordingplan

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

const (
	sqlserverPartitionMethod        = "RANGE RIGHT"
	sqlserverPartitionFunctionType  = "datetime2(6)"
	sqlserverPartitionColumn        = "bucket_start"
	sqlserverPartitionFilegroup     = "PRIMARY"
	sqlserverPartitionFunctionPfx   = "pf_"
	sqlserverPartitionSchemePfx     = "ps_"
	sqlserverPartitionOwnerProperty = "gateway.recording.partition_owner"
	sqlserverPartitionFutureValue   = "MAXVALUE"
	sqlserverPartitionPastName      = "p_past"
	sqlserverPartitionFutureName    = "p_future"
	sqlserverNativeMonthlyCount     = 4
	sqlserverPartitionBoundaryCount = sqlserverNativeMonthlyCount + 2
	sqlServerPartitionTimeLayout    = "2006-01-02T15:04:05.000000"
)

// ErrSQLServerPartitionUnavailable is the stable partition_unavailable class:
// native_monthly was explicitly selected but the destination cannot prove the
// required owned partition objects. It never falls back to plain storage.
var ErrSQLServerPartitionUnavailable = errors.New("partition_unavailable")

// SQLServerPartitionUnavailableError carries the safe reasons behind the
// partition_unavailable class.
type SQLServerPartitionUnavailableError struct {
	Reasons []string
}

func (e *SQLServerPartitionUnavailableError) Error() string { return "partition_unavailable" }

func (e *SQLServerPartitionUnavailableError) Is(target error) bool {
	return target == ErrSQLServerPartitionUnavailable
}

// SQLServerPartitionBoundary is one ordered range partition. LessThan is the
// canonical UTC datetime2(6) boundary; the final catch-all uses MAXVALUE.
type SQLServerPartitionBoundary struct {
	Name     string `json:"name"`
	LessThan string `json:"less_than"`
	Ordinal  int    `json:"ordinal"`
}

// SQLServerPartitionEvidence is the frozen native partition definition of one
// managed SQL Server group table. The receipt table never carries this value.
type SQLServerPartitionEvidence struct {
	Method       string                       `json:"method"`
	Column       string                       `json:"column"`
	FunctionName string                       `json:"function_name"`
	SchemeName   string                       `json:"scheme_name"`
	Filegroup    string                       `json:"filegroup"`
	Partitions   []SQLServerPartitionBoundary `json:"partitions"`
}

// SQLServerPartitionPreview exposes the pending or already-present partition
// projection in a schema preview without an operation ID or connection detail.
type SQLServerPartitionPreview struct {
	Action       string                       `json:"action"`
	FunctionName string                       `json:"function_name"`
	SchemeName   string                       `json:"scheme_name"`
	Filegroup    string                       `json:"filegroup"`
	Partitions   []SQLServerPartitionBoundary `json:"partitions"`
}

// SQLServerNativeMonthlyObjectNames returns the only partition object names a
// scoped table may own. Names are database-scoped in SQL Server, so deriving
// them from the shared scope-derived table name prevents adopting a DBA object.
func SQLServerNativeMonthlyObjectNames(tableName string) (functionName, schemeName string, err error) {
	if !validGroupIdentifier(tableName) {
		return "", "", fmt.Errorf("%w: sqlserver partition table name", ErrInvalidGroupSchemaLayout)
	}
	return sqlserverPartitionFunctionPfx + tableName, sqlserverPartitionSchemePfx + tableName, nil
}

// NewSQLServerNativeMonthlyPartition builds the canonical current UTC month
// plus the next three complete calendar months. RANGE RIGHT yields a past
// catch-all, four finite monthly partitions, and a future catch-all.
func NewSQLServerNativeMonthlyPartition(at time.Time, functionName, schemeName string) (*SQLServerPartitionEvidence, error) {
	if at.IsZero() {
		return nil, fmt.Errorf("%w: sqlserver native monthly anchor is zero", ErrInvalidGroupSchemaLayout)
	}
	if !validGroupIdentifier(functionName) || !validGroupIdentifier(schemeName) {
		return nil, fmt.Errorf("%w: sqlserver partition object name", ErrInvalidGroupSchemaLayout)
	}
	utc := at.UTC()
	if !sqlServerPartitionYearInRange(utc.Year()) {
		return nil, fmt.Errorf("%w: sqlserver native monthly anchor is outside datetime2 range", ErrInvalidGroupSchemaLayout)
	}
	monthStart := time.Date(utc.Year(), utc.Month(), 1, 0, 0, 0, 0, time.UTC)
	partitions := make([]SQLServerPartitionBoundary, 0, sqlserverPartitionBoundaryCount)
	partitions = append(partitions, SQLServerPartitionBoundary{
		Name: sqlserverPartitionPastName, LessThan: monthStart.Format(sqlServerPartitionTimeLayout), Ordinal: 1,
	})
	for offset := range sqlserverNativeMonthlyCount {
		month := monthStart.AddDate(0, offset, 0)
		next := monthStart.AddDate(0, offset+1, 0)
		if !sqlServerPartitionYearInRange(next.Year()) {
			return nil, fmt.Errorf("%w: sqlserver native monthly boundary is outside datetime2 range", ErrInvalidGroupSchemaLayout)
		}
		partitions = append(partitions, SQLServerPartitionBoundary{
			Name:     fmt.Sprintf("p_%04d%02d", month.Year(), month.Month()),
			LessThan: next.Format(sqlServerPartitionTimeLayout),
			Ordinal:  offset + 2,
		})
	}
	partitions = append(partitions, SQLServerPartitionBoundary{
		Name: sqlserverPartitionFutureName, LessThan: sqlserverPartitionFutureValue, Ordinal: sqlserverPartitionBoundaryCount,
	})
	evidence := &SQLServerPartitionEvidence{
		Method: sqlserverPartitionMethod, Column: sqlserverPartitionColumn,
		FunctionName: functionName, SchemeName: schemeName, Filegroup: sqlserverPartitionFilegroup,
		Partitions: partitions,
	}
	if err := validateSQLServerPartitionOwnership(evidence); err != nil {
		return nil, err
	}
	return evidence, nil
}

func sqlServerPartitionYearInRange(year int) bool { return year >= 1 && year <= 9999 }

func validateSQLServerPartitionOwnership(evidence *SQLServerPartitionEvidence) error {
	if evidence == nil {
		return nil
	}
	if evidence.Method != sqlserverPartitionMethod {
		return fmt.Errorf("%w: sqlserver partition method must be RANGE RIGHT", ErrInvalidGroupSchemaLayout)
	}
	if evidence.Column != sqlserverPartitionColumn {
		return fmt.Errorf("%w: sqlserver partition column must be bucket_start", ErrInvalidGroupSchemaLayout)
	}
	if evidence.Filegroup != sqlserverPartitionFilegroup {
		return fmt.Errorf("%w: sqlserver partition filegroup must be PRIMARY", ErrInvalidGroupSchemaLayout)
	}
	if !validGroupIdentifier(evidence.FunctionName) || !validGroupIdentifier(evidence.SchemeName) {
		return fmt.Errorf("%w: sqlserver partition object name", ErrInvalidGroupSchemaLayout)
	}
	if len(evidence.Partitions) != sqlserverPartitionBoundaryCount {
		return fmt.Errorf("%w: sqlserver partition must have past, four monthly, and future ranges", ErrInvalidGroupSchemaLayout)
	}
	past, err := parseSQLServerPartitionBoundary(evidence.Partitions[0].LessThan)
	if err != nil || evidence.Partitions[0].Name != sqlserverPartitionPastName || evidence.Partitions[0].Ordinal != 1 {
		return fmt.Errorf("%w: sqlserver past partition boundary", ErrInvalidGroupSchemaLayout)
	}
	for offset := range sqlserverNativeMonthlyCount {
		month := past.AddDate(0, offset, 0)
		next := past.AddDate(0, offset+1, 0)
		want := SQLServerPartitionBoundary{
			Name:     fmt.Sprintf("p_%04d%02d", month.Year(), month.Month()),
			LessThan: next.Format(sqlServerPartitionTimeLayout), Ordinal: offset + 2,
		}
		if evidence.Partitions[offset+1] != want {
			return fmt.Errorf("%w: sqlserver monthly partition names, boundaries, or ordinals", ErrInvalidGroupSchemaLayout)
		}
	}
	if evidence.Partitions[sqlserverPartitionBoundaryCount-1] != (SQLServerPartitionBoundary{
		Name: sqlserverPartitionFutureName, LessThan: sqlserverPartitionFutureValue, Ordinal: sqlserverPartitionBoundaryCount,
	}) {
		return fmt.Errorf("%w: sqlserver future partition catch-all", ErrInvalidGroupSchemaLayout)
	}
	return nil
}

func validateSQLServerPartitionEvidence(evidence *SQLServerPartitionEvidence, tableName string) error {
	if evidence == nil {
		return nil
	}
	if err := validateSQLServerPartitionOwnership(evidence); err != nil {
		return err
	}
	functionName, schemeName, err := SQLServerNativeMonthlyObjectNames(tableName)
	if err != nil || evidence.FunctionName != functionName || evidence.SchemeName != schemeName {
		return fmt.Errorf("%w: sqlserver partition objects must be scoped to %s", ErrInvalidGroupSchemaLayout, tableName)
	}
	return nil
}

func parseSQLServerPartitionBoundary(value string) (time.Time, error) {
	parsed, err := time.Parse(sqlServerPartitionTimeLayout, value)
	if err != nil || parsed.UTC().Format(sqlServerPartitionTimeLayout) != value || !sqlServerPartitionYearInRange(parsed.Year()) {
		return time.Time{}, fmt.Errorf("not a canonical UTC datetime2(6) month boundary")
	}
	if parsed.Day() != 1 || parsed.Hour() != 0 || parsed.Minute() != 0 || parsed.Second() != 0 || parsed.Nanosecond() != 0 {
		return time.Time{}, fmt.Errorf("not the first instant of a UTC month")
	}
	return parsed.UTC(), nil
}

func cloneSQLServerPartitionEvidence(evidence *SQLServerPartitionEvidence) *SQLServerPartitionEvidence {
	if evidence == nil {
		return nil
	}
	cloned := *evidence
	cloned.Partitions = slices.Clone(evidence.Partitions)
	return &cloned
}

// CloneSQLServerPartitionEvidence returns an independent native partition
// definition for persisted or preview layouts.
func CloneSQLServerPartitionEvidence(evidence *SQLServerPartitionEvidence) *SQLServerPartitionEvidence {
	return cloneSQLServerPartitionEvidence(evidence)
}

func sameSQLServerPartitionEvidence(left, right *SQLServerPartitionEvidence) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return left.Method == right.Method && left.Column == right.Column && left.FunctionName == right.FunctionName &&
		left.SchemeName == right.SchemeName && left.Filegroup == right.Filegroup &&
		slices.Equal(left.Partitions, right.Partitions)
}

func sqlServerPartitionDDL(layout GroupSchemaLayout, evidence *SQLServerPartitionEvidence) ([]string, error) {
	if evidence == nil {
		return nil, nil
	}
	if err := validateSQLServerPartitionEvidence(evidence, layout.TableName); err != nil {
		return nil, err
	}
	values := make([]string, 0, sqlserverNativeMonthlyCount+1)
	for _, partition := range evidence.Partitions[:sqlserverPartitionBoundaryCount-1] {
		boundary, err := parseSQLServerPartitionBoundary(partition.LessThan)
		if err != nil {
			return nil, fmt.Errorf("%w: sqlserver partition boundary: %w", ErrInvalidGroupSchemaLayout, err)
		}
		values = append(values, "N'"+boundary.Format("2006-01-02T15:04:05.000000")+"'")
	}
	function := quoteSQLServerIdentifier(evidence.FunctionName)
	scheme := quoteSQLServerIdentifier(evidence.SchemeName)
	owner := quoteSQLServerLiteral(layout.SQLServer.Schema + "." + layout.TableName)
	return []string{
		"CREATE PARTITION FUNCTION " + function + " (" + sqlserverPartitionFunctionType + ") AS " +
			sqlserverPartitionMethod + " FOR VALUES (" + strings.Join(values, ", ") + ");",
		"EXEC sys.sp_addextendedproperty @name=N'" + sqlserverPartitionOwnerProperty + "', @value=N'" + owner +
			"', @level0type=N'PARTITION FUNCTION', @level0name=N'" + evidence.FunctionName + "';",
		"CREATE PARTITION SCHEME " + scheme + " AS PARTITION " + function + " ALL TO ([" + evidence.Filegroup + "]);",
		"EXEC sys.sp_addextendedproperty @name=N'" + sqlserverPartitionOwnerProperty + "', @value=N'" + owner +
			"', @level0type=N'PARTITION SCHEME', @level0name=N'" + evidence.SchemeName + "';",
	}, nil
}

func quoteSQLServerLiteral(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func validateSQLServerNativeMonthlyGroupKey(columns []GroupSchemaColumn) error {
	for _, column := range columns {
		if column.Name == sqlserverPartitionColumn && column.SQLType == sqlserverTimeSQLType && column.PrimaryKey && !column.Nullable {
			return nil
		}
	}
	return fmt.Errorf("%w: sqlserver native partition requires a non-null datetime2(6) bucket_start in the primary key", ErrInvalidGroupSchemaLayout)
}
