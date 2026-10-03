//go:build f_write_group_fixture

package main

import (
	"errors"
	"os"
)

const fixtureWorkerNodeEnv = "F_FIXTURE_WORKER_NODE"

// The only alternate identity is the secondary acceptance worker. Production
// treats a new incarnation of the same node as a restart of a dead owner.
func preflightFixtureWorkerNode() error {
	value := os.Getenv(fixtureWorkerNodeEnv)
	if value != "" && value != "secondary" {
		return errors.New("F_FIXTURE_WORKER_NODE must be empty or secondary")
	}
	return nil
}

func fixtureWorkerNodeID(defaultNode string) string {
	if os.Getenv(fixtureWorkerNodeEnv) == "secondary" {
		return "f-acceptance-secondary"
	}
	return defaultNode
}
