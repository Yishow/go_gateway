package recordingplan

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

const (
	mysqlPartitionMethodRangeColumns = "RANGE COLUMNS"
	mysqlPartitionExpression         = "bucket_start"
	mysqlPartitionMaxValue           = "MAXVALUE"
	mysqlNativeMonthlyPartitionCount = 4
	mysqlPartitionBoundaryCount      = mysqlNativeMonthlyPartitionCount + 2
)

// MySQLPartitionEvidence is the frozen native partition definition confirmed
// for one managed MySQL group table.
type MySQLPartitionEvidence struct {
	Method     string                   `json:"method"`
	Expression string                   `json:"expression"`
	Partitions []MySQLPartitionBoundary `json:"partitions"`
}

// MySQLPartitionBoundary is one ordered RANGE COLUMNS partition boundary.
type MySQLPartitionBoundary struct {
	Name     string `json:"name"`
	LessThan string `json:"less_than"`
	Ordinal  int    `json:"ordinal"`
}

// NewMySQLNativeMonthlyPartition creates the fixed past/current/future layout
// anchored to the current UTC month. Its four finite monthly partitions cover
// that month and the following three; the final MAXVALUE partition catches any
// later data until a separately confirmed maintenance change is implemented.
func NewMySQLNativeMonthlyPartition(at time.Time) (*MySQLPartitionEvidence, error) {
	if at.IsZero() {
		return nil, fmt.Errorf("%w: mysql native monthly partition time is zero; provide a UTC timestamp", ErrInvalidGroupSchemaLayout)
	}
	utc := at.UTC()
	if !mysqlDateTimeYearInRange(utc.Year()) {
		return nil, mysqlPartitionDateTimeRangeError("anchor", utc)
	}
	monthStart := time.Date(utc.Year(), utc.Month(), 1, 0, 0, 0, 0, time.UTC)
	boundaries := make([]MySQLPartitionBoundary, 0, mysqlPartitionBoundaryCount)
	boundaries = append(boundaries, MySQLPartitionBoundary{
		Name: "p_past", LessThan: monthStart.Format(time.RFC3339), Ordinal: 1,
	})
	for offset := range mysqlNativeMonthlyPartitionCount {
		month := monthStart.AddDate(0, offset, 0)
		boundary := monthStart.AddDate(0, offset+1, 0)
		if !mysqlDateTimeYearInRange(boundary.Year()) {
			return nil, mysqlPartitionDateTimeRangeError("monthly", boundary)
		}
		boundaries = append(boundaries, MySQLPartitionBoundary{
			Name:     fmt.Sprintf("p_%04d%02d", month.Year(), month.Month()),
			LessThan: boundary.Format(time.RFC3339), Ordinal: offset + 2,
		})
	}
	boundaries = append(boundaries, MySQLPartitionBoundary{
		Name: "p_future", LessThan: mysqlPartitionMaxValue, Ordinal: mysqlPartitionBoundaryCount,
	})
	evidence := &MySQLPartitionEvidence{
		Method: mysqlPartitionMethodRangeColumns, Expression: mysqlPartitionExpression, Partitions: boundaries,
	}
	return normalizeMySQLPartitionEvidence(evidence)
}

func mysqlDateTimeYearInRange(year int) bool { return year >= 1000 && year <= 9999 }

func mysqlPartitionDateTimeRangeError(boundary string, at time.Time) error {
	return fmt.Errorf("%w: mysql native monthly partition %s boundary %s is outside MySQL DATETIME range (1000-01-01 through 9999-12-31)",
		ErrInvalidGroupSchemaLayout, boundary, at.UTC().Format("2006-01-02"))
}

func cloneMySQLPartitionEvidence(evidence *MySQLPartitionEvidence) *MySQLPartitionEvidence {
	if evidence == nil {
		return nil
	}
	cloned := *evidence
	cloned.Partitions = slices.Clone(evidence.Partitions)
	return &cloned
}

func normalizeMySQLPartitionEvidence(evidence *MySQLPartitionEvidence) (*MySQLPartitionEvidence, error) {
	cloned := cloneMySQLPartitionEvidence(evidence)
	if cloned == nil {
		return nil, nil
	}
	if err := validateMySQLPartitionEvidence(cloned); err != nil {
		return nil, err
	}
	return cloned, nil
}

func validateMySQLPartitionEvidence(evidence *MySQLPartitionEvidence) error {
	invalid := func(reason string) error {
		return fmt.Errorf("%w: mysql native monthly partition %s", ErrInvalidGroupSchemaLayout, reason)
	}
	if evidence.Method != mysqlPartitionMethodRangeColumns {
		return invalid("method must be RANGE COLUMNS")
	}
	if evidence.Expression != mysqlPartitionExpression {
		return invalid("expression must be bucket_start")
	}
	if len(evidence.Partitions) != mysqlPartitionBoundaryCount {
		return invalid("must contain p_past, four monthly partitions, and p_future")
	}
	past, err := parseMySQLMonthBoundary(evidence.Partitions[0].LessThan)
	if err != nil {
		return invalid("past boundary must be an RFC3339 UTC month start")
	}
	if past.Day() != 1 || past.Hour() != 0 || past.Minute() != 0 || past.Second() != 0 || past.Nanosecond() != 0 {
		return invalid("past boundary must be the first instant of a UTC month")
	}
	if evidence.Partitions[0] != (MySQLPartitionBoundary{Name: "p_past", LessThan: past.Format(time.RFC3339), Ordinal: 1}) {
		return invalid("past partition name, boundary, or ordinal is invalid")
	}
	for offset := range mysqlNativeMonthlyPartitionCount {
		month := past.AddDate(0, offset, 0)
		boundary := past.AddDate(0, offset+1, 0)
		if !mysqlDateTimeYearInRange(month.Year()) || !mysqlDateTimeYearInRange(boundary.Year()) {
			return mysqlPartitionDateTimeRangeError("monthly", boundary)
		}
		expected := MySQLPartitionBoundary{
			Name:     fmt.Sprintf("p_%04d%02d", month.Year(), month.Month()),
			LessThan: boundary.Format(time.RFC3339), Ordinal: offset + 2,
		}
		if evidence.Partitions[offset+1] != expected {
			return invalid("monthly partition names, ordinals, and contiguous UTC boundaries are required")
		}
	}
	future := MySQLPartitionBoundary{Name: "p_future", LessThan: mysqlPartitionMaxValue, Ordinal: mysqlPartitionBoundaryCount}
	if evidence.Partitions[mysqlPartitionBoundaryCount-1] != future {
		return invalid("future partition must be the final MAXVALUE catch-all")
	}
	return nil
}

func parseMySQLMonthBoundary(value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil || parsed.UTC().Format(time.RFC3339) != value {
		return time.Time{}, fmt.Errorf("not a canonical UTC RFC3339 timestamp")
	}
	if !mysqlDateTimeYearInRange(parsed.Year()) {
		return time.Time{}, fmt.Errorf("outside MySQL DATETIME range")
	}
	return parsed.UTC(), nil
}

func validateMySQLNativeMonthlyGroupColumns(columns []GroupSchemaColumn) error {
	for _, column := range columns {
		if column.Name != mysqlPartitionExpression {
			continue
		}
		if column.SQLType != mysqlTimeSQLType || column.Nullable || !column.PrimaryKey {
			return fmt.Errorf("%w: mysql native monthly partition requires non-null DATETIME(6) bucket_start in the primary key", ErrInvalidGroupSchemaLayout)
		}
		return nil
	}
	return fmt.Errorf("%w: mysql native monthly partition requires bucket_start", ErrInvalidGroupSchemaLayout)
}

func mysqlNativeMonthlyPartitionDDL(evidence *MySQLPartitionEvidence) (string, error) {
	partitions := make([]string, 0, len(evidence.Partitions))
	for _, partition := range evidence.Partitions {
		boundary := partition.LessThan
		if boundary != mysqlPartitionMaxValue {
			parsed, err := parseMySQLMonthBoundary(boundary)
			if err != nil {
				return "", fmt.Errorf("%w: mysql native monthly partition boundary: %w", ErrInvalidGroupSchemaLayout, err)
			}
			boundary = "'" + parsed.Format("2006-01-02 15:04:05") + "'"
		}
		partitions = append(partitions, "PARTITION "+quoteMySQLIdentifier(partition.Name)+
			" VALUES LESS THAN ("+boundary+")")
	}
	return " PARTITION BY " + mysqlPartitionMethodRangeColumns + " (" +
		quoteMySQLIdentifier(mysqlPartitionExpression) + ") (" + strings.Join(partitions, ", ") + ")", nil
}

func sameMySQLPartitionEvidence(left, right *MySQLPartitionEvidence) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return left.Method == right.Method && left.Expression == right.Expression && slices.Equal(left.Partitions, right.Partitions)
}
