package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

func TestSuccessorSchemaTargetRequiresPrivateAuthority(t *testing.T) {
	for _, mode := range []string{"plain", "mysql", "sqlserver", "timescale"} {
		t.Run(mode, func(t *testing.T) {
			scope := recordingplan.SchemaPreviewScope{GroupLayout: &recordingplan.GroupSchemaLayout{}}
			if mode == "mysql" {
				scope.GroupLayout.MySQL = &recordingplan.MySQLGroupSchemaOptions{}
			}
			if mode == "sqlserver" {
				scope.GroupLayout.SQLServer = &recordingplan.SQLServerGroupSchemaOptions{}
			}
			target := &checkedGroupSchemaTarget{inspector: dbtarget.NewManagedTableInspector(nil, nil)}
			_, err := target.targetInspector(scope)(t.Context(), "table")
			require.ErrorIs(t, err, recordingplan.ErrPreviewTokenStale)
			_, err = target.executeGroupSchemaStatements(t.Context(), scope, []string{"CREATE TABLE never_run (id INT)"})
			require.ErrorIs(t, err, recordingplan.ErrPreviewTokenStale)
			if mode == "sqlserver" {
				_, err = target.executeSQLServerGroupSchemaStatements(t.Context(), scope, nil, nil)
				require.ErrorIs(t, err, recordingplan.ErrPreviewTokenStale)
			}
			if mode == "timescale" {
				_, err = target.timescaleEvidence(t.Context(), scope, false)
				require.ErrorIs(t, err, recordingplan.ErrPreviewTokenStale)
			}
			sentinel := errors.New("authority rejected")
			calls := 0
			target.revalidate = func(context.Context) error { calls++; return sentinel }
			_, err = target.targetInspector(scope)(t.Context(), "table")
			require.ErrorIs(t, err, sentinel)
			require.Equal(t, 1, calls)
		})
	}
}

func TestOrdinarySchemaDelegatesKeepTheirOwnGuard(t *testing.T) {
	h := NewStudioV2WorkspaceWriteGroupSchemaHandler(nil, nil, nil, nil)
	h.inspector = dbtarget.NewManagedTableInspector(nil, nil)
	scope := recordingplan.SchemaPreviewScope{GroupLayout: &recordingplan.GroupSchemaLayout{}}
	for _, mode := range []string{"plain", "mysql", "sqlserver"} {
		t.Run(mode, func(t *testing.T) {
			_, err := h.targetInspector("g", WriteGroupSchemaRequest{}, scope)(t.Context(), "never_probe")
			require.ErrorIs(t, err, workspace.ErrWriteGroupServiceUnavailable)
			if mode == "mysql" {
				_, err = h.mysqlGroupTargetInspector("g", WriteGroupSchemaRequest{}, scope)(t.Context(), "never_probe")
			}
			if mode == "sqlserver" {
				_, err = h.sqlServerGroupTargetInspector("g", WriteGroupSchemaRequest{}, scope)(t.Context(), "never_probe")
			}
			require.ErrorIs(t, err, workspace.ErrWriteGroupServiceUnavailable)
			_, err = h.executeGroupSchemaStatements(t.Context(), "g", WriteGroupSchemaRequest{}, scope, nil)
			require.ErrorIs(t, err, workspace.ErrWriteGroupServiceUnavailable)
			_, err = h.executeSQLServerGroupSchemaStatements(t.Context(), "g", WriteGroupSchemaRequest{}, scope, nil, nil)
			require.ErrorIs(t, err, workspace.ErrWriteGroupServiceUnavailable)
		})
	}
}
