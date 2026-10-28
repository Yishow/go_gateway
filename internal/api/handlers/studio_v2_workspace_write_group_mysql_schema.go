package handlers

import (
	"context"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"

	"github.com/go-sql-driver/mysql"
)

func mysqlSchemaResumeRequired(token *recordingplan.SchemaPreviewToken, operation *recordingplan.SchemaOperation) bool {
	return token != nil && token.GroupLayout != nil && token.GroupLayout.MySQL != nil && operation != nil &&
		(operation.Status == recordingplan.SchemaOperationPartial || operation.Status == recordingplan.SchemaOperationUnknown)
}

func (h *StudioV2WorkspaceWriteGroupSchemaHandler) executeGroupSchemaStatements(ctx context.Context, id string, req WriteGroupSchemaRequest, scope recordingplan.SchemaPreviewScope, statements []string) (recordingplan.SchemaExecution, error) {
	target := h.checkedTarget(id, req, scope)
	return target.executeGroupSchemaStatements(ctx, scope, statements)
}

func (h *StudioV2WorkspaceWriteGroupSchemaHandler) mysqlGroupTargetInspector(id string, req WriteGroupSchemaRequest, scope recordingplan.SchemaPreviewScope) recordingplan.TargetInspector {
	target := h.checkedTarget(id, req, scope)
	return target.mysqlGroupTargetInspector(scope)
}

func mirrorMySQLPartitionEvidence(evidence *dbtarget.MySQLPartitionEvidence) *recordingplan.MySQLPartitionEvidence {
	if evidence == nil {
		return nil
	}
	var partitions []recordingplan.MySQLPartitionBoundary
	if evidence.Partitions != nil {
		partitions = make([]recordingplan.MySQLPartitionBoundary, len(evidence.Partitions))
		for i, boundary := range evidence.Partitions {
			partitions[i] = recordingplan.MySQLPartitionBoundary{
				Name:     boundary.Name,
				LessThan: boundary.LessThan,
				Ordinal:  boundary.Ordinal,
			}
		}
	}
	return &recordingplan.MySQLPartitionEvidence{
		Method:     evidence.Method,
		Expression: evidence.Expression,
		Partitions: partitions,
	}
}

func mysqlGroupSchemaError(err error) error {
	if errors.Is(err, dbtarget.ErrMySQLManagedCapabilityChanged) {
		return fmt.Errorf("%w: %w", recordingplan.ErrPreviewTokenStale, err)
	}
	if errors.Is(err, dbtarget.ErrMySQLManagedIndexUnusable) {
		return fmt.Errorf("%w: %w", recordingplan.ErrIncompatibleExistingTable, err)
	}
	var driverError *mysql.MySQLError
	if errors.Is(err, dbtarget.ErrMySQLManagedTableForbidden) ||
		(errors.As(err, &driverError) && (driverError.Number == 1044 || driverError.Number == 1142 || driverError.Number == 1143 || driverError.Number == 1227)) {
		return fmt.Errorf("%w: %w", recordingplan.ErrTargetPermissionDenied, err)
	}
	if errors.Is(err, dbtarget.ErrMySQLManagedTableUnconfirmed) {
		return fmt.Errorf("%w: %w", recordingplan.ErrTargetInspectionUnconfirmed, err)
	}
	return err
}
