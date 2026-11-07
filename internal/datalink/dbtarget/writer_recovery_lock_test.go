package dbtarget

import (
	"sync"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestWriterRecoveryRejectsInvalidLock(t *testing.T) {
	for _, lock := range []any{"invalid", (*sync.Mutex)(nil)} {
		writer := &Writer{}
		writer.recoveryLocks.Store("connector", lock)
		err := writer.recoverWriterTable(t.Context(), nil, &schema.DatabaseConnector{ID: "connector"}, "", "readings")
		require.ErrorContains(t, err, "invalid writer recovery lock")
	}
}
