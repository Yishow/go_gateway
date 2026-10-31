package point

import (
	"context"
	"database/sql"
	"errors"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sqltx"
	"reflect"
	"time"
)

type formatRecoveryAcceptor interface {
	AcceptFormatRecovery(context.Context, FormatRecoveryRequest, *FormatRecoveryEvidence, func() error, ...FormatRecoverySnapshotValidator) error
}

// AcceptFormatRecovery pins an already-durable source snapshot. Its transaction
// performs no row changes: the real format/Share commits precede this boundary.
// The callback may publish/activate only via a nonblocking lifecycle TryLock.
// Cleanup error forbids success and requires revoking only that callback's own
// token/lease; samples accepted against the pinned known source remain historical.
func (s *Service) AcceptFormatRecovery(ctx context.Context, req FormatRecoveryRequest, evidence *FormatRecoveryEvidence, stage func() error, validators ...FormatRecoverySnapshotValidator) error {
	repo, ok := s.repo.(formatRecoveryAcceptor)
	if !ok {
		return ErrFormatUnavailable
	}
	return repo.AcceptFormatRecovery(ctx, req, evidence, stage, validators...)
}

func (r *SQLRepository) AcceptFormatRecovery(ctx context.Context, req FormatRecoveryRequest, evidence *FormatRecoveryEvidence, stage func() error, validators ...FormatRecoverySnapshotValidator) (resultErr error) {
	if evidence == nil || evidence.TargetInputs == nil || stage == nil || evidence.WorkspaceWitness == "" || evidence.ReceiptFingerprint == "" {
		return ErrFormatInvalid
	}
	if err := validateFormatRecoveryRequest(req); err != nil {
		return err
	}
	owned, err := sqltx.Begin(ctx, r.db, nil)
	if err != nil {
		return ErrFormatUnavailable
	}
	defer func() {
		if err := owned.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			resultErr = formatTransactionFailure(errors.Join(resultErr, err))
		}
	}()
	// Internal gateway source storage is SQLite. A zero-row write obtains its
	// writer reservation without changing configuration or firing UPDATE triggers.
	// Do not pretend this is a portable row-lock guarantee for external targets.
	var sqliteVersion string
	if err = owned.Tx.QueryRowContext(ctx, `SELECT sqlite_version()`).Scan(&sqliteVersion); err != nil || sqliteVersion == "" {
		return ErrFormatUnavailable
	}
	if _, err = owned.Tx.ExecContext(ctx, `UPDATE system_settings SET value=value WHERE 0`); err != nil {
		return formatTransactionFailure(err)
	}
	current, err := r.readFormatRecoveryInTx(ctx, owned.Tx, req)
	if err != nil {
		return err
	}
	if current.TargetInputs.WitnessDigest != evidence.TargetInputs.WitnessDigest || current.Profile != evidence.Profile || current.DeviceStatus != evidence.DeviceStatus || current.RuntimeRequired != evidence.RuntimeRequired || current.WorkspaceWitness != evidence.WorkspaceWitness || current.ReceiptFingerprint != evidence.ReceiptFingerprint || current.Resolution != evidence.Resolution || !sameRecoveryPoint(current.Point, evidence.Point) {
		return ErrFormatConflict
	}
	if err = validateRecoverySnapshot(ctx, owned.Tx, validators); err != nil {
		return err
	}
	// The zero-write reservation remains held through this short publication.
	// The callback must never block waiting on locks held by configuration writers.
	// No SQL/provider I/O may occur under metadata/scheduler locks.
	if err = stage(); err != nil {
		return err
	}
	if err = owned.Commit(); err != nil {
		return formatTransactionFailure(err)
	}
	return nil
}

func sameRecoveryPoint(a, b schema.Point) bool {
	a.LastValue = nil
	a.LastReadAt = nil
	a.LastError = ""
	a.CreatedAt = time.Time{}
	a.UpdatedAt = time.Time{}
	b.LastValue = nil
	b.LastReadAt = nil
	b.LastError = ""
	b.CreatedAt = time.Time{}
	b.UpdatedAt = time.Time{}
	return reflect.DeepEqual(a, b)
}
