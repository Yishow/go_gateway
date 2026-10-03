//go:build !f_write_group_fixture

package main

import (
	"context"
	"database/sql"
	"log"

	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/runtime"
)

// groupFixtureHooks keeps the acceptance-only seam out of normal binaries.
// The default implementation deliberately has no listener or control state.
type groupFixtureHooks interface {
	pipelineConfig(nodeID string, ownerID string) grouppipeline.Config
	sampleSink(*grouppipeline.Pipeline) runtime.SampleSink
	configure(context.Context, *sql.DB, *gatewayServices) (func() error, error)
}

type defaultGroupFixture struct{}

func newGroupFixture() groupFixtureHooks { return defaultGroupFixture{} }

func preflightGroupFixtureDatabase() error { return nil }

func (defaultGroupFixture) pipelineConfig(nodeID, ownerID string) grouppipeline.Config {
	return grouppipeline.Config{
		NodeID:  nodeID,
		Owner:   nodeID + "/" + ownerID,
		OnError: logGroupPipelineError,
	}
}

func (defaultGroupFixture) sampleSink(pipe *grouppipeline.Pipeline) runtime.SampleSink {
	return pipe
}

func (defaultGroupFixture) configure(context.Context, *sql.DB, *gatewayServices) (func() error, error) {
	return func() error { return nil }, nil
}

var _ groupFixtureHooks = defaultGroupFixture{}

func logGroupPipelineError(err error) {
	log.Printf("寫入群組 pipeline 錯誤: %v", err)
}
