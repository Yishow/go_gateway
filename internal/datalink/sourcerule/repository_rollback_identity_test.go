package sourcerule

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"go-gateway/internal/datalink/schema"
	"testing"

	"github.com/stretchr/testify/require"
)

// Extend the existing rollback fixture only for the revision read's original failure path.
func (c *rollbackTestConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	if c.execErr != nil {
		return nil, c.execErr
	}
	return nil, errors.New("rollback fixture has no successful query rows")
}

func TestRuleRepositoryRollbackIgnoresOnlyBareTxDone(t *testing.T) {
	primary := errors.New("statement failed")
	cleanup := errors.New("rollback driver cause")
	wrapped := fmt.Errorf("driver wrapped terminal: %w", sql.ErrTxDone)
	joined := errors.Join(sql.ErrTxDone, cleanup)
	for _, operation := range []string{"create-links", "replace-links", "snapshots", "snapshots-at-revision"} {
		t.Run(operation, func(t *testing.T) {
			for _, tc := range []struct {
				name     string
				rollback error
				retained bool
			}{{"nil", nil, false}, {"bare", sql.ErrTxDone, false}, {"wrapped", wrapped, true}, {"joined", joined, true}} {
				t.Run(tc.name, func(t *testing.T) {
					conn := &rollbackTestConn{execErr: primary, rollbackErr: tc.rollback}
					db := sql.OpenDB(rollbackTestConnector{conn: conn})
					t.Cleanup(func() { require.NoError(t, db.Close()) })
					repo := NewSQLRepository(db)
					links := []*schema.SourceRuleLink{{ID: "link", RuleID: "rule", PointID: "point", Address: "40001"}}
					snapshots := []*schema.SourceRuleCandidateSnapshot{{SourceRuleID: "rule", RevisionID: "revision", CandidateType: schema.SourceRuleCandidateTypeTags, Payload: "{}"}}
					var err error
					switch operation {
					case "create-links":
						err = repo.CreateLinks(t.Context(), links)
					case "replace-links":
						err = repo.ReplaceLinks(t.Context(), "rule", links)
					case "snapshots":
						err = repo.ReplaceCandidateSnapshots(t.Context(), snapshots)
					default:
						err = repo.ReplaceCandidateSnapshotsAtRevision(t.Context(), snapshots, "revision")
					}
					require.ErrorIs(t, err, primary)
					require.Equal(t, 1, conn.rollbackCalls)
					require.Equal(t, tc.retained, errors.Is(err, sql.ErrTxDone))
					if tc.retained {
						require.ErrorIs(t, err, tc.rollback)
					}
					require.Equal(t, tc.name == "joined", errors.Is(err, cleanup))
				})
			}
		})
	}
}
