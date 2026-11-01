package lostcommit

import (
	"database/sql/driver"
	"errors"
	"slices"
	"testing"
)

func TestStatementObserverPreservesRejectOrdering(t *testing.T) {
	for _, operation := range []string{"exec", "query", "prepare", "prepare-context"} {
		t.Run(operation, func(t *testing.T) {
			faults := &Faults{}
			inner := &purePrepareConn{}
			conn := &lostConn{Conn: inner, faults: faults}
			order := make([]string, 0, 3)
			const query = "SELECT observer_contract"
			denied := errors.New("statement denied")
			queryDenied := errors.New("query denied")
			observe := func(got string) {
				if got != query {
					t.Errorf("observer changed query: %q", got)
				}
				order = append(order, "observe")
			}
			reject := func(got string) error {
				order = append(order, "reject")
				if got == query {
					return denied
				}
				return nil
			}
			rejectQuery := func(got string, _ []driver.NamedValue) error {
				order = append(order, "query")
				if got == query {
					return queryDenied
				}
				return nil
			}
			faults.ObserveStatement.Store(&observe)
			faults.Reject.Store(&reject)
			faults.RejectQuery.Store(&rejectQuery)
			var err error
			var rows driver.Rows
			switch operation {
			case "exec":
				_, err = conn.ExecContext(t.Context(), query, nil)
			case "query":
				rows, err = conn.QueryContext(t.Context(), query, nil)
			case "prepare":
				_, err = conn.Prepare(query)
			case "prepare-context":
				_, err = conn.PrepareContext(t.Context(), query)
			}
			if rows != nil {
				if closeErr := rows.Close(); closeErr != nil {
					t.Fatal(closeErr)
				}
				t.Fatal("rejected query returned rows")
			}
			if !errors.Is(err, denied) || inner.lastQuery != "" {
				t.Fatalf("Reject lost priority or reached driver: %v / %q", err, inner.lastQuery)
			}
			if !slices.Equal(order, []string{"observe", "reject"}) {
				t.Fatalf("wrong observer/Reject order: %v", order)
			}
			order = order[:0]
			faults.Reject.Store(nil)
			if operation == "query" {
				rows, err = conn.QueryContext(t.Context(), query, nil)
				if rows != nil {
					if closeErr := rows.Close(); closeErr != nil {
						t.Fatal(closeErr)
					}
					t.Fatal("query rejection returned rows")
				}
				if !errors.Is(err, queryDenied) || !slices.Equal(order, []string{"observe", "query"}) {
					t.Fatalf("RejectQuery lost priority: %v / %v", err, order)
				}
			}
			faults.ObserveStatement.Store(nil)
			faults.RejectQuery.Store(nil)
			stmt, err := conn.PrepareContext(t.Context(), query)
			if err != nil {
				t.Fatal(err)
			}
			if err := stmt.Close(); err != nil {
				t.Fatal(err)
			}
			if inner.lastQuery != query {
				t.Fatal("nil observer changed ordinary prepare")
			}
		})
	}
}

func TestStatementObserverDoesNotBypassExecError(t *testing.T) {
	faults := &Faults{}
	observed := false
	observe := func(string) { observed = true }
	faults.ObserveStatement.Store(&observe)
	injected := errors.New("exec fault")
	faults.ExecError.Store(&injected)
	conn := &lostConn{Conn: &purePrepareConn{}, faults: faults}
	_, err := conn.ExecContext(t.Context(), "UPDATE observer_contract", nil)
	if !errors.Is(err, injected) || observed {
		t.Fatalf("ExecError ordering changed: %v / observed=%v", err, observed)
	}
	faults.ExecError.Store(nil)
	_, err = conn.ExecContext(t.Context(), "UPDATE observer_contract", nil)
	if !errors.Is(err, driver.ErrSkip) || !observed {
		t.Fatalf("observer must precede unchanged fallback: %v / observed=%v", err, observed)
	}
}
