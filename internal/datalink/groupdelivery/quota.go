package groupdelivery

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"modernc.org/sqlite"
)

// QuotaScope names what a quota decision applies to.
type QuotaScope string

const (
	QuotaScopeGlobal QuotaScope = "global"
	QuotaScopeGroup  QuotaScope = "group"
)

// Quota states.
const (
	QuotaOK      = "ok"
	QuotaWarning = "warning"
	QuotaHard    = "hard_limit"
)

// Reasons attached to a refused intake.
const (
	ReasonQuotaHardLimit = "quota-hard-limit"
	ReasonDiskFull       = "disk-full"
)

// Default thresholds, matching the legacy queue monitor: warn at 80% and
// suspend intake at 95% of the configured maximum.
const (
	defaultWarningRatio = 0.80
	defaultHardRatio    = 0.95

	sqliteFullCode = 13
)

// QuotaConfig is the explicit capacity policy for accepted-but-undelivered
// data. A zero maximum means that scope is not limited; at least one must be set.
type QuotaConfig struct {
	GlobalMaxBytes int64
	GroupMaxBytes  int64
	WarningRatio   float64
	HardRatio      float64
}

// ErrQuotaExceeded matches every refusal caused by capacity.
var ErrQuotaExceeded = errors.New("durable intake refused for capacity")

// ErrInvalidQuota marks a capacity policy that cannot be enforced.
var ErrInvalidQuota = errors.New("durable quota invalid")

// QuotaError says exactly what was refused and why.
type QuotaError struct {
	Reason    string
	Scope     QuotaScope
	GroupID   string
	UsedBytes int64
	MaxBytes  int64
}

func (e *QuotaError) Error() string { return "durable intake refused: " + e.Reason }

// Is makes errors.Is(err, ErrQuotaExceeded) true.
func (e *QuotaError) Is(target error) bool { return target == ErrQuotaExceeded }

// QuotaStatus is the measured state for one group and for the store.
type QuotaStatus struct {
	Configured     bool
	State          string
	Scope          QuotaScope
	GroupID        string
	UsedBytes      int64
	MaxBytes       int64
	IntakeRefused  bool
	LossRiskNotice string
}

// WithQuota returns a store that enforces the capacity policy.
func (s *Store) WithQuota(config QuotaConfig) (*Store, error) {
	if config.GlobalMaxBytes < 0 || config.GroupMaxBytes < 0 || (config.GlobalMaxBytes == 0 && config.GroupMaxBytes == 0) {
		return nil, fmt.Errorf("%w: set a positive global or group maximum", ErrInvalidQuota)
	}
	if config.WarningRatio == 0 {
		config.WarningRatio = defaultWarningRatio
	}
	if config.HardRatio == 0 {
		config.HardRatio = defaultHardRatio
	}
	if config.WarningRatio <= 0 || config.HardRatio > 1 || config.WarningRatio > config.HardRatio {
		return nil, fmt.Errorf("%w: ratios must satisfy 0 < warning <= hard <= 1", ErrInvalidQuota)
	}
	clone := *s
	clone.quota = &config
	return &clone, nil
}

type queryer interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// usageStatement sums the payload bytes of accepted data that has not been
// delivered: unconsumed journal samples and every outbox row that is not
// committed. Finished data is excluded because it can be reclaimed.
func usageStatement(perGroup bool) string {
	samples, outbox := "", ""
	if perGroup {
		samples, outbox = " AND group_id = ?", " AND group_id = ?"
	}
	return `SELECT
		(SELECT COALESCE(SUM(LENGTH(CAST(payload AS BLOB))), 0) FROM wg_delivery_samples WHERE consumed = 0` + samples + `) +
		(SELECT COALESCE(SUM(LENGTH(CAST(payload AS BLOB))), 0) FROM wg_delivery_outbox WHERE state != 'sql_committed'` + outbox + `)`
}

func measureUsage(ctx context.Context, q queryer, groupID string) (int64, error) {
	var used int64
	var err error
	if groupID == "" {
		err = q.QueryRowContext(ctx, usageStatement(false)).Scan(&used)
	} else {
		err = q.QueryRowContext(ctx, usageStatement(true), groupID, groupID).Scan(&used)
	}
	if err != nil {
		return 0, fmt.Errorf("measure durable usage: %w", err)
	}
	return used, nil
}

// Usage returns the bytes of accepted data that has not been delivered yet,
// store-wide for an empty groupID or for one group.
func (s *Store) Usage(ctx context.Context, groupID string) (int64, error) {
	return measureUsage(ctx, s.db, groupID)
}

// scopeState classifies one scope's usage. Intake is refused once usage
// reaches the hard threshold, so the status and the refusal always agree.
// Accepted data can overshoot that threshold by at most the one sample that
// crossed it; keep HardRatio below 1 when the maximum must never be exceeded.
func (c QuotaConfig) scopeState(used, limit int64) string {
	switch {
	case limit <= 0:
		return QuotaOK
	case float64(used) >= c.HardRatio*float64(limit):
		return QuotaHard
	case float64(used) >= c.WarningRatio*float64(limit):
		return QuotaWarning
	}
	return QuotaOK
}

// checkQuota refuses new data when a scope is at or beyond its hard threshold.
func (s *Store) checkQuota(ctx context.Context, q queryer, key GroupKey) error {
	if s.quota == nil {
		return nil
	}
	scopes := []struct {
		scope   QuotaScope
		groupID string
		max     int64
	}{
		{QuotaScopeGroup, key.GroupID, s.quota.GroupMaxBytes},
		{QuotaScopeGlobal, "", s.quota.GlobalMaxBytes},
	}
	for _, scope := range scopes {
		if scope.max <= 0 {
			continue
		}
		used, err := s.cachedUsage(ctx, q, scope.groupID)
		if err != nil {
			return err
		}
		if s.quota.scopeState(used, scope.max) == QuotaHard {
			return &QuotaError{Reason: ReasonQuotaHardLimit, Scope: scope.scope, GroupID: scope.groupID, UsedBytes: used, MaxBytes: scope.max}
		}
	}
	return nil
}

// QuotaStatus measures usage now and reports the most constrained scope.
func (s *Store) QuotaStatus(ctx context.Context, key GroupKey) (QuotaStatus, error) {
	if s.quota == nil {
		return QuotaStatus{}, nil
	}
	best := QuotaStatus{Configured: true, State: QuotaOK, Scope: QuotaScopeGlobal}
	rank := map[string]int{QuotaOK: 0, QuotaWarning: 1, QuotaHard: 2}
	for _, scope := range []struct {
		scope   QuotaScope
		groupID string
		max     int64
	}{
		{QuotaScopeGroup, key.GroupID, s.quota.GroupMaxBytes},
		{QuotaScopeGlobal, "", s.quota.GlobalMaxBytes},
	} {
		if scope.max <= 0 {
			continue
		}
		used, err := measureUsage(ctx, s.db, scope.groupID)
		if err != nil {
			return QuotaStatus{}, err
		}
		state := s.quota.scopeState(used, scope.max)
		if rank[state] > rank[best.State] || best.MaxBytes == 0 {
			best = QuotaStatus{
				Configured: true, State: state, Scope: scope.scope, GroupID: scope.groupID,
				UsedBytes: used, MaxBytes: scope.max,
			}
		}
	}
	switch best.State {
	case QuotaHard:
		best.IntakeRefused = true
		best.LossRiskNotice = "new readings for this scope are not recorded until space is available; already accepted data is kept"
	case QuotaWarning:
		best.LossRiskNotice = "capacity is running low; intake will be refused at the hard limit"
	}
	return best, nil
}

// isDiskFull recognizes SQLite's out-of-space result.
func isDiskFull(err error) bool {
	var liteErr *sqlite.Error
	return errors.As(err, &liteErr) && liteErr.Code()&0xff == sqliteFullCode
}

func diskFullError(used int64) *QuotaError {
	return &QuotaError{Reason: ReasonDiskFull, Scope: QuotaScopeGlobal, UsedBytes: used}
}

// Reclaim deletes only data that is finished: consumed journal samples,
// committed or operator-skipped outbox rows and per-bucket status history older
// than the retention. It never
// touches pending, sending, retrying, blocked, unknown or quarantined data.
func (s *Store) Reclaim(ctx context.Context, retention time.Duration) (int64, error) {
	cutoff := timestamp(s.now().Add(-retention))
	var removed int64
	err := s.inTx(ctx, "reclaim", func(tx *sql.Tx) error {
		samples, err := tx.ExecContext(ctx, `DELETE FROM wg_delivery_samples WHERE consumed = 1`)
		if err != nil {
			return fmt.Errorf("reclaim consumed samples: %w", err)
		}
		outbox, err := tx.ExecContext(ctx, `
			DELETE FROM wg_delivery_outbox
			WHERE (state = 'sql_committed' AND committed_at != '' AND committed_at <= ?)
			   OR (state = 'operator_skipped' AND updated_at <= ?)`, cutoff, cutoff)
		if err != nil {
			return fmt.Errorf("reclaim finished outbox rows: %w", err)
		}
		buckets, err := tx.ExecContext(ctx, `DELETE FROM wg_delivery_buckets WHERE created_at <= ?`, cutoff)
		if err != nil {
			return fmt.Errorf("reclaim old bucket status: %w", err)
		}
		for _, result := range []sql.Result{samples, outbox, buckets} {
			if n, err := result.RowsAffected(); err == nil {
				removed += n
			}
		}
		return nil
	})
	s.usage.invalidate()
	return removed, err
}

// QuotaStatusView is the API shape of a QuotaStatus.
type QuotaStatusView struct {
	Configured     bool       `json:"configured"`
	State          string     `json:"state"`
	Scope          QuotaScope `json:"scope"`
	UsedBytes      int64      `json:"used_bytes"`
	MaxBytes       int64      `json:"max_bytes"`
	IntakeRefused  bool       `json:"intake_refused"`
	LossRiskNotice string     `json:"loss_risk_notice,omitempty"`
}

// View converts the measured status for API use.
func (q QuotaStatus) View() QuotaStatusView {
	state := q.State
	if !q.Configured {
		state = "unconfigured"
	}
	return QuotaStatusView{
		Configured: q.Configured, State: state, Scope: q.Scope, UsedBytes: q.UsedBytes, MaxBytes: q.MaxBytes,
		IntakeRefused: q.IntakeRefused, LossRiskNotice: q.LossRiskNotice,
	}
}

// usageCache keeps recent usage measurements so intake does not rescan the
// whole backlog for every accepted sample. Accepted payload sizes are added as
// they arrive, so it never under-counts accepted data; freed space (delivery,
// reclaim) is noticed at the next refresh, which makes it conservative.
type usageCache struct {
	mu           sync.Mutex
	ttl          time.Duration
	entries      map[string]usageEntry
	measurements atomic.Int64
}

type usageEntry struct {
	bytes int64
	at    time.Time
}

const defaultUsageTTL = time.Second

func newUsageCache() *usageCache {
	return &usageCache{ttl: defaultUsageTTL, entries: make(map[string]usageEntry)}
}

func (c *usageCache) get(groupID string, now time.Time) (int64, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.entries[groupID]
	if !ok || now.Sub(entry.at) >= c.ttl {
		return 0, false
	}
	return entry.bytes, true
}

func (c *usageCache) put(groupID string, bytes int64, now time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[groupID] = usageEntry{bytes: bytes, at: now}
}

// add accounts for newly accepted bytes in every cached scope.
func (c *usageCache) add(groupID string, bytes int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, key := range []string{"", groupID} {
		if entry, ok := c.entries[key]; ok {
			entry.bytes += bytes
			c.entries[key] = entry
		}
	}
}

func (c *usageCache) invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]usageEntry)
}

// cachedUsage returns a recent measurement or measures and remembers one.
func (s *Store) cachedUsage(ctx context.Context, q queryer, groupID string) (int64, error) {
	now := s.now()
	if used, ok := s.usage.get(groupID, now); ok {
		return used, nil
	}
	used, err := measureUsage(ctx, q, groupID)
	if err != nil {
		return 0, err
	}
	s.usage.measurements.Add(1)
	s.usage.put(groupID, used, now)
	return used, nil
}
