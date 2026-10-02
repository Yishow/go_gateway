package workspace

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/dbtarget"

	"github.com/stretchr/testify/require"
)

type mutatingWriteGroupReadinessInspector struct {
	db         *sql.DB
	mutate     func(context.Context, *sql.DB) error
	inspection *dbtarget.TableInspection
}

func (i *mutatingWriteGroupReadinessInspector) InspectTable(ctx context.Context, _, _, _ string) (*dbtarget.TableInspection, error) {
	if i.mutate != nil {
		mutate := i.mutate
		i.mutate = nil
		if err := mutate(ctx, i.db); err != nil {
			return nil, err
		}
	}
	return i.inspection, nil
}

func TestWriteGroupReadinessRejectsLiveRevisionChangesDuringInspection(t *testing.T) {
	for _, test := range []struct {
		name      string
		mutate    func(context.Context, *sql.DB, writeGroupFixture) error
		issueCode string
	}{
		{
			name: "connector identity revision",
			mutate: func(ctx context.Context, db *sql.DB, fixture writeGroupFixture) error {
				_, err := db.ExecContext(ctx, `UPDATE database_connectors SET identity_revision = ? WHERE id = ?`, "connector-revision-2", fixture.connectorID)
				return err
			},
			issueCode: "connector-revision-stale",
		},
		{
			name: "source metadata",
			mutate: func(ctx context.Context, db *sql.DB, fixture writeGroupFixture) error {
				_, err := db.ExecContext(ctx, `UPDATE tags SET unit = ? WHERE id = ?`, "changed-unit", fixture.tagID)
				return err
			},
			issueCode: "source-revision-stale",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()
			dbPath := filepath.Join(t.TempDir(), "readiness-snapshot.db")
			db := openWorkspaceTestDB(t, dbPath)
			defer db.Close()
			_, err := db.ExecContext(ctx, `PRAGMA journal_mode = WAL`)
			require.NoError(t, err)
			service, fixture, created := newReadinessGroup(ctx, t, db)
			beforeGroup, err := service.Get(ctx, created.Group.ID)
			require.NoError(t, err)
			beforeWorkspace, err := service.workspaceSvc.GetOrCreate(ctx)
			require.NoError(t, err)

			inspector := &mutatingWriteGroupReadinessInspector{
				db: db,
				mutate: func(ctx context.Context, db *sql.DB) error {
					return test.mutate(ctx, db, fixture)
				},
				inspection: &dbtarget.TableInspection{
					Status: dbtarget.TableInspectionExists, Schema: "persisted-schema", Table: "raw_values",
					Columns: []dbtarget.ColumnInfo{{Name: "temperature", DataType: "REAL"}},
				},
			}
			service.WithTableInspector(inspector)

			readiness, err := service.Readiness(ctx, created.Group.ID)
			require.NoError(t, err)
			require.False(t, readiness.ConfigReady)
			require.False(t, readiness.SchemaReady)
			require.False(t, readiness.Ready)
			require.Contains(t, readinessIssueCodes(readiness.Issues), test.issueCode)

			afterGroup, err := service.Get(ctx, created.Group.ID)
			require.NoError(t, err)
			afterWorkspace, err := service.workspaceSvc.GetOrCreate(ctx)
			require.NoError(t, err)
			require.Equal(t, beforeGroup, afterGroup)
			require.Equal(t, beforeWorkspace, afterWorkspace)
		})
	}
}
