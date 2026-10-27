package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func TestMySQLManagedSchemaShapeClassification(t *testing.T) {
	physical := errors.Join(ErrMySQLManagedPhysicalMismatch, ErrGroupInsertUnsupported)
	require.True(t, IsSchemaShapeError(physical))
	require.False(t, IsMissingRelationError(physical))

	missing := &mysqldriver.MySQLError{Number: 1146, Message: "table does not exist"}
	require.True(t, IsSchemaShapeError(missing))
	require.True(t, IsMissingRelationError(missing))

	missingColumn := &mysqldriver.MySQLError{Number: 1054, Message: "unknown column"}
	require.True(t, IsSchemaShapeError(missingColumn))
	require.False(t, IsMissingRelationError(missingColumn))

	for _, number := range []uint16{1044, 1045, 1142, 1062, 1205, 1213, 2006, 2013} {
		err := &mysqldriver.MySQLError{Number: number}
		require.False(t, IsSchemaShapeError(err), number)
		require.False(t, IsMissingRelationError(err), number)
	}
	require.False(t, IsSchemaShapeError(context.DeadlineExceeded))
	require.False(t, IsMissingRelationError(context.DeadlineExceeded))
}

func TestOpenMySQLManagedOperationsRequireExpectedBinding(t *testing.T) {
	_, err := InspectOpenMySQLManagedTable(t.Context(), &sql.DB{}, nil, "gw", "t")
	require.ErrorIs(t, err, ErrMySQLManagedServerMismatch)

	err = ExecuteOpenMySQLManagedStatement(t.Context(), &sql.DB{}, &MySQLRecordingCapability{}, "gw", "CREATE TABLE `t` (`x` INT)")
	require.ErrorIs(t, err, ErrMySQLManagedServerMismatch)

	err = ExecuteOpenMySQLManagedStatement(t.Context(), &sql.DB{}, &MySQLRecordingCapability{
		ConnectorID: "c", ConnectorRevision: "r", Database: "gw",
	}, "gw", "ALTER TABLE `t` ADD COLUMN `x` INT")
	require.ErrorIs(t, err, ErrMySQLManagedStatementRefused)
}

func TestOpenMySQLBindingInspectionPreservesForbidden(t *testing.T) {
	actual := &MySQLRecordingCapability{
		ConnectorID: "c", ConnectorRevision: "r", Database: "gw",
		ServerFamily: "mysql", ServerVersion: "8.4.5", CapabilityRevision: "cap",
		Connected: true, MetadataStatus: "verified", TransactionalEngine: true,
	}
	result := openMySQLBindingInspection(actual, "gw", "t", fmt.Errorf("reader: %w", ErrMySQLManagedTableForbidden))
	require.Equal(t, TableInspectionForbidden, result.Inspection.Status)
	require.Equal(t, "permission_denied", result.Inspection.Reason)

	result = openMySQLBindingInspection(actual, "gw", "t", &mysqldriver.MySQLError{Number: 1142})
	require.Equal(t, TableInspectionForbidden, result.Inspection.Status)
	require.Equal(t, "permission_denied", result.Inspection.Reason)
}
