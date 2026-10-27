package recordingplan

// VerifyManagedSQLServerCatalog checks a previously read SQL Server group/receipt
// snapshot against the confirmed canonical layout and returns the same catalog
// digest schema confirmation records. It performs no database calls.
func VerifyManagedSQLServerCatalog(layout GroupSchemaLayout, database string, group, receipt TargetTableInspection) (string, error) {
	if _, err := GenerateGroupSchemaDDL(dialectSQLServer, layout); err != nil {
		return "", err
	}
	found, err := inspectSQLServerGroupSnapshot(layout, database, group, receipt)
	if err != nil {
		return "", err
	}
	if found.pendingGroup || found.pendingReceipt {
		return "", ErrTargetInspectionUnconfirmed
	}
	return schemaTablesDigest(found.tables), nil
}
