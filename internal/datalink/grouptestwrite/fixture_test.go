package grouptestwrite

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/testutil/lostcommit"

	_ "modernc.org/sqlite"
)

const (
	testWorkspaceID = "ws-1"
	testRevision    = "setup-1"
	testConnector   = "db-1"
	testConnRev     = "identity-1"
	testGroupID     = "group-1"
)

// fakeGroups serves one saved group.
type fakeGroups struct {
	mu    sync.Mutex
	group *workspace.WriteGroup
}

func (f *fakeGroups) Get(_ context.Context, id string) (*workspace.WriteGroupSaveResult, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.group == nil || f.group.ID != id {
		return nil, workspace.ErrWriteGroupNotFound
	}
	clone := *f.group
	return &workspace.WriteGroupSaveResult{WorkspaceRevision: testRevision, Group: &clone}, nil
}

type fakeTags map[string]schema.DataType

func (f fakeTags) GetByID(_ context.Context, id string) (*schema.Tag, error) {
	dataType, ok := f[id]
	if !ok {
		return nil, errors.New("tag not found")
	}
	return &schema.Tag{ID: id, DataType: dataType}, nil
}

// fakeDestinations opens the real SQLite destination file, optionally through
// fault injection, and can be taken offline or made to wait.
type fakeDestinations struct {
	path     string
	faults   *lostcommit.Faults
	revision string
	offline  atomic.Bool
	gate     chan struct{} // when set, OpenDestination waits for it to close
	opened   atomic.Int32
}

func (f *fakeDestinations) GetByID(_ context.Context, id string) (*schema.DatabaseConnector, error) {
	if id != testConnector {
		return nil, dbtarget.ErrConnectorNotFound
	}
	return &schema.DatabaseConnector{ID: id, Kind: schema.DatabaseConnectorKindSQLite, IdentityRevision: f.revision, Enabled: true}, nil
}

func (f *fakeDestinations) OpenDestination(ctx context.Context, connectorID, expectedRevision string) (*dbtarget.OpenedDestination, error) {
	if f.gate != nil {
		select {
		case <-f.gate:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	if f.offline.Load() {
		return nil, dbtarget.ErrDestinationUnreachable
	}
	if connectorID != testConnector || expectedRevision != f.revision {
		return nil, dbtarget.ErrDestinationBlocked
	}
	f.opened.Add(1)
	db := lostcommit.OpenFaults(f.path, f.faults)
	db.SetMaxOpenConns(1)
	return &dbtarget.OpenedDestination{DB: db, Kind: schema.DatabaseConnectorKindSQLite, Close: db.Close}, nil
}

// fakeInspector reads real column metadata from the SQLite file.
type fakeInspector struct{ db *sql.DB }

func (f fakeInspector) InspectTable(ctx context.Context, _, _, table string) (*dbtarget.TableInspection, error) {
	rows, err := f.db.QueryContext(ctx, `SELECT name, type, "notnull" FROM pragma_table_info(?)`, table)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var columns []dbtarget.ColumnInfo
	for rows.Next() {
		var name, dataType string
		var notNull int
		if err := rows.Scan(&name, &dataType, &notNull); err != nil {
			return nil, err
		}
		columns = append(columns, dbtarget.ColumnInfo{Name: name, DataType: dataType, Nullable: notNull == 0})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(columns) == 0 {
		return &dbtarget.TableInspection{Status: dbtarget.TableInspectionMissing, Table: table}, nil
	}
	return &dbtarget.TableInspection{Status: dbtarget.TableInspectionExists, Table: table, Columns: columns}, nil
}

// harness is one gateway with a real SQLite destination holding a production
// neighbor row next to where the test row goes.
type harness struct {
	t        *testing.T
	svc      *Service
	ledger   *recordingplan.Service
	repo     recordingplan.Repository
	groups   *fakeGroups
	dest     *fakeDestinations
	inspect  fakeInspector
	faults   *lostcommit.Faults
	probe    *sql.DB
	path     string
	ws       Workspace
	receipts bool
}

type harnessOptions struct {
	receipt bool
	repo    recordingplan.Repository
}

func newHarness(t *testing.T, opts harnessOptions) *harness {
	t.Helper()
	path := filepath.Join(t.TempDir(), "destination.db")
	probe, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	probe.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = probe.Close() })
	for _, ddl := range []string{
		`CREATE TABLE readings (entity TEXT, temperature REAL, running INTEGER, batch TEXT, big INTEGER, prov TEXT, note TEXT)`,
		`INSERT INTO readings (entity, temperature, running, batch, big, prov, note) VALUES ('line-a', 20.5, 0, 'B-1', 7, '[]', 'production')`,
	} {
		if _, err := probe.ExecContext(t.Context(), ddl); err != nil {
			t.Fatal(err)
		}
	}
	if opts.receipt {
		if err := dbtarget.CreateEffectReceiptTable(t.Context(), probe, schema.DatabaseConnectorKindSQLite, ""); err != nil {
			t.Fatal(err)
		}
	}
	return attachHarness(t, path, probe, opts)
}

// attachHarness builds a gateway around an existing destination file; another
// process can attach to the same files.
func attachHarness(t *testing.T, path string, probe *sql.DB, opts harnessOptions) *harness {
	t.Helper()
	repo := opts.repo
	if repo == nil {
		repo = recordingplan.NewMemoryRepository()
	}
	faults := &lostcommit.Faults{Lose: &atomic.Bool{}}
	h := &harness{
		t: t, repo: repo, ledger: recordingplan.NewService(repo), faults: faults, probe: probe, path: path,
		ws: Workspace{ID: testWorkspaceID, Revision: testRevision}, receipts: opts.receipt,
		dest:    &fakeDestinations{path: path, faults: faults, revision: testConnRev},
		inspect: fakeInspector{db: probe},
	}
	dedupe := ""
	if opts.receipt {
		dedupe = "receipt"
	}
	h.groups = &fakeGroups{group: &workspace.WriteGroup{
		ID: testGroupID, WorkspaceID: testWorkspaceID, Revision: "rev-1", Status: workspace.WriteGroupStatusReady,
		Members: []workspace.WriteGroupMember{
			{DeviceID: "d1", PointID: "p1", TagID: "t-temp", TargetColumn: "temperature", Required: true},
			{DeviceID: "d1", PointID: "p2", TagID: "t-run", TargetColumn: "running", Required: true},
			{DeviceID: "d1", PointID: "p3", TagID: "t-batch", TargetColumn: "batch", Required: true},
			{DeviceID: "d1", PointID: "p4", TagID: "t-big", TargetColumn: "big", Required: true},
		},
		Destination: workspace.WriteGroupDestination{
			ConnectorID: testConnector, ConnectorRevision: testConnRev, Database: path, TableSchema: "main", TableName: "readings",
			StorageStrategy: workspace.WriteGroupStorageStrategyCustom,
		},
		RowPolicy:   workspace.WriteGroupRowPolicy{IntervalSeconds: 10, EntityKeyColumn: "entity", ProvenanceColumn: "prov"},
		WritePolicy: workspace.WriteGroupWritePolicy{DedupeCapability: dedupe},
	}}
	tags := fakeTags{"t-temp": schema.DataTypeFloat64, "t-run": schema.DataTypeBool, "t-batch": schema.DataTypeString, "t-big": schema.DataTypeUint64}
	h.svc = New(Dependencies{Groups: h.groups, Tags: tags, Destinations: h.dest, Inspector: h.inspect, Ledger: h.ledger}, Config{})
	return h
}

func (h *harness) count(where string) int {
	h.t.Helper()
	var n int
	if err := h.probe.QueryRowContext(h.t.Context(), `SELECT COUNT(*) FROM readings WHERE `+where).Scan(&n); err != nil {
		h.t.Fatal(err)
	}
	return n
}

func (h *harness) productionRow() (note string, temperature float64) {
	h.t.Helper()
	if err := h.probe.QueryRowContext(h.t.Context(), `SELECT note, temperature FROM readings WHERE entity = 'line-a'`).Scan(&note, &temperature); err != nil {
		h.t.Fatal(err)
	}
	return note, temperature
}

func (h *harness) preview() *Preview {
	h.t.Helper()
	preview, err := h.svc.Preview(h.t.Context(), h.ws, testGroupID)
	if err != nil {
		h.t.Fatalf("preview: %v", err)
	}
	return preview
}

func (h *harness) confirm(preview *Preview) (*Outcome, error) {
	return h.svc.Confirm(h.t.Context(), h.ws, Confirmation{Token: preview.Token, OperationID: preview.OperationID})
}
