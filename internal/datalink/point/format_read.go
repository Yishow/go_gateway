package point

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/sqltx"
)

// FormatReadScope names the operator's source scope, not database setup scope.
type FormatReadScope struct {
	WorkspaceID string `json:"workspace_id"`
	DeviceID    string `json:"device_id"`
	PointID     string `json:"point_id"`
	RuleID      string `json:"rule_id"`
}

// FormatReadResult is a whitelist: no device credentials or volatile readings.
// An explicitly empty PointRevision is an observed legacy revision, not missing proof.
type FormatReadResult struct {
	FormatReadScope
	WorkspaceRevision       string                   `json:"workspace_revision"`
	WorkspaceRevisionOrigin string                   `json:"workspace_revision_origin"`
	RuleRevision            string                   `json:"rule_revision"`
	PointRevision           string                   `json:"point_revision"`
	ConnectionRevision      string                   `json:"connection_revision"`
	Editable                bool                     `json:"editable"`
	Resolution              decoding.PointResolution `json:"resolution"`
}

type formatReader interface {
	ReadFormat(context.Context, FormatReadScope) (*FormatReadResult, error)
}

func (s *Service) ReadFormat(ctx context.Context, scope FormatReadScope) (*FormatReadResult, error) {
	reader, ok := s.repo.(formatReader)
	if !ok {
		return nil, ErrFormatUnavailable
	}
	result, err := reader.ReadFormat(ctx, scope)
	if err != nil {
		return nil, err
	}
	editable := false
	if gate, ok := s.formatInvalidator.(interface {
		CheckPointFormatScope(context.Context, FormatScope) error
	}); ok && result.Editable {
		editable = gate.CheckPointFormatScope(ctx, FormatScope{WorkspaceID: result.WorkspaceID, DeviceID: result.DeviceID, PointID: result.PointID, RuleID: result.RuleID, ExpectedWorkspaceRevision: result.WorkspaceRevision, ExpectedRuleRevision: result.RuleRevision, ExpectedPointRevision: &result.PointRevision}) == nil
	}
	result.Editable = editable
	return result, nil
}

// ReadFormat uses one read snapshot and never initializes revision or receipt state.
func (r *SQLRepository) ReadFormat(ctx context.Context, scope FormatReadScope) (out *FormatReadResult, resultErr error) {
	for _, v := range []string{scope.WorkspaceID, scope.DeviceID, scope.PointID} {
		if v == "" || len(v) > 256 || strings.TrimSpace(v) != v {
			return nil, ErrFormatInvalid
		}
	}
	if len(scope.RuleID) > 256 || strings.TrimSpace(scope.RuleID) != scope.RuleID {
		return nil, ErrFormatInvalid
	}
	owned, err := sqltx.Begin(ctx, r.db, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, ErrFormatUnavailable
	}
	defer func() {
		if err := owned.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			out = nil
			resultErr = ErrFormatUnavailable
		}
	}()
	if scope.RuleID == "" {
		err = owned.Tx.QueryRowContext(ctx, `SELECT rule_id FROM source_rule_links WHERE point_id=?`, scope.PointID).Scan(&scope.RuleID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFormatUnavailable
		}
	}
	snapshot, err := r.formatSnapshot(ctx, owned.Tx, FormatMutationRequest{FormatScope: FormatScope{WorkspaceID: scope.WorkspaceID, DeviceID: scope.DeviceID, PointID: scope.PointID, RuleID: scope.RuleID}})
	if err != nil {
		return nil, err
	}
	revision, raw, err := loadFormatWorkspaceRevision(ctx, owned.Tx, scope.WorkspaceID)
	if err != nil {
		return nil, err
	}
	var config struct {
		DataFormat string `json:"data_format"`
	}
	if json.Unmarshal([]byte(snapshot.connectionConfigRaw), &config) != nil {
		return nil, ErrFormatUnavailable
	}
	snapshot.defaults.ConnectionFormat = config.DataFormat
	snapshot.defaults.ConnectionRevision = connector.ReadConnectionRevision(snapshot.connectionConfigRaw)
	snapshot.defaults.WorkspaceRevision = revision.Revision
	snapshot.defaults.RuleID = scope.RuleID
	resolved, err := decoding.ResolvePoint(*snapshot.point, snapshot.defaults)
	if err != nil {
		return nil, ErrFormatInvalid
	}
	origin := "persisted"
	if raw == "" {
		origin = "initial"
	}
	return &FormatReadResult{FormatReadScope: scope, WorkspaceRevision: revision.Revision, WorkspaceRevisionOrigin: origin, RuleRevision: snapshot.ruleRevision, PointRevision: snapshot.point.RevisionID, ConnectionRevision: resolved.ConnectionRevision, Editable: !revision.Dirty, Resolution: resolved}, nil
}
