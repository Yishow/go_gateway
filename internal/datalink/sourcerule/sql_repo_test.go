package sourcerule

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func setupSQLRepoDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	migrator := datalink.NewMigrator()
	require.NoError(t, migrator.Migrate(db))

	return db
}

func TestSQLRepository_CreateAndList(t *testing.T) {
	db := setupSQLRepoDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()
	now := time.Now().UTC()

	rule := &schema.SourceRule{
		ID:               "rule-1",
		DeviceID:         "device-1",
		StartAddress:     "40001",
		Count:            2,
		DataType:         schema.DataTypeInt16,
		NamingPrefix:     "SRC",
		Enabled:          true,
		Locked:           false,
		Origin:           "manual",
		SkippedAddresses: `[]`,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	require.NoError(t, repo.Create(ctx, rule))

	links := []*schema.SourceRuleLink{
		{
			ID:        "link-1",
			RuleID:    "rule-1",
			Address:   "40001",
			PointID:   "point-1",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        "link-2",
			RuleID:    "rule-1",
			Address:   "40002",
			PointID:   "point-2",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	require.NoError(t, repo.CreateLinks(ctx, links))

	items, err := repo.List(ctx, ListFilter{})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "rule-1", items[0].ID)

	storedLinks, err := repo.ListLinks(ctx, "rule-1")
	require.NoError(t, err)
	assert.Len(t, storedLinks, 2)
	assert.Equal(t, "40001", storedLinks[0].Address)
}
