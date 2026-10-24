package recordingplan

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// StorageRepairAction is the ledger action of an automatic restore.
const StorageRepairAction = "storage_repair"

// Stable codes of a restore that is refused. They are the codes delivery shows.
const (
	// StorageCodeRepairBlocked: restoring is not allowed or not possible now,
	// and needs a person to repair the destination or confirm the schema again.
	StorageCodeRepairBlocked = "schema_repair_blocked"
	// StorageCodeOwnerMismatch: what holds the table's name is not the managed
	// table this scope owns, or the local evidence of ownership is missing.
	StorageCodeOwnerMismatch = "storage_owner_mismatch"
)

// Stable reasons behind a refused restore. They name facts, never values.
const (
	StorageReasonConfirmationMissing   = "confirmation-missing"
	StorageReasonOwnershipMismatch     = "ownership-mismatch"
	StorageReasonRepairNeedsConfirm    = "repair-needs-confirmation"
	StorageReasonConnectorRevision     = "connector-revision-changed"
	StorageReasonTableOccupied         = "table-occupied"
	StorageReasonRecoveryDisabled      = "recovery-disabled"
	StorageReasonPolicyUnconfirmed     = "recovery-policy-unconfirmed"
	StorageReasonReceiptsMissing       = "receipts-missing"
	StorageReasonTimescaleUnavailable  = "timescale-unavailable"
	StorageReasonTimescaleMoved        = "timescale-extension-moved"
	StorageReasonPermissionDenied      = "permission-denied"
	StorageReasonStatementFailed       = "statement-failed"
	StorageReasonVerificationMismatch  = "verification-mismatch"
	StorageReasonUnknownResult         = "unknown-ddl-result"
	StorageReasonAttemptAbandoned      = "repair_attempt_abandoned"
	storageRestoreMaxAttempts          = 3
	storageRepairNextActionConfirm     = "check the destination, then preview and confirm the schema again"
	storageRepairNextActionPermission  = "grant permission to create tables in the target schema, then retry the blocked rows"
	storageRepairNextActionOccupied    = "move or remove the table that holds the managed name, then preview and confirm the schema again"
	storageRepairNextActionReceipts    = "restore the receipt table from its backup, or preview and confirm the schema again"
	storageRepairNextActionReconnect   = "confirm the schema again for the current connector revision"
	storageRepairNextActionTimescale   = "restore the TimescaleDB extension, then retry the blocked rows"
	storageRepairNextActionStatement   = "fix the reported cause, then retry the blocked rows"
	storageRepairNextActionUnconfirmed = "confirm the schema again so the table is checked"
)

// ErrStorageRepairUnsupported means the repository cannot run storage repairs.
var ErrStorageRepairUnsupported = errors.New("storage repair is not supported by this repository")

// StorageRestoreOutcome says what a restore attempt came to.
type StorageRestoreOutcome string

const (
	// StorageRestored: the confirmed table was missing and is back, verified.
	StorageRestored StorageRestoreOutcome = "restored"
	// StorageNotMissing: the catalog shows a table that matches the confirmed
	// layout; nothing was changed and the failure was not about a missing table.
	StorageNotMissing StorageRestoreOutcome = "not_missing"
	// StorageInProgress: another execution is restoring this scope.
	StorageInProgress StorageRestoreOutcome = "in_progress"
	// StorageBlocked: restoring is refused or failed; Code and Reason say why.
	StorageBlocked StorageRestoreOutcome = "blocked"
)

// StorageRestoreTarget is the opened destination a restore inspects and changes.
type StorageRestoreTarget struct {
	Inspect TargetInspector
	Execute SchemaExecutor
	// CanCreate checks, without changing anything, that the connected role may
	// create the table. Nil means the statement itself is the check.
	CanCreate func(ctx context.Context) (bool, error)
}

// StorageDenial is a reason the frozen group itself gives why its table must
// not be restored (the policy is off or was never confirmed, the group's own
// confirmation is missing, the frozen destination does not match the item). It
// only matters once the catalog shows the table is really missing.
type StorageDenial struct{ Code, Reason string }

// StorageRestoreRequest names the scope whose table delivery found missing and
// what the frozen group expects of it.
type StorageRestoreRequest struct {
	Scope             StorageScope
	Dialect           string
	ConnectorRevision string
	// TableName and OwnerColumn are what the frozen group derives for the scope.
	TableName, OwnerColumn string
	// GroupID is the group whose delivery noticed; used for diagnostics only.
	GroupID string
	// Deny, when set, is applied once the catalog proves the table missing.
	Deny *StorageDenial
	// FallbackLayout is the layout the frozen group derives. It is used only to
	// look at the catalog when no confirmation froze a layout; it never
	// authorizes a restore.
	FallbackLayout *GroupSchemaLayout
}

// StorageRestoreResult is the outcome of a restore attempt.
type StorageRestoreResult struct {
	Outcome   StorageRestoreOutcome
	Code      string
	Reason    string
	Operation *SchemaOperation
}

type storageRepairDetail struct {
	Kind         string       `json:"kind"`
	Generation   int          `json:"generation"`
	ToGeneration int          `json:"to_generation"`
	GroupID      string       `json:"group_id,omitempty"`
	ConnectorID  string       `json:"connector_id"`
	Table        string       `json:"table"`
	Scope        StorageScope `json:"scope"`
	Dialect      string       `json:"dialect"`
	OwnerColumn  string       `json:"owner_column"`
}

// RestoreManagedStorage restores the confirmed shared table of a scope when,
// and only when, the destination catalog proves it is missing. Every condition
// comes from local evidence frozen at confirmation (layout, generation,
// connector revision, ownership) and from the catalog now, never from the
// failure that prompted the look. Concurrent attempts, in this or another
// process, meet in the operation ledger: one claims the attempt, the others
// wait for its result. The statements are the confirmed layout's group table,
// its hypertable call and indexes; the receipt table is never created here,
// because an empty one would make effects whose receipts were lost look as if
// they had never been committed.
//
// The attempt is recorded in the operation ledger before any statement runs,
// and ended together with its lineage change in one local transaction. If the
// owner never comes back, the next attempt resolves the operation from the
// catalog alone and never runs a statement it cannot prove is still needed.
func (s *Service) RestoreManagedStorage(ctx context.Context, req StorageRestoreRequest, target StorageRestoreTarget) (StorageRestoreResult, error) {
	repo, ok := s.repo.(storageRepairRepository)
	if !ok {
		return StorageRestoreResult{}, ErrStorageRepairUnsupported
	}
	for range storageRestoreMaxAttempts {
		lineage, err := ReadStorageLineage(ctx, repo.storageDB(), req.Scope)
		if err != nil && !errors.Is(err, ErrStorageLineageNotFound) {
			return StorageRestoreResult{}, err
		}
		if lineage == nil {
			lineage = &StorageLineage{StorageScope: req.Scope, Generation: 1}
		}
		// Look at the catalog first with the confirmed layout (or, with no
		// confirmation, the one the frozen group derives). Whatever the local
		// evidence says only matters once the table is shown to be missing or
		// occupied; a table that is there and right is nobody's concern here.
		layout := lineage.Layout
		if layout == nil {
			layout = req.FallbackLayout
		}
		if layout == nil {
			return s.refuse(ctx, repo, req, lineage, StorageCodeOwnerMismatch, StorageReasonConfirmationMissing)
		}
		inspected, refusal, err := inspectForRestore(ctx, req.Dialect, *layout, target)
		if err != nil {
			return StorageRestoreResult{}, err
		}
		if refusal != nil {
			return s.refuse(ctx, repo, req, lineage, refusal.code, refusal.reason)
		}
		if !inspected.pendingGroup {
			return StorageRestoreResult{Outcome: StorageNotMissing}, nil
		}
		if req.Deny != nil {
			return s.refuse(ctx, repo, req, lineage, req.Deny.Code, req.Deny.Reason)
		}
		if code, reason := lineageRefusal(lineage, req); code != "" {
			return s.refuse(ctx, repo, req, lineage, code, reason)
		}
		if inspected.pendingReceipt {
			return s.refuse(ctx, repo, req, lineage, StorageCodeRepairBlocked, StorageReasonReceiptsMissing)
		}
		if target.CanCreate != nil {
			allowed, err := target.CanCreate(ctx)
			if err != nil {
				return StorageRestoreResult{}, fmt.Errorf("check table creation privilege: %w", err)
			}
			if !allowed {
				return s.refuse(ctx, repo, req, lineage, StorageCodeRepairBlocked, StorageReasonPermissionDenied)
			}
		}
		statements, err := groupSchemaStatements(req.Dialect, *layout, inspected.tables)
		if err != nil {
			return StorageRestoreResult{}, err
		}
		if len(statements) == 0 {
			return StorageRestoreResult{}, errors.New("storage restore produced no statements")
		}

		// A restore that finished while this one was looking moved the
		// generation; what was seen is then out of date, so look again.
		if current, err := ReadStorageLineage(ctx, repo.storageDB(), req.Scope); err != nil {
			return StorageRestoreResult{}, err
		} else if current.Generation != lineage.Generation {
			continue
		}
		operation, outcome, err := s.claimStorageRepair(ctx, repo, req, lineage, *layout)
		if err != nil {
			return StorageRestoreResult{}, err
		}
		switch outcome {
		case ClaimAcquired:
			return s.runStorageRepair(ctx, repo, req, lineage, *layout, operation, statements, target)
		case ClaimInProgress, ClaimScopeBusy:
			if operation != nil && operation.Action == StorageRepairAction && operation.Status.IsActive() && schemaOperationLeaseExpired(operation) {
				if err := s.resolveStaleStorageRepair(ctx, repo, req, operation, target); err != nil {
					return StorageRestoreResult{}, err
				}
				continue
			}
			return StorageRestoreResult{Outcome: StorageInProgress, Operation: operation}, nil
		default:
			// Another execution finished this very attempt between the count and
			// the claim. Read the evidence again; the next attempt gets a new identity.
			continue
		}
	}
	return StorageRestoreResult{Outcome: StorageInProgress}, nil
}

// lineageRefusal checks the local evidence a restore stands on.
func lineageRefusal(lineage *StorageLineage, req StorageRestoreRequest) (code, reason string) {
	switch {
	case !lineage.Confirmed():
		return StorageCodeOwnerMismatch, StorageReasonConfirmationMissing
	case lineage.TableName != req.TableName || lineage.OwnerColumn != req.OwnerColumn ||
		lineage.Layout.TableName != lineage.TableName || lineage.Layout.OwnerColumn != lineage.OwnerColumn:
		return StorageCodeOwnerMismatch, StorageReasonOwnershipMismatch
	case lineage.RepairState != "":
		return StorageCodeRepairBlocked, StorageReasonRepairNeedsConfirm
	case lineage.ConnectorRevision != req.ConnectorRevision:
		return StorageCodeRepairBlocked, StorageReasonConnectorRevision
	}
	return "", ""
}

type restoreRefusal struct{ code, reason string }

// inspectForRestore reads the catalog and turns what makes a restore unsafe
// into a refusal. An inspection that could not be confirmed is an error, not a
// refusal: nothing is known, so nothing is decided.
func inspectForRestore(ctx context.Context, dialect string, layout GroupSchemaLayout, target StorageRestoreTarget) (groupInspection, *restoreRefusal, error) {
	inspected, err := inspectGroupTables(ctx, dialect, layout, target.Inspect)
	var unavailable *TimescaleUnavailableError
	switch {
	case err == nil:
		return inspected, nil, nil
	case errors.As(err, &unavailable):
		return groupInspection{}, &restoreRefusal{StorageCodeRepairBlocked, StorageReasonTimescaleUnavailable}, nil
	case errors.Is(err, ErrPreviewTokenStale):
		return groupInspection{}, &restoreRefusal{StorageCodeRepairBlocked, StorageReasonTimescaleMoved}, nil
	case errors.Is(err, ErrIncompatibleExistingTable):
		return groupInspection{}, &restoreRefusal{StorageCodeOwnerMismatch, StorageReasonTableOccupied}, nil
	}
	return groupInspection{}, nil, err
}

// refuse records the refusal as a durable diagnostic and returns it.
func (s *Service) refuse(ctx context.Context, repo storageRepairRepository, req StorageRestoreRequest, lineage *StorageLineage, code, reason string) (StorageRestoreResult, error) {
	generation := 0
	if lineage != nil {
		generation = lineage.Generation
	}
	if err := RecordStorageRefusal(ctx, repo.storageDB(), req.Scope, generation, req.GroupID, code, reason, time.Now()); err != nil {
		return StorageRestoreResult{}, err
	}
	return StorageRestoreResult{Outcome: StorageBlocked, Code: code, Reason: reason}, nil
}
