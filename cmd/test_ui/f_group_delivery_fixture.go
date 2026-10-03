//go:build f_write_group_fixture

package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go-gateway/internal/datalink/dbtarget"

	"modernc.org/sqlite"
)

const (
	fixtureClosureFunctionName = "f_write_group_fixture_closure_hold"
	fixtureFaultLocalReceipt   = "local_receipt_failure"
	fixtureFaultClosureHold    = "closure_hold"
)

var errFixtureClosureRegistration = sqlite.RegisterScalarFunction(
	fixtureClosureFunctionName,
	1,
	fixtureClosureScalar,
)

var (
	errFixtureFaultGroupUnknown = errors.New("fixture group is unknown")
	errFixtureFaultKindUnknown  = errors.New("fixture fault kind is unsupported")
	errFixtureFaultUnavailable  = errors.New("fixture fault controller is unavailable")
)

type fixtureFaultRequest struct {
	Kind    string `json:"kind"`
	GroupID string `json:"group_id"`
	Enabled bool   `json:"enabled"`
}

type fixtureFaultState struct {
	Kind        string `json:"kind"`
	GroupID     string `json:"group_id"`
	ConnectorID string `json:"connector_id"`
	Enabled     bool   `json:"enabled"`
	Reached     uint64 `json:"reached"`
	Holding     bool   `json:"holding"`
	Commits     uint64 `json:"commits"`
}

type fixtureFaultEntry struct {
	groupID        string
	connectorID    string
	kind           string
	enabled        atomic.Bool
	reached        atomic.Uint64
	closureGate    *fixtureClosureGate
	receiptTrigger string
	closureTrigger string
}

type fixtureFaultController struct {
	db          *sql.DB
	deps        *gatewayServices
	destination *dbtarget.FixtureDestinationFaultController
	now         func() time.Time
	opMu        sync.Mutex
	mu          sync.RWMutex
	entries     map[string]*fixtureFaultEntry
}

type fixtureClosureGate struct {
	groupID string
	release chan struct{}
	once    sync.Once
	reached atomic.Uint64
	holding atomic.Bool
}

var fixtureClosureRegistry = struct {
	sync.RWMutex
	gates map[string]*fixtureClosureGate
}{gates: make(map[string]*fixtureClosureGate)}

func newFixtureFaultController(db *sql.DB, deps *gatewayServices, now func() time.Time) *fixtureFaultController {
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &fixtureFaultController{
		db:          db,
		deps:        deps,
		now:         now,
		destination: dbtarget.NewFixtureDestinationFaultController(),
		entries:     make(map[string]*fixtureFaultEntry),
	}
}

func (c *fixtureFaultController) install() error {
	if c == nil || c.db == nil || c.deps == nil || c.deps.writeGroups == nil {
		return errFixtureFaultUnavailable
	}
	if errFixtureClosureRegistration != nil {
		return errors.New("fixture closure function registration failed")
	}
	dbtarget.InstallFixtureDestinationController(c.destination)
	return nil
}

func (c *fixtureFaultController) configure(ctx context.Context, request fixtureFaultRequest) (fixtureFaultState, error) {
	if c == nil || c.db == nil || c.deps == nil || c.deps.writeGroups == nil {
		return fixtureFaultState{}, errFixtureFaultUnavailable
	}
	kind := strings.TrimSpace(request.Kind)
	groupID := strings.TrimSpace(request.GroupID)
	if !fixtureFaultKindSupported(kind) || groupID == "" {
		return fixtureFaultState{}, errFixtureFaultKindUnknown
	}
	result, err := c.deps.writeGroups.ResolveAppliedAt(ctx, groupID, c.now())
	if err != nil || result == nil || result.Group == nil {
		return fixtureFaultState{}, errFixtureFaultGroupUnknown
	}
	group := result.Group
	groupID = result.GroupID
	connectorID := strings.TrimSpace(group.Destination.ConnectorID)
	if connectorID == "" {
		return fixtureFaultState{}, errFixtureFaultGroupUnknown
	}

	c.opMu.Lock()
	defer c.opMu.Unlock()
	entry := c.entry(groupID)
	if err := c.disableEntry(ctx, entry); err != nil {
		return fixtureFaultState{}, errors.New("fixture fault could not be cleared")
	}
	c.mu.Lock()
	entry.connectorID = connectorID
	entry.kind = kind
	c.mu.Unlock()
	entry.enabled.Store(false)

	if !request.Enabled {
		if isTargetFaultKind(kind) {
			if _, err := c.destination.Configure(connectorID, dbtarget.FixtureDestinationFaultKind(kind), false); err != nil {
				return fixtureFaultState{}, errors.New("fixture target fault could not be cleared")
			}
		}
		var dropErr error
		switch kind {
		case fixtureFaultLocalReceipt:
			dropErr = c.dropTrigger(ctx, fixtureTriggerName("f_fixture_receipt", groupID))
		case fixtureFaultClosureHold:
			releaseFixtureClosureGate(groupID)
			dropErr = c.dropTrigger(ctx, fixtureTriggerName("f_fixture_closure", groupID))
		}
		if dropErr != nil {
			return fixtureFaultState{}, errors.New("fixture fault could not be cleared")
		}
		return c.state(entry), nil
	}

	switch kind {
	case string(dbtarget.FixtureDestinationFaultTargetCommitResponseLost), string(dbtarget.FixtureDestinationFaultTargetCommitHold):
		if _, err := c.destination.Configure(connectorID, dbtarget.FixtureDestinationFaultKind(kind), true); err != nil {
			return fixtureFaultState{}, errors.New("fixture target fault could not be enabled")
		}
	case fixtureFaultLocalReceipt:
		trigger, err := c.createReceiptFailureTrigger(ctx, groupID)
		if err != nil {
			return fixtureFaultState{}, errors.New("fixture receipt fault could not be enabled")
		}
		c.mu.Lock()
		entry.receiptTrigger = trigger
		c.mu.Unlock()
	case fixtureFaultClosureHold:
		gate := &fixtureClosureGate{groupID: groupID, release: make(chan struct{})}
		installFixtureClosureGate(gate)
		trigger, err := c.createClosureHoldTrigger(ctx, groupID)
		if err != nil {
			removeFixtureClosureGate(groupID, gate)
			return fixtureFaultState{}, errors.New("fixture closure fault could not be enabled")
		}
		c.mu.Lock()
		entry.closureGate, entry.closureTrigger = gate, trigger
		c.mu.Unlock()
	default:
		return fixtureFaultState{}, errFixtureFaultKindUnknown
	}
	entry.enabled.Store(true)
	return c.state(entry), nil
}

func (c *fixtureFaultController) disableEntry(ctx context.Context, entry *fixtureFaultEntry) error {
	if entry == nil {
		return nil
	}
	c.mu.Lock()
	kind, connectorID := entry.kind, entry.connectorID
	gate, receiptTrigger, closureTrigger := entry.closureGate, entry.receiptTrigger, entry.closureTrigger
	entry.closureGate, entry.receiptTrigger, entry.closureTrigger = nil, "", ""
	c.mu.Unlock()
	entry.enabled.Store(false)
	if isTargetFaultKind(kind) && connectorID != "" {
		if _, err := c.destination.Configure(connectorID, dbtarget.FixtureDestinationFaultKind(kind), false); err != nil {
			return err
		}
	}
	if gate != nil {
		removeFixtureClosureGate(entry.groupID, gate)
		gate.releaseGate()
	}
	var dropErr error
	if receiptTrigger != "" {
		dropErr = errors.Join(dropErr, c.dropTrigger(ctx, receiptTrigger))
	}
	if closureTrigger != "" {
		dropErr = errors.Join(dropErr, c.dropTrigger(ctx, closureTrigger))
	}
	return dropErr
}

func (c *fixtureFaultController) close(ctx context.Context) error {
	if c == nil {
		return nil
	}
	c.opMu.Lock()
	defer c.opMu.Unlock()
	c.mu.RLock()
	entries := make([]*fixtureFaultEntry, 0, len(c.entries))
	for _, entry := range c.entries {
		entries = append(entries, entry)
	}
	c.mu.RUnlock()
	var closeErr error
	for _, entry := range entries {
		closeErr = errors.Join(closeErr, c.disableEntry(ctx, entry))
	}
	c.destination.Close()
	dbtarget.UninstallFixtureDestinationController(c.destination)
	return closeErr
}

func (c *fixtureFaultController) snapshot() []fixtureFaultState {
	if c == nil {
		return nil
	}
	c.mu.RLock()
	states := make([]fixtureFaultState, 0, len(c.entries))
	for _, entry := range c.entries {
		states = append(states, c.stateLocked(entry))
	}
	c.mu.RUnlock()
	slices.SortFunc(states, func(left, right fixtureFaultState) int {
		return strings.Compare(left.GroupID, right.GroupID)
	})
	return states
}

func (c *fixtureFaultController) state(entry *fixtureFaultEntry) fixtureFaultState {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.stateLocked(entry)
}

func (c *fixtureFaultController) stateLocked(entry *fixtureFaultEntry) fixtureFaultState {
	state := fixtureFaultState{
		Kind: entry.kind, GroupID: entry.groupID, ConnectorID: entry.connectorID,
		Enabled: entry.enabled.Load(), Reached: entry.reached.Load(),
	}
	if gate := entry.closureGate; gate != nil {
		state.Reached = gate.reached.Load()
		state.Holding = gate.holding.Load()
	}
	if isTargetFaultKind(entry.kind) {
		for _, target := range c.destination.Snapshot() {
			if target.ConnectorID != entry.connectorID {
				continue
			}
			state.Enabled = target.Enabled
			state.Reached = target.Reached
			state.Holding = target.Holding
			state.Commits = target.Commits
			break
		}
	}
	return state
}

func (c *fixtureFaultController) entry(groupID string) *fixtureFaultEntry {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := c.entries[groupID]
	if entry == nil {
		entry = &fixtureFaultEntry{groupID: groupID}
		c.entries[groupID] = entry
	}
	return entry
}

func (c *fixtureFaultController) createReceiptFailureTrigger(ctx context.Context, groupID string) (string, error) {
	name := fixtureTriggerName("f_fixture_receipt", groupID)
	if _, err := c.db.ExecContext(ctx, "DROP TRIGGER IF EXISTS "+quoteFixtureIdentifier(name)); err != nil {
		return "", err
	}
	//nolint:gosec // fixture SQL is limited to quoted trigger names and one quoted group literal.
	statement := fmt.Sprintf(`CREATE TRIGGER %s
BEFORE INSERT ON "wg_delivery_receipts"
WHEN EXISTS (SELECT 1 FROM "wg_delivery_outbox" WHERE "effect_key" = NEW."effect_key" AND "group_id" = %s)
BEGIN SELECT RAISE(ABORT, 'fixture local receipt failure'); END`, quoteFixtureIdentifier(name), quoteFixtureLiteral(groupID))
	if _, err := c.db.ExecContext(ctx, statement); err != nil {
		return "", err
	}
	return name, nil
}

func (c *fixtureFaultController) createClosureHoldTrigger(ctx context.Context, groupID string) (string, error) {
	name := fixtureTriggerName("f_fixture_closure", groupID)
	if _, err := c.db.ExecContext(ctx, "DROP TRIGGER IF EXISTS "+quoteFixtureIdentifier(name)); err != nil {
		return "", err
	}
	//nolint:gosec // fixture SQL is limited to quoted trigger/function names and one quoted group literal.
	statement := fmt.Sprintf(`CREATE TRIGGER %s
BEFORE INSERT ON "wg_delivery_checkpoints"
WHEN NEW."group_id" = %s
BEGIN SELECT %s(NEW."group_id"); END`, quoteFixtureIdentifier(name), quoteFixtureLiteral(groupID), quoteFixtureIdentifier(fixtureClosureFunctionName))
	if _, err := c.db.ExecContext(ctx, statement); err != nil {
		return "", err
	}
	return name, nil
}

func (c *fixtureFaultController) dropTrigger(ctx context.Context, name string) error {
	if strings.TrimSpace(name) == "" {
		return nil
	}
	_, err := c.db.ExecContext(ctx, "DROP TRIGGER IF EXISTS "+quoteFixtureIdentifier(name))
	return err
}

func fixtureFaultKindSupported(kind string) bool {
	switch kind {
	case string(dbtarget.FixtureDestinationFaultTargetCommitResponseLost), string(dbtarget.FixtureDestinationFaultTargetCommitHold), fixtureFaultLocalReceipt, fixtureFaultClosureHold:
		return true
	default:
		return false
	}
}

func isTargetFaultKind(kind string) bool {
	return kind == string(dbtarget.FixtureDestinationFaultTargetCommitResponseLost) || kind == string(dbtarget.FixtureDestinationFaultTargetCommitHold)
}

func fixtureTriggerName(prefix, groupID string) string {
	digest := sha256.Sum256([]byte(groupID))
	return prefix + "_" + hex.EncodeToString(digest[:8])
}

func quoteFixtureIdentifier(value string) string {
	return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
}

func quoteFixtureLiteral(value string) string {
	return `'` + strings.ReplaceAll(value, `'`, `''`) + `'`
}

func installFixtureClosureGate(gate *fixtureClosureGate) {
	fixtureClosureRegistry.Lock()
	old := fixtureClosureRegistry.gates[gate.groupID]
	fixtureClosureRegistry.gates[gate.groupID] = gate
	fixtureClosureRegistry.Unlock()
	if old != nil {
		old.releaseGate()
	}
}

func removeFixtureClosureGate(groupID string, expected *fixtureClosureGate) {
	fixtureClosureRegistry.Lock()
	if fixtureClosureRegistry.gates[groupID] == expected {
		delete(fixtureClosureRegistry.gates, groupID)
	}
	fixtureClosureRegistry.Unlock()
}

func releaseFixtureClosureGate(groupID string) {
	fixtureClosureRegistry.Lock()
	gate := fixtureClosureRegistry.gates[groupID]
	delete(fixtureClosureRegistry.gates, groupID)
	fixtureClosureRegistry.Unlock()
	if gate != nil {
		gate.releaseGate()
	}
}

func (g *fixtureClosureGate) releaseGate() { g.once.Do(func() { close(g.release) }) }

func fixtureClosureScalar(_ *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("fixture closure function argument error")
	}
	var groupID string
	switch value := args[0].(type) {
	case string:
		groupID = value
	case []byte:
		groupID = string(value)
	default:
		return nil, errors.New("fixture closure function group error")
	}
	fixtureClosureRegistry.RLock()
	gate := fixtureClosureRegistry.gates[groupID]
	fixtureClosureRegistry.RUnlock()
	if gate == nil {
		return nil, nil
	}
	gate.reached.Add(1)
	gate.holding.Store(true)
	<-gate.release
	gate.holding.Store(false)
	return nil, nil
}

func fixtureClosureRegistrationError() error { return errFixtureClosureRegistration }
