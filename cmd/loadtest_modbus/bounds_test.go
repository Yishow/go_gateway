package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStartServersRejectsValuesBeforeOpeningListeners(t *testing.T) {
	for _, counts := range [][2]int{{-1, 1}, {1, -1}, {2, 0}, {1, 1}} {
		servers, err := startServers(counts[0], counts[1], 65535)
		require.Error(t, err)
		require.Nil(t, servers)
	}
}
