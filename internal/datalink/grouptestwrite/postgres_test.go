package grouptestwrite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/testutil/lostcommit"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

// Live PostgreSQL test-write fixtures; skipped without POSTGRES_DSN. Each test
// uses its own schema and roles, dropped afterwards.

type pgDestinations struct {
	dsn      string
	revision string
	lose     *atomic.Bool
}

func (f *pgDestinations) GetByID(_ context.Context, id string) (*schema.DatabaseConnector, error) {
	return &schema.DatabaseConnector{ID: id, Kind: schema.DatabaseConnectorKindPostgres, IdentityRevision: f.revision, Enabled: true}, nil
}

func (f *pgDestinations) OpenDestination(_ context.Context, _, _ string) (*dbtarget.OpenedDestination, error) {
	db := lostcommit.Wrap(stdlib.GetDefaultDriver(), f.dsn, f.lose)
	return &dbtarget.OpenedDestination{DB: db, Kind: schema.DatabaseConnectorKindPostgres, Close: db.Close}, nil
}

type pgInspector struct {
	db         *sql.DB
	schemaName string
}

func (f pgInspector) InspectTable(ctx context.Context, _, _, table string) (*dbtarget.TableInspection, error) {
	rows, err := f.db.QueryContext(ctx, `SELECT column_name, data_type, is_nullable FROM information_schema.columns
		WHERE table_schema = $1 AND table_name = $2 ORDER BY ordinal_position`, f.schemaName, table)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var columns []dbtarget.ColumnInfo
	for rows.Next() {
		var name, dataType, nullable string
		if err := rows.Scan(&name, &dataType, &nullable); err != nil {
			return nil, err
		}
		columns = append(columns, dbtarget.ColumnInfo{Name: name, DataType: dataType, Nullable: nullable == "YES"})
	}
	if len(columns) == 0 {
		return &dbtarget.TableInspection{Status: dbtarget.TableInspectionMissing, Table: table}, nil
	}
	return &dbtarget.TableInspection{Status: dbtarget.TableInspectionExists, Table: table, Columns: columns}, rows.Err()
}

type pgHarness struct {
	t          *testing.T
	svc        *Service
	ws         Workspace
	admin      *sql.DB
	schemaName string
	dest       *pgDestinations
	groups     *fakeGroups
	ledger     *recordingplan.Service
}

func pgDSN(t *testing.T) string {
	t.Helper()
	dsn := strings.TrimSpace(os.Getenv("POSTGRES_DSN"))
	if dsn == "" {
		t.Skip("Skipping PostgreSQL test-write test: POSTGRES_DSN not set")
	}
	return dsn
}

// dsnAs swaps the user and password of a key/value DSN.
func dsnAs(dsn, user, password string) string {
	var kept []string
	for _, part := range strings.Fields(dsn) {
		if !strings.HasPrefix(part, "user=") && !strings.HasPrefix(part, "password=") {
			kept = append(kept, part)
		}
	}
	return strings.Join(append(kept, "user="+user, "password="+password), " ")
}

func newPGHarness(t *testing.T, receipt bool) *pgHarness {
	t.Helper()
	dsn := pgDSN(t)
	admin, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = admin.Close() })
	require.NoError(t, admin.PingContext(t.Context()))

	schemaName := fmt.Sprintf("gw_tw_%d", time.Now().UnixNano())
	_, err = admin.ExecContext(t.Context(), `CREATE SCHEMA "`+schemaName+`"`)
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = admin.ExecContext(context.Background(), `DROP SCHEMA "`+schemaName+`" CASCADE`) })
	_, err = admin.ExecContext(t.Context(), `CREATE TABLE "`+schemaName+`".readings (
		entity TEXT, temperature DOUBLE PRECISION, running BOOLEAN, batch TEXT, big BIGINT, prov JSONB, note TEXT)`)
	require.NoError(t, err)
	_, err = admin.ExecContext(t.Context(), `INSERT INTO "`+schemaName+`".readings (entity, temperature, running, batch, big, prov, note)
		VALUES ('line-a', 20.5, false, 'B-1', 7, '[]', 'production')`)
	require.NoError(t, err)
	dedupe := ""
	if receipt {
		require.NoError(t, dbtarget.CreateEffectReceiptTable(t.Context(), admin, schema.DatabaseConnectorKindPostgres, schemaName))
		dedupe = "receipt"
	}
	h := &pgHarness{t: t, admin: admin, schemaName: schemaName, ws: Workspace{ID: testWorkspaceID, Revision: testRevision}}
	h.dest = &pgDestinations{dsn: dsn, revision: testConnRev, lose: &atomic.Bool{}}
	h.ledger = recordingplan.NewService(recordingplan.NewMemoryRepository())
	h.groups = &fakeGroups{group: &workspace.WriteGroup{
		ID: testGroupID, WorkspaceID: testWorkspaceID, Revision: "rev-1", Status: workspace.WriteGroupStatusReady,
		Members: []workspace.WriteGroupMember{
			{DeviceID: "d1", PointID: "p1", TagID: "t-temp", TargetColumn: "temperature", Required: true},
			{DeviceID: "d1", PointID: "p2", TagID: "t-run", TargetColumn: "running", Required: true},
			{DeviceID: "d1", PointID: "p3", TagID: "t-batch", TargetColumn: "batch", Required: true},
			{DeviceID: "d1", PointID: "p4", TagID: "t-big", TargetColumn: "big", Required: true},
		},
		Destination: workspace.WriteGroupDestination{
			ConnectorID: testConnector, ConnectorRevision: testConnRev, Database: "gwtest", TableSchema: schemaName, TableName: "readings",
			StorageStrategy: workspace.WriteGroupStorageStrategyCustom,
		},
		RowPolicy:   workspace.WriteGroupRowPolicy{IntervalSeconds: 10, EntityKeyColumn: "entity", ProvenanceColumn: "prov"},
		WritePolicy: workspace.WriteGroupWritePolicy{DedupeCapability: dedupe},
	}}
	tags := fakeTags{"t-temp": schema.DataTypeFloat64, "t-run": schema.DataTypeBool, "t-batch": schema.DataTypeString, "t-big": schema.DataTypeUint64}
	h.svc = New(Dependencies{Groups: h.groups, Tags: tags, Destinations: h.dest, Inspector: pgInspector{db: admin, schemaName: schemaName}, Ledger: h.ledger}, Config{})
	return h
}

func (h *pgHarness) count(where string) int {
	h.t.Helper()
	var n int
	require.NoError(h.t, h.admin.QueryRowContext(h.t.Context(), `SELECT COUNT(*) FROM "`+h.schemaName+`".readings WHERE `+where).Scan(&n))
	return n
}

func (h *pgHarness) confirmPreview() (*Preview, *Outcome) {
	h.t.Helper()
	preview, err := h.svc.Preview(h.t.Context(), h.ws, testGroupID)
	require.NoError(h.t, err)
	outcome, err := h.svc.Confirm(h.t.Context(), h.ws, Confirmation{Token: preview.Token, OperationID: preview.OperationID})
	require.NoError(h.t, err)
	return preview, outcome
}

func TestPostgresTestWriteVerifiesTypedValuesAndCleansOnlyItsOwnRow(t *testing.T) {
	for name, receipt := range map[string]bool{"no dedupe": false, "receipt dedupe": true} {
		t.Run(name, func(t *testing.T) {
			h := newPGHarness(t, receipt)
			_, outcome := h.confirmPreview()
			op := outcome.Operation
			require.Equal(t, WriteVerified, op.WriteOutcome, "%+v", op)
			require.Equal(t, CleanupCleaned, op.CleanupStatus, "%+v", op)
			require.Equal(t, 1, h.count(`TRUE`), "only the production row remains")
			require.Equal(t, 1, h.count(`entity = 'line-a' AND note = 'production'`))
			if receipt {
				var receipts int
				require.NoError(t, h.admin.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM "`+h.schemaName+`".gw_effect_receipts`).Scan(&receipts))
				require.Zero(t, receipts)
			}
		})
	}
}

func TestPostgresTestWriteSeparatesMissingReadAndDeletePermissions(t *testing.T) {
	h := newPGHarness(t, false)
	role := fmt.Sprintf("gw_tw_ro_%d", time.Now().UnixNano())
	_, err := h.admin.ExecContext(t.Context(), `CREATE ROLE "`+role+`" LOGIN PASSWORD 'pw'`)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = h.admin.ExecContext(context.Background(), `DROP OWNED BY "`+role+`"`)
		_, _ = h.admin.ExecContext(context.Background(), `DROP ROLE IF EXISTS "`+role+`"`)
	})
	for _, grant := range []string{
		`GRANT USAGE ON SCHEMA "` + h.schemaName + `" TO "` + role + `"`,
		`GRANT INSERT ON "` + h.schemaName + `".readings TO "` + role + `"`,
	} {
		_, err := h.admin.ExecContext(t.Context(), grant)
		require.NoError(t, err)
	}
	h.dest.dsn = dsnAs(pgDSN(t), role, "pw")

	_, outcome := h.confirmPreview()
	op := outcome.Operation
	require.Equal(t, WriteUnverified, op.WriteOutcome, "a committed write that cannot be read back is not verified: %+v", op)
	require.Equal(t, reasonReadbackDenied, op.Reason)
	require.Equal(t, CleanupFailed, op.CleanupStatus, "a role that cannot read or delete cannot clean: %+v", op)
	require.Equal(t, reasonCleanupDenied, op.CleanupReason)
	require.Equal(t, 1, h.count(`entity LIKE 'gw-test-%'`), "the row still exists and is reported as such")
	require.Equal(t, 1, h.count(`entity = 'line-a'`))
}

func TestPostgresTestWriteLostCommitResponseWritesOnce(t *testing.T) {
	h := newPGHarness(t, true)
	h.dest.lose.Store(true)
	preview, err := h.svc.Preview(t.Context(), h.ws, testGroupID)
	require.NoError(t, err)
	outcome, err := h.svc.Confirm(t.Context(), h.ws, Confirmation{Token: preview.Token, OperationID: preview.OperationID})
	require.NoError(t, err)
	op := outcome.Operation
	// Both commits lose their response: the write is settled by readback and the
	// receipt, the cleanup stays honestly unknown.
	require.Equal(t, WriteVerified, op.WriteOutcome, "%+v", op)
	require.Equal(t, CleanupUnknown, op.CleanupStatus, "%+v", op)
	require.Zero(t, h.count(`entity LIKE 'gw-test-%'`), "the cleanup did commit although its response was lost")
}
