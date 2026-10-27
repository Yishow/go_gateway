package recordingplan

import (
	"maps"
	"slices"
)

func cloneMySQLTableEvidence(evidence *TargetMySQLTableEvidence) *TargetMySQLTableEvidence {
	if evidence == nil {
		return nil
	}
	cloned := *evidence
	cloned.PrimaryKey = slices.Clone(evidence.PrimaryKey)
	cloned.ColumnCollations = maps.Clone(evidence.ColumnCollations)
	if cloned.ColumnCollations == nil {
		cloned.ColumnCollations = make(map[string]string)
	}
	cloned.ColumnDefaults = make(map[string]*string, len(evidence.ColumnDefaults))
	for key, value := range evidence.ColumnDefaults {
		if value != nil {
			defaultValue := *value
			cloned.ColumnDefaults[key] = &defaultValue
		} else {
			cloned.ColumnDefaults[key] = nil
		}
	}
	cloned.JSONColumns = maps.Clone(evidence.JSONColumns)
	if cloned.JSONColumns == nil {
		cloned.JSONColumns = make(map[string]bool)
	}
	cloned.Partitions = cloneMySQLPartitionEvidence(evidence.Partitions)
	return &cloned
}

func cloneTargetTableInspection(state TargetTableInspection) TargetTableInspection {
	cloned := state
	cloned.Columns = slices.Clone(state.Columns)
	cloned.Indexes = slices.Clone(state.Indexes)
	for i := range cloned.Indexes {
		cloned.Indexes[i].Columns = slices.Clone(state.Indexes[i].Columns)
	}
	cloned.ColumnTypes = maps.Clone(state.ColumnTypes)
	if cloned.ColumnTypes == nil {
		cloned.ColumnTypes = make(map[string]string)
	}
	cloned.ColumnNullable = maps.Clone(state.ColumnNullable)
	if cloned.ColumnNullable == nil {
		cloned.ColumnNullable = make(map[string]bool)
	}
	cloned.ColumnPrimaryKey = maps.Clone(state.ColumnPrimaryKey)
	if cloned.ColumnPrimaryKey == nil {
		cloned.ColumnPrimaryKey = make(map[string]bool)
	}
	cloned.MySQL = cloneMySQLTableEvidence(state.MySQL)
	if state.Timescale != nil {
		timescale := *state.Timescale
		timescale.Reasons = slices.Clone(state.Timescale.Reasons)
		if state.Timescale.Hypertable != nil {
			hypertable := *state.Timescale.Hypertable
			timescale.Hypertable = &hypertable
		}
		cloned.Timescale = &timescale
	}
	return cloned
}
