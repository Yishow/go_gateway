package workspace

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"go-gateway/internal/datalink/sqltx"
	"path/filepath"
	"testing"

	"modernc.org/sqlite"
)

type replayCleanupConnector struct{ path string }

func (c replayCleanupConnector) Driver() driver.Driver { return &sqlite.Driver{} }
func (c replayCleanupConnector) Connect(context.Context) (driver.Conn, error) {
	base, err := c.Driver().Open(c.path)
	if err != nil {
		return nil, err
	}
	return replayCleanupConn{base}, nil
}

type replayCleanupConn struct{ driver.Conn }

func (c replayCleanupConn) Begin() (driver.Tx, error) {
	tx, err := c.Conn.Begin() //nolint:staticcheck // Exercise database/sql legacy driver fallback with injected terminal rollback failure.
	if err != nil {
		return nil, err
	}
	return replayCleanupTx{tx}, nil
}
func (c replayCleanupConn) ExecContext(ctx context.Context, q string, args []driver.NamedValue) (driver.Result, error) {
	if q == "ROLLBACK" {
		return nil, errors.New("test injected cleanup I/O failure")
	}
	return c.Conn.(driver.ExecerContext).ExecContext(ctx, q, args)
}

type replayCleanupTx struct{ driver.Tx }

func (t replayCleanupTx) Rollback() error { return errors.New("test injected rollback I/O failure") }
func TestCanonicalReplayCleanupFailureIsNotAbsence(t *testing.T) {
	for _, kind := range []string{"mixed", "circuit", "basic"} {
		t.Run(kind, func(t *testing.T) {
			h := newCircuitProfileHarness(t, filepath.Join(t.TempDir(), "replay.db"))
			var path string
			if err := h.db.QueryRowContext(t.Context(), "SELECT file FROM pragma_database_list WHERE name='main'").Scan(&path); err != nil {
				t.Fatal(err)
			}
			db := sql.OpenDB(replayCleanupConnector{path})
			defer db.Close()
			db.SetMaxOpenConns(1)
			h.service.repo = NewSQLWriteGroupRepository(db)
			candidate := mixedDraft(h, WriteGroupMemberCalculationLatest, 5, "candidate")
			candidate.WorkspaceID = h.fixture.workspaceID
			if kind == "circuit" {
				candidate = h.draft(WriteGroupProfileInstantaneous, "circuit", "Circuit", h.primary())
				candidate.WorkspaceID = h.fixture.workspaceID
			}
			var replay func() (*WriteGroupSaveResult, bool, error)
			var initial error
			if kind == "basic" {
				candidate = newBasicManagedDraft(h.fixture)
				_, initial = h.service.repo.findBasicManaged(t.Context(), h.fixture.workspaceID, h.fixture.deviceID)
				replay = func() (*WriteGroupSaveResult, bool, error) {
					return h.service.replayBasicManaged(t.Context(), h.fixture.workspaceID, h.fixture.deviceID, candidate)
				}
			} else {
				key, err := circuitProfileKeyOf(candidate)
				if err != nil {
					t.Fatal(err)
				}
				_, initial = h.service.repo.findProfileGroup(t.Context(), key)
				replay = func() (*WriteGroupSaveResult, bool, error) {
					return h.service.replayCircuitProfile(t.Context(), key, candidate)
				}
			}
			if !errors.Is(initial, sqltx.ErrConnectionUnusable) || !errors.Is(initial, ErrWriteGroupNotFound) {
				t.Fatalf("precondition: %v", initial)
			}
			result, found, err := replay()
			if result != nil || found || !errors.Is(err, sqltx.ErrConnectionUnusable) {
				t.Fatalf("cleanup failure swallowed: result=%v found=%v err=%v", result, found, err)
			}
			mutation := h.mutation(t.Context(), candidate)
			switch kind {
			case "basic":
				_, err = h.service.EnsureBasicManaged(t.Context(), h.fixture.deviceID, mutation)
			case "circuit":
				_, err = h.service.EnsureCircuitProfile(t.Context(), mutation)
			default:
				_, err = h.service.Create(t.Context(), mutation)
			}
			if !errors.Is(err, sqltx.ErrConnectionUnusable) {
				t.Fatalf("create path continued after cleanup failure: %v", err)
			}
			var count int
			if err := h.db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM write_groups").Scan(&count); err != nil {
				t.Fatal(err)
			}
			if count != 0 {
				t.Fatal("cleanup failure created a group")
			}
		})
	}
}
