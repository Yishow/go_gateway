package point

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sqltx"
)

var (
	ErrFormatScope       = errors.New("point format is outside the requested workspace or source")
	ErrFormatConflict    = errors.New("point format revision or request conflicts; reload before editing")
	ErrFormatInvalid     = errors.New("point format requires an explicit scope, revisions and supported intent")
	ErrFormatUnavailable = errors.New("point format persistence is unavailable or uncertain; reload operation before retrying")
)

// FormatScope binds the operator's visible snapshot, never an inferred latest revision.
// A pointer distinguishes an explicit empty legacy point revision from omission.
type FormatScope struct {
	WorkspaceID               string  `json:"workspace_id" validate:"required"`
	DeviceID                  string  `json:"device_id" validate:"required"`
	RuleID                    string  `json:"rule_id,omitempty"`
	PointID                   string  `json:"point_id" validate:"required"`
	ExpectedWorkspaceRevision string  `json:"expected_workspace_revision" validate:"required"`
	ExpectedRuleRevision      string  `json:"expected_rule_revision,omitempty"`
	ExpectedPointRevision     *string `json:"expected_point_revision" validate:"required"`
	RequestID                 string  `json:"request_id" validate:"required"`
}

// FormatMutationRequest is shared by the new endpoint and the legacy PUT bridge.
type FormatMutationRequest struct {
	freshGateError error
	FormatScope
	FormatMode       schema.PointFormatMode `json:"format_mode" validate:"required" enums:"inherit,override"`
	ConfiguredFormat string                 `json:"configured_format" validate:"required"`
	Patch            *UpdatePointRequest    `json:"patch,omitempty"`
}

// FormatMutationResult records configuration only; it never claims live application.
type FormatMutationResult struct {
	InvalidationStatus string `json:"invalidation_status"`
	replay             bool
	Point              schema.Point             `json:"point"`
	Resolution         decoding.PointResolution `json:"resolution"`
	WorkspaceRevision  string                   `json:"workspace_revision"`
	RuntimeApplyStatus string                   `json:"runtime_apply_status"`
}

// FormatMutationStore must atomically commit intent, CAS and durable replay receipt.
type FormatMutationStore interface {
	MutateFormat(context.Context, FormatMutationRequest) (*FormatMutationResult, error)
}

// FormatRuleGuard reuses the rule lifecycle's existing synchronization boundary.
type FormatRuleGuard interface {
	AcquireRuleMutation(context.Context, string) (context.Context, func())
}

// SetFormatRuleGuard is startup wiring, not a dynamic configuration operation.
func (s *Service) SetFormatRuleGuard(guard FormatRuleGuard) { s.formatRuleGuard = guard }

var formatRequestID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.:-]{0,127}$`)

func validateFormatRequest(req FormatMutationRequest) error {
	for _, id := range []string{req.WorkspaceID, req.DeviceID, req.PointID, req.ExpectedWorkspaceRevision} {
		if id == "" || len(id) > 256 || strings.TrimSpace(id) != id {
			return ErrFormatInvalid
		}
	}
	if req.ExpectedPointRevision == nil || len(*req.ExpectedPointRevision) > 256 || !formatRequestID.MatchString(req.RequestID) {
		return ErrFormatInvalid
	}
	if len(req.RuleID) > 256 || len(req.ExpectedRuleRevision) > 256 || strings.TrimSpace(req.RuleID) != req.RuleID {
		return ErrFormatInvalid
	}
	if (req.RuleID == "") != (req.ExpectedRuleRevision == "") {
		return ErrFormatInvalid
	}
	if req.FormatMode != schema.PointFormatInherit && req.FormatMode != schema.PointFormatOverride {
		return ErrFormatInvalid
	}
	if req.FormatMode == schema.PointFormatInherit && req.ConfiguredFormat != "" {
		return ErrFormatInvalid
	}
	if req.FormatMode == schema.PointFormatOverride {
		switch req.ConfiguredFormat {
		case "ABCD", "BADC", "CDAB", "DCBA":
		default:
			return ErrFormatInvalid
		}
	}
	return nil
}

// MutateFormat never uses the unguarded internal Update path or activates runtime.
func (s *Service) MutateFormat(ctx context.Context, req FormatMutationRequest) (*FormatMutationResult, error) {
	if err := validateFormatRequest(req); err != nil {
		return nil, err
	}
	store, ok := s.repo.(FormatMutationStore)
	if !ok {
		return nil, ErrFormatUnavailable
	}
	if req.RuleID != "" {
		if s.formatRuleGuard == nil {
			return nil, ErrFormatUnavailable
		}
		var release func()
		ctx, release = s.formatRuleGuard.AcquireRuleMutation(ctx, req.RuleID)
		defer release()
	}
	if s.formatInvalidator == nil {
		return nil, ErrFormatUnavailable
	}
	completion, ok := s.repo.(formatInvalidationStore)
	if !ok {
		return nil, ErrFormatUnavailable
	}
	if gate, ok := s.formatInvalidator.(interface {
		CheckPointFormatScope(context.Context, FormatScope) error
	}); ok {
		req.freshGateError = gate.CheckPointFormatScope(ctx, req.FormatScope)
		if req.freshGateError != nil && !errors.Is(req.freshGateError, ErrFormatScope) && !errors.Is(req.freshGateError, ErrFormatConflict) && !errors.Is(req.freshGateError, ErrFormatInvalid) {
			req.freshGateError = ErrFormatUnavailable
		}
	}
	result, err := store.MutateFormat(ctx, req)
	if err != nil {
		var failed *formatTransactionError
		if (errors.Is(err, sqltx.ErrCommitOutcomeUnknown) || errors.Is(err, sqltx.ErrConnectionUnusable)) && errors.As(err, &failed) && failed.attempted != nil {
			attempted := failed.attempted
			cleanupErr := s.formatInvalidator.InvalidatePointFormat(ctx, FormatInvalidation{WorkspaceID: req.WorkspaceID, ExpectedWorkspaceRevision: req.ExpectedWorkspaceRevision, WorkspaceRevision: attempted.WorkspaceRevision, PointID: attempted.Point.ID, ExpectedPointRevision: *req.ExpectedPointRevision, PointRevision: attempted.Point.RevisionID})
			return nil, formatTransactionFailure(errors.Join(err, cleanupErr))
		}

		return nil, err
	}
	if result.replay {
		result.replay = false
		if result.InvalidationStatus != "invalidated" {
			return nil, ErrFormatUnavailable
		}
		return result, nil
	}
	invalidation := FormatInvalidation{WorkspaceID: req.WorkspaceID, ExpectedWorkspaceRevision: req.ExpectedWorkspaceRevision, WorkspaceRevision: result.WorkspaceRevision, PointID: req.PointID, ExpectedPointRevision: *req.ExpectedPointRevision, PointRevision: result.Point.RevisionID}
	if s.formatInvalidator.InvalidatePointFormat(ctx, invalidation) != nil {
		return nil, ErrFormatUnavailable
	}
	if err := completion.CompleteFormatInvalidation(ctx, req, result); err != nil {
		return nil, formatTransactionFailure(err)
	}
	return result, nil
}
