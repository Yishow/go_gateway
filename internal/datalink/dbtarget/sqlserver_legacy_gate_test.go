package dbtarget

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestSQLServerLegacyWritesStayBlockedBeforeOpen(t *testing.T) {
	config, err := json.Marshal(sqlServerConnectorConfigForTest())
	require.NoError(t, err)
	connector := &schema.DatabaseConnector{Kind: schema.DatabaseConnectorKindSQLServer, ConnectionConfig: string(config)}
	opens := 0
	previous := openExternalDBManagerFunc
	openExternalDBManagerFunc = func(schema.DatabaseConnectorKind, ConnectionConfig) (*datalinkbase.DBManager, error) {
		opens++
		return nil, errors.New("unexpected destination open")
	}
	t.Cleanup(func() { openExternalDBManagerFunc = previous })
	writer := &Writer{}
	t.Run("single mapping", func(t *testing.T) {
		err := writer.writeMapping(t.Context(), connector, &schema.DatabaseTargetMapping{TableName: "unconfirmed_table", ColumnName: "value"}, 1, time.Now())
		require.ErrorIs(t, err, ErrManagedDestinationInvalid)
		require.Zero(t, opens)
	})
	t.Run("grouped mapping", func(t *testing.T) {
		err := writer.flushGroupedBucket(t.Context(), &groupedWriteBucket{Connector: connector, Values: map[string]any{"value": 1}})
		require.ErrorIs(t, err, ErrManagedDestinationInvalid)
		require.Zero(t, opens)
	})
}
