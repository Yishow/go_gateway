package workspace

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type writeGroupReadinessInspector struct {
	inspection *dbtarget.TableInspection
	calls      int
}

func (i *writeGroupReadinessInspector) InspectTable(
	context.Context, string, string, string,
) (*dbtarget.TableInspection, error) {
	i.calls++
	return i.inspection, nil
}

func newReadinessGroup(ctx context.Context, t *testing.T, db *sql.DB) (*WriteGroupService, writeGroupFixture, *WriteGroupSaveResult) {
	t.Helper()
	fixture := seedWriteGroupFixture(ctx, t, db)
	_, err := db.ExecContext(ctx, `
		UPDATE database_connectors
		SET kind = ?, connection_config = ?
		WHERE id = ?
	`, schema.DatabaseConnectorKindPostgres, `{"database":"persisted-db","schema":"persisted-schema"}`, fixture.connectorID)
	require.NoError(t, err)
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	group := newBasicWriteGroup(fixture)
	group.RowPolicy.IntervalSeconds = 15
	created, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: group,
	})
	require.NoError(t, err)
	return service, fixture, created
}

func TestWriteGroupReadinessVerifiesSavedTableWithoutMutation(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture, created := newReadinessGroup(ctx, t, db)
	// The connector identity is deliberately unchanged: this test changes only
	// the saved connector kind/configuration used by the read-only inspector seam.
	inspector := &writeGroupReadinessInspector{inspection: &dbtarget.TableInspection{
		Status: dbtarget.TableInspectionExists, Schema: "persisted-schema", Table: "raw_values",
		Columns: []dbtarget.ColumnInfo{{Name: "temperature", DataType: "REAL"}},
	}}
	service.WithTableInspector(inspector)

	before, err := service.repo.Get(ctx, fixture.workspaceID, created.Group.ID)
	require.NoError(t, err)
	readiness, err := service.Readiness(ctx, created.Group.ID)
	require.NoError(t, err)
	require.True(t, readiness.ConfigReady)
	require.True(t, readiness.SchemaReady)
	require.True(t, readiness.Ready)
	require.NotEmpty(t, readiness.SchemaDigest)
	require.Equal(t, 1, inspector.calls)

	after, err := service.repo.Get(ctx, fixture.workspaceID, created.Group.ID)
	require.NoError(t, err)
	require.Equal(t, before.AppliedRevision, after.AppliedRevision)
	require.Equal(t, before.Destination.SchemaRevision, after.Destination.SchemaRevision)
	require.Empty(t, after.Destination.SchemaDigest)
	require.Equal(t, before.Members, after.Members)
}

func TestWriteGroupReadinessBlocksMissingDestinationColumn(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, _, created := newReadinessGroup(ctx, t, db)
	service.WithTableInspector(&writeGroupReadinessInspector{inspection: &dbtarget.TableInspection{
		Status: dbtarget.TableInspectionExists, Schema: "persisted-schema", Table: "raw_values",
		Columns: []dbtarget.ColumnInfo{{Name: "other", DataType: "REAL"}},
	}})

	readiness, err := service.Readiness(ctx, created.Group.ID)
	require.NoError(t, err)
	require.True(t, readiness.ConfigReady)
	require.False(t, readiness.SchemaReady)
	require.False(t, readiness.Ready)
	require.Contains(t, readinessIssueCodes(readiness.Issues), "destination-column-missing")
}

func TestWriteGroupReadinessRejectsUnknownColumnTypeAndScope(t *testing.T) {
	group := &WriteGroup{
		ID:          "group-1",
		Destination: WriteGroupDestination{TableSchema: "main", TableName: "raw_values"},
		Members:     []WriteGroupMember{{TagID: "tag-1", TargetColumn: "temperature"}},
	}
	ready, issues := evaluateWriteGroupInspection(group, map[string]schema.DataType{
		"tag-1": schema.DataTypeFloat32,
	}, &dbtarget.TableInspection{
		Status:  dbtarget.TableInspectionExists,
		Columns: []dbtarget.ColumnInfo{{Name: "temperature", DataType: "mystery_integer"}},
	})
	require.False(t, ready)
	require.Contains(t, readinessIssueCodes(issues), "schema-scope-mismatch")

	ready, issues = evaluateWriteGroupInspection(group, map[string]schema.DataType{
		"tag-1": schema.DataTypeFloat32,
	}, &dbtarget.TableInspection{
		Status: dbtarget.TableInspectionExists, Schema: "main", Table: "raw_values",
		Columns: []dbtarget.ColumnInfo{{Name: "temperature", DataType: "mystery_integer"}},
	})
	require.False(t, ready)
	require.Contains(t, readinessIssueCodes(issues), "destination-column-type-unverified")
}

func TestWriteGroupReadinessDoesNotCreateMissingSQLiteTarget(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	target := filepath.Join(t.TempDir(), "missing-target.db")
	_, err := db.ExecContext(ctx, `
		UPDATE database_connectors
		SET connection_config = ?
		WHERE id = ?
	`, `{"dsn":"`+target+`","database":"`+target+`","schema":"main"}`, fixture.connectorID)
	require.NoError(t, err)
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	group := newBasicWriteGroup(fixture)
	group.RowPolicy.IntervalSeconds = 15
	created, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: group,
	})
	require.NoError(t, err)
	inspector := &writeGroupReadinessInspector{inspection: &dbtarget.TableInspection{
		Status: dbtarget.TableInspectionExists, Schema: "main", Table: "raw_values",
		Columns: []dbtarget.ColumnInfo{{Name: "temperature", DataType: "REAL"}},
	}}
	service.WithTableInspector(inspector)

	readiness, err := service.Readiness(ctx, created.Group.ID)
	require.NoError(t, err)
	require.False(t, readiness.SchemaReady)
	require.False(t, readiness.Ready)
	require.Contains(t, readinessIssueCodes(readiness.Issues), "schema-target-missing")
	require.Zero(t, inspector.calls)
	require.NoFileExists(t, target)
}

func TestWriteGroupReadinessUsesRealSQLiteReadOnlyInspector(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	target := filepath.Join(t.TempDir(), "existing-target.db")
	targetDB, err := sql.Open("sqlite", target)
	require.NoError(t, err)
	_, err = targetDB.ExecContext(ctx, `CREATE TABLE raw_values (temperature REAL NOT NULL)`)
	require.NoError(t, err)
	require.NoError(t, targetDB.Close())
	_, err = db.ExecContext(ctx, `
		UPDATE database_connectors
		SET connection_config = ?
		WHERE id = ?
	`, `{"dsn":"`+target+`","database":"`+target+`","schema":"main"}`, fixture.connectorID)
	require.NoError(t, err)

	workspaceSvc := NewService(NewSQLRepository(db))
	connectorService := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db)).WithTableInspector(
		dbtarget.NewReadOnlyTableInspector(connectorService),
	)
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	group := newBasicWriteGroup(fixture)
	group.RowPolicy.IntervalSeconds = 15
	created, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: group,
	})
	require.NoError(t, err)
	require.Nil(t, created.Group.Members[0].MeasurementID)
	var manualTargetCount int
	require.NoError(t, db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM database_target_mappings WHERE connector_id = ?`, fixture.connectorID,
	).Scan(&manualTargetCount))
	require.Zero(t, manualTargetCount)

	before, err := service.Get(ctx, created.Group.ID)
	require.NoError(t, err)
	readiness, err := service.Readiness(ctx, created.Group.ID)
	require.NoError(t, err)
	require.True(t, readiness.ConfigReady)
	require.True(t, readiness.SchemaReady)
	require.True(t, readiness.Ready)
	require.NotEmpty(t, readiness.SchemaDigest)

	after, err := service.Get(ctx, created.Group.ID)
	require.NoError(t, err)
	require.Equal(t, before.Group.AppliedRevision, after.Group.AppliedRevision)
	require.Empty(t, after.Group.AppliedRevision)
	require.Equal(t, before.Group.Destination.SchemaRevision, after.Group.Destination.SchemaRevision)
	require.Empty(t, after.Group.Destination.SchemaRevision)
	require.Empty(t, after.Group.Destination.SchemaDigest)
	require.Equal(t, before.Group.Members, after.Group.Members)
}

func TestWriteGroupReadinessBlocksUnsupportedBasicPolicies(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture, created := newReadinessGroup(ctx, t, db)
	current, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	candidate := cloneWriteGroup(created.Group)
	candidate.RowPolicy.IntervalSeconds = 0
	candidate.WritePolicy.Mode = "latest"
	updated, err := service.Update(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID:               fixture.workspaceID,
		ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedGroupRevision:     created.Group.Revision,
		ExpectedConnectorRevision: "connector-revision-1",
		Group:                     candidate,
	})
	require.NoError(t, err)

	readiness, err := service.Readiness(ctx, updated.Group.ID)
	require.NoError(t, err)
	require.False(t, readiness.ConfigReady)
	require.False(t, readiness.Ready)
	require.Contains(t, readinessIssueCodes(readiness.Issues), "interval-required")
	require.Contains(t, readinessIssueCodes(readiness.Issues), "write-policy-blocked")
	require.Empty(t, updated.Group.AppliedRevision)
}

func TestWriteGroupReadinessReportsStaleSourceAndConnector(t *testing.T) {
	for _, test := range []struct {
		name  string
		apply func(context.Context, *sql.DB, writeGroupFixture) error
		code  string
	}{
		{
			name: "source",
			apply: func(ctx context.Context, db *sql.DB, fixture writeGroupFixture) error {
				_, err := db.ExecContext(ctx, `UPDATE tags SET display_name = ? WHERE id = ?`, "Changed", fixture.tagID)
				return err
			},
			code: "source-revision-stale",
		},
		{
			name: "connector",
			apply: func(ctx context.Context, db *sql.DB, fixture writeGroupFixture) error {
				_, err := db.ExecContext(ctx, `UPDATE database_connectors SET identity_revision = ? WHERE id = ?`, "connector-revision-2", fixture.connectorID)
				return err
			},
			code: "connector-revision-stale",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()
			db := openWorkspaceTestDB(t, ":memory:")
			defer db.Close()
			service, fixture, created := newReadinessGroup(ctx, t, db)
			require.NoError(t, test.apply(ctx, db, fixture))

			readiness, err := service.Readiness(ctx, created.Group.ID)
			require.NoError(t, err)
			require.False(t, readiness.ConfigReady)
			require.False(t, readiness.Ready)
			require.Contains(t, readinessIssueCodes(readiness.Issues), test.code)
		})
	}
}

func readinessIssueCodes(issues []ReadinessIssue) []string {
	codes := make([]string, 0, len(issues))
	for _, issue := range issues {
		codes = append(codes, issue.Code)
	}
	return codes
}
