package dbtarget

import (
	"context"
	"fmt"
	"strings"
)

const sqlServerLongTestCollation = "Latin1_General_100_BIN2"

func readSQLServerGroupLongTestRows(ctx context.Context, q groupLongQueryer, ref GroupLongTestRowRef, columns []string, lock bool) ([][]any, error) {
	selected := append([]string{}, columns...)
	for _, mandatory := range []string{longColRecordID, longColWorkspaceID, longColGroupID, longColGroupRev, longColProvenance, ref.OwnerColumn} {
		found := false
		for _, column := range selected {
			if strings.EqualFold(column, mandatory) {
				found = true
				break
			}
		}
		if !found {
			selected = append(selected, mandatory)
		}
	}
	quoted := make([]string, len(selected))
	for index, column := range selected {
		quoted[index] = quoteIdentifier(ref.Kind, column)
	}
	query, args := sqlServerGroupLongTestSelect(ref, quoted, lock)
	rows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("read SQL Server long test rows: %w", err)
	}
	defer func() { _ = rows.Close() }()

	index := make(map[string]int, len(selected))
	for i, column := range selected {
		index[strings.ToLower(column)] = i
	}
	requested := make([]int, len(columns))
	for i, column := range columns {
		at, ok := index[strings.ToLower(column)]
		if !ok {
			return nil, fmt.Errorf("%w: requested SQL Server readback column is missing", ErrOwnedRowRefInvalid)
		}
		requested[i] = at
	}
	expected := make(map[string]string, len(ref.RecordIDs))
	for i, id := range ref.RecordIDs {
		expected[id] = ref.SampleIDs[i]
	}
	var out [][]any
	seenRecords := make(map[string]struct{}, MaxGroupLongTestRows)
	for rows.Next() {
		values := make([]any, len(selected))
		dest := make([]any, len(selected))
		for i := range dest {
			dest[i] = &values[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, fmt.Errorf("scan SQL Server long test row: %w", err)
		}
		if err := validateGroupLongTestRow(values, index, expected, ref); err != nil {
			return nil, err
		}
		recordID, ok := groupLongText(values[index[longColRecordID]])
		if !ok {
			return nil, fmt.Errorf("%w: SQL Server long test record identity is unreadable", ErrOwnedRowRefInvalid)
		}
		if _, duplicate := seenRecords[recordID]; duplicate {
			return nil, fmt.Errorf("%w: duplicate SQL Server long test record", ErrOwnedRowRefInvalid)
		}
		seenRecords[recordID] = struct{}{}
		if len(out) >= MaxGroupLongTestRows {
			return nil, fmt.Errorf("%w: SQL Server long test row limit exceeded", ErrOwnedRowRefInvalid)
		}
		row := make([]any, len(requested))
		for i, at := range requested {
			row[i] = normalizeGroupLongValue(values[at])
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read SQL Server long test rows: %w", err)
	}
	return out, nil
}

func sqlServerGroupLongTestSelect(ref GroupLongTestRowRef, quotedColumns []string, lock bool) (query string, args []any) {
	args = make([]any, len(ref.RecordIDs))
	conditions := make([]string, len(ref.RecordIDs))
	for index, recordID := range ref.RecordIDs {
		args[index] = recordID
		conditions[index] = sqlServerLongExactRecord(quoteIdentifier(ref.Kind, longColRecordID), fmt.Sprintf("@p%d", index+1))
	}
	lockHint := ""
	if lock {
		lockHint = " WITH (UPDLOCK, HOLDLOCK)"
	}
	// Fetch by the bounded operation-derived IDs alone. Scope, canonical owner
	// and exact JSON sample identity are then verified in Go; filtering on them
	// here would hide a conflicting row from cleanup preflight.
	query = fmt.Sprintf("SELECT TOP (%d) %s FROM %s%s WHERE (%s)",
		MaxGroupLongTestRows+1, strings.Join(quotedColumns, ", "),
		qualifiedTableName(ref.Kind, ref.SchemaName, ref.TableName), lockHint, strings.Join(conditions, " OR "))
	return query, args
}

func sqlServerLongExactText(expression, parameter string) string {
	collated := expression + " COLLATE " + sqlServerLongTestCollation
	bound := parameter + " COLLATE " + sqlServerLongTestCollation
	return "(" + collated + " = " + bound + " AND DATALENGTH(" + expression + ") = DATALENGTH(" + parameter + "))"
}

func sqlServerLongExactRecord(expression, parameter string) string {
	converted := "CONVERT(nvarchar(64), " + expression + ")"
	return sqlServerLongExactText(converted, parameter)
}
