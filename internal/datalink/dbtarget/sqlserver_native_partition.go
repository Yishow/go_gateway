package dbtarget

const sqlServerNativePrimaryFilegroup = "PRIMARY"

// SQLServerPartitionBoundary describes one canonical range boundary in the
// frozen SQL Server native monthly partition proof.
type SQLServerPartitionBoundary struct {
	Name     string
	LessThan string
	Ordinal  int
}

// SQLServerManagedPartitionEvidence is the inspected native partition shape
// for an owned SQL Server group table. A verified absence is represented by
// PartitionInspected on SQLServerManagedTableEvidence with NativeMonthly nil.
type SQLServerManagedPartitionEvidence struct {
	Method         string
	Column         string
	FunctionName   string
	SchemeName     string
	Filegroup      string
	Partitions     []SQLServerPartitionBoundary
	FunctionExists bool
	SchemeExists   bool
	Exclusive      bool
	IndexesAligned bool
}
