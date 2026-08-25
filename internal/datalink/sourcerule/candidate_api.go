package sourcerule

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/schema"
)

// CandidateSnapshotView is the API projection of all candidate groups for a
// source-rule revision.
type CandidateSnapshotView struct {
	SourceRuleID       string           `json:"source_rule_id"`
	WorkspaceID        string           `json:"workspace_id,omitempty"`
	WorkspaceRevision  string           `json:"workspace_revision,omitempty"`
	RevisionID         string           `json:"revision_id"`
	Tags               CandidateSetView `json:"tags"`
	DatabaseOutputs    CandidateSetView `json:"database_outputs"`
	LocalModbusOutputs CandidateSetView `json:"local_modbus_outputs"`
}

// CandidateSetView describes one candidate group and its review state.
type CandidateSetView struct {
	Status     schema.SourceRuleCandidateStatus `json:"status"`
	Reason     string                           `json:"reason,omitempty"`
	Candidates []any                            `json:"candidates"`
}

// TagCandidateView adds persisted mapping review state to a tag candidate.
type TagCandidateView struct {
	schema.SourceRuleTagCandidate
	Status               schema.MappingStatus `json:"status"`
	LastAppliedSignature string               `json:"last_applied_signature,omitempty"`
	BlockingReason       string               `json:"blocking_reason,omitempty"`
}

// GetCandidateView returns the persisted candidate review projection for a rule.
func (s *Service) GetCandidateView(ctx context.Context, ruleID string) (*CandidateSnapshotView, error) {
	rule, err := s.repo.GetByID(ctx, ruleID)
	if err != nil {
		return nil, fmt.Errorf("取得來源規則失敗: %w", err)
	}
	snapshots, err := s.repo.ListCandidateSnapshots(ctx, rule.ID, rule.RevisionID)
	if err != nil {
		return nil, fmt.Errorf("列出來源規則候選快照失敗: %w", err)
	}
	if shouldRecomputeDatabaseCandidateView(snapshots) {
		snapshots, err = s.recomputeDatabaseCandidateView(ctx, rule, snapshots)
		if err != nil {
			return nil, err
		}
	}
	return s.composeCandidateView(ctx, rule, snapshots)
}

// GetCandidateViewAtRevision reads only the revision selected by the caller.
// It rechecks the rule after the snapshot read so a revision transition cannot
// leak a mixed response.
func (s *Service) GetCandidateViewAtRevision(ctx context.Context, ruleID, revisionID string, scope *CandidateScopeRequest) (*CandidateSnapshotView, error) {
	if scope != nil && s.CandidateScopeConfigured() {
		if err := s.ValidateCandidateScope(ctx, ruleID, *scope); err != nil {
			return nil, err
		}
	}
	rule, err := s.repo.GetByID(ctx, ruleID)
	if err != nil {
		return nil, fmt.Errorf("取得來源規則失敗: %w", err)
	}
	if rule.RevisionID != revisionID {
		return nil, candidateRevisionConflict()
	}
	snapshots, err := s.repo.ListCandidateSnapshots(ctx, rule.ID, revisionID)
	if err != nil {
		return nil, fmt.Errorf("列出來源規則候選快照失敗: %w", err)
	}
	if current, getErr := s.repo.GetByID(ctx, ruleID); getErr != nil {
		return nil, getErr
	} else if current.RevisionID != revisionID {
		return nil, candidateRevisionConflict()
	}
	if shouldRecomputeDatabaseCandidateView(snapshots) {
		snapshots, err = s.recomputeDatabaseCandidateView(ctx, rule, snapshots)
		if err != nil {
			return nil, err
		}
	}
	return s.composeCandidateView(ctx, rule, snapshots)
}

// RecomputeCandidateView rebuilds candidate snapshots and returns their review
// projection for a rule.
func (s *Service) RecomputeCandidateView(ctx context.Context, ruleID string) (*CandidateSnapshotView, error) {
	rule, err := s.repo.GetByID(ctx, ruleID)
	if err != nil {
		return nil, fmt.Errorf("取得來源規則失敗: %w", err)
	}
	links, err := s.repo.ListLinks(ctx, rule.ID)
	if err != nil {
		return nil, fmt.Errorf("列出來源規則連結失敗: %w", err)
	}
	if err := s.persistCandidateSnapshots(ctx, rule, links); err != nil {
		return nil, err
	}
	snapshots, err := s.repo.ListCandidateSnapshots(ctx, rule.ID, rule.RevisionID)
	if err != nil {
		return nil, fmt.Errorf("列出來源規則候選快照失敗: %w", err)
	}
	return s.composeCandidateView(ctx, rule, snapshots)
}

// RecomputeCandidateViewAtRevision builds and persists candidates only for the
// selected rule revision. Scope is checked again immediately before the CAS.
func (s *Service) RecomputeCandidateViewAtRevision(ctx context.Context, ruleID, revisionID string, scope *CandidateScopeRequest) (*CandidateSnapshotView, error) {
	if scope != nil && s.CandidateScopeConfigured() {
		if err := s.ValidateCandidateScope(ctx, ruleID, *scope); err != nil {
			return nil, err
		}
	}
	rule, err := s.repo.GetByID(ctx, ruleID)
	if err != nil {
		return nil, fmt.Errorf("取得來源規則失敗: %w", err)
	}
	if rule.RevisionID != revisionID {
		return nil, candidateRevisionConflict()
	}
	links, err := s.repo.ListLinks(ctx, rule.ID)
	if err != nil {
		return nil, fmt.Errorf("列出來源規則連結失敗: %w", err)
	}
	if scope != nil && s.CandidateScopeConfigured() {
		if err := s.ValidateCandidateScope(ctx, ruleID, *scope); err != nil {
			return nil, err
		}
	}
	if err := s.persistCandidateSnapshotsAtRevision(ctx, rule, links, revisionID); err != nil {
		return nil, err
	}
	snapshots, err := s.repo.ListCandidateSnapshots(ctx, rule.ID, revisionID)
	if err != nil {
		return nil, fmt.Errorf("列出來源規則候選快照失敗: %w", err)
	}
	if current, getErr := s.repo.GetByID(ctx, ruleID); getErr != nil {
		return nil, getErr
	} else if current.RevisionID != revisionID {
		return nil, candidateRevisionConflict()
	}
	return s.composeCandidateView(ctx, rule, snapshots)
}

func (s *Service) persistCandidateSnapshotsAtRevision(ctx context.Context, rule *schema.SourceRule, links []*schema.SourceRuleLink, revisionID string) error {
	if rule.RevisionID != revisionID {
		return candidateRevisionConflict()
	}
	snapshots, err := s.buildCandidateSnapshots(ctx, rule, links)
	if err != nil {
		return err
	}
	if err := s.replaceCandidateSnapshotsAtRevision(ctx, snapshots, revisionID); err != nil {
		return fmt.Errorf("儲存來源規則候選快照失敗: %w", err)
	}
	tagCandidates, err := decodeCurrentTagCandidates(snapshots)
	if err != nil {
		return err
	}
	if err := s.markStaleTagReviewDecisions(ctx, rule, tagCandidates); err != nil {
		return err
	}
	return s.refreshLocalModbusConflictSnapshots(ctx)
}

func (s *Service) replaceCandidateSnapshotsAtRevision(ctx context.Context, snapshots []*schema.SourceRuleCandidateSnapshot, revisionID string) error {
	if repo, ok := s.repo.(CandidateSnapshotCASRepository); ok {
		return repo.ReplaceCandidateSnapshotsAtRevision(ctx, snapshots, revisionID)
	}
	return s.repo.ReplaceCandidateSnapshots(ctx, snapshots)
}

func candidateRevisionConflict() error {
	return &modbusshare.Error{Code: modbusshare.ErrCodeRevisionConflict, Message: candidateRevisionConflictMessage, Retryable: true, Action: candidateRetryAction}
}

func shouldRecomputeDatabaseCandidateView(snapshots []*schema.SourceRuleCandidateSnapshot) bool {
	for _, snapshot := range snapshots {
		if snapshot == nil || snapshot.CandidateType != schema.SourceRuleCandidateTypeDatabaseOutputs {
			continue
		}
		return snapshot.Status == schema.SourceRuleCandidateStatusBlocked
	}
	return false
}

func (s *Service) recomputeDatabaseCandidateView(
	ctx context.Context,
	rule *schema.SourceRule,
	snapshots []*schema.SourceRuleCandidateSnapshot,
) ([]*schema.SourceRuleCandidateSnapshot, error) {
	nextSnapshots, err := s.buildDatabaseCandidateSnapshots(ctx, rule, snapshots)
	if err != nil {
		return nil, err
	}
	if err := s.repo.ReplaceCandidateSnapshots(ctx, nextSnapshots); err != nil {
		return nil, fmt.Errorf("儲存來源規則候選快照失敗: %w", err)
	}
	refreshedSnapshots, err := s.repo.ListCandidateSnapshots(ctx, rule.ID, rule.RevisionID)
	if err != nil {
		return nil, fmt.Errorf("列出來源規則候選快照失敗: %w", err)
	}
	return refreshedSnapshots, nil
}

func (s *Service) buildDatabaseCandidateSnapshots(
	ctx context.Context,
	rule *schema.SourceRule,
	snapshots []*schema.SourceRuleCandidateSnapshot,
) ([]*schema.SourceRuleCandidateSnapshot, error) {
	var tagsSnapshot *schema.SourceRuleCandidateSnapshot
	var databaseSnapshot *schema.SourceRuleCandidateSnapshot
	var localModbusSnapshot *schema.SourceRuleCandidateSnapshot
	for _, snapshot := range snapshots {
		if snapshot == nil {
			continue
		}
		switch snapshot.CandidateType {
		case schema.SourceRuleCandidateTypeTags:
			tagsSnapshot = cloneCandidateSnapshot(snapshot)
		case schema.SourceRuleCandidateTypeDatabaseOutputs:
			databaseSnapshot = cloneCandidateSnapshot(snapshot)
		case schema.SourceRuleCandidateTypeLocalModbusOutputs:
			localModbusSnapshot = cloneCandidateSnapshot(snapshot)
		}
	}
	if tagsSnapshot == nil {
		return nil, fmt.Errorf("來源規則 %s revision %s 缺少 tags 候選快照", rule.ID, rule.RevisionID)
	}
	if databaseSnapshot == nil {
		return nil, fmt.Errorf("來源規則 %s revision %s 缺少 database outputs 候選快照", rule.ID, rule.RevisionID)
	}
	if localModbusSnapshot == nil {
		return nil, fmt.Errorf("來源規則 %s revision %s 缺少 local modbus outputs 候選快照", rule.ID, rule.RevisionID)
	}

	tagCandidates, err := decodeCandidatePayload[schema.SourceRuleTagCandidate](tagsSnapshot.Payload)
	if err != nil {
		return nil, fmt.Errorf("解析來源規則 tag 候選快照失敗: %w", err)
	}
	databaseCandidates, databaseStatus, databaseReason, err := s.buildDatabaseOutputCandidates(ctx, rule, tagCandidates)
	if err != nil {
		return nil, err
	}
	databasePayload, err := marshalCandidateSnapshotPayload(databaseCandidates)
	if err != nil {
		return nil, err
	}

	databaseSnapshot.Payload = databasePayload
	databaseSnapshot.Status = databaseStatus
	databaseSnapshot.Reason = databaseReason
	databaseSnapshot.GeneratedAt = time.Now().UTC()

	return []*schema.SourceRuleCandidateSnapshot{
		tagsSnapshot,
		databaseSnapshot,
		localModbusSnapshot,
	}, nil
}

func (s *Service) composeCandidateView(ctx context.Context, rule *schema.SourceRule, snapshots []*schema.SourceRuleCandidateSnapshot) (*CandidateSnapshotView, error) {
	snapshotIndex := make(map[schema.SourceRuleCandidateType]*schema.SourceRuleCandidateSnapshot, len(snapshots))
	for _, snapshot := range snapshots {
		if snapshot == nil {
			continue
		}
		snapshotIndex[snapshot.CandidateType] = snapshot
	}

	tagsSnapshot, ok := snapshotIndex[schema.SourceRuleCandidateTypeTags]
	if !ok {
		return nil, fmt.Errorf("來源規則 %s revision %s 缺少 tags 候選快照", rule.ID, rule.RevisionID)
	}
	databaseSnapshot, ok := snapshotIndex[schema.SourceRuleCandidateTypeDatabaseOutputs]
	if !ok {
		return nil, fmt.Errorf("來源規則 %s revision %s 缺少 database outputs 候選快照", rule.ID, rule.RevisionID)
	}
	localModbusSnapshot, ok := snapshotIndex[schema.SourceRuleCandidateTypeLocalModbusOutputs]
	if !ok {
		return nil, fmt.Errorf("來源規則 %s revision %s 缺少 local modbus outputs 候選快照", rule.ID, rule.RevisionID)
	}

	tagCandidates, err := decodeCandidatePayload[schema.SourceRuleTagCandidate](tagsSnapshot.Payload)
	if err != nil {
		return nil, fmt.Errorf("解析來源規則 tag 候選快照失敗: %w", err)
	}
	tagViews := make([]any, 0, len(tagCandidates))
	for _, candidate := range tagCandidates {
		view, err := s.buildTagCandidateView(ctx, candidate)
		if err != nil {
			return nil, err
		}
		tagViews = append(tagViews, view)
	}

	databaseCandidates, err := decodeCandidatePayload[schema.SourceRuleDatabaseOutputCandidate](databaseSnapshot.Payload)
	if err != nil {
		return nil, fmt.Errorf("解析來源規則 database output 候選快照失敗: %w", err)
	}
	localModbusCandidates, err := decodeCandidatePayload[schema.SourceRuleLocalModbusOutputCandidate](localModbusSnapshot.Payload)
	if err != nil {
		return nil, fmt.Errorf("解析來源規則 local modbus 候選快照失敗: %w", err)
	}

	return &CandidateSnapshotView{
		SourceRuleID: rule.ID,
		RevisionID:   rule.RevisionID,
		Tags: CandidateSetView{
			Status:     tagsSnapshot.Status,
			Reason:     tagsSnapshot.Reason,
			Candidates: tagViews,
		},
		DatabaseOutputs: CandidateSetView{
			Status:     databaseSnapshot.Status,
			Reason:     databaseSnapshot.Reason,
			Candidates: sliceToAny(databaseCandidates),
		},
		LocalModbusOutputs: CandidateSetView{
			Status:     localModbusSnapshot.Status,
			Reason:     localModbusSnapshot.Reason,
			Candidates: sliceToAny(localModbusCandidates),
		},
	}, nil
}

func (s *Service) buildTagCandidateView(ctx context.Context, candidate schema.SourceRuleTagCandidate) (*TagCandidateView, error) {
	view := &TagCandidateView{
		SourceRuleTagCandidate: candidate,
		Status:                 schema.MappingStatusDraft,
	}
	if candidate.MappingID == nil || s.mappingSvc == nil {
		return view, nil
	}

	mappingRecord, err := s.mappingSvc.GetByID(ctx, *candidate.MappingID)
	if err != nil {
		return nil, fmt.Errorf("取得來源規則候選映射失敗: %w", err)
	}
	view.Status = mappingRecord.Status
	if view.Status == "" {
		if mappingRecord.Enabled {
			view.Status = schema.MappingStatusActive
		} else {
			view.Status = schema.MappingStatusDraft
		}
	}
	view.LastAppliedSignature = mappingRecord.LastAppliedSignature
	if view.LastAppliedSignature == "" {
		steps, err := decodeTransformPipeline(mappingRecord.TransformPipeline)
		if err != nil {
			return nil, fmt.Errorf("解析來源規則候選既有映射轉換管線失敗: %w", err)
		}
		view.LastAppliedSignature, err = mappingCandidateSignature(steps)
		if err != nil {
			return nil, err
		}
	}
	view.BlockingReason = mappingRecord.BlockingReason
	return view, nil
}

func decodeCandidatePayload[T any](payload string) ([]T, error) {
	var envelope struct {
		Candidates []T `json:"candidates"`
	}
	if err := json.Unmarshal([]byte(payload), &envelope); err != nil {
		return nil, err
	}
	if envelope.Candidates == nil {
		return []T{}, nil
	}
	return envelope.Candidates, nil
}

func sliceToAny[T any](values []T) []any {
	items := make([]any, 0, len(values))
	for _, value := range values {
		items = append(items, value)
	}
	return items
}
