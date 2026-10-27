package dbtarget

import (
	"maps"
	"slices"
)

// CloneSQLServerManagedRowProof returns an independent snapshot, including
// nested index columns and default pointers. Nil and verified-empty evidence
// retain their distinct meanings.
func CloneSQLServerManagedRowProof(proof *SQLServerManagedRowProof) *SQLServerManagedRowProof {
	if proof == nil {
		return nil
	}
	cloned := *proof
	cloned.Group = cloneSQLServerManagedInspection(proof.Group)
	cloned.Receipt = cloneSQLServerManagedInspection(proof.Receipt)
	return &cloned
}

func cloneSQLServerManagedInspection(physical *SQLServerManagedInspection) *SQLServerManagedInspection {
	if physical == nil {
		return nil
	}
	cloned := *physical
	cloned.Indexes = cloneSQLServerIndexes(physical.Indexes)
	if physical.Inspection != nil {
		inspection := *physical.Inspection
		inspection.Columns = slices.Clone(inspection.Columns)
		inspection.Indexes = cloneSQLServerIndexes(inspection.Indexes)
		cloned.Inspection = &inspection
	}
	if physical.Evidence != nil {
		evidence := *physical.Evidence
		evidence.NativeMonthly = cloneSQLServerManagedPartitionEvidence(evidence.NativeMonthly)
		evidence.PrimaryKey = slices.Clone(evidence.PrimaryKey)
		evidence.ColumnCollations = maps.Clone(evidence.ColumnCollations)
		evidence.ColumnDefaults = maps.Clone(evidence.ColumnDefaults)
		for name, value := range evidence.ColumnDefaults {
			if value != nil {
				independent := *value
				evidence.ColumnDefaults[name] = &independent
			}
		}
		cloned.Evidence = &evidence
	}
	return &cloned
}

func cloneSQLServerManagedPartitionEvidence(evidence *SQLServerManagedPartitionEvidence) *SQLServerManagedPartitionEvidence {
	if evidence == nil {
		return nil
	}
	cloned := *evidence
	cloned.Partitions = slices.Clone(evidence.Partitions)
	return &cloned
}

func cloneSQLServerIndexes(indexes []IndexInfo) []IndexInfo {
	cloned := slices.Clone(indexes)
	for i := range cloned {
		cloned[i].Columns = slices.Clone(cloned[i].Columns)
	}
	return cloned
}
