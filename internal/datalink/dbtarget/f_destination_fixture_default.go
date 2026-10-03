//go:build !f_write_group_fixture

package dbtarget

import (
	"database/sql"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"
)

// FixtureDestinationFaultKind identifies a fault that can be installed only
// by the disposable write-group fixture binary.
type FixtureDestinationFaultKind string

const (
	FixtureDestinationFaultTargetCommitResponseLost FixtureDestinationFaultKind = "target_commit_response_lost"
	FixtureDestinationFaultTargetCommitHold         FixtureDestinationFaultKind = "target_commit_hold"
)

// FixtureDestinationFaultState is the safe, non-credential state exposed by
// the fixture controller.
type FixtureDestinationFaultState struct {
	ConnectorID string                      `json:"connector_id"`
	Kind        FixtureDestinationFaultKind `json:"kind"`
	Enabled     bool                        `json:"enabled"`
	Reached     uint64                      `json:"reached"`
	Holding     bool                        `json:"holding"`
	Commits     uint64                      `json:"commits"`
}

// FixtureDestinationFaultController is a no-op in a normal production build.
type FixtureDestinationFaultController struct{}

// NewFixtureDestinationFaultController returns the normal-build no-op.
func NewFixtureDestinationFaultController() *FixtureDestinationFaultController {
	return &FixtureDestinationFaultController{}
}

// InstallFixtureDestinationController is intentionally inert in production.
func InstallFixtureDestinationController(*FixtureDestinationFaultController) {}

// UninstallFixtureDestinationController is intentionally inert in production.
func UninstallFixtureDestinationController(*FixtureDestinationFaultController) {}

// Configure is intentionally inert in production.
func (*FixtureDestinationFaultController) Configure(string, FixtureDestinationFaultKind, bool) (FixtureDestinationFaultState, error) {
	return FixtureDestinationFaultState{}, nil
}

// Snapshot returns no fixture faults in production.
func (*FixtureDestinationFaultController) Snapshot() []FixtureDestinationFaultState { return nil }

// Close is intentionally inert in production.
func (*FixtureDestinationFaultController) Close() {}

func openFixtureDestination(string, schema.DatabaseConnectorKind, ConnectionConfig, *datalink.DBManager) (*sql.DB, func() error, error) {
	return nil, nil, nil
}
