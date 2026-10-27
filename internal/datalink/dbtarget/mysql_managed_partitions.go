package dbtarget

import (
	"context"
	"database/sql"
	"slices"
	"strings"
	"time"
)

const (
	mysqlPartitionMethodRangeColumns = "RANGE COLUMNS"
	mysqlPartitionExpression         = "bucket_start"
	mysqlPartitionMaxValue           = "MAXVALUE"
	mysqlPartitionEngineInnoDB       = "InnoDB"
)

// MySQLPartitionBoundary is one RANGE COLUMNS boundary in catalog order.
type MySQLPartitionBoundary struct {
	Name     string `json:"name"`
	LessThan string `json:"less_than"`
	Ordinal  int    `json:"ordinal"`
}

// MySQLPartitionEvidence is the frozen native monthly partition proof.
type MySQLPartitionEvidence struct {
	Method     string                   `json:"method"`
	Expression string                   `json:"expression"`
	Partitions []MySQLPartitionBoundary `json:"partitions"`
}

type mysqlPartitionRow struct {
	name, method, expression, description sql.NullString
	ordinal                               sql.NullInt64
	subpartitionName, subpartitionMethod  sql.NullString
	subpartitionExpression                sql.NullString
	engine                                sql.NullString
}

func readMySQLManagedPartitions(ctx context.Context, q mysqlManagedQuerier, database, table, engine string) (*MySQLPartitionEvidence, error) {
	rows, err := q.QueryContext(ctx, `
		SELECT p.PARTITION_NAME, p.PARTITION_METHOD, p.PARTITION_EXPRESSION, p.PARTITION_DESCRIPTION,
		       p.PARTITION_ORDINAL_POSITION, p.SUBPARTITION_NAME, p.SUBPARTITION_METHOD, p.SUBPARTITION_EXPRESSION,
		       t.ENGINE
		FROM information_schema.PARTITIONS p
		JOIN information_schema.TABLES t
		  ON t.TABLE_SCHEMA = p.TABLE_SCHEMA AND t.TABLE_NAME = p.TABLE_NAME
		WHERE p.TABLE_SCHEMA = ? AND p.TABLE_NAME = ?
		ORDER BY p.PARTITION_ORDINAL_POSITION`, database, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var partitions []mysqlPartitionRow
	for rows.Next() {
		var row mysqlPartitionRow
		if err := rows.Scan(&row.name, &row.method, &row.expression, &row.description, &row.ordinal,
			&row.subpartitionName, &row.subpartitionMethod, &row.subpartitionExpression, &row.engine); err != nil {
			return nil, err
		}
		partitions = append(partitions, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return mysqlManagedPartitionEvidence(partitions, engine)
}

func mysqlManagedPartitionEvidence(rows []mysqlPartitionRow, engine string) (*MySQLPartitionEvidence, error) {
	if len(rows) == 0 {
		return nil, mysqlManagedUnconfirmed("mysql partition catalog is empty")
	}
	if len(rows) == 1 && mysqlPartitionRowEmpty(rows[0]) {
		return nil, nil // an unpartitioned table has one catalog row with all fields NULL
	}
	if !strings.EqualFold(strings.TrimSpace(engine), mysqlPartitionEngineInnoDB) {
		return nil, mysqlManagedUnconfirmed("partitioned table is not InnoDB")
	}

	evidence := &MySQLPartitionEvidence{Method: mysqlPartitionMethodRangeColumns, Expression: mysqlPartitionExpression}
	seenNames := make(map[string]struct{}, len(rows))
	var previous time.Time
	hasFinite := false
	maxValueSeen := false
	for index, row := range rows {
		if mysqlPartitionRowEmpty(row) {
			return nil, mysqlManagedUnconfirmed("mysql partition catalog is incomplete")
		}
		name := strings.TrimSpace(row.name.String)
		if name == "" {
			return nil, mysqlManagedUnconfirmed("mysql partition name is missing")
		}
		key := strings.ToLower(name)
		if _, duplicate := seenNames[key]; duplicate {
			return nil, mysqlManagedUnconfirmed("mysql partition name is duplicated")
		}
		seenNames[key] = struct{}{}
		if !strings.EqualFold(strings.TrimSpace(row.engine.String), mysqlPartitionEngineInnoDB) {
			return nil, mysqlManagedUnconfirmed("mysql partition row is not InnoDB")
		}
		if !strings.EqualFold(strings.TrimSpace(row.method.String), mysqlPartitionMethodRangeColumns) {
			return nil, mysqlManagedUnconfirmed("mysql partition method is not RANGE COLUMNS")
		}
		expression, err := canonicalMySQLPartitionExpression(row.expression.String)
		if err != nil {
			return nil, err
		}
		if expression != evidence.Expression {
			return nil, mysqlManagedUnconfirmed("mysql partition expression is not bucket_start")
		}
		if row.subpartitionName.Valid && strings.TrimSpace(row.subpartitionName.String) != "" ||
			row.subpartitionMethod.Valid && strings.TrimSpace(row.subpartitionMethod.String) != "" ||
			row.subpartitionExpression.Valid && strings.TrimSpace(row.subpartitionExpression.String) != "" {
			return nil, mysqlManagedUnconfirmed("mysql subpartitions are not supported")
		}
		if !row.ordinal.Valid || row.ordinal.Int64 != int64(index+1) {
			return nil, mysqlManagedUnconfirmed("mysql partition ordinals are incomplete")
		}
		lessThan, boundaryTime, isMax, err := normalizeMySQLPartitionBoundary(row.description.String)
		if err != nil {
			return nil, err
		}
		if isMax {
			if index != len(rows)-1 || maxValueSeen {
				return nil, mysqlManagedUnconfirmed("mysql MAXVALUE partition must be last and unique")
			}
			maxValueSeen = true
		} else {
			if maxValueSeen {
				return nil, mysqlManagedUnconfirmed("mysql partition follows MAXVALUE")
			}
			if hasFinite && !previous.Before(boundaryTime) {
				return nil, mysqlManagedUnconfirmed("mysql partition boundaries are not strictly increasing")
			}
			previous, hasFinite = boundaryTime, true
		}
		evidence.Partitions = append(evidence.Partitions, MySQLPartitionBoundary{Name: name, LessThan: lessThan, Ordinal: index + 1})
	}
	if !hasFinite || !maxValueSeen {
		return nil, mysqlManagedUnconfirmed("mysql partition catalog has no finite boundary or MAXVALUE catch-all")
	}
	return evidence, nil
}

func mysqlPartitionRowEmpty(row mysqlPartitionRow) bool {
	emptyString := func(value sql.NullString) bool { return !value.Valid || strings.TrimSpace(value.String) == "" }
	return emptyString(row.name) && emptyString(row.method) && emptyString(row.expression) && emptyString(row.description) &&
		!row.ordinal.Valid && emptyString(row.subpartitionName) && emptyString(row.subpartitionMethod) && emptyString(row.subpartitionExpression)
}

func canonicalMySQLPartitionExpression(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if len(value) >= 2 && value[0] == '`' && value[len(value)-1] == '`' {
		value = value[1 : len(value)-1]
	}
	if !strings.EqualFold(value, mysqlPartitionExpression) {
		return "", mysqlManagedUnconfirmed("mysql partition expression is foreign")
	}
	return mysqlPartitionExpression, nil
}

func normalizeMySQLPartitionBoundary(raw string) (canonical string, parsed time.Time, isMax bool, err error) {
	value := strings.TrimSpace(raw)
	if strings.EqualFold(value, mysqlPartitionMaxValue) {
		return mysqlPartitionMaxValue, time.Time{}, true, nil
	}
	if strings.HasPrefix(value, "'") || strings.HasSuffix(value, "'") {
		if len(value) < 2 || value[0] != '\'' || value[len(value)-1] != '\'' {
			return "", time.Time{}, false, mysqlManagedUnconfirmed("mysql partition boundary quoting is invalid")
		}
		value = strings.TrimSpace(value[1 : len(value)-1])
	}
	for _, layout := range []string{"2006-01-02 15:04:05.999999999", "2006-01-02 15:04:05", time.RFC3339Nano} {
		if parsed, parseErr := time.ParseInLocation(layout, value, time.UTC); parseErr == nil {
			parsed = parsed.UTC()
			if parsed.Nanosecond()%1000 != 0 {
				return "", time.Time{}, false, mysqlManagedUnconfirmed("mysql partition boundary exceeds microseconds")
			}
			return parsed.Format(time.RFC3339Nano), parsed, false, nil
		}
	}
	return "", time.Time{}, false, mysqlManagedUnconfirmed("mysql partition boundary is invalid")
}

func cloneMySQLPartitionEvidence(partitions *MySQLPartitionEvidence) *MySQLPartitionEvidence {
	if partitions == nil {
		return nil
	}
	clone := *partitions
	clone.Partitions = slices.Clone(partitions.Partitions)
	return &clone
}

func sameMySQLPartitionEvidence(left, right *MySQLPartitionEvidence) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return left.Method == right.Method && left.Expression == right.Expression &&
		slices.Equal(left.Partitions, right.Partitions)
}
