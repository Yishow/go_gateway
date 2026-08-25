package handlers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSafeRuntimeApplyMessageDoesNotExposeRawException(t *testing.T) {
	message := safeRuntimeApplyMessage("dial tcp plc.internal:502: dsn=postgres://secret", "failed")
	require.Equal(t, "runtime apply failed", message)
	require.NotContains(t, message, "plc.internal")
	require.NotContains(t, message, "postgres://secret")
}
