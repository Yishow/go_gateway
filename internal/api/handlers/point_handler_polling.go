package handlers

import (
	"context"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
)

const (
	pollQualityBad       = 0
	pollQualityUncertain = 64
	pollQualityGood      = 192
)

type pointManualPoller interface {
	PollNow(pointIDs []string) []collector.CollectedValue
}

type pointMappingLister interface {
	List(ctx context.Context, filter mapping.ListFilter) ([]*schema.Mapping, error)
}

type pointPollingGroupGetter interface {
	GetByID(ctx context.Context, id string) (*schema.PollingGroup, error)
}

func (h *PointHandler) WithPolling(
	poller pointManualPoller,
	mappingLister pointMappingLister,
	pollingGroupSvc pointPollingGroupGetter,
) *PointHandler {
	h.manualPoller = poller
	h.mappingLister = mappingLister
	h.pollingGroupSvc = pollingGroupSvc
	return h
}

func (h *PointHandler) pollPoint(ctx context.Context, pointID string) (PollResult, error) {
	pt, err := h.svc.GetByID(ctx, pointID)
	if err != nil {
		return PollResult{}, err
	}

	manualResults := h.pollNow([]string{pointID})
	if cv, ok := manualResults[pointID]; ok {
		if err := h.svc.UpdateReadResult(ctx, pointID, cv.Value, cv.Error); err != nil {
			return PollResult{}, err
		}
		return h.buildPollResultFromCollectedValue(ctx, pt, cv), nil
	}

	result := h.buildPollResultFromPoint(ctx, pt)
	if h.manualPoller != nil {
		result.Error = mergePollErrors(result.Error, "Poll result unavailable")
		result.Quality = pollQualityBad
	}
	return result, nil
}

func (h *PointHandler) pollBatch(ctx context.Context, pointIDs []string) ([]PollResult, error) {
	manualResults := h.pollNow(pointIDs)
	if len(manualResults) > 0 {
		updates := make([]point.ReadResultUpdate, 0, len(manualResults))
		for _, pointID := range pointIDs {
			cv, ok := manualResults[pointID]
			if !ok {
				continue
			}
			updates = append(updates, point.ReadResultUpdate{
				PointID: pointID,
				Value:   cv.Value,
				Error:   cv.Error,
			})
		}
		if len(updates) > 0 {
			if err := h.svc.BatchUpdateReadResult(ctx, updates); err != nil {
				return nil, err
			}
		}
	}

	results := make([]PollResult, 0, len(pointIDs))
	for _, pointID := range pointIDs {
		pt, err := h.svc.GetByID(ctx, pointID)
		if err != nil {
			results = append(results, PollResult{
				PointID: pointID,
				Error:   "Point not found",
			})
			continue
		}

		if cv, ok := manualResults[pointID]; ok {
			results = append(results, h.buildPollResultFromCollectedValue(ctx, pt, cv))
			continue
		}

		result := h.buildPollResultFromPoint(ctx, pt)
		if h.manualPoller != nil {
			result.Error = mergePollErrors(result.Error, "Poll result unavailable")
			result.Quality = pollQualityBad
		}
		results = append(results, result)
	}

	return results, nil
}

func (h *PointHandler) pollNow(pointIDs []string) map[string]collector.CollectedValue {
	if h.manualPoller == nil || len(pointIDs) == 0 {
		return map[string]collector.CollectedValue{}
	}

	pollerValue := reflect.ValueOf(h.manualPoller)
	if pollerValue.Kind() == reflect.Pointer && pollerValue.IsNil() {
		return map[string]collector.CollectedValue{}
	}

	values := h.manualPoller.PollNow(pointIDs)
	results := make(map[string]collector.CollectedValue, len(values))
	for _, value := range values {
		results[value.PointID] = value
	}
	return results
}

func (h *PointHandler) buildPollResultFromCollectedValue(
	ctx context.Context,
	pt *schema.Point,
	cv collector.CollectedValue,
) PollResult {
	observedAt := cv.Timestamp
	if observedAt.IsZero() {
		observedAt = time.Now()
	}

	transformedValue, transformErr := h.resolveTransformedValue(ctx, pt.ID, cv.Value)
	quality := pollQualityCode(cv.Quality, cv.Error)
	if transformErr != "" && quality == pollQualityGood {
		quality = pollQualityUncertain
	}

	return PollResult{
		PointID:          pt.ID,
		Value:            cv.Value,
		TransformedValue: transformedValue,
		Timestamp:        observedAt.Format(time.RFC3339),
		Quality:          quality,
		Stale:            h.isPointStale(ctx, pt, observedAt, time.Now()),
		Error:            mergePollErrors(cv.Error, transformErr),
	}
}

func (h *PointHandler) buildPollResultFromPoint(ctx context.Context, pt *schema.Point) PollResult {
	observedAt := time.Time{}
	if pt.LastReadAt != nil {
		observedAt = *pt.LastReadAt
	}

	rawValue := decodeStoredPointValue(pt.LastValue)
	transformedValue, transformErr := h.resolveTransformedValue(ctx, pt.ID, rawValue)

	quality := pollQualityCode("", pt.LastError)
	if pt.LastReadAt == nil && pt.LastError == "" {
		quality = pollQualityUncertain
	}
	if transformErr != "" && quality == pollQualityGood {
		quality = pollQualityUncertain
	}

	return PollResult{
		PointID:          pt.ID,
		Value:            rawValue,
		TransformedValue: transformedValue,
		Timestamp: func() string {
			if observedAt.IsZero() {
				return ""
			}
			return observedAt.Format(time.RFC3339)
		}(),
		Quality: quality,
		Stale:   h.isPointStale(ctx, pt, observedAt, time.Now()),
		Error:   mergePollErrors(pt.LastError, transformErr),
	}
}

func (h *PointHandler) resolveTransformedValue(ctx context.Context, pointID string, rawValue any) (any, string) {
	if rawValue == nil || h.mappingLister == nil {
		return rawValue, ""
	}

	enabled := true
	mappings, err := h.mappingLister.List(ctx, mapping.ListFilter{
		PointID: &pointID,
		Enabled: &enabled,
		Limit:   1000,
	})
	if err != nil {
		return rawValue, err.Error()
	}
	if len(mappings) == 0 {
		return rawValue, ""
	}

	sort.SliceStable(mappings, func(i, j int) bool {
		if mappings[i].CreatedAt.Equal(mappings[j].CreatedAt) {
			return mappings[i].ID > mappings[j].ID
		}
		return mappings[i].CreatedAt.After(mappings[j].CreatedAt)
	})

	primary := mappings[0]
	if !shouldRunPollPipeline(primary.TransformPipeline) {
		return rawValue, ""
	}

	out, err := mapping.ExecutePipeline(rawValue, primary.TransformPipeline)
	if err != nil {
		return rawValue, err.Error()
	}
	return out.CurrentValue, ""
}

func (h *PointHandler) isPointStale(
	ctx context.Context,
	pt *schema.Point,
	observedAt time.Time,
	now time.Time,
) bool {
	if observedAt.IsZero() {
		return true
	}

	interval := time.Second
	if pt.PollingGroupID != nil && h.pollingGroupSvc != nil {
		group, err := h.pollingGroupSvc.GetByID(ctx, *pt.PollingGroupID)
		if err == nil && group.IntervalMs > 0 {
			interval = time.Duration(group.IntervalMs) * time.Millisecond
		}
	}

	return now.Sub(observedAt) > interval
}

func decodeStoredPointValue(rawValue *string) any {
	if rawValue == nil {
		return nil
	}

	var decoded any
	if err := json.Unmarshal([]byte(*rawValue), &decoded); err == nil {
		return decoded
	}

	return *rawValue
}

func shouldRunPollPipeline(rawPipeline string) bool {
	trimmed := strings.TrimSpace(rawPipeline)
	return trimmed != "" && trimmed != "[]"
}

func pollQualityCode(flag schema.QualityFlag, errMsg string) int {
	if errMsg != "" {
		return pollQualityBad
	}

	switch flag {
	case schema.QualityBad:
		return pollQualityBad
	case schema.QualityUncertain:
		return pollQualityUncertain
	default:
		return pollQualityGood
	}
}

func mergePollErrors(parts ...string) string {
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		filtered = append(filtered, part)
	}
	return strings.Join(filtered, "; ")
}
