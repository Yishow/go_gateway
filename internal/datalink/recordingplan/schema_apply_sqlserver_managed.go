package recordingplan

import (
	"context"
)

// SQLServerSchemaApplyTarget is the narrow local SQL Server confirmation
// target. Execute is supplied only by the saved-group SQL Server adapter; this
// path is separate from the generic ApplySchemaPreview verification gate.
type SQLServerSchemaApplyTarget struct {
	Inspect TargetInspector
	Execute func(context.Context, *SchemaPreviewToken, []string) (SchemaExecution, error)
}

// ApplySQLServerSchemaPreview reuses the existing confirmed preview/ledger/
// catalog verification engine while accepting only the explicit SQL Server
// typed target. Generic ApplySchemaPreview remains closed for SQL Server.
func (s *Service) ApplySQLServerSchemaPreview(ctx context.Context, token *SchemaPreviewToken, target SQLServerSchemaApplyTarget) (*SchemaOperation, ClaimOutcome, error) {
	if token == nil || token.GroupLayout == nil || token.GroupLayout.SQLServer == nil ||
		token.Dialect != dialectSQLServer || target.Inspect == nil || target.Execute == nil {
		return nil, "", ErrSchemaAdapterUnverified
	}
	execute := func(ctx context.Context, statements []string) (SchemaExecution, error) {
		return target.Execute(ctx, token, statements)
	}
	return s.applySchemaPreview(ctx, token, SchemaApplyTarget{Inspect: target.Inspect, Execute: execute}, true)
}
