//go:build f_write_group_fixture

package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFixtureWorkerIdentityUsesADistinctNodeForTheSecondaryProcess(t *testing.T) {
	t.Setenv("F_FIXTURE_WORKER_NODE", "secondary")
	config := newGroupFixture().pipelineConfig("primary", "incarnation")
	require.Equal(t, "f-acceptance-secondary", config.NodeID)
	require.Equal(t, "f-acceptance-secondary/incarnation", config.Owner)
}

func TestFixtureWorkerNodeRejectsAnArbitraryIdentityBeforeDatabaseOpen(t *testing.T) {
	fixtureOwnedDatabasePath(t, fixtureMarkerValue)
	t.Setenv(fixtureWorkerNodeEnv, "arbitrary-worker")
	require.Error(t, preflightGroupFixtureDatabase())
}

func TestFixtureWorkerNodeLeavesDefaultIdentityUnchanged(t *testing.T) {
	t.Setenv(fixtureWorkerNodeEnv, "")
	config := newGroupFixture().pipelineConfig("primary", "incarnation")
	require.Equal(t, "primary", config.NodeID)
	require.Equal(t, "primary/incarnation", config.Owner)
}
