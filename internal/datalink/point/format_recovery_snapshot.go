package point

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"strings"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sqltx"
)

// FormatRecoveryRequest binds recovery to the original operation's saved result,
// not to the pre-save revisions and never to a newly fetched latest revision.
type FormatRecoveryRequest struct {
	FormatScope
	ExpectedConnectionRevision string `json:"expected_connection_revision"`
	ExpectedSettingsRevision   string `json:"expected_settings_revision"`
}

// FormatRecoveryEvidence is internal evidence only, not permission to reinstall.
// Consumers must still condition installation on their captured runtime proof.
type FormatRecoveryEvidence struct {
	TargetInputs       *FormatRecoveryTargetInputs `json:"-"`
	Profile            connector.ReadProfile       `json:"-"`
	DeviceStatus       schema.DeviceStatus         `json:"-"`
	RuntimeRequired    bool                        `json:"-"`
	Point              schema.Point                `json:"-"`
	Resolution         decoding.PointResolution    `json:"-"`
	WorkspaceWitness   string                      `json:"-"`
	ReceiptFingerprint string                      `json:"-"`
}

type formatRecoveryReader interface {
	ReadFormatRecovery(context.Context, FormatRecoveryRequest, ...FormatRecoverySnapshotValidator) (*FormatRecoveryEvidence, error)
}

// ReadFormatRecovery delegates to the authoritative read-only repository snapshot.
func (s *Service) ReadFormatRecovery(ctx context.Context, req FormatRecoveryRequest, validators ...FormatRecoverySnapshotValidator) (*FormatRecoveryEvidence, error) {
	reader, ok := s.repo.(formatRecoveryReader)
	if !ok {
		return nil, ErrFormatUnavailable
	}
	return reader.ReadFormatRecovery(ctx, req, validators...)
}

// ReadFormatRecovery performs no writes, initialization, callbacks, or rebase.
func (r *SQLRepository) ReadFormatRecovery(ctx context.Context, req FormatRecoveryRequest, validators ...FormatRecoverySnapshotValidator) (out *FormatRecoveryEvidence, resultErr error) {
	if err := validateFormatRecoveryRequest(req); err != nil {
		return nil, err
	}
	owned, err := sqltx.Begin(ctx, r.db, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, ErrFormatUnavailable
	}
	defer func() {
		if err := owned.Rollback(); err != nil && err != sql.ErrTxDone {
			out = nil
			resultErr = formatTransactionFailure(err)
		}
	}()
	out, resultErr = r.readFormatRecoveryInTx(ctx, owned.Tx, req)
	if resultErr == nil {
		resultErr = validateRecoverySnapshot(ctx, owned.Tx, validators)
	}
	if resultErr != nil {
		out = nil
	}
	return
}

func validateFormatRecoveryRequest(req FormatRecoveryRequest) error {
	if err := validateFormatRequest(FormatMutationRequest{FormatScope: req.FormatScope, FormatMode: schema.PointFormatInherit}); err != nil {
		return err
	}
	for _, v := range []string{req.ExpectedConnectionRevision, req.ExpectedSettingsRevision, *req.ExpectedPointRevision} {
		if v == "" || len(v) > 256 || strings.TrimSpace(v) != v {
			return ErrFormatInvalid
		}
	}
	return nil
}

func (r *SQLRepository) readFormatRecoveryInTx(ctx context.Context, tx *sql.Tx, req FormatRecoveryRequest) (*FormatRecoveryEvidence, error) {
	snapshot, err := r.formatSnapshot(ctx, tx, FormatMutationRequest{FormatScope: req.FormatScope})
	if err != nil {
		return nil, err
	}
	revision, workspaceRaw, err := loadFormatWorkspaceRevision(ctx, tx, req.WorkspaceID)
	if err != nil {
		return nil, err
	}
	var cleanliness struct {
		Dirty *bool `json:"dirty"`
	}
	if json.Unmarshal([]byte(workspaceRaw), &cleanliness) != nil || cleanliness.Dirty == nil {
		return nil, ErrFormatUnavailable
	}
	if revision.Dirty || workspaceRaw == "" {
		return nil, ErrFormatUnavailable
	}
	if revision.Revision != req.ExpectedWorkspaceRevision || snapshot.point.RevisionID != *req.ExpectedPointRevision || snapshot.ruleRevision != req.ExpectedRuleRevision || snapshot.defaults.ConnectionRevision != req.ExpectedConnectionRevision {
		return nil, ErrFormatConflict
	}
	var deviceStatus schema.DeviceStatus
	if tx.QueryRowContext(ctx, `SELECT status FROM devices WHERE id=?`, req.DeviceID).Scan(&deviceStatus) != nil {
		return nil, ErrFormatUnavailable
	}
	switch deviceStatus {
	case schema.DeviceStatusActive, schema.DeviceStatusDraft, schema.DeviceStatusDisabled:
	default:
		return nil, ErrFormatUnavailable
	}
	var settingsRaw string
	if tx.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key='modbus_share'`).Scan(&settingsRaw) != nil {
		return nil, ErrFormatUnavailable
	}
	var settings struct {
		Revision string `json:"settings_revision"`
	}
	if json.Unmarshal([]byte(settingsRaw), &settings) != nil || settings.Revision == "" {
		return nil, ErrFormatUnavailable
	}
	if settings.Revision != req.ExpectedSettingsRevision {
		return nil, ErrFormatConflict
	}
	var rawReceipt string
	if tx.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key=?`, formatReceiptKey(req.WorkspaceID, req.RequestID)).Scan(&rawReceipt) != nil {
		return nil, ErrFormatUnavailable
	}
	var receipt formatReceipt
	if json.Unmarshal([]byte(rawReceipt), &receipt) != nil || !validFormatReceipt(receipt, req.FormatScope) || receipt.Result.InvalidationStatus != "invalidated" {
		return nil, ErrFormatUnavailable
	}
	if receipt.Result.Point.RevisionID != *req.ExpectedPointRevision || receipt.Result.WorkspaceRevision != req.ExpectedWorkspaceRevision || !sameRecoveryPoint(receipt.Result.Point, *snapshot.point) {
		return nil, ErrFormatConflict
	}
	var connection struct {
		DataFormat string `json:"data_format"`
	}
	if json.Unmarshal([]byte(snapshot.connectionConfigRaw), &connection) != nil {
		return nil, ErrFormatUnavailable
	}
	snapshot.defaults.ConnectionFormat = connection.DataFormat
	snapshot.defaults.WorkspaceRevision = revision.Revision
	// formatSnapshot has already proven this exact rule owns the point in this tx.
	snapshot.defaults.RuleID = req.RuleID
	profile, err := connector.ResolveReadProfile(*snapshot.point, snapshot.defaults)
	resolved := profile.Resolution()
	if err != nil {
		return nil, ErrFormatInvalid
	}
	if resolved != receipt.Result.Resolution || frozenFormatContext(snapshot.defaults.Protocol, snapshot.defaults.RuleOwned, resolved) != receipt.Context {
		return nil, ErrFormatConflict
	}
	target, err := ReadFormatRecoveryTargetInputs(ctx, tx, snapshot.point)
	if err != nil {
		return nil, err
	}
	digest := sha256.Sum256([]byte(rawReceipt))
	p := clonePoint(snapshot.point)
	p.LastValue = nil
	p.LastReadAt = nil
	p.LastError = ""
	return &FormatRecoveryEvidence{TargetInputs: target, Profile: profile, DeviceStatus: deviceStatus, RuntimeRequired: deviceStatus == schema.DeviceStatusActive && p.Enabled, Point: *p, Resolution: resolved, WorkspaceWitness: workspaceRaw, ReceiptFingerprint: hex.EncodeToString(digest[:])}, nil
}
