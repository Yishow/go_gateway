package sourcerule

import (
	"context"
	"encoding/json"
	"fmt"

	"go-gateway/internal/datalink/schema"
)

type CandidateSnapshotView struct {
	SourceRuleID       string           `json:"source_rule_id"`
	RevisionID         string           `json:"revision_id"`
	Tags               CandidateSetView `json:"tags"`
	DatabaseOutputs    CandidateSetView `json:"database_outputs"`
	LocalModbusOutputs CandidateSetView `json:"local_modbus_outputs"`
}

type CandidateSetView struct {
	Status     schema.SourceRuleCandidateStatus `json:"status"`
	Reason     string                           `json:"reason,omitempty"`
	Candidates []any                            `json:"candidates"`
}

type TagCandidateView struct {
	schema.SourceRuleTagCandidate
	Status               schema.MappingStatus `json:"status"`
	LastAppliedSignature string               `json:"last_applied_signature,omitempty"`
	BlockingReason       string               `json:"blocking_reason,omitempty"`
}

func (s *Service) GetCandidateView(ctx context.Context, ruleID string) (*CandidateSnapshotView, error) {
	rule, err := s.repo.GetByID(ctx, ruleID)
	if err != nil {
		return nil, fmt.Errorf("取得來源規則失敗: %w", err)
	}
	snapshots, err := s.repo.ListCandidateSnapshots(ctx, rule.ID, rule.RevisionID)
	if err != nil {
		return nil, fmt.Errorf("列出來源規則候選快照失敗: %w", err)
	}
	return s.composeCandidateView(ctx, rule, snapshots)
}

func (s *Service) RecomputeCandidateView(ctx context.Context, ruleID string) (*CandidateSnapshotView, error) {
	rule, err := s.repo.GetByID(ctx, ruleID)
	if err != nil {
		return nil, fmt.Errorf("取得來源規則失敗: %w", err)
	}
	links, err := s.repo.ListLinks(ctx, rule.ID)
	if err != nil {
		return nil, fmt.Errorf("列出來源規則連結失敗: %w", err)
	}
	snapshots, err := s.buildCandidateSnapshots(ctx, rule, links)
	if err != nil {
		return nil, err
	}
	if err := s.repo.ReplaceCandidateSnapshots(ctx, snapshots); err != nil {
		return nil, fmt.Errorf("儲存來源規則候選快照失敗: %w", err)
	}
	return s.composeCandidateView(ctx, rule, snapshots)
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
	localModbusCandidates, err := decodeCandidatePayload[map[string]any](localModbusSnapshot.Payload)
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
