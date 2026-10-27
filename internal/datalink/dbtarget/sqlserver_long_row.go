package dbtarget

import (
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

// sqlServerLongColumnUsable accepts only the fixed long-table representation.
// In particular TIMESTAMP/ROWVERSION is never an event time in SQL Server.
func sqlServerLongColumnUsable(column ColumnInfo, role longColumnRole) bool {
	declared := strings.ToUpper(strings.Join(strings.Fields(column.DataType), ""))
	switch role {
	case longTime:
		return declared == sqlServerDateTime2SQLType
	case longText:
		if column.Name == longColRecordID {
			return declared == sqlServerHashSQLType
		}
		return declared == "NVARCHAR(191)"
	case longProvenanceRole, longSlotText:
		return declared == sqlServerNvarcharMaxSQLType
	case longSlotInt:
		return declared == sqlServerBigintSQLType
	case longSlotDecimal, longSlotDelta:
		return declared == "DECIMAL(38,18)"
	case longSlotFloat:
		return declared == "FLOAT(53)"
	case longSlotBool:
		return declared == sqlServerBitSQLType
	}
	return false
}

// SQL Server pads strings for equality even with BIN2 collation. Refuse IDs
// ending in a space so two different source identities cannot compare equal.
// NVARCHAR lengths count UTF-16 units, including two units for astral runes.
func validateSQLServerLongIdentity(spec LongRowSpec) []LayoutIssue {
	if spec.Dialect != SQLDialectSQLServer {
		return nil
	}
	type identity struct{ column, member, value string }
	values := []identity{
		{longColWorkspaceID, "", spec.WorkspaceID},
		{longColGroupID, "", spec.GroupID},
		{longColGroupRev, "", spec.GroupRevision},
		{longColProfileKind, "", spec.ProfileKind},
	}
	for _, member := range spec.Members {
		values = append(values,
			identity{longColDeviceID, member.MemberKey, member.DeviceID},
			identity{longColPointID, member.MemberKey, member.PointID},
			identity{longColTagID, member.MemberKey, member.TagID},
			identity{longColCircuitID, member.MemberKey, member.CircuitID},
			identity{longColMeasurement, member.MemberKey, member.MeasurementID},
		)
		if member.Counter != nil {
			values = append(values,
				identity{longColSeriesEpoch, member.MemberKey, member.Counter.SeriesEpoch},
				identity{longColUnit, member.MemberKey, member.Counter.Unit},
			)
		}
	}
	var issues []LayoutIssue
	for _, value := range values {
		if !utf8.ValidString(value.value) || strings.HasSuffix(value.value, " ") || len(utf16.Encode([]rune(value.value))) > 191 {
			issues = append(issues, LayoutIssue{Code: issueIdentityUnsupported, Column: value.column, MemberKey: value.member})
		}
	}
	return issues
}
