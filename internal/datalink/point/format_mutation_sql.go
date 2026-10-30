package point

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sqltx"
)

const formatReceiptPrefix = "point_format_receipt:v1:"

type formatReceipt struct {
	Context     decoding.FormatContext `json:"context"`
	Version     int                    `json:"version"`
	WorkspaceID string                 `json:"workspace_id"`
	RequestID   string                 `json:"request_id"`
	PayloadHash string                 `json:"payload_hash"`
	Result      FormatMutationResult   `json:"result"`
}

// MutateFormat commits the local configuration and replay outcome in one tx.
// External runtime application is deliberately not performed in this transaction.
func (r *SQLRepository) MutateFormat(ctx context.Context, req FormatMutationRequest) (out *FormatMutationResult, resultErr error) {
	if err := validateFormatRequest(req); err != nil {
		return nil, err
	}
	if err := validateFormatPatch(req.Patch); err != nil {
		return nil, err
	}
	payload, _ := json.Marshal(req)
	hash := sha256.Sum256(payload)
	digest := hex.EncodeToString(hash[:])
	receiptKey := formatReceiptKey(req.WorkspaceID, req.RequestID)
	ownedTx, err := sqltx.Begin(ctx, r.db, nil)
	if err != nil {
		return nil, ErrFormatUnavailable
	}
	tx := ownedTx.Tx
	var attempted *FormatMutationResult
	defer func() {
		if err := ownedTx.Rollback(); err != nil && err != sql.ErrTxDone {
			out = nil
			resultErr = formatAttemptFailure(err, attempted)
		}
	}()
	snapshot, err := r.formatSnapshot(ctx, tx, req)
	if err != nil {
		return nil, err
	}
	// Scope is checked even on replay; stale expected revisions are not re-applied.
	var rawReceipt string
	err = tx.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key=?`, receiptKey).Scan(&rawReceipt)
	if err == nil {
		var receipt formatReceipt
		if json.Unmarshal([]byte(rawReceipt), &receipt) != nil || receipt.Version != 1 || receipt.WorkspaceID != req.WorkspaceID || receipt.RequestID != req.RequestID {
			return nil, ErrFormatUnavailable
		}
		if receipt.PayloadHash != digest {
			return nil, ErrFormatConflict
		}
		if !validFormatReceipt(receipt, req.FormatScope) {
			return nil, ErrFormatUnavailable
		}
		receipt.Result.replay = true
		return &receipt.Result, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, ErrFormatUnavailable
	}
	if req.freshGateError != nil {
		return nil, req.freshGateError
	}
	if snapshot.point.RevisionID != *req.ExpectedPointRevision || snapshot.ruleRevision != req.ExpectedRuleRevision {
		return nil, ErrFormatConflict
	}
	wsRevision, rawRevision, err := loadFormatWorkspaceRevision(ctx, tx, req.WorkspaceID)
	if err != nil {
		return nil, err
	}
	if wsRevision.Revision != req.ExpectedWorkspaceRevision {
		return nil, ErrFormatConflict
	}
	if wsRevision.Dirty {
		return nil, ErrFormatUnavailable
	}
	var connection struct {
		DataFormat string `json:"data_format"`
	}
	if json.Unmarshal([]byte(snapshot.connectionConfigRaw), &connection) != nil {
		return nil, ErrFormatUnavailable
	}
	snapshot.defaults.ConnectionFormat = connection.DataFormat
	candidate := clonePoint(snapshot.point)
	if req.Patch != nil {
		if err := applyPointUpdate(candidate, *req.Patch); err != nil {
			return nil, ErrFormatInvalid
		}
	}
	if snapshot.defaults.RuleOwned && (candidate.Address != snapshot.point.Address || candidate.Function != snapshot.point.Function) {
		return nil, ErrFormatInvalid
	}
	candidateValue, err := decoding.WithIntent(*candidate, decoding.IntentEdit{Mode: &req.FormatMode, ConfiguredFormat: &req.ConfiguredFormat})
	if err != nil {
		return nil, ErrFormatInvalid
	}
	candidate = &candidateValue
	candidate.DataFormat = req.ConfiguredFormat
	// A mixed legacy PUT is one guarded update, including its non-format fields.
	if candidate.RevisionID == snapshot.point.RevisionID {
		candidate.RevisionID, err = common.NewUUID()
		if err != nil {
			return nil, ErrFormatUnavailable
		}
	}
	candidate.UpdatedAt = time.Now().UTC()
	resolution, err := decoding.ResolvePoint(*candidate, snapshot.defaults)
	if err != nil {
		return nil, ErrFormatInvalid
	}
	nextWorkspace, err := common.NewUUID()
	if err != nil {
		return nil, ErrFormatUnavailable
	}
	attempted = &FormatMutationResult{Point: *candidate, WorkspaceRevision: nextWorkspace, RuntimeApplyStatus: "pending"}
	if err = writePointConfigRevision(ctx, tx, candidate, *req.ExpectedPointRevision); err != nil {
		return nil, err
	}
	if err = advanceFormatWorkspace(ctx, tx, req.WorkspaceID, rawRevision, nextWorkspace); err != nil {
		return nil, err
	}
	resolution.WorkspaceRevision = nextWorkspace
	// Stored receipts exclude volatile readings/errors and connection configuration.
	candidate.LastValue = nil
	candidate.LastReadAt = nil
	candidate.LastError = ""
	result := FormatMutationResult{InvalidationStatus: "pending", Point: *candidate, Resolution: resolution, WorkspaceRevision: nextWorkspace, RuntimeApplyStatus: "pending"}
	receipt := formatReceipt{Context: frozenFormatContext(snapshot.defaults.Protocol, snapshot.defaults.RuleOwned, resolution), Version: 1, WorkspaceID: req.WorkspaceID, RequestID: req.RequestID, PayloadHash: digest, Result: result}
	encoded, err := json.Marshal(receipt)
	if err != nil {
		return nil, ErrFormatUnavailable
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO system_settings(key,value,updated_at) VALUES(?,?,?)`, receiptKey, string(encoded), time.Now().UTC()); err != nil {
		return nil, ErrFormatUnavailable
	}
	if err = ownedTx.Commit(); err != nil {
		return nil, formatAttemptFailure(err, attempted)
	}
	return &result, nil
}

func validateFormatPatch(patch *UpdatePointRequest) error {
	if patch == nil {
		return nil
	}
	if patch.FormatScope != (FormatScope{}) || patch.DataFormat != nil || patch.FormatMode != nil || patch.ConfiguredFormat != nil {
		return ErrFormatInvalid
	}
	if patch.Mode != nil && *patch.Mode != schema.PointModeReadOnly && *patch.Mode != schema.PointModeReadWrite {
		return ErrFormatInvalid
	}
	return nil
}

func formatReceiptKey(workspaceID, requestID string) string {
	keyBytes, _ := json.Marshal([]string{workspaceID, requestID})
	digest := sha256.Sum256(keyBytes)
	return formatReceiptPrefix + hex.EncodeToString(digest[:])
}

func validFormatReceipt(receipt formatReceipt, scope FormatScope) bool {
	r := receipt.Result
	if receipt.Version != 1 || receipt.WorkspaceID != scope.WorkspaceID || receipt.RequestID != scope.RequestID ||
		r.Point.ID != scope.PointID || r.Point.DeviceID != scope.DeviceID || r.Point.RevisionID == "" || r.WorkspaceRevision == "" ||
		r.Resolution.PointRevision != r.Point.RevisionID || r.Resolution.WorkspaceRevision != r.WorkspaceRevision ||
		(r.InvalidationStatus != "pending" && r.InvalidationStatus != "invalidated") || r.RuntimeApplyStatus != "pending" ||
		r.Point.FormatVersion != schema.PointFormatVersion || (r.Point.FormatMode != schema.PointFormatInherit && r.Point.FormatMode != schema.PointFormatOverride) ||
		r.Point.DataFormat != r.Point.FormatOverride ||
		r.Point.LastValue != nil || r.Point.LastReadAt != nil || r.Point.LastError != "" ||
		receipt.Context.RuleOwned != (scope.RuleID != "") {
		return false
	}
	if receipt.Context.RuleOwned && receipt.Context.RuleRevision == "" {
		return false
	}
	if !receipt.Context.RuleOwned && receipt.Context.RuleRevision != "" {
		return false
	}
	hash, err := hex.DecodeString(receipt.PayloadHash)
	if err != nil || len(hash) != sha256.Size {
		return false
	}
	// Recompute only against the receipt's frozen canonical defaults, never the
	// current connection/rule. A later default edit cannot rewrite an old outcome.
	resolved, err := decoding.ResolvePoint(r.Point, receipt.Context)
	return err == nil && resolved == r.Resolution
}

func frozenFormatContext(protocol schema.ProtocolType, ruleOwned bool, resolution decoding.PointResolution) decoding.FormatContext {
	context := decoding.FormatContext{Protocol: protocol, RuleOwned: ruleOwned, WorkspaceRevision: resolution.WorkspaceRevision, RuleRevision: resolution.RuleRevision, ConnectionRevision: resolution.ConnectionRevision}
	switch resolution.Source {
	case decoding.SourceRule:
		context.RuleFormat = resolution.EffectiveFormat
	case decoding.SourceConnection:
		context.ConnectionFormat = resolution.EffectiveFormat
	}
	return context
}
