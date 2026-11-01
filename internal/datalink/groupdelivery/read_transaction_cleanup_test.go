package groupdelivery

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"testing"

	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

var errReadCleanupQuery = errors.New("read fixture query failure")
var errReadCleanupRollback = errors.New("rollback failed at /private/fixture.db: SELECT secret_password FROM private_tokens")

// This driver performs no SQL or external I/O. It exposes the existing
// database/sql transaction error boundary without changing the Store seam.
type readCleanupDriver struct{ state *readCleanupState }
type readCleanupState struct {
	empty                               bool
	rollbackErr                         error
	begins, queries, commits, rollbacks int
}

func (d readCleanupDriver) Open(string) (driver.Conn, error) { return readCleanupConn(d), nil }

type readCleanupConnector struct{ state *readCleanupState }

func (c readCleanupConnector) Connect(context.Context) (driver.Conn, error) {
	return readCleanupConn(c), nil
}
func (c readCleanupConnector) Driver() driver.Driver { return readCleanupDriver(c) }

type readCleanupConn struct{ state *readCleanupState }

func (c readCleanupConn) Prepare(string) (driver.Stmt, error) { return nil, errReadCleanupQuery }
func (c readCleanupConn) Close() error                        { return nil }
func (c readCleanupConn) Begin() (driver.Tx, error) {
	c.state.begins++
	return readCleanupTx(c), nil
}
func (c readCleanupConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return c.Begin()
}
func (c readCleanupConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	c.state.queries++
	if c.state.empty {
		return readCleanupRows{}, nil
	}
	return nil, errReadCleanupQuery
}

type readCleanupTx struct{ state *readCleanupState }

func (tx readCleanupTx) Commit() error   { tx.state.commits++; return nil }
func (tx readCleanupTx) Rollback() error { tx.state.rollbacks++; return tx.state.rollbackErr }

func TestReadTransactionCleanupPreservesRefusalAndRollbackFailure(t *testing.T) {
	cases := []struct {
		name    string
		storage bool
		read    func(*testing.T, *Store) error
	}{
		{"ownership", false, func(t *testing.T, s *Store) error {
			got, err := s.OwnershipView(t.Context(), testKey)
			require.Equal(t, OwnershipView{}, got)
			return err
		}},
		{"interest", false, func(t *testing.T, s *Store) error {
			got, err := s.HasCutoverSource(t.Context(), OwnershipSource{WorkspaceID: testKey.WorkspaceID, DeviceID: "d", PointID: "p", TagID: "t"})
			require.False(t, got)
			return err
		}},
		{"resolve", false, func(t *testing.T, s *Store) error {
			got, err := s.ResolveCutoverOwnership(t.Context(), OwnershipSource{WorkspaceID: testKey.WorkspaceID, DeviceID: "d", PointID: "p", TagID: "t"}, at(1))
			require.Equal(t, CutoverOwnership{}, got)
			return err
		}},
		{"journals", true, func(t *testing.T, s *Store) error {
			got, err := s.OpenJournals(t.Context())
			require.Nil(t, got)
			return err
		}},
		{"scope", true, func(t *testing.T, s *Store) error {
			return validateHandoffScope(t.Context(), s.DB(), testKey.GroupID, testKey.GroupRevision)
		}},
		{"recognize", true, func(t *testing.T, s *Store) error {
			got, err := s.RecognizeTransferredSample(t.Context(), testKey, snapshot.Sample{SampleID: "sample"})
			require.Equal(t, HandoffRecognition{}, got)
			return err
		}},
	}
	for _, tc := range cases {
		for _, cleanupFails := range []bool{false, true} {
			name := tc.name + "/cleanup-ok"
			if cleanupFails {
				name = tc.name + "/cleanup-failed"
			}
			t.Run(name, func(t *testing.T) {
				state := &readCleanupState{}
				if cleanupFails {
					state.rollbackErr = errReadCleanupRollback
				}
				db := sql.OpenDB(readCleanupConnector{state})
				t.Cleanup(func() { require.NoError(t, db.Close()) })
				err := tc.read(t, NewStore(db))
				require.ErrorIs(t, err, ErrOwnershipUnavailable, "original typed refusal remains reachable")
				if tc.storage {
					var storage *handoffStorageError
					require.ErrorAs(t, err, &storage)
					require.ErrorIs(t, err, errReadCleanupQuery)
				}
				if cleanupFails {
					require.ErrorIs(t, err, errReadCleanupRollback)
					require.ErrorIs(t, err, ErrOwnershipUnavailable)
					var cleanupStorage *handoffStorageError
					require.ErrorAs(t, err, &cleanupStorage)
					require.NotContains(t, err.Error(), "/private/fixture.db")
					require.NotContains(t, err.Error(), "SELECT secret_password")
					require.NotContains(t, err.Error(), "private_tokens")
				} else {
					require.NotErrorIs(t, err, errReadCleanupRollback)
				}
				require.Equal(t, 1, state.begins)
				require.Equal(t, 1, state.queries)
				require.Zero(t, state.commits)
				require.Equal(t, 1, state.rollbacks)
			})
		}
	}
}

type readCleanupRows struct{}

func (readCleanupRows) Columns() []string         { return []string{"unused-empty-column"} }
func (readCleanupRows) Close() error              { return nil }
func (readCleanupRows) Next([]driver.Value) error { return io.EOF }

func TestReadTransactionMissRequiresSuccessfulCleanup(t *testing.T) {
	for _, fail := range []bool{false, true} {
		name := "clean"
		if fail {
			name = "cleanup-failed"
		}
		t.Run(name, func(t *testing.T) {
			state := &readCleanupState{empty: true}
			if fail {
				state.rollbackErr = errReadCleanupRollback
			}
			db := sql.OpenDB(readCleanupConnector{state})
			t.Cleanup(func() { require.NoError(t, db.Close()) })
			got, err := NewStore(db).RecognizeTransferredSample(t.Context(), testKey, snapshot.Sample{SampleID: "not-present"})
			require.Equal(t, HandoffRecognition{}, got)
			if fail {
				require.ErrorIs(t, err, errReadCleanupRollback)
				require.ErrorIs(t, err, ErrOwnershipUnavailable)
				var cleanupStorage *handoffStorageError
				require.ErrorAs(t, err, &cleanupStorage)
				require.NotContains(t, err.Error(), "/private/fixture.db")
				require.NotContains(t, err.Error(), "SELECT secret_password")
				require.NotContains(t, err.Error(), "private_tokens")
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, 1, state.begins)
			require.Equal(t, 1, state.queries)
			require.Zero(t, state.commits, "verified miss retains original rollback boundary")
			require.Equal(t, 1, state.rollbacks)
		})
	}
}
