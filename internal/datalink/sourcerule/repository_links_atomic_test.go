package sourcerule

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
	"modernc.org/sqlite"
)

var linkPauseFunctionSequence atomic.Int64

func atomicLinkSQLDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "links.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	db.SetMaxOpenConns(4)
	_, err = db.ExecContext(t.Context(), `PRAGMA journal_mode=WAL`)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `CREATE TABLE source_rule_links(id TEXT PRIMARY KEY,rule_id TEXT NOT NULL,address TEXT NOT NULL,point_id TEXT NOT NULL UNIQUE,tag_id TEXT,mapping_id TEXT,created_at DATETIME NOT NULL,updated_at DATETIME NOT NULL,UNIQUE(rule_id,address))`)
	require.NoError(t, err)
	return db
}

func atomicLinkFixtures() []*schema.SourceRuleLink {
	now := time.Now()
	return []*schema.SourceRuleLink{{ID: "link1", RuleID: "rule1", Address: "40001", PointID: "point1", CreatedAt: now, UpdatedAt: now}, {ID: "link2", RuleID: "rule1", Address: "40002", PointID: "point2", CreatedAt: now, UpdatedAt: now}}
}

func TestSQLReplaceLinksFailureRollsBackWholeSetAndReadersNeverSeeDeleteGap(t *testing.T) {
	entered, resume := make(chan struct{}), make(chan struct{})
	name := fmt.Sprintf("pause_link_replace_%d", linkPauseFunctionSequence.Add(1))
	require.NoError(t, sqlite.RegisterScalarFunction(name, 0, func(*sqlite.FunctionContext, []driver.Value) (driver.Value, error) {
		close(entered)
		<-resume
		return nil, errors.New("injected second insert failure")
	}))
	db := atomicLinkSQLDB(t)
	repo := NewSQLRepository(db)
	require.NoError(t, repo.CreateLinks(t.Context(), atomicLinkFixtures()))
	before, err := repo.ListLinks(t.Context(), "rule1")
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `CREATE TRIGGER pause_link_insert BEFORE INSERT ON source_rule_links WHEN NEW.address='40002' BEGIN SELECT `+name+`(); END`)
	require.NoError(t, err)
	next := cloneSourceRuleLinks(before)
	next[0].MappingID = common.Ptr("new-mapping")
	done := make(chan error, 1)
	go func() { done <- repo.ReplaceLinks(context.Background(), "rule1", next) }()
	select {
	case <-entered:
	case <-time.After(3 * time.Second):
		close(resume)
		t.Fatal("did not pause after transactional DELETE and first INSERT")
	}
	readerCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	during, readErr := repo.ListLinks(readerCtx, "rule1")
	close(resume)
	require.NoError(t, readErr)
	require.Equal(t, before, during, "a reader must see all original links during unfinished replacement")
	require.ErrorContains(t, <-done, "injected second insert failure")
	after, err := repo.ListLinks(t.Context(), "rule1")
	require.NoError(t, err)
	require.Equal(t, before, after)
}

func TestSQLReplaceLinksUniqueFailurePreservesExistingLinks(t *testing.T) {
	db := atomicLinkSQLDB(t)
	repo := NewSQLRepository(db)
	links := atomicLinkFixtures()
	require.NoError(t, repo.CreateLinks(t.Context(), links))
	foreign := *links[1]
	foreign.ID = "foreign"
	foreign.RuleID = "rule2"
	foreign.PointID = "foreign-point"
	require.NoError(t, repo.CreateLinks(t.Context(), []*schema.SourceRuleLink{&foreign}))
	before, err := repo.ListLinks(t.Context(), "rule1")
	require.NoError(t, err)
	next := cloneSourceRuleLinks(before)
	next[1].PointID = foreign.PointID
	require.ErrorContains(t, repo.ReplaceLinks(t.Context(), "rule1", next), "UNIQUE")
	after, err := repo.ListLinks(t.Context(), "rule1")
	require.NoError(t, err)
	require.Equal(t, before, after)
	require.NoError(t, repo.ReplaceLinks(t.Context(), "rule1", nil))
	after, err = repo.ListLinks(t.Context(), "rule1")
	require.NoError(t, err)
	require.Empty(t, after)
}

func TestMemoryReplaceLinksRejectsInvalidSetWithoutChangingOriginal(t *testing.T) {
	repo := NewMemoryRepository()
	links := atomicLinkFixtures()
	require.NoError(t, repo.CreateLinks(t.Context(), links))
	before, err := repo.ListLinks(t.Context(), "rule1")
	require.NoError(t, err)
	invalid := cloneSourceRuleLinks(links)
	invalid[1].PointID = invalid[0].PointID
	require.Error(t, repo.ReplaceLinks(t.Context(), "rule1", invalid))
	after, err := repo.ListLinks(t.Context(), "rule1")
	require.NoError(t, err)
	require.Equal(t, before, after)
	next := cloneSourceRuleLinks(links)
	next[0].MappingID = common.Ptr("mapping-new")
	require.NoError(t, repo.ReplaceLinks(t.Context(), "rule1", next))
	next[0].MappingID = common.Ptr("external-edit")
	after, err = repo.ListLinks(t.Context(), "rule1")
	require.NoError(t, err)
	require.Equal(t, "mapping-new", *after[0].MappingID)
	require.NoError(t, repo.ReplaceLinks(t.Context(), "rule1", nil))
	after, err = repo.ListLinks(t.Context(), "rule1")
	require.NoError(t, err)
	require.Empty(t, after)
}
