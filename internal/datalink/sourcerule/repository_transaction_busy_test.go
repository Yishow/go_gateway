package sourcerule

import (
	"database/sql"
	"go-gateway/internal/datalink/schema"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRuleRepositoryCommitBusyCleansNativeTransaction(t *testing.T) {
	for _, operation := range []string{"create-links", "replace-links", "candidates", "candidate-cas"} {
		t.Run(operation, func(t *testing.T) {
			db := atomicLinkSQLDB(t)
			db.SetMaxOpenConns(1)
			_, err := db.ExecContext(t.Context(), `PRAGMA journal_mode=DELETE; PRAGMA busy_timeout=0; CREATE TABLE source_rules(id TEXT PRIMARY KEY,revision_id TEXT); INSERT INTO source_rules VALUES('rule1','rev'); CREATE TABLE source_rule_candidate_snapshots(source_rule_id TEXT,revision_id TEXT,candidate_type TEXT,payload TEXT,status TEXT,reason TEXT,generated_at DATETIME); INSERT INTO source_rule_candidate_snapshots VALUES('rule1','rev','tag','old','ready','',CURRENT_TIMESTAMP)`)
			require.NoError(t, err)
			repo := NewSQLRepository(db)
			links := atomicLinkFixtures()
			if operation != "create-links" {
				require.NoError(t, repo.CreateLinks(t.Context(), links))
			}
			var seq int
			var name, path string
			require.NoError(t, db.QueryRowContext(t.Context(), `PRAGMA database_list`).Scan(&seq, &name, &path))
			reader, err := sql.Open("sqlite", path)
			require.NoError(t, err)
			defer reader.Close()
			readTx, err := reader.BeginTx(t.Context(), nil)
			require.NoError(t, err)
			defer readTx.Rollback()
			var before int
			require.NoError(t, readTx.QueryRowContext(t.Context(), `SELECT count(*) FROM source_rule_links`).Scan(&before))
			snapshots := []*schema.SourceRuleCandidateSnapshot{{SourceRuleID: "rule1", RevisionID: "rev", CandidateType: "tag", Payload: "new", Status: "ready", GeneratedAt: time.Now()}}
			switch operation {
			case "create-links":
				err = repo.CreateLinks(t.Context(), links)
			case "replace-links":
				err = repo.ReplaceLinks(t.Context(), "rule1", nil)
			case "candidates":
				err = repo.ReplaceCandidateSnapshots(t.Context(), snapshots)
			case "candidate-cas":
				err = repo.ReplaceCandidateSnapshotsAtRevision(t.Context(), snapshots, "rev")
			}
			require.Error(t, err)
			var after int
			require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count(*) FROM source_rule_links`).Scan(&after))
			require.Equal(t, before, after)
			var payload string
			require.NoError(t, db.QueryRowContext(t.Context(), `SELECT payload FROM source_rule_candidate_snapshots`).Scan(&payload))
			require.Equal(t, "old", payload)
			require.NoError(t, readTx.Rollback())
			next, err := db.BeginTx(t.Context(), nil)
			require.NoError(t, err)
			require.NoError(t, next.Rollback())
		})
	}
}
