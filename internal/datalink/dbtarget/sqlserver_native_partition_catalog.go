package dbtarget

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

const sqlServerNativeBoundaryTimeLayout = "2006-01-02T15:04:05.000000"

type sqlServerNativeObjectFacts struct {
	functionID, schemeID, schemeFunctionID sql.NullInt64
	right                                  sql.NullBool
	functionType, parameterType            sql.NullString
	scale                                  sql.NullInt64
	parameterCount, schemeCount            int64
	functionOwner, schemeOwner             sql.NullString
	visible                                sql.NullBool
}

func readSQLServerManagedPartitions(ctx context.Context, q sqlServerManagedCatalogQuerier, schemaName, tableName string, objectID int64) (*SQLServerManagedPartitionEvidence, bool, error) {
	if tableName == EffectReceiptTable {
		return nil, false, nil
	}
	if schemaName == "" || tableName == "" || strings.TrimSpace(tableName) != tableName {
		return nil, false, ErrSQLServerManagedTableUnconfirmed
	}
	function, scheme := "pf_"+tableName, "ps_"+tableName
	var objects sqlServerNativeObjectFacts
	err := q.QueryRowContext(ctx, sqlServerNativeObjectsQuery, function, scheme, "gateway.recording.partition_owner").Scan(
		&objects.functionID, &objects.right, &objects.functionType, &objects.parameterType, &objects.scale,
		&objects.parameterCount, &objects.schemeID, &objects.schemeFunctionID,
		&objects.functionOwner, &objects.schemeOwner, &objects.schemeCount, &objects.visible,
	)
	if err != nil {
		return nil, false, err
	}
	if !objects.visible.Valid || !objects.visible.Bool {
		return nil, false, ErrSQLServerManagedTableUnconfirmed
	}
	if !objects.functionID.Valid && !objects.schemeID.Valid {
		return nil, true, nil
	}
	evidence := &SQLServerManagedPartitionEvidence{
		Column: "bucket_start", FunctionName: function, SchemeName: scheme,
		FunctionExists: objects.functionID.Valid && objects.functionID.Int64 > 0,
		SchemeExists:   objects.schemeID.Valid && objects.schemeID.Int64 > 0,
	}
	if !evidence.FunctionExists || !objects.right.Valid || strings.TrimRight(objects.functionType.String, " ") != "R" ||
		objects.parameterType.String != "datetime2" || !objects.scale.Valid || objects.scale.Int64 != 6 || objects.parameterCount != 1 {
		return evidence, true, nil
	}
	evidence.Method = "RANGE LEFT"
	if objects.right.Bool {
		evidence.Method = sqlServerNativeRangeRight
	}
	evidence.Partitions, err = readSQLServerNativeBoundaries(ctx, q, objects.functionID.Int64)
	if err != nil {
		return nil, false, err
	}
	owner := schemaName + "." + tableName
	functionOwned := objects.functionOwner.Valid && objects.functionOwner.String == owner
	if !evidence.SchemeExists {
		// A function with no referencing scheme can be safely completed only
		// when its own marker is scoped to this table and no other scheme uses it.
		evidence.Exclusive = functionOwned && objects.schemeCount == 0
		return evidence, true, nil
	}
	if !objects.schemeFunctionID.Valid || objects.schemeFunctionID.Int64 != objects.functionID.Int64 {
		return evidence, true, nil
	}
	primary, err := readSQLServerNativeDestinations(ctx, q, objects.schemeID.Int64, len(evidence.Partitions))
	if err != nil {
		return nil, false, err
	}
	if primary {
		evidence.Filegroup = sqlServerNativePrimaryFilegroup
	}
	exclusive, aligned, err := readSQLServerNativeIndexAlignment(ctx, q, objectID, objects.schemeID.Int64)
	if err != nil {
		return nil, false, err
	}
	evidence.Exclusive = exclusive && objects.schemeCount == 1 && functionOwned &&
		objects.schemeOwner.Valid && objects.schemeOwner.String == owner
	evidence.IndexesAligned = aligned
	return evidence, true, nil
}

func readSQLServerNativeBoundaries(ctx context.Context, q sqlServerManagedCatalogQuerier, functionID int64) ([]SQLServerPartitionBoundary, error) {
	rows, err := q.QueryContext(ctx, sqlServerNativeBoundsQuery, functionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []SQLServerPartitionBoundary
	var previous time.Time
	for rows.Next() {
		var ordinal, parameter int
		var at time.Time
		if err := rows.Scan(&ordinal, &parameter, &at); err != nil {
			return nil, err
		}
		at = at.UTC()
		if ordinal != len(result)+1 || parameter != 1 || at.Year() < 1 || at.Year() > 9999 ||
			at.Day() != 1 || at.Hour() != 0 || at.Minute() != 0 || at.Second() != 0 || at.Nanosecond() != 0 ||
			(!previous.IsZero() && !at.After(previous)) {
			return nil, ErrSQLServerManagedTableUnconfirmed
		}
		name := sqlServerNativePastPartitionName
		if !previous.IsZero() {
			name = "p_" + previous.Format("200601")
			if !at.Equal(previous.AddDate(0, 1, 0)) {
				name = "p_gap_" + at.Format("200601")
			}
		}
		result = append(result, SQLServerPartitionBoundary{Name: name, LessThan: at.Format(sqlServerNativeBoundaryTimeLayout), Ordinal: ordinal})
		previous = at
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if len(result) < 5 {
		return nil, ErrSQLServerManagedTableUnconfirmed
	}
	return append(result, SQLServerPartitionBoundary{Name: sqlServerNativeFuturePartitionName, LessThan: "MAXVALUE", Ordinal: len(result) + 1}), nil
}

func readSQLServerNativeDestinations(ctx context.Context, q sqlServerManagedCatalogQuerier, schemeID int64, partitionCount int) (bool, error) {
	rows, err := q.QueryContext(ctx, sqlServerNativeDestinationsQuery, schemeID)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	count, primary := 0, true
	for rows.Next() {
		var ordinal int
		var name, kind string
		if err := rows.Scan(&ordinal, &name, &kind); err != nil {
			return false, err
		}
		count++
		primary = primary && ordinal == count && name == sqlServerNativePrimaryFilegroup && kind == "FG"
	}
	if err := rows.Err(); err != nil {
		return false, err
	}
	if err := rows.Close(); err != nil {
		return false, err
	}
	// SQL Server may retain one additional destination as NEXT USED.
	return primary && (count == partitionCount || count == partitionCount+1), nil
}

func readSQLServerNativeIndexAlignment(ctx context.Context, q sqlServerManagedCatalogQuerier, objectID, schemeID int64) (exclusive, aligned bool, err error) {
	rows, err := q.QueryContext(ctx, sqlServerNativeIndexesQuery, objectID, schemeID)
	if err != nil {
		return false, false, err
	}
	defer rows.Close()
	exclusive, aligned, ownCount := true, objectID > 0, 0
	clusteredUnique := false
	seen := make(map[string]bool)
	for rows.Next() {
		var table, index, space, partitionKeys, uniqueKeys int64
		var unique bool
		if err := rows.Scan(&table, &index, &space, &partitionKeys, &uniqueKeys, &unique); err != nil {
			return false, false, err
		}
		identity := fmt.Sprintf("%d:%d", table, index)
		if seen[identity] || table <= 0 || index < 0 {
			return false, false, ErrSQLServerManagedTableUnconfirmed
		}
		seen[identity] = true
		if table != objectID {
			exclusive = false
			continue
		}
		ownCount++
		clusteredUnique = clusteredUnique || (index == 1 && unique && uniqueKeys == 1)
		aligned = aligned && index > 0 && space == schemeID && partitionKeys == 1 && (!unique || uniqueKeys == 1)
	}
	if err := rows.Err(); err != nil {
		return false, false, err
	}
	if err := rows.Close(); err != nil {
		return false, false, err
	}
	return exclusive, aligned && ownCount == 2 && clusteredUnique, nil
}
