package dbtarget

import (
	"database/sql"
	"errors"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type sqlServerDestinationTestSession struct {
	db     *sql.DB
	closed int
}

func (s *sqlServerDestinationTestSession) DB() *sql.DB  { return s.db }
func (s *sqlServerDestinationTestSession) Close() error { s.closed++; return s.db.Close() }

func TestSQLServerDestinationKeepsUnverifiedWritesClosed(t *testing.T) {
	script := sqlServerRecordingMetadataFixture()
	script.createTable = false // unverified permission
	session := &sqlServerDestinationTestSession{db: openSQLServerRecordingMetadataScript(t, script)}
	previous := openSQLServerManagedSessionFunc
	openSQLServerManagedSessionFunc = func(schema.DatabaseConnectorKind, ConnectionConfig) (sqlServerManagedSession, error) {
		return session, nil
	}
	t.Cleanup(func() { openSQLServerManagedSessionFunc = previous })
	config := sqlServerConnectorConfigForTest()
	config["database"], config["schema"] = script.database, script.expectedSchema
	connector := &schema.DatabaseConnector{ID: "sql-connector", IdentityRevision: "revision-1", Kind: schema.DatabaseConnectorKindSQLServer}
	opened, err := openSQLServerDeliveryDestination(t.Context(), connector, config)
	require.ErrorIs(t, err, ErrDestinationBlocked)
	require.Nil(t, opened)
	require.Equal(t, 1, session.closed)
	require.Len(t, script.queries, 4, "the opener reuses only backend-derived read-only metadata")
	require.Zero(t, script.execCalls)
	require.Zero(t, script.beginCalls)
}

func TestSQLServerDestinationOpensWhenVerified(t *testing.T) {
	script := sqlServerRecordingMetadataFixture()
	session := &sqlServerDestinationTestSession{db: openSQLServerRecordingMetadataScript(t, script)}
	previous := openSQLServerManagedSessionFunc
	openSQLServerManagedSessionFunc = func(schema.DatabaseConnectorKind, ConnectionConfig) (sqlServerManagedSession, error) {
		return session, nil
	}
	t.Cleanup(func() { openSQLServerManagedSessionFunc = previous })
	config := sqlServerConnectorConfigForTest()
	config["database"], config["schema"] = script.database, script.expectedSchema
	connector := &schema.DatabaseConnector{ID: "sql-connector", IdentityRevision: "revision-1", Kind: schema.DatabaseConnectorKindSQLServer}
	opened, err := openSQLServerDeliveryDestination(t.Context(), connector, config)
	require.NoError(t, err)
	require.NotNil(t, opened)
	require.True(t, opened.SQLServer.SupportsManagedLong)
}

func TestSQLServerDestinationFailureMessagesStaySafe(t *testing.T) {
	previous := openSQLServerManagedSessionFunc
	t.Cleanup(func() { openSQLServerManagedSessionFunc = previous })
	for _, tc := range []struct {
		message string
		want    error
	}{
		{"x509 certificate invalid for private-endpoint password=hidden", ErrDestinationBlocked},
		{"network unavailable private-endpoint hidden", ErrDestinationUnreachable},
	} {
		t.Run(tc.want.Error(), func(t *testing.T) {
			openSQLServerManagedSessionFunc = func(schema.DatabaseConnectorKind, ConnectionConfig) (sqlServerManagedSession, error) {
				return nil, errors.New(tc.message)
			}
			_, err := openSQLServerDeliveryDestination(t.Context(), &schema.DatabaseConnector{Kind: schema.DatabaseConnectorKindSQLServer}, sqlServerConnectorConfigForTest())
			require.ErrorIs(t, err, tc.want)
			require.NotContains(t, err.Error(), "private-endpoint")
			require.NotContains(t, err.Error(), "hidden")
		})
	}
}
