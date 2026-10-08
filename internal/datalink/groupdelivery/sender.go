package groupdelivery

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/delivery"
	"go-gateway/internal/datalink/schema"
)

// Target is an opened destination connection.
type Target struct {
	DB   *sql.DB
	Kind schema.DatabaseConnectorKind
	// Close releases the connection once the attempt is finished; optional.
	Close func() error
}

// ErrTargetBlocked marks a destination that cannot be used until repaired
// (revision changed, credential invalid, identity mismatch).
var ErrTargetBlocked = errors.New("destination blocked")

// ErrLocalReceiptFailed means the destination committed but the local state
// update did not; the item is left for interrupted-delivery resolution.
var ErrLocalReceiptFailed = errors.New("local delivery receipt failed")

// TargetResolver opens the destination bound to an item's frozen identity.
type TargetResolver interface {
	Resolve(ctx context.Context, item OutboxItem) (Target, error)
}

// Safe error codes stored with an item.
// Defaults for the per-attempt bounds. They are limits, not performance claims.
const (
	defaultDeliveryTimeout = 30 * time.Second
	defaultSettleTimeout   = 5 * time.Second
	defaultLeaseMargin     = 30 * time.Second
)

const (
	codeInsertFailed       = "insert-failed"
	codeCommitAmbiguous    = "commit-ambiguous"
	codeIdentityConflict   = "destination-identity-conflict"
	codeUnsupported        = "destination-unsupported"
	codeTargetBlocked      = "target-blocked"
	codeRowRejected        = "destination-rejected-row"
	codeTargetUnusable     = "target-unusable"
	codeTargetUnavailable  = "target-unavailable"
	codePayloadMismatch    = "payload-digest-mismatch"
	codePayloadUndecodable = "payload-undecodable"
)

// SenderConfig bounds retries; Backoff is injectable for tests.
type SenderConfig struct {
	// Owner identifies this worker incarnation in claims, as "<node>/<incarnation>".
	Owner string
	// LeaseTTL keeps a claim valid through the destination attempt and local
	// settlement. An undersized value is raised with the existing lease margin.
	LeaseTTL time.Duration
	// DeliveryTimeout bounds one destination attempt.
	DeliveryTimeout time.Duration
	// SettleTimeout bounds the local state update that follows a destination
	// attempt; it is independent of the caller's cancellation so a committed
	// effect is recorded even while the process is shutting down.
	SettleTimeout time.Duration
	// MaxRetries caps attempts of a failing row; zero (the default) means no
	// cap. Transient outages (destination offline, timeouts) must never turn a
	// whole backlog into `blocked` that nothing releases, so production retries
	// with capped backoff and reserves `blocked` for causes that need repair.
	MaxRetries int
	// Backoff returns the next state (retrying or blocked) and due time for the
	// given retry count. Nil uses delivery.CalculateBackoff.
	Backoff func(retryCount int) (state string, next time.Time)
	Now     func() time.Time
}

// DeliveryResult is the state an item ended in.
type DeliveryResult struct {
	State string
	Code  string
	// Already is true when the item had been committed before this call.
	Already bool
}

// Sender delivers outbox rows to their frozen destinations.
type Sender struct {
	owner           string
	leaseTTL        time.Duration
	deliveryTimeout time.Duration
	settleTimeout   time.Duration
	store           *Store
	targets         TargetResolver
	config          SenderConfig
	backoff         func(int) (string, time.Time)
	nowValue        func() time.Time
}

// NewSender builds a sender.
func NewSender(store *Store, targets TargetResolver, config SenderConfig) *Sender {
	config.MaxRetries = max(config.MaxRetries, 0)
	backoff := config.Backoff
	if backoff == nil {
		limit := config.MaxRetries
		backoff = func(retry int) (string, time.Time) {
			status, next := delivery.CalculateBackoff(retry, limit)
			return string(status), next
		}
	}
	now := config.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	if config.Owner == "" {
		config.Owner = "gateway/default"
	}
	if config.DeliveryTimeout <= 0 {
		config.DeliveryTimeout = defaultDeliveryTimeout
	}
	if config.SettleTimeout <= 0 {
		config.SettleTimeout = defaultSettleTimeout
	}
	if minimum := config.DeliveryTimeout + config.SettleTimeout; config.LeaseTTL < minimum {
		// A lease shorter than the attempt plus settlement would let another
		// worker take over while the current worker is recording its result.
		config.LeaseTTL = minimum + defaultLeaseMargin
	}
	return &Sender{
		owner: config.Owner, leaseTTL: config.LeaseTTL, deliveryTimeout: config.DeliveryTimeout,
		settleTimeout: config.SettleTimeout, store: store, targets: targets, config: config,
		backoff: backoff, nowValue: now,
	}
}

// Deliver attempts one outbox item. Items that are not pending or retrying are
// never sent: unknown and blocked items wait for an operator, committed items
// are already done.
func (s *Sender) Deliver(ctx context.Context, effectKey string) (DeliveryResult, error) {
	item, err := s.store.GetOutbox(ctx, effectKey)
	if err != nil {
		return DeliveryResult{}, err
	}
	switch item.State {
	case StateCommitted:
		return DeliveryResult{State: StateCommitted, Already: true}, nil
	case StatePending, StateRetrying:
	default:
		return DeliveryResult{State: item.State, Code: item.LastErrorCode}, ErrNotDeliverable
	}
	sum := sha256.Sum256(item.Payload)
	if hex.EncodeToString(sum[:]) != item.PayloadDigest {
		return s.block(ctx, effectKey, codePayloadMismatch)
	}
	row, err := DecodeRowPayload(item.Payload)
	if err != nil {
		return s.block(ctx, effectKey, codePayloadUndecodable)
	}
	claim, err := s.store.BeginDelivery(ctx, effectKey, s.owner, s.leaseTTL)
	if err != nil {
		return DeliveryResult{State: item.State}, err
	}
	// The destination attempt is bounded and follows the caller's cancellation;
	// local settlement starts only after Resolve/InsertGroupRow returns and has
	// its own budget, independent of cancellation during process shutdown.
	attemptCtx, cancelAttempt := context.WithTimeout(ctx, s.deliveryTimeout)
	defer cancelAttempt()

	target, err := s.targets.Resolve(attemptCtx, item)
	if err != nil {
		settleCtx, cancelSettle := context.WithTimeout(context.WithoutCancel(ctx), s.settleTimeout)
		defer cancelSettle()
		if errors.Is(err, ErrTargetBlocked) {
			return s.blockClaimed(settleCtx, claim, codeTargetBlocked)
		}
		return s.retry(settleCtx, claim, codeTargetUnavailable)
	}
	if target.Close != nil {
		defer func() { _ = target.Close() }() //nolint:errcheck // releasing a finished connection; nothing to recover
	}
	_, err = dbtarget.InsertGroupRow(attemptCtx, target.DB, dbtarget.GroupInsertRequest{
		Kind: target.Kind, SchemaName: item.Destination.TableSchema, TableName: item.Destination.TableName,
		Row: row, PayloadDigest: item.PayloadDigest, Strategy: item.Destination.DedupeCapability,
		RecordKeyColumn: item.Destination.RecordKeyColumn, CommittedAt: timestamp(s.nowValue()),
	})
	settleCtx, cancelSettle := context.WithTimeout(context.WithoutCancel(ctx), s.settleTimeout)
	defer cancelSettle()
	return s.settle(settleCtx, claim, item, err)
}

func (s *Sender) settle(ctx context.Context, claim Claim, item OutboxItem, insertErr error) (DeliveryResult, error) {
	var groupErr *dbtarget.GroupInsertError
	switch {
	case insertErr == nil:
		if err := s.store.CompleteDelivery(ctx, claim, item.PayloadDigest); err != nil {
			if errors.Is(err, ErrFenced) {
				return DeliveryResult{State: StateSending}, err
			}
			// The destination committed; leave the item in sending so interrupted
			// delivery resolution, not a blind retry, decides what happens next.
			return DeliveryResult{State: StateSending}, fmt.Errorf("%w: %w", ErrLocalReceiptFailed, err)
		}
		return DeliveryResult{State: StateCommitted}, nil
	case errors.Is(insertErr, dbtarget.ErrReceiptDigestMismatch), errors.Is(insertErr, dbtarget.ErrExistingRowDiffers):
		return s.blockClaimed(ctx, claim, codeIdentityConflict)
	case errors.Is(insertErr, dbtarget.ErrGroupInsertUnsupported):
		return s.blockClaimed(ctx, claim, codeUnsupported)
	case errors.As(insertErr, &groupErr) && groupErr.Ambiguous():
		if item.Destination.DedupeCapability == DedupeReceipt || item.Destination.DedupeCapability == DedupeUniqueKey {
			return s.retry(ctx, claim, codeCommitAmbiguous)
		}
		if err := s.store.MarkUnknown(ctx, claim, codeCommitAmbiguous); err != nil {
			return DeliveryResult{State: StateSending}, err
		}
		return DeliveryResult{State: StateUnknown, Code: codeCommitAmbiguous}, nil
	}
	switch dbtarget.ClassifyInsertError(insertErr) {
	case dbtarget.InsertErrorRow:
		// The destination can never accept this row: isolate it, keep its
		// payload and let only its own partition wait.
		if err := s.store.MarkQuarantined(ctx, claim, codeRowRejected); err != nil {
			return DeliveryResult{State: StateSending}, err
		}
		return DeliveryResult{State: StateQuarantined, Code: codeRowRejected}, nil
	case dbtarget.InsertErrorTarget:
		return s.blockClaimed(ctx, claim, codeTargetUnusable)
	}
	return s.retry(ctx, claim, codeInsertFailed)
}

func (s *Sender) retry(ctx context.Context, claim Claim, code string) (DeliveryResult, error) {
	item, err := s.store.GetOutbox(ctx, claim.EffectKey)
	if err != nil {
		return DeliveryResult{}, err
	}
	state, next := s.backoff(item.RetryCount + 1)
	if err := s.store.MarkRetry(ctx, claim, state, code, next); err != nil {
		return DeliveryResult{State: StateSending}, err
	}
	return DeliveryResult{State: state, Code: code}, nil
}

func (s *Sender) blockClaimed(ctx context.Context, claim Claim, code string) (DeliveryResult, error) {
	if err := s.store.BlockClaimed(ctx, claim, code); err != nil {
		return DeliveryResult{}, err
	}
	return DeliveryResult{State: StateBlocked, Code: code}, nil
}

func (s *Sender) block(ctx context.Context, effectKey, code string) (DeliveryResult, error) {
	if err := s.store.MarkBlocked(ctx, effectKey, code); err != nil {
		return DeliveryResult{}, err
	}
	return DeliveryResult{State: StateBlocked, Code: code}, nil
}
