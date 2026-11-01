package point

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type terminalGuardConnector struct {
	receiptErr, rollbackErr error
	receipts, rollbacks     int
}

func (c *terminalGuardConnector) Connect(context.Context) (driver.Conn, error) {
	return &terminalGuardConn{owner: c}, nil
}
func (c *terminalGuardConnector) Driver() driver.Driver { return terminalGuardDriver{owner: c} }

type terminalGuardDriver struct{ owner *terminalGuardConnector }

func (d terminalGuardDriver) Open(string) (driver.Conn, error) {
	return &terminalGuardConn{owner: d.owner}, nil
}

type terminalGuardConn struct {
	driver.Conn
	owner *terminalGuardConnector
}

func (*terminalGuardConn) Close() error { return nil }
func (c *terminalGuardConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return terminalGuardTx{owner: c.owner}, nil
}
func (c *terminalGuardConn) QueryContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	switch {
	case strings.Contains(query, "key='studio_v2_workspace'"):
		return &terminalGuardRows{values: []driver.Value{`{"id":"workspace","ordered_device_ids":["device"]}`}}, nil
	case strings.Contains(query, "FROM points WHERE id=?"):
		return &terminalGuardRows{values: []driver.Value{"point", "device", "point", nil, "HR0", nil, "uint32", "ABCD", "override", "ABCD", int64(schema.PointFormatVersion), "point-revision", "read", nil, nil, nil, nil, true, "2026-11-01T00:00:00Z", "2026-11-01T00:00:00Z"}}, nil
	case strings.Contains(query, "SELECT protocol,connection_config"):
		return &terminalGuardRows{values: []driver.Value{"modbus_tcp", `{}`}}, nil
	case strings.Contains(query, "FROM source_rule_links"):
		return &terminalGuardRows{}, nil
	case strings.Contains(query, "SELECT value FROM system_settings WHERE key=?"):
		c.owner.receipts++
		if c.owner.receiptErr != nil {
			return nil, c.owner.receiptErr
		}
		return &terminalGuardRows{}, nil
	default:
		return nil, errors.New("unexpected terminal guard fixture query")
	}
}

type terminalGuardTx struct{ owner *terminalGuardConnector }

func (terminalGuardTx) Commit() error      { return nil }
func (tx terminalGuardTx) Rollback() error { tx.owner.rollbacks++; return tx.owner.rollbackErr }

type terminalGuardRows struct {
	values []driver.Value
	done   bool
}

func (r *terminalGuardRows) Columns() []string { return make([]string, len(r.values)) }
func (*terminalGuardRows) Close() error        { return nil }
func (r *terminalGuardRows) Next(dest []driver.Value) error {
	if r.done || len(r.values) == 0 {
		return io.EOF
	}
	r.done = true
	copy(dest, r.values)
	return nil
}
func TestFormatInspectForeignRollbackOnlyBareTerminalIsBenign(t *testing.T) {
	cause := errors.New("private cleanup failure")
	for _, tc := range []struct {
		name     string
		rollback error
		benign   bool
	}{
		{"nil", nil, true}, {"bare", sql.ErrTxDone, true},
		{"wrapped", fmt.Errorf("private terminal wrapper: %w", sql.ErrTxDone), false},
		{"joined-cause", errors.Join(sql.ErrTxDone, cause), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			owner := &terminalGuardConnector{rollbackErr: tc.rollback}
			db := sql.OpenDB(owner)
			t.Cleanup(func() { require.NoError(t, db.Close()) })
			result, err := NewSQLRepository(db).InspectFormatMutation(t.Context(), FormatScope{WorkspaceID: "workspace", DeviceID: "device", PointID: "point", RequestID: "request"})
			require.Nil(t, result)
			require.Equal(t, 1, owner.receipts)
			require.Equal(t, 1, owner.rollbacks)
			if tc.benign {
				require.ErrorIs(t, err, ErrPointNotFound)
			} else {
				require.ErrorIs(t, err, ErrFormatUnavailable)
				require.ErrorIs(t, err, tc.rollback)
				require.NotContains(t, err.Error(), "private")
			}
		})
	}
}
func TestFormatInspectOnlyBareMissingReceiptIsNotFound(t *testing.T) {
	cause := errors.New("private receipt read failure")
	for _, tc := range []struct {
		name    string
		readErr error
		missing bool
	}{
		{"empty-rows", nil, true}, {"bare", sql.ErrNoRows, true},
		{"wrapped", fmt.Errorf("private missing wrapper: %w", sql.ErrNoRows), false},
		{"joined-cause", errors.Join(sql.ErrNoRows, cause), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			owner := &terminalGuardConnector{receiptErr: tc.readErr}
			db := sql.OpenDB(owner)
			t.Cleanup(func() { require.NoError(t, db.Close()) })
			result, err := NewSQLRepository(db).InspectFormatMutation(t.Context(), FormatScope{WorkspaceID: "workspace", DeviceID: "device", PointID: "point", RequestID: "request"})
			require.Nil(t, result)
			require.Equal(t, 1, owner.receipts)
			require.Equal(t, 1, owner.rollbacks)
			if tc.missing {
				require.ErrorIs(t, err, ErrPointNotFound)
			} else {
				require.ErrorIs(t, err, ErrFormatUnavailable)
				require.NotContains(t, err.Error(), "private")
			}
		})
	}
}
