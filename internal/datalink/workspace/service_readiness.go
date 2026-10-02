package workspace

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
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

type readinessTagService interface {
	GetByID(ctx context.Context, id string) (*schema.Tag, error)
}

type readinessMappingRecordService interface {
	GetByID(ctx context.Context, id string) (*schema.Mapping, error)
	List(ctx context.Context, filter mapping.ListFilter) ([]*schema.Mapping, error)
}

type readinessResolvedBinding struct {
	tagID        string
	mappingID    string
	tagValid     bool
	mappingValid bool
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

// WithReadinessSetupReaders configures persisted tag and mapping readers used to
// verify source-rule links still point at existing Step 3 rows.
func (s *Service) WithReadinessSetupReaders(tagSvc readinessTagService, mappingSvc readinessMappingRecordService) *Service {
	s.readinessTags = tagSvc
	s.readinessLinkRecords = mappingSvc
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

	connectorID := strings.TrimSpace(record.DatabaseConnectorID)
	effectiveConnectorID := connectorID
	var connectorRecord *schema.DatabaseConnector
	if connectorID != "" {
		if s.readinessConnectors == nil {
			return nil, ErrReadinessUnavailable
		}
		connector, err := s.readinessConnectors.GetByID(ctx, connectorID)
		if err != nil {
			if isMissingDatabaseConnectorError(err) {
				issues = append(issues, ReadinessIssue{
					Code:     "database-connector-missing",
					Severity: ReadinessSeverityBlocking,
					Step:     ReadinessStep4,
					Scope:    connectorID,
					Message:  "workspace database connector binding no longer exists",
				})
				effectiveConnectorID = ""
			} else {
				return nil, fmt.Errorf("evaluate workspace database connector readiness %s: %w", connectorID, err)
			}
		} else {
			connectorRecord = connector
			if issue, ok := databaseConnectorReadinessIssue(connector); ok {
				issues = append(issues, issue)
			}
		}
	}

	coverage, groupIssues, err := s.canonicalWriteGroupReadiness(ctx, record)
	if err != nil {
		return nil, err
	}
	issues = append(issues, groupIssues...)
	downstreamIssues, err := s.downstreamReadinessIssues(ctx, record.OrderedDeviceIDs, effectiveConnectorID, coverage)
	if err != nil {
		return nil, err
	}
	issues = append(issues, downstreamIssues...)
	rowGroupIssues, err := s.databaseRowGroupReadinessIssues(ctx, record, connectorRecord, coverage)
	if err != nil {
		return nil, err
	}
	issues = append(issues, rowGroupIssues...)

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

func (s *Service) resolveReadinessBinding(ctx context.Context, link *schema.SourceRuleLink) (readinessResolvedBinding, error) {
	binding := readinessResolvedBinding{
		tagID:        trimOptionalString(link.TagID),
		mappingID:    trimOptionalString(link.MappingID),
		tagValid:     trimOptionalString(link.TagID) != "",
		mappingValid: trimOptionalString(link.MappingID) != "",
	}
	pointID := strings.TrimSpace(link.PointID)
	if pointID == "" || s.readinessTags == nil || s.readinessLinkRecords == nil {
		return binding, nil
	}
	binding.tagValid = false
	binding.mappingValid = false

	if mappingRecord, tagRecord, ok, err := s.resolveReadinessDirectMapping(ctx, link); err != nil {
		return readinessResolvedBinding{}, err
	} else if ok {
		binding.mappingID = mappingRecord.ID
		binding.tagID = tagRecord.ID
		binding.mappingValid = true
		binding.tagValid = true
		return binding, nil
	}

	if mappingRecord, tagRecord, ok, err := s.resolveReadinessRecoveredMapping(ctx, link, pointID); err != nil {
		return readinessResolvedBinding{}, err
	} else if ok {
		binding.mappingID = mappingRecord.ID
		binding.tagID = tagRecord.ID
		binding.mappingValid = true
		binding.tagValid = true
		return binding, nil
	}

	if tagRecord, ok, err := s.resolveReadinessDirectTag(ctx, link); err != nil {
		return readinessResolvedBinding{}, err
	} else if ok {
		binding.tagID = tagRecord.ID
		binding.tagValid = true
	}

	return binding, nil
}

func (s *Service) resolveReadinessDirectMapping(ctx context.Context, link *schema.SourceRuleLink) (*schema.Mapping, *schema.Tag, bool, error) {
	mappingID := trimOptionalString(link.MappingID)
	if mappingID == "" {
		return nil, nil, false, nil
	}
	mappingRecord, err := s.readinessLinkRecords.GetByID(ctx, mappingID)
	if err != nil {
		if errors.Is(err, mapping.ErrMappingNotFound) {
			return nil, nil, false, nil
		}
		return nil, nil, false, fmt.Errorf("load readiness mapping %s: %w", mappingID, err)
	}
	return s.matchReadinessMapping(ctx, link, mappingRecord)
}

func (s *Service) resolveReadinessRecoveredMapping(ctx context.Context, link *schema.SourceRuleLink, pointID string) (*schema.Mapping, *schema.Tag, bool, error) {
	records, err := s.readinessLinkRecords.List(ctx, mapping.ListFilter{PointID: &pointID})
	if err != nil {
		return nil, nil, false, fmt.Errorf("list readiness mappings for point %s: %w", pointID, err)
	}
	matches := make([]struct {
		mapping *schema.Mapping
		tag     *schema.Tag
	}, 0, len(records))
	for _, record := range records {
		mappingRecord, tagRecord, ok, matchErr := s.matchReadinessMapping(ctx, link, record)
		if matchErr != nil {
			return nil, nil, false, matchErr
		}
		if ok {
			matches = append(matches, struct {
				mapping *schema.Mapping
				tag     *schema.Tag
			}{mapping: mappingRecord, tag: tagRecord})
		}
	}
	if len(matches) != 1 {
		return nil, nil, false, nil
	}
	return matches[0].mapping, matches[0].tag, true, nil
}

func (s *Service) matchReadinessMapping(ctx context.Context, link *schema.SourceRuleLink, mappingRecord *schema.Mapping) (*schema.Mapping, *schema.Tag, bool, error) {
	if mappingRecord == nil {
		return nil, nil, false, nil
	}
	if strings.TrimSpace(mappingRecord.PointID) != strings.TrimSpace(link.PointID) {
		return nil, nil, false, nil
	}
	tagRecord, err := s.readinessTags.GetByID(ctx, mappingRecord.TagID)
	if err != nil {
		if errors.Is(err, tag.ErrTagNotFound) {
			return nil, nil, false, nil
		}
		return nil, nil, false, fmt.Errorf("load readiness tag %s: %w", mappingRecord.TagID, err)
	}
	if !isReadinessRuleManagedTagOwnedBy(tagRecord, link.RuleID, link.Address) && !isLegacyReadinessDirectBinding(link, mappingRecord, tagRecord) {
		return nil, nil, false, nil
	}
	return mappingRecord, tagRecord, true, nil
}

func (s *Service) resolveReadinessDirectTag(ctx context.Context, link *schema.SourceRuleLink) (*schema.Tag, bool, error) {
	tagID := trimOptionalString(link.TagID)
	if tagID == "" {
		return nil, false, nil
	}
	tagRecord, err := s.readinessTags.GetByID(ctx, tagID)
	if err != nil {
		if errors.Is(err, tag.ErrTagNotFound) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("load readiness tag %s: %w", tagID, err)
	}
	if !isReadinessRuleManagedTagOwnedBy(tagRecord, link.RuleID, link.Address) {
		return nil, false, nil
	}
	return tagRecord, true, nil
}

func isReadinessRuleManagedTagOwnedBy(tagRecord *schema.Tag, ruleID, address string) bool {
	if tagRecord == nil || strings.TrimSpace(tagRecord.Labels) == "" {
		return false
	}
	var labels map[string]string
	if err := json.Unmarshal([]byte(tagRecord.Labels), &labels); err != nil {
		return false
	}
	return labels["source"] == "source-rule" &&
		labels["source_rule_id"] == strings.TrimSpace(ruleID) &&
		strings.EqualFold(strings.TrimSpace(labels["source_rule_address"]), strings.TrimSpace(address))
}

func isLegacyReadinessDirectBinding(link *schema.SourceRuleLink, mappingRecord *schema.Mapping, tagRecord *schema.Tag) bool {
	if link == nil || mappingRecord == nil || tagRecord == nil {
		return false
	}
	if strings.TrimSpace(tagRecord.Labels) != "" {
		return false
	}
	if link.MappingID == nil || link.TagID == nil {
		return false
	}
	return strings.TrimSpace(*link.MappingID) == strings.TrimSpace(mappingRecord.ID) &&
		strings.TrimSpace(*link.TagID) == strings.TrimSpace(tagRecord.ID) &&
		strings.TrimSpace(mappingRecord.PointID) == strings.TrimSpace(link.PointID) &&
		strings.TrimSpace(mappingRecord.TagID) == strings.TrimSpace(tagRecord.ID)
}

func trimOptionalString(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func isMissingDatabaseConnectorError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "資料庫連接器不存在")
}
