//go:build f_write_group_fixture

package main

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"go-gateway/internal/datalink/workspace"
)

func (f *groupFixture) pause(ctx context.Context) error {
	f.opMu.Lock()
	defer f.opMu.Unlock()
	if f.deps == nil || f.deps.pollingGroup == nil || f.deps.scheduler == nil || f.deps.runtime == nil {
		return errors.New("fixture production services are unavailable")
	}
	if err := f.deps.runtime.PauseFixture(ctx); err != nil {
		return errors.New("fixture runtime pause failed")
	}
	groups, err := f.deps.pollingGroup.List(ctx)
	if err != nil {
		return errors.New("could not list persisted polling groups")
	}
	for _, group := range groups {
		if group != nil {
			f.deps.scheduler.RemovePollingGroup(group.ID)
		}
	}
	f.paused.Store(true)
	return nil
}

func (f *groupFixture) poll(ctx context.Context, request fixturePollRequest) (fixturePollResponse, error) {
	f.opMu.Lock()
	defer f.opMu.Unlock()
	if !f.paused.Load() {
		return fixturePollResponse{}, errors.New("pause is required before manual poll")
	}
	groupID := strings.TrimSpace(request.GroupID)
	if groupID == "" {
		return fixturePollResponse{}, errors.New("group_id is required")
	}
	if f.deps == nil || f.pipe == nil || f.deps.writeGroups == nil || f.deps.point == nil || f.deps.runtime == nil {
		return fixturePollResponse{}, errors.New("fixture production services are unavailable")
	}
	result, err := f.deps.writeGroups.Get(ctx, groupID)
	if err != nil || result == nil || result.Group == nil {
		return fixturePollResponse{}, errors.New("unknown canonical write group")
	}
	group := result.Group
	if group.AppliedRevision == "" || group.Status == workspace.WriteGroupStatusDeleted {
		return fixturePollResponse{}, errors.New("write group is not applied")
	}
	memberDevices, err := fixtureWriteGroupMemberDevices(group)
	if err != nil {
		return fixturePollResponse{}, err
	}
	allowed := make(map[string]struct{}, len(memberDevices))
	for pointID := range memberDevices {
		allowed[pointID] = struct{}{}
	}
	pointIDs, err := fixturePointSubset(request.PointIDs, allowed)
	if err != nil {
		return fixturePollResponse{}, err
	}
	if len(pointIDs) == 0 {
		return fixturePollResponse{}, errors.New("canonical write group has no members")
	}
	for _, pointID := range pointIDs {
		record, pointErr := f.deps.point.GetByID(ctx, pointID)
		if pointErr != nil || record == nil {
			return fixturePollResponse{}, errors.New("canonical write-group point is unavailable")
		}
		if !record.Enabled || record.DeviceID != memberDevices[pointID] {
			return fixturePollResponse{}, errors.New("canonical write-group point identity is invalid")
		}
	}
	acquisitionID := strings.TrimSpace(request.AcquisitionID)
	if len(acquisitionID) > 128 {
		return fixturePollResponse{}, errors.New("acquisition_id is too long")
	}
	if request.AcquisitionID != "" && acquisitionID == "" {
		return fixturePollResponse{}, errors.New("acquisition_id is invalid")
	}
	if acquisitionID != "" {
		f.setForcedAcquisitionID(acquisitionID)
		defer f.setForcedAcquisitionID("")
	}
	if err := f.pipe.Reconcile(ctx); err != nil {
		return fixturePollResponse{}, errors.New("pipeline reconcile failed")
	}
	if !f.capture.hasCapacity() {
		return fixturePollResponse{}, errors.New("manual poll capture failed")
	}
	before := f.capture.length()
	values := f.deps.scheduler.PollNowContext(ctx, pointIDs)
	for _, value := range values {
		if err := f.deps.runtime.AcceptFixtureCollectedValue(ctx, value); err != nil {
			return fixturePollResponse{}, errors.New("manual poll capture failed")
		}
	}
	after := f.capture.length()
	indices := make([]int, 0, after-before)
	for index := before; index < after; index++ {
		indices = append(indices, index)
	}
	return fixturePollResponse{Indices: indices, ResultCount: len(values)}, nil
}

func fixtureWriteGroupMemberDevices(group *workspace.WriteGroup) (map[string]string, error) {
	if group == nil {
		return nil, errors.New("canonical write group is missing")
	}
	memberDevices := make(map[string]string, len(group.Members))
	for _, member := range group.Members {
		if strings.TrimSpace(member.PointID) == "" || strings.TrimSpace(member.DeviceID) == "" {
			return nil, errors.New("canonical write-group member identity is incomplete")
		}
		if previous, exists := memberDevices[member.PointID]; exists && previous != member.DeviceID {
			return nil, errors.New("canonical write-group member identity conflicts")
		}
		memberDevices[member.PointID] = member.DeviceID
	}
	return memberDevices, nil
}

func fixturePointSubset(requested []string, allowed map[string]struct{}) ([]string, error) {
	if len(requested) == 0 {
		pointIDs := make([]string, 0, len(allowed))
		for pointID := range allowed {
			pointIDs = append(pointIDs, pointID)
		}
		slices.Sort(pointIDs)
		return pointIDs, nil
	}
	pointIDs := make([]string, 0, len(requested))
	seen := make(map[string]struct{}, len(requested))
	for _, raw := range requested {
		pointID := strings.TrimSpace(raw)
		if pointID == "" {
			return nil, errors.New("point_ids contains an empty point")
		}
		if _, duplicate := seen[pointID]; duplicate {
			return nil, errors.New("point_ids contains a duplicate point")
		}
		seen[pointID] = struct{}{}
		if _, ok := allowed[pointID]; !ok {
			return nil, fmt.Errorf("point is not a persisted member of group")
		}
		pointIDs = append(pointIDs, pointID)
	}
	return pointIDs, nil
}

func (f *groupFixture) release(ctx context.Context, indices []int) (fixtureReleaseResponse, error) {
	f.opMu.Lock()
	defer f.opMu.Unlock()
	if f.pipe == nil {
		return fixtureReleaseResponse{}, errors.New("fixture pipeline is unavailable")
	}
	if len(indices) == 0 {
		return fixtureReleaseResponse{}, errors.New("indices is required")
	}
	seen := make(map[int]struct{}, len(indices))
	f.capture.mu.Lock()
	for _, index := range indices {
		if _, duplicate := seen[index]; duplicate {
			f.capture.mu.Unlock()
			return fixtureReleaseResponse{}, errors.New("indices contains a duplicate capture")
		}
		seen[index] = struct{}{}
		if index < 0 || index >= len(f.capture.samples) {
			f.capture.mu.Unlock()
			return fixtureReleaseResponse{}, errors.New("unknown capture index")
		}
	}
	f.capture.mu.Unlock()
	results := make([]fixtureReleaseResult, 0, len(indices))
	for _, index := range indices {
		envelope, err := f.capture.sample(index)
		if err != nil {
			return fixtureReleaseResponse{}, err
		}
		err = f.pipe.AcceptSample(ctx, envelope)
		f.capture.mu.Lock()
		capture := &f.capture.samples[index]
		capture.attempts++
		if err == nil {
			capture.status = "accepted"
			capture.lastReason = ""
		} else {
			capture.status = "rejected"
			reasons := releaseReasons(err)
			capture.lastReason = strings.Join(reasons, ",")
		}
		f.capture.mu.Unlock()
		result := fixtureReleaseResult{Index: index}
		if err == nil {
			result.Accepted = true
			result.Outcome = "accepted"
		} else {
			result.Accepted = false
			result.Outcome = "rejected"
			reasons := releaseReasons(err)
			result.Reason = strings.Join(reasons, ",")
		}
		results = append(results, result)
	}
	return fixtureReleaseResponse{Results: results}, nil
}

func (f *groupFixture) tick(ctx context.Context) error {
	f.opMu.Lock()
	defer f.opMu.Unlock()
	if f.pipe == nil {
		return errors.New("fixture pipeline is unavailable")
	}
	if err := f.pipe.Reconcile(ctx); err != nil {
		return err
	}
	before := f.errors.Load()
	f.pipe.TickAll(ctx)
	if f.errors.Load() > before {
		return errors.New("pipeline tick failed")
	}
	return nil
}
