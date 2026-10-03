//go:build f_write_group_fixture

package dbtarget

import (
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/testutil/lostcommit"
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

type fixtureDestinationFault struct {
	connectorID string
	kind        FixtureDestinationFaultKind
	enabled     atomic.Bool
	reached     atomic.Uint64
	holding     atomic.Bool
	commits     atomic.Uint64
	faults      *lostcommit.Faults
	barrier     *fixtureDestinationBarrier
}

type fixtureDestinationBarrier struct {
	release chan struct{}
	once    sync.Once
}

func newFixtureDestinationBarrier() *fixtureDestinationBarrier {
	return &fixtureDestinationBarrier{release: make(chan struct{})}
}

func (b *fixtureDestinationBarrier) wait() { <-b.release }

func (b *fixtureDestinationBarrier) close() { b.once.Do(func() { close(b.release) }) }

// FixtureDestinationFaultController owns all target-side fixture fault state.
// It is deliberately separate from the production delivery writer.
type FixtureDestinationFaultController struct {
	mu    sync.RWMutex
	items map[string]*fixtureDestinationFault
}

var fixtureDestinationController atomic.Pointer[FixtureDestinationFaultController]

// NewFixtureDestinationFaultController returns an empty tagged controller.
func NewFixtureDestinationFaultController() *FixtureDestinationFaultController {
	return &FixtureDestinationFaultController{items: make(map[string]*fixtureDestinationFault)}
}

// InstallFixtureDestinationController makes c visible to future destination
// opens. Passing nil removes the process-local fixture controller.
func InstallFixtureDestinationController(c *FixtureDestinationFaultController) {
	fixtureDestinationController.Store(c)
}

// UninstallFixtureDestinationController removes c only when it is still the
// active controller, so one fixture cleanup cannot clear another fixture.
func UninstallFixtureDestinationController(c *FixtureDestinationFaultController) {
	fixtureDestinationController.CompareAndSwap(c, nil)
}

// Configure enables or disables one connector-scoped target fault.
func (c *FixtureDestinationFaultController) Configure(connectorID string, kind FixtureDestinationFaultKind, enabled bool) (FixtureDestinationFaultState, error) {
	connectorID = strings.TrimSpace(connectorID)
	if connectorID == "" {
		return FixtureDestinationFaultState{}, errors.New("fixture connector is required")
	}
	if kind != FixtureDestinationFaultTargetCommitResponseLost && kind != FixtureDestinationFaultTargetCommitHold {
		return FixtureDestinationFaultState{}, errors.New("unsupported fixture destination fault")
	}
	if c == nil {
		return FixtureDestinationFaultState{}, errors.New("fixture destination controller is unavailable")
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	item := c.items[connectorID]
	if item == nil {
		item = &fixtureDestinationFault{
			connectorID: connectorID,
			faults:      &lostcommit.Faults{Lose: &atomic.Bool{}},
		}
		c.items[connectorID] = item
	}
	if item.barrier != nil {
		item.barrier.close()
		item.barrier = nil
	}
	item.kind = kind
	item.enabled.Store(enabled)
	item.faults.Lose.Store(enabled && kind == FixtureDestinationFaultTargetCommitResponseLost)
	item.faults.AfterCommit.Store(nil)
	if enabled {
		var barrier *fixtureDestinationBarrier
		if kind == FixtureDestinationFaultTargetCommitHold {
			barrier = newFixtureDestinationBarrier()
			item.barrier = barrier
		}
		hook := func() {
			item.reached.Add(1)
			item.commits.Add(1)
			if barrier == nil {
				return
			}
			item.holding.Store(true)
			barrier.wait()
			item.holding.Store(false)
		}
		item.faults.AfterCommit.Store(&hook)
	}
	return item.state(), nil
}

// Snapshot returns deterministic, safe target fault state.
func (c *FixtureDestinationFaultController) Snapshot() []FixtureDestinationFaultState {
	if c == nil {
		return nil
	}
	c.mu.RLock()
	states := make([]FixtureDestinationFaultState, 0, len(c.items))
	for _, item := range c.items {
		states = append(states, item.state())
	}
	c.mu.RUnlock()
	slices.SortFunc(states, func(left, right FixtureDestinationFaultState) int {
		return strings.Compare(left.ConnectorID, right.ConnectorID)
	})
	return states
}

// Close releases all target hold barriers and disables all target faults.
func (c *FixtureDestinationFaultController) Close() {
	if c == nil {
		return
	}
	c.mu.Lock()
	for _, item := range c.items {
		item.enabled.Store(false)
		item.faults.Lose.Store(false)
		item.faults.AfterCommit.Store(nil)
		if item.barrier != nil {
			item.barrier.close()
			item.barrier = nil
		}
	}
	c.items = make(map[string]*fixtureDestinationFault)
	c.mu.Unlock()
}

func (f *fixtureDestinationFault) state() FixtureDestinationFaultState {
	return FixtureDestinationFaultState{
		ConnectorID: f.connectorID,
		Kind:        f.kind,
		Enabled:     f.enabled.Load(),
		Reached:     f.reached.Load(),
		Holding:     f.holding.Load(),
		Commits:     f.commits.Load(),
	}
}

func (c *FixtureDestinationFaultController) active(connectorID string) *fixtureDestinationFault {
	c.mu.RLock()
	item := c.items[connectorID]
	if item == nil || !item.enabled.Load() {
		c.mu.RUnlock()
		return nil
	}
	c.mu.RUnlock()
	return item
}

func openFixtureDestination(connectorID string, kind schema.DatabaseConnectorKind, config ConnectionConfig, manager *datalinkbase.DBManager) (*sql.DB, func() error, error) {
	controller := fixtureDestinationController.Load()
	if controller == nil {
		return nil, nil, nil
	}
	item := controller.active(strings.TrimSpace(connectorID))
	if item == nil {
		return nil, nil, nil
	}
	if manager == nil || manager.DB() == nil {
		return nil, nil, errors.New("fixture destination manager is unavailable")
	}
	dbConfig, err := buildExternalDBConfig(kind, config)
	if err != nil {
		return nil, nil, fmt.Errorf("fixture destination configuration is invalid: %w", err)
	}
	wrapped := lostcommit.WrapFaults(manager.DB().Driver(), dbConfig.DSN, item.faults)
	closeFn := func() error {
		return errors.Join(wrapped.Close(), manager.Close())
	}
	return wrapped, closeFn, nil
}
