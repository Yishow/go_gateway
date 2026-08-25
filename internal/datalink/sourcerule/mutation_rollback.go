package sourcerule

import (
	"context"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
)

// sourceRuleSnapshot is the durable state needed to compensate a mutation
// when the subsequent runtime Share projection cannot be committed.
type sourceRuleSnapshot struct {
	rule      *schema.SourceRule
	links     []*schema.SourceRuleLink
	snapshots []*schema.SourceRuleCandidateSnapshot
	points    map[string]*schema.Point
	tags      map[string]*schema.Tag
	mappings  map[string]*schema.Mapping
	share     *ShareRuntimeSnapshot
}

func (s *Service) snapshotRule(ctx context.Context, id string) (*sourceRuleSnapshot, error) {
	rule, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	links, err := s.repo.ListLinks(ctx, id)
	if err != nil {
		return nil, err
	}
	snapshots, err := s.repo.ListCandidateSnapshots(ctx, id, rule.RevisionID)
	if err != nil {
		return nil, err
	}
	snapshot := &sourceRuleSnapshot{rule: cloneRule(rule), links: cloneSourceRuleLinks(links), snapshots: cloneCandidateSnapshotsForRollback(snapshots), points: make(map[string]*schema.Point), tags: make(map[string]*schema.Tag), mappings: make(map[string]*schema.Mapping)}
	for _, link := range links {
		if link == nil {
			continue
		}
		pointRecord, pointErr := s.pointSvc.GetByID(ctx, link.PointID)
		if pointErr != nil {
			return nil, fmt.Errorf("snapshot source-rule point %s: %w", link.PointID, pointErr)
		}
		snapshot.points[link.PointID] = clonePointForRollback(pointRecord)
		if s.tagSvc != nil && link.TagID != nil && *link.TagID != "" {
			tagRecord, tagErr := s.tagSvc.GetByID(ctx, *link.TagID)
			if tagErr != nil && !errors.Is(tagErr, tag.ErrTagNotFound) {
				return nil, fmt.Errorf("snapshot source-rule tag %s: %w", *link.TagID, tagErr)
			}
			if tagErr == nil {
				snapshot.tags[*link.TagID] = cloneTagForRollback(tagRecord)
			}
		}
		if s.mappingSvc != nil && link.MappingID != nil && *link.MappingID != "" {
			mappingRecord, mappingErr := s.mappingSvc.GetByID(ctx, *link.MappingID)
			if mappingErr != nil && !errors.Is(mappingErr, mapping.ErrMappingNotFound) {
				return nil, fmt.Errorf("snapshot source-rule mapping %s: %w", *link.MappingID, mappingErr)
			}
			if mappingErr == nil {
				snapshot.mappings[*link.MappingID] = cloneMappingForRollback(mappingRecord)
			}
		}
	}
	if share, ok := s.shareRuntimeSync.(ShareRuntimeSnapshotter); ok {
		shareSnapshot, shareErr := share.SnapshotSourceRuleShare(ctx, RuntimeReconcileRequest{Operation: RuntimeReconcileOperationUpdate, Scope: RuntimeReconcileScope{RuleID: rule.ID, DeviceID: rule.DeviceID}})
		if shareErr != nil {
			return nil, fmt.Errorf("snapshot Share projection: %w", shareErr)
		}
		snapshot.share = &shareSnapshot
	}
	return snapshot, nil
}

func clonePointForRollback(record *schema.Point) *schema.Point { return clonePointSnapshot(record) }
func clonePointSnapshot(record *schema.Point) *schema.Point {
	if record == nil {
		return nil
	}
	copyRecord := *record
	if record.PollingGroupID != nil {
		value := *record.PollingGroupID
		copyRecord.PollingGroupID = &value
	}
	if record.LastReadAt != nil {
		value := *record.LastReadAt
		copyRecord.LastReadAt = &value
	}
	if record.LastValue != nil {
		value := *record.LastValue
		copyRecord.LastValue = &value
	}
	return &copyRecord
}
func cloneTagForRollback(record *schema.Tag) *schema.Tag {
	if record == nil {
		return nil
	}
	copyRecord := *record
	return &copyRecord
}
func cloneMappingForRollback(record *schema.Mapping) *schema.Mapping {
	if record == nil {
		return nil
	}
	copyRecord := *record
	return &copyRecord
}

func cloneCandidateSnapshotsForRollback(snapshots []*schema.SourceRuleCandidateSnapshot) []*schema.SourceRuleCandidateSnapshot {
	if len(snapshots) == 0 {
		return nil
	}
	cloned := make([]*schema.SourceRuleCandidateSnapshot, 0, len(snapshots))
	for _, snapshot := range snapshots {
		if snapshot == nil {
			continue
		}
		copySnapshot := *snapshot
		cloned = append(cloned, &copySnapshot)
	}
	return cloned
}

// restoreRuleSnapshot compensates the durable rule/candidate projection. The
// runtime reconciler owns runtime rollback; this restores its authoritative
// input so the next reconcile sees the state that preceded the failed save.
func (s *Service) restoreRuleSnapshot(ctx context.Context, snapshot *sourceRuleSnapshot) error {
	if snapshot == nil || snapshot.rule == nil {
		return fmt.Errorf("source rule snapshot is empty")
	}
	current, currentErr := s.repo.GetByID(ctx, snapshot.rule.ID)
	if currentErr != nil && !errors.Is(currentErr, ErrSourceRuleNotFound) {
		return fmt.Errorf("read current source rule before rollback: %w", currentErr)
	}
	var restoreErrs []error
	currentLinks, linksErr := s.repo.ListLinks(ctx, snapshot.rule.ID)
	if linksErr != nil {
		return fmt.Errorf("read current source-rule links before rollback: %w", linksErr)
	}
	oldPointIDs := make(map[string]struct{}, len(snapshot.points))
	for id := range snapshot.points {
		oldPointIDs[id] = struct{}{}
	}
	for _, link := range currentLinks {
		if link == nil {
			continue
		}
		if _, keep := oldPointIDs[link.PointID]; !keep {
			if err := s.cleanupRuleLinkResources(ctx, currentOrSnapshotRule(current, snapshot), link); err != nil {
				restoreErrs = append(restoreErrs, err)
			}
		}
	}
	if currentErr == nil && current.RevisionID != snapshot.rule.RevisionID {
		if err := s.repo.DeleteCandidateSnapshots(ctx, current.ID, current.RevisionID); err != nil {
			restoreErrs = append(restoreErrs, fmt.Errorf("remove failed source-rule candidate snapshot: %w", err))
		}
	}
	if currentErr != nil {
		if err := s.repo.Create(ctx, cloneRule(snapshot.rule)); err != nil {
			restoreErrs = append(restoreErrs, fmt.Errorf("restore source rule: %w", err))
		}
	} else if err := s.repo.Update(ctx, cloneRule(snapshot.rule)); err != nil {
		restoreErrs = append(restoreErrs, fmt.Errorf("restore source rule: %w", err))
	}
	if err := s.repo.DeleteLinks(ctx, snapshot.rule.ID); err != nil {
		restoreErrs = append(restoreErrs, fmt.Errorf("restore source-rule links: %w", err))
	}
	if len(snapshot.links) > 0 {
		if err := s.repo.CreateLinks(ctx, cloneSourceRuleLinks(snapshot.links)); err != nil {
			restoreErrs = append(restoreErrs, fmt.Errorf("restore source-rule links: %w", err))
		}
	}
	if len(snapshot.snapshots) > 0 {
		if err := s.repo.ReplaceCandidateSnapshots(ctx, cloneCandidateSnapshotsForRollback(snapshot.snapshots)); err != nil {
			restoreErrs = append(restoreErrs, fmt.Errorf("restore source-rule candidate snapshots: %w", err))
		}
	}
	for _, record := range snapshot.points {
		if err := s.pointSvc.Restore(ctx, record); err != nil {
			restoreErrs = append(restoreErrs, err)
		}
		if s.runtimeSync != nil {
			s.runtimeSync.UpsertPoint(record)
		}
	}
	for _, record := range snapshot.tags {
		if s.tagSvc != nil {
			if err := s.tagSvc.Restore(ctx, record); err != nil {
				restoreErrs = append(restoreErrs, err)
			}
		}
	}
	for _, record := range snapshot.mappings {
		if s.mappingSvc != nil {
			if err := s.mappingSvc.Restore(ctx, record); err != nil {
				restoreErrs = append(restoreErrs, err)
			}
		}
	}
	if snapshot.share != nil {
		if share, ok := s.shareRuntimeSync.(ShareRuntimeSnapshotter); ok {
			if _, err := share.RestoreSourceRuleShare(ctx, *snapshot.share); err != nil {
				restoreErrs = append(restoreErrs, err)
			}
		}
	}
	return errors.Join(restoreErrs...)
}

func currentOrSnapshotRule(current *schema.SourceRule, snapshot *sourceRuleSnapshot) *schema.SourceRule {
	if current != nil {
		return current
	}
	return snapshot.rule
}

func inverseRuntimeOperation(operation RuntimeReconcileOperation) (RuntimeReconcileOperation, bool) {
	switch operation {
	case RuntimeReconcileOperationCreate:
		return RuntimeReconcileOperationDelete, true
	case RuntimeReconcileOperationUpdate:
		return RuntimeReconcileOperationUpdate, true
	case RuntimeReconcileOperationDelete:
		return RuntimeReconcileOperationCreate, true
	case RuntimeReconcileOperationEnable:
		return RuntimeReconcileOperationDisable, true
	case RuntimeReconcileOperationDisable:
		return RuntimeReconcileOperationEnable, true
	default:
		return "", false
	}
}

func (s *Service) compensateFailedShareMutation(ctx context.Context, snapshot *sourceRuleSnapshot, operation RuntimeReconcileOperation, mutationErr error) error {
	rollbackErr := s.restoreRuleSnapshot(ctx, snapshot)
	if snapshot != nil && snapshot.rule != nil {
		if inverse, ok := inverseRuntimeOperation(operation); ok {
			outcome := s.reconcileSourceRuleRuntime(ctx, inverse, snapshot.rule.ID, snapshot.rule.DeviceID)
			if outcome.Status != RuntimeReconcileStatusAligned && outcome.Status != RuntimeReconcileStatusDeferred {
				rollbackErr = errors.Join(rollbackErr, fmt.Errorf("source-rule runtime compensation did not align: %s", outcome.Status))
			}
			if compensator, ok := s.shareRuntimeSync.(ShareRuntimeCompensator); ok {
				shareOutcome, shareErr := compensator.CompensateSourceRuleShare(ctx, RuntimeReconcileRequest{Operation: inverse, Scope: RuntimeReconcileScope{RuleID: snapshot.rule.ID, DeviceID: snapshot.rule.DeviceID}})
				if shareErr != nil || (shareOutcome.Status != RuntimeReconcileStatusAligned && shareOutcome.Status != RuntimeReconcileStatusDeferred) {
					rollbackErr = errors.Join(rollbackErr, shareErr, fmt.Errorf("source-rule Share compensation did not align: %s", shareOutcome.Status))
				}
			}
		}
	}
	if rollbackErr != nil {
		return s.dirtyUnknown(ctx)
	}
	return mutationErr
}

func (s *Service) dirtyUnknown(ctx context.Context) error {
	if failCloser, ok := s.shareRuntimeSync.(ShareRuntimeFailCloser); ok {
		failCloser.FailClosed(ctx)
	}
	return dirtyUnknownSourceRuleError()
}

func dirtyUnknownSourceRuleError() error {
	return &modbusshare.Error{Code: modbusshare.ErrCodeDirtyUnknown, Message: "source-rule mutation compensation is uncertain", Retryable: true, DirtyState: "dirty_unknown", Action: "run source-rule recovery reconcile"}
}
