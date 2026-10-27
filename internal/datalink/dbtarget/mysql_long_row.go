package dbtarget

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"go-gateway/internal/datalink/measurement"
)

// ValidateSampleValue checks bounded exact SQL value representability before intake.
// Other dialects retain their existing intake behavior.
func (l *LongRowLayout) ValidateSampleValue(value measurement.ExactValue) error {
	if l == nil {
		return blockedf("long-row layout is unavailable")
	}
	if l.spec.Dialect != SQLDialectMySQL && l.spec.Dialect != SQLDialectSQLServer {
		return nil
	}
	slot := longSlotFor(value.Type())
	if slot == "" {
		return blockedf("exact value has no supported long-row slot")
	}
	column, ok := l.columns[slot]
	if !ok {
		return blockedf("verified long-row value slot is unavailable")
	}
	if _, err := EncodeExactValue(l.spec.Dialect, column.DataType, value); err != nil {
		if errors.Is(err, ErrExactSQLBlocked) {
			return err
		}
		return fmt.Errorf("sample value: %w", errors.Join(ErrExactSQLBlocked, err))
	}
	return nil
}

// validateMySQLLongIdentity checks only strings that the fixed long row binds
// into bounded VARCHAR columns. Catalog ownership and collation are separate proofs.
func validateMySQLLongIdentity(spec LongRowSpec) []LayoutIssue {
	if spec.Dialect != SQLDialectMySQL {
		return nil
	}
	columns := make(map[string]ColumnInfo, len(spec.Columns))
	for _, column := range spec.Columns {
		columns[strings.ToLower(column.Name)] = column
	}

	type scopedValue struct {
		column    string
		memberKey string
		value     string
	}
	values := []scopedValue{
		{column: longColWorkspaceID, value: spec.WorkspaceID},
		{column: longColGroupID, value: spec.GroupID},
		{column: longColGroupRev, value: spec.GroupRevision},
		{column: longColProfileKind, value: spec.ProfileKind},
	}
	for _, member := range spec.Members {
		values = append(values,
			scopedValue{column: longColCircuitID, memberKey: member.MemberKey, value: member.CircuitID},
			scopedValue{column: longColDeviceID, memberKey: member.MemberKey, value: member.DeviceID},
			scopedValue{column: longColPointID, memberKey: member.MemberKey, value: member.PointID},
			scopedValue{column: longColTagID, memberKey: member.MemberKey, value: member.TagID},
			scopedValue{column: longColMeasurement, memberKey: member.MemberKey, value: member.MeasurementID},
			scopedValue{column: longColValueType, memberKey: member.MemberKey, value: string(member.Type)},
		)
		if member.Counter != nil {
			values = append(values,
				scopedValue{column: longColSeriesEpoch, memberKey: member.MemberKey, value: member.Counter.SeriesEpoch},
				scopedValue{column: longColUnit, memberKey: member.MemberKey, value: member.Counter.Unit},
			)
		}
	}

	var issues []LayoutIssue
	for _, item := range values {
		if item.value == "" {
			continue
		}
		column, ok := columns[strings.ToLower(item.column)]
		if !ok {
			continue // The layout constructor reports missing columns separately.
		}
		limit, bounded := mysqlBoundedVarcharLength(column.DataType)
		if !bounded || utf8.RuneCountInString(item.value) > limit {
			issues = append(issues, LayoutIssue{
				Code: issueIdentityUnsupported, Column: item.column, MemberKey: item.memberKey,
			})
		}
	}
	return issues
}

func mysqlBoundedVarcharLength(declared string) (int, bool) {
	column, err := parseDeclaredType(declared)
	if err != nil || column.base != mysqlTypeVarchar || !column.hasP || column.hasS || column.p < 1 {
		return 0, false
	}
	return column.p, true
}
