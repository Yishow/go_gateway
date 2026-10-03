//go:build !f_write_group_fixture

package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalBinaryRetryEnvironmentIsIgnored(t *testing.T) {
	t.Setenv(fixtureMaxRetriesEnv, "2")
	require.NoError(t, preflightFixtureRetries())
	require.Zero(t, groupDeliverySenderConfig().MaxRetries)
}
