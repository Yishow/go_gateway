package sourcerule

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"go-gateway/internal/datalink/schema"
)

const localModbusConflictReasonPrefix = "local 5020 register conflict"

type localModbusConflictSnapshot struct {
	ruleID     string
	revisionID string
	snapshots  []*schema.SourceRuleCandidateSnapshot
	candidates []schema.SourceRuleLocalModbusOutputCandidate
}

type localModbusConflictCandidateKey struct {
	ruleID      string
	revisionID  string
	candidateID string
}

type localModbusConflictCandidateRef struct {
	key     localModbusConflictCandidateKey
	ruleID  string
	address string
	start   int
	end     int
}

func (s *Service) refreshLocalModbusConflictSnapshots(ctx context.Context) error {
	rules, err := s.repo.List(ctx, ListFilter{})
	if err != nil {
		return fmt.Errorf("列出來源規則失敗: %w", err)
	}

	conflictSnapshots, err := s.listLocalModbusConflictSnapshots(ctx, rules)
	if err != nil {
		return err
	}
	if len(conflictSnapshots) == 0 {
		return nil
	}

	reasonsByCandidate := buildLocalModbusConflictReasons(conflictSnapshots)
	generatedAt := time.Now().UTC()
	for _, conflictSnapshot := range conflictSnapshots {
		nextSnapshots, err := applyLocalModbusConflictReasons(conflictSnapshot, reasonsByCandidate, generatedAt)
		if err != nil {
			return err
		}
		if err := s.repo.ReplaceCandidateSnapshots(ctx, nextSnapshots); err != nil {
			return fmt.Errorf("儲存 local modbus 衝突候選快照失敗: %w", err)
		}
	}
	return nil
}

func (s *Service) listLocalModbusConflictSnapshots(
	ctx context.Context,
	rules []*schema.SourceRule,
) ([]localModbusConflictSnapshot, error) {
	sort.Slice(rules, func(i, j int) bool {
		if rules[i] == nil {
			return false
		}
		if rules[j] == nil {
			return true
		}
		return rules[i].ID < rules[j].ID
	})

	result := make([]localModbusConflictSnapshot, 0, len(rules))
	for _, rule := range rules {
		if rule == nil || strings.TrimSpace(rule.RevisionID) == "" {
			continue
		}

		snapshots, err := s.repo.ListCandidateSnapshots(ctx, rule.ID, rule.RevisionID)
		if err != nil {
			return nil, fmt.Errorf("列出來源規則候選快照失敗: %w", err)
		}

		localSnapshot := currentLocalModbusSnapshot(snapshots)
		if localSnapshot == nil {
			continue
		}
		candidates, err := decodeCandidatePayload[schema.SourceRuleLocalModbusOutputCandidate](localSnapshot.Payload)
		if err != nil {
			return nil, fmt.Errorf("解析 local modbus 候選快照失敗: %w", err)
		}

		result = append(result, localModbusConflictSnapshot{
			ruleID:     rule.ID,
			revisionID: rule.RevisionID,
			snapshots:  cloneCandidateSnapshots(snapshots),
			candidates: candidates,
		})
	}
	return result, nil
}

func buildLocalModbusConflictReasons(
	conflictSnapshots []localModbusConflictSnapshot,
) map[localModbusConflictCandidateKey][]string {
	refs := make([]localModbusConflictCandidateRef, 0)
	for _, conflictSnapshot := range conflictSnapshots {
		for _, candidate := range conflictSnapshot.candidates {
			normalized := normalizeLocalModbusConflictCandidate(candidate)
			if !shouldEvaluateLocalModbusConflictCandidate(normalized) {
				continue
			}
			registerCount := normalized.RegisterCount
			if registerCount <= 0 {
				registerCount = 1
			}
			start := int(*normalized.Register)
			refs = append(refs, localModbusConflictCandidateRef{
				key: localModbusConflictCandidateKey{
					ruleID:      conflictSnapshot.ruleID,
					revisionID:  conflictSnapshot.revisionID,
					candidateID: normalized.ID,
				},
				ruleID:  conflictSnapshot.ruleID,
				address: normalized.Address,
				start:   start,
				end:     start + registerCount - 1,
			})
		}
	}

	reasonsByCandidate := make(map[localModbusConflictCandidateKey][]string)
	for i := 0; i < len(refs); i++ {
		for j := i + 1; j < len(refs); j++ {
			overlapStart, overlapEnd, ok := localModbusConflictRange(refs[i], refs[j])
			if !ok {
				continue
			}
			reasonsByCandidate[refs[i].key] = append(
				reasonsByCandidate[refs[i].key],
				fmt.Sprintf("%s: overlaps rule %s address %s on registers %d-%d", localModbusConflictReasonPrefix, refs[j].ruleID, refs[j].address, overlapStart, overlapEnd),
			)
			reasonsByCandidate[refs[j].key] = append(
				reasonsByCandidate[refs[j].key],
				fmt.Sprintf("%s: overlaps rule %s address %s on registers %d-%d", localModbusConflictReasonPrefix, refs[i].ruleID, refs[i].address, overlapStart, overlapEnd),
			)
		}
	}

	for key, reasons := range reasonsByCandidate {
		sort.Strings(reasons)
		reasonsByCandidate[key] = compactSortedStrings(reasons)
	}
	return reasonsByCandidate
}

func applyLocalModbusConflictReasons(
	conflictSnapshot localModbusConflictSnapshot,
	reasonsByCandidate map[localModbusConflictCandidateKey][]string,
	generatedAt time.Time,
) ([]*schema.SourceRuleCandidateSnapshot, error) {
	nextSnapshots := cloneCandidateSnapshots(conflictSnapshot.snapshots)
	localSnapshot := currentLocalModbusSnapshot(nextSnapshots)
	if localSnapshot == nil {
		return nil, fmt.Errorf("來源規則 %s revision %s 缺少 local modbus 候選快照", conflictSnapshot.ruleID, conflictSnapshot.revisionID)
	}

	candidates := make([]schema.SourceRuleLocalModbusOutputCandidate, 0, len(conflictSnapshot.candidates))
	blockingReason := ""
	for _, candidate := range conflictSnapshot.candidates {
		nextCandidate := normalizeLocalModbusConflictCandidate(candidate)
		key := localModbusConflictCandidateKey{
			ruleID:      conflictSnapshot.ruleID,
			revisionID:  conflictSnapshot.revisionID,
			candidateID: nextCandidate.ID,
		}
		if reasons := reasonsByCandidate[key]; len(reasons) > 0 {
			nextCandidate.Status = schema.SourceRuleLocalModbusOutputStatusBlockedConflict
			nextCandidate.BlockingReason = strings.Join(reasons, "; ")
			if blockingReason == "" {
				blockingReason = nextCandidate.BlockingReason
			}
		}
		candidates = append(candidates, nextCandidate)
	}

	payload, err := marshalCandidateSnapshotPayload(candidates)
	if err != nil {
		return nil, err
	}
	localSnapshot.Payload = payload
	localSnapshot.GeneratedAt = generatedAt
	if localSnapshot.Status == schema.SourceRuleCandidateStatusReady && blockingReason != "" {
		localSnapshot.Status = schema.SourceRuleCandidateStatusBlocked
		localSnapshot.Reason = blockingReason
	}
	return nextSnapshots, nil
}

func normalizeLocalModbusConflictCandidate(candidate schema.SourceRuleLocalModbusOutputCandidate) schema.SourceRuleLocalModbusOutputCandidate {
	if candidate.Status == schema.SourceRuleLocalModbusOutputStatusBlockedConflict ||
		strings.HasPrefix(candidate.BlockingReason, localModbusConflictReasonPrefix) {
		candidate.Status = schema.SourceRuleLocalModbusOutputStatusDeferred
		candidate.BlockingReason = ""
	}
	return candidate
}

func shouldEvaluateLocalModbusConflictCandidate(candidate schema.SourceRuleLocalModbusOutputCandidate) bool {
	if candidate.Register == nil || candidate.RegisterCount <= 0 {
		return false
	}
	if candidate.Status == schema.SourceRuleLocalModbusOutputStatusOutOfSync {
		return false
	}
	return strings.TrimSpace(candidate.BlockingReason) == ""
}

func localModbusConflictRange(
	left localModbusConflictCandidateRef,
	right localModbusConflictCandidateRef,
) (start, end int, overlaps bool) {
	if left.end < right.start || right.end < left.start {
		return 0, 0, false
	}
	overlapStart := maxInt(left.start, right.start)
	overlapEnd := minInt(left.end, right.end)
	return overlapStart, overlapEnd, true
}

func cloneCandidateSnapshots(snapshots []*schema.SourceRuleCandidateSnapshot) []*schema.SourceRuleCandidateSnapshot {
	result := make([]*schema.SourceRuleCandidateSnapshot, 0, len(snapshots))
	for _, snapshot := range snapshots {
		result = append(result, cloneCandidateSnapshot(snapshot))
	}
	return result
}

func compactSortedStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	result := []string{values[0]}
	for _, value := range values[1:] {
		if value == result[len(result)-1] {
			continue
		}
		result = append(result, value)
	}
	return result
}

func maxInt(left, right int) int {
	if left > right {
		return left
	}
	return right
}

func minInt(left, right int) int {
	if left < right {
		return left
	}
	return right
}
