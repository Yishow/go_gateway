package settings

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSQLRepository_InvalidJSONIsTypedStorageError(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewSQLRepository(db)
	_, err := db.ExecContext(context.Background(), `INSERT INTO system_settings (key, value, updated_at) VALUES (?, ?, ?)`, "broken", "not-json", time.Now())
	require.NoError(t, err)

	_, err = repo.Get(context.Background(), "broken")
	var storageErr *StorageError
	require.ErrorAs(t, err, &storageErr)
	require.Equal(t, "decode", storageErr.Operation)
}
