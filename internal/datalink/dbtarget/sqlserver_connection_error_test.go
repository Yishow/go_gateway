package dbtarget

import (
	"errors"
	"fmt"
	"testing"

	"go-gateway/internal/datalink/schema"

	mssql "github.com/microsoft/go-mssqldb"
	"github.com/stretchr/testify/require"
)

func TestSQLServerConnectorSafeErrors(t *testing.T) {
	t.Parallel()
	for name, source := range map[string]error{
		"certificate":    errors.New("TLS certificate invalid for private.fixture.invalid: test-only-password"),
		"authentication": mssql.Error{Number: 18456, Message: "Login failed for user private-login; test-only-password"},
		"transport":      errors.New("dial private.fixture.invalid:1433: test-only-password"),
	} {
		t.Run(name, func(t *testing.T) {
			message := connectorErrorMessage(schema.DatabaseConnectorKindSQLServer, source)
			require.NotContains(t, message, "private.fixture.invalid")
			require.NotContains(t, message, "private-login")
			require.NotContains(t, message, "test-only-password")
			require.NotEmpty(t, message)
		})
	}
}

func TestSQLServerConnectorAuthenticationClassification(t *testing.T) {
	t.Parallel()
	err := fmt.Errorf("wrapped: %w", mssql.Error{Number: 18456, Message: "Login failed"})
	require.Equal(t, schema.DatabaseConnectorStatusAuthFailed, classifyConnectorError(err))
	other := errors.New("ordinary transport outage")
	require.Equal(t, schema.DatabaseConnectorStatusUnreachable, classifyConnectorError(other))
}
