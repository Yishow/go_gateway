//go:build !f_write_group_fixture

package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalPipelineIgnoresTheFixtureSecondaryNode(t *testing.T) {
	t.Setenv("F_FIXTURE_WORKER_NODE", "secondary")
	config := newGroupFixture().pipelineConfig("primary", "incarnation")
	require.Equal(t, "primary", config.NodeID)
	require.Equal(t, "primary/incarnation", config.Owner)
}
