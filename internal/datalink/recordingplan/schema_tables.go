package recordingplan

import (
	"strings"
)

// managedTableSpec lists the columns and expected type families a compatible managed table must expose.
type managedTableSpec struct {
	name             string
	columns          []string
	columnCategories map[string]string
}

func makeCategoryMap(specs ...string) map[string]string {
	res := make(map[string]string)
	for i := 0; i+1 < len(specs); i += 2 {
		cat := specs[i]
		for _, col := range strings.Fields(specs[i+1]) {
			res[col] = cat
		}
	}
	return res
}

// managedTableSpecs mirrors the columns GenerateManagedSchemaDDL creates.
var managedTableSpecs = []managedTableSpec{
	{
		name: "samples",
		columns: strings.Fields(`workspace_id plan_id stream_id record_id measurement_id series_epoch
			observed_at received_at quality quality_reason value_type val_num val_dec val_str val_bool is_test`),
		columnCategories: makeCategoryMap(
			"text", "workspace_id plan_id stream_id record_id measurement_id series_epoch quality quality_reason value_type val_dec val_str",
			"time", "observed_at received_at",
			"float", "val_num",
			"bool", "val_bool is_test",
		),
	},
	{
		name: "intervals",
		columns: strings.Fields(`workspace_id plan_id stream_id record_id measurement_id series_epoch
			interval_start interval_end calculation_revision sample_count mean_val min_val max_val
			quantity_delta known_subtotal is_estimated is_complete is_test`),
		columnCategories: makeCategoryMap(
			"text", "workspace_id plan_id stream_id record_id measurement_id series_epoch calculation_revision",
			"time", "interval_start interval_end",
			"int", "sample_count",
			"float", "mean_val min_val max_val quantity_delta known_subtotal",
			"bool", "is_estimated is_complete is_test",
		),
	},
	{
		name: "events",
		columns: strings.Fields(`workspace_id plan_id stream_id record_id measurement_id series_epoch
			observed_at event_type state_from state_to duration_ms message is_test`),
		columnCategories: makeCategoryMap(
			"text", "workspace_id plan_id stream_id record_id measurement_id series_epoch event_type state_from state_to message",
			"time", "observed_at",
			"int", "duration_ms",
			"bool", "is_test",
		),
	},
	{
		name: "snapshots",
		columns: strings.Fields(`workspace_id plan_id stream_id record_id batch_id trigger_id
			observed_at completeness payload_json is_test`),
		columnCategories: makeCategoryMap(
			"text", "workspace_id plan_id stream_id record_id batch_id trigger_id completeness",
			"time", "observed_at",
			"json", "payload_json",
			"bool", "is_test",
		),
	},
	{
		name: "definitions",
		columns: strings.Fields(`workspace_id plan_id measurement_id definition_revision series_epoch
			quantity unit semantic_kind schema_json created_at`),
		columnCategories: makeCategoryMap(
			"text", "workspace_id plan_id measurement_id definition_revision series_epoch quantity unit semantic_kind",
			"time", "created_at",
			"json", "schema_json",
		),
	},
	{
		name: "receipts",
		columns: strings.Fields(`workspace_id plan_id destination_id batch_id last_record_id
			delivered_at record_count`),
		columnCategories: makeCategoryMap(
			"text", "workspace_id plan_id destination_id batch_id last_record_id",
			"time", "delivered_at",
			"int", "record_count",
		),
	},
}

// statementsForTables keeps the CREATE TABLE and index statements of pending tables.
func statementsForTables(statements []string, prefix string, pending map[string]bool) []string {
	kept := make([]string, 0, len(statements))
	for _, stmt := range statements {
		for _, spec := range managedTableSpecs {
			table := prefix + spec.name
			if pending[spec.name] && (strings.Contains(stmt, "EXISTS "+table+" (") || strings.Contains(stmt, " ON "+table+"(")) {
				kept = append(kept, stmt)
				break
			}
		}
	}
	return kept
}

func isColumnTypeCompatible(expectedCategory, actualType string) bool {
	actual := strings.ToLower(strings.TrimSpace(actualType))
	if actual == "" {
		return true
	}
	switch expectedCategory {
	case "text", "json":
		return strings.Contains(actual, "char") || strings.Contains(actual, "text") ||
			strings.Contains(actual, "clob") || strings.Contains(actual, "json") ||
			strings.Contains(actual, "string") || strings.Contains(actual, "uuid")
	case "float":
		return strings.Contains(actual, "real") || strings.Contains(actual, "float") ||
			strings.Contains(actual, "double") || strings.Contains(actual, "numeric") ||
			strings.Contains(actual, "decimal") || strings.Contains(actual, "num")
	case "int":
		return strings.Contains(actual, "int") || strings.Contains(actual, "serial") ||
			strings.Contains(actual, "numeric") || strings.Contains(actual, "decimal")
	case "bool":
		return strings.Contains(actual, "bool") || strings.Contains(actual, "bit") ||
			strings.Contains(actual, "int") || strings.Contains(actual, "tinyint")
	case "time":
		return strings.Contains(actual, "time") || strings.Contains(actual, "date") ||
			strings.Contains(actual, "text") // SQLite stores datetime as TEXT
	default:
		return true
	}
}

// incompatibleColumnType reports the first column whose actual type does not
// satisfy its expected category; empty results mean every type is compatible.
func incompatibleColumnType(categories, actualTypes map[string]string) (incompatibleCol, incompatibleType string) {
	if len(actualTypes) == 0 {
		return "", ""
	}
	for col, category := range categories {
		if actual, ok := actualTypes[col]; ok && actual != "" {
			if !isColumnTypeCompatible(category, actual) {
				return col, actual
			}
		}
	}
	return "", ""
}
