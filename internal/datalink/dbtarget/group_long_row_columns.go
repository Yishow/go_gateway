package dbtarget

import (
	"regexp"

	"go-gateway/internal/datalink/measurement"
)

var longOwnerPattern = regexp.MustCompile(`^_gw_owner_[0-9a-f]{24}$`)

// longOwnerColumn finds the single owner marker carried by the verified table.
// The runtime proof separately binds its name to the group's storage scope.
func longOwnerColumn(dialect SQLDialect, columns []ColumnInfo) (string, bool) {
	var owner string
	for _, column := range columns {
		if !longOwnerPattern.MatchString(column.Name) {
			continue
		}
		if owner != "" || column.Nullable || column.PrimaryKey || !isTextColumn(dialect, column) {
			return "", false
		}
		owner = column.Name
	}
	return owner, owner != ""
}

// Fixed managed_long_v1 column names. They are the documented schema contract.
const (
	longColBucketStart  = "bucket_start"
	longColBucketEnd    = "bucket_end"
	longColRecordID     = "record_id"
	longColWorkspaceID  = "workspace_id"
	longColGroupID      = "group_id"
	longColGroupRev     = "group_revision"
	longColCircuitID    = "circuit_id"
	longColDeviceID     = "device_id"
	longColPointID      = "point_id"
	longColTagID        = "tag_id"
	longColMeasurement  = "measurement_id"
	longColSeriesEpoch  = "series_epoch"
	longColProfileKind  = "profile_kind"
	longColValueType    = "value_type"
	longColUnit         = "unit"
	longColStatus       = "status"
	longColQuality      = "quality"
	longColObservedAt   = "observed_at"
	longColReceivedAt   = "received_at"
	longColValueInt     = "value_int"
	longColValueDecimal = "value_decimal"
	longColValueFloat   = "value_float"
	longColValueBool    = "value_bool"
	longColValueText    = "value_text"
	longColDeltaValue   = "delta_value"
	longColDeltaFrom    = "delta_from"
	longColDeltaTo      = "delta_to"
	longColProvenance   = "provenance"

	// longTimeText is the SQLite form of every long-table time: UTC with a fixed
	// microsecond fraction, so text order is time order.
	longTimeText = "2006-01-02T15:04:05.000000Z"

	longStatusOK = "ok"
	// longStatusPartial marks a good value that covers only part of its window.
	longStatusPartial = "partial"
)

// longColumnRole says how one fixed column is written.
type longColumnRole int

const (
	longText longColumnRole = iota
	longTime
	longProvenanceRole
	longSlotInt
	longSlotDecimal
	longSlotFloat
	longSlotBool
	longSlotText
	longSlotDelta
)

var longColumnRoles = []struct {
	name     string
	role     longColumnRole
	required bool // NOT NULL in the fixed layout
}{
	{longColBucketStart, longTime, true}, {longColBucketEnd, longTime, true},
	{longColRecordID, longText, true}, {longColWorkspaceID, longText, true},
	{longColGroupID, longText, true}, {longColGroupRev, longText, true},
	{longColCircuitID, longText, false}, {longColDeviceID, longText, true},
	{longColPointID, longText, true}, {longColTagID, longText, true},
	{longColMeasurement, longText, false}, {longColSeriesEpoch, longText, false},
	{longColProfileKind, longText, true}, {longColValueType, longText, true},
	{longColUnit, longText, false}, {longColStatus, longText, true}, {longColQuality, longText, false},
	{longColObservedAt, longTime, false}, {longColReceivedAt, longTime, false},
	{longColValueInt, longSlotInt, false}, {longColValueDecimal, longSlotDecimal, false},
	{longColValueFloat, longSlotFloat, false}, {longColValueBool, longSlotBool, false},
	{longColValueText, longSlotText, false}, {longColDeltaValue, longSlotDelta, false},
	{longColDeltaFrom, longTime, false}, {longColDeltaTo, longTime, false},
	{longColProvenance, longProvenanceRole, true},
}

func longColumnUsable(dialect SQLDialect, column ColumnInfo, role longColumnRole) bool {
	if dialect == SQLDialectSQLServer {
		return sqlServerLongColumnUsable(column, role)
	}
	exact := func(types ...measurement.ExactType) bool {
		for _, typ := range types {
			if _, err := exactStrategyFor(dialect, column.DataType, typ); err != nil {
				return false
			}
		}
		return true
	}
	switch role {
	case longText:
		return isTextColumn(dialect, column)
	case longTime:
		if dialect == SQLDialectMySQL {
			return mysqlLongDatetime6Column(column)
		}
		return bucketColumnKind(dialect, column) != bucketUnsupported
	case longProvenanceRole:
		if dialect == SQLDialectMySQL {
			return mysqlLongProvenanceColumn(column)
		}
		return isProvenanceColumn(dialect, column)
	case longSlotInt:
		return exact(measurement.ExactInt64)
	case longSlotDecimal:
		return exact(measurement.ExactDecimal, measurement.ExactUint64)
	case longSlotFloat:
		return exact(measurement.ExactFloat64)
	case longSlotBool:
		return exact(measurement.ExactBool)
	case longSlotText:
		return exact(measurement.ExactText)
	case longSlotDelta:
		return exact(measurement.ExactDecimal)
	}
	return false
}

func mysqlLongDatetime6Column(column ColumnInfo) bool {
	declared, err := parseDeclaredType(column.DataType)
	return err == nil && declared.base == mysqlTypeDatetime && declared.hasP && declared.p == 6 && !declared.hasS
}

func mysqlLongProvenanceColumn(column ColumnInfo) bool {
	declared, err := parseDeclaredType(column.DataType)
	return err == nil && (declared.base == "JSON" || declared.base == "LONGTEXT")
}

// longSlotFor names the one value column an exact type is stored in ("" when none). uint64 and
// decimal share the exact decimal slot so neither passes through a float.
func longSlotFor(typ measurement.ExactType) string {
	switch typ {
	case measurement.ExactInt64:
		return longColValueInt
	case measurement.ExactUint64, measurement.ExactDecimal:
		return longColValueDecimal
	case measurement.ExactFloat64:
		return longColValueFloat
	case measurement.ExactBool:
		return longColValueBool
	case measurement.ExactText:
		return longColValueText
	}
	return ""
}
