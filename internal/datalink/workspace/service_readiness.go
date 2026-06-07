package workspace

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"
)

var ErrReadinessUnavailable = errors.New("workspace readiness unavailable")

type readinessDeviceService interface {
	GetByID(ctx context.Context, id string) (*schema.Device, error)
	CheckReadiness(ctx context.Context, id string) (*schema.DeviceReadiness, error)
}

type readinessSourceRuleService interface {
	ListByDeviceIDs(ctx context.Context, deviceIDs []string) ([]*schema.SourceRule, error)
	ListLinks(ctx context.Context, ruleID string) ([]*schema.SourceRuleLink, error)
}

type readinessConnectorService interface {
	GetByID(ctx context.Context, id string) (*schema.DatabaseConnector, error)
}

type readinessDBTargetService interface {
	List(ctx context.Context, filter dbtarget.TargetMappingListFilter) ([]*schema.DatabaseTargetMapping, error)
}

// ReadinessSeverity identifies whether a readiness issue blocks activation or only warns the operator.
type ReadinessSeverity string

const (
	ReadinessSeverityBlocking ReadinessSeverity = "blocking"
	ReadinessSeverityWarning  ReadinessSeverity = "warning"
)

// ReadinessStep identifies the Studio V2 step that owns a readiness issue.
type ReadinessStep string

const (
	ReadinessStep1 ReadinessStep = "Step 1"
	ReadinessStep2 ReadinessStep = "Step 2"
	ReadinessStep3 ReadinessStep = "Step 3"
	ReadinessStep4 ReadinessStep = "Step 4"
)

// ReadinessIssue is one normalized readiness finding for the persisted workspace snapshot.
type ReadinessIssue struct {
	Code     string            `json:"code"`
	Severity ReadinessSeverity `json:"severity"`
	Step     ReadinessStep     `json:"step"`
	Scope    string            `json:"scope"`
	Message  string            `json:"message"`
}

// ReadinessSummary aggregates normalized readiness issues for the persisted workspace snapshot.
type ReadinessSummary struct {
	Ready         bool             `json:"ready"`
	BlockingCount int              `json:"blocking_count"`
	WarningCount  int              `json:"warning_count"`
	Issues        []ReadinessIssue `json:"issues"`
}

// ReadinessBlockedError reports that a requested workspace operation is blocked by readiness issues.
type ReadinessBlockedError struct {
	Operation string
	Summary   *ReadinessSummary
}

// Error implements the error interface.
func (e *ReadinessBlockedError) Error() string {
	operation := strings.TrimSpace(e.Operation)
	if operation == "" {
		operation = "operation"
	}
	return fmt.Sprintf("workspace readiness blocked %s", operation)
}

// BlockingIssues returns only the blocking issues from the readiness summary.
func (e *ReadinessBlockedError) BlockingIssues() []ReadinessIssue {
	if e == nil || e.Summary == nil || len(e.Summary.Issues) == 0 {
		return nil
	}
	issues := make([]ReadinessIssue, 0, e.Summary.BlockingCount)
	for _, issue := range e.Summary.Issues {
		if issue.Severity == ReadinessSeverityBlocking {
			issues = append(issues, issue)
		}
	}
	return issues
}

// WithReadinessServices configures persisted workspace readers used by readiness evaluation.
func (s *Service) WithReadinessServices(deviceSvc readinessDeviceService, ruleSvc readinessSourceRuleService, connectorSvc readinessConnectorService, dbTargetSvc readinessDBTargetService) *Service {
	s.readinessDevices = deviceSvc
	s.readinessRules = ruleSvc
	s.readinessConnectors = connectorSvc
	s.readinessMappings = dbTargetSvc
	return s
}

// Readiness evaluates the persisted workspace snapshot into one normalized readiness summary.
func (s *Service) Readiness(ctx context.Context) (*ReadinessSummary, error) {
	if s.readinessDevices == nil {
		return nil, ErrReadinessUnavailable
	}

	record, err := s.GetOrCreate(ctx)
	if err != nil {
		return nil, err
	}

	issues := make([]ReadinessIssue, 0, len(record.OrderedDeviceIDs)+1)
	if len(record.OrderedDeviceIDs) == 0 {
		issues = append(issues, ReadinessIssue{
			Code:     "workspace-device-missing",
			Severity: ReadinessSeverityBlocking,
			Step:     ReadinessStep1,
			Scope:    "workspace",
			Message:  "workspace has no persisted devices",
		})
	}
	for _, deviceID := range record.OrderedDeviceIDs {
		readiness, err := s.readinessDevices.CheckReadiness(ctx, deviceID)
		if err != nil {
			return nil, fmt.Errorf("evaluate workspace device readiness %s: %w", deviceID, err)
		}
		if issue, ok := deviceReadinessIssue(deviceID, readiness); ok {
			issues = append(issues, issue)
		}
	}

	downstreamIssues, err := s.downstreamReadinessIssues(ctx, record.OrderedDeviceIDs, strings.TrimSpace(record.DatabaseConnectorID))
	if err != nil {
		return nil, err
	}
	issues = append(issues, downstreamIssues...)

	connectorID := strings.TrimSpace(record.DatabaseConnectorID)
	if connectorID != "" {
		if s.readinessConnectors == nil {
			return nil, ErrReadinessUnavailable
		}
		connector, err := s.readinessConnectors.GetByID(ctx, connectorID)
		if err != nil {
			return nil, fmt.Errorf("evaluate workspace database connector readiness %s: %w", connectorID, err)
		}
		if issue, ok := databaseConnectorReadinessIssue(connector); ok {
			issues = append(issues, issue)
		}
	}

	summary := &ReadinessSummary{Issues: issues}
	for _, issue := range issues {
		switch issue.Severity {
		case ReadinessSeverityBlocking:
			summary.BlockingCount++
		case ReadinessSeverityWarning:
			summary.WarningCount++
		}
	}
	summary.Ready = summary.BlockingCount == 0

	return summary, nil
}

func (s *Service) downstreamReadinessIssues(ctx context.Context, deviceIDs []string, connectorID string) ([]ReadinessIssue, error) {
	if s.readinessRules == nil || len(deviceIDs) == 0 {
		return nil, nil
	}

	rules, err := s.readinessRules.ListByDeviceIDs(ctx, deviceIDs)
	if err != nil {
		return nil, fmt.Errorf("list workspace source rules for readiness: %w", err)
	}
	issues := make([]ReadinessIssue, 0, len(rules))
	for _, rule := range rules {
		if rule == nil {
			continue
		}
		links, err := s.readinessRules.ListLinks(ctx, rule.ID)
		if err != nil {
			return nil, fmt.Errorf("list workspace source rule links for readiness %s: %w", rule.ID, err)
		}
		for _, link := range links {
			issue, ok, err := s.downstreamLinkIssue(ctx, connectorID, link)
			if err != nil {
				return nil, err
			}
			if ok {
				issues = append(issues, issue)
			}
		}
	}

	return issues, nil
}

func (s *Service) downstreamLinkIssue(ctx context.Context, connectorID string, link *schema.SourceRuleLink) (ReadinessIssue, bool, error) {
	if link == nil {
		return ReadinessIssue{}, false, nil
	}

	pointID := strings.TrimSpace(link.PointID)
	if pointID == "" {
		return ReadinessIssue{
			Code:     "point-missing",
			Severity: ReadinessSeverityBlocking,
			Step:     ReadinessStep2,
			Scope:    strings.TrimSpace(link.RuleID),
			Message:  "source rule link is missing its derived point",
		}, true, nil
	}

	if link.TagID == nil || strings.TrimSpace(*link.TagID) == "" {
		return ReadinessIssue{
			Code:     "tag-missing",
			Severity: ReadinessSeverityBlocking,
			Step:     ReadinessStep2,
			Scope:    pointID,
			Message:  "derived point is missing its persisted tag",
		}, true, nil
	}

	if link.MappingID == nil || strings.TrimSpace(*link.MappingID) == "" {
		return ReadinessIssue{
			Code:     "mapping-missing",
			Severity: ReadinessSeverityBlocking,
			Step:     ReadinessStep3,
			Scope:    pointID,
			Message:  "derived point is missing its persisted mapping",
		}, true, nil
	}

	if connectorID == "" {
		return ReadinessIssue{}, false, nil
	}
	if s.readinessMappings == nil {
		return ReadinessIssue{}, false, ErrReadinessUnavailable
	}

	tagID := strings.TrimSpace(*link.TagID)
	rows, err := s.readinessMappings.List(ctx, dbtarget.TargetMappingListFilter{
		ConnectorID: &connectorID,
		TagID:       &tagID,
	})
	if err != nil {
		return ReadinessIssue{}, false, fmt.Errorf("list database targets for readiness tag %s: %w", tagID, err)
	}
	if len(rows) == 0 {
		return ReadinessIssue{
			Code:     "database-target-missing",
			Severity: ReadinessSeverityBlocking,
			Step:     ReadinessStep4,
			Scope:    pointID,
			Message:  "derived point is missing its persisted database target",
		}, true, nil
	}

	return ReadinessIssue{}, false, nil
}

func deviceReadinessIssue(deviceID string, readiness *schema.DeviceReadiness) (ReadinessIssue, bool) {
	if readiness == nil || readiness.ActivationAllowed {
		return ReadinessIssue{}, false
	}

	code := "device-activation-blocked"
	switch {
	case readiness.ConnectStatus == schema.ReadinessStageStatusSuccess && readiness.ProbeStatus != schema.ReadinessStageStatusSuccess:
		code = "device-probe-required"
	case readiness.ConnectStatus != schema.ReadinessStageStatusSuccess:
		code = "device-connect-required"
	}

	messageInputs := append([]string{}, readiness.BlockingReasons...)
	messageInputs = append(messageInputs, readiness.AvailabilityReason, "device is not ready for activation")
	message := firstReadinessMessage(messageInputs...)
	return ReadinessIssue{
		Code:     code,
		Severity: ReadinessSeverityBlocking,
		Step:     ReadinessStep1,
		Scope:    strings.TrimSpace(deviceID),
		Message:  message,
	}, true
}

func databaseConnectorReadinessIssue(connector *schema.DatabaseConnector) (ReadinessIssue, bool) {
	if connector == nil || !connector.Enabled {
		return ReadinessIssue{}, false
	}

	code := ""
	switch connector.Status {
	case schema.DatabaseConnectorStatusUnreachable:
		code = "database-connector-unreachable"
	case schema.DatabaseConnectorStatusAuthFailed:
		code = "database-connector-auth-failed"
	case schema.DatabaseConnectorStatusError:
		code = "database-connector-error"
	default:
		return ReadinessIssue{}, false
	}

	message := firstReadinessMessage(connector.LastCheckError, "database connector requires attention")
	return ReadinessIssue{
		Code:     code,
		Severity: ReadinessSeverityWarning,
		Step:     ReadinessStep4,
		Scope:    strings.TrimSpace(connector.ID),
		Message:  message,
	}, true
}

func firstReadinessMessage(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
