package recordingplan

import "fmt"

// checkPreservedSQLServerCatalog binds objects the confirmed operation promised
// to leave unchanged. Newly created objects acquire their first catalog identity
// after execution, so comparing them with the prior missing state would refuse
// a successful CREATE or lost-acknowledgement recovery.
func checkPreservedSQLServerCatalog(before, after []SchemaPreviewTable) error {
	actual := make(map[string]SchemaPreviewTable, len(after))
	for _, table := range after {
		actual[table.Name] = table
	}
	for _, expected := range before {
		if expected.Action == SchemaTableActionCreate {
			continue
		}
		current, found := actual[expected.Name]
		if expected.Action != SchemaTableActionUnchanged || expected.CatalogDigest == "" || !found ||
			current.Action != SchemaTableActionUnchanged || current.CatalogDigest != expected.CatalogDigest {
			return fmt.Errorf("%w: %w", ErrPreviewTokenStale, ErrSchemaTargetChanged)
		}
	}
	return nil
}
